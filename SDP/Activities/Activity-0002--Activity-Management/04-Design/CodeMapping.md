# Activity model to implementation mapping

Target concern: b3e45cab-c629-490b-9192-3ed2b4628387:00002. This authored mapping
supplements the preliminary Blueprint's honest generic code-conformance unknown.
It does not claim an accepted SDL release or automatic code synthesis.

All paths below are relative to SDPTool/ unless a repository is stated explicitly.
The three modeled units are responsibility groups inside one SdpToolHost process;
they are not three new deployed services or wrappers around an old implementation.

| SDL responsibility | Concrete owner / boundary | Verification |
| --- | --- | --- |
| ActivityService / CreatePlannedActivity | activities/create.go, Service.Create | records_test; activity_integration_test CLI/action journeys |
| CreateActivitySession | Service.CreateSession; real identity/phase/plan/milestone/predecessor/handoff validation | records_test Session cases and actual CLI create |
| ValidateActivityTransition | activities/lifecycle.go, Transition/Update; replay.go validates persisted chain/batches | lifecycle_test concurrency, dependency, disposition and stale cases |
| AdoptActivityRecords | activities/admission.go, AdoptAdmission; explicit store inventory and canonical event | admission_test and root bind/adopt/hold journey |
| CheckImplementationScope | activities.CheckAdmission, governance/activity_admission.go | direct Start/Prepare/Activate and final dispatch tests |
| RequestExecutionQuiescence / ActivityExecutionDispositionPort | PendingLifecycle and injected governance.ActivityExecutionInspector | pending/exact-send/cross-store tests; actual controller callback test |
| ActivityRecordRepository / PublishActivityRecordSet | activities/persistence.go, projecthistory.Transaction.AppendBatch | per-file crash faults, exact receipt retry, external-conflict refusal |
| RetainDescriptorRevisions | Activities/.descriptors hash store, immutable old pins and state-preserving update event | replacement interruption and real manual0.1 upgrade-copy tests |
| RecoverActivityRecordOperation | RecoverLifecycle, retained intent/preimages and fixed event/result bytes | stopped transaction and pending quiescence retry tests |
| ActivityProjectionService / ProjectActivityNavigation | activities/projection.go; root activity_navigation/expansion.go final protocol | source/legacy/invalid-state cases, paged native tree tests |
| ProjectActivityTimeline | activities/projection.go and visual.go; explicit causal intervals/marker joins | unknown/clock reversal, unrelated-event exclusion, filters and XML escaping |
| CanonicalManagementHistory / ActivityEventStream | ProjectManagement/Ledger.ndjson via shared projecthistory | shared CAS/batch/reservation tests; full mixed-history validators |
| ActivityDocumentStore / ActivityDocuments / ActivitySessionDocuments | descriptor/mandate/phase/root Sessions homes; .operations recovery metadata | layout/README checks, snapshots and mutable Session log tests |
| ExternalCliCaller / ExternalCoordinatorCaller | activity_cli.go and action_activities.go over activities.Operations | strict input/attribution and registered-operation parity |
| ExternalViewer / ActivityViews | ActivityView / ExpandActivity source revision and private bundle | native XFMD model-free selection and atomic deferred merge |

ActivityCommandShape and related SDL message field declarations are bounded
structural placeholders. activities/types.go, versioned schemas and Activities.md
specify executable fields/constraints. Changed reviewed input revisions are refused;
Blueprint lifecycle/freshness remains owned by its existing service.

The collaborating XFMD consumer is independently mapped under its existing
FUNC-032 SDP Project Navigation ownership. It renders final producer roots, targets
and pages; it neither replays Activity events nor classifies lifecycle state.
Separate worktree evidence lives in XFMD SDP/05--Implementation/SDP-Navigation/
evidence/Activity-A3-001. Existing native resource protections remain explicit.

No alternate writable Activity ledger, duplicate PG execution state or transitional
CLI-internal bridge remains. .admission store records are locators, not run state;
reading an Activity catalog does not silently enable implementation admission.
Read-only review/leader/tool investigation remains outside implementation quiescence.
