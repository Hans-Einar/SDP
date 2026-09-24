# R2 — concrete documentation conflicts

Inventory basis: [R1 document map](../R1/Documentation-index.md). Distinguish observed errors from unresolved language/process choices. Evidence comes from existing G-phase implementations/checks; R2 is documentation maintenance.

| ID | Document / observed conflict | Basis for handling | R2-M2 |
| --- | --- | --- | --- |
| D01 | SDL KanBan says code still lives in SystemDesignLanguage | R1 migration map and SDL/go | Corrected to SDL |
| D02 | SDL definition points to isolated parser prototype; docs index labels all studies/ exploratory | Sections 1–11 plus registered profiles (section 12 is open) are design-core 0.5; 13–16 are candidates; Go parser is active | Clarified mixed authority; removed obsolete parser entry point |
| D03 | Data/Channel profiles present 0.3/0.4 and old Python test counts as current | Go design-core 0.5 includes V2–V4; G4/G6/G5 evidence | Separated active rules from dated V2–V4 evidence |
| D04 | Delivery-plan profile says all G statuses remain planned | Architecture model and G evidence; implementation-status remains an explicit source assertion | Labeled V4 snapshot; corrected global current-status claim |
| D05 | Navigation/notation design says G6 is unimplemented | SDL/go/documents, query, broker, reader, G6 evidence and G7 launcher | Labeled original design with implementation/limit links |
| D06 | SDUI runtime contract calls binding planned and leaves rules open that now exist in Go | SDUI/go/runtime/README, SDL action-core/bridge, G3/G4 | Replaced with concise current boundary overview linking package contracts |
| D07 | SDUI layout/composition proposals call Python active and layout/runtime pending | SDUI language, go-layout-contract, architecture, G1–G3 | Kept design background; explicitly identified active profile and differences |
| D08 | prototype-widgets gives runnable instructions for deleted Python tools | G5-M4, Go SVG CLI and sdl-document | Removed obsolete commands; kept dated artifact description |
| D09 | SDUI lacks a docs index; checkpoint is described as current authority | README, active Go contracts, implementation plan; checkpoint is dated history | Delivered document map and clear active entry points |
| D10 | Port inventory uses obsolete Go path; IR study suggests no runtime exists | R1, action-core 0.1 and G4 | Corrected path; distinguished broad IR study from delivered action runtime |

Further work: broader candidate semantics, SDP process/template profile and complete substantive harmonization. R2 does not claim to resolve every statement in over 7000 study lines. Change generated viewpoints through tools when sources change.

## Verified outcome

D01 also covers KB-SDP-001's obsolete baseline. D07 corrects the claim that font units were undecided: Go uses DIP. D10 includes the source-tree study's old 0.1/Python entry point; broad workspace/IR remains a candidate, while bounded data/Channel/action profiles are implemented.

SDUI has a [complete map](../../../SDUI/docs/README.md) covering all 17 subject documents. The [SDL entry point](../../../SDL/docs/README.md) explicitly identifies active core and proposals in the definition. Runtime documentation points to one source per package responsibility; the original proposal remains in Git.

V2–V4 historical counts and original G6 sketches are not rewritten as fresh tests. R2 identifies age/authority and points to active entry points.
