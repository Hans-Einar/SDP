# SDP shared schemas, payload and validation

SDPTool is the sole installation and upgrade engine. Use `gh sdp install` or
`gh sdp upgrade`, or invoke the compiled SDPTool directly. See
[SDPTool](../SDPTool/README.md) and [installation records](../SDPTool/install/Records.md).

The current release inventory is [SDPTool/profiles/payload.json](../SDPTool/profiles/payload.json).
Go builds and signs the immutable release descriptor. No legacy profile artifact
or alternate installation script is built or executed by the current toolchain.

Toolkit retains shared schemas, managed payload, Python document/schema validators,
and CLI helper scripts. The old install-v1 manifest, schemas and expected records
remain data for interpreting historical contracts; they do not select the current
payload or provide an executable installer. Legacy templates are archived under
Template/legacy. Project-owned current templates live in Template/sdp-root.

Validation: `python3 Toolkit/scripts/validate_sdp.py` and Go tests in SDPTool.
The CI pipeline runs Go installation, signed-package and recovery checks.
Historical implementation and test results remain dated records, not instructions
to execute retired tools. Published old assets remain immutable.
