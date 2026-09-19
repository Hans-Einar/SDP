# Dataset, Database and Datagram: data origins and contract-defined message families

Current consolidation: [Checkpoint #1](checkpoint%231/README.md), 2026-09-18.
It adds the candidate Value ownership/binding model and reconsiders ControlSet
membership within the full abstraction hierarchy. This proposal remains a
historical input, not an implemented extension of the core grammar.

Date: 2026-09-17  
Status: revised extension proposal E16, aligned with the later E17–E18 study. Not promoted
into the general language, executable corpus or parser.

## 1. Owner clarification and corrected model

The owner intends one internal APT data model that can represent both Classic
and a selected newer StanForD format. Each imported entry identifies its source
format/version and can hold common and format-specific fields. Absent values
remain explicitly unset. Consumers receive defined projections of this data
without needing to understand its source file encoding or storage layout.

This corrects the earlier proposal in this document: **a Dataset is not defined
as a collection of Datagrams**. Dataset entries belong to its internal data
contract. The owner's further clarification makes a declared Datagram a
contract-bound **message family**, linked to a logical Dataset. Its contract
defines the permitted variants and what their payloads mean; each transmitted
Datagram value selects one variant. One family can expose different parts of
the same Dataset without a separate SDL Datagram declaration per payload shape.

A Dataset can supply several families, and one family definition can be used
for multiple runtime instances of its declared Dataset. There is no required
one-to-one correspondence between entries, transmitted values or payload shapes.
The family-to-source relation remains explicit and typed.

A Database is an optional addressable data base: retained data can be queried
and changed through its declared operations, independently of unsolicited arrival.
It can be realized by files or a remote device data-access facade. Units may hold
live/transient Datasets without a Database. SQL is a useful analogy, not a required
engine or public API. Queryability, persistence across restart and current
availability are separate contract properties; none implies the other two.

The later [ControlSet and data-access study](SDL-ControlSets-Layer-Boundaries-and-Data-Access.md)
records the owner's decision that MessageSet is compiler output. It also introduces
ControlSet/Command candidates and evaluates layer boundaries. Those decisions
supersede this proposal's former authored MessageSet membership.

The agreed direction and the proposed grammar are distinct: the owner clarified
the storage/projection model; the exact signatures and constraints below remain
language candidates. No actual newer-standard schema or implemented wire API is
asserted by these illustrative names.

## 2. Vocabulary and grammatical classes

| Term | Type / grammatical role | Proposed meaning |
|---|---|---|
| Database | Noun type | An addressable retained data base exposing declared query/store operations over Datasets. Retention, durability and availability are explicit; SQL and a separate process are not implied. |
| Dataset | Noun type | A logically bounded data collection with a governing internal contract, independent of source encoding and storage engine. |
| Entry | Contract-level concept, not yet a new SDL declaration kind | An identified item within a Dataset, with revision and fields according to that Dataset's contract. |
| Datagram | Noun type | A declared message family governed by a contract defining its variants, header and payload semantics, linked to a logical Dataset. A runtime value selects one variant and identifies its source instance. |
| Functionality | Existing noun type | Owns the responsibility for selecting, validating and projecting data into a Datagram value. SQL, in-memory code or another implementation may realize it. |
| Contract | Existing noun type | Defines applicable shape, meaning, identity, presence, units, compatibility and validity constraints. It may initially be explicitly incomplete. |
| `upholds` | Contract conformance relation | Dataset/Datagram values must satisfy the referenced contract. This states an obligation, not proof. |
| `owns` | Typed responsibility relation | Existing Unit → Functionality; proposed additional Unit → Database. These do not permit reversed or arbitrary ownership. |
| `holds` | Typed holding relation | Unit → Dataset or Database → Dataset. Storage allocation does not grant remote access. |
| `from` | Declaration clause | Datagram → Dataset: identifies the family's logical data origin, independently of any particular runtime entry or storage location. |
| `projects … into …` | Typed ternary relation | Functionality × Dataset × Datagram: refines a producer responsibility for variants of that family; it must agree with the family's declared source. |
| MessageSet | Generated artifact | Per-connection message catalog derived from Datagram families, ControlSets, participant roles, modes and contracts; no independently authored membership. |

A name is a proper-name identifier with one declared type. It may be a subject
or object only in compatible argument positions. Database ownership, holding
and projection are different relations; no synonyms are introduced for them.

State still describes a lifecycle situation, including currentness and pending
work. Representation remains a UI semantic projection with its own context.
Neither is automatically the underlying Dataset or its mutable alias.

## 3. Contract-bearing declarations and candidate grammar

Dataset and Datagram retain mandatory contract references:

```text
contract AptDatasetContract.
contract AptDataContract.
dataset AptData upholds contract AptDatasetContract.
datagram AptMessages upholds contract AptDataContract from dataset AptData.
```

Reject bare `dataset AptData.` and `datagram AptMessages.` in this candidate.
The source clause is mandatory for the linked-Dataset model now proposed; it
links a design definition, not a physical address or a runtime instance.
Each declaration introduces one symbol plus its mandatory conformance relation.
A Contract can be declared early and marked open, allowing incomplete design
work while preventing executable lowering that requires its missing semantics.

Database initially uses a simple declaration. Its schema comes from the held
Datasets; transaction, recovery and durability policy still need explicit
completion for an executable profile. Do not silently infer any of those policies.

```ebnf
databaseDeclaration = "database", identifier, "." ;
datasetDeclaration = "dataset", identifier, "upholds", contractReference,
                     { "and", contractReference }, "." ;
datagramDeclaration = "datagram", identifier, "upholds", contractReference,
                      { "and", contractReference }, "from", datasetReference, "." ;
contractReference = [ "contract" ], identifier ;
datasetReference = [ "dataset" ], identifier ;
holding = reference, "holds", reference, "." ;
ownership = reference, "owns", reference, "." ;
projection = reference, "projects", reference, "into", reference, "." ;
```

`reference` is the earlier optional-type reference construction. The proposed
keyword set gains `database`, `dataset` and `datagram`; all used existing kinds,
including Contract and Functionality, must also be registered.

Proposed complete signatures:

```text
owns     : (UnitCompatible × Functionality) | (UnitCompatible × Database)
holds    : (UnitCompatible × Dataset) | (Database × Dataset)
upholds  : (Dataset | Datagram) × Contract
from     : Datagram × Dataset
projects : Functionality × Dataset × Datagram
```

These are unions of complete signatures, not a Cartesian product. The earlier
candidate `Dataset holds Datagram` is withdrawn. Dataset entry shape belongs in
its data contract; Datagram production belongs to a Functionality.

Keep one immediate Database owner and one immediate Dataset holder per allocation
context. A containing Container's accountability is derived when a nested Unit
owns the Database or holds the Dataset. Do not duplicate immediate assignments
to both parent and child. Reusable Library definitions still require explicit
per-application allocation; they are not global store instances.

`and` retains conjunction: all referenced contract constraints apply
simultaneously. Duplicate contract references are rejected. Supporting alternative
source formats is expressed by the common internal model and import mappings,
not by weakening conjunction into “either.”

## 4. One APT model, optional Database

Illustrative declaration-first fragment:

```text
container BuckingService.
contract AptDatasetContract.
contract AptDataContract.
database AptFileDb.
dataset AptData upholds contract AptDatasetContract.
datagram AptMessages upholds contract AptDataContract from dataset AptData.
functionality ProjectAptPriceMatrix.
functionality ProjectAptFileType.

container BuckingService owns database AptFileDb.
database AptFileDb holds dataset AptData.
container BuckingService owns functionality ProjectAptPriceMatrix.
container BuckingService owns functionality ProjectAptFileType.
functionality ProjectAptPriceMatrix projects dataset AptData into datagram AptMessages.
functionality ProjectAptFileType projects dataset AptData into datagram AptMessages.
```

Here `AptDataContract` defines variants such as `price_matrix` and `source_info`.
Those are contract-local variant keys, not additional SDL Datagram declarations.
The two projection responsibilities can produce those respective variants; their
exact variant mappings belong in contract/behavior refinement, not in the
structural `projects` edge alone. A coarse model can defer those Functionality
details while already declaring the family and source.

An alternative allocation without a Database uses:

```text
container BuckingService holds dataset AptData.
```

Use that holding assignment instead of `AptFileDb holds AptData`, not alongside
it for the same allocation. No Database is required merely to generate a Datagram.

The example assigns responsibility at Container granularity. In the actual MVP1
draft, AptDomain is already a Unit inside BuckingService. Detailed refinement
can place Database ownership and projection Functionalities there without adding
a new APT service or moving domain authority into a UI or transport adapter.

“BuckingService need not know the format” means its common domain algorithms and
consumers can depend on the common data/projection contract. The import adapter
still understands the file format. When it is inside BuckingService, it is that
bounded part of the Container that retains format knowledge; the architecture
does not make decoding disappear.

## 5. Internal entries and missing values

A candidate internal entry contract includes:

- Entry identity, source standard/version, and source artifact identity/digest.
- Dataset/entry revision and import lineage.
- Common semantic fields and explicitly defined format-specific fields.
- Presence for each optional value; stable identities and units for matrix axes.
- Validation and import diagnostics independent of absent optional values.

These are candidate contract requirements, not an assertion about fields in an
existing StanForD standard. The specific newer standard/version remains
`NewStanfordCodecGap`.

One logical Dataset schema can accommodate both forms. This does not require a
single physical table with every possible column: storage can later use tables,
documents, tagged extensions or memory while preserving the same logical contract.

Propose one canonical presence vocabulary:

```text
presence = set | not_set
```

This is a candidate value-domain notation, not executable SDL field-assignment
syntax. When `presence = set`, the value must satisfy its declared type and
constraints. When `presence = not_set`, no value is supplied; zero, empty text
or false must not be fabricated as a substitute.

Distinguish presence from quality/currentness. A present value can be stale or
invalid. Malformed input, parser failures or unsupported conversions must not
silently become an ordinary optional `not_set`. Where consumers need the reason,
the contract can require a structured cause such as source absence, unavailable
derivation or unsupported conversion. The exact reason vocabulary is a separate
bounded contract decision, not a set of interchangeable absence spellings.

For retained historical APT data, preserve enough source information for the
promised audit/round-trip behavior. Do not claim exact reproduction of source
bytes from a normalized entry unless its contract and retained artifact support it.

## 6. Datagram family, variants and projection

`AptMessages` identifies the family and `AptDataContract` defines its variants.
`ProjectAptPriceMatrix` can select an entry/revision and produce a `price_matrix`
value of that family; `ProjectAptFileType` can produce `source_info`. A consumer
needs those contract meanings, not physical tables or source-format parsing.

The SQL analogy is “select the contracted fields for an identified entry,” but
the public operation remains domain-oriented. Do not expose arbitrary SQL or
physical table names as the Channel interface. Selection parameters, units,
filtering, matrix axes, missing-value mapping and failure behavior belong in the
projection's contract/behavior.

A source-independent Datagram requires source-independent **meaning**, not only
identical field names. The import/projection layer must preserve or explicitly
convert quantity, unit, scale, axis meaning and other relevant semantics.
Incompatible information must remain distinguishable; returning an unqualified
number merely because both formats have a “price” field is insufficient.

This matters particularly for MVP1: relative APT machine instructions and
owner-value/settlement economics remain different concepts. A common projection
must not silently convert one into the other. Format-specific knowledge belongs
in the responsible adapter, not scattered through unrelated consumers.

The `projects` fact supplies structural traceability, not a runnable query body.
Executable lowering still needs the selection expression, field mapping, policy
for absent values and implementation/IR binding. Multiple Datasets may eventually
feed a projection, but join/snapshot semantics must be defined before a
multi-input production is added. No Cartesian join is inferred from several edges.

The revised model requires a logical source Dataset for every declared data family.
That does not require prior persistence: a Dataset may be transient, externally
provided, assembled from current observations or derived from other datasets.
For the new family to be constructible, its semantic payload inputs must be
available in that logical Dataset or derivable through its declared mappings.
Generated envelope metadata comes from explicitly defined runtime context.

“All required data exists in the Dataset” is a sufficiency obligation, not a
promise that every variant is available at every instant. A partial Dataset may
support one variant and lack inputs for another. The contract must define
presence rules, production preconditions and unavailable-result behavior.
A compiler must not fill missing payload fields from undeclared side inputs.

The later ControlSet decision resolves the command case: Command invocations
have argument/result contracts and do not require a fabricated Dataset. A compiler
can generate their request/result envelopes separately from Dataset-backed data
Datagrams. Only the latter inherit this proposal's mandatory logical source rule.

### 6.1 Source identity and envelope

The owner's packet/header analogy is modeled as an application-level envelope;
SDL does not adopt a transport's addressing, ordering or framing semantics.

A candidate conceptual header carries:

| Information | Purpose |
|---|---|
| Datagram contract identity and version | Select the permitted family shapes and meanings. |
| Variant key | Select one payload schema within that contract. |
| Source Dataset contract identity and version | Identify the logical source data model. |
| Source Dataset instance ID | Distinguish two independent sources that uphold the same contract. |
| Entry ID and source revision/snapshot, where required | Correlate values and consistent projections. |
| Producer identity/session, where required | Distinguish who produced the value from which Dataset supplied its data. |
| Sequence/correlation information, as required by the interaction | Support declared ordering, replay or request/reply behavior. |

These are semantic requirements. A wire profile may carry them directly or use
negotiated immutable references; SDL need not repeat full contract schemas in
every message. Contract identity/version resolution must be unambiguous.

A header claiming a contract is not evidence that its body satisfies it. The
receiver validates the payload against the selected known variant and enforces
the Channel/participant rules. Likewise, forwarding or projecting data must not
rewrite its logical source to the forwarding endpoint. Dataset revision, producer
session and transport delivery sequence remain separate identities.

Contract declarations identify design definitions; Dataset instance and entry
identities belong to runtime allocation. Neither the Dataset name in SDL nor
its contract ID alone can distinguish two machines or two loaded APT datasets.

### 6.2 P1000 witness without invented protocol facts

Illustrative candidate declarations:

```text
contract P1000DatasetContract.
contract P1000DataContract.
dataset P1000Dataset upholds contract P1000DatasetContract.
datagram P1000Data upholds contract P1000DataContract from dataset P1000Dataset.
```

The family contract defines all permitted variants. Their exact tags, layouts
and meanings must come from the P1000 evidence or an explicitly designed adapter
contract; this sketch invents none. A logical source may represent the device's
data without asserting that the physical device runs a database or SDL runtime.

Preserve MVP1's current boundary: SIM/PGW exchange Target records with Machine,
and Machine owns P1000 decoding and calibration. The SDL envelope can wrap
evidence-backed wire observations without claiming PGW already exposes decoded
sensor semantics. If a semantic Dataset is populated after Machine decoding,
that is a distinct derived Dataset with traceable provenance. The new terms
must not silently move decoding authority into the gateway.

### 6.3 What SDL must supply for executable IR

Keep detail at the appropriate level, while keeping it reachable:

| Layer | Minimum relevant information |
|---|---|
| SDL design | Dataset and family identities; source relation; governing contracts; accountable holders/producers; Channel and participant roles. |
| Referenced contract | Exact versioned variants, header/payload schemas, meaning, units, presence and admissible combinations; no opaque unspecified variant behavior. |
| Execution refinement/profile | Dataset instances/initialization; variant selection and payload construction or declared adapter binding; state effects, validation, delivery and failure behavior. |
| Generated IR | Resolved typed family/variant operations, source identities, runtime checks and traceability to the design/contract/binding. |

A structural design can be useful before all execution refinements exist.
A runnable target must diagnose missing information rather than synthesize
business meaning. Imported complete contracts and bound implementations may
supply the detail; SDL does not need to become SQL or a general-purpose language.

This is the control boundary for generalization: abstract details may be deferred
or supplied through references, but each fact required by the chosen execution
profile must have an explicit, compatible source. A parser accepting the family
name proves neither source sufficiency nor executable completeness.

## 7. Consumer acceptance and consistent provenance

A Datagram is valid if it satisfies its public contract. A consumer separately
decides whether that valid value is sufficient for a specific capability.

For example, a UI may display missing matrix cells, while a planning capability
may require particular complete axes/cells. Both can receive the same valid
Datagram. Their acceptance conditions must be explicit and testable; “the receiver
decides” must not become permission to interpret missing values arbitrarily.

Consequently, a field required by the Datagram contract cannot be `not_set`.
An optional field may be absent in a valid Datagram, yet still prevent a
particular operation from proceeding. Parsing/contract validity and fitness for
an operation are different checks.

An optional `source_info` variant can expose source format/version when needed.
It must correlate with a `price_matrix` value using at least the same entry identity
and relevant dataset/entry revision, or an equivalent snapshot token defined by
the contract. Those keys can reside in a common envelope.

Otherwise, a concurrent import or edit could combine a price matrix from entry
revision A with a source-format response from revision B. Returning two individually
valid Datagrams would then still mislead the consumer. Define coherent snapshot
selection, explicit revision mismatch and unavailable historical revision outcomes.
A local Database transaction alone does not guarantee consistency across
separate network replies.

Consumers need a logical source Dataset instance identity and its governing
contract, plus entry/snapshot identifiers where required. They do not need
physical storage locations, database/table names or the holder's private layout.
The source contract identifies meaning; the instance ID identifies which data.

## 8. Channel boundary and generated MessageSet

A Channel can transfer families whose origins are different logical Datasets.
Authors describe the families, governing contracts, endpoint participation and
allowed variants. MessageSet is derived output; there is no separately authored
`message-set` declaration or `includes` membership in this revised proposal.

Illustrative candidate fragment, assuming all names are declared:

```text
AptTransferContract allows AptMessages.
AptTransferContract allows MachineObservations.
BuckingDataUpstream follows AptTransferContract.
container BuckingService uses channel BuckingDataUpstream as sender of datagram AptMessages.
container BuckingWeb uses channel BuckingDataUpstream as receiver of datagram AptMessages.
```

Here the two families have independently declared Dataset origins.
`allows` now takes a Datagram family in the candidate; this replaces its
authored-MessageSet use. The second family still needs explicit endpoint roles
before an executable connection can deliver it. Allowing a family does not
automatically grant every participant every variant.

Contract/variant restrictions, endpoint roles and execution mode determine the
generated MessageSet. ControlSet contracts additionally contribute their typed
request/result messages. Source identities, reply routing and delivery policies
must survive that compilation; generated membership cannot widen authority.

The public boundary can expose a Dataset contract and logical source identity
while keeping database/table locations private. On-demand queries additionally
need selection, correlation, revision and failure semantics. Output Datagrams
alone do not create an operation for requesting them.

Do not add `Channel sends Datagram` as a competing sender relation: endpoints
send, while Channels constrain permitted transfers. Existing external-schema
references must be reconciled with family/variant contracts before executable use.

The 66-file baseline still contains the earlier MessageSet objects. They have
not been silently deleted or treated as already migrated. Their contracts and
role restrictions must be preserved when a later language revision migrates the
corpus to generated catalogs. The general language definition/parser are unchanged.

## 9. Source organization and declaration order

Proposed organization, leaving current corpus files in place:

```text
SDL/MVP1/
  System.design
  Contracts/
    AptData.design
    AptTransfer.design
  Datasets/
    AptData.design
  Datagrams/
    AptMessages.design
  Channels/
    BuckingDataUpstream.design
  BuckingService/
    AptDomain.design
    AptStorage.design
    AptProjections.design
```

Database declaration/allocation belongs with its responsible Unit; a Databases
folder is optional navigation. File placement never creates a declaration,
ownership fact or implicit remote API.

Preserve System entry ownership of Channel declarations and participant topology.
Detail files enrich those identities rather than redeclaring them. Public
exports can expose Datagram/contract information while keeping storage private.

The owner requests declaration before use. Full examples above obey that rule.
The current exercise instead collects declarations across files before linking.
The proposal requires contracts before data declarations and data declarations
before holding/projection/participation uses in an explicitly ordered expansion.
Do not use filesystem enumeration as source ordering.

Before promotion, reconcile that rule with the existing whole-source declaration
pass and choose its scope. This document does not silently retrofit order
requirements to the current 66-file model.

## 10. Acceptance cases and promotion gates

| Case | Proposed result |
|---|---|
| Dataset/Datagram without a Contract | Reject the declaration. |
| Named, intentionally opaque Contract | Accept incomplete design; block execution dependent on its missing semantics. |
| Unit owns Database; Database holds Dataset | Valid storage responsibility chain. |
| Unit directly holds Dataset | Valid without a Database. |
| Both allocations hold the same Dataset context independently | Reject ambiguous immediate holding. |
| Dataset holds Datagram | Reject the withdrawn storage-membership relation. |
| Functionality projects Dataset into Datagram | Valid only when the source agrees with the family's `from` relation; variant mapping remains necessary for execution. |
| Datagram contract defines several variants | One family declaration suffices; each value identifies its variant. |
| Datagram names a Dataset contract but no instance | Type identity is known; runtime provenance is incomplete. |
| Header selects an unknown contract version or variant | Reject or follow an explicit compatibility policy; never guess a payload shape. |
| Reversed or incorrectly annotated projection arguments | Type error. |
| Classic and newer-format entries inhabit one internal model | Valid when the Dataset contract explicitly represents their semantics and provenance. |
| Optional value absent with `not_set` | Valid if its contract permits absence; consumer preconditions still apply. |
| Invalid source silently converted to optional absence | Reject that mapping. |
| Receiver cannot operate with missing optional fields | Report unmet capability preconditions; do not fabricate values. |
| Matrix and format response use different entry revisions | Reject their combination or request a coherent snapshot. |
| Two stored formats have similarly named but incompatible quantities | Require explicit semantic mapping or report unsupported information. |
| Channel reveals family/source contracts and logical identity, but no storage locations | Valid encapsulated interface. |
| Query output is defined but request selection/error contract is absent | Query interaction remains incomplete. |
| Dataset or Database is declared | Do not infer a SQL engine, reset-surviving persistence, concurrency policy, full CRUD or a physical process. |

Before adopting E16 together with E17–E18: finalize signatures and order; choose exact source standards;
define a small internal APT entry contract; define one APT family with
price-matrix and source-info variants; specify header/presence/snapshot semantics;
test P1000 observations and a transient command case; and trace one valid, missing,
invalid and concurrent-change path through their projections. Add negative type
examples and contract tests when extending the general grammar/parser.

## 11. Relation to existing work

- [Exercise findings and extensions](MVP1-SDL-Exercise-Findings-and-Extensions.md):
  E16 records this data-holding/projection gap and the owner's clarification.
- [Whole-system corpus](../experiments/mvp1_sdl/README.md): baseline remains
  unchanged; this proposal does not claim new model coverage or executable support.
- [Language definition](Design-Language-Definition.md): existing types,
  optional type annotations, ownership and Channel participation.
- [Source tree study](SDL-Source-Tree-and-Compilation-Study.md): explicit source
  membership, declaration identity and compilation phases.

Revision note: this version replaces the earlier Dataset-as-Datagram-collection
proposal. It preserves mandatory data contracts and typed references while
separating storage entries from public projections. The next owner clarification
further groups payload shapes as contract-defined variants of one Dataset-linked
Datagram family. It supersedes the earlier requirement to give each payload
shape its own SDL Datagram declaration. There is one current proposal here,
not two interchangeable dialects.

The subsequent ControlSet/data-access clarification makes MessageSet generated,
separates commands from Dataset-backed data families, and requires explicit
addressability/retention/availability distinctions for Database.
