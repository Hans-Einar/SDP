# PGL3-C6 retained evidence

The final full SDPTool race/vet runner is `checks.py`; original private paths are
retained, so adjust those paths on another host. `candidate.json` binds baseline
plus tested source/test/schema bytes. `checks.json` and logs are the final pass.
The final output-only error schema correction is explicitly recorded in the candidate.

`mcp-wire.json` comes from a compiled stdio coordinator server and includes actual
view, report acknowledgement/result and tool list. `console-wire.json` preserves
the exact event objects emitted by the compiled console, including invalid-input
error, attributed reply, newly arriving report and acknowledgement. No native
model task is used. `binaries.json` identifies private development builds; the
stable installation and global client configuration are unchanged.

Run `python3 PATH/verify_wire.py` from the repository root. It resolves adjacent
canonical schema files locally through a registry; it does not fetch remote schemas.
`schema-checks.json` records ten validated final objects.

`console-capture-before-correction.json` is rejected preliminary evidence: a test
helper introduced null optional fields by re-encoding decoded events. It is not the
actual wire format and is not used to establish conformance. The final helper retains
raw emitted JSON; final full-suite and schema checks were performed after correction.
