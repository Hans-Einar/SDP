# Native PGL3-C8 evidence

Actual Codex 0.164.0 / gpt-6-astra high in the owner-authorized pg_lab project.
`native-summary.json` pins source, binaries, lab baseline and scope. `native/`
contains filtered transport traces, fault markers, real CLI projections and context,
external behavior checks and retained failed diagnostics. No raw host/auth/reasoning
logs are retained. Process-local wrappers disable unrelated Apps/plugins/delegation.

`verify_native.py` validates captured JSON and exact recovery requests. `checks.py`
runs the affected product race/vet checks using the workstation toolchain. The lab's
`./pg` helper reads current retained state without launching a model.

Retained preparation scripts are records of the actual cache-root execution. They
are not safe reset commands: do not execute them in this evidence directory or reuse
an uncertain store. For a new native trial, use a clean fixture at the recorded lab
baseline and a new private runtime directory, then explicitly select model execution.
The native fault proxy is maintained in SDPTool/labs/project-governance.

The native worker run is verifying with one report receipt; no acceptance, canonical
publication, independent review or full PG readiness is asserted.
