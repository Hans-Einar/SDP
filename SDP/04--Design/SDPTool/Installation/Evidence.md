# IPD design evidence

Evidence distinguishes parser/model consistency, generated projections and authored
contract review from installer execution. No Go installation or live XFMD upgrade
is performed. The current tool was built from SDL/go at planning baseline db7d82d
with the existing Go 1.27.1 linux/amd64 toolchain into /tmp/ipd-sdl.

## IPD-1-M1 — installation responsibility model

The canonical model adds the project maintainer, four install/upgrade/adoption/
recovery use cases, SDPTool installation units, the external gh-sdp process and
external release artifact service. The latter is a collaborating service boundary,
not a new SDP-owned server. Five new delivery activities are explicitly planned.
Library/package boundaries are not modeled as deployed containers.

Commands: `sdl format`, `sdl check`, `sdl ast` on the canonical SDPTool.design.
The first draft exposed declaration-before-fact and canonical-order requirements;
the formatter normalized it, then check/AST passed with no diagnostics and 90
symbols. A set comparison preserves every statement in the db7d82d model.
The AST is parser output, not a manually constructed JSON model. Detailed channels,
records and scenarios belong to IPD-2; no protocol implementation is inferred.

