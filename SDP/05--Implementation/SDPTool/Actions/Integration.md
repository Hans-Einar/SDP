# ACT0 integration evidence

Owner authorization: Session0011 T010, 2026-10-18, explicitly includes main merges
for integration and subsequent complete Actions delivery. No release is selected.

Pinned remote main: 04f88ff. Inputs: release 3475367 (includes program PR lineage),
blueprint/actions f64a5ba and PG de4b795 (PGL2-C1 included). Isolated worktree:
/tmp/sdp-session11-integration, branch sdp/session11-integration. Combined code and
history candidate: df29a12. Later changes in this milestone are documentation only.

Resolved shared CLI dispatch, string-option recognition, presentation registry,
README/Contract, current owning cards and model/management/traceability consistency.
Both program run and governance dispatch survive; blueprint handlers/presenters
survive. Published 2.2.0 lineage includes terminal and signal behavior corrections.

[History import](../../../ProjectManagement/History/Session0011/README.md) gives
exact original bytes, mapping and verification. Actions is now KB-SDP-053, PLR1 is
MAINT-SDP-0016. Released subjects retain their IDs. Canonical validation passes with
66 cards, 45 management records, 4 lineage operations and 600 events at integration.

Verification on combined code:

- `go -C SDPTool test ./...`: pass, including ModelGovernance, blueprint assignment,
  controller/MCP, runnable-program and CLI suites. [Log](evidence/ACT0/sdptool-tests.log).
- `go -C SDPTool vet ./...`: pass. [Log](evidence/ACT0/vet.log).
- ProjectManagement validator, Toolkit validator, archived import verification
  and git diff --check: pass.
- Broader SDL/SDUI suites are running separately and recorded when complete.
- Same-context integration review only; no independent agent review or new native
  GUI/provider execution is claimed. Existing PG/blueprint owner gates remain open.

Excluded: all dirty primary/PG/XFMD/other-worktree files; coordinator Session0012
and its overview remain owned by that branch; independent viewpoint work remains
outside this candidate. No reset, stash or deletion in another worktree. Remote
unreported changes are not claimed absent. Input heads were rechecked before merge.

ACT0 is green for the Actions contract/registry implementation on this integrated
baseline under the owner's integration assignment. This does not reserve files
against later PG work: record the boundary in the handoff and refresh before final
integration. No source release archive or signature was changed.
