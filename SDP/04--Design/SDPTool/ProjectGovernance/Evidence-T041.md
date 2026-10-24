# T041 — Session execution design and SDL evidence

Scope: owner-selected concretization under Session0007 / PLAN-SDP-0018 / KB052,
with KB038 and KB037. Baseline code: 7850fdfb83d3590eab843f85d2e3ba7ffc68b652.
No product implementation, lab agent execution, model release or independent review
is claimed. [Contract](Session-Execution.md), [plan](../../../05--Implementation/SDPTool/ProjectGovernance/Plan.md)
and [model area](../../../SDL/ProjectGovernance/SDPTool/SessionWorkflow/README.md)
are the delivery. Canonical work records remain active.

## Observed validation

| Check | Actual result | Limit |
| --- | --- | --- |
| Build existing SDL and SDPTool tools | Both pass with cached Go 1.27.1, local toolchain/module cache | Initial unqualified go was unavailable on PATH; explicit cached toolchain used without global changes |
| ModelGovernance initial WORK and source checkpoint | Created artifact 29973936-d34d-4751-b15a-66c2d0f64a2b, observed baseline checkpoint and frozen PROPOSAL | New concern has no accepted release; source-observed baseline is not owner-approved NOW |
| SDL strict check of final WORK | design-core/0.6, four sources, valid, no warnings | Structure/scenarios only; no executable state-machine predicates or runtime proof |
| Model checkpoint/status/history | Head :00002, clean live/source digest and readable lineage | Target remains mutable, implementation status planned |
| Diagnostic Blueprint | Generated NOW proposal -> TARGET WORK, revision 65947524a753af731e475e044a8d78eec393d359fdef5a0e3fdb865dd0aa36a6 | ConstraintsPass=false because no phase-specific change permissions were granted; no ready assignment/review/acceptance claim |
| Static viewpoints | 91 generated files in /tmp/sdp-t041-views | Generated from validated facts; no XFMD UI trial or runtime simulation |
| Exact source/generated hashes | Four target source files and 107 generated files verified | Binary identity participates in Blueprint revision; a different tool build may change output identity |
| Model operation journals | Four local journals all state done | Local backups/journals retained on disk but not committed; tracked checkpoints/frozen lineage retain model history |

The first authoring check reported noncanonical ordering. The target was formatted
through the SDL tool and rechecked successfully. The frozen observed proposal was
captured before formatting; ModelGovernance and Blueprint compilation validate its
semantics, while the final WORK meets the strict canonical check. No frozen bytes
were changed to hide that sequence.

## Reproduction and exact records

Build from SDPTool and SDL/go respectively using
`/home/warloc/.cache/sdp-workstation/go1.27.1/go/bin/go`, GOTOOLCHAIN=local,
GOCACHE=/home/warloc/.cache/sdp-workstation/gocache and
GOMODCACHE=/home/warloc/.cache/sdp-workstation/gomodcache:

```sh
go build -o /tmp/sdp-t041-sdptool ./cmd/sdptool
go build -o /tmp/sdp-t041-sdl ./cmd/sdl
```

The model-area README supplies exact check/view/Blueprint commands. Source and tool
hashes and generated output manifests are in [manifest.json](evidence-T041/manifest.json).
Actual outputs: [model check](evidence-T041/model-check.json),
[checkpoint](evidence-T041/model-commit.json), [status](evidence-T041/model-status.json),
[history](evidence-T041/model-history.json), [Blueprint result](evidence-T041/blueprint-result.json)
and [check record](evidence-T041/checks.json). The
[task](evidence-T041/blueprint-task.json) deliberately grants no implementation
permissions. Generated preview/view directories are disposable; retained source,
model lineage and task input reproduce their meaning. They are not installation
or publication artifacts.

## Remaining gates

Constrain and review a phase-specific implementation handoff against the selected
model/task revision, then implement PGF1-C1 and subsequent scoped increments.
PG-A25–30 are specified, not executed. Runtime profile/schema support, checked
KanBan moves, shared CLI/MCP operations, complete recovery, native integration and
independent PG review remain open. Maintain WORK through implementation; integrate
and validate it, freeze CANDIDATE and release with actual acceptance and named code
verification. No software publication, repository merge or Ponsse work is implied.
