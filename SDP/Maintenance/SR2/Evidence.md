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

Client gh-sdp 0.2.1 is published under its own bounded SPS-008 release:
https://github.com/Hans-Einar/gh-sdp/releases/tag/v0.2.1
Actual publication timestamp: 2026-10-01T13:35:16Z. Annotated tag object
5a93c3b8a51e62c9c3558227b3cf166f6f90f505 peels to the reviewed candidate
315ec1de9dbe9aee990baae47a1ac2a7a07a0702. It selects the immutable bootstrap
v0.0.0-20261001092620-93517ad98cd1. Root downloaded all three assets and checked
their bytes against the independently reviewed client package. Binary SHA-256:
230440b0bbeb5af33f98b91bd227b26f78e3deefba9e23dee1e1508e75061dca.

External client evidence: VER-SPS-008 and REV-SPS-008-001 in Hans-Einar/gh-sdp.
The review independently checks the exact candidate, full packaged race tests,
isolated gh routing, vet, immutable module and public default from a fresh cache.
No findings remain. The client has no CI workflow; its native gate passes and no
client CI result is claimed. Upstream CI is separately successful.
Main branches, global installed extension and XFMD remain unchanged.

GitHub CLI already checks for an executed extension's updates once per 24 hours
and writes the notice to stderr. It can be disabled with
GH_NO_EXTENSION_UPDATE_NOTIFIER. No custom notification code is needed.
Source: https://cli.github.com/manual/gh_help_environment
That variable was not enabled in this agent's environment. A prior same-day check
can delay the owner's first notification; it does not mean the release is absent.
We avoid invoking the owner's global gh sdp during publication, preserving their
opportunity to test. gh checks the client release, not its pinned SDPTool version.

The paired client is now published. The owner may run `gh sdp --version` to use
the still-installed version and observe any notice, then `gh extension upgrade sdp`
and `gh sdp --version` to verify SDPTool 2.1.0. Existing Sessions browsing does not
require `gh sdp . upgrade`; native XFMD dynamic tabs and a refresh are still needed.

## Independent client publication closeout

External REV-SPS-008-002 independently verified the actual annotated tag,
public release, all downloaded assets, fresh-cache production default and matching
engine/descriptor. Its disposition is approved/pass with no remediation findings.
Root inspected that report and independently checked downloaded client bytes and
latest-release selection. The client coordinator owns linking its report and
final administrative closeout in Hans-Einar/gh-sdp. Both publications are complete.

The prepared candidate/review files under SR1 retain their historical preparation
status. This SR2 record is the authority for subsequent actual publication; no
released note sections, signed descriptors or historical ledger bytes were edited.
