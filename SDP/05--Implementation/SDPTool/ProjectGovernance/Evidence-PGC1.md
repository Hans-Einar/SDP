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

## PGC1-C — final candidate checks

Code candidate **98b449c** with Go 1.27.1 linux/amd64. Later closeout changes are
documentation/evidence/management only; Go and application-schema bytes are unchanged.
[Final command record](evidence-PGC1/final-checks.json), retained logs and
[binary/artifact manifest](evidence-PGC1/manifest.json) bind the result.

- Full SDPTool `test -race -count=1 ./...` passes: 15 tested packages, including
  facade, core/host/controller/MCP, new compatibility policy and installation paths.
- Full SDPTool `go vet ./...`, Toolkit validator, 85 Toolkit tests, six management
  tests, management replay and exact Session0011 history-import verification pass.
- Final binaries built from 98b449c connect through the real PG controller to actual
  Codex 0.164.0. Initialize/account inspection and idle ephemeral MCP discovery pass;
  actual/minimum version 0.164.0 and sdp-codex-wire/1 are reported. No model turn is
  submitted and no private operational state is created by this probe.
- Sixty captured wire objects validate: 14 against exported default Codex 0.164.0
  schemas, 46 against current governance schemas. Includes current request fields,
  skills, approvals, attachment/recovery, attention and legacy application observations.
- Fresh source-archived owner-lab preparation passes without launching Codex or a
  model. Candidate metadata correctly says stable >=0.164.0 with bounded wire
  compatibility; fixture run remains ready. Existing owner lab remains untouched.

See [schema validation](evidence-PGC1/final-schema-validation.json),
[real no-model probe](evidence-PGC1/live-codex-0.164-probe.log) and
[lab summary](evidence-PGC1/lab-preparation-summary.json).

PGC1 delivers forward **eligibility and bounded wire compatibility**, with actual
0.164.0 negotiation and synthetic compatible-future/refusal tests. It does not
certify every future host, actual model task/approval/recovery on 0.164.0, independent
review, full PG leadership or publication. Those broader gates remain open. Existing
native 0.161.0 evidence remains historical and is not relabeled 0.164.0 evidence.

[Final record checks](final-record-checks-PGC1.json) confirm local links, retained artifact hashes, ledger-prefix preservation and post-closeout Toolkit/management consistency.
