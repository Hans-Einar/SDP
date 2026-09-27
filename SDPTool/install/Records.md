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

## Apply and lock placement

Apply re-derives the reviewed plan under an advisory lock keyed by physical project
root in the private platform cache (`sdptool/locks`). Keeping the lock outside the
project allows a clean target to be re-inspected before creating SDP itself; lock
ownership ends automatically if the process exits. Journals/backups remain under
SDP/.sdp-operations. The Linux implementation is exercised here; other native
platform support must not be inferred from compilation. A legacy pending journal
blocks Go apply and requires a separately assessed recovery or migration. The
current distribution does not ship or invoke the retired engine.

Operations use root-confined Go filesystem handles for writes and atomic replacement
from operation-owned temporary files. Backups are verified before each mutation.
Apply may create parent directories; no recursive directory deletion is used.
Project documents and unrelated executable source files are never treated as
release-owned. No rollback or filesystem power-loss guarantee is advertised.

## Distribution selection and development fixtures

`--release PATH-OR-HTTPS-URL` selects an exact descriptor; `SDP_RELEASE` supplies
an override for both entry points. The compiled default selects an immutable release descriptor URL; see
bootstrap.DefaultRelease for the candidate's exact version. Detached signatures are JSON with keyId (SHA-256 of public key) and
base64 signature. Production keys must ship in the reviewed bootstrap module;
the selected publisher public key ships with v0.2.0. `--test-key FILE` / `SDP_TEST_KEY` explicitly enables
an Ed25519 public test key (base64 file), producing **test-signed** provenance,
never signed production provenance. The saved proof stores this distinction.

`--offline` / `SDP_OFFLINE=true` permits only verified cached distribution bytes.
`SDP_CACHE_DIR` selects a private cache directory. Immutable selector entries bind
an exact descriptor digest; a changed selector is an error, not silent latest
selection. Engine inputs are reverified independently of client bootstrap.
Saved apply/resume do not fetch or change input selection. Previous descriptors
are recovered by receipt digest from the verified external cache; development
operators can also supply `--previous-artifact` with explicit unreleased consent.

Build the engine-neutral development fixture with `go run ./tools/profile`, using
an explicit source commit, output path and packaged binary. profiles/five-phase.json selects the native profiles/payload.json inventory.
Go builds descriptors directly; no legacy artifact generator or installer is used.
For authorized publication, --release 1.0.0 --sign-key PRIVATE_FILE requires
a clean exact Git HEAD, a binary, and the compiled-in publisher key. It writes a
detached signature but does not publish. Keep the private key outside repositories.
