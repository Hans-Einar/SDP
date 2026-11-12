# SDPTool restructuring mandate

Activity: ACTIVITY-SDP-0001. Owner authority: project owner instruction,
2026-11-12, recorded in Session0016 T6 and Session0017 T001.

## Need

SDPTool must retain a coherent architecture across agent handoffs. Duplicate
operation definitions and CLI-mediated internal calls allow temporary solutions
to become permanent architecture. Analysis must survive into implementation as
explicit contracts and executable checks, with all required cleanup completed.

## Desired outcome

Complete PLAN-SDP-0026 for existing SDPTool behavior: inventoried operations and
consumers, verified NOW and TARGET, typed application boundaries, clear state
ownership, tested compatibility and recovery, and removal of every migration-only
path. The final integrated candidate must have independent review and the actual
owner acceptance disposition. Unfinished cleanup cannot be renamed completion.

## Scope and constraints

Restructure existing functionality only. No Activity commands, new public actions,
new browsing behavior or new Blueprint product capabilities are delivered here.
Activity handling is the next separate, manually created Activity. After that
support is verified, the browsing change is the first Activity created by SDPTool.
Verification checks and authored Blueprint companion documents support this
migration; they are not new shipped product capabilities.

Use the existing shared management ledger. Keep system requirements and SDL in
one canonical location. Preserve history and external contracts; any behavior
change needs an explicit owner decision. No automatic migration, agent stopping,
publishing or release authority is inferred from an Activity state change.

## Route and completion

[PLAN-SDP-0026](../../../05--Implementation/SDPTool/ArchitectureMigration/Plan.md)
is the authoritative plan. KB-SDP-062 is the selected card. Session0017 records
manual coordination and subsequent execution. The plan owns current milestone
status; Activity activation alone does not imply implementation or acceptance.

Completion requires the plan's complete finite inventory, verified target rules,
all required behavior/integration evidence, zero migration-only paths and owner
acceptance. The removed Activity feature milestone belongs to the next Activity
and is explicitly transferred, not claimed as delivered.
