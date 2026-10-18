# Native input tooling

`x11_input.py` delivers actual X11/XTest input to an explicitly named window on
an explicitly supplied isolated display. It uses system libX11/libXtst through
ctypes; it does not require xdotool or Python Xlib. It never chooses the owner's
active display automatically. Use a separate Xvfb display and XDG configuration
and cache paths for native acceptance applications.

## Tool smoke check — Session0010 T002

Built a temporary native Fyne button probe with Go 1.27.1, `-tags desktop`, using
the SDUI module dependencies. Started Xvfb `:189`, screen 1000x700x24, no TCP.
The temporary probe used title `SDUI XTest probe` and separate XDG paths under
`/tmp/sdui-x11probe`. Exact tool invocations from repository root:

```sh
python3 SDP/05--Implementation/SDUI/Widgets/native/x11_input.py --display :189 --title 'SDUI XTest probe' locate
python3 SDP/05--Implementation/SDUI/Widgets/native/x11_input.py --display :189 --title 'SDUI XTest probe' click 150 100
python3 SDP/05--Implementation/SDUI/Widgets/native/x11_input.py --display :189 --title 'SDUI XTest probe' key Tab space
```

Observed window origin (0,0); click produced `native-activation:1`; Tab followed
by space produced `native-activation:2`. A space sent before keyboard focus did not
activate the button. The process was terminated after the check.

This proves the input route can reach the native driver. It is not collection,
viewport, reload, IME or WCI1 product acceptance evidence. Product acceptance must
supply its own source/binary inventory, fixture logs, actions and inspected images.
