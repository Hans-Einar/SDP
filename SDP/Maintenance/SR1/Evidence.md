# Sessions release preparation evidence

## Installed consumer baseline

Read-only invocation from XFMD: `gh sdp --version --json` reports SDPTool 2.0.0,
schema sdptool/0.2, source d304261c90066a86b2d8ffcaa2115517ea339b05.
`gh sdp . discover --json` returns roots files/sdl/kanban/sdui. Resolving IDs
through navigation.nodes yields directory/tab/tab/tab. This confirms XFMD can
adopt dynamic tabs before the Sessions-enabled tool is published. No live project
files or globally installed extension were changed.

## Preparation boundary

Product implementation 2419853 is tested under PLAN-SDP-0015. SR1 freezes additive
2.1.0 notes and generated log, retains Framework 2.0.0 and the identical payload
inventory, adds the original signed 2.0.0 descriptor as a supported predecessor,
and selects the future 2.1.0 bootstrap URL. That URL is not yet a published asset.
No tag/publication or consumer installation is claimed. Exact candidate evidence
will be appended after the clean package, independent review and CI checks.

Independent review identified an inherited duplicate-node path when WalkDir
reports an unreadable directory twice. SR1 corrects that within the accepted
unavailable/unique-node contract and adds real permission regressions for Sessions
root/nested directories. The predecessor expectation adds the exact original
2.0.0 digest; all previous digests remain checked.

The corrected full SDPTool Go suite passes. Candidate packaging and publication
checks remain separate; this commit establishes the source for those checks.

## Exact candidate and executed checks

Source: `93517ad98cd188c0debeb1d0f3d36d123c6e4a3b` (clean detached checkout).
[Candidate.json](Candidate.json) records package hashes and runtime identity.
Only Linux/amd64 is packaged in this preparation, matching the current publisher
platform. No additional platform release asset is claimed.

- Full SDPTool suite and race suite pass; bootstrap race suite passes.
- 85 Toolkit tests and repository/management validators pass.
- Actual packaged-child signing/install test passes (test signer in that test).
- Production-signed descriptor rehearsals use the compiled trusted publisher key,
  without any test-key override. All four original signed predecessors (0.2.0,
  0.2.1, 1.0.0, 2.0.0) upgrade successfully, both fresh and with owner Sessions
  content. Ledger prefixes, owner files and inert legacy navigation files survive;
  repeated upgrade is a no-op. A fresh installation creates no navigation registry.
  See [Signed-upgrades.json](Signed-upgrades.json).
- A Git archive without Git metadata builds the development descriptor with the
  explicit source commit. The same eight upgrade/preservation/no-op cases and
  fresh install pass: [Archive-upgrades.json](Archive-upgrades.json). This is
  archive/development evidence, separate from production signature evidence.
- Installed `gh sdp` bootstraps the actual production-signed candidate via an
  explicit local SDP_RELEASE and isolated cache, reporting exact 2.1.0 identity.
  Its published default still selects 2.0.0; no global extension was updated.
- Read-only discovery of the actual XFMD worktree using the candidate returns
  Sessions as an available tab. Its installation facts exactly match installed
  2.0.0-tool discovery before the trial. No project migration or writes occurred.
- Independent review confirms all 64 payload sources and payload.json match
  v2.0.0 byte-for-byte. Framework remains 2.0.0.

PR: https://github.com/Hans-Einar/SDP/pull/51
Exact candidate CI passed both contracts and go-installation-linux at
https://github.com/Hans-Einar/SDP/actions/runs/36842826943 .
The package is staged locally; no release download exists yet. Final independent
preparation disposition is recorded below.

## Consumer handoff

`navigation.roots` contains node IDs. Resolve each ID in `navigation.nodes`, use
nodes with `kind: "tab"` as tabs and their label as captions; Files remains a
separate directory root. Do not hardcode three tab slots. Use each node's children
and typed open/generate targets. Installed 2.0.0 already supports this contract.

Once a new gh-sdp client selecting published SDP 2.1.0 is available, the normal
consumer command is `gh extension upgrade sdp`, followed by discovery/refresh.
No `gh sdp . upgrade` is required merely to browse an existing SDP/Sessions tree.
A project upgrade remains a separate operation if template/process changes are
wanted. Current project receipt and currently selected tool version may differ.

## Preparation limits

No publication authority is inferred from the request to prepare. No tag,
GitHub Release, main merge, client release, global extension update or live XFMD
upgrade was performed. Future publication must verify the selected exact candidate
and assets again, publish them, reconcile actual identities, and prepare/publish a
client with the corresponding default before advertising the normal update route.

## Independent review and closeout

[REV-SDPTOOL-SR1](Review.md) approves the exact product candidate for preparation,
with both initial findings resolved and no remaining material findings. It
independently checks the signed artifact, hashes/payloads, runtime identity and
successful CI. This later closeout commit changes records only; the packaged and
reviewed product remains the exact commit in Candidate.json, not an implied
package rebuilt from the records commit. SR1-M2 and MAINT-SDP-0012 are complete
within their preparation scope. Public release identities remain unset.
