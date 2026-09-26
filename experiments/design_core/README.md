# Design-core — port complete

G4/G5 replaced the Python design-core 0.5 experiment with the [Go parser](../../SDL/go/README.md). The [language definition](../../SDL/docs/profiles/SDL-Structural-Core-Profile.md) and [example](examples/ui-ownership.design) remain. 151 frozen parser cases live in SDL/go/parser/testdata; 170 diagram projections with original source live in SDL/go/viewpoint/testdata.

No active Python CLI/fallback. Historical code remains in Git before G5-M4. From the SDP root: `go -C SDL/go run ./cmd/sdl check ../../experiments/design_core/examples/ui-ownership.design`.
