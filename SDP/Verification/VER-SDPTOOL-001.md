# VER-SDPTOOL-001 — initial feature design

Subject: DES-SDPTOOL-001 / SdpTool. Milestone: TF1-M1. Date: 2026-09-25.
Scope: local process structure, design-model and document consistency only.
This is not executable, native GUI or end-to-end preview acceptance evidence.

## Artifacts and authority

The owner's five-phase decision is applied in SDP/01--Mandate through
05--Implementation. [Requirements](../02--Requirements/SDPTool.md),
[architecture/model](../03--Architecture/SDPTool.md),
[detailed design](../04--Design/SDPTool.md) and
[plan](../05--Implementation/SDPTool.md) share one SDP project, KanBan and ledger.
The plan moved from Toolkit/SDPTool; no nested SDP or parallel Features tree remains.
Distributed templates and discovery metadata are unchanged.

## Go model check and generated projections — passed

Built the SDL CLI from the current checkout using Go 1.27.1 into a temporary
verification directory. Used that binary, rather than the older installed binary:

```sh
sdl check SDP/03--Architecture/SDPTool.design
sdl ast SDP/03--Architecture/SDPTool.design
sdl viewpoints SDP/03--Architecture/SDPTool.design --format navigator --project sdp --output TEMP/navigation
sdl viewpoints SDP/03--Architecture/SDPTool.design --format static --viewpoint VP01,VP02,VP06 --monolithic --project sdp --output TEMP/review
```

Canonical formatting was produced by SDL's formatter. Check returned valid=true
with no diagnostics; AST generation succeeded. Source SHA-256:
`7ed598ad171ea5e883082d2f593b8472169514e056fbd7969b79467714e89de2`.
Navigator export contained 15 files and no diagram files. The selected static
review generated 18 diagram descriptions/Mermaid sources across VP01/VP02/VP06.
The generated combined report contains SdpTool. Outputs are temporary verification
artifacts, not a second maintained design or manually written viewpoints.

The model has eight owned Functionality responsibilities and planned Activities.
No implementation-status is implemented/verified. Interfaces describe proposed
boundaries; concrete preview/configuration contracts are still to be designed.

## Board, traceability and links

Checks passed: 22 cards/102 events across the SDP/SDL/SDUI boards replay with
correct placement and preserved prefixes from fbd434a. Local documentation checks
validate 2,237 file links and 130 fragments and preserve 105 frozen records/ledger
prefixes and 574 generated outputs. New Traceability paths and event envelope
validate; its prior ledger bytes are retained. Relations schema checks report the
same nine pre-existing ID errors as fbd434a and no new errors. KB-SDP-011 remains
their owner; this is not a claim that full Toolkit conformance passes.

XFMD's two producer-card links were updated to the active #017 address; no product
code changed. Its board replays 15 cards/51 events with preserved prefix from
646f595. Existing blueprint/link and symbol checks pass (37 blueprints, 70
requirements and 279 callees). Both repository diffs pass whitespace checks.

## Remaining work

The active card is ready for the saved-file preview operation contract, followed
by its Go facade implementation. General discovery, unsaved-buffer preview and
native XFMD integration retain separate milestones/owners. Full runtime tests and
GUI verification will be required for those deliveries, not inferred from this
parser/projection check.
