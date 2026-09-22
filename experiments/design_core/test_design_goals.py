"""V1 traceability, allocation cardinality and source-preserving AST checks."""
import unittest

import design_core as dc

SOURCE = '''language design-core version 0.5.
actor Author.
usecase Review.
usecase Prototype.
feature Checking.
feature Preview.
functionality ParseSource.
unit Frontend.
container Cli.
container Gui.
mode Inspect.
mode Live.
Author pursues Review.
Author pursues Prototype.
Checking supports Review.
Checking supports Prototype.
Preview supports Prototype.
ParseSource contributes-to Checking.
ParseSource contributes-to Preview.
ParseSource contributes-to Review.
Frontend owns ParseSource.
ParseSource allocated-to Cli in mode Inspect.
ParseSource allocated-to Gui in mode Live.
'''


class GoalTests(unittest.TestCase):
    def test_many_to_many_and_direct_contribution_roundtrip(self):
        model = dc.parse(SOURCE)
        self.assertEqual(dc.validate(model), [])
        canonical = dc.canonicalize(model)
        self.assertEqual(dc.check(canonical)[1], [])
        self.assertEqual(dc.canonicalize(dc.parse(canonical)), canonical)
        self.assertEqual({dc.sentence(s) for s in model.statements},
                         {dc.sentence(s) for s in dc.parse(canonical).statements})

    def test_allocation_preserves_argument_spans_and_json_type(self):
        model = dc.parse(SOURCE)
        allocation = model.statements[-1]
        self.assertIsInstance(allocation, dc.Allocation)
        for arg, expected in [(allocation.subject, 'ParseSource'),
                              (allocation.container, 'Gui'), (allocation.mode, 'Live')]:
            self.assertEqual(SOURCE[arg.span.start:arg.span.end], expected)
        self.assertEqual(dc.to_json(allocation)['node'], 'Allocation')
        self.assertEqual(SOURCE[allocation.span.start:allocation.span.end],
                         'ParseSource allocated-to Gui in mode Live.')

    def test_same_mode_rejects_multiple_containers(self):
        errors = dc.validate(dc.parse(SOURCE + 'ParseSource allocated-to Gui in mode Inspect.\n'))
        self.assertEqual([d.code for d in errors], ['ALLOCATION_CARDINALITY'])
        self.assertEqual(errors[0].span.line, len(SOURCE.splitlines()) + 1)

    def test_duplicate_allocation_is_not_a_second_instance(self):
        errors = dc.validate(dc.parse(SOURCE + 'ParseSource allocated-to Cli in mode Inspect.\n'))
        self.assertEqual([d.code for d in errors], ['DUPLICATE_FACT'])

    def test_each_allocation_argument_is_typed(self):
        for fact, code in [
            ('Frontend allocated-to Cli in mode Inspect.', 'SUBJECT_TYPE_MISMATCH'),
            ('ParseSource allocated-to Frontend in mode Inspect.', 'OBJECT_TYPE_MISMATCH'),
            ('ParseSource allocated-to Cli in mode Author.', 'QUALIFIER_TYPE_MISMATCH'),
        ]:
            with self.subTest(fact=fact):
                self.assertEqual([d.code for d in dc.validate(dc.parse(SOURCE + fact))], [code])

    def test_missing_allocation_names_and_mode(self):
        errors = dc.validate(dc.parse('language design-core version 0.5.\nMissing allocated-to Unknown in mode Absent.\n'))
        self.assertEqual([d.code for d in errors], ['UNDECLARED_NAME'] * 3)
        with self.assertRaises(dc.ParseError):
            dc.parse(SOURCE + 'ParseSource allocated-to Cli.')

    def test_capability_is_not_a_feature(self):
        text = SOURCE.replace('actor Author.', 'actor Author.\ncapability Ability.')
        for fact in ['ParseSource contributes-to Ability.', 'Ability supports Review.']:
            self.assertTrue(any(d.code.endswith('_TYPE_MISMATCH') for d in dc.validate(dc.parse(text + fact))))

    def test_incomplete_goal_model_is_structurally_valid(self):
        self.assertEqual(dc.check('language design-core version 0.5.\nactor Author.\nfeature Missing.\nusecase Review.\n')[1], [])

    def test_allocation_never_replaces_owner_or_changes_containment(self):
        errors = dc.validate(dc.parse(SOURCE.replace('Frontend owns ParseSource.\n', '')))
        self.assertEqual([d.code for d in errors], ['OWNERSHIP_CARDINALITY'])
        self.assertFalse(any(isinstance(s, dc.Relation) and s.verb == 'contains' for s in dc.parse(SOURCE).statements))


if __name__ == '__main__':
    unittest.main()
