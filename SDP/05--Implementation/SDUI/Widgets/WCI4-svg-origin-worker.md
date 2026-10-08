# WCI4 SVG origin and uniform-scale correction — host worker

Status: product files frozen; component verification complete. Mendel independently approved the final four host files for the correction commit/package3. The coordinator reports the eight-path combined delta committed as `dcf2a74`; final native/package/fullsuite acceptance remains pending. No commit, branch, dependency, canonical record, phase-candidate or original-workspace mutation was performed by this worker.

## Candidate and ownership

Detached worktree: `/tmp/wci4-svg-origin-fix`.
Baseline: `af422b1b1e26bc3d62107f13e74f3c03e00b0a8b`.
Role: SDP Worker (reused SDP 1.1.1 / Worker 2.0.0 and applicable SDUI instructions).
The coordinator explicitly authorized the nonzero-origin derivative and subsequently the single-argument scale correction. Root transform order follows the superseding owner disposition: `displayScale * rootTransform * originTranslation * child`.

Exactly four host product/test/doc paths changed (manifest below). Lorentz owns the separate SDL fixture changes visible in this worktree; this patch excludes them. The coordinator subsequently reported guarded integration into phase/original; those copies were not performed by this worker.

## Repair

Pinned Fyne/oksvg applied nonzero viewBox offsets without the required display scaling, producing incorrect actual pixels. The pinned decoder also interprets `scale(s)` as `scale(s, 0)`, erasing visible geometry. Both were reproduced through public `canvas.Image` rasterization.

The backend callback now prepares/cache-validates one immutable native resource per original resource digest before `PreparePreviews` freezes outcomes. Nonzero origins become a zero-origin root with original presentation/transform on an outer group and negative-origin translation on its inner group. Single-argument scale is expanded in actual root/descendant transform attributes, including zero-origin resources. Ordinary unchanged representations preserve exact bytes. Other child XML is retained. No dependency fork, external renderer, scene model, runtime callback or new API was added.

Provider bytes, identity, digest and original admission budget remain authoritative. Native cache names hash actual derived bytes. The derived representation passes closed-subset/finite-geometry validation and strict native decoding before outcome/capability freeze. The derivative has the existing 4 MiB backend bound: additional wrapper/scale bytes can trigger the declared label/reject outcome before freeze. This explicitly accepted bound is documented and tested. Resize retains the prepared native resource; close drops owned references. No GPU-cache erasure claim.

## Commands and results

Working directory for Go commands: `/tmp/wci4-svg-origin-fix/SDUI/go`.
Prefix: `GOWORK=off GOCACHE=/tmp/sdp-wci0-gocache GOFLAGS=-mod=readonly`.
Environment: `go version go1.27.1 linux/amd64`. No network/sandbox bypass or module mutation.

1. Exact final raster tests against baseline backend using `/tmp/wci4-origin-baseline-overlay.json`:
   `go test -overlay /tmp/wci4-origin-baseline-overlay.json ./host/fynehost -run '^TestPreparedPreviewOriginNativePixels$' -count=1`
   EXPECTED FAIL, exit 1, 1.733s; nine affected cases fail, unchanged ordinary geometry passes. This overlay substitutes only baseline `document_previews.go`; it does not mutate product sources. Receipt: `/tmp/wci4-origin-baseline-final-tests.log`, `.exit`.
2. Corrected targeted tests before final empty-document addition:
   `go test ./host/fynehost -run '^TestPreparedPreviewOrigin' -count=1`
   PASS, exit 0, 2.163s. Receipt: `/tmp/wci4-origin-targeted.log`.
3. Existing and new preview regressions:
   `go test -race ./host/fynehost -run '^TestPreparedPreview' -count=1`
   PASS, exit 0, 70.509s. Receipt: `/tmp/wci4-origin-race.log`, `.exit`. Compiled before the final test-only empty-document addition; production bytes match final manifest.
4. Final exact host candidate, including persistent empty/self-closing cases:
   `go test -race ./host/fynehost -run '^TestPreparedPreviewOrigin' -count=1`
   PASS, exit 0, 34.469s. Receipt: `/tmp/wci4-origin-final-race.log`, `.exit`.
5. `git diff --check -- SDUI/go/host/fynehost`: PASS.

Persistent raster comparisons cover positive, negative and mixed origins; root and noncommuting descendant transforms; uniform single-argument scaling with positive/zero origins, exponent/sign/whitespace spellings; inherited paint/opacity; ordinary shape; and three display scales (0.5, 1, 1.885). Every actual RGBA pixel is compared against an independent zero-origin geometry oracle. Persistent empty-document coverage includes self-closing and explicit empty roots, origins of both signs and single-argument scale. Additional tests cover original identity/bytes, derived cache naming and immutable access, resize/close ownership, finite source whose derived coordinates overflow, and near-limit source admitted while derivative label/reject is truthful.

Reviewer separately reported 27 empty/self-closing/metadata combinations passing race (1.105s), with overlay/evidence under `/tmp/wci4-reviewer-origin`; this is reviewer-reported evidence, not a worker-run command.

## Frozen source manifest

```text
413719c6f98cb05f9e62fb49c9c3429858b5728435df9c110447c900ca26e15a  SDUI/go/host/fynehost/README.md
5fe69cf9dc5b9c3637e6358c43b009e090d5da1d8024af41f606924e60ae2f5f  SDUI/go/host/fynehost/document_previews.go
8f0f9578fddb31c13a22b3366ac132c63c34c352723ff70f55819840f71d5130  SDUI/go/host/fynehost/preview_native_svg.go
d94eda0c1a62cc111a81b0d2a17dfcc416ea0db808f25d047e5aaa3c1f9cbc80  SDUI/go/host/fynehost/preview_origin_test.go

```

Manifest artifact: `/tmp/wci4-origin-host.sha256`.
Manifest SHA256: `cba26b6b973b651ed46de82751c8f4c9f43fc9ad964239b1758c2ab155657cab`.
Patch artifact: `/tmp/wci4-origin-host.patch` (includes both new Go files, excludes SDL).
Patch SHA256: `b10042b9e15387c85b21841a6fea9d80f21fc34f426fd0aab7ef12cf59c7d6d7`.

## Remaining integration evidence

No unresolved worker component failure. Final source correction review is approved; whole M2/package approval remains pending. Actual OS green-pixel bounds for the retained `scale(.5)` fixture, complete SDUI/SDL suites, package assembly and final integration acceptance remain coordinator responsibilities. The coordinator has already committed the combined delta as `dcf2a74`. This correction does not claim broader renderer parity, dark/system-theme fidelity, OS accessibility delivery or native acceptance from headless component tests. Existing light-profile and lifecycle limitations remain unchanged.

Final handoff: all eight detached source hashes (four host-owned, four Lorentz-owned) verified against `/tmp/wci4-origin-delta.json`. Both worker narrow race processes finished with exit 0; no worker test/build process remains active and no further run is planned. Product bytes remain frozen.
