# L1 — English documentation

Owner instruction, 2026-09-24: translate documentation before further backlog
work, and require English in AGENTS.md, including all KanBan documentation.
Branch: `sdp/phase-l1-english-documentation`, from R2 `918fa46`.

| Milestone | Scope | Completion evidence |
| --- | --- | --- |
| L1-M1 | English policy; translate all three KanBan boards and their workflow documents | Meaning/metadata/link review; ledger prefix preservation |
| L1-M2 | Translate remaining maintained documentation, including dated narrative evidence and research | Translation inventory, source/code/link preservation, terminology review |
| L1-M3 | English generated documentation from its generator; regenerate current views | Go tests, generator output manifests, English output checks |

Translation preserves decisions, unresolved questions, names, source citations,
versions and test claims. It must not silently adopt a proposed language rule.
Frozen fixtures and append-only old machine records retain original bytes; new
records and explanatory prose are English. Intentional localized UI strings are
example data, not an alternative documentation language. Generated artifacts are
changed only by rebuilding from translated generator text.

Local translation models were evaluated but their drafts were not adopted; the
translation was written and reviewed directly. No document text was sent to an
external translation service. Original revisions remain in Git. A source/hash inventory
records translated narrative snapshots without rewriting their old machine evidence.

After L1, K5 adds a simple CardState convention and shell listing/installer.
Then R3 resolves the remaining KB-SDP-001/010 deliverables or presents concrete
choices for owner gate review. These later phases do not delay English policy.

L1-M1 and L1-M2 are complete. L1-M3 is next: generator strings and current outputs.

L1-M3 is complete. All three milestones are delivered; see Evidence.md.
