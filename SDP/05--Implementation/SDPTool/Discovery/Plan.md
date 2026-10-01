# Source-owned discovery and release — ImplementationPlan

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0014 |
| project | SDP |
| state | completed |
| PlanType | ImplementationPlan |
| BranchPolicy | phase |
| CommitPolicy | milestone |
| Systems | SDPTOOL, SDL, SDUI |
| source | KB-SDP-046; owner authorizes Session 0003, implementation and release in one turn |

## Outcome and authority

Discover the selected project's SDP area without navigation.json, return source/
directory-derived navigation, and release a usable Go SDPTool plus updated gh-sdp.
Owner will upgrade XFMD manually. Include completed source composition, output,
Sessions and release-log deliveries in the combined SDP 2.0.0 release. No native
XFMD application changes, filesystem watcher, language extension or live project
upgrade is included. Release is authorized; main merge is not selected. Publish
reviewed release-branch commits, preserving main and unrelated work.

[Session 0003](../../../Sessions/session-%230003--SDP_discovery_and_release.md)
owns the route and turns. KB044 release preparation supplies prior evidence, not
certification of this candidate. Client SPS-007 owns its separate repository work.

## Selected design

- Recognize an explicitly selected SDP directory or its immediate SDP child, with
  no ancestor/repository guessing. Inspect the area with bounded filesystem traversal.
  A directory named SDP is a discovery area, not proof of installed conformance.
- Sources are discovered throughout the SDP area. Actual directories/files remain
  browsable, including old examples. Skip internal Git/operation data; do not follow
  symlinks during enumeration. Show symlinks as entries without dereferencing them.
- Parse supported SDL/SDUI with their owning Go packages. Headers select the profile;
  syntax-only 0.6 fragments remain visible with context-required state. Each System
  root is compiled independently; unrelated Systems are never concatenated.
  Invalid/unsupported sources are visible with diagnostics, not fatal to other roots.
- Project inventory and model IDs are derived, never read from navigation.json.
  IDs use normalized relative source path plus hash to avoid basename collisions;
  renaming a source changes its identity, content changes affect its revision.
- `discover --json` returns one sdptool/0.2 snapshot with inventory and navigation.
  `tree` uses the same discovery service; unfiltered navigation shows all sources/
  Systems. --model selects one semantic model. Viewpoint generation stays on demand.
  Viewer owns memory buffering, file watching, request ordering and refresh.
- Preserve current manifest/receipt reporting and installation-in-progress gates.
  Existing project manifest supplies project name; folder identity is the fallback.
  Board path follows SDP/KanBan. Plans are discovered; ambiguous view-ip selection
  requires an explicit plan path rather than an authored default registry.
- Retire installer creation/use of navigation.json. Preserve old project-owned
  files as inert historical data during upgrades; don't erase user content or
  silently migrate model lists into another sidecar. Remove current repo registration
  authority and schema/guidance, preserving immutable evidence.
- Preserve typed targets, original-source revisions, output containment, stale
  selection rejection, installation signatures and preview/apply/recovery semantics.

## Phases

| Phase | Milestone | Acceptance | State |
| --- | --- | --- | --- |
| DS1 | DS1-M1 | Plan/Session and discovery snapshot, all source roots and errors, source-owned IDs | completed |
| DS2 | DS2-M1 | Selection/viewer/installer/docs integration, native client immutable dependency preparation | completed |
| DS3 | DS3-M1 | Real XFMD read-only discovery, disposable upgrades, full tests and independent review | completed |
| DS4 | DS4-M1 | Clean production package, exact-candidate CI, signed SDP release and public verified gh-sdp release | completed |
| DS5 | DS5-M1 | Truthful publication reconciliation and manual upgrade instructions, Session closeout | completed |

## Git, versions and evidence

Start at c78831b. Stack working phase branches under sdp/session-0003; commit each
milestone and push finished phases. Preserve untracked sourceinput and Node files.
Use a clean isolated worktree for packaging; no deleting unrelated content.
Go 1.27.1, GOMAXPROCS=2 and -p 2. Follow Toolkit/docs/ReleaseChecklist.md and the
installed client release gate. SDP 2.0.0 reflects removed registration/default JSON;
gh-sdp 0.2.0 is independently versioned. SDL language remains 0.6; SDUI remains 0.2.

Test fresh/manual/installed SDP areas, root vs SDP selection, single/multiple
Systems, fragments and invalid files, duplicate basenames, edits/add/remove/rename,
ignored stale navigation.json, no writes and read-only trees, symlink/resource
bounds, revision-bound selection, human/JSON parity and installer transitions.
Independent product/release review uses fresh contexts. Verify original signed
predecessors, immutable package/descriptor bytes and public downloads before
claiming release. Record all commands/results/limitations in Evidence.md.

## Closeout

DS1–DS5 delivered in Session 0003. Publication.json and Manual-upgrade.md bind
actual releases and owner commands. Main is unchanged; phase branches are pushed.
KB044/046 are completed; KB017 and automatic Session/runtime/native-viewer work
retain their existing separate scope. No live XFMD or global extension mutation.
