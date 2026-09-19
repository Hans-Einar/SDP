# ControlSets, layer boundaries and addressable data

Current discussion: [Checkpoint #1](checkpoint%231/README.md), 2026-09-18.
The combined Commands-and-Values recommendation below is now an alternative
under reconsideration. The checkpoint explores separate owned Values and a
command-only ControlSet without claiming that the owner has adopted that change.

Date: 2026-09-17  
Status: design study and candidate extensions E17–E18. Owner decisions and
recommendations are distinguished below. No changes to the general grammar,
parser, executable corpus, product code or hardware behavior.

## 1. Decisions and questions from the owner

Selected direction:

- MessageSet should be compiler output rather than a separately maintained
  authoring object. One Channel can carry Datagram families from several Datasets.
- Data and control contracts should describe communication between Units across
  horizontal layers as well as between Containers.
- A ControlSet describes Commands with arguments and results.
- Database is an implementation-independent data access abstraction: files and
  remotely addressable P1000 data are possible realizations; SQL is not required.
- Live data can arrive as Datagrams from a Dataset without a Database.

Still exploratory:

- Whether ControlSet also groups exposed Values.
- Whether an internal connection should be called Channel, Pipe or Wire.
- Exact Database persistence/availability guarantees and which CRUD operations
  each Database offers.
- Exact syntax and runtime realization of the new contracts.

Recommendations below are proposals, not claims that these choices are already
adopted or implemented. The earlier Dataset/Datagram proposal is updated to remove
authored MessageSet membership; its baseline corpus remains historical evidence.

## 2. Evidence from other models

### AADL: distinguish the boundary from its execution binding

SEI describes AADL models whose data/event-data ports are mapped into partition
communication resources, while functional code is supplied separately. This is
evidence for preserving communication semantics in an architecture model and
choosing an execution realization later; it is not proof of SDL semantics.
[SEI: AADL code generation](https://www.sei.cmu.edu/blog/aadl-code-generation-for-avionics-systems/).

SDL recommendation: use explicit typed endpoints/connections and layer-crossing
rules. Do not infer a queue, thread or network merely because data crosses an
architectural boundary. A direct local call and an asynchronous delivery can share
a payload contract but still have different scheduling, ownership and failure
contracts.

### AsyncAPI: messages, participants and replies are distinct concerns

AsyncAPI describes messages associated with Channels and operations that send or
receive selected messages. Request/reply is an explicit interaction, including
reply routing/correlation. SDL can borrow this separation without adopting its
document format or maintaining a second authored MessageSet.
[Channels](https://www.asyncapi.com/docs/concepts/asyncapi-document/adding-channels),
[operations](https://www.asyncapi.com/docs/concepts/asyncapi-document/adding-operations),
[request/reply](https://www.asyncapi.com/docs/tutorials/getting-started/request-reply).

SDL recommendation: authors declare families, operations, roles and contracts;
the compiler derives the effective message catalog for each connection/context.

### OPC UA: values and operations can coexist without becoming interchangeable

OPC UA's modeling guidance distinguishes Methods with inputs/outputs from
Variables. It specifically discourages emulating a Method by writing inputs into
Variables and later fetching outputs from other Variables. Its Argument type
describes input/output parameters.
[Methods](https://reference.opcfoundation.org/specs/OPC-10000-3/a-6),
[Arguments](https://reference.opcfoundation.org/specs/OPC-10000-3/8.6).

SDL recommendation: a ControlSet can group Commands and Value bindings, but each
member retains its own semantics. An assignment to a Value must not secretly mean
an unrelated procedure call. Conversely, a Command result must remain correlated
with that invocation, not just overwrite a shared “last result” value.

These are distinct kinds of prior art: AADL is an architecture language, AsyncAPI
an interface-description specification, and OPC UA an information/service model.
None defines the SDL terms Dataset, Datagram or ControlSet in the proposed sense.

## 3. What HEOS and HSX already contribute

Read-only inspection pinned to:

- HEOS: `c41995ecc503684f4bb5200ae8f8da3665d57e1c`.
- HSX: `e374da88f4dd470bad2d8ec6a2a14f1ce367e40e`.

### HEOS observations

[ICommand](https://github.com/Hans-Einar/HEOS/blob/c41995ecc503684f4bb5200ae8f8da3665d57e1c/src/ICommand.h)
and [command_system.h](https://github.com/Hans-Einar/HEOS/blob/c41995ecc503684f4bb5200ae8f8da3665d57e1c/src/command_system.h)
show named command handlers, an optional argument parser, execution with Arguments,
and separate output/status handling.

[value_system.h](https://github.com/Hans-Einar/HEOS/blob/c41995ecc503684f4bb5200ae8f8da3665d57e1c/src/value_system.h)
contains value identity/group metadata, read-only/read-write modes, scalar/collection
forms and load/save hooks. Inspected excerpts also contain EEPROM-specific
implementations. These are reusable conceptual precedents, not generic SDL storage
semantics or permission to run hardware operations.

### HSX observations and an evidence limit

The [Value & Command architecture](https://github.com/Hans-Einar/HSX/blob/e374da88f4dd470bad2d8ec6a2a14f1ce367e40e/main/03--Architecture/03.04--ValCmd.md)
and [detailed design](https://github.com/Hans-Einar/HSX/blob/e374da88f4dd470bad2d8ec6a2a14f1ce367e40e/main/04--Design/04.04--ValCmd.md)
describe grouped identity, separate value/command registries and external
transport bindings. The design separates notification/persistence mechanisms
from the registry and describes owner-specific access.

The inspected [value header](https://github.com/Hans-Einar/HSX/blob/e374da88f4dd470bad2d8ec6a2a14f1ce367e40e/include/hsx_value.h)
declares get/set/subscription and explicit persistence operations/modes.
The [command header](https://github.com/Hans-Einar/HSX/blob/e374da88f4dd470bad2d8ec6a2a14f1ce367e40e/include/hsx_command.h)
separates command registration/call/async call and status/result conventions.

Do not present all HSX material as one fully verified contract:
[the interface document](https://github.com/Hans-Einar/HSX/blob/e374da88f4dd470bad2d8ec6a2a14f1ce367e40e/docs/hsx_value_interface.md)
emphasizes zero-argument commands and async status replies; it is not evidence
that arbitrary typed command arguments/results already work. It also differs from
the inspected header in value function-ID details. No HSX tests or runtime were
executed here. SDL should adopt the conceptual separation, not copy ABI numbers,
f16 representation or unresolved documentation differences.

The useful common idea is a discoverable surface containing typed Values and
Commands, with identity, access rules and bindings to transport/persistence.

## 4. Proposed responsibilities and vocabulary

| Concept | Design responsibility |
|---|---|
| Dataset | Logical data model/source and its instance identity; may be transient or retained. |
| Datagram | Contract-defined family of transferable data variants tied to a logical Dataset. |
| Database | Addressable retained data available through declared query/read and, where supported, store/change operations. |
| ControlSet | Contract-bound interaction surface grouping Commands and, optionally, bindings to exposed Values. It does not own a second copy of domain state. |
| Command | Invocable operation with argument/result contracts, authority, preconditions, effects and lifecycle. |
| Value member | Contract member bound to Dataset-backed state/projection, with explicit read/observe/write capabilities. A separate global SDL Value noun is not required yet. |
| Channel | Logical connection between explicitly identified endpoints, usable across Container or layer boundaries. |
| MessageSet | Generated catalog of the messages permitted by a connection's resolved contracts, endpoints, roles and mode. |

A Command is not an owning Unit or the Functionality that realizes its behavior.
Several commands may share implementation responsibilities; a Functionality may
serve both a command and an internal scenario. Keep public operation contracts
separate from source-call bindings.

For the candidate surface, use the keyword spelling `control-set`, consistent
with existing multiword type keywords. This is a proposed spelling, not a parser
change. Use `command` for its invocable members. Arguments/results can initially
be defined by the referenced contract; do not force a global noun for every field.

## 5. ControlSet with Commands and Values

Recommend allowing both member categories, with explicit behavior:

| Member | Example meaning | Essential distinctions |
|---|---|---|
| Readable Value | Current selected APT identity | Currentness/quality, identity and snapshot. |
| Observable Value | Measured length updates | Subscription, initial baseline, loss/coalescing and ordering. |
| Writable Value | Requested setting | Validation, access, revision check and committed value; a request is not immediate shared-memory assignment. |
| Command | Activate a selected APT definition | Arguments, domain preconditions, side effects and invocation-specific result. |
| Query command | Retrieve an identified entry | Selection/revision contract; read-only effects where promised. |

A ControlSet Value refers to existing Dataset state or a declared projection;
it must not become a duplicate authoritative value with independent updates.
If the same state is also published as a Datagram, both paths must agree on
identity, revision, units and currentness.

A Command contract should specify:

- Argument and result schemas, including the legitimate empty/void cases.
- Domain failure versus rejected invocation versus transport failure.
- Immediate completion or an explicit accepted/pending/completed lifecycle.
- Invocation identity, retry/deduplication and cancellation/timeout semantics.
- Responsible handler, preconditions, effects and postconditions.

A timeout does not prove that the operation did not occur. A transport
acknowledgment does not prove domain completion. A bare notification should not
be disguised as a command with a fabricated result.

No extra “command Dataset” is required just to make commands legal. Commands are
a separate semantic interaction; request/result envelopes can be generated from
ControlSet contracts. The Dataset-origin rule applies to data Datagram families,
not to every control envelope the compiler emits.

## 6. MessageSet as compiler output

Authoring identifies Channel topology, supported Datagram families and exposed
ControlSets, with explicit participant roles and contract restrictions. A Channel
can therefore carry live observations from one Dataset and query responses from
another, alongside command traffic, without a separately authored grouping.

For an illustrative connection, a generated catalog might contain:

| Generated entry | Derived from | Origin distinction |
|---|---|---|
| Live observation variants | A live-data Datagram contract | Transient observation Dataset instance. |
| Parameter result variants | An addressable-data Datagram contract | Queryable parameter Dataset instance/revision. |
| Query request and failure response | ControlSet command contract | Invocation context, not a fabricated data store. |
| Update request and completion | ControlSet command contract | Requested mutation and correlated result. |

This table is explanatory, not SDL syntax or a P1000 wire specification.

The compiler must derive **per-role/per-mode permitted traffic**, not grant every
participant the union of all variants. Include contract versions, source families,
direction, reply association and relevant delivery policy. Diagnose incompatible
schemas, duplicate/colliding identities, missing reply routes and ambiguous
variant restrictions. Generated catalog identifiers must be stable and traceable
to their source declarations.

The catalog can become part of a Container boundary export and SDL IR. Rebuilding
must replace it deterministically; users should not hand-edit it as another source
of truth. Existing corpus MessageSets must be deliberately migrated later, with
their old sender/receiver restrictions preserved.

## 7. Between layers: recommend one Channel concept

Recommend keeping **Channel** as the single logical connection concept, with
explicit endpoint and boundary scope. This is a recommendation, since the owner
has not chosen between Channel, Pipe and Wire.

Use the existing Interface notion for a Unit's contracted interaction surface.
A connection relates those surfaces. Existing physical Ports/adapters need not
be repurposed into another synonym for Interface.

This avoids three nouns distinguished only by where a line is drawn. Pipe/Wire
could become distinct types later if they acquire real semantics, such as a
specific stream protocol or execution binding. They should not be alternative
spellings for the same connection.

Layer scope requires more than placing files into folders:

1. Identify the participating Units and the layer boundary crossed.
2. Declare data and control contracts visible at that boundary.
3. Require crossing interactions to use those contracts; do not allow an
   unmodeled read/write of another layer's mutable state.
4. Define permitted dependency directions and allowed exceptions/cycles.
5. Bind the connection to execution separately.

The current exercise's Layer is a grouping view and can span library ownership.
It is not yet an enforceable encapsulation boundary. Explicit authoritative
boundary membership/exposure rules are therefore necessary before checking
“all inter-layer values use Dataset/Datagram.” Folder location is insufficient.

This rule need not create a Dataset or Channel for each local variable/function
argument. It applies to declared architectural boundaries. Several related values
can share one Dataset and one family contract.

A local binding may lower to a typed call, immutable reference or in-memory
message. A cross-process binding may require serialization and delivery machinery.
Neither crossing a layer nor writing `async` automatically creates an OS thread.
Compiler optimization may remove copies while preserving contract semantics,
lifetime/aliasing constraints and traceable boundaries.

## 8. Database, live data and persistence

The owner's “always available” describes on-demand access instead of dependence
on the next unsolicited arrival. Treat it as a declared access capability with
failure behavior, not an assertion that a resource can never disconnect or fail.

Three independent questions must be answerable:

| Dimension | Examples |
|---|---|
| Access | Query current data, retrieve a named entry, store a change, observe new arrivals. |
| Retention/durability | Retained during one process/device session, preserved across reconnect, persisted across restart/power loss. |
| Current availability | Reachable, disconnected, missing entry, permission denied, invalid/unavailable revision. |

For now, define Database as an addressable retained data base and require an
explicit retention/durability contract. This includes the owner's persistent
filesystem case without pretending all readable P1000 state survives reset.
If the owner ultimately reserves Database strictly for restart-persistent data,
queryable volatile device state will need a different access classification.
That narrower definition is still open; queryability alone cannot prove persistence.

Examples:

| Case | Model |
|---|---|
| Directory of retained APT files | Database implemented by filesystem; files can instantiate Datasets, and contracted projections yield Datagrams. |
| P1000 retained parameters | Remote addressable data facade; supported reads/writes and actual retention require protocol evidence. |
| Unsolicited P1000 observations | Transient Dataset identity plus incoming Datagram family; no query capability inferred. |
| Receiver caches the latest observation | Explicit receiver-held Dataset with freshness policy; cached value is not a new measurement. |
| Receiver records the stream for replay | A new retained/queryable data base with source lineage and declared gaps/retention. |

One device may expose both queryable state and unsolicited observations. They
need different interaction contracts, but do not automatically require different
Dataset identities. Use separate Datasets when ownership, schema, retention or
consistency domains differ; preserve a shared identity when both access modes
refer to the same logical data. This is a modeling decision, not a naming rule.

A transient Dataset may describe the live data source even if the receiver never
materializes the complete state. Losing a message does not imply the source can
re-send or answer a query. Buffering, replay and recovery are separate promises.

## 9. CRUD and remote resources

CRUD is a useful checklist, not mandatory functionality for every Database:

- Query/read must specify identity, selection, projection and revision semantics.
- Create/store must define validation, assigned identity and commit behavior.
- Update must define concurrency/preconditions and accepted result.
- Delete must be explicitly supported; a file deletion is different from resetting
  a fixed device parameter.

A read-only resource is legitimate. A device with fixed parameter slots may
support read/update but neither create nor delete. Expose only the operations
its actual contract supports.

Database owns no extra domain authority. A Unit owning a logical proxy does not
own the remote device's data or bypass its access rules. Separate logical
resource ownership, physical authority and local mirror state in the design.

For P1000, query/update messages, meanings and persistence need repository/capture
evidence. A DSL declaration cannot establish them. This study creates no serial
transmission path and changes no existing simulator-only boundary.

## 10. What must survive into SDL IR

| Design fact | Required executable refinement |
|---|---|
| Data family and source | Variant/schema binding, source instance and revision, payload construction/validation. |
| Command and ControlSet | Typed invocation/result, handler binding, lifecycle, failure and retry rules. |
| Value exposure | Authoritative Dataset binding, read/write/observe behavior, revision and freshness. |
| Channel boundary | Endpoints, permitted families/operations, scheduling/lifetime/delivery profile. |
| Database access | Supported operations, retention, selection/commit semantics and adapter. |
| Layer restrictions | Checkable boundary membership and traceability through any optimized local binding. |
| Generated MessageSet | Derived catalog with stable identity, source mappings and role/mode restrictions. |

Dataset, Datagram and Command definitions can be abstract early. A target execution
profile must report missing contracts/bindings rather than manufacture a storage
engine, query endpoint, retry policy or default argument. The generated IR should
preserve these distinctions even if several operations share the same wire format.

## 11. Focused next witnesses

Before expanding the entire grammar, specify three small scenarios:

1. **Live measurement across UI layers:** a Dataset-backed observation becomes a
   Representation update. Check stale/missing data, direction and no hidden mutable
   cross-layer access.
2. **Queryable retained APT:** filesystem Dataset, query projection and a
   revision-checked update Command; interrupted storage has an explicit outcome.
3. **Device with both access modes:** unsolicited observation plus an offline
   simulated query/update surface. Distinguish correlation, retry, cache freshness
   and unknown retention. Do not invent P1000 opcodes.

Acceptance failures to exercise: unauthorized variant in generated catalog;
ControlSet Value not bound to an owner; command result treated as latest-value
telemetry; timeout treated as cancellation; querying a live-only source; implicit
Create/Delete on fixed device state; layer access bypassing its declared interface;
and unresolved persistence incorrectly reported as durable.

## 12. Integration boundary

This study extends [the Dataset/Datagram proposal](SDL-Datasets-Datagrams-and-Data-Contracts.md).
It records E17 for ControlSets/generated catalogs/layer boundaries and E18 for
addressable versus live data. The [exercise findings](MVP1-SDL-Exercise-Findings-and-Extensions.md)
index those candidates. The existing 66-file corpus and general EBNF remain
unchanged until the candidate semantics are reconciled and selected.
