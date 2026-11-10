# R221-K1 — Linux publisher handoff

Owner 2026-11-11 selected an independent Linux key because the laptop is not
always available. This is trust/bootstrap preparation, not a published release.

## Preflight

Linux SDP is the sole local worktree, clean at 7f8bb4e before changes, on
sdp/release-2.2.1. gh-sdp is clean main at d5097b9. Active shared cards concern
other tooling/UI/governance work; inspected source/status shows no dirty overlap
with bootstrap trust, profile signing or R221. The owner's explicit instruction
owns the release signing prerequisite. Verdict green for this narrow scope.
Unknown work outside these observed checkouts is not claimed absent. Laptop
checkouts remain unchanged. Existing R221 package and review are preserved.

## Decision and boundary

Use a second compiled Ed25519 trust anchor, not laptop key copying, GitHub SSH
key reuse or a test-trust override. Linux owns the private key; only the public
key and a signed non-release probe are repository content. Retain the old key.
A newly built gh-sdp must carry the changed bootstrap before Linux-signed
descriptors can be selected. Installed clients do not acquire trust automatically.

## Evidence

Observed on Linux, 2026-11-11:
- Go 1.26.10, linux/amd64.
- `go test -race ./...` and `go vet ./...` in SDPTool/bootstrap: passed.
- `go test -race ./tools/profile ./install` and matching vet in SDPTool: passed.
- ProjectManagement validator: 73 cards, 48 management records, 729 events; passed.
- `git diff --check`: passed.
- Linux private-key generation used Go crypto/rand and Ed25519; 0600 file,
  0700 publisher directory. Only a public signed probe left the host.
- Linux key ID: `2c947457cffe015f870168cb332242c49606fdea378a13d875579544c2989c03`.
- Positive Linux proof and tampered-data, truncated-signature, wrong-known-key,
  unknown-key and invalid-provenance regressions passed.
Client build, exact signed package and independent review are pending.
No private key backup, release publication or end-user client upgrade is claimed.
