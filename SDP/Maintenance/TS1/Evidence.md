# TS1 verification and handoff

The TS1-M3 section below supersedes the earlier authoring-path arrangement.
Earlier installation evidence retains the paths used when it was captured.

Candidate: baseline 8cdb7ba plus uncommitted TS1 edits on sdp/mvp1-source-ui.
The [installation record](installation-check.json) identifies the full baseline,
development descriptor digest and exact input payload hashes. No release was made.

## Checks performed

- `python3 -m unittest Toolkit.tests.test_process_profile -v`: two checks pass;
  deterministic/neutral generated profile and invalid-input rejection.
- `python3 Toolkit/scripts/validate_sdp.py --repo .`: passes. The legacy inventory
  explicitly excludes its new development-only directory README.
- `python3 SDP/ProjectManagement/validate.py`: validates current history and records.
- Existing L1 document check, with Markdown enumeration limited to Git-indexed
  files plus new TS1 documents to exclude unrelated node_modules: 105 frozen
  records/prefixes, 574 generated outputs, 2,754 links and 134 fragments pass.
  All 55 local Markdown links in the fresh installed project resolve.
- `git diff --check`: passes.
- Go 1.27.1 `go build ./cmd/sdptool` and `go run ./tools/profile` succeed using
  temporary caches, current sources and `--release development-ts1`. No PowerShell.
- Actual Go CLI install preview/apply in a disposable project: completed; installed
  SDL README/AGENTS exactly match source bytes.
- Edit the installed SDL README, then preview/apply the same development descriptor:
  no-change, and project-specific prose preserved.
- Reconstruct a development predecessor descriptor from HEAD's inventory/payloads,
  without changing historical bytes. Install it, edit its project README, then
  upgrade with an explicit development successor/predecessor digest: completed,
  four file actions (two managed guidance updates, two missing SDL guides), custom
  project overview preserved. This is a development rehearsal, not a published
  predecessor compatibility declaration or consumer rollout.

## Reproduction

From SDPTool, use Go 1.26+ to build `./cmd/sdptool` and run `./tools/profile`
with `--repo .. --output <new-file> --source-commit <HEAD> --release development-ts1`.
Use a disposable project and writable XDG_CACHE_HOME, then run the built executable:

```text
sdptool <temporary-project> install --artifact <descriptor> --allow-unreleased --plan-output <new-plan> --json
sdptool <temporary-project> install --apply <plan> --json
```

For a predecessor rehearsal, preserve each HEAD inventory file's exact bytes and
hash in a development descriptor. The development successor must explicitly list
that descriptor's SHA256 in upgradesFrom. Supply `--previous-artifact` on upgrade
preview, then apply the saved plan. Never alter a production descriptor/signature.

## Outcome and limits

TS1-M1/M2 documentation and installation checks are delivered in the working tree.
The new source tree contains guidance only, without fabricated model files. Current
five-phase templates now describe source/document authority and useful phase roles;
retained seven-phase templates remain legacy inputs with a clear directory notice.
Managed instructions route agents to the source guide when present.

Existing project templates are preserved during upgrade; review and reconcile them
explicitly to adopt the new prose. Existing SDL model migration remains KB-SDP-020;
experimental syntax/source-set support remains KB-SDL-005. No parser, runtime,
application, published release or consumer project was changed.

Git metadata is read-only in this session. The required TS1 phase commit/push is
pending; the Maintenance record remains active for that concrete closeout.
Unrelated SDL/go/sourceinput and node/package files are excluded from this work.

## TS1-M3 — canonical template root

Owner follow-up relocates the current template from Template/profiles/five-phase
into Template/sdp-root and removes Template/profiles. Both legacy roots now live
under Template/legacy/install-v1. The archive version names the retained install-v1
interface, not a historical complete release. All 22 legacy payload files retain
exact bytes; see [hash inventory](legacy-template-hashes.json), whose keys are
relative to the archive root. The previous turn's development-only directory
notice is replaced by the archive-level README.

Four neutral seeds needed by the current profile are copied into current sdp-root;
current installation has no dependency on legacy templates. The legacy inventory
uses the archived copies. Its exclusion list distinguishes the current templates
from that old interface. The old project-root template had no Go-profile consumer,
so only its archived copy remains. Managed root AGENTS still comes from Toolkit.

Validation:

- Current Go descriptor rebuilt with identical development identity and baseline:
  **byte-identical** to the pre-relocation descriptor, including every payload,
  hash, ownership and installed destination. Previous fresh-install and preservation
  checks therefore cover the exact same payload after this source-only relocation.
- `python3 -m unittest Toolkit.tests.test_process_profile Toolkit.tests.test_validate_sdp -q`:
  all 83 checks pass, including repository contracts and deterministic profile build.
- `python3 Toolkit/conformance/install-v1/run_conformance.py --validate-only`:
  all 19 scenario packages validate. Updated expected plans change exactly 113
  source-path fields, with no action/destination/policy changes. These are live
  conformance expectations, not edits to frozen historical evidence. The example
  plan and one PowerShell contract-test input receive matching source paths.
  Actual PowerShell scenario execution and independent conformance review were
  not performed in this session; do not claim runtime replay or release approval.
- L1 document check using Git-indexed extant Markdown plus new work documents:
  105 frozen records/prefixes, 574 generated outputs, 2,752 links and 134 fragments
  pass. Unrelated node_modules remains outside this check.
- Management validation and whitespace checks pass. Git commit/push remains pending
  under the same read-only Git metadata restriction; no release or consumer update.
