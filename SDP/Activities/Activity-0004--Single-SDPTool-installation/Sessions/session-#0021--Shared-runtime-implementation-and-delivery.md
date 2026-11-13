# Shared runtime implementation and delivery

## Roadmap

C0 selected model/contract → C1 shared runtime registration → C2 exact extension
payload → C3 native migration → C4 independent qualification → C5 publication.
Current step: C0/C1; bounded C2/C3 implementation proceeds against the selected contract.

| Field | Value |
| --- | --- |
| Session reference | SESSION-SDP-0021 |
| Status | active |
| Activity | ACTIVITY-SDP-0004 |
| Phase | 05-Implementation |
| Phase contributions | 03-Requirements, 04-Design, 05-Implementation, 06-Verification |
| Plan | PLAN-SDP-0029 |

## Goal and authority

Owner: execute Activity4. Deliver one shared SDPTool engine through gh extension
and offline XFMD consumption. Plan and mandate select full qualified delivery,
without main merges or unrelated project upgrades.

## Concurrent-work preflight

SDP432c292 on sdp/shared-tool-installation contains the planned Activity4 study;
only the persistent project-history lock is unrelated untracked state. gh-sdp
b177e14 is clean. Isolated XFMD140be2b has no tracked modifications; its .deps
link is preserved. Primary XFMD main has unrelated dirty work and is excluded.
Activity1/2/3 are complete; actual service transitions select Activity4 and create
this Session. KB063 is the only card selected by this plan. No identity collision
or observed overlapping work affects the three bounded source areas. Green for
this scope; unknown remote work is not asserted absent.

## Work journal — 2026-11-13

Loaded SDP entry, master, change-analysis, architect, planning, versioning and
release routines. The implementation contract selects exact canonical-engine
republishing, explicit setup, direct native lookup and no runtime bootstrap.
Bounded delegated workers own gh-sdp and the isolated XFMD checkout; main owns
canonical engine/model, integration and release coordination. Fresh independent
review remains required. SDUI preview/check move into a shared library; the native
prototype helper is embedded in the same release executable and materialized
locally by content hash, with no download/update route of its own.

Actual Activity update/session-create/activate receipts are retained by the service
in ../.operations and management history. This journal is a manual work summary.

## Linked records

- [Implementation plan](../05-Implementation/Plan.md)
- [Runtime contract](../04-Design/ImplementationContract.md)
- [KB063](../../../KanBan/active/%23063--Change--Single-SDPTool-installation.md)
