# S1 — language source organization

Owner instruction: 2026-09-25. Updated direction: shared root SDP process area;
the initial separate-tree destination is superseded. Owner-selected three-system boundaries
are recorded in the [system study](../../03--Architecture/System-Boundaries-study.md).
PM1 relocated this plan and consolidated the boards. Model/source migration has not run. Primary: [KB-SDL-004](../../KanBan/superseded/%23004--SDL--Change--Language-source-organization.md).
Base: 1a3f5e8; branch: sdp/phase-s1-language-source-organization.

| Milestone | Acceptance | State |
| --- | --- | --- |
| S1-M0 | Record owner-selected system boundaries, shared root process ownership and the XFMD bootstrap handoff | Delivered, process/design decision only |
| S1-M1 | Explicit Go source-set input profile, per-file parsing/provenance, combined validation, CLI/projection/broker integration, bounded failure and revision tests | Unfinished draft; resume using selected three-system boundaries |
| S1-M2 | Adopt shared root phase homes for both language areas, split the mixed model without semantic loss, migrate active consumers and preserve historical evidence, verify navigation and links | Planned |

The original source-tree study informs this bounded tool input profile. No new
language namespaces, imports, includes, System keyword, recursive discovery or
repository extraction is selected. SDL/SDUI implementation remains in their Go
modules. The shared historical G plan remains one authority.

Evidence will record actual checks at each milestone. Keep the old source byte
sequence available as a dated migration baseline, not a parallel maintained model.

## S1-M0 evidence

The owner selected SDL, SDUI and SDPTool as separate systems in root SDP, and
XFMD as a separate system with its own process area. The
[verification record](../../Verification/VER-S1-M0.md) covers this decision
and the independently committed XFMD process bootstrap. Go input-loader drafts
remain uncommitted and untested; S1-M1/M2 are not delivered.

## Scrum-0001 successor planning

The mixed primary was fully split into [KB-SDL-005](../../KanBan/backlog/%23005--SDL--Change--System-and-source-sets.md)
and [KB-SDP-020](../../KanBan/backlog/%23020--Change--Shared-design-source-organization.md).
Use those bounded contracts for subsequent phases. The rows above retain S1's
historical planning state; they do not authorize continuing the abandoned separate
language process layout. PM1 delivered only board/process housekeeping.
