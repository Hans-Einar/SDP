#!/usr/bin/env python3
"""Behavior checks for directory scope, metadata and safe CLI installation."""
import os
import errno
import pty
import subprocess
import tempfile
import unittest
from pathlib import Path

CLI = Path(__file__).resolve().parent


class KanBanCLI(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory(prefix='kanban cli ')
        self.root = Path(self.tmp.name)
        self.board = self.root / 'KanBan'
        self.board.mkdir()

    def tearDown(self):
        self.tmp.cleanup()

    def card(self, folder, number, state, tail=''):
        directory = self.board / folder
        directory.mkdir(exist_ok=True)
        path = directory / f'#{number:03}--Example with spaces.md'
        path.write_text(f'# Example\n\n| Field | Value |\n| --- | --- |\n| id | KB-T-{number:03} |\n| CardState | {state} |\n\n'+tail)
        return path

    def run_cli(self, *args, cwd=None):
        return subprocess.run(['bash', str(CLI / 'kanban.sh'), *args], cwd=cwd or self.board,
                              text=True, capture_output=True)

    def test_scope_alias_order_and_quoted_filenames(self):
        self.card('backlog', 2, 'queued')
        self.card('active', 1, 'gate-review')
        (self.board / 'examples').mkdir()
        (self.board / 'examples/#999.md').write_text('| Field | Value |\n| CardState | ready |\n')
        status = self.run_cli('status')
        self.assertEqual(status.returncode, 0, status.stderr)
        self.assertEqual(status.stdout, self.run_cli('state').stdout)
        self.assertLess(status.stdout.index('gate-review:'), status.stdout.index('queued:'))
        self.assertIn('Example\\ with\\ spaces.md', status.stdout)
        self.assertNotIn('999', status.stdout)
        active = self.run_cli('status', cwd=self.board / 'active')
        self.assertIn('gate-review:', active.stdout)
        self.assertNotIn('queued:', active.stdout)
        self.assertEqual(status.stdout, self.run_cli('state', str(self.board), cwd=self.root).stdout)

    def test_missing_cards_and_body_metadata_are_not_matches(self):
        (self.board / '#001.md').write_text('# No metadata\n\n```\n| CardState | ready |\n```\n')
        (self.board / '#002.md').write_text('# No metadata\n\n```\n| Field | Value |\n| CardState | ready |\n```\n')
        (self.board / 'Card-template.md').write_text('| Field | Value |\n| CardState | ready |\n')
        result = self.run_cli('status')
        self.assertEqual(result.returncode, 1)
        self.assertEqual(result.stderr, 'no KanBan cards found\n')
        self.assertEqual(result.stdout, '')

    def test_bad_duplicate_and_empty_state(self):
        p = self.card('active', 1, 'unsupported')
        for content in [p.read_text(), p.read_text().replace('unsupported', ''),
                        p.read_text().replace('| CardState | unsupported |', '| CardState | ready |\n| CardState | queued |')]:
            p.write_text(content)
            result = self.run_cli('status')
            self.assertEqual(result.returncode, 2)
            self.assertIn('invalid or duplicate CardState:', result.stderr)
            self.assertEqual(result.stdout, '')

    def test_body_table_does_not_override_metadata(self):
        self.card('active', 1, 'in-progress', '| Field | Value |\n| CardState | completed |\n')
        result = self.run_cli('state')
        self.assertEqual(result.returncode, 0)
        self.assertIn('in-progress:', result.stdout)
        self.assertNotIn('completed:', result.stdout)

    def test_arguments(self):
        self.assertEqual(self.run_cli('--help').returncode, 0)
        for args in [(), ('bad',), ('state', 'missing'), ('status', '.', 'extra')]:
            self.assertEqual(self.run_cli(*args).returncode, 2)

    def test_terminal_links_encode_complete_paths_and_close_each_link(self):
        original = self.card('active', 1, 'ready')
        special = original.with_name('#001--space æ % ? [x]\t\x1b.md')
        original.rename(special)
        self.card('backlog', 2, 'queued')
        master, slave = pty.openpty()
        try:
            process = subprocess.Popen(
                ['bash', str(CLI / 'kanban.sh'), 'status'], cwd=self.board,
                stdout=slave, stderr=subprocess.PIPE, env=dict(os.environ, TERM='xterm-kitty'))
            os.close(slave)
            slave = None
            data = bytearray()
            while True:
                try:
                    chunk = os.read(master, 4096)
                except OSError as error:
                    if error.errno != errno.EIO:
                        raise
                    break
                if not chunk:
                    break
                data.extend(chunk)
            _, stderr = process.communicate(timeout=5)
            self.assertEqual(process.returncode, 0, stderr)
        finally:
            os.close(master)
            if slave is not None:
                os.close(slave)
        opening = b'\x1b]8;;' + special.as_uri().encode() + b'\x1b\\'
        closing = b'\x1b]8;;\x1b\\'
        self.assertIn(opening, data)
        self.assertEqual(data.count(closing), 2)
        self.assertEqual(data.count(b'\x1b'), 8)  # No raw filename escape leaks.
        plain = self.run_cli('status')
        self.assertEqual(plain.returncode, 0, plain.stderr)
        self.assertNotIn('\x1b', plain.stdout)
        self.assertNotIn('file://', plain.stdout)

    def test_installer_copies_all_scripts_preserves_previous_and_is_idempotent(self):
        target = self.root / 'bin area'
        target.mkdir()
        old = target / 'kanban'
        old.write_text('previous unrelated command\n')
        cmd = ['bash', str(CLI / 'install-cli.sh'), '--bin-dir', str(target)]
        first = subprocess.run(cmd, text=True, capture_output=True)
        self.assertEqual(first.returncode, 0, first.stderr)
        backups = list(target.glob('kanban.backup.*'))
        self.assertEqual(len(backups), 1)
        self.assertEqual(backups[0].read_text(), 'previous unrelated command\n')
        for source in CLI.glob('*.sh'):
            installed = target / source.stem
            self.assertEqual(installed.read_bytes(), source.read_bytes())
            self.assertTrue(os.access(installed, os.X_OK))
        second = subprocess.run(cmd, text=True, capture_output=True)
        self.assertEqual(second.returncode, 0, second.stderr)
        self.assertEqual(list(target.glob('kanban.backup.*')), backups)
        self.card('backlog', 1, 'queued')
        result = subprocess.run([str(target / 'kanban'), 'status'], cwd=self.board, text=True, capture_output=True)
        self.assertEqual(result.returncode, 0, result.stderr)
        # The installed installer can refresh from an explicit checkout path.
        refresh = subprocess.run([str(target / 'install-cli'), '--source', str(CLI), '--bin-dir', str(target)], text=True, capture_output=True)
        self.assertEqual(refresh.returncode, 0, refresh.stderr)

    def test_installer_rejects_symlink_without_touching_target(self):
        target = self.root / 'bin'
        target.mkdir()
        victim = self.root / 'original'
        victim.write_text('preserve\n')
        (target / 'kanban').symlink_to(victim)
        result = subprocess.run(['bash', str(CLI / 'install-cli.sh'), '--bin-dir', str(target)], text=True, capture_output=True)
        self.assertEqual(result.returncode, 2)
        self.assertEqual(victim.read_text(), 'preserve\n')


if __name__ == '__main__':
    unittest.main(verbosity=2)
