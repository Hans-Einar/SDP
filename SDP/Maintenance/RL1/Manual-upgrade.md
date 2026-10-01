# Manual XFMD upgrade handoff

Not executable as a normal published upgrade yet: SDP 1.1.0 is proposed and
unpublished. [ReleaseChecklist](ReleaseChecklist.md) owns readiness. Latest actual
release is SDP 1.0.0; running the default now will select that existing release.

There are two different updates:

- `gh extension upgrade sdp` updates the gh-sdp executable. It does not migrate
  the current project's SDP directory. Its default SDP release is compiled/pinned.
- `gh sdp upgrade` previews a project upgrade. Applying a saved plan is explicit.
  `gh update sdp` is not the project upgrade command.

After SDP 1.1.0 publication and signed asset verification, an existing gh-sdp client
can select the new descriptor explicitly, even before a client-default patch:

```sh
cd /home/warloc/git/xfmd-sdl-navigation
export SDP_RELEASE=https://github.com/Hans-Einar/SDP/releases/download/v1.1.0/sdp-release.json
sdp_upgrade_plan_dir=$(mktemp -d)
gh sdp upgrade --plan-output "$sdp_upgrade_plan_dir/plan.json" --json
```

Read the preview. The current installation must match a supported signed predecessor.
Expected new files include SDP/Sessions/README.md (if missing), Session-template.md
and SDP/Framework/ReleaseChecklist.md. Managed agent guidance refreshes; authored
Sessions, source models and populated project documents stay unchanged. No custom
adoption manifest is needed for a valid known release receipt. If preview reports
conflicts, investigate them; do not bypass them or rewrite the receipt.

When the preview is acceptable:

```sh
gh sdp upgrade --apply "$sdp_upgrade_plan_dir/plan.json" --json
gh sdp --version
gh sdp upgrade --json
```

The last preview should be noChange=true. Inspect the installed receipt for the
exact target version/digest and sdp.sessions.manual.v1. Check the generated
Maintenance entry and retained ledger prefix. An explicit SDP_RELEASE can be
unset when an independently verified client default selects the same release.
These commands are a future handoff, not evidence they were run in XFMD.
