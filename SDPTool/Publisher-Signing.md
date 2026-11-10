# SDP publisher signing

## Current ownership and trust

Owner decision, 2026-11-11: Linux is the always-available build/signing host.
The intermittently available laptop is not a required release dependency.
GitHub authentication (Git SSH/API credentials) is separate from distribution
authentication. No private signing key belongs in Git, a release asset, a CI log
or the client binary.

The bootstrap accepts these independent Ed25519 publishers:

| Host | Public-key SHA-256 / signature keyId |
| --- | --- |
| Existing laptop | bfc4838a7fbe7e7e3e84f86b9b435e69e6e503f4caf01a92319147e9ab701618 |
| Linux build host | 2c947457cffe015f870168cb332242c49606fdea378a13d875579544c2989c03 |

Linux public key (base64):
`3c4v3py8MxlsZ0zGp0T2Wtmcl2rXQ4f4WEM/YBShBiw=`

Private key on Linux: `~/.config/sdp/publisher/ed25519-linux-2026-11.key`.
It is base64-encoded Go Ed25519 private-key material, in a regular 0600 file
under a 0700 publisher directory. It is not an OpenSSH or OpenPGP key.
The existing laptop key has not been copied or revoked. No backup/recovery copy
of the new key has been established by this work.

## Build and sign

Use the existing clean-source package/profile workflow; pass the Linux key to
`go run ./tools/profile --sign-key "$HOME/.config/sdp/publisher/ed25519-linux-2026-11.key"`
from SDPTool together with the selected root, version, source, binary and output
arguments documented by `go run ./tools/profile --help`. This is a signing-option
example, not a complete release command. The profile builder checks exact clean
HEAD, signs exact descriptor bytes, verifies the compiled trust anchor, and
refuses to overwrite a signature. Verify all assets and the signed consumer path
before publication. Never reuse a descriptor or binary from before the trust edit
as evidence of the new executable.

## Client rollout and compatibility

Both bootstrap and SDPTool must contain the new trust anchor. gh-sdp imports a
pinned bootstrap module: update that immutable dependency and build a new client.
The initial bridge candidate retained SDP2.2.0. The owner subsequently authorized
publication and installation on Linux, Beinerud and the laptop, plus Linux project
upgrades. The selected gh-sdp0.2.3 release now defaults to SDP2.2.1 with both keys.
Publish and verify the engine distribution first, then the paired client; consumers
must update the client before selecting Linux-signed releases.

Old clients intentionally reject a Linux-only signature. Publishing a public key
on GitHub, setting up SSH, or supplying a descriptor cannot teach them new trust.
Do not use SDP_TEST_KEY to bypass this production migration. Existing signed
releases remain immutable and continue to verify through the retained old key.
No multi-signature format, remote signer, revocation or automatic key discovery
has been introduced.

See the [installation contract](../SDP/04--Design/SDPTool/Installation/Contract.md)
and [R221 signing handoff](../SDP/Maintenance/R221/Linux-Signing.md).
