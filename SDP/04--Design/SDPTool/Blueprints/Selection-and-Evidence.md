# BP2-B — selection and evidence contract candidate

2026-10-07, PLAN-SDP-0001 / KB050 / Session0008. This defines a bounded proposed
producer algorithm. The executable experiment uses toolkit ASTs; it is not a
production blueprint command, language extension or complete runtime impact proof.

## Identity and difference

Compare two checked snapshots of the same modeled System and compatible profiles.
An element key is System identity + declaration kind + qualified name; the current
0.6 profile has one System namespace. A kind change is remove/add. A rename is
remove/add until an explicit reviewed identity mapping exists. Reject an unexplained
System identity change instead of guessing correspondence. Artifact UUID, commit,
source digest and parser revision identify the input snapshots independently.

Compare typed facts including all semantic operands, role, mode, property/value,
message and contract fields. Ignore source positions and include-path spelling for
semantic comparison, but retain their per-side provenance and source hashes. A
property change appears as removed/added facts; pairing into MODIFIED is display
convenience only. Source-only formatting/moves still change source identity. Do not
lose removed facts by building an index from TARGET alone.

## Deterministic selection

1. Build a union of NOW and TARGET declarations/facts, retaining separate source
   origins for both sides. Seed every added/removed/changed element and every
   endpoint of a changed fact. Add explicit authored context with its reason.
2. Apply the selected relation-policy version to a sorted work queue. Record every
   included fact and its endpoints. Visit each element/fact once per closure rule;
   cycles terminate without pretending the underlying design is a tree.
3. Include enclosing ownership/containment and allocation context, capability
   realizations/providers and directly connected contract participants. The pilot
   conservatively closes these relations in both directions. Such inclusion means
   potentially relevant context, not proof that the implementation must change.
4. Treat a selected channel and its protocol, permitted messages, payload contracts,
   fields, properties and request/reply links as an atomic obligation group on each
   side. Include participants and modes. Never hide a removed peer or field through
   deduplication across snapshots. Assess missing sender/receiver separately for NOW
   and TARGET: union membership must not conceal an incomplete TARGET.
5. System and mode labels provide context but do not trigger inclusion of every
   sibling or unrelated participant. Retain visible cut edges and the explicit
   excluded inventory, with reasons. An authored inclusion may expand a previously
   excluded branch, as BuckingUI does in the pilot.
6. Record missing dependency/provider/code mappings as UNKNOWN. Unsupported relation
   policies are errors for an executable assignment, not silently ignored facts.
   A diagram radius is a display limit, never a completeness argument.
7. Resource limits fail before publishing a complete assignment bundle. An optional
   diagnostic preview may identify the overflow, but cannot be labeled complete.
   Caller-supplied expected input identities must match the captured snapshots.

| Relation family | Selection behavior | Meaning not inferred |
| --- | --- | --- |
| owns / contains / allocation | Include ownership and execution-placement context; stop expansion through System/mode labels | Runtime call order or code ownership |
| realizes / provides | Include selected capability and its modeled providers/contributors | Every provider is affected in implementation |
| channel participation | Include channel, message, mode and both endpoints on each side | Connectivity or delivery guarantees beyond the model |
| upholds / permits / replies-to / has-field / properties | Include full modeled contract closure | Conformance to real protobuf/database implementation |
| other semantic relations | Require an explicit versioned policy; report unsupported until defined | Generic graph reachability proves impact |

Authored exclusions are proposed scope boundaries. If an exclusion intersects a
changed endpoint or an atomic contract group, fail with a scope conflict; do not
silently drop required context. Remote systems are boundary references unless an
explicit pinned model snapshot is supplied. No checkout or GitHub lookup is required
to acknowledge an external dependency. Mark unavailable definitions UNKNOWN.

## Completeness and the observed parser gap

The experiment removed MachineService's ReadBaseline receiver. The real frontend
still accepted the model. Selection therefore checks endpoint coverage itself and
reports TARGET missing channel peer. This is a blueprint evidence obligation, not
an assertion that the language validator has a defect. Even paired modeled endpoints
do not prove a runtime provider is available or healthy.

Selection can establish closure only over known facts and supported policies. The
MVP1 specimen lacks calibration-to-reduction calls, detailed browser delivery,
external protocol semantics and model-to-code mappings. These remain UNKNOWN even
if every vertex in the reduced model appears in the view. An assignment may proceed
only with an explicit human disposition of relevant unknowns; no automatic pass.

## PRESERVE versus CHANGE

Selection and permission are separate. CalibrateMeasurement ownership changes from
MachineService to CalibrationProcessor inside the same container; allocation stays
with MachineService. PRESERVE includes ReduceMachineState/PublishMachineBaseline
ownership and the boundary contract. The negative specimen moves ReduceMachineState
to BuckingWeb. It parses successfully, but fails an authored ownership assertion.
A code-only violation with unchanged SDL cannot be detected by this structural check.

Constraints carry ID, author/source, scope, predicate or referenced procedure,
NOW expectation, TARGET expectation and evidence status. Conflicting CHANGE/PRESERVE
rules block assignment publication pending resolution, rather than favoring whichever
was loaded last. Preserve rules must identify whether an enclosing container's
boundary or all of its internals are protected.

## Code tags and verification references — coordination with KB004

Tags are locators, not proof. A proposed mapping record contains System/element key,
code revision + file digest, repository-relative path and symbol/region selector,
mapping role, provenance and reviewer disposition. It may cite existing inline tags;
BP2 adopts no new mandatory tag grammar. Multiple implementations require distinct
mapping roles; absent, stale or ambiguous mappings remain UNKNOWN. File/line numbers
alone are insufficient across edits. Do not infer authority from matching names.

A check reference identifies check ID, baseline (NOW or TARGET), source digest,
code revision/content digest, exact command and fixture/config inputs, environment,
observed result, evidence artifact hash and the obligation it actually verifies.
Use not-run, pass, fail, unavailable and stale distinctly. A pass against NOW cannot
be reused as a TARGET pass without demonstrated applicability. Parser validation,
structural obligation checks, behavioral tests and independent review remain separate
rows. Results belong in existing Traceability records linked to the assignment/card;
no replacement execution ledger is introduced here.

The pilot deliberately has no code mapping. Model-only release is not implemented
NOW. No real machine protocol execution, Go behavior test or general SDL code-tag
verifier is delivered by this milestone. KB004 and KB-SDL-006 retain those scopes.

## Evidence and production handoff

`experiments/blueprint_mvp1/selection_probe.py` reads JSON from the real `sdl check`
and `sdl ast` commands built from commit 043c59c. It demonstrates the pilot relation
subset; it does not parse SDL independently. The tests include a synthetic cycle
case explicitly labeled as normalized-graph evidence, not new language support.
The retained summary and raw test output are in Selection-evidence.json and
Selection-tests.txt. Initial missing-peer rejection expectations were disproved
and corrected to an explicit selection diagnostic, without changing the parser.

Production Go must consume captured source bytes and checked in-memory snapshots,
not the experiment's separate filesystem CLI reads. It must also implement System
identity/profile checks, full supported relation policies, bounded fact/byte counts,
constraint conflict resolution and evidence publication before claiming a safe
assignment producer. The experiment proves no hostile-input hardening or concurrency
atomicity. BP2-C owns the API/publisher contract and implementation slices.
