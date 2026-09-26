# GIP implementation evidence

Linux x86_64; Go 1.27.1. Evidence identifies HEAD plus the milestone diff;
the introducing commit binds that diff. Live XFMD remains untouched.

## GIP-1-M1

2026-09-26T23:10:42Z. Executable closed records, five canonical golden fixtures, bounded strict JSON/YAML decoding, portable path/root inspection and receipt 3.0 reader delivered. Go package suite passes; management validator passes. No apply engine or signed distribution claimed.

## GIP-1-M2

2026-09-26T23:16:10Z. Direct packaged install/upgrade previews delivered for explicit development artifacts, with deterministic plans, manual adoption, preserving relocation/link rebasing, old/current/new managed comparison and no-change detection. Go race suite passes. Packaged preview produced six actions and left its project empty. Signed resolution and apply remain later milestones; therefore the broader SDL InstallationPlanSlice stays planned.

## GIP-2-M1

2026-09-26T23:20:20Z. Root-bound saved-plan apply, advisory lock, drift revalidation, byte backups, atomic replacements and journaled receipt/history/report publication delivered. Clean install and known upgrade tests pass with race detection. Actual packaged apply completed and discovery reads schema 3.0. Interruption matrix remains GIP-2-M2.

## GIP-2-M2

2026-09-26T23:24:02Z. Process-exit fault matrix passes every preparation/backup/write/checkpoint/completion boundary for clean, known-upgrade and manual-adoption cases. Reserved finalization bytes/IDs and ledger prefixes remain stable; edited targets and legacy pending operations block. Full Go race suite passes (installation package 48.3 seconds). Model apply/recovery activities updated and review regenerated through SDL. No power-loss or rollback guarantee.

## GIP-3-M1

2026-09-26T23:31:03Z. Shared stdlib bootstrap verifies Ed25519 test descriptors, platform binary digest/size and advertised protocol, with private immutable cache and offline checks. Engine independently verifies signed input, records test-signed versus production provenance and resolves previous descriptors by digest. Go profile explicitly reuses files-only legacy inventory without PowerShell prerequisites. Packaged signed install passes; bootstrap race tests pass. No production keys, stable catalog or published release exists.

### GIP-3 integration review and hardening

Independent gh-sdp review found a FIFO-open hang in shared bootstrap input reading.
Both bootstrap and engine now reject nonregular objects before opening, retaining
post-open bounds checks. A native FIFO regression and signed negative tests for
protocol, platform, binary size/path/digest and missing protocol advertisement pass.
HTTP/probe timeouts exist but are not claimed as separately exercised timing tests.

Full engine race suite passes after management-history predecessor/card validation,
cache publication hardening and incoming Markdown ancestor consistency fixes.
The updated version test explicitly expects receipt 3.0 and the now-implemented
installation protocol; this is the authorized compatibility addition, not a waived
legacy obligation. Clean project AGENTS.md is preserved automatically; unequal
AGENTS-project.md collisions still block. A fresh XFMD copied baseline preview/apply
also succeeded during integration; its reproducible paired trial belongs to GIP-4.

## GIP-3-M2

2026-09-26T23:42:57Z. Thin gh-sdp client implemented in its own repository with shared bootstrap pinned to fb79727. Independent review REV-SPS-003-002 approves client fca8480 plus clean engine fb79727 on Linux after resolving FIFO and evidence findings. Actual gh sdp/direct plan and apply parity, offline behavior, race and vet pass. Client implementation does not duplicate installation policy. Wider engine and XFMD evidence remain GIP-4.

## GIP-4-M1

2026-09-26T23:48:18Z. Fresh XFMD b95a4bb baseline verified unchanged. Real direct sdptool and isolated gh sdp produce identical plans and successfully apply the same root-bound plan on a disposable copy: 138 actions, 194 preserved paths. Full-copy process exit followed by actual gh resume and repeat/no-change pass. Original history prefix and root instruction backup are exact; only the inspected KanBan link changes in AGENTS-project.md. Legacy comparison matches all 138 actions, with JSON key order normalized for board/navigation and receipt 3.0 intentionally separate.

### GIP-4 closeout hardening

A caught failure before the initial journal is published now returns preparation
failure (exit 4), with no invented recoverable operation ID. Only that operation's
temporary file and empty setup directories are removed; no reviewed file actions
have run. Limit errors retain exit 2. Existing operation identities cannot be
silently overwritten. The injected publication-failure regression proves unchanged
project content and successful fresh retry; normal apply/recovery checks still pass.
This is distinct from failures after publication, which retain their journal and
forward-resume identity. Final exact-candidate runs follow this source correction.

## Final candidate verification and reproduction

Product source: 07335b0. Packaged from a clean detached checkout, then the temporary
worktree was removed. The unrelated untracked SDL/go/sourceinput draft stayed in
the working checkout and was never inspected, staged or modified. That draft had
made earlier Go VCS metadata dirty, so those binaries were superseded explicitly
by the clean build rather than misrepresented as exact commits.

Final engine and client identities, review scope and rollout limitations are in
[Rollout-and-Retirement.md](Rollout-and-Retirement.md). Runtime code is unchanged by
the subsequent documentation/evidence closeout commit.

Commands (Go 1.27.1, Linux amd64):

```sh
go -C SDPTool/bootstrap test -race -count=1 ./...
go -C SDPTool/bootstrap vet ./...
SDP_TEST_BINARY=/tmp/gip4-exact-package/sdptool \
  go -C SDPTool test -race -count=1 ./...
go -C SDPTool vet ./...
SDP_TEST_XFMD=/home/warloc/git/xfmd-sdl-navigation \
SDP_TEST_BINARY=/tmp/gip4-exact-package/sdptool \
GH_SDP_BINARY=/tmp/gip-gh-sdp \
SDP_TEST_ARTIFACT=/tmp/gip4-exact-package/release.json \
SDP_TEST_EVIDENCE=/tmp/gip4-final-xfmd.json \
  go -C SDPTool test ./install -run '^TestXFMDDisposableAdoption$' -count=1 -v
python3 SDP/ProjectManagement/validate.py
python3 Toolkit/scripts/validate_sdp.py --mode toolkit
python3 SDP/Maintenance/L1/verify_documents.py
python3 SDP/04--Design/SDPTool/Installation/review.py --sdl /tmp/ipd-sdl --check
```

Use local tool paths appropriate to the checkout when reproducing. The profile
fixture is generated with `go -C SDPTool run ./tools/profile --repo ..` and explicit
source-commit, packaged binary and output parameters. It is unsigned development
input; the test harness generates an ephemeral test signature and isolates GitHub
CLI config, extension data and distribution cache. It writes only a disposable
project copy. Tests without the required candidate/source environment variables
skip that external fixture; an ordinary package test run is not a claim it ran.

[Final engine log](evidence/Final-engine-tests.txt): all packages pass with race
detection; installation suite 56.96 seconds, including the process-exit matrix and
actual packaged signed-test install. Bootstrap race/vet and engine vet also pass.
[Final XFMD log](evidence/Final-XFMD-trial.txt): 30.52 seconds, real direct/gh entry
points, full-copy exit/resume and no-change repeat. The
[machine record](evidence/Final-XFMD-trial.json) contains exact source observations,
adoption manifest, action hashes, candidate hashes and assertions. Temporary roots
in that record are provenance, not reusable live destinations.

The retained PowerShell engine was used only for a read-only comparison capture:
`Install-SDP.ps1 -ProjectRoot /home/warloc/git/xfmd-sdl-navigation -ProfileArtifact
Toolkit/profiles/five-phase.artifact.json -PlanJson -ForceManagedFiles`.
No PowerShell command applied to live XFMD. The
[comparison record](evidence/Final-legacy-comparison.json) pins that capture hash
and the Go evidence hash. Reproduce the comparison with
SDPTool/conformance/compare_legacy.py; it accepts captured JSON and runs no installer.
The older M1 trial/comparison records are preserved as historical evidence.

All current canonical management/schema/document checks pass. Generated review
files retain their generator's final blank line; whitespace checks use the previously
documented blank-at-EOF exception for generated output only. No generated Markdown
or Mermaid was manually edited. No independent full-core review, native Windows/macOS
acceptance, release publication, merge or live project mutation is claimed.

## GIP-4-M2

2026-09-26T23:54:50Z. Exact clean candidate 07335b0 and reviewed client fca8480 pass final packaged/race/XFMD verification. Delivered rollout proposal and legacy retirement map; retained only compatibility/recovery and shared inputs whose consumers are not yet replaced. Production keys/default release, native Windows/macOS and live upgrade remain unselected; broader XFMD SDL pilot explicitly deferred. All selected GIP milestones delivered.
