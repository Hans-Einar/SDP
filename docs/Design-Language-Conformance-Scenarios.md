# Design-language conformance scenarios

Date: 2026-09-15

Status: candidate language test plan with worked sentence cases. The owner has
selected the goal of a compact, expressive, unambiguous language without dialects
or synonymous canonical forms. The concrete syntax below remains a proposal.
The initial study did not execute a parser, formal conformance suite or product
implementation.

Related: [grammar exploration](Vocabulary-and-Grammar-Exploration.md),
[MVP1 example](MVP1-Design-Language-Example.md), and
[behavior/source traceability](Scenarios-State-and-Implementation-Traceability.md).

Update 2026-09-16: [Design language definition](Design-Language-Definition.md)
is now the primary working syntax/type definition. Its bounded `design-core 0.1`
supports structural ownership/realization and mode-qualified dependencies.
The observation templates below and the full flow cases remain extension
proposals, not accepted alternative syntax. The subsequent
[parser prototype](../experiments/design_core/README.md) tests the bounded core,
including its complete definition examples. This does not establish conformance
for the full behavioral scenarios below.

Update 2026-09-17: the owner accepts optional explicit type annotations on future
typed references. See [definition Section 15](Design-Language-Definition.md#15-explicit-type-annotations-and-checking-without-annotations).
One semantic meaning and canonical export are still required; valid authored
source may include or omit these checked assertions. Neighboring sentences never
supply implicit types. This changes the earlier single-source-spelling goal,
not the supported grammar of the existing 0.1 parser.

## 1. What we are testing now

Use design scenarios to discover which distinctions the language must express.
For each distinction, select one canonical construction, define its interpretation
and create cases that must be rejected. The unit of progress is a resolved
semantic distinction and its test cases, not a growing list of synonyms.

The owner can continue using ordinary short prompts. The agent translates them
into canonical statements and exposes unresolved meaning instead of silently
choosing a convenient interpretation. The canonical language itself remains strict.

Separate three goals:

1. **One interpretation per accepted sentence:** participant types, scope,
   conditions and effects are explicit enough to exclude competing readings.
2. **One canonical form per modeled fact:** fixed vocabulary, direction, clause
   order, identifiers and serialization define its normalized export. Optional
   source type annotations normalize to that same fact after validation.
3. **One maintained identity across views:** text, diagrams, tables and source
   bindings reference the same facts instead of creating independent copies.

This is a goal for a defined, bounded semantic model. It is not a claim that we
can automatically recognize every equivalent program or architecture. Two
different designs may legitimately achieve the same user outcome; those are
different realizations, not linguistic synonyms to collapse.

## 2. Candidate canonical-expression rules

| Concern | Proposed rule |
|---|---|
| Vocabulary | One registered term per concept and one verb per defined relation. No author-defined synonyms inside design sentences. |
| Direction | Use the registered subject/verb/object roles. Do not store both a relation and its passive/inverse wording. |
| Identity | One stable identifier per model object in its declared namespace. Renaming a display label does not create a new object. |
| Qualification | Fixed qualifier names and order, with declared types and scope. Required qualifiers cannot be implied from nearby prose. |
| Type annotations | Optional prefix type assertions improve local readability; resolve and check declarations in all cases. Annotation presence cannot change fact identity or excuse duplicate facts. |
| Properties | A closed property vocabulary for each subject type. A property cannot duplicate a relation under a second spelling. |
| Collections | Unordered sets have one serialization order and no duplicates. Ordered steps preserve explicit order; sorting them would change meaning. |
| Logic | Initially prefer references to named Conditions with a bounded predicate grammar over unrestricted natural-English expressions. |
| Missing information | Missing, unknown, false and prohibited remain distinct. Omission cannot invent a default obligation or a negative fact. |
| Provenance | Authority status, model revision and source/evidence attach to statements under one defined scheme. They are not inferred from tense. |
| Extension | Add a construct only after a scenario demonstrates meaning the existing vocabulary cannot express adequately. Update all examples together. |

Readable explanations may surround canonical blocks, but they do not override
the statements or introduce hidden requirements. If two explanations disagree
about a sentence, the sentence or its definition needs revision.

The existing grammar exploration and MVP1 prose remain research examples, not
a competing dialect. A later accepted profile must migrate their canonical
examples explicitly; it must not advertise every historical form as valid.

## 3. Scenario portfolio

These are language-design cases. New conditions used to stress the language are
hypothetical until adopted in the product's own design authority.

| Case | Design situation | Distinction the language must preserve | Required artifact |
|---|---|---|---|
| L01 | One UI unit owns schema validation; several local responsibilities realize layout replacement. | Ownership, containment, contribution and provided offer are different relations. | Typed structural statements and rejected synonyms/type errors. |
| L02 | Live inspection needs fresh observations; cached inspection remains available offline. | Consumption versus required dependency, qualified by mode. | One scoped dependency form; an explicit account of optional/degraded behavior. |
| L03 | Candidate layout B is validated, then its Composition contract changes before publication. | Past result versus current guard; entry condition versus commit condition. | Branch/state description with exact revisions and check point; no implicit global snapshot. |
| L04 | Measurements advance while a layout is prepared; a large matrix update is delayed. | Causal order, concurrency, delivery policy and durable versus supersedable data. | Control/data dependencies and frame conditions; clear completion/join policy where needed. |
| L05 | An APT edit times out, but the domain may already have accepted it. | Command identity, execution phase, observed outcome and recovery. | Correlated observation statements; rejection of timeout-as-cancellation and retry-as-new-action. |
| L06 | Two views show one length in different display units. | Semantic quantity versus display formatting; one Representation versus two view instances. | Typed value/unit/identity facts and no duplicated domain subscription. |
| L07 | The headless UI runtime moves from browser-local execution to a Go Host. | Logical ownership versus process/deployment binding; changed boundary contract without automatic Feature renaming. | Before/after allocation with preserved behavior and explicit new communication obligations. |
| L08 | An event has no handler; an existing source symbol silently ignores it. | Desired diagnostic behavior, actual implementation, missing realization and evidence coverage. | Role-specific source binding and a declared gap, rather than a false completeness claim. |

Together these cases test the selected MVP1 design's structure, state, behavior,
failure, quantities, deployment and traceability. They do not require building
all those product paths before assessing the language.

## 4. Worked canonical sentence cases

These proposed forms deliberately expose missing grammar work. Symbol names are
example model identities whose types are stated below. They are not runtime
observations or a validated parser input corpus.

### L01: ownership and realization

`PresentationManager` is a Unit; `ValidateBindings` and `PublishSchema` are
Functionalities; `LayoutReplacement` is a Capability.

Canonical candidates:

```text
PresentationManager owns ValidateBindings.
PresentationManager owns PublishSchema.
ValidateBindings realizes LayoutReplacement.
PublishSchema realizes LayoutReplacement.
```

Expected interpretation: one immediate owner per Functionality; two contributions
to one capability. These facts do not prove complete realization or publish an
external interface.

| Reject | Reason |
|---|---|
| `ValidateBindings is owned by PresentationManager.` | Passive synonym of the registered ownership form. |
| `PresentationManager is responsible for ValidateBindings.` | Competing verb phrase for the same ownership fact. |
| `PresentationManager has owner-of = ValidateBindings.` | A property used to encode ownership a second way. |
| `LayoutReplacement owns ValidateBindings.` | Subject type does not satisfy the ownership signature. |
| `ValidateBindings enables LayoutReplacement.` | `enables` has no registered meaning; it might mean contribution, guard satisfaction or permission. |

The existing word `realizes` means **contributes to realization** here. If a
complete-realization claim is needed, it needs a separate evidence-bearing
concept, not a second interpretation of this verb.

### L02: mode-scoped dependency

`OperatorUI` is a Unit; `LiveInspection` is a Capability; `ObservationInterface`
is an Interface; `Live` is a declared operating mode in the relevant scope.

Canonical candidates:

```text
OperatorUI consumes ObservationInterface.
LiveInspection requires ObservationInterface in mode Live.
```

The second sentence requires an extension to the initial structural grammar:
one explicit dependency-scope qualifier. A future profile should represent an
unconditional dependency with one explicitly defined scope, rather than allowing
both omitted and explicit qualifiers to mean the same thing.

Reject `OperatorUI needs observations` as a design sentence: actor, contract,
necessity and mode are insufficiently defined. Reject `utilizes` as a synonym
for `consumes`. Do not replace `requires` with `consumes`: the facts are different.

Important negative test: absence of a requirement for cached inspection is not
proof that it works offline. Its supported mode and necessary local data must
be described separately. Runtime dependency checks also need a defined treatment
of unavailable and unknown conditions.

### L03/L05: present knowledge versus past execution

The canonical model should store observations explicitly. A candidate sentence
template, outside the initial grammar, is:

```text
observation <id> reports <subject>.<property> = <value> at <time> from <evidence>.
```

Illustrative candidates:

```text
observation Obs42 reports ValidationRun42.phase = completed at T42 from Trace42.
observation Obs43 reports ValidationRun42.result = rejected at T43 from Trace43.
observation Obs44 reports EditRun9.observedOutcome = unknown at T44 from ClientTrace9.
```

The run definition must bind its activity, input subject/revision and validation
basis. Times and evidence identifiers must resolve to typed records. Rejection
can be a result of completed validation; `unknown` is the client's knowledge,
not an asserted terminal execution phase. An observation is a sourced claim,
not automatically accepted truth merely because it has this syntax.

Reject `B was validated`, `B is ready` and `the edit failed` as canonical substitutes:
they omit the result, relevant revision, guard or evidence. Natural-language
narratives may explain the explicit records but cannot create another encoding.

This case remains **semantically incomplete** for automatic validation until
run/observation identity, property types, time scope and conflicting-observation
rules are defined. Recording that gap is a result of the language exercise.

## 5. The first detailed flow exercise: L03 plus L04

Use layout replacement while live readings continue. The minimum facts that
must survive every representation are:

1. The request identifies candidate B and Composition C.
2. Validation checks B against the identified structural revision of C.
3. A successful validation result is tied to those revisions.
4. Publication requires that the relevant structural revision still matches at
   the declared commit point; checking only at activity entry is insufficient.
5. A mismatch follows an explicit rejection/revalidation branch. Which branch
   to adopt is a design choice, not a synonym the translator may choose freely.
6. Measurement revisions may advance concurrently; they must not be confused
   with the Composition's structural contract revision.
7. Rejection preserves the previous valid Presentation and identity/command
   continuity, while still allowing legitimate measurement changes.
8. Publication and actual display are separate observable outcomes.

Required grammar questions include: how to name an activity's input bindings;
how to express a guard and its evaluation point; how to reference exact versus
current revisions; and how to express permitted changes versus preserved state.
Do not declare this flow canonical by hiding those meanings in English prose.

The commit condition describes required coherence. It does not imply a global
transaction across independent services. The design must identify the owner and
mechanism that can satisfy it, or record the realization as unresolved.

## 6. How a scenario validates the language

For each case, keep its source intent, necessary facts, canonical statements,
forbidden alternatives, expected interpretation and outstanding gaps together.

Evaluate the following separately:

| Check | Pass condition |
|---|---|
| Expression | All necessary distinctions fit registered constructs; no decisive meaning is hidden in prose. |
| Interpretation | Reading a statement produces the intended participants, scope and obligations without guessing. |
| Rejection | Synonyms, missing required qualifiers, type errors and ambiguous references receive a specific rejection reason. |
| Canonical form | Re-encoding a supported fact produces the identical canonical representation, including set order and qualifiers. |
| View preservation | Text, diagram and typed record retain identities and semantics; a filtered view identifies its deliberate omissions. |
| Change impact | A changed guard/route identifies affected model facts and source bindings, including unresolved links. |
| Independence | A reader/context without the author's explanation reaches the intended interpretation. |

For supported facts, the intended later property is:

```text
decode(encode(fact)) = fact
encode(decode(canonicalSentence)) = canonicalSentence
```

This requires a defined fact model and encoding rules. It does not mean a
diagram snapshot can recover facts intentionally hidden by that view, nor that
we have proved semantic equivalence for arbitrary programs.

The present work is a manual design of cases and expected outcomes. No independent
reader trial, round-trip engine or formal rejection validator has yet passed
these checks. The worked examples already expose gaps; they are not advertised
as completed conformance tests.

## 7. Preventing dialects as the language grows

Maintain one versioned vocabulary/profile with canonical terms, signatures,
qualifiers, defaults, ordering and examples. Projects instantiate that profile
with domain identities; they do not redefine `owns`, `requires` or `completed`.
Domain vocabulary can add a concept such as APT without changing core semantics.

Before adding a word, demonstrate the missing distinction through a case. Then
check whether the proposal duplicates a relation/property already present,
update positive and negative cases, and reconcile existing examples. Old syntax
can have a deliberate migration path, but should not remain an undocumented
second way to write current models. Incoming colloquial prompts are interpreted
at the boundary; aliases do not leak into canonical storage.

Do not compress sentences by omitting scope needed for correctness. Prefer shared
named contracts/conditions and references to long repeated prose. Compactness
should reduce duplicated meaning, not require readers to guess it.

The immediate work sequence is L01/L02 for structural precision, L03/L04 for
flow/state semantics, then L05–L08 to challenge failure, identity and realization.
Only after those cases stabilize should we select parser/model tooling. A product
implementation pilot can follow under its own authorization and acceptance scope.
