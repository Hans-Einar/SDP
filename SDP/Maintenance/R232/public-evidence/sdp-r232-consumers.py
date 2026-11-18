#!/usr/bin/env python3
"""Disposable public-release verification. Never modifies the live gh extension."""
from pathlib import Path
import argparse, hashlib, json, os, shutil, subprocess

ROOT = Path('/home/warloc/git/gh-sdp')
LIVE = Path('/home/warloc/.local/share/gh/extensions/gh-sdp')
GH = shutil.which('gh')
UNSHARE = shutil.which('unshare')
OLD_VERSION = '2.3.1'
NEW_VERSION = '2.3.2'
def sha(b): return hashlib.sha256(b).hexdigest()
def snapshot(path):
    return {str(p.relative_to(path)): sha(p.read_bytes()) for p in sorted(path.rglob('*')) if p.is_file()}

def main():
    p=argparse.ArgumentParser();p.add_argument('--output',required=True);p.add_argument('--package');p.add_argument('--canonical-descriptor');p.add_argument('--baseline-only',action='store_true');a=p.parse_args()
    output=Path(a.output).resolve();assert output==Path('/tmp/sdp-r232-public-consumers'), 'use isolated selected /tmp output';output.mkdir(parents=True,exist_ok=True)
    assert GH and UNSHARE, 'gh and unshare required'
    before=snapshot(LIVE)
    receipt=json.loads((ROOT/'SDP/Sprints/Sprint-012/evidence/package.json').read_text())
    expected_old=receipt['sha256']
    assert before['gh-sdp']==expected_old
    manifest=(LIVE/'manifest.yml').read_text()
    assert 'tag: v2.3.1' in manifest and 'ispinned: false' in manifest
    # Capture credentials into process memory only. Never log/store the token.
    token=os.environ.get('GH_TOKEN') or os.environ.get('GITHUB_TOKEN')
    if not token:
        credential=subprocess.run([GH,'auth','token','--hostname','github.com'],capture_output=True,text=True)
        assert credential.returncode==0, 'current GitHub authentication unavailable'
        token=credential.stdout.strip()
    assert token, 'empty GitHub credential'
    network={k:os.environ[k] for k in ['HTTP_PROXY','HTTPS_PROXY','ALL_PROXY','NO_PROXY','SSL_CERT_FILE','SSL_CERT_DIR'] if k in os.environ}
    logs=[]
    def env(home):
        return dict(network,HOME=str(home),PATH=os.environ.get('PATH','/usr/bin:/bin'),GH_TOKEN=token,GH_CONFIG_DIR=str(home/'gh-config'),XDG_DATA_HOME=str(home/'data'),XDG_CACHE_HOME=str(home/'cache'),GH_PROMPT_DISABLED='1',GH_NO_UPDATE_NOTIFIER='1',GH_NO_EXTENSION_UPDATE_NOTIFIER='1')
    def run(args,environment,cwd,name):
        command=[str(x) for x in args]
        result=subprocess.run(command,env=environment,cwd=cwd,capture_output=True,text=True,timeout=180)
        # The command contains no secrets; GH_DEBUG is never enabled.
        logs.append(dict(name=name,command=command,exitCode=result.returncode,stdout=result.stdout,stderr=result.stderr))
        (output/'commands.json').write_text(json.dumps(logs,indent=2)+'\n')
        assert result.returncode==0, name+' failed: '+result.stderr[:2000]
        return result.stdout
    try:
        control=output/'control';control.mkdir(exist_ok=True);control_env=env(control)
        release=json.loads(run([GH,'api','repos/Hans-Einar/gh-sdp/releases/tags/v2.3.1'],control_env,control,'old-public-release'))
        old_asset=next(x for x in release['assets'] if x['name']=='gh-sdp-linux-amd64')
        assert old_asset['digest']=='sha256:'+expected_old
        old_download=control/'old-download';old_download.mkdir(exist_ok=True)
        run([GH,'release','download','v2.3.1','--repo','Hans-Einar/gh-sdp','--pattern','gh-sdp-linux-amd64','--dir',old_download,'--clobber'],control_env,control,'old-public-asset')
        assert sha((old_download/'gh-sdp-linux-amd64').read_bytes())==expected_old
        (output/'baseline.json').write_text(json.dumps(dict(liveSnapshot=before,publishedTag='v2.3.1',publishedAssetSHA256=expected_old,liveUnchanged=snapshot(LIVE)==before),indent=2)+'\n')
        if a.baseline_only:return
        assert a.package and a.canonical_descriptor, '--package and --canonical-descriptor required for published 2.3.2 tests'
        package=Path(a.package).resolve(); package_manifest=json.loads((package/'gh-sdp.manifest.json').read_text())
        expected_engine=package_manifest['sha256'];expected_source=package_manifest['engineCommit']
        assert package_manifest['version']==package_manifest['engineVersion']==NEW_VERSION
        assert package_manifest['provenance']=='signed' and len(package_manifest['distributionCommit'])==40
        descriptor=Path(a.canonical_descriptor).read_bytes();canonical=json.loads(descriptor)
        assert sha(descriptor)==package_manifest['descriptorSHA256']
        assert descriptor==(package/'sdp-release.json').read_bytes()
        assert canonical['release']==NEW_VERSION and canonical['sourceCommit']==expected_source
        asset=next(x for x in canonical['binaries'] if x['platform']=='linux/amd64')
        assert asset['sha256']==expected_engine and asset['size']==package_manifest['size']
        public=json.loads(run([GH,'api','repos/Hans-Einar/gh-sdp/releases/tags/v2.3.2'],control_env,control,'new-public-release'))
        assert not public['draft'] and not public['prerelease'] and public['published_at']
        asset=next(x for x in public['assets'] if x['name']=='gh-sdp-linux-amd64')
        assert asset['digest']=='sha256:'+expected_engine
        # Guard against any newer default tag before normal unpinned operations.
        latest=json.loads(run([GH,'api','repos/Hans-Einar/gh-sdp/releases/latest'],control_env,control,'public-latest'))
        assert latest['tag_name']=='v2.3.2'
        outcomes=[]
        for kind in ['upgrade','clean-install']:
            home=output/kind;home.mkdir();environment=env(home)
            project=home/'project';(project/'SDP').mkdir(parents=True);(project/'AGENTS.md').write_text('# Project-owned instructions\nPreserve exact bytes.\n')
            original=snapshot(project)
            extension=home/'data/gh/extensions/gh-sdp'
            if kind=='upgrade':
                extension.parent.mkdir(parents=True)
                shutil.copytree(LIVE,extension)
                assert snapshot(extension)==before, 'baseline clone must preserve actual bytes and manifest'
                cloned_manifest=extension/'manifest.yml'
                cloned_manifest.write_text(cloned_manifest.read_text().replace(str(LIVE/'gh-sdp'),str(extension/'gh-sdp')))
                old_version=json.loads(run([extension/'gh-sdp','--version','--json'],environment,project,kind+'-old-version'))
                assert old_version['version']==OLD_VERSION and old_version['revision']==receipt['engineCommit']
                run([extension/'gh-sdp','setup','--json'],environment,project,kind+'-previous-setup')
                registered=home/'.local/bin/sdptool';previous_target=registered.resolve(strict=True)
                assert registered.is_symlink() and previous_target==extension/'gh-sdp'
                assert sha(registered.read_bytes())==expected_old
                run([GH,'extension','upgrade','sdp'],environment,project,kind+'-normal-public-upgrade')
                # Check the pre-existing registration immediately after gh upgrade,
                # before another setup or gh sdp invocation.
                assert registered.resolve(strict=True)==previous_target
                assert sha(registered.read_bytes())==expected_engine
                replaced=json.loads(run([registered,'--version','--json'],environment,project,kind+'-registered-version-before-setup'))
                assert replaced['version']==NEW_VERSION and replaced['revision']==expected_source
            else:
                run([GH,'extension','install' ,'Hans-Einar/gh-sdp'],environment,project,kind+'-public-install')
            assert sha((extension/'gh-sdp').read_bytes())==expected_engine
            # No gh sdp invocation has occurred since the update/install.
            run([GH,'sdp','setup','--json'],environment,project,kind+'-setup')
            alias=home/'.local/bin/sdptool';target=alias.resolve(strict=True)
            assert alias.is_symlink() and target==extension/'gh-sdp'
            gh_version=json.loads(run([GH,'sdp','--version','--json'],environment,project,kind+'-gh-version'))
            direct_version=json.loads(run([alias,'--version','--json'],environment,project,kind+'-direct-version'))
            assert gh_version==direct_version and direct_version['version']==NEW_VERSION and direct_version['revision']==expected_source
            assert 'runtime.setup/1' in direct_version['capabilities'] and 'sdui.source/1' in direct_version['capabilities']
            offline=dict(HOME=str(home),PATH='/nonexistent',XDG_CACHE_HOME=str(home/'offline-empty-cache'),SDP_OFFLINE='true',LANG='C.UTF-8')
            def disconnected(args,name):
                return run([UNSHARE,'--user','--map-root-user','--net',alias]+args,offline,project,kind+'-offline-'+name)
            assert json.loads(disconnected(['--version','--json'],'version'))==direct_version
            discovery=json.loads(disconnected(['discover','--json'],'discover'));assert discovery['operation']=='discover'
            source=home/'page.sdui';source.write_text('sdui 0.2; page=[button("OK")];')
            checked=json.loads(disconnected(['sdui','check','--json','-source',str(source),'-entry','page'],'sdui-check'));assert checked['operation']=='sdui-check'
            bundle=home/'preview';bundle.mkdir()
            preview=json.loads(disconnected(['sdui','preview','--json','-source',str(source),'-entry','page','-output',str(bundle)],'sdui-preview'));assert preview['operation']=='sdui-preview' and 'OK' in Path(preview['entry']).read_text()
            assert snapshot(project)==original, 'public tool install/update mutated project'
            assert not (home/'offline-empty-cache/sdptool/distribution').exists()
            outcomes.append(dict(case=kind,version=direct_version,engineSHA256=sha(alias.read_bytes()),registeredTarget=str(target),networkNamespaceIsolated=True,authAbsent=True,ghAbsentFromPATH=True,projectUnchanged=True))
        (output/'results.json').write_text(json.dumps(dict(publicTag='v2.3.2',outcomes=outcomes,liveInstallationUnchanged=snapshot(LIVE)==before,packageManifest=package_manifest),indent=2)+'\n')
    finally:
        assert snapshot(LIVE)==before, 'LIVE EXTENSION CHANGED'

if __name__=='__main__':main()
