# XFMD blueprint integration handoff

Status: blueprint implementation originally inspected at d2cc760; Actions/JSON
producer delivery is recorded in Session0011 on 2026-10-18. This is not a published-release or native-GUI acceptance
claim. Session0008 T013 records this handoff.

## Start here

- [Producer and command documentation](blueprints/README.md): task JSON, generation,
  immutable retention, assessment and limits.
- [Consumer contract](Contract.md): navigation envelope and BPI3-M2b assignment nodes.
- [Lifecycle contract](../SDP/04--Design/SDPTool/Blueprints/Assignment-Lifecycle.md):
  transitions, request schema, authority, retry and evidence requirements.
- [Executable lifecycle trial](../SDP/05--Implementation/SDPTool/Blueprints/Trial-BPI3-M2b/README.md):
  generated requests/results and discovery from the compiled CLI.
- [Delivery evidence](../SDP/05--Implementation/SDPTool/Blueprints/Evidence-BPI3-M2b.md).

## Actual integration boundary

SDPTool exposes reusable Go packages; its CLI calls them. This does not expose a
network server or a C/C++ ABI automatically. XFMD's inspected working source uses
`src/application/sdp/SdpToolJob.cpp`: an asynchronous subprocess with argument
vectors, JSON output, timeout/cancellation and diagnostics. Its default command is
`gh sdp`; `XFMD_SDP_TOOL` selects one executable path instead. No shell command
string should be constructed from document content.

`SdpNavigation.cpp` requests discovery; `SdpProtocol.cpp` decodes the protocol;
`SdpPanel.cpp` creates tabs from navigation root nodes with `kind: "tab"`.
Thus the current consumer has the generic foundation to show the Blueprints tab
when the selected producer returns it. This inspection is not an end-to-end test
of the installed XFMD build against this unreleased producer.

The wrapper forwards CLI arguments to a verified SDPTool binary. Replace
`sdptool` with `gh sdp` in the examples only when that wrapper selects a distribution
containing these commands. Updating the wrapper alone need not change the engine.

## Commands and ownership

Paths and model names below are examples, not an instruction to mutate XFMD.
Generation uses a ModelGovernance area; assignment operations use the SDP area.

```sh
sdptool /path/to/models model create blueprint from release:0.1.0 to work:Calibration --entry System.design --task /path/to/task.json --catalogue /path/to/project/SDP/Blueprints --json
sdptool /path/to/project discover --json
sdptool /path/to/project/SDP model assignment list --json
sdptool model assess blueprint --bundle /path/to/retained-revision --evidence /path/to/evidence.json --json
sdptool /path/to/project/SDP model assignment apply --request /path/to/request.json --as ACTOR --authority controller --json
```

Use `--output DIRECTORY` instead of `--catalogue DIRECTORY` for a replaceable
preview. WORK previews are preliminary; generation does not freeze WORK or assign
an agent. Task permissions, retained revision and assignment identity are distinct.
Assessment can emit a blocked result with exit code 3: preserve its JSON response.

Discovery derives the Blueprints root, retained task/revision entries and nonempty
assignment-state groups from bundles and canonical ProjectManagement history.
No navigation registration file or GUI-owned work-state database is needed.
Assignment nodes expose assignmentId, assignmentRevision, workState, assignee,
sourceFreshness, readinessStatus and evidenceStatus. Multiple assignments can use
one revision; do not collapse them into one task status. Historical state survives
missing inputs, with diagnostics and unavailable targets where appropriate.
Use returned open targets and revision hashes. Refresh after mutations and on
relevant filesystem changes. Never edit generated bundles to change workflow state.

The lifecycle includes draft, ready, assigned, in-progress, review, completed,
on-hold, canceled and superseded. Read the exact transition contract before adding
buttons. Requests bind expectedEvent and a stable event ID for safe retries.
The local CLI's --as/--authority are attribution, not remote authentication. A future
adapter must establish the caller independently. Lifecycle writes currently require
Linux; parsing or blueprint generation does not prove implementation conformance.

## Dynamic actions and JSON input — implemented

[Actions.md](Actions.md) is the current producer contract and command guide.
`sdptool [PROJECT] actions --json` advertises supported actions and typed input
schemas; `sdptool invoke --request -|FILE` executes one JSON request. The registry
contains discovery plus blueprint generation, assessment, assignment list/apply.
Existing human CLI routes share the domain handlers. Caller attribution remains
separate from JSON. Catalogue support and contextual availability are distinct.

Native XFMD still owns menu/Preferences/toolbar implementation under external
KB-XFMD-030. Its previously inspected job transport must accept bounded stdin and
retain structured stdout on nonzero exit. The headless consumer and actual gh-sdp
trial prove the producer/transport; they do not prove a new native XFMD build.
Unsupported action IDs remain in saved preferences, hidden on old producers;
unavailable supported actions stay disabled with a reason. Unknown icons use a
generic fallback. No executable snippets are accepted from catalogue metadata.

[Session0011 evidence](../SDP/05--Implementation/SDPTool/Actions/Evidence.md) pins
the tested engine and wrapper separately. The Actions card is now KB-SDP-053;
KB-SDP-051 remains the published runnable-program card after identity reconciliation.

## Testing before main or release

For immediate consumer development, build this branch and select its executable
for one XFMD process. Run the matching XFMD working build; an older installed build
may not implement this override or dynamic tabs.

```sh
go -C /tmp/sdp-session11-integration/SDPTool build -o /tmp/sdptool-blueprint-dev ./cmd/sdptool
/tmp/sdptool-blueprint-dev /home/warloc/git/xfmd-sdl-navigation discover --json
XFMD_SDP_TOOL=/tmp/sdptool-blueprint-dev /path/to/working/xfmd
```

This tests the JSON integration without replacing stable gh-sdp or upgrading the
project installation. A project without retained bundles/assignments will not show
a populated blueprint tree; use the documented trial for representative fixtures.
Record source commit and binary digest with the test result.

GitHub CLI supports local extensions (`gh extension install .`, with a locally built
executable at the extension repository root) and tagged binary releases selected
with `--pin`. See the [official install contract](https://cli.github.com/manual/gh_extension_install).
A branch name alone does not build a Go extension remotely. Alpha/beta binary builds
can be published as explicitly selected tagged prereleases before a main merge;
that still requires actual release assets and separate publication authorization.

For the complete `gh sdp` path, the current gh-sdp README/main.go supports
`SDP_RELEASE` (signed descriptor), `SDP_TEST_KEY` (nonproduction public key) and
`SDP_CACHE_DIR` (separate cache). Prepare a branch-built signed test distribution,
then scope these settings to test invocations. The descriptor and signature/assets
must exist; setting SDP_RELEASE to a source branch or executable is not sufficient.
Test extension installation can use isolated GH_CONFIG_DIR and XDG_DATA_HOME.
Session0011 now verifies this path with a local test-signed distribution and
isolated extension/cache configuration. No normal installation or publication
was changed. See the linked evidence for the exact reproducible trial.
