import copy
import importlib.util
import json
from pathlib import Path
import unittest

ROOT = Path(__file__).resolve().parents[2]
s = importlib.util.spec_from_file_location('profile', ROOT/'Toolkit/scripts/build_process_profile.py')
m = importlib.util.module_from_spec(s)
s.loader.exec_module(m)

class ProfileBuild(unittest.TestCase):
    def setUp(self):
        self.config = json.loads((ROOT/'Toolkit/profiles/five-phase.json').read_text())

    def test_reproducible_and_neutral(self):
        one = m.canonical(m.build(self.config))
        self.assertEqual(one, m.canonical(m.build(copy.deepcopy(self.config))))
        self.assertEqual(one, (ROOT/'Toolkit/profiles/five-phase.artifact.json').read_bytes())
        self.assertNotIn(b'MAINT-SDP-0003', one)
        self.assertNotIn(b'KB-SDP-028', one)

    def test_invalid_inputs(self):
        mutations = [
            lambda c: c.update(operation='shell'),
            lambda c: c['files'].append(c['files'][0]),
            lambda c: c['files'][0].update(source='../../secret'),
            lambda c: c['files'][0].update(source='Skills/missing.md'),
            lambda c: c['files'][0].update(destination='SDP/../escape'),
            lambda c: c['files'][0].update(destination='SDP/CON'),
            lambda c: c['files'][0].update(ownership='overwrite'),
            lambda c: c['relocations'].append(c['relocations'][0]),
            lambda c: c['relocations'][0].update(to='SDP/02--Study/nested'),
            lambda c: c.update(prerequisites=['execute anything']),
            lambda c: c['files'][0].update(destination='SDP/.sdp-operations/install.lock'),
            lambda c: c['files'][0].update(destination='SDP/Framework/installed-toolkit.manifest.yaml'),
            lambda c: c['files'][0].update(destination='arbitrary.py'),
            lambda c: c['files'][0].update(destination='SDP/Framework/installed-toolkit.manifest.yaml/child.md'),
            lambda c: c['relocations'][0].update(to='SDP/.sdp-operations/imported'),
            lambda c: c['relocations'][0].update(to='SDP/ProjectManagement'),
            lambda c: c['relocations'][0].update(to='SDP/Framework/installed-toolkit.manifest.yaml/child'),
        ]
        for change in mutations:
            c = copy.deepcopy(self.config)
            change(c)
            with self.subTest(c=c), self.assertRaises(ValueError):
                m.build(c)

if __name__ == '__main__':
    unittest.main()
