# Manual XFMD upgrade to SDP 2.0.0

First update the native GitHub CLI extension. This changes the client and its
compiled default release; it does not upgrade the project's installed SDP area.

```sh
gh extension upgrade sdp
gh sdp --version --json
```

Expected engine facts: version 2.0.0, schema sdptool/0.2, revision
`d304261c90066a86b2d8ffcaa2115517ea339b05`. An explicit SDP_RELEASE environment
override selects another descriptor; retain it only if that is intentional.

Then preview from the active XFMD worktree, review the printed actions, and apply
that saved plan. Choose a new plan filename if this one already exists.

```sh
cd /home/warloc/git/xfmd-sdl-navigation
gh sdp . upgrade --plan-output /tmp/xfmd-sdp-2.0.0-plan.json
gh sdp . upgrade --apply /tmp/xfmd-sdp-2.0.0-plan.json
gh sdp . discover
gh sdp . tree
```

The preview writes only the explicitly selected plan outside the project. Apply
performs the reviewed upgrade and records signed release 2.0.0 in installed facts.
If preview reports conflicts or a pending operation, follow its diagnostics;
do not bypass the plan's guards. Source files and existing project documentation
are preserved. No model-registration command is needed.

Machine clients use `gh sdp . discover --json` and read its derived `inventory`,
`sources` and inline `navigation`. Source identities are returned by discovery,
not fixed aliases. They must support sdptool/0.2 and explicitly request --json.
SDPTool supplies a snapshot; native XFMD integration owns buffering, filesystem
watching, debounce and manual refresh. This release does not implement those
native viewer behaviors. Detail generation remains on demand.

No live XFMD upgrade or global extension update was performed during Session 0003.
