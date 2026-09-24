# SDL language candidates and open conformance work

Status: **candidates/open work**, not additional design-core 0.5 rules.
Moved from the definition's sections 12–16 in R3-M2, 2026-09-24. Original section
numbers are retained for provenance. Owner decisions explicitly recorded below
remain decisions about future scope; they do not imply parser support.
The [active structural profile](../profiles/SDL-Structural-Core-Profile.md) and
its registered additions define current structural support. In particular, the
[Channel profile](../profiles/SDL-Channel-Scenario-Profile.md) is implemented;
the broader topology/workspace/behavior proposals below do not supersede it.
Executable actions use [action-core](../profiles/SDL-Executable-Action-Profile.md),
not automatic execution of all candidate IR concepts.

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
language evaluation remain open. See the [prototype guide](../../../experiments/design_core/README.md)
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
