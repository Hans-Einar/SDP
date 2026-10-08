# WCI4 helper recipe/IME preparation pilot

This is a pre-WCI4-product pilot from clean detached source 90a94b5, not the final
provider/package candidate. The actual unmodified XFMD tools/build_sdui_tools.py
built both private helpers from /tmp/wci4-packaging-baseline/SDUI/go into an isolated
stage, with GOWORK=off, GOFLAGS=-mod=readonly and the shared Go cache. Build exit 0;
module information selects the local patched GLFW from the SDUI root. No external
consumer checkout, installed helper or owner desktop configuration changed.

Attempt 3 passes eight actual configured IBus/XIM assertions on sdui-fyne. Starting
from accepted A, a dirty B suffix then composition Return does not accept AB; plain
Escape restores A. Composition Escape preserves dirty AB, and ordinary Escape then
restores A. Ordinary Return subsequently accepts AB中. External clipboard reads
replace a sentinel and inspect exact bytes. Native window close exits zero. Fixture
stderr is empty; launcher D-Bus/portal warnings remain explicit. The inspected
original-resolution image shows the composed text. No SDL or new-family launcher
capability is claimed by this basic 0.2 helper test.

Attempt 1 timed out finding the window although its captured OS image shows the
helper. Read-only XFetchName inspection showed Latin-1 WM_NAME bytes (middle dot),
while the old harness decoded invalid UTF-8 with replacement characters. The X11
helper now falls back to Latin-1 when UTF-8 decoding fails; ASCII names are unchanged.
Attempt 2 was intentionally stopped during that harness correction, not a product
failure. Attempt 3 uses the corrected title lookup. All attempts remain retained.

The final WCI4 package must rebuild from its exact final source and repeat these
checks; this pilot neither closes WCI4-A06/A07 nor establishes publication.
