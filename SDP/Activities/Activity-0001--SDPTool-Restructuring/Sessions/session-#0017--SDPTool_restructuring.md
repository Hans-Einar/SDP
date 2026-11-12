# Session 0017 — SDPTool restructuring coordination

| State | Step | Outcome / authority |
| --- | --- | --- |
| completed | S1 Manual Activity bootstrap | Owner-selected Activity registered and activated; future features split out |
| completed | S2 Baseline, target and migration | PLAN-SDP-0026 M0–M3 executed; product candidate 527cc90 |
| completed | S3 Whole-candidate qualification | Exact-candidate broad suite, six native gates and independent product review pass |
| completed | S4 Acceptance | Owner explicitly accepts restructuring and authorizes Activity0002 |

| Field | Value |
| --- | --- |
| Session reference | SESSION-SDP-0017 |
| Activity | ACTIVITY-SDP-0001 |
| Phase contributions | mandate, study, requirements, design, implementation, verification |
| Status | completed |
| Primary card | KB-SDP-062 |
| Plan | PLAN-SDP-0026 |
| Snapshot date | 2026-11-12 |
| Current step | S4 |
| Execution authority | Owner authorizes execution through Activity0001; no publication or acceptance inferred |

## Goal and handoff

Establish the first Activity and carry the complete restructuring plan into its
baseline/design work. [Mandate](../01-Mandate/Mandate.md) owns scope;
[plan](../../../05--Implementation/SDPTool/ArchitectureMigration/Plan.md) owns
milestones. Earlier discussion remains in
[Session0016](../../../Sessions/session-%230016--Activity_design.md).

## T001 — 2026-11-12 manual creation and sequencing

Manual work summary; host IDs and precise work start/end unknown. Owner selects
three successive Activities: restructuring now, Activity support next (manual
creation), browsing third (first tool-created Activity). Used SDP Planning,
Architect context and Traceability. Created the first Activity only, its mandate,
standard README, phase references and this Session. Defined the minimal manual
record contract and appended actual creation/activation with current timestamps.
Revised PLAN-SDP-0026 to exclude all new functionality. KB062 is ready; KB061 moves
to onHold; KB060 remains queued. No SDL/Go changes or generated outputs.

Preflight for record writes: one SDP worktree at a96affd on
sdp/sdl-browsing-planning, existing dirty work is this conversation's planning.
No existing Activity/Session0017 identity or prior Activity event was found. Scoped
record-write verdict green; remote and unrelated running work are not established.
Ledger selection does not claim to stop agents. Product execution requires its
own renewed concurrent-work preflight.

[Bootstrap evidence](../06-Verification/Bootstrap.md) records validation and current
reader limits. Next is M0, not Activity feature implementation. Plan lifecycle is
still planned; Activity activation denotes the selected need, not completed design
or started agent dispatch. No Git commit, push or release.

## T002 — Activity execution and technical qualification

Owner authorizes working through Activity0001. Loaded SDP, Master and Worker,
reusing Architect/Planning/Vertical Refactor/Traceability instructions. Preflight:
one SDP worktree, baseline a96affd; dirty files all originate in this conversation's
planning/manual bootstrap; no existing product edits. Checked sibling XFMD status:
unrelated RPM/process work is excluded. No conflict in selected SDPTool/SDL model
scope. Remote concurrent work remains outside observation. Scoped verdict green.
Created execution branch sdp/sdptool-architecture-migration; preserve the setup as
its first checkpoint. M0 baseline and complete operation inventory underway.

Execution continuation: d0bf6dc preserves bootstrap/authority; a32d0d6 pins M0/M1
inventory, validated NOW/TARGET and Blueprint; c47adb5 migrates Session/CLI/action
policy and MCP framing; 527cc90 completes root routing and typed program execution.
Bounded inventory and Worker assignments used the selected model/contract. A
separate reviewer approved product changes and verified negative architecture
specimens. No product functionality or generated viewpoint was added manually.

Baseline and migration broad suites passed, including governance recovery. Final
root tests, compiled Session CLI/MCP journeys, app/MCP/history race checks, vet,
bootstrap and gh-sdp checks passed. 187 compiled CLI comparisons have zero wire/exit
differences; all 258 baseline library signatures remain. Five native XFMD producer
gates passed. SDUI native initially timed out on baseline and candidate because the
local build lacks its required helper executables. With an isolated helper
installation, the unchanged test passes with both baseline and final producer. The final frozen-source broad suite passed
against 527cc90. Evidence and outstanding owner disposition are in ../06-Verification.

Temporary baseline worktree /tmp/sdp-activity1-baseline-checkout was this assignment's
read-only test checkout and was removed after passing tests and a clean status;
no conflicting work was adopted. No push/merge/release or
owner acceptance is claimed. Plan, inventory, traceability and card are reconciled to the passed candidate;
KB062 enters gate-review. Next: owner acceptance of the restructuring result.
Activity0002 is not created or started.

## T003 — Owner-authorized branch push

Owner requests pushing sdp/sdptool-architecture-migration to GitHub for local
inspection. This explicitly authorizes branch publication after technical delivery
22e2637; it does not grant merge, release or owner acceptance of Activity0001.
Preflight: correct branch, clean worktree, origin points to Hans-Einar/SDP.
Product code and verification are unchanged. Push this Session update with the
delivery and verify the remote branch head. S4 remains awaiting owner acceptance.

## T004 — Owner acceptance and successor selection

Owner explicitly approves Activity0001 and instructs starting Activity0002. Close
PLAN0026, KB062 and Activity0001 using actual acceptance time; preserve all technical
evidence and original event bytes. Next Activity is registered manually for Activity
support with a new Session. No merge or release is authorized by this disposition.
