# PGI4-C3 — Checked result publication and review coordination

Session0007 T037 / PLAN-SDP-0018 / KB038 and KB052. Baseline ef757ae; current
working branch and milestone commit policy. [Preflight](preflight-PGI4-C3.json)
records the green PG-local scope. Exact candidate:
[evidence-PGI4-C3/candidate.json](evidence-PGI4-C3/candidate.json).

## Delivered boundary

A trusted controller can prepare a reviewable six-file proposal for a current
passed review and its attributed passing tests, inspect it, apply it, or abandon
an unapplied proposal. Publication records a bounded result in the selected
Session/card/plan, a new evidence document and canonical management/Traceability
history. Original document and ledger bytes are preserved as prefixes. Card,
plan, Session lifecycle and authored roadmap remain unchanged. Run completion and
owner acceptance are separate, still unimplemented dispositions in this increment.

Profile 0.12 retains exact publication intent. Per-project locking coordinates
separate PG stores; the pending-store pointer prevents another store from taking
over interrupted effects. Persistence precedes file effects. Preimage/postimage
comparison, fsync and exact retry provide bounded recovery, not multi-file atomicity
or protection against uncooperative same-user editors. Unknown external bytes stop
rather than being overwritten. Reads expose historical status and current validity.

The project validator reuses existing installation history checks, with a distinct
canonical read policy for already-recorded 0.1-to-0.2 typed-plan adoption. Installation
retains its existing stricter policy, covered by a regression test. PG additionally
checks supported payload fields, counters/times, index, references and grouping.
This is validation of the selected shared-management profile, not arbitrary
project-specific policy or semantic truth of prose.

Coordination version 2 adds compact review notices with observed context, verdict,
summary/finding count, applicability and publication status/applicability. Terminal
reviews enter the existing serialized leader queue; content deduplication avoids
repeat delivery while changes in applicability create a new notice. Existing MCP
capabilities remain coordination_read and report_acknowledge for coordinators;
there is no model-callable publication or owner-approval capability.

## Evidence and limits

- [Final checks](evidence-PGI4-C3/checks.json): final binary builds, full SDPTool
  race suite and vet; command/environment and raw logs retained. No product code
  changed after these checks. One schema-only correction preceded the final rerun;
  earlier passing checks remain explicitly preliminary.
- [Observed schemas](evidence-PGI4-C3/schema-checks.json): core request/result,
  publication CLI and coordinator CLI/MCP output validated against the final schemas.
- [Compiled lab checks](evidence-PGI4-C3/lab-checks.json): final publication
  `97a79b56-b20c-455c-b06d-80c710cab9fe`; six exact postimages unchanged on retry,
  one card/one plan/four management events and one verification event pass canonical
  validation. Original 63 lab files stay unchanged; README and the read-only `pg`
  helper are intentionally extended. The first publication trial preserved all 65
  prior files before those two command/document updates.
- [Actual MCP exchange](evidence-PGI4-C3/mcp-protocol.json): initialized local MCP,
  complete two-tool coordinator inventory and read-only review/publication result.
  This is local transport evidence, not another native Codex run.
- [Retained native review notices](evidence-PGI4-C3/retained-native-review-notices.json):
  C2 actual positive and negative review results now appear in coordination. Their
  original stores remain unchanged. The new publication fixture uses synthetic
  protocol review observations; it does not claim fresh native review independence.
- [Old reader](evidence-PGI4-C3/old-reader.json): actual C2 executable refuses the
  new publication store without changing HEAD. Current writers retain the profile.

Core tests exercise crashes after durable intent, after each of six files, and
before/after the final store head; original project prefixes and exact postimages
are checked. Failed intent-save recovery/abandon, conflicting external edits,
source/test/authority drift, symlinks, active reviews, contradictory failed tests,
unsupported minimal history, invalid canonical payload/index/references and another
store's pending claim are covered. These are deterministic process-boundary tests,
not physical power-loss tests. Leader payload/deduplication is tested without a new
native leader turn; C8's native delivery evidence remains historical and separate.

## Remaining work

PGI4 final run/lifecycle disposition is still open: publication does not automatically
close work or rebase the old assignment after canonical input bytes change. Broader
PGL2/PGL4 SAD/owner flows, parallel/active-host/capacity checks and PGI5 independent
integrated verification remain. No independent review of this PG implementation,
merge, release, Blueprint restart or Ponsse MVP1 execution is claimed.
