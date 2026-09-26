# SDPTOOL-VER-T4-M1

Milestone: T4-M1. Card: KB-SDP-026. Sprint: SPR-SDP-0001.
Management event: EVT-PM-SDP-000021. Candidate: files in the introducing milestone commit;
Git records the exact candidate. No self-referential commit hash is invented.

## Delivery

Published executable consumer fixtures/examples and completed producer review with strict metadata/default-model fixes, bounded inventory and verified error preservation.

## Verification

Go 1.27.1 go test -race ./... and go vet ./... passed. The compiled CLI consumer verifies discovery, graph references, selection bundle contents, stale rejection, fresh regeneration, unsupported operations and configured viewer arguments. Real repository tree/Concept1/SVG trials succeeded; details and limits are in REVIEW-SDP-0001. Final checks passed: 35 cards, five management records, three lineage operations and 212 events; all six Sprint members are completed. Management tests and 15 lineage negatives pass. Navigation schemas validate; 2,397 local links/130 fragments resolve, with 105 frozen records/prefixes and 574 generated artifacts preserved. Pre-Sprint management bytes remain an exact prefix. git diff --check passes.

## Boundaries

Producer checks do not establish native XFMD behavior. Deferred work remains in
KB-SDP-017 and the Sprint record. No installer or skills adoption is claimed.
