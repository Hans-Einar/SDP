# Pinned GLFW X11 input filtering

WCI3-M2 selects the local sibling `glfw` module at upstream version
`v0.1.0-pre.1.0.20260707082822-2a407d02d01a`. Its 145 source/license files are
preserved with exactly one conditional change in `glfw/src/x11_window.c`.
`source-inventory.json` records upstream and selected SHA-256 for every file;
`xim-filter.patch` reproduces that change. This is a source dependency, not a
new public release or a replacement of the machine's module cache.

The upstream X11 event path called `XFilterEvent` but delivered the key callback
and advanced its duplicate-event timestamp before consulting the filtered result.
An IME-consumed Return could therefore submit the previous text before the
composition committed. The patch guards both key delivery and timestamp update
with `!filtered`. Character lookup keeps its existing filtered guard. Other
platform sources and ordinary unfiltered-key behavior remain upstream bytes.

The reviewed mechanism experiment is recorded in
../../../SDP/04--Design/SDUI/Widgets/ime-probe from the repository root. Actual
extended SDUI/SDL composition, cancellation and subsequent ordinary submission
remain required native evidence; this source inventory alone proves no IME support.
The selected scope is X11/XIM with the tested configured engine, not all IMEs or
Wayland. A future upstream fix can replace this copy only after equivalent native
regression evidence and an explicit dependency update.

SDUI/go and SDL/go explicitly replace `github.com/go-gl/glfw/v3.4/glfw` with this
copy. Go replacements are not transitive. A consumer compiling generated native
code in a different module must explicitly select the same reviewed dependency.
The inspected XFMD helper builder compiles from SDUI/go, so that recipe uses this
replacement when supplied this source tree; package evidence must still verify
the actual staged binaries and build information. SDPTool itself does not compile
Fyne and needs no unused replacement. No executable is built during UI viewing.

Retain all upstream notices when distributing source or binaries: top-level
`glfw/LICENSE`, bundled `glfw/glfw/LICENSE.md`, and the notices embedded in bundled
third-party source. WCI4 owns the actual package license/native-library inventory.
Run `python SDUI/third_party/glfw-policy/verify.py` from the repository root to
check every selected file, absence of extras, and the exact patch reversal against
the upstream hashes. This does not fetch dependencies or modify their bytes.
