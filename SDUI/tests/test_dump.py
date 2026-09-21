import json
import os
from pathlib import Path
import subprocess
import sys
import unittest

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROOT/'src'))
from sdui import parse, SduiError
from sdui.normalize import normalize
from sdui.dump import gui_dump, markdown_lines, cell_width


class DumpTest(unittest.TestCase):
    def dump(self, source, columns=80):
        return gui_dump(normalize(parse('sdui 0.2; P=' + source + ';'))['P'], columns)

    def test_rows_columns_and_groups_are_visible(self):
        text = self.dump('[<"LEFT", "RIGHT">; <"BELOW">]*b')
        row = next(line for line in text.splitlines() if 'LEFT' in line)
        self.assertIn('RIGHT', row)
        self.assertLess(text.index('RIGHT'), text.index('BELOW'))
        self.assertTrue(all(cell_width(line) <= 80 for line in text.splitlines()[1:]))

    def test_raw_markdown_and_mermaid_fences(self):
        source = '# Header\n**bold**\n```mermaid\ngraph LR; SECRET-->B\n```\nAfter'
        text = self.dump('[' + json.dumps(source) + ']')
        self.assertIn('# Header', text)
        self.assertIn('**bold**', text)
        self.assertIn('[Mermaid utelatt]', text)
        self.assertNotIn('SECRET', text)
        self.assertIn('After', text)
        literal = '````text\n```mermaid\nA-->B\n```\n````'
        self.assertEqual(markdown_lines(literal), literal.splitlines())
        self.assertEqual(markdown_lines('~~~mermaid\nA\n~~~~\nB'), ['[Mermaid utelatt]', 'B'])
        self.assertEqual(markdown_lines('```mermaid\nA'), ['[Mermaid utelatt]'])

    def test_no_terminal_escapes_and_unicode_cell_bounds(self):
        text = self.dump('["\\u001b[31m红色 é"; "\\u202eabc"]*b')
        self.assertNotIn('\x1b', text)
        self.assertNotIn('\u202e', text)
        self.assertIn('\\u001b', text)
        self.assertIn('红色 é', text)
        self.assertTrue(all(cell_width(line) <= 80 for line in text.splitlines()[1:]))

    def test_reference_layout_and_finite_space_errors(self):
        for source in ['["A"]', '[<"A", "B">]*b {16:9,<->}']:
            self.assertIn('A', self.dump(source))
        with self.assertRaises(SduiError) as error:
            self.dump('[<"A","B","C","D","E","F","G","H">]*b', 20)
        self.assertEqual(error.exception.code, 'dump-space')
        with self.assertRaises(SduiError):
            self.dump('[]', 0)

    def test_concept1_has_six_boxes_and_source_weights(self):
        tree = parse((ROOT/'examples/concept1-bucking.sdui').read_text())
        bucking = normalize(tree)['bucking']
        top, middle, track = [row[0] for row in bucking.rows]
        self.assertEqual([dict(node.layout)['y'] for node in (top,middle,track)], ['15fr','45fr','40fr'])
        self.assertEqual([dict(node.layout)['x'] for node in middle.rows[0]], ['25fr','50fr','25fr'])
        boxes = list(top.rows[0]) + list(middle.rows[0]) + [track]
        self.assertEqual([node.path.split('/')[-1] for node in boxes],
                         ['length','diameter','selection','suggestions','currentStem','stemTrack'])
        actual = gui_dump(bucking, 160)
        self.assertEqual(actual, (ROOT/'examples/concept1-bucking.dump.txt').read_text())
        self.assertNotIn('Rot --> Kapp1', actual)
        self.assertTrue(all(cell_width(line) <= 160 for line in actual.splitlines()[1:]))

    def test_cli_dump_entry_and_no_partial_error_output(self):
        env = dict(os.environ, PYTHONPATH=str(ROOT/'src'))
        run = lambda *args: subprocess.run([sys.executable,'-m','sdui',str(ROOT/'examples/concept1-bucking.sdui'),*args],capture_output=True,env=env)
        good = run('--format','dump','--entry','bucking')
        self.assertEqual(good.returncode, 0, good.stderr)
        self.assertIn(b'Apteringsforslag', good.stdout)
        for args in [('--format','dump'), ('--format','dump','--entry','navigation'),
                     ('--format','dump','--entry','bucking','--columns','20')]:
            result = run(*args)
            self.assertEqual(result.returncode, 2, result.stderr)
            self.assertEqual(result.stdout, b'')
            self.assertIn('error', json.loads(result.stderr))
