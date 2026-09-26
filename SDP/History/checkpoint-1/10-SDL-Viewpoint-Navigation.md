# Checkpoint #1 — G6 document navigation

Date: 2026-09-22. Owner requests navigable viewpoint documents and on-click generation with XFMD navigator/main panels. **Historical design basis:** G6-M1–M6 are now delivered; [supplement 11](11-Go-Implementation-and-Navigation.md) records implementation. This preserves G6-D1/D2 planning.

[Design contract](../../../SDL/docs/integration/SDL-Navigable-Viewpoints-Design.md) covers directory proposals, URI/adapters, windows, publication/lifetime. [Generated plan](../../../SDUI/design/viewpoints/implementation.md) derives responsibilities, owners, dependencies and scenarios from [SDL](../../../SDUI/design/architecture.design). SDL tools regenerated the [combined report](../../../SDUI/design/viewpoints/viewpoints.md).

M1–M4: static navigation → on-demand generation → XFMD panels → optional service. G6 depends on G4-M1 structural parsing, not domain execution/native generation; G5-M3 consumes its export. D2 adds M5 notation/M6 explicit class profiles. [Levels/notation](../../../SDL/docs/integration/SDL-Viewpoint-Levels-and-Notation.md) covers navigator-only, A0–A5, type catalogues and Mode/State.

Navigators select model facts; full/individual pages share projection, not interpretation of generated Markdown. SDL generates documents; XFMD owns windows/panels/routing. At this design milestone export format remained unchanged. Fyne remains the interactive SDUI host.

Readers are locally configured adapter IDs, not arbitrary Markdown commands. Capture target windows on click. At this stage --active/URI forms were proposals. File readers receive small immutable image/document packages in private runtime directories; pure RAM transfer is later work. No daemon started here.

SDL 0.5 models plans, responsibilities and typed messages without grammar changes. SelectedViewOpened models the path; InvalidViewSelectionRejected/ViewProjectionFailed model pre-open rejection. Models do not execute protocols or validate every field/lease/resource limit.

Verified: 121 tests (61 SDL, 24 tool, 36 SDUI), 463 declarations, 1424 facts, 116 responsibilities, 24 planned milestones. 170 SVGs with identical re-export/source checks; phase/opening sequences visually inspected. VP07 reports 18 missing allocations, including three shared DocumentBrowsing responsibilities; no complete deployment claim. [Machine report](../../../SDL/tools/verification.json).

Historical reproduction command requires the original revision:

```sh
python3 SystemDesignLanguage/tools/verify_design.py --phase G6-D2 --renderer /home/warloc/git/mermaid-rs-renderer/target/debug/mmdr
```

The [branch stack](../../Development-Branch-Stack.md) records G6-D1/D2 above V4 as design delivery, not implementation-phase completion.
