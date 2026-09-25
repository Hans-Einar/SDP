# SK1 execution evidence

Source baseline aabb359; M1 content/contract commit 7df5afe. This is Maintenance
MAINT-SDP-0002, authorized before the next installer Scrum. No consumer repository
or global agent configuration was upgraded. The unrelated SDL sourceinput draft
is excluded. The Toolkit remains unreleased.

## SK1-M2 — canonical adoption and installation

Thirteen prepared roles and two shared references now live under root Skills/.
Both older maintained source collections are removed. Thirteen relative symlinks
in .agents/skills expose the project collection; installed copies retain their
existing .codex/skills destination. The manifest inventories every installed role
and reference; README is development-only. Native name/description and nested
string metadata coexist with legacy installed-v1 validation. Current canonical
and capability-v2 installations must satisfy native metadata. Skill versions,
AGENTS contract 2.0.0, examples and explicit capabilities agree.

The same-Toolkit-version preflight fails without mutation if declared skill
versions/inventory differ and ForceManagedFiles is absent. The existing forced
path backs up changed managed files and preserves project-owned content. This
avoids generated version facts getting ahead of installed instructions. It does
not perform generic process-directory upgrades.

Commands and results:

- All thirteen skill-creator quick_validate checks pass; installed references
  and source inventory are also exercised by Toolkit/tests/test_skills.py.
- Python unittest discovery: 82 tests, 81 pass; the single repository-contract
  test reports exactly the [38-error baseline](toolkit-baseline.txt). No new
  validation errors; old consuming-project fixtures retain their exact bytes.
- PowerShell 7.6.6 full Toolkit/tests/Install-SDP.Tests.ps1 passes on Linux.
  A stale required-exclusion test now removes the actual SDP exclusion instead
  of the already-retired root document name, restoring that negative test.
- [Upgrade probe](verify_skill_upgrade.py) passes clean/repeat installation,
  complete baseline managed-file restoration, unforced mutation-free failure,
  forced upgrade with customization backup, new router/reference creation,
  AGENTS replacement, project-content preservation, repeat idempotence and
  external-link rejection. The old input comes from aabb359 Git objects.
- The original seventeen conformance scenarios pass against the regenerated,
  independently inspected plan authorities. The two added skill-transition
  cases also pass the non-regenerating comparison run: nineteen total scenarios.
  Their failure/force outcomes received independent review as well.

The PowerShell binary was unpacked only under /tmp from Microsoft's official
PowerShell GitHub release v7.6.6 for testing; no system package/configuration was
changed. No Windows-specific runtime result is claimed by this Linux run.

Independent review of the actual M2 diff and eleven changed expected plans found
two evidence/documentation issues: historical draft links pointed at adopted
content, and the first upgrade fixture retained new roles/AGENTS. Both were fixed
and the reviewer rechecked them. Historical references now pin aabb359; the probe
restores all legacy managed copies and removes newly introduced destinations.
The reviewer approved the corrected scope; owner acceptance remains separate.

## SK1-M3 — catalog and actual use

[Catalog record](catalog-evidence.json) and [reproduction helper](verify_catalog.py)
identify the committed M2 source and every adopted skill/reference hash. Codex CLI
0.156.1 discovers all thirteen enabled roles without duplicates in four contexts:
repository root, nested cwd, a separate Git worktree and an actual fresh Toolkit
installation using .codex/skills. A nested independent Git repository does not
inherit the parent's collection. No user configuration was changed. The original
four-case discovery probe also reproduces its retained baseline from Git bytes.

[Behavior record](behavior-evaluation.json) retains fixture inputs, final hashes,
diffs, prompts and reported read paths. Four independent fresh contexts received
realistic isolated tasks, without this thread's conclusions or expected answers:

| Trial | Actual selected skills and result |
| --- | --- |
| Bounded formatter fix | sdp → Worker/Verifier; changes only formatter, tests and the existing Maintenance record; three tests pass; missing labels fixed, other IDs preserved |
| Read-only evidence audit | sdp → Auditor/Verifier; identifies stale candidate evidence and unrun application checks; no edits or false demand for a Sprint/installed Framework |
| Owner-choice advice | sdp → Architect; recommends the supported temporary-file route, distinguishes advice from approval, identifies remaining decisions; no implementation |
| Trivial arithmetic | No skill loaded or lifecycle work started; correct answer |

Explicit source paths were supplied for the first three trials. The tests prove
observed explicit loading/routing and bounded behavior; they do not prove every
possible implicit prompt match. Read paths are agent-reported, supported by the
returned task artifacts, rather than an exported native-host telemetry trace.
Other unexercised roles have metadata/reference/catalog evidence, not claimed
individual behavioral certification. Sustained compliance, other IDE/desktop
hosts and Windows installation remain outside these measured results.

The coordinating session also read the adopted router, Worker and Traceability
instructions for this Maintenance. Independent M2 review approved the corrected
implementation/evidence and both old and new conformance authorities. This is
agent review, not owner acceptance or publication of a release.

## Final structural checks

Management validation: 35 cards, seven management records, three lineage operations
and 227 events. Four management test groups and 15 lineage negative cases pass.
All pre-SK1 ledger bytes remain exact prefixes; 105 frozen records/prefixes and
574 generated artifacts are preserved. Documentation check resolves 2,380 local
file links and 130 fragments. Full Toolkit errors remain byte-for-byte the
recorded baseline; git diff --check passes. The phase has separate M1/M2/M3
commits and retains the owner's branch stack.
