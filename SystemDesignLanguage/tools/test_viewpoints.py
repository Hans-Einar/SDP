import json
from pathlib import Path
import subprocess
import sys
import tempfile
import unittest

from viewpoints import Views, core

ROOT = Path(__file__).resolve().parents[2]
CLI = Path(__file__).with_name('sdl.py')


def source(extra=''):
    return core.canonicalize(core.parse('''language design-core version 0.4.
container Host.
unit Parser.
functionality ParseSource.
capability SourceModel.
interface InputPort.
mode Inspection.
Host contains Parser.
Parser owns ParseSource.
Parser provides SourceModel.
ParseSource realizes SourceModel.
Parser consumes InputPort.
SourceModel requires InputPort in mode Inspection.
''' + extra))


class ViewpointTests(unittest.TestCase):
    def test_projects_only_declared_relationships_with_source_provenance(self):
        views = Views(source())
        facts = {r['id']: r for r in views.facts}
        for d in views.diagrams:
            for a, verb, b, fact in d.edges:
                self.assertEqual((a, verb, b), (facts[fact]['subject'], facts[fact]['verb'], facts[fact]['object']))
                self.assertGreater(facts[fact]['line'], 1)
            self.assertTrue(set(views.node_map(d)) == {f'n_{n}' for n in d.nodes})
        self.assertNotIn('sequenceDiagram', views.markdown())
        self.assertNotIn('```mermaid\npacket\n', views.markdown())
        self.assertIn('NO_SCENARIOS', {g['code'] for g in views.gaps})
        self.assertIn('Capability er ikke Feature', views.markdown())

    def test_selection_and_changed_source_control_output(self):
        original = Views(source(), ['VP02']).markdown()
        changed_source = core.canonicalize(core.parse(source().replace('Parser', 'Compiler')))
        changed = Views(changed_source, ['VP02']).markdown()
        self.assertIn('Compiler', changed)
        self.assertNotIn('Parser', changed)
        self.assertNotEqual(original, changed)
        self.assertNotIn('## VP03', original)
        self.assertNotIn('## VP11', original)
        with self.assertRaises(ValueError):
            Views(source(), ['ImaginaryView'])

    def test_all_facts_remain_in_inventory_including_unrealized_functionality(self):
        text = source().replace('capability SourceModel.', 'functionality Extra.\ncapability SourceModel.')
        text += 'Parser owns Extra.\n'
        views = Views(core.canonicalize(core.parse(text)))
        self.assertIn('Parser owns Extra', views.markdown())
        self.assertEqual(len(views.facts), len(views.model.statements))
        self.assertNotIn('Extra -->|realizes|', views.markdown())

    def test_cli_rejects_unsupported_profile_without_writing_output(self):
        with tempfile.TemporaryDirectory() as tmp:
            p = Path(tmp) / 'source.design'; out = Path(tmp) / 'views'
            p.write_text('language sdl-mvp1-exercise version 0.4.\n')
            result = subprocess.run([sys.executable, str(CLI), 'viewpoints', str(p), '--output', str(out)], capture_output=True)
            self.assertEqual(result.returncode, 2)
            self.assertFalse(out.exists())

    def test_cli_regenerates_same_artifacts_without_renderer(self):
        with tempfile.TemporaryDirectory() as tmp:
            p = Path(tmp) / 'source.design'; out = Path(tmp) / 'views'
            p.write_text(source())
            cmd = [sys.executable, str(CLI), 'viewpoints', str(p), '--output', str(out), '--viewpoint', 'VP02']
            subprocess.run(cmd, check=True, capture_output=True)
            before = {str(f.relative_to(out)): f.read_bytes() for f in out.rglob('*') if f.is_file()}
            subprocess.run(cmd, check=True, capture_output=True)
            self.assertEqual(before, {str(f.relative_to(out)): f.read_bytes() for f in out.rglob('*') if f.is_file()})
            report = json.loads((out / 'manifest.json').read_text())
            self.assertEqual([v['id'] for v in report['viewpoints']], ['VP02'])
            self.assertEqual(len(report['diagrams']), 2)
            self.assertIsNone(report['renderer'])

    def test_selection_removes_stale_generated_files_but_preserves_notes(self):
        with tempfile.TemporaryDirectory() as tmp:
            p = Path(tmp) / 'source.design'; out = Path(tmp) / 'views'
            p.write_text(source())
            cmd = [sys.executable, str(CLI), 'viewpoints', str(p), '--output', str(out)]
            subprocess.run(cmd, check=True, capture_output=True)
            (out / 'notes.txt').write_text('Keep my notes')
            self.assertTrue(list((out / 'diagrams').glob('VP03-*')))
            subprocess.run(cmd + ['--viewpoint', 'VP02'], check=True, capture_output=True)
            self.assertFalse(list((out / 'diagrams').glob('VP03-*')))
            self.assertEqual((out / 'notes.txt').read_text(), 'Keep my notes')

    def test_failed_renderer_preserves_previous_report(self):
        with tempfile.TemporaryDirectory() as tmp:
            p = Path(tmp) / 'source.design'; out = Path(tmp) / 'views'
            p.write_text(source())
            cmd = [sys.executable, str(CLI), 'viewpoints', str(p), '--output', str(out)]
            subprocess.run(cmd, check=True, capture_output=True)
            before = {str(f.relative_to(out)): f.read_bytes() for f in out.rglob('*') if f.is_file()}
            result = subprocess.run(cmd + ['--renderer', str(Path(tmp) / 'missing-renderer')], capture_output=True)
            self.assertEqual(result.returncode, 2)
            self.assertEqual(before, {str(f.relative_to(out)): f.read_bytes() for f in out.rglob('*') if f.is_file()})


if __name__ == '__main__':
    unittest.main()
