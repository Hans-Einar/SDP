# Activity command model — T3 evidence

Recorded 2026-11-12 for SESSION-SDP-0016 and KB-SDP-061. Scope is authored SDL
and command design, not executable Activity support or independent approval.

## Candidate and tools

Repository base: `a96affd5378b1cee1a2926a3808bb9e2f2dccbb9`.
Branch: `sdp/sdl-browsing-planning`; planning and model changes are uncommitted.
Tools built successfully from this checkout with Go 1.26.10:

```sh
# From SDPTool/
go build -o /tmp/sdp-activity-design-sdptool ./cmd/sdptool
# From SDL/go/
go build -o /tmp/sdp-activity-design-sdl ./cmd/sdl
```

Local binary SHA-256 identities:

- `/tmp/sdp-activity-design-sdptool`: `3c447fc61cae1981d16dcdf32da1749ef2b8864fe9f8e0f996d220df0bf7890e`.
- `/tmp/sdp-activity-design-sdl`: `72cd16738f97cf49d891ea6d53b9bda1a45b95c2004600ad4db1b089b9d3af85`.

## Source verification

From the repository root:

```sh
/tmp/sdp-activity-design-sdl check SDP/SDL/ProjectGovernance/SDPTool/ActivityWorkflow/WORK--ActivityCommands/System.design
/tmp/sdp-activity-design-sdptool SDP/SDL/ProjectGovernance/SDPTool/ActivityWorkflow model status work:ActivityCommands --json
/tmp/sdp-activity-design-sdptool . discover --navigation lazy-sdl --json
```

- The first authoring check reported noncanonical form. Sources were canonicalized
  using the SDL formatter; the final check passes with `valid: true`, no warnings,
  profile `design-core/0.6`, five source files and six authored scenarios.
- Composed graph revision: `97fb12cbeeccf6593307a7a42f55743a5982be7245ff59d2eeda7cb26abfad70`.
- ModelGovernance WORK identity: `b3e45cab-c629-490b-9192-3ed2b4628387`.
- Checkpoint: `b3e45cab-c629-490b-9192-3ed2b4628387:00001`; clean and preliminary.
- Source inventory/live digest: `95cf3330eb60e9427517d924bcd4acda0dde5643ad0d5971d2cb0d40790c1e38`.
  This inventory digest is distinct from the composed graph revision.
- Lazy discovery finds System.design as validated and preliminary. Its four included
  files remain context-required when discovered alone under the existing browsing
  behavior. That is not a model validation failure; browsing changes remain queued.
- Existing SessionLifecycle WORK remains clean at :00003 with unchanged live digest
  `3ca418bedb135935605b75a7aa88647ef5fe01098777b323a69e30b8b946e495`.

## Limits

The checkpoint preserves source history; it is not a Git commit, accepted model,
release or publication. No viewpoints or Blueprint bundles were exported. No Go
product source changed. Runtime atomicity, permissions, single-active enforcement,
agent disposition, timeline rendering and executable command availability remain
future implementation/verification work. Parser success does not prove them.
