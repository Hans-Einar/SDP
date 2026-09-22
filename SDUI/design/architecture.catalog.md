# Parser/runtime — generert ansvarsoversikt

Generert fra [architecture.design](architecture.design) med eksisterende design-core-parser.
Dette er en visning av målstruktur, ikke implementasjonsbevis. Ikke rediger denne filen manuelt.

| Unit / Container | Underliggende Units | Eide Functionality-er | Brukte grensesnitt |
| --- | --- | --- | --- |
| CommandLineHost | — | ComposeHeadlessExport, WriteGeneratedArtifacts | ExportSinkPort, PreparedFramePort, SduiFrontendPort, SourceSnapshotPort |
| ContentServices | DiagramProvider, MarkdownProvider, ResourceStore | — | — |
| DevelopmentTools | DiagnosticReporter, GoBuildRunner, GoCodeGenerator, ReloadCoordinator, SourceLoader, SourceWatcher | RetireReplacedPythonEntryPoints, VerifyNativeBehaviorParity | — |
| DiagnosticReporter | — | ReportBindingDiagnostics, ReportSourceDiagnostics | — |
| DiagramProvider | — | PrepareDiagramResource | DiagramEnginePort, ResourcePort |
| DocumentBroker | — | KeyViewRevision, RejectStaleViewResults, ServeViewRequests, ValidateViewRequest | — |
| DomainStateMigrator | — | CheckDomainStateCompatibility, MigrateOrResetDomainState | DomainStatePort, SdlModelPort |
| FyneBackend | — | HandleFocusAndTextInput, PublishPresentation, ReconcileWidgets, ReleaseNativeWidgets | PreparedFramePort, UiSessionPort |
| FyneHost | — | ComposeInteractiveSession, ScheduleUiPublication | DomainBindingPort, ReloadPort, SdlFrontendPort, SduiFrontendPort, SourceSnapshotPort, UiSessionPort, WidgetBackendPort |
| GoBuildRunner | — | BuildGeneratedApplication, RestartChangedGoProgram | BuildToolPort, GeneratedArtifactPort |
| GoCodeGenerator | — | GenerateBindingRegistration, GenerateModelConstructors, PreserveHandwrittenSources | ExecutionProfilePort, SdlModelPort, SduiModelPort |
| GoDomainImplementation | — | PerformDomainOperation | — |
| MarkdownProvider | — | MeasureMarkdownContent, PrepareMarkdown | DiagramPort, MeasurementPort, ResourcePort |
| ReloadCoordinator | — | KeepLastValidModels, PrepareCandidateModels, PublishModelGeneration, RetirePreviousGeneration | BindingReloadPort, DiagnosticPort, SdlFrontendPort, SdlReloadPort, SduiFrontendPort, SourceSnapshotPort, UiReloadPort |
| ResourceStore | — | ReleaseVisualResources, ValidateVisualResources | — |
| SdlDispatcher | — | CancelPendingActions, CorrelateActionResult, InvokeRegisteredFunction, ValidateActionInput | DomainFunctionPort, DomainStatePort |
| SdlExecutionGate | — | CheckExecutionCompleteness | DiagnosticPort, SdlModelPort |
| SdlFrontend | SdlLexer, SdlNormalizer, SdlParser, SdlValidator | CoordinateSdlCompilation | DiagnosticPort, SourceSnapshotPort |
| SdlFunctionRegistry | — | CheckFunctionSignatures, RegisterDomainFunctions | DomainFunctionPort |
| SdlLexer | — | TokenizeSdlSource | — |
| SdlLibrary | SdlFrontend, SdlRuntime, SdlViewpointGenerator | — | — |
| SdlNormalizer | — | NormalizeSdlModel, PreserveSdlSourceMap | — |
| SdlParser | — | BuildSdlAst | — |
| SdlRuntime | DomainStateMigrator, SdlDispatcher, SdlExecutionGate, SdlFunctionRegistry, SdlStateStore | CloseSdlInstance, CreateSdlInstance | DomainFunctionPort, SdlModelPort |
| SdlStateStore | — | ManageDomainState, SnapshotDomainState | — |
| SdlUiBindingAdapter | — | ConnectTypedWidgetHandles, DisconnectBindings, PublishDomainUpdates, ResolveCallbackSymbols, RouteDomainBindings | DiagnosticPort, SdlExecutionPort, UiSessionPort |
| SdlValidator | — | ResolveSdlSymbols, ValidateSdlProfile, ValidateSdlStructure | — |
| SdlViewpointGenerator | — | ComposeViewPackage, ExportViewpointMarkdown, GenerateViewNavigation, PreserveViewAnchors, ProjectSdlViewpoints, ProjectSelectedView, TraceViewpointFacts | — |
| SduiDispatcher | — | CorrelateUiResult, DispatchUiEvent, RejectStaleUiEvent, ValidateUiEvent | DomainBindingPort, UiStatePort |
| SduiFrontend | SduiLexer, SduiNormalizer, SduiParser, SduiValidator | CoordinateSduiCompilation | DiagnosticPort, SourceSnapshotPort |
| SduiInstanceStore | — | ManageWidgetIdentities, ProjectUiGeneration, RevokeWidgetGenerations, SnapshotUiState | — |
| SduiLayout | — | AllocateGeometry, BuildPreparedFrame, ComputeClipping, MeasureUiContent, ResolveAncestorDimensions | ContentProviderPort, MeasurementPort, UiSnapshotPort |
| SduiLexer | — | TokenizeSduiSource | — |
| SduiLibrary | SduiFrontend, SduiLayout, SduiPresentation, SduiRuntime | — | — |
| SduiNormalizer | — | ExpandUiDefinitions, PreserveUiRegions, PreserveUiSourceMap | — |
| SduiParser | — | BuildSduiAst | — |
| SduiPresentation | — | ComposeMarkdownDocument, ExportConsoleSnapshot, ExportSvgSnapshot | PreparedFramePort |
| SduiPropertyStore | — | ApplyPropertyBatch, TrackInputDraft, ValidatePropertyBatch | UiStatePort |
| SduiRuntime | SduiDispatcher, SduiInstanceStore, SduiPropertyStore, UiStateReconciler | CloseUiInstance, CreateUiInstance | DomainBindingPort, SduiModelPort |
| SduiValidator | — | ResolveUiNames, ValidateRelativeFormatting, ValidateSymbolicBindings, ValidateWidgetArguments | — |
| SourceLoader | — | DesignSourceArchive, IdentifySourceRevision, ReadBoundedSources | SourceInputPort |
| SourceWatcher | — | CoalesceSourceChanges, ObserveSourceChanges | FileChangePort |
| UiStateReconciler | — | MatchCompatibleWidgets, PreserveCompatibleUiState, ResetIncompatibleUiState | SduiModelPort, UiStatePort |
| ViewArtifactStore | — | EvictUnusedViewBundles, PublishViewBundle, RetainVisibleViewBundle | — |
| ViewServiceHost | DocumentBroker, ViewArtifactStore, ViewerLaunchAdapter | — | — |
| ViewerLaunchAdapter | — | DispatchViewOpen, ResolveConfiguredViewer | — |
| XfmdDocumentHost | — | CaptureNavigationTarget, RouteDocumentToPane | — |
