# PGI1-M1 — core and compiled CLI evidence

Date: 2026-10-06. Plan: PLAN-SDP-0018. Design: SDPTOOL-DES-0001.
Result: the selected core/CLI increment is delivered and passes bounded verification.
This does not complete KB038 or the full owner-to-Codex Maintenance workflow.

## Candidate and environment

Base: `043c59cc50c33d9f2a4e8ba48340984ab5c1bff3` (completed ModelGovernance).
[Candidate manifest](candidate-PGI1.json) records every changed SDPTool file and
SHA-256, test outputs, toolchain and environment. The file-set digest is
`281359a2233a36b761fa498426333af72204eb9bbf41aa2e5da3c35b65b40e38`.
This identifies tested uncommitted implementation on the base; the subsequent
milestone commit also includes documentation and lifecycle evidence. It does not
mislabel tests as having run on a preexisting final commit.

Toolchain: local Go 1.27.1 linux/amd64. Dependencies were already cached; GOPROXY
was off and GOTOOLCHAIN local. GOCACHE and the successful test run's XDG cache and
runtime directories were under `/tmp/sdp-pgi1-8SeezS`. No module versions changed.
The original worktree's Git metadata is read-only in this host. Work and milestone
commit therefore use an independent clone at `/tmp/sdp-pgi1-8SeezS/repo`, branch
`sdp/project-governance-pilot`. Returned workspace files remain uncommitted on its
existing branch; unrelated Session0006, SDL sourceinput and npm files are preserved.
No push, merge, account change or release belongs to this delivery.

## Executed checks

| Check | Outcome and observation |
| --- | --- |
| `go -C SDPTool test ./governance` | Passed core fixtures after the interrupted ledger-check compile fix |
| `go -C SDPTool test . -run TestCompiledGovernanceJourney -v` | Passed; builds actual executable and invokes each operation in a fresh process |
| `go -C SDPTool test -race ./...` with temporary XDG cache/runtime | Passed all SDPTool packages, including ModelGovernance, installation, profile and presentation regressions |
| `go -C SDPTool vet ./...` | Passed, no output |
| `python3 SDPTool/governance/check_schema.py <capture>` | Passed all 21 observed CLI requests/results against the local JSON schema |
| ProjectManagement/Toolkit validation | Passed: 61 cards, 38 management records, four lineage operations and 502 events; Toolkit passed |
| Candidate hashes, ledger byte prefixes, changed links and diff whitespace | Checked before delivery; original committed history retained |

The first full race run failed in existing installer/profile/preview tests because
inherited `/home/warloc/.cache` and `/run/user/1000` paths were read-only. Governance
and ModelGovernance tests passed in that run. Repeating with writable temporary
XDG directories passed the whole suite. No assertions were removed or weakened.
The manifest preserves the failed log digest and the successful package outputs.

Reproduce from the candidate using the selected Go executable and the environment
in the manifest. Create private temporary XDG directories before running the full
suite. Set `PGI_SCHEMA_EVIDENCE` to a temporary directory to emit
`cli-wire-cases.json` from the compiled journey, then run the Python checker above.
The schema checker uses the already installed jsonschema package. Exact capture
bytes vary with temporary roots and random identities; the manifest pins this run.

## Acceptance mapping at PGI1 scope

| Design case | Evidence now | Remaining beyond this increment |
| --- | --- | --- |
| PG-A01 | Status route returns no-work; context/route create no state directory or project files | Client answer and conversation capture |
| PG-A02 | Active Maintenance card/plan/Session plus canonical management chain resolve; one request creates one run, retry returns same receipt | Owner interaction and fresh Master handoff |
| PG-A03 | Missing/ambiguous/incompatible/unavailable/conflict routes distinguished; uncovered start rejected | Installed method catalog and procedure-gap publication |
| PG-A08 | Exact retries, changed-key-content refusal, old/new HEAD recovery, orphan generation handling and competing revision writers | Physical power-loss behavior is not claimed |
| PG-A10 | Pinned serialized routine/digest; unknown/forged role fields and untrusted Caller refuse; local controller can act within current coverage | Actual host identity and routine upgrade integration |
| PG-A18 | Bounded UTF-8/JSON, duplicate and case-alias rejection, unknown fields, unsafe/symlinked paths, finite state capacity | MCP/app-server frames, cursors and paging are not implemented |

Additional core evidence: stale authority/source inputs block advancement; unrelated
file edits do not. Interrupt/cancel can stop after ordinary coverage drift. Run and
Session identity survive a new Service/CLI process. The compiled journey returns
through interrupt/resume/start/report and refuses complete with evidence_missing.
Every project fixture file, including canonical ledger and authored Session,
retains its pre-run digest. Model artifacts are neither required nor changed.

Tests and readable wire definitions live in
[the core](../../../../SDPTool/governance/README.md),
[core tests](../../../../SDPTool/governance/core_test.go),
[compiled CLI journey](../../../../SDPTool/governance_cli_test.go) and
[schema](../../../../SDPTool/governance/schema.json).

## Implementation choices and limits

- Local operational snapshots do not replace ProjectManagement or Traceability.
  Only the selected subjects' canonical chains/placement are checked at runtime;
  full repository validation remains with the existing validator. No project
  records are published by these operations.
- The root-hash directory is a lookup location; the persisted UUID is the checkout
  identity. Registration happens only on explicit bind. The whole canonical
  ledger is conservatively pinned, so another ledger append requires reconciliation
  before advancing the current run. Unrelated ordinary files remain outside that pin.
- CLI invocation supplies trusted local Caller authority. Referenced text does not
  itself prove owner approval. Arbitrary same-user shell access is outside this
  boundary; future agent adapters must establish real caller attribution.
- The routine's verify/review/record requirements remain pinned, but PGI1 cannot
  satisfy them or mark work completed. Resume restores only local routine state;
  no live host attempt, transport replay, managed-account access or independent
  reviewer provenance is established.
- Linux flock, immutable generation writes, file/directory sync and HEAD replacement
  are exercised on local fixtures. Injected failures cover publication boundaries,
  not physical power loss, remote filesystems or non-Linux mutation.
- Review here is a same-context implementation self-check. Current host instructions
  prohibit unsolicited delegation; no independent Reviewer was spawned or claimed.
  Independent complete-pilot review remains an explicit PGI5 obligation.

## Disposition

PGI1-M1 is delivered. PLAN-SDP-0018 stays active; KB038 is active/ready at the
PGI2 handoff. KB036/037/042/043 remain separate backlog capabilities. Session0007
records the interrupted turn and resumed outcome. The next bounded increment is
thin MCP/controller protocol adapters with actual compatibility negotiation;
PGI3–PGI5 retain client execution, evidence publication and integrated verification.
