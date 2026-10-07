# WCI1 final native input evidence

All four runs use binary SHA-256
`d97eac538dbb80816da863e07e3ec4aa4bb6ef0416b62c2e37211591a8badbd8`,
built from the source inventory in ../../candidate-WCI1.json. Go 1.27.1,
Fyne 2.8.1 desktop driver, isolated Xvfb :189 at 1000x700x24, no window
manager, separate XDG config/cache per run. Binary build command from SDL/go:

```sh
go build -tags desktop -o /tmp/sdui-wci1-native-final ./examples/collections/cmd/native
```

Run from repository root, substituting basic/empty/nested/lifecycle:

```sh
python3 SDP/05--Implementation/SDUI/Widgets/native/verify_wci1.py \
  --binary /tmp/sdui-wci1-native-final --display :189 \
  --out /tmp/new-wci1-native-basic --variant basic
```

Basic: 16 checks; empty: 8; nested: 20; lifecycle: 10. All 54 passed,
including exactly-once session teardown in each run. All stderr files are empty.
Actions and checks are NDJSON; raw fixture events are losslessly gzip compressed.
PNG files are OS captures, not Fyne Canvas.Capture. The coordinator inspected
recovery text, clipped rows, independent gutters, fixed sibling and hidden-tree
painting. No WM decoration/IME or future widget-family evidence is claimed.

Motion delivery receives a 100 ms settle before wheel/drag. Provider completion
and state transitions use explicit fixture events/barriers, never timed guesses.
See ../Pilot-WCI1.md for earlier failures and their corrections. Hidden native
controls retain internal diagnostic row geometry: visibility is proven through
widget state and actual painting, not absence from the diagnostic row cache.
