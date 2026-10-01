# SR2 publication evidence

## SDP 2.1.0

Owner publication authority: direct request on 2026-10-01. Selected source:
`93517ad98cd188c0debeb1d0f3d36d123c6e4a3b`. Annotated tag object `bfb14a6d90d9af78fe314a927548d8f1cca6888a` peels to that exact commit.
GitHub release: https://github.com/Hans-Einar/SDP/releases/tag/v2.1.0
Actual publication timestamp: 2026-10-01T13:24:49Z (stable, not draft).

Before publication, all prepared package hashes, clean exact checkout, executable
identity, generated release logs and absent local/remote tag were rechecked.
No product source changed after the independent SR1 approval. Exact-candidate CI
and subsequent preparation-record CI passed. [Publication.json](Publication.json)
records real asset identities; every downloaded asset matches the SR1 candidate.

The downloaded executable successfully previewed and applied a fresh temporary
installation using the public HTTPS descriptor and production trust. Its receipt
reports 2.1.0, signed provenance and the exact descriptor/source identity:
[Public-install.json](Public-install.json). This is a temporary trial only.

Descriptor SHA-256: `d8436398ff6ed1cf74b7754106792075203dc3745f91d19c56450951b7302dde`.
Binary SHA-256: `691c021751cefb35c7f92c902545cd6c5d71ccadfcde07ae507a6ce026c8a823`.

## Client and notification handoff

Client gh-sdp 0.2.1 is being prepared in its own bounded SPS-008 release with
immutable bootstrap from the selected source. Its release is recorded separately
once actual. Main branches, global installed extension and XFMD are unchanged.

GitHub CLI already checks for an executed extension's updates once per 24 hours
and writes the notice to stderr. It can be disabled with
GH_NO_EXTENSION_UPDATE_NOTIFIER. No custom notification code is needed.
Source: https://cli.github.com/manual/gh_help_environment
That variable was not enabled in this agent's environment. A prior same-day check
can delay the owner's first notification; it does not mean the release is absent.
We avoid invoking the owner's global gh sdp during publication, preserving their
opportunity to test. gh checks the client release, not its pinned SDPTool version.

Once the paired client is published, the owner may run `gh sdp --version` to use
the still-installed version and observe any notice, then `gh extension upgrade sdp`
and `gh sdp --version` to verify SDPTool 2.1.0. Existing Sessions browsing does not
require `gh sdp . upgrade`; native XFMD dynamic tabs and a refresh are still needed.
