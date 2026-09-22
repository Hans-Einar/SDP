# SDL — generert implementasjonsplan

Denne rapporten er generert av SDL-verktøyet fra den validerte designkilden.
Status er modellens påstand. Planlagt Go-/UI-arbeid er ikke implementert av denne rapporten.

[SDL-kilde](../architecture.design) · [Alle viewpoints](printout.md)

109 Functionality-er har eksplisitt milepæl-/ansvarskobling; 0 mangler.

## G1FrontendPort

Status: **planned**. Kilde: f0273.

Leveransebidrag: StructuralModelInspection (f0272).

![Aktivitetsinndeling: G1FrontendPort](diagrams/VP06-detail-G1FrontendPort.svg)

### G1M1ParserAndAst

Status: **planned**. Kilde: f0278.

Forutsetninger: Ingen eksplisitt deklarert.

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| BuildSduiAst | SduiParser | f0274, f0988 |
| IdentifySourceRevision | SourceLoader | f0275, f1054 |
| ReadBoundedSources | SourceLoader | f0276, f1055 |
| TokenizeSduiSource | SduiLexer | f0277, f0970 |

Eksempelbane: UiCompilationAccepted. Kobling: f1127.

![Scenario: UiCompilationAccepted — modus SourceInspection](diagrams/VP08-UiCompilationAccepted.svg)

### G1M2ValidationAndNormalization

Status: **planned**. Kilde: f0289.

Forutsetninger: G1M1ParserAndAst (f0288).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| CoordinateSduiCompilation | SduiFrontend | f0280, f0933 |
| ExpandUiDefinitions | SduiNormalizer | f0281, f0982 |
| PreserveUiRegions | SduiNormalizer | f0282, f0983 |
| PreserveUiSourceMap | SduiNormalizer | f0283, f0984 |
| ResolveUiNames | SduiValidator | f0284, f1014 |
| ValidateRelativeFormatting | SduiValidator | f0285, f1015 |
| ValidateSymbolicBindings | SduiValidator | f0286, f1016 |
| ValidateWidgetArguments | SduiValidator | f0287, f1017 |

Eksempelbane: UiCompilationAccepted. Kobling: f1128.

[Scenariofigur](diagrams/VP08-UiCompilationAccepted.svg)

### G1M3Concept1AndDumps

Status: **planned**. Kilde: f0294.

Forutsetninger: G1M2ValidationAndNormalization (f0293).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| ExportConsoleSnapshot | SduiPresentation | f0291, f0994 |
| ReportSourceDiagnostics | DiagnosticReporter | f0292, f0177 |

## G2LayoutAndPresentation

Status: **planned**. Kilde: f0297.

Leveransebidrag: InteractiveUiPreview (f0296).

![Aktivitetsinndeling: G2LayoutAndPresentation](diagrams/VP06-detail-G2LayoutAndPresentation.svg)

### G2M1RelativeMeasurement

Status: **planned**. Kilde: f0303.

Forutsetninger: G1M2ValidationAndNormalization (f0302).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| AllocateGeometry | SduiLayout | f0298, f0960 |
| ComputeClipping | SduiLayout | f0299, f0962 |
| MeasureUiContent | SduiLayout | f0300, f0963 |
| ResolveAncestorDimensions | SduiLayout | f0301, f0964 |

Eksempelbane: InteractiveFramePrepared. Kobling: f0523.

![Scenario: InteractiveFramePrepared — modus UiPreview](diagrams/VP08-InteractiveFramePrepared.svg)

### G2M2SharedSvgGeometry

Status: **planned**. Kilde: f0308.

Forutsetninger: G2M1RelativeMeasurement (f0307).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| BuildPreparedFrame | SduiLayout | f0305, f0961 |
| ExportSvgSnapshot | SduiPresentation | f0306, f0995 |

Eksempelbane: StaticFrameExported. Kobling: f1071.

![Scenario: StaticFrameExported — modus StaticExport](diagrams/VP08-StaticFrameExported.svg)

### G2M3FyneInteractions

Status: **planned**. Kilde: f0317.

Forutsetninger: G2M2SharedSvgGeometry (f0316).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| ComposeInteractiveSession | FyneHost | f0310, f0263 |
| HandleFocusAndTextInput | FyneBackend | f0311, f0244 |
| PublishPresentation | FyneBackend | f0312, f0245 |
| ReconcileWidgets | FyneBackend | f0313, f0246 |
| ReleaseNativeWidgets | FyneBackend | f0314, f0247 |
| ScheduleUiPublication | FyneHost | f0315, f0264 |

Eksempelbane: InteractiveFramePrepared. Kobling: f0524.

[Scenariofigur](diagrams/VP08-InteractiveFramePrepared.svg)

Eksempelbane: UnboundLocalAction. Kobling: f1205.

![Scenario: UnboundLocalAction — modus UiPreview](diagrams/VP08-UnboundLocalAction.svg)

### G2M4RichContent

Status: **planned**. Kilde: f0325.

Forutsetninger: G2M3FyneInteractions (f0324).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| MeasureMarkdownContent | MarkdownProvider | f0319, f0575 |
| PrepareDiagramResource | DiagramProvider | f0320, f0181 |
| PrepareMarkdown | MarkdownProvider | f0321, f0576 |
| ReleaseVisualResources | ResourceStore | f0322, f0793 |
| ValidateVisualResources | ResourceStore | f0323, f0794 |

## G3UiRuntimeAndReload

Status: **planned**. Kilde: f0360.

Leveransebidrag: LiveModelReload (f0359).

![Aktivitetsinndeling: G3UiRuntimeAndReload](diagrams/VP06-detail-G3UiRuntimeAndReload.svg)

### G3M1TypedUiSession

Status: **planned**. Kilde: f0342.

Forutsetninger: G1M2ValidationAndNormalization (f0340), G2M3FyneInteractions (f0341).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| ApplyPropertyBatch | SduiPropertyStore | f0327, f1000 |
| CloseUiInstance | SduiRuntime | f0328, f1010 |
| CorrelateUiResult | SduiDispatcher | f0329, f0917 |
| CreateUiInstance | SduiRuntime | f0330, f1011 |
| DispatchUiEvent | SduiDispatcher | f0331, f0918 |
| ManageWidgetIdentities | SduiInstanceStore | f0332, f0950 |
| ProjectUiGeneration | SduiInstanceStore | f0333, f0951 |
| RejectStaleUiEvent | SduiDispatcher | f0334, f0919 |
| RevokeWidgetGenerations | SduiInstanceStore | f0335, f0952 |
| SnapshotUiState | SduiInstanceStore | f0336, f0953 |
| TrackInputDraft | SduiPropertyStore | f0337, f1001 |
| ValidatePropertyBatch | SduiPropertyStore | f0338, f1002 |
| ValidateUiEvent | SduiDispatcher | f0339, f0920 |

### G3M2CandidatePublication

Status: **planned**. Kilde: f0351.

Forutsetninger: G3M1TypedUiSession (f0350).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| CoalesceSourceChanges | SourceWatcher | f0344, f1058 |
| KeepLastValidModels | ReloadCoordinator | f0345, f0738 |
| ObserveSourceChanges | SourceWatcher | f0346, f1059 |
| PrepareCandidateModels | ReloadCoordinator | f0347, f0739 |
| PublishModelGeneration | ReloadCoordinator | f0348, f0740 |
| RetirePreviousGeneration | ReloadCoordinator | f0349, f0741 |

Eksempelbane: UiModelReloadAccepted. Kobling: f1165.

![Scenario: UiModelReloadAccepted — modus LiveEditing](diagrams/VP08-UiModelReloadAccepted.svg)

Eksempelbane: UiModelReloadRejected. Kobling: f1175.

![Scenario: UiModelReloadRejected — modus LiveEditing](diagrams/VP08-UiModelReloadRejected.svg)

### G3M3CompatibleState

Status: **planned**. Kilde: f0357.

Forutsetninger: G3M2CandidatePublication (f0356).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| MatchCompatibleWidgets | UiStateReconciler | f0353, f1191 |
| PreserveCompatibleUiState | UiStateReconciler | f0354, f1192 |
| ResetIncompatibleUiState | UiStateReconciler | f0355, f1193 |

Eksempelbane: UiModelReloadAccepted. Kobling: f1166.

[Scenariofigur](diagrams/VP08-UiModelReloadAccepted.svg)

## G4SdlRuntimeAndBinding

Status: **planned**. Kilde: f0404.

Leveransebidrag: TypedDomainBinding (f0403).

![Aktivitetsinndeling: G4SdlRuntimeAndBinding](diagrams/VP06-detail-G4SdlRuntimeAndBinding.svg)

### G4M1SdlFrontend

Status: **planned**. Kilde: f0369.

Forutsetninger: Ingen eksplisitt deklarert.

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| BuildSdlAst | SdlParser | f0361, f0872 |
| CoordinateSdlCompilation | SdlFrontend | f0362, f0845 |
| NormalizeSdlModel | SdlNormalizer | f0363, f0869 |
| PreserveSdlSourceMap | SdlNormalizer | f0364, f0870 |
| ResolveSdlSymbols | SdlValidator | f0365, f0901 |
| TokenizeSdlSource | SdlLexer | f0366, f0853 |
| ValidateSdlProfile | SdlValidator | f0367, f0902 |
| ValidateSdlStructure | SdlValidator | f0368, f0903 |

Eksempelbane: SdlModelReloadAccepted. Kobling: f0862.

![Scenario: SdlModelReloadAccepted — modus BoundLiveEditing](diagrams/VP08-SdlModelReloadAccepted.svg)

### G4M2TypedExecution

Status: **planned**. Kilde: f0384.

Forutsetninger: G4M1SdlFrontend (f0383).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| CancelPendingActions | SdlDispatcher | f0371, f0826 |
| CheckExecutionCompleteness | SdlExecutionGate | f0372, f0837 |
| CheckFunctionSignatures | SdlFunctionRegistry | f0373, f0850 |
| CloseSdlInstance | SdlRuntime | f0374, f0881 |
| CorrelateActionResult | SdlDispatcher | f0375, f0827 |
| CreateSdlInstance | SdlRuntime | f0376, f0882 |
| InvokeRegisteredFunction | SdlDispatcher | f0377, f0828 |
| ManageDomainState | SdlStateStore | f0378, f0885 |
| PerformDomainOperation | GoDomainImplementation | f0379, f0505 |
| RegisterDomainFunctions | SdlFunctionRegistry | f0380, f0851 |
| SnapshotDomainState | SdlStateStore | f0381, f0886 |
| ValidateActionInput | SdlDispatcher | f0382, f0829 |

### G4M3UiDomainBinding

Status: **planned**. Kilde: f0394.

Forutsetninger: G3M1TypedUiSession (f0392), G4M2TypedExecution (f0393).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| ConnectTypedWidgetHandles | SdlUiBindingAdapter | f0386, f0891 |
| DisconnectBindings | SdlUiBindingAdapter | f0387, f0892 |
| PublishDomainUpdates | SdlUiBindingAdapter | f0388, f0893 |
| ReportBindingDiagnostics | DiagnosticReporter | f0389, f0176 |
| ResolveCallbackSymbols | SdlUiBindingAdapter | f0390, f0894 |
| RouteDomainBindings | SdlUiBindingAdapter | f0391, f0895 |

Eksempelbane: BoundActionAccepted. Kobling: f0028.

![Scenario: BoundActionAccepted — modus BoundExecution](diagrams/VP08-BoundActionAccepted.svg)

Eksempelbane: BoundActionRejected. Kobling: f0041.

![Scenario: BoundActionRejected — modus BoundExecution](diagrams/VP08-BoundActionRejected.svg)

### G4M4DomainReload

Status: **planned**. Kilde: f0401.

Forutsetninger: G3M3CompatibleState (f0399), G4M3UiDomainBinding (f0400).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| CheckDomainStateCompatibility | DomainStateMigrator | f0396, f0212 |
| MigrateOrResetDomainState | DomainStateMigrator | f0397, f0213 |
| RestartChangedGoProgram | GoBuildRunner | f0398, f0488 |

Eksempelbane: SdlModelReloadAccepted. Kobling: f0863.

[Scenariofigur](diagrams/VP08-SdlModelReloadAccepted.svg)

## G5NativeGeneration

Status: **planned**. Kilde: f0431.

Leveransebidrag: DesignDocumentation (f0429), NativeGoAssembly (f0430).

![Aktivitetsinndeling: G5NativeGeneration](diagrams/VP06-detail-G5NativeGeneration.svg)

### G5M1GeneratedGo

Status: **planned**. Kilde: f0410.

Forutsetninger: G4M4DomainReload (f0409).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| BuildGeneratedApplication | GoBuildRunner | f0405, f0487 |
| GenerateBindingRegistration | GoCodeGenerator | f0406, f0495 |
| GenerateModelConstructors | GoCodeGenerator | f0407, f0496 |
| PreserveHandwrittenSources | GoCodeGenerator | f0408, f0497 |

Eksempelbane: NativeProgramBuilt. Kobling: f0604.

![Scenario: NativeProgramBuilt — modus NativeBuild](diagrams/VP08-NativeProgramBuilt.svg)

### G5M2BehaviorParity

Status: **planned**. Kilde: f0414.

Forutsetninger: G5M1GeneratedGo (f0413).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| VerifyNativeBehaviorParity | DevelopmentTools | f0412, f0174 |

### G5M3DocumentationExport

Status: **planned**. Kilde: f0422.

Forutsetninger: G2M4RichContent (f0419), G5M2BehaviorParity (f0420), G6M1StaticNavigation (f0421).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| ComposeHeadlessExport | CommandLineHost | f0416, f0092 |
| ComposeMarkdownDocument | SduiPresentation | f0417, f0993 |
| WriteGeneratedArtifacts | CommandLineHost | f0418, f0093 |

Eksempelbane: StaticFrameExported. Kobling: f1072.

[Scenariofigur](diagrams/VP08-StaticFrameExported.svg)

### G5M4RetirePython

Status: **planned**. Kilde: f0427.

Forutsetninger: G1M3Concept1AndDumps (f0425), G5M3DocumentationExport (f0426).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| RetireReplacedPythonEntryPoints | DevelopmentTools | f0424, f0173 |

## G6NavigableDocumentation

Status: **planned**. Kilde: f0463.

Leveransebidrag: NavigableDesignDocumentation (f0462).

![Aktivitetsinndeling: G6NavigableDocumentation](diagrams/VP06-detail-G6NavigableDocumentation.svg)

### G6M1StaticNavigation

Status: **planned**. Kilde: f0439.

Forutsetninger: G4M1SdlFrontend (f0438).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| ComposeViewPackage | SdlViewpointGenerator | f0432, f0905 |
| ExportViewpointMarkdown | SdlViewpointGenerator | f0433, f0906 |
| GenerateViewNavigation | SdlViewpointGenerator | f0434, f0907 |
| PreserveViewAnchors | SdlViewpointGenerator | f0435, f0908 |
| ProjectSdlViewpoints | SdlViewpointGenerator | f0436, f0909 |
| TraceViewpointFacts | SdlViewpointGenerator | f0437, f0911 |

Eksempelbane: SelectedViewOpened. Kobling: f1025.

![Scenario: SelectedViewOpened — modus DocumentBrowsing](diagrams/VP08-SelectedViewOpened.svg)

### G6M2OnDemandViews

Status: **planned**. Kilde: f0446.

Forutsetninger: G6M1StaticNavigation (f0445).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| KeyViewRevision | DocumentBroker | f0441, f0191 |
| ProjectSelectedView | SdlViewpointGenerator | f0442, f0910 |
| PublishViewBundle | ViewArtifactStore | f0443, f1245 |
| ValidateViewRequest | DocumentBroker | f0444, f0194 |

Eksempelbane: InvalidViewSelectionRejected. Kobling: f0533.

![Scenario: InvalidViewSelectionRejected — modus DocumentBrowsing](diagrams/VP08-InvalidViewSelectionRejected.svg)

Eksempelbane: SelectedViewOpened. Kobling: f1026.

[Scenariofigur](diagrams/VP08-SelectedViewOpened.svg)

Eksempelbane: ViewProjectionFailed. Kobling: f1320.

![Scenario: ViewProjectionFailed — modus DocumentBrowsing](diagrams/VP08-ViewProjectionFailed.svg)

### G6M3XfmdNavigation

Status: **planned**. Kilde: f0453.

Forutsetninger: G6M2OnDemandViews (f0452).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| CaptureNavigationTarget | XfmdDocumentHost | f0448, f1379 |
| DispatchViewOpen | ViewerLaunchAdapter | f0449, f1369 |
| ResolveConfiguredViewer | ViewerLaunchAdapter | f0450, f1370 |
| RouteDocumentToPane | XfmdDocumentHost | f0451, f1380 |

Eksempelbane: SelectedViewOpened. Kobling: f1027.

[Scenariofigur](diagrams/VP08-SelectedViewOpened.svg)

### G6M4SessionPublication

Status: **planned**. Kilde: f0460.

Forutsetninger: G6M3XfmdNavigation (f0459).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| EvictUnusedViewBundles | ViewArtifactStore | f0455, f1244 |
| RejectStaleViewResults | DocumentBroker | f0456, f0192 |
| RetainVisibleViewBundle | ViewArtifactStore | f0457, f1246 |
| ServeViewRequests | DocumentBroker | f0458, f0193 |

Eksempelbane: InvalidViewSelectionRejected. Kobling: f0534.

[Scenariofigur](diagrams/VP08-InvalidViewSelectionRejected.svg)

Eksempelbane: SelectedViewOpened. Kobling: f1028.

[Scenariofigur](diagrams/VP08-SelectedViewOpened.svg)

Eksempelbane: ViewProjectionFailed. Kobling: f1321.

[Scenariofigur](diagrams/VP08-ViewProjectionFailed.svg)

## Udekket modellansvar

Ingen deklarerte Functionality-er mangler addresses-kobling. Dette beviser ikke full kravdekning.

Runtime-semantikk, full AST-/ABI-schema og fysisk layoutmåling må fortsatt realiseres og testes i implementasjonsfasene.
