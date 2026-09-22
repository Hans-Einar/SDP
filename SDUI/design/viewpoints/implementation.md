# SDL — generert implementasjonsplan

Denne rapporten er generert av SDL-verktøyet fra den validerte designkilden.
Status er modellens påstand. Planlagt Go-/UI-arbeid er ikke implementert av denne rapporten.

[SDL-kilde](../architecture.design) · [Alle viewpoints](printout.md)

94 Functionality-er har eksplisitt milepæl-/ansvarskobling; 0 mangler.

## G1FrontendPort

Status: **planned**. Kilde: f0244.

Leveransebidrag: StructuralModelInspection (f0243).

![Aktivitetsinndeling: G1FrontendPort](diagrams/VP06-detail-G1FrontendPort.svg)

### G1M1ParserAndAst

Status: **planned**. Kilde: f0249.

Forutsetninger: Ingen eksplisitt deklarert.

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| BuildSduiAst | SduiParser | f0245, f0871 |
| IdentifySourceRevision | SourceLoader | f0246, f0916 |
| ReadBoundedSources | SourceLoader | f0247, f0917 |
| TokenizeSduiSource | SduiLexer | f0248, f0853 |

Eksempelbane: UiCompilationAccepted. Kobling: f0988.

![Scenario: UiCompilationAccepted — modus SourceInspection](diagrams/VP08-UiCompilationAccepted.svg)

### G1M2ValidationAndNormalization

Status: **planned**. Kilde: f0260.

Forutsetninger: G1M1ParserAndAst (f0259).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| CoordinateSduiCompilation | SduiFrontend | f0251, f0816 |
| ExpandUiDefinitions | SduiNormalizer | f0252, f0865 |
| PreserveUiRegions | SduiNormalizer | f0253, f0866 |
| PreserveUiSourceMap | SduiNormalizer | f0254, f0867 |
| ResolveUiNames | SduiValidator | f0255, f0897 |
| ValidateRelativeFormatting | SduiValidator | f0256, f0898 |
| ValidateSymbolicBindings | SduiValidator | f0257, f0899 |
| ValidateWidgetArguments | SduiValidator | f0258, f0900 |

Eksempelbane: UiCompilationAccepted. Kobling: f0989.

[Scenariofigur](diagrams/VP08-UiCompilationAccepted.svg)

### G1M3Concept1AndDumps

Status: **planned**. Kilde: f0265.

Forutsetninger: G1M2ValidationAndNormalization (f0264).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| ExportConsoleSnapshot | SduiPresentation | f0262, f0877 |
| ReportSourceDiagnostics | DiagnosticReporter | f0263, f0171 |

## G2LayoutAndPresentation

Status: **planned**. Kilde: f0268.

Leveransebidrag: InteractiveUiPreview (f0267).

![Aktivitetsinndeling: G2LayoutAndPresentation](diagrams/VP06-detail-G2LayoutAndPresentation.svg)

### G2M1RelativeMeasurement

Status: **planned**. Kilde: f0274.

Forutsetninger: G1M2ValidationAndNormalization (f0273).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| AllocateGeometry | SduiLayout | f0269, f0843 |
| ComputeClipping | SduiLayout | f0270, f0845 |
| MeasureUiContent | SduiLayout | f0271, f0846 |
| ResolveAncestorDimensions | SduiLayout | f0272, f0847 |

Eksempelbane: InteractiveFramePrepared. Kobling: f0462.

![Scenario: InteractiveFramePrepared — modus UiPreview](diagrams/VP08-InteractiveFramePrepared.svg)

### G2M2SharedSvgGeometry

Status: **planned**. Kilde: f0279.

Forutsetninger: G2M1RelativeMeasurement (f0278).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| BuildPreparedFrame | SduiLayout | f0276, f0844 |
| ExportSvgSnapshot | SduiPresentation | f0277, f0878 |

Eksempelbane: StaticFrameExported. Kobling: f0933.

![Scenario: StaticFrameExported — modus StaticExport](diagrams/VP08-StaticFrameExported.svg)

### G2M3FyneInteractions

Status: **planned**. Kilde: f0288.

Forutsetninger: G2M2SharedSvgGeometry (f0287).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| ComposeInteractiveSession | FyneHost | f0281, f0234 |
| HandleFocusAndTextInput | FyneBackend | f0282, f0215 |
| PublishPresentation | FyneBackend | f0283, f0216 |
| ReconcileWidgets | FyneBackend | f0284, f0217 |
| ReleaseNativeWidgets | FyneBackend | f0285, f0218 |
| ScheduleUiPublication | FyneHost | f0286, f0235 |

Eksempelbane: InteractiveFramePrepared. Kobling: f0463.

[Scenariofigur](diagrams/VP08-InteractiveFramePrepared.svg)

Eksempelbane: UnboundLocalAction. Kobling: f1066.

![Scenario: UnboundLocalAction — modus UiPreview](diagrams/VP08-UnboundLocalAction.svg)

### G2M4RichContent

Status: **planned**. Kilde: f0296.

Forutsetninger: G2M3FyneInteractions (f0295).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| MeasureMarkdownContent | MarkdownProvider | f0290, f0503 |
| PrepareDiagramResource | DiagramProvider | f0291, f0175 |
| PrepareMarkdown | MarkdownProvider | f0292, f0504 |
| ReleaseVisualResources | ResourceStore | f0293, f0687 |
| ValidateVisualResources | ResourceStore | f0294, f0688 |

## G3UiRuntimeAndReload

Status: **planned**. Kilde: f0331.

Leveransebidrag: LiveModelReload (f0330).

![Aktivitetsinndeling: G3UiRuntimeAndReload](diagrams/VP06-detail-G3UiRuntimeAndReload.svg)

### G3M1TypedUiSession

Status: **planned**. Kilde: f0313.

Forutsetninger: G1M2ValidationAndNormalization (f0311), G2M3FyneInteractions (f0312).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| ApplyPropertyBatch | SduiPropertyStore | f0298, f0883 |
| CloseUiInstance | SduiRuntime | f0299, f0893 |
| CorrelateUiResult | SduiDispatcher | f0300, f0800 |
| CreateUiInstance | SduiRuntime | f0301, f0894 |
| DispatchUiEvent | SduiDispatcher | f0302, f0801 |
| ManageWidgetIdentities | SduiInstanceStore | f0303, f0833 |
| ProjectUiGeneration | SduiInstanceStore | f0304, f0834 |
| RejectStaleUiEvent | SduiDispatcher | f0305, f0802 |
| RevokeWidgetGenerations | SduiInstanceStore | f0306, f0835 |
| SnapshotUiState | SduiInstanceStore | f0307, f0836 |
| TrackInputDraft | SduiPropertyStore | f0308, f0884 |
| ValidatePropertyBatch | SduiPropertyStore | f0309, f0885 |
| ValidateUiEvent | SduiDispatcher | f0310, f0803 |

### G3M2CandidatePublication

Status: **planned**. Kilde: f0322.

Forutsetninger: G3M1TypedUiSession (f0321).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| CoalesceSourceChanges | SourceWatcher | f0315, f0920 |
| KeepLastValidModels | ReloadCoordinator | f0316, f0634 |
| ObserveSourceChanges | SourceWatcher | f0317, f0921 |
| PrepareCandidateModels | ReloadCoordinator | f0318, f0635 |
| PublishModelGeneration | ReloadCoordinator | f0319, f0636 |
| RetirePreviousGeneration | ReloadCoordinator | f0320, f0637 |

Eksempelbane: UiModelReloadAccepted. Kobling: f1026.

![Scenario: UiModelReloadAccepted — modus LiveEditing](diagrams/VP08-UiModelReloadAccepted.svg)

Eksempelbane: UiModelReloadRejected. Kobling: f1036.

![Scenario: UiModelReloadRejected — modus LiveEditing](diagrams/VP08-UiModelReloadRejected.svg)

### G3M3CompatibleState

Status: **planned**. Kilde: f0328.

Forutsetninger: G3M2CandidatePublication (f0327).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| MatchCompatibleWidgets | UiStateReconciler | f0324, f1052 |
| PreserveCompatibleUiState | UiStateReconciler | f0325, f1053 |
| ResetIncompatibleUiState | UiStateReconciler | f0326, f1054 |

Eksempelbane: UiModelReloadAccepted. Kobling: f1027.

[Scenariofigur](diagrams/VP08-UiModelReloadAccepted.svg)

## G4SdlRuntimeAndBinding

Status: **planned**. Kilde: f0375.

Leveransebidrag: TypedDomainBinding (f0374).

![Aktivitetsinndeling: G4SdlRuntimeAndBinding](diagrams/VP06-detail-G4SdlRuntimeAndBinding.svg)

### G4M1SdlFrontend

Status: **planned**. Kilde: f0340.

Forutsetninger: Ingen eksplisitt deklarert.

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| BuildSdlAst | SdlParser | f0332, f0762 |
| CoordinateSdlCompilation | SdlFrontend | f0333, f0735 |
| NormalizeSdlModel | SdlNormalizer | f0334, f0759 |
| PreserveSdlSourceMap | SdlNormalizer | f0335, f0760 |
| ResolveSdlSymbols | SdlValidator | f0336, f0791 |
| TokenizeSdlSource | SdlLexer | f0337, f0743 |
| ValidateSdlProfile | SdlValidator | f0338, f0792 |
| ValidateSdlStructure | SdlValidator | f0339, f0793 |

Eksempelbane: SdlModelReloadAccepted. Kobling: f0752.

![Scenario: SdlModelReloadAccepted — modus BoundLiveEditing](diagrams/VP08-SdlModelReloadAccepted.svg)

### G4M2TypedExecution

Status: **planned**. Kilde: f0355.

Forutsetninger: G4M1SdlFrontend (f0354).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| CancelPendingActions | SdlDispatcher | f0342, f0716 |
| CheckExecutionCompleteness | SdlExecutionGate | f0343, f0727 |
| CheckFunctionSignatures | SdlFunctionRegistry | f0344, f0740 |
| CloseSdlInstance | SdlRuntime | f0345, f0771 |
| CorrelateActionResult | SdlDispatcher | f0346, f0717 |
| CreateSdlInstance | SdlRuntime | f0347, f0772 |
| InvokeRegisteredFunction | SdlDispatcher | f0348, f0718 |
| ManageDomainState | SdlStateStore | f0349, f0775 |
| PerformDomainOperation | GoDomainImplementation | f0350, f0444 |
| RegisterDomainFunctions | SdlFunctionRegistry | f0351, f0741 |
| SnapshotDomainState | SdlStateStore | f0352, f0776 |
| ValidateActionInput | SdlDispatcher | f0353, f0719 |

### G4M3UiDomainBinding

Status: **planned**. Kilde: f0365.

Forutsetninger: G3M1TypedUiSession (f0363), G4M2TypedExecution (f0364).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| ConnectTypedWidgetHandles | SdlUiBindingAdapter | f0357, f0781 |
| DisconnectBindings | SdlUiBindingAdapter | f0358, f0782 |
| PublishDomainUpdates | SdlUiBindingAdapter | f0359, f0783 |
| ReportBindingDiagnostics | DiagnosticReporter | f0360, f0170 |
| ResolveCallbackSymbols | SdlUiBindingAdapter | f0361, f0784 |
| RouteDomainBindings | SdlUiBindingAdapter | f0362, f0785 |

Eksempelbane: BoundActionAccepted. Kobling: f0028.

![Scenario: BoundActionAccepted — modus BoundExecution](diagrams/VP08-BoundActionAccepted.svg)

Eksempelbane: BoundActionRejected. Kobling: f0041.

![Scenario: BoundActionRejected — modus BoundExecution](diagrams/VP08-BoundActionRejected.svg)

### G4M4DomainReload

Status: **planned**. Kilde: f0372.

Forutsetninger: G3M3CompatibleState (f0370), G4M3UiDomainBinding (f0371).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| CheckDomainStateCompatibility | DomainStateMigrator | f0367, f0188 |
| MigrateOrResetDomainState | DomainStateMigrator | f0368, f0189 |
| RestartChangedGoProgram | GoBuildRunner | f0369, f0427 |

Eksempelbane: SdlModelReloadAccepted. Kobling: f0753.

[Scenariofigur](diagrams/VP08-SdlModelReloadAccepted.svg)

## G5NativeGeneration

Status: **planned**. Kilde: f0404.

Leveransebidrag: DesignDocumentation (f0402), NativeGoAssembly (f0403).

![Aktivitetsinndeling: G5NativeGeneration](diagrams/VP06-detail-G5NativeGeneration.svg)

### G5M1GeneratedGo

Status: **planned**. Kilde: f0381.

Forutsetninger: G4M4DomainReload (f0380).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| BuildGeneratedApplication | GoBuildRunner | f0376, f0426 |
| GenerateBindingRegistration | GoCodeGenerator | f0377, f0434 |
| GenerateModelConstructors | GoCodeGenerator | f0378, f0435 |
| PreserveHandwrittenSources | GoCodeGenerator | f0379, f0436 |

Eksempelbane: NativeProgramBuilt. Kobling: f0532.

![Scenario: NativeProgramBuilt — modus NativeBuild](diagrams/VP08-NativeProgramBuilt.svg)

### G5M2BehaviorParity

Status: **planned**. Kilde: f0385.

Forutsetninger: G5M1GeneratedGo (f0384).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| VerifyNativeBehaviorParity | DevelopmentTools | f0383, f0168 |

### G5M3DocumentationExport

Status: **planned**. Kilde: f0395.

Forutsetninger: G2M4RichContent (f0393), G5M2BehaviorParity (f0394).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| ComposeHeadlessExport | CommandLineHost | f0387, f0090 |
| ComposeMarkdownDocument | SduiPresentation | f0388, f0876 |
| ExportViewpointMarkdown | SdlViewpointGenerator | f0389, f0795 |
| ProjectSdlViewpoints | SdlViewpointGenerator | f0390, f0796 |
| TraceViewpointFacts | SdlViewpointGenerator | f0391, f0797 |
| WriteGeneratedArtifacts | CommandLineHost | f0392, f0091 |

Eksempelbane: StaticFrameExported. Kobling: f0934.

[Scenariofigur](diagrams/VP08-StaticFrameExported.svg)

### G5M4RetirePython

Status: **planned**. Kilde: f0400.

Forutsetninger: G1M3Concept1AndDumps (f0398), G5M3DocumentationExport (f0399).

| Ansvar | Logisk eier | Kilde-ID-er |
| --- | --- | --- |
| RetireReplacedPythonEntryPoints | DevelopmentTools | f0397, f0167 |

## Udekket modellansvar

Ingen deklarerte Functionality-er mangler addresses-kobling. Dette beviser ikke full kravdekning.

Runtime-semantikk, full AST-/ABI-schema og fysisk layoutmåling må fortsatt realiseres og testes i implementasjonsfasene.
