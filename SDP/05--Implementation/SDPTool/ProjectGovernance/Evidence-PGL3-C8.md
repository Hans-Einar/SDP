# PGL3-C8 — Native coordination and recovery in pg_lab

Owner-authorized actual-model test, Session0007 T034 / PLAN-SDP-0018 / KB-SDP-052.
[Preflight](preflight-PGL3-C8.json), [native result](evidence-PGL3-C8/native-summary.json),
[controller contract](../../../../SDPTool/governancecontroller/README.md#serialized-project-leader-turns--pgl3-c7).

## Outcome

The bounded native workflow passed in `/home/warloc/git/pg_lab`, on a dedicated
`test/pgl3-c8-native-leader` branch. Actual Codex 0.164.0 and gpt-6-astra with high
reasoning performed three Project Leader turns and one bounded worker turn using
the existing managed account. Model/effort were observed in native thread/start
responses, not inferred from fixture labels. No helper/reviewer agents were spawned.

The worker added FormatLabel and tests in exactly two permitted files. It reported
through assignment MCP. Project tests and five externally supplied behavior cases
passed. All 38 baseline regular files remained unchanged during model execution.
The Project Leader received owner input queued during its first running turn,
read current coordination state and acknowledged the worker report through MCP.
The run remains verifying; receipt never advanced it or granted owner acceptance.
The lab's pinned coverage records remain fixtures so current inspection stays useful.

## Actual recovery evidence

| Native case | Retained observation | Result |
| --- | --- | --- |
| Owner message during leader execution | [Immutable enqueue-generation proof](evidence-PGL3-C8/enqueue-proof.json) | Message queued while dispatch was running, selected by next decision turn |
| Completed read response lost | Proxy withheld the real completed thread/read response and closed transport | New process recovered the same mapped turn; only initialize/thread-read requests |
| Submission response lost | Proxy withheld actual turn/start response, waited for native completion, then closed transport | Core remained sending without a turn ID; unique original input correlation recovered the completed turn |
| Report processing | [Native inbox](evidence-PGL3-C8/native/report-inbox.json) and report receipt | Exactly one acknowledgement; assignment still verifying |
| Model scope | Protected file hashes and worker file hashes in native result | Only format.go and format_test.go added; external verification passed |

Each recovery trace contains no thread/start or turn/start. Each original faulted
execution contains exactly one turn/start. The proxy never fabricates a model reply,
changes content, edits the operational store or suppresses a real uncertainty by
resetting state. It records allowlisted transport metadata, not private reasoning,
authentication or raw tool output. Faults waited for actual completion: attachment
to a model still executing after host death is a separate, unproved case.

## Discovery and correction

The original C7 candidate stopped before model submission when default Codex included
307 Apps tools alongside the two PG tools: the response was 1,163,770 bytes, exceeding
the existing 1 MiB frame bound. A process-local wrapper disabling Apps, remote plugins
and delegation reduced the measured inventory to exactly two PG tools and 51,041
bytes. Global settings and the transport bound are unchanged. Local CLI feature
inspection and [official configuration guidance](https://learn.chatgpt.com/docs/config-file/config-basic)
support the per-process feature overrides; the actual measured inventory is retained.

The restart then revealed a genuine C7 lifetime defect: an unsent bound thread may
not remain loaded when its owned host process exits. The controller now tracks threads
created by its current instance and safely supersedes an unsent bound thread from
an earlier instance. Pending messages are reselected in a fresh thread. Durable
sending/running claims still require exact recovery and can never be superseded.
A targeted regression proves both halves. The abandoned unsent preparation and failed
logs remain retained; no store migration, fabricated acceptance or resend workaround.

## Candidate and verification boundary

Product baseline df5728aff607d1c1a5b5ec964f95a4d2f57f0e7e plus source/test/proxy
hashes and actual binary hashes in [native-summary.json](evidence-PGL3-C8/native-summary.json).
The lab baseline is f66e65ec522ed19c8dc31c4216bf28697cada2f0, before native worker edits.
The coordinator later added a read-only `./pg status|leader|reports` helper and the
lab result document; those additions are not attributed to a model.

[Checks](evidence-PGL3-C8/checks.json) record race tests and vet for the affected
controller, host and command packages. The unchanged core retains C7's full-suite
baseline evidence. [Schema checks](evidence-PGL3-C8/schema-checks.json) validate eleven
actual native output/context objects; [verify_native.py](evidence-PGL3-C8/verify_native.py)
also checks model/host observations and absence of resubmission during recovery.
The maintained [native proxy](../../../../SDPTool/labs/project-governance/native_host_proxy.py)
provides deterministic fault injection against the actual host.

This is one bounded worker/leader laboratory workflow, not two-Master concurrency,
sustained-hours capacity, independent review or full ProjectGovernance readiness.
SAD/scope closure, capacity/compaction, active-host recovery and canonical evidence
publication remain parent-plan gates. No merge/release or production owner acceptance.
Blueprint pause, SDL-only context direction and Ponsse MVP1 deferral are preserved.
