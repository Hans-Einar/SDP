# Type inventory — A5

[Up](index.md)

The same model ID is used at every level; object levels are unspecified.

## activity (8)

### ActionsDelivery

### BlueprintAssessmentDelivery

### BlueprintCatalogueDelivery

### BlueprintLifecycleDelivery

### BlueprintPreviewDelivery

### ExistingInstallFacade

### ExistingReadFacade

### RoutineCoreDelivery


## actor (3)

### ProjectMaintainer

### ProjectReader

### ProjectSupervisor


## capability (6)

### BlueprintAssignmentControl

### DiagnosticBlueprint

### ProjectInventory

### ReviewedInstallation

### RoutineCoordination

### TypedLocalInvocation


## channel (1)

### RoutineCalls


## container (1)

### SdpToolHost


## contract (3)

### RoutineCallsProtocol

### RoutineOutcome

### RoutineRequest


## database (1)

### OperationalRunStore


## dataset (1)

### RunSnapshots


## feature (8)

### BlueprintAssignments

### BlueprintPreview

### DiscoverableActions

### GovernedProjectOperations

### InstallationFacade

### ProjectFacade

### RetainedBlueprintNavigation

### ScopedBlueprintReadiness


## field (11)

### RoutineOutcomeCoverage

### RoutineOutcomeDisposition

### RoutineOutcomeEventCursor

### RoutineOutcomeRevision

### RoutineRequestAction

### RoutineRequestCheckoutIdentity

### RoutineRequestEvidenceReferences

### RoutineRequestExpectedRevision

### RoutineRequestOperationIdentity

### RoutineRequestProjectRoot

### RoutineRequestRunIdentity


## functionality (13)

### AssessBlueprintEvidence

### BrowseRetainedBlueprints

### CheckTransitionAuthority

### DescribeRegisteredProject

### DescribeSupportedActions

### GenerateBlueprintPreview

### GovernBlueprintAssignment

### InvokeReviewedInstallation

### InvokeTypedAction

### ReconcileRunRevision

### ResolveApplicableRoutine

### ResolveCheckoutAuthority

### VerifyRetainedBlueprints


## interface (6)

### CanonicalProjectHistoryLibrary

### ExistingInstallationLibrary

### ExistingSdlProjectionLibrary

### ModelSnapshotLibrary

### ProjectRecords

### TrustedAssignmentBinding


## message (3)

### RequestTransition

### TransitionAccepted

### TransitionRefused


## mode (3)

### InstallationHandling

### ProjectReading

### RoutineHandling


## scenario (2)

### RoutineAccepted

### RoutineStaleOrUnauthorized


## unit (12)

### ActionRegistry

### BlueprintAssessment

### BlueprintCatalogue

### BlueprintFacade

### BlueprintLifecycle

### ContextFacade

### ExistingInstallationFacade

### ExistingNavigationFacade

### ExternalCommandAdapter

### JsonInvocationAdapter

### RoutineCoordinator

### RunRepository


## usecase (4)

### BrowseProjectFacts

### CoordinateGovernedWork

### InspectProposedModelChange

### InstallReviewedProcess
