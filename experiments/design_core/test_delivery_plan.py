import unittest
import design_core as core


class PlanTests(unittest.TestCase):
    def test_plan_is_typed_and_does_not_claim_runtime_execution(self):
        source = '''language design-core version 0.5.
activity Phase.
activity Milestone.
feature Outcome.
functionality ParseSource.
unit Parser.
Phase delivers Outcome.
Milestone refines Phase.
Milestone addresses ParseSource.
Milestone has implementation-status = planned.
Parser owns ParseSource.
'''
        self.assertEqual(core.validate(core.parse(source)), [])
        canonical = core.canonicalize(core.parse(source))
        self.assertEqual(core.check(canonical)[1], [])
        invalid = source.replace('Milestone addresses ParseSource.', 'Phase addresses Parser.')
        self.assertIn('OBJECT_TYPE_MISMATCH', [d.code for d in core.validate(core.parse(invalid))])

    def test_dependency_cycles_and_wrong_status_are_rejected(self):
        source = 'language design-core version 0.5.\nactivity A.\nactivity B.\nA depends-on B.\nB depends-on A.\n'
        self.assertIn('STRUCTURE_CYCLE', [d.code for d in core.validate(core.parse(source))])
        source = 'language design-core version 0.5.\nunit Parser.\nParser has implementation-status = planned.\n'
        self.assertIn('PROPERTY_TYPE_MISMATCH', [d.code for d in core.validate(core.parse(source))])


if __name__ == '__main__':
    unittest.main()
