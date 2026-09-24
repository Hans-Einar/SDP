# SDUI — worktrees, extraction and language direction

**Historical mapping:** inspections/advice date from September 19–20. Go/Fyne and
independent SVG were selected on September 21; see [architecture](target-architecture.md).
Documentation translation does not recheck source/worktree status. FOX/C ABI are
not active deliverables.

Mapped 2026-09-19, corrected 2026-09-20. This source inventory preserves the owner's
input after chat/worktree reorganization. It implements no code. The owner rejected
legacy SDUI/BoxUI compatibility. Versions/library/API names below are proposals unless
stated otherwise.

## 1. Owner input at the time

- Independent SDUI prototype language, initially with XFMD as host.
- Separate BoxUI from Mermaid; proposed sdui-fox-renderer and a layering review.
- Reuse the SDUI parser, never copy it; libsdui was a proposed library name.
- Defined SDL/SDUI runtime ABI and host hooks for native FOX widgets, callbacks
  and property updates.
- Free text as Markdown, including Mermaid within frames: desired coverage,
  not evidence of implemented support.
- General frames `[]`, widget contents `<>`, layout `{}`; BoxUI as a selected
  frame implementation.
- `*box` selects an alternative implementation; `*b` is an unambiguous abbreviation,
  not multiplication/repetition.

Owner sketch (not valid SDUI 0.1; localized example text retained):

```text
[heading='boxui heading', [<"# gyldig markdown"; button('OK') > ]*box, [ ]*b ; [ ]*b ]*b
```

## 2. Git snapshot

New sessions must recheck status before changes. This inspection used worktree lists,
status/log/merge-base/diffs and source code. No branches were switched/deleted.
SDP-vNow was clean on sdp-vNow 9ad4324 before this note. SDL/MVP1 was preserved in
e051fe2, SDUI in 3becdb1.

### Mermaid worktrees

| Absolute directory | Branch / HEAD | Content and observed status |
| --- | --- | --- |
| /home/warloc/git/mermaid-rs-renderer | master, afab5e9 | Clean. Mermaid extensions, measurement/routing and diagrams; no BoxUI module. |
| /home/warloc/git/mermaid-rs-renderer-boxui | phase/boxui-045-core, 4bfd179 | Clean. Shared mandate/design and early core: model, JSON parser, measured layout, SVG, six documented core tests. |
| /home/warloc/git/mermaid-rs-renderer-boxui-implementation | phase/boxui-046-implementation, 61a85b6 | Clean. R1–R4 pinned by XFMD: strict source/prepare parsing, typed model, snapshots, geometry, SVG, control maps, child-SVG composition and limits. No FOX/SDL runtime. |
| /tmp/xfmd-mrr-measurements | detached, 3726ccb | 23 staged files, 784 added lines: premeasured text and cooperative time budgets/checkpoints in diagram layout. Separate from BoxUI; preserve. |

Both BoxUI branches share 7076cac; core has one exclusive commit, implementation
four. **Implementation does not descend from core.** Do not mix their parsers/layout.
Producer handoff and XFMD integration identify implementation as the integration source.
feature/boxui-extension points to design base 7076cac without a registered worktree.

Measurement work later appears as 38b6018 and extensions in 3aad3fb. The staged /tmp
patch has a different stable patch ID from 38b6018; complete redundancy is unproven.
No cleanup was performed/recommended without a separate comparison. `gh pr list
--state open` returned none for the renderer; that does not prove integration.
Local master had no BoxUI.

### XFMD worktrees

| Absolute directory | Branch / HEAD | Content and observed status |
| --- | --- | --- |
| /home/warloc/git/xfmd | main, c245fd9 | Clean, no BoxUI integration; Mermaid pin 589517a. |
| /home/warloc/git/xfmd-boxui | sprint/003/phase/048-boxui-verification, a7495b8 | Clean. Native controls, local synthetic activity, session/ledger, static PDF, BoxUI fence. Full Mermaid pin 61a85b6; PR #37 open/unmerged. |
| /tmp/xfmd-p31-lifecycle | phase/p31-mermaid-svg, 35e70b7 | Clean. Older SVG/Cairo integration and Pango cleanup; already ancestor of main/BoxUI, no new local changes. |

[XFMD PR #37](https://github.com/Hans-Einar/xfmd/pull/37).
No reversal from Mermaid master is needed. The then-proposed extraction would reuse
implementation code and redirect XFMD to an independent producer. Historical
branches/worktrees remain provenance; the later Go decision supersedes that delivery.

## 3. Concrete extraction boundaries

Paths below are relative to the named worktree, not SDP. These are historical
recommendations, superseded where the current Go architecture differs.

| Owner | Sources | Proposed disposition at inspection |
| --- | --- | --- |
| Mermaid implementation | src/boxui/{model,parse,validate,frame,layout,svg,embedded}.rs, mod.rs | Reuse suitable parts in a FOX-independent SDUI core; no legacy JSON/wire compatibility path. |
| Mermaid implementation | tests/boxui.rs, examples/boxui_*, SDP/06--Container-Design/contracts, fixtures/evidence | Preserve provenance/license; port relevant tests, retain old contracts as historical comparisons. |
| Mermaid implementation | src/lib.rs, Cargo manifest/lock | BoxUI added via pub mod boxui and roxmltree; an independent crate needs actual dependencies, not the entire Mermaid core. |
| XFMD BoxUI | src/interpreter/mermaid/rust/src/boxui.rs | Thin mermaid_rs_renderer::boxui parse adapter; redirect producer dependency. |
| XFMD BoxUI | src/renderer/diagram/rust/src/boxui.rs | Prepare adapter sharing measurement/cancellation; separate UI from Mermaid adapter. |
| XFMD BoxUI | src/application/composition/mermaid/src/boxui.rs | C ABI parse/prepare/free, buffer release and panic boundary; not owned by renderer fork. Preserve function during transition. |
| XFMD BoxUI | src/contracts/boxui, src/interpreter/boxui, src/renderer/boxui | Model, fence extraction, host adapters and placement; migrate with producer. |
| XFMD BoxUI | src/application/boxui/BoxUiSession.*, BoxUiPreparation.* | Identity, snapshots, ledger and preparation; evaluate runtime extraction separately, do not copy application coordination into layout. |
| XFMD BoxUI | src/application/boxui/SyntheticActivity.* | Local simulation/test adapter, not SDL runtime. |
| XFMD BoxUI | src/application/adapters/FoxBoxUiOverlay.*, FoxBoxUiInput.* | FOX reuse candidates coupled to FoxRenderHost/BoxUiSession; require an explicit host boundary before independence. |

Rust modules use local super modules and standard/support libraries. Parent preparation
does not parse Mermaid; the host supplies prepared child SVGs. Examples/integration
tests may use Mermaid as a test consumer. Ordinary Mermaid/treemap retain ownership
and regressions.

## 4. Owner correction, 2026-09-20

The proposed first delivery preserving BoxUI 0.1/BX-HOST is withdrawn. Reuse suitable
code directly in the new model, port examples and remove replaced execution paths.
No Mermaid master change is needed. Historical worktrees are not active product
dependencies. [Plan](implementation-plan.md) owns phases/milestones/cleanup;
[architecture](target-architecture.md) owns ports/libraries;
[layout proposal](layout-language-proposal.md) develops frames, widgets, Markdown,
modifiers and symbols. At the date of this mapping it was planned, not parser support.
This note freezes no old ABI/wire format.

## 5. SDP pilot

The numbered structure was piloted in both Mermaid BoxUI worktrees' SDP directories:
01--Mandate/01-01--Mandate.md, 00--Project/00-02--SDP-Structure.md and sdp-project.json.
The former Mermaid-extension mandate is superseded by independent SDUI and creates
no new constraint.

## 8. Verification boundary

The preceding review ran SDUI 15/15, design-core 31/31 and MVP1 inventory checks;
these do not prove the new language/runtime/rendering. This mapping ran on 61a85b6:

```sh
cargo test --offline --locked --no-default-features --test boxui
```

18/18 passed, with nine dead-code warnings from ordinary routing modules. Tests cover
deterministic frames, control geometry, child composition, resources and cancellation.
They do not prove an independent crate exists. The worktree stayed clean. Any extraction
must run its own regressions and identify consumer pins. No new GUI/PDF test or installation
was part of this inspection.
