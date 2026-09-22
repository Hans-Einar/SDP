# SDL-viewpoints — generert modellrapport

Generert fra validert SDL. Struktur og designpåstander, ikke observert kjøring.
Ingen håndskrevet arkitekturfakta er lagt til av generatoren.

| Viewpoint | Status | Grunnlag / mangel |
| --- | --- | --- |
| VP01 — Bruksmål og sporbarhet | Tilgjengelig | pursues, supports og contributes-to; modellens omfang, uten oppdiktet System-grense. |
| VP02 — Arkitektur og logisk inndeling | Tilgjengelig | Container/Unit og contains; bibliotekstruktur er ikke en deployment-allokering. |
| VP03 — Ansvar og kapabiliteter over arkitekturen | Tilgjengelig | owns, realizes og provides. Capability er ikke Feature. |
| VP04 — Grensesnitt og samarbeid | Tilgjengelig | consumes viser bruk; ingen tilbyder, Channel eller kjørbar meldingsflyt utledes. |
| VP05 — Avhengigheter per modus | Tilgjengelig | requires in mode; modi har ingen implisitt arv. |
| VP06 — Aktivitetsdetaljering | Tilgjengelig | refines er detaljering, ikke rekkefølge eller tilstandsoverganger. |
| VP07 — Features over arkitekturen | Tilgjengelig | contributes-to, owns og eksplisitt allocated-to per modus. Uspesifisert allokering vises som hull. |
| VP08 — Channel-kontrakter og sekvenser | Kan ikke genereres | Channel, deltakere, meldingskontrakt og ordnede scenario-/protokollsteg mangler. |
| VP09 — Dataset, Datagram og persistent Database | Tilgjengelig | Eksplisitte holdere, kilde, kontrakter, varianter, felt og projeksjoner. |
| VP10 — Datagram-koding og packet | Tilgjengelig | Kun closed kontrakt med validert Encoding og eksplisitte bitplasseringer. |
| VP11 — Egenskaper, sporbarhet og modellhull | Tilgjengelig | Deklarasjoner og alle fakta med kildeposisjoner; støttegrenser beholdes. |

## VP01 — Bruksmål og sporbarhet

Bruksmålskartene viser Actors, støttende Features og direkte Functionality-bidrag.
De etterfølgende Feature-kartene detaljerer bidragene med samme modellidentiteter.
Oppdelingen endrer ingen relasjoner og innfører ingen System-grense.

### Bruksmål: BuildNativeProduct

```mermaid
flowchart LR
    n_BuildNativeProduct["BuildNativeProduct (usecase)"]
    n_DesignAuthor["DesignAuthor (actor)"]
    n_NativeGoAssembly["NativeGoAssembly (feature)"]
    n_DesignAuthor -->|pursues| n_BuildNativeProduct
    n_NativeGoAssembly -->|supports| n_BuildNativeProduct
```

Kildegrunnlag: f0059, f0185.

### Bruksmål: EditRunningPrototype

```mermaid
flowchart LR
    n_DesignAuthor["DesignAuthor (actor)"]
    n_EditRunningPrototype["EditRunningPrototype (usecase)"]
    n_LiveModelReload["LiveModelReload (feature)"]
    n_StructuralModelInspection["StructuralModelInspection (feature)"]
    n_DesignAuthor -->|pursues| n_EditRunningPrototype
    n_LiveModelReload -->|supports| n_EditRunningPrototype
    n_StructuralModelInspection -->|supports| n_EditRunningPrototype
```

Kildegrunnlag: f0060, f0168, f0445.

### Bruksmål: InspectModels

```mermaid
flowchart LR
    n_DesignAuthor["DesignAuthor (actor)"]
    n_DesignReviewer["DesignReviewer (actor)"]
    n_InspectModels["InspectModels (usecase)"]
    n_StructuralModelInspection["StructuralModelInspection (feature)"]
    n_TraceViewpointFacts["TraceViewpointFacts (functionality)"]
    n_DesignAuthor -->|pursues| n_InspectModels
    n_DesignReviewer -->|pursues| n_InspectModels
    n_StructuralModelInspection -->|supports| n_InspectModels
    n_TraceViewpointFacts -->|contributes-to| n_InspectModels
```

Kildegrunnlag: f0061, f0066, f0446, f0456.

### Bruksmål: PrototypeUserInterface

```mermaid
flowchart LR
    n_DesignAuthor["DesignAuthor (actor)"]
    n_DesignReviewer["DesignReviewer (actor)"]
    n_InteractiveUiPreview["InteractiveUiPreview (feature)"]
    n_PrototypeUserInterface["PrototypeUserInterface (usecase)"]
    n_DesignAuthor -->|pursues| n_PrototypeUserInterface
    n_DesignReviewer -->|pursues| n_PrototypeUserInterface
    n_InteractiveUiPreview -->|supports| n_PrototypeUserInterface
```

Kildegrunnlag: f0062, f0067, f0161.

### Bruksmål: PublishDesignDocumentation

```mermaid
flowchart LR
    n_DesignAuthor["DesignAuthor (actor)"]
    n_DesignDocumentation["DesignDocumentation (feature)"]
    n_DesignReviewer["DesignReviewer (actor)"]
    n_PublishDesignDocumentation["PublishDesignDocumentation (usecase)"]
    n_DesignAuthor -->|pursues| n_PublishDesignDocumentation
    n_DesignDocumentation -->|supports| n_PublishDesignDocumentation
    n_DesignReviewer -->|pursues| n_PublishDesignDocumentation
```

Kildegrunnlag: f0063, f0065, f0068.

### Bruksmål: TryDomainInteraction

```mermaid
flowchart LR
    n_DesignAuthor["DesignAuthor (actor)"]
    n_TryDomainInteraction["TryDomainInteraction (usecase)"]
    n_TypedDomainBinding["TypedDomainBinding (feature)"]
    n_DesignAuthor -->|pursues| n_TryDomainInteraction
    n_TypedDomainBinding -->|supports| n_TryDomainInteraction
```

Kildegrunnlag: f0064, f0461.

### Functionality-bidrag til Feature: DesignDocumentation

```mermaid
flowchart LR
    n_ComposeMarkdownDocument["ComposeMarkdownDocument (functionality)"]
    n_DesignDocumentation["DesignDocumentation (feature)"]
    n_ExportSvgSnapshot["ExportSvgSnapshot (functionality)"]
    n_ExportViewpointMarkdown["ExportViewpointMarkdown (functionality)"]
    n_ProjectSdlViewpoints["ProjectSdlViewpoints (functionality)"]
    n_TraceViewpointFacts["TraceViewpointFacts (functionality)"]
    n_WriteGeneratedArtifacts["WriteGeneratedArtifacts (functionality)"]
    n_ComposeMarkdownDocument -->|contributes-to| n_DesignDocumentation
    n_ExportSvgSnapshot -->|contributes-to| n_DesignDocumentation
    n_ExportViewpointMarkdown -->|contributes-to| n_DesignDocumentation
    n_ProjectSdlViewpoints -->|contributes-to| n_DesignDocumentation
    n_TraceViewpointFacts -->|contributes-to| n_DesignDocumentation
    n_WriteGeneratedArtifacts -->|contributes-to| n_DesignDocumentation
```

Kildegrunnlag: f0043, f0111, f0115, f0221, f0455, f0507.

### Functionality-bidrag til Feature: InteractiveUiPreview

```mermaid
flowchart LR
    n_AllocateGeometry["AllocateGeometry (functionality)"]
    n_BuildPreparedFrame["BuildPreparedFrame (functionality)"]
    n_ComposeInteractiveSession["ComposeInteractiveSession (functionality)"]
    n_HandleFocusAndTextInput["HandleFocusAndTextInput (functionality)"]
    n_InteractiveUiPreview["InteractiveUiPreview (feature)"]
    n_ReconcileWidgets["ReconcileWidgets (functionality)"]
    n_AllocateGeometry -->|contributes-to| n_InteractiveUiPreview
    n_BuildPreparedFrame -->|contributes-to| n_InteractiveUiPreview
    n_ComposeInteractiveSession -->|contributes-to| n_InteractiveUiPreview
    n_HandleFocusAndTextInput -->|contributes-to| n_InteractiveUiPreview
    n_ReconcileWidgets -->|contributes-to| n_InteractiveUiPreview
```

Kildegrunnlag: f0001, f0012, f0040, f0156, f0237.

### Functionality-bidrag til Feature: LiveModelReload

```mermaid
flowchart LR
    n_KeepLastValidModels["KeepLastValidModels (functionality)"]
    n_LiveModelReload["LiveModelReload (feature)"]
    n_ObserveSourceChanges["ObserveSourceChanges (functionality)"]
    n_PrepareCandidateModels["PrepareCandidateModels (functionality)"]
    n_PreserveCompatibleUiState["PreserveCompatibleUiState (functionality)"]
    n_ProjectUiGeneration["ProjectUiGeneration (functionality)"]
    n_PublishModelGeneration["PublishModelGeneration (functionality)"]
    n_ReportSourceDiagnostics["ReportSourceDiagnostics (functionality)"]
    n_KeepLastValidModels -->|contributes-to| n_LiveModelReload
    n_ObserveSourceChanges -->|contributes-to| n_LiveModelReload
    n_PrepareCandidateModels -->|contributes-to| n_LiveModelReload
    n_PreserveCompatibleUiState -->|contributes-to| n_LiveModelReload
    n_ProjectUiGeneration -->|contributes-to| n_LiveModelReload
    n_PublishModelGeneration -->|contributes-to| n_LiveModelReload
    n_ReportSourceDiagnostics -->|contributes-to| n_LiveModelReload
```

Kildegrunnlag: f0166, f0197, f0201, f0206, f0225, f0231, f0262.

### Functionality-bidrag til Feature: NativeGoAssembly

```mermaid
flowchart LR
    n_BuildGeneratedApplication["BuildGeneratedApplication (functionality)"]
    n_GenerateBindingRegistration["GenerateBindingRegistration (functionality)"]
    n_GenerateModelConstructors["GenerateModelConstructors (functionality)"]
    n_NativeGoAssembly["NativeGoAssembly (feature)"]
    n_PreserveHandwrittenSources["PreserveHandwrittenSources (functionality)"]
    n_BuildGeneratedApplication -->|contributes-to| n_NativeGoAssembly
    n_GenerateBindingRegistration -->|contributes-to| n_NativeGoAssembly
    n_GenerateModelConstructors -->|contributes-to| n_NativeGoAssembly
    n_PreserveHandwrittenSources -->|contributes-to| n_NativeGoAssembly
```

Kildegrunnlag: f0008, f0136, f0139, f0209.

### Functionality-bidrag til Feature: StructuralModelInspection

```mermaid
flowchart LR
    n_BuildSdlAst["BuildSdlAst (functionality)"]
    n_BuildSduiAst["BuildSduiAst (functionality)"]
    n_ReportSourceDiagnostics["ReportSourceDiagnostics (functionality)"]
    n_StructuralModelInspection["StructuralModelInspection (feature)"]
    n_ValidateSdlStructure["ValidateSdlStructure (functionality)"]
    n_ValidateWidgetArguments["ValidateWidgetArguments (functionality)"]
    n_BuildSdlAst -->|contributes-to| n_StructuralModelInspection
    n_BuildSduiAst -->|contributes-to| n_StructuralModelInspection
    n_ReportSourceDiagnostics -->|contributes-to| n_StructuralModelInspection
    n_ValidateSdlStructure -->|contributes-to| n_StructuralModelInspection
    n_ValidateWidgetArguments -->|contributes-to| n_StructuralModelInspection
```

Kildegrunnlag: f0015, f0020, f0263, f0492, f0502.

### Functionality-bidrag til Feature: TypedDomainBinding

```mermaid
flowchart LR
    n_ConnectTypedWidgetHandles["ConnectTypedWidgetHandles (functionality)"]
    n_DispatchUiEvent["DispatchUiEvent (functionality)"]
    n_InvokeRegisteredFunction["InvokeRegisteredFunction (functionality)"]
    n_PublishDomainUpdates["PublishDomainUpdates (functionality)"]
    n_ResolveCallbackSymbols["ResolveCallbackSymbols (functionality)"]
    n_TypedDomainBinding["TypedDomainBinding (feature)"]
    n_ConnectTypedWidgetHandles -->|contributes-to| n_TypedDomainBinding
    n_DispatchUiEvent -->|contributes-to| n_TypedDomainBinding
    n_InvokeRegisteredFunction -->|contributes-to| n_TypedDomainBinding
    n_PublishDomainUpdates -->|contributes-to| n_TypedDomainBinding
    n_ResolveCallbackSymbols -->|contributes-to| n_TypedDomainBinding
```

Kildegrunnlag: f0047, f0098, f0163, f0228, f0268.


## VP02 — Arkitektur og logisk inndeling

Container er en erklært runtimegrense. Unit-røtter viser logisk struktur.
contains angir ikke deployment. Eksplisitt Functionality-allokering vises per modus i VP07.

### Arkitekturrøtter — ingen kobling/allokering er utledet

```mermaid
flowchart LR
    n_CommandLineHost["CommandLineHost (container)"]
    n_ContentServices["ContentServices (unit)"]
    n_DevelopmentTools["DevelopmentTools (unit)"]
    n_FyneBackend["FyneBackend (unit)"]
    n_FyneHost["FyneHost (container)"]
    n_GoDomainImplementation["GoDomainImplementation (unit)"]
    n_SdlLibrary["SdlLibrary (unit)"]
    n_SdlUiBindingAdapter["SdlUiBindingAdapter (unit)"]
    n_SduiLibrary["SduiLibrary (unit)"]
```

Kildegrunnlag: Kun deklarasjoner.

### Logisk inndeling: ContentServices

```mermaid
flowchart LR
    n_ContentServices["ContentServices (unit)"]
    n_DiagramProvider["DiagramProvider (unit)"]
    n_MarkdownProvider["MarkdownProvider (unit)"]
    n_ResourceStore["ResourceStore (unit)"]
    n_ContentServices -->|contains| n_DiagramProvider
    n_ContentServices -->|contains| n_MarkdownProvider
    n_ContentServices -->|contains| n_ResourceStore
```

Kildegrunnlag: f0049, f0050, f0051.

### Logisk inndeling: DevelopmentTools

```mermaid
flowchart LR
    n_DevelopmentTools["DevelopmentTools (unit)"]
    n_DiagnosticReporter["DiagnosticReporter (unit)"]
    n_GoBuildRunner["GoBuildRunner (unit)"]
    n_GoCodeGenerator["GoCodeGenerator (unit)"]
    n_ReloadCoordinator["ReloadCoordinator (unit)"]
    n_SourceLoader["SourceLoader (unit)"]
    n_SourceWatcher["SourceWatcher (unit)"]
    n_DevelopmentTools -->|contains| n_DiagnosticReporter
    n_DevelopmentTools -->|contains| n_GoBuildRunner
    n_DevelopmentTools -->|contains| n_GoCodeGenerator
    n_DevelopmentTools -->|contains| n_ReloadCoordinator
    n_DevelopmentTools -->|contains| n_SourceLoader
    n_DevelopmentTools -->|contains| n_SourceWatcher
```

Kildegrunnlag: f0082, f0083, f0084, f0085, f0086, f0087.

### Logisk inndeling: SdlFrontend

```mermaid
flowchart LR
    n_SdlFrontend["SdlFrontend (unit)"]
    n_SdlLexer["SdlLexer (unit)"]
    n_SdlNormalizer["SdlNormalizer (unit)"]
    n_SdlParser["SdlParser (unit)"]
    n_SdlValidator["SdlValidator (unit)"]
    n_SdlFrontend -->|contains| n_SdlLexer
    n_SdlFrontend -->|contains| n_SdlNormalizer
    n_SdlFrontend -->|contains| n_SdlParser
    n_SdlFrontend -->|contains| n_SdlValidator
```

Kildegrunnlag: f0301, f0302, f0303, f0304.

### Logisk inndeling: SdlLibrary

```mermaid
flowchart LR
    n_SdlFrontend["SdlFrontend (unit)"]
    n_SdlLibrary["SdlLibrary (unit)"]
    n_SdlRuntime["SdlRuntime (unit)"]
    n_SdlViewpointGenerator["SdlViewpointGenerator (unit)"]
    n_SdlLibrary -->|contains| n_SdlFrontend
    n_SdlLibrary -->|contains| n_SdlRuntime
    n_SdlLibrary -->|contains| n_SdlViewpointGenerator
```

Kildegrunnlag: f0313, f0314, f0315.

### Logisk inndeling: SdlRuntime

```mermaid
flowchart LR
    n_DomainStateMigrator["DomainStateMigrator (unit)"]
    n_SdlDispatcher["SdlDispatcher (unit)"]
    n_SdlExecutionGate["SdlExecutionGate (unit)"]
    n_SdlFunctionRegistry["SdlFunctionRegistry (unit)"]
    n_SdlRuntime["SdlRuntime (unit)"]
    n_SdlStateStore["SdlStateStore (unit)"]
    n_SdlRuntime -->|contains| n_DomainStateMigrator
    n_SdlRuntime -->|contains| n_SdlDispatcher
    n_SdlRuntime -->|contains| n_SdlExecutionGate
    n_SdlRuntime -->|contains| n_SdlFunctionRegistry
    n_SdlRuntime -->|contains| n_SdlStateStore
```

Kildegrunnlag: f0325, f0326, f0327, f0328, f0329.

### Logisk inndeling: SduiFrontend

```mermaid
flowchart LR
    n_SduiFrontend["SduiFrontend (unit)"]
    n_SduiLexer["SduiLexer (unit)"]
    n_SduiNormalizer["SduiNormalizer (unit)"]
    n_SduiParser["SduiParser (unit)"]
    n_SduiValidator["SduiValidator (unit)"]
    n_SduiFrontend -->|contains| n_SduiLexer
    n_SduiFrontend -->|contains| n_SduiNormalizer
    n_SduiFrontend -->|contains| n_SduiParser
    n_SduiFrontend -->|contains| n_SduiValidator
```

Kildegrunnlag: f0362, f0363, f0364, f0365.

### Logisk inndeling: SduiLibrary

```mermaid
flowchart LR
    n_SduiFrontend["SduiFrontend (unit)"]
    n_SduiLayout["SduiLayout (unit)"]
    n_SduiLibrary["SduiLibrary (unit)"]
    n_SduiPresentation["SduiPresentation (unit)"]
    n_SduiRuntime["SduiRuntime (unit)"]
    n_SduiLibrary -->|contains| n_SduiFrontend
    n_SduiLibrary -->|contains| n_SduiLayout
    n_SduiLibrary -->|contains| n_SduiPresentation
    n_SduiLibrary -->|contains| n_SduiRuntime
```

Kildegrunnlag: f0385, f0386, f0387, f0388.

### Logisk inndeling: SduiRuntime

```mermaid
flowchart LR
    n_SduiDispatcher["SduiDispatcher (unit)"]
    n_SduiInstanceStore["SduiInstanceStore (unit)"]
    n_SduiPropertyStore["SduiPropertyStore (unit)"]
    n_SduiRuntime["SduiRuntime (unit)"]
    n_UiStateReconciler["UiStateReconciler (unit)"]
    n_SduiRuntime -->|contains| n_SduiDispatcher
    n_SduiRuntime -->|contains| n_SduiInstanceStore
    n_SduiRuntime -->|contains| n_SduiPropertyStore
    n_SduiRuntime -->|contains| n_UiStateReconciler
```

Kildegrunnlag: f0411, f0412, f0413, f0414.


## VP03 — Ansvar og kapabiliteter over arkitekturen

### Bidrag til kapabilitet: BoundInteraction

```mermaid
flowchart LR
    n_BoundInteraction["BoundInteraction (capability)"]
    n_ConnectTypedWidgetHandles["ConnectTypedWidgetHandles (functionality)"]
    n_DisconnectBindings["DisconnectBindings (functionality)"]
    n_PublishDomainUpdates["PublishDomainUpdates (functionality)"]
    n_ResolveCallbackSymbols["ResolveCallbackSymbols (functionality)"]
    n_RouteDomainBindings["RouteDomainBindings (functionality)"]
    n_SdlUiBindingAdapter["SdlUiBindingAdapter (unit)"]
    n_ConnectTypedWidgetHandles -->|realizes| n_BoundInteraction
    n_DisconnectBindings -->|realizes| n_BoundInteraction
    n_PublishDomainUpdates -->|realizes| n_BoundInteraction
    n_ResolveCallbackSymbols -->|realizes| n_BoundInteraction
    n_RouteDomainBindings -->|realizes| n_BoundInteraction
    n_SdlUiBindingAdapter -->|owns| n_ConnectTypedWidgetHandles
    n_SdlUiBindingAdapter -->|owns| n_DisconnectBindings
    n_SdlUiBindingAdapter -->|owns| n_PublishDomainUpdates
    n_SdlUiBindingAdapter -->|owns| n_ResolveCallbackSymbols
    n_SdlUiBindingAdapter -->|owns| n_RouteDomainBindings
```

Kildegrunnlag: f0048, f0096, f0229, f0269, f0284, f0340, f0341, f0342, f0343, f0344.

### Tilbydere av kapabilitet: BoundInteraction

```mermaid
flowchart LR
    n_BoundInteraction["BoundInteraction (capability)"]
    n_SdlUiBindingAdapter["SdlUiBindingAdapter (unit)"]
    n_SdlUiBindingAdapter -->|provides| n_BoundInteraction
```

Kildegrunnlag: f0345.

### Bidrag til kapabilitet: DevelopmentReload

```mermaid
flowchart LR
    n_CoalesceSourceChanges["CoalesceSourceChanges (functionality)"]
    n_DevelopmentReload["DevelopmentReload (capability)"]
    n_KeepLastValidModels["KeepLastValidModels (functionality)"]
    n_ObserveSourceChanges["ObserveSourceChanges (functionality)"]
    n_PrepareCandidateModels["PrepareCandidateModels (functionality)"]
    n_PublishModelGeneration["PublishModelGeneration (functionality)"]
    n_ReloadCoordinator["ReloadCoordinator (unit)"]
    n_RetirePreviousGeneration["RetirePreviousGeneration (functionality)"]
    n_SourceWatcher["SourceWatcher (unit)"]
    n_CoalesceSourceChanges -->|realizes| n_DevelopmentReload
    n_KeepLastValidModels -->|realizes| n_DevelopmentReload
    n_ObserveSourceChanges -->|realizes| n_DevelopmentReload
    n_PrepareCandidateModels -->|realizes| n_DevelopmentReload
    n_PublishModelGeneration -->|realizes| n_DevelopmentReload
    n_ReloadCoordinator -->|owns| n_KeepLastValidModels
    n_ReloadCoordinator -->|owns| n_PrepareCandidateModels
    n_ReloadCoordinator -->|owns| n_PublishModelGeneration
    n_ReloadCoordinator -->|owns| n_RetirePreviousGeneration
    n_RetirePreviousGeneration -->|realizes| n_DevelopmentReload
    n_SourceWatcher -->|owns| n_CoalesceSourceChanges
    n_SourceWatcher -->|owns| n_ObserveSourceChanges
```

Kildegrunnlag: f0030, f0167, f0198, f0202, f0233, f0253, f0254, f0255, f0256, f0280, f0440, f0441.

### Tilbydere av kapabilitet: DevelopmentReload

```mermaid
flowchart LR
    n_DevelopmentReload["DevelopmentReload (capability)"]
    n_DevelopmentTools["DevelopmentTools (unit)"]
    n_ReloadCoordinator["ReloadCoordinator (unit)"]
    n_SourceWatcher["SourceWatcher (unit)"]
    n_DevelopmentTools -->|provides| n_DevelopmentReload
    n_ReloadCoordinator -->|provides| n_DevelopmentReload
    n_SourceWatcher -->|provides| n_DevelopmentReload
```

Kildegrunnlag: f0088, f0257, f0442.

### Bidrag til kapabilitet: DomainOperations

```mermaid
flowchart LR
    n_DomainOperations["DomainOperations (capability)"]
    n_GoDomainImplementation["GoDomainImplementation (unit)"]
    n_PerformDomainOperation["PerformDomainOperation (functionality)"]
    n_GoDomainImplementation -->|owns| n_PerformDomainOperation
    n_PerformDomainOperation -->|realizes| n_DomainOperations
```

Kildegrunnlag: f0153, f0199.

### Tilbydere av kapabilitet: DomainOperations

```mermaid
flowchart LR
    n_DomainOperations["DomainOperations (capability)"]
    n_GoDomainImplementation["GoDomainImplementation (unit)"]
    n_GoDomainImplementation -->|provides| n_DomainOperations
```

Kildegrunnlag: f0154.

### Bidrag til kapabilitet: ExecutableDesign

```mermaid
flowchart LR
    n_CancelPendingActions["CancelPendingActions (functionality)"]
    n_CheckDomainStateCompatibility["CheckDomainStateCompatibility (functionality)"]
    n_CheckExecutionCompleteness["CheckExecutionCompleteness (functionality)"]
    n_CheckFunctionSignatures["CheckFunctionSignatures (functionality)"]
    n_CloseSdlInstance["CloseSdlInstance (functionality)"]
    n_CorrelateActionResult["CorrelateActionResult (functionality)"]
    n_CreateSdlInstance["CreateSdlInstance (functionality)"]
    n_DomainStateMigrator["DomainStateMigrator (unit)"]
    n_ExecutableDesign["ExecutableDesign (capability)"]
    n_InvokeRegisteredFunction["InvokeRegisteredFunction (functionality)"]
    n_ManageDomainState["ManageDomainState (functionality)"]
    n_MigrateOrResetDomainState["MigrateOrResetDomainState (functionality)"]
    n_RegisterDomainFunctions["RegisterDomainFunctions (functionality)"]
    n_SdlDispatcher["SdlDispatcher (unit)"]
    n_SdlExecutionGate["SdlExecutionGate (unit)"]
    n_SdlFunctionRegistry["SdlFunctionRegistry (unit)"]
    n_SdlRuntime["SdlRuntime (unit)"]
    n_SdlStateStore["SdlStateStore (unit)"]
    n_SnapshotDomainState["SnapshotDomainState (functionality)"]
    n_ValidateActionInput["ValidateActionInput (functionality)"]
    n_CancelPendingActions -->|realizes| n_ExecutableDesign
    n_CheckDomainStateCompatibility -->|realizes| n_ExecutableDesign
    n_CheckExecutionCompleteness -->|realizes| n_ExecutableDesign
    n_CheckFunctionSignatures -->|realizes| n_ExecutableDesign
    n_CloseSdlInstance -->|realizes| n_ExecutableDesign
    n_CorrelateActionResult -->|realizes| n_ExecutableDesign
    n_CreateSdlInstance -->|realizes| n_ExecutableDesign
    n_DomainStateMigrator -->|owns| n_CheckDomainStateCompatibility
    n_DomainStateMigrator -->|owns| n_MigrateOrResetDomainState
    n_InvokeRegisteredFunction -->|realizes| n_ExecutableDesign
    n_ManageDomainState -->|realizes| n_ExecutableDesign
    n_MigrateOrResetDomainState -->|realizes| n_ExecutableDesign
    n_RegisterDomainFunctions -->|realizes| n_ExecutableDesign
    n_SdlDispatcher -->|owns| n_CancelPendingActions
    n_SdlDispatcher -->|owns| n_CorrelateActionResult
    n_SdlDispatcher -->|owns| n_InvokeRegisteredFunction
    n_SdlDispatcher -->|owns| n_ValidateActionInput
    n_SdlExecutionGate -->|owns| n_CheckExecutionCompleteness
    n_SdlFunctionRegistry -->|owns| n_CheckFunctionSignatures
    n_SdlFunctionRegistry -->|owns| n_RegisterDomainFunctions
    n_SdlRuntime -->|owns| n_CloseSdlInstance
    n_SdlRuntime -->|owns| n_CreateSdlInstance
    n_SdlStateStore -->|owns| n_ManageDomainState
    n_SdlStateStore -->|owns| n_SnapshotDomainState
    n_SnapshotDomainState -->|realizes| n_ExecutableDesign
    n_ValidateActionInput -->|realizes| n_ExecutableDesign
```

Kildegrunnlag: f0024, f0025, f0026, f0027, f0028, f0055, f0057, f0102, f0103, f0164, f0170, f0184, f0241, f0290, f0291, f0292, f0293, f0297, f0308, f0309, f0330, f0331, f0334, f0335, f0428, f0483.

### Tilbydere av kapabilitet: ExecutableDesign

```mermaid
flowchart LR
    n_DomainStateMigrator["DomainStateMigrator (unit)"]
    n_ExecutableDesign["ExecutableDesign (capability)"]
    n_SdlDispatcher["SdlDispatcher (unit)"]
    n_SdlExecutionGate["SdlExecutionGate (unit)"]
    n_SdlFunctionRegistry["SdlFunctionRegistry (unit)"]
    n_SdlLibrary["SdlLibrary (unit)"]
    n_SdlRuntime["SdlRuntime (unit)"]
    n_SdlStateStore["SdlStateStore (unit)"]
    n_DomainStateMigrator -->|provides| n_ExecutableDesign
    n_SdlDispatcher -->|provides| n_ExecutableDesign
    n_SdlExecutionGate -->|provides| n_ExecutableDesign
    n_SdlFunctionRegistry -->|provides| n_ExecutableDesign
    n_SdlLibrary -->|provides| n_ExecutableDesign
    n_SdlRuntime -->|provides| n_ExecutableDesign
    n_SdlStateStore -->|provides| n_ExecutableDesign
```

Kildegrunnlag: f0104, f0294, f0298, f0310, f0316, f0332, f0336.

### Bidrag til kapabilitet: InteractiveSession

```mermaid
flowchart LR
    n_ApplyPropertyBatch["ApplyPropertyBatch (functionality)"]
    n_CloseUiInstance["CloseUiInstance (functionality)"]
    n_CorrelateUiResult["CorrelateUiResult (functionality)"]
    n_CreateUiInstance["CreateUiInstance (functionality)"]
    n_DispatchUiEvent["DispatchUiEvent (functionality)"]
    n_InteractiveSession["InteractiveSession (capability)"]
    n_ManageWidgetIdentities["ManageWidgetIdentities (functionality)"]
    n_MatchCompatibleWidgets["MatchCompatibleWidgets (functionality)"]
    n_PreserveCompatibleUiState["PreserveCompatibleUiState (functionality)"]
    n_RejectStaleUiEvent["RejectStaleUiEvent (functionality)"]
    n_ResetIncompatibleUiState["ResetIncompatibleUiState (functionality)"]
    n_RevokeWidgetGenerations["RevokeWidgetGenerations (functionality)"]
    n_SduiDispatcher["SduiDispatcher (unit)"]
    n_SduiInstanceStore["SduiInstanceStore (unit)"]
    n_SduiPropertyStore["SduiPropertyStore (unit)"]
    n_SduiRuntime["SduiRuntime (unit)"]
    n_SnapshotUiState["SnapshotUiState (functionality)"]
    n_TrackInputDraft["TrackInputDraft (functionality)"]
    n_UiStateReconciler["UiStateReconciler (unit)"]
    n_ValidatePropertyBatch["ValidatePropertyBatch (functionality)"]
    n_ValidateUiEvent["ValidateUiEvent (functionality)"]
    n_ApplyPropertyBatch -->|realizes| n_InteractiveSession
    n_CloseUiInstance -->|realizes| n_InteractiveSession
    n_CorrelateUiResult -->|realizes| n_InteractiveSession
    n_CreateUiInstance -->|realizes| n_InteractiveSession
    n_DispatchUiEvent -->|realizes| n_InteractiveSession
    n_ManageWidgetIdentities -->|realizes| n_InteractiveSession
    n_MatchCompatibleWidgets -->|realizes| n_InteractiveSession
    n_PreserveCompatibleUiState -->|realizes| n_InteractiveSession
    n_RejectStaleUiEvent -->|realizes| n_InteractiveSession
    n_ResetIncompatibleUiState -->|realizes| n_InteractiveSession
    n_RevokeWidgetGenerations -->|realizes| n_InteractiveSession
    n_SduiDispatcher -->|owns| n_CorrelateUiResult
    n_SduiDispatcher -->|owns| n_DispatchUiEvent
    n_SduiDispatcher -->|owns| n_RejectStaleUiEvent
    n_SduiDispatcher -->|owns| n_ValidateUiEvent
    n_SduiInstanceStore -->|owns| n_ManageWidgetIdentities
    n_SduiInstanceStore -->|owns| n_RevokeWidgetGenerations
    n_SduiInstanceStore -->|owns| n_SnapshotUiState
    n_SduiPropertyStore -->|owns| n_ApplyPropertyBatch
    n_SduiPropertyStore -->|owns| n_TrackInputDraft
    n_SduiPropertyStore -->|owns| n_ValidatePropertyBatch
    n_SduiRuntime -->|owns| n_CloseUiInstance
    n_SduiRuntime -->|owns| n_CreateUiInstance
    n_SnapshotUiState -->|realizes| n_InteractiveSession
    n_TrackInputDraft -->|realizes| n_InteractiveSession
    n_UiStateReconciler -->|owns| n_MatchCompatibleWidgets
    n_UiStateReconciler -->|owns| n_PreserveCompatibleUiState
    n_UiStateReconciler -->|owns| n_ResetIncompatibleUiState
    n_ValidatePropertyBatch -->|realizes| n_InteractiveSession
    n_ValidateUiEvent -->|realizes| n_InteractiveSession
```

Kildegrunnlag: f0004, f0029, f0056, f0058, f0099, f0172, f0179, f0207, f0242, f0265, f0281, f0355, f0356, f0357, f0358, f0369, f0371, f0372, f0405, f0406, f0407, f0415, f0416, f0429, f0460, f0479, f0480, f0481, f0484, f0499.

### Tilbydere av kapabilitet: InteractiveSession

```mermaid
flowchart LR
    n_InteractiveSession["InteractiveSession (capability)"]
    n_SduiDispatcher["SduiDispatcher (unit)"]
    n_SduiInstanceStore["SduiInstanceStore (unit)"]
    n_SduiLibrary["SduiLibrary (unit)"]
    n_SduiPropertyStore["SduiPropertyStore (unit)"]
    n_SduiRuntime["SduiRuntime (unit)"]
    n_UiStateReconciler["UiStateReconciler (unit)"]
    n_SduiDispatcher -->|provides| n_InteractiveSession
    n_SduiInstanceStore -->|provides| n_InteractiveSession
    n_SduiLibrary -->|provides| n_InteractiveSession
    n_SduiPropertyStore -->|provides| n_InteractiveSession
    n_SduiRuntime -->|provides| n_InteractiveSession
    n_UiStateReconciler -->|provides| n_InteractiveSession
```

Kildegrunnlag: f0359, f0373, f0389, f0408, f0417, f0482.

### Bidrag til kapabilitet: MeasuredPresentation

```mermaid
flowchart LR
    n_AllocateGeometry["AllocateGeometry (functionality)"]
    n_BuildPreparedFrame["BuildPreparedFrame (functionality)"]
    n_ComputeClipping["ComputeClipping (functionality)"]
    n_MeasureUiContent["MeasureUiContent (functionality)"]
    n_MeasuredPresentation["MeasuredPresentation (capability)"]
    n_ResolveAncestorDimensions["ResolveAncestorDimensions (functionality)"]
    n_SduiLayout["SduiLayout (unit)"]
    n_AllocateGeometry -->|realizes| n_MeasuredPresentation
    n_BuildPreparedFrame -->|realizes| n_MeasuredPresentation
    n_ComputeClipping -->|realizes| n_MeasuredPresentation
    n_MeasureUiContent -->|realizes| n_MeasuredPresentation
    n_ResolveAncestorDimensions -->|realizes| n_MeasuredPresentation
    n_SduiLayout -->|owns| n_AllocateGeometry
    n_SduiLayout -->|owns| n_BuildPreparedFrame
    n_SduiLayout -->|owns| n_ComputeClipping
    n_SduiLayout -->|owns| n_MeasureUiContent
    n_SduiLayout -->|owns| n_ResolveAncestorDimensions
```

Kildegrunnlag: f0002, f0013, f0045, f0181, f0266, f0377, f0378, f0379, f0380, f0381.

### Tilbydere av kapabilitet: MeasuredPresentation

```mermaid
flowchart LR
    n_MeasuredPresentation["MeasuredPresentation (capability)"]
    n_SduiLayout["SduiLayout (unit)"]
    n_SduiLibrary["SduiLibrary (unit)"]
    n_SduiLayout -->|provides| n_MeasuredPresentation
    n_SduiLibrary -->|provides| n_MeasuredPresentation
```

Kildegrunnlag: f0382, f0390.

### Bidrag til kapabilitet: NativeInteraction

```mermaid
flowchart LR
    n_ComposeInteractiveSession["ComposeInteractiveSession (functionality)"]
    n_FyneBackend["FyneBackend (unit)"]
    n_FyneHost["FyneHost (container)"]
    n_HandleFocusAndTextInput["HandleFocusAndTextInput (functionality)"]
    n_NativeInteraction["NativeInteraction (capability)"]
    n_PublishPresentation["PublishPresentation (functionality)"]
    n_ReconcileWidgets["ReconcileWidgets (functionality)"]
    n_ReleaseNativeWidgets["ReleaseNativeWidgets (functionality)"]
    n_ScheduleUiPublication["ScheduleUiPublication (functionality)"]
    n_ComposeInteractiveSession -->|realizes| n_NativeInteraction
    n_FyneBackend -->|owns| n_HandleFocusAndTextInput
    n_FyneBackend -->|owns| n_PublishPresentation
    n_FyneBackend -->|owns| n_ReconcileWidgets
    n_FyneBackend -->|owns| n_ReleaseNativeWidgets
    n_FyneHost -->|owns| n_ComposeInteractiveSession
    n_FyneHost -->|owns| n_ScheduleUiPublication
    n_HandleFocusAndTextInput -->|realizes| n_NativeInteraction
    n_PublishPresentation -->|realizes| n_NativeInteraction
    n_ReconcileWidgets -->|realizes| n_NativeInteraction
    n_ReleaseNativeWidgets -->|realizes| n_NativeInteraction
    n_ScheduleUiPublication -->|realizes| n_NativeInteraction
```

Kildegrunnlag: f0041, f0120, f0121, f0122, f0123, f0132, f0133, f0157, f0234, f0239, f0243, f0287.

### Tilbydere av kapabilitet: NativeInteraction

```mermaid
flowchart LR
    n_FyneBackend["FyneBackend (unit)"]
    n_FyneHost["FyneHost (container)"]
    n_NativeInteraction["NativeInteraction (capability)"]
    n_FyneBackend -->|provides| n_NativeInteraction
    n_FyneHost -->|provides| n_NativeInteraction
```

Kildegrunnlag: f0124, f0134.

### Bidrag til kapabilitet: NativeRealization

```mermaid
flowchart LR
    n_BuildGeneratedApplication["BuildGeneratedApplication (functionality)"]
    n_GenerateBindingRegistration["GenerateBindingRegistration (functionality)"]
    n_GenerateModelConstructors["GenerateModelConstructors (functionality)"]
    n_GoBuildRunner["GoBuildRunner (unit)"]
    n_GoCodeGenerator["GoCodeGenerator (unit)"]
    n_NativeRealization["NativeRealization (capability)"]
    n_PreserveHandwrittenSources["PreserveHandwrittenSources (functionality)"]
    n_RestartChangedGoProgram["RestartChangedGoProgram (functionality)"]
    n_BuildGeneratedApplication -->|realizes| n_NativeRealization
    n_GenerateBindingRegistration -->|realizes| n_NativeRealization
    n_GenerateModelConstructors -->|realizes| n_NativeRealization
    n_GoBuildRunner -->|owns| n_BuildGeneratedApplication
    n_GoBuildRunner -->|owns| n_RestartChangedGoProgram
    n_GoCodeGenerator -->|owns| n_GenerateBindingRegistration
    n_GoCodeGenerator -->|owns| n_GenerateModelConstructors
    n_GoCodeGenerator -->|owns| n_PreserveHandwrittenSources
    n_PreserveHandwrittenSources -->|realizes| n_NativeRealization
    n_RestartChangedGoProgram -->|realizes| n_NativeRealization
```

Kildegrunnlag: f0009, f0137, f0140, f0143, f0144, f0149, f0150, f0151, f0210, f0279.

### Tilbydere av kapabilitet: NativeRealization

```mermaid
flowchart LR
    n_GoBuildRunner["GoBuildRunner (unit)"]
    n_GoCodeGenerator["GoCodeGenerator (unit)"]
    n_NativeRealization["NativeRealization (capability)"]
    n_GoBuildRunner -->|provides| n_NativeRealization
    n_GoCodeGenerator -->|provides| n_NativeRealization
```

Kildegrunnlag: f0145, f0152.

### Bidrag til kapabilitet: RichContent

```mermaid
flowchart LR
    n_DiagramProvider["DiagramProvider (unit)"]
    n_MarkdownProvider["MarkdownProvider (unit)"]
    n_MeasureMarkdownContent["MeasureMarkdownContent (functionality)"]
    n_PrepareDiagramResource["PrepareDiagramResource (functionality)"]
    n_PrepareMarkdown["PrepareMarkdown (functionality)"]
    n_ReleaseVisualResources["ReleaseVisualResources (functionality)"]
    n_ResourceStore["ResourceStore (unit)"]
    n_RichContent["RichContent (capability)"]
    n_ValidateVisualResources["ValidateVisualResources (functionality)"]
    n_DiagramProvider -->|owns| n_PrepareDiagramResource
    n_MarkdownProvider -->|owns| n_MeasureMarkdownContent
    n_MarkdownProvider -->|owns| n_PrepareMarkdown
    n_MeasureMarkdownContent -->|realizes| n_RichContent
    n_PrepareDiagramResource -->|realizes| n_RichContent
    n_PrepareMarkdown -->|realizes| n_RichContent
    n_ReleaseVisualResources -->|realizes| n_RichContent
    n_ResourceStore -->|owns| n_ReleaseVisualResources
    n_ResourceStore -->|owns| n_ValidateVisualResources
    n_ValidateVisualResources -->|realizes| n_RichContent
```

Kildegrunnlag: f0094, f0176, f0177, f0180, f0203, f0204, f0244, f0276, f0277, f0500.

### Tilbydere av kapabilitet: RichContent

```mermaid
flowchart LR
    n_ContentServices["ContentServices (unit)"]
    n_DiagramProvider["DiagramProvider (unit)"]
    n_MarkdownProvider["MarkdownProvider (unit)"]
    n_ResourceStore["ResourceStore (unit)"]
    n_RichContent["RichContent (capability)"]
    n_ContentServices -->|provides| n_RichContent
    n_DiagramProvider -->|provides| n_RichContent
    n_MarkdownProvider -->|provides| n_RichContent
    n_ResourceStore -->|provides| n_RichContent
```

Kildegrunnlag: f0052, f0095, f0178, f0278.

### Bidrag til kapabilitet: SdlSourceModel

```mermaid
flowchart LR
    n_BuildSdlAst["BuildSdlAst (functionality)"]
    n_CoordinateSdlCompilation["CoordinateSdlCompilation (functionality)"]
    n_NormalizeSdlModel["NormalizeSdlModel (functionality)"]
    n_PreserveSdlSourceMap["PreserveSdlSourceMap (functionality)"]
    n_ResolveSdlSymbols["ResolveSdlSymbols (functionality)"]
    n_SdlFrontend["SdlFrontend (unit)"]
    n_SdlLexer["SdlLexer (unit)"]
    n_SdlNormalizer["SdlNormalizer (unit)"]
    n_SdlParser["SdlParser (unit)"]
    n_SdlSourceModel["SdlSourceModel (capability)"]
    n_SdlValidator["SdlValidator (unit)"]
    n_TokenizeSdlSource["TokenizeSdlSource (functionality)"]
    n_ValidateSdlProfile["ValidateSdlProfile (functionality)"]
    n_ValidateSdlStructure["ValidateSdlStructure (functionality)"]
    n_BuildSdlAst -->|realizes| n_SdlSourceModel
    n_CoordinateSdlCompilation -->|realizes| n_SdlSourceModel
    n_NormalizeSdlModel -->|realizes| n_SdlSourceModel
    n_PreserveSdlSourceMap -->|realizes| n_SdlSourceModel
    n_ResolveSdlSymbols -->|realizes| n_SdlSourceModel
    n_SdlFrontend -->|owns| n_CoordinateSdlCompilation
    n_SdlLexer -->|owns| n_TokenizeSdlSource
    n_SdlNormalizer -->|owns| n_NormalizeSdlModel
    n_SdlNormalizer -->|owns| n_PreserveSdlSourceMap
    n_SdlParser -->|owns| n_BuildSdlAst
    n_SdlValidator -->|owns| n_ResolveSdlSymbols
    n_SdlValidator -->|owns| n_ValidateSdlProfile
    n_SdlValidator -->|owns| n_ValidateSdlStructure
    n_TokenizeSdlSource -->|realizes| n_SdlSourceModel
    n_ValidateSdlProfile -->|realizes| n_SdlSourceModel
    n_ValidateSdlStructure -->|realizes| n_SdlSourceModel
```

Kildegrunnlag: f0018, f0053, f0191, f0213, f0272, f0305, f0311, f0318, f0319, f0321, f0346, f0347, f0348, f0449, f0490, f0495.

### Tilbydere av kapabilitet: SdlSourceModel

```mermaid
flowchart LR
    n_SdlFrontend["SdlFrontend (unit)"]
    n_SdlLexer["SdlLexer (unit)"]
    n_SdlLibrary["SdlLibrary (unit)"]
    n_SdlNormalizer["SdlNormalizer (unit)"]
    n_SdlParser["SdlParser (unit)"]
    n_SdlSourceModel["SdlSourceModel (capability)"]
    n_SdlValidator["SdlValidator (unit)"]
    n_SdlFrontend -->|provides| n_SdlSourceModel
    n_SdlLexer -->|provides| n_SdlSourceModel
    n_SdlLibrary -->|provides| n_SdlSourceModel
    n_SdlNormalizer -->|provides| n_SdlSourceModel
    n_SdlParser -->|provides| n_SdlSourceModel
    n_SdlValidator -->|provides| n_SdlSourceModel
```

Kildegrunnlag: f0306, f0312, f0317, f0320, f0322, f0349.

### Bidrag til kapabilitet: SduiSourceModel

```mermaid
flowchart LR
    n_BuildSduiAst["BuildSduiAst (functionality)"]
    n_CoordinateSduiCompilation["CoordinateSduiCompilation (functionality)"]
    n_ExpandUiDefinitions["ExpandUiDefinitions (functionality)"]
    n_PreserveUiRegions["PreserveUiRegions (functionality)"]
    n_PreserveUiSourceMap["PreserveUiSourceMap (functionality)"]
    n_ResolveUiNames["ResolveUiNames (functionality)"]
    n_SduiFrontend["SduiFrontend (unit)"]
    n_SduiLexer["SduiLexer (unit)"]
    n_SduiNormalizer["SduiNormalizer (unit)"]
    n_SduiParser["SduiParser (unit)"]
    n_SduiSourceModel["SduiSourceModel (capability)"]
    n_SduiValidator["SduiValidator (unit)"]
    n_TokenizeSduiSource["TokenizeSduiSource (functionality)"]
    n_ValidateRelativeFormatting["ValidateRelativeFormatting (functionality)"]
    n_ValidateSymbolicBindings["ValidateSymbolicBindings (functionality)"]
    n_ValidateWidgetArguments["ValidateWidgetArguments (functionality)"]
    n_BuildSduiAst -->|realizes| n_SduiSourceModel
    n_CoordinateSduiCompilation -->|realizes| n_SduiSourceModel
    n_ExpandUiDefinitions -->|realizes| n_SduiSourceModel
    n_PreserveUiRegions -->|realizes| n_SduiSourceModel
    n_PreserveUiSourceMap -->|realizes| n_SduiSourceModel
    n_ResolveUiNames -->|realizes| n_SduiSourceModel
    n_SduiFrontend -->|owns| n_CoordinateSduiCompilation
    n_SduiLexer -->|owns| n_TokenizeSduiSource
    n_SduiNormalizer -->|owns| n_ExpandUiDefinitions
    n_SduiNormalizer -->|owns| n_PreserveUiRegions
    n_SduiNormalizer -->|owns| n_PreserveUiSourceMap
    n_SduiParser -->|owns| n_BuildSduiAst
    n_SduiValidator -->|owns| n_ResolveUiNames
    n_SduiValidator -->|owns| n_ValidateRelativeFormatting
    n_SduiValidator -->|owns| n_ValidateSymbolicBindings
    n_SduiValidator -->|owns| n_ValidateWidgetArguments
    n_TokenizeSduiSource -->|realizes| n_SduiSourceModel
    n_ValidateRelativeFormatting -->|realizes| n_SduiSourceModel
    n_ValidateSymbolicBindings -->|realizes| n_SduiSourceModel
    n_ValidateWidgetArguments -->|realizes| n_SduiSourceModel
```

Kildegrunnlag: f0023, f0054, f0108, f0216, f0219, f0275, f0366, f0383, f0393, f0394, f0395, f0397, f0419, f0420, f0421, f0422, f0452, f0487, f0498, f0505.

### Tilbydere av kapabilitet: SduiSourceModel

```mermaid
flowchart LR
    n_SduiFrontend["SduiFrontend (unit)"]
    n_SduiLexer["SduiLexer (unit)"]
    n_SduiLibrary["SduiLibrary (unit)"]
    n_SduiNormalizer["SduiNormalizer (unit)"]
    n_SduiParser["SduiParser (unit)"]
    n_SduiSourceModel["SduiSourceModel (capability)"]
    n_SduiValidator["SduiValidator (unit)"]
    n_SduiFrontend -->|provides| n_SduiSourceModel
    n_SduiLexer -->|provides| n_SduiSourceModel
    n_SduiLibrary -->|provides| n_SduiSourceModel
    n_SduiNormalizer -->|provides| n_SduiSourceModel
    n_SduiParser -->|provides| n_SduiSourceModel
    n_SduiValidator -->|provides| n_SduiSourceModel
```

Kildegrunnlag: f0367, f0384, f0391, f0396, f0398, f0423.

### Bidrag til kapabilitet: SourceDiagnostics

```mermaid
flowchart LR
    n_DiagnosticReporter["DiagnosticReporter (unit)"]
    n_ReportBindingDiagnostics["ReportBindingDiagnostics (functionality)"]
    n_ReportSourceDiagnostics["ReportSourceDiagnostics (functionality)"]
    n_SourceDiagnostics["SourceDiagnostics (capability)"]
    n_DiagnosticReporter -->|owns| n_ReportBindingDiagnostics
    n_DiagnosticReporter -->|owns| n_ReportSourceDiagnostics
    n_ReportBindingDiagnostics -->|realizes| n_SourceDiagnostics
    n_ReportSourceDiagnostics -->|realizes| n_SourceDiagnostics
```

Kildegrunnlag: f0089, f0090, f0259, f0264.

### Tilbydere av kapabilitet: SourceDiagnostics

```mermaid
flowchart LR
    n_DiagnosticReporter["DiagnosticReporter (unit)"]
    n_SourceDiagnostics["SourceDiagnostics (capability)"]
    n_DiagnosticReporter -->|provides| n_SourceDiagnostics
```

Kildegrunnlag: f0091.

### Bidrag til kapabilitet: SourceLoading

```mermaid
flowchart LR
    n_IdentifySourceRevision["IdentifySourceRevision (functionality)"]
    n_ReadBoundedSources["ReadBoundedSources (functionality)"]
    n_SourceLoader["SourceLoader (unit)"]
    n_SourceLoading["SourceLoading (capability)"]
    n_IdentifySourceRevision -->|realizes| n_SourceLoading
    n_ReadBoundedSources -->|realizes| n_SourceLoading
    n_SourceLoader -->|owns| n_IdentifySourceRevision
    n_SourceLoader -->|owns| n_ReadBoundedSources
```

Kildegrunnlag: f0158, f0235, f0436, f0437.

### Tilbydere av kapabilitet: SourceLoading

```mermaid
flowchart LR
    n_SourceLoader["SourceLoader (unit)"]
    n_SourceLoading["SourceLoading (capability)"]
    n_SourceLoader -->|provides| n_SourceLoading
```

Kildegrunnlag: f0438.

### Bidrag til kapabilitet: StaticDocumentation

```mermaid
flowchart LR
    n_CommandLineHost["CommandLineHost (container)"]
    n_ComposeHeadlessExport["ComposeHeadlessExport (functionality)"]
    n_ComposeMarkdownDocument["ComposeMarkdownDocument (functionality)"]
    n_ExportConsoleSnapshot["ExportConsoleSnapshot (functionality)"]
    n_ExportSvgSnapshot["ExportSvgSnapshot (functionality)"]
    n_SduiPresentation["SduiPresentation (unit)"]
    n_StaticDocumentation["StaticDocumentation (capability)"]
    n_WriteGeneratedArtifacts["WriteGeneratedArtifacts (functionality)"]
    n_CommandLineHost -->|owns| n_ComposeHeadlessExport
    n_CommandLineHost -->|owns| n_WriteGeneratedArtifacts
    n_ComposeHeadlessExport -->|realizes| n_StaticDocumentation
    n_ComposeMarkdownDocument -->|realizes| n_StaticDocumentation
    n_ExportConsoleSnapshot -->|realizes| n_StaticDocumentation
    n_ExportSvgSnapshot -->|realizes| n_StaticDocumentation
    n_SduiPresentation -->|owns| n_ComposeMarkdownDocument
    n_SduiPresentation -->|owns| n_ExportConsoleSnapshot
    n_SduiPresentation -->|owns| n_ExportSvgSnapshot
    n_WriteGeneratedArtifacts -->|realizes| n_StaticDocumentation
```

Kildegrunnlag: f0035, f0036, f0038, f0044, f0109, f0112, f0400, f0401, f0402, f0508.

### Tilbydere av kapabilitet: StaticDocumentation

```mermaid
flowchart LR
    n_CommandLineHost["CommandLineHost (container)"]
    n_SduiLibrary["SduiLibrary (unit)"]
    n_SduiPresentation["SduiPresentation (unit)"]
    n_StaticDocumentation["StaticDocumentation (capability)"]
    n_CommandLineHost -->|provides| n_StaticDocumentation
    n_SduiLibrary -->|provides| n_StaticDocumentation
    n_SduiPresentation -->|provides| n_StaticDocumentation
```

Kildegrunnlag: f0037, f0392, f0403.


## VP05 — Avhengigheter per modus

### Nødvendige porter i modus: BoundExecution

```mermaid
flowchart LR
    n_BoundInteraction["BoundInteraction (capability)"]
    n_DomainFunctionPort["DomainFunctionPort (interface)"]
    n_ExecutableDesign["ExecutableDesign (capability)"]
    n_SdlExecutionPort["SdlExecutionPort (interface)"]
    n_UiSessionPort["UiSessionPort (interface)"]
    n_BoundInteraction -->|requires| n_SdlExecutionPort
    n_BoundInteraction -->|requires| n_UiSessionPort
    n_ExecutableDesign -->|requires| n_DomainFunctionPort
```

Kildegrunnlag: f0005, f0006, f0105.

### Nødvendige porter i modus: BoundLiveEditing

```mermaid
flowchart LR
    n_BindingReloadPort["BindingReloadPort (interface)"]
    n_DevelopmentReload["DevelopmentReload (capability)"]
    n_FileChangePort["FileChangePort (interface)"]
    n_SdlReloadPort["SdlReloadPort (interface)"]
    n_SourceSnapshotPort["SourceSnapshotPort (interface)"]
    n_UiReloadPort["UiReloadPort (interface)"]
    n_DevelopmentReload -->|requires| n_BindingReloadPort
    n_DevelopmentReload -->|requires| n_FileChangePort
    n_DevelopmentReload -->|requires| n_SdlReloadPort
    n_DevelopmentReload -->|requires| n_SourceSnapshotPort
    n_DevelopmentReload -->|requires| n_UiReloadPort
```

Kildegrunnlag: f0074, f0075, f0077, f0078, f0080.

### Nødvendige porter i modus: LiveEditing

```mermaid
flowchart LR
    n_DevelopmentReload["DevelopmentReload (capability)"]
    n_FileChangePort["FileChangePort (interface)"]
    n_SourceSnapshotPort["SourceSnapshotPort (interface)"]
    n_UiReloadPort["UiReloadPort (interface)"]
    n_DevelopmentReload -->|requires| n_FileChangePort
    n_DevelopmentReload -->|requires| n_SourceSnapshotPort
    n_DevelopmentReload -->|requires| n_UiReloadPort
```

Kildegrunnlag: f0076, f0079, f0081.

### Nødvendige porter i modus: NativeBuild

```mermaid
flowchart LR
    n_BuildToolPort["BuildToolPort (interface)"]
    n_GeneratedArtifactPort["GeneratedArtifactPort (interface)"]
    n_NativeRealization["NativeRealization (capability)"]
    n_NativeRealization -->|requires| n_BuildToolPort
    n_NativeRealization -->|requires| n_GeneratedArtifactPort
```

Kildegrunnlag: f0187, f0188.

### Nødvendige porter i modus: RichDocument

```mermaid
flowchart LR
    n_ContentProviderPort["ContentProviderPort (interface)"]
    n_DiagramEnginePort["DiagramEnginePort (interface)"]
    n_RichContent["RichContent (capability)"]
    n_RichContent -->|requires| n_ContentProviderPort
    n_RichContent -->|requires| n_DiagramEnginePort
```

Kildegrunnlag: f0282, f0283.

### Nødvendige porter i modus: SourceInspection

```mermaid
flowchart LR
    n_SdlSourceModel["SdlSourceModel (capability)"]
    n_SduiSourceModel["SduiSourceModel (capability)"]
    n_SourceSnapshotPort["SourceSnapshotPort (interface)"]
    n_SdlSourceModel -->|requires| n_SourceSnapshotPort
    n_SduiSourceModel -->|requires| n_SourceSnapshotPort
```

Kildegrunnlag: f0333, f0418.

### Nødvendige porter i modus: StaticExport

```mermaid
flowchart LR
    n_ExportSinkPort["ExportSinkPort (interface)"]
    n_MeasuredPresentation["MeasuredPresentation (capability)"]
    n_MeasurementPort["MeasurementPort (interface)"]
    n_PreparedFramePort["PreparedFramePort (interface)"]
    n_StaticDocumentation["StaticDocumentation (capability)"]
    n_MeasuredPresentation -->|requires| n_MeasurementPort
    n_StaticDocumentation -->|requires| n_ExportSinkPort
    n_StaticDocumentation -->|requires| n_PreparedFramePort
```

Kildegrunnlag: f0182, f0443, f0444.

### Nødvendige porter i modus: UiPreview

```mermaid
flowchart LR
    n_MeasuredPresentation["MeasuredPresentation (capability)"]
    n_MeasurementPort["MeasurementPort (interface)"]
    n_NativeInteraction["NativeInteraction (capability)"]
    n_WidgetBackendPort["WidgetBackendPort (interface)"]
    n_MeasuredPresentation -->|requires| n_MeasurementPort
    n_NativeInteraction -->|requires| n_WidgetBackendPort
```

Kildegrunnlag: f0183, f0186.


## VP06 — Aktivitetsdetaljering

### Detaljert aktivitet → overordnet aktivitet

```mermaid
flowchart LR
    n_BuildNativeRealization["BuildNativeRealization (activity)"]
    n_ExportUiDocumentation["ExportUiDocumentation (activity)"]
    n_InspectDesignSource["InspectDesignSource (activity)"]
    n_InspectSdlSource["InspectSdlSource (activity)"]
    n_InspectSduiSource["InspectSduiSource (activity)"]
    n_RealizeDesign["RealizeDesign (activity)"]
    n_ReloadBoundModels["ReloadBoundModels (activity)"]
    n_ReloadDesignSession["ReloadDesignSession (activity)"]
    n_ReloadUiModel["ReloadUiModel (activity)"]
    n_RunBoundUiAction["RunBoundUiAction (activity)"]
    n_RunDesignSession["RunDesignSession (activity)"]
    n_RunUnboundUiPreview["RunUnboundUiPreview (activity)"]
    n_BuildNativeRealization -->|refines| n_RealizeDesign
    n_ExportUiDocumentation -->|refines| n_InspectDesignSource
    n_InspectSdlSource -->|refines| n_InspectDesignSource
    n_InspectSduiSource -->|refines| n_InspectDesignSource
    n_ReloadBoundModels -->|refines| n_ReloadDesignSession
    n_ReloadUiModel -->|refines| n_ReloadDesignSession
    n_RunBoundUiAction -->|refines| n_RunDesignSession
    n_RunUnboundUiPreview -->|refines| n_RunDesignSession
```

Kildegrunnlag: f0010, f0113, f0159, f0160, f0245, f0258, f0285, f0286.


## VP07 — Features over arkitekturen

### Feature: DesignDocumentation — modus SourceInspection

```mermaid
flowchart LR
    n_CommandLineHost["CommandLineHost (container)"]
    n_ComposeMarkdownDocument["ComposeMarkdownDocument (functionality)"]
    n_DesignDocumentation["DesignDocumentation (feature)"]
    n_ExportSvgSnapshot["ExportSvgSnapshot (functionality)"]
    n_ExportViewpointMarkdown["ExportViewpointMarkdown (functionality)"]
    n_ProjectSdlViewpoints["ProjectSdlViewpoints (functionality)"]
    n_SdlViewpointGenerator["SdlViewpointGenerator (unit)"]
    n_SduiPresentation["SduiPresentation (unit)"]
    n_TraceViewpointFacts["TraceViewpointFacts (functionality)"]
    n_WriteGeneratedArtifacts["WriteGeneratedArtifacts (functionality)"]
    n_CommandLineHost -->|owns| n_WriteGeneratedArtifacts
    n_ComposeMarkdownDocument -->|contributes-to| n_DesignDocumentation
    n_ExportSvgSnapshot -->|contributes-to| n_DesignDocumentation
    n_ExportViewpointMarkdown -->|contributes-to| n_DesignDocumentation
    n_ProjectSdlViewpoints -->|contributes-to| n_DesignDocumentation
    n_SdlViewpointGenerator -->|owns| n_ExportViewpointMarkdown
    n_SdlViewpointGenerator -->|owns| n_ProjectSdlViewpoints
    n_SdlViewpointGenerator -->|owns| n_TraceViewpointFacts
    n_SduiPresentation -->|owns| n_ComposeMarkdownDocument
    n_SduiPresentation -->|owns| n_ExportSvgSnapshot
    n_TraceViewpointFacts -->|allocated-to| n_CommandLineHost
    n_TraceViewpointFacts -->|contributes-to| n_DesignDocumentation
    n_WriteGeneratedArtifacts -->|contributes-to| n_DesignDocumentation
```

Kildegrunnlag: f0036, f0043, f0111, f0115, f0221, f0350, f0351, f0352, f0400, f0402, f0453, f0455, f0507.

### Feature: DesignDocumentation — modus StaticExport

```mermaid
flowchart LR
    n_CommandLineHost["CommandLineHost (container)"]
    n_ComposeMarkdownDocument["ComposeMarkdownDocument (functionality)"]
    n_DesignDocumentation["DesignDocumentation (feature)"]
    n_ExportSvgSnapshot["ExportSvgSnapshot (functionality)"]
    n_ExportViewpointMarkdown["ExportViewpointMarkdown (functionality)"]
    n_ProjectSdlViewpoints["ProjectSdlViewpoints (functionality)"]
    n_SdlViewpointGenerator["SdlViewpointGenerator (unit)"]
    n_SduiPresentation["SduiPresentation (unit)"]
    n_TraceViewpointFacts["TraceViewpointFacts (functionality)"]
    n_WriteGeneratedArtifacts["WriteGeneratedArtifacts (functionality)"]
    n_CommandLineHost -->|owns| n_WriteGeneratedArtifacts
    n_ComposeMarkdownDocument -->|allocated-to| n_CommandLineHost
    n_ComposeMarkdownDocument -->|contributes-to| n_DesignDocumentation
    n_ExportSvgSnapshot -->|allocated-to| n_CommandLineHost
    n_ExportSvgSnapshot -->|contributes-to| n_DesignDocumentation
    n_ExportViewpointMarkdown -->|allocated-to| n_CommandLineHost
    n_ExportViewpointMarkdown -->|contributes-to| n_DesignDocumentation
    n_ProjectSdlViewpoints -->|allocated-to| n_CommandLineHost
    n_ProjectSdlViewpoints -->|contributes-to| n_DesignDocumentation
    n_SdlViewpointGenerator -->|owns| n_ExportViewpointMarkdown
    n_SdlViewpointGenerator -->|owns| n_ProjectSdlViewpoints
    n_SdlViewpointGenerator -->|owns| n_TraceViewpointFacts
    n_SduiPresentation -->|owns| n_ComposeMarkdownDocument
    n_SduiPresentation -->|owns| n_ExportSvgSnapshot
    n_TraceViewpointFacts -->|allocated-to| n_CommandLineHost
    n_TraceViewpointFacts -->|contributes-to| n_DesignDocumentation
    n_WriteGeneratedArtifacts -->|allocated-to| n_CommandLineHost
    n_WriteGeneratedArtifacts -->|contributes-to| n_DesignDocumentation
```

Kildegrunnlag: f0036, f0042, f0043, f0110, f0111, f0114, f0115, f0220, f0221, f0350, f0351, f0352, f0400, f0402, f0454, f0455, f0506, f0507.

### Feature: InteractiveUiPreview — modus UiPreview

```mermaid
flowchart LR
    n_AllocateGeometry["AllocateGeometry (functionality)"]
    n_BuildPreparedFrame["BuildPreparedFrame (functionality)"]
    n_ComposeInteractiveSession["ComposeInteractiveSession (functionality)"]
    n_FyneBackend["FyneBackend (unit)"]
    n_FyneHost["FyneHost (container)"]
    n_HandleFocusAndTextInput["HandleFocusAndTextInput (functionality)"]
    n_InteractiveUiPreview["InteractiveUiPreview (feature)"]
    n_ReconcileWidgets["ReconcileWidgets (functionality)"]
    n_SduiLayout["SduiLayout (unit)"]
    n_AllocateGeometry -->|allocated-to| n_FyneHost
    n_AllocateGeometry -->|contributes-to| n_InteractiveUiPreview
    n_BuildPreparedFrame -->|allocated-to| n_FyneHost
    n_BuildPreparedFrame -->|contributes-to| n_InteractiveUiPreview
    n_ComposeInteractiveSession -->|allocated-to| n_FyneHost
    n_ComposeInteractiveSession -->|contributes-to| n_InteractiveUiPreview
    n_FyneBackend -->|owns| n_HandleFocusAndTextInput
    n_FyneBackend -->|owns| n_ReconcileWidgets
    n_FyneHost -->|owns| n_ComposeInteractiveSession
    n_HandleFocusAndTextInput -->|allocated-to| n_FyneHost
    n_HandleFocusAndTextInput -->|contributes-to| n_InteractiveUiPreview
    n_ReconcileWidgets -->|allocated-to| n_FyneHost
    n_ReconcileWidgets -->|contributes-to| n_InteractiveUiPreview
    n_SduiLayout -->|owns| n_AllocateGeometry
    n_SduiLayout -->|owns| n_BuildPreparedFrame
```

Kildegrunnlag: f0000, f0001, f0011, f0012, f0039, f0040, f0120, f0122, f0132, f0155, f0156, f0236, f0237, f0377, f0378.

### Feature: LiveModelReload — modus LiveEditing

```mermaid
flowchart LR
    n_DiagnosticReporter["DiagnosticReporter (unit)"]
    n_FyneHost["FyneHost (container)"]
    n_KeepLastValidModels["KeepLastValidModels (functionality)"]
    n_LiveModelReload["LiveModelReload (feature)"]
    n_ObserveSourceChanges["ObserveSourceChanges (functionality)"]
    n_PrepareCandidateModels["PrepareCandidateModels (functionality)"]
    n_PreserveCompatibleUiState["PreserveCompatibleUiState (functionality)"]
    n_ProjectUiGeneration["ProjectUiGeneration (functionality)"]
    n_PublishModelGeneration["PublishModelGeneration (functionality)"]
    n_ReloadCoordinator["ReloadCoordinator (unit)"]
    n_ReportSourceDiagnostics["ReportSourceDiagnostics (functionality)"]
    n_SduiInstanceStore["SduiInstanceStore (unit)"]
    n_SourceWatcher["SourceWatcher (unit)"]
    n_UiStateReconciler["UiStateReconciler (unit)"]
    n_DiagnosticReporter -->|owns| n_ReportSourceDiagnostics
    n_KeepLastValidModels -->|allocated-to| n_FyneHost
    n_KeepLastValidModels -->|contributes-to| n_LiveModelReload
    n_ObserveSourceChanges -->|allocated-to| n_FyneHost
    n_ObserveSourceChanges -->|contributes-to| n_LiveModelReload
    n_PrepareCandidateModels -->|allocated-to| n_FyneHost
    n_PrepareCandidateModels -->|contributes-to| n_LiveModelReload
    n_PreserveCompatibleUiState -->|allocated-to| n_FyneHost
    n_PreserveCompatibleUiState -->|contributes-to| n_LiveModelReload
    n_ProjectUiGeneration -->|allocated-to| n_FyneHost
    n_ProjectUiGeneration -->|contributes-to| n_LiveModelReload
    n_PublishModelGeneration -->|allocated-to| n_FyneHost
    n_PublishModelGeneration -->|contributes-to| n_LiveModelReload
    n_ReloadCoordinator -->|owns| n_KeepLastValidModels
    n_ReloadCoordinator -->|owns| n_PrepareCandidateModels
    n_ReloadCoordinator -->|owns| n_PublishModelGeneration
    n_ReportSourceDiagnostics -->|allocated-to| n_FyneHost
    n_ReportSourceDiagnostics -->|contributes-to| n_LiveModelReload
    n_SduiInstanceStore -->|owns| n_ProjectUiGeneration
    n_SourceWatcher -->|owns| n_ObserveSourceChanges
    n_UiStateReconciler -->|owns| n_PreserveCompatibleUiState
```

Kildegrunnlag: f0090, f0165, f0166, f0196, f0197, f0200, f0201, f0205, f0206, f0224, f0225, f0230, f0231, f0253, f0254, f0255, f0261, f0262, f0370, f0441, f0480.

### Feature: LiveModelReload — modus SourceInspection

```mermaid
flowchart LR
    n_CommandLineHost["CommandLineHost (container)"]
    n_DiagnosticReporter["DiagnosticReporter (unit)"]
    n_KeepLastValidModels["KeepLastValidModels (functionality)"]
    n_LiveModelReload["LiveModelReload (feature)"]
    n_ObserveSourceChanges["ObserveSourceChanges (functionality)"]
    n_PrepareCandidateModels["PrepareCandidateModels (functionality)"]
    n_PreserveCompatibleUiState["PreserveCompatibleUiState (functionality)"]
    n_ProjectUiGeneration["ProjectUiGeneration (functionality)"]
    n_PublishModelGeneration["PublishModelGeneration (functionality)"]
    n_ReloadCoordinator["ReloadCoordinator (unit)"]
    n_ReportSourceDiagnostics["ReportSourceDiagnostics (functionality)"]
    n_SduiInstanceStore["SduiInstanceStore (unit)"]
    n_SourceWatcher["SourceWatcher (unit)"]
    n_UiStateReconciler["UiStateReconciler (unit)"]
    n_DiagnosticReporter -->|owns| n_ReportSourceDiagnostics
    n_KeepLastValidModels -->|contributes-to| n_LiveModelReload
    n_ObserveSourceChanges -->|contributes-to| n_LiveModelReload
    n_PrepareCandidateModels -->|contributes-to| n_LiveModelReload
    n_PreserveCompatibleUiState -->|contributes-to| n_LiveModelReload
    n_ProjectUiGeneration -->|contributes-to| n_LiveModelReload
    n_PublishModelGeneration -->|contributes-to| n_LiveModelReload
    n_ReloadCoordinator -->|owns| n_KeepLastValidModels
    n_ReloadCoordinator -->|owns| n_PrepareCandidateModels
    n_ReloadCoordinator -->|owns| n_PublishModelGeneration
    n_ReportSourceDiagnostics -->|allocated-to| n_CommandLineHost
    n_ReportSourceDiagnostics -->|contributes-to| n_LiveModelReload
    n_SduiInstanceStore -->|owns| n_ProjectUiGeneration
    n_SourceWatcher -->|owns| n_ObserveSourceChanges
    n_UiStateReconciler -->|owns| n_PreserveCompatibleUiState
```

Kildegrunnlag: f0090, f0166, f0197, f0201, f0206, f0225, f0231, f0253, f0254, f0255, f0260, f0262, f0370, f0441, f0480.

### Feature: NativeGoAssembly — modus NativeBuild

```mermaid
flowchart LR
    n_BuildGeneratedApplication["BuildGeneratedApplication (functionality)"]
    n_CommandLineHost["CommandLineHost (container)"]
    n_GenerateBindingRegistration["GenerateBindingRegistration (functionality)"]
    n_GenerateModelConstructors["GenerateModelConstructors (functionality)"]
    n_GoBuildRunner["GoBuildRunner (unit)"]
    n_GoCodeGenerator["GoCodeGenerator (unit)"]
    n_NativeGoAssembly["NativeGoAssembly (feature)"]
    n_PreserveHandwrittenSources["PreserveHandwrittenSources (functionality)"]
    n_BuildGeneratedApplication -->|allocated-to| n_CommandLineHost
    n_BuildGeneratedApplication -->|contributes-to| n_NativeGoAssembly
    n_GenerateBindingRegistration -->|allocated-to| n_CommandLineHost
    n_GenerateBindingRegistration -->|contributes-to| n_NativeGoAssembly
    n_GenerateModelConstructors -->|allocated-to| n_CommandLineHost
    n_GenerateModelConstructors -->|contributes-to| n_NativeGoAssembly
    n_GoBuildRunner -->|owns| n_BuildGeneratedApplication
    n_GoCodeGenerator -->|owns| n_GenerateBindingRegistration
    n_GoCodeGenerator -->|owns| n_GenerateModelConstructors
    n_GoCodeGenerator -->|owns| n_PreserveHandwrittenSources
    n_PreserveHandwrittenSources -->|allocated-to| n_CommandLineHost
    n_PreserveHandwrittenSources -->|contributes-to| n_NativeGoAssembly
```

Kildegrunnlag: f0007, f0008, f0135, f0136, f0138, f0139, f0143, f0149, f0150, f0151, f0208, f0209.

### Feature: StructuralModelInspection — modus LiveEditing

```mermaid
flowchart LR
    n_BuildSdlAst["BuildSdlAst (functionality)"]
    n_BuildSduiAst["BuildSduiAst (functionality)"]
    n_DiagnosticReporter["DiagnosticReporter (unit)"]
    n_FyneHost["FyneHost (container)"]
    n_ReportSourceDiagnostics["ReportSourceDiagnostics (functionality)"]
    n_SdlParser["SdlParser (unit)"]
    n_SdlValidator["SdlValidator (unit)"]
    n_SduiParser["SduiParser (unit)"]
    n_SduiValidator["SduiValidator (unit)"]
    n_StructuralModelInspection["StructuralModelInspection (feature)"]
    n_ValidateSdlStructure["ValidateSdlStructure (functionality)"]
    n_ValidateWidgetArguments["ValidateWidgetArguments (functionality)"]
    n_BuildSdlAst -->|contributes-to| n_StructuralModelInspection
    n_BuildSduiAst -->|contributes-to| n_StructuralModelInspection
    n_DiagnosticReporter -->|owns| n_ReportSourceDiagnostics
    n_ReportSourceDiagnostics -->|allocated-to| n_FyneHost
    n_ReportSourceDiagnostics -->|contributes-to| n_StructuralModelInspection
    n_SdlParser -->|owns| n_BuildSdlAst
    n_SdlValidator -->|owns| n_ValidateSdlStructure
    n_SduiParser -->|owns| n_BuildSduiAst
    n_SduiValidator -->|owns| n_ValidateWidgetArguments
    n_ValidateSdlStructure -->|contributes-to| n_StructuralModelInspection
    n_ValidateWidgetArguments -->|contributes-to| n_StructuralModelInspection
```

Kildegrunnlag: f0015, f0020, f0090, f0261, f0263, f0321, f0348, f0397, f0422, f0492, f0502.

### Feature: StructuralModelInspection — modus SourceInspection

```mermaid
flowchart LR
    n_BuildSdlAst["BuildSdlAst (functionality)"]
    n_BuildSduiAst["BuildSduiAst (functionality)"]
    n_CommandLineHost["CommandLineHost (container)"]
    n_DiagnosticReporter["DiagnosticReporter (unit)"]
    n_ReportSourceDiagnostics["ReportSourceDiagnostics (functionality)"]
    n_SdlParser["SdlParser (unit)"]
    n_SdlValidator["SdlValidator (unit)"]
    n_SduiParser["SduiParser (unit)"]
    n_SduiValidator["SduiValidator (unit)"]
    n_StructuralModelInspection["StructuralModelInspection (feature)"]
    n_ValidateSdlStructure["ValidateSdlStructure (functionality)"]
    n_ValidateWidgetArguments["ValidateWidgetArguments (functionality)"]
    n_BuildSdlAst -->|allocated-to| n_CommandLineHost
    n_BuildSdlAst -->|contributes-to| n_StructuralModelInspection
    n_BuildSduiAst -->|allocated-to| n_CommandLineHost
    n_BuildSduiAst -->|contributes-to| n_StructuralModelInspection
    n_DiagnosticReporter -->|owns| n_ReportSourceDiagnostics
    n_ReportSourceDiagnostics -->|allocated-to| n_CommandLineHost
    n_ReportSourceDiagnostics -->|contributes-to| n_StructuralModelInspection
    n_SdlParser -->|owns| n_BuildSdlAst
    n_SdlValidator -->|owns| n_ValidateSdlStructure
    n_SduiParser -->|owns| n_BuildSduiAst
    n_SduiValidator -->|owns| n_ValidateWidgetArguments
    n_ValidateSdlStructure -->|allocated-to| n_CommandLineHost
    n_ValidateSdlStructure -->|contributes-to| n_StructuralModelInspection
    n_ValidateWidgetArguments -->|allocated-to| n_CommandLineHost
    n_ValidateWidgetArguments -->|contributes-to| n_StructuralModelInspection
```

Kildegrunnlag: f0014, f0015, f0019, f0020, f0090, f0260, f0263, f0321, f0348, f0397, f0422, f0491, f0492, f0501, f0502.

### Feature: TypedDomainBinding — modus BoundExecution

```mermaid
flowchart LR
    n_ConnectTypedWidgetHandles["ConnectTypedWidgetHandles (functionality)"]
    n_DispatchUiEvent["DispatchUiEvent (functionality)"]
    n_FyneHost["FyneHost (container)"]
    n_InvokeRegisteredFunction["InvokeRegisteredFunction (functionality)"]
    n_PublishDomainUpdates["PublishDomainUpdates (functionality)"]
    n_ResolveCallbackSymbols["ResolveCallbackSymbols (functionality)"]
    n_SdlDispatcher["SdlDispatcher (unit)"]
    n_SdlUiBindingAdapter["SdlUiBindingAdapter (unit)"]
    n_SduiDispatcher["SduiDispatcher (unit)"]
    n_TypedDomainBinding["TypedDomainBinding (feature)"]
    n_ConnectTypedWidgetHandles -->|allocated-to| n_FyneHost
    n_ConnectTypedWidgetHandles -->|contributes-to| n_TypedDomainBinding
    n_DispatchUiEvent -->|allocated-to| n_FyneHost
    n_DispatchUiEvent -->|contributes-to| n_TypedDomainBinding
    n_InvokeRegisteredFunction -->|allocated-to| n_FyneHost
    n_InvokeRegisteredFunction -->|contributes-to| n_TypedDomainBinding
    n_PublishDomainUpdates -->|allocated-to| n_FyneHost
    n_PublishDomainUpdates -->|contributes-to| n_TypedDomainBinding
    n_ResolveCallbackSymbols -->|allocated-to| n_FyneHost
    n_ResolveCallbackSymbols -->|contributes-to| n_TypedDomainBinding
    n_SdlDispatcher -->|owns| n_InvokeRegisteredFunction
    n_SdlUiBindingAdapter -->|owns| n_ConnectTypedWidgetHandles
    n_SdlUiBindingAdapter -->|owns| n_PublishDomainUpdates
    n_SdlUiBindingAdapter -->|owns| n_ResolveCallbackSymbols
    n_SduiDispatcher -->|owns| n_DispatchUiEvent
```

Kildegrunnlag: f0046, f0047, f0097, f0098, f0162, f0163, f0227, f0228, f0267, f0268, f0292, f0340, f0342, f0343, f0356.

### Modellhull i dette utsnittet

| Identitet | Modus | Mangel |
| --- | --- | --- |
| ComposeMarkdownDocument | SourceInspection | Container-allokering er uspesifisert for bidrag til DesignDocumentation. |
| ExportSvgSnapshot | SourceInspection | Container-allokering er uspesifisert for bidrag til DesignDocumentation. |
| ExportViewpointMarkdown | SourceInspection | Container-allokering er uspesifisert for bidrag til DesignDocumentation. |
| ProjectSdlViewpoints | SourceInspection | Container-allokering er uspesifisert for bidrag til DesignDocumentation. |
| WriteGeneratedArtifacts | SourceInspection | Container-allokering er uspesifisert for bidrag til DesignDocumentation. |
| KeepLastValidModels | SourceInspection | Container-allokering er uspesifisert for bidrag til LiveModelReload. |
| ObserveSourceChanges | SourceInspection | Container-allokering er uspesifisert for bidrag til LiveModelReload. |
| PrepareCandidateModels | SourceInspection | Container-allokering er uspesifisert for bidrag til LiveModelReload. |
| PreserveCompatibleUiState | SourceInspection | Container-allokering er uspesifisert for bidrag til LiveModelReload. |
| ProjectUiGeneration | SourceInspection | Container-allokering er uspesifisert for bidrag til LiveModelReload. |
| PublishModelGeneration | SourceInspection | Container-allokering er uspesifisert for bidrag til LiveModelReload. |
| BuildSdlAst | LiveEditing | Container-allokering er uspesifisert for bidrag til StructuralModelInspection. |
| BuildSduiAst | LiveEditing | Container-allokering er uspesifisert for bidrag til StructuralModelInspection. |
| ValidateSdlStructure | LiveEditing | Container-allokering er uspesifisert for bidrag til StructuralModelInspection. |
| ValidateWidgetArguments | LiveEditing | Container-allokering er uspesifisert for bidrag til StructuralModelInspection. |


## VP09 — Dataset, Datagram og persistent Database

### Dataopprinnelse og holder: DesignSourceDocuments

```mermaid
flowchart LR
    n_DesignSourceArchive["DesignSourceArchive (database)"]
    n_DesignSourceDocuments["DesignSourceDocuments (dataset)"]
    n_DesignSourceRecord["DesignSourceRecord (contract)"]
    n_SourceLoader["SourceLoader (unit)"]
    n_DesignSourceArchive -->|holds| n_DesignSourceDocuments
    n_DesignSourceDocuments -->|upholds| n_DesignSourceRecord
    n_SourceLoader -->|owns| n_DesignSourceArchive
```

Kildegrunnlag: f0069, f0070, f0435.

### Dataopprinnelse og holder: UiSessionState

```mermaid
flowchart LR
    n_SduiInstanceStore["SduiInstanceStore (unit)"]
    n_UiGenerationContract["UiGenerationContract (contract)"]
    n_UiGenerationNotices["UiGenerationNotices (datagram)"]
    n_UiSessionRecord["UiSessionRecord (contract)"]
    n_UiSessionState["UiSessionState (dataset)"]
    n_SduiInstanceStore -->|holds| n_UiSessionState
    n_UiGenerationNotices -->|from| n_UiSessionState
    n_UiGenerationNotices -->|upholds| n_UiGenerationContract
    n_UiSessionState -->|upholds| n_UiSessionRecord
```

Kildegrunnlag: f0368, f0466, f0467, f0476.

### Kontraktstruktur: DesignSourceRecord

```mermaid
flowchart LR
    n_DesignSourceRecord["DesignSourceRecord (contract)"]
    n_SourceDocumentRevision["SourceDocumentRevision (field)"]
    n_SourceDocumentText["SourceDocumentText (field)"]
    n_DesignSourceRecord -->|has-field| n_SourceDocumentRevision
    n_DesignSourceRecord -->|has-field| n_SourceDocumentText
```

Kildegrunnlag: f0072, f0073.

### Kontraktstruktur: UiGenerationContract

```mermaid
flowchart LR
    n_NoticeGeneration["NoticeGeneration (field)"]
    n_NoticeVersion["NoticeVersion (field)"]
    n_UiGenerationChanged["UiGenerationChanged (variant)"]
    n_UiGenerationContract["UiGenerationContract (contract)"]
    n_UiGenerationChanged -->|has-field| n_NoticeGeneration
    n_UiGenerationContract -->|defines| n_UiGenerationChanged
    n_UiGenerationContract -->|has-field| n_NoticeVersion
```

Kildegrunnlag: f0462, f0463, f0465.

### Kontraktstruktur: UiSessionRecord

```mermaid
flowchart LR
    n_SessionDraft["SessionDraft (field)"]
    n_SessionGeneration["SessionGeneration (field)"]
    n_UiSessionRecord["UiSessionRecord (contract)"]
    n_UiSessionRecord -->|has-field| n_SessionDraft
    n_UiSessionRecord -->|has-field| n_SessionGeneration
```

Kildegrunnlag: f0474, f0475.


## VP10 — Datagram-koding og packet

### Packet: UiGenerationWire / UiGenerationChanged — big-endian, most-significant-first

```mermaid
packet
    0-15: "NoticeVersion"
    16-79: "NoticeGeneration"
```

Kildegrunnlag: f0192, f0193, f0194, f0195, f0462, f0463, f0464, f0465, f0468, f0469, f0470, f0471, f0472.


### VP09 — felt og kontraktegenskaper

| Modellfaktum | Kilde-ID |
| --- | --- |
| DesignSourceRecord has completeness = closed. | f0071 |
| NoticeGeneration has presence = required. | f0192 |
| NoticeGeneration has value-type = unsigned. | f0193 |
| NoticeVersion has presence = required. | f0194 |
| NoticeVersion has value-type = unsigned. | f0195 |
| SessionDraft has presence = optional. | f0424 |
| SessionDraft has value-type = text. | f0425 |
| SessionGeneration has presence = required. | f0426 |
| SessionGeneration has value-type = unsigned. | f0427 |
| SourceDocumentRevision has presence = required. | f0430 |
| SourceDocumentRevision has value-type = unsigned. | f0431 |
| SourceDocumentText has presence = required. | f0432 |
| SourceDocumentText has value-type = text. | f0433 |
| UiGenerationContract has completeness = closed. | f0464 |
| UiGenerationWire has bit-order = most-significant-first. | f0469 |
| UiGenerationWire has byte-order = big-endian. | f0470 |
| UiSessionRecord has completeness = closed. | f0473 |

### VP09 — projeksjonsansvar

| Functionality | Dataset | Datagram-familie | Faktum |
| --- | --- | --- | --- |
| ProjectUiGeneration | UiSessionState | UiGenerationNotices | f0226 |

## VP04 — Grensesnittbruk

Tabellen dekker alle consumes-fakta. Portnavn er ikke Channel-kontrakter eller tilbyderkoblinger.

| Unit / Container | Interface | Faktum | Kildelinje |
| --- | --- | --- | --- |
| CommandLineHost | ExportSinkPort | f0031 | 260 |
| CommandLineHost | PreparedFramePort | f0032 | 261 |
| CommandLineHost | SduiFrontendPort | f0033 | 262 |
| CommandLineHost | SourceSnapshotPort | f0034 | 263 |
| DiagramProvider | DiagramEnginePort | f0092 | 321 |
| DiagramProvider | ResourcePort | f0093 | 322 |
| DomainStateMigrator | DomainStatePort | f0100 | 329 |
| DomainStateMigrator | SdlModelPort | f0101 | 330 |
| FyneBackend | PreparedFramePort | f0118 | 347 |
| FyneBackend | UiSessionPort | f0119 | 348 |
| FyneHost | DomainBindingPort | f0125 | 354 |
| FyneHost | ReloadPort | f0126 | 355 |
| FyneHost | SdlFrontendPort | f0127 | 356 |
| FyneHost | SduiFrontendPort | f0128 | 357 |
| FyneHost | SourceSnapshotPort | f0129 | 358 |
| FyneHost | UiSessionPort | f0130 | 359 |
| FyneHost | WidgetBackendPort | f0131 | 360 |
| GoBuildRunner | BuildToolPort | f0141 | 370 |
| GoBuildRunner | GeneratedArtifactPort | f0142 | 371 |
| GoCodeGenerator | ExecutionProfilePort | f0146 | 375 |
| GoCodeGenerator | SdlModelPort | f0147 | 376 |
| GoCodeGenerator | SduiModelPort | f0148 | 377 |
| MarkdownProvider | DiagramPort | f0173 | 402 |
| MarkdownProvider | MeasurementPort | f0174 | 403 |
| MarkdownProvider | ResourcePort | f0175 | 404 |
| ReloadCoordinator | BindingReloadPort | f0246 | 475 |
| ReloadCoordinator | DiagnosticPort | f0247 | 476 |
| ReloadCoordinator | SdlFrontendPort | f0248 | 477 |
| ReloadCoordinator | SdlReloadPort | f0249 | 478 |
| ReloadCoordinator | SduiFrontendPort | f0250 | 479 |
| ReloadCoordinator | SourceSnapshotPort | f0251 | 480 |
| ReloadCoordinator | UiReloadPort | f0252 | 481 |
| SdlDispatcher | DomainFunctionPort | f0288 | 517 |
| SdlDispatcher | DomainStatePort | f0289 | 518 |
| SdlExecutionGate | DiagnosticPort | f0295 | 524 |
| SdlExecutionGate | SdlModelPort | f0296 | 525 |
| SdlFrontend | DiagnosticPort | f0299 | 528 |
| SdlFrontend | SourceSnapshotPort | f0300 | 529 |
| SdlFunctionRegistry | DomainFunctionPort | f0307 | 536 |
| SdlRuntime | DomainFunctionPort | f0323 | 552 |
| SdlRuntime | SdlModelPort | f0324 | 553 |
| SdlUiBindingAdapter | DiagnosticPort | f0337 | 566 |
| SdlUiBindingAdapter | SdlExecutionPort | f0338 | 567 |
| SdlUiBindingAdapter | UiSessionPort | f0339 | 568 |
| SduiDispatcher | DomainBindingPort | f0353 | 582 |
| SduiDispatcher | UiStatePort | f0354 | 583 |
| SduiFrontend | DiagnosticPort | f0360 | 589 |
| SduiFrontend | SourceSnapshotPort | f0361 | 590 |
| SduiLayout | ContentProviderPort | f0374 | 603 |
| SduiLayout | MeasurementPort | f0375 | 604 |
| SduiLayout | UiSnapshotPort | f0376 | 605 |
| SduiPresentation | PreparedFramePort | f0399 | 628 |
| SduiPropertyStore | UiStatePort | f0404 | 633 |
| SduiRuntime | DomainBindingPort | f0409 | 638 |
| SduiRuntime | SduiModelPort | f0410 | 639 |
| SourceLoader | SourceInputPort | f0434 | 663 |
| SourceWatcher | FileChangePort | f0439 | 668 |
| UiStateReconciler | SduiModelPort | f0477 | 706 |
| UiStateReconciler | UiStatePort | f0478 | 707 |

## VP11 — Egenskaper og fullstendig faktaregister

Registeret inkluderer alle fakta, også de som ikke har en egen tegning.

| ID | Utsagn | Kildelinje |
| --- | --- | --- |
| f0000 | AllocateGeometry allocated-to FyneHost in mode UiPreview. | 229 |
| f0001 | AllocateGeometry contributes-to InteractiveUiPreview. | 230 |
| f0002 | AllocateGeometry realizes MeasuredPresentation. | 231 |
| f0003 | ApplyPropertyBatch has state-retention = stateful. | 232 |
| f0004 | ApplyPropertyBatch realizes InteractiveSession. | 233 |
| f0005 | BoundInteraction requires SdlExecutionPort in mode BoundExecution. | 234 |
| f0006 | BoundInteraction requires UiSessionPort in mode BoundExecution. | 235 |
| f0007 | BuildGeneratedApplication allocated-to CommandLineHost in mode NativeBuild. | 236 |
| f0008 | BuildGeneratedApplication contributes-to NativeGoAssembly. | 237 |
| f0009 | BuildGeneratedApplication realizes NativeRealization. | 238 |
| f0010 | BuildNativeRealization refines RealizeDesign. | 239 |
| f0011 | BuildPreparedFrame allocated-to FyneHost in mode UiPreview. | 240 |
| f0012 | BuildPreparedFrame contributes-to InteractiveUiPreview. | 241 |
| f0013 | BuildPreparedFrame realizes MeasuredPresentation. | 242 |
| f0014 | BuildSdlAst allocated-to CommandLineHost in mode SourceInspection. | 243 |
| f0015 | BuildSdlAst contributes-to StructuralModelInspection. | 244 |
| f0016 | BuildSdlAst has repeatability = deterministic. | 245 |
| f0017 | BuildSdlAst has state-retention = stateless. | 246 |
| f0018 | BuildSdlAst realizes SdlSourceModel. | 247 |
| f0019 | BuildSduiAst allocated-to CommandLineHost in mode SourceInspection. | 248 |
| f0020 | BuildSduiAst contributes-to StructuralModelInspection. | 249 |
| f0021 | BuildSduiAst has repeatability = deterministic. | 250 |
| f0022 | BuildSduiAst has state-retention = stateless. | 251 |
| f0023 | BuildSduiAst realizes SduiSourceModel. | 252 |
| f0024 | CancelPendingActions realizes ExecutableDesign. | 253 |
| f0025 | CheckDomainStateCompatibility realizes ExecutableDesign. | 254 |
| f0026 | CheckExecutionCompleteness realizes ExecutableDesign. | 255 |
| f0027 | CheckFunctionSignatures realizes ExecutableDesign. | 256 |
| f0028 | CloseSdlInstance realizes ExecutableDesign. | 257 |
| f0029 | CloseUiInstance realizes InteractiveSession. | 258 |
| f0030 | CoalesceSourceChanges realizes DevelopmentReload. | 259 |
| f0031 | CommandLineHost consumes ExportSinkPort. | 260 |
| f0032 | CommandLineHost consumes PreparedFramePort. | 261 |
| f0033 | CommandLineHost consumes SduiFrontendPort. | 262 |
| f0034 | CommandLineHost consumes SourceSnapshotPort. | 263 |
| f0035 | CommandLineHost owns ComposeHeadlessExport. | 264 |
| f0036 | CommandLineHost owns WriteGeneratedArtifacts. | 265 |
| f0037 | CommandLineHost provides StaticDocumentation. | 266 |
| f0038 | ComposeHeadlessExport realizes StaticDocumentation. | 267 |
| f0039 | ComposeInteractiveSession allocated-to FyneHost in mode UiPreview. | 268 |
| f0040 | ComposeInteractiveSession contributes-to InteractiveUiPreview. | 269 |
| f0041 | ComposeInteractiveSession realizes NativeInteraction. | 270 |
| f0042 | ComposeMarkdownDocument allocated-to CommandLineHost in mode StaticExport. | 271 |
| f0043 | ComposeMarkdownDocument contributes-to DesignDocumentation. | 272 |
| f0044 | ComposeMarkdownDocument realizes StaticDocumentation. | 273 |
| f0045 | ComputeClipping realizes MeasuredPresentation. | 274 |
| f0046 | ConnectTypedWidgetHandles allocated-to FyneHost in mode BoundExecution. | 275 |
| f0047 | ConnectTypedWidgetHandles contributes-to TypedDomainBinding. | 276 |
| f0048 | ConnectTypedWidgetHandles realizes BoundInteraction. | 277 |
| f0049 | ContentServices contains DiagramProvider. | 278 |
| f0050 | ContentServices contains MarkdownProvider. | 279 |
| f0051 | ContentServices contains ResourceStore. | 280 |
| f0052 | ContentServices provides RichContent. | 281 |
| f0053 | CoordinateSdlCompilation realizes SdlSourceModel. | 282 |
| f0054 | CoordinateSduiCompilation realizes SduiSourceModel. | 283 |
| f0055 | CorrelateActionResult realizes ExecutableDesign. | 284 |
| f0056 | CorrelateUiResult realizes InteractiveSession. | 285 |
| f0057 | CreateSdlInstance realizes ExecutableDesign. | 286 |
| f0058 | CreateUiInstance realizes InteractiveSession. | 287 |
| f0059 | DesignAuthor pursues BuildNativeProduct. | 288 |
| f0060 | DesignAuthor pursues EditRunningPrototype. | 289 |
| f0061 | DesignAuthor pursues InspectModels. | 290 |
| f0062 | DesignAuthor pursues PrototypeUserInterface. | 291 |
| f0063 | DesignAuthor pursues PublishDesignDocumentation. | 292 |
| f0064 | DesignAuthor pursues TryDomainInteraction. | 293 |
| f0065 | DesignDocumentation supports PublishDesignDocumentation. | 294 |
| f0066 | DesignReviewer pursues InspectModels. | 295 |
| f0067 | DesignReviewer pursues PrototypeUserInterface. | 296 |
| f0068 | DesignReviewer pursues PublishDesignDocumentation. | 297 |
| f0069 | DesignSourceArchive holds DesignSourceDocuments. | 298 |
| f0070 | DesignSourceDocuments upholds DesignSourceRecord. | 299 |
| f0071 | DesignSourceRecord has completeness = closed. | 300 |
| f0072 | DesignSourceRecord has-field SourceDocumentRevision. | 301 |
| f0073 | DesignSourceRecord has-field SourceDocumentText. | 302 |
| f0074 | DevelopmentReload requires BindingReloadPort in mode BoundLiveEditing. | 303 |
| f0075 | DevelopmentReload requires FileChangePort in mode BoundLiveEditing. | 304 |
| f0076 | DevelopmentReload requires FileChangePort in mode LiveEditing. | 305 |
| f0077 | DevelopmentReload requires SdlReloadPort in mode BoundLiveEditing. | 306 |
| f0078 | DevelopmentReload requires SourceSnapshotPort in mode BoundLiveEditing. | 307 |
| f0079 | DevelopmentReload requires SourceSnapshotPort in mode LiveEditing. | 308 |
| f0080 | DevelopmentReload requires UiReloadPort in mode BoundLiveEditing. | 309 |
| f0081 | DevelopmentReload requires UiReloadPort in mode LiveEditing. | 310 |
| f0082 | DevelopmentTools contains DiagnosticReporter. | 311 |
| f0083 | DevelopmentTools contains GoBuildRunner. | 312 |
| f0084 | DevelopmentTools contains GoCodeGenerator. | 313 |
| f0085 | DevelopmentTools contains ReloadCoordinator. | 314 |
| f0086 | DevelopmentTools contains SourceLoader. | 315 |
| f0087 | DevelopmentTools contains SourceWatcher. | 316 |
| f0088 | DevelopmentTools provides DevelopmentReload. | 317 |
| f0089 | DiagnosticReporter owns ReportBindingDiagnostics. | 318 |
| f0090 | DiagnosticReporter owns ReportSourceDiagnostics. | 319 |
| f0091 | DiagnosticReporter provides SourceDiagnostics. | 320 |
| f0092 | DiagramProvider consumes DiagramEnginePort. | 321 |
| f0093 | DiagramProvider consumes ResourcePort. | 322 |
| f0094 | DiagramProvider owns PrepareDiagramResource. | 323 |
| f0095 | DiagramProvider provides RichContent. | 324 |
| f0096 | DisconnectBindings realizes BoundInteraction. | 325 |
| f0097 | DispatchUiEvent allocated-to FyneHost in mode BoundExecution. | 326 |
| f0098 | DispatchUiEvent contributes-to TypedDomainBinding. | 327 |
| f0099 | DispatchUiEvent realizes InteractiveSession. | 328 |
| f0100 | DomainStateMigrator consumes DomainStatePort. | 329 |
| f0101 | DomainStateMigrator consumes SdlModelPort. | 330 |
| f0102 | DomainStateMigrator owns CheckDomainStateCompatibility. | 331 |
| f0103 | DomainStateMigrator owns MigrateOrResetDomainState. | 332 |
| f0104 | DomainStateMigrator provides ExecutableDesign. | 333 |
| f0105 | ExecutableDesign requires DomainFunctionPort in mode BoundExecution. | 334 |
| f0106 | ExpandUiDefinitions has repeatability = deterministic. | 335 |
| f0107 | ExpandUiDefinitions has state-retention = stateless. | 336 |
| f0108 | ExpandUiDefinitions realizes SduiSourceModel. | 337 |
| f0109 | ExportConsoleSnapshot realizes StaticDocumentation. | 338 |
| f0110 | ExportSvgSnapshot allocated-to CommandLineHost in mode StaticExport. | 339 |
| f0111 | ExportSvgSnapshot contributes-to DesignDocumentation. | 340 |
| f0112 | ExportSvgSnapshot realizes StaticDocumentation. | 341 |
| f0113 | ExportUiDocumentation refines InspectDesignSource. | 342 |
| f0114 | ExportViewpointMarkdown allocated-to CommandLineHost in mode StaticExport. | 343 |
| f0115 | ExportViewpointMarkdown contributes-to DesignDocumentation. | 344 |
| f0116 | ExportViewpointMarkdown has repeatability = deterministic. | 345 |
| f0117 | ExportViewpointMarkdown has state-retention = stateless. | 346 |
| f0118 | FyneBackend consumes PreparedFramePort. | 347 |
| f0119 | FyneBackend consumes UiSessionPort. | 348 |
| f0120 | FyneBackend owns HandleFocusAndTextInput. | 349 |
| f0121 | FyneBackend owns PublishPresentation. | 350 |
| f0122 | FyneBackend owns ReconcileWidgets. | 351 |
| f0123 | FyneBackend owns ReleaseNativeWidgets. | 352 |
| f0124 | FyneBackend provides NativeInteraction. | 353 |
| f0125 | FyneHost consumes DomainBindingPort. | 354 |
| f0126 | FyneHost consumes ReloadPort. | 355 |
| f0127 | FyneHost consumes SdlFrontendPort. | 356 |
| f0128 | FyneHost consumes SduiFrontendPort. | 357 |
| f0129 | FyneHost consumes SourceSnapshotPort. | 358 |
| f0130 | FyneHost consumes UiSessionPort. | 359 |
| f0131 | FyneHost consumes WidgetBackendPort. | 360 |
| f0132 | FyneHost owns ComposeInteractiveSession. | 361 |
| f0133 | FyneHost owns ScheduleUiPublication. | 362 |
| f0134 | FyneHost provides NativeInteraction. | 363 |
| f0135 | GenerateBindingRegistration allocated-to CommandLineHost in mode NativeBuild. | 364 |
| f0136 | GenerateBindingRegistration contributes-to NativeGoAssembly. | 365 |
| f0137 | GenerateBindingRegistration realizes NativeRealization. | 366 |
| f0138 | GenerateModelConstructors allocated-to CommandLineHost in mode NativeBuild. | 367 |
| f0139 | GenerateModelConstructors contributes-to NativeGoAssembly. | 368 |
| f0140 | GenerateModelConstructors realizes NativeRealization. | 369 |
| f0141 | GoBuildRunner consumes BuildToolPort. | 370 |
| f0142 | GoBuildRunner consumes GeneratedArtifactPort. | 371 |
| f0143 | GoBuildRunner owns BuildGeneratedApplication. | 372 |
| f0144 | GoBuildRunner owns RestartChangedGoProgram. | 373 |
| f0145 | GoBuildRunner provides NativeRealization. | 374 |
| f0146 | GoCodeGenerator consumes ExecutionProfilePort. | 375 |
| f0147 | GoCodeGenerator consumes SdlModelPort. | 376 |
| f0148 | GoCodeGenerator consumes SduiModelPort. | 377 |
| f0149 | GoCodeGenerator owns GenerateBindingRegistration. | 378 |
| f0150 | GoCodeGenerator owns GenerateModelConstructors. | 379 |
| f0151 | GoCodeGenerator owns PreserveHandwrittenSources. | 380 |
| f0152 | GoCodeGenerator provides NativeRealization. | 381 |
| f0153 | GoDomainImplementation owns PerformDomainOperation. | 382 |
| f0154 | GoDomainImplementation provides DomainOperations. | 383 |
| f0155 | HandleFocusAndTextInput allocated-to FyneHost in mode UiPreview. | 384 |
| f0156 | HandleFocusAndTextInput contributes-to InteractiveUiPreview. | 385 |
| f0157 | HandleFocusAndTextInput realizes NativeInteraction. | 386 |
| f0158 | IdentifySourceRevision realizes SourceLoading. | 387 |
| f0159 | InspectSdlSource refines InspectDesignSource. | 388 |
| f0160 | InspectSduiSource refines InspectDesignSource. | 389 |
| f0161 | InteractiveUiPreview supports PrototypeUserInterface. | 390 |
| f0162 | InvokeRegisteredFunction allocated-to FyneHost in mode BoundExecution. | 391 |
| f0163 | InvokeRegisteredFunction contributes-to TypedDomainBinding. | 392 |
| f0164 | InvokeRegisteredFunction realizes ExecutableDesign. | 393 |
| f0165 | KeepLastValidModels allocated-to FyneHost in mode LiveEditing. | 394 |
| f0166 | KeepLastValidModels contributes-to LiveModelReload. | 395 |
| f0167 | KeepLastValidModels realizes DevelopmentReload. | 396 |
| f0168 | LiveModelReload supports EditRunningPrototype. | 397 |
| f0169 | ManageDomainState has state-retention = stateful. | 398 |
| f0170 | ManageDomainState realizes ExecutableDesign. | 399 |
| f0171 | ManageWidgetIdentities has state-retention = stateful. | 400 |
| f0172 | ManageWidgetIdentities realizes InteractiveSession. | 401 |
| f0173 | MarkdownProvider consumes DiagramPort. | 402 |
| f0174 | MarkdownProvider consumes MeasurementPort. | 403 |
| f0175 | MarkdownProvider consumes ResourcePort. | 404 |
| f0176 | MarkdownProvider owns MeasureMarkdownContent. | 405 |
| f0177 | MarkdownProvider owns PrepareMarkdown. | 406 |
| f0178 | MarkdownProvider provides RichContent. | 407 |
| f0179 | MatchCompatibleWidgets realizes InteractiveSession. | 408 |
| f0180 | MeasureMarkdownContent realizes RichContent. | 409 |
| f0181 | MeasureUiContent realizes MeasuredPresentation. | 410 |
| f0182 | MeasuredPresentation requires MeasurementPort in mode StaticExport. | 411 |
| f0183 | MeasuredPresentation requires MeasurementPort in mode UiPreview. | 412 |
| f0184 | MigrateOrResetDomainState realizes ExecutableDesign. | 413 |
| f0185 | NativeGoAssembly supports BuildNativeProduct. | 414 |
| f0186 | NativeInteraction requires WidgetBackendPort in mode UiPreview. | 415 |
| f0187 | NativeRealization requires BuildToolPort in mode NativeBuild. | 416 |
| f0188 | NativeRealization requires GeneratedArtifactPort in mode NativeBuild. | 417 |
| f0189 | NormalizeSdlModel has repeatability = deterministic. | 418 |
| f0190 | NormalizeSdlModel has state-retention = stateless. | 419 |
| f0191 | NormalizeSdlModel realizes SdlSourceModel. | 420 |
| f0192 | NoticeGeneration has presence = required. | 421 |
| f0193 | NoticeGeneration has value-type = unsigned. | 422 |
| f0194 | NoticeVersion has presence = required. | 423 |
| f0195 | NoticeVersion has value-type = unsigned. | 424 |
| f0196 | ObserveSourceChanges allocated-to FyneHost in mode LiveEditing. | 425 |
| f0197 | ObserveSourceChanges contributes-to LiveModelReload. | 426 |
| f0198 | ObserveSourceChanges realizes DevelopmentReload. | 427 |
| f0199 | PerformDomainOperation realizes DomainOperations. | 428 |
| f0200 | PrepareCandidateModels allocated-to FyneHost in mode LiveEditing. | 429 |
| f0201 | PrepareCandidateModels contributes-to LiveModelReload. | 430 |
| f0202 | PrepareCandidateModels realizes DevelopmentReload. | 431 |
| f0203 | PrepareDiagramResource realizes RichContent. | 432 |
| f0204 | PrepareMarkdown realizes RichContent. | 433 |
| f0205 | PreserveCompatibleUiState allocated-to FyneHost in mode LiveEditing. | 434 |
| f0206 | PreserveCompatibleUiState contributes-to LiveModelReload. | 435 |
| f0207 | PreserveCompatibleUiState realizes InteractiveSession. | 436 |
| f0208 | PreserveHandwrittenSources allocated-to CommandLineHost in mode NativeBuild. | 437 |
| f0209 | PreserveHandwrittenSources contributes-to NativeGoAssembly. | 438 |
| f0210 | PreserveHandwrittenSources realizes NativeRealization. | 439 |
| f0211 | PreserveSdlSourceMap has repeatability = deterministic. | 440 |
| f0212 | PreserveSdlSourceMap has state-retention = stateless. | 441 |
| f0213 | PreserveSdlSourceMap realizes SdlSourceModel. | 442 |
| f0214 | PreserveUiRegions has repeatability = deterministic. | 443 |
| f0215 | PreserveUiRegions has state-retention = stateless. | 444 |
| f0216 | PreserveUiRegions realizes SduiSourceModel. | 445 |
| f0217 | PreserveUiSourceMap has repeatability = deterministic. | 446 |
| f0218 | PreserveUiSourceMap has state-retention = stateless. | 447 |
| f0219 | PreserveUiSourceMap realizes SduiSourceModel. | 448 |
| f0220 | ProjectSdlViewpoints allocated-to CommandLineHost in mode StaticExport. | 449 |
| f0221 | ProjectSdlViewpoints contributes-to DesignDocumentation. | 450 |
| f0222 | ProjectSdlViewpoints has repeatability = deterministic. | 451 |
| f0223 | ProjectSdlViewpoints has state-retention = stateless. | 452 |
| f0224 | ProjectUiGeneration allocated-to FyneHost in mode LiveEditing. | 453 |
| f0225 | ProjectUiGeneration contributes-to LiveModelReload. | 454 |
| f0226 | ProjectUiGeneration projects UiSessionState into UiGenerationNotices. | 455 |
| f0227 | PublishDomainUpdates allocated-to FyneHost in mode BoundExecution. | 456 |
| f0228 | PublishDomainUpdates contributes-to TypedDomainBinding. | 457 |
| f0229 | PublishDomainUpdates realizes BoundInteraction. | 458 |
| f0230 | PublishModelGeneration allocated-to FyneHost in mode LiveEditing. | 459 |
| f0231 | PublishModelGeneration contributes-to LiveModelReload. | 460 |
| f0232 | PublishModelGeneration has state-retention = stateful. | 461 |
| f0233 | PublishModelGeneration realizes DevelopmentReload. | 462 |
| f0234 | PublishPresentation realizes NativeInteraction. | 463 |
| f0235 | ReadBoundedSources realizes SourceLoading. | 464 |
| f0236 | ReconcileWidgets allocated-to FyneHost in mode UiPreview. | 465 |
| f0237 | ReconcileWidgets contributes-to InteractiveUiPreview. | 466 |
| f0238 | ReconcileWidgets has state-retention = stateful. | 467 |
| f0239 | ReconcileWidgets realizes NativeInteraction. | 468 |
| f0240 | RegisterDomainFunctions has state-retention = stateful. | 469 |
| f0241 | RegisterDomainFunctions realizes ExecutableDesign. | 470 |
| f0242 | RejectStaleUiEvent realizes InteractiveSession. | 471 |
| f0243 | ReleaseNativeWidgets realizes NativeInteraction. | 472 |
| f0244 | ReleaseVisualResources realizes RichContent. | 473 |
| f0245 | ReloadBoundModels refines ReloadDesignSession. | 474 |
| f0246 | ReloadCoordinator consumes BindingReloadPort. | 475 |
| f0247 | ReloadCoordinator consumes DiagnosticPort. | 476 |
| f0248 | ReloadCoordinator consumes SdlFrontendPort. | 477 |
| f0249 | ReloadCoordinator consumes SdlReloadPort. | 478 |
| f0250 | ReloadCoordinator consumes SduiFrontendPort. | 479 |
| f0251 | ReloadCoordinator consumes SourceSnapshotPort. | 480 |
| f0252 | ReloadCoordinator consumes UiReloadPort. | 481 |
| f0253 | ReloadCoordinator owns KeepLastValidModels. | 482 |
| f0254 | ReloadCoordinator owns PrepareCandidateModels. | 483 |
| f0255 | ReloadCoordinator owns PublishModelGeneration. | 484 |
| f0256 | ReloadCoordinator owns RetirePreviousGeneration. | 485 |
| f0257 | ReloadCoordinator provides DevelopmentReload. | 486 |
| f0258 | ReloadUiModel refines ReloadDesignSession. | 487 |
| f0259 | ReportBindingDiagnostics realizes SourceDiagnostics. | 488 |
| f0260 | ReportSourceDiagnostics allocated-to CommandLineHost in mode SourceInspection. | 489 |
| f0261 | ReportSourceDiagnostics allocated-to FyneHost in mode LiveEditing. | 490 |
| f0262 | ReportSourceDiagnostics contributes-to LiveModelReload. | 491 |
| f0263 | ReportSourceDiagnostics contributes-to StructuralModelInspection. | 492 |
| f0264 | ReportSourceDiagnostics realizes SourceDiagnostics. | 493 |
| f0265 | ResetIncompatibleUiState realizes InteractiveSession. | 494 |
| f0266 | ResolveAncestorDimensions realizes MeasuredPresentation. | 495 |
| f0267 | ResolveCallbackSymbols allocated-to FyneHost in mode BoundExecution. | 496 |
| f0268 | ResolveCallbackSymbols contributes-to TypedDomainBinding. | 497 |
| f0269 | ResolveCallbackSymbols realizes BoundInteraction. | 498 |
| f0270 | ResolveSdlSymbols has repeatability = deterministic. | 499 |
| f0271 | ResolveSdlSymbols has state-retention = stateless. | 500 |
| f0272 | ResolveSdlSymbols realizes SdlSourceModel. | 501 |
| f0273 | ResolveUiNames has repeatability = deterministic. | 502 |
| f0274 | ResolveUiNames has state-retention = stateless. | 503 |
| f0275 | ResolveUiNames realizes SduiSourceModel. | 504 |
| f0276 | ResourceStore owns ReleaseVisualResources. | 505 |
| f0277 | ResourceStore owns ValidateVisualResources. | 506 |
| f0278 | ResourceStore provides RichContent. | 507 |
| f0279 | RestartChangedGoProgram realizes NativeRealization. | 508 |
| f0280 | RetirePreviousGeneration realizes DevelopmentReload. | 509 |
| f0281 | RevokeWidgetGenerations realizes InteractiveSession. | 510 |
| f0282 | RichContent requires ContentProviderPort in mode RichDocument. | 511 |
| f0283 | RichContent requires DiagramEnginePort in mode RichDocument. | 512 |
| f0284 | RouteDomainBindings realizes BoundInteraction. | 513 |
| f0285 | RunBoundUiAction refines RunDesignSession. | 514 |
| f0286 | RunUnboundUiPreview refines RunDesignSession. | 515 |
| f0287 | ScheduleUiPublication realizes NativeInteraction. | 516 |
| f0288 | SdlDispatcher consumes DomainFunctionPort. | 517 |
| f0289 | SdlDispatcher consumes DomainStatePort. | 518 |
| f0290 | SdlDispatcher owns CancelPendingActions. | 519 |
| f0291 | SdlDispatcher owns CorrelateActionResult. | 520 |
| f0292 | SdlDispatcher owns InvokeRegisteredFunction. | 521 |
| f0293 | SdlDispatcher owns ValidateActionInput. | 522 |
| f0294 | SdlDispatcher provides ExecutableDesign. | 523 |
| f0295 | SdlExecutionGate consumes DiagnosticPort. | 524 |
| f0296 | SdlExecutionGate consumes SdlModelPort. | 525 |
| f0297 | SdlExecutionGate owns CheckExecutionCompleteness. | 526 |
| f0298 | SdlExecutionGate provides ExecutableDesign. | 527 |
| f0299 | SdlFrontend consumes DiagnosticPort. | 528 |
| f0300 | SdlFrontend consumes SourceSnapshotPort. | 529 |
| f0301 | SdlFrontend contains SdlLexer. | 530 |
| f0302 | SdlFrontend contains SdlNormalizer. | 531 |
| f0303 | SdlFrontend contains SdlParser. | 532 |
| f0304 | SdlFrontend contains SdlValidator. | 533 |
| f0305 | SdlFrontend owns CoordinateSdlCompilation. | 534 |
| f0306 | SdlFrontend provides SdlSourceModel. | 535 |
| f0307 | SdlFunctionRegistry consumes DomainFunctionPort. | 536 |
| f0308 | SdlFunctionRegistry owns CheckFunctionSignatures. | 537 |
| f0309 | SdlFunctionRegistry owns RegisterDomainFunctions. | 538 |
| f0310 | SdlFunctionRegistry provides ExecutableDesign. | 539 |
| f0311 | SdlLexer owns TokenizeSdlSource. | 540 |
| f0312 | SdlLexer provides SdlSourceModel. | 541 |
| f0313 | SdlLibrary contains SdlFrontend. | 542 |
| f0314 | SdlLibrary contains SdlRuntime. | 543 |
| f0315 | SdlLibrary contains SdlViewpointGenerator. | 544 |
| f0316 | SdlLibrary provides ExecutableDesign. | 545 |
| f0317 | SdlLibrary provides SdlSourceModel. | 546 |
| f0318 | SdlNormalizer owns NormalizeSdlModel. | 547 |
| f0319 | SdlNormalizer owns PreserveSdlSourceMap. | 548 |
| f0320 | SdlNormalizer provides SdlSourceModel. | 549 |
| f0321 | SdlParser owns BuildSdlAst. | 550 |
| f0322 | SdlParser provides SdlSourceModel. | 551 |
| f0323 | SdlRuntime consumes DomainFunctionPort. | 552 |
| f0324 | SdlRuntime consumes SdlModelPort. | 553 |
| f0325 | SdlRuntime contains DomainStateMigrator. | 554 |
| f0326 | SdlRuntime contains SdlDispatcher. | 555 |
| f0327 | SdlRuntime contains SdlExecutionGate. | 556 |
| f0328 | SdlRuntime contains SdlFunctionRegistry. | 557 |
| f0329 | SdlRuntime contains SdlStateStore. | 558 |
| f0330 | SdlRuntime owns CloseSdlInstance. | 559 |
| f0331 | SdlRuntime owns CreateSdlInstance. | 560 |
| f0332 | SdlRuntime provides ExecutableDesign. | 561 |
| f0333 | SdlSourceModel requires SourceSnapshotPort in mode SourceInspection. | 562 |
| f0334 | SdlStateStore owns ManageDomainState. | 563 |
| f0335 | SdlStateStore owns SnapshotDomainState. | 564 |
| f0336 | SdlStateStore provides ExecutableDesign. | 565 |
| f0337 | SdlUiBindingAdapter consumes DiagnosticPort. | 566 |
| f0338 | SdlUiBindingAdapter consumes SdlExecutionPort. | 567 |
| f0339 | SdlUiBindingAdapter consumes UiSessionPort. | 568 |
| f0340 | SdlUiBindingAdapter owns ConnectTypedWidgetHandles. | 569 |
| f0341 | SdlUiBindingAdapter owns DisconnectBindings. | 570 |
| f0342 | SdlUiBindingAdapter owns PublishDomainUpdates. | 571 |
| f0343 | SdlUiBindingAdapter owns ResolveCallbackSymbols. | 572 |
| f0344 | SdlUiBindingAdapter owns RouteDomainBindings. | 573 |
| f0345 | SdlUiBindingAdapter provides BoundInteraction. | 574 |
| f0346 | SdlValidator owns ResolveSdlSymbols. | 575 |
| f0347 | SdlValidator owns ValidateSdlProfile. | 576 |
| f0348 | SdlValidator owns ValidateSdlStructure. | 577 |
| f0349 | SdlValidator provides SdlSourceModel. | 578 |
| f0350 | SdlViewpointGenerator owns ExportViewpointMarkdown. | 579 |
| f0351 | SdlViewpointGenerator owns ProjectSdlViewpoints. | 580 |
| f0352 | SdlViewpointGenerator owns TraceViewpointFacts. | 581 |
| f0353 | SduiDispatcher consumes DomainBindingPort. | 582 |
| f0354 | SduiDispatcher consumes UiStatePort. | 583 |
| f0355 | SduiDispatcher owns CorrelateUiResult. | 584 |
| f0356 | SduiDispatcher owns DispatchUiEvent. | 585 |
| f0357 | SduiDispatcher owns RejectStaleUiEvent. | 586 |
| f0358 | SduiDispatcher owns ValidateUiEvent. | 587 |
| f0359 | SduiDispatcher provides InteractiveSession. | 588 |
| f0360 | SduiFrontend consumes DiagnosticPort. | 589 |
| f0361 | SduiFrontend consumes SourceSnapshotPort. | 590 |
| f0362 | SduiFrontend contains SduiLexer. | 591 |
| f0363 | SduiFrontend contains SduiNormalizer. | 592 |
| f0364 | SduiFrontend contains SduiParser. | 593 |
| f0365 | SduiFrontend contains SduiValidator. | 594 |
| f0366 | SduiFrontend owns CoordinateSduiCompilation. | 595 |
| f0367 | SduiFrontend provides SduiSourceModel. | 596 |
| f0368 | SduiInstanceStore holds UiSessionState. | 597 |
| f0369 | SduiInstanceStore owns ManageWidgetIdentities. | 598 |
| f0370 | SduiInstanceStore owns ProjectUiGeneration. | 599 |
| f0371 | SduiInstanceStore owns RevokeWidgetGenerations. | 600 |
| f0372 | SduiInstanceStore owns SnapshotUiState. | 601 |
| f0373 | SduiInstanceStore provides InteractiveSession. | 602 |
| f0374 | SduiLayout consumes ContentProviderPort. | 603 |
| f0375 | SduiLayout consumes MeasurementPort. | 604 |
| f0376 | SduiLayout consumes UiSnapshotPort. | 605 |
| f0377 | SduiLayout owns AllocateGeometry. | 606 |
| f0378 | SduiLayout owns BuildPreparedFrame. | 607 |
| f0379 | SduiLayout owns ComputeClipping. | 608 |
| f0380 | SduiLayout owns MeasureUiContent. | 609 |
| f0381 | SduiLayout owns ResolveAncestorDimensions. | 610 |
| f0382 | SduiLayout provides MeasuredPresentation. | 611 |
| f0383 | SduiLexer owns TokenizeSduiSource. | 612 |
| f0384 | SduiLexer provides SduiSourceModel. | 613 |
| f0385 | SduiLibrary contains SduiFrontend. | 614 |
| f0386 | SduiLibrary contains SduiLayout. | 615 |
| f0387 | SduiLibrary contains SduiPresentation. | 616 |
| f0388 | SduiLibrary contains SduiRuntime. | 617 |
| f0389 | SduiLibrary provides InteractiveSession. | 618 |
| f0390 | SduiLibrary provides MeasuredPresentation. | 619 |
| f0391 | SduiLibrary provides SduiSourceModel. | 620 |
| f0392 | SduiLibrary provides StaticDocumentation. | 621 |
| f0393 | SduiNormalizer owns ExpandUiDefinitions. | 622 |
| f0394 | SduiNormalizer owns PreserveUiRegions. | 623 |
| f0395 | SduiNormalizer owns PreserveUiSourceMap. | 624 |
| f0396 | SduiNormalizer provides SduiSourceModel. | 625 |
| f0397 | SduiParser owns BuildSduiAst. | 626 |
| f0398 | SduiParser provides SduiSourceModel. | 627 |
| f0399 | SduiPresentation consumes PreparedFramePort. | 628 |
| f0400 | SduiPresentation owns ComposeMarkdownDocument. | 629 |
| f0401 | SduiPresentation owns ExportConsoleSnapshot. | 630 |
| f0402 | SduiPresentation owns ExportSvgSnapshot. | 631 |
| f0403 | SduiPresentation provides StaticDocumentation. | 632 |
| f0404 | SduiPropertyStore consumes UiStatePort. | 633 |
| f0405 | SduiPropertyStore owns ApplyPropertyBatch. | 634 |
| f0406 | SduiPropertyStore owns TrackInputDraft. | 635 |
| f0407 | SduiPropertyStore owns ValidatePropertyBatch. | 636 |
| f0408 | SduiPropertyStore provides InteractiveSession. | 637 |
| f0409 | SduiRuntime consumes DomainBindingPort. | 638 |
| f0410 | SduiRuntime consumes SduiModelPort. | 639 |
| f0411 | SduiRuntime contains SduiDispatcher. | 640 |
| f0412 | SduiRuntime contains SduiInstanceStore. | 641 |
| f0413 | SduiRuntime contains SduiPropertyStore. | 642 |
| f0414 | SduiRuntime contains UiStateReconciler. | 643 |
| f0415 | SduiRuntime owns CloseUiInstance. | 644 |
| f0416 | SduiRuntime owns CreateUiInstance. | 645 |
| f0417 | SduiRuntime provides InteractiveSession. | 646 |
| f0418 | SduiSourceModel requires SourceSnapshotPort in mode SourceInspection. | 647 |
| f0419 | SduiValidator owns ResolveUiNames. | 648 |
| f0420 | SduiValidator owns ValidateRelativeFormatting. | 649 |
| f0421 | SduiValidator owns ValidateSymbolicBindings. | 650 |
| f0422 | SduiValidator owns ValidateWidgetArguments. | 651 |
| f0423 | SduiValidator provides SduiSourceModel. | 652 |
| f0424 | SessionDraft has presence = optional. | 653 |
| f0425 | SessionDraft has value-type = text. | 654 |
| f0426 | SessionGeneration has presence = required. | 655 |
| f0427 | SessionGeneration has value-type = unsigned. | 656 |
| f0428 | SnapshotDomainState realizes ExecutableDesign. | 657 |
| f0429 | SnapshotUiState realizes InteractiveSession. | 658 |
| f0430 | SourceDocumentRevision has presence = required. | 659 |
| f0431 | SourceDocumentRevision has value-type = unsigned. | 660 |
| f0432 | SourceDocumentText has presence = required. | 661 |
| f0433 | SourceDocumentText has value-type = text. | 662 |
| f0434 | SourceLoader consumes SourceInputPort. | 663 |
| f0435 | SourceLoader owns DesignSourceArchive. | 664 |
| f0436 | SourceLoader owns IdentifySourceRevision. | 665 |
| f0437 | SourceLoader owns ReadBoundedSources. | 666 |
| f0438 | SourceLoader provides SourceLoading. | 667 |
| f0439 | SourceWatcher consumes FileChangePort. | 668 |
| f0440 | SourceWatcher owns CoalesceSourceChanges. | 669 |
| f0441 | SourceWatcher owns ObserveSourceChanges. | 670 |
| f0442 | SourceWatcher provides DevelopmentReload. | 671 |
| f0443 | StaticDocumentation requires ExportSinkPort in mode StaticExport. | 672 |
| f0444 | StaticDocumentation requires PreparedFramePort in mode StaticExport. | 673 |
| f0445 | StructuralModelInspection supports EditRunningPrototype. | 674 |
| f0446 | StructuralModelInspection supports InspectModels. | 675 |
| f0447 | TokenizeSdlSource has repeatability = deterministic. | 676 |
| f0448 | TokenizeSdlSource has state-retention = stateless. | 677 |
| f0449 | TokenizeSdlSource realizes SdlSourceModel. | 678 |
| f0450 | TokenizeSduiSource has repeatability = deterministic. | 679 |
| f0451 | TokenizeSduiSource has state-retention = stateless. | 680 |
| f0452 | TokenizeSduiSource realizes SduiSourceModel. | 681 |
| f0453 | TraceViewpointFacts allocated-to CommandLineHost in mode SourceInspection. | 682 |
| f0454 | TraceViewpointFacts allocated-to CommandLineHost in mode StaticExport. | 683 |
| f0455 | TraceViewpointFacts contributes-to DesignDocumentation. | 684 |
| f0456 | TraceViewpointFacts contributes-to InspectModels. | 685 |
| f0457 | TraceViewpointFacts has repeatability = deterministic. | 686 |
| f0458 | TraceViewpointFacts has state-retention = stateless. | 687 |
| f0459 | TrackInputDraft has state-retention = stateful. | 688 |
| f0460 | TrackInputDraft realizes InteractiveSession. | 689 |
| f0461 | TypedDomainBinding supports TryDomainInteraction. | 690 |
| f0462 | UiGenerationChanged has-field NoticeGeneration. | 691 |
| f0463 | UiGenerationContract defines UiGenerationChanged. | 692 |
| f0464 | UiGenerationContract has completeness = closed. | 693 |
| f0465 | UiGenerationContract has-field NoticeVersion. | 694 |
| f0466 | UiGenerationNotices from UiSessionState. | 695 |
| f0467 | UiGenerationNotices upholds UiGenerationContract. | 696 |
| f0468 | UiGenerationWire encodes UiGenerationChanged. | 697 |
| f0469 | UiGenerationWire has bit-order = most-significant-first. | 698 |
| f0470 | UiGenerationWire has byte-order = big-endian. | 699 |
| f0471 | UiGenerationWire places NoticeGeneration at 16 bits 64. | 700 |
| f0472 | UiGenerationWire places NoticeVersion at 0 bits 16. | 701 |
| f0473 | UiSessionRecord has completeness = closed. | 702 |
| f0474 | UiSessionRecord has-field SessionDraft. | 703 |
| f0475 | UiSessionRecord has-field SessionGeneration. | 704 |
| f0476 | UiSessionState upholds UiSessionRecord. | 705 |
| f0477 | UiStateReconciler consumes SduiModelPort. | 706 |
| f0478 | UiStateReconciler consumes UiStatePort. | 707 |
| f0479 | UiStateReconciler owns MatchCompatibleWidgets. | 708 |
| f0480 | UiStateReconciler owns PreserveCompatibleUiState. | 709 |
| f0481 | UiStateReconciler owns ResetIncompatibleUiState. | 710 |
| f0482 | UiStateReconciler provides InteractiveSession. | 711 |
| f0483 | ValidateActionInput realizes ExecutableDesign. | 712 |
| f0484 | ValidatePropertyBatch realizes InteractiveSession. | 713 |
| f0485 | ValidateRelativeFormatting has repeatability = deterministic. | 714 |
| f0486 | ValidateRelativeFormatting has state-retention = stateless. | 715 |
| f0487 | ValidateRelativeFormatting realizes SduiSourceModel. | 716 |
| f0488 | ValidateSdlProfile has repeatability = deterministic. | 717 |
| f0489 | ValidateSdlProfile has state-retention = stateless. | 718 |
| f0490 | ValidateSdlProfile realizes SdlSourceModel. | 719 |
| f0491 | ValidateSdlStructure allocated-to CommandLineHost in mode SourceInspection. | 720 |
| f0492 | ValidateSdlStructure contributes-to StructuralModelInspection. | 721 |
| f0493 | ValidateSdlStructure has repeatability = deterministic. | 722 |
| f0494 | ValidateSdlStructure has state-retention = stateless. | 723 |
| f0495 | ValidateSdlStructure realizes SdlSourceModel. | 724 |
| f0496 | ValidateSymbolicBindings has repeatability = deterministic. | 725 |
| f0497 | ValidateSymbolicBindings has state-retention = stateless. | 726 |
| f0498 | ValidateSymbolicBindings realizes SduiSourceModel. | 727 |
| f0499 | ValidateUiEvent realizes InteractiveSession. | 728 |
| f0500 | ValidateVisualResources realizes RichContent. | 729 |
| f0501 | ValidateWidgetArguments allocated-to CommandLineHost in mode SourceInspection. | 730 |
| f0502 | ValidateWidgetArguments contributes-to StructuralModelInspection. | 731 |
| f0503 | ValidateWidgetArguments has repeatability = deterministic. | 732 |
| f0504 | ValidateWidgetArguments has state-retention = stateless. | 733 |
| f0505 | ValidateWidgetArguments realizes SduiSourceModel. | 734 |
| f0506 | WriteGeneratedArtifacts allocated-to CommandLineHost in mode StaticExport. | 735 |
| f0507 | WriteGeneratedArtifacts contributes-to DesignDocumentation. | 736 |
| f0508 | WriteGeneratedArtifacts realizes StaticDocumentation. | 737 |

### Deklarasjonsregister

| Identitet | Type | Kildelinje |
| --- | --- | --- |
| AllocateGeometry | functionality | 2 |
| ApplyPropertyBatch | functionality | 3 |
| BindingReloadPort | interface | 4 |
| BoundExecution | mode | 5 |
| BoundInteraction | capability | 6 |
| BoundLiveEditing | mode | 7 |
| BuildGeneratedApplication | functionality | 8 |
| BuildNativeProduct | usecase | 9 |
| BuildNativeRealization | activity | 10 |
| BuildPreparedFrame | functionality | 11 |
| BuildSdlAst | functionality | 12 |
| BuildSduiAst | functionality | 13 |
| BuildToolPort | interface | 14 |
| CancelPendingActions | functionality | 15 |
| CheckDomainStateCompatibility | functionality | 16 |
| CheckExecutionCompleteness | functionality | 17 |
| CheckFunctionSignatures | functionality | 18 |
| CloseSdlInstance | functionality | 19 |
| CloseUiInstance | functionality | 20 |
| CoalesceSourceChanges | functionality | 21 |
| CommandLineHost | container | 22 |
| ComposeHeadlessExport | functionality | 23 |
| ComposeInteractiveSession | functionality | 24 |
| ComposeMarkdownDocument | functionality | 25 |
| ComputeClipping | functionality | 26 |
| ConnectTypedWidgetHandles | functionality | 27 |
| ContentProviderPort | interface | 28 |
| ContentServices | unit | 29 |
| CoordinateSdlCompilation | functionality | 30 |
| CoordinateSduiCompilation | functionality | 31 |
| CorrelateActionResult | functionality | 32 |
| CorrelateUiResult | functionality | 33 |
| CreateSdlInstance | functionality | 34 |
| CreateUiInstance | functionality | 35 |
| DesignAuthor | actor | 36 |
| DesignDocumentation | feature | 37 |
| DesignReviewer | actor | 38 |
| DesignSourceArchive | database | 39 |
| DesignSourceDocuments | dataset | 40 |
| DesignSourceRecord | contract | 41 |
| DevelopmentReload | capability | 42 |
| DevelopmentTools | unit | 43 |
| DiagnosticPort | interface | 44 |
| DiagnosticReporter | unit | 45 |
| DiagramEnginePort | interface | 46 |
| DiagramPort | interface | 47 |
| DiagramProvider | unit | 48 |
| DisconnectBindings | functionality | 49 |
| DispatchUiEvent | functionality | 50 |
| DomainBindingPort | interface | 51 |
| DomainFunctionPort | interface | 52 |
| DomainOperations | capability | 53 |
| DomainStateMigrator | unit | 54 |
| DomainStatePort | interface | 55 |
| EditRunningPrototype | usecase | 56 |
| ExecutableDesign | capability | 57 |
| ExecutionProfilePort | interface | 58 |
| ExpandUiDefinitions | functionality | 59 |
| ExportConsoleSnapshot | functionality | 60 |
| ExportSinkPort | interface | 61 |
| ExportSvgSnapshot | functionality | 62 |
| ExportUiDocumentation | activity | 63 |
| ExportViewpointMarkdown | functionality | 64 |
| FileChangePort | interface | 65 |
| FyneBackend | unit | 66 |
| FyneHost | container | 67 |
| GenerateBindingRegistration | functionality | 68 |
| GenerateModelConstructors | functionality | 69 |
| GeneratedArtifactPort | interface | 70 |
| GoBuildRunner | unit | 71 |
| GoCodeGenerator | unit | 72 |
| GoDomainImplementation | unit | 73 |
| HandleFocusAndTextInput | functionality | 74 |
| IdentifySourceRevision | functionality | 75 |
| InspectDesignSource | activity | 76 |
| InspectModels | usecase | 77 |
| InspectSdlSource | activity | 78 |
| InspectSduiSource | activity | 79 |
| InteractiveSession | capability | 80 |
| InteractiveUiPreview | feature | 81 |
| InvokeRegisteredFunction | functionality | 82 |
| KeepLastValidModels | functionality | 83 |
| LiveEditing | mode | 84 |
| LiveModelReload | feature | 85 |
| ManageDomainState | functionality | 86 |
| ManageWidgetIdentities | functionality | 87 |
| MarkdownProvider | unit | 88 |
| MatchCompatibleWidgets | functionality | 89 |
| MeasureMarkdownContent | functionality | 90 |
| MeasureUiContent | functionality | 91 |
| MeasuredPresentation | capability | 92 |
| MeasurementPort | interface | 93 |
| MigrateOrResetDomainState | functionality | 94 |
| NativeBuild | mode | 95 |
| NativeGoAssembly | feature | 96 |
| NativeInteraction | capability | 97 |
| NativeRealization | capability | 98 |
| NormalizeSdlModel | functionality | 99 |
| NoticeGeneration | field | 100 |
| NoticeVersion | field | 101 |
| ObserveSourceChanges | functionality | 102 |
| PerformDomainOperation | functionality | 103 |
| PrepareCandidateModels | functionality | 104 |
| PrepareDiagramResource | functionality | 105 |
| PrepareMarkdown | functionality | 106 |
| PreparedFramePort | interface | 107 |
| PreserveCompatibleUiState | functionality | 108 |
| PreserveHandwrittenSources | functionality | 109 |
| PreserveSdlSourceMap | functionality | 110 |
| PreserveUiRegions | functionality | 111 |
| PreserveUiSourceMap | functionality | 112 |
| ProjectSdlViewpoints | functionality | 113 |
| ProjectUiGeneration | functionality | 114 |
| PrototypeUserInterface | usecase | 115 |
| PublishDesignDocumentation | usecase | 116 |
| PublishDomainUpdates | functionality | 117 |
| PublishModelGeneration | functionality | 118 |
| PublishPresentation | functionality | 119 |
| ReadBoundedSources | functionality | 120 |
| RealizeDesign | activity | 121 |
| ReconcileWidgets | functionality | 122 |
| RegisterDomainFunctions | functionality | 123 |
| RejectStaleUiEvent | functionality | 124 |
| ReleaseNativeWidgets | functionality | 125 |
| ReleaseVisualResources | functionality | 126 |
| ReloadBoundModels | activity | 127 |
| ReloadCoordinator | unit | 128 |
| ReloadDesignSession | activity | 129 |
| ReloadPort | interface | 130 |
| ReloadUiModel | activity | 131 |
| ReportBindingDiagnostics | functionality | 132 |
| ReportSourceDiagnostics | functionality | 133 |
| ResetIncompatibleUiState | functionality | 134 |
| ResolveAncestorDimensions | functionality | 135 |
| ResolveCallbackSymbols | functionality | 136 |
| ResolveSdlSymbols | functionality | 137 |
| ResolveUiNames | functionality | 138 |
| ResourcePort | interface | 139 |
| ResourceStore | unit | 140 |
| RestartChangedGoProgram | functionality | 141 |
| RetirePreviousGeneration | functionality | 142 |
| RevokeWidgetGenerations | functionality | 143 |
| RichContent | capability | 144 |
| RichDocument | mode | 145 |
| RouteDomainBindings | functionality | 146 |
| RunBoundUiAction | activity | 147 |
| RunDesignSession | activity | 148 |
| RunUnboundUiPreview | activity | 149 |
| ScheduleUiPublication | functionality | 150 |
| SdlDispatcher | unit | 151 |
| SdlExecutionGate | unit | 152 |
| SdlExecutionPort | interface | 153 |
| SdlFrontend | unit | 154 |
| SdlFrontendPort | interface | 155 |
| SdlFunctionRegistry | unit | 156 |
| SdlLexer | unit | 157 |
| SdlLibrary | unit | 158 |
| SdlModelPort | interface | 159 |
| SdlNormalizer | unit | 160 |
| SdlParser | unit | 161 |
| SdlReloadPort | interface | 162 |
| SdlRuntime | unit | 163 |
| SdlSourceModel | capability | 164 |
| SdlStateStore | unit | 165 |
| SdlUiBindingAdapter | unit | 166 |
| SdlValidator | unit | 167 |
| SdlViewpointGenerator | unit | 168 |
| SduiDispatcher | unit | 169 |
| SduiFrontend | unit | 170 |
| SduiFrontendPort | interface | 171 |
| SduiInstanceStore | unit | 172 |
| SduiLayout | unit | 173 |
| SduiLexer | unit | 174 |
| SduiLibrary | unit | 175 |
| SduiModelPort | interface | 176 |
| SduiNormalizer | unit | 177 |
| SduiParser | unit | 178 |
| SduiPresentation | unit | 179 |
| SduiPropertyStore | unit | 180 |
| SduiRuntime | unit | 181 |
| SduiSourceModel | capability | 182 |
| SduiValidator | unit | 183 |
| SessionDraft | field | 184 |
| SessionGeneration | field | 185 |
| SnapshotDomainState | functionality | 186 |
| SnapshotUiState | functionality | 187 |
| SourceDiagnostics | capability | 188 |
| SourceDocumentRevision | field | 189 |
| SourceDocumentText | field | 190 |
| SourceInputPort | interface | 191 |
| SourceInspection | mode | 192 |
| SourceLoader | unit | 193 |
| SourceLoading | capability | 194 |
| SourceSnapshotPort | interface | 195 |
| SourceWatcher | unit | 196 |
| StaticDocumentation | capability | 197 |
| StaticExport | mode | 198 |
| StructuralModelInspection | feature | 199 |
| TokenizeSdlSource | functionality | 200 |
| TokenizeSduiSource | functionality | 201 |
| TraceViewpointFacts | functionality | 202 |
| TrackInputDraft | functionality | 203 |
| TryDomainInteraction | usecase | 204 |
| TypedDomainBinding | feature | 205 |
| UiGenerationChanged | variant | 206 |
| UiGenerationContract | contract | 207 |
| UiGenerationNotices | datagram | 208 |
| UiGenerationWire | encoding | 209 |
| UiPreview | mode | 210 |
| UiReloadPort | interface | 211 |
| UiSessionPort | interface | 212 |
| UiSessionRecord | contract | 213 |
| UiSessionState | dataset | 214 |
| UiSnapshotPort | interface | 215 |
| UiStatePort | interface | 216 |
| UiStateReconciler | unit | 217 |
| ValidateActionInput | functionality | 218 |
| ValidatePropertyBatch | functionality | 219 |
| ValidateRelativeFormatting | functionality | 220 |
| ValidateSdlProfile | functionality | 221 |
| ValidateSdlStructure | functionality | 222 |
| ValidateSymbolicBindings | functionality | 223 |
| ValidateUiEvent | functionality | 224 |
| ValidateVisualResources | functionality | 225 |
| ValidateWidgetArguments | functionality | 226 |
| WidgetBackendPort | interface | 227 |
| WriteGeneratedArtifacts | functionality | 228 |
