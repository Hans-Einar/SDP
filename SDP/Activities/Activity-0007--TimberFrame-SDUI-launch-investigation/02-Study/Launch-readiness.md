# TimberFrame standalone SDUI launch investigation

| Field | Value |
| --- | --- |
| id | SDUI-STUDY-TF-LAUNCH |
| Status | completed |

Owner request: clone Hans-Einar/freecad-timberframe into ~/git and reproduce the
reported `unsupported-pane` at `CaseParameters/settings (164:2)`. This is a bounded
study, not authorization for product changes or a new release. Primary card:
KB-SDUI-007; Activity ACTIVITY-SDP-0007; Session SESSION-SDP-0025.

## Exact inputs and reproduction

- Consumer clone: `/home/warloc/git/freecad-timberframe`, main at
  `9b7af85320b584c9c987a2aca40d4cfb46dbb016`, unchanged working tree.
- Installed `sdptool` and `gh sdp`: 2.3.2, canonical revision
  `f7a6c9f4793b55de5d0017ae74210135f2a4d365`.
- Source: `SDP/SDL/TimberFrame/WORK--Design/SDUI/FreeCADHost/TimberFrame.sdui`.
  SHA256 `d261086dd2d3c171ef41b784dd972d55b2aa562f73f25ea79b903fccd73e1efe`.
- Entry: `CaseParameters`; line 164 declares `settings=tabs(...)`.

From the consumer repository:

```sh
sdptool sdui check --json -source SDP/SDL/TimberFrame/WORK--Design/SDUI/FreeCADHost/TimberFrame.sdui -entry CaseParameters
sdptool sdui run -entry CaseParameters -watch=false SDP/SDL/TimberFrame/WORK--Design/SDUI/FreeCADHost/TimberFrame.sdui --json
```

Both fail with the exact unsupported-pane path and source position. Check returns
code `io`; run's native helper emits the same cause and the facade reports child
exit status 2. Static `sdui preview` for the same entry succeeds. Evidence retains
actual stdout/stderr and source revision. No display is needed to reproduce the
rejection: admission fails before window creation.

## Cause and governing boundary

`SDUI/go/prototype/check.go`, `prototype.Check`, deliberately rejects composition
nodes in the standalone 0.3 prototype, returning missing host tabs/1. This is
covered by `prototype/panes_test.go` and documented in `SDUI/docs/profile-0.3.md`,
Structural descriptions and SVG. `SDPTool/Contract.md` maps public `sdui check/run`
to this standalone route. A valid SDUI source or static preview is not proof of
native readiness for every host. The 2.3.2 bundle-capacity change does not alter
this admission policy.

SDUI also has a more complete `fynehost.DocumentHost`. The consumer already owns
`tools/preview_gui.sh` and `tools/sdui-preview/main.go`, which use this host and
explicitly supply choices.json and schematic SVG resources. Its GUI README names
this as the intended prototype route. The model therefore exercises supported
DocumentHost functionality that is absent from the packaged standalone route.
Besides tabs, CaseParameters contains number/checkbox/select declarations; merely
removing the tabs refusal would expose additional unsupported controls/providers.

## Proposed next boundary

Use the existing project review launcher for this prototype. A future general
SDPTool launch improvement should reuse DocumentHost, explicitly declare provider
and resource inputs, and report readiness for the selected host. Do not remove
admission guards or flatten the model into another UI merely to bypass the error.
No such implementation is selected by this reproduction request. FreeCAD/CAD,
solver behavior and production Qt integration are outside the study.

## Method and scope

Loaded SDP entrypoint, Change Analysis and document workflow, plus the consumer's
entrypoint and project instructions. Reused Session/Activity and shared KanBan
contracts. One main agent; no separate review or broad test campaign. Local
preflight: producer 73f513d with only its existing ledger lock, completed Activities
1–6, historical clean release worktrees; consumer fresh clean main. Existing
PG/blueprint cards are unrelated. Writes are isolated evidence and new study
records on sdp/timberframe-launch-study; remote writers remain unobserved. No
product/model edits or Framework migrations. The initial gh clone chose an SSH
route that stalled; explicit authenticated HTTPS clone succeeded without global
configuration changes.

The Activity API rejected the SDUI-namespace card in sourceCards because the
project namespace is SDP. No API or board was changed; the primary card link is
retained in the mandate and Session instead. Record creation took place after
initial read-only reproduction; this is a work summary, not an exact transcript.

## Observed project-host result

The unmodified project launcher successfully built and ran Analysis on an isolated
Xvfb display, using Go 1.27.1 and producer source 73f513d. Exit code 0; actual native
window ready; six native pane transitions and three split resize/restore checks
passed. These are the existing programmatic host smoke checks, not OS mouse
injection, screenshot review or FreeCAD integration acceptance. No live desktop or
CAD session was touched; source trees remain unchanged. Only the Analysis entry
was tested, which includes CaseParameters and its nested settings tabs.

From `/home/warloc/git/freecad-timberframe`, the intended project route is:

```sh
./tools/preview_gui.sh -entry Analysis
```

The reproduction is complete. Generalizing public standalone launch to DocumentHost
and explicit fixture providers is a proposed follow-up, not implemented here.
