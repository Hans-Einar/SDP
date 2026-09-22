# Design language definition

Version: **design-core 0.5 — draft**

Date: 2026-09-22

Go-port G4-M1 implementert 2026-09-22: [frontend og bevis](../SystemDesignLanguage/go/README.md).
Ingen nye språkregler innføres av porten.

V1 mål/bidrag/allokering implementert 2026-09-22. Seksjon 13 og videre beskriver
framtidig språkarbeid; gjeldende implementert grammatikk er `design-core 0.5`.

Status: primary working definition of the proposed language. It consolidates the
structural vocabulary and adds a precise grammar, type rules and canonical form.
It is not an installed SDP contract or a completed behavioral language. An isolated
[parser and validator prototype](../experiments/design_core/README.md) implements
this bounded core, including AST inspection and canonical formatting.

## 1. Authority and scope

For new canonical design sentences, this document takes precedence over the
exploratory syntax in [the grammar study](Vocabulary-and-Grammar-Exploration.md)
and [the MVP1 example](MVP1-Design-Language-Example.md). Their source reasoning
remains useful; an earlier example does not add an alternative language form.
[Conformance scenarios](Design-Language-Conformance-Scenarios.md) supply the
questions and expected rejection cases used to develop this definition.

The owner requires one unambiguous meaning and one canonical representation per
modeled fact, without synonyms for relations. The 2026-09-17 clarification permits
optional explicit type annotations in future source syntax (Section 15); these
do not create different semantic facts or depend on neighboring sentences.
Version 0.5 defines a bounded structural core: names/types,
containment, ownership, capability contribution/offers, interface use, mode-scoped
dependency, activity refinement, Actor/UseCase/Feature traceability, mode-scoped
Functionality allocation and two Functionality properties. Section 12
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

### 5.1 V1 — mål, bidrag og allokering

`pursues`, `supports` og `contributes-to` bevarer mange-til-mange-forhold.
Et direkte Functionality-bidrag til UseCase krever ingen Feature. En Feature
kan mangle bidrag, og en UseCase kan mangle Actor; modellen er da ufullstendig,
men strukturelt gyldig. Viewpoints viser hull uten å opprette mellomledd.
Feature er ikke subtype eller synonym for Capability.

`allocated-to` lagres som `Allocation(subject, container, mode, span)` i AST,
med posisjonerte Identifier-argumenter. Alle tre argumenter typekontrolleres.
Flere forskjellige Containers for samme Functionality/Mode avvises med
`ALLOCATION_CARDINALITY`. Samme Functionality kan plasseres forskjellig i ulike
modi; slike plasseringer er alternative kontekster, ikke samtidige instanser.
Kopier eller flere samtidige realiseringer trenger egne Functionality-identiteter.
Duplikatfakta avvises som ellers. Manglende allokering er uspesifisert.
Ingen allokering arves fra `contains`, eier, navn eller en annen Mode. Eierskap
er logisk ansvar; allokering er kjørekontekst. Denne profilen implementerer ikke
System, deployer, prosesser, startrekkefølge, modusaktivering eller instansbinding.

VP01 viser deklarerte mål og direkte/Feature-formidlede bidrag som en merket
flowchart, delt i bruksmålskart og Feature-bidragskart for lesbarhet.
Modellen er visningens omfang; ingen formell System-grense oppdiktes.
VP07 lager ett Feature-utsnitt per eksplisitt allokeringsmodus, med bidrag,
logiske eiere og Container-plasseringer. Uallokerte bidrag vises også og listes
som hull i den viste modusen. Alle piler tilsvarer nøyaktig ett kildefaktum;
modus fremgår av tegningstittel og faktaregister, uten implisitt arv.

0.2 erstatter den aktive 0.1-profilen. Lokale parserprøver og SDL/SDUI-modellen
er portert; ingen gammel parser eller fallback beholdes. Historiske checkpoint-
og MVP1-kilder er ikke automatisk portert eller erklært kompatible.

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

V2 utvider de lukkede type-, relasjons- og egenskapstabellene med den
[normative data-/wireprofilen](SDL-Data-Contract-Profile.md). Denne inngår i 0.3;
Dataset/Datagram/Database er nå implementerte begreper i denne avgrensningen.

V3s [Channel-/scenarioprofil](SDL-Channel-Scenario-Profile.md) inngår i 0.4.
Den avgrenser de videre kandidatene i seksjon 13–15; de blir ikke automatisk språkregler.

V4s [leveranseplanprofil](SDL-Delivery-Plan-Profile.md) inngår i 0.5. Den utvider
relasjonstabellen og Activity-egenskapene med addresses/delivers/depends-on og
implementation-status. Dette er planfakta, ikke runtime-utføring.

## 7. Core productions and registered profile extensions

Data-, Channel- og leveranseplanprofilene lenket ovenfor er normative tillegg
til disse kjerneproduksjonene; deres type-/verb-/egenskapsregistre utvider de
avgrensede tabellene nedenfor. Bare den samlede 0.5-profilen er aktiv.

EBNF notation: quoted strings are literal tokens, comma is concatenation,
`|` is choice, and braces mean zero or more repetitions. Whitespace separates
tokens; Section 9 fixes canonical serialization. Identifiers follow Section 3.

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

### 10.3 Bruksmål, bidrag og alternative kjørekontekster

Direkte og Feature-formidlet bidrag er begge eksplisitte fakta. Eierskapet er
uendret når samme ansvar plasseres i forskjellige Containers i ulike modi.

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

To forskjellige Containers i samme Mode ville feilet med
`ALLOCATION_CARDINALITY`. Å bytte Feature med Capability i `contributes-to`
ville feilet med `OBJECT_TYPE_MISMATCH`.

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

## 12. Open extensions and conformance work

| Scenario pressure | Missing definition |
|---|---|
| L01 complete capability offer | Interface exposure, contract conditions and complete-realization evidence. |
| L02 runtime applicability | Mode ownership/activation predicates and unconditional dependency representation. |
| L03 validation/publication race | Typed inputs, revisions, named Conditions, guard evaluation point and commit scope. |
| L04 concurrent updates | Sequence/fork/join, ordering, frame conditions and delivery-policy contracts. |
| L05 uncertain command outcome | Run, Event and Observation types; correlation, lifecycle/result distinctions and time semantics. |
| L06 quantities and shared views | Typed values/units and separate semantic, Representation and view identities. |
| L07 placement changes | Complete deployment and reusable-definition versus runtime-instance bindings; bounded mode-scoped responsibility allocation is implemented. |
| L08 implementation evidence | Revision-specific SourceSymbol and evidence types with role-specific bindings. |

Add each extension by specifying its concept, grammatical class, argument types,
semantics, canonical form and positive/negative cases here. An explanatory word
appearing in a scenario does not automatically become a language keyword.

The [SDL source-tree and compilation study](SDL-Source-Tree-and-Compilation-Study.md)
proposes future workspace scopes, multi-file linking, public exports and blueprint
generation. These require an explicit language extension; version 0.5 retains
its standalone-model semantics and does not resolve names from neighboring files.

Version 0.5 supports executable checks for ownership reversal, other core type
errors, and the structural subset of L01/L02. The prototype tests parse all three
complete examples above and check their semantics and canonical layout. Additional
tests cover invalid forms, structural constraints and structural text/AST round
trips. This is not a proof that the EBNF and implementation accept identical
languages. Full scenario expressiveness, diagram/source round trips and independent
language evaluation remain open. See the [prototype guide](../experiments/design_core/README.md)
for commands, AST structure and implementation limits.

## 13. Planned workspace extension: one declared System

Owner decision, 2026-09-17: the compilation entry is an explicitly selected source
file. `System.design` is a useful filename convention, not a required name. The
source declares the modeled System, and the first workspace language version
permits exactly one declared System per complete compilation.

This section records the next extension; it does not alter the `design-core 0.5`
grammar or claim parser support. The owner's root/cardinality decision is selected.
The type integration and canonical details below are proposed for that extension.

### 13.1 Declaration and typed containment

Add the model type `System`, declared by `system Name.`. The declared name is a
proper-name identifier, just like a Unit name; its type is System. File and
directory names neither declare nor rename it.

Keep one declaration and one binary fact per sentence. With the existing keyword
case and punctuation conventions, the owner's sketch becomes this future-language
fragment, with the future version header deliberately omitted:

```text
system MVP1.
container BuckingService.
container MachineService.
container StemService.
container UI.
MVP1 contains BuckingService.
MVP1 contains MachineService.
MVP1 contains StemService.
MVP1 contains UI.
```

These are illustrative identities, not a registration or deployment decision for
MVP1. In particular, StemService extraction and the UI's host/renderer allocation
need their own design facts. Container declarations occur once; a separate
boundary file references/enriches the declared identity rather than redeclaring it.

Recommend System as a distinct root type, **not a subtype of Unit**. Extend the
allowed argument pairs of `contains`:

```text
contains : (System × Container) | (Unit × Unit)
```

This is a union of two complete signatures, not the unrestricted Cartesian
product `(System | Unit) × (Container | Unit)`. Thus `MVP1 contains UI.` is valid,
while `UI contains MVP1.` and `MVP1 contains ValidateBindings.` are type errors.
A System has no containment parent. Existing unique-parent and cycle rules still
apply across the linked model. Containment describes membership, not deployment
on the same machine, invocation order or communication.

Do not make `MVP1 owns ValidateBindings.` legal by treating System as a generic
Unit. System-level capability offers and logical subsystem membership need their
own explicit signature decisions; this increment does not infer those permissions.

### 13.2 Entry and compilation rules

- Select the entry file explicitly. Renaming it from `System.design` to
  `Ponsse.design` preserves model identity and meaning.
- Put the single System declaration in that entry file. Other included files do
  not declare a separate System or repeat the root declaration.
- Count System declarations across the complete source set. Zero or more than one
  is an error, even when the extra declaration uses a different name. A repeated
  declaration of the same name is also a duplicate declaration error.
- The source set comes from explicit membership rules. `MVP1 contains UI.` does
  not implicitly load `UI.design` or every file under `UI/`.
- Check root/cardinality after collecting the declared source set. Parsing an
  individual layer file must not fail merely because that file has no System.
- An independently checked fragment is not a complete system build. External
  interface references do not require another System declaration.

For future canonical entry layout, recommend the version header followed by the
System declaration, then other declarations and facts in their defined order.
This is a versioned change from 0.2's declaration ordering. Comma-separated
declaration/containment lists are not an additional canonical form; the short
owner sketch can be translated into the individual facts above.

### 13.3 Acceptance cases for the extension

| Case | Expected result |
|---|---|
| Entry file with any filename declares one System; included layer files declare none | Root/cardinality check passes. |
| Entry filename or containing directory changes without changing declarations | Same System identity. |
| Source set has no System declaration | Missing-System error. |
| Two included files declare distinct Systems | Multiple-System error identifying both declarations. |
| Two included files declare the same System name | Multiple declarations rejected; do not silently merge roots. |
| The only System declaration is in a non-entry file | Entry-declaration error. |
| System contains a declared Container | Valid argument types. |
| Container contains System, or System owns Functionality | Invalid argument types. |
| A Container is declared in both the entry and its boundary file | Duplicate declaration error. |
| Entry contains a Container but omits the file defining it from the source set | Unresolved-name error; no implicit file discovery. |

Implement these cases with the explicit workspace/version extension. Existing
0.2 fixtures retain their current behavior, including standalone models without
a System declaration.

## 14. Proposed Channels, participants and message direction

Owner input, 2026-09-17: declare the system's Channels and participating Containers
in the entry file, including MachineService-to-BuckingService/UI observations and
UI-to-MachineService traffic. The owner explicitly supplied pseudocode, not a
request to accept competing `uses` forms. The following is a design proposal,
not an implemented extension or an adopted protocol for MVP1.

Subsequent owner decision: optional type annotations may precede the participant,
Channel and MessageSet references. Section 15 supersedes the earlier rejection of
`uses channel Name` and defines its precise meaning. Relation wording and argument
order remain fixed.

### 14.1 Keep the concepts separate

Retain **Channel** as the reusable logical communication arrangement from the
[Channel study](Feature-Functionality-and-Channel-Study.md#5-reusable-channels-without-accidental-coupling).
An Interface describes an exposed collaboration boundary. A Contract defines
allowed messages, meaning and obligations. A binding maps the logical interaction
to concrete communication mechanisms. A Channel is neither a Capability nor
necessarily a socket, queue or connection.

The entry file declares Channels and connects participants. Detailed contract
definitions may be in explicitly included files or authoritative external schemas;
placing topology at the root need not put every message field there. Each Channel
and contract definition still has one authoritative declaration. Generated
Container exports project the same relationships rather than restating them by hand.

Interpret the owner's examples as communication between Containers inside the
one declared System. Describing external interfaces remains possible in principle;
declaring additional Systems remains outside the selected first workspace scope.

### 14.2 Names are descriptive; roles define direction

`MachineServiceUpstream` can mean the collaboration toward the UI, and
`MachineServiceDownstream` the collaboration toward the target adapter. Those
names are valid design choices, but their position-relative meaning must not be
encoded in compiler heuristics. Moving a diagram node cannot reverse an edge.

Prefer semantic names when the content is known: `MachineObservations` for an
observation-only Channel, `MachineControl` for a command/result collaboration,
or `TargetCommunication` for the service/adapter protocol boundary. These are
examples, not synonyms registered for one modeled Channel. Select one declared
identity for each actual arrangement. A broad existing service contract need not
be split merely to use directional names.

A Channel can support multiple directions. A participant can send one class of
messages and receive another. Client/server, provider/consumer and sender/receiver
are not interchangeable: a provider can receive requests and send responses.
Likewise, a receive role does not imply the whole Container requires that Channel
in every operating mode or has authority to mutate domain state.

### 14.3 One candidate participation construction

Use one participation construction with optional type annotations:

```text
Participant uses ChannelName as sender of MessageSetName.
Participant uses ChannelName as receiver of MessageSetName.
```

The two lines are the shortest form of the construction. In this extension's
EBNF, square brackets mean an optional element:

```text
participation = reference, "uses", reference, "as", directionRole,
                "of", reference, "." ;
reference = [ typeKeyword ], identifier ;
typeKeyword = kind | "system" | "channel" | "message-set" ;
directionRole = "sender" | "receiver" ;
```

This is a proposed production only. Its semantic arguments are
`Container × Channel × DirectionRole × MessageSet`. `kind` refers to the type
keywords in Section 7. The grammar recognizes annotations; semantic checks verify
their compatibility with the argument position and the declaration. A known but
incompatible annotation is a type error, not an unknown keyword. Here `uses` has
the precise meaning of declared messaging participation, not a generic dependency
relation. `uses X as channel` remains invalid because `as` introduces a direction
role, not a type annotation.

`MessageSet` is a proposed named set of contract-defined message types. It is
not a runtime message, Event, role or free-text topic label. Its membership,
schema references and association with a Channel contract must be explicit before
full validation is possible. Every referenced name needs a declaration/resolution;
examples cannot silently invent messages because their names sound plausible.

Illustrative future fragment, using the owner's Channel names and omitting
Container declarations, version header and detailed contracts:

```text
channel MachineServiceUpstream.
channel UIDownstream.
message-set MachineCommands.
message-set MachineObservations.
BuckingService uses MachineServiceUpstream as receiver of MachineObservations.
MachineService uses MachineServiceUpstream as sender of MachineObservations.
MachineService uses UIDownstream as receiver of MachineCommands.
UI uses MachineServiceUpstream as receiver of MachineObservations.
UI uses UIDownstream as sender of MachineCommands.
```

The `of` qualifier is deliberately mandatory in this candidate. Omitting it when
a Channel happens to carry one message set would add a second canonical form and
make later contract growth change the meaning of the old role declaration.
Participation is a design claim, not evidence of a transmission or a delivery.

This is consistent with separating channels from operations and their send/receive
perspective in AsyncAPI, although SDL's proposed syntax and logical Channel scope
are its own. AsyncAPI also models replies separately; choosing a channel alone
does not define a complete request/reply interaction.
[AsyncAPI operations](https://www.asyncapi.com/docs/concepts/asyncapi-document/adding-operations),
[AsyncAPI reply information](https://www.asyncapi.com/docs/concepts/asyncapi-document/reply-info).

### 14.4 Conditions that topology alone cannot establish

- **Fan-out versus competing receivers:** for the observation example, the intended
  design is independent delivery to BuckingService and UI under their respective
  declared delivery policies. Two receiver declarations alone must not imply a
  shared work queue, nor guarantee identical sampling or lossless delivery.
- **Return messages:** UI command submission needs a separately specified outcome
  direction, correlation and routing policy. A contract may put command and
  result message sets on the same logical Channel; that requires explicit role
  declarations for each. Do not assume results are broadcast to every participant.
- **Routing:** membership does not imply every sender sends every message to every
  receiver. Filters, addressed replies and permitted sender/receiver pairs belong
  to the contract. Overlapping message sets cannot provide conflicting policies.
- **Authority and effects:** sending a command does not grant permission to execute
  it. Domain owners retain validation and accepted-state authority.
- **Concrete path:** logical UI participation does not remove the recorded browser
  adaptation or other required intermediaries. The MVP1 binding must preserve or
  explicitly reassign BWEB's role; no direct browser-to-machine API is adopted here.
- **Real/simulated target:** describe target observations, commands and acknowledgments
  with their own directions. Selection between P1000 and a simulator belongs to
  explicit mode/allocation/binding rules; declaring both must not imply sending to
  both. Neither a topology nor a diagram authorizes hardware transmission.

### 14.5 Required rejection and acceptance cases

| Case | Expected result |
|---|---|
| Declared Containers send/receive an allowed MessageSet through a Channel | Participation types pass; contract completeness is checked separately. |
| Same Container sends requests and receives results on one Channel | Valid with explicit message sets and a compatible interaction contract. |
| A Functionality occupies the Container participant position | Type error in this inter-container construction. |
| Interface or Capability used where Channel is required | Type error; no implicit conversion. |
| Message set is unresolved or absent from the Channel contract | Reject the unresolved/incompatible reference. |
| `uses channel X as receiver of message-set Y` with compatible declarations | Accept explicit type annotations; same semantic arguments as the short form. |
| `uses X as channel`, missing `of`, or role `both` | Reject alternate/underspecified forms. |
| Same participation fact appears twice | Duplicate-fact error. |
| Change only a Channel name's upstream/downstream wording | Direction remains determined by participation, not name spelling. |
| Multiple receivers without a specified distribution/routing policy | Topology can be inspected; a complete communication contract is not established. |
| Add command sending but omit required outcome/correlation semantics | Report an incomplete interaction contract, not a verified command path. |

Before adopting this grammar, complete one observation contract and one
request/result contract using it. Specify MessageSet definitions, Channel-contract
association, role compatibility and distribution rules in this language definition
together. The current parser remains a `design-core 0.5` structural parser.

## 15. Explicit type annotations and checking without annotations

Owner decision, 2026-09-17: all three source forms below are permitted in the
planned participation extension. Repeating a type makes a remote reference easier
to understand without looking up its declaration. Each annotation is optional
independently; relation words and argument order do not change.

```text
container MachineService uses channel MachineServiceUpstream as sender of message-set MachineObservations.
BuckingService uses channel MachineServiceUpstream as receiver of MachineObservations.
UI uses MachineServiceUpstream as receiver of MachineObservations.
```

These three lines describe three participants, not the same fact. To compare
annotated and unannotated representations of one fact, hold every identifier and
role fixed. For example, removing annotations from the first line preserves its
MachineService sender participation.

### 15.1 Annotation is an assertion, not a declaration

In a relation, `container MachineService` asserts a type for a reference. It does
not introduce MachineService, redeclare it, cast it, select an overload or alter
its actual type. An authoritative declaration is still required. By contrast,
`container MachineService.` is a declaration because the sentence ends there.

A future parser must distinguish these productions by the following tokens; it
cannot assume that every sentence starting with a type keyword is a declaration.
The existing requirement to place declarations before statements is unaffected.
Do not replace the 0.2 parser's behavior without an explicit version extension.

There are three distinct inputs to checking a reference:

| Type information | Source | Purpose |
|---|---|---|
| Expected type | Registered argument signature | What kind of object this relation position accepts. |
| Annotated type, optional | The submitted sentence | The author's explicit assertion about this reference. |
| Actual declared type | Resolved declaration | What kind of object the identifier really denotes. |

Use the established subtype rules when checking compatibility. The actual type
must satisfy both the signature and any annotation. An annotation never broadens
what the relation permits. For example, `unit` can describe a declared Container
under `Container <: Unit`, but it does not make a declared plain Unit acceptable
in a Container-only participation position.

### 15.2 Early checks and resolved checks

1. Parse the reference and retain its optional annotation and source span.
2. An early semantic check can reject an annotation that cannot match the argument
   signature. For example, `interface X` in a Channel position is incompatible
   before resolving X. If compatibility depends on a more specific declared type,
   defer that conclusion rather than rejecting a potentially valid subtype.
3. Resolve all names in the declared scope, regardless of annotations.
4. Check the actual type against the signature and any annotation. A `channel X`
   annotation cannot rescue X if its declaration is Interface.
5. Normalize the successful relation into the same typed semantic fact as its
   unannotated equivalent.

The names for the two styles are **explicit type annotations** and **type checking
without explicit annotations**. Expected argument types are derived from the
signature; actual types come from declarations. This is not an inference engine
that invents object types. "Early" describes when an error can be identified, not
a separate kind of truth or a guarantee that parsing alone proves type validity.

No type is borrowed from a preceding sentence. Reordering unrelated structural
statements cannot change reference resolution, role or type. Mixed annotation
density is a reading aid, never implicit grammatical context.

### 15.3 Source syntax, AST and canonical semantic output

- The source AST preserves whether each annotation was written and where, enabling
  faithful diagnostics and an authoring formatter that retains the chosen detail.
- The resolved model records reference identity and actual type independently of
  annotation spelling. Keep source provenance separately from semantic identity.
- Canonical exported participation text uses the shortest unannotated form;
  declarations retain the types. Machine-readable semantic output retains resolved
  types as well. Annotated and unannotated versions of one fact normalize equally.
- Duplicate-fact checking uses resolved arguments and role, ignoring annotation
  presence. Repeating a fact in another annotation style is still a duplicate.
- In the future extension, permitted annotations are valid source, not a
  `NONCANONICAL_FORM` error. Checking authored source and requesting/checking the
  canonical export are distinct operations. Version 0.5's existing canonical
  checks remain unchanged.

The annotation mechanism is intended to be reusable for typed references in later
productions. Section 14 specifies its first concrete use. It does not silently
add annotated references to every existing 0.2 production. The principle is one
meaning and one canonical semantic representation, with optional source-level
type assertions, rather than an unrestricted collection of sentence variants.

### 15.4 Conformance cases to implement with the extension

| Case | Expected result |
|---|---|
| All three references annotated, only Channel annotated, or none annotated | Accept when declarations and contracts agree. |
| All eight annotation-presence combinations for the same participation | Same normalized semantic fact; source AST records the different annotations. |
| Annotation names Interface in the Channel position | Early annotation/signature mismatch, even before resolving the identifier. |
| `channel X` but X is declared Interface | Resolved declaration/signature and annotation mismatch; do not cast X. |
| Fully annotated references with no declarations | Unresolved names; annotations do not create objects. |
| `unit X` where a Container participant is expected; X is declared Container | Valid subtype compatibility after resolution. |
| The same sentence with X declared plain Unit | Invalid participant type despite the compatible Unit annotation. |
| Annotated and unannotated copies of the same participation in one model | Duplicate-fact error. |
| Move a short statement away from an annotated statement above it | Same resolution and semantics. |
| Prefix type keyword followed by a relation rather than a period | Parse as an annotated reference in a statement, not a declaration. |

## 16. Future executable SDL IR

Owner direction, 2026-09-17: explore lowering a fully resolved design into a
versioned executable intermediate representation that can run in an interpreter,
receive data through Channels and interact through ports/adapters, console streams
and an external renderer connected at the Presentation boundary.

This is an execution target, not a new surface-language dialect. Full parsing and
type checking are prerequisites but do not define missing behavior. An executable
profile must specify state, handlers, transitions, data operations, message
delivery, concurrency, time and external effects before lowering is valid.

Record `async` as a concept to define, not a currently legal property. Asynchronous
completion does not require a thread for each Container or Channel. Execution
contexts, scheduling and OS deployment need separate rules. The
[executable IR and runtime study](SDL-Executable-IR-and-Runtime-Study.md) develops
this direction, an initial bounded profile and acceptance scenarios. No interpreter,
runtime, I/O adapter or new parser support is delivered by this documentation.
