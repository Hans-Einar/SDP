#!/usr/bin/env python3
"""Probe real SDL/SDUI capability boundaries; no producer files are modified."""
import argparse
import hashlib
import json
import subprocess
from pathlib import Path

HERE = Path(__file__).resolve().parent


def main():
    parser = argparse.ArgumentParser()
    parser.add_argument("--sdl", required=True)
    parser.add_argument("--sdui", required=True)
    parser.add_argument("--producer-commit", required=True)
    args = parser.parse_args()
    fixtures = HERE / "fixtures"
    outputs = HERE / "outputs"
    fixtures.mkdir(exist_ok=True)
    outputs.mkdir(exist_ok=True)
    core = "language design-core version 0.5.\n"
    cases = [
        ("sdl-core", "sdl", core + "container XfmdDesktop.\n", "check", 0, None),
        ("sdl-system", "sdl", core + "system XFMD.\n", "check", 1, "UNSUPPORTED_SYNTAX"),
        ("sdl-include", "sdl", core + 'include "other.design".\n', "check", 1, "UNSUPPORTED_SYNTAX"),
        ("sdl-requirement", "sdl", core + "requirement UR001.\n", "check", 1, "UNSUPPORTED_SYNTAX"),
        ("sdui-controls", "sdui", 'sdui 0.2; page=[go=button("Open"), path=input("Path",value="README.md")];\n', "svg", 0, None),
        ("sdui-scroll-ast", "sdui", 'sdui 0.2; page=[button("Open")] {overflow-y=scroll};\n', "ast", 0, None),
        ("sdui-scroll-svg", "sdui", 'sdui 0.2; page=[button("Open")] {overflow-y=scroll};\n', "svg", 2, "unsupported-scroll"),
        ("sdui-missing-binding", "sdui", 'sdui 0.2; ref: domain "missing.sdl"; page=[go=button("Open",callback=domain.document.@open)];\n', "ast", 0, None),
        ("sdui-icon", "sdui", 'sdui 0.2; page=[button("Open",icon="open")];\n', "ast", 2, "widget-property"),
        ("sdui-checked", "sdui", 'sdui 0.2; page=[button("Markdown",checked=true)];\n', "ast", 2, "widget-property"),
        ("sdui-context-event", "sdui", 'sdui 0.2; ref: domain "missing.sdl"; page=[go=button("Open",onContext=domain.path.@menu)];\n', "ast", 2, "widget-property"),
    ]
    for widget in ("tree", "tabs", "splitter", "list", "menu", "dialog", "textarea", "slider", "checkbox", "select"):
        cases.append(("sdui-widget-" + widget, "sdui", f'sdui 0.2; page=[{widget}("Example")];\n', "ast", 2, "widget-kind"))
    records = []
    for name, tool, source, operation, expected, diagnostic in cases:
        path = fixtures / (name + (".design" if tool == "sdl" else ".sdui"))
        path.write_text(source)
        binary = str(Path(getattr(args, tool)).resolve())
        command = [binary, operation, str(path)] if tool == "sdl" else [binary, str(path), "--format", operation, "--entry", "page"]
        result = subprocess.run(command, capture_output=True, text=True, timeout=30)
        (outputs / (name + ".stdout")).write_text(result.stdout)
        (outputs / (name + ".stderr")).write_text(result.stderr)
        matched = result.returncode == expected and (diagnostic is None or diagnostic in result.stdout + result.stderr)
        records.append({"case": name, "command": command, "fixtureSha256": hashlib.sha256(path.read_bytes()).hexdigest(),
                        "exitCode": result.returncode, "expectedExitCode": expected,
                        "expectedDiagnostic": diagnostic, "expectationMet": matched})
    report = {"producerCommit": args.producer_commit,
              "binaries": {name: {"path": str(Path(getattr(args, name)).resolve()),
                                   "sha256": hashlib.sha256(Path(getattr(args, name)).read_bytes()).hexdigest()}
                           for name in ("sdl", "sdui")},
              "cases": records}
    (HERE / "probe-results.json").write_text(json.dumps(report, indent=2) + "\n")
    for record in records:
        print(("PASS" if record["expectationMet"] else "FAIL"), record["case"], "exit", record["exitCode"])
    raise SystemExit(0 if all(r["expectationMet"] for r in records) else 1)


if __name__ == "__main__":
    main()
