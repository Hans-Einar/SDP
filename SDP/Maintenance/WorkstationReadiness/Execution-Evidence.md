# WR2 execution evidence — consolidated SDP and PG MCP readiness

Session0013 T002 / MAINT-SDP-0018 / KB-SDP-056, 23 October 2026.
This is same-context verification of the authorized consolidation round, not
independent review, native agent acceptance, a release or full PG completion.

## Candidate and environment

- Source: commit 4174c68 on sdp/workstation-completion-assessment, based on
  consolidated main bfe3b40. Its subsequent changes are current documentation and
  management/evidence only; Go implementation and runtime schemas are unchanged.
- Go 1.27.1 linux/amd64, privately provisioned in the workstation cache. Official
  distribution archive SHA-256 verified before extraction; [record](evidence/WR2/toolchain.json).
  No system package, default PATH or stable gh-sdp engine was replaced.
- Three private binaries: sdptool, sdp-controller and sdp-mcp. Actual version probe
  identifies development-WR2 / 4174c68; [build and hashes](evidence/WR2/checks.json).
- Codex CLI/app-server 0.164.0. The production PG adapter remains pinned to
  0.161.0. No version guard was removed and no model task was submitted.
- Private trial directories live under /home/warloc/.cache/sdp-workstation/session0013-WR2;
  project files, operational state, toolchain and module/build caches are not added
  to the product tree. [Artifact manifest](evidence/WR2/manifest.json) binds retained logs/harnesses.

## C1 / WR2-M1 — candidate build and regression checks

| Executed command family | Result | Boundary |
| --- | --- | --- |
| SDPTool build of sdptool, sdp-controller, sdp-mcp with development identity | PASS | Current local candidate, not a published package. |
| go -C SDPTool test -race -count=1 ./... | PASS; 14 tested packages | Core, facade, model/Blueprint, controller/host/MCP and installation tests; protocol hosts in tests include fixtures. |
| go -C SDL/go test -race -count=1 ./... | PASS; 21 tested packages | Existing parser, source graph, Blueprint, projection/runtime/package regressions. |
| go -C SDUI/go test -race -count=1 ./... | PASS; 14 tested packages | Existing semantic/runtime/presentation/host regressions, including test doubles; no fresh native GUI inspection. |
| go -C SDPTool/bootstrap test -race -count=1 ./... | PASS; one tested package | Bootstrap module regressions. |
| go vet ./... for SDPTool, SDL/go and SDUI/go | PASS | Static checks on those three modules. |
| Current binary --version and actions catalogue | PASS | Six actions including project.model-tree. The installed signed 2.2.0 engine remains separate. |
| Current controller probe against installed Codex | EXPECTED REFUSAL, exit 2 | incompatible_profile: require codex-cli 0.161.0; confirms a current execution blocker rather than a passing host probe. |

All command lines, exact exit codes, durations and logs are in
[checks.json](evidence/WR2/checks.json) and [runtime probes](evidence/WR2/runtime-probes.json).
The 50 tested-package results are not a count of individual test cases. No passing
fixture test is presented as current native model execution.

## C2 / WR2-M2 — current authority and responsibility consolidation

[Scoped preflight](evidence/preflight-WR2-M2.json) records the clean candidate and
bounded documentation write set. No demonstrated semantic duplicate implementation
required a code refactor. The following current sources are reconciled:

- [State responsibility map](../../04--Design/SDPTool/ProjectGovernance/State-Boundaries.md)
  separates model artifact, retained Blueprint, canonical assignment, prepared
  package, decision, run, host attempt, report and canonical project outcome.
- [MCP README](../../../SDPTool/governancemcp/README.md) now lists both implemented
  attention tools and explains three unbound versus seven fully bound tools.
- [Orchestration design](../../04--Design/SDPTool/ProjectGovernance/Orchestration.md)
  no longer presents the old five-tool count or already-reconciled branch collisions
  as current unresolved facts. Global cross-worktree allocation remains unfinished.
- [Actions usage](../../../SDPTool/Actions.md) now includes the implemented
  project.model-tree operation and lazy discovery parameters already specified by
  the producer contract. No new runtime feature or privilege was introduced.

Historical evidence/journals, owner pauses, unrelated product files and prior
accepted decisions remain intact. The current entrypoints link the new study
instead of directing readers to nonexistent old development worktrees.

## C3 / WR2-M3 — integrated checks at their actual level

| Trial | Actual result | What remains unverified |
| --- | --- | --- |
| Existing provider-backed design trial with real compiled binaries and checked 0.6 specimen | Creates real ModelGovernance WORKs, retained Blueprint and canonical assignment; prepares PG package; current review digest is deterministic; uncommitted TARGET drift blocks provider review without writes/approval; exact restoration restores the basis | Fixture assignment; no native implementation, real owner decision or full project adoption. |
| Actual compiled unbound MCP over stdio | Negotiates 2025-11-25; lists context_read/route_evaluate/run_read; reads this project's coverage with no conditions; status route is no-work; unbound run and unavailable mutations refuse; operational state not created | Direct protocol client rather than a model choosing tools. |
| Actual Codex 0.164.0 app-server discovery and read calls | Ephemeral idle thread discovers server; three host-mediated tool calls read coverage/status and confirm unbound-run refusal; no PG state/config changes and no turn/start | Does not validate the pinned PG execution controller, model submission, approval interaction or native recovery. |
| Fully bound compiled worker MCP | Seven tools present; context, events and attention reads succeed; exact duplicate report gives same receipt; run becomes verifying; complete and controller start refuse; all fixture preimages preserved | Host/attempt binding scripted explicitly; report text is fixture evidence, not actual implementation or owner acceptance. |
| Existing full race/regression suites | Current provider freshness, activation, attention/reply/continuation, replay/refusal and recovery checks pass | Independent review and a whole current native owner-to-accepted-result workflow remain outstanding. |

Retained artifacts: [provider trial](evidence/WR2/provider-trial.json),
[preliminary input refusals](evidence/WR2/provider-input-attempts.json),
[direct MCP](evidence/WR2/mcp-readonly.json),
[bound worker MCP](evidence/WR2/mcp-bound-worker.json),
[actual Codex reads](evidence/WR2/codex-mcp-readcalls.json).
The trial reuses the existing reduced specimen; it does not create or change a
Ponsse pilot model. Preliminary profile/input failures were retained, not hidden.

## Revised C4 / WR2-M4 — MCP adoption study

[MCP adoption study](MCP-Adoption-Study.md) delivers usable modes, measured blockers,
a concrete [inspection configuration](mcp-readonly.example.toml) and recommended
next increments. The TOML parses and Codex accepts the configuration through
command-line-only overrides; [validation](evidence/WR2/config-validation.json).
No persistent global MCP registration was installed. Current host schema export
is summarized in [schema inventory](evidence/WR2/schema-inventory.json); this is
not a completed protocol migration assessment.

An additional [implementation-context check](evidence/WR2-implementation-context.json)
confirms that unbound context_read's shared service remains Maintenance-specific:
PG's active ImplementationPlan produces explicit coverage conditions. The successful
WR2 selection was active at observation time; it becomes historical at closeout.
This does not certify generic project/leader inspection or justify wrapper plans.

## Claims and remaining work

The authorized consolidation and adoption investigation can finish with these
bounded results. The selected Go/source integration checks pass, current docs are
reconciled and local inspection is demonstrated. Full PG/MCP execution readiness
is **not verified**: native 0.164.0 execution compatibility, leader report consumption,
SAD/coordination, independent review, canonical outcome publication and capacity/
recovery gates retain their owning PG plan/cards. Do not close those gates from
this evidence or claim this round finishes the broader SDP system.

The owner-requested MVP1 deferral is preserved. Next bounded PG work is current-host
compatibility and then leader/control/report/evidence integration, as detailed in
the adoption study. This round adds no merge, release, application migration or
acceptance of the existing paused Blueprint pilot.

## Final WR2-M5 record checks

[Final record checks](evidence/WR2-final-checks.json) confirm resolved local links,
valid JSON, unchanged retained runtime artifact hashes and original ledger prefix.
Final management replay passes with 70 cards, 47 records, four lineage operations
and 637 events; Toolkit validation passes. The 85 Toolkit and six management tests
pass on this round's documentation candidate. Closure does not extend the runtime
or independent-review claims above.
