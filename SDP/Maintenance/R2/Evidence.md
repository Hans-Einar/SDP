# R2 — verification evidence

## R2-M1

Read R1 baseline/current entry points; registered ten concrete conflicts in Findings. KanBan schema/replay/locations and local links pass: 16 cards, 27 events, 1907 local Markdown targets. Historical ledger prefixes, 574 generated artifacts and R1 migration boundaries preserved. `git diff --check` passes. No code or generated model changes.

## R2-M2

Checked document changes against SDL parser/vocabulary, bridge.Plan, SDUI Handle and existing G3/G4/G6/G5 evidence. No grammar, Go source/module, installation-manifest or generated-viewpoint changes.

CLI spot checks from repository root using local Go 1.27.1:

```sh
go -C SDL/go run ./cmd/sdl check ../../SDUI/design/architecture.design
go -C SDL/go run ./cmd/sdl action-check examples/echo.sdl
go -C SDL/go run ./cmd/sdl class-check examples/runtime-classes.sdl
go -C SDUI/go run ./cmd/sdui ../examples/concept1-bucking.sdui --format svg --entry bucking -o /tmp/concept1.svg
```

All three SDL profiles returned valid=true. SVG used a unique temporary output directory: valid XML/SVG, viewBox 0 0 1920 1200, 402669 bytes. This is structural CLI verification, not a new visual/native GUI trial. Old G-phase test counts were not rerun. Full Go regression was not repeated for documentation-only changes.

R1 checks local Markdown targets, 574 unchanged generated artifacts, historical fingerprints, Go sources and installation policy. K3 schema/replay and physical locations check all three boards; K4 ledger prefixes remain byte-preserved.

Limits: no anchor/external URL checks. No new current-status claims about other repositories/PRs; G6 external evidence is dated. KB-SDP-011's nine older Toolkit ID mismatches were neither changed nor retested. Phase/template profiles and broader candidate semantics remain open.

Final checks passed: 16 cards / 29 events across three boards; 1973 local Markdown targets; all 17 SDUI subject documents indexed. Git diff against K4 confirms unchanged History, Go, design, Template and installation files. R2 contains only Markdown and append-only ledger events. `git diff --check` passes.
