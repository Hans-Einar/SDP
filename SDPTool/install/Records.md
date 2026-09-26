# Installation records — protocol 1

The executable schemas are the exported structs, strict decoder and validators in
records.go. All fields are required, including explicit nulls for nullable values.
The decoder rejects duplicate and case-colliding keys, unknown keys, YAML aliases,
anchors, explicit tags, noninteger numbers and multiple documents. JSON is the
canonical saved representation; YAML is also accepted for authored adoption input.
Canonical bytes have sorted keys, UTF-8, no HTML escaping and one final LF.

| Record | Schema | Purpose |
| --- | --- | --- |
| Descriptor | sdp-release-descriptor/1 | Publisher inventory and inline payload, explicit retirements, supported predecessor digests and platform binary hashes |
| Adoption | sdp-adoption/1 | Physical project root, exact observed snapshot, target descriptor digest, preserving moves and named managed refresh exceptions |
| Receipt | 3.0 | Release/digest/profile/provenance lookup facts and completing operation; no authoritative inventory |
| Plan | sdp-install-plan/1 | Exact inputs, observations, actions, conflicts and deterministic digest |
| Journal | sdp-install-journal/1 | Exact plan, reserved finalization bytes, step progress and integrity |

Input carries the exact source bytes and absolute source path. Null previous/adoption
means no input, not a zero digest. Development artifacts require explicit consent;
signed input verification is added in GIP-3. Receipt schema 3.0 readers precede writers;
1.0/2.0 readers remain available. Legacy receipts need an explicit adoption manifest
until their original authoritative inventory is available: a version label is not
sufficient ownership evidence.

Descriptor/adoption limits are 16 MiB; plan/journal limits are 64 MiB. Per-file
payload/observation limit is 64 MiB, total 512 MiB, and 100,000 observed entries.
Inline payloads must also fit the containing record limit. This initial transport
is intentionally bounded rather than promising all limit maxima simultaneously.
Release destinations are confined to managed root instructions/skills and SDP.
Engine-owned receipts, journals, board, navigation and histories cannot be replaced
by a release payload. Missing templates initialize only. Paths use portable slash
notation; symlinks, special files, escapes, case collisions and Windows device
names are rejected. Inspection includes process paths and incoming Markdown links;
excluded dependency/build trees remain explicitly outside link-rewrite coverage.

The engine supports Go 1.26 or newer. Native platform verification, trusted release
keys and published catalog selection are separate from schema availability.
