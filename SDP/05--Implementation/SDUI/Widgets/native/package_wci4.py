#!/usr/bin/env python3
"""Prepare an isolated, source-inventoried WCI4 candidate; never install or publish.

Uses the existing SDPTool packaging entrypoint and the supplied consumer's helper
recipe. Records actual commands, dependency modules, notices and binary build facts.
Native/protocol acceptance is a separate operation on the resulting exact binaries.
"""
import argparse
import gzip
import hashlib
import json
import os
from pathlib import Path
import re
import shutil
import subprocess
import time


def sha(path):
    return hashlib.sha256(Path(path).read_bytes()).hexdigest()


def write(path, value):
    path.parent.mkdir(parents=True, exist_ok=True)
    path.write_text(json.dumps(value, indent=2, sort_keys=True) + '\n')


def inventory(root):
    names = subprocess.check_output(['git', 'ls-files', '-z', '--cached', '--others', '--exclude-standard', '--', 'SDUI', 'SDL', 'SDPTool'], cwd=root).split(b'\0')
    result = {}
    for raw in sorted(set(names)):
        if not raw:
            continue
        name = os.fsdecode(raw)
        p = root / name
        if p.is_file() and '__pycache__' not in p.parts:
            result[name] = sha(p)
    return result


def objects(data):
    decoder = json.JSONDecoder()
    offset = 0
    while offset < len(data):
        while offset < len(data) and data[offset].isspace():
            offset += 1
        if offset == len(data):
            break
        value, offset = decoder.raw_decode(data, offset)
        yield value


def main():
    p = argparse.ArgumentParser(description=__doc__)
    p.add_argument('--source', required=True, type=Path)
    p.add_argument('--consumer', required=True, type=Path)
    p.add_argument('--out', required=True, type=Path)
    p.add_argument('--gocache', required=True, type=Path)
    a = p.parse_args()
    root, consumer, out = a.source.resolve(), a.consumer.resolve(), a.out.resolve()
    out.mkdir(parents=True, exist_ok=False)
    evidence = out / 'evidence'
    evidence.mkdir()
    env = dict(os.environ, GOWORK='off', GOFLAGS='-mod=readonly', GOCACHE=str(a.gocache.resolve()), GOTOOLCHAIN='local')
    # An ambient SDP version/tool override would misidentify the existing recipe.
    env.pop('SDP_VERSION', None)
    env.pop('SDP_GO', None)
    receipts = []

    def run(name, argv, cwd=root, expected=0):
        begin = time.time()
        result = subprocess.run(list(map(str, argv)), cwd=cwd, env=env, stdout=subprocess.PIPE, stderr=subprocess.STDOUT)
        log = evidence / (name + '.log.gz')
        with gzip.open(log, 'wb') as f:
            f.write(result.stdout)
        receipts.append(dict(name=name, argv=list(map(str, argv)), cwd=str(cwd), exitCode=result.returncode,
                             elapsedSeconds=round(time.time()-begin, 3), log=log.relative_to(out).as_posix(),
                             rawLogSHA256=hashlib.sha256(result.stdout).hexdigest(),
                             environment={k: env[k] for k in ('GOWORK','GOFLAGS','GOCACHE','GOTOOLCHAIN')}))
        write(evidence / 'commands.json', receipts)
        print(name, 'exit', result.returncode, flush=True)
        if result.returncode != expected:
            raise RuntimeError(f'{name}: exit {result.returncode}; see {log}')
        return result.stdout.decode(errors='strict')

    consumer_paths = ('tools/build_sdui_tools.py','tools/run_with_xvfb.py','cmake/Sdp.cmake','cmake/Tests.cmake',
                      'src/application/sdp/SduiTools.cpp','src/application/sdp/SduiWorkflow.cpp',
                      'src/application/sdp/SdpProtocol.cpp','tests/gui/SduiWorkflowGuiTest.cpp')
    consumer_before = {name: sha(consumer/name) for name in consumer_paths}
    write(evidence/'consumer-source-before.json',consumer_before)
    before = inventory(root)
    write(evidence / 'source-before.json', before)
    recipe = consumer / 'tools/build_sdui_tools.py'
    write(evidence / 'provenance.json', dict(sourceRoot=str(root), sourceHEAD=run('source-head', ['git','rev-parse','HEAD']).strip(),
        sourceStatus=run('source-status', ['git','status','--porcelain']), consumerRoot=str(consumer),
        consumerHEAD=run('consumer-head', ['git','rev-parse','HEAD'], consumer).strip(),
        consumerRecipe=str(recipe), consumerRecipeSHA256=sha(recipe), packageRecipeSHA256=sha(root/'SDPTool/package.sh'),
        harnessSHA256=sha(__file__), scope='Local candidate preparation; no installation, publication or licensing grant.'))
    run('compiler', ['go','version'])
    run('compiler-environment', ['go','env','GOOS','GOARCH','CGO_ENABLED','CC','GOVERSION'])
    run('glfw-patch-verification', ['python3',root/'SDUI/third_party/glfw-policy/verify.py'])
    run('sdptool-package', ['bash',root/'SDPTool/package.sh',out/'producer'])
    (out/'bin').mkdir()
    shutil.copy2(out/'producer/sdptool', out/'bin/sdptool')
    run('consumer-helper-recipe', ['python3',recipe,'--source',root/'SDUI/go','--output',out/'libexec/xfmd'])
    run('sdui-cli', ['go','build','-o',out/'bin/sdui','./cmd/sdui'], root/'SDUI/go')
    fixtures = ('collections','panes','commands','values','text','previews')
    for family in fixtures:
        if not (root/'SDL/go/examples'/family/'cmd/native/main.go').is_file():
            raise RuntimeError('Missing connected fixture source: '+family)
    shutil.copytree(root/'SDUI/third_party/glfw',out/'sources/dependencies/glfw')
    shutil.copytree(root/'SDUI/third_party/glfw-policy',out/'sources/dependencies/glfw-policy')
    launches = []
    for family in fixtures:
        src = root/'SDL/go/examples'/family
        for file in sorted(src.rglob('*')):
            if file.is_file() and file.suffix in ('.go','.md','.sdl','.sdui'):
                dest = out/'sources/fixtures'/family/file.relative_to(src)
                dest.parent.mkdir(parents=True,exist_ok=True)
                shutil.copy2(file,dest)
        launches.append(dict(family=family,argv=['libexec/sdui/'+family],environment={'FYNE_THEME':'light'},
            source='sources/fixtures/'+family,scope='Connected acceptance application with explicit SDL bindings/providers. Use an isolated display/configuration.'))
    write(out/'launches.json',launches)
    for family in fixtures:
        run('fixture-'+family, ['go','build','-tags','desktop','-o',out/'libexec/sdui'/family,'./examples/'+family+'/cmd/native'], root/'SDL/go')

    modules = {}
    targets = [('sdptool','SDPTool',[],['./cmd/sdptool']),
               ('helpers','SDUI/go',['-tags','desktop'],['./cmd/sdui','./cmd/sdui-preview','./cmd/sdui-fyne']),
               ('fixtures','SDL/go',['-tags','desktop'],['./examples/'+n+'/cmd/native' for n in fixtures])]
    for name, sub, flags, packages in targets:
        run('module-graph-'+name, ['go','list','-m','-json','all'], root/sub)
        data = run('compiled-dependencies-'+name, ['go','list','-deps','-json',*flags,*packages], root/sub)
        for pkg in objects(data):
            module = pkg.get('Module')
            if not module:
                continue
            actual = module.get('Replace', module)
            key = (module['Path'], module.get('Version',''), actual.get('Dir',''))
            record = modules.setdefault(key, dict(module=module['Path'],version=module.get('Version',''),directory=actual.get('Dir',''),
                replacement=module.get('Replace'),buildRoots=[],licenses=[]))
            if sub not in record['buildRoots']:
                record['buildRoots'].append(sub)
    missing = []
    for number, record in enumerate(sorted(modules.values(), key=lambda r:r['module'])):
        directory = Path(record['directory'])
        if record['module'].startswith('github.com/Hans-Einar/SDP/'):
            record['licenseScope'] = 'Project source; no repository-wide license grant asserted by this preparation.'
            continue
        files = sorted(f for f in directory.rglob('*') if f.is_file() and re.match(r'^(licen[cs]e|copying|notice|copyright)([._-].*)?$', f.name, re.I))
        for f in files:
            relative = f.relative_to(directory)
            dest = out/'licenses'/f'{number:02d}'/relative
            dest.parent.mkdir(parents=True, exist_ok=True)
            shutil.copy2(f,dest)
            record['licenses'].append(dict(source=str(f),path=dest.relative_to(out).as_posix(),sha256=sha(dest)))
        if not files:
            missing.append(record['module'])
    goroot = run('goroot',['go','env','GOROOT']).strip()
    shutil.copy2(Path(goroot)/'LICENSE',out/'licenses/Go-LICENSE')
    write(evidence/'dependency-licenses.json',dict(modules=sorted(modules.values(),key=lambda r:r['module']),missing=missing,
        goLicenseSHA256=sha(out/'licenses/Go-LICENSE'),nativeLibraries='Recorded as system dependencies; their binaries are not bundled.'))
    if missing:
        raise RuntimeError('Missing third-party license records: '+str(missing))

    binaries = [out/'bin/sdptool',out/'bin/sdui',out/'libexec/xfmd/sdui-preview',out/'libexec/xfmd/sdui-fyne'] + [out/'libexec/sdui'/n for n in fixtures]
    facts = []
    for binary in binaries:
        name = binary.relative_to(out).as_posix().replace('/','-')
        build = run(name+'-build-info',['go','version','-m',binary])
        native = subprocess.run(['ldd',str(binary)],env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT)
        (evidence/(name+'-ldd.txt')).write_bytes(native.stdout)
        if b'not found' in native.stdout or (native.returncode != 0 and b'not a dynamic executable' not in native.stdout and b'statically linked' not in native.stdout):
            raise RuntimeError('Unresolved native dependency inspection for '+str(binary))
        facts.append(dict(path=binary.relative_to(out).as_posix(),sha256=sha(binary),bytes=binary.stat().st_size,
                          buildInfo=build,lddExit=native.returncode,lddLog=name+'-ldd.txt'))
    write(evidence/'binaries.json',facts)
    consumer_after = {name: sha(consumer/name) for name in consumer_paths}
    write(evidence/'consumer-source-after.json',consumer_after)
    if consumer_before != consumer_after:
        raise RuntimeError('Inspected consumer inputs changed during package preparation.')
    after = inventory(root)
    write(evidence/'source-after.json',after)
    if before != after:
        raise RuntimeError('Candidate source changed during package build; package not accepted.')
    (out/'README.txt').write_text('WCI4 local candidate preparation, not a published release.\n'
        'bin/sdptool is copied byte-for-byte from the existing producer package.\n'
        'libexec/xfmd contains helpers from the unmodified consumer build recipe.\n'
        'libexec/sdui contains connected acceptance applications, not standalone launcher support.\n'
        'Native/protocol verification is recorded separately against these binary hashes.\n'
        'Native preview verification uses FYNE_THEME=light; dark/system-theme contrast is not verified.\n'
        'Third-party notices are copied for compiled Go dependencies, including the pinned GLFW tree.\n'
        'System native libraries are required and inventoried, not bundled.\n'
        'No project-wide license grant, redistribution approval or OS screen-reader support is asserted.\n')
    write(evidence/'package-inventory.json',{str(p.relative_to(out)):sha(p) for p in sorted(out.rglob('*')) if p.is_file()})
    print('Prepared',out,flush=True)

if __name__ == '__main__':
    main()
