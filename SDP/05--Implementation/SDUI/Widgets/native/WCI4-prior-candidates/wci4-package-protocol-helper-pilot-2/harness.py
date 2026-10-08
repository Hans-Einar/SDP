#!/usr/bin/env python3
"""Bounded WCI4 packaged-command acceptance; Python stdlib only, never launch a GUI.

python3 WCI4-package-protocol.py --package STAGE --out NEW_DIRECTORY \
    [--source PATH_TO_wci4-all-families.sdui]
Use --pilot-helpers only for an explicitly incomplete helper rehearsal. Final main
must run without that flag against the matching four-binary stage. Codegen has no
command in this package: associate the exact-candidate constructor suite separately.
"""
import argparse
import hashlib
import json
import os
from pathlib import Path
import re
import subprocess
import sys
import time

BINARIES = {"sdptool": "bin/sdptool", "sdui": "bin/sdui",
            "preview": "libexec/xfmd/sdui-preview", "native": "libexec/xfmd/sdui-fyne"}
FAMILY_PATHS = ["grouped", "cmd", "legacyButton", "decorated", "toggle", "bound",
                "navigation", "entries", "panes", "pair", "actions", "context",
                "modal", "palette", "open", "flag", "slider", "choice", "amount",
                "basic", "single", "multi", "legacySVG", "image", "plain", "reading",
                "hiddenReading"]


def digest(data):
    return hashlib.sha256(data).hexdigest()


def require(ok, message):
    if not ok:
        raise AssertionError(message)


def json_write(path, value):
    path.write_text(json.dumps(value, ensure_ascii=False, indent=2) + "\n", encoding="utf-8")


class Run:
    def __init__(self, args):
        self.args, self.out, self.package = args, args.out.resolve(), args.package.resolve()
        require(self.out != self.package and self.package not in self.out.parents,
                "--out must be outside the read-only package")
        self.out.mkdir(parents=True, exist_ok=True)
        require(not list(self.out.iterdir()), "--out must be empty; previous evidence is never overwritten")
        (self.out / "harness.py").write_bytes(Path(__file__).read_bytes())
        self.logs = self.out / "commands"
        self.logs.mkdir()
        self.cases, self.commands, self.binaries = [], [], {}
        self.env = dict(os.environ)
        # Neither negative native test may reach app.New/ShowAndRun. Also prevent
        # accidental connection to an existing user display if a binary regresses.
        for name in ("DISPLAY", "WAYLAND_DISPLAY"):
            self.env.pop(name, None)
        for key, relative in BINARIES.items():
            path = self.package / relative
            present = path.is_file() and os.access(path, os.X_OK)
            self.binaries[key] = {"path": str(path), "resolved": str(path.resolve()),
                                  "present": present,
                                  "sha256": digest(path.read_bytes()) if present else None}

    def command(self, name, key, arguments, error=None):
        binary = self.binaries[key]
        require(binary["present"], "missing executable: " + binary["path"])
        argv = [binary["path"], *map(str, arguments)]
        number = len(self.commands) + 1
        prefix = self.logs / (f"{number:03d}-" + name)
        started = time.time()
        record = {"name": name, "argv": argv, "cwd": str(self.out),
                  "binarySha256": digest(Path(argv[0]).read_bytes()),
                  "displayEnvironment": "DISPLAY/WAYLAND_DISPLAY removed",
                  "expected": "success" if error is None else "nonzero with " + error}
        stdout, stderr, code, timed_out = b"", b"", None, False
        try:
            result = subprocess.run(argv, cwd=self.out, env=self.env,
                                    capture_output=True, timeout=30, check=False)
            stdout, stderr, code = result.stdout, result.stderr, result.returncode
        except subprocess.TimeoutExpired as exc:
            stdout, stderr, timed_out = exc.stdout or b"", exc.stderr or b"", True
        except OSError as exc:
            stderr = str(exc).encode("utf-8")
        prefix.with_name(prefix.name + ".stdout").write_bytes(stdout)
        prefix.with_name(prefix.name + ".stderr").write_bytes(stderr)
        record.update(exit=code, timeout=timed_out, seconds=time.time() - started,
                      stdout=str(prefix.with_name(prefix.name + ".stdout").relative_to(self.out)),
                      stderr=str(prefix.with_name(prefix.name + ".stderr").relative_to(self.out)),
                      stdoutSha256=digest(stdout), stderrSha256=digest(stderr),
                      binaryUnchanged=digest(Path(argv[0]).read_bytes()) == binary["sha256"])
        json_write(prefix.with_name(prefix.name + ".json"), record)
        self.commands.append(record)
        require(not timed_out and code is not None, name + ": command failed to complete")
        require(record["binaryUnchanged"], name + ": binary changed during acceptance")
        if error is None:
            require(code == 0, name + f": exit {code}; see retained stderr")
        else:
            require(code != 0 and stdout == b"" and error in stderr.decode("utf-8", "replace"),
                    name + ": expected sourced rejection, nonzero exit and empty stdout")
        return stdout

    def case(self, name, function):
        first = len(self.commands)
        try:
            function()
            self.cases.append({"name": name, "status": "PASS", "commands": [first + 1, len(self.commands)]})
        except Exception as exc:
            self.cases.append({"name": name, "status": "FAIL", "error": str(exc),
                               "commands": [first + 1, len(self.commands)]})

    def source(self, name, data, project=False):
        base = self.out / (name + " project with spaces")
        area = base / "SDP" if project else base
        area.mkdir(parents=True)
        if project:
            (area / "SDP-project.manifest.yaml").write_text(
                "schemaVersion: '1.0'\nproject:\n  name: packaged-widget-fixture\n", encoding="utf-8")
        source = area / "Page å with spaces.sdui"
        source.write_bytes(data)
        return base, source

    def combined(self, name, data, profile):
        _, source = self.source(name, data)
        output = self.out / (name + " Combined output with spaces")
        output.mkdir()
        argv = ["-source", source, "-entry", "Main", "-revision", digest(data), "-output", output]
        result = json.loads(self.command(name + "-combined", "preview", argv))
        require(result["schema"] == "sdptool/0.2" and result["operation"] == "sdui-preview"
                and result["revision"] == digest(data), "wrong helper protocol/revision")
        require(Path(result["entry"]) == output / "entry.md" and Path(result["directory"]) == output,
                "helper response paths differ from actual artifact")
        metadata = json.loads((output / "sdui.json").read_bytes())
        if profile == "sdui/0.2":
            require(set(result) == {"schema", "operation", "entry", "directory", "revision"}
                    and set(metadata) == {"schema", "revision", "source", "frame", "spanUnits", "outputs"},
                    "legacy envelope keys changed")
        else:
            require(result["profile"] == metadata["profile"] == profile, "wrong source profile")
        require(metadata["schema"] == "sdui-combined/1" and metadata["revision"] == digest(data)
                and metadata["source"] == str(source) and metadata["frame"] == "Main"
                and metadata["spanUnits"] == "UTF-8 bytes; end exclusive", "wrong snapshot metadata")
        require((output / "source.sdui").read_bytes() == data, "source snapshot bytes changed")
        for filename in ("source.sdui", "entry.md"):
            require(metadata["outputs"][filename] == digest((output / filename).read_bytes()), "artifact hash mismatch")
        text = (output / "entry.md").read_text(encoding="utf-8")
        links = re.findall(r"\[L(\d+):(\d+)\]\(sdui-source://(\d+)/(\d+)\)", text)
        require(links, "Combined source links absent")
        spans = []
        for line, column, start, end in links:
            start, end, line, column = map(int, (start, end, line, column))
            require(0 <= start < end <= len(data), "invalid source byte range")
            selected = data[start:end].decode("utf-8")
            prefix = data[:start].decode("utf-8")
            require(line == prefix.count("\n") + 1 and column == len(prefix.rsplit("\n", 1)[-1]) + 1,
                    "source link line/column does not address its UTF-8 snapshot")
            spans.append({"start": start, "end": end, "line": line, "column": column, "text": selected})
        json_write(output / "verified-source-spans.json", spans)
        require(any("å" in s["text"] or "🙂" in s["text"] for s in spans), "UTF-8 source link not exercised")
        prior = {p.name: p.read_bytes() for p in output.iterdir()}
        source.write_bytes(data + b"\n")
        self.command(name + "-combined-stale", "preview", argv, "stale")
        require(prior == {p.name: p.read_bytes() for p in output.iterdir()}, "stale helper replaced artifacts")
        source.write_bytes(data)
        return text

    def discovery(self, name, data, profile):
        project, source = self.source(name, data, project=True)
        result = json.loads(self.command(name + "-discover", "sdptool", ["--json", project, "discover"]))
        require(result["schema"] == "sdptool/0.2" and result["status"] == "valid", "wrong discovery envelope")
        models = result["inventory"]["sdui"]
        require(len(models) == 1 and models[0]["profile"] == profile, "wrong discovered profile/inventory")
        require(Path(models[0]["source"]).name == source.name, "discovery lost spaced source path")
        sources = result["sources"]
        require(len(sources) == 1 and sources[0]["state"] == "validated", "source not validated")
        target_nodes = [n for n in result["navigation"]["nodes"]
                        if n["kind"] == "frame" and n["label"] == "Main"]
        require(len(target_nodes) == 1, "Main source target absent")
        target = target_nodes[0]["target"]
        require(target["operation"] == "sdui-preview" and target["revision"] == digest(data)
                and target["model"] == models[0]["id"] and Path(target["path"]) == source,
                "wrong static delegation identity")
        output = self.out / (name + " UIPreview output with spaces")
        argv = ["--json", project, "sdui-preview", "--model", target["model"], "--entry", "Main",
                "--revision", digest(data), "--output", output]
        preview = json.loads(self.command(name + "-ui-preview", "sdptool", argv))
        require(preview["schema"] == "sdptool/0.2" and preview["profile"] == profile
                and preview["revision"] == digest(data) and preview["operation"] == "sdui-preview",
                "UIPreview profile/protocol/revision mismatch")
        require(json.loads((output / "sdptool.json").read_bytes()) == preview, "published result differs")
        require(Path(preview["entry"]) == output / "entry.md", "UIPreview entry path mismatch")
        # Compare real executable delegation with the owning SDUI Markdown export.
        direct = self.command(name + "-direct-markdown", "sdui", [source, "--format", "markdown", "--entry", "Main"])
        require((output / "entry.md").read_bytes() == direct, "UIPreview diverges from SDUI Markdown")
        before = {str(p.relative_to(output)): p.read_bytes() for p in output.rglob("*") if p.is_file()}
        source.write_bytes(data + b"\n")
        self.command(name + "-ui-stale", "sdptool", argv, "stale")
        after = {str(p.relative_to(output)): p.read_bytes() for p in output.rglob("*") if p.is_file()}
        require(after == before, "stale UIPreview replaced artifact")
        source.write_bytes(data)
        return source, direct.decode("utf-8")

    def native_checks(self):
        data = b'sdui 0.2; Main=[field=input("Text",value="old");button("OK")];'
        _, source = self.source("native legacy", data)
        args = ["-source", source, "-entry", "Main", "-revision", digest(data), "-check"]
        ready = json.loads(self.command("legacy-prototype-check", "preview", args))
        require(ready["schema"] == "sdptool/0.2" and ready["operation"] == "sdui-check"
                and ready["status"] == "prototype" and ready["revision"] == digest(data)
                and "profile" not in ready and "no SDL runtime" in ready["diagnostic"],
                "legacy prototype check claimed wrong readiness")
        self.command("preview-stale-check", "preview", ["-source", source, "-entry", "Main",
                     "-revision", "0" * 64, "-check"], "stale")
        self.command("native-stale", "native", ["-entry", "Main", "-watch=false", "-revision", "0" * 64, source], "stale")
        unsupported = {
            "collection": ('tree("Nodes")', "unsupported-provider"),
            "pane": ('tabs("Tabs")[one=page("One")[]]', "unsupported-pane"),
            "interaction": ('dialog("Dialog")[]', "unsupported-interaction"),
            "value": ('checkbox("Flag")', "unsupported-value"),
            "choice": ('select("Choice")', "unsupported-provider"),
            "text": ('input("Text",multiline=false)', "unsupported-text"),
            "svg": ('svg(art.Image.@resource,description="Figure",fallback="label")', "unsupported-preview"),
            "markdown": ('markdown("Plain",description="Reading",fallback="label")', "unsupported-preview"),
        }
        for name, (body, error) in unsupported.items():
            data = ('sdui 0.3; ref: art "unopened"; Main=[p=' + body + ' {visible=false}];').encode()
            _, path = self.source("unsupported " + name, data)
            for binary, args in [("preview", ["-source", path, "-entry", "Main", "-revision", digest(data), "-check"]),
                                 ("native", ["-entry", "Main", "-watch=false", "-revision", digest(data), path])]:
                self.command(name + "-" + binary + "-unsupported", binary, args, error)
                stderr = (self.out / self.commands[-1]["stderr"]).read_text(encoding="utf-8")
                require("Main/p" in stderr and re.search(r"\d+:\d+", stderr), "unsupported adapter diagnostic lost source location")

    def source_exports(self, name, source, profile, all_families=False):
        ast = json.loads(self.command(name + "-ast", "sdui", [source, "--entry", "Main", "--format", "ast"]))
        require(ast["astFormat"] == profile.replace("sdui/", "sdui-ast/")
                and ast["validation"] == "local-profile" and ast["document"]["profile"] == profile, "AST profile mismatch")
        dump = self.command(name + "-dump", "sdui", [source, "--entry", "Main", "--format", "dump"]).decode("utf-8")
        if all_families:
            serialized = json.dumps(ast["document"], ensure_ascii=False)
            for kind in ("tree", "list", "tabs", "page", "split", "command", "menu", "menuGroup", "item",
                         "separator", "dialog", "checkbox", "slider", "select", "number", "svg", "markdown"):
                require('"' + kind + '"' in serialized, "AST missing family " + kind)
            require("resources not supplied; source intent only" in dump and "provider data not supplied" in dump,
                    "all-family dump claims supplied state")
        else:
            svg = self.command(name + "-legacy-svg", "sdui", [source, "--entry", "Main", "--format", "svg"])
            require(svg.startswith(b"<svg") and b"</svg>" in svg, "legacy static SVG unavailable")

    def rejected_exports(self, families):
        bodies = {"all-families": (families, "unsupported-interaction-export")}
        for name, body in (("svg", 'svg(art.Image.@resource,description="Figure",fallback="label")'),
                           ("markdown", 'markdown("Plain",description="Reading",fallback="reject")')):
            bodies[name] = (('sdui 0.3; ref: art "unopened"; Main=[p=' + body + ' {visible=false}];').encode(),
                            "unsupported-resource-export")
        for name, (data, code) in bodies.items():
            _, source = self.source("static rejection " + name, data)
            target = self.out / (name + " prior output.svg")
            old = b"prior artifact\n"
            target.write_bytes(old)
            self.command(name + "-atomic-export", "sdui", [source, "--entry", "Main", "--format", "svg", "-o", target], code)
            require(target.read_bytes() == old, "failed SVG export overwrote old artifact")
            diagnostic = json.loads((self.out / self.commands[-1]["stderr"]).read_bytes())["error"]
            require(diagnostic["code"] == code and diagnostic["span"]["line"] > 0
                    and "Main/" in diagnostic["message"], "unsourced SVG rejection")

    def finish(self, missing):
        changed = [k for k, b in self.binaries.items() if b["present"] and
                   (not Path(b["path"]).is_file() or digest(Path(b["path"]).read_bytes()) != b["sha256"])]
        failed = any(c["status"] == "FAIL" for c in self.cases) or bool(changed)
        status = "FAIL" if failed else "INCOMPLETE" if missing or self.args.pilot_helpers else "PASS"
        artifacts = {str(p.relative_to(self.out)): digest(p.read_bytes())
                     for p in sorted(self.out.rglob("*")) if p.is_file()}
        summary = {"schema": "wci4-package-protocol/1", "status": status,
                   "package": str(self.package), "pilotHelpersOnly": self.args.pilot_helpers,
                   "harnessSha256": digest(Path(__file__).read_bytes()),
                   "source": str(self.args.source.resolve()), "sourceSha256": self.source_hash,
                   "binaries": self.binaries, "missingBinaries": missing, "changedBinaries": changed,
                   "cases": self.cases, "commands": self.commands, "artifacts": artifacts,
                   "codegen": {"status": "NOT_PROVIDED", "reason": "No codegen command in the four-binary package",
                               "requiredSeparateEvidence": "Exact-candidate SDUI codegen suite and emitted-constructor execution"},
                   "scope": "Package protocol only; no OS launch, connected host, IME, installation or publication proof"}
        json_write(self.out / "summary.json", summary)
        print(json.dumps({"status": status, "cases": len(self.cases), "commands": len(self.commands),
                          "summary": str(self.out / "summary.json")}, ensure_ascii=False))
        return 1 if failed else 2 if status == "INCOMPLETE" else 0


def main():
    parser = argparse.ArgumentParser(description=__doc__, formatter_class=argparse.RawDescriptionHelpFormatter)
    parser.add_argument("--package", required=True, type=Path)
    parser.add_argument("--out", required=True, type=Path)
    parser.add_argument("--source", type=Path, default=Path(__file__).resolve().parent /
                        "SDUI/go/parser/testdata/wci4-all-families.sdui")
    parser.add_argument("--pilot-helpers", action="store_true", help="Run helper cases only; INCOMPLETE or FAIL, never final PASS")
    args = parser.parse_args()
    run = Run(args)
    families = args.source.read_bytes()
    run.source_hash = digest(families)
    (run.out / "wci4-all-families.sdui").write_bytes(families)
    missing = [key for key, item in run.binaries.items() if not item["present"]]
    if missing and not args.pilot_helpers:
        return run.finish(missing)
    for version in ("0.2", "0.3"):
        # Multibyte characters precede subsequent components/use sites, exposing
        # rune-versus-byte errors in source links without requiring native glyphs.
        data = ('sdui ' + version + '; Leaf=<field=input("Blå",value="å🙂"),"π prose">; Main=[a=Leaf;b=Leaf];').encode()
        run.case(version + " Combined UTF-8 snapshot", lambda v=version, d=data: run.combined("profile " + v, d, "sdui/" + v))
        if not args.pilot_helpers:
            run.case(version + " discovery/UIPreview/AST", lambda v=version, d=data: profile_case(run, v, d))
    run.case("all-family Combined source intent", lambda: family_combined(run, families))
    run.case("prototype supported and native negative protocols", run.native_checks)
    if not args.pilot_helpers:
        def all_families():
            source, text = run.discovery("all families", families, "sdui/0.3")
            require("resources not supplied; source intent only" in text and "provider data not supplied" in text,
                    "UIPreview implies connected resources")
            run.source_exports("all-families", source, "sdui/0.3", all_families=True)
        run.case("all-family discovery/structural AST", all_families)
        run.case("public SVG rejection preserves artifact", lambda: run.rejected_exports(families))
    return run.finish(missing)


def profile_case(run, version, data):
    source, _ = run.discovery("profile " + version, data, "sdui/" + version)
    ast = json.loads(run.command(version + "-unicode-ast", "sdui", [source, "--format", "ast"]))
    require(ast["astFormat"] == "sdui-ast/" + version and ast["document"]["profile"] == "sdui/" + version,
            "UTF-8 AST profile mismatch")
    _, simple = run.source("legacy export " + version,
                          ('sdui ' + version + '; Main=[input("Text");button("OK");"Plain"];').encode())
    run.source_exports("legacy " + version, simple, "sdui/" + version)


def family_combined(run, families):
    text = run.combined("all families", families, "sdui/0.3")
    for name in FAMILY_PATHS:
        require("Main/" + name in text, "Combined omits family path " + name)
    require("resources not supplied; source intent only" in text and "provider data not supplied" in text,
            "Combined implies connected resources")


if __name__ == "__main__":
    try:
        sys.exit(main())
    except (OSError, AssertionError) as exc:
        print("harness setup failed: " + str(exc), file=sys.stderr)
        sys.exit(2)
