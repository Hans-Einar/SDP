# RGS1 scenario challenge and future acceptance cases

These are analytical walkthroughs of the [proposed contract](Study.md), performed
by the study author. They are **not executed tests**, an independent review or
proof that an engine/host integration exists. Historical inputs are summarized
in [Evidence.md](Evidence.md); future test fixtures must be sanitized and isolated.
Expected outcomes below are proposed acceptance conditions, not observed runtime
results. All requirement IDs refer to the Study's R01–R16 table.

## S01 — a release request changes implementation scope

**Basis:** E01, recent canonical-template and Go-only work. A request to remove
remaining PowerShell arrives during a release plan.

**Route:** RT14 → RT05/RT07 → RT10 → RT11/RT12 → RT13. Recover the active plan,
compare the new request with its scope and public compatibility requirements,
then revise or separate work before changing affected entry points.

**Expected:** reuse existing owner authorization where it covers the work; pin
candidate and installed payload separately. The view shows the changed scope,
which evidence became stale and what publication is authorized. A passing build
on an earlier candidate cannot satisfy the revised release gate. The study must
not retrospectively claim the earlier real work had no plans.

**Negative case:** “cleanup” classification bypasses version/consumer analysis.
Reject the dependent release transition; continue unaffected investigation.
**Requirements:** R02, R04. **Design consequence:** classify effects, not vocabulary.

## S02 — installing SDP is not adopting all of its working method

**Basis:** E02. XFMD receives a new signed SDP payload while its local application
workflow and adoption work remain separately authoritative.

**Route:** RT03/RT10/RT13. Resolve the installed receipt, project instructions and
uncommitted process state. Detect the adoption distinction before routing product
work using a newly available procedure.

**Expected:** display installed version, actual refinement/authority and any
conflict. Preserve project-owned work during upgrade. A previous running routine
retains its definition; a conflicting refinement stops dependent adoption.

**Negative case:** a version number alone marks all project practices migrated.
**Requirements:** R01, R10. **Design consequence:** capability and adoption are
separate facts; test both fresh installation and upgrade with a local overlay.

## S03 — a question during active work

**Basis:** the owner's mixture of factual questions and implementation steering.
The owner asks what MCP means while a selected maintenance task is active.

**Route:** RT01. Answer with current sources. Preserve the active task and resume
it if the conversation still requests ongoing work.

**Expected:** no duplicate card, plan or new implementation activation. A later
“build an MCP adapter” is classified as action and assessed against scope.

**Negative case:** every mention of MCP activates the adapter roadmap, or every
question requires approval of a plan. **Requirements:** R02, R16.
**Design consequence:** a legitimate no-project-work result is essential.

## S04 — “continue Ponsse” after a context switch

**Basis:** E03. Concept1 records a stop and outstanding acceptance; MVP1 is a
different work line with separate authority.

**Route:** RT14. Recover explicit project, candidate and selected work. Compare
what the new instruction authorizes with the exact outstanding decision. Ask a
bounded question only if the ambiguity affects the next action.

**Expected:** no silent switch from Concept1 to MVP1, no claim that manual
acceptance occurred, and no physical-I/O operation. Independent authorized work
can proceed. The monitor shows why a particular step waits and who can resolve it.

**Negative case:** session restart is treated as either new blanket approval or
loss of all valid earlier authorization. **Requirements:** R01, R04, R14.
**Design consequence:** store permission scope and stop reasons, not just “paused”.

## S05 — MVP1 assignment from a partly unsupported model

**Basis:** E03 and the SDL experiment. The owner wants a BuckingUI change whose
surrounding model uses an experimental profile.

**Route:** RT04/RT05 → RT07/RT08. Determine which SDL facts can actually be loaded,
what static SDUI previews establish and what code/contract context must be supplied
another way. Coordinate with BP2 without silently activating its implementation.

**Expected:** the bundle labels supported facts, authored intent and unknown
coverage. Show affected SimulatorUI/shared consumers if established; do not invent
a complete dependency graph. A narrowly bounded assignment can proceed if its
required context is independently available and selected authority permits it.

**Negative case:** unsupported syntax becomes “no dependencies”, or the entire
project is blocked until every experimental feature is implemented.
**Requirements:** R03, R05, R06. **Design consequence:** partial coverage must be
first-class and readiness must be specific to the selected assignment.

## S06 — a stored assignment says draft, GitHub says merged

**Basis:** E06's verified TerrainAnalyzer PR 21 status contradiction.

**Route:** RT01/RT14 → RT10 if reconciliation is selected. Read the stored
disposition and a fresh remote observation with timestamp and PR/head identity.

**Expected:** show both facts and their sources. Do not infer an owner approval,
rewrite old events or mark unrelated conditions fulfilled. Reconciliation records
what actually changed, preserving the old record's historical meaning.

**Negative case:** dashboard green means “all accepted” because GitHub says merged.
**Requirements:** R01, R07. **Design consequence:** observations and dispositions
must have separate owners and timestamps.

## S07 — review produces findings outside its authorized scope

**Basis:** E07, GrassPhenology's review-only record.

**Route:** RT12 → RT02/RT06 for findings. Close the review deliverable with its
actual disposition; leave unselected fixes as linked future work.

**Expected:** no code mutation merely because a fix is obvious. Review completion
is visible without claiming product conformance or forcing a capability card to
stay active forever. A later owner selection can activate a bounded fix plan.

**Negative case:** automated transition from review finding to implementation
bypasses selection, or self-check is relabeled independent review.
**Requirements:** R04, R07, R14. **Design consequence:** explicit review outcomes
and subsequent work selection are different events.

## S08 — one update affects another derived view

**Basis:** E05's reported agro-crm semantic inconsistency. This is a design challenge,
not a claim that the defect remains in current main.

**Route:** RT04/RT05 → RT07/RT08/RT09 → RT11. Identify producers, consumers,
derived-state rules and cancellation/restart behavior. Add contract-level checks
alongside the local function test where those effects matter.

**Expected:** the ready check names the invariant across consumers. A changed
shared contract invalidates affected earlier evidence. A local test alone cannot
satisfy a gate requiring the cross-view behavior.

**Negative case:** “only this function changed” is accepted as evidence of local
impact without examining its consumers. **Requirements:** R05, R08.
**Design consequence:** evidence level follows the intended outcome.

## S09 — emulator and DAP integration

**Basis:** E08. The adapter passes tests against a fake runtime and references a
provider card that is not checked out locally.

**Route:** RT05/RT08/RT11. Resolve the adapter contract and required integration
level. Preserve the external card as unverified. Obtain the named compatible
provider candidate or report required peer evidence unavailable.

**Expected:** fake tests remain valid for what they prove; they do not complete a
real-runtime requirement. Record both candidates for integration. The agent cannot
supply its own approval in place of a required reviewer. Generic emulator scope
must not acquire P1000 board policy as a convenience fix.

**Negative case:** local missing card means broken reference, or fabricated peer
context satisfies readiness. **Requirements:** R05, R06, R07, R08, R14, R15.
**Design consequence:** external references and required external evidence differ.

## S10 — HSX work is on a non-default branch

**Basis:** E04. Default main and a debugger refactor branch have different process
and integration context; multiple later branches exist.

**Route:** RT01/RT14 → RT05. Resolve the selected checkout and its authority, then
inspect required peers at named revisions. Report unresolved selection honestly.

**Expected:** neither default main nor the most recently named branch is assumed
current merely because it is easy to find. Pinned private sources are accessible
only through the caller's authorized context; the observer gets minimal summaries.

**Negative case:** a repository-wide dashboard displays a single unqualified
“current assignment” from whichever branch was indexed first.
**Requirements:** R01, R15. **Design consequence:** logical project identity is
separate from worktree/candidate identity.

## S11 — no suitable routine, or only a retrieval failure?

**Basis:** the owner's request for discovery of missing procedures; E10 boundaries.
Three fixtures: no procedure for a needed operation; an existing procedure cannot
be read; two applicable procedures conflict. A fourth is an ordinary factual query.

**Route:** RT15 only for the governed cases. Record match candidates and why they
failed. Search an existing gap before creating one. RT01 handles the factual query.

**Expected:** four distinct outcomes. No automated permanent procedure creation.
Read-only investigation continues where allowed; dependent mutation waits if its
conditions cannot be established. A recorded bounded exception is not universal
permission for later runs.

**Negative case:** HTTP failure creates a new routine; keyword similarity grants
publication; every unusual question becomes a gap card.
**Requirements:** R03, R11, R13. **Design consequence:** explicit negative matches
and typed resolution failures are part of the contract.

## S12 — host/service failure and observer reconnect

**Basis:** E10 documented MCP/host interfaces, not a tested integration here.
The observer disconnects, the MCP service restarts, and a command tool remains
available outside the service. Include later input to a persistent shell session.

**Route:** RT14 with the declared host fallback. Restore persisted instance state,
refresh relevant observations and resume the event cursor or snapshot.

**Expected:** observer loss changes no work state. Unavailable enforcement is
visible. Tested guarded operations reject unmet prerequisites; unguarded tool
paths are reported as coverage limits. No universal-interception claim from one
successful PreToolUse test. Service recovery does not silently replay effects.

**Negative case:** the service outage causes all project work to be considered
approved, or makes a disconnected worker look completed.
**Requirements:** R03, R09, R11, R12. **Design consequence:** test both actual host
capture and bypass/continuation paths; a colored view is not enforcement evidence.

## S13 — parallel work and a crash between record updates

**Basis:** proposed operational requirements; no existing generic engine claimed.
Two authorized child assignments target the same run revision. One transition is
retried after a crash between durable state and management-record projection.

**Route:** RT14 recovery. Validate expected revision and idempotency key; reconcile
by operation identity before accepting further dependent transitions.

**Expected:** exactly one logical effect per accepted operation, no lost progress,
no duplicate PM event and no fabricated reviewer result. Stale competing updates
are rejected/rebased explicitly. The observer distinguishes pending projection,
waiting review and historical completed steps. Closing it has no side effect.

**Negative case:** rewriting the management ledger to match the latest in-memory
state hides the interruption. **Requirements:** R07, R08, R09, R12, R14.
**Design consequence:** crash-point tests and role checks are required, not just
happy-path JSON schema validation.

## S14 — procedure upgrade while work is active

**Basis:** current signed installation principles extended as a proposal. A project
has a local routine refinement and an active run pinned to an earlier base version.

**Route:** RT10/RT13/RT16. Check schema/capability and refinement compatibility;
preflight the proposed upgrade without rewriting the active run definition.

**Expected:** preserve project-owned refinements. The current run continues pinned
or migrates through an explicit disposition with affected evidence reassessed.
Historical replay still resolves the old version. Distribution failure leaves a
recoverable installation and does not turn partial adoption into success.

**Negative case:** fetch latest routine text for an old instance and reinterpret
its passed steps under new rules. **Requirements:** R08, R10.
**Design consequence:** compatibility of definitions and instances is distinct from
compatibility of installed file paths.

## S15 — repeated steering without process overload

**Basis:** the owner's requests often evolve during sustained work. Ten messages
clarify one authorized milestone, then one introduces a separate feature.

**Route:** RT14/RT07 for covered clarification; RT02/RT05 for the new feature.
Compare intended effects, not message count. Retain valid earlier permission.

**Expected:** one coherent current plan, material revisions recorded, no ten-card
explosion and no repeated permission prompt for already covered work. Separate
scope stays captured until selected. The monitor explains the distinction briefly.

**Negative case:** silently implement the new feature, or block all work until the
owner learns the process identifiers. **Requirements:** R02, R04, R16.
**Design consequence:** the service should supply context to the user, not require
context reconstruction by the user.

## S16 — repeated rework signals a larger problem

**Basis:** portfolio reviews and the owner's repeated drift experience. Several
local fixes recur with changing ownership assumptions or manual process rescue.

**Route:** RT05/RT15, possibly RT06. Gather linked evidence and compare explanations:
missing procedure, wrong design boundary, weak verification, stale authority or
normal isolated exceptions. Recommend a bounded study/refactor if justified.

**Expected:** no arbitrary numerical threshold automatically starts a refactor or
creates permanent policy. The owner sees convergence problems and a proposed
choice before another narrow fix consumes the same effort. Preserve usable local
work and do not impose a full architecture redesign without selection.

**Negative case:** all recurrence is diagnosed as “agent tunnel vision”, or every
exception produces a unique routine. **Requirements:** R05, R13, R16.
**Design consequence:** measure causes, cost and false positives during the pilot.

## Findings from the analytical challenge

The proposal needs multiple dimensions of state, not one “done” flag (S06/S07/S09).
Procedure absence must be separated from retrieval and compatibility failures
(S05/S11/S12). Version-pinned runs and recoverable projections are necessary for
credible restart and monitoring (S12–S14). Existing authorization and no-work paths
are necessary to avoid turning guidance into constant interruption (S03/S15).
Readiness must tolerate explicitly bounded partial models without claiming whole
system knowledge (S05).

These consequences are incorporated in Study.md and Routine-Catalog.md. The next
DesignPlan should turn a **small subset** into executable fixtures first: S01,
S03, S11, S12, S13 and S14 for one maintenance workflow. Add real application and
blueprint integration cases later. Retain all sixteen as the coverage map; do not
claim they passed merely because this document states their expected outcomes.
