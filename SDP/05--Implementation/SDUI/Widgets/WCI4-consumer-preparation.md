# WCI4 consumer preparation — read-only reconnaissance

Session0010 T003, during WCI3 implementation. This is preparation evidence, not
WCI4 implementation, packaging or external XFMD acceptance. The actual available
consumer root is `/home/warloc/git/xfmd-sdl-navigation`; its dirty work is untouched.
[Inventory](xfmd-consumer-inspection.json) pins the inspected source bytes and HEAD.

The actual build recipe `tools/build_sdui_tools.py` builds `sdui-preview` and
`-tags desktop` `sdui-fyne` from an explicit SDUI/go source argument into an explicit
output. It does not have its own Go module root. Consequently, a selected GLFW
replacement in SDUI/go will apply to this recipe; it must still be verified in
compiled module/build info and a real composition run. At final preparation, invoke
with `GOWORK=off`, `GOFLAGS=-mod=readonly` and a new isolated output directory.

`cmake/Sdp.cmake` consumes `XFMD_SDUI_TOOLS_DIR`, defaulting to the build directory's
`libexec/xfmd`, and installs both executables optionally into `libexec/xfmd`. Its Go
license payload comes from consumer `LICENSES/go/`. An optional install clause is
not proof that a produced archive contains current helpers or matching licenses.
Producer preparation must explicitly inventory the two files and licenses.

`SduiTools.cpp` resolves absolute executable overrides `XFMD_SDUI_PREVIEW` and
`XFMD_SDUI_FYNE`, then executable-relative `libexec/xfmd` locations. No PATH search
or compile-on-launch is part of this route. `SduiWorkflow.cpp` sends source, entry
and revision to `sdui-preview -check`; launch readiness requires the matching reply.
`SdpProtocol.cpp` expects the existing `sdptool/0.2` envelope. This is independent
of 0.3 SDUI source support and does not call for a protocol version bump.

The actual `SduiWorkflowGuiTest` covers 0.2 combined composition, source spans,
launch readiness, missing/failed helper, window independence and stale source.
It unsets the Fyne override during its test, so merely setting both environment
variables is insufficient for a faithful staged test: its executable-relative
helper location must also contain the candidate. Any staging uses copied test
binaries/new directories, never replacement of the consumer's installed helpers.
All-family 0.3 connected acceptance remains the producer application's real SDL
fixture because standalone helpers cannot manufacture providers or bindings.

Next at WCI4: recheck consumer hashes, build matching helpers and SDPTool from the
final candidate, inspect licenses/native dependencies/build info, exercise actual
protocol and native consumer fixtures, and record payload versus external package
integration separately. No external application source change or publication is
inferred from this inspection.

## Reviewed launcher boundary

Coordinator and independent reviewer checked the whole-card outcome against the
selected contracts. The 0.3 runnable route is DocumentHost with real application
providers/SDL bindings, demonstrated by the connected fixtures. Existing standalone
`sdui-fyne`/`prototype.Check` uses RuntimeView and explicitly rejects 0.3 panes,
commands and scalar adapters, including provider-free controls. This is a separate
host-adapter limit, not merely a missing-provider diagnostic.

WCI4 must stage and name the actual runnable connected route alongside the matching
producer/helpers, and show supported/fallback/unsupported rows truthfully. It may
not claim XFMD Launch or the legacy standalone helper exposes the new families.
The selected inventory does not require migrating every launcher to DocumentHost;
no such route change is silently added during packaging. A matching helper payload
and successful protocol test establish compatibility, not expanded launcher support.
