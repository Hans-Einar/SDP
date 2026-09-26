# Type inventory — A5

[Up](index.md)

The same model ID is used at every level; object levels are unspecified.

## activity (42)

### BuildNativeRealization

### ExportUiDocumentation

### G1FrontendPort

### G1M1ParserAndAst

### G1M2ValidationAndNormalization

### G1M3Concept1AndDumps

### G2LayoutAndPresentation

### G2M1RelativeMeasurement

### G2M2SharedSvgGeometry

### G2M3FyneInteractions

### G2M4RichContent

### G3M1TypedUiSession

### G3M2CandidatePublication

### G3M3CompatibleState

### G3UiRuntimeAndReload

### G4M1SdlFrontend

### G4M2TypedExecution

### G4M3UiDomainBinding

### G4M4DomainReload

### G4SdlRuntimeAndBinding

### G5M1GeneratedGo

### G5M2BehaviorParity

### G5M3DocumentationExport

### G5M4RetirePython

### G5NativeGeneration

### G6M1StaticNavigation

### G6M2OnDemandViews

### G6M3XfmdNavigation

### G6M4SessionPublication

### G6M5SemanticNotation

### G6M6ClassViews

### G6NavigableDocumentation

### InspectDesignSource

### InspectSdlSource

### InspectSduiSource

### RealizeDesign

### ReloadBoundModels

### ReloadDesignSession

### ReloadUiModel

### RunBoundUiAction

### RunDesignSession

### RunUnboundUiPreview


## actor (2)

### DesignAuthor

### DesignReviewer


## capability (14)

### BoundInteraction

### DevelopmentReload

### DomainOperations

### ExecutableDesign

### InteractiveSession

### MeasuredPresentation

### NativeInteraction

### NativeRealization

### RichContent

### SdlSourceModel

### SduiSourceModel

### SourceDiagnostics

### SourceLoading

### StaticDocumentation


## channel (22)

### FramePresentationCalls

### GoBuildCalls

### GoDomainCalls

### GoGenerationCalls

### LayoutCalls

### ModelReloadCalls

### NativeUiActions

### SdlActionCalls

### SdlCompilationCalls

### SvgExportCalls

### UiAstCalls

### UiCompilationCalls

### UiDomainActions

### UiGenerationEvents

### UiNormalizationCalls

### UiTokenizationCalls

### UiValidationCalls

### ViewDisplayCalls

### ViewLaunchCalls

### ViewNavigationCalls

### ViewProjectionCalls

### ViewPublicationCalls


## container (4)

### CommandLineHost

### FyneHost

### ViewServiceHost

### XfmdDocumentHost


## contract (45)

### ActionArguments

### ActionOutcome

### AstArtifactContract

### DesignSourceRecord

### FramePresentationCallsProtocol

### GeneratedGoContract

### GoBuildCallsProtocol

### GoDomainCallsProtocol

### GoGenerationCallsProtocol

### LayoutArguments

### LayoutCallsProtocol

### ModelReloadCallsProtocol

### NativeBuildContract

### NativeUiActionsProtocol

### NormalizedModelContract

### PreparedFrameContract

### PresentationOutcome

### ReloadArguments

### ReloadOutcome

### SdlActionCallsProtocol

### SdlCompilationCallsProtocol

### SvgDocumentContract

### SvgExportCallsProtocol

### TokenArtifactContract

### UiAstCallsProtocol

### UiCompilationCallsProtocol

### UiDomainActionsProtocol

### UiGenerationContract

### UiGenerationEventsProtocol

### UiNormalizationCallsProtocol

### UiSessionRecord

### UiTokenizationCallsProtocol

### UiValidationCallsProtocol

### ValidationOutcomeContract

### ViewBundleContract

### ViewDisplayCallsProtocol

### ViewFailureContract

### ViewLaunchCallsProtocol

### ViewNavigationCallsProtocol

### ViewOpenContract

### ViewOpenedContract

### ViewProjectionCallsProtocol

### ViewPublicationCallsProtocol

### ViewReferenceContract

### ViewSelectionContract


## database (1)

### DesignSourceArchive


## datagram (1)

### UiGenerationNotices


## dataset (2)

### DesignSourceDocuments

### UiSessionState


## encoding (1)

### UiGenerationWire


## feature (8)

### DesignDocumentation

### InteractiveUiPreview

### LiveModelReload

### NativeGoAssembly

### NavigableDesignDocumentation

### StructuralModelInspection

### TypedDesignInspection

### TypedDomainBinding


## field (54)

### ActionGeneration

### ActionInputText

### ActionOutputText

### ActionStatusCode

### ActionSymbol

### AstArtifact

### FailedViewRequestId

### GeneratedGoSources

### LayoutModelArtifact

### LayoutViewportHeight

### LayoutViewportWidth

### ModelDiagnostics

### ModelIsValid

### NativeBuildDiagnostics

### NativeBuildSucceeded

### NormalizedModelArtifact

### NoticeGeneration

### NoticeVersion

### OpenViewConsumerId

### OpenViewEntryPath

### OpenViewLeaseId

### OpenViewPaneId

### OpenViewRequestId

### OpenViewRevision

### OpenViewWindowId

### OpenedViewRequestId

### OpenedViewRevision

### PreparedFrameArtifact

### PresentationReady

### ReloadDiagnostic

### ReloadPublishedGeneration

### ReloadSourceRevision

### ReloadSourceText

### SessionDraft

### SessionGeneration

### SourceDocumentRevision

### SourceDocumentText

### SvgDocumentText

### TokenArtifact

### ViewBundleBytes

### ViewBundleRevision

### ViewConsumerId

### ViewEntryPath

### ViewFailureCode

### ViewFailureDiagnostic

### ViewLeaseId

### ViewManifestPath

### ViewPaneId

### ViewProjectId

### ViewPublishedRevision

### ViewRequestId

### ViewSelector

### ViewSourceRevision

### ViewWindowId


## functionality (116)

### AllocateGeometry

### ApplyPropertyBatch

### ApplySemanticNotation

### BuildGeneratedApplication

### BuildPreparedFrame

### BuildSdlAst

### BuildSduiAst

### CancelPendingActions

### CaptureNavigationTarget

### CheckDomainStateCompatibility

### CheckExecutionCompleteness

### CheckFunctionSignatures

### CloseSdlInstance

### CloseUiInstance

### CoalesceSourceChanges

### ComposeHeadlessExport

### ComposeInteractiveSession

### ComposeMarkdownDocument

### ComposeViewPackage

### ComputeClipping

### ConnectTypedWidgetHandles

### CoordinateSdlCompilation

### CoordinateSduiCompilation

### CorrelateActionResult

### CorrelateUiResult

### CreateSdlInstance

### CreateUiInstance

### DisconnectBindings

### DispatchUiEvent

### DispatchViewOpen

### EvictUnusedViewBundles

### ExpandUiDefinitions

### ExportConsoleSnapshot

### ExportModelInventories

### ExportSvgSnapshot

### ExportViewpointMarkdown

### GenerateBindingRegistration

### GenerateModelConstructors

### GenerateViewNavigation

### HandleFocusAndTextInput

### IdentifySourceRevision

### IndexViewpointLevels

### InvokeRegisteredFunction

### KeepLastValidModels

### KeyViewRevision

### ManageDomainState

### ManageWidgetIdentities

### MatchCompatibleWidgets

### MeasureMarkdownContent

### MeasureUiContent

### MigrateOrResetDomainState

### NormalizeSdlModel

### ObserveSourceChanges

### PerformDomainOperation

### PrepareCandidateModels

### PrepareDiagramResource

### PrepareMarkdown

### PreserveCompatibleUiState

### PreserveHandwrittenSources

### PreserveSdlSourceMap

### PreserveUiRegions

### PreserveUiSourceMap

### PreserveViewAnchors

### ProjectClassViews

### ProjectSdlViewpoints

### ProjectSelectedView

### ProjectUiGeneration

### PublishDomainUpdates

### PublishModelGeneration

### PublishPresentation

### PublishViewBundle

### ReadBoundedSources

### ReconcileWidgets

### RegisterDomainFunctions

### RejectStaleUiEvent

### RejectStaleViewResults

### ReleaseNativeWidgets

### ReleaseVisualResources

### ReportBindingDiagnostics

### ReportSourceDiagnostics

### ResetIncompatibleUiState

### ResolveAncestorDimensions

### ResolveCallbackSymbols

### ResolveConfiguredViewer

### ResolveSdlSymbols

### ResolveUiNames

### RestartChangedGoProgram

### RetainVisibleViewBundle

### RetirePreviousGeneration

### RetireReplacedPythonEntryPoints

### RevokeWidgetGenerations

### RouteDocumentToPane

### RouteDomainBindings

### ScheduleUiPublication

### SelectRelationshipViews

### ServeViewRequests

### SnapshotDomainState

### SnapshotUiState

### TokenizeSdlSource

### TokenizeSduiSource

### TraceViewpointFacts

### TrackInputDraft

### ValidateActionInput

### ValidateClassRelations

### ValidatePropertyBatch

### ValidateRelativeFormatting

### ValidateSdlProfile

### ValidateSdlStructure

### ValidateSymbolicBindings

### ValidateUiEvent

### ValidateViewRequest

### ValidateVisualResources

### ValidateWidgetArguments

### VerifyDiagramCapabilities

### VerifyNativeBehaviorParity

### WriteGeneratedArtifacts


## interface (30)

### BindingReloadPort

### BuildToolPort

### ContentProviderPort

### DiagnosticPort

### DiagramEnginePort

### DiagramPort

### DomainBindingPort

### DomainFunctionPort

### DomainStatePort

### ExecutionProfilePort

### ExportSinkPort

### FileChangePort

### GeneratedArtifactPort

### MeasurementPort

### PreparedFramePort

### ReloadPort

### ResourcePort

### SdlExecutionPort

### SdlFrontendPort

### SdlModelPort

### SdlReloadPort

### SduiFrontendPort

### SduiModelPort

### SourceInputPort

### SourceSnapshotPort

### UiReloadPort

### UiSessionPort

### UiSnapshotPort

### UiStatePort

### WidgetBackendPort


## message (48)

### BoundActionRequest

### BoundActionResult

### BuildGoRequest

### BuildGoResult

### BuildUiAstRequest

### BuildUiAstResult

### CompileSdlRequest

### CompileSdlResult

### CompileUiRejected

### CompileUiRequest

### CompileUiResult

### DisplayViewRequest

### DomainActionRequest

### DomainActionResult

### ExportSvgRequest

### ExportSvgResult

### GenerateGoRequest

### GenerateGoResult

### LaunchViewRequest

### LayoutRequest

### LayoutResult

### NormalizeUiRequest

### NormalizeUiResult

### PresentFrameRequest

### PresentFrameResult

### ProjectViewRequest

### PublishViewRequest

### ReloadPublished

### ReloadRejected

### ReloadRequest

### SdlActionRequest

### SdlActionResult

### SelectViewRequest

### TokenizeUiRequest

### TokenizeUiResult

### UiActionRejected

### UiActionRequest

### UiActionResult

### ValidateUiRequest

### ValidateUiResult

### ViewBundleResult

### ViewDisplayResult

### ViewLaunchResult

### ViewOpenedResult

### ViewProjectionRejected

### ViewReferenceResult

### ViewRequestRejected

### ViewTargetUnavailable


## mode (9)

### BoundExecution

### BoundLiveEditing

### DocumentBrowsing

### LiveEditing

### NativeBuild

### RichDocument

### SourceInspection

### StaticExport

### UiPreview


## scenario (13)

### BoundActionAccepted

### BoundActionRejected

### InteractiveFramePrepared

### InvalidViewSelectionRejected

### NativeProgramBuilt

### SdlModelReloadAccepted

### SelectedViewOpened

### StaticFrameExported

### UiCompilationAccepted

### UiModelReloadAccepted

### UiModelReloadRejected

### UnboundLocalAction

### ViewProjectionFailed


## unit (43)

### ContentServices

### DevelopmentTools

### DiagnosticReporter

### DiagramProvider

### DocumentBroker

### DomainStateMigrator

### FyneBackend

### GoBuildRunner

### GoCodeGenerator

### GoDomainImplementation

### MarkdownProvider

### ReloadCoordinator

### ResourceStore

### SdlDispatcher

### SdlExecutionGate

### SdlFrontend

### SdlFunctionRegistry

### SdlLexer

### SdlLibrary

### SdlNormalizer

### SdlParser

### SdlRuntime

### SdlStateStore

### SdlUiBindingAdapter

### SdlValidator

### SdlViewpointGenerator

### SduiDispatcher

### SduiFrontend

### SduiInstanceStore

### SduiLayout

### SduiLexer

### SduiLibrary

### SduiNormalizer

### SduiParser

### SduiPresentation

### SduiPropertyStore

### SduiRuntime

### SduiValidator

### SourceLoader

### SourceWatcher

### UiStateReconciler

### ViewArtifactStore

### ViewerLaunchAdapter


## usecase (7)

### BrowseDesignViews

### BuildNativeProduct

### EditRunningPrototype

### InspectModels

### PrototypeUserInterface

### PublishDesignDocumentation

### TryDomainInteraction


## variant (1)

### UiGenerationChanged
