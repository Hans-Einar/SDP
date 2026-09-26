"""Historical identity compatibility without weakening relation integrity."""
import copy
import json
import unittest
from pathlib import Path
import yaml
from jsonschema import Draft202012Validator
from test_validate_sdp import VALIDATE
ROOT = Path(__file__).resolve().parents[2]

class TraceabilityCompatibilityTests(unittest.TestCase):
    def test_work_identity_categories_agree(self):
        relations=json.loads((ROOT/'Toolkit/schemas/relations.schema.json').read_text())
        current=json.loads((ROOT/'Toolkit/schemas/current-index.schema.json').read_text())
        cases=[('sprints','sprintId',['Sprint-001','SPR-SDP-005'],['SPR-SDP-5','SLC-SDP-005-001']),
               ('iterations','iterationId',['SPI-001','ITR-SDP-005-001'],['ITR-SDP-005','ITR-sdp-005-001']),
               ('slices','sliceId',['SPS-001','SLC-SDP-005-001'],['SLC-SDP-005','SDPTOOL-VER-T1-M1'])]
        for category,definition,valid,invalid in cases:
            for value in valid+invalid:
                with self.subTest(category=category,value=value):
                    self.assertEqual(Draft202012Validator(relations['$defs'][category]['propertyNames']).is_valid(value),value in valid)
                    self.assertEqual(Draft202012Validator(current['$defs'][definition]).is_valid(value),value in valid)

    def test_real_records_keep_typed_evidence_and_reciprocity(self):
        records=yaml.safe_load((ROOT/'SDP/Traceability/Relations.yaml').read_text())
        schema=json.loads((ROOT/'Toolkit/schemas/relations.schema.json').read_text())
        Draft202012Validator(schema).validate(records)
        self.assertNotIn('SDPTOOL-VER-P0-M1',records['slices'])
        self.assertIn('SDPTOOL-VER-P0-M1',records['verification'])
        self.assertIn('VER-SPS-001',records['slices']['SPS-001']['verification'])
        self.assertIn('REV-SPS-001-007',records['slices']['SPS-001']['reviews'])
        self.assertEqual(VALIDATE.validate_relations_semantics(records,ROOT/'SDP','test',ROOT),[])
        broken=copy.deepcopy(records)
        broken['slices']['SPS-001']['reviews'].remove('REV-SPS-001-007')
        self.assertTrue(any('missing reverse link' in x for x in VALIDATE.validate_relations_semantics(broken,ROOT/'SDP','test',ROOT)))
        broken=copy.deepcopy(records);broken['verification']['SDPTOOL-VER-P0-M1']['slice']='SPS-999'
        self.assertTrue(any('dangling ID' in x for x in VALIDATE.validate_relations_semantics(broken,ROOT/'SDP','test',ROOT)))

if __name__=='__main__':unittest.main()
