# SDPTOOL-VER-T2-M1

Milestone: T2-M1. Card: KB-SDP-023. Sprint: SPR-SDP-0001.
Management event: EVT-PM-SDP-000014. Candidate: files in the introducing milestone commit;
Git records the exact candidate. No self-referential commit hash is invented.

## Delivery

Implemented read-only project recognition for explicit repo/SDP-area selection, strict bounded JSON registration, confined paths, capability declarations and referenced YAML installation facts.

## Verification

Go race tests pass for repo/SDP equivalence, explicit path CLI, no parent guessing, optional capabilities, invalid/duplicate/trailing JSON, null lists, unknown schemas/profiles, escaping/symlink paths and malformed/missing/unsupported referenced manifests. Source/metadata bytes are preserved; installed facts are declared, not falsely certified conforming.

## Boundaries

Producer checks do not establish native XFMD behavior. Deferred work remains in
KB-SDP-017 and the Sprint record. No installer or skills adoption is claimed.
