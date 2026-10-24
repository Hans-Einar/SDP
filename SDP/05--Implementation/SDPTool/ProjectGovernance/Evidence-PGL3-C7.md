# PGL3-C7 — Serialized leader delivery and mapped recovery

Session0007 T033, PLAN-SDP-0018 and KB-SDP-052.
[Design](../../../04--Design/SDPTool/ProjectGovernance/Orchestration.md#pgl3-c7--serialized-project-leader-wake-up-and-mapped-recovery-2026-10-24),
[preflight](preflight-PGL3-C7.json),
[commands and bounds](../../../../SDPTool/governancecontroller/README.md#serialized-project-leader-turns--pgl3-c7).

## Candidate and delivered behavior

Baseline dca5c4cea87d386bc60f7690fabe9fc0b28a72a5 plus
[candidate source hashes](evidence-PGL3-C7/candidate.json), on the existing
sdp/codex-forward-compatibility branch with milestone commits. Only this scoped
continuation was observed in the registered worktree. No global host configuration,
stable installation, source model or external repository is changed.

Operational profile 0.9 adds durable attributed owner messages and one active leader
dispatch to the existing atomic store. Leader receipts share the global operation
namespace and preserve earlier assignment/preparation/attention/report records.
Inspection is read-only; message enqueue never launches a model. Separate effect
leases leave owner enqueue and worker operations available.

An explicitly enabled bounded loop prepares short decision turns, opens a fresh
persistent read-only thread, verifies exclusive coordinator MCP inventory and sends
exact skill inputs. Pending owner messages precede bounded changed report/question
notices. Durable claims precede external effects. New owner input or changed selected
skill/notice facts supersede safely unsent work. Uncertain opening/sending cannot
be replaced automatically. Recovery requires original mapped thread/cwd/host version
and exactly one matching client message plus exact payload/skills in full history.
Only visible completed agent messages enter the answer. Host completion is delivery,
not assignment acceptance, semantic consumption, independent review or approval.

## Evidence matrix

| Obligation | Actual evidence | Result boundary |
| --- | --- | --- |
| Queue persistence, retry and no implicit assignment | Core and compiled CLI inspect/message/retry tests | Unbound inspect creates no store; explicit enqueue binds; stable operation/content retries preserve receipts |
| Serialization and owner priority | Service/controller tests, independent owner enqueue during send | One active dispatch/effect lease; stale preparation supersedes; new input remains queued for next turn |
| Uncertain effects | Controller fixtures for lost thread/start and turn/start replies | No replacement opening; sending recovers exactly or refuses without resend |
| Exact recovery | Missing/duplicate/conflicting input, partial history, wrong host version/cwd/thread/turn/status tests | All refuse while preserving unresolved state; allowlisted answer excludes private reasoning |
| Restricted coordinator launch | Actual host-client protocol fixture plus exclusive-inventory tests | Read-only/never, explicit skill inputs and exactly coordinator server; extra servers/pagination refuse |
| Restarted executable workflow | Compiled controller plus recorded app-server subprocess | Lost send, owner console enqueue, new-process recovery, next owner message delivery; one original thread/start and turn/start before recovery |
| Profile coexistence and old-reader refusal | Report delivery/acknowledgement tests and actual retained C6 binary | No run advancement/profile downgrade; old binary refuses unknown leader state |
| Bounded data | Core message-capacity/encoded batch tests and output schema checks | 64 messages, 16 dispatches, 256 receipts; 64 KiB payload/answer, up to eight messages/notices |
| Regression and candidate integrity | Retained full SDPTool race/vet, source/artifact hashes, validators | See machine results; no independent-review claim |

The protocol subprocess is explicitly the test executable, never installed Codex.
It reports compatible future version 0.165.0 to exercise the existing >=0.164.0
policy. This is protocol simulation, not a native-host compatibility certification.
The original lost-response test expected raw server error text; corrected it to
expect the existing sanitized `host_rpc_error` contract, without changing product
error handling. Encoded message batching was refined to fit the byte limit without
blocking delivery of a smaller leading batch. Notice fingerprint normalization is
separate from actionable revisions retained in context.

## Reproduction and limitations

Run [checks.py](evidence-PGL3-C7/checks.py) for the workstation toolchain, retained
compiled workflow captures and actual previous-reader refusal; then run
[verify_wire.py](evidence-PGL3-C7/verify_wire.py) for output/context schemas.
[checks.json](evidence-PGL3-C7/checks.json),
[schema-checks.json](evidence-PGL3-C7/schema-checks.json) and
[manifest.json](evidence-PGL3-C7/manifest.json) record actual outcomes and identities.
Fixture paths/UUIDs/times in captures are observed synthetic-run values, not project
work IDs. The private binaries are separately hashed; they are not a release.

No native model/subagent task or independent review was performed. Actual native
leader scheduling, recovery with a real host, complete SAD/scope-change workflow,
sustained-hours capacity/compaction and canonical evidence publication remain parent
plan gates. The loop cannot automatically resolve an unmapped thread/start; reaching
lifetime bounds refuses further writes. It does not execute worker launch/approval
proposals from prose. Manual Sessions without registered runs remain outside the
operational work projection. SDL-only direction, Blueprint pause, deferred KB059
code/runtime ideas and Ponsse MVP1 deferral remain unchanged. No merge or release.
