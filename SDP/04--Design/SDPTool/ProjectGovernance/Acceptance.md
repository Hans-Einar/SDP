# PGD1 — Pilot acceptance and requirement mapping

These are executable acceptance specifications for the implementation handoff,
not passing runtime tests. [Design](Design.md) owns behavior; cards and RGS2 own
the source requirements. Use an isolated project fixture with one real bounded
Maintenance change, existing work records, and a deliberately unrelated dirty file.

| Case | Governing source | Trigger and required observable result |
| --- | --- | --- |
| PG-A01 Direct answer | KB036; RGS2 explain versus act | Submit status/unrelated question: answer through the read route; no new card, plan or run |
| PG-A02 Existing work | KB036/038 | Submit covered change: resolve current card/plan/authority once; show routine and exact assignment; no duplicate plan or approval request |
| PG-A03 Missing or ambiguous route | KB036/037 | Remove routine, make it unreadable, incompatible or conflicting: four distinct dispositions; dependent mutation rejected; repeat does not create duplicate gap cards |
| PG-A04 Bounded execution | KB038; RGS2 role model | Fresh Master receives governing context, owned paths and return conditions; preserve unrelated edits; observed child identities and context rules match the claimed role mapping |
| PG-A05 Scope change | KB036/037 | Steer with an out-of-scope addition during work: retain original goal, record delta, stop dependent progress, reconcile in-flight effects and revise coverage |
| PG-A06 Evidence and rework | KB037 | Completion without required evidence refuses; review finding returns for rework; self-check cannot satisfy required independent review |
| PG-A07 Candidate changes | KB037; Traceability | Modify relevant source after verification: evidence remains historical but cannot satisfy the current gate; unrelated excluded file change alone must not invalidate scoped evidence |
| PG-A08 Core retry/crash | KB037 | Repeat same operation/key before and after simulated crash: same receipt/revision; changed arguments with same key refuse; orphan generation does not advance head |
| PG-A09 Record publication recovery | KB037/042 | Crash between journal/card/ledger writes: reconcile exactly once; preserve original ledger prefix; changed external preimage stops with explicit conflict |
| PG-A10 Run version and authority | KB036/037 | Upgrade routine or submit forged owner role: old run stays pinned; forged input cannot grant rights; existing authorization still works within scope |
| PG-A11 Observer and controller | KB038 | Close observer: work continues. Lose controller transport: state becomes uncertain; reconnect to mapped work without duplicate turn submission |
| PG-A12 Conversation capture | KB038/042 | Multiple owner inputs, steering, commentary, final and interruption: retain observed distinctions and IDs; no fabricated final response; replay completed items without duplication |
| PG-A13 Session continuity | KB042 | Resume on a new thread/attempt: same Session/assignment, new attempt; current next step and source revisions recovered; authored text and old entries preserved |
| PG-A14 Honest timing | KB043 input | Clock reset or disconnected interval: epoch/gap explicit; no invented active time or live cursor; static output exposes as-of position |
| PG-A15 Compatibility and auth | KB038; RGS2 | Pinned app-server and chosen MCP SDK negotiate in a real local probe; managed account path exercised without API-key fallback; unsupported capabilities refuse explicitly |
| PG-A16 Return disposition | KB037/038/042 | Completed model turn alone does not close card; exact candidate/test/review result and remaining work are recorded; owner-reserved acceptance requires actual provenance |
| PG-A17 Model boundary | Owner ProjectGovernance scope; KB049/050 | Ordinary maintenance works without model store/blueprints; optional unavailable model reference stays explicit; no model artifact is mutated by process bookkeeping |
| PG-A18 Bounded protocol | Existing SDPTool strict input conventions | Duplicate/unknown fields, escaping paths, oversized frames and invalid cursors refuse without corrupting state; finite page limits enforced |

## Evidence ladder

1. Pure-core state/fixture tests establish guards, persistence and replay only.
2. Adapter conformance tests establish equivalent results through CLI and MCP;
   record the selected protocol/SDK versions and actual negotiation.
3. A real controller/app-server trial establishes thread/item/auth behavior and
   reconnect coverage. Schema export alone cannot satisfy this level.
4. One complete Maintenance task establishes useful owner-to-evidence behavior.
   Independent review uses a context independent of implementation, the governing
   intent and exact candidate. If host provenance is unavailable, record the gap.

For every result record candidate revision/digests, command/environment, inputs,
observed outputs, verdict and limitations. Do not label simulated protocol fixtures
as a live Codex trial or an agent's declared skill use as measured compliance.
PG-A15 requires real account/host conditions; inability to run it is an explicit
remaining integration condition, not permission to mark the entire pilot complete.

## Fixture preservation and review

Use pre/post digests for unrelated files, original ledger prefix and authored Session
content. Include negative/recovery paths alongside the successful workflow. Reuse
project management validators for publication output. The full 16-family catalog,
publication routines and timeline renderer remain deferred after a successful pilot.
