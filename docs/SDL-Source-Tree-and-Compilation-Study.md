# SDL source tree and compilation study

Date: 2026-09-17

Status: architecture proposal prompted by the owner's source-tree sketch. This
study does not change the accepted `design-core 0.1` syntax, implement a workspace
compiler, create files in MVP1, or settle its remaining deployment decisions.

Owner follow-up, 2026-09-17: use an explicitly declared System as the root, allow
one System per compilation, and treat `System.design` as a convention rather than
a fixed entry filename. The planned rules are recorded in
[the language definition](Design-Language-Definition.md#13-planned-workspace-extension-one-declared-system).

Related: [language definition](Design-Language-Definition.md),
[parser prototype](../experiments/design_core/README.md),
[MVP1 design example](MVP1-Design-Language-Example.md), and
[behavior/source traceability](Scenarios-State-and-Implementation-Traceability.md).

## 1. Target and governing distinction

The proposed target is a maintained design source tree that describes a system
at several levels. A compiler resolves its declarations and relationships into
one checked model, then produces boundary contracts, blueprints and navigable
implementation maps. Agents can inspect a user outcome, open the responsible
units and inspect the relevant implementation and design reasons.

**The source tree organizes authoring; the model graph expresses the system.**
A source file, namespace, design layer, accountable Unit and runtime Container
are different concepts. They often align, but that alignment must be declared
and checked rather than silently inferred from directory nesting.

For example, moving a Presentation declaration to a smaller source file must not
rename the modeled Unit, transfer its ownership, change its runtime allocation,
or break an existing implementation binding. A Feature can involve several
Containers without being copied into every Container directory.

SDL is the proposed umbrella name for this work. `design-core 0.1` remains the
implemented language version; the new name does not introduce another spelling
of its sentences or make the extensions below executable.

## 2. Proposed authoring layout

This is an illustrative future tree under the Ponsse repository, not a set of
existing files or valid new declaration syntax:

```text
MVP1/
  SDL/
    MVP1/
      System.design
      UI/
        Boundary.design
        PonsseDomain.design
        Representation.design
        Composition.design
        Presentation.design
        Renderer.design
        Transport.design
        Delivery.design
        Commands.design
      MachineService/
        Boundary.design
        ...
      BuckingService/
        Boundary.design
        ...
      Contracts/
        ...
      Features/
        InspectLiveState.design
      Scenarios/
        ReplacePresentation.design
      Decisions/
        ...
      Views/
        ...
```

Use lowercase `.design` consistently. Use `Renderer.design` for the selected
Renderer responsibility; `Render` would suggest an action rather than the named
architectural responsibility. Start with one file per coherent concern and split
only when useful. The additional support files represent responsibilities already
identified in the UI study; they are not additional mandatory horizontal layers.

| Location | Proposed responsibility |
|---|---|
| Entry file, conventionally `System.design` | Single System declaration, language version, source membership and declared package scopes. It may declare Container identities once; child boundary files then reference those identities. |
| `UI/Boundary.design` | Logical UI identity/kind, declared layer structure, public offers, external uses and boundary obligations. Mark runtime allocation separately when known. |
| Layer files | Internal Unit/Functionality declarations and their owned facts, layer membership and relevant dependency rules. File membership alone does not establish ownership. |
| Support files | Transport/recovery, delivery policy and command-correlation responsibilities, placed into explicitly declared layers or supporting units. |
| `Contracts/` | System-level contract definitions or explicit bindings to authoritative external schemas. Preserve the existing contract-governance owner; do not copy schemas into a new competing registry. |
| `Features/`, `Scenarios/` | Cross-system outcomes and behavioral interactions, referencing participant identities. Local subflows may live with their owning package. Each definition exists once. |
| `Decisions/` | Decision/reason records, affected identities and links to existing SDP authority. Existing records can be referenced instead of duplicated. |
| `Views/` | Authored selection/layout instructions for generated views. They may select model facts but cannot create architecture facts. Exact format remains open. |

Not every project needs every directory, and folder names are not new language
keywords. The root entry should enumerate included files or bounded source groups
explicitly. Discovery order cannot affect meaning. Generated outputs, editor
backups and unlisted drafts must not silently join a build. Report unlisted
`.design` files so forgotten source is visible.

Select the entry file explicitly; neither `System.design` nor the enclosing
directory name supplies System identity. A future directory-based command needs
an explicitly configured entry, not a guess among source files. The System
declaration belongs in the entry, and exactly one must exist across the complete
source set. A layer file need not contain a System declaration.

`MVP1 contains UI.` expresses model membership; it does not locate UI source
files. Source membership and containment remain separate. The
[planned System extension](Design-Language-Definition.md#13-planned-workspace-extension-one-declared-system)
defines the proposed typed form and rejection cases without changing version 0.1.

The owner's next topology proposal places Channel declarations and Container
participation in the entry file as well. Detailed contracts can be referenced
from their authoritative included sources. The
[Channel proposal](Design-Language-Definition.md#14-proposed-channels-participants-and-message-direction)
separates channel identity, message-scoped sender/receiver roles, distribution
semantics and concrete bindings. Direction is not inferred from an upstream or
downstream suffix, and topology is not a complete communication contract.

### 2.1 UI boundary versus runtime Container

The four selected concerns remain Representation -> Composition -> Presentation
-> Renderer. The earlier MVP1 study leaves exact host/browser/native allocation
open. A Go headless host and a browser renderer can occupy separate runtime
applications while sharing one logical UI design package.

Therefore, `UI/` can already be a useful source boundary. Its descriptor must say
whether it represents a logical subsystem or a runtime Container. If the final
deployment splits UIHost and BrowserRenderer, record both Containers and their
contract without renaming all logical UI identities. C4 likewise defines a
Container as an application or data store, not an arbitrary source folder.
[C4 Container definition](https://c4model.com/abstractions/container).

`PonsseDomain.design` should describe Ponsse-specific UI projections and adapters.
Machine, APT and bucking truth stays with the recorded domain owners. Generic
Representation, Composition and Presentation machinery must not acquire Ponsse
domain dependencies merely because all five files share a directory.

An architectural layer is a grouping with dependency constraints, not necessarily
an accountable Unit or a running process. Do not use `owns` to mean layer
membership. Also, the selected arrows do not imply one synchronous data pipeline:
measurements, presentation replacement and returned user intents have distinct
interactions and policies.

## 3. Names, scopes and cross-file resolution

Recommend the following rules before choosing concrete import syntax:

1. Each design object has one authoritative declaration in one declared package
   scope. Files within that scope can be split without changing its identity.
2. Identity is independent of the physical file path. Initially use declared
   scope plus declared local name. A namespace move or rename is an explicit
   identity migration; a file move within the same scope is not.
3. Conceptually, a name such as `MVP1.UI.PresentationManager` is qualified by
   declared scope. This dotted spelling is explanatory here: `design-core 0.1`
   does not accept dotted identifiers. The next grammar must choose one form.
4. Unqualified references resolve only within the declared local scope.
   Cross-scope references use the single chosen qualified form and an explicit
   dependency declaration. No wildcard imports, implicit parent-scope search or
   user-defined aliases in the first workspace version.
5. Duplicate declarations are errors across the complete scope. Unresolved names
   do not become implicit external interfaces. External contracts need declared,
   versioned identities even when their implementations are outside the build.
6. Public visibility is explicit. A consumer of another Container may use its
   exported boundary, not reach into its private Units because their files happen
   to be available in the same checkout. An enclosing scenario can refer to
   internal behavior for analysis without granting implementation access to it.
7. A fact has one authoritative source location. Reference it from other files
   instead of repeating it. Cross-container interactions belong to their declared
   contract/scenario owner; consuming packages own their local dependency facts.

Source membership, name visibility, design dependency and runtime communication
are separate relationships. Importing a name does not mean a runtime message is
sent. Likewise, two Containers exchanging messages need not form an invalid
compilation cycle. Resolve declarations before checking relationship-specific
constraints; apply acyclicity to containment/refinement and explicitly constrained
layer dependencies, not to every edge in the system graph.

The workspace extension must be an explicit new language version with migrated
fixtures. Existing `0.1` models remain standalone models with their existing
namespace rules. Concatenating them or resolving their unknown names from nearby
files would silently change their meaning.

## 4. Parser, linker and compiler responsibilities

The next architecture should retain a distinct syntax tree for each file:

```text
Root entry + declared source set
  -> parse each file -> file ASTs with source locations
  -> collect declarations -> package/global symbol index
  -> resolve references and visibility
  -> validate types, boundaries, ownership and declared constraints
  -> resolved system model (typed graph)
  -> public boundary exports / blueprints / impact and coverage reports
```

The owner also proposes an executable output path. After behavioral completeness
checks for an execution profile, the resolved model could lower into a versioned
IR for an interpreter with typed external adapters and a Presentation boundary.
This requires additional semantics beyond successful parsing/linking; see the
[executable IR and runtime study](SDL-Executable-IR-and-Runtime-Study.md).

An AST preserves the structure of written statements. The resolved model connects
references to their actual declarations and adds derived indexes while retaining
source provenance. There may be a workspace AST root containing the file ASTs,
but this tree does not replace the system graph. LLVM's introductory compiler
tutorial similarly distinguishes source constructs represented by an AST from
later compilation stages; SDL's graph and outputs are our proposed application
of that general separation.
[LLVM parser/AST tutorial](https://llvm.org/docs/tutorial/MyFirstLanguageFrontend/LangImpl02.html).

| Stage | Must establish | Must not claim |
|---|---|---|
| Load/parse | Known version, explicit inputs, valid sentence structure, file/line/column provenance. | A name resolves or an architecture is valid. |
| Link | Unique identities, resolvable references, valid visibility and dependency versions. | A referenced capability is correctly implemented. |
| Validate | Typed relations, consistent ownership, allowed dependencies and applicable completeness checks. | Unmodeled runtime behavior is verified. |
| Export | Public projection is closed over its referenced contract types and belongs to this input revision. | A generated contract is a deployed API or a granted permission. |
| Generate views | Every displayed fact has model identity/provenance; aggregation is explicit. | Hidden detail does not exist, or a diagram proves behavior. |

A preprocessor is not needed merely to assemble files. Avoid textual inclusion,
macros and conditional rewriting at this stage: they introduce another place
where meaning and diagnostics can diverge. A loader plus linker can perform the
needed assembly without altering source text.

For reproducibility, record compiler/language version, dependency versions and
content hashes of every declared input. Deterministic output ordering must not
reorder behavioral steps. Distinguish a complete workspace check from a partial
package check with external dependencies. Unresolved requirements or absent source
bindings can be reported as explicit coverage gaps where the selected validation
target permits them; they cannot be presented as implementation completeness.

## 5. The generated container boundary

The proposed `UI.container` is useful as a **compiled public boundary artifact**.
Generate it after linking and semantic validation, rather than before those checks.
An authored `Boundary.design` states the public facts; the export gathers and
checks them without requiring a second handwritten interface description.

Keep outputs outside the scanned source tree, for example:

```text
MVP1/build/SDL/MVP1/
  model.json
  exports/
    UI.container
  blueprints/
    UI.mmd
    UI.map.json
  coverage.json
```

These filenames are proposed artifact roles, not implemented schemas. Export
`.container` only for an actual Container; a logical UI subsystem needs a
differently typed model view until allocation is decided. The final encoding can
be a versioned JSON schema without inventing a second handwritten language.

The public artifact should contain:

- Container identity, model/contract versions and an input digest.
- Explicitly exposed capabilities and the interfaces through which they are used.
- Public operations, messages, value types, units and identity rules, including
  all transitively referenced public contract definitions or versioned references.
- Consumed interfaces and mode/condition-scoped required dependencies.
- Relevant ordering, freshness, delivery, failure, correlation and authority
  obligations. These are necessary to collaborate correctly, not private code.
- Provenance for exported claims and any unresolved public contract obligations.

It should exclude private Functionality, layout helpers, implementation symbols
and incidental internal dependencies. If an exposed operation needs a private
type, fail the export or require an explicit public contract type; do not leak
all internals to make the export complete. Preserve centrally governed contracts
as references to their authoritative definitions.

When compiling the full system from source, exports are outputs, not an alternate
input competing with that same source. A separately compiled consumer can load
a pinned export instead. Reject mixed source/export definitions of the same
identity and stale or incompatible artifacts. A contract version policy needs
semantic compatibility rules; a digest alone establishes identity, not compatibility.

Build outputs can normally be regenerated and ignored by Git; reviewed release
artifacts may be retained deliberately. Neither policy changes their derived
authority. Generation must not overwrite authored files or publish partial outputs
as a successful full build.

## 6. Blueprints and decorated Mermaid

Define a blueprint as a **selected view of the resolved model, with traceability**.
Keep its machine-readable graph separate from its visual rendering. The same
model should support system/container maps, layer dependencies, capability
realization, scenario flows and implementation impact views.

Recommended first output bundle:

| Artifact | Responsibility |
|---|---|
| `model.json` | Versioned resolved model and provenance, independent of a diagram renderer. |
| `UI.mmd` | Standard Mermaid syntax with deterministic diagram identifiers and readable labels. |
| `UI.map.json` | Diagram node/edge identifiers mapped to model facts, source locations, view filters and model digest. Aggregated nodes map to multiple identities explicitly. |

This accommodates the owner's decorated-diagram idea without depending on a
custom Mermaid dialect. Mermaid supports comments ignored by its parser and node
links/callbacks, with interaction dependent on the host's security configuration.
A short generated comment can identify the map artifact. Ordinary Mermaid shows
the drawing; an SDL-aware viewer reads the map to provide deeper navigation.
[Mermaid flowchart documentation](https://mermaid.js.org/syntax/flowchart.html).

The separate map remains authoritative for additional metadata; do not rely on
comments surviving other tools or arbitrary metadata being interpreted by every
Mermaid renderer. Pin and test the renderer version used for generated diagrams.
Generating valid Mermaid and proving a working interactive viewer are separate
acceptance checks. No Mermaid renderer was executed for this study.

Initially, edits go into SDL and regenerate views. Later graphical editing must
translate an explicit model change back to authoritative source, validate it and
show a source diff. A hand-edited `.mmd` file must not silently become a second
design authority. Not every diagram view contains enough information to reconstruct
the full model, so arbitrary diagram round trips cannot be promised.

## 7. Navigation down to implementation and design reasons

The intended navigation chain is:

```text
Feature -> Scenario step -> Activity -> Capability/contract
        -> responsible Functionality -> owning Unit/layer/Container
        -> implementation binding -> revision-specific source symbol
```

Each arrow needs a defined relationship, not just matching labels. A source binding
identifies repository, revision, symbol and realization role; one Functionality
may need several symbols and one shared symbol may contribute to several facts.
The existing inspected MVP1 symbols are evidence candidates, not fabricated
implementations of the target Presentation contract.

Preserve four different kinds of information:

| Information | Example | Consequence for tools |
|---|---|---|
| Design obligation | Replacing Presentation preserves Representation identities. | A change touching the related steps/units must surface this constraint. |
| Decision and reason | Domain truth stays with its owner; UI exposes projections and intents. | Show the reason and the decision it supersedes, not only the latest box diagram. |
| Intended realization | A given Functionality is assigned to PresentationManager. | Support impact analysis without claiming it is implemented. |
| Observed evidence | A test or inspected symbol at a revision supports part of the obligation. | Show coverage, freshness and missing realization separately. |

Facts and decisions need stable identities and explicit status such as selected,
proposed or superseded. The exact provenance schema remains to be defined; link
existing SDP records instead of inventing another decision ledger. Changed source
invalidates or qualifies dependent evidence rather than leaving a green diagram.
Impact reports are conservative assistance, not a complete inferred call graph:
dynamic dispatch, messages and missing bindings remain visible uncertainties.

For example, inserting revalidation before Presentation publication should expose
the Composition contract revision, candidate lifecycle, commit guard, renderer
completion distinction and their bindings. It must not select every machine
function merely because measurements eventually reach the same screen.

## 8. What is missing before modeling MVP1 in full

| Capability | Present in 0.1 | Required definition |
|---|---|---|
| Parse structural sentences | Yes; one standalone model. | Preserve existing conformance while adding file-aware diagnostics. |
| System/package/layer structure | Unit and Container types only. | System identity, package scopes, Layer membership and permitted dependencies. |
| Multi-file compilation | No. | Root source membership, declaration collection, qualified references, visibility and version rules. |
| Reusable Channels and public contracts | Named Interface, `consumes`, scoped `requires`. | Contract content/ownership, exposure, participant roles, protocol and compatibility. A Channel is not an import or one connection per Feature. |
| Container exports | No. | Explicit visibility, public type closure, artifact schema and dependency identity. |
| Features and scenario flow | Activity and `refines` only. | Feature/Scenario/step relationships, sequence, branch, concurrency, guards and outcomes. |
| State and time | Named Mode only. | State variables, revisions, predicates, event/observation identities and check/commit semantics. |
| Implementation navigation | No. | Revision-specific SourceSymbol, role-specific bindings and evidence coverage. |
| Decision preservation | Surrounding study prose. | Fact IDs, authority/status, rationale references and supersession. |
| Blueprints | No generator. | View selection, graph/map schemas, source maps and renderer acceptance tests. |

The tree helps identify these gaps; it cannot supply their semantics implicitly.
Do not solve missing contracts with free-text properties or invented meanings
for existing verbs. Add each construct to the single language definition with
positive and negative examples before implementing it.

## 9. Recommended sequence and acceptance cases

**First: agree the workspace model at the level of this study.** Resolve declared
scopes, file-independent identity, ownership of shared facts, public visibility,
layer semantics and source/output separation. Keep detailed artifact schemas and
deployment choices open where they are not required for the first experiment.

**Second: implement one small multi-file language increment.** Split a structural
UI sample across a boundary and two layer files, with a declared external consumer.
Use one authoritative declaration per object. Exercise the loader, linker and
cross-file semantic checks before attempting a complete MVP1 transcription.

Required cases:

| Case | Expected result |
|---|---|
| Move a declaration between files in the same declared scope | Same model identity and facts; only source provenance changes. |
| Change filesystem discovery order | Same resolved model and canonical artifacts. |
| Duplicate a name across files | Error pointing to both declarations. |
| Refer to a private Unit from another Container's implementation | Visibility error; the boundary contract remains usable. |
| Add a second owner or a containment cycle across files | Same rejection as within one file. |
| Load an old public export with a newer provider source for the same identity | Explicit conflict/staleness result, never silent preference. |
| Generate a boundary that references an unexported private type | Reject an incomplete public contract. |
| Render a blueprint or inspect its metadata | Every node/edge resolves to current model facts or explicit aggregation. |
| Omit a source binding or leave a design decision unresolved | Visible coverage/status gap; no false implementation-complete result. |
| Keep an unlisted `.design` file beside included sources | Report it; do not silently change the compilation unit. |

**Third: generate one structural blueprint and public boundary projection.**
Use the smallest model that tests the relevant semantics. Test clean rebuilds,
deterministic output, artifact provenance and the diagram renderer.

**Fourth: complete the Presentation replacement scenario.** Add the behavioral
constructs demanded by invalid, stale and concurrent paths. Then test whether a
flow change exposes affected responsibilities, decisions and real/missing source
bindings. This ties back to L03/L04/L08 in the
[conformance study](Design-Language-Conformance-Scenarios.md).

This sequence establishes the destination before broad language growth, while
keeping each increment small enough to falsify. It does not require building a
general code generator or describing all of MVP1 before finding whether the
source model, contracts and blueprints preserve the intended meaning.
