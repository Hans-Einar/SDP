# Traceability

Project traceability normally contains:

- `CurrentIndex.yaml` — actual current release and development coordinates
- `Relations.yaml` — links among requirements, design, work, review, verification,
  migrations and releases
- `Ledger.ndjson` — append-only event history, one valid JSON object per line

Release events conform to `Toolkit/schemas/release-event.schema.json`. Only append
state transitions that occurred. In particular, do not record
`release-tag-created` or `release-published` before the real tag or GitHub Release
exists. Corrections are new events; historical lines are never rewritten.

Release IDs use `REL-X.Y.Z`. Fix IDs use `FIX-X.Y.Z-NNN` or, for emergency work,
`HOTFIX-X.Y.Z-NNN`.

## Management cutover — 2026-09-25

[ProjectManagement](../ProjectManagement/README.md) now owns scheduling, card,
Scrum/Sprint, Maintenance, CodeReview and Refactor lifecycle. Keep pre-cutover
management events here as historical bytes; do not duplicate new ones here.
Record actual system .design/code changes and their evidence with references to
management work IDs and event IDs. New IDs put system first (SDL-DES-0001,
SDUI-REQ-0001, SDPTOOL-VER-0001); project-wide tooling uses SDP explicitly.
Old IDs are not renamed. KB-SDP-011 still owns nine legacy Toolkit ID failures;
this local convention is not a claim that all distributed schemas accept it.

A Sprint touching system code/design is represented here through realization
records that reference its SPR-SDP ID and exact management event, system/model
revision, commit and verification. Sprint planning/start/end belongs only in the
management ledger. One completed Sprint does not mark every Feature verified.
