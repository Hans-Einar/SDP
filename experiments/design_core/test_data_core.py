import subprocess
import sys
import unittest
import design_core as dc

SOURCE = '''language design-core version 0.3.
unit Store.
database Archive.
dataset State.
contract Record.
field Value.
datagram Updates.
contract Family.
variant Changed.
field Revision.
field Payload.
encoding Wire.
functionality Project.
Store owns Archive.
Archive holds State.
State upholds Record.
Record has completeness = closed.
Record has-field Value.
Value has value-type = unsigned.
Value has presence = optional.
Updates upholds Family.
Updates from State.
Family has completeness = closed.
Family defines Changed.
Family has-field Revision.
Changed has-field Payload.
Revision has value-type = unsigned.
Revision has presence = required.
Payload has value-type = bytes.
Payload has presence = required.
Wire encodes Changed.
Wire has byte-order = big-endian.
Wire has bit-order = most-significant-first.
Wire places Revision at 0 bits 16.
Wire places Payload at 16 bits 32.
Store owns Project.
Project projects State into Updates.
'''


class DataTests(unittest.TestCase):
    def errors(self, text):
        return [d.code for d in dc.validate(dc.parse(text))]

    def test_contract_projection_packet_roundtrip_and_cli(self):
        canonical = dc.canonicalize(dc.parse(SOURCE))
        self.assertEqual(dc.check(canonical)[1], [])
        self.assertEqual(dc.canonicalize(dc.parse(canonical)), canonical)
        result = subprocess.run([sys.executable, dc.__file__, 'check', '-'],
                                input=canonical, capture_output=True, text=True)
        self.assertEqual(result.returncode, 0, result.stdout + result.stderr)

    def test_projection_and_placement_spans(self):
        model = dc.parse(SOURCE)
        for s in model.statements:
            if isinstance(s, (dc.Projection, dc.Placement)):
                self.assertEqual(SOURCE[s.span.start:s.span.end], dc.sentence(s))
            if isinstance(s, dc.Placement):
                self.assertEqual(SOURCE[s.width.span.start:s.width.span.end], str(s.width.value))

    def test_database_is_optional_but_holder_and_contract_are_explicit(self):
        source = SOURCE.replace('database Archive.\n', '').replace('Store owns Archive.\n', '').replace('Archive holds State.', 'Store holds State.')
        self.assertEqual(self.errors(source), [])
        for fact in ['Store holds State.\n', 'State upholds Record.\n', 'Updates from State.\n']:
            self.assertIn('DATA_CARDINALITY', self.errors(source.replace(fact, '')))

    def test_wrong_projection_source_and_wrong_types(self):
        source = SOURCE.replace('dataset State.', 'dataset State.\ndataset Other.') + 'Store holds Other.\nOther upholds Record.\n'
        self.assertIn('PROJECTION_SOURCE_MISMATCH', self.errors(source.replace('projects State', 'projects Other')))
        self.assertIn('OBJECT_TYPE_MISMATCH', self.errors(SOURCE.replace('Updates from State.', 'Updates from Archive.')))

    def test_field_and_variant_ownership_are_unique(self):
        for extra in ['Record has-field Payload.', 'Record defines Changed.']:
            self.assertIn('DATA_CARDINALITY', self.errors(SOURCE + extra))

    def test_open_contract_cannot_produce_packet(self):
        self.assertIn('OPEN_ENCODING', self.errors(SOURCE.replace('Family has completeness = closed', 'Family has completeness = open')))
        no_wire = '\n'.join(s for s in SOURCE.splitlines() if not s.startswith(('Wire ', 'encoding Wire'))) + '\n'
        self.assertEqual(self.errors(no_wire.replace('Family has completeness = closed', 'Family has completeness = open')), [])

    def test_packet_rejects_gaps_overlap_zero_and_duplicate_placement(self):
        for replacement in ['at 15 bits 32', 'at 17 bits 32', 'at 16 bits 0']:
            self.assertIn('ENCODING_RANGE', self.errors(SOURCE.replace('at 16 bits 32', replacement)))
        self.assertIn('ENCODING_FIELDS', self.errors(SOURCE + 'Wire places Payload at 48 bits 32.'))
        self.assertIn('ENCODING_FIELDS', self.errors(SOURCE.replace('Wire places Revision at 0 bits 16.\n', '')))

    def test_packet_requires_supported_types_presence_and_order(self):
        for old, new in [('value-type = bytes', 'value-type = text'),
                         ('value-type = bytes', 'value-type = boolean'),
                         ('Payload has presence = required', 'Payload has presence = optional'),
                         ('at 16 bits 32', 'at 16 bits 31')]:
            self.assertIn('ENCODING_TYPE', self.errors(SOURCE.replace(old, new)))
        self.assertIn('MISSING_CONTRACT_PROPERTY', self.errors(SOURCE.replace('Wire has byte-order = big-endian.\n', '')))

    def test_bad_numeric_syntax_has_structured_failure(self):
        for value in ['-1', '1.5', '01', '999999999999999999999999']:
            with self.assertRaises(dc.ParseError):
                dc.parse(SOURCE.replace('at 16 bits', 'at ' + value + ' bits'))


if __name__ == '__main__':
    unittest.main()
