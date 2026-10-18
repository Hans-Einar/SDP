# WCI3-M2 frontend worker delivery

Status: component implementation complete and frozen for integration/independent
review; no native/IME or whole-stage acceptance claim. SDP Worker and shared document
workflow reused. Work is uncommitted on sdui/widgets-wci3 in the isolated clone.
Phase HEAD: ea49991f2065679c93e39fe02a3d458fe72803df.
Canonical ORIGINAL Values-and-text inspected SHA-256:
8199c275fecff3f6d03835540a3e586dd418ff70470f641590d2d75ace8f51be.
Main owns management, canonical architecture/runtime requirements, dependencies and
OS-native proof. Concurrent other-lane changes were preserved; none is claimed here.

## Delivered behavior and API

* parser.InputPolicy and InputOptions(*Instance) (InputPolicy,error) are exported.
  Extended is true only for explicit new argument presence, including false/empty.
  PlaceholderSet preserves omission; effective placeholder falls back to text.
  Booleans default false only in returned policy, never inserted into source.
* .3 input admits exactly multiline/readOnly/required booleans and placeholder
  string alongside old text/value/callback. .2 rejects new arguments. Blank labels
  and required-empty remain source-valid. Runtime owns content/validation rules.
* Strict selected-root resolution checks normalized input shape/types, including
  malformed constructed Literal payloads. Extended definition roots can own a
  callback using their definition name, consistent with new scalar definitions;
  anonymous callbacks reject. Basic input callback naming rules remain unchanged.
* Existing code generation needs no new emitter or AST fields. A compiled generated
  constructor harness verifies exact Unicode/newlines, false/empty/absent arguments,
  source spans, callback refs, independent reuse and fresh constructor copies.
  Version03, sdui-ast/0.3 and sdui-go-model/2 identities remain unchanged.
* Text/Markdown/combined descriptions identify static input properties and source
  initial text, including hidden controls, path/span and unverified bindings.
* Public SVG rejects every extended input with unsupported-text-export before
  writing an artifact, even hidden or with supplied geometry. Prepared native
  omission still requires the exact input inventory; wrong/missing/extra entries
  reject. No broader SkipControls behavior or new widget kind was introduced.
* Standalone prototype explicitly reports unsupported-text with source/use
  provenance, no readiness status. It has no extended document-host adapter; this
  does not require a legacy-launcher migration. Connected preparation/native proof
  is owned by host/bridge/main and cannot be inferred from source acceptance.

No existing exported Go signature changed. New API is additive. ResolveInteractions
now rejects malformed .3 constructed input schemas. Public extended SVG and standalone
preview intentionally reject newly admitted source; valid legacy paths remain intact.
The .2 SVG path was left on its original checks. No runtime, layout, host, SDL,
module, product architecture/requirements, prior evidence, or management file edited.

## Verification

PASS from SDUI/go:

```text
go test -race ./parser ./codegen ./presentation ./svg ./prototype ./cmd/sdui
```

This includes A01 positive/negative schemas, argument-presence/default distinction,
.2 rejection, malformed AST/normalized models and original diagnostic spans;
compiled generated constructors and existing .2 generated-byte/provenance regression;
A09 static descriptions and source-linked SVG/native-inventory boundaries;
standalone unsupported preflight; CLI retained-artifact failure behavior.
All six packages passed on the final scoped product files below.

Independent baseline comparison: archived SDUI from ea49991f into a temporary
non-workspace directory, built its cmd/sdui, and compared current cmd/sdui bytes for
`Main=[field=input("Label",value="å🙂")];` under both .2 and basic .3. AST, dump,
Markdown and SVG all matched byte-for-byte (eight comparisons). Baseline hashes
are now fixed assertions in cmd/sdui/input_test.go. Existing generated .2 bytes are
also compared against the committed SDL generated-model fixture; it was not edited.
The original M1 negative tests for four future input properties were replaced by
selected M2 positive/presence tests and .2 rejection tests; textarea still rejects.

Scoped git diff --check passed. Temporary test failures were incorrect test width
(800 instead of supported <=400) and nonexistent CLI composition format; tests were
corrected to actual APIs, with combined composition covered through presentation.
No product behavior was weakened to accommodate those test assumptions.

## Coordination and remaining proof

Actual InputOptions API was published directly to runtime James, layout Gibbs,
host Noether, bridge Lorentz and main before dependent implementation. Their state,
measurement, bridge and native adapters use their own files; no competing stubs.
The memo's new checkpoint distinguishes implemented APIs from historical proposals.
No unresolved frontend contract departure remains. Required-empty reload, CR/LF,
exact typed Commit and native rollback/history rules are the reviewed canonical
policies implemented/proved in the responsible lanes, not frontend evidence.
A02/A06 runtime/real-SDL and A07/A08 native editing/OS IME remain integrated proof;
A09 lifecycle and publication also require runtime/host checks. Component race tests
cannot certify OS input or all-platform IME. Main/reviewer must assess the assembled
candidate; this report makes no release/merge or whole-stage delivery claim.

## Exact lane freeze

19 payload paths below, plus this report (20 written paths total). Report excluded
from its own digest. Hashes are SHA-256 of file bytes, relative to clone root.
This manifest covers only this lane, not the concurrently changing whole workspace.

```text
467df50aca1163c7f258e8d9db3959dfb38b81dd98218028a56b5b3a90899180  SDUI/docs/go-generation.md
fe625b8b62d1b8f9d1ece9eac064576454e189b9c5ecb8d2f0663454de947930  SDUI/docs/profile-0.3.md
12a5c2247a17436088f26889a2437cf95af64a16562a4327e9d63b6d580ce45e  SDUI/go/cmd/sdui/input_test.go
0f49e467ace3acb3aaa990e0f2f40ab39cf426d42eeaaaaeddbd1a7d829a84ca  SDUI/go/codegen/input_test.go
38c105f5a831877137e4ddc2a044b0323a8adc8a987ceab287359af54e600f83  SDUI/go/parser/input.go
a8e829e1a2500bf767558e53164befcb6cbe07d68f807c702a0c159bc1d8b97c  SDUI/go/parser/input_test.go
11135e5edf95ff87f7a494870759c510b184c912fc1b97ed74a212c128712f11  SDUI/go/parser/interaction_identity.go
4e97c5677db796aa50ac7d8004ed860b5f101be362134527863217f352d8fc3b  SDUI/go/parser/scalars_test.go
695d91e2051f899ade8f9640091d5c917ec3d20cb1066bf0e0fe7585185520d0  SDUI/go/parser/validate.go
5d7f64f0fb80a2c137a57caf7cc9e8cc8031a3493ec234c22c2bdb5544fbbbfc  SDUI/go/presentation/dump.go
fefc2828e61609fad788b0509a3c3ed07cfa4da8693152e48e83dc0b11d275d2  SDUI/go/presentation/input.go
aa0eba1991bece4f936b9ba259309c42d46c4b46c5e319dd1329f06015c46508  SDUI/go/presentation/input_test.go
24506139fbf8b4ccf87f0c4c120f5a5906e1c3defb1e179259feb42c1fbe390d  SDUI/go/presentation/markdown.go
7e92f154c2e9f9886ad1bbb4cc285892d54a6c699f182e39c6bff830c1f0aaa7  SDUI/go/prototype/check.go
0a7c49c21a35f52bb7fbc85d2efa2fa2c47c6c74e1a1fc0e5bbb036bfe37fce4  SDUI/go/prototype/input_test.go
5a515532477bf90dc727078fef0f641ee30517558e2c00359f68d990b4405738  SDUI/go/svg/check.go
1049c6af37f7a9862222d3e5a4007a29dee5739843bf7372828e7631306fecb9  SDUI/go/svg/input_test.go
68bf4283d3f8100a9271e2ade1f34470169529726ed051d54fcdcbbdbb87524b  SDUI/grammar/sdui-0.3.ebnf
56c125b25a09f1e2f2af1cb0439ed54e8d8e53267d403707b60a205318dd7532  WCI3-M2-frontend-API.md
```

Manifest SHA-256 (UTF-8 lines above, two spaces between digest/path, final newline):
1b21d3c65e7d0644f99be2337865ea8bbd7866e82454e4555e180a83579a83a8
