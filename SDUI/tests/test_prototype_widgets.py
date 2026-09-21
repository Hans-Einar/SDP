from pathlib import Path
import sys
import unittest
import xml.etree.ElementTree as ET

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROOT/'src'))
from sdui import parse
from sdui.normalize import normalize
from sdui.prototype_widgets import PrototypeWidget, prototype_widgets, svg_widget, html_widget
from sdui.prototype_html import widget_gallery_html


class PrototypeWidgetTest(unittest.TestCase):
    def controls(self, body):
        return prototype_widgets(normalize(parse('sdui 0.2; P=' + body + ';'))['P'])

    def test_same_source_controls_are_used_and_callbacks_are_not_transferred(self):
        source = 'sdui 0.2; ref: domain "missing.sdl"; P=[ok=button("OK", callback=domain.o.@run), value=input("Navn", value="Æøå")];'
        controls = prototype_widgets(normalize(parse(source))['P'])
        self.assertEqual(controls['P/ok'].kind, 'button')
        self.assertEqual(controls['P/value'].value, 'Æøå')
        self.assertFalse(hasattr(controls['P/ok'], 'callback'))
        html = widget_gallery_html(list(controls.values()))
        self.assertNotIn('domain.o', html)
        self.assertNotIn('missing.sdl', html)

    def test_visibility_and_disabled_ancestor(self):
        controls = self.controls('[<ok=button("OK"), field=input("Navn")> {enabled=false}; hidden=button("Hidden") {visible=false}]')
        self.assertEqual(len(controls), 2)
        self.assertTrue(all(not x.enabled for x in controls.values()))
        self.assertEqual(next(x.value for x in controls.values() if x.kind=='input'), '')
        self.assertIn(' disabled', html_widget(next(iter(controls.values())), 1))

    def test_svg_states_and_field_values(self):
        b = PrototypeWidget('P/ok', 'button', 'OK', '', True)
        normal, _ = svg_widget(b, 0, 0, 120)
        pressed, _ = svg_widget(b, 0, 0, 120, pressed=True)
        self.assertNotEqual(normal, pressed)
        self.assertIn('y="3"', pressed)
        disabled = PrototypeWidget('P/no', 'button', 'No', '', False)
        self.assertEqual(svg_widget(disabled,0,0,120), svg_widget(disabled,0,0,120,pressed=True,focused=True))
        field = PrototypeWidget('P/name', 'input', 'Navn', 'Ola', True)
        svg, height = svg_widget(field,0,0,200,focused=True)
        self.assertEqual(height, 66)
        self.assertIn('Ola', svg)
        self.assertIn('stroke-width="2"', svg)
        with self.assertRaises(ValueError):
            svg_widget(field, 0, 0, 20)
        with self.assertRaises(ValueError):
            svg_widget(field, float('nan'), 0, 120)

    def test_literal_labels_and_values_do_not_become_markup(self):
        value = '<script>alert("x")</script> & "value"'
        w = PrototypeWidget('P/"name', 'input', '<input>', value, True)
        svg = svg_widget(w, 0, 0, 250)[0]
        ET.fromstring(svg)
        self.assertNotIn('<script>', svg)
        html = html_widget(w, 0)
        self.assertIn('&lt;script&gt;', html)
        self.assertNotIn('<script>', html)
        self.assertIn('value="&lt;script&gt;', html)

    def test_gallery_has_print_mirrors_and_no_submit_or_external_resources(self):
        controls = self.controls('[ok=button("OK"), field=input("Navn", value="Ola")]')
        html = widget_gallery_html(list(controls.values()))
        self.assertIn('beforeprint', html)
        self.assertIn('textContent = field.querySelector(\'input\').value', html)
        self.assertIn('.print-value { display:block;', html)
        self.assertIn('.sdui-field input { display:none }', html)
        self.assertIn('type="button"', html)
        for unwanted in ['type="submit"', '<form', 'fetch(', '<script src=', 'localStorage']:
            self.assertNotIn(unwanted, html)

    def test_generated_preview_shapes_are_static_and_have_real_source_values(self):
        namespace = {'s': 'http://www.w3.org/2000/svg'}
        for name in ['prototype-controls.svg', 'concept1-bucking.widgets.svg']:
            tree = ET.parse(ROOT/'examples'/name)
            self.assertTrue(tree.findall('.//s:svg[@data-widget="button"]', namespace))
            fields = tree.findall('.//s:svg[@data-widget="input"]', namespace)
            self.assertEqual(len(fields), 3)
            self.assertIn('Ola Nordmann', ''.join(tree.getroot().itertext()))
            self.assertFalse(tree.findall('.//s:script', namespace))
            self.assertFalse(tree.findall('.//s:foreignObject', namespace))
