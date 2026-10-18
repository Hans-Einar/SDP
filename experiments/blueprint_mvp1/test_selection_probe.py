import copy
import os
from pathlib import Path
import shutil
import tempfile
import unittest

from selection_probe import load, select

class SelectionTests(unittest.TestCase):
    @classmethod
    def setUpClass(cls):
        cls.tool = os.environ['BP2_SDL']
        cls.root = Path(__file__).resolve().parent
        cls.now = load(cls.tool, cls.root / 'NOW')
        cls.target = load(cls.tool, cls.root / 'TARGET')

    def test_union_retains_removed_owner_and_atomic_channel_contract(self):
        r = select(self.now, self.target)
        self.assertTrue(any('CalibrateMeasurement' in k and 'MachineService' in k for k in r['removed']))
        self.assertTrue({'MachineService', 'BuckingWeb', 'BaselineRequest', 'Subscriber',
                         'BaselinePayload', 'CalibrationState', 'MachineBaseline', 'ReadBaseline'} <= set(r['nodes']))
        self.assertEqual(r['added'], sorted(r['added']))

    def test_authored_context_and_unknown_frontier(self):
        r = select(self.now, self.target)
        self.assertNotIn('BuckingUI', r['nodes'])
        self.assertIn('BuckingUI', r['excluded'])
        self.assertTrue(r['frontier'])
        r = select(self.now, self.target, ['BuckingUI'])
        self.assertIn('authored context', r['reasons']['BuckingUI'])
        self.assertTrue(r['unknowns'])

    def test_cycle_terminates_and_is_deterministic(self):
        # Synthetic normalized graph case, not evidence for new SDL syntax.
        m = {'nodes': {'a': 'unit', 'b': 'unit'}, 'facts': {
            'a-b': {'refs': ['a', 'b']}, 'b-a': {'refs': ['b', 'a']}}, 'revision': 'fixture'}
        first = select(m, m, ['a'])
        self.assertEqual(first, select(m, m, ['a']))
        self.assertEqual(first['nodes'], ['a', 'b'])

    def test_overflow_refuses_partial_success(self):
        with self.assertRaisesRegex(ValueError, 'overflow'):
            select(self.now, self.target, limit=2)

    def test_stale_input(self):
        with self.assertRaisesRegex(ValueError, 'stale'):
            select(self.now, self.target, expected=('old', self.target['revision']))

    def test_missing_endpoint_is_exposed_despite_parser_acceptance(self):
        with tempfile.TemporaryDirectory() as tmp:
            path = Path(tmp) / 'model'
            shutil.copytree(self.root / 'NOW', path)
            p = path / 'Contracts/Machine.design'
            p.write_text(p.read_text().replace('MachineService uses MachineServiceUpstream as receiver of ReadBaseline in mode Inspection.\n', ''))
            missing = load(self.tool, path)
            r = select(self.now, missing)
            self.assertTrue(any('TARGET missing channel peer' in g for g in r['unknowns']))

    def test_valid_sdl_can_violate_preserved_ownership(self):
        with tempfile.TemporaryDirectory() as tmp:
            path = Path(tmp) / 'model'
            shutil.copytree(self.root / 'TARGET', path)
            p = path / 'Containers/MachineService.design'
            p.write_text(p.read_text().replace('MachineService owns ReduceMachineState.', 'BuckingWeb owns ReduceMachineState.'))
            import json, subprocess
            formatted = json.loads(subprocess.check_output([self.tool, 'format', str(path / 'System.design'), '--file-map'], text=True))
            for name, text in formatted['files'].items():
                (path / name).write_text(text)
            drift = load(self.tool, path)
            # Authored PRESERVE obligation, not a parser invariant.
            def preserved(m):
                return any(f['fact'].get('verb') == 'owns' and
                           f['fact'].get('subject', {}).get('name') == 'MachineService' and
                           f['fact'].get('object', {}).get('name') == 'ReduceMachineState'
                           for f in m['facts'].values())
            self.assertTrue(preserved(self.target))
            self.assertFalse(preserved(drift))
            self.assertIn('BuckingWeb', select(self.now, drift)['nodes'])

    def test_source_locations_do_not_create_semantic_differences(self):
        other = copy.deepcopy(self.now)
        for f in other['facts'].values():
            f['source']['line'] += 7
        r = select(self.now, other)
        self.assertEqual(r['added'], [])
        self.assertEqual(r['removed'], [])

if __name__ == '__main__':
    unittest.main()
