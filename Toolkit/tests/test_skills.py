"""Native/distribution metadata and installed-reference regression checks."""
import importlib.util,json,sys,tempfile,unittest,re
from pathlib import Path
import yaml
ROOT=Path(__file__).resolve().parents[2]
spec=importlib.util.spec_from_file_location('skills_validate',ROOT/'Toolkit/scripts/validate_sdp.py')
v=importlib.util.module_from_spec(spec);sys.modules[spec.name]=v;spec.loader.exec_module(v)

class Skills(unittest.TestCase):
    def test_canonical_inventory_and_portable_references(self):
        manifest=yaml.safe_load((ROOT/'SDP.manifest.yaml').read_text())
        contract=json.loads((ROOT/'Toolkit/SDP-install.manifest.json').read_text())
        roles=list((ROOT/'Skills').glob('*/SKILL.md'))
        self.assertEqual(len(roles),14)
        self.assertEqual({p.parent.name for p in roles},set(manifest['skills']))
        with tempfile.TemporaryDirectory() as t:
            installed=Path(t)
            for e in contract['entries']:
                if e.get('source','').startswith('Skills/'):
                    p=installed/e['destination'];p.parent.mkdir(parents=True,exist_ok=True)
                    p.write_bytes((ROOT/e['source']).read_bytes())
            for p in roles:
                self.assertIn('name',v.parse_front_matter(p))
                self.assertEqual(v.validate_skill_metadata(p,p.parent.name,manifest['skills'][p.parent.name],manifest['toolkit']['version']),[])
                self.assertEqual((ROOT/'.agents/skills'/p.parent.name).resolve(),p.parent)
            for p in installed.rglob('*.md'):
                for link in re.findall(r'\]\(([^)\s]+)\)',p.read_text()):
                    if '://' not in link:self.assertTrue((p.parent/link).is_file(),(p,link))

    def test_native_metadata_errors_and_legacy_compatibility(self):
        base=yaml.safe_load((ROOT/'Skills/sdp/SKILL.md').read_text().split('---',2)[1])
        with tempfile.TemporaryDirectory() as t:
            p=Path(t)/'SKILL.md'
            for mutate,expected in [
                (lambda d:d.pop('description'),'description'),
                (lambda d:d.update(name='wrong'),'name must equal'),
                (lambda d:d['metadata'].update(skillVersion=2),'must be strings'),
                (lambda d:d.update(skillVersion='1.0.0'),'mixed v1/v2'),
                (lambda d:d['metadata'].update(capabilities=''),'non-empty list'),
            ]:
                d=yaml.safe_load(yaml.safe_dump(base));mutate(d)
                p.write_text('---\n'+yaml.safe_dump(d)+'---\nbody\n')
                self.assertTrue(any(expected in e for e in v.validate_skill_metadata(p,'sdp','1.0.0','0.2.0')),expected)
            legacy={'skillId':'sdp','skillVersion':'1.0.0','minimumToolkitVersion':'0.2.0','capabilities':['sdp.route'],'compatibilityNotes':'Historical fixture.'}
            p.write_text('---\n'+yaml.safe_dump(legacy)+'---\n')
            self.assertEqual(v.validate_skill_metadata(p,'sdp','1.0.0','0.2.0'),[])
            legacy['skillVersion']='2.0.0';p.write_text('---\n'+yaml.safe_dump(legacy)+'---\n')
            self.assertTrue(any('native metadata required' in e for e in v.validate_skill_metadata(p,'sdp','2.0.0','0.2.0',require_native=True)))

if __name__=='__main__':unittest.main()
