# TS1 release and XFMD upgrade handoff

Current update, 2026-09-28: the owner has removed the access restriction and
authorized execution. Earlier restrictions below are historical; RP3 owns release
and consumer rollout.

Owner requests committing this delivery so a release and XFMD upgrade can follow.
Current session permissions prohibit Git metadata writes and writes to gh-sdp or
xfmd-sdl-navigation. Do not bypass those restrictions. No commit, merge, tag,
production descriptor or consumer upgrade is claimed.

## Manifest chain

Toolkit/profiles/five-phase.json is the current explicit payload inventory and
now reads Template/sdp-root. SDPTool/profiles/five-phase.json selects that inventory
and the supported predecessor descriptor digests. Release building emits immutable
sdp-release.json with paths, payloads, ownership and hashes, plus its signature.
The installed receipt identifies the exact release/descriptor; it does not supply
a mutable project-owned authoritative inventory.

Read-only observation of xfmd-sdl-navigation: signed release 0.2.0, receipt schema
3.0, descriptor edc0c72101a437c6e12c40a081ef59ae41cf0db32bcaedb73824ec48495aaee5.
The earlier manual bootstrap has already been adopted. A fresh custom adoption
manifest is not indicated by this receipt. Source configuration also now admits
published 0.2.1 descriptor 66590e8e967ede6b36d8fa45cdbee1cd80f69505cca698b0e4a9bd960842735a;
its hash was checked against the locally retained published descriptor. A new
signed upgrade rehearsal is still required before release.

## Remaining delivery

1. Commit the selected TS1 files on sdp/mvp1-source-ui, excluding unrelated
   SDL/go/sourceinput and node/package files. Review branch-wide scope before
   integration; the branch also contains the earlier MVP1 pilot and backlog work.
2. Select release identities and prepare coherent notes/defaults/manifest records.
   Run exact-candidate release gates, including the legacy conformance requirements.
   Build/sign from a clean committed candidate and preserve old published assets.
3. Publish the SDP release, then update and release gh-sdp's bootstrap dependency
   so its default selects that release. Existing clients can also select an explicit
   release; their compiled default does not follow latest automatically.
4. Reinspect XFMD immediately before upgrade. At this observation it has uncommitted
   KanBan README/history changes and a new KB019; preserve that work. Generate a
   reviewed plan using the signed target and the receipt's verified predecessor,
   then apply the exact plan and verify the receipt, new SDL guidance and navigation.

Template relocation changes publisher source paths, not installed destinations.
The new SDP/SDL guide files initialize when missing. Existing phase READMEs and
other project-owned prose are preserved: a normal upgrade does not rewrite them
or move existing .design models. Reconcile that documentation explicitly if the
owner wants the full new prose in XFMD. No application code change is selected.

## Required session access

Writable SDP-vNow working files and the real shared Git directory
/home/warloc/git/SDP/.git are needed to commit. Client release and consumer upgrade
also require their respective checkout/Git write access and the normal authorized
release/network facilities. Current approval policy is never, so this session
cannot request escalation or change its own permission profile.
