# WCI4 frontend/export worker handoff

Status: bounded WCI4-M1 lane complete and frozen for integration/independent review.
No commit, branch change, module edit, product publication or native acceptance claim.

## Authority and candidate

- SDP entrypoint and Worker/document workflow reused; current selected original
  Providers-and-packaging contract and final frontend/provider/layout/host API
  memos plus WCI4-assignment-audit read. Session/PM/canonical updates belong to main.
- Clone: /tmp/sdp-sdui-widgets, sdui/widgets-wci4.
- Phase HEAD: 90a94b5daa4b36635be102858b2ce4a8012150c6, with the exact frontend
  file hashes below; other lanes' concurrent changes are excluded from this inventory.
- Original canonical Providers-and-packaging SHA-256: 91285ab840cc7ddab2c571a513899c969ee109781d7ff4ff276cfae254b7a4d8.
- Environment: go1.27.1 linux/amd64. Tests are component/source-consumer evidence,
  not OS-native input, package installation or independent approval evidence.

## Delivered behavior and actual API

- `parser.PreviewOptions(*Instance) (PreviewPolicy,error)` and
  `PreviewPolicy{Explicit bool; Description,Fallback string}` are on disk and were
  announced directly to James, Gibbs, Noether and main. No provider/runtime imports.
- Explicit SVG validates source/description/fallback with description OR fallback
  presence as opt-in. Both are then required; source member resource and declared
  module alias are checked. Caption, @resource alone and .3 profile do not opt in.
- New Markdown call requires exactly typed text/description/fallback. Source AST
  remains widget/Arguments; normalization uses markdown/empty Widget/Text plus all
  three arguments, preserving spans, reuse identity and independent maps. Bare
  Markdown and old SVG stay unchanged. Description/text limits and UTF-8 validated;
  unsupported HTML/images remain provider decisions. Strict selected-root checks
  include hidden/closed descendants; no schema I/O or live outcome fields.
- `svg.PreparedPreviewRenderer` embeds existing ContentRenderer and the pure
  `CheckPreview(*parser.Instance) error` contract. Native explicit Markdown checks
  the frozen per-path outcome through this interface; no preparation or rendering
  during SVG preflight. There is no Markdown inventory entry.
- Native explicit SVG omission requires SkipControls plus exactly path -> svg for
  an actually prepared control. Legacy SVG still paints its placeholder; extra or
  invented legacy-SVG inventory entries reject. Omission does not skip descendants.
- Public SVG rejects every explicit preview with unsupported-resource-export and
  original path/span, even hidden/label policy/supplied Content/geometry. No supplied
  resource public exporter was added. Existing output survives rejected CLI export.
- Structural Dump/Markdown/Combined describe reference or quoted source text,
  caption, policy, description, original span and hidden intent, explicitly saying
  resources not supplied; source intent only. No rendered/native readiness claim.
- Standalone prototype/check reports unsupported-preview with Host svg-resource/1
  or markdown/1 and source/reuse provenance. No legacy-launcher migration.
- Generic codegen needed no product changes: real emitted all-family constructors
  compile/run, Document and Root round-trip including exact source spans and reuse,
  and strict selected-root resolution succeeds. Version03/UIProfile/UIASTFormat
  remain sdui-go-model/2, sdui/0.3, sdui-ast/0.3.
- Added source-only consumer tests, not duplicate implementation: sdui-preview
  delegates Combined; SDPTool UIPreview delegates Markdown. Both truthfully retain
  .3 profile; helper protocol remains sdptool/0.2. Source snapshots/hashes, stale
  revisions, invalid-schema discovery and artifact preservation are exercised.
  No SDPTool or sdui-preview production source change was necessary.
- Updated only assigned English profile/generation docs and .3 grammar comment.

## Compatibility and review boundaries

Public API additions are PreviewPolicy, PreviewOptions and PreparedPreviewRenderer.
The native renderer's supported set gains opt-in SVG (exact NativeControls inventory)
and explicit Markdown (matching prepared Content). SkipControls alone still grants
no omission. Legacy SVG cannot be silently dropped. Existing ContentRenderer API is
unchanged; the added interface is required only by opted-in native Markdown.

Frozen .2 schemas, empty Instance.Profile omission, legacy .3 forms and generator
versions are preserved. Tests compare AST/dump/Markdown/SVG and generated Go hashes
for legacy SVG/bare Markdown in both profiles against a separately built git archive
of 90a94b5, not snapshots of the new implementation. Existing input frozen hashes,
constructor provenance and prior family tests remain intact. No material contract
change or competing provider/host stub was introduced.

## Verification and acceptance trace

PASS from SDUI/go:

```
go test -race ./parser ./codegen ./presentation ./svg ./prototype ./cmd/sdui ./cmd/sdui-preview
go test -race ./cmd/sdui-preview
```

The second command checks the final added source-snapshot/hash/stale-revision proof.
PASS from SDPTool:

```
go test -mod=readonly -race -run TestSDUI -count=1 .
```

PASS: scoped git diff --check; Go files formatted. Generated-constructor tests compile
and execute separate generated packages. The earlier concurrent preparation/provider
Check-method compile gap resolved when the provider implementation landed; final
commands pass without frontend stubs. Development test errors were corrected before
these passing checks.

| Acceptance | Frontend proof and remaining ownership |
| --- | --- |
| WCI4-A01 | Positive/negative typed schema, limits/UTF-8, no-I/O refs, normalized Markdown, independent reuse, forbidden interactive roles, actual generated constructors and legacy bytes. Complete for this lane. |
| WCI4-A03 | Schema/identity shape cannot be converted to fallback; public rejection preserves old artifact. Actual provider identity/digest/content fallback and bundle publication proofs belong to provider/host. |
| WCI4-A04 | Pure CheckPreview preflight probe checks both equal-text paths with different policies including hidden nodes, rejects missing/mismatched outcomes, and invokes no renderer. Immutable real outcomes and zero StateGate/ticket rerenders belong to provider/host. |
| WCI4-A05 | Shared all-family source fixture covers frame/grouping, legacy/decorated/toggle/shared-command buttons, tree/list/scroll, tabs/page/split, menu/group/item/separator/context, modal/nonmodal dialog, four scalars, basic/single/multiline input, legacy/explicit SVG and bare/explicit Markdown. Selected-root inventory, emitted constructors, public support/rejection table, native inventory, structural output and discovery pass. |
| WCI4-A06 | Direct source-consumer commands/protocol, .3 metadata, .2 legacy shape, revision/snapshot provenance and stale rejection pass. Actual staged packaged binaries/build roots/licenses and connected route proof remain main-owned. |
| WCI4-A02/A07/A08 | No new native, IME or full-candidate claim. Full race/affected suites, exact-candidate package/native proof and independent Mendel review remain integration obligations. |

No unresolved frontend implementation issue is known. Freeze product edits now;
reviewer findings or coordinator instruction may reopen this bounded lane. No other
lane files, canonical/Session/PM records, modules or previous evidence were edited.

## Exact changed paths and SHA-256 freeze

The manifest lists all 23 frontend-owned changed/new files (excluding this report).
Existing unrelated root reports and other lanes' changes are not part of this lane.

```text
2f4fba0b26f86f08b98196c16f0a48d3c15f5d4c850f8ed5744987ee0bcdfa8c  SDPTool/sdui_preview_families_test.go
4f272dac5604ae2ead49f1ed23572cb6fe7090484749948688dd283ed27326f8  SDUI/docs/go-generation.md
0461819670d70ad030f1a5c28e17d264c54556bdfb3f8e243b537aa1a12936e1  SDUI/docs/profile-0.3.md
b0fbd293e427960ac59434c7f33de326c0c0cd3c4518284958bb535bdf2b349f  SDUI/go/cmd/sdui-preview/preview_test.go
2e3f88eb58c77540d0efdbe9b872c060594f55b76225b10bfc457b836b4812f1  SDUI/go/cmd/sdui/preview_test.go
8a20e79f72a2d95b6551b467e9577e50b6046916cc3124ed1ffaa499a4fc745f  SDUI/go/codegen/preview_test.go
77c5b8ea2b765d4bfb0adfb6777d421b24a59623e2baa044485b064a727a75a8  SDUI/go/parser/all_families_test.go
ba85a06f73993a509d0ede7da724581a93e6e165edc250066b80e6f0762de368  SDUI/go/parser/interaction_identity.go
2f81e6d8455a9e879e0650c5f8670b35b3c0159f7cdfd3df3e123111885ad311  SDUI/go/parser/normalize.go
2e543495f82f58016ec3cf2b22b96a1499d273db09caae2b06bb45d2c75fe533  SDUI/go/parser/preview.go
3b50d3326737ed48ae0e1fb29c7331d658f333289ee884d75c0ee5ad4a8e77c3  SDUI/go/parser/preview_test.go
1681e0610e41df77917740052b81b5b1630a8c36baad9388e1ab33fc2b1c87e0  SDUI/go/parser/testdata/wci4-all-families.sdui
034f9d137ad7a9458880ef78c20f9347312ba0abf119e6b74e59ff76ca97511e  SDUI/go/parser/validate.go
b9f22e46c3c82bedf00bfea64bd09e5abc938f1862c46aa9451d2197af908300  SDUI/go/presentation/dump.go
7fce86b10f1267d037a129a432a3892075762a1e01c4c608ca70e490ea1957d0  SDUI/go/presentation/markdown.go
5ac0edc7606c7f69aea9d22b7775d4aa65dfb7736c1a72928221326090e7393d  SDUI/go/presentation/preview.go
6b68c29c2d1c4c291aa801393ae8c23ee0f6409b9da38dc6464fa46eb094a1f4  SDUI/go/presentation/preview_test.go
24b8413b5f345cb8350300c369eb297f94389c07b7aa0c210069f982b8b1e274  SDUI/go/prototype/check.go
f9541f8b0691e30ccaf1e5f042592f7cc8ff29126b05017a04538e48fe93e8ca  SDUI/go/prototype/preview_test.go
cb3ff0642d1a426fe972fae70f2d02770ef41e56f933c9be9a992ea35757b462  SDUI/go/svg/check.go
2986b1aff7aacc93d87ea7094b6ab18bad1b29d922876c904bd524d5c21a9acd  SDUI/go/svg/preview_test.go
1a58cce07ce727fc9dc263fb5af51c0e755bdf64a640bace9269a00ffaa14e68  SDUI/go/svg/render.go
b74111e212726fac4dd6885bdbc1140f923f502cedbeef6c65bc3a5f0ea04a0d  SDUI/grammar/sdui-0.3.ebnf
```
