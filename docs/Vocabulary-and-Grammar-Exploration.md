# Exploration: From vocabulary to a design grammar

Date: 2026-09-15

Current language definition: [Design language definition, design-core 0.1](Design-Language-Definition.md).
That document is the primary working definition. The vocabulary and syntax here
preserve the exploration that led to it; they do not define alternate core forms.

Status: grammar proposal for discussion. The owner has agreed with the preceding
capability/functionality and scenario/state refinements. The sentence forms,
qualifier definitions and grammar fragment below are new proposals, not an
adopted language, parser or execution model.

Read with [the vocabulary study](Feature-Functionality-and-Channel-Study.md) and
[behavior-to-source traceability](Scenarios-State-and-Implementation-Traceability.md).
The [MVP1 worked example](MVP1-Design-Language-Example.md) applies this vocabulary
to the overall system and then opens the selected UI design into detailed
responsibilities, scenarios and implementation bindings.

Subsequent owner constraint: the canonical design language must have one
expression per modeled meaning, without dialects or synonyms. The
[conformance scenarios](Design-Language-Conformance-Scenarios.md) now guide its
development. Earlier natural-English alternatives in this exploration illustrate
distinctions; they are not multiple permitted canonical sentence forms.

## 1. Three things a useful language needs

**Vocabulary** defines the concepts. **Grammar** defines how expressions combine.
**Semantics** defines what each valid expression means, including its scope and
consequences. A grammatically valid sentence can still describe an invalid design
or make a claim unsupported by the implementation.

The grammatical analogy is useful if it produces precise questions and
relationships. It should not require the owner to write formal sentences.
The agent translates an ordinary short prompt into reviewable model statements,
keeps uncertainty visible, and links the statements to design and code evidence.

Prefer a small controlled English expression of shared model concepts first.
The same meaning could later have textual, diagrammatic and structured-data
representations. English word endings should not become the identity of a model
object or its only machine-readable semantics.

## 2. Word classes and their modeling roles

| Grammar concept | Proposed modeling role | Examples |
|---|---|---|
| Common nouns | Types of design objects. | Unit, Container, Capability, Functionality, Activity, Contract, State, Event. |
| Proper names | Identified objects, distinct from their types. | UIHost, ValidateSchema, LayoutReplacement. |
| Relational verbs | Structural, dependency and realization facts. | owns, provides, consumes, requires, realizes, refines. |
| Action verbs | Operations performed by a responsible actor/unit in a behavior. | validate, publish, subscribe, reject, persist. |
| Adjectives | Defined properties, constraints or classifications of an object. | stateful, deterministic, idempotent, replaceable. |
| Adverbial phrases | Manner, scope, timing or frequency of a behavior. | within a declared deadline; once per request identity; concurrently with measurement processing. |
| Prepositions | Roles and context in a relationship. | through a contract; from a provider; to a consumer; within a Container; against revision r. |
| Quantifiers | Cardinality and scope of rules. | each, exactly one, at least one, no. |
| Modal verbs | Ability, obligation, permission or recommendation. | can, must, may, should, must not. |
| Conjunctions | Logical composition or flow composition, explicitly distinguished. | and, or, if, then, otherwise, while. |
| Tense and aspect | Temporal descriptions of identified executions and observations. | validates; is validating; has completed validation; is scheduled to validate. |

A noun can name behavior: `Validation` may identify an Activity or Functionality.
That is not a contradiction. Distinguish the definition from an occurrence:
`ValidateSchema` is reusable behavior; `ValidationRun42` is one execution of it.
Similarly distinguish an event type from an event occurrence, and a Container
definition from a deployed instance.

## 3. Verbs should have signatures

Instead of permitting any noun–verb–noun combination, specify participant types
and the meaning of each relation. An initial set could be:

| Relation pattern | Meaning / constraint |
|---|---|
| Unit owns Functionality | One immediate accountable owner per model revision; Container containment does not duplicate ownership. |
| Functionality realizes Capability | Local behavior contributes to an ability; the link alone does not prove that all necessary contributions exist. |
| Unit provides Capability | An offer at a declared boundary, completed by its interface/contract and conditions. |
| Unit consumes Interface | A collaborator actually uses that interface; required versus optional use is separate. |
| Capability requires Interface | The ability depends on a contract under stated modes/conditions, independent of a particular provider instance. |
| Activity refines Activity | One behavior is a more detailed realization of another, preserving its stated obligations. |
| Functionality binds-to SourceSymbol | A revision-specific implementation binding with a declared role and evidence. |

These are proposed signatures, not a complete relation inventory. Container is
a kind of Unit in this interpretation. Internal containment and deployment are
separate relationships. Capability dependencies may later need additional
signatures for resources or other capabilities; do not overload `requires` with
an unspecified target type.

Use one canonical direction and one canonical sentence form:
`UIHost owns ValidateSchema.` Passive `ValidateSchema is owned by UIHost` is not
a permitted alternative design sentence. This supersedes the earlier suggestion
to permit both renderings. Natural-language prompts may use either phrasing, but
the agent must resolve them to the same canonical fact before writing the model.
A relational statement
such as `UIHost provides LayoutReplacement` is not an instruction to run it.

For action verbs, specify actor, inputs, outputs/effects, affected state and
completion meaning. `Publish` is incomplete if we do not know whether it means
queued locally, delivered, applied, or acknowledged. Multi-participant relations
need named roles rather than an ambiguous series of arrows.

## 4. Adjectives can describe functionality without creating a class explosion

Use independent **facets** rather than a separate object type for every adjective
combination. A functionality can be stateful, deterministic and idempotent at
the same time. Pure/stateless, synchronous/asynchronous and durable/ephemeral
answer different questions.

| Facet | Candidate words | Required precision |
|---|---|---|
| Purpose | validating, transforming, coordinating, querying | A descriptive category; derive it from the actual responsibility, not the function name alone. |
| State | stateful, stateless | Retained state across which invocations, owned by whom? |
| Effects | pure, read-only, mutating | Effects on which domain; are logging, I/O and nondeterminism included? Stateless does not imply pure. |
| Repeatability | deterministic | Same explicit inputs and modeled initial state produce the same defined outcomes under stated environmental assumptions. |
| Repetition | idempotent | Repeating which operation with which identity preserves which effect? It need not return the same response or execute only once. |
| Completion | synchronous, asynchronous | At which interface, relative to which completion signal? This is often a property of an operation/binding rather than an entire unit. |
| Retention | durable, ephemeral | Which state/effect survives which failure or restart? Not a blanket adjective for an entire algorithm. |
| Exposure | internal, provided | Relative to which boundary and contract? |
| Currentness | fresh, stale, unknown | Relative to which revision, observation and freshness rule? Usually runtime data/state properties. |

Adjectives should map to defined properties or predicates and, where relevant,
testable obligations. Words such as fast, robust or safe are incomplete until
their acceptance meaning is specified. Not every property belongs on
Functionality: attach it to the state, operation, interface or delivery policy
whose behavior it actually constrains.

## 5. Tense, aspect, modality and execution state

The owner's tense idea is useful for readable descriptions, but it spans several
independent dimensions:

| Sentence | Intended interpretation |
|---|---|
| UIHost validates schemas. | General responsibility in a particular design revision; no execution is asserted. |
| UIHost can validate this schema. | An ability claim under specified conditions; no permission or execution is implied. |
| UIHost must validate before publication. | An obligation governing a flow; not evidence that validation happened. |
| UIHost may publish this candidate. | Permission or optional behavior, with its scope explicitly defined. |
| Run42 is scheduled to validate candidate B. | Planned execution; cancellation or failure can prevent it. |
| Run42 is validating candidate B. | Observed execution phase as of a stated observation. |
| Run42 completed validation of B revision 7. | A historical completion occurrence; its result still needs to be stated. |
| B revision 7 passed validation against contract C revision 3. | A result with an explicit subject, revision and validation basis. |

Do not infer `B is currently valid` merely from a past validation event. B or C
may have changed, authority may have expired, or publication may require another
guard. Completed validation can also produce a rejection. “Will validate” is
ambiguous between a prediction, a requirement and a plan; use the precise form.

Under the owner's canonical-language constraint, tense variations above remain
explanatory prose. A canonical runtime statement must explicitly identify its
observation, subject, phase/result, time/revision and evidence; prose tense alone
must not supply a second encoding of those facts.

For a pilot, store the behavior reference, run identity, lifecycle phase, result,
subject/contract revisions and time/evidence explicitly. Proposed lifecycle
labels might include planned, running, completed, failed and cancelled, but
they need a transition contract before becoming normative states. Keep a client's
unknown outcome separate from the execution's actual lifecycle phase.

Keep three axes separate:

1. Design authority: proposed, owner-selected, superseded.
2. Implementation evidence: planned binding, inspected source, tested candidate.
3. Execution: this run's phase and observed outcome.

Past tense does not make a design accepted, and an accepted design does not make
its described future behavior implemented. Timestamps also do not establish a
total causal order across independent processes; use the relevant stream/causality
and revision identities.

## 6. Other grammar features with practical value

**Quantification:** “Each Functionality has exactly one immediate owner” has a
different meaning from “Some Functionality has an owner.” Quantifier scope must
be explicit. “Each consumer uses a provider” does not establish that all consumers
use the same provider.

**Logic and flow:** Boolean `A and B` says both predicates hold; it does not mean
execute A then B. Use explicit sequence, choice, fork/join and waiting constructs
for behavior. Similarly, `after` establishes order, not causation. Parallel
branches require a completion/join policy.

**Conditions:** `when` can introduce a trigger, `if` a guard, `unless` an exception,
and `while` a maintained condition or concurrent interval. Choose one precise
meaning per construction; do not import all ambiguities of ordinary English.

**Negation and unknowns:** “No completion event is recorded” is not “the operation
did not complete.” An omitted relation is not an explicit prohibition. Identify
where a model is deliberately complete enough for absence to have meaning.

**References and scope:** pronouns such as “it” or “the current one” are convenient
in prompts but should resolve to explicit identities in the model. Scope inputs,
state variables and revisions within the appropriate activity/run/unit.

**Provenance and rationale:** `because`, `according to` and `observed at` can expose
reason, authority and evidence links. They are metadata about a statement, not
automatic proof that it is correct. Preserve the source and status rather than
turning an agent's explanation into an owner decision.

## 7. A small grammar fragment to discuss

This EBNF-style fragment covers declarations and structural facts only. It is
not a full grammar for scenarios, predicates, temporal rules or observations;
no parser has been implemented or exercised.

```text
model         = { declaration | relation | qualification } ;
declaration   = kind, identifier, "." ;
relation      = identifier, relationVerb, identifier, "." ;
qualification = identifier, "has", property, "=", value, "." ;

kind          = "unit" | "container" | "functionality" | "capability"
              | "interface" | "activity" | "source-symbol" ;
relationVerb  = "owns" | "realizes" | "provides" | "consumes"
              | "requires" | "refines" | "binds-to" ;
```

For this fragment, identifiers and property/value names are simple symbolic
tokens. Exact lexical rules, namespaces and literal types remain to be defined.
Object existence, relation signatures, allowed facets and cardinality are
semantic checks, not established merely by parsing the sentence.

An illustrative design fragment, not an assertion about deployed MVP1:

```text
unit UIHost.
functionality ValidateSchema.
capability LayoutReplacement.

UIHost owns ValidateSchema.
ValidateSchema realizes LayoutReplacement.
UIHost provides LayoutReplacement.
ValidateSchema has repeatability = deterministic.
```

This states a contribution and a property claim; it does not fully define or
prove LayoutReplacement. Its contract, other contributions and evidence remain
necessary. At this stage, rule/flow templates can stay separately readable:

```text
RULE: Publication must use a candidate validated against the current contract revision.
FLOW: Prepare candidate, then validate; on rejection retain the active presentation.
OBSERVATION: Run42 rejected candidate B revision 7 against contract C revision 3.
```

These three sentences describe obligation, designed behavior and an illustrative
observation respectively. Their labels prevent confusing a desired flow with
evidence of execution. They are not accepted syntax of the fragment above.

## 8. Existing work to reuse

OMG's **SBVR** is particularly relevant to this question. It separates noun
concepts, verb concepts with participant roles, and propositions/rules, while
distinguishing meaning from its expression. It is a conceptual reference for
typed statements, not an implementation of our complete behavior/source model.
Relevant sections are 8.1 and 11.1–11.2.
[SBVR 1.5](https://www.omg.org/spec/SBVR/1.5/PDF).

**Attempto Controlled English** demonstrates restricted English with nouns,
verbs, adjectives and adverbs. The inspected ACE 6.0 rules use simple present
tense; it should not be presented as already supporting the proposed temporal
wording. The useful lesson is controlled interpretation, not adopting all its
syntax or limitations. [ACE construction rules](https://attempto.ifi.uzh.ch/site/docs/ace/6.0/ace_constructionrules.html).

**RFC 2119/8174** provide an established convention for requirement strength.
If adopted explicitly, uppercase MUST, SHOULD and MAY have defined normative
meanings. Keep ability (`can`) separate from optionality/permission (`may`). This
study does not grant new authority by displaying those words in examples.
[RFC 2119](https://www.rfc-editor.org/rfc/rfc2119.html),
[RFC 8174](https://www.rfc-editor.org/rfc/rfc8174.html).

UML/SysML and the earlier state-machine research still supply behavior modeling
concepts. Controlled sentences would express the same model facts and link to
those behavior views, not require an independent duplicate of each diagram.

## 9. Next evaluation and limits

Use the [conformance scenarios](Design-Language-Conformance-Scenarios.md) as the
next bounded task. The owner selected scenario-driven language validation, not
MVP1 implementation or immediate parser construction.

Take a small set of statements about the selected MVP1 UI direction. For each,
record subject type, relation/action, participants, qualifiers, authority and
evidence. Include deliberate ambiguities: asynchronous relative to what,
validated which revision, requires in which mode, and “may” as possibility versus
permission. Test that a reader or agent exposes those ambiguities instead of
silently assigning meaning.

Try expressing the same statement as controlled English, a graph and a typed
record; check whether they preserve identity and semantics without duplicate
maintenance. Then test one scenario branch and its source bindings. Select
encoding/tooling only after that comparison identifies what is actually needed.

Primary references were inspected and local links/fence pairing checked. No
grammar parser, model validator or behavioral pilot was run. Product code and
skill candidates remain unchanged. This is an incremental language exploration,
not a replacement for vNow's installed process or an adoption of vNext tooling.
