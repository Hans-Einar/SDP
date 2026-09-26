# SDPTOOL-VER-P0-M1

Milestone: P0-M1. Card: KB-SDP-021. Sprint: SPR-SDP-0001.
Management event: EVT-PM-SDP-000011. Candidate: files in the introducing milestone commit;
Git records the exact candidate. No self-referential commit hash is invented.

## Delivery

Go saved-file preview implemented using existing SDL projection, query and bundle services, with compact default selection, structured diagnostics, cancellation and revision checks, guarded publication and caller-owned resources.

## Verification

Go 1.27.1: go test -race ./... passed (fresh-source regeneration, stale/error preservation, malformed/profile-invalid input, missing renderer, invalid selection, containment including symlinks, output conflicts and cancellation). A prebuilt facade generated VP02-roots Mermaid and SVG from SDPTool.design through the real mmdr executable. No GUI visual acceptance is claimed. Board, link and whitespace checks are run before commit.

## Boundaries

Producer checks do not establish native XFMD behavior. Deferred work remains in
KB-SDP-017 and the Sprint record. No installer or skills adoption is claimed.
