# Shared runtime implementation contract

Owner selects the single-update-path design and Activity4 execution. Canonical
SDPTool source owns all tool behavior. gh-sdp distributes the exact canonical
engine binary under the GitHub CLI extension asset name. It must not ship a
second runtime launcher that chooses/downloads another engine. Prefer verifying
and republishing the signed canonical engine bytes over maintaining a duplicate
main/source build in the extension repository. Package provenance links both
repositories while engine version/commit remain identical across entrypoints.

## Runtime registration

Add `sdptool setup` / `gh sdp setup`: explicit initial registration, never a
side effect of discovery or read-only commands. Linux public entrypoint is
$HOME/.local/bin/sdptool, a symbolic link to the actual current executable path.
The target path survives ordinary gh binary upgrades. Do not hard-code GitHub's
extension-directory layout. The command must identify source/entrypoint/version,
be idempotent, refuse overwriting unmanaged files/foreign links, diagnose PATH
shadowing and absent PATH entries, and avoid partial success claims. The GUI
resolves the documented public path independently of shell PATH. Explicit developer
overrides remain explicit and observable, never automatic old-engine fallback.
No network, gh process, credentials or project mutation is needed for registration
or local browsing. Setup invoked from transient/deleted/unsupported sources must
fail clearly. Version and schema remain separate; no bootstrap selection is hidden
behind `--version` or another ordinary operation.

Expose registration through the existing operation/action authority, not a second
CLI-only implementation. Keep project install/upgrade preview/apply semantics
and signed descriptor verification in their existing domain. Release selection
for project framework data must not select a different executing engine. Remove
the old gh runtime bootstrap route as part of this migration; retain shared
signature/descriptor facilities still used by the installer with clear ownership.

## Native consumer and helper inventory

XFMD executes the registered SDPTool directly, validates protocol/capabilities,
and diagnoses missing/incompatible registration with actionable setup guidance.
Remove private engine lookup, gh fallback and engine shipping from the next RPM.
Do not silently execute an arbitrary older PATH installation. Keep source revision,
rendering, history, cancellation and existing explicit developer-test overrides.

Inventory SDL/SDUI helper calls. All tool operations must be reachable through
SDPTool/gh sdp; required managed helpers share its versioned distribution and update
owner. Application-specific declared programs remain project-owned applications,
not a second SDP updater. Return a real helper gap before deleting a supported
workflow. Canonical SDK/engine stays in SDP; native rendering stays in XFMD.

## Qualification and delivery

Publish/verify the shared engine and extension before successor XFMD removal of
its bundled copy. Use actual isolated gh install/upgrade and clean target-platform
RPM checks, with the shared engine established by the supported registration.
Verify offline native browsing with gh absent, auth unavailable, and network disabled;
not merely mocked command strings. Compare hashes/source versions, not only labels.
Retain signed publisher/provenance and project-preservation guarantees. Independent
review must cover migration boundaries and public verification. No main merge or
silent removal of arbitrary user executables is authorized.

## Selected SDUI dispatch

`sdptool sdui preview -source PATH -entry ENTRY -revision SHA -output DIR` and
`sdui check -source PATH -entry ENTRY -revision SHA` reuse the SDUI-owned preview
library and existing JSON contracts. `sdui run -entry ENTRY -revision SHA
-watch=true PATH` starts the native prototype. The Fyne executable is embedded in
the production engine; extraction is local, hash-verified, atomic and private.
It is neither separately updated nor fetched on invocation. Development builds
without the desktop payload fail clearly on run; release qualification must use
the bundled payload. Explicit native development overrides remain observable.

Setup and these host operations use the common action registry with explicitly
project-independent scope. Existing actions continue requiring project resolution;
adding host scope must not weaken their context checks.

The SDUI source service exposes typed `preview.Options` / `preview.Generate`.
Both SDPTool actions and CLI call that service directly. The retained development
`sdui-preview` executable only parses flags and serializes the same result; no
new core operation reconstructs a legacy CLI command line to access the service.
