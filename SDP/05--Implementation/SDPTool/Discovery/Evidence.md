# Session 0003 evidence

## DS1-M1 — source discovery

Candidate: c78831b plus DS1-M1 commit contents; Go 1.27.1, Linux amd64.
GOMAXPROCS=2 and -p 2. Root SDPTool and bootstrap module tests pass. Added real
filesystem tests cover composed Systems, fragments, invalid/future SDL, broken
SDUI, duplicate filenames, source changes, rename/removal, stale registration
inertness, contained paths and source limits. Compiled consumer journey preserves
revision-bound selections and output guards using source-derived model IDs.
Existing semantic service tests explicitly construct Go inventories; they no
longer pretend that a persisted register supplies discovery. New compiled consumer
and filesystem tests own the actual discovery boundary.

Built /tmp/sdptool-session3. Read-only discovery against the actual XFMD checkout
succeeds without registration edits. /tmp/sdp-session3-xfmd-discovery.json contains
the temporary result; summary remains reproducible from the CLI. No XFMD file was
written and no generated navigation sidecar is required.

Independent review, installer removal, final packaging and publication remain
pending. This evidence is not a release-readiness claim.
