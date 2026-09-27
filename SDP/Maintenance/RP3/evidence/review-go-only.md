# RP3 Go-only independent review

Review date: 2026-09-28. Role: independent SDP Reviewer.
Candidate: 705df7e3d550629f13b8e9b0044fee8f6b1b3bed.
Compared against the previously reviewed 56919a19edcbd6e672c8a26661d87d16828be4f8
candidate and inspected surrounding current Go installation and validation code.

## Authority and disposition

The owner's revised scope in MAINT-SDP-0010 removes all PowerShell execution,
including archived executable copies and shell-driven tests, and selects SDP
1.0.0. This is an authorized retirement of source-distributed entry points, not
a claim that the retired engines satisfy their former runtime contracts.
The earlier review.md applies only to the superseded 0.2.2 candidate. No merge,
release or XFMD upgrade is established by either source review.

Bounded approval of the Go-only source transition. No unresolved material source
finding was identified. Current Go installation, preservation, trust and recovery
obligations remain in force. Release completion still requires the runtime and
publication evidence below.

## Inspected changes and results

- The native SDPTool/profiles/payload.json files array is exactly equal to all
  61 entries of the preceding Toolkit/profiles/five-phase.json inventory,
  including ordering, source, destination and ownership. Every listed source
  exists. This move removes legacy artifact metadata from current authoring
  without dropping a payload or changing its installation policy.
- Descriptor construction now uses the existing strict install.Decode function
  for the inventory and checks sdp-payload-inventory/1. Inspection confirmed
  duplicate/case-colliding keys, missing/unknown fields and inappropriate value
  shapes are rejected by the shared decoder. Relative source validation and
  bounded regular-file reads remain; SafeAbsolute checks path components for
  symlinks. Output still passes ValidateDescriptor, which constrains destination
  scope, ownership, duplicates/case collisions, reserved engine paths, content
  hashes and size limits. Production signing still requires clean exact HEAD,
  a packaged binary and a trusted publisher key.
- Current template completeness and reproducibility tests exercise the real
  authoring tree. Independently reran the completeness, reproducibility and
  exact published-predecessor tests on this committed candidate:
  `GOMAXPROCS=2 /home/warloc/.local/share/sdp-toolchains/go1.27.1/go/bin/go test -p 1 ./tools/profile -run 'TestCurrentTemplateInventoryComplete|TestCurrentDescriptorReproducible|TestProfileRetainsExplicitReleasePredecessor' -count=1`
  from SDPTool. Result: pass. Concurrency was bounded for the owner's laptop.
- Inspected the source tree and executable Python/Go/shell/CI paths: no remaining
  .ps1/.psm1/.psd1 files or calls to a retired engine were found. The L1 document
  verifier's exception is confined to the explicitly retired archived script;
  it requires that file to remain absent. Removed SK1 scripts no longer invoke
  the retired engine; historical implementation remains reachable through Git.
- CI retains both Go module race suites, actual packaged-child signed install,
  and descriptor building from current authoring sources. The descriptor command
  resolves repository and packaged-binary paths consistently with `go -C SDPTool`.
  The file-inventory gate rejects reintroduced PowerShell file extensions.
- Existing Go tests remain for known upgrade/local edits, unknown adoption,
  root-instruction preservation, drift rejection, locking, pending legacy
  journals, signature/cache failures and process-exit recovery. The recovery
  matrix covers clean, known-upgrade and manual-adoption transactions across
  backup/write/journal boundaries and verifies repeated recovery is stable.
  Retiring tests of the old engine does not remove these current-engine checks.
- The retained Python validators still inspect historical install-v1 scenario
  and expected JSON data. Documentation now distinguishes those data checks
  from current Go runtime evidence. Current skill portability checks consume
  the native payload inventory.
- Both published 0.2.0/0.2.1 predecessor digests remain explicitly declared;
  bootstrap selects immutable v1.0.0. The installed protocol, receipt/profile
  identities and ownership rules remain unchanged. Release notes disclose
  removal of source-distributed installer/recovery entry points and explicitly
  retain blocking for pending legacy operations. The major-version choice is
  recorded as the owner's decision rather than inferred implementation authority.

## Minor documentation follow-up

SDPTool/install/Records.md still uses `--release 0.2.0` in its production-signing
example. Update this to 1.0.0 or an explicit selected-version placeholder before
the release documentation is finalized. This example does not alter the actual
compiled default or descriptor compatibility.

## Evidence limits and release prerequisites

The reviewer inspected coordinator logs showing all root SDPTool Go packages
passing without race instrumentation and 85 Python tests passing. The coordinator
reports the separate bootstrap suite also passes. Only the three targeted profile
tests above were independently rerun by this reviewer; neither full-suite race
success nor CI completion is claimed here.

Attach passing Go-only CI, including race and actual packaged-child tests, to the
final candidate. A test-signed fixture is not production provenance. Build/sign
the exact clean merged source, rehearse supported predecessor upgrades with the
actual 1.0.0 descriptor, and retain evidence of project-content preservation and
repeat no-op behavior. Record actual immutable publication, gh-sdp client identity,
XFMD receipt and navigation results before closing RP3. Earlier 0.2.2 package,
preview or conformance results are superseded evidence, not final release proof.

This review wrote only this report and made no implementation changes.
