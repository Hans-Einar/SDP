import json
import os
from pathlib import Path
import subprocess
import sys
import unittest

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROOT / 'src'))
from sdui import parse
from sdui.normalize import normalize
from sdui.markdown_dump import markdown_dump, content_lines, code_span


class MarkdownDumpTest(unittest.TestCase):
    def dump(self, component):
        root = normalize(parse('sdui 0.2; P=' + component + ';'))['P']
        return markdown_dump(root, 80)

    def test_markdown_is_renderable_source_not_escaped_text(self):
        source = '## Title\n\n**bold**\n\n- One\n- Two\n\n| A | B |\n| - | - |\n| 1 | 2 |'
        text = self.dump('[' + json.dumps(source) + ']*b')
        content = text.split('## Innhold\n', 1)[1]
        self.assertIn('> > ## Title', content)
        self.assertIn('> > **bold**', content)
        self.assertIn('> > - One', content)
        self.assertIn('> > | 1 | 2 |', content)
        self.assertNotIn('\\*\\*bold', content)

    def test_mermaid_is_omitted_and_other_code_is_preserved(self):
        source = '```mermaid\nSECRET-->B\n```\n\n```python\nprint("kept")\n```'
        text = self.dump('[' + json.dumps(source) + ']')
        self.assertNotIn('SECRET', text)
        self.assertIn('[Mermaid utelatt]', text)
        self.assertIn('> > ```python', text)
        self.assertIn('> > print("kept")', text)

    def test_fence_containment_and_literal_widget_labels(self):
        self.assertEqual(content_lines('```python\nx=1'), ['```python', 'x=1', '```'])
        self.assertEqual(content_lines('~~~text\nx\n~~~~'), ['~~~text', 'x', '~~~~'])
        text = self.dump('["````text\\nkept\\n````"; button("`<button>x</button>`")]*b')
        self.assertIn('`````text\nSDUI GUI dump', text)
        self.assertIn('**Knapp:** `` `<button>x</button>` ``', text)
        self.assertEqual(code_span('`x``y`'), '``` `x``y` ```')
        self.assertNotIn('\x1b', code_span('\x1b[31m'))

    def test_regions_rows_visibility_and_static_widgets(self):
        text = self.dump('[header="HEADER", <button("OK"), input("Name", value="Ada")>; "HIDDEN" {visible=false}, footer="FOOTER"]')
        content = text.split('## Innhold\n', 1)[1]
        self.assertLess(content.index('HEADER'), content.index('**Knapp:**'))
        self.assertLess(content.index('**Inndata:**'), content.index('FOOTER'))
        self.assertIn('2 komponenter fra venstre mot høyre', content)
        self.assertNotIn('HIDDEN', text)
        self.assertIn('` Ada `', content)
        self.assertNotIn('<input', text)
        self.assertNotIn('<button', text)

    def test_concept1_snapshot(self):
        tree = parse((ROOT / 'examples/concept1-bucking.sdui').read_text())
        text = markdown_dump(normalize(tree)['bucking'])
        self.assertEqual(text, (ROOT / 'examples/concept1-bucking.dump.md').read_text())
        self.assertIn('> > > > | 1  | Sagtømmer | 430 cm | 26 cm  |', text)
        self.assertNotIn('Rot --> Kapp1', text)

    def test_cli_stdout_output_and_errors(self):
        env = dict(os.environ, PYTHONPATH=str(ROOT / 'src'))
        cmd = [sys.executable, '-m', 'sdui', str(ROOT / 'examples/concept1-bucking.sdui'),
               '--format', 'markdown']
        good = subprocess.run(cmd + ['--entry', 'bucking'], capture_output=True, env=env)
        self.assertEqual(good.returncode, 0, good.stderr)
        self.assertEqual(good.stdout, (ROOT / 'examples/concept1-bucking.dump.md').read_bytes())
        for args in [[], ['--entry', 'missing'], ['--entry', 'bucking', '--columns', '20'],
                     ['--entry', 'bucking', '--syntax-only']]:
            result = subprocess.run(cmd + args, capture_output=True, env=env)
            self.assertEqual(result.returncode, 2, result.stderr)
            self.assertEqual(result.stdout, b'')
