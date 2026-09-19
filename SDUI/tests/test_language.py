import json
import os
from pathlib import Path
import random
import subprocess
import sys
import tempfile
import unittest

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROOT / 'src'))
from sdui import SduiError, parse, to_data, validate
from sdui.ast import Reference

MINIMAL = 'sdui 0.1; Page = {main = [{greeting = text("Hello")}]};'
EXAMPLE = (ROOT / 'examples/main-page.sdui').read_text(encoding='utf-8')


class LanguageTest(unittest.TestCase):
    def valid(self, source):
        tree = parse(source)
        validate(tree)
        return tree

    def invalid(self, source, code=None):
        with self.assertRaises(SduiError) as caught:
            self.valid(source)
        if code:
            self.assertEqual(caught.exception.code, code)
        self.assertGreaterEqual(caught.exception.span.start, 0)
        return caught.exception

    def test_owner_example_preserves_boxes_rows_and_bindings(self):
        tree = self.valid(EXAMPLE)
        self.assertEqual(tree.references[0].path, 'some_SDL_file.sdl')
        definition = tree.definitions[0]
        self.assertEqual(definition.name, 'BoxUIDefinition')
        self.assertEqual([x.name for x in definition.root.children], ['top', 'svgBox'])
        left, right = definition.root.children[0].children
        self.assertEqual([w.name for w in left.content.rows[0].widgets], ['button1', 'button2'])
        self.assertEqual(len(right.content.rows), 2)
        field = right.content.rows[1].widgets[0]
        self.assertEqual(field.name, 'input1_boxui')
        callback = field.arguments[-1].value
        self.assertIsInstance(callback, Reference)
        self.assertEqual((callback.module, callback.object, callback.member), ('sdlFile', 'input1_sdl', 'callback'))
        self.assertEqual((tree.connections[0].definition, tree.connections[0].widget), ('BoxUIDefinition', 'input1_boxui'))
        self.assertEqual(tree.connections[0].object, 'input1_sdl')

    def test_static_comments_quotes_and_unicode(self):
        tree = self.valid((ROOT / 'examples/static.sdui').read_text())
        rows = tree.definitions[0].root.content.rows
        self.assertEqual(rows[0].widgets[0].arguments[0].value.value, 'Hello # this is text, not a comment')
        self.assertEqual(rows[1].widgets[0].arguments[0].value.value, "Quoted: 'Hei' and Unicode: æ")
        self.assertEqual(len(rows), 3)

    def test_utf8_byte_spans_and_crlf(self):
        source = MINIMAL.replace('Hello', 'Æøå').replace('Page', '\r\nPage')
        tree = self.valid(source)
        literal = tree.definitions[0].root.content.rows[0].widgets[0].arguments[0].value
        self.assertEqual(source.encode()[literal.span.start:literal.span.end].decode(), '"Æøå"')
        self.assertEqual(literal.span.line, 2)
        self.assertEqual(literal.span.column, source.split('\n')[1].index('"') + 1)
        error = self.invalid(source.replace('Æøå', 'Æøå') + '?', 'lexical')
        self.assertEqual(error.span.start, len(source.encode()))

    def test_anonymous_group_and_reusable_names_across_definitions(self):
        tree = self.valid('sdui 0.1; A={main=[axis="row", [child=[{a=text("A")}]]]}; B={main=[{a=text("B")}]};')
        self.assertIsNone(tree.definitions[0].root.children[0].name)

    def test_version_is_exact_and_trailing_input_rejected(self):
        for version in ['0.2', '0.10', '1e-1', '-0.1']:
            self.invalid(MINIMAL.replace('0.1', version), 'version')
        self.invalid(MINIMAL + 'garbage', 'syntax')
        self.invalid(MINIMAL.replace('sdui 0.1;', ''), 'syntax')

    def test_local_symbol_errors(self):
        self.invalid(EXAMPLE.replace('callback=sdlFile.', 'callback=missing.', 1), 'unknown-module')
        self.invalid(EXAMPLE.replace('BoxUIDefinition.input1_boxui);', 'Unknown.input1_boxui);'), 'unknown-definition')
        self.invalid(EXAMPLE.replace('BoxUIDefinition.input1_boxui);', 'BoxUIDefinition.main);'), 'unknown-widget')
        self.invalid(EXAMPLE.replace('button2 =', 'button1 ='), 'duplicate-node')
        self.invalid(EXAMPLE + 'sdlFile.input1_sdl.setHandle(BoxUIDefinition.input1_boxui);', 'duplicate-connection')
        self.invalid(EXAMPLE.replace('ref: sdlFile', 'ref: sdlFile "x"; ref: sdlFile', 1), 'duplicate-module')
        self.invalid(MINIMAL + MINIMAL.split(';', 1)[1], 'duplicate-definition')

    def test_widget_contracts(self):
        for code, contents in [
            ('widget-kind', 'w=unknown("x")'),
            ('widget-property', 'w=text(text="x", typo="y")'),
            ('duplicate-argument', 'w=text("x", text="y")'),
            ('argument-order', 'w=text(text="x", "y")'),
            ('argument-type', 'w=input(text=12)'),
            ('argument-type', 'w=button("OK", callback="execute")'),
            ('missing-argument', 'w=text()'),
            ('missing-argument', 'w=svg(label="x")'),
        ]:
            with self.subTest(contents=contents):
                self.invalid('sdui 0.1; Page={main=[{' + contents + '}]};', code)

    def test_box_properties_and_layout_constraints(self):
        for props, code in [('weight=0,', 'property-type'), ('min=10,max=2,', 'size-range'),
                            ('title=4,', 'property-type'), ('handle="main",', 'box-property'),
                            ('title="A",title="B",', 'duplicate-property'),
                            ('axis="row",', 'axis-scope')]:
            with self.subTest(props=props):
                self.invalid(MINIMAL.replace('main = [', 'main = [' + props), code)
        self.invalid('sdui 0.1; Page={main=[]};', 'box-body')
        self.invalid('sdui 0.1; Page={main=[child=[{x=text("x")}], {y=text("y")}]};', 'box-body')

    def test_no_arbitrary_code_or_ambiguous_shorthand(self):
        for suffix in ['sdlFile.run();', 'sdlFile.input.value("x");', '__import__("os");']:
            self.invalid(MINIMAL + suffix, 'syntax')
        self.invalid(EXAMPLE.replace('sdlFile.button1_sdl.@callback', 'sdlFile.@button1_sdl.@callback'), 'syntax')
        self.invalid(MINIMAL.replace('text("Hello")', 'text("x" + "y")'), 'lexical')
        self.invalid(MINIMAL.replace('text("Hello")', 'text("x",)'), 'syntax')
        self.invalid(MINIMAL.replace('text("Hello")', 'text("x");;'), 'syntax')

    def test_lexical_failures(self):
        for text in ['"unfinished', '"bad\\q"', '"\\uD800"', '"\\u0000"', '"a\nb"', '1e999', '"\\uZZZZ"']:
            self.invalid(MINIMAL.replace('"Hello"', text), 'lexical')
        self.invalid(MINIMAL + '\ud800', 'encoding')

    def test_bounded_work(self):
        self.invalid('#' + 'x' * 262144, 'source-limit')
        self.invalid(';' * 50001, 'token-limit')
        self.invalid('sdui 0.1; P={root=' + '[' * 65 + ']' * 65 + '};', 'depth-limit')
        many = ','.join(f'w{i}=text("x")' for i in range(2048))
        self.invalid('sdui 0.1; P={root=[{' + many + '}]};', 'node-limit')
        self.invalid(MINIMAL.replace('"Hello"', ','.join('"x"' for _ in range(33))), 'argument-limit')

    def test_deterministic_ast_and_syntax_semantics_separation(self):
        data = to_data(parse(EXAMPLE))
        self.assertEqual(data, to_data(parse(EXAMPLE)))
        golden = json.loads((ROOT / 'examples/main-page.ast.json').read_text())
        self.assertEqual(golden['document'], data)
        unknown = MINIMAL.replace('text(', 'futureWidget(')
        tree = parse(unknown)
        self.assertEqual(tree.definitions[0].root.content.rows[0].widgets[0].kind, 'futureWidget')
        with self.assertRaises(SduiError):
            validate(tree)

    def test_truncations_and_malformed_inputs_have_structured_errors(self):
        for offset in range(1, len(MINIMAL)):
            with self.subTest(offset=offset):
                self.invalid(MINIMAL[:offset])
        randomizer = random.Random(42)
        for _ in range(300):
            text = ''.join(randomizer.choices('[]{}(),;=.:@abc012#"\n ', k=randomizer.randrange(1, 200)))
            try:
                parse(text)
            except SduiError:
                pass  # IndexError/RecursionError/ValueError would fail this test.

    def test_cli_files_and_source_protection(self):
        env = dict(os.environ, PYTHONPATH=str(ROOT / 'src'))
        with tempfile.TemporaryDirectory() as directory:
            source = Path(directory) / 'example.sdui'
            output = Path(directory) / 'result.json'
            source.write_text(MINIMAL)
            cmd = [sys.executable, '-m', 'sdui', str(source), '-o', str(output)]
            result = subprocess.run(cmd, capture_output=True, env=env)
            self.assertEqual(result.returncode, 0, result.stderr)
            self.assertEqual(json.loads(output.read_text())['document']['profile'], 'sdui/0.1')
            result = subprocess.run(cmd[:-1] + [str(source)], capture_output=True, env=env)
            self.assertEqual(result.returncode, 3)
            self.assertEqual(source.read_text(), MINIMAL)
            result = subprocess.run([sys.executable, '-m', 'sdui', str(source) + '.missing'],
                                    capture_output=True, env=env)
            self.assertEqual(result.returncode, 3)

    def test_cli_json_exit_codes_and_no_module_loading(self):
        env = dict(os.environ, PYTHONPATH=str(ROOT / 'src'))
        run = lambda data, *args: subprocess.run([sys.executable, '-m', 'sdui', '-', *args],
            input=data, capture_output=True, env=env)
        result = run(EXAMPLE.encode())  # Referenced .sdl file deliberately does not exist.
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(json.loads(result.stdout)['astFormat'], 'sdui-ast/0.1')
        for data, code in [(b'\xff', 'encoding'), (b'sdui', 'syntax'), (b'x' * 262145, 'source-limit')]:
            result = run(data)
            self.assertEqual(result.returncode, 2)
            self.assertEqual(result.stdout, b'')
            self.assertEqual(json.loads(result.stderr)['error']['code'], code)
        result = run(MINIMAL.replace('text(', 'newWidget(').encode(), '--syntax-only')
        self.assertEqual(result.returncode, 0)
        self.assertEqual(json.loads(result.stdout)['validation'], 'syntax-only')


if __name__ == '__main__':
    unittest.main()
