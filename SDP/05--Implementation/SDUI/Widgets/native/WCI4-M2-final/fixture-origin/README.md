# Bounded SVG-origin fixture handoff

Only `/tmp/wci4-svg-origin-fix/SDL/go/examples/previews/**` was edited.
Detached HEAD `af422b1b1e26bc3d62107f13e74f3c03e00b0a8b`; no commits,
branch/module/host/product API/SDL semantic changes. Phase and original untouched.
Reused SDP Worker; Main owns Session, native proof and final matching packages.

Added positive-origin, mixed-origin and root-transform to existing ResourceSlot.
The existing resource condition command accepts them without a controller change.
Default and negative-origin bytes remain unchanged. All new colors are #238248.
SVG-origin-bounds.md gives exact S × Rroot × Torigin × Cchildren pixel expectations,
including image fit, clip, device scale and edge-antialiasing interpretation.

Actual command, from detached SDL/go, with GOWORK=off and GOFLAGS=-mod=readonly:

```
go test -mod=readonly -race ./examples/previews -run '^TestOriginResourceSlots$' -count=1 -v
```

PASS 1.054s, exit0. Four subtests validate source dimensions/digest/source identity,
closed SVG subset, declared rectangle/color/transform matching the pixel oracle,
and real resource condition admission for Hero/ReportSVG/DetailSVG. No app, host,
window, rasterizer, native build or full fixture suite was run. This establishes
resource/condition validity, not rendered/native correctness. Main owns OS assertions.

command.json records exact argv/cwd/compiler/environment/HEAD; test.log and exit.txt
are original subprocess captures. before.json and after.json match for all four
source/docs files below. No further edits are planned in this bounded delta.

| File under examples/previews | SHA256 |
| --- | --- |
| `resources.go` | `ddda535931b6d64f328bd3cf34203604a37d91d855a03d4a64fc979a5fb92efb` |
| `origin_resources_test.go` | `2580359cc45ad7ec47c9885287e7cfc9c0b2ef825cfb3c8f1935ba5b3a0caed2` |
| `README.md` | `9819c82651dd1c6e34de0788d0582ecb032aeb73e6225599aac45356176c9c67` |
| `SVG-origin-bounds.md` | `21cd02ac99dab47c98bed5e48cba5c01daa61c4f1b1eb4d8bbba2431fa87e199` |

Test log SHA256: `4209ad18ee58becfcc7e111ffa37e28046e073b62f0648afd802566e4c613494`.
