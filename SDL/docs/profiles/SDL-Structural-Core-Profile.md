# SDL structural core — design-core 0.5

Working profile, implemented by the [Go frontend](../../go/README.md).
Profile date: 2026-09-22. Authority separation: R3-M2, 2026-09-24.
This moves the former definition's sections 1–11 without adding grammar rules.
The registered [data](SDL-Data-Contract-Profile.md),
[Channel](SDL-Channel-Scenario-Profile.md) and
[delivery-plan](SDL-Delivery-Plan-Profile.md) additions complete design-core 0.5.
Execution and classes use separate action-core/class-core profiles.

This is the active bounded working definition, not approval of all research
proposals or proof of implemented domain behavior. Open questions and future
syntax now have their own [candidate document](../studies/SDL-Language-Candidates.md).
The [former address](../studies/Design-Language-Definition.md) remains a navigation
map for historical references, not a second definition.

## 1. Authority and scope

For new canonical design sentences, this document takes precedence over the
exploratory syntax in [the grammar study](../studies/Vocabulary-and-Grammar-Exploration.md)
and [the MVP1 example](../studies/MVP1-Design-Language-Example.md). Their source reasoning
remains useful; an earlier example does not add an alternative language form.
[Conformance scenarios](../studies/Design-Language-Conformance-Scenarios.md) supply the
questions and expected rejection cases used to develop this definition.

The owner requires one unambiguous meaning and one canonical representation per
modeled fact, without synonyms for relations. The 2026-09-17 clarification permits
optional explicit type annotations in future source syntax ([candidate section 15](../studies/SDL-Language-Candidates.md#15-explicit-type-annotations-and-checking-without-annotations)); these
do not create different semantic facts or depend on neighboring sentences.
Version 0.5 defines a bounded structural core: names/types,
containment, ownership, capability contribution/offers, interface use, mode-scoped
dependency, activity refinement, Actor/UseCase/Feature traceability, mode-scoped
Functionality allocation and two Functionality properties. [Candidate section 12](../studies/SDL-Language-Candidates.md#12-open-extensions-and-conformance-work)
lists unsupported constructs. Unsupported syntax is rejected, not guessed.

This draft is the single working definition; it does not imply owner acceptance
of every new syntax decision. Changes to the definition must update its examples
and conformance cases together. The language version is independent of SDP's
product version, skill versions and MVP1's release version.

## 2. Word class, model type and sentence role

These are three different classifications:

| Classification | What it answers | Example |
|---|---|---|
| Grammatical class | What kind of expression is this? | `PresentationManager` is an identifier used as a proper name, occupying a noun-phrase position. |
| Model type | What kind of design object does the name identify? | `PresentationManager` identifies a `Unit`. |
| Sentence role | Which argument position does it fill in this statement? | Subject in `PresentationManager owns ValidateBindings.`; object in `UIHost contains PresentationManager.` |

`Unit` and `Functionality` are noun concepts/model types. `unit` and
`functionality` are their declaration keywords. `owns` is a registered relational
verb. `stateful` is an adjective-like property value; it is not an object name.

A name such as `ValidateBindings` looks like a verb phrase in ordinary English.
In a declaration it names a Functionality and therefore acts as a proper name
in design sentences. The language does not infer types from naming conventions.

Object declarations establish both the name's grammatical class and its model
type. Do not repeat `noun` beside every declaration: that would duplicate a fact
already defined by the declaration rule. A future vocabulary extension must
register the word class and typed usage of any new kind of expression centrally.

Subject and object are not permanent attributes of an identifier. Both Unit and
Functionality names can occupy either role **where the relation signature permits
their type**. A concrete statement's role is determined by its position.

## 3. Lexical rules and declarations

Each model starts with the exact header:

```text
language design-core version 0.5.
```

Then declare every identifier before the statements that use it:

```text
unit PresentationManager.
functionality ValidateBindings.
```

Rules:

- Identifiers match `[A-Z][A-Za-z0-9]*`. They are case-sensitive ASCII names.
- Keywords and relation/property/value tokens use the exact lowercase spelling
  defined here. No alternative capitalization, passive forms or synonyms exist.
- One identifier has exactly one declaration and one declared type in a model.
- A declaration identifies a design object, not a runtime instance or a run.
- Version 0.5 has one namespace per model. Imports and cross-model name resolution
  are unsupported; references cannot silently resolve to another file's names.
- A Container is declared with `container`, not once as `unit` and again as
  `container`. Its Unit compatibility follows from the type system.
- Statements end in `.`. Canonical layout is specified in Section 9.
- Descriptive labels and comments may surround a canonical model in Markdown;
  they are not statements inside the model and cannot supply missing semantics.

### 3.1 Symbol table

Declarations conceptually create this lookup table:

| Identifier | Grammatical class | Declared type | Compatible with |
|---|---|---|---|
| UIHost | Proper-name identifier | Container | Container, Unit |
| PresentationManager | Proper-name identifier | Unit | Unit |
| ValidateBindings | Proper-name identifier | Functionality | Functionality |
| LayoutReplacement | Proper-name identifier | Capability | Capability |

This is how the language learns what an unfamiliar user-chosen name means.
A name does not become an adjective or a different model type because it occurs
in a different sentence position.

## 4. Type vocabulary

| Declaration | Model type | Meaning in this core |
|---|---|---|
| `unit Name.` | Unit | Bounded design responsibility; a deployment claim is not implied. |
| `container Name.` | Container | Unit explicitly identified as an application/data-store runtime boundary. |
| `functionality Name.` | Functionality | Cohesive local behavioral responsibility, with one immediate accountable Unit owner. |
| `capability Name.` | Capability | An ability whose conditions and realization may be detailed by the wider design. |
| `interface Name.` | Interface | A named collaboration boundary; the declaration alone is not a complete contract. |
| `activity Name.` | Activity | A named behavior that may have a more detailed realization. This core does not define its execution. |
| `actor Name.` | Actor | External stakeholder role pursuing a UseCase; not an owning Unit. |
| `usecase Name.` | UseCase | Named observable goal, not a scenario or ordered execution. |
| `feature Name.` | Feature | Named offered behavior supporting UseCases; distinct from Capability. |
| `mode Name.` | Mode | A named operating context used to scope a dependency or allocation. Its runtime activation predicate is outside this core. |

The only subtype rule is `Container <: Unit`. No other implicit type conversions
are allowed. For example, a Capability cannot act as an owning Unit. An Activity
does not become a Functionality merely because both concern behavior.

Functionality's architectural boundary remains local to one Container where
runtime allocation is known. This core can describe units before allocation is
complete; passing its checks does not prove complete runtime allocation. Shared
library definitions retain one logical owner. Explicit Functionality allocation
identifies one Container per Mode; instance multiplicity remains future work.

## 5. Relation vocabulary and argument signatures

The general syntactic shape is **subject, verb, object**, with a fixed qualifier
where required. The semantic signature supplies the allowed participant types.

| Canonical construction | Subject type | Object type | Meaning and constraints |
|---|---|---|---|
| `A contains B.` | Unit | Unit | Immediate logical containment. Not deployment or a claim of state ownership. Each contained Unit has at most one immediate parent. No self-containment or cycles. |
| `A owns B.` | Unit | Functionality | A is the immediate accountable owner of B. Each declared Functionality has exactly one owner. |
| `A realizes B.` | Functionality | Capability | A contributes to realizing B. Several contributions may exist; the link does not assert that their set is complete or verified. |
| `A provides B.` | Unit | Capability | A declares an offered ability at its boundary. Contract details and verification are not established by this fact alone. |
| `A consumes B.` | Unit | Interface | A is designed to use B. This is not an observation that a request has occurred, nor a statement that B is mandatory. |
| `A requires B in mode M.` | Capability | Interface | B is necessary for A under the identified Mode M. M must be declared as Mode. It does not assert that A is possible whenever B is present. |
| `A pursues B.` | Actor | UseCase | A pursues goal B. Many-to-many; no implicit actor or system boundary. |
| `A supports B.` | Feature | UseCase | A supports B. Many-to-many, without a completeness claim. |
| `A contributes-to B.` | Functionality | Feature or UseCase | Direct contribution, including a direct goal contribution without a Feature. |
| `A allocated-to B in mode M.` | Functionality | Container | Explicit runtime location in declared Mode M. At most one Container per Functionality/Mode. Does not change logical ownership. |
| `A refines B.` | Activity | Activity | A is a more detailed realization of B. Direction is detailed to abstract. No self-refinement or cycles. Behavioral preservation needs later evidence. |

Verb signatures can be written compactly as:

```text
contains : Unit × Unit
owns     : Unit × Functionality
realizes : Functionality × Capability
provides : Unit × Capability
consumes : Unit × Interface
requires : Capability × Interface × Mode
refines  : Activity × Activity
pursues  : Actor × UseCase
supports : Feature × UseCase
contributes-to : Functionality × (Feature | UseCase)
allocated-to   : Functionality × Container × Mode
```

These signatures are definitions in this document, not alternate model syntax.
For `requires`, Mode fills a third argument expressed by the fixed `in mode`
qualifier. Mode names are model-global identities in 0.2; use distinct identities
for distinct contexts. There is no implicit universal Mode, no missing-mode
default, and no supported unconditional dependency shorthand in this core.

`contains` and `owns` are deliberately different. A parent Unit does not receive
a second direct `owns` fact for every Functionality of its child. An ancestor
view can derive indirect containment without storing those derived edges as
new immediate relations. Provided capabilities do not automatically propagate
up the containment hierarchy.

### 5.1 V1 — goals, contributions and allocation

`pursues`, `supports` and `contributes-to` preserve many-to-many relationships. Direct Functionality contributions to UseCase require no Feature. Features may lack contributions and UseCases may lack Actors: incomplete but structurally valid. Viewpoints expose gaps without inventing intermediaries. Feature is neither a Capability subtype nor synonym.

`allocated-to` is stored as `Allocation(subject, container, mode, span)` in AST, with source positions for all three identifiers. Multiple distinct Containers for the same Functionality/Mode are rejected with `ALLOCATION_CARDINALITY`. A Functionality may be allocated differently across modes; these are alternative contexts, not concurrent instances. Copies/concurrent realizations need separate Functionality identities. Duplicate facts are rejected normally; absent allocation is unspecified. No allocation inherits from `contains`, ownership, names or another Mode. Ownership is logical responsibility; allocation is execution context. This profile implements no System, deployer, processes, startup ordering, mode activation or instance binding.

VP01 shows declared goals and direct/Feature-mediated contributions as labeled flowcharts, split into use-case and Feature-contribution maps for readability. Model scope defines view scope; no formal System boundary is invented. VP07 produces a Feature selection per explicit allocation mode, with contributions, logical owners and Container placements. Unallocated contributions remain visible and are listed as gaps for that mode. Every arrow maps to exactly one source fact; diagram titles/fact registers identify modes without implicit inheritance.

0.2 replaces the active 0.1 profile. Local parser tests and SDL/SDUI models were ported; no old parser/fallback remains. Historical checkpoint/MVP1 sources are not automatically ported or declared compatible.

## 6. Adjectives as typed properties

Adjective-like values use one canonical property construction:

```text
ValidateBindings has repeatability = deterministic.
```

| Property | Allowed subject | Allowed values | Meaning |
|---|---|---|---|
| `state-retention` | Functionality | `stateful`, `stateless` | Whether the responsibility retains its own state across invocations. This says nothing by itself about external state or side effects. |
| `repeatability` | Functionality | `deterministic`, `nondeterministic` | Whether the same explicit inputs, initial modeled state and environmental assumptions constrain the specified outputs/effects to the same result. Nondeterministic permits variation; it is not shorthand for untested. |

Each subject/property pair may have at most one value. Missing properties mean
unspecified, not false and not a default. Assigning a value is a design claim;
no property assignment proves the implementation satisfies it.

The property vocabulary is closed. Reject `owner-of`, `fast`, `async` or arbitrary
project-specific keys. They cannot encode an existing relation a second way.
Useful concepts such as idempotence, completion, freshness and durability remain
pending until their scope/conditions and validation rules are defined.

Reject adjective-prefix alternatives such as `deterministic functionality
ValidateBindings.` They would create a second way to encode the same property.

V2 extends the closed type/relation/property tables through the [normative data/wire profile](SDL-Data-Contract-Profile.md), included in 0.3. Dataset/Datagram/Database are implemented within that scope.

V3's [Channel/scenario profile](SDL-Channel-Scenario-Profile.md) joins 0.4. It bounds the broader section 13–15 candidates; they do not automatically become language rules.

V4's [delivery-plan profile](SDL-Delivery-Plan-Profile.md) joins 0.5, adding addresses/delivers/depends-on and implementation-status to relations/Activity properties. These are plan facts, not runtime execution.

## 7. Core productions and registered profile extensions

The linked data, Channel and delivery-plan profiles are normative extensions to these core productions. Their type/verb/property registries extend the bounded tables below. Only the combined 0.5 profile is active.

```text
model = header, { declaration }, { statement } ;

header = "language", "design-core", "version", "0.5", "." ;

declaration = kind, identifier, "." ;
kind = "unit" | "container" | "functionality" | "capability"
     | "interface" | "activity" | "mode" | "actor" | "usecase" | "feature" ;

statement = binaryRelation | dependency | allocation | propertyAssignment
          | projection | placement | participation | step ;
(* V2 kinds, binary verbs, properties and productions: SDL-Data-Contract-Profile.md *)

binaryRelation = identifier, binaryVerb, identifier, "." ;
binaryVerb = "contains" | "owns" | "realizes" | "provides"
           | "consumes" | "refines" | "pursues" | "supports" | "contributes-to" ;

dependency = identifier, "requires", identifier, "in", "mode", identifier, "." ;

allocation = identifier, "allocated-to", identifier, "in", "mode", identifier, "." ;

propertyAssignment = identifier, "has", property, "=", propertyValue, "." ;
property = "state-retention" | "repeatability" ;
propertyValue = "stateful" | "stateless" | "deterministic" | "nondeterministic" ;

identifier = uppercaseLetter, { letter | digit } ;
uppercaseLetter = "A" | "B" | "C" | "D" | "E" | "F" | "G" | "H" | "I"
                | "J" | "K" | "L" | "M" | "N" | "O" | "P" | "Q" | "R"
                | "S" | "T" | "U" | "V" | "W" | "X" | "Y" | "Z" ;
letter = uppercaseLetter | lowercaseLetter ;
lowercaseLetter = "a" | "b" | "c" | "d" | "e" | "f" | "g" | "h" | "i"
                | "j" | "k" | "l" | "m" | "n" | "o" | "p" | "q" | "r"
                | "s" | "t" | "u" | "v" | "w" | "x" | "y" | "z" ;
digit = "0" | "1" | "2" | "3" | "4" | "5" | "6" | "7" | "8" | "9" ;
```

The grammar permits identifiers in subject/object positions. **Type checking**
then determines whether their declared types fit that verb. A grammar without
the symbol table and signatures cannot reject reversed ownership merely from
its sentence shape. Property/value compatibility is also a semantic check:
`state-retention = deterministic` has the token shape but the wrong value type.

## 8. Validation stages and diagnostics

A conforming implementation must distinguish these checks:

1. Recognize the supported language version and token spellings.
2. Parse declarations and sentence structure.
3. Build the symbol table; reject duplicate declarations and unresolved names.
4. Check every argument and property against its registered type signature.
5. Check uniqueness, ownership cardinality and containment/refinement cycles.
6. Check canonical serialization, including duplicate facts and ordering.

The prototype implements the following diagnostic identifiers. Diagnostics carry
source spans; syntax, type and structural errors remain distinct.

| Diagnostic | Example / reason |
|---|---|
| `UNSUPPORTED_SYNTAX` | `ValidateBindings is owned by PresentationManager.` uses an unregistered form. |
| `UNSUPPORTED_VERSION` | The header requests a version other than `0.5`. |
| `UNDECLARED_NAME` | A statement references `Validator` without a declaration. |
| `DUPLICATE_DECLARATION` | One identifier is declared twice, even with the same type. |
| `SUBJECT_TYPE_MISMATCH` | `ValidateBindings owns PresentationManager.` expects Unit, receives Functionality. |
| `OBJECT_TYPE_MISMATCH` | The same statement expects Functionality as object, receives Unit. |
| `QUALIFIER_TYPE_MISMATCH` | The `in mode` reference names a Unit rather than a Mode. |
| `PROPERTY_TYPE_MISMATCH` | A Unit receives a Functionality-only property, or a property receives the wrong value kind. |
| `OWNERSHIP_CARDINALITY` | A Functionality has no owner or two different direct owners. |
| `ALLOCATION_CARDINALITY` | A Functionality is allocated to two Containers in the same Mode. |
| `CONTAINMENT_CARDINALITY` | A Unit has two different immediate parents. |
| `STRUCTURE_CYCLE` | Containment or refinement points to itself or forms a cycle. |
| `DUPLICATE_FACT` | The exact relation, or the same property assignment, is repeated. |
| `PROPERTY_CONFLICT` | One subject/property has two different values. |
| `NONCANONICAL_FORM` | Recognizable statements use noncanonical spacing or serialization order. |

Given the declarations in Section 3, reversed ownership has two argument errors.
The message should identify the positions, expected types and declared types.
Do not silently reverse the sentence: that would change the submitted design.

Passing these checks means **structurally valid design-core 0.5**, not correct
architecture, complete contracts, implemented behavior or verified product safety.

## 9. Canonical representation

The canonical document has the header first, declarations second, statements
third. Declarations are sorted by identifier using ASCII order. Statement lines
are sorted by their complete canonical ASCII text. Structural statements are an
unordered fact set, so this order does not express execution sequence.

Use one space between tokens, except no space before the final period. One
sentence occupies one line. Use LF line endings, one final newline, no blank
lines and no comments inside the canonical model. Duplicate facts are rejected,
not silently dropped. Display tables/diagrams may reorder facts for readability
but must retain identity and cannot be passed off as canonical source text.

This chooses one serialization for the same explicitly identified fact set.
It does not establish that differently named objects or different realizations
are semantically identical. Identity reconciliation is deliberate model work.
Behavioral sequence will need a different explicit construct; sorting future
workflow steps would not be a valid canonicalization rule.

## 10. Complete examples

### 10.1 Ownership and both sentence positions

This example is a logical decomposition. Declaring UIHost as Unit does not
settle whether it will be a separate Go runtime Container.

```design-core
language design-core version 0.5.
capability LayoutReplacement.
unit PresentationManager.
unit UIHost.
functionality ValidateBindings.
PresentationManager owns ValidateBindings.
UIHost contains PresentationManager.
UIHost provides LayoutReplacement.
ValidateBindings has repeatability = deterministic.
ValidateBindings realizes LayoutReplacement.
```

`PresentationManager` is subject of `owns`, object of `contains`.
`ValidateBindings` is object of `owns`, subject of `realizes` and `has`.
Both names act as noun phrases throughout. Their different model types explain
why swapping them around `owns` is invalid.

The deterministic property is illustrative, not a verification claim about MVP1.
Only one contributing Functionality is shown; full capability realization and
external interface exposure remain to be described.

### 10.2 Scoped consumption and dependency

```design-core
language design-core version 0.5.
capability LiveInspection.
interface ObservationInterface.
mode OperatorLive.
unit OperatorUI.
LiveInspection requires ObservationInterface in mode OperatorLive.
OperatorUI consumes ObservationInterface.
OperatorUI provides LiveInspection.
```

Consumption and necessity are different facts. The model does not claim cached
inspection works offline, nor that the whole application requires this interface
in every mode. The Mode is a named design context, not a runtime-state predicate
that this core can evaluate.

### 10.3 Negative sentence cases

With the declarations above, reject:

```text
ValidateBindings owns PresentationManager.
PresentationManager is responsible for ValidateBindings.
UIHost has owner-of = ValidateBindings.
LiveInspection requires ObservationInterface.
LiveInspection requires ObservationInterface in mode OperatorUI.
PresentationManager has repeatability = deterministic.
ValidateBindings has state-retention = deterministic.
```

Reasons, respectively: reversed types; unregistered synonym; relation disguised
as an unknown property; missing mandatory Mode; incorrect qualifier type;
incorrect property subject type; incorrect property value type. The lines are
independent negative examples, not one complete model.

### 10.4 Use cases, contributions and alternative execution contexts

Direct and Feature-mediated contributions are explicit facts. Ownership remains
unchanged when the same responsibility is allocated to different Containers in different modes.

```design-core
language design-core version 0.5.
actor Author.
container CommandLineHost.
usecase InspectModel.
mode Inspection.
container InteractiveHost.
mode LiveEditing.
feature StructuralChecking.
functionality ValidateModel.
unit Validator.
Author pursues InspectModel.
StructuralChecking supports InspectModel.
ValidateModel allocated-to CommandLineHost in mode Inspection.
ValidateModel allocated-to InteractiveHost in mode LiveEditing.
ValidateModel contributes-to InspectModel.
ValidateModel contributes-to StructuralChecking.
Validator owns ValidateModel.
```

Two distinct Containers in one Mode would fail with `ALLOCATION_CARDINALITY`.
Replacing Feature with Capability in `contributes-to` would fail with `OBJECT_TYPE_MISMATCH`.

## 11. Updating the other documents

This definition supersedes the exploration's unrestricted generic relation and
property forms for canonical core sentences:

- `contains` and Mode declarations now have explicit meanings and signatures.
- `requires` always carries the fixed `in mode` qualifier in this version.
- `has` accepts only the closed typed property table in Section 6.
- `binds-to`, `source-symbol` declarations and observation sentences remain
  proposed extensions; they are not supported core syntax.
- Active/passive alternatives, pronouns and informal tense variants are not
  canonical core sentences.

Research prose can discuss all these concepts. Only blocks explicitly claiming
`design-core 0.5` compliance are required to satisfy this definition. Earlier
blocks are exploratory fragments, not backward-compatible alternative syntax.
