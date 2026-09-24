# R1 — organize the SDP repository

The owner activated KB-SDP-001 on 2026-09-24 and requested immediate housekeeping. Branch: `sdp/phase-r1-repository-organization`, from K3 `431e47e`.

## Goal and boundaries

One physical home for templates, project records and each language. Preserve current installation/language semantics and historical evidence. KB-SDP-010 owns editorial consolidation/conflict resolution. Directory migration creates no GitHub repositories, submodules or sdptool implementation.

## Target structure

| Area | Responsibility |
| --- | --- |
| Template/ | Sole source for neutral project-root/sdp-root templates |
| Toolkit/ | Installer, schemas, verification, managed payload/skills; contract documents in docs/ |
| SDP/ | Development of SDP itself: KanBan, Traceability, Sprints, review/verification, studies and maintenance |
| SDL/ | Existing SDL code/grammar, documentation, scripts and its own SDP/KanBan |
| SDUI/ | Existing SDUI code/documentation and its own SDP/KanBan |
| docs/ | Short shared entry point and SDP process documents |
| SDP/History/ | Labeled checkpoint and earlier bootstrap/process drafts |
| experiments/ | Named candidate exercises, not production profiles |
| examples/ | Toolkit contract examples, still manifest-managed |

Checkpoint #1 is a shared dated discussion/implementation snapshot, not solely an SDL specification. Gather it under SDP/History/checkpoint-1 and link from SDL/SDUI. Active language profiles belong to SDL. Preserve historical source-index/hash manifests; the migration map records new addresses without falsifying old evidence.

Numbered root files are empty templates matching installed templates except final newlines. Deduplicate them against Template. Preserve existing template phase/destination names during this migration: renumbering/per-phase studies must be designed with KB-SDL-001's requirements profile. Do not distribute parser-invalid .design templates. This remains in KB-SDP-001 after physical organization.

## Milestones

| ID | Delivery | Verification | Status |
| --- | --- | --- | --- |
| R1-M1 | Inventory, ownership, migration map; activate 001, register 010 | All sources exist; duplicate templates compared; baseline recorded | Delivered |
| R1-M2 | Template and project records; archive old bootstrap | Toolkit regressions, conformance, baseline validator, traceability and links | Delivered |
| R1-M3 | Gather SDL; distribute docs; indexes and working entry points | Go race tests, CLI/launcher, unchanged generated manifests, links | Delivered |

## Migration rules

The [migration map](Migration-map.json) records each source, destination, milestone and old hash. The [document inventory](Documents.json) classifies every Markdown file formerly in docs; classification guides ownership/reading, not new language authority. Documentation-index.md provides the readable map.

Repair current Markdown links and executable paths. Preserve ledger lines, dated hash manifests, port fixtures and generated viewpoints; generated content must still come from generators. Retain Go module identity during directory migration; choose final module/repository identity during separate extraction. Move installer source addresses while preserving destinations/behavior. Conformance expectations change only for declared source addresses.

Known baseline: Toolkit validator reports nine ID-format mismatches in older Issue #5 records. Do not hide them, renumber arbitrarily or confuse them with new failures. PowerShell is unavailable locally; CI runs native Windows installer tests. Local structural checks are not native PowerShell evidence.

Commit each completed milestone; push the delivered phase head without merging. KB-SDP-001 remains active while phase/template profiles and final repository extraction need decisions; physical housekeeping does not complete language/process design.
