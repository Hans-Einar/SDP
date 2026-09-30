# SSD1 acceptance and implementation handoff

These are required **future implementation checks**, not passing new-profile tests.
[Evidence](Evidence.md) distinguishes today's baseline probes from these obligations.
The contract is [Contract.md](Contract.md); the runnable-input candidate is
[Frontend.design-set.json](examples/frontend/Frontend.design-set.json).

| ID | Input / operation | Expected outcome and proof boundary |
| --- | --- | --- |
| SS01 | Split Frontend fixture | One System, all baseline declarations/facts plus System declaration/membership; all member ASTs and source spans retained. Compare independently with the original 0.5 model. |
| SS02 | Move ParseStructuralModel and its facts between members | Same identities and semantic graph; changed provenance/revision; no duplicated or lost facts. |
| SS03 | Reorder sources in manifest / filesystem enumeration | Same semantic graph, canonical file map and deterministic diagnostics. Raw manifest revision changes; unlisted file edits do not. |
| SS04 | Duplicate StructuralParser declaration in Features | Duplicate diagnostic identifies both original file spans. Identical declarations are not deduplicated. |
| SS05 | Replace SourceTextPort use with undeclared MissingPort | Error at original use; no lookup in neighbor directories or another registered System. |
| SS06 | Remove entry System; add second System; declare it in member | Each invalid variant rejected with correct spans. A fragment syntax parse alone is not a successful System check. |
| SS07 | Missing/wrong System-container membership; System as allocated-to target | Explicit membership/type diagnostic, no invented runtime container; unit containment validation remains separate. |
| SS08 | Add cross-file unit containment/refinement cycle or second owner | Existing whole-model semantic rejection with participating source locations; ordinary two-way runtime communication is not a generic import-cycle error. |
| SS09 | Mixed 0.5/0.6, action-core/class-core member or unknown version | Reject incompatible profiles before generation; no stripped/replaced headers. Existing standalone profiles keep their own commands. |
| SS10 | Missing file, duplicate member, entry repeated in sources, hard-link alias | Manifest/file diagnostic; no successful partial model. |
| SS11 | Traversal, absolute path, symlink beneath root, directory/device, over-limit bytes/tokens/files | Bounded loader rejection. Nested manifest is not a source, so no recursive include cycle exists. Verify failure does not touch inputs/output. |
| SS12 | Unknown/duplicate JSON keys, wrong types, trailing data, invalid UTF-8 | Strict manifest rejection with useful location/path; no last-key-wins behavior. |
| SS13 | CRLF, UTF-8 invalid character and late EOF in a member | Diagnostics have original path and byte span; byte-column policy explicit. No synthetic concatenated position. |
| SS14 | Global references valid but individual member cannot validate alone | Set check succeeds; per-file canonical format works after global validation. Formatting twice gives the same file map and never moves facts between files. |
| SS15 | Qualified/private/export/import attempt or independently declared second System | Explicit unsupported syntax/boundary error. No claim of cross-System linking, exports or privacy. |
| SS16 | Check → tree → selected Markdown | All use the same candidate revision and System identity; source maps identify original declarations/relations. Consumer assertions derive from known fixture facts, not output snapshots alone. |
| SS17 | Change non-entry member while renderer is paused, with and without cache hit | Recheck rejects stale delivery. Existing lease remains valid for its old immutable revision; no new “current” success. Also test missing member and changed manifest. |
| SS18 | Old tree revision used after member edit | SDPTool selection fails before replacing output; new tree/selection succeeds with new revision. |
| SS19 | Concurrent requests, cancellation and failed open | Existing latest-sequence/cancel/lease behavior remains; no stale request opens merely because a cache entry exists. |
| SS20 | Output aliases manifest, member or their ancestor | Reject before replacing/removing anything. Verify content bytes, not just error strings. |
| SS21 | Relocate identical set to another absolute checkout | Same candidate revision; provenance uses portable relative paths, display root remains caller context. |
| SS22 | Old 0.5 CLI AST/format/check/view and registered navigation | Existing contracts and supported fixtures unchanged; 0.5 does not silently accept System or sets. |
| SS23 | Open a member .design directly versus registered source set | Direct preview reports missing System/unresolved references as appropriate; registered set works. No implicit nearest manifest or changed meaning of direct preview. |
| SS24 | Mismatched registration System/profile, unsupported consumer capability | Reject before output; never emit hardcoded design-core/0.5 provenance for 0.6 inputs. |
| SS25 | Change parser/projector/renderer identity without changing inputs | Output cache is invalidated by tool/schema identity while input revision remains an input digest. |
| SS26 | Reorder scenarios containing steps 2 and 10 across files | Semantic and generated sequence order remains numeric; canonical lexical ordering cannot alter modeled execution order. |

## S2 implementation sequence recommendation

One vertical producer/consumer feature, with milestone boundaries that leave current
single-file tools working:

1. Introduce file-aware pure frontend structures and profile dispatch, preserving
   0.5 APIs. Define checked-model ownership and source-aware diagnostic envelope.
2. Add 0.6 System semantics and explicit set loading with limits, per-file format
   and complete-graph validation. Unit and CLI checks exercise original spans.
3. Migrate viewpoint/document entrypoints, broker and SDPTool to the same checked
   result. Version changed outputs/capabilities and prove non-entry invalidation.
4. Run fixture and current ecosystem standalone regression/consumer checks, then
   independent review and bounded real-model pilot. Keep unsupported consumers
   explicitly single-file until migrated; do not silently downgrade them.

S2 must decide phase names, commits and exact API/schema compatibility from this
contract, and retain the card's phase-branch commitment. It must also make an
explicit disposition of the card's later public/cross-System obligation before
claiming the entire card complete. An implementation milestone can close without
claiming that later capability. No new global installation or release is implied.

## Migration decision for KB-SDP-020

Keep current authoritative models and navigation registrations in place. Use the
Frontend copy for first evidence, then migrate one real System only after its
identity/fact comparison and navigation workflow pass. Existing README links to
external owners remain documentary references. Do not pretend they are compiled
public imports. Keep the Landscape independently checked and clearly scoped.
The experimental MVP1 profile needs additional language work even after these
source-set checks pass. It is not the initial acceptance fixture.

## Completion levels

- **S1 complete:** implementable design, bounded examples, inspected-code evidence,
  explicit scope/migration disposition and truthful Session handoff.
- **Source-set increment complete:** SS01–SS26 pass on the implementation candidate,
  with explicitly scoped unsupported consumers and required independent review.
- **KB-SDL-005 complete:** all its selected obligations delivered or explicitly
  transferred with preserved lineage; public linking may not simply disappear.
- **XFMD integration complete:** evidence from its actual native consumer, owned by
  its workstream. Producer tests alone do not establish this claim.
