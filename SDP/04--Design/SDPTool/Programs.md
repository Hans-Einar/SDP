# SDPTool application discovery and launch

## SDPTOOL-DES-RSP1

Realizes SDPTOOL-REQ-RSP1. The contract lives in SDPTool/Programs.md beside the
implementation. SDPTool owns an optional explicit project program catalog,
read-only availability validation, typed discovery output and foreground process
launch. The application owns interpretation, Go bindings, providers and native
window lifecycle. XFMD is a future consumer of this same API.

Decision source: owner Session0010 T010. This replaces the suggested XFMD-specific
launch integration and creates a narrow exception to the earlier no-executable-
metadata policy. Ordinary SDUI source/preview files still do not execute commands.
The explicit run operation is the execution boundary; discovery is observational.

Chosen representation: SDP/programs.json with stable program ID, label, SDUI
source/entry and argv. A separate catalog avoids new SDUI grammar and supports
multiple programs sharing sources. A command such as make run reuses an existing
application composition root. Inferring commands from any Makefile, duplicating
registrations in XFMD, and introducing a Go plugin loader were rejected because
they do not establish the source/application association or add unnecessary runtime
ownership. Source/entry existence and command availability are static readiness,
not proof of every application binding or successful compilation.

Verification: source-preserving invalid catalog behavior, literal argv/cwd,
revision selection, exit/cancellation, real gh-sdp bootstrap routing and native
widget-lab SDL -> Go return. PLAN-SDP-0023 records exact candidate evidence.
