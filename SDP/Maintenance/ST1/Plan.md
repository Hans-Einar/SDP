# MAINT-SDP-0006 — root SDPTool and installation ownership

| Field | Value |
| --- | --- |
| id | MAINT-SDP-0006 |
| project | SDP |
| state | completed |
| PlanType | MaintenancePlan |
| BranchPolicy | current |
| CommitPolicy | phase |
| source | Owner instruction 2026-09-26: root SDPTool, thin gh-sdp, withdraw competing Toolkit work |
| Systems | SDPTOOL, SDP |

## Outcome and authority

Move the existing Go SDPTool module from Toolkit/SDPTool to root SDPTool. It is
already the common SDP command entry point; future install/upgrade execution
belongs here, with gh-sdp as a thin client. Withdraw KB-SDP-018's separate Toolkit
audit and correct KB-SDP-033's competing-engine direction. Preserve completed
installation evidence and outstanding XFMD adoption. This direct owner instruction
authorizes the relocation and documentation changes, not a live installation,
release or new gh-sdp implementation.

## Baseline and scope

Baseline d105291757c9b1147ac57970fde35114cb207cd1; working branch
sdp/study-xfmd-install-adoption. Preserve unrelated untracked SDL/go/sourceinput.
SDPTool currently provides discovery, navigation, previews and viewer coordination.
Its installation.go reads facts and journals; install/upgrade mutation is still
implemented in PowerShell under Toolkit. Do not advertise a Go installer as done.
Toolkit becomes a legacy/transition area; retained schemas, profile payload,
validators and fixtures remain required until explicitly migrated. Do not delete
these inputs or break existing installation workflows.

Update module/import/build paths, repository consumers, current instructions and
relative Markdown links. Keep frozen fixtures, hash-pinned evidence and append-only
historical bytes unchanged. Historical path observations remain dated observations.
Use no compatibility symlink that would create a second apparent source home.

## Phase and milestone

| Phase | Milestone | Acceptance | State |
| --- | --- | --- | --- |
| ST1 | ST1-M1 | Root module builds/packages; Go tests and current document/management checks pass; consumers use root path; ownership and card dispositions agree | Completed |

## Git policy

Use the current working branch and one delivery commit for this single phase.
Preserve earlier plan branch commitments. Existing phase-push authorization applies;
no merge or release publication. The introducing commit records the tested diff.

## Verification and remaining work

Run Go race tests, packaged CLI version/help/discovery/navigation smoke checks,
project-management replay, documentation/frozen-record validation and applicable
Toolkit regression checks. Record actual commands/results here before completion.
Go install/upgrade implementation, distribution protocol and the real XFMD upgrade
remain future selected work under SDPTool and KB-SDP-033.

## Delivered evidence — ST1-M1

Tested candidate: baseline d105291757c9b1147ac57970fde35114cb207cd1 plus this
phase's working diff; the introducing commit records the final source. Go 1.27.1
linux/amd64 from the existing local toolchain. No tools were downloaded.

- Moved all 29 tracked SDPTool files; 22 retain identical bytes. Seven changed
  only for source/module/import/build locations, one repository-root test lookup,
  Markdown references and current ownership documentation. All testdata bytes
  are preserved. No compatibility directory/symlink remains at Toolkit/SDPTool.
- `go test -race ./...` from root SDPTool passes. CLI help now names the new
  contract path. The package build passes with its relocated linker symbol;
  version reports the actual baseline plus dirty suffix, not a release.
- Packaged CLI help/version, repository/SDP-area discovery (equal results), tree
  and saved .design preview succeed. Output is generated from the existing model;
  no generated viewpoint files were edited.
- `python3 -m unittest discover -s Toolkit/tests -p 'test_*.py' -v`: 106 tests,
  19 skipped when optional PowerShell/prebuilt consumer inputs were absent.
  The installed-consumer integration was then explicitly run with existing
  `/tmp/sk1-pwsh/pwsh` and the newly packaged binary: passes on a disposable
  project. Its initial invocation used a wrong test class name and was corrected
  to `test_process_install.ProcessInstall.test_installed_consumer_workflow`.
  This is compatibility evidence, not a live XFMD upgrade or full fault-matrix run.
- Toolkit validator, project-management replay, document/link and historical-byte
  checks pass. Existing 105 frozen records/prefixes and 574 generated outputs are
  unchanged. `git diff --check` passes.

KB-SDP-018 is canceled by owner instruction, not falsely completed. KB-SDP-033
retains XFMD adoption with SDPTool ownership; KB-SDP-017 stays ready for its
independent feature scope. REQ-SDPTOOL-007 records the unimplemented installation
target without adding it to the earlier SDL model as if designed or delivered.
Current links, CI package invocation, discovery probe path and module replacements
use root SDPTool. Dated historical paths and hash-pinned review inventories remain
unchanged; the retained Toolkit still supplies required inputs and the current
PowerShell installer.

No Go install/upgrade command, gh-sdp implementation, published release, global
binary replacement, native XFMD change or live upgrade is claimed. A concrete
Go installation migration plan under SDPTool remains the next installation task.
