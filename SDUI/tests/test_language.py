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
from sdui.formatting import formatting
from sdui.normalize import normalize

MINIMAL = 'sdui 0.2; Page = [<greeting="Hello">];'
EXAMPLE = (ROOT / 'examples/main-page.sdui').read_text()


class LanguageTest(unittest.TestCase):
    def valid(self, source):
        tree = parse(source)
        validate(tree)
        return tree

    def invalid(self, source, code=None):
        with self.assertRaises(SduiError) as caught:
            self.valid(source)
        if code:
            self.assertEqual(caught.exception.code, code, str(caught.exception))
        self.assertGreaterEqual(caught.exception.span.start, 0)
        return caught.exception

    def test_ported_example_bindings_and_instance_paths(self):
        tree = self.valid(EXAMPLE)
        self.assertEqual(tree.references[0].path, 'some_SDL_file.sdl')
        conn = tree.connections[0]
        self.assertEqual(conn.path, ('top', 'rightTop', 'input1_boxui'))
        self.assertEqual(conn.object, 'input1_sdl')
        root = normalize(tree)['BoxUIDefinition']
        field = root.rows[0][0].rows[0][1].rows[0][0].rows[1][0]
        callback = dict(field.arguments)['callback']
        self.assertIsInstance(callback, Reference)
        self.assertEqual(callback.module, 'sdlFile')
        self.assertEqual(callback.member, 'callback')

    def test_nested_groups_and_local_rows_preserved(self):
        tree = self.valid('sdui 0.2; P=[<<"A", "B"> {v<}; <"C", <"D"> {font=12}>> {<->}];')
        outer = tree.definitions[0].root.rows[0].items[0]
        self.assertEqual([len(r.items) for r in outer.rows], [1, 1])
        self.assertEqual(len(outer.rows[0].items[0].rows[0].items), 2)
        inner = outer.rows[1].items[0].rows[0].items[1]
        self.assertEqual(inner.kind, 'group')
        self.assertEqual(formatting(inner), {'font': 12})
        resolved = normalize(tree)['P'].rows[0][0]
        self.assertEqual(resolved.rows[1][0].rows[0][1].kind, 'group')
        self.assertNotEqual(resolved.rows[0][0].path, resolved.rows[1][0].path)

    def test_formatting_binds_before_separator(self):
        tree = self.valid('sdui 0.2; G=<"a">; P=[G {font=12}, <ok=button("OK") {font=10}> {>-<}]*b {<->};')
        children = tree.definitions[1].root.rows[0].items
        self.assertEqual(children[0].kind, 'use')
        self.assertEqual(formatting(children[0]), {'font': 12})
        self.assertEqual(formatting(children[1]), {'x': 'content'})
        self.invalid('sdui 0.2; P=["A", {font=12} "B"];', 'syntax')
        self.invalid('sdui 0.2; P=["A"; {font=12}];', 'syntax')
        self.invalid('sdui 0.2; P=[] {,};', 'syntax')
        self.invalid('sdui 0.2; P=[] {font=10,};', 'syntax')
        self.invalid('sdui 0.2; P=[] {<->}*b;', 'syntax')

    def test_regions_and_empty_containers(self):
        tree = self.valid('sdui 0.2; P=[header="## Header", body=<>, footer=[]]*box;')
        root = normalize(tree)['P']
        self.assertEqual([key for key, _ in root.regions], ['header', 'body', 'footer'])
        self.assertEqual(root.rows, ())
        self.invalid('sdui 0.2; P=[header="a", header="b"];', 'duplicate-region')
        self.invalid('sdui 0.2; P=[body=<>, "extra"];', 'body-conflict')
        self.invalid('sdui 0.2; P=[<[]>];', 'group-content')
        self.invalid('sdui 0.2; F=[]; P=[<F>];', 'group-content')

    def test_references_forward_reuse_and_cycles(self):
        tree = self.valid('sdui 0.2; P=[left=G {font=12}, right=G]; G=<ok=button("OK")>;')
        root = normalize(tree)['P']
        left, right = root.rows[0]
        self.assertEqual((left.rows[0][0].path, right.rows[0][0].path), ('P/left/ok', 'P/right/ok'))
        self.assertEqual(dict(left.layout)['font'], 12)
        self.assertNotIn('font', dict(right.layout))
        self.invalid('sdui 0.2; P=[missing];', 'unknown-definition')
        self.invalid('sdui 0.2; A=<B>; B=<A>;', 'reference-cycle')
        self.invalid('sdui 0.2; G=<>; P=[G,G];', 'duplicate-node')

    def test_static_comments_quotes_and_unicode(self):
        tree = self.valid((ROOT / 'examples/static.sdui').read_text())
        rows = tree.definitions[0].root.rows[0].items[1].rows
        self.assertEqual(rows[0].items[0].text.value, 'Hello # this is text, not a comment')
        self.assertEqual(rows[1].items[0].text.value, "Quoted: 'Hei' and Unicode: æ")
        self.assertEqual(len(rows), 3)

    def test_raw_markdown_and_utf8_spans(self):
        source = 'sdui 0.2;\r\nP=["""# Æøå\n  ```mermaid\ngraph LR; A-->B\n  ```\n"""];'
        tree = self.valid(source)
        text = tree.definitions[0].root.rows[0].items[0].text
        self.assertTrue(text.value.startswith('# Æøå\n  ```'))
        self.assertEqual(text.span.line, 2)
        self.assertEqual(text.span.column, 4)
        self.assertEqual(source.encode()[text.span.start:text.span.end].decode(), '"""' + text.value + '"""')
        err = self.invalid(source + '?', 'lexical')
        self.assertEqual(err.span.start, len(source.encode()))

    def test_version_is_exact_no_legacy(self):
        for version in ['0.1', '0.20', '2e-1', '-0.2']:
            self.invalid(MINIMAL.replace('0.2', version), 'version')
        self.invalid('sdui 0.2; P={main=[{w=text("old")} ]};', 'syntax')
        self.invalid(MINIMAL + 'garbage', 'syntax')

    def test_local_symbol_errors(self):
        self.invalid(EXAMPLE.replace('callback=sdlFile.', 'callback=missing.', 1), 'unknown-module')
        self.invalid(EXAMPLE.replace('BoxUIDefinition.top.rightTop.input1_boxui);', 'Unknown.input1_boxui);'), 'unknown-definition')
        self.invalid(EXAMPLE.replace('BoxUIDefinition.top.rightTop.input1_boxui);', 'BoxUIDefinition.top);'), 'unknown-widget')
        self.invalid(EXAMPLE.replace('button2=', 'button1='), 'duplicate-node')
        self.invalid(EXAMPLE + 'sdlFile.input1_sdl.setHandle(BoxUIDefinition.top.rightTop.input1_boxui);', 'duplicate-connection')
        self.invalid(EXAMPLE.replace('ref: sdlFile', 'ref: sdlFile "x"; ref: sdlFile', 1), 'duplicate-module')
        self.invalid(MINIMAL + MINIMAL.split(';', 1)[1], 'duplicate-definition')

    def test_widget_contracts(self):
        for code, item in [
            ('widget-kind', 'w=unknown("x")'), ('widget-kind', 'w=text("old")'),
            ('widget-property', 'w=button("x", typo="y")'),
            ('duplicate-argument', 'w=button("x", label="y")'),
            ('argument-order', 'w=button(label="x", "y")'),
            ('argument-type', 'w=input(text=12)'),
            ('argument-type', 'w=button("OK", callback="execute")'),
            ('missing-argument', 'w=button()'), ('missing-argument', 'w=svg(label="x")'),
        ]:
            with self.subTest(item=item):
                self.invalid('sdui 0.2; P=[<' + item + '>];', code)

    def test_anonymous_unbound_widgets_and_wrap_constraints(self):
        self.valid('sdui 0.2; P=[<"# Markdown"; button("OK"), button("Cancel")>]*b;')
        self.invalid('sdui 0.2; ref: m "missing.sdl"; P=[button("OK", callback=m.o.@run)];', 'widget-name')
        self.invalid('sdui 0.2; P=[<"A" {x=1fr}, "B"> {wrap=wrap}];', 'wrap-layout')

        self.invalid('sdui 0.2; G=<"A" {x=1fr}>; P=[G {wrap=wrap}];', 'wrap-layout')

    def test_shapes_ratio_and_relative_dimensions(self):
        for a, b in [('^<', '<^'), ('>^', '^>'), ('v<', '<v'), ('>v', 'v>'), ('->', '>-')]:
            left = self.valid('sdui 0.2; P=[] {' + a + '};').definitions[0].root
            right = self.valid('sdui 0.2; P=[] {' + b + '};').definitions[0].root
            self.assertEqual(formatting(left), formatting(right))
        self.valid('sdui 0.2; P=[] {16:9, <->, font=10};')
        self.valid('sdui 0.2; P=[] {4:3, scale-y=0.5};')
        self.valid('sdui 0.2; P=[] {scale=0.5, padding=(0.01,0.02,0.01,0.02)};')
        for rules, code in [('16:9,scale=1', 'ratio-axis'), ('16:9,<->,^|v', 'ratio-axis'),
                            ('16:9,x=1fr', 'ratio-axis'), ('0:9', 'layout-type'),
                            ('font=0', 'layout-type'), ('font=1fr', 'layout-type'),
                            ('width=100', 'layout-property'), ('x=100', 'layout-type'),
                            ('min-x=0.8,max-x=0.2', 'size-range'),
                            ('<->,>-<', 'layout-conflict'), ('gap=-1', 'layout-type'),
                            ('<v,v<', 'layout-conflict'), ('heading="old"', 'layout-property')]:
            with self.subTest(rules=rules):
                self.invalid('sdui 0.2; P=[] {' + rules + '};', code)
        self.invalid('sdui 0.2; P=<"A"> {16:9};', 'ratio-scope')
        self.invalid('sdui 0.2; G=<"A">; P=[G {16:9}];', 'ratio-scope')

    def test_no_arbitrary_code_or_empty_rows(self):
        for suffix in ['sdlFile.run();', 'sdlFile.input.value("x");', '__import__("os");']:
            self.invalid(MINIMAL + suffix, 'syntax')
        for body in ['"x" + "y"', 'ok=button("x",)', '"x";;', '"x",']:
            self.invalid('sdui 0.2; P=[<' + body + '>];')

    def test_lexical_failures(self):
        for text in ['"unfinished', '"bad\\q"', '"\\uD800"', '"\\u0000"', '"a\nb"', '1e999', '"\\uZZZZ"']:
            self.invalid(MINIMAL.replace('"Hello"', text), 'lexical')
        self.invalid(MINIMAL + '\ud800', 'encoding')

    def test_bounded_work(self):
        self.invalid('#' + 'x' * 262144, 'source-limit')
        self.invalid(';' * 50001, 'token-limit')
        self.invalid('sdui 0.2; P=' + '[' * 65 + ']' * 65 + ';', 'depth-limit')
        many = ','.join(f'w{i}=button("x")' for i in range(2048))
        self.invalid('sdui 0.2; P=[<' + many + '>];', 'node-limit')
        self.invalid('sdui 0.2; P=[w=button(' + ','.join('"x"' for _ in range(33)) + ')];', 'argument-limit')
        definitions = 'A0=<"x">;' + ''.join(f'A{i}=<left=A{i-1},right=A{i-1}>;' for i in range(1, 15))
        self.invalid('sdui 0.2;' + definitions, 'expansion-limit')

    def test_deterministic_ast_and_syntax_semantics_separation(self):
        data = to_data(parse(EXAMPLE))
        self.assertEqual(data, to_data(parse(EXAMPLE)))
        self.assertEqual(json.loads((ROOT / 'examples/main-page.ast.json').read_text())['document'], data)
        tree = parse('sdui 0.2; P=[w=futureWidget()];')
        self.assertEqual(tree.definitions[0].root.rows[0].items[0].widget, 'futureWidget')
        with self.assertRaises(SduiError):
            validate(tree)

    def test_truncations_and_malformed_inputs_are_structured(self):
        for offset in range(1, len(MINIMAL)):
            self.invalid(MINIMAL[:offset])
        randomizer = random.Random(42)
        for _ in range(500):
            text = ''.join(randomizer.choices('[]{}<>*^|(),;=.:@abc012#"\n ', k=randomizer.randrange(1, 200)))
            try:
                parse(text)
            except SduiError:
                pass

    def test_cli_source_protection_and_exit_codes(self):
        env = dict(os.environ, PYTHONPATH=str(ROOT / 'src'))
        run = lambda data, *args: subprocess.run([sys.executable, '-m', 'sdui', '-', *args],
            input=data, capture_output=True, env=env)
        result = run(EXAMPLE.encode())
        self.assertEqual(result.returncode, 0, result.stderr)
        self.assertEqual(json.loads(result.stdout)['astFormat'], 'sdui-ast/0.2')
        for data, code in [(b'\xff', 'encoding'), (b'sdui', 'syntax'), (b'x' * 262145, 'source-limit')]:
            result = run(data)
            self.assertEqual(result.returncode, 2)
            self.assertEqual(result.stdout, b'')
            self.assertEqual(json.loads(result.stderr)['error']['code'], code)
        result = run(b'sdui 0.2; P=[w=futureWidget()];', '--syntax-only')
        self.assertEqual(result.returncode, 0)
        with tempfile.TemporaryDirectory() as directory:
            source = Path(directory)/'example.sdui'
            source.write_text(MINIMAL)
            result = subprocess.run([sys.executable, '-m', 'sdui', str(source), '-o', str(source)], capture_output=True, env=env)
            self.assertEqual(result.returncode, 3)
            self.assertEqual(source.read_text(), MINIMAL)
