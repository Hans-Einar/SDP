# Cross-repository patterns

## Evidence boundary

- **OBSERVED:** The study considered 35 owner-visible repositories and produced
  reports for the 17 with qualifying commits. The comparison matrix partitions
  the corpus into nine default-tree project SDP implementations, the canonical
  SDP method repository, four repositories with SDP only in non-default/open
  evidence, and three with no SDP evidence.
- **OBSERVED:** Default-branch truth is not the same as recent activity. HSX,
  HEOS, map_tracker and RadarData contain important current SDP evidence only on
  non-default branches or open PRs. Their practices remain in-flight evidence.
- **OWNER DIRECTION:** Issue #5 fixes the strategic destination: evolutionary
  Feature/Refactor-owned work, GitHub Issues as assignment boundaries, bounded
  Codex Masters and a supervisory Steering Group. Frequency in the historical
  corpus does not override that direction.
- **INFERENCE:** The useful question is therefore not “which folder layout won
  by majority?” It is which observed mechanisms improved boundedness,
  architecture, review, evidence and recovery, and which mechanisms repeatedly
  drifted or added ceremony.

The exact row-level comparison is in [EvidenceMatrix.md](EvidenceMatrix.md).

## Adoption and evolution clusters

| Cluster | Repositories | Observed implication |
|---|---|---|
| Issue-authoritative delivery accepted on default | ActionCam, Lyndata, TerrainAnalyzer | Strong evidence for Issue scope, exact baseline, branch/PR, evidence and stop boundaries; still vulnerable to oversized Issues and stale copied state. |
| First-class Feature delivery accepted on default | Lyndata, TerrainAnalyzer | Feature can own later capability and local design evolution; mappings to Issues and Slices are not reliably 1:1. |
| First-class Refactor execution accepted on default | GrassPhenology, SharedUI | Refactor is a useful peer owner for structural change and local Study/design; `weight_app_flutter` and HSX strengthen this only in unmerged evidence. |
| Sprint-first pre-manifest execution | GrassPhenology, LogClassifier, map_tracker, SharedUI, weight_app_flutter and others | Bounded Slice contracts often work, but mandatory Sprint/Iteration wrappers and mixed work types create semantic and ID overhead. |
| Branch-only or incomplete adoption | HEOS, HSX, map_tracker, RadarData | Migration must distinguish default truth, in-flight evidence and accepted practice. |
| No SDP | LogParser, python, tplink | A proportional bootstrap is necessary; retroactive reconstruction would invent evidence. |

## Pattern 1 — bounded contracts and handoffs help; copied live status does not

- **OBSERVED:** ActionCam, GrassPhenology, LogClassifier, Lyndata, map_tracker,
  SDP-Analyzer, SharedUI, TerrainAnalyzer and weight_app_flutter preserve enough
  goal, scope, invariants, evidence, residual risk and next-step context for a
  fresh session to recover work.
- **OBSERVED:** ActionCam, GrassPhenology, gh-sdp, Lyndata, SharedUI,
  TerrainAnalyzer, weight_app_flutter and SDP-Analyzer also contain stale or
  contradictory Handoff/current-state prose.
- **INFERENCE:** Handoff is valuable as a timestamped narrative of rationale,
  risk, limitations and next decision. It is unreliable as another mutable
  registry of Issue/PR/head/check/status fields.

## Pattern 2 — adversarial review materially improves outcomes

- **OBSERVED:** At least seven default-tree reports record concrete
  review-caused rework: gh-sdp, GrassPhenology, LogClassifier, Lyndata,
  SDP-Analyzer, TerrainAnalyzer and weight_app_flutter. Findings include lost
  updates, persistence/cancellation defects, stale targets, architecture
  contradictions and candidate-specific evidence errors.
- **OBSERVED:** CI or local tests were often green before review found material
  defects. Lyndata's SLC-007 also shows the reverse failure: owner closure and an
  “approved” index coexist with unresolved Medium findings.
- **OBSERVED:** Native GitHub reviewer identity is usually absent. The same
  `Hans-Einar` account publishes most repository and GitHub evidence, so fresh
  context is a process attestation rather than externally proven identity.
- **INFERENCE:** The durable contract must separate test/verification,
  independent review and Steering/owner acceptance. None substitutes for the
  others, and a changed head invalidates the prior review disposition.

## Pattern 3 — horizontal foundation plus later local evolution works

- **OBSERVED:** GrassPhenology, LogClassifier, Lyndata, SharedUI,
  TerrainAnalyzer, weight_app_flutter and SDP-Analyzer revise requirements or
  design after implementation evidence. ActionCam also evolves requirements as
  physical facts change.
- **OBSERVED:** The initial horizontal boundaries in Lyndata, SharedUI,
  TerrainAnalyzer, weight_app_flutter and SDP-Analyzer constrain later vertical
  delivery without requiring the future product to be fully designed at
  bootstrap.
- **OBSERVED:** Later decisions are sometimes written only into Sprint/index
  prose or silently mutate global phase documents. That weakens the audit trail
  for why accepted architecture changed.
- **INFERENCE:** Later Feature/Refactor-local Study, requirement and design
  deltas need a durable home plus explicit refinement/supersession links to the
  project foundation.

## Pattern 4 — the Slice contract is valuable, but the term is overloaded

- **OBSERVED:** Strong bounded Slice contracts occur in ActionCam,
  GrassPhenology, LogClassifier, Lyndata, map_tracker, SharedUI,
  TerrainAnalyzer and weight_app_flutter.
- **OBSERVED:** Repositories also call activation, baseline, publication,
  verification, governance-closeout, tooling and component-layer steps Slices.
  ActionCam's SLC-026 and SDP PR #4's SPS-001 became large workstreams.
- **INFERENCE:** “bounded task” and “vertical Slice” are not synonyms. A Slice
  should normally be the smallest independently verifiable end-to-end outcome.
  Studies, preparation, release and governance gates can remain bounded work
  without being mislabelled vertical Slices.

## Pattern 5 — Issue-first delivery is strong but not self-bounding

- **OBSERVED:** ActionCam, Lyndata and TerrainAnalyzer are the clearest accepted
  default-tree examples of Issue-authoritative execution with exact baseline,
  branch/draft PR, evidence and stop/acceptance boundaries.
- **OBSERVED:** ActionCam's Issue #59 still grew to 35 commits, 357 files, six
  Sprints and 27 Slices. HSX uses Issue authority but leaves accepted-in-flight
  work spread across unmerged branches. RadarData created its draft PR before
  its Issue.
- **INFERENCE:** An Issue becomes a useful assignment boundary only when its
  scope, work owner, expected evidence, stop condition and split policy are
  explicit and enforced by the Issue Master.

## Pattern 6 — duplicated mutable state is the dominant failure mode

- **OBSERVED:** Thirteen reports record mutable-state or contract drift among
  GitHub, CurrentAssignment, CurrentIndex, project manifest, Relations, Ledger,
  Features/Sprints/Refactors, review, verification, README/Handoff or executable
  schema: ActionCam, gh-sdp, GrassPhenology, HEOS, HSX, LogClassifier, Lyndata,
  RadarData, SDP, SharedUI, TerrainAnalyzer, weight_app_flutter and
  SDP-Analyzer.
- **OBSERVED:** All three default-tree CurrentAssignment examples—ActionCam,
  Lyndata and TerrainAnalyzer—are stale or contradictory. RadarData's unmerged
  compact assignment already disagrees with copied verification heads.
- **OBSERVED:** Less ceremony alone is not sufficient: SharedUI has no
  traceability triad yet still has contradictory active/draft/completed prose,
  while no-SDP repositories cannot reconstruct assignment or acceptance.
- **INFERENCE:** GitHub-native facts should be observed/generated rather than
  hand-authored in many files. Repository-authored state should focus on
  authority, decisions, semantics, exceptions, evidence meaning and stop gates.

## Pattern 7 — Sprint and Iteration are often redundant, sometimes useful

- **OBSERVED:** ActionCam, gh-sdp, LogClassifier, Lyndata, SDP-Analyzer and
  weight_app_flutter use one Iteration for every observed Sprint; several other
  repositories show the same pattern repeatedly.
- **OBSERVED:** HSX uses Iterations as meaningful programme/refreeze stages;
  SharedUI and parts of TerrainAnalyzer use multiple Iterations or Slices for
  dependency/learning transitions.
- **INFERENCE:** The evidence rejects mandatory Sprint/Iteration identity, not
  all scheduling or coordination groups. A grouping should exist only when it
  represents a real timebox, multi-Slice objective or replanning decision.

## Pattern 8 — one Study is normal; independent Studies are conditional

- **OBSERVED:** gh-sdp, Lyndata, SDP-Analyzer, SharedUI and weight_app_flutter
  use one coherent broad or local Study effectively.
- **OBSERVED:** ActionCam and non-default HSX show the strongest multi-Study,
  synthesis, review and convergence/refreeze practice. Their uncertainty was
  genuinely separable and, in ActionCam, physically risky.
- **INFERENCE:** Canonical SDP should support one Study, multiple independent
  Studies and a convergence gate, but should require the extra machinery only
  when uncertainty, evidence owners, safety or parallelism justify it.

## Pattern 9 — append-only history is useful when it records material change

- **OBSERVED:** Append-only corrections preserve withdrawn claims, review
  rework, block/unblock decisions and changed evidence in ActionCam,
  GrassPhenology, Lyndata, TerrainAnalyzer and other dense repositories.
- **OBSERVED:** High event counts and manual lifecycle logging correlate with
  drift, backdated/imported events and even future-dated HSX events. The merged
  SDP schema currently accepts release events only, while prose calls for work
  transitions; Issue #5 exposed that contradiction in CI.
- **INFERENCE:** The Ledger should retain only decision-, acceptance-,
  recovery-, correction-, merge- and release-significant transitions under a
  versioned general envelope. Routine GitHub and agent activity should be
  derived.

## Evidence constraints on the recommendation

- **OWNER DIRECTION:** Feature/Refactor ownership, universal Issue assignments,
  bounded Issue Masters and ChatGPT Steering Group are mandated destinations,
  not majority-observed successes.
- **RECOMMENDATION:** Separate reusable `steering-group` and `sdp` skills and the
  exact Issue-to-release graph are study recommendations, not owner-decided
  implementation details.
- **OBSERVED:** Only two default repositories exercise Feature ownership, two
  exercise first-class Refactor ownership and none exercise a first-class Fix
  record on default. Only three show accepted Issue-authoritative delivery.
- **INFERENCE:** The proposed workflow must therefore be treated as an
  evidence-informed next contract requiring schemas, pilots and migration—not
  as a description of an already universal practice.
- **INFERENCE:** ActionCam's physical-safety details, HSX's debugger domain
  model, weight_app_flutter's Bluetooth/field gates and other domain mechanics
  stay project-local. Canonical SDP carries their generic authorization,
  bounded evidence, review, recovery and stop patterns.
