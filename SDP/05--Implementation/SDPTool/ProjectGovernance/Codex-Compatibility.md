# PGC1 — Codex 0.164.0 and compatible later hosts

Owner Session0007 T029, 24 October 2026. Selected under PLAN-SDP-0018 / KB057.

## Intent, choices and boundary

Replace the single-release 0.161.0 connection pin with a minimum stable version
0.164.0 and an explicitly bounded wire profile. Later stable versions are eligible
without a code edit for each release, provided the operations used by PG satisfy
the existing handshake, required response and authority/recovery contracts. Being
newer is not evidence of tested compatibility or permission to retry an uncertain
operation. Unknown additional fields remain ignored; malformed/unsupported required
fields and RPC methods fail explicitly. No version spoofing or silent downgrade.

Alternatives: another exact 0.164.0 pin repeats the maintenance problem; accepting
any version without preserving actual host identity corrupts provenance. A minimum
version with per-operation contract checks preserves forward eligibility and honest
failure. It does not promise support for arbitrary future breaking API changes.
Export actual 0.164.0 schemas and check the used requests; fixtures cover compatible
future versions and missing/changed required responses. Keep optional/native-role
observations unavailable where they cannot be established.

The observed version belongs to the connection and persisted attempt, not a global
constant. Stdio executable version must agree with initialize userAgent; socket
initialization uses the same minimum and observed identity with existing peer checks.
Report actual observed version in probe/provenance. Recovery/reattachment must not
relabel old attempts as the new host; read-only legacy history remains readable.
Keep legacy 0.160.1/0.161.0 observation validation for retained history/core fixtures,
while new network/process connections require at least 0.164.0.

## Selected milestones

| Milestone | Acceptance | State |
| --- | --- | --- |
| PGC1-A Contract and preflight | Source impact, owner intent, floor/profile and provenance/recovery boundary documented | delivered |
| PGC1-B Adapter/core compatibility | Stable version range, actual connection/attempt provenance, handshake/response checks and compatible-future/refusal tests | delivered; bounded implementation/scoped race checks pass |
| PGC1-C Candidate evidence | Scoped and full affected race/vet, real 0.164.0 no-model probe/MCP discovery and actual request-schema checks; precise native/independent-review limits | delivered; exact 98b449c checks pass at declared boundary |

No new leader MCP rights, report scheduler, evidence publication, default model,
account mode, sandbox permission, Ponsse pilot/model edits, release or merge selected.
Use one agent; native model/child trials and independent review remain distinct from
protocol fixtures and read-only live probes. No native task is silently launched by
these checks. The original PG acceptance gaps remain with PLAN-SDP-0018.

[Implementation/evidence](Evidence-PGC1.md) records PGC1-B. PGC1-C candidate checks now pass at the declared boundary; eligible future versions are not declared tested native hosts.
