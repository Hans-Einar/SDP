# WCI2-M1 runtime worker — implementation ready for integration review

Worker: James. Loaded/reused SDP 1.1.1, Worker 2.0.0 and document-workflow. Scope: tabs/page/split runtime plus runtime documentation only, API memo and this report. No M2 command/menu/dialog state. No host edits; document.go was explicitly returned to Noether before this milestone. No commits, branch changes or management writes. Main owns Session0010 S3 projection, integration/native evidence and independent review.

## Candidate and coordinated contract

Clone `/tmp/sdp-sdui-widgets`, branch `sdui/widgets-wci2`, base HEAD `53d031e9a9a6e10bdcf56a2ac7b6e863adc81a9d`. WCI1 delivery c39b330 and evidence 53d031e are preserved. Governing clone Panes-and-commands.md selects M1 after independent review of substantive cf93dea6; later status text records authorization. Read-only seam proposals/dispositions remain `/tmp/WCI2-runtime-seams.md`.

`WCI2-runtime-API.md` was published before dependent implementation and sent main, Dalton, Gibbs, Noether and Lorentz. Frontend supplies composition kinds and parser.PaneChildren; runtime does not invent page naming rules. Layout consumes exact Snapshot types and returns SplitGeometry. Bridge owns SDL TabPageID/TabPreviousPageID mapping; preparation discovers CallbackOwners separately from legacy Widgets. Noether acknowledged the final ticket handshake. Clarifications sent directly: disabled visible expanded splits still require geometry; direct page intent remains separate from tabs ancestor activity.

## Delivered runtime behavior

- Concrete copied Tabs/Page/Split states, stable handles and direct page IDs; silent SelectPage, strict SetSplitProportion, CollapseSplit/RestoreSplit, header/divider Focus and remembered EnterPage. Existing Apply updates pane intent/labels atomically; fallback handles hidden/disabled/removed selection and empty body.
- Effective descendant activity without erasing intent. Hidden drafts/data/offsets/focus memory survive; descendant provider tokens revoke only on accepted hide. Provider contexts remain host-owned. Reveal never resets AutoLoadPending or revives a canceled load; existing explicit retry remains bounded. Native viewport updates to inactive pane descendants reject.
- Typed ActivatePage/AdjustSplit event envelopes and one tabs InteractionHandler. Geometry-only probe precedes domain invocation; same-page/programmatic/fallback/reload do not call it. Sequence consumption remains truthful on failure. Full post-consumption state and target value/draft guards reject stale results. Handler payloads are copied. Nested interaction dispatch rejects even across a reentrant legacy Reload. Independently accepted legacy mutations remain accepted and retain their own presentation ticket.
- One typed PresentationGate measurement authority with legacy StateGate compatibility for models without splits. Validate every visible expanded split's finite bounds and exact effective clamp, including disabled splits. Collapsed state retains saved proportion, hidden minima are suspended by layout, failed restore retains previous state. Layout imports remain absent from runtime.
- Speculative gate runs on a disposable copied candidate, so neither preparation nor speculative clamps can leak into final state. Finalized Snapshot goes to PresentationPrepare; its private ticket promotes accepted native pending only after runtime swap. Failed/stale preparation discards only its own ticket; consumed sequence alone cannot publish. Publish is contractually non-failing and non-reentrant. Host supplies source/bundle/size guards and actual resource/native application.
- Detached Successor retains compatible pane identities, hidden drafts/offsets, selection/focus, relative split state and global generation/sequence/request watermarks. Removed/recreated pages get new identity. Axis/kind changes reset split; collapsible=false requires expanded preparation. No prior interaction handler or presentation hook is copied. Legacy Reload retains matching bindings/hooks and .2/WCI1 behavior.

Maintained docs: `SDUI/go/runtime/README.md` and `SDUI/docs/runtime-contract.md` describe actual M1 APIs and typed bridge selectors with explicit native evidence limits. No future-family completion claims.

## Verification

Commands run from `SDUI/go`:

- `go test ./runtime -count=1` — PASS, including the original WCI1/.2 suite and new pane regressions.
- `go test -race ./runtime -count=1` — PASS (1.277s).
- `go vet ./runtime` — PASS.
- Root `git diff --check -- SDUI/go/runtime SDUI/docs/runtime-contract.md` — PASS.

Twenty-two new test functions plus subcases exercise meaningful lifecycle boundaries: copied snapshots; preserved hidden draft/offset; request cancellation and single explicit restart; page intent/fallback/zero eligible; measured strict versus user clamp; positive-minimum collapse/impossible restore; disabled visible geometry; native final ticket matching accepted state; callback probe failure with live request intact; callback error/stale result retaining accepted legacy draft or reload; reentrant dispatch refusal; final resource failure with succeeded outcome; stale preparation discarding only its own ticket; malformed/stale envelope; draft/reserved target conflicts; unknown outcome; malformed geometry; focus fallback and collection memory; compatible/changed/removed successor identities; abandoned candidate isolation; and speculative clamp isolation.

The runtime's numerical gate tests use controlled measured bounds; Gibbs owns real measured-minimum/layout evidence. Ticket tests observe lifecycle callbacks and accepted snapshots, not Fyne rendering. This worker ran no M1 OS input, end-to-end native acceptance or independent gate. Main/Noether own those remaining integration checks. Concurrent lane changes are outside the hashes below.

## Owned candidate SHA-256

```text
f6ced7f0ef649b043a787effc2702e1026460033b3a5bbfa2497ba0898191f04  SDUI/docs/runtime-contract.md
1e737e9dc9a27381b23ae3d158ad24eb2eb50548009e7d368630e704779b7fc4  SDUI/go/runtime/README.md
155730690beab707278fc9e6a5ef1f945c8bdac045d233c88619869dd6add0b8  SDUI/go/runtime/collection_mutation.go
480ed3b71c36df2191dd9a064915fe918519b586a6996858a29edf0e1924c263  SDUI/go/runtime/events.go
b0ebb0777d181878c9862edf2ea609b5bd312d74bc128e98cb470cbfa97a40a3  SDUI/go/runtime/interactions.go
d22a5c49f06102beb3fecedc33b7c2df2a7d9853701304de52119f740baa96a5  SDUI/go/runtime/pane_activity.go
8df0d59b3fc26ce51ca94272eb42dd595dc38fb5b3536b9eaa2c73c6f8868442  SDUI/go/runtime/pane_interaction_test.go
8de9f6c053d2ab0223b9360eb4e189b4e36117b7701a8c29d7cadb097e779440  SDUI/go/runtime/pane_mutation.go
7334c388d85fb9c47e2e6e38b79ce15e12163a4590b8b224fff212e6845c710f  SDUI/go/runtime/pane_successor.go
72cfca2f709c8fa462cc1ade2c5bef6136ab87832f4ace83115767d12d960c45  SDUI/go/runtime/pane_successor_test.go
313f86f3e251f8e9a38fb6141a9994e866f9582add33f2aea0673b9bfc202085  SDUI/go/runtime/pane_validation_test.go
7a97d72a6d7b83558da6d6072884c706c862ec3a780ad95c09f7ca8d3c05e34c  SDUI/go/runtime/panes.go
cfd3f53103267e3a648c8a14c2563b3264ffd06ef15cb7cb0ba95271eda9f318  SDUI/go/runtime/panes_test.go
726d28a669e9eeb5e82d2ed99fe16358ed641fbc5f42eac264087700298a7271  SDUI/go/runtime/presentation.go
2c40ebbdeb17519c34b26e8098ca9753dd5ccd11e7b969f0ae368d7c6d36704f  SDUI/go/runtime/properties.go
9818c55db58b0b68bfbf4ade61634b5b5bcab101694d4007c84e33e8ba229fb2  SDUI/go/runtime/reload.go
5f9535d7c2bda1d7c76fc6b6fda4550237be6c73f64352063e9635d94e0daccc  SDUI/go/runtime/session.go
0a8dac35b0ab372886f0f6196fa231e4c1e139bf48b560353530f376424d5f6b  SDUI/go/runtime/state.go
a3ff520370d6859efa7d67e28c113e7d42c7332ab124b42a8ab0f97ab8889d5e  SDUI/go/runtime/successor.go
80eaafead9f5f2672c7528df34079a2f6c7c58aabb7be64d57b605f1cfc12884  SDUI/go/runtime/types.go
```

Dependency/API source SHA-256 at report time:

```text
f3c8c59944c47b232dd774f28d7e5506571c2bc03a464ed05a8081f3ae81c885  SDUI/go/parser/panes.go
502992d0ba01d60f542b0f7d1f8c8ed90b85ec3a898e05f08b18b3c0acc6b4f0  SDUI/go/parser/normalize.go
c2052d1647ff70244ce5a9de4a0773b7817be29bda61988ca83da14fb9e54da8  SDP/04--Design/SDUI/Widgets/Panes-and-commands.md
c0373687f642623126f33bf9eb11223e1c71e0d9cd3f09c946a2ddbb5e0b671c  WCI2-runtime-API.md
```
