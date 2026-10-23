# PGC1 — Codex compatibility evidence

Session0007 T029 / PLAN-SDP-0018 / KB057, 24 October 2026.

## Implemented boundary — PGC1-B

New public host connections admit stable numeric Codex versions >= 0.164.0 under
the bounded existing wire contract. Executable and initialize identity must agree;
private sockets use observed server identity and existing peer/lifetime checks.
Unknown optional fields are tolerated, while required operation/response failures
refuse. Later eligible versions are not automatically workflow-certified.

Probe output, controller observations and recovery bindings preserve actual version;
unmapped replay cannot invent it. Different-version recovery/reattachment refuses
without relabeling historical attempts. Legacy 0.160.1/0.161.0 observations remain
readable and valid in core history. Version schemas and lab metadata agree with
current policy, with archived engine requirements derived from selected source.
No role authority, reporting acceptance or permission policy is widened.

Scoped race tests pass for codexcompat, governancehost, governancecontroller,
governance and cmd/sdp-controller. A subsequent MCP/core-host schema check passes
through governancemcp, codexcompat and governancehost after updating the shared
host-version schema. New tests cover a compatible future 0.190.3 host/socket,
actual version persistence, different-version recovery refusal without mutation,
handshake mismatch, missing required fields/methods, null/unknown result states,
legacy replay identity and unsupported/preview version rejection.

Eleven captured work/skill/approval request/response objects validate against actual
exported default Codex 0.164.0 schemas. A development binary's real 0.164.0 account/
idle-thread/MCP discovery probe succeeds with modelSubmitted=false. Final committed
candidate probe, full affected-module race/vet and artifact hashes are PGC1-C and
will be recorded separately below. No native model task or independent review is
claimed by these checks.

The first Python schema check found that a bare regex end anchor accepts a terminal
newline; the final shared schema adds a portable newline rejection and the runtime
already refuses it. A brief uncommitted lookahead draft was replaced before any
release because the Go regex engine does not support it. The final MCP SDK schema
transport check passes; no hidden failing gate is relabeled successful.

## Preserved open work

Current-host native task/approval/recovery workflow evidence, independent review,
leader report consumption/coordination, canonical publication and sustained capacity
remain separate PG gates. No Ponsse/MVP1 changes, stable install, merge or release.
