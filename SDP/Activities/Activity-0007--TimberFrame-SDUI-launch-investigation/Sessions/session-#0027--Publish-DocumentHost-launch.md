# Publish DocumentHost launch

| Field | Value |
| --- | --- |
| Session reference | SESSION-SDP-0027 |
| Status | completed |
| Activity | ACTIVITY-SDP-0007 |
| Phase | 06-Verification |
| Phase contributions | 06-Verification |
| Plan | MAINT-SDP-0023 |

## Goal

Publish signed SDPTool and gh-sdp2.4.0 and verify public upgrade in isolation

## Linked records

- plan: [MAINT-SDP-0023](../../../Maintenance/R240/Plan.md)

## T001 — Release authorization and preflight

Owner requests publication of the completed DocumentHost launcher. MAINT0023 selects
2.4.0 as an additive minor release. Plan preflight records clean historical worktrees,
shared write ownership and preserved ledger lock. Loaded SDP, Release, Versioning
and Planning; reused shared document and lifecycle routines. Product candidate735c865
and model0.1.0 are already verified; only release metadata/default/predecessor changes
are selected. One extension worker and one consolidated review conserve quota.

| State | Step | Outcome |
| --- | --- | --- |
| completed | R240-M1 | Signed candidate, exact CI and independent review pass |
| completed | R240-M2 | Both public releases and isolated upgrade verified |

Next: exact signed package, tests/CI, review, authorized publication and closeout.

## T001 qualification checkpoint

Signed clean candidate52888aa and same-engine extensionf7763ac are prepared.
All ten supported predecessor rehearsals, source archive, focused profile/install,
bootstrap,90Toolkit and actual signed-child checks pass. Older superseded CI run
was canceled to avoid duplicate testing. Exact final CI44287504558 and final native
package observation are pending. Same review context will cover publication receipts.

## T001 publication — current work summary

The authorized SDPTool/gh-sdp2.4.0 is published and all downloaded assets verified.
Exact CI and consolidated review pass. All ten signed predecessor upgrades and
actual isolated public extension upgrade/fresh install/default Framework checks
pass. The signed native package opens TimberFrame and preserves/reloads its view.
No live extension, real consumer Framework or main merge was performed. Release,
Versioning, Planning and SDP routines governed this delivery; Reviewer was loaded
by the independent reviewer. Both roadmap steps are completed. Activity0007's
completed implementation remains intact and its release contribution is complete.
Next: owner may run `gh extension upgrade sdp`; no selected release work remains.
This is a work summary, not a captured final response.

Closeout changes only records/evidence. Toolkit and management validators pass;
the closeout commit skips a duplicate full CI run. Exact published candidate
52888aa retains successful CI44287504558 and all production package evidence.
