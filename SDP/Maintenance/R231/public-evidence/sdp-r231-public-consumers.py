#!/usr/bin/env python3
"""Disposable public-release verification. Never modifies the live gh extension."""
from pathlib import Path
import argparse, hashlib, json, os, shutil, subprocess

ROOT = Path('/home/warloc/git/gh-sdp')
LIVE = Path('/home/warloc/.local/share/gh/extensions/gh-sdp')
GH = shutil.which('gh')
UNSHARE = shutil.which('unshare')
EXPECTED_OLD = '5631dc2e849da5598606394319451be94f6023e0eb7e9b4edc6c82e1647e39c5'
EXPECTED_ENGINE = '04f78532a4e7ae8653df5b1e623e1407312713a90043f2ac36e7f68e0dc0495b'
EXPECTED_SOURCE = '84d5e46012fdc45a3690bd0b26b617501fd1caa7'
def sha(b): return hashlib.sha256(b).hexdigest()
def snapshot(path):
    return {str(p.relative_to(path)): sha(p.read_bytes()) for p in sorted(path.rglob('*')) if p.is_file()}

def main():
    p=argparse.ArgumentParser();p.add_argument('--output',required=True);p.add_argument('--package');p.add_argument('--baseline-only',action='store_true');a=p.parse_args()
    output=Path(a.output).resolve();output.mkdir(parents=True,exist_ok=False)
    assert GH and UNSHARE, 'gh and unshare required'
    before=snapshot(LIVE)
    receipt=json.loads((ROOT/'SDP/Sprints/Sprint-010/Publication.json').read_text())
    assert before['gh-sdp']==EXPECTED_OLD==receipt['assetSHA256']['gh-sdp-linux-amd64']
    manifest=(LIVE/'manifest.yml').read_text()
    assert 'tag: v0.2.3' in manifest and 'ispinned: false' in manifest
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
        control=output/'control';control.mkdir();control_env=env(control)
        release=json.loads(run([GH,'api','repos/Hans-Einar/gh-sdp/releases/tags/v0.2.3'],control_env,control,'old-public-release'))
        old_asset=next(x for x in release['assets'] if x['name']=='gh-sdp-linux-amd64')
        assert old_asset['digest']=='sha256:'+EXPECTED_OLD
        old_download=control/'old-download';old_download.mkdir()
        run([GH,'release','download','v0.2.3','--repo','Hans-Einar/gh-sdp','--pattern','gh-sdp-linux-amd64','--dir',old_download],control_env,control,'old-public-asset')
        assert sha((old_download/'gh-sdp-linux-amd64').read_bytes())==EXPECTED_OLD
        (output/'baseline.json').write_text(json.dumps(dict(liveSnapshot=before,publishedTag='v0.2.3',publishedAssetSHA256=EXPECTED_OLD,liveUnchanged=snapshot(LIVE)==before),indent=2)+'\n')
        if a.baseline_only:return
        assert a.package, '--package required for published 2.3.1 tests'
        package=Path(a.package).resolve(); package_manifest=json.loads((package/'gh-sdp.manifest.json').read_text())
        assert package_manifest['distributionCommit']=='ae7aa216c4343c6d4904078c8545cd1d2d8f7a64'
        assert package_manifest['provenance']=='signed' and package_manifest['descriptor']=='https://github.com/Hans-Einar/SDP/releases/download/v2.3.1/sdp-release.json'
        assert package_manifest['engineCommit']==EXPECTED_SOURCE and package_manifest['sha256']==EXPECTED_ENGINE
        public=json.loads(run([GH,'api','repos/Hans-Einar/gh-sdp/releases/tags/v2.3.1'],control_env,control,'new-public-release'))
        assert not public['draft'] and not public['prerelease'] and public['published_at']
        asset=next(x for x in public['assets'] if x['name']=='gh-sdp-linux-amd64')
        assert asset['digest']=='sha256:'+EXPECTED_ENGINE
        # Guard against any newer default tag before normal unpinned operations.
        latest=json.loads(run([GH,'api','repos/Hans-Einar/gh-sdp/releases/latest'],control_env,control,'public-latest'))
        assert latest['tag_name']=='v2.3.1'
        outcomes=[]
        for kind in ['upgrade','upgrade-2.3.0','clean-install']:
            home=output/kind;home.mkdir();environment=env(home)
            project=home/'project';(project/'SDP').mkdir(parents=True);(project/'AGENTS.md').write_text('# Project-owned instructions\nPreserve exact bytes.\n')
            original=snapshot(project)
            extension=home/'data/gh/extensions/gh-sdp'
            if kind=='upgrade':
                extension.parent.mkdir(parents=True)
                shutil.copytree(LIVE,extension)
                assert snapshot(extension)==before, 'baseline clone must preserve actual bytes and manifest'
                run([GH,'extension','upgrade','sdp'],environment,project,kind+'-normal-public-upgrade')
            elif kind=='upgrade-2.3.0':
                previous=Path('/tmp/activity4-public-consumer-results/clean-install/data/gh/extensions/gh-sdp')
                assert sha((previous/'gh-sdp').read_bytes())=='2cfcf12f8bfb0d6fa8d11391db874f68d9bc9210c85974ff3cecee50f915ced6'
                assert 'tag: v2.3.0' in (previous/'manifest.yml').read_text() and 'ispinned: false' in (previous/'manifest.yml').read_text()
                extension.parent.mkdir(parents=True)
                shutil.copytree(previous,extension)
                manifest=extension/'manifest.yml'
                manifest.write_text(manifest.read_text().replace(str(previous/'gh-sdp'),str(extension/'gh-sdp')))
                run([extension/'gh-sdp','setup','--json'],environment,project,kind+'-previous-setup')
                previous_target=(home/'.local/bin/sdptool').resolve(strict=True)
                assert previous_target==extension/'gh-sdp'
                run([GH,'extension','upgrade','sdp'],environment,project,kind+'-normal-public-upgrade')
                assert (home/'.local/bin/sdptool').resolve(strict=True)==previous_target
                assert sha((home/'.local/bin/sdptool').read_bytes())==EXPECTED_ENGINE
            else:
                run([GH,'extension','install' ,'Hans-Einar/gh-sdp'],environment,project,kind+'-public-install')
            assert sha((extension/'gh-sdp').read_bytes())==EXPECTED_ENGINE
            # No gh sdp invocation has occurred since the update/install.
            run([GH,'sdp','setup','--json'],environment,project,kind+'-setup')
            alias=home/'.local/bin/sdptool';target=alias.resolve(strict=True)
            assert alias.is_symlink() and target==extension/'gh-sdp'
            gh_version=json.loads(run([GH,'sdp','--version','--json'],environment,project,kind+'-gh-version'))
            direct_version=json.loads(run([alias,'--version','--json'],environment,project,kind+'-direct-version'))
            assert gh_version==direct_version and direct_version['version']=='2.3.1' and direct_version['revision']==EXPECTED_SOURCE
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
        (output/'results.json').write_text(json.dumps(dict(publicTag='v2.3.1',outcomes=outcomes,liveInstallationUnchanged=snapshot(LIVE)==before,packageManifest=package_manifest),indent=2)+'\n')
    finally:
        assert snapshot(LIVE)==before, 'LIVE EXTENSION CHANGED'

if __name__=='__main__':main()
