import unittest
from viewpoints import Views, core
from test_data_core import SOURCE


class DataViewTests(unittest.TestCase):
    def test_packet_uses_exact_ranges_and_placement_source_facts(self):
        views = Views(core.canonicalize(core.parse(SOURCE)), ['VP09', 'VP10'])
        packet = next(d for d in views.diagrams if d.kind == 'packet')
        self.assertEqual(packet.mermaid(views.kinds), 'packet\n    0-15: "Revision"\n    16-47: "Payload"\n')
        facts = {r['id']: r for r in views.facts}
        for element in packet.elements:
            record = facts[element['fact']]
            self.assertEqual(element['first_bit'], record['offset'])
            self.assertEqual(element['last_bit'], record['offset'] + record['width'] - 1)
        self.assertIn('Project | State | Updates', views.markdown())
        self.assertIn('Archive', views.markdown())

    def test_missing_encoding_gives_gap_without_fake_packet(self):
        text = '\n'.join(line for line in SOURCE.splitlines() if not line.startswith(('Wire ', 'encoding Wire'))) + '\n'
        views = Views(core.canonicalize(core.parse(text)), ['VP10'])
        self.assertFalse(views.diagrams)
        self.assertIn('NO_ENCODING', {g['code'] for g in views.gaps})

    def test_changed_width_changes_packet_and_invalid_layout_is_rejected(self):
        text = SOURCE.replace('at 16 bits 32', 'at 16 bits 64')
        self.assertIn('16-79', Views(core.canonicalize(core.parse(text)), ['VP10']).markdown())
        with self.assertRaises(ValueError):
            Views(SOURCE.replace('at 16 bits 32', 'at 8 bits 32'))

    def test_every_projection_argument_and_property_remains_in_inventory(self):
        views = Views(core.canonicalize(core.parse(SOURCE)), ['VP11'])
        for statement in views.model.statements:
            self.assertIn(core.sentence(statement), views.markdown())


if __name__ == '__main__':
    unittest.main()
