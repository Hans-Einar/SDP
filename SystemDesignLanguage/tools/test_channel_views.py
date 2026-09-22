import unittest
from viewpoints import Views, core
from test_channel_core import SOURCE


class ChannelViewTests(unittest.TestCase):
    def test_sequence_is_source_ordered_and_each_arrow_has_contract_proof(self):
        views = Views(core.canonicalize(core.parse(SOURCE)), ['VP08'])
        diagram = next(d for d in views.diagrams if d.kind == 'sequence')
        self.assertEqual([e['ordinal'] for e in diagram.elements], [1, 2])
        syntax = diagram.mermaid(views.kinds)
        self.assertIn('n_Client->>n_Server: 1: Request', syntax)
        self.assertIn('n_Server-->>n_Client: 2: Result', syntax)
        facts = {f['id']: f for f in views.facts}
        for e in diagram.elements:
            self.assertEqual(facts[e['fact']]['ordinal'], e['ordinal'])
            self.assertTrue({'uses', 'permits', 'upholds', 'step'} <= {facts[f]['verb'] for f in e['proof']})
        self.assertEqual(len(views.message_sets), 2)
        result = next(r for r in views.message_sets if r['message'] == 'Result')
        self.assertEqual(result['senders'], ['Server'])
        self.assertEqual(result['receivers'], ['Client'])

    def test_numeric_order_survives_lexical_canonicalization(self):
        source = '\n'.join(line for line in SOURCE.splitlines() if ' step ' not in line)+'\n'
        for i in range(1, 14, 2):
            source += f'Accepted step {i} sends Request from Client to Server via Service.\n'
            source += f'Accepted step {i+1} sends Result from Server to Client via Service reply-to {i}.\n'
        views = Views(core.canonicalize(core.parse(source)), ['VP08'])
        d = views.diagrams[0]
        self.assertEqual([e['ordinal'] for e in d.elements], list(range(1, 15)))
        self.assertLess(d.mermaid(views.kinds).index(': 2:'), d.mermaid(views.kinds).index(': 10:'))

    def test_channel_without_scenario_does_not_invent_sequence(self):
        source = '\n'.join(line for line in SOURCE.splitlines() if not line.startswith(('scenario Accepted.', 'Accepted ')))+'\n'
        views = Views(core.canonicalize(core.parse(source)), ['VP08'])
        self.assertFalse(views.diagrams)
        self.assertEqual(len(views.message_sets), 2)
        self.assertIn('NO_SCENARIOS', {g['code'] for g in views.gaps})

    def test_incomplete_participation_is_reported_without_guessing_receiver(self):
        source = '\n'.join(line for line in SOURCE.splitlines() if not line.startswith(('scenario Accepted.', 'Accepted ', 'Client uses Service as receiver'))) + '\n'
        views = Views(core.canonicalize(core.parse(source)), ['VP08'])
        result = next(r for r in views.message_sets if r['message'] == 'Result')
        self.assertEqual(result['receivers'], [])
        self.assertIn('INCOMPLETE_PARTICIPATION', {g['code'] for g in views.gaps})

    def test_selection_excludes_sequence_and_messageset(self):
        views = Views(core.canonicalize(core.parse(SOURCE)), ['VP09'])
        self.assertEqual(views.message_sets, [])
        self.assertNotIn('sequenceDiagram', views.markdown())


if __name__ == '__main__':
    unittest.main()
