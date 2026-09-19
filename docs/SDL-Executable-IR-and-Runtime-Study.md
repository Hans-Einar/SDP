# Executable SDL IR and runtime study

Date: 2026-09-17

Status: proposed execution direction, not an implemented interpreter or a complete
execution specification. The owner requests that this direction inform language
development. Instructions, schemas, scheduling rules and supported adapter
contracts remain to be established through small conformance examples.

Related: [language definition](Design-Language-Definition.md#16-future-executable-sdl-ir),
[source tree and compilation](SDL-Source-Tree-and-Compilation-Study.md),
[MVP1 design example](MVP1-Design-Language-Example.md), and
[conformance scenarios](Design-Language-Conformance-Scenarios.md).

## 1. Intended execution target

The owner proposes generating an intermediate language from a complete design,
executing it in an interpreter, injecting data through Channels and observing
state changes, console output or a real UI renderer. Ports/adapters mediate the
outside world. This can make design behavior testable before its final native
implementation exists and expose assumptions that structural diagrams cannot.

Recommend a **versioned SDL execution IR** with precisely defined operations. It
can have a readable textual serialization, but need not be a general-purpose
language or generated Python/Go source. Its instructions should be more explicit
about execution than the design source, not allow the compiler to invent missing
behavior. "Standardized" initially means a documented project-owned contract with
conformance cases, not an existing industry standard named SDL IR.

```text
SDL files -> ASTs -> linked and checked system model
                        |-> structural blueprints / public exports
                        |-> execution-profile validation
                            -> executable IR -> interpreter
                                                |-> state and event trace
                                                |-> typed ports <-> adapters
                                                |-> Presentation contract <-> renderer
```

Keep the structural model and execution IR distinct. An architecture model can
contain an unimplemented Functionality. An interpreter cannot execute its name
as though that supplied an algorithm. Lowering needs a defined body, a versioned
external implementation binding, or an explicitly marked test double. Missing
realization must produce an execution-coverage error, not a successful no-op or
generated guess. A bounded runnable slice can use declared doubles; its report
must identify their scope and assumptions.

## 2. What parsing cannot supply

Before a selected scenario can execute, its transitive behavior and dependencies
need enough information to answer these questions:

| Concern | Minimum executable meaning |
|---|---|
| Data | Message schemas, units, identity, validity, value operations and error cases. |
| State | Owned variables, initial values, lifecycle, reset and permitted mutations. |
| Behavior | Handler inputs, transitions, guards, outputs and failure outcomes. |
| Communication | Routing, fan-out versus competing delivery, ordering scope, queue limits, overflow and correlation. |
| Completion | Accepted, started, completed, failed and unknown outcomes where applicable; cancellation semantics. |
| Time | Clock source, deadlines, timer ordering and treatment of late events. |
| Concurrency | Atomic steps, yield/wait points, state isolation and permitted interleavings. |
| Effects | Typed port contracts, allowed adapter bindings, failures and recording policy. |
| Evidence | Source/model identities, IR/runtime versions and inputs needed to explain a run. |

Typed value input/output, commands and streams cover many interfaces. They do not
alone specify transaction boundaries, persistence, timing, ownership or algorithms.
The initial goal should be event-driven application coordination and UI scenarios,
not a claim to model every hard real-time, numerical or operating-system workload.

## 3. Async and execution contexts

Do not define `async` as "create a thread." It must identify a semantic distinction,
such as completion occurring after submission and observed through an explicit
result. The adjective/property definition must specify its subject type,
completion/error contract and interaction with waiting and cancellation. `async`
remains unsupported syntax until those rules are defined.

Concurrency permits overlapping progress; parallel execution is one possible
realization. For example, Python's asyncio documentation describes cooperative
scheduling in which an event loop runs one task at a time. This is an existence
example, not a decision to implement SDL with asyncio.
[Python task documentation](https://docs.python.org/3.12/library/asyncio-task.html).

Proposed initial model:

- A Container has isolated modeled state and explicitly assigned logical execution
  contexts. The first profile can serialize state-mutating steps per context;
  neither a folder nor a Container automatically requires its own OS thread.
- A Channel is a routing/delivery contract, with queues where that contract needs
  them. It does not automatically receive a worker thread or processing behavior.
- A scheduler advances enabled handlers and delivery events. Defined atomic steps
  run to completion; async operations have explicit suspension/completion points.
  Waiting cannot hold an implicit global system lock.
- Independent progress, starvation and blocking need scheduling rules. Later
  thread/process mappings must preserve the observable contract.

For the first simulator, prefer a deterministic scheduler, virtual clock and
declared tie-breaking rule. Record inputs, scheduler choices and random seeds
when used. Repeatability also requires controlling or recording external
observations. One deterministic run does not explore all schedules: vary relevant
interleavings, including validation versus publication while measurements advance.

SCXML is an established reference for event-driven state-machine execution,
including run-to-completion semantics. Evaluate it when defining transition
semantics or a possible backend; it does not establish a complete mapping of
SDL's contract, source-binding and UI requirements by itself.
[W3C SCXML specification](https://www.w3.org/TR/scxml/).

## 4. Ports and adapters for external interaction

A **port** is the runtime-facing typed interaction boundary. An **adapter** connects
it to a fixture, file, console, transport or renderer. Align these terms with
Interface/Channel before making them SDL keywords. A Channel connects logical
participants; a port exposes a participant's boundary; a selected adapter realizes
access to an environment. They must not become synonyms for one declaration.

Recommend that all external effects cross declared ports/adapters, including
console interaction. `stdin`, `stdout` and `stderr` are familiar console bindings,
not privileged paths that bypass modeling and recording. Keep application stdout,
application diagnostics and interpreter diagnostics distinguishable.

Do not force every boundary into a generic byte stream:

| Interaction | Required distinctions |
|---|---|
| Byte/text stream | Encoding, framing, end-of-stream, partial reads/writes and disconnects. |
| Typed message stream | Message schema, delivery/order policy and backpressure. |
| Request/result | Identity, acceptance, completion, timeout, retry and unknown outcome. |
| Presentation exchange | Schema revisions, values, intents, resynchronization and displayed-state observations. |

These are proposed contract categories, not finalized syntax or a requirement to
build every adapter now. Byte adapters need framing/decoding before bytes satisfy
a typed message contract. File adapters need explicit read/write scope, mode,
lifetime and failures; persistence is not implied by a successful send.

Use fixture/replay and in-memory adapters first. Time and randomness must also be
controlled services or recorded inputs for repeatability. A blocked adapter must
not prevent unrelated simulated contexts from progressing indefinitely.

Serial I/O can be an eventual adapter category, but the Ponsse project authorizes
offline replay/simulation and forbids live serial transmission and actuation.
The initial profile must have no default hardware-output path. Captured serial
bytes can be replayed through the decoding contract; a live adapter requires a
separate authorized milestone. This study adds no hardware tooling.

## 5. Renderer outside the interpreter engine

The proposed boundary fits the selected MVP1 direction: interpreted behavior
maintains Representation/Composition/Presentation, while an external renderer
adapter consumes the Presentation contract and returns typed intents. Hosting
the renderer outside the interpreter does not place it outside the modeled
System; runtime hosting and logical membership are distinct boundaries.

The contract needs more than command/value strings:

- Initial schema/state, stable Representation and view identities, and revisions.
- Updates, layout replacement, supported schema features and compatibility.
- User intent identity/correlation, preserving domain authority.
- Delivery/backpressure, reconnect and coherent resynchronization.
- Observations separating schema publication from actual rendering success.

An actual React or Fyne renderer can be reused only if it implements this contract,
possibly through a host-specific adapter. Framework choice alone does not provide
compatibility. The renderer owns its UI loop and framework threading requirements;
SDL should not prescribe an extra renderer thread by default.

Use a headless contract adapter first, then a real renderer witness. Headless
checks can establish schema/value/intent behavior but not viewport layout, focus,
accessibility or rendering performance. Both must exercise the same declared
contract; a separate imitation protocol would weaken the test.

## 6. Initial executable profile and IR

Recommend an explicitly versioned execution profile: a supported subset with
additional completeness rules, not a dialect with alternative fact spellings.
It must not silently execute every syntactically valid model.

Entry criteria:

1. Workspace linking, type and annotation checks pass.
2. The selected scope has initial state and complete reachable handlers, or
   explicitly identified test doubles/external implementations.
3. Every used Channel has message, routing, delivery and failure semantics.
4. Every enabled effect has an allowed adapter; clocks and nondeterminism are accounted for.
5. Scheduling/atomicity and lifecycle rules cover stop, reset and error handling.
   Resource bounds produce visible outcomes rather than silent loss.
6. IR instructions retain source/model provenance so diagnostics and traces link
   back to scenarios, statements and any implementation bindings.

An artifact could initially be called `execution.ir.json`; that is a proposed
filename, not a frozen encoding. Define schema, instruction semantics and version
compatibility through conformance cases. Avoid arbitrary embedded host-language
code in the first profile. A domain algorithm can use an explicit typed, versioned
function binding, with real implementation versus substitute recorded.

Traces should preserve cause/correlation, logical time, relevant revisions,
transitions, emissions and failures. Tests assert independently specified outcomes
and invariants, not merely the interpreter's own output. Replay of a recorded
adapter response establishes a different claim from a fresh external interaction.

## 7. First experiments

Use a small runnable design slice:

1. Inject a typed machine observation from replay. Independent consumers update
   state under explicit policies. Check that required information reaches both
   consumers without accidental competition for one queue.
2. Replace a Presentation while observations continue. Exercise invalid bindings
   and a Composition revision change between validation/publication. Assert retained
   Representation identity and the specified stale-candidate outcome.
3. Route an intent through a simulated domain handler and correlate its result.
   Inject timeout/late-result cases; do not equate timeout with cancellation or
   generate a new command identity for an uncertain retry.
4. Attach a headless Presentation adapter, then one real renderer implementing
   the contract. Observe publication, delivery and visual completion separately.

Require source-to-IR traceability, malformed-message rejection, visible missing
handlers, repeatable virtual-time traces and bounded delivery under load. Test
different relevant event orders, not just one baseline. Later thread/process
backends should satisfy the same observable contracts.

This direction informs language extensions now, while multi-file linking and
communication contracts remain immediate prerequisites. Implement the first IR
against one fully specified scenario before claiming general executable-design
support.

## 8. Completeness feedback and future compiler targets

Owner follow-up, 2026-09-17: executable IR generation should expose omissions
that prevent building the described system. A later compiler could lower SDL IR
to LLVM IR and support further implementation targets. This is a direction to
investigate, not a selected backend or an implementation claim.

### 8.1 What a failed build can teach us

Execution-profile validation can turn missing design into concrete diagnostics:
an input has no handler, a handler has no executable body, a state has no initial
value, a required result path is absent, a Channel lacks a delivery policy, or an
effect lacks an adapter. Point back to the relevant SDL fact and identify which
execution rule cannot be satisfied. Do not choose an undocumented default to make
the build succeed.

The limit matters: a compiler can demand what the language/profile requires. It
cannot discover an unstated user requirement simply because compilation succeeds
or fails. A legal handler can still implement the wrong behavior. Successful
lowering establishes sufficient specified structure for that execution target,
not requirements completeness, architectural quality or correctness for all runs.
Retain scenario tests, invariants, explicit assumptions and owner-facing review.

This gives a useful feedback loop: a concrete scenario exposes a missing
distinction; the language defines it; the compiler checks its presence and
consistency; interpreter traces test its consequences. Avoid adding executable
detail solely to satisfy a tool if the correct result is an explicit modeling gap.

### 8.2 LLVM as an optional execution backend

A feasible future path is SDL execution IR -> LLVM IR -> native object/machine
code or JIT execution, with an SDL runtime supplying the chosen message,
scheduling and adapter operations. LLVM's code generator translates LLVM IR into
target machine code; this is different from generating maintainable source in
arbitrary high-level languages.
[LLVM code generator](https://llvm.org/docs/CodeGenerator.html).

The work would primarily be an SDL lowering/compiler front end using LLVM's
existing targets, not a new machine backend inside LLVM. Channels, state ownership,
async completion and Presentation are not automatically implemented by emitting
LLVM instructions. Define how they map to generated state machines/functions and
runtime calls, including memory ownership, data layout, ABI, error propagation,
numeric semantics, lifecycle and scheduling. Preserve source maps and a separate
model index when optimizations change the generated code's structure.

If the objective is source generation in Go, TypeScript or another implementation
language, prefer a separate generator from the higher-level SDL IR, where these
design concepts are still explicit. Do not depend on recovering them from
optimized LLVM IR. Dedicated lower-level-to-source tools may be possible, but LLVM
is not a universal automatic inverse of all language front ends.

```text
Checked SDL execution IR
  |-> reference interpreter
  |-> LLVM lowering + runtime -> compiled/JIT execution
  |-> language-specific generator + runtime -> target source
```

All paths are proposed. A backend must reject unsupported semantics rather than
silently approximate them. Compare required observable outcomes, ordering and
invariants against the interpreter under the same recorded inputs. Internal steps
and physical timing need not match unless the contract requires them. Differential
tests help detect divergence but cannot alone prove the shared specification is right.

### 8.3 Evaluate MLIR before building extensive lowering infrastructure

MLIR supports staged lowering from domain-oriented representations toward LLVM
IR; its Toy tutorial demonstrates that progression. It is therefore a candidate
implementation framework for an SDL IR, not a mandatory first dependency.
[MLIR Toy tutorial](https://mlir.llvm.org/docs/Tutorials/Toy/).

MLIR's term "dialect" describes an internal set of operations/types. An internal
SDL dialect would not authorize synonyms or alternative semantics in the authored
SDL language. The externally visible language and its contracts remain singular.

For a concrete source-output example, MLIR's EmitC dialect supports conversion
toward C/C++ output. This does not automatically translate arbitrary SDL or LLVM
programs: supported operations and explicit conversions are still needed.
[MLIR EmitC](https://mlir.llvm.org/docs/Dialects/EmitC/).

First make one scenario executable with the small interpreter profile. Then
evaluate direct LLVM lowering versus MLIR on that same scenario, measuring the
semantic mapping, runtime requirements, diagnostics/provenance and maintenance
cost. This keeps LLVM open as an implementation route while making runnable
design and reliable feedback the immediate deliverable.
