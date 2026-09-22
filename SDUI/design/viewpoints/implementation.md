# SDL — generert implementasjonsplan

Denne rapporten er generert av SDL-verktøyet fra den validerte designkilden.
Status er modellens påstand. Planlagt Go-/UI-arbeid er ikke implementert av denne rapporten.

[SDL-kilde](../architecture.design) · [Alle viewpoints](printout.md)

116 Functionality-er har eksplisitt milepæl-/ansvarskobling; 0 mangler.

## G1FrontendPort

Status: **planned**. Kilde: f0277.

Leveransebidrag: StructuralModelInspection (f0276).

![Aktivitetsinndeling: G1FrontendPort](diagrams/VP06-detail-G1FrontendPort.svg)

### G1M1ParserAndAst

Status: **planned**. Kilde: f0282.

Forutsetninger: Ingen eksplisitt deklarert.

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| BuildSduiAst | SduiParser | f0278, f1017 |
| IdentifySourceRevision | SourceLoader | f0279, f1085 |
| ReadBoundedSources | SourceLoader | f0280, f1086 |
| TokenizeSduiSource | SduiLexer | f0281, f0999 |

Eksempelbane: UiCompilationAccepted. Kobling: f1160.

![Scenario: UiCompilationAccepted — modus SourceInspection](diagrams/VP08-UiCompilationAccepted.svg)

### G1M2ValidationAndNormalization

Status: **planned**. Kilde: f0293.

Forutsetninger: G1M1ParserAndAst (f0292).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| CoordinateSduiCompilation | SduiFrontend | f0284, f0962 |
| ExpandUiDefinitions | SduiNormalizer | f0285, f1011 |
| PreserveUiRegions | SduiNormalizer | f0286, f1012 |
| PreserveUiSourceMap | SduiNormalizer | f0287, f1013 |
| ResolveUiNames | SduiValidator | f0288, f1043 |
| ValidateRelativeFormatting | SduiValidator | f0289, f1044 |
| ValidateSymbolicBindings | SduiValidator | f0290, f1045 |
| ValidateWidgetArguments | SduiValidator | f0291, f1046 |

Eksempelbane: UiCompilationAccepted. Kobling: f1161.

[Scenariofigur](diagrams/VP08-UiCompilationAccepted.svg)

### G1M3Concept1AndDumps

Status: **planned**. Kilde: f0298.

Forutsetninger: G1M2ValidationAndNormalization (f0297).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| ExportConsoleSnapshot | SduiPresentation | f0295, f1023 |
| ReportSourceDiagnostics | DiagnosticReporter | f0296, f0179 |

## G2LayoutAndPresentation

Status: **planned**. Kilde: f0301.

Leveransebidrag: InteractiveUiPreview (f0300).

![Aktivitetsinndeling: G2LayoutAndPresentation](diagrams/VP06-detail-G2LayoutAndPresentation.svg)

### G2M1RelativeMeasurement

Status: **planned**. Kilde: f0307.

Forutsetninger: G1M2ValidationAndNormalization (f0306).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| AllocateGeometry | SduiLayout | f0302, f0989 |
| ComputeClipping | SduiLayout | f0303, f0991 |
| MeasureUiContent | SduiLayout | f0304, f0992 |
| ResolveAncestorDimensions | SduiLayout | f0305, f0993 |

Eksempelbane: InteractiveFramePrepared. Kobling: f0543.

![Scenario: InteractiveFramePrepared — modus UiPreview](diagrams/VP08-InteractiveFramePrepared.svg)

### G2M2SharedSvgGeometry

Status: **planned**. Kilde: f0312.

Forutsetninger: G2M1RelativeMeasurement (f0311).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| BuildPreparedFrame | SduiLayout | f0309, f0990 |
| ExportSvgSnapshot | SduiPresentation | f0310, f1024 |

Eksempelbane: StaticFrameExported. Kobling: f1102.

![Scenario: StaticFrameExported — modus StaticExport](diagrams/VP08-StaticFrameExported.svg)

### G2M3FyneInteractions

Status: **planned**. Kilde: f0321.

Forutsetninger: G2M2SharedSvgGeometry (f0320).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| ComposeInteractiveSession | FyneHost | f0314, f0267 |
| HandleFocusAndTextInput | FyneBackend | f0315, f0248 |
| PublishPresentation | FyneBackend | f0316, f0249 |
| ReconcileWidgets | FyneBackend | f0317, f0250 |
| ReleaseNativeWidgets | FyneBackend | f0318, f0251 |
| ScheduleUiPublication | FyneHost | f0319, f0268 |

Eksempelbane: InteractiveFramePrepared. Kobling: f0544.

[Scenariofigur](diagrams/VP08-InteractiveFramePrepared.svg)

Eksempelbane: UnboundLocalAction. Kobling: f1238.

![Scenario: UnboundLocalAction — modus UiPreview](diagrams/VP08-UnboundLocalAction.svg)

### G2M4RichContent

Status: **planned**. Kilde: f0329.

Forutsetninger: G2M3FyneInteractions (f0328).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| MeasureMarkdownContent | MarkdownProvider | f0323, f0595 |
| PrepareDiagramResource | DiagramProvider | f0324, f0183 |
| PrepareMarkdown | MarkdownProvider | f0325, f0596 |
| ReleaseVisualResources | ResourceStore | f0326, f0815 |
| ValidateVisualResources | ResourceStore | f0327, f0816 |

## G3UiRuntimeAndReload

Status: **planned**. Kilde: f0364.

Leveransebidrag: LiveModelReload (f0363).

![Aktivitetsinndeling: G3UiRuntimeAndReload](diagrams/VP06-detail-G3UiRuntimeAndReload.svg)

### G3M1TypedUiSession

Status: **planned**. Kilde: f0346.

Forutsetninger: G1M2ValidationAndNormalization (f0344), G2M3FyneInteractions (f0345).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| ApplyPropertyBatch | SduiPropertyStore | f0331, f1029 |
| CloseUiInstance | SduiRuntime | f0332, f1039 |
| CorrelateUiResult | SduiDispatcher | f0333, f0946 |
| CreateUiInstance | SduiRuntime | f0334, f1040 |
| DispatchUiEvent | SduiDispatcher | f0335, f0947 |
| ManageWidgetIdentities | SduiInstanceStore | f0336, f0979 |
| ProjectUiGeneration | SduiInstanceStore | f0337, f0980 |
| RejectStaleUiEvent | SduiDispatcher | f0338, f0948 |
| RevokeWidgetGenerations | SduiInstanceStore | f0339, f0981 |
| SnapshotUiState | SduiInstanceStore | f0340, f0982 |
| TrackInputDraft | SduiPropertyStore | f0341, f1030 |
| ValidatePropertyBatch | SduiPropertyStore | f0342, f1031 |
| ValidateUiEvent | SduiDispatcher | f0343, f0949 |

### G3M2CandidatePublication

Status: **planned**. Kilde: f0355.

Forutsetninger: G3M1TypedUiSession (f0354).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| CoalesceSourceChanges | SourceWatcher | f0348, f1089 |
| KeepLastValidModels | ReloadCoordinator | f0349, f0760 |
| ObserveSourceChanges | SourceWatcher | f0350, f1090 |
| PrepareCandidateModels | ReloadCoordinator | f0351, f0761 |
| PublishModelGeneration | ReloadCoordinator | f0352, f0762 |
| RetirePreviousGeneration | ReloadCoordinator | f0353, f0763 |

Eksempelbane: UiModelReloadAccepted. Kobling: f1198.

![Scenario: UiModelReloadAccepted — modus LiveEditing](diagrams/VP08-UiModelReloadAccepted.svg)

Eksempelbane: UiModelReloadRejected. Kobling: f1208.

![Scenario: UiModelReloadRejected — modus LiveEditing](diagrams/VP08-UiModelReloadRejected.svg)

### G3M3CompatibleState

Status: **planned**. Kilde: f0361.

Forutsetninger: G3M2CandidatePublication (f0360).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| MatchCompatibleWidgets | UiStateReconciler | f0357, f1224 |
| PreserveCompatibleUiState | UiStateReconciler | f0358, f1225 |
| ResetIncompatibleUiState | UiStateReconciler | f0359, f1226 |

Eksempelbane: UiModelReloadAccepted. Kobling: f1199.

[Scenariofigur](diagrams/VP08-UiModelReloadAccepted.svg)

## G4SdlRuntimeAndBinding

Status: **planned**. Kilde: f0408.

Leveransebidrag: TypedDomainBinding (f0407).

![Aktivitetsinndeling: G4SdlRuntimeAndBinding](diagrams/VP06-detail-G4SdlRuntimeAndBinding.svg)

### G4M1SdlFrontend

Status: **planned**. Kilde: f0373.

Forutsetninger: Ingen eksplisitt deklarert.

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| BuildSdlAst | SdlParser | f0365, f0894 |
| CoordinateSdlCompilation | SdlFrontend | f0366, f0867 |
| NormalizeSdlModel | SdlNormalizer | f0367, f0891 |
| PreserveSdlSourceMap | SdlNormalizer | f0368, f0892 |
| ResolveSdlSymbols | SdlValidator | f0369, f0923 |
| TokenizeSdlSource | SdlLexer | f0370, f0875 |
| ValidateSdlProfile | SdlValidator | f0371, f0925 |
| ValidateSdlStructure | SdlValidator | f0372, f0926 |

Eksempelbane: SdlModelReloadAccepted. Kobling: f0884.

![Scenario: SdlModelReloadAccepted — modus BoundLiveEditing](diagrams/VP08-SdlModelReloadAccepted.svg)

### G4M2TypedExecution

Status: **planned**. Kilde: f0388.

Forutsetninger: G4M1SdlFrontend (f0387).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| CancelPendingActions | SdlDispatcher | f0375, f0848 |
| CheckExecutionCompleteness | SdlExecutionGate | f0376, f0859 |
| CheckFunctionSignatures | SdlFunctionRegistry | f0377, f0872 |
| CloseSdlInstance | SdlRuntime | f0378, f0903 |
| CorrelateActionResult | SdlDispatcher | f0379, f0849 |
| CreateSdlInstance | SdlRuntime | f0380, f0904 |
| InvokeRegisteredFunction | SdlDispatcher | f0381, f0850 |
| ManageDomainState | SdlStateStore | f0382, f0907 |
| PerformDomainOperation | GoDomainImplementation | f0383, f0523 |
| RegisterDomainFunctions | SdlFunctionRegistry | f0384, f0873 |
| SnapshotDomainState | SdlStateStore | f0385, f0908 |
| ValidateActionInput | SdlDispatcher | f0386, f0851 |

### G4M3UiDomainBinding

Status: **planned**. Kilde: f0398.

Forutsetninger: G3M1TypedUiSession (f0396), G4M2TypedExecution (f0397).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| ConnectTypedWidgetHandles | SdlUiBindingAdapter | f0390, f0913 |
| DisconnectBindings | SdlUiBindingAdapter | f0391, f0914 |
| PublishDomainUpdates | SdlUiBindingAdapter | f0392, f0915 |
| ReportBindingDiagnostics | DiagnosticReporter | f0393, f0178 |
| ResolveCallbackSymbols | SdlUiBindingAdapter | f0394, f0916 |
| RouteDomainBindings | SdlUiBindingAdapter | f0395, f0917 |

Eksempelbane: BoundActionAccepted. Kobling: f0030.

![Scenario: BoundActionAccepted — modus BoundExecution](diagrams/VP08-BoundActionAccepted.svg)

Eksempelbane: BoundActionRejected. Kobling: f0043.

![Scenario: BoundActionRejected — modus BoundExecution](diagrams/VP08-BoundActionRejected.svg)

### G4M4DomainReload

Status: **planned**. Kilde: f0405.

Forutsetninger: G3M3CompatibleState (f0403), G4M3UiDomainBinding (f0404).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| CheckDomainStateCompatibility | DomainStateMigrator | f0400, f0214 |
| MigrateOrResetDomainState | DomainStateMigrator | f0401, f0215 |
| RestartChangedGoProgram | GoBuildRunner | f0402, f0506 |

Eksempelbane: SdlModelReloadAccepted. Kobling: f0885.

[Scenariofigur](diagrams/VP08-SdlModelReloadAccepted.svg)

## G5NativeGeneration

Status: **planned**. Kilde: f0435.

Leveransebidrag: DesignDocumentation (f0433), NativeGoAssembly (f0434).

![Aktivitetsinndeling: G5NativeGeneration](diagrams/VP06-detail-G5NativeGeneration.svg)

### G5M1GeneratedGo

Status: **planned**. Kilde: f0414.

Forutsetninger: G4M4DomainReload (f0413).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| BuildGeneratedApplication | GoBuildRunner | f0409, f0505 |
| GenerateBindingRegistration | GoCodeGenerator | f0410, f0513 |
| GenerateModelConstructors | GoCodeGenerator | f0411, f0514 |
| PreserveHandwrittenSources | GoCodeGenerator | f0412, f0515 |

Eksempelbane: NativeProgramBuilt. Kobling: f0624.

![Scenario: NativeProgramBuilt — modus NativeBuild](diagrams/VP08-NativeProgramBuilt.svg)

### G5M2BehaviorParity

Status: **planned**. Kilde: f0418.

Forutsetninger: G5M1GeneratedGo (f0417).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| VerifyNativeBehaviorParity | DevelopmentTools | f0416, f0176 |

### G5M3DocumentationExport

Status: **planned**. Kilde: f0426.

Forutsetninger: G2M4RichContent (f0423), G5M2BehaviorParity (f0424), G6M1StaticNavigation (f0425).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| ComposeHeadlessExport | CommandLineHost | f0420, f0094 |
| ComposeMarkdownDocument | SduiPresentation | f0421, f1022 |
| WriteGeneratedArtifacts | CommandLineHost | f0422, f0095 |

Eksempelbane: StaticFrameExported. Kobling: f1103.

[Scenariofigur](diagrams/VP08-StaticFrameExported.svg)

### G5M4RetirePython

Status: **planned**. Kilde: f0431.

Forutsetninger: G1M3Concept1AndDumps (f0429), G5M3DocumentationExport (f0430).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| RetireReplacedPythonEntryPoints | DevelopmentTools | f0428, f0175 |

## G6NavigableDocumentation

Status: **planned**. Kilde: f0481.

Leveransebidrag: NavigableDesignDocumentation (f0479), TypedDesignInspection (f0480).

![Aktivitetsinndeling: G6NavigableDocumentation](diagrams/VP06-detail-G6NavigableDocumentation.svg)

### G6M1StaticNavigation

Status: **planned**. Kilde: f0445.

Forutsetninger: G4M1SdlFrontend (f0444).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| ComposeViewPackage | SdlViewpointGenerator | f0436, f0929 |
| ExportModelInventories | SdlViewpointGenerator | f0437, f0930 |
| ExportViewpointMarkdown | SdlViewpointGenerator | f0438, f0931 |
| GenerateViewNavigation | SdlViewpointGenerator | f0439, f0932 |
| IndexViewpointLevels | SdlViewpointGenerator | f0440, f0933 |
| PreserveViewAnchors | SdlViewpointGenerator | f0441, f0934 |
| ProjectSdlViewpoints | SdlViewpointGenerator | f0442, f0936 |
| TraceViewpointFacts | SdlViewpointGenerator | f0443, f0939 |

Eksempelbane: SelectedViewOpened. Kobling: f1056.

![Scenario: SelectedViewOpened — modus DocumentBrowsing](diagrams/VP08-SelectedViewOpened.svg)

### G6M2OnDemandViews

Status: **planned**. Kilde: f0453.

Forutsetninger: G6M1StaticNavigation (f0452).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| KeyViewRevision | DocumentBroker | f0447, f0193 |
| ProjectSelectedView | SdlViewpointGenerator | f0448, f0937 |
| PublishViewBundle | ViewArtifactStore | f0449, f1282 |
| SelectRelationshipViews | SdlViewpointGenerator | f0450, f0938 |
| ValidateViewRequest | DocumentBroker | f0451, f0196 |

Eksempelbane: InvalidViewSelectionRejected. Kobling: f0553.

![Scenario: InvalidViewSelectionRejected — modus DocumentBrowsing](diagrams/VP08-InvalidViewSelectionRejected.svg)

Eksempelbane: SelectedViewOpened. Kobling: f1057.

[Scenariofigur](diagrams/VP08-SelectedViewOpened.svg)

Eksempelbane: ViewProjectionFailed. Kobling: f1357.

![Scenario: ViewProjectionFailed — modus DocumentBrowsing](diagrams/VP08-ViewProjectionFailed.svg)

### G6M3XfmdNavigation

Status: **planned**. Kilde: f0460.

Forutsetninger: G6M2OnDemandViews (f0459).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| CaptureNavigationTarget | XfmdDocumentHost | f0455, f1416 |
| DispatchViewOpen | ViewerLaunchAdapter | f0456, f1406 |
| ResolveConfiguredViewer | ViewerLaunchAdapter | f0457, f1407 |
| RouteDocumentToPane | XfmdDocumentHost | f0458, f1417 |

Eksempelbane: SelectedViewOpened. Kobling: f1058.

[Scenariofigur](diagrams/VP08-SelectedViewOpened.svg)

### G6M4SessionPublication

Status: **planned**. Kilde: f0467.

Forutsetninger: G6M3XfmdNavigation (f0466).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| EvictUnusedViewBundles | ViewArtifactStore | f0462, f1281 |
| RejectStaleViewResults | DocumentBroker | f0463, f0194 |
| RetainVisibleViewBundle | ViewArtifactStore | f0464, f1283 |
| ServeViewRequests | DocumentBroker | f0465, f0195 |

Eksempelbane: InvalidViewSelectionRejected. Kobling: f0554.

[Scenariofigur](diagrams/VP08-InvalidViewSelectionRejected.svg)

Eksempelbane: SelectedViewOpened. Kobling: f1059.

[Scenariofigur](diagrams/VP08-SelectedViewOpened.svg)

Eksempelbane: ViewProjectionFailed. Kobling: f1358.

[Scenariofigur](diagrams/VP08-ViewProjectionFailed.svg)

### G6M5SemanticNotation

Status: **planned**. Kilde: f0472.

Forutsetninger: G6M1StaticNavigation (f0471).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| ApplySemanticNotation | SdlViewpointGenerator | f0469, f0928 |
| VerifyDiagramCapabilities | SdlViewpointGenerator | f0470, f0940 |

### G6M6ClassViews

Status: **planned**. Kilde: f0477.

Forutsetninger: G6M5SemanticNotation (f0476).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| ProjectClassViews | SdlViewpointGenerator | f0474, f0935 |
| ValidateClassRelations | SdlValidator | f0475, f0924 |

## Udekket modellansvar

Ingen deklarerte Functionality-er mangler addresses-kobling. Dette beviser ikke full kravdekning.

Runtime-semantikk, full AST-/ABI-schema og fysisk layoutmåling må fortsatt realiseres og testes i implementasjonsfasene.
