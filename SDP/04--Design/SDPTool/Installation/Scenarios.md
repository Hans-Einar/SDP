# Installation design scenarios and acceptance

These are design walkthroughs against [Contract.md](Contract.md) and the
[SDL model](../../../03--Architecture/SDPTool.design). They are not Go installer
test results. Parser validation proves modeled declarations, permitted messages,
participants and reply correlation; it does not execute the proposed workflow.

| Case | Input / trigger | Required observable result | SDL responsibility / illustrative scenario | Future implementation evidence |
| --- | --- | --- | --- | --- |
| IC01 clean install | Empty project, verified release | Preview creates no project files; reviewed apply creates neutral process, receipt and one Maintenance outcome | InstallationPlanPrepared; InstallationApplied | Temporary-project before/after hashes, reader discovery and history |
| IC02 known upgrade | Valid receipt, authenticated old and target inventories | Explicit changes with preserved project bytes; correct new identity | InstallationPlanPrepared / InstallationPlanner | Old/current/new fixture and repeat |
| IC03 manual XFMD adoption | Unknown version, inspected manual baseline and reviewed adoption manifest | No invented old release; preserve documents, root instructions and exact history prefix; explicit board relocation | ValidateAdoptionBaseline; ManualAdoptionTrial | Fresh copied XFMD inventory, conflicts and disposable full operation |
| IC04 drift | Change an input or target after preview | No apply of a stale plan; error 3; unchanged files before first mutation | InstallationDriftRejected | Change each bound input separately; compare target hashes |
| IC05 local edits | Modified managed file and project-owned content | Conflict unless reviewed managed refresh names it; project content preserved and backed up when relocated | InstallationBaselineInspector / InstallationExecutor | Distinct managed/project ownership fixtures |
| IC06 invalid input | Unsafe path, case collision, symlink, unsupported schema or unknown ownership | Fail before writes with explicit code/path; no shell hook execution | InstallationCoordinator / InstallationBaselineInspector | Host-native adversarial inputs, canonical-root checks |
| IC07 interruption | Process exits after backup, mutation or checkpoint | Pending operation visible; resume verifies exact bytes and publishes completion once | InstallationResumed / InstallationJournal | Inject exits at every boundary; verify no duplicate event/receipt |
| IC08 no change | Same verified release and unchanged installed content | No actions, no new report/events/timestamps | BuildInstallationPlan / PublishInstalledProcessFacts | Two-run byte equality |
| IC09 incompatible client | Binary lacks installation protocol/capability | Wrapper rejects before delegation, with actionable incompatibility | SdpToolBootstrapPrepared / GhSdpLauncher | Stub binaries with missing/wrong protocol; unchanged project |
| IC10 missing/untrusted artifact | Missing offline cache, bad signature/hash/platform or unknown key | Fail before execution/write; no silent latest-release fallback | ReleaseArtifactService / InstallationReleaseResolver | Corrupt/unsigned/offline fixtures, exact selected identity |
| IC11 finalization interruption | Receipt written but completion checkpoint missing | Discovery reports incomplete; reserved bytes/event IDs reused on resume | InstallationResumed / InstallationRecorder | Exit at each finalization step and check history prefix |
| IC12 competing or old operation | Lock held or legacy pending journal | No new mutation; report pending ID and required matching engine | InstallationJournal | Contention and explicit legacy-journal rejection |

## Refreshed XFMD baseline and walkthrough

Read-only inspection on 2026-09-27: xfmd-sdl-navigation at
b95a4bbd5ef43c9d7ed27270b58eeff2701a0cb8, branch
sprint/009/phase/057-favorites, clean. Board/history remain under
SDP/Agents/KanBan; no project/installed manifest or navigation registration is
present. The older KB-SDP-033 snapshots refer to earlier commits and must not
be silently reused as current input.

1. Observe the exact process area, root instructions, incoming Markdown references,
   file types and hashes. Record the unknown old version and explicit scan exclusions.
2. Select a real target descriptor (none is published by this work). Build a draft
   adoption map that preserves original ledger bytes, maps the board home and
   separates managed instructions from project-owned customization.
3. Run preview against a disposable copy, producing a root-bound plan. Review AGENTS
   preservation collisions, non-Markdown path references and every source removal.
4. In the future Go trial, apply only that copy's plan. Verify process discovery,
   receipt provenance, Maintenance reporting, byte-preserved project history and a
   no-change repeat; inject interruption/drift cases from the table.
5. Only after that evidence and explicit live selection, observe the live worktree
   again and make a new plan for its physical root. A temporary-copy plan cannot
   authorize the live project. Published release selection remains separate.

The existing PowerShell disposable-snapshot fixture is useful reference evidence
for preservation and ordering. It does not demonstrate the proposed Go engine,
signed distribution, new receipt schema or thin wrapper. The IPD review must keep
these differences visible rather than treating a legacy fixture pass as parity.

## Known model limits and review findings

- The five sequence scenarios show valid logical exchanges, not filesystem calls,
  loops over files, JSON schema validation, hashing or signature verification.
- InstallationDriftRejected is a pre-mutation failure path. It does not imply that
  every error is detected by the coordinator rather than a delegated unit.
- Direct CLI entry bypasses GhSdpLauncher; the illustrated process boundary is the
  wrapper path. Internal channels remain in-process Go calls.
- Database means persistent file/catalog storage. There is no SQL server requirement.
- JSON manifests/receipts have no fixed bit offsets; a VP10 packet is deliberately
  absent. No unrelated binary encoding is invented to make that viewpoint nonempty.
- Scoped SDL field names describe logical identity; nested record constraints and
  exact wire schemas remain the implementation contract's job.
- The initial payload/protocol draft uses bounded opaque record bytes. The first
  implementation milestone must freeze explicit JSON schemas and reject unknown
  fields before claiming protocol compatibility. No full schema coverage is claimed
  from those bytes appearing in a checked model.
- Independent review and owner acceptance have not occurred. This is a recorded
  author walkthrough with machine-checked model consistency, not invented approval.
