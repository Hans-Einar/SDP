# WCI2-M1 SDL bridge and panes fixture — worker handoff

Status: bounded lane complete, ready for independent review/integration; holding
at coordinator request after persisting this report. This is not stage acceptance.

## Authority and candidate

Coordinator explicitly reassigned Architect to **SDP Worker**, WCI2-M1 only.
Loaded/reused sdp 1.1.1, sdp-worker 2.0.0 and document-workflow. Governing contract:
`SDP/04--Design/SDUI/Widgets/Panes-and-commands.md`, reviewed substantive hash
cf93dea68510be05a626a3953159fc375ce04f55b18562c35eb3699d0859b839 plus recorded
M1 transition; implementation seams in `WCI2-runtime-API.md`,
`WCI2-frontend-API.md` and `/tmp/WCI2-runtime-seams.md`.

Workspace: `/tmp/sdp-sdui-widgets`, branch `sdui/widgets-wci2`.
HEAD: `53d031e9a9a6e10bdcf56a2ac7b6e863adc81a9d` plus uncommitted concurrent M1
work. The nine lane files below were rehashed when persisting this handoff and
match the tested/reported lane bytes. Go: `go1.27.1 linux/amd64`.

Exclusive implementation scope was SDL/go/bridge and NEW SDL/go/examples/panes.
No collections, runtime, host, frontend, dependencies, management records, commits
or branch changes by this worker. This root handoff was separately authorized.
Main owns integration, native evidence and Session updates. M2 remains unselected.

## Delivered behavior

- Closed text-only TabPageID/TabPreviousPageID selectors, tabs callback-owner
  discovery through runtime CallbackOwners, and BindInteraction installation only
  after complete plan preflight. Legacy Widgets and bridge handler behavior remain.
- Runtime-owned direct-page identities are consumed without a parallel resolver.
  Shared/reused tabs bind independently; failed rebind preserves prior handlers.
- InteractionReply preserves domain success through post-Execute receiver/draft
  or revision conflicts; execution errors/malformed output report unknown.
  Stale UI replies never replay domain execution.
- Real SDL Page fixture, tab-owned inputs/tree scroll, positive split minima,
  retained drafts/offsets across selection and compatible reload, controlled lazy
  provider barriers, action/resource/reload failure injection and JSON inspection.
  Native CLI uses actual input for user gestures; no synthetic activate command.

## Verification performed by this worker

Working directory for Go commands: `/tmp/sdp-sdui-widgets/SDL/go`.

```sh
go test -race -count=1 ./bridge ./examples/panes ./examples/collections
go build -tags desktop -o /tmp/wci2-panes-native ./examples/panes/cmd/native
```

Both passed. Final race run reported bridge 1.638s, panes 39.137s and existing
collections 21.959s. Earlier scoped bridge and provider-barrier race runs also
passed. From clone root:

```sh
git diff --check -- SDL/go/bridge SDL/go/examples/panes
```

Passed. Tests establish real SDL IDs/callback counts; zero-call and atomic
preflight rejection; silent selection/reload; reused owner independence; preserved
outcomes and newer receiver drafts; rejected final geometry/resources without
speculative promotion; scroll/draft/split retention; failing reload preservation;
late provider completion and cleanup. Headless fixture tests are not native proof.
Initial local fixture syntax errors were corrected; partial-runtime and host
measurement integration blocks cleared before this successful final run.

Built binary SHA-256:
`45e89e08ddc501e89d961275883c037f5f7114a230187e1f83ee8efb3ec075d5`.
The binary hash was rechecked at handoff. Other lanes were evolving concurrently;
this is HEAD-plus-working-tree evidence, not an unchanged committed candidate.
Main must freeze/rebuild the integrated candidate for final acceptance.

## Exact nine-file SHA-256 inventory

```text
b20d3624a2d02e16fa517460f5e05bf7d05a95379ae087389713702e7c18b1f4  SDL/go/bridge/bind.go
dc92c8ff4cb5b9c8a3ffbd19c147682e1c2a7a2f39d5f12acf7ec8999c4f2b0b  SDL/go/bridge/interaction.go
6cc47c696ab3d6e9dd08e581502c44704e015dcc251d013aabcf49aef2e54e5e  SDL/go/bridge/panes_test.go
9d1602e82869595262f3b9dcb43c0d27d6e11af2ba640b2631a03fa3b07002e2  SDL/go/bridge/values.go
40ebb89507b47da8aaf40fc4cf880a3ce2cb35a52d261b835614755992e68e12  SDL/go/examples/panes/README.md
4ef282ddab3d290c1ec9036414afc4f07c26d16caa60f6abc0f8d819e7d5f9ac  SDL/go/examples/panes/cmd/native/main.go
6563a4115b83b2ea9b93566c56687945442a81dbf4c82dcde578c0b553778363  SDL/go/examples/panes/controls.go
fdbcbe8750f0ce9804e90641f02330b74b08a0a93c34c12bde823fc805abc9e1  SDL/go/examples/panes/fixture.go
d2e92ec03e2bd842d280068de51b7082f3e23be298e9a962d653a0240e885122  SDL/go/examples/panes/fixture_test.go
```

Scope SHA-256: `afadc05bee678e00c4ae2bda8e64e77276bc883d8747741ca7789eb6b49a4ddb`.
Algorithm: SHA-256 of the nine path-sorted lines above, each exactly
`hash`, two ASCII spaces, repository-relative path, LF (including final LF).
This handoff file is excluded to avoid a self-referential hash.

## Latest coordinator-reported native status and limits

Main reports native pilot passes for tab callbacks, drafts/inactive offsets,
eligibility, error/invalid outcomes and draft conflict. The resource-failure
harness initially retained the injected preview Dirty state; main corrected it
through actual Tab/Escape recovery. The product guard behaved correctly, and no
host bug is confirmed from that incident. This is attributed coordinator feedback,
not a native run or independent verification performed by this worker. Final
native checks and independent integrated review remain main's responsibility.

No command/menu/dialog selectors, acceptance framework, WCI3 values, release or
inventory closure is delivered here. No native/WM conformance claim follows from
this lane's tests or binary build. Next action: hold for coordinator review or a
concrete bounded correction; no further implementation started.
