# PGI3-C5 — bounded launch provenance

PLAN-SDP-0018 / KB-SDP-038, 2026-10-09. Candidate: parent `4044f37f7c416f5fe3e79f793d938f1bfed355d8` plus
[the recorded product/test file set](live-PGI3-C5.json), digest `fbb1c76d25495d2ae6ed86f1f5a67362cbb27dc85e3573599d6acb04e248c377`.
Same-context implementation and verification; no independent review, owner acceptance,
full PGI3 delivery or release is claimed.

## Delivered behavior

The controller persists a bounded launch-context observation before submitting work.
The read-only `sdp-controller provenance` command exposes the controller-assigned
Master role, governing source snapshots, host-reported session/context metadata and
project skill availability. It works without a Codex executable or a host lease.
After an observation epoch changes, the snapshot is explicitly historical.

Host instruction paths do not prove consumed file bytes. A local SHA-256 snapshot
identifies locally inspected bytes. Enabled skill discovery does not prove invocation
or compliance: useEvidence remains not-observed. Missing optional protocol metadata
is unavailable, not fabricated empty evidence. External skill bodies are not read or
persisted. Discovery does not change enablement, account configuration, launch text,
MCP permissions or recovery correlation. The observation is controller-only and
cannot grant agent authority. See the [design](../../../04--Design/SDPTool/ProjectGovernance/Design.md)
and [CLI instructions](../../../../SDPTool/governancecontroller/README.md).

## Actual workflow and candidate boundary

Installed Codex 0.161.0 and the [official app-server reference](https://learn.chatgpt.com/docs/app-server)
provide instructionSources and skills/list discovery. The disposable project copies
the maintained Skills collection and normal discovery symlinks; it is not another
maintained source. Initial fixture setup refused a dot-prefixed .agents assignment
exclusion before a run or model turn existed. The corrected fixture preserves those
links explicitly and keeps only two Go files writable; the product path contract
is unchanged. Earlier trials passed intermediate candidates; trial4 is the retained
final runtime-source trial.

Trial4 fixed Unicode label trimming, passed driver Go tests, preserved every other
tracked fixture preimage and discovery link, and completed one explicit report
approval scenario. The host completed; the routine reached verifying, not completed.
The snapshot reports 14 enabled project skills with matching digests, 60 external
skills and zero discovery errors. It reports both global and project instruction
paths. Native role and parent fields are empty; source is the literal host value
vscode. None establishes a native role hierarchy. The agent's statements about
reading skills remain self-reported messages, not observed invocation evidence.

After this actual run, captured schema validation found an unresolved Source
reference in the new definition. It was replaced with the existing inline source
schema, and the checker now resolves all local references, including unused ones.
All final runtime Go sources were used by trial4, but this schema correction changes
the embedded binary identity. Rebuilt binaries were not submitted to another model
turn: their offline provenance projection is equal to trial4, all captured records
validate, and final affected-package race tests including the compiled controller
CLI pass. The JSON records both binary identities and this precise boundary.

## Verification and remaining work

Final affected-package race tests, compiled controller CLI, vet, 11 captured
application cases, 13 installed-host protocol cases and 21 live records/projections
pass. Tests cover malformed/missing discovery fields, project confinement, duplicate
paths, invalid claims, idempotent observations, read-only projection, unavailable
metadata and historical epochs. A negative fixture proves that an unused unresolved
schema reference is rejected. No full-repository regression claim is made.

The accompanying JSON retains source identities, actual mapped events/provenance,
pre/postimages, discovery links, explicit driver decision, outputs and reproduction
scripts. Management/Toolkit validation and final ledger/link integrity are recorded
there after documentation updates. Historical bytes and concurrent shared work are
preserved; known private C1 event collisions still need integration reconciliation.

Next is explicit skill invocation/correlation and applicable native role mapping
under current delegation rules. Text-only lost-ack matching must not be silently
extended to skill inputs. PGI4 retains checked evidence/publication; PGI5 retains
independent integrated verification. No delegation, merge, push or release is selected.
