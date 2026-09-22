import unittest

from viewpoints import Views, core


def model(extra=''):
    return core.canonicalize(core.parse('''language design-core version 0.5.
actor User.
usecase Goal.
usecase OtherGoal.
feature First.
feature Second.
functionality Work.
functionality Unplaced.
unit Owner.
container Host.
container OtherHost.
mode Cli.
mode Gui.
User pursues Goal.
First supports Goal.
First supports OtherGoal.
Second supports Goal.
Work contributes-to First.
Work contributes-to Second.
Work contributes-to Goal.
Unplaced contributes-to First.
Owner owns Work.
Owner owns Unplaced.
Host contains Owner.
Work allocated-to Host in mode Cli.
Work allocated-to OtherHost in mode Gui.
''' + extra))


class GoalViewTests(unittest.TestCase):
    def test_many_to_many_and_direct_contribution_have_exact_source_facts(self):
        views = Views(model())
        edges = {(a, b, c) for d in views.diagrams if d.ident.startswith('VP01-')
                 for a, b, c, _ in d.edges}
        self.assertTrue({('Work', 'contributes-to', 'Goal'),
                         ('First', 'supports', 'Goal'), ('Second', 'supports', 'Goal'),
                         ('Work', 'contributes-to', 'First'), ('Work', 'contributes-to', 'Second')} <= edges)
        facts = {r['id']: r for r in views.facts}
        for d in views.diagrams:
            for a, verb, b, fact in d.edges:
                self.assertEqual((a, verb, b), (facts[fact]['subject'], facts[fact]['verb'], facts[fact]['object']))

    def test_modes_do_not_leak_into_each_other_and_gaps_are_per_mode(self):
        views = Views(model(), ['VP07'])
        diagrams = {d.ident: d for d in views.diagrams}
        for mode, host, absent in [('Cli', 'Host', 'OtherHost'), ('Gui', 'OtherHost', 'Host')]:
            d = diagrams['VP07-First-' + mode]
            self.assertIn(host, d.nodes)
            self.assertNotIn(absent, d.nodes)
            self.assertIn('Owner', d.nodes)
            self.assertIn('Unplaced', d.nodes)
            self.assertTrue(any(g['model_id'] == 'Unplaced' and g['mode'] == mode for g in views.gaps))
        self.assertTrue(all(g['viewpoint'] == 'VP07' for g in views.gaps))

    def test_contains_never_implies_allocation(self):
        text = '\n'.join(line for line in model().splitlines() if ' allocated-to ' not in line) + '\n'
        views = Views(text, ['VP07'])
        self.assertTrue(all(d.ident.endswith('-unallocated') for d in views.diagrams))
        self.assertTrue(all('Host' not in d.nodes and 'OtherHost' not in d.nodes for d in views.diagrams))
        self.assertTrue(any(g['model_id'] == 'Work' and g['mode'] is None for g in views.gaps))

    def test_declared_but_unconnected_goals_and_features_are_visible(self):
        text = 'language design-core version 0.5.\nactor Alone.\nfeature Empty.\nusecase Missing.\n'
        views = Views(text, ['VP01', 'VP07'])
        self.assertEqual({d.ident for d in views.diagrams},
                         {'VP01-Missing', 'VP01-feature-Empty', 'VP07-Empty-unallocated'})
        self.assertTrue({'NO_GOAL', 'NO_ACTOR', 'NO_CONTRIBUTION'} <= {g['code'] for g in views.gaps})
        self.assertIn('Alone', views.markdown())

    def test_changing_allocation_changes_projection_and_inventory_keeps_mode(self):
        before = Views(model(), ['VP07', 'VP11'])
        text = model().replace('Work allocated-to Host in mode Cli.', 'Work allocated-to OtherHost in mode Cli.')
        after = Views(core.canonicalize(core.parse(text)), ['VP07', 'VP11'])
        self.assertNotEqual(before.markdown(), after.markdown())
        self.assertIn('Work allocated-to OtherHost in mode Cli', after.markdown())
        d = next(d for d in after.diagrams if d.ident == 'VP07-First-Cli')
        self.assertNotIn('Host', d.nodes)

    def test_old_profile_and_conflicting_allocation_cannot_be_projected(self):
        with self.assertRaises(ValueError):
            Views(model().replace('version 0.5', 'version 0.1'))
        with self.assertRaises(ValueError):
            Views(model() + 'Work allocated-to OtherHost in mode Cli.\n')


if __name__ == '__main__':
    unittest.main()
