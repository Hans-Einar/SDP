# PGI3-M1 — terminal controller implementation checkpoint

Date: 2026-10-07. PLAN-SDP-0018, KB-SDP-038, Session0007 T007.
Disposition: **implemented controller candidate; PGI3 remains in progress**.
Live model execution, executing-turn recovery and native role provenance are not
verified. No independent review, owner acceptance, merge, push or release is claimed.

## Candidate and scope

The base is PGI2 commit 12542b2 in `/tmp/sdp-pgi1-8SeezS/repo`, branch
`sdp/project-governance-pilot`. [Candidate manifest](candidate-PGI3-C1.json)
identifies the exact changed implementation, schema, tests and package documentation,
commands, toolchain, output identities and limits. Tests run on an uncommitted
candidate; the checkpoint commit also contains these process records.

The host now permits writes only under SDL and /tmp. Code and the maintained
Session/card/plan/ledger changes are together in the isolated clone. They have **not**
been copied back to the shared SDPTool/SDP directories. Shared PGI2 and other agents'
changes remain untouched. A later import must recheck destination preimages and
ledger IDs; do not overwrite concurrent records or assume the shared journal is current.

## Implemented behavior

- Controller-only exact owner input and assignment payload are persisted before
  sending, with prepared/sending/acknowledged/uncertain states and stable correlation.
  Unresolved sends refuse automatic retry. Source drift, lease contention, stale
  revision, attempted input rewrite and untrusted provenance refuse.
- The terminal runs an existing covered Maintenance assignment, verifies ChatGPT
  mode and the five scoped MCP tools, and starts a fresh bounded Master turn.
  Explicit sandbox selection defaults to read-only; no model/account override or
  native delegation is granted. This candidate has no Steering/PM classifier.
- Line-oriented steering, interrupt and scope-change preserve the assignment and
  Session. Scope-change records replan-required before asking the host to interrupt.
  Host outcome remains distinct from routine completion and accepted evidence.
- Separate status reads are side-effect-free. Input EOF and closed output observers
  do not cancel work. History recovery retains an epoch gap and replays completed
  items without duplicates; it never resubmits or resumes a model turn.
- Owner input text/payload are bounded to 64 KiB each and 64 inputs per run, with
  inherited state/receipt limits. Run summaries explicitly mark omitted content.
  Last-observation timestamps describe recorded state, not live-process guarantees.

See the [controller contract](../../../../SDPTool/governancecontroller/README.md)
for invocation, exact boundaries and operational limits.

## Evidence

- Final frozen candidate: the complete SDPTool race suite and `go vet ./...` pass.
- 28 observed application wire cases pass: 21 CLI, four host/core and three
  dispatch cases. Four captured work RPC requests match the installed default
  Codex 0.160.1 schemas (thread/start, turn/start, turn/steer, turn/interrupt).
- The source manifest pins 18 changed implementation/schema/test/doc files. The
  shared copies still matched their PGI2 base preimages when inspected; no shared
  files were written. Full tests apply to this isolated candidate only.
- Project-management validation passes: 61 cards, 38 management records, four
  lineage operations and 509 events. Toolkit, new links, ledger byte-prefix
  preservation, candidate hashes and whitespace checks pass.

Crash fixtures interrupt before/after immutable HEAD publication and recover the
same exact input through a fresh Service. A lost response after possible execution
never produces a second host submission. Drift prevents prepared dispatch. Tests
also cover managed-account refusal, controller contention, in-scope steering,
scope-stop, cross-thread refusal, unchanged Session identity and unchanged project
bytes. The compiled command test uses a **recorded fake app-server subprocess**,
not Codex; it exercises run/status/recover and closing the output observer.

The actual rebuilt controller probe against installed Codex 0.160.1 exits 2 with
`reconciliation_required: connection ended: EOF`, before initialization completes.
A separate filtered startup diagnostic exits 1 with no protocol output and matches
read-only-filesystem errors; raw host logs and credentials were not retained. The
exact failing path is not established. No model turn was sent. PGI2's
successful real idle negotiation remains historical evidence from the prior host
permissions; it does not validate PGI3 submission or recovery in this environment.

The first schema check found that the shared Result operation enum omitted the
new dispatch operation. The enum was corrected, including existing host operation
names, and captured outputs were rechecked. No test obligation was removed.

## Remaining work and next step

PGI3 cannot be marked delivered on these fixtures. Re-establish a functioning
app-server in a permitted host, run one useful bounded real Maintenance task,
exercise owner steering and interruption/history recovery, and measure actual role
provenance. Current host delegation restrictions remain in force; instruction text
is not proof of a native Master/Worker/Reviewer hierarchy. Deliberate resolution
of an ambiguous unmapped turn and full executing-turn reattachment remain explicit
gaps; do not bypass them by manual state edits or duplicate runs.

PGI4 still owns candidate-bound evidence and checked Session/card publication;
PGI5 owns independent integrated verification. KB038 remains in-progress for this
unfinished PGI3 work, not solely as a pointer to a later unrelated capability.
Related KB036/037/042/043 remain backlog; ModelGovernance, Blueprints and full KanBan
TUI are unchanged and are not added as prerequisites.
