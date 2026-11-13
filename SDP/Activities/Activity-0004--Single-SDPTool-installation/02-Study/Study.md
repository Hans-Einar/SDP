# One shared SDPTool runtime — initial study

## Owner decision and scope

The owner rejects the private SDPTool shipped inside XFMD0.10 and requires one
installation/update route through `gh extension install/upgrade sdp`. All SDP
operations should be reachable through `gh sdp`; native consumers must invoke
`sdptool` without running gh, authenticating with GitHub or accessing the network
for local browsing. This supersedes the private-bundle distribution choice for
the next design. It does not rewrite the delivered Activity3/RPM evidence.

KB-SDP-063 and ACTIVITY-SDP-0004 capture this need. This pass is architecture
analysis, not product implementation, a complete implementation plan or publication.

## Observed facts

- SDP source9db0d3d exposes the actual engine as `SDPTool.RunWithInput`; its CLI
  main is already a thin process entrypoint. The current bootstrap instead pins
  an independent signed release and obtains another executable in a private cache.
- gh-sdp sourceb177e14 uses that bootstrap and executes the resolved child. It
  does not currently contain the engine. Its default remains SDP2.2.1.
- XFMD source140be2b selects explicit XFMD_SDP_TOOL, then its private helper,
  then `gh sdp`. Both the private helper and gh fallback conflict with the new
  desired normal runtime path.
- GitHub CLI binary extension install downloads an executable and writes its
  manifest; upgrade reuses that installer. It does not execute extension-provided
  post-install hooks. Therefore it cannot automatically run our PATH-registration
  code merely because the extension was installed/upgraded.

Primary upstream evidence, inspected2026-11-13:
[install manual](https://cli.github.com/manual/gh_extension_install) and
[Manager.installBin / upgradeBinExtension](https://github.com/cli/cli/blob/trunk/pkg/cmd/extension/manager.go).

## Recommended architecture

Ship the canonical SDPTool engine as the gh-sdp extension executable, calling the
same engine entrypoint rather than downloading another executable at runtime.
Expose a stable standalone entrypoint to that exact file: on Linux a user-local
`sdptool` symbolic link. `gh sdp`, direct `sdptool`, and XFMD then execute one
installed engine. Naming the file gh-sdp does not require gh to execute it.
The source remains owned by SDP; the extension repository is the distribution
adapter and must not fork engine behavior.

Because gh provides no post-install hook, recommend one explicit initial local
registration operation (proposed new command `gh sdp setup`) to create the public
entrypoint using the executable's actual installed location. This command does
not exist yet. After registration, ordinary extension upgrades replace the same
payload and the link continues to resolve it; XFMD needs no gh invocation.
Avoid silently installing aliases as a side effect of read-only commands.

Define one public runtime-location contract, including GUI environments without
shell PATH and Windows/macOS equivalents. Do not hard-code gh's private storage
layout into XFMD. On Linux, a documented user-local public path can be resolved
directly; setup must diagnose conflicting existing PATH executables and avoid
overwriting unmanaged installations. The former /usr/libexec/xfmd/sdptool is
removed by the next selected RPM migration, not by editing an already-published RPM.

Other required SDP helpers/resources belong to one coherent, versioned distribution
with the same update owner. Their exact packaging/activation needs an inventory;
this study does not claim all SDL/SDUI native tools fit into one Go executable.
Local browsing must be complete offline after installation/registration. Any
materialization of embedded helpers must use the current distribution's exact
bytes, not perform independent version selection or background downloads.

Use capability/protocol compatibility for XFMD, with a useful missing/incompatible
runtime diagnostic. Make binary path, release version and commit observable.
Keep schema identifiers such as sdptool/0.2 separate from product version display.
Project framework upgrades remain explicit project operations through `gh sdp`;
a runtime upgrade must not mutate every project installation.

## Alternatives and migration boundary

Retaining today's bootstrap plus a symlink into its versioned cache would still
permit gh's selected engine and XFMD's linked engine to diverge until another
bootstrap call. It also permits an offline consumer to miss a just-downloaded
extension update. That is unsuitable for the requested normal update behavior.
A separate SDPTool RPM creates another update authority and is also unsuitable.

Before implementation: model the selected runtime/registration boundary and
inventory required helpers, then define the complete migration and release plan.
Publish/qualify the shared runtime before replacing XFMD's bundled-helper package.
Retire normal bootstrap engine-selection/cache behavior in the same migration,
without treating development overrides or legacy executables as automatic fallback.
Do not delete arbitrary existing user binaries or historical signed releases.

Acceptance must cover: first installation; registration and PATH collisions;
extension upgrade while XFMD exists; identical path/version across all entrypoints;
XFMD browsing with gh unavailable, authentication removed and network disabled;
missing/incompatible runtime; failed/interrupted upgrade and recovery; removal of
the old package helper; preserved project documents/settings; and explicit support
boundaries for non-Linux platforms. Package signatures/provenance must remain
verified under the revised distribution, not assumed from old bootstrap evidence.

## Checkpoint

Known checkouts: SDP9db0d3d, gh-sdpb177e14 and isolated XFMD140be2b. Product trees
were not edited. SDP's existing runtime lock and XFMD's existing .deps cache link
are preserved. New branch sdp/shared-tool-installation contains only need/study
records; no Activity is active and no implementation admission is selected.
Unknown remote/concurrent work is outside this read-only preflight's coverage.
