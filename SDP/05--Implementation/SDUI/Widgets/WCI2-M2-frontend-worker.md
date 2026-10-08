# WCI2-M2 frontend/export worker handoff

**Status:** bounded lane implemented; scoped verification passed; ready for independent integrated review. No native or whole-M2 acceptance claim.

Role: SDP Worker, reused SDP entrypoint/document workflow and worker routine. Authority: coordinator-selected M2 under reviewed original Panes-and-commands, substantive review 9a9c8710; coordinator approved strict selected-root resolution and host/layout agreed native per-canvas context. Work stayed in the exclusive clone frontend/export files. No commits, branch changes, management/Session changes, original-workspace changes or WCI3/WCI4 code by this worker.

Starting M1 baseline: 403c5402c89532e7bf3130776b27f2300e7612fb. Shared HEAD at final verification: a4f2c22435d91fe07935b8b9d5fcf46fafc6e36b (M1 evidence/handoff commit advanced by coordinator). This is dirty-candidate evidence, identified by the frontend file hashes below, not an unchanged-commit test claim. Concurrent runtime/layout/host/bridge edits were preserved and excluded from this lane's file list.

## Delivered behavior

- Exact unreleased .3 command/button/menu/menuGroup/item/separator/dialog schemas and composition bodies, closed argument types, labels/names, toggle/exclusive/context/effect/key rules, menu placement, nonvisual declarations and visual split-child checks. .2 syntax/output remains frozen.
- Explicit M2 button opt-in only for command/toggle/checked/exclusive/key/context/target/effect. Basic callback/icon/tooltip buttons retain legacy routing; no default injection changes that distinction.
- Expansion records innermost definition-instance scope in the sole closed derived Arguments key `$scope`; original references/spans/use chains remain unchanged. Strict pure ResolveInteractions returns typed Scope/Command/Target/Dialog paths, validates hidden content and resolves relative/entry-absolute references without runtime lexical lookup. Independent reused instances keep their own group scope; absolute references deliberately share.
- ResolveDialogField shares the frontend named-path index and resolves only a dialog-owned input, including anonymous-container projection and excluding nested dialogs. Serialized interaction type/shape errors and malformed scopes reject without source I/O.
- Codegen validates the selected root, retains sdui-go-model/2, UIProfile/UIASTFormat and unchanged AST fields. Executed generated constructors reproduce source normalization, identities, spans/use chains and independent mutable models.
- Structural text/composition describes command behavior, source identity, hidden declarations and initially closed dialogs without live menu/state/acceptance claims. Unknown kinds and unresolved selected references reject before artifacts.
- Public SVG explicitly rejects M2 even with supplied geometry; .3 button icon/tooltip decorations also reject rather than disappear. Native SkipControls requires exact button/menu/dialog adapter inventory; menu children belong to their prepared menu and nonvisual commands have no fabricated native control. InteractionRoot provides strict full-entry resolution for an auxiliary canvas and requires exact subtree pointer membership. Inventory remains scoped; extra/wrong/missing entries reject.
- Standalone prototype explicitly rejects unsupported command/menu/dialog adapters, including hidden declarations. CLI failure preserves the previous output; AST local-profile remains schema/placement-only rather than a connected-readiness claim.

## API and compatibility

Frozen typed signatures and canonical clarification are in WCI2-M2-frontend-API.md. Added parser exports: InteractionIdentity, ResolveInteractions, IsCommandButton, IsAuxiliary and ResolveDialogField. Added svg.Options.InteractionRoot for native canvas context; existing named-field options and default .2 output are unchanged. Go callers using unkeyed svg.Options literals must account for the added exported field.

No Node/Instance fields or tags, source profile, AST envelope or generator version changed. No derived arguments are added to legacy/M1 nodes. Existing SkipControls already required NativeControls in WCI1; M2 extends its supported set to prepared menus/dialogs, without restoring public silent omission. New .3 interaction labels/schema/key constraints follow the reviewed contract. No new capability dimension or dependency.

Normalize returns reusable definitions as templates; it does not pretend an unresolved absolute reference is executable in isolation. Every selected-entry consumer must use the frontend strict resolver before validity/readiness. This is wired directly in this lane's generator, prototype and exports; runtime and preparation integration is owned by James/Noether. Static AST local-profile validates local shape. Main must integrate the proposed canonical paragraph in the API memo; no canonical document was changed here.

## Verification

Environment: Go 1.27.1 linux/amd64. Final commands passed:

- From SDUI/go: `go test -race ./parser ./codegen ./presentation ./svg ./prototype ./cmd/sdui`.
- From SDL/go: `go test ./codegen` (affected generator/provenance regression).
- `git diff --check`.

The first broad scoped run encountered transient missing runtime methods while James implemented them; no temporary stubs or other-lane edits were made. Subsequent normal and final race runs passed.

WCI2-A01: positive/negative schemas, named/reused/alias scope, relative/absolute identity, missing/wrong/ambiguous targets, hidden invalid refs, explicit legacy boundary, original provenance, serialized invalid arguments, actual generated constructors and .2 exact generated-byte/source-AST/text fixtures.

WCI2-A04/A05 frontend portion: canonical presentation ownership, initial exclusive/key conflicts, effect/context restrictions, menu/group/submenu placement and captured target agreement. No native dispatch claim.

WCI2-A06/A07 frontend portion: owned field mapping through anonymous containers, nested-dialog exclusion and invalid field paths. No SDL execution/transaction claim.

WCI2-A09: truthful structural descriptions; supplied/hidden M2 SVG rejection; exact native inventory/full-entry subtree membership; standalone unsupported-adapter preflight; CLI atomic failure. Existing source AST parity, concept text snapshots, legacy generated bytes and reuse provenance tests remain passing. Unknown-future composition fixtures now use `futurePane` because dialog is selected M2 vocabulary.

Remaining: coordinator's full integrated/race/native proof and independent review, including actual menus/dialogs, SDL effects/Accept, capture generations and atomic publication. Platform-reserved/normalized key conflicts, live state, icon preparation and native readiness stay with their designated owners. Frontend lane freezes after this handoff; no optional framework or future family implementation is proposed.

## Exact touched frontend files and hashes

SHA-256 at handoff (report excluded from its own hash):

| File | SHA-256 |
| --- | --- |
| SDUI/docs/go-generation.md | 1dab1017e3237d01f8057705df439989c2d1a8875b4cbcafd3964f3224d74ffd |
| SDUI/docs/profile-0.3.md | b8d452c37eee26bbbf0b6c5eb34c83ccd8be67d8768b0f1072c50e96f1b0ee2b |
| SDUI/go/cmd/sdui/interactions_test.go | 0b222d8731c1f0997fd33f4881e4feb6a68c10a1cd211d058431d70f9fabd8cc |
| SDUI/go/codegen/generate.go | 8b356af69ac312168540791374e3ac9bd3ad3b5d662b8e23614c7934022d4f3a |
| SDUI/go/codegen/interactions_test.go | 861dd3617cc4fb46979ab42abde88c8b08645fcdf80ae02cd336d8514410ad87 |
| SDUI/go/parser/interaction_identity.go | af9a45028b5630da8e19c7b4276b219ea57aead0d0a29817fdbd261003f64ecd |
| SDUI/go/parser/interaction_schema.go | c088ae962a25cb002abcc78b90dc5f59ca8ba371a8793c8fa77aed103fd41439 |
| SDUI/go/parser/interaction_validation.go | 4ecd7827cc57566ad6056b901797d908d65ed1b876dab3f7c8db13720d7e2661 |
| SDUI/go/parser/interactions_test.go | c545442e6c004d6a5d6f67c4e94b06373e8d44dabbe189e57d7205aad19be45a |
| SDUI/go/parser/normalize.go | b4d7fcb5a56e1ff8c382b57f90fd51d540cefb609791828052dca3a0f18b35b6 |
| SDUI/go/parser/panes.go | d5903901dc9d40df91b88704451f1b7541bdbe30b4ebe2b35aa64440256b86c4 |
| SDUI/go/parser/panes_test.go | 975f1d03eab4002c65cac869fac52f3baea292ed2c05a8513d88c8e34ebefd8d |
| SDUI/go/parser/validate.go | 83a1ae196e4844eabf58f99aa065c94b9eefaff6b47da561f4c317d60fa0b134 |
| SDUI/go/presentation/check.go | 04e2db01a9b0bb54aaafa847935d84ade7ac2c50770b4d1920a760b67b6d9a56 |
| SDUI/go/presentation/dump.go | ad8414aaa97b3512f6c4b2f4bab1f11043903afe3f22060afaa47314146aa480 |
| SDUI/go/presentation/interactions.go | 25ec5623f43dc58c685499e230d2d086126ae570380be7e6fca000a53cae196f |
| SDUI/go/presentation/interactions_test.go | 10796e09b18865941ef9cb10a2803373f330142aaef6f7b2771d2b609e49d89d |
| SDUI/go/presentation/markdown.go | 4f5da65a6011b6444de86257956640b92f8133c9a9fe72d4417acdc347958cc9 |
| SDUI/go/presentation/panes_test.go | 8eceb1c0ab24139eb907536c7e06519bae15cfcc30230fa8491a3069d1350d8f |
| SDUI/go/prototype/check.go | c57b4e4c0cf6049abc4fe3f158090f7ae7316e0798cfa8d10f57e36e5c4d2aa4 |
| SDUI/go/prototype/interactions_test.go | 49a83080b32f1a20cbb63f91d970823eadc020c1a773d66bb10f5d25c5ed471e |
| SDUI/go/svg/check.go | 9005d46a487db7afb8f0523f6276329284b4c3b339206cc489a85f8ee3353eb0 |
| SDUI/go/svg/interactions_test.go | 908e18cb4fb64df2684788da997fa04a8585f67f8a6e717d10b45ec50e480573 |
| SDUI/go/svg/panes_test.go | f33764c87e49738185e243b7222f136eda3f29aaa9ebd20ca002d4673ffe4aa0 |
| SDUI/go/svg/render.go | 7fb8e908085bb6c28593a4eb6572b9532177a4041e14bd1e42847cfbd35f5e8a |
| SDUI/grammar/sdui-0.3.ebnf | 5d1bf10b1d1cc5cb787333e1dfd32b63886e6430baa658f6b79a87358f71d26d |
| WCI2-M2-frontend-API.md | 711e5cc4cb932e53181f4e0804c75439e2d83b8da1acbb2535f062e1eaf12ac1 |
