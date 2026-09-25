# K8 — backlog consolidation

Owner request: 2026-09-25. Base: `29d5828`, phase branch
`sdp/phase-k8-backlog-consolidation`, stacked on K7.

## Milestone K8-M1

Status: delivered and verified; proposed product capabilities remain backlog.

Review every backlog card across registered SDP, SDL and SDUI boards. Consolidate
overlapping primary work without losing scope/history, capture the owner's native
XFMD navigation direction, update affected references/indexes and append ledger
events. Validate lineage, current placement, CardState and documentation links.
Commit the milestone and push the phase; no runtime/GUI implementation or merge.

## Review of all 12 original backlog cards

| Card | Disposition | Reason / dependency |
| --- | --- | --- |
| KB-SDP-002 | Full merge with #016 into KB-SDP-017 | Project-aware facade and discovery/viewer contract form one first navigation delivery. Preserve advanced planning/history work as later milestones. |
| KB-SDP-016 | Full merge with #002 into KB-SDP-017 | Same project context and consumer; owner's tab/tree clarification removes the reason for separate discovery planning. |
| KB-SDP-003 | Keep Idea | Time-axis graph, zoom and historical revisions differ from the initial native status tree; may consume shared card history services later. |
| KB-SDP-004 | Keep Proposal | Evidence semantics must be usable independently of CLI/viewer choices. Dependency for #017's advanced plan/roadmap work. |
| KB-SDP-005 | Keep Ref to KB-SDL-002 | Preserves local viewpoint impact and weak-link evidence exclusion; not a duplicate primary or merge source. |
| KB-SDP-006 | Keep Ref to KB-SDL-001 | Preserves SDP template/profile dependency; requirements narrative can progress independently of visual path routing. |
| KB-SDP-011 | Keep Bug | Existing ID/CI defect requires a historical compatibility decision; merging into future evidence redesign would obscure its acceptance. |
| KB-SDP-014 | Keep Proposal | Versioned KanBan contract/distribution also serves standalone XFMD without SDP. #017 consumes compatibility, not ownership of that contract. |
| KB-SDL-001 | Keep Proposal | Stakeholders, actor kinds and narratives have distinct language/profile acceptance. Depends on #002's basic weak-link semantics, not every visual traversal capability. |
| KB-SDL-002 | Keep Proposal | Weak relationships and ordered visual traversal also serve architecture/design independently of requirements narratives. Coordinate grammar with SDL #001 without forcing one large delivery. |
| KB-SDL-003 | Retarget Ref to KB-SDP-017 | Same stable Ref ID; clarify SDL-owned catalog/model tree/projection services. Rename filename's primary number and record move. |
| KB-SDUI-001 | Retarget Ref to KB-SDP-017 | Same stable Ref ID; capture SDUI subtab boundary, preserving Go/Fyne/SVG direction. Rename filename's primary number and record move. |

Result: 11 backlog cards (7 SDP, 3 SDL, 1 SDUI), including four intentional
cross-project Refs. One queued primary, KB-SDP-017. Two superseded primary sources
remain available with full transfer lineage. No split is necessary. KB-SDP-010
remains active at non-blocking gate-review.

## Ownership and external handoff

The [consolidated card](../../KanBan/active/%23017--Proposal--sdptool-and-project-navigation.md)
is the current producer-side planning entry. XFMD's separate KB-XFMD-014 was read
for protocol/consumer alignment, but that repository and its card are not edited
here. Its old bare file URL for #016 now points to the previous location; the
XFMD follow-up should retarget it to #017 and adopt the owner's subtab/tree
clarification. No message has been sent to another agent.

The launcher accepts a model argument/environment overrides but fixes project ID
`sdui`; therefore the proposal records this actual limitation rather than calling
it completely unconfigurable. Existing metadata filenames remain proposals.

## Verification and completion

See [evidence](Evidence.md). This maintenance milestone is documentation and board
consolidation only; it does not establish valid discovery schemas, implement
sdptool or verify native XFMD behavior.
