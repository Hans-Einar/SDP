# MVP1: whole-system SDL authoring exercise

Date: 2026-09-17  
Status: proposed design corpus, with selected and inherited obligations explicitly marked.  
Language: `sdl-mvp1-exercise 0.1`, an isolated candidate extension profile.

For the later system-wide conceptual consolidation, see
[Checkpoint #1](../../SDP/History/checkpoint-1/README.md). Its Value/ControlSet and
Functionality proposals do not change the meaning or validation of this baseline.

Start at [System.design](SDL/MVP1/System.design). It declares the one System,
registered constituents, Channels, endpoint roles and modes, and explicitly
includes the other 65 source files. The filename is a convention for this exercise;
the proposed language does not require that entry filename.

This is a full-scope **design authoring exercise**, not a running system. It covers
all 17 registered constituents and maps all 332 source requirement identifiers
(99 system and 233 local). It contains 109 Functionality definitions, 27 Capabilities,
14 Features, 20 scenario graphs and 11 Channels. These counts establish inventory
coverage, not completeness of every predicate, payload, algorithm or implementation.
The 23 named Gaps keep those missing decisions visible.

## Purpose and authority

The owner requested a system-wide trial before extending the parser. Necessary
language additions were first recorded in
[findings and local extensions](../../SDL/docs/studies/MVP1-SDL-Exercise-Findings-and-Extensions.md).
The general [language definition](../../SDL/docs/studies/Design-Language-Definition.md) and
`design-core 0.1` implementation have not been changed for this exercise.
These files cannot be passed to that parser as conforming inputs.

Ponsse evidence is pinned to
[`7589a5811a21ec6bb13979612060ceff0f690df8`](https://github.com/Hans-Einar/ponsse/tree/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1).
The integration branch head was checked on September 17 and still matched.
The owner's selection of the September 10 UI direction is recorded separately
from the older repository requirements; historical discussion headers do not
cancel that owner decision.

The active Concept1 checkout lacks this MVP1 tree and has unrelated work.
Consequently the trial lives here in SDP. A later transfer into the correct Ponsse
checkout can preserve the `SDL/MVP1` tree without changing identities.

## How to read the system

1. Read [System.design](SDL/MVP1/System.design) for boundaries and communication.
   MachineServiceUpstream carries both request and response MessageSets; upstream
   names describe topology, while sender/receiver states direction. There is no
   private Channel per Feature.
2. Read [Features.design](SDL/MVP1/Features.design), then a scenario below.
   Action Steps call named Functionalities, with explicit Container context for
   reusable library code; branch and terminal Steps may be control-only.
   `receives`/`emits` connect selected steps to Channels.
3. Follow `owns` to the accountable Unit, `realizes` to its Capability and
   `provides` to a boundary. Layer grouping adds a view of responsibility,
   not a second owner.
4. Read [State.design](SDL/MVP1/State.design) and
   [domain value types](SDL/MVP1/Contracts/DomainValues.design) for state,
   identities and revision distinctions. These are typed design descriptions,
   not executable initializers or a full value-expression language.
5. Read the relevant contract and
   [common invariants](SDL/MVP1/Contracts/Common.design). Wire-schema references
   identify existing authority; English obligations do not generate a complete
   communication runtime.
6. Consult [decisions](SDL/MVP1/Governance/Decisions.design),
   [gaps](SDL/MVP1/Governance/Gaps.design) and
   [partial source bindings](SDL/MVP1/Governance/Bindings.design) before drawing
   implementation conclusions.

All names are system-global in this trial. A name identifies a design definition,
not automatically one runtime instance. Stateful Library definitions require
per-application/session allocation, which remains a named language gap.
Declarations may follow uses in another included file; linking needs the complete
source set. An optional type annotation asserts the declared type of a reference;
it neither declares a second object nor derives its meaning from the preceding line.

## Constituent coverage

The [registry](https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/CONSTITUENTS.md)
contains nine Containers and eight Libraries. The model preserves all of them.

| Registry code | Modeled constituent / responsibility entry | Kind | Pinned source location |
|---|---|---|---|
| CTR | [ContractRegistry](SDL/MVP1/Libraries/Contracts.design) | library | [contracts](https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/contracts) |
| SVK | [ServiceKit](SDL/MVP1/Libraries/ServiceKit.design) | library | [shared/servicekit](https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/shared/servicekit) |
| LBK | [LabKit](SDL/MVP1/Libraries/LabKit.design) | library | [shared/labkit](https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/shared/labkit) |
| UIR | [UIRuntime](SDL/MVP1/UI/Representation.design) | library | [shared/ui-runtime](https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/shared/ui-runtime) |
| BOX | [BoxUI](SDL/MVP1/UI/Presentation.design) | library | [shared/box-ui](https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/shared/box-ui) |
| TPR | [TaperModels](SDL/MVP1/Libraries/TaperModels.design) | library | [packages/taper-models](https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/packages/taper-models) |
| BCK | [BuckingCore](SDL/MVP1/Libraries/BuckingCore.design) | library | [packages/bucking-core](https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/packages/bucking-core) |
| SFC | [StanfordClassic](SDL/MVP1/Libraries/StanfordClassic.design) | library | [packages/stanford-classic](https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/packages/stanford-classic) |
| SIM | [P1000Simulator](SDL/MVP1/P1000Simulator/Domain.design) | container | [services/p1000-simulator](https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/services/p1000-simulator) |
| SUI | [SimulatorUI](SDL/MVP1/UI/SimulatorDomain.design) | container | [web/simulator-ui](https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/web/simulator-ui) |
| MCH | [MachineService](SDL/MVP1/MachineService/Domain.design) | container | [services/machine-service](https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/services/machine-service) |
| MLB | [MachineLab](SDL/MVP1/MachineLab/Diagnostics.design) | container | [services/machine-lab](https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/services/machine-lab) |
| BKS | [BuckingService](SDL/MVP1/BuckingService/Planning.design) | container | [services/bucking-service](https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/services/bucking-service) |
| BLB | [BuckingLab](SDL/MVP1/BuckingLab/Diagnostics.design) | container | [services/bucking-lab](https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/services/bucking-lab) |
| BWEB | [BuckingWeb](SDL/MVP1/BuckingWeb/Adaptation.design) | container | [services/bucking-web](https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/services/bucking-web) |
| BUI | [BuckingUI](SDL/MVP1/UI/PonsseDomain.design) | container | [web/bucking-ui](https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/web/bucking-ui) |
| PGW | [P1000Gateway](SDL/MVP1/P1000Gateway/Domain.design) | container | [services/p1000-gateway](https://github.com/Hans-Einar/ponsse/blob/7589a5811a21ec6bb13979612060ceff0f690df8/MVP1/services/p1000-gateway) |

StemProfile and AptDomain are explicit Units within BuckingService; extraction
into new services is open. The logical UI subsystem groups BuckingUI, SimulatorUI
and renderer endpoints without inventing an accepted Go UI Host deployment.
BWEB retains independent Machine and Bucking connections. Lab services remain
optional observers. PhysicalPassive describes a disabled-write observational
configuration; modeled future TX obligations do not enable hardware output.

## UI horizontal structure

| Layer / concern | Source | Responsibility |
|---|---|---|
| Boundary | [Boundary](SDL/MVP1/UI/Boundary.design) | Logical subsystem and horizontal grouping. |
| Ponsse domain | [Operator](SDL/MVP1/UI/PonsseDomain.design), [Simulator](SDL/MVP1/UI/SimulatorDomain.design) | App-owned projections and intentions; simulator truth stays isolated. |
| Representation | [Representation](SDL/MVP1/UI/Representation.design) | Typed semantic identity, value quality and incremental collections. |
| Composition | [Composition](SDL/MVP1/UI/Composition.design) | Context lifetime and shared references. |
| Presentation | [Presentation](SDL/MVP1/UI/Presentation.design) | Validated schema, bindings, formatting and publication. |
| Renderer | [Renderer](SDL/MVP1/UI/Renderer.design) | Trusted views and actual render observation. |
| Support | [Transport](SDL/MVP1/UI/Transport.design), [Delivery](SDL/MVP1/UI/Delivery.design), [Commands/diagnostics](SDL/MVP1/UI/CommandsDiagnostics.design) | Independent sessions, bounded work, recovery, commands and diagnostic routing. |
| Concrete design objects | [Views](SDL/MVP1/UI/Views.design) | Two views of one value, stem series, APT matrix, app Compositions and Presentation schemas. |

The selected sequence is Representation → Composition → Presentation → Renderer.
The UIR/BOX responsibility assignment is a proposed realization of that choice;
it does not claim the packages already implement it. Renderer Steps using a
BuckingUI context identify the logical application they serve, not a settled
process location. The current grammar cannot fully allocate remote/embedded
renderer instances; `UiHostingGap` and `LibraryAllocationGap` record that limit.

## Contract coverage

| Contract source | Boundary |
|---|---|
| [Target](SDL/MVP1/Contracts/Target.design) | Machine ↔ selected SIM or passive PGW. |
| [Machine](SDL/MVP1/Contracts/Machine.design) | Machine ↔ Bucking, BWEB, MachineLab and CLI. |
| [Bucking](SDL/MVP1/Contracts/Bucking.design) | Bucking ↔ BWEB, BuckingLab and CLI. |
| [Lab](SDL/MVP1/Contracts/Lab.design) | Lab control/artifact inspection. |
| [Simulator control](SDL/MVP1/Contracts/SimulatorControl.design) | SUI/CLI ↔ simulator control, separate from Target. |
| [Operator browser](SDL/MVP1/Contracts/OperatorBrowser.design) | BUI ↔ BWEB, with explicit schema gaps. |
| [Presentation](SDL/MVP1/Contracts/Presentation.design) | Logical UI ↔ renderer; payload/recovery contract is a proposal. |
| [Permissive](SDL/MVP1/Contracts/Permissive.design) | External presence/permissive signal; incomplete contract. |
| [Physical serial](SDL/MVP1/Contracts/PhysicalSerial.design) | Passive physical observation only in this exercise. |
| [Common](SDL/MVP1/Contracts/Common.design) | Authority, identity, durability, isolation and responsiveness obligations. |
| [Ports/adapters](SDL/MVP1/Ports.design) | Physical observation, replay/files and console boundaries. |

MessageSets reference authoritative envelopes where available. Their variants
may have different delivery and authority rules. A declared sender is not granted
every command in an envelope, and a receiver is not promised every event. Exact
routing, capacity and recovery contracts remain necessary before runtime lowering.

## Scenario index

These are contract-level scenario graphs with named preconditions, branch
conditions and terminal outcomes. They are broader than happy-path sketches,
but not exhaustive state-machine exploration. Branch descriptions are English,
not executable predicates. `runs-alongside` permits overlap without specifying
a scheduler, thread count or interleaving proof.

[Functionality coverage](functionality-coverage.csv) lists all 109 definitions,
their owner, Capability, source file and direct scenario calls. The scenarios
directly call 73; the other 36 have responsibility/contract descriptions but no
explicit scenario call yet. This is a concrete behavior-coverage limit rather
than an assumption that unshown internal calls exist.

| Scenario | Design pressure exercised |
|---|---|
| [InspectMachine](SDL/MVP1/Scenarios/InspectMachine.design) | Target → Machine → BWEB → operator projection; unsupported evidence and unavailable calibration. |
| [PlanCurrentStem](SDL/MVP1/Scenarios/PlanCurrentStem.design) | Coherent samples → model → frozen inputs → K1 → currentness check; stale result retains prior plan. |
| [ComputeArchive](SDL/MVP1/Scenarios/ComputeArchive.design) | Background K-best, bounded work and cancellation; overlaps current-stem planning. |
| [ImportApt](SDL/MVP1/Scenarios/ImportApt.design) | Supported format/source identity → neutral APT domain → draft; importing is not activation. |
| [EditAptCell](SDL/MVP1/Scenarios/EditAptCell.design) | Typed draft edit with revision/authority validation; no renderer-owned business state. |
| [StoreApt](SDL/MVP1/Scenarios/StoreApt.design) | Domain-owned store/export; accepted draft and encoding compatibility stay distinct. |
| [ActivateApt](SDL/MVP1/Scenarios/ActivateApt.design) | Explicit activation, lifecycle and economics preconditions; import flag remains a contract gap. |
| [ReplacePresentation](SDL/MVP1/Scenarios/ReplacePresentation.design) | Prepare/validate/publish/reconcile/observe; invalid and stale candidates versus post-publication renderer failure. |
| [RecoverUiBaseline](SDL/MVP1/Scenarios/RecoverUiBaseline.design) | Shadow baseline and stream cursor commit; source and delivery sequence domains remain distinct. |
| [RecoverCommand](SDL/MVP1/Scenarios/RecoverCommand.design) | Same-ID retry/lookup; known, still unknown, conflict and expired outcomes remain distinct. |
| [SimulatorWireAction](SDL/MVP1/Scenarios/SimulatorWireAction.design) | Simulator control plane changes target behavior; hidden truth remains isolated. |
| [SemanticOperatorAction](SDL/MVP1/Scenarios/SemanticOperatorAction.design) | Normalized Machine command through BWEB; connection-bound authority. |
| [ReplayDiagnostics](SDL/MVP1/Scenarios/ReplayDiagnostics.design) | Explicit artifact/target selection and isolated replay; diagnostics do not own domain truth. |
| [GlobalStopFlow](SDL/MVP1/Scenarios/GlobalStopFlow.design) | Identified software stop, downstream coordination and retained diagnostic access. |
| [DegradedOperation](SDL/MVP1/Scenarios/DegradedOperation.design) | Machine remains observable through BWEB while Bucking is unavailable. |
| [CommitProduction](SDL/MVP1/Scenarios/CommitProduction.design) | Validate lineage and persist first; optional Lab absence/failure cannot invalidate the durable fact. |
| [RecoverTarget](SDL/MVP1/Scenarios/RecoverTarget.design) | Reconnect, explicit recovery range/choice and separate transport delivery/P1000 acknowledgment. |
| [UiDiagnosticRouting](SDL/MVP1/Scenarios/UiDiagnosticRouting.design) | Unhandled routing versus handler failure; bounded and non-recursive diagnostics. |
| [ShutdownSystem](SDL/MVP1/Scenarios/ShutdownSystem.design) | Stop owned work, drain/account for pending durable state, release subscriptions. |
| [ShareUiValue](SDL/MVP1/Scenarios/ShareUiValue.design) | Two Views bind the same Representation; disposing one must preserve the other. |

A committed production fact survives failed delivery or absent Lab evidence.
Likewise a published Presentation is not necessarily displayed. These distinctions
must survive future AST/IR translation; flattening them into a generic success/error
would lose essential design meaning.

## Requirements, decisions and evidence

[Requirements.design](SDL/MVP1/Governance/Requirements.design) mirrors source
obligations and links them to accountable design objects.
[requirement-coverage.csv](requirement-coverage.csv) is the corresponding review
index. This mirror is derived evidence, not a replacement for the authoritative
Ponsse documents. `inherited` means taken from the pinned source; it does not
upgrade an upstream proposed requirement to accepted.

Two earlier obligations are explicitly superseded for this target:
`MVP1-SYS-REQ-050` and `BOX-REQ-009`, whose optional-renderer wording conflicts
with the selected UI direction. Their original text remains available and their
replacement decision is linked. Other preserved obligations still constrain the
new model. A broad requirement link to a Container says where the obligation
belongs; it does not establish clause-by-clause sufficiency.

[source-inventory.json](source-inventory.json) records the revision, registry,
all expected requirement IDs and SHA-256 for 31 cited source files. Four observed
UI mechanism bindings are deliberately partial; two target bindings are explicitly
open. The model does not assert 109 implemented source functions.

## Findings and limits

The [findings document](../../SDL/docs/studies/MVP1-SDL-Exercise-Findings-and-Extensions.md)
records 15 extension candidates and the concrete modeling mistakes/omissions found.
The subsequent [Dataset/Datagram proposal](../../SDL/docs/studies/SDL-Datasets-Datagrams-and-Data-Contracts.md)
adds candidate E16 for internal datasets, optional databases and Dataset-linked
Datagram families with contract-defined variants, source identity, missing values
and snapshot consistency. It has not yet been applied
to this corpus or its inventory audit; the 15 candidates above describe the
original modeling pass.
The later [ControlSet/layer/data-access study](../../SDL/docs/studies/SDL-ControlSets-Layer-Boundaries-and-Data-Access.md)
records E17–E18, including the owner-selected move from authored MessageSets to
generated catalogs. Baseline `.design` files still use their original MessageSets;
their migration is future work, not a completed compiler feature.
The [23 Gap objects](SDL/MVP1/Governance/Gaps.design) cover:

- UI hosting, stateful library allocation, optional Stem/APT service extraction.
- Presentation/browser/permissive contracts and the APT edit/store/activation lifecycle.
- Presentation publication races, renderer failure, transient edits and collection taxonomy.
- Delivery budgets, journal scope and full production provenance.
- Unsupported/new StanForD codecs, P1000 evidence and the physical-write gate.
- Executable predicates, numerical/codec bindings and source-level realization links.
- System shutdown allocation, quiescence order, deadlines and failure-time cleanup.

Inventory coverage is complete relative to the pinned registry and requirement IDs.
Behavioral precision and executable coverage remain incomplete. The current
corpus cannot build a runnable SDL IR, establish physical safety, or verify the
accepted visual layout. Those limits are explicit design work, not permission
for a compiler to invent defaults.

## Reproducible authoring audit

From the SDP repository root:

```sh
python3 experiments/mvp1_sdl/audit_inventory.py
python3 experiments/mvp1_sdl/audit_inventory.py --source-root /path/to/pinned/ponsse
```

Python 3.6+ and the standard library suffice. The optional source root must contain
the pinned `MVP1/` tree. This checks source hashes and independently extracts the
requirement IDs as well as checking the authored inventory.

[audit_inventory.py](audit_inventory.py) checks entry membership, declaration
uniqueness, identifier references outside quoted text, containment cycles,
single Functionality ownership, channel membership for explicit traffic,
reachable scenario steps/terminal paths, constituent coverage and requirement
mapping. It does **not** implement the SDL grammar, type checker, payload
validator or interpreter. In particular, finding one path to an outcome is not
proof that execution terminates under every branch choice.

[Recorded audit result](inventory-audit.json): pass, 66 files; detailed declaration
counts are recorded in the report, along with 332 requirement mappings,
64 participation facts and 36 explicit traffic facts.
All 31 source hashes were checked. No product code, physical I/O or renderer
execution was part of this exercise.

## Next bounded language decision

Use the APT edit/activation path and Presentation replacement path to decide
which of E01–E15 should enter the general language. Define executable guards,
typed operations and failure/commit semantics before claiming IR completeness.
Then add positive and negative parser/linker tests for each adopted construct.
Do not promote the entire exercise profile simply because this corpus uses it.
