# GIP delivery, rollout proposal and legacy retirement boundary

## Delivered candidate

The Go installation workflow is implemented. The delivered product source is
SDP commit 07335b0 (GIP-4 preparation-failure correction included), with the thin
client product fca8480 and bootstrap module fb79727. gh-sdp closeout is df59b7b;
its later commits contain review/evidence only. Final evidence is recorded in
[Evidence.md](Evidence.md) and the [XFMD trial](evidence/Final-XFMD-trial.json).

| Candidate | SHA-256 |
| --- | --- |
| Packaged SDPTool, clean isolated build | 0d984b548f47b6ecb3b5049a7b2c24268822cf1032fdaeb373bb7574745fc18e |
| Packaged gh-sdp, reviewed product | f7fa971d950591512919c65ede5d009d65cce2d2d4c25c33b91ecfd640af50f8 |

Installation defaults to preview. Apply consumes a saved, root-bound plan;
resume uses the retained journal. The client executes a verified local Go binary,
not GitHub-hosted code or a PowerShell engine. No global gh-sdp extension was
installed by this work: actual gh routing was exercised with isolated config/data.

## Scenario evidence

| Cases | Evidence delivered |
| --- | --- |
| IC01/02 | Clean install, known descriptor upgrade and receipt 3.0 discovery; actual packaged signed-test install |
| IC03 | Fresh XFMD b95a4bb copied baseline; equal direct/gh plans, apply through both, explicit adoption mapping |
| IC04/05 | Root/input/project drift and managed edits reject; named refresh and instruction preservation collision cases |
| IC06 | Closed strict records, bounds, path confinement/case/symlink/nonregular rejection, including native FIFO regression |
| IC07/11 | Child process exits at preparation and every backup/write/checkpoint/completion boundary for clean, known and manual fixtures; full XFMD copy resumed through actual gh |
| IC08 | Repeat preview/apply/resume preserves bytes and IDs; no new completion event/report |
| IC09/10 | Signed test descriptor, wrong signer/hash/platform/protocol/path/size, private cache corruption and offline rejection; actual incompatible packaged binary |
| IC12 | Concurrent Go installer lock and legacy pending journal block; no legacy journal reinterpretation |

Forward recovery is proven at the explicit process-exit boundaries in the tests.
It is not a claim of whole-tree rollback, power-loss durability or native testing
on other operating systems. HTTP and probe timeouts are implemented; timing tests
and a production network publisher deployment were not part of this delivery.
Independent review covers the thin client and its bootstrap integration. The core
engine checks and full XFMD trial are author-run verification, not independent
approval of every installer path.

## Fresh XFMD findings

The observed live worktree remained clean at b95a4bbd5ef43c9d7ed27270b58eeff2701a0cb8.
The trial plans 138 file actions and preserves 194 paths. The original KanBan ledger
is an exact prefix of the new shared history. Project instructions move into
AGENTS-project.md with only the inspected KanBan board link rebased; the byte-exact
original remains in the operation backup. Project-owned phase documents remain
initialize-if-missing, with relocation/link changes separately recorded.

The [legacy comparison](evidence/Final-legacy-comparison.json) matches all 138 ordinary
actions. Board/navigation JSON values are equal; only key ordering differs. The
legacy facts 2.0 write is intentionally excluded because Go receipt 3.0 is a
separate journaled finalization with descriptor/provenance identity. No unrelated
legacy golden fixture was rewritten to force this result.

The manifest embedded in the evidence is the **actual disposable-copy manifest**.
Its temporary root and plan digest are evidence, not a live upgrade plan. The
AllowReferenceWarnings disposition is specific to this copy experiment. Before a
live operation, inspect non-Markdown path references and excluded dependency/build
areas; neither this installer nor the trial claims semantic rewrites of source code.

## Concrete next delivery proposal — not executed

1. Select a production process-release identity and Ed25519 publisher key, using
   the versioning/release workflow. Build/pin the candidate payload and binary;
   ship the trusted public key in reviewed client/standalone distributions. The
   current development descriptor and ephemeral test key must not be relabeled
   as a production release. A published default stable catalog is not available.
2. Initially advertise only the verified Linux platform. Native Windows/macOS
   install/apply/recovery and executable-cache tests are prerequisites for claims
   there. Select binary distribution and an immutable descriptor URL explicitly.
3. Re-inspect the **current** xfmd-sdl-navigation worktree. Generate a new adoption
   manifest and preview tied to that physical root and the selected signed release.
   Compare every action, the AGENTS instruction handoff, native consumer paths and
   KanBan/history mapping. The historical manifest must not be reused blindly.
4. With the verified release configured, the intended review command is:

   ```sh
   gh sdp upgrade --manifest /tmp/xfmd-upgrade.yaml \
     --plan-output /tmp/xfmd-plan.json --json
   ```

   Only after explicit live selection run `gh sdp upgrade --apply /tmp/xfmd-plan.json`.
   Retain the returned journal/backup/report identity; use `--resume ID` for a
   journaled incomplete operation. A caught pre-publication failure has no resumable
   ID and may be retried after resolving the reported input/filesystem problem.
5. Verify live project discovery and native XFMD navigation with XFMD's owner/agent.
   That application integration is not included in a copied filesystem trial.
   Begin the deferred XFMD SDL-modeling pilot as a separately bounded assignment.

KB-SDP-033 holds this concrete rollout review. The broader SDL-modeling pilot is
explicitly deferred; it is neither implemented by GIP nor silently discarded.
The implementation plan may close without granting publication or live mutation.

## Legacy ownership and retirement map

| Retained area | Why retained / removal prerequisite |
| --- | --- |
| Toolkit/scripts/Install-SDP.ps1 and Process-Install.ps1 | Recovery of schema 2.0 pending journals and existing consumer workflows. Do not mix engines during an operation. Retire only after consumer inventory and pending-operation migration/recovery are accounted for. |
| Toolkit profiles, payload and Template sources | Shared authored inputs still used by the explicit Go files-only profile. Move them only with a source ownership/path migration, not by copying a second maintained collection. |
| Toolkit schema 1.0/2.0 and frozen conformance fixtures | Compatibility/read evidence. Keep historical bytes; these are not competing new APIs. |
| Legacy Windows/PowerShell CI | Native compatibility coverage not yet replaced by Go Windows execution. Removing it would hide a known verification gap. |
| Go installation/profile code | Sole new install/upgrade policy and execution owner. |
| Shared bootstrap module | Sole bootstrap implementation consumed by gh-sdp and SDPTool; no migration policy. |
| gh-sdp legacy Study | Historical accepted research; current SPS-003 records explicitly supersede its client-owned apply allocation. |

No legacy execution file is deleted by GIP: current evidence does not cover all
of those remaining consumers/platforms. This is a deliberate retirement boundary,
not a hidden production dependency of the new Go user path.
