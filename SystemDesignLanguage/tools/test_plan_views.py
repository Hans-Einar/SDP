import unittest
from viewpoints import Views, core
from plan_views import implementation_markdown


class PlanViewTests(unittest.TestCase):
    def source(self, extra=''):
        return core.canonicalize(core.parse('''language design-core version 0.5.
activity Phase.
activity First.
activity Second.
functionality ParseSource.
functionality Missing.
feature Authoring.
unit Parser.
Phase delivers Authoring.
Phase has implementation-status = planned.
First refines Phase.
First addresses ParseSource.
First has implementation-status = planned.
Second refines Phase.
Second depends-on First.
Second has implementation-status = planned.
Parser owns ParseSource.
Parser owns Missing.
''' + extra))

    def test_plan_is_projected_from_arbitrary_names_and_shows_missing_coverage(self):
        views = Views(self.source(), ['VP06'])
        output = implementation_markdown(views)
        self.assertIn('## Phase', output)
        self.assertIn('### First', output)
        self.assertIn('ParseSource | Parser', output)
        self.assertIn('**planned**', output)
        self.assertIn('Missing', output)
        self.assertNotIn('G1', output)
        self.assertNotIn('architecture.design', output)
        facts = {r['id']: r for r in views.facts}
        for d in views.diagrams:
            for a, verb, b, fact in d.edges:
                self.assertEqual((a, verb, b), (facts[fact]['subject'], facts[fact]['verb'], facts[fact]['object']))

    def test_modified_plan_changes_report_and_can_close_coverage(self):
        output = implementation_markdown(Views(self.source('Second addresses Missing.\n')))
        self.assertIn('0 mangler', output)
        self.assertIn('| Missing | Parser |', output)


if __name__ == '__main__':
    unittest.main()
