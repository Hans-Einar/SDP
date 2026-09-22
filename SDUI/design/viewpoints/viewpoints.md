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
| VP06 — Aktiviteter og leveranseplan | Tilgjengelig | refines, addresses, delivers og depends-on; planstatus er en eksplisitt kildepåstand. |
| VP07 — Features over arkitekturen | Tilgjengelig | contributes-to, owns og eksplisitt allocated-to per modus. Uspesifisert allokering vises som hull. |
| VP08 — Channel-kontrakter og sekvenser | Tilgjengelig | Eksplisitte scenario-steg validert mot permits, deltakelse, modus og request/resultat-korrelasjon. |
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

Kildegrunnlag: f0138, f0528.

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

Kildegrunnlag: f0139, f0495, f0940.

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

Kildegrunnlag: f0140, f0145, f0941, f0968.

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

Kildegrunnlag: f0141, f0146, f0469.

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

Kildegrunnlag: f0142, f0144, f0147.

### Bruksmål: TryDomainInteraction

```mermaid
flowchart LR
    n_DesignAuthor["DesignAuthor (actor)"]
    n_TryDomainInteraction["TryDomainInteraction (usecase)"]
    n_TypedDomainBinding["TypedDomainBinding (feature)"]
    n_DesignAuthor -->|pursues| n_TryDomainInteraction
    n_TypedDomainBinding -->|supports| n_TryDomainInteraction
```

Kildegrunnlag: f0143, f0973.

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

Kildegrunnlag: f0122, f0202, f0206, f0599, f0967, f1104.

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

Kildegrunnlag: f0018, f0062, f0119, f0455, f0615.

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

Kildegrunnlag: f0474, f0562, f0566, f0584, f0603, f0609, f0673.

### Functionality-bidrag til Feature: NativeGoAssembly

```mermaid
flowchart LR
    n_BuildGeneratedApplication["BuildGeneratedApplication (functionality)"]
    n_GenerateBindingRegistration["GenerateBindingRegistration (functionality)"]
    n_GenerateModelConstructors["GenerateModelConstructors (functionality)"]
    n_NativeGoAssembly["NativeGoAssembly (feature)"]
    n_PreserveHandwrittenSources["PreserveHandwrittenSources (functionality)"]
    n_RetireReplacedPythonEntryPoints["RetireReplacedPythonEntryPoints (functionality)"]
    n_VerifyNativeBehaviorParity["VerifyNativeBehaviorParity (functionality)"]
    n_BuildGeneratedApplication -->|contributes-to| n_NativeGoAssembly
    n_GenerateBindingRegistration -->|contributes-to| n_NativeGoAssembly
    n_GenerateModelConstructors -->|contributes-to| n_NativeGoAssembly
    n_PreserveHandwrittenSources -->|contributes-to| n_NativeGoAssembly
    n_RetireReplacedPythonEntryPoints -->|contributes-to| n_NativeGoAssembly
    n_VerifyNativeBehaviorParity -->|contributes-to| n_NativeGoAssembly
```

Kildegrunnlag: f0053, f0406, f0414, f0587, f0693, f1102.

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

Kildegrunnlag: f0065, f0070, f0674, f1079, f1094.

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

Kildegrunnlag: f0126, f0179, f0471, f0606, f0679.


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

Kildegrunnlag: f0128, f0129, f0130.

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

Kildegrunnlag: f0161, f0162, f0163, f0164, f0165, f0166.

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

Kildegrunnlag: f0731, f0732, f0733, f0734.

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

Kildegrunnlag: f0745, f0746, f0747.

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

Kildegrunnlag: f0766, f0767, f0768, f0769, f0770.

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

Kildegrunnlag: f0812, f0813, f0814, f0815.

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

Kildegrunnlag: f0857, f0858, f0859, f0860.

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

Kildegrunnlag: f0889, f0890, f0891, f0892.


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

Kildegrunnlag: f0127, f0177, f0607, f0680, f0697, f0781, f0782, f0783, f0784, f0785.

### Tilbydere av kapabilitet: BoundInteraction

```mermaid
flowchart LR
    n_BoundInteraction["BoundInteraction (capability)"]
    n_SdlUiBindingAdapter["SdlUiBindingAdapter (unit)"]
    n_SdlUiBindingAdapter -->|provides| n_BoundInteraction
```

Kildegrunnlag: f0786.

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

Kildegrunnlag: f0085, f0475, f0563, f0567, f0611, f0634, f0635, f0636, f0637, f0691, f0920, f0921.

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

Kildegrunnlag: f0169, f0638, f0922.

### Bidrag til kapabilitet: DomainOperations

```mermaid
flowchart LR
    n_DomainOperations["DomainOperations (capability)"]
    n_GoDomainImplementation["GoDomainImplementation (unit)"]
    n_PerformDomainOperation["PerformDomainOperation (functionality)"]
    n_GoDomainImplementation -->|owns| n_PerformDomainOperation
    n_PerformDomainOperation -->|realizes| n_DomainOperations
```

Kildegrunnlag: f0444, f0564.

### Tilbydere av kapabilitet: DomainOperations

```mermaid
flowchart LR
    n_DomainOperations["DomainOperations (capability)"]
    n_GoDomainImplementation["GoDomainImplementation (unit)"]
    n_GoDomainImplementation -->|provides| n_DomainOperations
```

Kildegrunnlag: f0445.

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

Kildegrunnlag: f0079, f0080, f0081, f0082, f0083, f0134, f0136, f0188, f0189, f0472, f0497, f0511, f0619, f0716, f0717, f0718, f0719, f0727, f0740, f0741, f0771, f0772, f0775, f0776, f0908, f1070.

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

Kildegrunnlag: f0190, f0720, f0728, f0742, f0748, f0773, f0777.

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

Kildegrunnlag: f0021, f0084, f0135, f0137, f0180, f0499, f0506, f0585, f0620, f0676, f0694, f0800, f0801, f0802, f0803, f0833, f0835, f0836, f0883, f0884, f0885, f0893, f0894, f0909, f0972, f1052, f1053, f1054, f1071, f1086.

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

Kildegrunnlag: f0804, f0837, f0861, f0886, f0895, f1055.

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

Kildegrunnlag: f0019, f0063, f0124, f0508, f0677, f0843, f0844, f0845, f0846, f0847.

### Tilbydere av kapabilitet: MeasuredPresentation

```mermaid
flowchart LR
    n_MeasuredPresentation["MeasuredPresentation (capability)"]
    n_SduiLayout["SduiLayout (unit)"]
    n_SduiLibrary["SduiLibrary (unit)"]
    n_SduiLayout -->|provides| n_MeasuredPresentation
    n_SduiLibrary -->|provides| n_MeasuredPresentation
```

Kildegrunnlag: f0848, f0862.

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

Kildegrunnlag: f0120, f0215, f0216, f0217, f0218, f0234, f0235, f0456, f0612, f0617, f0621, f0700.

### Tilbydere av kapabilitet: NativeInteraction

```mermaid
flowchart LR
    n_FyneBackend["FyneBackend (unit)"]
    n_FyneHost["FyneHost (container)"]
    n_NativeInteraction["NativeInteraction (capability)"]
    n_FyneBackend -->|provides| n_NativeInteraction
    n_FyneHost -->|provides| n_NativeInteraction
```

Kildegrunnlag: f0219, f0236.

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

Kildegrunnlag: f0054, f0407, f0415, f0426, f0427, f0434, f0435, f0436, f0588, f0690.

### Tilbydere av kapabilitet: NativeRealization

```mermaid
flowchart LR
    n_GoBuildRunner["GoBuildRunner (unit)"]
    n_GoCodeGenerator["GoCodeGenerator (unit)"]
    n_NativeRealization["NativeRealization (capability)"]
    n_GoBuildRunner -->|provides| n_NativeRealization
    n_GoCodeGenerator -->|provides| n_NativeRealization
```

Kildegrunnlag: f0428, f0437.

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

Kildegrunnlag: f0175, f0503, f0504, f0507, f0568, f0569, f0622, f0687, f0688, f1092.

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

Kildegrunnlag: f0131, f0176, f0505, f0689.

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

Kildegrunnlag: f0068, f0132, f0547, f0591, f0683, f0735, f0743, f0759, f0760, f0762, f0791, f0792, f0793, f0956, f1077, f1082.

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

Kildegrunnlag: f0736, f0744, f0749, f0761, f0763, f0794.

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

Kildegrunnlag: f0073, f0133, f0194, f0594, f0597, f0686, f0816, f0853, f0865, f0866, f0867, f0871, f0897, f0898, f0899, f0900, f0959, f1074, f1085, f1097.

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

Kildegrunnlag: f0817, f0854, f0863, f0868, f0872, f0901.

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

Kildegrunnlag: f0170, f0171, f0670, f0675.

### Tilbydere av kapabilitet: SourceDiagnostics

```mermaid
flowchart LR
    n_DiagnosticReporter["DiagnosticReporter (unit)"]
    n_SourceDiagnostics["SourceDiagnostics (capability)"]
    n_DiagnosticReporter -->|provides| n_SourceDiagnostics
```

Kildegrunnlag: f0172.

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

Kildegrunnlag: f0457, f0613, f0916, f0917.

### Tilbydere av kapabilitet: SourceLoading

```mermaid
flowchart LR
    n_SourceLoader["SourceLoader (unit)"]
    n_SourceLoading["SourceLoading (capability)"]
    n_SourceLoader -->|provides| n_SourceLoading
```

Kildegrunnlag: f0918.

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

Kildegrunnlag: f0090, f0091, f0117, f0123, f0195, f0203, f0876, f0877, f0878, f1105.

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

Kildegrunnlag: f0092, f0864, f0879.


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

Kildegrunnlag: f0050, f0051, f0191.

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

Kildegrunnlag: f0153, f0154, f0156, f0157, f0159.

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

Kildegrunnlag: f0155, f0158, f0160.

### Nødvendige porter i modus: NativeBuild

```mermaid
flowchart LR
    n_BuildToolPort["BuildToolPort (interface)"]
    n_GeneratedArtifactPort["GeneratedArtifactPort (interface)"]
    n_NativeRealization["NativeRealization (capability)"]
    n_NativeRealization -->|requires| n_BuildToolPort
    n_NativeRealization -->|requires| n_GeneratedArtifactPort
```

Kildegrunnlag: f0538, f0539.

### Nødvendige porter i modus: RichDocument

```mermaid
flowchart LR
    n_ContentProviderPort["ContentProviderPort (interface)"]
    n_DiagramEnginePort["DiagramEnginePort (interface)"]
    n_RichContent["RichContent (capability)"]
    n_RichContent -->|requires| n_ContentProviderPort
    n_RichContent -->|requires| n_DiagramEnginePort
```

Kildegrunnlag: f0695, f0696.

### Nødvendige porter i modus: SourceInspection

```mermaid
flowchart LR
    n_SdlSourceModel["SdlSourceModel (capability)"]
    n_SduiSourceModel["SduiSourceModel (capability)"]
    n_SourceSnapshotPort["SourceSnapshotPort (interface)"]
    n_SdlSourceModel -->|requires| n_SourceSnapshotPort
    n_SduiSourceModel -->|requires| n_SourceSnapshotPort
```

Kildegrunnlag: f0774, f0896.

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

Kildegrunnlag: f0509, f0929, f0930.

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

Kildegrunnlag: f0510, f0529.


## VP06 — Aktiviteter og leveranseplan

### Aktivitetsrøtter

```mermaid
flowchart LR
    n_G1FrontendPort["G1FrontendPort (activity)"]
    n_G2LayoutAndPresentation["G2LayoutAndPresentation (activity)"]
    n_G3UiRuntimeAndReload["G3UiRuntimeAndReload (activity)"]
    n_G4SdlRuntimeAndBinding["G4SdlRuntimeAndBinding (activity)"]
    n_G5NativeGeneration["G5NativeGeneration (activity)"]
    n_InspectDesignSource["InspectDesignSource (activity)"]
    n_RealizeDesign["RealizeDesign (activity)"]
    n_ReloadDesignSession["ReloadDesignSession (activity)"]
    n_RunDesignSession["RunDesignSession (activity)"]
```

Kildegrunnlag: Kun deklarasjoner.

### Aktivitetsinndeling: G1FrontendPort

```mermaid
flowchart LR
    n_G1FrontendPort["G1FrontendPort (activity)"]
    n_G1M1ParserAndAst["G1M1ParserAndAst (activity)"]
    n_G1M2ValidationAndNormalization["G1M2ValidationAndNormalization (activity)"]
    n_G1M3Concept1AndDumps["G1M3Concept1AndDumps (activity)"]
    n_StructuralModelInspection["StructuralModelInspection (feature)"]
    n_G1FrontendPort -->|delivers| n_StructuralModelInspection
    n_G1M1ParserAndAst -->|refines| n_G1FrontendPort
    n_G1M2ValidationAndNormalization -->|refines| n_G1FrontendPort
    n_G1M3Concept1AndDumps -->|refines| n_G1FrontendPort
```

Kildegrunnlag: f0243, f0250, f0261, f0266.

### Aktivitetsinndeling: G2LayoutAndPresentation

```mermaid
flowchart LR
    n_G2LayoutAndPresentation["G2LayoutAndPresentation (activity)"]
    n_G2M1RelativeMeasurement["G2M1RelativeMeasurement (activity)"]
    n_G2M2SharedSvgGeometry["G2M2SharedSvgGeometry (activity)"]
    n_G2M3FyneInteractions["G2M3FyneInteractions (activity)"]
    n_G2M4RichContent["G2M4RichContent (activity)"]
    n_InteractiveUiPreview["InteractiveUiPreview (feature)"]
    n_G2LayoutAndPresentation -->|delivers| n_InteractiveUiPreview
    n_G2M1RelativeMeasurement -->|refines| n_G2LayoutAndPresentation
    n_G2M2SharedSvgGeometry -->|refines| n_G2LayoutAndPresentation
    n_G2M3FyneInteractions -->|refines| n_G2LayoutAndPresentation
    n_G2M4RichContent -->|refines| n_G2LayoutAndPresentation
```

Kildegrunnlag: f0267, f0275, f0280, f0289, f0297.

### Aktivitetsinndeling: G3UiRuntimeAndReload

```mermaid
flowchart LR
    n_G3M1TypedUiSession["G3M1TypedUiSession (activity)"]
    n_G3M2CandidatePublication["G3M2CandidatePublication (activity)"]
    n_G3M3CompatibleState["G3M3CompatibleState (activity)"]
    n_G3UiRuntimeAndReload["G3UiRuntimeAndReload (activity)"]
    n_LiveModelReload["LiveModelReload (feature)"]
    n_G3M1TypedUiSession -->|refines| n_G3UiRuntimeAndReload
    n_G3M2CandidatePublication -->|refines| n_G3UiRuntimeAndReload
    n_G3M3CompatibleState -->|refines| n_G3UiRuntimeAndReload
    n_G3UiRuntimeAndReload -->|delivers| n_LiveModelReload
```

Kildegrunnlag: f0314, f0323, f0329, f0330.

### Aktivitetsinndeling: G4SdlRuntimeAndBinding

```mermaid
flowchart LR
    n_G4M1SdlFrontend["G4M1SdlFrontend (activity)"]
    n_G4M2TypedExecution["G4M2TypedExecution (activity)"]
    n_G4M3UiDomainBinding["G4M3UiDomainBinding (activity)"]
    n_G4M4DomainReload["G4M4DomainReload (activity)"]
    n_G4SdlRuntimeAndBinding["G4SdlRuntimeAndBinding (activity)"]
    n_TypedDomainBinding["TypedDomainBinding (feature)"]
    n_G4M1SdlFrontend -->|refines| n_G4SdlRuntimeAndBinding
    n_G4M2TypedExecution -->|refines| n_G4SdlRuntimeAndBinding
    n_G4M3UiDomainBinding -->|refines| n_G4SdlRuntimeAndBinding
    n_G4M4DomainReload -->|refines| n_G4SdlRuntimeAndBinding
    n_G4SdlRuntimeAndBinding -->|delivers| n_TypedDomainBinding
```

Kildegrunnlag: f0341, f0356, f0366, f0373, f0374.

### Aktivitetsinndeling: G5NativeGeneration

```mermaid
flowchart LR
    n_DesignDocumentation["DesignDocumentation (feature)"]
    n_G5M1GeneratedGo["G5M1GeneratedGo (activity)"]
    n_G5M2BehaviorParity["G5M2BehaviorParity (activity)"]
    n_G5M3DocumentationExport["G5M3DocumentationExport (activity)"]
    n_G5M4RetirePython["G5M4RetirePython (activity)"]
    n_G5NativeGeneration["G5NativeGeneration (activity)"]
    n_NativeGoAssembly["NativeGoAssembly (feature)"]
    n_G5M1GeneratedGo -->|refines| n_G5NativeGeneration
    n_G5M2BehaviorParity -->|refines| n_G5NativeGeneration
    n_G5M3DocumentationExport -->|refines| n_G5NativeGeneration
    n_G5M4RetirePython -->|refines| n_G5NativeGeneration
    n_G5NativeGeneration -->|delivers| n_DesignDocumentation
    n_G5NativeGeneration -->|delivers| n_NativeGoAssembly
```

Kildegrunnlag: f0382, f0386, f0396, f0401, f0402, f0403.

### Aktivitetsinndeling: InspectDesignSource

```mermaid
flowchart LR
    n_ExportUiDocumentation["ExportUiDocumentation (activity)"]
    n_InspectDesignSource["InspectDesignSource (activity)"]
    n_InspectSdlSource["InspectSdlSource (activity)"]
    n_InspectSduiSource["InspectSduiSource (activity)"]
    n_ExportUiDocumentation -->|refines| n_InspectDesignSource
    n_InspectSdlSource -->|refines| n_InspectDesignSource
    n_InspectSduiSource -->|refines| n_InspectDesignSource
```

Kildegrunnlag: f0204, f0458, f0459.

### Aktivitetsinndeling: RealizeDesign

```mermaid
flowchart LR
    n_BuildNativeRealization["BuildNativeRealization (activity)"]
    n_RealizeDesign["RealizeDesign (activity)"]
    n_BuildNativeRealization -->|refines| n_RealizeDesign
```

Kildegrunnlag: f0060.

### Aktivitetsinndeling: ReloadDesignSession

```mermaid
flowchart LR
    n_ReloadBoundModels["ReloadBoundModels (activity)"]
    n_ReloadDesignSession["ReloadDesignSession (activity)"]
    n_ReloadUiModel["ReloadUiModel (activity)"]
    n_ReloadBoundModels -->|refines| n_ReloadDesignSession
    n_ReloadUiModel -->|refines| n_ReloadDesignSession
```

Kildegrunnlag: f0626, f0669.

### Aktivitetsinndeling: RunDesignSession

```mermaid
flowchart LR
    n_RunBoundUiAction["RunBoundUiAction (activity)"]
    n_RunDesignSession["RunDesignSession (activity)"]
    n_RunUnboundUiPreview["RunUnboundUiPreview (activity)"]
    n_RunBoundUiAction -->|refines| n_RunDesignSession
    n_RunUnboundUiPreview -->|refines| n_RunDesignSession
```

Kildegrunnlag: f0698, f0699.

### Planlagt ansvar: G1M1ParserAndAst

```mermaid
flowchart LR
    n_BuildSduiAst["BuildSduiAst (functionality)"]
    n_G1M1ParserAndAst["G1M1ParserAndAst (activity)"]
    n_IdentifySourceRevision["IdentifySourceRevision (functionality)"]
    n_ReadBoundedSources["ReadBoundedSources (functionality)"]
    n_SduiLexer["SduiLexer (unit)"]
    n_SduiParser["SduiParser (unit)"]
    n_SourceLoader["SourceLoader (unit)"]
    n_TokenizeSduiSource["TokenizeSduiSource (functionality)"]
    n_G1M1ParserAndAst -->|addresses| n_BuildSduiAst
    n_G1M1ParserAndAst -->|addresses| n_IdentifySourceRevision
    n_G1M1ParserAndAst -->|addresses| n_ReadBoundedSources
    n_G1M1ParserAndAst -->|addresses| n_TokenizeSduiSource
    n_SduiLexer -->|owns| n_TokenizeSduiSource
    n_SduiParser -->|owns| n_BuildSduiAst
    n_SourceLoader -->|owns| n_IdentifySourceRevision
    n_SourceLoader -->|owns| n_ReadBoundedSources
```

Kildegrunnlag: f0245, f0246, f0247, f0248, f0853, f0871, f0916, f0917.

### Planlagt ansvar: G1M2ValidationAndNormalization

```mermaid
flowchart LR
    n_CoordinateSduiCompilation["CoordinateSduiCompilation (functionality)"]
    n_ExpandUiDefinitions["ExpandUiDefinitions (functionality)"]
    n_G1M2ValidationAndNormalization["G1M2ValidationAndNormalization (activity)"]
    n_PreserveUiRegions["PreserveUiRegions (functionality)"]
    n_PreserveUiSourceMap["PreserveUiSourceMap (functionality)"]
    n_ResolveUiNames["ResolveUiNames (functionality)"]
    n_SduiFrontend["SduiFrontend (unit)"]
    n_SduiNormalizer["SduiNormalizer (unit)"]
    n_SduiValidator["SduiValidator (unit)"]
    n_ValidateRelativeFormatting["ValidateRelativeFormatting (functionality)"]
    n_ValidateSymbolicBindings["ValidateSymbolicBindings (functionality)"]
    n_ValidateWidgetArguments["ValidateWidgetArguments (functionality)"]
    n_G1M2ValidationAndNormalization -->|addresses| n_CoordinateSduiCompilation
    n_G1M2ValidationAndNormalization -->|addresses| n_ExpandUiDefinitions
    n_G1M2ValidationAndNormalization -->|addresses| n_PreserveUiRegions
    n_G1M2ValidationAndNormalization -->|addresses| n_PreserveUiSourceMap
    n_G1M2ValidationAndNormalization -->|addresses| n_ResolveUiNames
    n_G1M2ValidationAndNormalization -->|addresses| n_ValidateRelativeFormatting
    n_G1M2ValidationAndNormalization -->|addresses| n_ValidateSymbolicBindings
    n_G1M2ValidationAndNormalization -->|addresses| n_ValidateWidgetArguments
    n_SduiFrontend -->|owns| n_CoordinateSduiCompilation
    n_SduiNormalizer -->|owns| n_ExpandUiDefinitions
    n_SduiNormalizer -->|owns| n_PreserveUiRegions
    n_SduiNormalizer -->|owns| n_PreserveUiSourceMap
    n_SduiValidator -->|owns| n_ResolveUiNames
    n_SduiValidator -->|owns| n_ValidateRelativeFormatting
    n_SduiValidator -->|owns| n_ValidateSymbolicBindings
    n_SduiValidator -->|owns| n_ValidateWidgetArguments
```

Kildegrunnlag: f0251, f0252, f0253, f0254, f0255, f0256, f0257, f0258, f0816, f0865, f0866, f0867, f0897, f0898, f0899, f0900.

### Planlagt ansvar: G1M3Concept1AndDumps

```mermaid
flowchart LR
    n_DiagnosticReporter["DiagnosticReporter (unit)"]
    n_ExportConsoleSnapshot["ExportConsoleSnapshot (functionality)"]
    n_G1M3Concept1AndDumps["G1M3Concept1AndDumps (activity)"]
    n_ReportSourceDiagnostics["ReportSourceDiagnostics (functionality)"]
    n_SduiPresentation["SduiPresentation (unit)"]
    n_DiagnosticReporter -->|owns| n_ReportSourceDiagnostics
    n_G1M3Concept1AndDumps -->|addresses| n_ExportConsoleSnapshot
    n_G1M3Concept1AndDumps -->|addresses| n_ReportSourceDiagnostics
    n_SduiPresentation -->|owns| n_ExportConsoleSnapshot
```

Kildegrunnlag: f0171, f0262, f0263, f0877.

### Planlagt ansvar: G2M1RelativeMeasurement

```mermaid
flowchart LR
    n_AllocateGeometry["AllocateGeometry (functionality)"]
    n_ComputeClipping["ComputeClipping (functionality)"]
    n_G2M1RelativeMeasurement["G2M1RelativeMeasurement (activity)"]
    n_MeasureUiContent["MeasureUiContent (functionality)"]
    n_ResolveAncestorDimensions["ResolveAncestorDimensions (functionality)"]
    n_SduiLayout["SduiLayout (unit)"]
    n_G2M1RelativeMeasurement -->|addresses| n_AllocateGeometry
    n_G2M1RelativeMeasurement -->|addresses| n_ComputeClipping
    n_G2M1RelativeMeasurement -->|addresses| n_MeasureUiContent
    n_G2M1RelativeMeasurement -->|addresses| n_ResolveAncestorDimensions
    n_SduiLayout -->|owns| n_AllocateGeometry
    n_SduiLayout -->|owns| n_ComputeClipping
    n_SduiLayout -->|owns| n_MeasureUiContent
    n_SduiLayout -->|owns| n_ResolveAncestorDimensions
```

Kildegrunnlag: f0269, f0270, f0271, f0272, f0843, f0845, f0846, f0847.

### Planlagt ansvar: G2M2SharedSvgGeometry

```mermaid
flowchart LR
    n_BuildPreparedFrame["BuildPreparedFrame (functionality)"]
    n_ExportSvgSnapshot["ExportSvgSnapshot (functionality)"]
    n_G2M2SharedSvgGeometry["G2M2SharedSvgGeometry (activity)"]
    n_SduiLayout["SduiLayout (unit)"]
    n_SduiPresentation["SduiPresentation (unit)"]
    n_G2M2SharedSvgGeometry -->|addresses| n_BuildPreparedFrame
    n_G2M2SharedSvgGeometry -->|addresses| n_ExportSvgSnapshot
    n_SduiLayout -->|owns| n_BuildPreparedFrame
    n_SduiPresentation -->|owns| n_ExportSvgSnapshot
```

Kildegrunnlag: f0276, f0277, f0844, f0878.

### Planlagt ansvar: G2M3FyneInteractions

```mermaid
flowchart LR
    n_ComposeInteractiveSession["ComposeInteractiveSession (functionality)"]
    n_FyneBackend["FyneBackend (unit)"]
    n_FyneHost["FyneHost (container)"]
    n_G2M3FyneInteractions["G2M3FyneInteractions (activity)"]
    n_HandleFocusAndTextInput["HandleFocusAndTextInput (functionality)"]
    n_PublishPresentation["PublishPresentation (functionality)"]
    n_ReconcileWidgets["ReconcileWidgets (functionality)"]
    n_ReleaseNativeWidgets["ReleaseNativeWidgets (functionality)"]
    n_ScheduleUiPublication["ScheduleUiPublication (functionality)"]
    n_FyneBackend -->|owns| n_HandleFocusAndTextInput
    n_FyneBackend -->|owns| n_PublishPresentation
    n_FyneBackend -->|owns| n_ReconcileWidgets
    n_FyneBackend -->|owns| n_ReleaseNativeWidgets
    n_FyneHost -->|owns| n_ComposeInteractiveSession
    n_FyneHost -->|owns| n_ScheduleUiPublication
    n_G2M3FyneInteractions -->|addresses| n_ComposeInteractiveSession
    n_G2M3FyneInteractions -->|addresses| n_HandleFocusAndTextInput
    n_G2M3FyneInteractions -->|addresses| n_PublishPresentation
    n_G2M3FyneInteractions -->|addresses| n_ReconcileWidgets
    n_G2M3FyneInteractions -->|addresses| n_ReleaseNativeWidgets
    n_G2M3FyneInteractions -->|addresses| n_ScheduleUiPublication
```

Kildegrunnlag: f0215, f0216, f0217, f0218, f0234, f0235, f0281, f0282, f0283, f0284, f0285, f0286.

### Planlagt ansvar: G2M4RichContent

```mermaid
flowchart LR
    n_DiagramProvider["DiagramProvider (unit)"]
    n_G2M4RichContent["G2M4RichContent (activity)"]
    n_MarkdownProvider["MarkdownProvider (unit)"]
    n_MeasureMarkdownContent["MeasureMarkdownContent (functionality)"]
    n_PrepareDiagramResource["PrepareDiagramResource (functionality)"]
    n_PrepareMarkdown["PrepareMarkdown (functionality)"]
    n_ReleaseVisualResources["ReleaseVisualResources (functionality)"]
    n_ResourceStore["ResourceStore (unit)"]
    n_ValidateVisualResources["ValidateVisualResources (functionality)"]
    n_DiagramProvider -->|owns| n_PrepareDiagramResource
    n_G2M4RichContent -->|addresses| n_MeasureMarkdownContent
    n_G2M4RichContent -->|addresses| n_PrepareDiagramResource
    n_G2M4RichContent -->|addresses| n_PrepareMarkdown
    n_G2M4RichContent -->|addresses| n_ReleaseVisualResources
    n_G2M4RichContent -->|addresses| n_ValidateVisualResources
    n_MarkdownProvider -->|owns| n_MeasureMarkdownContent
    n_MarkdownProvider -->|owns| n_PrepareMarkdown
    n_ResourceStore -->|owns| n_ReleaseVisualResources
    n_ResourceStore -->|owns| n_ValidateVisualResources
```

Kildegrunnlag: f0175, f0290, f0291, f0292, f0293, f0294, f0503, f0504, f0687, f0688.

### Planlagt ansvar: G3M1TypedUiSession

```mermaid
flowchart LR
    n_ApplyPropertyBatch["ApplyPropertyBatch (functionality)"]
    n_CloseUiInstance["CloseUiInstance (functionality)"]
    n_CorrelateUiResult["CorrelateUiResult (functionality)"]
    n_CreateUiInstance["CreateUiInstance (functionality)"]
    n_DispatchUiEvent["DispatchUiEvent (functionality)"]
    n_G3M1TypedUiSession["G3M1TypedUiSession (activity)"]
    n_ManageWidgetIdentities["ManageWidgetIdentities (functionality)"]
    n_ProjectUiGeneration["ProjectUiGeneration (functionality)"]
    n_RejectStaleUiEvent["RejectStaleUiEvent (functionality)"]
    n_RevokeWidgetGenerations["RevokeWidgetGenerations (functionality)"]
    n_SduiDispatcher["SduiDispatcher (unit)"]
    n_SduiInstanceStore["SduiInstanceStore (unit)"]
    n_SduiPropertyStore["SduiPropertyStore (unit)"]
    n_SduiRuntime["SduiRuntime (unit)"]
    n_SnapshotUiState["SnapshotUiState (functionality)"]
    n_TrackInputDraft["TrackInputDraft (functionality)"]
    n_ValidatePropertyBatch["ValidatePropertyBatch (functionality)"]
    n_ValidateUiEvent["ValidateUiEvent (functionality)"]
    n_G3M1TypedUiSession -->|addresses| n_ApplyPropertyBatch
    n_G3M1TypedUiSession -->|addresses| n_CloseUiInstance
    n_G3M1TypedUiSession -->|addresses| n_CorrelateUiResult
    n_G3M1TypedUiSession -->|addresses| n_CreateUiInstance
    n_G3M1TypedUiSession -->|addresses| n_DispatchUiEvent
    n_G3M1TypedUiSession -->|addresses| n_ManageWidgetIdentities
    n_G3M1TypedUiSession -->|addresses| n_ProjectUiGeneration
    n_G3M1TypedUiSession -->|addresses| n_RejectStaleUiEvent
    n_G3M1TypedUiSession -->|addresses| n_RevokeWidgetGenerations
    n_G3M1TypedUiSession -->|addresses| n_SnapshotUiState
    n_G3M1TypedUiSession -->|addresses| n_TrackInputDraft
    n_G3M1TypedUiSession -->|addresses| n_ValidatePropertyBatch
    n_G3M1TypedUiSession -->|addresses| n_ValidateUiEvent
    n_SduiDispatcher -->|owns| n_CorrelateUiResult
    n_SduiDispatcher -->|owns| n_DispatchUiEvent
    n_SduiDispatcher -->|owns| n_RejectStaleUiEvent
    n_SduiDispatcher -->|owns| n_ValidateUiEvent
    n_SduiInstanceStore -->|owns| n_ManageWidgetIdentities
    n_SduiInstanceStore -->|owns| n_ProjectUiGeneration
    n_SduiInstanceStore -->|owns| n_RevokeWidgetGenerations
    n_SduiInstanceStore -->|owns| n_SnapshotUiState
    n_SduiPropertyStore -->|owns| n_ApplyPropertyBatch
    n_SduiPropertyStore -->|owns| n_TrackInputDraft
    n_SduiPropertyStore -->|owns| n_ValidatePropertyBatch
    n_SduiRuntime -->|owns| n_CloseUiInstance
    n_SduiRuntime -->|owns| n_CreateUiInstance
```

Kildegrunnlag: f0298, f0299, f0300, f0301, f0302, f0303, f0304, f0305, f0306, f0307, f0308, f0309, f0310, f0800, f0801, f0802, f0803, f0833, f0834, f0835, f0836, f0883, f0884, f0885, f0893, f0894.

### Planlagt ansvar: G3M2CandidatePublication

```mermaid
flowchart LR
    n_CoalesceSourceChanges["CoalesceSourceChanges (functionality)"]
    n_G3M2CandidatePublication["G3M2CandidatePublication (activity)"]
    n_KeepLastValidModels["KeepLastValidModels (functionality)"]
    n_ObserveSourceChanges["ObserveSourceChanges (functionality)"]
    n_PrepareCandidateModels["PrepareCandidateModels (functionality)"]
    n_PublishModelGeneration["PublishModelGeneration (functionality)"]
    n_ReloadCoordinator["ReloadCoordinator (unit)"]
    n_RetirePreviousGeneration["RetirePreviousGeneration (functionality)"]
    n_SourceWatcher["SourceWatcher (unit)"]
    n_G3M2CandidatePublication -->|addresses| n_CoalesceSourceChanges
    n_G3M2CandidatePublication -->|addresses| n_KeepLastValidModels
    n_G3M2CandidatePublication -->|addresses| n_ObserveSourceChanges
    n_G3M2CandidatePublication -->|addresses| n_PrepareCandidateModels
    n_G3M2CandidatePublication -->|addresses| n_PublishModelGeneration
    n_G3M2CandidatePublication -->|addresses| n_RetirePreviousGeneration
    n_ReloadCoordinator -->|owns| n_KeepLastValidModels
    n_ReloadCoordinator -->|owns| n_PrepareCandidateModels
    n_ReloadCoordinator -->|owns| n_PublishModelGeneration
    n_ReloadCoordinator -->|owns| n_RetirePreviousGeneration
    n_SourceWatcher -->|owns| n_CoalesceSourceChanges
    n_SourceWatcher -->|owns| n_ObserveSourceChanges
```

Kildegrunnlag: f0315, f0316, f0317, f0318, f0319, f0320, f0634, f0635, f0636, f0637, f0920, f0921.

### Planlagt ansvar: G3M3CompatibleState

```mermaid
flowchart LR
    n_G3M3CompatibleState["G3M3CompatibleState (activity)"]
    n_MatchCompatibleWidgets["MatchCompatibleWidgets (functionality)"]
    n_PreserveCompatibleUiState["PreserveCompatibleUiState (functionality)"]
    n_ResetIncompatibleUiState["ResetIncompatibleUiState (functionality)"]
    n_UiStateReconciler["UiStateReconciler (unit)"]
    n_G3M3CompatibleState -->|addresses| n_MatchCompatibleWidgets
    n_G3M3CompatibleState -->|addresses| n_PreserveCompatibleUiState
    n_G3M3CompatibleState -->|addresses| n_ResetIncompatibleUiState
    n_UiStateReconciler -->|owns| n_MatchCompatibleWidgets
    n_UiStateReconciler -->|owns| n_PreserveCompatibleUiState
    n_UiStateReconciler -->|owns| n_ResetIncompatibleUiState
```

Kildegrunnlag: f0324, f0325, f0326, f1052, f1053, f1054.

### Planlagt ansvar: G4M1SdlFrontend

```mermaid
flowchart LR
    n_BuildSdlAst["BuildSdlAst (functionality)"]
    n_CoordinateSdlCompilation["CoordinateSdlCompilation (functionality)"]
    n_G4M1SdlFrontend["G4M1SdlFrontend (activity)"]
    n_NormalizeSdlModel["NormalizeSdlModel (functionality)"]
    n_PreserveSdlSourceMap["PreserveSdlSourceMap (functionality)"]
    n_ResolveSdlSymbols["ResolveSdlSymbols (functionality)"]
    n_SdlFrontend["SdlFrontend (unit)"]
    n_SdlLexer["SdlLexer (unit)"]
    n_SdlNormalizer["SdlNormalizer (unit)"]
    n_SdlParser["SdlParser (unit)"]
    n_SdlValidator["SdlValidator (unit)"]
    n_TokenizeSdlSource["TokenizeSdlSource (functionality)"]
    n_ValidateSdlProfile["ValidateSdlProfile (functionality)"]
    n_ValidateSdlStructure["ValidateSdlStructure (functionality)"]
    n_G4M1SdlFrontend -->|addresses| n_BuildSdlAst
    n_G4M1SdlFrontend -->|addresses| n_CoordinateSdlCompilation
    n_G4M1SdlFrontend -->|addresses| n_NormalizeSdlModel
    n_G4M1SdlFrontend -->|addresses| n_PreserveSdlSourceMap
    n_G4M1SdlFrontend -->|addresses| n_ResolveSdlSymbols
    n_G4M1SdlFrontend -->|addresses| n_TokenizeSdlSource
    n_G4M1SdlFrontend -->|addresses| n_ValidateSdlProfile
    n_G4M1SdlFrontend -->|addresses| n_ValidateSdlStructure
    n_SdlFrontend -->|owns| n_CoordinateSdlCompilation
    n_SdlLexer -->|owns| n_TokenizeSdlSource
    n_SdlNormalizer -->|owns| n_NormalizeSdlModel
    n_SdlNormalizer -->|owns| n_PreserveSdlSourceMap
    n_SdlParser -->|owns| n_BuildSdlAst
    n_SdlValidator -->|owns| n_ResolveSdlSymbols
    n_SdlValidator -->|owns| n_ValidateSdlProfile
    n_SdlValidator -->|owns| n_ValidateSdlStructure
```

Kildegrunnlag: f0332, f0333, f0334, f0335, f0336, f0337, f0338, f0339, f0735, f0743, f0759, f0760, f0762, f0791, f0792, f0793.

### Planlagt ansvar: G4M2TypedExecution

```mermaid
flowchart LR
    n_CancelPendingActions["CancelPendingActions (functionality)"]
    n_CheckExecutionCompleteness["CheckExecutionCompleteness (functionality)"]
    n_CheckFunctionSignatures["CheckFunctionSignatures (functionality)"]
    n_CloseSdlInstance["CloseSdlInstance (functionality)"]
    n_CorrelateActionResult["CorrelateActionResult (functionality)"]
    n_CreateSdlInstance["CreateSdlInstance (functionality)"]
    n_G4M2TypedExecution["G4M2TypedExecution (activity)"]
    n_GoDomainImplementation["GoDomainImplementation (unit)"]
    n_InvokeRegisteredFunction["InvokeRegisteredFunction (functionality)"]
    n_ManageDomainState["ManageDomainState (functionality)"]
    n_PerformDomainOperation["PerformDomainOperation (functionality)"]
    n_RegisterDomainFunctions["RegisterDomainFunctions (functionality)"]
    n_SdlDispatcher["SdlDispatcher (unit)"]
    n_SdlExecutionGate["SdlExecutionGate (unit)"]
    n_SdlFunctionRegistry["SdlFunctionRegistry (unit)"]
    n_SdlRuntime["SdlRuntime (unit)"]
    n_SdlStateStore["SdlStateStore (unit)"]
    n_SnapshotDomainState["SnapshotDomainState (functionality)"]
    n_ValidateActionInput["ValidateActionInput (functionality)"]
    n_G4M2TypedExecution -->|addresses| n_CancelPendingActions
    n_G4M2TypedExecution -->|addresses| n_CheckExecutionCompleteness
    n_G4M2TypedExecution -->|addresses| n_CheckFunctionSignatures
    n_G4M2TypedExecution -->|addresses| n_CloseSdlInstance
    n_G4M2TypedExecution -->|addresses| n_CorrelateActionResult
    n_G4M2TypedExecution -->|addresses| n_CreateSdlInstance
    n_G4M2TypedExecution -->|addresses| n_InvokeRegisteredFunction
    n_G4M2TypedExecution -->|addresses| n_ManageDomainState
    n_G4M2TypedExecution -->|addresses| n_PerformDomainOperation
    n_G4M2TypedExecution -->|addresses| n_RegisterDomainFunctions
    n_G4M2TypedExecution -->|addresses| n_SnapshotDomainState
    n_G4M2TypedExecution -->|addresses| n_ValidateActionInput
    n_GoDomainImplementation -->|owns| n_PerformDomainOperation
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
```

Kildegrunnlag: f0342, f0343, f0344, f0345, f0346, f0347, f0348, f0349, f0350, f0351, f0352, f0353, f0444, f0716, f0717, f0718, f0719, f0727, f0740, f0741, f0771, f0772, f0775, f0776.

### Planlagt ansvar: G4M3UiDomainBinding

```mermaid
flowchart LR
    n_ConnectTypedWidgetHandles["ConnectTypedWidgetHandles (functionality)"]
    n_DiagnosticReporter["DiagnosticReporter (unit)"]
    n_DisconnectBindings["DisconnectBindings (functionality)"]
    n_G4M3UiDomainBinding["G4M3UiDomainBinding (activity)"]
    n_PublishDomainUpdates["PublishDomainUpdates (functionality)"]
    n_ReportBindingDiagnostics["ReportBindingDiagnostics (functionality)"]
    n_ResolveCallbackSymbols["ResolveCallbackSymbols (functionality)"]
    n_RouteDomainBindings["RouteDomainBindings (functionality)"]
    n_SdlUiBindingAdapter["SdlUiBindingAdapter (unit)"]
    n_DiagnosticReporter -->|owns| n_ReportBindingDiagnostics
    n_G4M3UiDomainBinding -->|addresses| n_ConnectTypedWidgetHandles
    n_G4M3UiDomainBinding -->|addresses| n_DisconnectBindings
    n_G4M3UiDomainBinding -->|addresses| n_PublishDomainUpdates
    n_G4M3UiDomainBinding -->|addresses| n_ReportBindingDiagnostics
    n_G4M3UiDomainBinding -->|addresses| n_ResolveCallbackSymbols
    n_G4M3UiDomainBinding -->|addresses| n_RouteDomainBindings
    n_SdlUiBindingAdapter -->|owns| n_ConnectTypedWidgetHandles
    n_SdlUiBindingAdapter -->|owns| n_DisconnectBindings
    n_SdlUiBindingAdapter -->|owns| n_PublishDomainUpdates
    n_SdlUiBindingAdapter -->|owns| n_ResolveCallbackSymbols
    n_SdlUiBindingAdapter -->|owns| n_RouteDomainBindings
```

Kildegrunnlag: f0170, f0357, f0358, f0359, f0360, f0361, f0362, f0781, f0782, f0783, f0784, f0785.

### Planlagt ansvar: G4M4DomainReload

```mermaid
flowchart LR
    n_CheckDomainStateCompatibility["CheckDomainStateCompatibility (functionality)"]
    n_DomainStateMigrator["DomainStateMigrator (unit)"]
    n_G4M4DomainReload["G4M4DomainReload (activity)"]
    n_GoBuildRunner["GoBuildRunner (unit)"]
    n_MigrateOrResetDomainState["MigrateOrResetDomainState (functionality)"]
    n_RestartChangedGoProgram["RestartChangedGoProgram (functionality)"]
    n_DomainStateMigrator -->|owns| n_CheckDomainStateCompatibility
    n_DomainStateMigrator -->|owns| n_MigrateOrResetDomainState
    n_G4M4DomainReload -->|addresses| n_CheckDomainStateCompatibility
    n_G4M4DomainReload -->|addresses| n_MigrateOrResetDomainState
    n_G4M4DomainReload -->|addresses| n_RestartChangedGoProgram
    n_GoBuildRunner -->|owns| n_RestartChangedGoProgram
```

Kildegrunnlag: f0188, f0189, f0367, f0368, f0369, f0427.

### Planlagt ansvar: G5M1GeneratedGo

```mermaid
flowchart LR
    n_BuildGeneratedApplication["BuildGeneratedApplication (functionality)"]
    n_G5M1GeneratedGo["G5M1GeneratedGo (activity)"]
    n_GenerateBindingRegistration["GenerateBindingRegistration (functionality)"]
    n_GenerateModelConstructors["GenerateModelConstructors (functionality)"]
    n_GoBuildRunner["GoBuildRunner (unit)"]
    n_GoCodeGenerator["GoCodeGenerator (unit)"]
    n_PreserveHandwrittenSources["PreserveHandwrittenSources (functionality)"]
    n_G5M1GeneratedGo -->|addresses| n_BuildGeneratedApplication
    n_G5M1GeneratedGo -->|addresses| n_GenerateBindingRegistration
    n_G5M1GeneratedGo -->|addresses| n_GenerateModelConstructors
    n_G5M1GeneratedGo -->|addresses| n_PreserveHandwrittenSources
    n_GoBuildRunner -->|owns| n_BuildGeneratedApplication
    n_GoCodeGenerator -->|owns| n_GenerateBindingRegistration
    n_GoCodeGenerator -->|owns| n_GenerateModelConstructors
    n_GoCodeGenerator -->|owns| n_PreserveHandwrittenSources
```

Kildegrunnlag: f0376, f0377, f0378, f0379, f0426, f0434, f0435, f0436.

### Planlagt ansvar: G5M2BehaviorParity

```mermaid
flowchart LR
    n_DevelopmentTools["DevelopmentTools (unit)"]
    n_G5M2BehaviorParity["G5M2BehaviorParity (activity)"]
    n_VerifyNativeBehaviorParity["VerifyNativeBehaviorParity (functionality)"]
    n_DevelopmentTools -->|owns| n_VerifyNativeBehaviorParity
    n_G5M2BehaviorParity -->|addresses| n_VerifyNativeBehaviorParity
```

Kildegrunnlag: f0168, f0383.

### Planlagt ansvar: G5M3DocumentationExport

```mermaid
flowchart LR
    n_CommandLineHost["CommandLineHost (container)"]
    n_ComposeHeadlessExport["ComposeHeadlessExport (functionality)"]
    n_ComposeMarkdownDocument["ComposeMarkdownDocument (functionality)"]
    n_ExportViewpointMarkdown["ExportViewpointMarkdown (functionality)"]
    n_G5M3DocumentationExport["G5M3DocumentationExport (activity)"]
    n_ProjectSdlViewpoints["ProjectSdlViewpoints (functionality)"]
    n_SdlViewpointGenerator["SdlViewpointGenerator (unit)"]
    n_SduiPresentation["SduiPresentation (unit)"]
    n_TraceViewpointFacts["TraceViewpointFacts (functionality)"]
    n_WriteGeneratedArtifacts["WriteGeneratedArtifacts (functionality)"]
    n_CommandLineHost -->|owns| n_ComposeHeadlessExport
    n_CommandLineHost -->|owns| n_WriteGeneratedArtifacts
    n_G5M3DocumentationExport -->|addresses| n_ComposeHeadlessExport
    n_G5M3DocumentationExport -->|addresses| n_ComposeMarkdownDocument
    n_G5M3DocumentationExport -->|addresses| n_ExportViewpointMarkdown
    n_G5M3DocumentationExport -->|addresses| n_ProjectSdlViewpoints
    n_G5M3DocumentationExport -->|addresses| n_TraceViewpointFacts
    n_G5M3DocumentationExport -->|addresses| n_WriteGeneratedArtifacts
    n_SdlViewpointGenerator -->|owns| n_ExportViewpointMarkdown
    n_SdlViewpointGenerator -->|owns| n_ProjectSdlViewpoints
    n_SdlViewpointGenerator -->|owns| n_TraceViewpointFacts
    n_SduiPresentation -->|owns| n_ComposeMarkdownDocument
```

Kildegrunnlag: f0090, f0091, f0387, f0388, f0389, f0390, f0391, f0392, f0795, f0796, f0797, f0876.

### Planlagt ansvar: G5M4RetirePython

```mermaid
flowchart LR
    n_DevelopmentTools["DevelopmentTools (unit)"]
    n_G5M4RetirePython["G5M4RetirePython (activity)"]
    n_RetireReplacedPythonEntryPoints["RetireReplacedPythonEntryPoints (functionality)"]
    n_DevelopmentTools -->|owns| n_RetireReplacedPythonEntryPoints
    n_G5M4RetirePython -->|addresses| n_RetireReplacedPythonEntryPoints
```

Kildegrunnlag: f0167, f0397.

### Eksplisitte aktivitetsavhengigheter

```mermaid
flowchart LR
    n_G1M1ParserAndAst["G1M1ParserAndAst (activity)"]
    n_G1M2ValidationAndNormalization["G1M2ValidationAndNormalization (activity)"]
    n_G1M3Concept1AndDumps["G1M3Concept1AndDumps (activity)"]
    n_G2M1RelativeMeasurement["G2M1RelativeMeasurement (activity)"]
    n_G2M2SharedSvgGeometry["G2M2SharedSvgGeometry (activity)"]
    n_G2M3FyneInteractions["G2M3FyneInteractions (activity)"]
    n_G2M4RichContent["G2M4RichContent (activity)"]
    n_G3M1TypedUiSession["G3M1TypedUiSession (activity)"]
    n_G3M2CandidatePublication["G3M2CandidatePublication (activity)"]
    n_G3M3CompatibleState["G3M3CompatibleState (activity)"]
    n_G4M1SdlFrontend["G4M1SdlFrontend (activity)"]
    n_G4M2TypedExecution["G4M2TypedExecution (activity)"]
    n_G4M3UiDomainBinding["G4M3UiDomainBinding (activity)"]
    n_G4M4DomainReload["G4M4DomainReload (activity)"]
    n_G5M1GeneratedGo["G5M1GeneratedGo (activity)"]
    n_G5M2BehaviorParity["G5M2BehaviorParity (activity)"]
    n_G5M3DocumentationExport["G5M3DocumentationExport (activity)"]
    n_G5M4RetirePython["G5M4RetirePython (activity)"]
    n_G1M2ValidationAndNormalization -->|depends-on| n_G1M1ParserAndAst
    n_G1M3Concept1AndDumps -->|depends-on| n_G1M2ValidationAndNormalization
    n_G2M1RelativeMeasurement -->|depends-on| n_G1M2ValidationAndNormalization
    n_G2M2SharedSvgGeometry -->|depends-on| n_G2M1RelativeMeasurement
    n_G2M3FyneInteractions -->|depends-on| n_G2M2SharedSvgGeometry
    n_G2M4RichContent -->|depends-on| n_G2M3FyneInteractions
    n_G3M1TypedUiSession -->|depends-on| n_G1M2ValidationAndNormalization
    n_G3M1TypedUiSession -->|depends-on| n_G2M3FyneInteractions
    n_G3M2CandidatePublication -->|depends-on| n_G3M1TypedUiSession
    n_G3M3CompatibleState -->|depends-on| n_G3M2CandidatePublication
    n_G4M2TypedExecution -->|depends-on| n_G4M1SdlFrontend
    n_G4M3UiDomainBinding -->|depends-on| n_G3M1TypedUiSession
    n_G4M3UiDomainBinding -->|depends-on| n_G4M2TypedExecution
    n_G4M4DomainReload -->|depends-on| n_G3M3CompatibleState
    n_G4M4DomainReload -->|depends-on| n_G4M3UiDomainBinding
    n_G5M1GeneratedGo -->|depends-on| n_G4M4DomainReload
    n_G5M2BehaviorParity -->|depends-on| n_G5M1GeneratedGo
    n_G5M3DocumentationExport -->|depends-on| n_G2M4RichContent
    n_G5M3DocumentationExport -->|depends-on| n_G5M2BehaviorParity
    n_G5M4RetirePython -->|depends-on| n_G1M3Concept1AndDumps
    n_G5M4RetirePython -->|depends-on| n_G5M3DocumentationExport
```

Kildegrunnlag: f0259, f0264, f0273, f0278, f0287, f0295, f0311, f0312, f0321, f0327, f0354, f0363, f0364, f0370, f0371, f0380, f0384, f0393, f0394, f0398, f0399.


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

Kildegrunnlag: f0091, f0122, f0202, f0206, f0599, f0795, f0796, f0797, f0876, f0878, f0965, f0967, f1104.

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

Kildegrunnlag: f0091, f0121, f0122, f0201, f0202, f0205, f0206, f0598, f0599, f0795, f0796, f0797, f0876, f0878, f0966, f0967, f1103, f1104.

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

Kildegrunnlag: f0017, f0018, f0061, f0062, f0118, f0119, f0215, f0217, f0234, f0454, f0455, f0614, f0615, f0843, f0844.

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

Kildegrunnlag: f0171, f0473, f0474, f0561, f0562, f0565, f0566, f0583, f0584, f0602, f0603, f0608, f0609, f0634, f0635, f0636, f0672, f0673, f0834, f0921, f1053.

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

Kildegrunnlag: f0171, f0474, f0562, f0566, f0584, f0603, f0609, f0634, f0635, f0636, f0671, f0673, f0834, f0921, f1053.

### Feature: NativeGoAssembly — modus NativeBuild

```mermaid
flowchart LR
    n_BuildGeneratedApplication["BuildGeneratedApplication (functionality)"]
    n_CommandLineHost["CommandLineHost (container)"]
    n_DevelopmentTools["DevelopmentTools (unit)"]
    n_GenerateBindingRegistration["GenerateBindingRegistration (functionality)"]
    n_GenerateModelConstructors["GenerateModelConstructors (functionality)"]
    n_GoBuildRunner["GoBuildRunner (unit)"]
    n_GoCodeGenerator["GoCodeGenerator (unit)"]
    n_NativeGoAssembly["NativeGoAssembly (feature)"]
    n_PreserveHandwrittenSources["PreserveHandwrittenSources (functionality)"]
    n_RetireReplacedPythonEntryPoints["RetireReplacedPythonEntryPoints (functionality)"]
    n_VerifyNativeBehaviorParity["VerifyNativeBehaviorParity (functionality)"]
    n_BuildGeneratedApplication -->|allocated-to| n_CommandLineHost
    n_BuildGeneratedApplication -->|contributes-to| n_NativeGoAssembly
    n_DevelopmentTools -->|owns| n_RetireReplacedPythonEntryPoints
    n_DevelopmentTools -->|owns| n_VerifyNativeBehaviorParity
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
    n_RetireReplacedPythonEntryPoints -->|allocated-to| n_CommandLineHost
    n_RetireReplacedPythonEntryPoints -->|contributes-to| n_NativeGoAssembly
    n_VerifyNativeBehaviorParity -->|allocated-to| n_CommandLineHost
    n_VerifyNativeBehaviorParity -->|contributes-to| n_NativeGoAssembly
```

Kildegrunnlag: f0052, f0053, f0167, f0168, f0405, f0406, f0413, f0414, f0426, f0434, f0435, f0436, f0586, f0587, f0692, f0693, f1101, f1102.

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

Kildegrunnlag: f0065, f0070, f0171, f0672, f0674, f0762, f0793, f0871, f0900, f1079, f1094.

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

Kildegrunnlag: f0064, f0065, f0069, f0070, f0171, f0671, f0674, f0762, f0793, f0871, f0900, f1078, f1079, f1093, f1094.

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

Kildegrunnlag: f0125, f0126, f0178, f0179, f0470, f0471, f0605, f0606, f0678, f0679, f0718, f0781, f0783, f0784, f0801.

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


## VP08 — Channel-kontrakter og sekvenser

### Scenario: BoundActionAccepted — modus BoundExecution

```mermaid
sequenceDiagram
    participant n_FyneBackend as Fyne Backend
    participant n_SduiDispatcher as Sdui Dispatcher
    participant n_SdlUiBindingAdapter as Sdl Ui Binding Adapter
    participant n_SdlDispatcher as Sdl Dispatcher
    participant n_GoDomainImplementation as Go Domain Implementation
    participant n_SduiInstanceStore as Sdui Instance Store
    n_FyneBackend->>n_SduiDispatcher: 1: Ui Action Request (Native Ui Actions)
    n_SduiDispatcher->>n_SdlUiBindingAdapter: 2: Bound Action Request (Ui Domain Actions)
    n_SdlUiBindingAdapter->>n_SdlDispatcher: 3: Sdl Action Request (Sdl Action Calls)
    n_SdlDispatcher->>n_GoDomainImplementation: 4: Domain Action Request (Go Domain Calls)
    n_GoDomainImplementation-->>n_SdlDispatcher: 5: Domain Action Result (Go Domain Calls) reply-to 4
    n_SdlDispatcher-->>n_SdlUiBindingAdapter: 6: Sdl Action Result (Sdl Action Calls) reply-to 3
    n_SdlUiBindingAdapter-->>n_SduiDispatcher: 7: Bound Action Result (Ui Domain Actions) reply-to 2
    n_SduiDispatcher-->>n_FyneBackend: 8: Ui Action Result (Native Ui Actions) reply-to 1
    n_SduiInstanceStore->>n_FyneBackend: 9: Ui Generation Notices / Ui Generation Changed (Ui Generation Events)
```

Kildegrunnlag: f0026, f0027, f0028, f0029, f0030, f0031, f0032, f0033, f0034, f0035, f0036, f0037, f0038, f0045, f0046, f0047, f0048, f0049, f0181, f0182, f0183, f0184, f0185, f0223, f0224, f0225, f0440, f0442, f0443, f0446, f0448, f0540, f0543, f0544, f0701, f0703, f0704, f0705, f0706, f0707, f0708, f0709, f0721, f0722, f0723, f0724, f0787, f0788, f0789, f0790, f0805, f0807, f0808, f0809, f0838, f0977, f0978, f0979, f0980, f0981, f1006, f1008, f1009, f1014, f1016, f1018.

### Scenario: BoundActionRejected — modus BoundExecution

```mermaid
sequenceDiagram
    participant n_FyneBackend as Fyne Backend
    participant n_SduiDispatcher as Sdui Dispatcher
    n_FyneBackend->>n_SduiDispatcher: 1: Ui Action Request (Native Ui Actions)
    n_SduiDispatcher-->>n_FyneBackend: 2: Ui Action Rejected (Native Ui Actions) reply-to 1
```

Kildegrunnlag: f0039, f0040, f0041, f0042, f0043, f0044, f0222, f0224, f0540, f0542, f0543, f0805, f0806, f0974, f0975, f0976, f0977, f0978.

### Scenario: InteractiveFramePrepared — modus UiPreview

```mermaid
sequenceDiagram
    participant n_FyneHost as Fyne Host
    participant n_SduiLayout as Sdui Layout
    participant n_FyneBackend as Fyne Backend
    n_FyneHost->>n_SduiLayout: 1: Layout Request (Layout Calls)
    n_SduiLayout-->>n_FyneHost: 2: Layout Result (Layout Calls) reply-to 1
    n_FyneHost->>n_FyneBackend: 3: Present Frame Request (Frame Presentation Calls)
    n_FyneBackend-->>n_FyneHost: 4: Present Frame Result (Frame Presentation Calls) reply-to 3
```

Kildegrunnlag: f0209, f0211, f0212, f0220, f0221, f0237, f0238, f0241, f0242, f0460, f0461, f0462, f0463, f0464, f0465, f0466, f0467, f0468, f0480, f0482, f0483, f0486, f0487, f0488, f0489, f0490, f0574, f0575, f0576, f0577, f0578, f0850, f0852.

### Scenario: NativeProgramBuilt — modus NativeBuild

```mermaid
sequenceDiagram
    participant n_CommandLineHost as Command Line Host
    participant n_GoCodeGenerator as Go Code Generator
    participant n_GoBuildRunner as Go Build Runner
    n_CommandLineHost->>n_GoCodeGenerator: 1: Generate Go Request (Go Generation Calls)
    n_GoCodeGenerator-->>n_CommandLineHost: 2: Generate Go Result (Go Generation Calls) reply-to 1
    n_CommandLineHost->>n_GoBuildRunner: 3: Build Go Request (Go Build Calls)
    n_GoBuildRunner-->>n_CommandLineHost: 4: Build Go Result (Go Build Calls) reply-to 3
```

Kildegrunnlag: f0055, f0056, f0057, f0058, f0059, f0093, f0094, f0095, f0096, f0408, f0409, f0410, f0411, f0412, f0420, f0422, f0423, f0429, f0430, f0438, f0439, f0450, f0452, f0453, f0530, f0531, f0532, f0533, f0534, f0535, f0536, f0537.

### Scenario: SdlModelReloadAccepted — modus BoundLiveEditing

```mermaid
sequenceDiagram
    participant n_SourceWatcher as Source Watcher
    participant n_ReloadCoordinator as Reload Coordinator
    participant n_SdlFrontend as Sdl Frontend
    n_SourceWatcher->>n_ReloadCoordinator: 1: Reload Request (Model Reload Calls)
    n_ReloadCoordinator->>n_SdlFrontend: 2: Compile Sdl Request (Sdl Compilation Calls)
    n_SdlFrontend-->>n_ReloadCoordinator: 3: Compile Sdl Result (Sdl Compilation Calls) reply-to 2
    n_ReloadCoordinator-->>n_SourceWatcher: 4: Reload Published (Model Reload Calls) reply-to 1
```

Kildegrunnlag: f0104, f0105, f0106, f0107, f0108, f0516, f0518, f0520, f0639, f0641, f0645, f0646, f0655, f0656, f0657, f0663, f0664, f0710, f0712, f0713, f0737, f0738, f0750, f0751, f0752, f0753, f0754, f0755, f0756, f0757, f0758, f0923, f0927.

### Scenario: StaticFrameExported — modus StaticExport

```mermaid
sequenceDiagram
    participant n_CommandLineHost as Command Line Host
    participant n_SduiLayout as Sdui Layout
    participant n_SduiPresentation as Sdui Presentation
    n_CommandLineHost->>n_SduiLayout: 1: Layout Request (Layout Calls)
    n_SduiLayout-->>n_CommandLineHost: 2: Layout Result (Layout Calls) reply-to 1
    n_CommandLineHost->>n_SduiPresentation: 3: Export Svg Request (Svg Export Calls)
    n_SduiPresentation-->>n_CommandLineHost: 4: Export Svg Result (Svg Export Calls) reply-to 3
```

Kildegrunnlag: f0097, f0098, f0099, f0100, f0196, f0197, f0198, f0199, f0200, f0480, f0482, f0483, f0486, f0487, f0488, f0489, f0490, f0849, f0851, f0880, f0881, f0931, f0932, f0933, f0934, f0935, f0936, f0937, f0938, f0939, f0946, f0948, f0949.

### Scenario: UiCompilationAccepted — modus SourceInspection

```mermaid
sequenceDiagram
    participant n_CommandLineHost as Command Line Host
    participant n_SduiFrontend as Sdui Frontend
    participant n_SduiLexer as Sdui Lexer
    participant n_SduiParser as Sdui Parser
    participant n_SduiValidator as Sdui Validator
    participant n_SduiNormalizer as Sdui Normalizer
    n_CommandLineHost->>n_SduiFrontend: 1: Compile Ui Request (Ui Compilation Calls)
    n_SduiFrontend->>n_SduiLexer: 2: Tokenize Ui Request (Ui Tokenization Calls)
    n_SduiLexer-->>n_SduiFrontend: 3: Tokenize Ui Result (Ui Tokenization Calls) reply-to 2
    n_SduiFrontend->>n_SduiParser: 4: Build Ui Ast Request (Ui Ast Calls)
    n_SduiParser-->>n_SduiFrontend: 5: Build Ui Ast Result (Ui Ast Calls) reply-to 4
    n_SduiFrontend->>n_SduiValidator: 6: Validate Ui Request (Ui Validation Calls)
    n_SduiValidator-->>n_SduiFrontend: 7: Validate Ui Result (Ui Validation Calls) reply-to 6
    n_SduiFrontend->>n_SduiNormalizer: 8: Normalize Ui Request (Ui Normalization Calls)
    n_SduiNormalizer-->>n_SduiFrontend: 9: Normalize Ui Result (Ui Normalization Calls) reply-to 8
    n_SduiFrontend-->>n_CommandLineHost: 10: Compile Ui Result (Ui Compilation Calls) reply-to 1
```

Kildegrunnlag: f0074, f0075, f0076, f0077, f0078, f0102, f0103, f0112, f0113, f0114, f0115, f0116, f0548, f0549, f0550, f0551, f0552, f0818, f0819, f0821, f0825, f0826, f0827, f0828, f0829, f0830, f0831, f0855, f0856, f0869, f0870, f0873, f0874, f0902, f0903, f0960, f0961, f0962, f0963, f0964, f0982, f0984, f0985, f0986, f0987, f0988, f0989, f0990, f0991, f0992, f0993, f0994, f0995, f0996, f0997, f0998, f0999, f1000, f1001, f1004, f1005, f1042, f1044, f1045, f1056, f1058, f1059, f1060, f1062, f1063, f1087, f1088, f1089, f1090, f1091.

### Scenario: UiModelReloadAccepted — modus LiveEditing

```mermaid
sequenceDiagram
    participant n_SourceWatcher as Source Watcher
    participant n_ReloadCoordinator as Reload Coordinator
    participant n_SduiFrontend as Sdui Frontend
    participant n_SduiInstanceStore as Sdui Instance Store
    participant n_FyneBackend as Fyne Backend
    n_SourceWatcher->>n_ReloadCoordinator: 1: Reload Request (Model Reload Calls)
    n_ReloadCoordinator->>n_SduiFrontend: 2: Compile Ui Request (Ui Compilation Calls)
    n_SduiFrontend-->>n_ReloadCoordinator: 3: Compile Ui Result (Ui Compilation Calls) reply-to 2
    n_ReloadCoordinator-->>n_SourceWatcher: 4: Reload Published (Model Reload Calls) reply-to 1
    n_SduiInstanceStore->>n_FyneBackend: 5: Ui Generation Notices / Ui Generation Changed (Ui Generation Events)
```

Kildegrunnlag: f0112, f0113, f0114, f0115, f0116, f0226, f0516, f0518, f0520, f0640, f0642, f0648, f0649, f0655, f0656, f0657, f0663, f0664, f0820, f0824, f0839, f0924, f0928, f1001, f1004, f1005, f1014, f1016, f1018, f1024, f1025, f1026, f1027, f1028, f1029, f1030, f1031, f1032, f1033.

### Scenario: UiModelReloadRejected — modus LiveEditing

```mermaid
sequenceDiagram
    participant n_SourceWatcher as Source Watcher
    participant n_ReloadCoordinator as Reload Coordinator
    participant n_SduiFrontend as Sdui Frontend
    n_SourceWatcher->>n_ReloadCoordinator: 1: Reload Request (Model Reload Calls)
    n_ReloadCoordinator->>n_SduiFrontend: 2: Compile Ui Request (Ui Compilation Calls)
    n_SduiFrontend-->>n_ReloadCoordinator: 3: Compile Ui Rejected (Ui Compilation Calls) reply-to 2
    n_ReloadCoordinator-->>n_SourceWatcher: 4: Reload Rejected (Model Reload Calls) reply-to 1
```

Kildegrunnlag: f0109, f0110, f0111, f0112, f0113, f0516, f0519, f0520, f0640, f0644, f0647, f0649, f0660, f0661, f0662, f0663, f0664, f0820, f0822, f0926, f0928, f1001, f1003, f1004, f1034, f1035, f1036, f1037, f1038, f1039, f1040, f1041.

### Scenario: UnboundLocalAction — modus UiPreview

```mermaid
sequenceDiagram
    participant n_FyneHost as Fyne Host
    participant n_GoDomainImplementation as Go Domain Implementation
    n_FyneHost->>n_GoDomainImplementation: 1: Domain Action Request (Go Domain Calls)
    n_GoDomainImplementation-->>n_FyneHost: 2: Domain Action Result (Go Domain Calls) reply-to 1
```

Kildegrunnlag: f0181, f0182, f0183, f0184, f0185, f0239, f0240, f0440, f0442, f0443, f0447, f0449, f1064, f1065, f1066, f1067, f1068, f1069.


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

Kildegrunnlag: f0148, f0149, f0915.

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

Kildegrunnlag: f0832, f1017, f1018, f1049.

### Kontraktstruktur: ActionArguments

```mermaid
flowchart LR
    n_ActionArguments["ActionArguments (contract)"]
    n_ActionGeneration["ActionGeneration (field)"]
    n_ActionInputText["ActionInputText (field)"]
    n_ActionSymbol["ActionSymbol (field)"]
    n_ActionArguments -->|has-field| n_ActionGeneration
    n_ActionArguments -->|has-field| n_ActionInputText
    n_ActionArguments -->|has-field| n_ActionSymbol
```

Kildegrunnlag: f0001, f0002, f0003.

### Kontraktstruktur: ActionOutcome

```mermaid
flowchart LR
    n_ActionOutcome["ActionOutcome (contract)"]
    n_ActionOutputText["ActionOutputText (field)"]
    n_ActionStatusCode["ActionStatusCode (field)"]
    n_ActionOutcome -->|has-field| n_ActionOutputText
    n_ActionOutcome -->|has-field| n_ActionStatusCode
```

Kildegrunnlag: f0009, f0010.

### Kontraktstruktur: AstArtifactContract

```mermaid
flowchart LR
    n_AstArtifact["AstArtifact (field)"]
    n_AstArtifactContract["AstArtifactContract (contract)"]
    n_AstArtifactContract -->|has-field| n_AstArtifact
```

Kildegrunnlag: f0025.

### Kontraktstruktur: DesignSourceRecord

```mermaid
flowchart LR
    n_DesignSourceRecord["DesignSourceRecord (contract)"]
    n_SourceDocumentRevision["SourceDocumentRevision (field)"]
    n_SourceDocumentText["SourceDocumentText (field)"]
    n_DesignSourceRecord -->|has-field| n_SourceDocumentRevision
    n_DesignSourceRecord -->|has-field| n_SourceDocumentText
```

Kildegrunnlag: f0151, f0152.

### Kontraktstruktur: FramePresentationCallsProtocol

```mermaid
flowchart LR
    n_FramePresentationCallsProtocol["FramePresentationCallsProtocol (contract)"]
    n_PresentFrameRequest["PresentFrameRequest (message)"]
    n_PresentFrameResult["PresentFrameResult (message)"]
    n_FramePresentationCallsProtocol -->|permits| n_PresentFrameRequest
    n_FramePresentationCallsProtocol -->|permits| n_PresentFrameResult
```

Kildegrunnlag: f0211, f0212.

### Kontraktstruktur: GeneratedGoContract

```mermaid
flowchart LR
    n_GeneratedGoContract["GeneratedGoContract (contract)"]
    n_GeneratedGoSources["GeneratedGoSources (field)"]
    n_GeneratedGoContract -->|has-field| n_GeneratedGoSources
```

Kildegrunnlag: f0417.

### Kontraktstruktur: GoBuildCallsProtocol

```mermaid
flowchart LR
    n_BuildGoRequest["BuildGoRequest (message)"]
    n_BuildGoResult["BuildGoResult (message)"]
    n_GoBuildCallsProtocol["GoBuildCallsProtocol (contract)"]
    n_GoBuildCallsProtocol -->|permits| n_BuildGoRequest
    n_GoBuildCallsProtocol -->|permits| n_BuildGoResult
```

Kildegrunnlag: f0422, f0423.

### Kontraktstruktur: GoDomainCallsProtocol

```mermaid
flowchart LR
    n_DomainActionRequest["DomainActionRequest (message)"]
    n_DomainActionResult["DomainActionResult (message)"]
    n_GoDomainCallsProtocol["GoDomainCallsProtocol (contract)"]
    n_GoDomainCallsProtocol -->|permits| n_DomainActionRequest
    n_GoDomainCallsProtocol -->|permits| n_DomainActionResult
```

Kildegrunnlag: f0442, f0443.

### Kontraktstruktur: GoGenerationCallsProtocol

```mermaid
flowchart LR
    n_GenerateGoRequest["GenerateGoRequest (message)"]
    n_GenerateGoResult["GenerateGoResult (message)"]
    n_GoGenerationCallsProtocol["GoGenerationCallsProtocol (contract)"]
    n_GoGenerationCallsProtocol -->|permits| n_GenerateGoRequest
    n_GoGenerationCallsProtocol -->|permits| n_GenerateGoResult
```

Kildegrunnlag: f0452, f0453.

### Kontraktstruktur: LayoutArguments

```mermaid
flowchart LR
    n_LayoutArguments["LayoutArguments (contract)"]
    n_LayoutModelArtifact["LayoutModelArtifact (field)"]
    n_LayoutViewportHeight["LayoutViewportHeight (field)"]
    n_LayoutViewportWidth["LayoutViewportWidth (field)"]
    n_LayoutArguments -->|has-field| n_LayoutModelArtifact
    n_LayoutArguments -->|has-field| n_LayoutViewportHeight
    n_LayoutArguments -->|has-field| n_LayoutViewportWidth
```

Kildegrunnlag: f0477, f0478, f0479.

### Kontraktstruktur: LayoutCallsProtocol

```mermaid
flowchart LR
    n_LayoutCallsProtocol["LayoutCallsProtocol (contract)"]
    n_LayoutRequest["LayoutRequest (message)"]
    n_LayoutResult["LayoutResult (message)"]
    n_LayoutCallsProtocol -->|permits| n_LayoutRequest
    n_LayoutCallsProtocol -->|permits| n_LayoutResult
```

Kildegrunnlag: f0482, f0483.

### Kontraktstruktur: ModelReloadCallsProtocol

```mermaid
flowchart LR
    n_ModelReloadCallsProtocol["ModelReloadCallsProtocol (contract)"]
    n_ReloadPublished["ReloadPublished (message)"]
    n_ReloadRejected["ReloadRejected (message)"]
    n_ReloadRequest["ReloadRequest (message)"]
    n_ModelReloadCallsProtocol -->|permits| n_ReloadPublished
    n_ModelReloadCallsProtocol -->|permits| n_ReloadRejected
    n_ModelReloadCallsProtocol -->|permits| n_ReloadRequest
```

Kildegrunnlag: f0518, f0519, f0520.

### Kontraktstruktur: NativeBuildContract

```mermaid
flowchart LR
    n_NativeBuildContract["NativeBuildContract (contract)"]
    n_NativeBuildDiagnostics["NativeBuildDiagnostics (field)"]
    n_NativeBuildSucceeded["NativeBuildSucceeded (field)"]
    n_NativeBuildContract -->|has-field| n_NativeBuildDiagnostics
    n_NativeBuildContract -->|has-field| n_NativeBuildSucceeded
```

Kildegrunnlag: f0522, f0523.

### Kontraktstruktur: NativeUiActionsProtocol

```mermaid
flowchart LR
    n_NativeUiActionsProtocol["NativeUiActionsProtocol (contract)"]
    n_UiActionRejected["UiActionRejected (message)"]
    n_UiActionRequest["UiActionRequest (message)"]
    n_UiActionResult["UiActionResult (message)"]
    n_NativeUiActionsProtocol -->|permits| n_UiActionRejected
    n_NativeUiActionsProtocol -->|permits| n_UiActionRequest
    n_NativeUiActionsProtocol -->|permits| n_UiActionResult
```

Kildegrunnlag: f0542, f0543, f0544.

### Kontraktstruktur: NormalizedModelContract

```mermaid
flowchart LR
    n_NormalizedModelArtifact["NormalizedModelArtifact (field)"]
    n_NormalizedModelContract["NormalizedModelContract (contract)"]
    n_NormalizedModelContract -->|has-field| n_NormalizedModelArtifact
```

Kildegrunnlag: f0556.

### Kontraktstruktur: PreparedFrameContract

```mermaid
flowchart LR
    n_PreparedFrameArtifact["PreparedFrameArtifact (field)"]
    n_PreparedFrameContract["PreparedFrameContract (contract)"]
    n_PreparedFrameContract -->|has-field| n_PreparedFrameArtifact
```

Kildegrunnlag: f0573.

### Kontraktstruktur: PresentationOutcome

```mermaid
flowchart LR
    n_PresentationOutcome["PresentationOutcome (contract)"]
    n_PresentationReady["PresentationReady (field)"]
    n_PresentationOutcome -->|has-field| n_PresentationReady
```

Kildegrunnlag: f0580.

### Kontraktstruktur: ReloadArguments

```mermaid
flowchart LR
    n_ReloadArguments["ReloadArguments (contract)"]
    n_ReloadSourceRevision["ReloadSourceRevision (field)"]
    n_ReloadSourceText["ReloadSourceText (field)"]
    n_ReloadArguments -->|has-field| n_ReloadSourceRevision
    n_ReloadArguments -->|has-field| n_ReloadSourceText
```

Kildegrunnlag: f0624, f0625.

### Kontraktstruktur: ReloadOutcome

```mermaid
flowchart LR
    n_ReloadDiagnostic["ReloadDiagnostic (field)"]
    n_ReloadOutcome["ReloadOutcome (contract)"]
    n_ReloadPublishedGeneration["ReloadPublishedGeneration (field)"]
    n_ReloadOutcome -->|has-field| n_ReloadDiagnostic
    n_ReloadOutcome -->|has-field| n_ReloadPublishedGeneration
```

Kildegrunnlag: f0653, f0654.

### Kontraktstruktur: SdlActionCallsProtocol

```mermaid
flowchart LR
    n_SdlActionCallsProtocol["SdlActionCallsProtocol (contract)"]
    n_SdlActionRequest["SdlActionRequest (message)"]
    n_SdlActionResult["SdlActionResult (message)"]
    n_SdlActionCallsProtocol -->|permits| n_SdlActionRequest
    n_SdlActionCallsProtocol -->|permits| n_SdlActionResult
```

Kildegrunnlag: f0703, f0704.

### Kontraktstruktur: SdlCompilationCallsProtocol

```mermaid
flowchart LR
    n_CompileSdlRequest["CompileSdlRequest (message)"]
    n_CompileSdlResult["CompileSdlResult (message)"]
    n_SdlCompilationCallsProtocol["SdlCompilationCallsProtocol (contract)"]
    n_SdlCompilationCallsProtocol -->|permits| n_CompileSdlRequest
    n_SdlCompilationCallsProtocol -->|permits| n_CompileSdlResult
```

Kildegrunnlag: f0712, f0713.

### Kontraktstruktur: SvgDocumentContract

```mermaid
flowchart LR
    n_SvgDocumentContract["SvgDocumentContract (contract)"]
    n_SvgDocumentText["SvgDocumentText (field)"]
    n_SvgDocumentContract -->|has-field| n_SvgDocumentText
```

Kildegrunnlag: f0943.

### Kontraktstruktur: SvgExportCallsProtocol

```mermaid
flowchart LR
    n_ExportSvgRequest["ExportSvgRequest (message)"]
    n_ExportSvgResult["ExportSvgResult (message)"]
    n_SvgExportCallsProtocol["SvgExportCallsProtocol (contract)"]
    n_SvgExportCallsProtocol -->|permits| n_ExportSvgRequest
    n_SvgExportCallsProtocol -->|permits| n_ExportSvgResult
```

Kildegrunnlag: f0948, f0949.

### Kontraktstruktur: TokenArtifactContract

```mermaid
flowchart LR
    n_TokenArtifact["TokenArtifact (field)"]
    n_TokenArtifactContract["TokenArtifactContract (contract)"]
    n_TokenArtifactContract -->|has-field| n_TokenArtifact
```

Kildegrunnlag: f0953.

### Kontraktstruktur: UiAstCallsProtocol

```mermaid
flowchart LR
    n_BuildUiAstRequest["BuildUiAstRequest (message)"]
    n_BuildUiAstResult["BuildUiAstResult (message)"]
    n_UiAstCallsProtocol["UiAstCallsProtocol (contract)"]
    n_UiAstCallsProtocol -->|permits| n_BuildUiAstRequest
    n_UiAstCallsProtocol -->|permits| n_BuildUiAstResult
```

Kildegrunnlag: f0984, f0985.

### Kontraktstruktur: UiCompilationCallsProtocol

```mermaid
flowchart LR
    n_CompileUiRejected["CompileUiRejected (message)"]
    n_CompileUiRequest["CompileUiRequest (message)"]
    n_CompileUiResult["CompileUiResult (message)"]
    n_UiCompilationCallsProtocol["UiCompilationCallsProtocol (contract)"]
    n_UiCompilationCallsProtocol -->|permits| n_CompileUiRejected
    n_UiCompilationCallsProtocol -->|permits| n_CompileUiRequest
    n_UiCompilationCallsProtocol -->|permits| n_CompileUiResult
```

Kildegrunnlag: f1003, f1004, f1005.

### Kontraktstruktur: UiDomainActionsProtocol

```mermaid
flowchart LR
    n_BoundActionRequest["BoundActionRequest (message)"]
    n_BoundActionResult["BoundActionResult (message)"]
    n_UiDomainActionsProtocol["UiDomainActionsProtocol (contract)"]
    n_UiDomainActionsProtocol -->|permits| n_BoundActionRequest
    n_UiDomainActionsProtocol -->|permits| n_BoundActionResult
```

Kildegrunnlag: f1008, f1009.

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

Kildegrunnlag: f1010, f1011, f1013.

### Kontraktstruktur: UiGenerationEventsProtocol

```mermaid
flowchart LR
    n_UiGenerationEventsProtocol["UiGenerationEventsProtocol (contract)"]
    n_UiGenerationNotices["UiGenerationNotices (datagram)"]
    n_UiGenerationEventsProtocol -->|permits| n_UiGenerationNotices
```

Kildegrunnlag: f1016.

### Kontraktstruktur: UiNormalizationCallsProtocol

```mermaid
flowchart LR
    n_NormalizeUiRequest["NormalizeUiRequest (message)"]
    n_NormalizeUiResult["NormalizeUiResult (message)"]
    n_UiNormalizationCallsProtocol["UiNormalizationCallsProtocol (contract)"]
    n_UiNormalizationCallsProtocol -->|permits| n_NormalizeUiRequest
    n_UiNormalizationCallsProtocol -->|permits| n_NormalizeUiResult
```

Kildegrunnlag: f1044, f1045.

### Kontraktstruktur: UiSessionRecord

```mermaid
flowchart LR
    n_SessionDraft["SessionDraft (field)"]
    n_SessionGeneration["SessionGeneration (field)"]
    n_UiSessionRecord["UiSessionRecord (contract)"]
    n_UiSessionRecord -->|has-field| n_SessionDraft
    n_UiSessionRecord -->|has-field| n_SessionGeneration
```

Kildegrunnlag: f1047, f1048.

### Kontraktstruktur: UiTokenizationCallsProtocol

```mermaid
flowchart LR
    n_TokenizeUiRequest["TokenizeUiRequest (message)"]
    n_TokenizeUiResult["TokenizeUiResult (message)"]
    n_UiTokenizationCallsProtocol["UiTokenizationCallsProtocol (contract)"]
    n_UiTokenizationCallsProtocol -->|permits| n_TokenizeUiRequest
    n_UiTokenizationCallsProtocol -->|permits| n_TokenizeUiResult
```

Kildegrunnlag: f1058, f1059.

### Kontraktstruktur: UiValidationCallsProtocol

```mermaid
flowchart LR
    n_UiValidationCallsProtocol["UiValidationCallsProtocol (contract)"]
    n_ValidateUiRequest["ValidateUiRequest (message)"]
    n_ValidateUiResult["ValidateUiResult (message)"]
    n_UiValidationCallsProtocol -->|permits| n_ValidateUiRequest
    n_UiValidationCallsProtocol -->|permits| n_ValidateUiResult
```

Kildegrunnlag: f1062, f1063.

### Kontraktstruktur: ValidationOutcomeContract

```mermaid
flowchart LR
    n_ModelDiagnostics["ModelDiagnostics (field)"]
    n_ModelIsValid["ModelIsValid (field)"]
    n_ValidationOutcomeContract["ValidationOutcomeContract (contract)"]
    n_ValidationOutcomeContract -->|has-field| n_ModelDiagnostics
    n_ValidationOutcomeContract -->|has-field| n_ModelIsValid
```

Kildegrunnlag: f1099, f1100.


## VP10 — Datagram-koding og packet

### Packet: UiGenerationWire / UiGenerationChanged — big-endian, most-significant-first

```mermaid
packet
    0-15: "NoticeVersion"
    16-79: "NoticeGeneration"
```

Kildegrunnlag: f0557, f0558, f0559, f0560, f1010, f1011, f1012, f1013, f1019, f1020, f1021, f1022, f1023.


### VP09 — felt og kontraktegenskaper

| Modellfaktum | Kilde-ID |
| --- | --- |
| ActionArguments has completeness = closed. | f0000 |
| ActionGeneration has presence = required. | f0004 |
| ActionGeneration has value-type = unsigned. | f0005 |
| ActionInputText has presence = required. | f0006 |
| ActionInputText has value-type = text. | f0007 |
| ActionOutcome has completeness = closed. | f0008 |
| ActionOutputText has presence = required. | f0011 |
| ActionOutputText has value-type = text. | f0012 |
| ActionStatusCode has presence = required. | f0013 |
| ActionStatusCode has value-type = unsigned. | f0014 |
| ActionSymbol has presence = required. | f0015 |
| ActionSymbol has value-type = text. | f0016 |
| AstArtifact has presence = required. | f0022 |
| AstArtifact has value-type = bytes. | f0023 |
| AstArtifactContract has completeness = closed. | f0024 |
| DesignSourceRecord has completeness = closed. | f0150 |
| FramePresentationCallsProtocol has completeness = closed. | f0210 |
| GeneratedGoContract has completeness = closed. | f0416 |
| GeneratedGoSources has presence = required. | f0418 |
| GeneratedGoSources has value-type = bytes. | f0419 |
| GoBuildCallsProtocol has completeness = closed. | f0421 |
| GoDomainCallsProtocol has completeness = closed. | f0441 |
| GoGenerationCallsProtocol has completeness = closed. | f0451 |
| LayoutArguments has completeness = closed. | f0476 |
| LayoutCallsProtocol has completeness = closed. | f0481 |
| LayoutModelArtifact has presence = required. | f0484 |
| LayoutModelArtifact has value-type = bytes. | f0485 |
| LayoutViewportHeight has presence = required. | f0491 |
| LayoutViewportHeight has value-type = decimal. | f0492 |
| LayoutViewportWidth has presence = required. | f0493 |
| LayoutViewportWidth has value-type = decimal. | f0494 |
| ModelDiagnostics has presence = optional. | f0512 |
| ModelDiagnostics has value-type = text. | f0513 |
| ModelIsValid has presence = required. | f0514 |
| ModelIsValid has value-type = boolean. | f0515 |
| ModelReloadCallsProtocol has completeness = closed. | f0517 |
| NativeBuildContract has completeness = closed. | f0521 |
| NativeBuildDiagnostics has presence = optional. | f0524 |
| NativeBuildDiagnostics has value-type = text. | f0525 |
| NativeBuildSucceeded has presence = required. | f0526 |
| NativeBuildSucceeded has value-type = boolean. | f0527 |
| NativeUiActionsProtocol has completeness = closed. | f0541 |
| NormalizedModelArtifact has presence = required. | f0553 |
| NormalizedModelArtifact has value-type = bytes. | f0554 |
| NormalizedModelContract has completeness = closed. | f0555 |
| NoticeGeneration has presence = required. | f0557 |
| NoticeGeneration has value-type = unsigned. | f0558 |
| NoticeVersion has presence = required. | f0559 |
| NoticeVersion has value-type = unsigned. | f0560 |
| PreparedFrameArtifact has presence = required. | f0570 |
| PreparedFrameArtifact has value-type = bytes. | f0571 |
| PreparedFrameContract has completeness = closed. | f0572 |
| PresentationOutcome has completeness = closed. | f0579 |
| PresentationReady has presence = required. | f0581 |
| PresentationReady has value-type = boolean. | f0582 |
| ReloadArguments has completeness = closed. | f0623 |
| ReloadDiagnostic has presence = optional. | f0650 |
| ReloadDiagnostic has value-type = text. | f0651 |
| ReloadOutcome has completeness = closed. | f0652 |
| ReloadPublishedGeneration has presence = optional. | f0658 |
| ReloadPublishedGeneration has value-type = unsigned. | f0659 |
| ReloadSourceRevision has presence = required. | f0665 |
| ReloadSourceRevision has value-type = unsigned. | f0666 |
| ReloadSourceText has presence = required. | f0667 |
| ReloadSourceText has value-type = text. | f0668 |
| SdlActionCallsProtocol has completeness = closed. | f0702 |
| SdlCompilationCallsProtocol has completeness = closed. | f0711 |
| SessionDraft has presence = optional. | f0904 |
| SessionDraft has value-type = text. | f0905 |
| SessionGeneration has presence = required. | f0906 |
| SessionGeneration has value-type = unsigned. | f0907 |
| SourceDocumentRevision has presence = required. | f0910 |
| SourceDocumentRevision has value-type = unsigned. | f0911 |
| SourceDocumentText has presence = required. | f0912 |
| SourceDocumentText has value-type = text. | f0913 |
| SvgDocumentContract has completeness = closed. | f0942 |
| SvgDocumentText has presence = required. | f0944 |
| SvgDocumentText has value-type = text. | f0945 |
| SvgExportCallsProtocol has completeness = closed. | f0947 |
| TokenArtifact has presence = required. | f0950 |
| TokenArtifact has value-type = bytes. | f0951 |
| TokenArtifactContract has completeness = closed. | f0952 |
| UiAstCallsProtocol has completeness = closed. | f0983 |
| UiCompilationCallsProtocol has completeness = closed. | f1002 |
| UiDomainActionsProtocol has completeness = closed. | f1007 |
| UiGenerationContract has completeness = closed. | f1012 |
| UiGenerationEventsProtocol has completeness = closed. | f1015 |
| UiGenerationWire has bit-order = most-significant-first. | f1020 |
| UiGenerationWire has byte-order = big-endian. | f1021 |
| UiNormalizationCallsProtocol has completeness = closed. | f1043 |
| UiSessionRecord has completeness = closed. | f1046 |
| UiTokenizationCallsProtocol has completeness = closed. | f1057 |
| UiValidationCallsProtocol has completeness = closed. | f1061 |
| ValidationOutcomeContract has completeness = closed. | f1098 |

### VP09 — projeksjonsansvar

| Functionality | Dataset | Datagram-familie | Faktum |
| --- | --- | --- | --- |
| ProjectUiGeneration | UiSessionState | UiGenerationNotices | f0604 |


### VP08 — avledet MessageSet per Channel og modus

Generert fra permits og deltakelse, ikke en separat authored modell. Tom deltakelse er et hull.

| Channel | Mode | Message / Datagram | Sender | Receiver | Kilde-ID-er |
| --- | --- | --- | --- | --- | --- |
| FramePresentationCalls | UiPreview | PresentFrameRequest | FyneHost | FyneBackend | f0209, f0211, f0220, f0238, f0575 |
| FramePresentationCalls | UiPreview | PresentFrameResult | FyneBackend | FyneHost | f0209, f0212, f0221, f0237, f0578 |
| GoBuildCalls | NativeBuild | BuildGoRequest | CommandLineHost | GoBuildRunner | f0056, f0094, f0420, f0422, f0429 |
| GoBuildCalls | NativeBuild | BuildGoResult | GoBuildRunner | CommandLineHost | f0059, f0093, f0420, f0423, f0430 |
| GoDomainCalls | BoundExecution | DomainActionRequest | SdlDispatcher | GoDomainImplementation | f0182, f0440, f0442, f0446, f0722 |
| GoDomainCalls | BoundExecution | DomainActionResult | GoDomainImplementation | SdlDispatcher | f0185, f0440, f0443, f0448, f0721 |
| GoDomainCalls | UiPreview | DomainActionRequest | FyneHost | GoDomainImplementation | f0182, f0240, f0440, f0442, f0447 |
| GoDomainCalls | UiPreview | DomainActionResult | GoDomainImplementation | FyneHost | f0185, f0239, f0440, f0443, f0449 |
| GoGenerationCalls | NativeBuild | GenerateGoRequest | CommandLineHost | GoCodeGenerator | f0096, f0409, f0438, f0450, f0452 |
| GoGenerationCalls | NativeBuild | GenerateGoResult | GoCodeGenerator | CommandLineHost | f0095, f0412, f0439, f0450, f0453 |
| LayoutCalls | StaticExport | LayoutRequest | CommandLineHost | SduiLayout | f0098, f0480, f0482, f0487, f0849 |
| LayoutCalls | StaticExport | LayoutResult | SduiLayout | CommandLineHost | f0097, f0480, f0483, f0490, f0851 |
| LayoutCalls | UiPreview | LayoutRequest | FyneHost | SduiLayout | f0242, f0480, f0482, f0487, f0850 |
| LayoutCalls | UiPreview | LayoutResult | SduiLayout | FyneHost | f0241, f0480, f0483, f0490, f0852 |
| ModelReloadCalls | BoundLiveEditing | ReloadPublished | ReloadCoordinator | SourceWatcher | f0516, f0518, f0641, f0657, f0923 |
| ModelReloadCalls | BoundLiveEditing | ReloadRejected | ReloadCoordinator | SourceWatcher | f0516, f0519, f0643, f0662, f0925 |
| ModelReloadCalls | BoundLiveEditing | ReloadRequest | SourceWatcher | ReloadCoordinator | f0516, f0520, f0639, f0664, f0927 |
| ModelReloadCalls | LiveEditing | ReloadPublished | ReloadCoordinator | SourceWatcher | f0516, f0518, f0642, f0657, f0924 |
| ModelReloadCalls | LiveEditing | ReloadRejected | ReloadCoordinator | SourceWatcher | f0516, f0519, f0644, f0662, f0926 |
| ModelReloadCalls | LiveEditing | ReloadRequest | SourceWatcher | ReloadCoordinator | f0516, f0520, f0640, f0664, f0928 |
| NativeUiActions | BoundExecution | UiActionRejected | SduiDispatcher | FyneBackend | f0222, f0540, f0542, f0806, f0976 |
| NativeUiActions | BoundExecution | UiActionRequest | FyneBackend | SduiDispatcher | f0224, f0540, f0543, f0805, f0978 |
| NativeUiActions | BoundExecution | UiActionResult | SduiDispatcher | FyneBackend | f0223, f0540, f0544, f0807, f0981 |
| SdlActionCalls | BoundExecution | SdlActionRequest | SdlUiBindingAdapter | SdlDispatcher | f0701, f0703, f0706, f0723, f0788 |
| SdlActionCalls | BoundExecution | SdlActionResult | SdlDispatcher | SdlUiBindingAdapter | f0701, f0704, f0709, f0724, f0787 |
| SdlCompilationCalls | BoundLiveEditing | CompileSdlRequest | ReloadCoordinator | SdlFrontend | f0105, f0646, f0710, f0712, f0737 |
| SdlCompilationCalls | BoundLiveEditing | CompileSdlResult | SdlFrontend | ReloadCoordinator | f0108, f0645, f0710, f0713, f0738 |
| SvgExportCalls | StaticExport | ExportSvgRequest | CommandLineHost | SduiPresentation | f0100, f0197, f0880, f0946, f0948 |
| SvgExportCalls | StaticExport | ExportSvgResult | SduiPresentation | CommandLineHost | f0099, f0200, f0881, f0946, f0949 |
| UiAstCalls | SourceInspection | BuildUiAstRequest | SduiFrontend | SduiParser | f0075, f0819, f0873, f0982, f0984 |
| UiAstCalls | SourceInspection | BuildUiAstResult | SduiParser | SduiFrontend | f0078, f0818, f0874, f0982, f0985 |
| UiCompilationCalls | LiveEditing | CompileUiRejected | SduiFrontend | ReloadCoordinator | f0111, f0647, f0822, f1001, f1003 |
| UiCompilationCalls | LiveEditing | CompileUiRequest | ReloadCoordinator | SduiFrontend | f0113, f0649, f0820, f1001, f1004 |
| UiCompilationCalls | LiveEditing | CompileUiResult | SduiFrontend | ReloadCoordinator | f0116, f0648, f0824, f1001, f1005 |
| UiCompilationCalls | SourceInspection | CompileUiRejected | SduiFrontend | CommandLineHost | f0101, f0111, f0823, f1001, f1003 |
| UiCompilationCalls | SourceInspection | CompileUiRequest | CommandLineHost | SduiFrontend | f0103, f0113, f0821, f1001, f1004 |
| UiCompilationCalls | SourceInspection | CompileUiResult | SduiFrontend | CommandLineHost | f0102, f0116, f0825, f1001, f1005 |
| UiDomainActions | BoundExecution | BoundActionRequest | SduiDispatcher | SdlUiBindingAdapter | f0046, f0789, f0809, f1006, f1008 |
| UiDomainActions | BoundExecution | BoundActionResult | SdlUiBindingAdapter | SduiDispatcher | f0049, f0790, f0808, f1006, f1009 |
| UiGenerationEvents | BoundExecution | UiGenerationNotices | SduiInstanceStore | FyneBackend | f0225, f0838, f1014, f1016, f1018 |
| UiGenerationEvents | LiveEditing | UiGenerationNotices | SduiInstanceStore | FyneBackend | f0226, f0839, f1014, f1016, f1018 |
| UiNormalizationCalls | SourceInspection | NormalizeUiRequest | SduiFrontend | SduiNormalizer | f0549, f0827, f0869, f1042, f1044 |
| UiNormalizationCalls | SourceInspection | NormalizeUiResult | SduiNormalizer | SduiFrontend | f0552, f0826, f0870, f1042, f1045 |
| UiTokenizationCalls | SourceInspection | TokenizeUiRequest | SduiFrontend | SduiLexer | f0829, f0855, f0961, f1056, f1058 |
| UiTokenizationCalls | SourceInspection | TokenizeUiResult | SduiLexer | SduiFrontend | f0828, f0856, f0964, f1056, f1059 |
| UiValidationCalls | SourceInspection | ValidateUiRequest | SduiFrontend | SduiValidator | f0831, f0902, f1060, f1062, f1088 |
| UiValidationCalls | SourceInspection | ValidateUiResult | SduiValidator | SduiFrontend | f0830, f0903, f1060, f1063, f1091 |

## VP04 — Grensesnittbruk

Tabellen dekker alle consumes-fakta. Portnavn er ikke Channel-kontrakter eller tilbyderkoblinger.

| Unit / Container | Interface | Faktum | Kildelinje |
| --- | --- | --- | --- |
| CommandLineHost | ExportSinkPort | f0086 | 456 |
| CommandLineHost | PreparedFramePort | f0087 | 457 |
| CommandLineHost | SduiFrontendPort | f0088 | 458 |
| CommandLineHost | SourceSnapshotPort | f0089 | 459 |
| DiagramProvider | DiagramEnginePort | f0173 | 543 |
| DiagramProvider | ResourcePort | f0174 | 544 |
| DomainStateMigrator | DomainStatePort | f0186 | 556 |
| DomainStateMigrator | SdlModelPort | f0187 | 557 |
| FyneBackend | PreparedFramePort | f0213 | 583 |
| FyneBackend | UiSessionPort | f0214 | 584 |
| FyneHost | DomainBindingPort | f0227 | 597 |
| FyneHost | ReloadPort | f0228 | 598 |
| FyneHost | SdlFrontendPort | f0229 | 599 |
| FyneHost | SduiFrontendPort | f0230 | 600 |
| FyneHost | SourceSnapshotPort | f0231 | 601 |
| FyneHost | UiSessionPort | f0232 | 602 |
| FyneHost | WidgetBackendPort | f0233 | 603 |
| GoBuildRunner | BuildToolPort | f0424 | 794 |
| GoBuildRunner | GeneratedArtifactPort | f0425 | 795 |
| GoCodeGenerator | ExecutionProfilePort | f0431 | 801 |
| GoCodeGenerator | SdlModelPort | f0432 | 802 |
| GoCodeGenerator | SduiModelPort | f0433 | 803 |
| MarkdownProvider | DiagramPort | f0500 | 870 |
| MarkdownProvider | MeasurementPort | f0501 | 871 |
| MarkdownProvider | ResourcePort | f0502 | 872 |
| ReloadCoordinator | BindingReloadPort | f0627 | 997 |
| ReloadCoordinator | DiagnosticPort | f0628 | 998 |
| ReloadCoordinator | SdlFrontendPort | f0629 | 999 |
| ReloadCoordinator | SdlReloadPort | f0630 | 1000 |
| ReloadCoordinator | SduiFrontendPort | f0631 | 1001 |
| ReloadCoordinator | SourceSnapshotPort | f0632 | 1002 |
| ReloadCoordinator | UiReloadPort | f0633 | 1003 |
| SdlDispatcher | DomainFunctionPort | f0714 | 1084 |
| SdlDispatcher | DomainStatePort | f0715 | 1085 |
| SdlExecutionGate | DiagnosticPort | f0725 | 1095 |
| SdlExecutionGate | SdlModelPort | f0726 | 1096 |
| SdlFrontend | DiagnosticPort | f0729 | 1099 |
| SdlFrontend | SourceSnapshotPort | f0730 | 1100 |
| SdlFunctionRegistry | DomainFunctionPort | f0739 | 1109 |
| SdlRuntime | DomainFunctionPort | f0764 | 1134 |
| SdlRuntime | SdlModelPort | f0765 | 1135 |
| SdlUiBindingAdapter | DiagnosticPort | f0778 | 1148 |
| SdlUiBindingAdapter | SdlExecutionPort | f0779 | 1149 |
| SdlUiBindingAdapter | UiSessionPort | f0780 | 1150 |
| SduiDispatcher | DomainBindingPort | f0798 | 1168 |
| SduiDispatcher | UiStatePort | f0799 | 1169 |
| SduiFrontend | DiagnosticPort | f0810 | 1180 |
| SduiFrontend | SourceSnapshotPort | f0811 | 1181 |
| SduiLayout | ContentProviderPort | f0840 | 1210 |
| SduiLayout | MeasurementPort | f0841 | 1211 |
| SduiLayout | UiSnapshotPort | f0842 | 1212 |
| SduiPresentation | PreparedFramePort | f0875 | 1245 |
| SduiPropertyStore | UiStatePort | f0882 | 1252 |
| SduiRuntime | DomainBindingPort | f0887 | 1257 |
| SduiRuntime | SduiModelPort | f0888 | 1258 |
| SourceLoader | SourceInputPort | f0914 | 1284 |
| SourceWatcher | FileChangePort | f0919 | 1289 |
| UiStateReconciler | SduiModelPort | f1050 | 1420 |
| UiStateReconciler | UiStatePort | f1051 | 1421 |

## VP11 — Egenskaper og fullstendig faktaregister

Registeret inkluderer alle fakta, også de som ikke har en egen tegning.

| ID | Utsagn | Kildelinje |
| --- | --- | --- |
| f0000 | ActionArguments has completeness = closed. | 370 |
| f0001 | ActionArguments has-field ActionGeneration. | 371 |
| f0002 | ActionArguments has-field ActionInputText. | 372 |
| f0003 | ActionArguments has-field ActionSymbol. | 373 |
| f0004 | ActionGeneration has presence = required. | 374 |
| f0005 | ActionGeneration has value-type = unsigned. | 375 |
| f0006 | ActionInputText has presence = required. | 376 |
| f0007 | ActionInputText has value-type = text. | 377 |
| f0008 | ActionOutcome has completeness = closed. | 378 |
| f0009 | ActionOutcome has-field ActionOutputText. | 379 |
| f0010 | ActionOutcome has-field ActionStatusCode. | 380 |
| f0011 | ActionOutputText has presence = required. | 381 |
| f0012 | ActionOutputText has value-type = text. | 382 |
| f0013 | ActionStatusCode has presence = required. | 383 |
| f0014 | ActionStatusCode has value-type = unsigned. | 384 |
| f0015 | ActionSymbol has presence = required. | 385 |
| f0016 | ActionSymbol has value-type = text. | 386 |
| f0017 | AllocateGeometry allocated-to FyneHost in mode UiPreview. | 387 |
| f0018 | AllocateGeometry contributes-to InteractiveUiPreview. | 388 |
| f0019 | AllocateGeometry realizes MeasuredPresentation. | 389 |
| f0020 | ApplyPropertyBatch has state-retention = stateful. | 390 |
| f0021 | ApplyPropertyBatch realizes InteractiveSession. | 391 |
| f0022 | AstArtifact has presence = required. | 392 |
| f0023 | AstArtifact has value-type = bytes. | 393 |
| f0024 | AstArtifactContract has completeness = closed. | 394 |
| f0025 | AstArtifactContract has-field AstArtifact. | 395 |
| f0026 | BoundActionAccepted exercises TryDomainInteraction. | 396 |
| f0027 | BoundActionAccepted has completeness = closed. | 397 |
| f0028 | BoundActionAccepted illustrates G4M3UiDomainBinding. | 398 |
| f0029 | BoundActionAccepted runs-in BoundExecution. | 399 |
| f0030 | BoundActionAccepted step 1 sends UiActionRequest from FyneBackend to SduiDispatcher via NativeUiActions. | 400 |
| f0031 | BoundActionAccepted step 2 sends BoundActionRequest from SduiDispatcher to SdlUiBindingAdapter via UiDomainActions. | 401 |
| f0032 | BoundActionAccepted step 3 sends SdlActionRequest from SdlUiBindingAdapter to SdlDispatcher via SdlActionCalls. | 402 |
| f0033 | BoundActionAccepted step 4 sends DomainActionRequest from SdlDispatcher to GoDomainImplementation via GoDomainCalls. | 403 |
| f0034 | BoundActionAccepted step 5 sends DomainActionResult from GoDomainImplementation to SdlDispatcher via GoDomainCalls reply-to 4. | 404 |
| f0035 | BoundActionAccepted step 6 sends SdlActionResult from SdlDispatcher to SdlUiBindingAdapter via SdlActionCalls reply-to 3. | 405 |
| f0036 | BoundActionAccepted step 7 sends BoundActionResult from SdlUiBindingAdapter to SduiDispatcher via UiDomainActions reply-to 2. | 406 |
| f0037 | BoundActionAccepted step 8 sends UiActionResult from SduiDispatcher to FyneBackend via NativeUiActions reply-to 1. | 407 |
| f0038 | BoundActionAccepted step 9 sends UiGenerationNotices variant UiGenerationChanged from SduiInstanceStore to FyneBackend via UiGenerationEvents. | 408 |
| f0039 | BoundActionRejected exercises TryDomainInteraction. | 409 |
| f0040 | BoundActionRejected has completeness = closed. | 410 |
| f0041 | BoundActionRejected illustrates G4M3UiDomainBinding. | 411 |
| f0042 | BoundActionRejected runs-in BoundExecution. | 412 |
| f0043 | BoundActionRejected step 1 sends UiActionRequest from FyneBackend to SduiDispatcher via NativeUiActions. | 413 |
| f0044 | BoundActionRejected step 2 sends UiActionRejected from SduiDispatcher to FyneBackend via NativeUiActions reply-to 1. | 414 |
| f0045 | BoundActionRequest has message-kind = request. | 415 |
| f0046 | BoundActionRequest upholds ActionArguments. | 416 |
| f0047 | BoundActionResult has message-kind = result. | 417 |
| f0048 | BoundActionResult replies-to BoundActionRequest. | 418 |
| f0049 | BoundActionResult upholds ActionOutcome. | 419 |
| f0050 | BoundInteraction requires SdlExecutionPort in mode BoundExecution. | 420 |
| f0051 | BoundInteraction requires UiSessionPort in mode BoundExecution. | 421 |
| f0052 | BuildGeneratedApplication allocated-to CommandLineHost in mode NativeBuild. | 422 |
| f0053 | BuildGeneratedApplication contributes-to NativeGoAssembly. | 423 |
| f0054 | BuildGeneratedApplication realizes NativeRealization. | 424 |
| f0055 | BuildGoRequest has message-kind = request. | 425 |
| f0056 | BuildGoRequest upholds GeneratedGoContract. | 426 |
| f0057 | BuildGoResult has message-kind = result. | 427 |
| f0058 | BuildGoResult replies-to BuildGoRequest. | 428 |
| f0059 | BuildGoResult upholds NativeBuildContract. | 429 |
| f0060 | BuildNativeRealization refines RealizeDesign. | 430 |
| f0061 | BuildPreparedFrame allocated-to FyneHost in mode UiPreview. | 431 |
| f0062 | BuildPreparedFrame contributes-to InteractiveUiPreview. | 432 |
| f0063 | BuildPreparedFrame realizes MeasuredPresentation. | 433 |
| f0064 | BuildSdlAst allocated-to CommandLineHost in mode SourceInspection. | 434 |
| f0065 | BuildSdlAst contributes-to StructuralModelInspection. | 435 |
| f0066 | BuildSdlAst has repeatability = deterministic. | 436 |
| f0067 | BuildSdlAst has state-retention = stateless. | 437 |
| f0068 | BuildSdlAst realizes SdlSourceModel. | 438 |
| f0069 | BuildSduiAst allocated-to CommandLineHost in mode SourceInspection. | 439 |
| f0070 | BuildSduiAst contributes-to StructuralModelInspection. | 440 |
| f0071 | BuildSduiAst has repeatability = deterministic. | 441 |
| f0072 | BuildSduiAst has state-retention = stateless. | 442 |
| f0073 | BuildSduiAst realizes SduiSourceModel. | 443 |
| f0074 | BuildUiAstRequest has message-kind = request. | 444 |
| f0075 | BuildUiAstRequest upholds TokenArtifactContract. | 445 |
| f0076 | BuildUiAstResult has message-kind = result. | 446 |
| f0077 | BuildUiAstResult replies-to BuildUiAstRequest. | 447 |
| f0078 | BuildUiAstResult upholds AstArtifactContract. | 448 |
| f0079 | CancelPendingActions realizes ExecutableDesign. | 449 |
| f0080 | CheckDomainStateCompatibility realizes ExecutableDesign. | 450 |
| f0081 | CheckExecutionCompleteness realizes ExecutableDesign. | 451 |
| f0082 | CheckFunctionSignatures realizes ExecutableDesign. | 452 |
| f0083 | CloseSdlInstance realizes ExecutableDesign. | 453 |
| f0084 | CloseUiInstance realizes InteractiveSession. | 454 |
| f0085 | CoalesceSourceChanges realizes DevelopmentReload. | 455 |
| f0086 | CommandLineHost consumes ExportSinkPort. | 456 |
| f0087 | CommandLineHost consumes PreparedFramePort. | 457 |
| f0088 | CommandLineHost consumes SduiFrontendPort. | 458 |
| f0089 | CommandLineHost consumes SourceSnapshotPort. | 459 |
| f0090 | CommandLineHost owns ComposeHeadlessExport. | 460 |
| f0091 | CommandLineHost owns WriteGeneratedArtifacts. | 461 |
| f0092 | CommandLineHost provides StaticDocumentation. | 462 |
| f0093 | CommandLineHost uses GoBuildCalls as receiver of BuildGoResult in mode NativeBuild. | 463 |
| f0094 | CommandLineHost uses GoBuildCalls as sender of BuildGoRequest in mode NativeBuild. | 464 |
| f0095 | CommandLineHost uses GoGenerationCalls as receiver of GenerateGoResult in mode NativeBuild. | 465 |
| f0096 | CommandLineHost uses GoGenerationCalls as sender of GenerateGoRequest in mode NativeBuild. | 466 |
| f0097 | CommandLineHost uses LayoutCalls as receiver of LayoutResult in mode StaticExport. | 467 |
| f0098 | CommandLineHost uses LayoutCalls as sender of LayoutRequest in mode StaticExport. | 468 |
| f0099 | CommandLineHost uses SvgExportCalls as receiver of ExportSvgResult in mode StaticExport. | 469 |
| f0100 | CommandLineHost uses SvgExportCalls as sender of ExportSvgRequest in mode StaticExport. | 470 |
| f0101 | CommandLineHost uses UiCompilationCalls as receiver of CompileUiRejected in mode SourceInspection. | 471 |
| f0102 | CommandLineHost uses UiCompilationCalls as receiver of CompileUiResult in mode SourceInspection. | 472 |
| f0103 | CommandLineHost uses UiCompilationCalls as sender of CompileUiRequest in mode SourceInspection. | 473 |
| f0104 | CompileSdlRequest has message-kind = request. | 474 |
| f0105 | CompileSdlRequest upholds ReloadArguments. | 475 |
| f0106 | CompileSdlResult has message-kind = result. | 476 |
| f0107 | CompileSdlResult replies-to CompileSdlRequest. | 477 |
| f0108 | CompileSdlResult upholds ReloadOutcome. | 478 |
| f0109 | CompileUiRejected has message-kind = result. | 479 |
| f0110 | CompileUiRejected replies-to CompileUiRequest. | 480 |
| f0111 | CompileUiRejected upholds ReloadOutcome. | 481 |
| f0112 | CompileUiRequest has message-kind = request. | 482 |
| f0113 | CompileUiRequest upholds ReloadArguments. | 483 |
| f0114 | CompileUiResult has message-kind = result. | 484 |
| f0115 | CompileUiResult replies-to CompileUiRequest. | 485 |
| f0116 | CompileUiResult upholds ReloadOutcome. | 486 |
| f0117 | ComposeHeadlessExport realizes StaticDocumentation. | 487 |
| f0118 | ComposeInteractiveSession allocated-to FyneHost in mode UiPreview. | 488 |
| f0119 | ComposeInteractiveSession contributes-to InteractiveUiPreview. | 489 |
| f0120 | ComposeInteractiveSession realizes NativeInteraction. | 490 |
| f0121 | ComposeMarkdownDocument allocated-to CommandLineHost in mode StaticExport. | 491 |
| f0122 | ComposeMarkdownDocument contributes-to DesignDocumentation. | 492 |
| f0123 | ComposeMarkdownDocument realizes StaticDocumentation. | 493 |
| f0124 | ComputeClipping realizes MeasuredPresentation. | 494 |
| f0125 | ConnectTypedWidgetHandles allocated-to FyneHost in mode BoundExecution. | 495 |
| f0126 | ConnectTypedWidgetHandles contributes-to TypedDomainBinding. | 496 |
| f0127 | ConnectTypedWidgetHandles realizes BoundInteraction. | 497 |
| f0128 | ContentServices contains DiagramProvider. | 498 |
| f0129 | ContentServices contains MarkdownProvider. | 499 |
| f0130 | ContentServices contains ResourceStore. | 500 |
| f0131 | ContentServices provides RichContent. | 501 |
| f0132 | CoordinateSdlCompilation realizes SdlSourceModel. | 502 |
| f0133 | CoordinateSduiCompilation realizes SduiSourceModel. | 503 |
| f0134 | CorrelateActionResult realizes ExecutableDesign. | 504 |
| f0135 | CorrelateUiResult realizes InteractiveSession. | 505 |
| f0136 | CreateSdlInstance realizes ExecutableDesign. | 506 |
| f0137 | CreateUiInstance realizes InteractiveSession. | 507 |
| f0138 | DesignAuthor pursues BuildNativeProduct. | 508 |
| f0139 | DesignAuthor pursues EditRunningPrototype. | 509 |
| f0140 | DesignAuthor pursues InspectModels. | 510 |
| f0141 | DesignAuthor pursues PrototypeUserInterface. | 511 |
| f0142 | DesignAuthor pursues PublishDesignDocumentation. | 512 |
| f0143 | DesignAuthor pursues TryDomainInteraction. | 513 |
| f0144 | DesignDocumentation supports PublishDesignDocumentation. | 514 |
| f0145 | DesignReviewer pursues InspectModels. | 515 |
| f0146 | DesignReviewer pursues PrototypeUserInterface. | 516 |
| f0147 | DesignReviewer pursues PublishDesignDocumentation. | 517 |
| f0148 | DesignSourceArchive holds DesignSourceDocuments. | 518 |
| f0149 | DesignSourceDocuments upholds DesignSourceRecord. | 519 |
| f0150 | DesignSourceRecord has completeness = closed. | 520 |
| f0151 | DesignSourceRecord has-field SourceDocumentRevision. | 521 |
| f0152 | DesignSourceRecord has-field SourceDocumentText. | 522 |
| f0153 | DevelopmentReload requires BindingReloadPort in mode BoundLiveEditing. | 523 |
| f0154 | DevelopmentReload requires FileChangePort in mode BoundLiveEditing. | 524 |
| f0155 | DevelopmentReload requires FileChangePort in mode LiveEditing. | 525 |
| f0156 | DevelopmentReload requires SdlReloadPort in mode BoundLiveEditing. | 526 |
| f0157 | DevelopmentReload requires SourceSnapshotPort in mode BoundLiveEditing. | 527 |
| f0158 | DevelopmentReload requires SourceSnapshotPort in mode LiveEditing. | 528 |
| f0159 | DevelopmentReload requires UiReloadPort in mode BoundLiveEditing. | 529 |
| f0160 | DevelopmentReload requires UiReloadPort in mode LiveEditing. | 530 |
| f0161 | DevelopmentTools contains DiagnosticReporter. | 531 |
| f0162 | DevelopmentTools contains GoBuildRunner. | 532 |
| f0163 | DevelopmentTools contains GoCodeGenerator. | 533 |
| f0164 | DevelopmentTools contains ReloadCoordinator. | 534 |
| f0165 | DevelopmentTools contains SourceLoader. | 535 |
| f0166 | DevelopmentTools contains SourceWatcher. | 536 |
| f0167 | DevelopmentTools owns RetireReplacedPythonEntryPoints. | 537 |
| f0168 | DevelopmentTools owns VerifyNativeBehaviorParity. | 538 |
| f0169 | DevelopmentTools provides DevelopmentReload. | 539 |
| f0170 | DiagnosticReporter owns ReportBindingDiagnostics. | 540 |
| f0171 | DiagnosticReporter owns ReportSourceDiagnostics. | 541 |
| f0172 | DiagnosticReporter provides SourceDiagnostics. | 542 |
| f0173 | DiagramProvider consumes DiagramEnginePort. | 543 |
| f0174 | DiagramProvider consumes ResourcePort. | 544 |
| f0175 | DiagramProvider owns PrepareDiagramResource. | 545 |
| f0176 | DiagramProvider provides RichContent. | 546 |
| f0177 | DisconnectBindings realizes BoundInteraction. | 547 |
| f0178 | DispatchUiEvent allocated-to FyneHost in mode BoundExecution. | 548 |
| f0179 | DispatchUiEvent contributes-to TypedDomainBinding. | 549 |
| f0180 | DispatchUiEvent realizes InteractiveSession. | 550 |
| f0181 | DomainActionRequest has message-kind = request. | 551 |
| f0182 | DomainActionRequest upholds ActionArguments. | 552 |
| f0183 | DomainActionResult has message-kind = result. | 553 |
| f0184 | DomainActionResult replies-to DomainActionRequest. | 554 |
| f0185 | DomainActionResult upholds ActionOutcome. | 555 |
| f0186 | DomainStateMigrator consumes DomainStatePort. | 556 |
| f0187 | DomainStateMigrator consumes SdlModelPort. | 557 |
| f0188 | DomainStateMigrator owns CheckDomainStateCompatibility. | 558 |
| f0189 | DomainStateMigrator owns MigrateOrResetDomainState. | 559 |
| f0190 | DomainStateMigrator provides ExecutableDesign. | 560 |
| f0191 | ExecutableDesign requires DomainFunctionPort in mode BoundExecution. | 561 |
| f0192 | ExpandUiDefinitions has repeatability = deterministic. | 562 |
| f0193 | ExpandUiDefinitions has state-retention = stateless. | 563 |
| f0194 | ExpandUiDefinitions realizes SduiSourceModel. | 564 |
| f0195 | ExportConsoleSnapshot realizes StaticDocumentation. | 565 |
| f0196 | ExportSvgRequest has message-kind = request. | 566 |
| f0197 | ExportSvgRequest upholds PreparedFrameContract. | 567 |
| f0198 | ExportSvgResult has message-kind = result. | 568 |
| f0199 | ExportSvgResult replies-to ExportSvgRequest. | 569 |
| f0200 | ExportSvgResult upholds SvgDocumentContract. | 570 |
| f0201 | ExportSvgSnapshot allocated-to CommandLineHost in mode StaticExport. | 571 |
| f0202 | ExportSvgSnapshot contributes-to DesignDocumentation. | 572 |
| f0203 | ExportSvgSnapshot realizes StaticDocumentation. | 573 |
| f0204 | ExportUiDocumentation refines InspectDesignSource. | 574 |
| f0205 | ExportViewpointMarkdown allocated-to CommandLineHost in mode StaticExport. | 575 |
| f0206 | ExportViewpointMarkdown contributes-to DesignDocumentation. | 576 |
| f0207 | ExportViewpointMarkdown has repeatability = deterministic. | 577 |
| f0208 | ExportViewpointMarkdown has state-retention = stateless. | 578 |
| f0209 | FramePresentationCalls upholds FramePresentationCallsProtocol. | 579 |
| f0210 | FramePresentationCallsProtocol has completeness = closed. | 580 |
| f0211 | FramePresentationCallsProtocol permits PresentFrameRequest. | 581 |
| f0212 | FramePresentationCallsProtocol permits PresentFrameResult. | 582 |
| f0213 | FyneBackend consumes PreparedFramePort. | 583 |
| f0214 | FyneBackend consumes UiSessionPort. | 584 |
| f0215 | FyneBackend owns HandleFocusAndTextInput. | 585 |
| f0216 | FyneBackend owns PublishPresentation. | 586 |
| f0217 | FyneBackend owns ReconcileWidgets. | 587 |
| f0218 | FyneBackend owns ReleaseNativeWidgets. | 588 |
| f0219 | FyneBackend provides NativeInteraction. | 589 |
| f0220 | FyneBackend uses FramePresentationCalls as receiver of PresentFrameRequest in mode UiPreview. | 590 |
| f0221 | FyneBackend uses FramePresentationCalls as sender of PresentFrameResult in mode UiPreview. | 591 |
| f0222 | FyneBackend uses NativeUiActions as receiver of UiActionRejected in mode BoundExecution. | 592 |
| f0223 | FyneBackend uses NativeUiActions as receiver of UiActionResult in mode BoundExecution. | 593 |
| f0224 | FyneBackend uses NativeUiActions as sender of UiActionRequest in mode BoundExecution. | 594 |
| f0225 | FyneBackend uses UiGenerationEvents as receiver of UiGenerationNotices in mode BoundExecution. | 595 |
| f0226 | FyneBackend uses UiGenerationEvents as receiver of UiGenerationNotices in mode LiveEditing. | 596 |
| f0227 | FyneHost consumes DomainBindingPort. | 597 |
| f0228 | FyneHost consumes ReloadPort. | 598 |
| f0229 | FyneHost consumes SdlFrontendPort. | 599 |
| f0230 | FyneHost consumes SduiFrontendPort. | 600 |
| f0231 | FyneHost consumes SourceSnapshotPort. | 601 |
| f0232 | FyneHost consumes UiSessionPort. | 602 |
| f0233 | FyneHost consumes WidgetBackendPort. | 603 |
| f0234 | FyneHost owns ComposeInteractiveSession. | 604 |
| f0235 | FyneHost owns ScheduleUiPublication. | 605 |
| f0236 | FyneHost provides NativeInteraction. | 606 |
| f0237 | FyneHost uses FramePresentationCalls as receiver of PresentFrameResult in mode UiPreview. | 607 |
| f0238 | FyneHost uses FramePresentationCalls as sender of PresentFrameRequest in mode UiPreview. | 608 |
| f0239 | FyneHost uses GoDomainCalls as receiver of DomainActionResult in mode UiPreview. | 609 |
| f0240 | FyneHost uses GoDomainCalls as sender of DomainActionRequest in mode UiPreview. | 610 |
| f0241 | FyneHost uses LayoutCalls as receiver of LayoutResult in mode UiPreview. | 611 |
| f0242 | FyneHost uses LayoutCalls as sender of LayoutRequest in mode UiPreview. | 612 |
| f0243 | G1FrontendPort delivers StructuralModelInspection. | 613 |
| f0244 | G1FrontendPort has implementation-status = planned. | 614 |
| f0245 | G1M1ParserAndAst addresses BuildSduiAst. | 615 |
| f0246 | G1M1ParserAndAst addresses IdentifySourceRevision. | 616 |
| f0247 | G1M1ParserAndAst addresses ReadBoundedSources. | 617 |
| f0248 | G1M1ParserAndAst addresses TokenizeSduiSource. | 618 |
| f0249 | G1M1ParserAndAst has implementation-status = planned. | 619 |
| f0250 | G1M1ParserAndAst refines G1FrontendPort. | 620 |
| f0251 | G1M2ValidationAndNormalization addresses CoordinateSduiCompilation. | 621 |
| f0252 | G1M2ValidationAndNormalization addresses ExpandUiDefinitions. | 622 |
| f0253 | G1M2ValidationAndNormalization addresses PreserveUiRegions. | 623 |
| f0254 | G1M2ValidationAndNormalization addresses PreserveUiSourceMap. | 624 |
| f0255 | G1M2ValidationAndNormalization addresses ResolveUiNames. | 625 |
| f0256 | G1M2ValidationAndNormalization addresses ValidateRelativeFormatting. | 626 |
| f0257 | G1M2ValidationAndNormalization addresses ValidateSymbolicBindings. | 627 |
| f0258 | G1M2ValidationAndNormalization addresses ValidateWidgetArguments. | 628 |
| f0259 | G1M2ValidationAndNormalization depends-on G1M1ParserAndAst. | 629 |
| f0260 | G1M2ValidationAndNormalization has implementation-status = planned. | 630 |
| f0261 | G1M2ValidationAndNormalization refines G1FrontendPort. | 631 |
| f0262 | G1M3Concept1AndDumps addresses ExportConsoleSnapshot. | 632 |
| f0263 | G1M3Concept1AndDumps addresses ReportSourceDiagnostics. | 633 |
| f0264 | G1M3Concept1AndDumps depends-on G1M2ValidationAndNormalization. | 634 |
| f0265 | G1M3Concept1AndDumps has implementation-status = planned. | 635 |
| f0266 | G1M3Concept1AndDumps refines G1FrontendPort. | 636 |
| f0267 | G2LayoutAndPresentation delivers InteractiveUiPreview. | 637 |
| f0268 | G2LayoutAndPresentation has implementation-status = planned. | 638 |
| f0269 | G2M1RelativeMeasurement addresses AllocateGeometry. | 639 |
| f0270 | G2M1RelativeMeasurement addresses ComputeClipping. | 640 |
| f0271 | G2M1RelativeMeasurement addresses MeasureUiContent. | 641 |
| f0272 | G2M1RelativeMeasurement addresses ResolveAncestorDimensions. | 642 |
| f0273 | G2M1RelativeMeasurement depends-on G1M2ValidationAndNormalization. | 643 |
| f0274 | G2M1RelativeMeasurement has implementation-status = planned. | 644 |
| f0275 | G2M1RelativeMeasurement refines G2LayoutAndPresentation. | 645 |
| f0276 | G2M2SharedSvgGeometry addresses BuildPreparedFrame. | 646 |
| f0277 | G2M2SharedSvgGeometry addresses ExportSvgSnapshot. | 647 |
| f0278 | G2M2SharedSvgGeometry depends-on G2M1RelativeMeasurement. | 648 |
| f0279 | G2M2SharedSvgGeometry has implementation-status = planned. | 649 |
| f0280 | G2M2SharedSvgGeometry refines G2LayoutAndPresentation. | 650 |
| f0281 | G2M3FyneInteractions addresses ComposeInteractiveSession. | 651 |
| f0282 | G2M3FyneInteractions addresses HandleFocusAndTextInput. | 652 |
| f0283 | G2M3FyneInteractions addresses PublishPresentation. | 653 |
| f0284 | G2M3FyneInteractions addresses ReconcileWidgets. | 654 |
| f0285 | G2M3FyneInteractions addresses ReleaseNativeWidgets. | 655 |
| f0286 | G2M3FyneInteractions addresses ScheduleUiPublication. | 656 |
| f0287 | G2M3FyneInteractions depends-on G2M2SharedSvgGeometry. | 657 |
| f0288 | G2M3FyneInteractions has implementation-status = planned. | 658 |
| f0289 | G2M3FyneInteractions refines G2LayoutAndPresentation. | 659 |
| f0290 | G2M4RichContent addresses MeasureMarkdownContent. | 660 |
| f0291 | G2M4RichContent addresses PrepareDiagramResource. | 661 |
| f0292 | G2M4RichContent addresses PrepareMarkdown. | 662 |
| f0293 | G2M4RichContent addresses ReleaseVisualResources. | 663 |
| f0294 | G2M4RichContent addresses ValidateVisualResources. | 664 |
| f0295 | G2M4RichContent depends-on G2M3FyneInteractions. | 665 |
| f0296 | G2M4RichContent has implementation-status = planned. | 666 |
| f0297 | G2M4RichContent refines G2LayoutAndPresentation. | 667 |
| f0298 | G3M1TypedUiSession addresses ApplyPropertyBatch. | 668 |
| f0299 | G3M1TypedUiSession addresses CloseUiInstance. | 669 |
| f0300 | G3M1TypedUiSession addresses CorrelateUiResult. | 670 |
| f0301 | G3M1TypedUiSession addresses CreateUiInstance. | 671 |
| f0302 | G3M1TypedUiSession addresses DispatchUiEvent. | 672 |
| f0303 | G3M1TypedUiSession addresses ManageWidgetIdentities. | 673 |
| f0304 | G3M1TypedUiSession addresses ProjectUiGeneration. | 674 |
| f0305 | G3M1TypedUiSession addresses RejectStaleUiEvent. | 675 |
| f0306 | G3M1TypedUiSession addresses RevokeWidgetGenerations. | 676 |
| f0307 | G3M1TypedUiSession addresses SnapshotUiState. | 677 |
| f0308 | G3M1TypedUiSession addresses TrackInputDraft. | 678 |
| f0309 | G3M1TypedUiSession addresses ValidatePropertyBatch. | 679 |
| f0310 | G3M1TypedUiSession addresses ValidateUiEvent. | 680 |
| f0311 | G3M1TypedUiSession depends-on G1M2ValidationAndNormalization. | 681 |
| f0312 | G3M1TypedUiSession depends-on G2M3FyneInteractions. | 682 |
| f0313 | G3M1TypedUiSession has implementation-status = planned. | 683 |
| f0314 | G3M1TypedUiSession refines G3UiRuntimeAndReload. | 684 |
| f0315 | G3M2CandidatePublication addresses CoalesceSourceChanges. | 685 |
| f0316 | G3M2CandidatePublication addresses KeepLastValidModels. | 686 |
| f0317 | G3M2CandidatePublication addresses ObserveSourceChanges. | 687 |
| f0318 | G3M2CandidatePublication addresses PrepareCandidateModels. | 688 |
| f0319 | G3M2CandidatePublication addresses PublishModelGeneration. | 689 |
| f0320 | G3M2CandidatePublication addresses RetirePreviousGeneration. | 690 |
| f0321 | G3M2CandidatePublication depends-on G3M1TypedUiSession. | 691 |
| f0322 | G3M2CandidatePublication has implementation-status = planned. | 692 |
| f0323 | G3M2CandidatePublication refines G3UiRuntimeAndReload. | 693 |
| f0324 | G3M3CompatibleState addresses MatchCompatibleWidgets. | 694 |
| f0325 | G3M3CompatibleState addresses PreserveCompatibleUiState. | 695 |
| f0326 | G3M3CompatibleState addresses ResetIncompatibleUiState. | 696 |
| f0327 | G3M3CompatibleState depends-on G3M2CandidatePublication. | 697 |
| f0328 | G3M3CompatibleState has implementation-status = planned. | 698 |
| f0329 | G3M3CompatibleState refines G3UiRuntimeAndReload. | 699 |
| f0330 | G3UiRuntimeAndReload delivers LiveModelReload. | 700 |
| f0331 | G3UiRuntimeAndReload has implementation-status = planned. | 701 |
| f0332 | G4M1SdlFrontend addresses BuildSdlAst. | 702 |
| f0333 | G4M1SdlFrontend addresses CoordinateSdlCompilation. | 703 |
| f0334 | G4M1SdlFrontend addresses NormalizeSdlModel. | 704 |
| f0335 | G4M1SdlFrontend addresses PreserveSdlSourceMap. | 705 |
| f0336 | G4M1SdlFrontend addresses ResolveSdlSymbols. | 706 |
| f0337 | G4M1SdlFrontend addresses TokenizeSdlSource. | 707 |
| f0338 | G4M1SdlFrontend addresses ValidateSdlProfile. | 708 |
| f0339 | G4M1SdlFrontend addresses ValidateSdlStructure. | 709 |
| f0340 | G4M1SdlFrontend has implementation-status = planned. | 710 |
| f0341 | G4M1SdlFrontend refines G4SdlRuntimeAndBinding. | 711 |
| f0342 | G4M2TypedExecution addresses CancelPendingActions. | 712 |
| f0343 | G4M2TypedExecution addresses CheckExecutionCompleteness. | 713 |
| f0344 | G4M2TypedExecution addresses CheckFunctionSignatures. | 714 |
| f0345 | G4M2TypedExecution addresses CloseSdlInstance. | 715 |
| f0346 | G4M2TypedExecution addresses CorrelateActionResult. | 716 |
| f0347 | G4M2TypedExecution addresses CreateSdlInstance. | 717 |
| f0348 | G4M2TypedExecution addresses InvokeRegisteredFunction. | 718 |
| f0349 | G4M2TypedExecution addresses ManageDomainState. | 719 |
| f0350 | G4M2TypedExecution addresses PerformDomainOperation. | 720 |
| f0351 | G4M2TypedExecution addresses RegisterDomainFunctions. | 721 |
| f0352 | G4M2TypedExecution addresses SnapshotDomainState. | 722 |
| f0353 | G4M2TypedExecution addresses ValidateActionInput. | 723 |
| f0354 | G4M2TypedExecution depends-on G4M1SdlFrontend. | 724 |
| f0355 | G4M2TypedExecution has implementation-status = planned. | 725 |
| f0356 | G4M2TypedExecution refines G4SdlRuntimeAndBinding. | 726 |
| f0357 | G4M3UiDomainBinding addresses ConnectTypedWidgetHandles. | 727 |
| f0358 | G4M3UiDomainBinding addresses DisconnectBindings. | 728 |
| f0359 | G4M3UiDomainBinding addresses PublishDomainUpdates. | 729 |
| f0360 | G4M3UiDomainBinding addresses ReportBindingDiagnostics. | 730 |
| f0361 | G4M3UiDomainBinding addresses ResolveCallbackSymbols. | 731 |
| f0362 | G4M3UiDomainBinding addresses RouteDomainBindings. | 732 |
| f0363 | G4M3UiDomainBinding depends-on G3M1TypedUiSession. | 733 |
| f0364 | G4M3UiDomainBinding depends-on G4M2TypedExecution. | 734 |
| f0365 | G4M3UiDomainBinding has implementation-status = planned. | 735 |
| f0366 | G4M3UiDomainBinding refines G4SdlRuntimeAndBinding. | 736 |
| f0367 | G4M4DomainReload addresses CheckDomainStateCompatibility. | 737 |
| f0368 | G4M4DomainReload addresses MigrateOrResetDomainState. | 738 |
| f0369 | G4M4DomainReload addresses RestartChangedGoProgram. | 739 |
| f0370 | G4M4DomainReload depends-on G3M3CompatibleState. | 740 |
| f0371 | G4M4DomainReload depends-on G4M3UiDomainBinding. | 741 |
| f0372 | G4M4DomainReload has implementation-status = planned. | 742 |
| f0373 | G4M4DomainReload refines G4SdlRuntimeAndBinding. | 743 |
| f0374 | G4SdlRuntimeAndBinding delivers TypedDomainBinding. | 744 |
| f0375 | G4SdlRuntimeAndBinding has implementation-status = planned. | 745 |
| f0376 | G5M1GeneratedGo addresses BuildGeneratedApplication. | 746 |
| f0377 | G5M1GeneratedGo addresses GenerateBindingRegistration. | 747 |
| f0378 | G5M1GeneratedGo addresses GenerateModelConstructors. | 748 |
| f0379 | G5M1GeneratedGo addresses PreserveHandwrittenSources. | 749 |
| f0380 | G5M1GeneratedGo depends-on G4M4DomainReload. | 750 |
| f0381 | G5M1GeneratedGo has implementation-status = planned. | 751 |
| f0382 | G5M1GeneratedGo refines G5NativeGeneration. | 752 |
| f0383 | G5M2BehaviorParity addresses VerifyNativeBehaviorParity. | 753 |
| f0384 | G5M2BehaviorParity depends-on G5M1GeneratedGo. | 754 |
| f0385 | G5M2BehaviorParity has implementation-status = planned. | 755 |
| f0386 | G5M2BehaviorParity refines G5NativeGeneration. | 756 |
| f0387 | G5M3DocumentationExport addresses ComposeHeadlessExport. | 757 |
| f0388 | G5M3DocumentationExport addresses ComposeMarkdownDocument. | 758 |
| f0389 | G5M3DocumentationExport addresses ExportViewpointMarkdown. | 759 |
| f0390 | G5M3DocumentationExport addresses ProjectSdlViewpoints. | 760 |
| f0391 | G5M3DocumentationExport addresses TraceViewpointFacts. | 761 |
| f0392 | G5M3DocumentationExport addresses WriteGeneratedArtifacts. | 762 |
| f0393 | G5M3DocumentationExport depends-on G2M4RichContent. | 763 |
| f0394 | G5M3DocumentationExport depends-on G5M2BehaviorParity. | 764 |
| f0395 | G5M3DocumentationExport has implementation-status = planned. | 765 |
| f0396 | G5M3DocumentationExport refines G5NativeGeneration. | 766 |
| f0397 | G5M4RetirePython addresses RetireReplacedPythonEntryPoints. | 767 |
| f0398 | G5M4RetirePython depends-on G1M3Concept1AndDumps. | 768 |
| f0399 | G5M4RetirePython depends-on G5M3DocumentationExport. | 769 |
| f0400 | G5M4RetirePython has implementation-status = planned. | 770 |
| f0401 | G5M4RetirePython refines G5NativeGeneration. | 771 |
| f0402 | G5NativeGeneration delivers DesignDocumentation. | 772 |
| f0403 | G5NativeGeneration delivers NativeGoAssembly. | 773 |
| f0404 | G5NativeGeneration has implementation-status = planned. | 774 |
| f0405 | GenerateBindingRegistration allocated-to CommandLineHost in mode NativeBuild. | 775 |
| f0406 | GenerateBindingRegistration contributes-to NativeGoAssembly. | 776 |
| f0407 | GenerateBindingRegistration realizes NativeRealization. | 777 |
| f0408 | GenerateGoRequest has message-kind = request. | 778 |
| f0409 | GenerateGoRequest upholds NormalizedModelContract. | 779 |
| f0410 | GenerateGoResult has message-kind = result. | 780 |
| f0411 | GenerateGoResult replies-to GenerateGoRequest. | 781 |
| f0412 | GenerateGoResult upholds GeneratedGoContract. | 782 |
| f0413 | GenerateModelConstructors allocated-to CommandLineHost in mode NativeBuild. | 783 |
| f0414 | GenerateModelConstructors contributes-to NativeGoAssembly. | 784 |
| f0415 | GenerateModelConstructors realizes NativeRealization. | 785 |
| f0416 | GeneratedGoContract has completeness = closed. | 786 |
| f0417 | GeneratedGoContract has-field GeneratedGoSources. | 787 |
| f0418 | GeneratedGoSources has presence = required. | 788 |
| f0419 | GeneratedGoSources has value-type = bytes. | 789 |
| f0420 | GoBuildCalls upholds GoBuildCallsProtocol. | 790 |
| f0421 | GoBuildCallsProtocol has completeness = closed. | 791 |
| f0422 | GoBuildCallsProtocol permits BuildGoRequest. | 792 |
| f0423 | GoBuildCallsProtocol permits BuildGoResult. | 793 |
| f0424 | GoBuildRunner consumes BuildToolPort. | 794 |
| f0425 | GoBuildRunner consumes GeneratedArtifactPort. | 795 |
| f0426 | GoBuildRunner owns BuildGeneratedApplication. | 796 |
| f0427 | GoBuildRunner owns RestartChangedGoProgram. | 797 |
| f0428 | GoBuildRunner provides NativeRealization. | 798 |
| f0429 | GoBuildRunner uses GoBuildCalls as receiver of BuildGoRequest in mode NativeBuild. | 799 |
| f0430 | GoBuildRunner uses GoBuildCalls as sender of BuildGoResult in mode NativeBuild. | 800 |
| f0431 | GoCodeGenerator consumes ExecutionProfilePort. | 801 |
| f0432 | GoCodeGenerator consumes SdlModelPort. | 802 |
| f0433 | GoCodeGenerator consumes SduiModelPort. | 803 |
| f0434 | GoCodeGenerator owns GenerateBindingRegistration. | 804 |
| f0435 | GoCodeGenerator owns GenerateModelConstructors. | 805 |
| f0436 | GoCodeGenerator owns PreserveHandwrittenSources. | 806 |
| f0437 | GoCodeGenerator provides NativeRealization. | 807 |
| f0438 | GoCodeGenerator uses GoGenerationCalls as receiver of GenerateGoRequest in mode NativeBuild. | 808 |
| f0439 | GoCodeGenerator uses GoGenerationCalls as sender of GenerateGoResult in mode NativeBuild. | 809 |
| f0440 | GoDomainCalls upholds GoDomainCallsProtocol. | 810 |
| f0441 | GoDomainCallsProtocol has completeness = closed. | 811 |
| f0442 | GoDomainCallsProtocol permits DomainActionRequest. | 812 |
| f0443 | GoDomainCallsProtocol permits DomainActionResult. | 813 |
| f0444 | GoDomainImplementation owns PerformDomainOperation. | 814 |
| f0445 | GoDomainImplementation provides DomainOperations. | 815 |
| f0446 | GoDomainImplementation uses GoDomainCalls as receiver of DomainActionRequest in mode BoundExecution. | 816 |
| f0447 | GoDomainImplementation uses GoDomainCalls as receiver of DomainActionRequest in mode UiPreview. | 817 |
| f0448 | GoDomainImplementation uses GoDomainCalls as sender of DomainActionResult in mode BoundExecution. | 818 |
| f0449 | GoDomainImplementation uses GoDomainCalls as sender of DomainActionResult in mode UiPreview. | 819 |
| f0450 | GoGenerationCalls upholds GoGenerationCallsProtocol. | 820 |
| f0451 | GoGenerationCallsProtocol has completeness = closed. | 821 |
| f0452 | GoGenerationCallsProtocol permits GenerateGoRequest. | 822 |
| f0453 | GoGenerationCallsProtocol permits GenerateGoResult. | 823 |
| f0454 | HandleFocusAndTextInput allocated-to FyneHost in mode UiPreview. | 824 |
| f0455 | HandleFocusAndTextInput contributes-to InteractiveUiPreview. | 825 |
| f0456 | HandleFocusAndTextInput realizes NativeInteraction. | 826 |
| f0457 | IdentifySourceRevision realizes SourceLoading. | 827 |
| f0458 | InspectSdlSource refines InspectDesignSource. | 828 |
| f0459 | InspectSduiSource refines InspectDesignSource. | 829 |
| f0460 | InteractiveFramePrepared exercises PrototypeUserInterface. | 830 |
| f0461 | InteractiveFramePrepared has completeness = closed. | 831 |
| f0462 | InteractiveFramePrepared illustrates G2M1RelativeMeasurement. | 832 |
| f0463 | InteractiveFramePrepared illustrates G2M3FyneInteractions. | 833 |
| f0464 | InteractiveFramePrepared runs-in UiPreview. | 834 |
| f0465 | InteractiveFramePrepared step 1 sends LayoutRequest from FyneHost to SduiLayout via LayoutCalls. | 835 |
| f0466 | InteractiveFramePrepared step 2 sends LayoutResult from SduiLayout to FyneHost via LayoutCalls reply-to 1. | 836 |
| f0467 | InteractiveFramePrepared step 3 sends PresentFrameRequest from FyneHost to FyneBackend via FramePresentationCalls. | 837 |
| f0468 | InteractiveFramePrepared step 4 sends PresentFrameResult from FyneBackend to FyneHost via FramePresentationCalls reply-to 3. | 838 |
| f0469 | InteractiveUiPreview supports PrototypeUserInterface. | 839 |
| f0470 | InvokeRegisteredFunction allocated-to FyneHost in mode BoundExecution. | 840 |
| f0471 | InvokeRegisteredFunction contributes-to TypedDomainBinding. | 841 |
| f0472 | InvokeRegisteredFunction realizes ExecutableDesign. | 842 |
| f0473 | KeepLastValidModels allocated-to FyneHost in mode LiveEditing. | 843 |
| f0474 | KeepLastValidModels contributes-to LiveModelReload. | 844 |
| f0475 | KeepLastValidModels realizes DevelopmentReload. | 845 |
| f0476 | LayoutArguments has completeness = closed. | 846 |
| f0477 | LayoutArguments has-field LayoutModelArtifact. | 847 |
| f0478 | LayoutArguments has-field LayoutViewportHeight. | 848 |
| f0479 | LayoutArguments has-field LayoutViewportWidth. | 849 |
| f0480 | LayoutCalls upholds LayoutCallsProtocol. | 850 |
| f0481 | LayoutCallsProtocol has completeness = closed. | 851 |
| f0482 | LayoutCallsProtocol permits LayoutRequest. | 852 |
| f0483 | LayoutCallsProtocol permits LayoutResult. | 853 |
| f0484 | LayoutModelArtifact has presence = required. | 854 |
| f0485 | LayoutModelArtifact has value-type = bytes. | 855 |
| f0486 | LayoutRequest has message-kind = request. | 856 |
| f0487 | LayoutRequest upholds LayoutArguments. | 857 |
| f0488 | LayoutResult has message-kind = result. | 858 |
| f0489 | LayoutResult replies-to LayoutRequest. | 859 |
| f0490 | LayoutResult upholds PreparedFrameContract. | 860 |
| f0491 | LayoutViewportHeight has presence = required. | 861 |
| f0492 | LayoutViewportHeight has value-type = decimal. | 862 |
| f0493 | LayoutViewportWidth has presence = required. | 863 |
| f0494 | LayoutViewportWidth has value-type = decimal. | 864 |
| f0495 | LiveModelReload supports EditRunningPrototype. | 865 |
| f0496 | ManageDomainState has state-retention = stateful. | 866 |
| f0497 | ManageDomainState realizes ExecutableDesign. | 867 |
| f0498 | ManageWidgetIdentities has state-retention = stateful. | 868 |
| f0499 | ManageWidgetIdentities realizes InteractiveSession. | 869 |
| f0500 | MarkdownProvider consumes DiagramPort. | 870 |
| f0501 | MarkdownProvider consumes MeasurementPort. | 871 |
| f0502 | MarkdownProvider consumes ResourcePort. | 872 |
| f0503 | MarkdownProvider owns MeasureMarkdownContent. | 873 |
| f0504 | MarkdownProvider owns PrepareMarkdown. | 874 |
| f0505 | MarkdownProvider provides RichContent. | 875 |
| f0506 | MatchCompatibleWidgets realizes InteractiveSession. | 876 |
| f0507 | MeasureMarkdownContent realizes RichContent. | 877 |
| f0508 | MeasureUiContent realizes MeasuredPresentation. | 878 |
| f0509 | MeasuredPresentation requires MeasurementPort in mode StaticExport. | 879 |
| f0510 | MeasuredPresentation requires MeasurementPort in mode UiPreview. | 880 |
| f0511 | MigrateOrResetDomainState realizes ExecutableDesign. | 881 |
| f0512 | ModelDiagnostics has presence = optional. | 882 |
| f0513 | ModelDiagnostics has value-type = text. | 883 |
| f0514 | ModelIsValid has presence = required. | 884 |
| f0515 | ModelIsValid has value-type = boolean. | 885 |
| f0516 | ModelReloadCalls upholds ModelReloadCallsProtocol. | 886 |
| f0517 | ModelReloadCallsProtocol has completeness = closed. | 887 |
| f0518 | ModelReloadCallsProtocol permits ReloadPublished. | 888 |
| f0519 | ModelReloadCallsProtocol permits ReloadRejected. | 889 |
| f0520 | ModelReloadCallsProtocol permits ReloadRequest. | 890 |
| f0521 | NativeBuildContract has completeness = closed. | 891 |
| f0522 | NativeBuildContract has-field NativeBuildDiagnostics. | 892 |
| f0523 | NativeBuildContract has-field NativeBuildSucceeded. | 893 |
| f0524 | NativeBuildDiagnostics has presence = optional. | 894 |
| f0525 | NativeBuildDiagnostics has value-type = text. | 895 |
| f0526 | NativeBuildSucceeded has presence = required. | 896 |
| f0527 | NativeBuildSucceeded has value-type = boolean. | 897 |
| f0528 | NativeGoAssembly supports BuildNativeProduct. | 898 |
| f0529 | NativeInteraction requires WidgetBackendPort in mode UiPreview. | 899 |
| f0530 | NativeProgramBuilt exercises BuildNativeProduct. | 900 |
| f0531 | NativeProgramBuilt has completeness = closed. | 901 |
| f0532 | NativeProgramBuilt illustrates G5M1GeneratedGo. | 902 |
| f0533 | NativeProgramBuilt runs-in NativeBuild. | 903 |
| f0534 | NativeProgramBuilt step 1 sends GenerateGoRequest from CommandLineHost to GoCodeGenerator via GoGenerationCalls. | 904 |
| f0535 | NativeProgramBuilt step 2 sends GenerateGoResult from GoCodeGenerator to CommandLineHost via GoGenerationCalls reply-to 1. | 905 |
| f0536 | NativeProgramBuilt step 3 sends BuildGoRequest from CommandLineHost to GoBuildRunner via GoBuildCalls. | 906 |
| f0537 | NativeProgramBuilt step 4 sends BuildGoResult from GoBuildRunner to CommandLineHost via GoBuildCalls reply-to 3. | 907 |
| f0538 | NativeRealization requires BuildToolPort in mode NativeBuild. | 908 |
| f0539 | NativeRealization requires GeneratedArtifactPort in mode NativeBuild. | 909 |
| f0540 | NativeUiActions upholds NativeUiActionsProtocol. | 910 |
| f0541 | NativeUiActionsProtocol has completeness = closed. | 911 |
| f0542 | NativeUiActionsProtocol permits UiActionRejected. | 912 |
| f0543 | NativeUiActionsProtocol permits UiActionRequest. | 913 |
| f0544 | NativeUiActionsProtocol permits UiActionResult. | 914 |
| f0545 | NormalizeSdlModel has repeatability = deterministic. | 915 |
| f0546 | NormalizeSdlModel has state-retention = stateless. | 916 |
| f0547 | NormalizeSdlModel realizes SdlSourceModel. | 917 |
| f0548 | NormalizeUiRequest has message-kind = request. | 918 |
| f0549 | NormalizeUiRequest upholds AstArtifactContract. | 919 |
| f0550 | NormalizeUiResult has message-kind = result. | 920 |
| f0551 | NormalizeUiResult replies-to NormalizeUiRequest. | 921 |
| f0552 | NormalizeUiResult upholds NormalizedModelContract. | 922 |
| f0553 | NormalizedModelArtifact has presence = required. | 923 |
| f0554 | NormalizedModelArtifact has value-type = bytes. | 924 |
| f0555 | NormalizedModelContract has completeness = closed. | 925 |
| f0556 | NormalizedModelContract has-field NormalizedModelArtifact. | 926 |
| f0557 | NoticeGeneration has presence = required. | 927 |
| f0558 | NoticeGeneration has value-type = unsigned. | 928 |
| f0559 | NoticeVersion has presence = required. | 929 |
| f0560 | NoticeVersion has value-type = unsigned. | 930 |
| f0561 | ObserveSourceChanges allocated-to FyneHost in mode LiveEditing. | 931 |
| f0562 | ObserveSourceChanges contributes-to LiveModelReload. | 932 |
| f0563 | ObserveSourceChanges realizes DevelopmentReload. | 933 |
| f0564 | PerformDomainOperation realizes DomainOperations. | 934 |
| f0565 | PrepareCandidateModels allocated-to FyneHost in mode LiveEditing. | 935 |
| f0566 | PrepareCandidateModels contributes-to LiveModelReload. | 936 |
| f0567 | PrepareCandidateModels realizes DevelopmentReload. | 937 |
| f0568 | PrepareDiagramResource realizes RichContent. | 938 |
| f0569 | PrepareMarkdown realizes RichContent. | 939 |
| f0570 | PreparedFrameArtifact has presence = required. | 940 |
| f0571 | PreparedFrameArtifact has value-type = bytes. | 941 |
| f0572 | PreparedFrameContract has completeness = closed. | 942 |
| f0573 | PreparedFrameContract has-field PreparedFrameArtifact. | 943 |
| f0574 | PresentFrameRequest has message-kind = request. | 944 |
| f0575 | PresentFrameRequest upholds PreparedFrameContract. | 945 |
| f0576 | PresentFrameResult has message-kind = result. | 946 |
| f0577 | PresentFrameResult replies-to PresentFrameRequest. | 947 |
| f0578 | PresentFrameResult upholds PresentationOutcome. | 948 |
| f0579 | PresentationOutcome has completeness = closed. | 949 |
| f0580 | PresentationOutcome has-field PresentationReady. | 950 |
| f0581 | PresentationReady has presence = required. | 951 |
| f0582 | PresentationReady has value-type = boolean. | 952 |
| f0583 | PreserveCompatibleUiState allocated-to FyneHost in mode LiveEditing. | 953 |
| f0584 | PreserveCompatibleUiState contributes-to LiveModelReload. | 954 |
| f0585 | PreserveCompatibleUiState realizes InteractiveSession. | 955 |
| f0586 | PreserveHandwrittenSources allocated-to CommandLineHost in mode NativeBuild. | 956 |
| f0587 | PreserveHandwrittenSources contributes-to NativeGoAssembly. | 957 |
| f0588 | PreserveHandwrittenSources realizes NativeRealization. | 958 |
| f0589 | PreserveSdlSourceMap has repeatability = deterministic. | 959 |
| f0590 | PreserveSdlSourceMap has state-retention = stateless. | 960 |
| f0591 | PreserveSdlSourceMap realizes SdlSourceModel. | 961 |
| f0592 | PreserveUiRegions has repeatability = deterministic. | 962 |
| f0593 | PreserveUiRegions has state-retention = stateless. | 963 |
| f0594 | PreserveUiRegions realizes SduiSourceModel. | 964 |
| f0595 | PreserveUiSourceMap has repeatability = deterministic. | 965 |
| f0596 | PreserveUiSourceMap has state-retention = stateless. | 966 |
| f0597 | PreserveUiSourceMap realizes SduiSourceModel. | 967 |
| f0598 | ProjectSdlViewpoints allocated-to CommandLineHost in mode StaticExport. | 968 |
| f0599 | ProjectSdlViewpoints contributes-to DesignDocumentation. | 969 |
| f0600 | ProjectSdlViewpoints has repeatability = deterministic. | 970 |
| f0601 | ProjectSdlViewpoints has state-retention = stateless. | 971 |
| f0602 | ProjectUiGeneration allocated-to FyneHost in mode LiveEditing. | 972 |
| f0603 | ProjectUiGeneration contributes-to LiveModelReload. | 973 |
| f0604 | ProjectUiGeneration projects UiSessionState into UiGenerationNotices. | 974 |
| f0605 | PublishDomainUpdates allocated-to FyneHost in mode BoundExecution. | 975 |
| f0606 | PublishDomainUpdates contributes-to TypedDomainBinding. | 976 |
| f0607 | PublishDomainUpdates realizes BoundInteraction. | 977 |
| f0608 | PublishModelGeneration allocated-to FyneHost in mode LiveEditing. | 978 |
| f0609 | PublishModelGeneration contributes-to LiveModelReload. | 979 |
| f0610 | PublishModelGeneration has state-retention = stateful. | 980 |
| f0611 | PublishModelGeneration realizes DevelopmentReload. | 981 |
| f0612 | PublishPresentation realizes NativeInteraction. | 982 |
| f0613 | ReadBoundedSources realizes SourceLoading. | 983 |
| f0614 | ReconcileWidgets allocated-to FyneHost in mode UiPreview. | 984 |
| f0615 | ReconcileWidgets contributes-to InteractiveUiPreview. | 985 |
| f0616 | ReconcileWidgets has state-retention = stateful. | 986 |
| f0617 | ReconcileWidgets realizes NativeInteraction. | 987 |
| f0618 | RegisterDomainFunctions has state-retention = stateful. | 988 |
| f0619 | RegisterDomainFunctions realizes ExecutableDesign. | 989 |
| f0620 | RejectStaleUiEvent realizes InteractiveSession. | 990 |
| f0621 | ReleaseNativeWidgets realizes NativeInteraction. | 991 |
| f0622 | ReleaseVisualResources realizes RichContent. | 992 |
| f0623 | ReloadArguments has completeness = closed. | 993 |
| f0624 | ReloadArguments has-field ReloadSourceRevision. | 994 |
| f0625 | ReloadArguments has-field ReloadSourceText. | 995 |
| f0626 | ReloadBoundModels refines ReloadDesignSession. | 996 |
| f0627 | ReloadCoordinator consumes BindingReloadPort. | 997 |
| f0628 | ReloadCoordinator consumes DiagnosticPort. | 998 |
| f0629 | ReloadCoordinator consumes SdlFrontendPort. | 999 |
| f0630 | ReloadCoordinator consumes SdlReloadPort. | 1000 |
| f0631 | ReloadCoordinator consumes SduiFrontendPort. | 1001 |
| f0632 | ReloadCoordinator consumes SourceSnapshotPort. | 1002 |
| f0633 | ReloadCoordinator consumes UiReloadPort. | 1003 |
| f0634 | ReloadCoordinator owns KeepLastValidModels. | 1004 |
| f0635 | ReloadCoordinator owns PrepareCandidateModels. | 1005 |
| f0636 | ReloadCoordinator owns PublishModelGeneration. | 1006 |
| f0637 | ReloadCoordinator owns RetirePreviousGeneration. | 1007 |
| f0638 | ReloadCoordinator provides DevelopmentReload. | 1008 |
| f0639 | ReloadCoordinator uses ModelReloadCalls as receiver of ReloadRequest in mode BoundLiveEditing. | 1009 |
| f0640 | ReloadCoordinator uses ModelReloadCalls as receiver of ReloadRequest in mode LiveEditing. | 1010 |
| f0641 | ReloadCoordinator uses ModelReloadCalls as sender of ReloadPublished in mode BoundLiveEditing. | 1011 |
| f0642 | ReloadCoordinator uses ModelReloadCalls as sender of ReloadPublished in mode LiveEditing. | 1012 |
| f0643 | ReloadCoordinator uses ModelReloadCalls as sender of ReloadRejected in mode BoundLiveEditing. | 1013 |
| f0644 | ReloadCoordinator uses ModelReloadCalls as sender of ReloadRejected in mode LiveEditing. | 1014 |
| f0645 | ReloadCoordinator uses SdlCompilationCalls as receiver of CompileSdlResult in mode BoundLiveEditing. | 1015 |
| f0646 | ReloadCoordinator uses SdlCompilationCalls as sender of CompileSdlRequest in mode BoundLiveEditing. | 1016 |
| f0647 | ReloadCoordinator uses UiCompilationCalls as receiver of CompileUiRejected in mode LiveEditing. | 1017 |
| f0648 | ReloadCoordinator uses UiCompilationCalls as receiver of CompileUiResult in mode LiveEditing. | 1018 |
| f0649 | ReloadCoordinator uses UiCompilationCalls as sender of CompileUiRequest in mode LiveEditing. | 1019 |
| f0650 | ReloadDiagnostic has presence = optional. | 1020 |
| f0651 | ReloadDiagnostic has value-type = text. | 1021 |
| f0652 | ReloadOutcome has completeness = closed. | 1022 |
| f0653 | ReloadOutcome has-field ReloadDiagnostic. | 1023 |
| f0654 | ReloadOutcome has-field ReloadPublishedGeneration. | 1024 |
| f0655 | ReloadPublished has message-kind = result. | 1025 |
| f0656 | ReloadPublished replies-to ReloadRequest. | 1026 |
| f0657 | ReloadPublished upholds ReloadOutcome. | 1027 |
| f0658 | ReloadPublishedGeneration has presence = optional. | 1028 |
| f0659 | ReloadPublishedGeneration has value-type = unsigned. | 1029 |
| f0660 | ReloadRejected has message-kind = result. | 1030 |
| f0661 | ReloadRejected replies-to ReloadRequest. | 1031 |
| f0662 | ReloadRejected upholds ReloadOutcome. | 1032 |
| f0663 | ReloadRequest has message-kind = request. | 1033 |
| f0664 | ReloadRequest upholds ReloadArguments. | 1034 |
| f0665 | ReloadSourceRevision has presence = required. | 1035 |
| f0666 | ReloadSourceRevision has value-type = unsigned. | 1036 |
| f0667 | ReloadSourceText has presence = required. | 1037 |
| f0668 | ReloadSourceText has value-type = text. | 1038 |
| f0669 | ReloadUiModel refines ReloadDesignSession. | 1039 |
| f0670 | ReportBindingDiagnostics realizes SourceDiagnostics. | 1040 |
| f0671 | ReportSourceDiagnostics allocated-to CommandLineHost in mode SourceInspection. | 1041 |
| f0672 | ReportSourceDiagnostics allocated-to FyneHost in mode LiveEditing. | 1042 |
| f0673 | ReportSourceDiagnostics contributes-to LiveModelReload. | 1043 |
| f0674 | ReportSourceDiagnostics contributes-to StructuralModelInspection. | 1044 |
| f0675 | ReportSourceDiagnostics realizes SourceDiagnostics. | 1045 |
| f0676 | ResetIncompatibleUiState realizes InteractiveSession. | 1046 |
| f0677 | ResolveAncestorDimensions realizes MeasuredPresentation. | 1047 |
| f0678 | ResolveCallbackSymbols allocated-to FyneHost in mode BoundExecution. | 1048 |
| f0679 | ResolveCallbackSymbols contributes-to TypedDomainBinding. | 1049 |
| f0680 | ResolveCallbackSymbols realizes BoundInteraction. | 1050 |
| f0681 | ResolveSdlSymbols has repeatability = deterministic. | 1051 |
| f0682 | ResolveSdlSymbols has state-retention = stateless. | 1052 |
| f0683 | ResolveSdlSymbols realizes SdlSourceModel. | 1053 |
| f0684 | ResolveUiNames has repeatability = deterministic. | 1054 |
| f0685 | ResolveUiNames has state-retention = stateless. | 1055 |
| f0686 | ResolveUiNames realizes SduiSourceModel. | 1056 |
| f0687 | ResourceStore owns ReleaseVisualResources. | 1057 |
| f0688 | ResourceStore owns ValidateVisualResources. | 1058 |
| f0689 | ResourceStore provides RichContent. | 1059 |
| f0690 | RestartChangedGoProgram realizes NativeRealization. | 1060 |
| f0691 | RetirePreviousGeneration realizes DevelopmentReload. | 1061 |
| f0692 | RetireReplacedPythonEntryPoints allocated-to CommandLineHost in mode NativeBuild. | 1062 |
| f0693 | RetireReplacedPythonEntryPoints contributes-to NativeGoAssembly. | 1063 |
| f0694 | RevokeWidgetGenerations realizes InteractiveSession. | 1064 |
| f0695 | RichContent requires ContentProviderPort in mode RichDocument. | 1065 |
| f0696 | RichContent requires DiagramEnginePort in mode RichDocument. | 1066 |
| f0697 | RouteDomainBindings realizes BoundInteraction. | 1067 |
| f0698 | RunBoundUiAction refines RunDesignSession. | 1068 |
| f0699 | RunUnboundUiPreview refines RunDesignSession. | 1069 |
| f0700 | ScheduleUiPublication realizes NativeInteraction. | 1070 |
| f0701 | SdlActionCalls upholds SdlActionCallsProtocol. | 1071 |
| f0702 | SdlActionCallsProtocol has completeness = closed. | 1072 |
| f0703 | SdlActionCallsProtocol permits SdlActionRequest. | 1073 |
| f0704 | SdlActionCallsProtocol permits SdlActionResult. | 1074 |
| f0705 | SdlActionRequest has message-kind = request. | 1075 |
| f0706 | SdlActionRequest upholds ActionArguments. | 1076 |
| f0707 | SdlActionResult has message-kind = result. | 1077 |
| f0708 | SdlActionResult replies-to SdlActionRequest. | 1078 |
| f0709 | SdlActionResult upholds ActionOutcome. | 1079 |
| f0710 | SdlCompilationCalls upholds SdlCompilationCallsProtocol. | 1080 |
| f0711 | SdlCompilationCallsProtocol has completeness = closed. | 1081 |
| f0712 | SdlCompilationCallsProtocol permits CompileSdlRequest. | 1082 |
| f0713 | SdlCompilationCallsProtocol permits CompileSdlResult. | 1083 |
| f0714 | SdlDispatcher consumes DomainFunctionPort. | 1084 |
| f0715 | SdlDispatcher consumes DomainStatePort. | 1085 |
| f0716 | SdlDispatcher owns CancelPendingActions. | 1086 |
| f0717 | SdlDispatcher owns CorrelateActionResult. | 1087 |
| f0718 | SdlDispatcher owns InvokeRegisteredFunction. | 1088 |
| f0719 | SdlDispatcher owns ValidateActionInput. | 1089 |
| f0720 | SdlDispatcher provides ExecutableDesign. | 1090 |
| f0721 | SdlDispatcher uses GoDomainCalls as receiver of DomainActionResult in mode BoundExecution. | 1091 |
| f0722 | SdlDispatcher uses GoDomainCalls as sender of DomainActionRequest in mode BoundExecution. | 1092 |
| f0723 | SdlDispatcher uses SdlActionCalls as receiver of SdlActionRequest in mode BoundExecution. | 1093 |
| f0724 | SdlDispatcher uses SdlActionCalls as sender of SdlActionResult in mode BoundExecution. | 1094 |
| f0725 | SdlExecutionGate consumes DiagnosticPort. | 1095 |
| f0726 | SdlExecutionGate consumes SdlModelPort. | 1096 |
| f0727 | SdlExecutionGate owns CheckExecutionCompleteness. | 1097 |
| f0728 | SdlExecutionGate provides ExecutableDesign. | 1098 |
| f0729 | SdlFrontend consumes DiagnosticPort. | 1099 |
| f0730 | SdlFrontend consumes SourceSnapshotPort. | 1100 |
| f0731 | SdlFrontend contains SdlLexer. | 1101 |
| f0732 | SdlFrontend contains SdlNormalizer. | 1102 |
| f0733 | SdlFrontend contains SdlParser. | 1103 |
| f0734 | SdlFrontend contains SdlValidator. | 1104 |
| f0735 | SdlFrontend owns CoordinateSdlCompilation. | 1105 |
| f0736 | SdlFrontend provides SdlSourceModel. | 1106 |
| f0737 | SdlFrontend uses SdlCompilationCalls as receiver of CompileSdlRequest in mode BoundLiveEditing. | 1107 |
| f0738 | SdlFrontend uses SdlCompilationCalls as sender of CompileSdlResult in mode BoundLiveEditing. | 1108 |
| f0739 | SdlFunctionRegistry consumes DomainFunctionPort. | 1109 |
| f0740 | SdlFunctionRegistry owns CheckFunctionSignatures. | 1110 |
| f0741 | SdlFunctionRegistry owns RegisterDomainFunctions. | 1111 |
| f0742 | SdlFunctionRegistry provides ExecutableDesign. | 1112 |
| f0743 | SdlLexer owns TokenizeSdlSource. | 1113 |
| f0744 | SdlLexer provides SdlSourceModel. | 1114 |
| f0745 | SdlLibrary contains SdlFrontend. | 1115 |
| f0746 | SdlLibrary contains SdlRuntime. | 1116 |
| f0747 | SdlLibrary contains SdlViewpointGenerator. | 1117 |
| f0748 | SdlLibrary provides ExecutableDesign. | 1118 |
| f0749 | SdlLibrary provides SdlSourceModel. | 1119 |
| f0750 | SdlModelReloadAccepted exercises EditRunningPrototype. | 1120 |
| f0751 | SdlModelReloadAccepted has completeness = closed. | 1121 |
| f0752 | SdlModelReloadAccepted illustrates G4M1SdlFrontend. | 1122 |
| f0753 | SdlModelReloadAccepted illustrates G4M4DomainReload. | 1123 |
| f0754 | SdlModelReloadAccepted runs-in BoundLiveEditing. | 1124 |
| f0755 | SdlModelReloadAccepted step 1 sends ReloadRequest from SourceWatcher to ReloadCoordinator via ModelReloadCalls. | 1125 |
| f0756 | SdlModelReloadAccepted step 2 sends CompileSdlRequest from ReloadCoordinator to SdlFrontend via SdlCompilationCalls. | 1126 |
| f0757 | SdlModelReloadAccepted step 3 sends CompileSdlResult from SdlFrontend to ReloadCoordinator via SdlCompilationCalls reply-to 2. | 1127 |
| f0758 | SdlModelReloadAccepted step 4 sends ReloadPublished from ReloadCoordinator to SourceWatcher via ModelReloadCalls reply-to 1. | 1128 |
| f0759 | SdlNormalizer owns NormalizeSdlModel. | 1129 |
| f0760 | SdlNormalizer owns PreserveSdlSourceMap. | 1130 |
| f0761 | SdlNormalizer provides SdlSourceModel. | 1131 |
| f0762 | SdlParser owns BuildSdlAst. | 1132 |
| f0763 | SdlParser provides SdlSourceModel. | 1133 |
| f0764 | SdlRuntime consumes DomainFunctionPort. | 1134 |
| f0765 | SdlRuntime consumes SdlModelPort. | 1135 |
| f0766 | SdlRuntime contains DomainStateMigrator. | 1136 |
| f0767 | SdlRuntime contains SdlDispatcher. | 1137 |
| f0768 | SdlRuntime contains SdlExecutionGate. | 1138 |
| f0769 | SdlRuntime contains SdlFunctionRegistry. | 1139 |
| f0770 | SdlRuntime contains SdlStateStore. | 1140 |
| f0771 | SdlRuntime owns CloseSdlInstance. | 1141 |
| f0772 | SdlRuntime owns CreateSdlInstance. | 1142 |
| f0773 | SdlRuntime provides ExecutableDesign. | 1143 |
| f0774 | SdlSourceModel requires SourceSnapshotPort in mode SourceInspection. | 1144 |
| f0775 | SdlStateStore owns ManageDomainState. | 1145 |
| f0776 | SdlStateStore owns SnapshotDomainState. | 1146 |
| f0777 | SdlStateStore provides ExecutableDesign. | 1147 |
| f0778 | SdlUiBindingAdapter consumes DiagnosticPort. | 1148 |
| f0779 | SdlUiBindingAdapter consumes SdlExecutionPort. | 1149 |
| f0780 | SdlUiBindingAdapter consumes UiSessionPort. | 1150 |
| f0781 | SdlUiBindingAdapter owns ConnectTypedWidgetHandles. | 1151 |
| f0782 | SdlUiBindingAdapter owns DisconnectBindings. | 1152 |
| f0783 | SdlUiBindingAdapter owns PublishDomainUpdates. | 1153 |
| f0784 | SdlUiBindingAdapter owns ResolveCallbackSymbols. | 1154 |
| f0785 | SdlUiBindingAdapter owns RouteDomainBindings. | 1155 |
| f0786 | SdlUiBindingAdapter provides BoundInteraction. | 1156 |
| f0787 | SdlUiBindingAdapter uses SdlActionCalls as receiver of SdlActionResult in mode BoundExecution. | 1157 |
| f0788 | SdlUiBindingAdapter uses SdlActionCalls as sender of SdlActionRequest in mode BoundExecution. | 1158 |
| f0789 | SdlUiBindingAdapter uses UiDomainActions as receiver of BoundActionRequest in mode BoundExecution. | 1159 |
| f0790 | SdlUiBindingAdapter uses UiDomainActions as sender of BoundActionResult in mode BoundExecution. | 1160 |
| f0791 | SdlValidator owns ResolveSdlSymbols. | 1161 |
| f0792 | SdlValidator owns ValidateSdlProfile. | 1162 |
| f0793 | SdlValidator owns ValidateSdlStructure. | 1163 |
| f0794 | SdlValidator provides SdlSourceModel. | 1164 |
| f0795 | SdlViewpointGenerator owns ExportViewpointMarkdown. | 1165 |
| f0796 | SdlViewpointGenerator owns ProjectSdlViewpoints. | 1166 |
| f0797 | SdlViewpointGenerator owns TraceViewpointFacts. | 1167 |
| f0798 | SduiDispatcher consumes DomainBindingPort. | 1168 |
| f0799 | SduiDispatcher consumes UiStatePort. | 1169 |
| f0800 | SduiDispatcher owns CorrelateUiResult. | 1170 |
| f0801 | SduiDispatcher owns DispatchUiEvent. | 1171 |
| f0802 | SduiDispatcher owns RejectStaleUiEvent. | 1172 |
| f0803 | SduiDispatcher owns ValidateUiEvent. | 1173 |
| f0804 | SduiDispatcher provides InteractiveSession. | 1174 |
| f0805 | SduiDispatcher uses NativeUiActions as receiver of UiActionRequest in mode BoundExecution. | 1175 |
| f0806 | SduiDispatcher uses NativeUiActions as sender of UiActionRejected in mode BoundExecution. | 1176 |
| f0807 | SduiDispatcher uses NativeUiActions as sender of UiActionResult in mode BoundExecution. | 1177 |
| f0808 | SduiDispatcher uses UiDomainActions as receiver of BoundActionResult in mode BoundExecution. | 1178 |
| f0809 | SduiDispatcher uses UiDomainActions as sender of BoundActionRequest in mode BoundExecution. | 1179 |
| f0810 | SduiFrontend consumes DiagnosticPort. | 1180 |
| f0811 | SduiFrontend consumes SourceSnapshotPort. | 1181 |
| f0812 | SduiFrontend contains SduiLexer. | 1182 |
| f0813 | SduiFrontend contains SduiNormalizer. | 1183 |
| f0814 | SduiFrontend contains SduiParser. | 1184 |
| f0815 | SduiFrontend contains SduiValidator. | 1185 |
| f0816 | SduiFrontend owns CoordinateSduiCompilation. | 1186 |
| f0817 | SduiFrontend provides SduiSourceModel. | 1187 |
| f0818 | SduiFrontend uses UiAstCalls as receiver of BuildUiAstResult in mode SourceInspection. | 1188 |
| f0819 | SduiFrontend uses UiAstCalls as sender of BuildUiAstRequest in mode SourceInspection. | 1189 |
| f0820 | SduiFrontend uses UiCompilationCalls as receiver of CompileUiRequest in mode LiveEditing. | 1190 |
| f0821 | SduiFrontend uses UiCompilationCalls as receiver of CompileUiRequest in mode SourceInspection. | 1191 |
| f0822 | SduiFrontend uses UiCompilationCalls as sender of CompileUiRejected in mode LiveEditing. | 1192 |
| f0823 | SduiFrontend uses UiCompilationCalls as sender of CompileUiRejected in mode SourceInspection. | 1193 |
| f0824 | SduiFrontend uses UiCompilationCalls as sender of CompileUiResult in mode LiveEditing. | 1194 |
| f0825 | SduiFrontend uses UiCompilationCalls as sender of CompileUiResult in mode SourceInspection. | 1195 |
| f0826 | SduiFrontend uses UiNormalizationCalls as receiver of NormalizeUiResult in mode SourceInspection. | 1196 |
| f0827 | SduiFrontend uses UiNormalizationCalls as sender of NormalizeUiRequest in mode SourceInspection. | 1197 |
| f0828 | SduiFrontend uses UiTokenizationCalls as receiver of TokenizeUiResult in mode SourceInspection. | 1198 |
| f0829 | SduiFrontend uses UiTokenizationCalls as sender of TokenizeUiRequest in mode SourceInspection. | 1199 |
| f0830 | SduiFrontend uses UiValidationCalls as receiver of ValidateUiResult in mode SourceInspection. | 1200 |
| f0831 | SduiFrontend uses UiValidationCalls as sender of ValidateUiRequest in mode SourceInspection. | 1201 |
| f0832 | SduiInstanceStore holds UiSessionState. | 1202 |
| f0833 | SduiInstanceStore owns ManageWidgetIdentities. | 1203 |
| f0834 | SduiInstanceStore owns ProjectUiGeneration. | 1204 |
| f0835 | SduiInstanceStore owns RevokeWidgetGenerations. | 1205 |
| f0836 | SduiInstanceStore owns SnapshotUiState. | 1206 |
| f0837 | SduiInstanceStore provides InteractiveSession. | 1207 |
| f0838 | SduiInstanceStore uses UiGenerationEvents as sender of UiGenerationNotices in mode BoundExecution. | 1208 |
| f0839 | SduiInstanceStore uses UiGenerationEvents as sender of UiGenerationNotices in mode LiveEditing. | 1209 |
| f0840 | SduiLayout consumes ContentProviderPort. | 1210 |
| f0841 | SduiLayout consumes MeasurementPort. | 1211 |
| f0842 | SduiLayout consumes UiSnapshotPort. | 1212 |
| f0843 | SduiLayout owns AllocateGeometry. | 1213 |
| f0844 | SduiLayout owns BuildPreparedFrame. | 1214 |
| f0845 | SduiLayout owns ComputeClipping. | 1215 |
| f0846 | SduiLayout owns MeasureUiContent. | 1216 |
| f0847 | SduiLayout owns ResolveAncestorDimensions. | 1217 |
| f0848 | SduiLayout provides MeasuredPresentation. | 1218 |
| f0849 | SduiLayout uses LayoutCalls as receiver of LayoutRequest in mode StaticExport. | 1219 |
| f0850 | SduiLayout uses LayoutCalls as receiver of LayoutRequest in mode UiPreview. | 1220 |
| f0851 | SduiLayout uses LayoutCalls as sender of LayoutResult in mode StaticExport. | 1221 |
| f0852 | SduiLayout uses LayoutCalls as sender of LayoutResult in mode UiPreview. | 1222 |
| f0853 | SduiLexer owns TokenizeSduiSource. | 1223 |
| f0854 | SduiLexer provides SduiSourceModel. | 1224 |
| f0855 | SduiLexer uses UiTokenizationCalls as receiver of TokenizeUiRequest in mode SourceInspection. | 1225 |
| f0856 | SduiLexer uses UiTokenizationCalls as sender of TokenizeUiResult in mode SourceInspection. | 1226 |
| f0857 | SduiLibrary contains SduiFrontend. | 1227 |
| f0858 | SduiLibrary contains SduiLayout. | 1228 |
| f0859 | SduiLibrary contains SduiPresentation. | 1229 |
| f0860 | SduiLibrary contains SduiRuntime. | 1230 |
| f0861 | SduiLibrary provides InteractiveSession. | 1231 |
| f0862 | SduiLibrary provides MeasuredPresentation. | 1232 |
| f0863 | SduiLibrary provides SduiSourceModel. | 1233 |
| f0864 | SduiLibrary provides StaticDocumentation. | 1234 |
| f0865 | SduiNormalizer owns ExpandUiDefinitions. | 1235 |
| f0866 | SduiNormalizer owns PreserveUiRegions. | 1236 |
| f0867 | SduiNormalizer owns PreserveUiSourceMap. | 1237 |
| f0868 | SduiNormalizer provides SduiSourceModel. | 1238 |
| f0869 | SduiNormalizer uses UiNormalizationCalls as receiver of NormalizeUiRequest in mode SourceInspection. | 1239 |
| f0870 | SduiNormalizer uses UiNormalizationCalls as sender of NormalizeUiResult in mode SourceInspection. | 1240 |
| f0871 | SduiParser owns BuildSduiAst. | 1241 |
| f0872 | SduiParser provides SduiSourceModel. | 1242 |
| f0873 | SduiParser uses UiAstCalls as receiver of BuildUiAstRequest in mode SourceInspection. | 1243 |
| f0874 | SduiParser uses UiAstCalls as sender of BuildUiAstResult in mode SourceInspection. | 1244 |
| f0875 | SduiPresentation consumes PreparedFramePort. | 1245 |
| f0876 | SduiPresentation owns ComposeMarkdownDocument. | 1246 |
| f0877 | SduiPresentation owns ExportConsoleSnapshot. | 1247 |
| f0878 | SduiPresentation owns ExportSvgSnapshot. | 1248 |
| f0879 | SduiPresentation provides StaticDocumentation. | 1249 |
| f0880 | SduiPresentation uses SvgExportCalls as receiver of ExportSvgRequest in mode StaticExport. | 1250 |
| f0881 | SduiPresentation uses SvgExportCalls as sender of ExportSvgResult in mode StaticExport. | 1251 |
| f0882 | SduiPropertyStore consumes UiStatePort. | 1252 |
| f0883 | SduiPropertyStore owns ApplyPropertyBatch. | 1253 |
| f0884 | SduiPropertyStore owns TrackInputDraft. | 1254 |
| f0885 | SduiPropertyStore owns ValidatePropertyBatch. | 1255 |
| f0886 | SduiPropertyStore provides InteractiveSession. | 1256 |
| f0887 | SduiRuntime consumes DomainBindingPort. | 1257 |
| f0888 | SduiRuntime consumes SduiModelPort. | 1258 |
| f0889 | SduiRuntime contains SduiDispatcher. | 1259 |
| f0890 | SduiRuntime contains SduiInstanceStore. | 1260 |
| f0891 | SduiRuntime contains SduiPropertyStore. | 1261 |
| f0892 | SduiRuntime contains UiStateReconciler. | 1262 |
| f0893 | SduiRuntime owns CloseUiInstance. | 1263 |
| f0894 | SduiRuntime owns CreateUiInstance. | 1264 |
| f0895 | SduiRuntime provides InteractiveSession. | 1265 |
| f0896 | SduiSourceModel requires SourceSnapshotPort in mode SourceInspection. | 1266 |
| f0897 | SduiValidator owns ResolveUiNames. | 1267 |
| f0898 | SduiValidator owns ValidateRelativeFormatting. | 1268 |
| f0899 | SduiValidator owns ValidateSymbolicBindings. | 1269 |
| f0900 | SduiValidator owns ValidateWidgetArguments. | 1270 |
| f0901 | SduiValidator provides SduiSourceModel. | 1271 |
| f0902 | SduiValidator uses UiValidationCalls as receiver of ValidateUiRequest in mode SourceInspection. | 1272 |
| f0903 | SduiValidator uses UiValidationCalls as sender of ValidateUiResult in mode SourceInspection. | 1273 |
| f0904 | SessionDraft has presence = optional. | 1274 |
| f0905 | SessionDraft has value-type = text. | 1275 |
| f0906 | SessionGeneration has presence = required. | 1276 |
| f0907 | SessionGeneration has value-type = unsigned. | 1277 |
| f0908 | SnapshotDomainState realizes ExecutableDesign. | 1278 |
| f0909 | SnapshotUiState realizes InteractiveSession. | 1279 |
| f0910 | SourceDocumentRevision has presence = required. | 1280 |
| f0911 | SourceDocumentRevision has value-type = unsigned. | 1281 |
| f0912 | SourceDocumentText has presence = required. | 1282 |
| f0913 | SourceDocumentText has value-type = text. | 1283 |
| f0914 | SourceLoader consumes SourceInputPort. | 1284 |
| f0915 | SourceLoader owns DesignSourceArchive. | 1285 |
| f0916 | SourceLoader owns IdentifySourceRevision. | 1286 |
| f0917 | SourceLoader owns ReadBoundedSources. | 1287 |
| f0918 | SourceLoader provides SourceLoading. | 1288 |
| f0919 | SourceWatcher consumes FileChangePort. | 1289 |
| f0920 | SourceWatcher owns CoalesceSourceChanges. | 1290 |
| f0921 | SourceWatcher owns ObserveSourceChanges. | 1291 |
| f0922 | SourceWatcher provides DevelopmentReload. | 1292 |
| f0923 | SourceWatcher uses ModelReloadCalls as receiver of ReloadPublished in mode BoundLiveEditing. | 1293 |
| f0924 | SourceWatcher uses ModelReloadCalls as receiver of ReloadPublished in mode LiveEditing. | 1294 |
| f0925 | SourceWatcher uses ModelReloadCalls as receiver of ReloadRejected in mode BoundLiveEditing. | 1295 |
| f0926 | SourceWatcher uses ModelReloadCalls as receiver of ReloadRejected in mode LiveEditing. | 1296 |
| f0927 | SourceWatcher uses ModelReloadCalls as sender of ReloadRequest in mode BoundLiveEditing. | 1297 |
| f0928 | SourceWatcher uses ModelReloadCalls as sender of ReloadRequest in mode LiveEditing. | 1298 |
| f0929 | StaticDocumentation requires ExportSinkPort in mode StaticExport. | 1299 |
| f0930 | StaticDocumentation requires PreparedFramePort in mode StaticExport. | 1300 |
| f0931 | StaticFrameExported exercises PublishDesignDocumentation. | 1301 |
| f0932 | StaticFrameExported has completeness = closed. | 1302 |
| f0933 | StaticFrameExported illustrates G2M2SharedSvgGeometry. | 1303 |
| f0934 | StaticFrameExported illustrates G5M3DocumentationExport. | 1304 |
| f0935 | StaticFrameExported runs-in StaticExport. | 1305 |
| f0936 | StaticFrameExported step 1 sends LayoutRequest from CommandLineHost to SduiLayout via LayoutCalls. | 1306 |
| f0937 | StaticFrameExported step 2 sends LayoutResult from SduiLayout to CommandLineHost via LayoutCalls reply-to 1. | 1307 |
| f0938 | StaticFrameExported step 3 sends ExportSvgRequest from CommandLineHost to SduiPresentation via SvgExportCalls. | 1308 |
| f0939 | StaticFrameExported step 4 sends ExportSvgResult from SduiPresentation to CommandLineHost via SvgExportCalls reply-to 3. | 1309 |
| f0940 | StructuralModelInspection supports EditRunningPrototype. | 1310 |
| f0941 | StructuralModelInspection supports InspectModels. | 1311 |
| f0942 | SvgDocumentContract has completeness = closed. | 1312 |
| f0943 | SvgDocumentContract has-field SvgDocumentText. | 1313 |
| f0944 | SvgDocumentText has presence = required. | 1314 |
| f0945 | SvgDocumentText has value-type = text. | 1315 |
| f0946 | SvgExportCalls upholds SvgExportCallsProtocol. | 1316 |
| f0947 | SvgExportCallsProtocol has completeness = closed. | 1317 |
| f0948 | SvgExportCallsProtocol permits ExportSvgRequest. | 1318 |
| f0949 | SvgExportCallsProtocol permits ExportSvgResult. | 1319 |
| f0950 | TokenArtifact has presence = required. | 1320 |
| f0951 | TokenArtifact has value-type = bytes. | 1321 |
| f0952 | TokenArtifactContract has completeness = closed. | 1322 |
| f0953 | TokenArtifactContract has-field TokenArtifact. | 1323 |
| f0954 | TokenizeSdlSource has repeatability = deterministic. | 1324 |
| f0955 | TokenizeSdlSource has state-retention = stateless. | 1325 |
| f0956 | TokenizeSdlSource realizes SdlSourceModel. | 1326 |
| f0957 | TokenizeSduiSource has repeatability = deterministic. | 1327 |
| f0958 | TokenizeSduiSource has state-retention = stateless. | 1328 |
| f0959 | TokenizeSduiSource realizes SduiSourceModel. | 1329 |
| f0960 | TokenizeUiRequest has message-kind = request. | 1330 |
| f0961 | TokenizeUiRequest upholds ReloadArguments. | 1331 |
| f0962 | TokenizeUiResult has message-kind = result. | 1332 |
| f0963 | TokenizeUiResult replies-to TokenizeUiRequest. | 1333 |
| f0964 | TokenizeUiResult upholds TokenArtifactContract. | 1334 |
| f0965 | TraceViewpointFacts allocated-to CommandLineHost in mode SourceInspection. | 1335 |
| f0966 | TraceViewpointFacts allocated-to CommandLineHost in mode StaticExport. | 1336 |
| f0967 | TraceViewpointFacts contributes-to DesignDocumentation. | 1337 |
| f0968 | TraceViewpointFacts contributes-to InspectModels. | 1338 |
| f0969 | TraceViewpointFacts has repeatability = deterministic. | 1339 |
| f0970 | TraceViewpointFacts has state-retention = stateless. | 1340 |
| f0971 | TrackInputDraft has state-retention = stateful. | 1341 |
| f0972 | TrackInputDraft realizes InteractiveSession. | 1342 |
| f0973 | TypedDomainBinding supports TryDomainInteraction. | 1343 |
| f0974 | UiActionRejected has message-kind = result. | 1344 |
| f0975 | UiActionRejected replies-to UiActionRequest. | 1345 |
| f0976 | UiActionRejected upholds ActionOutcome. | 1346 |
| f0977 | UiActionRequest has message-kind = request. | 1347 |
| f0978 | UiActionRequest upholds ActionArguments. | 1348 |
| f0979 | UiActionResult has message-kind = result. | 1349 |
| f0980 | UiActionResult replies-to UiActionRequest. | 1350 |
| f0981 | UiActionResult upholds ActionOutcome. | 1351 |
| f0982 | UiAstCalls upholds UiAstCallsProtocol. | 1352 |
| f0983 | UiAstCallsProtocol has completeness = closed. | 1353 |
| f0984 | UiAstCallsProtocol permits BuildUiAstRequest. | 1354 |
| f0985 | UiAstCallsProtocol permits BuildUiAstResult. | 1355 |
| f0986 | UiCompilationAccepted exercises InspectModels. | 1356 |
| f0987 | UiCompilationAccepted has completeness = closed. | 1357 |
| f0988 | UiCompilationAccepted illustrates G1M1ParserAndAst. | 1358 |
| f0989 | UiCompilationAccepted illustrates G1M2ValidationAndNormalization. | 1359 |
| f0990 | UiCompilationAccepted runs-in SourceInspection. | 1360 |
| f0991 | UiCompilationAccepted step 1 sends CompileUiRequest from CommandLineHost to SduiFrontend via UiCompilationCalls. | 1361 |
| f0992 | UiCompilationAccepted step 10 sends CompileUiResult from SduiFrontend to CommandLineHost via UiCompilationCalls reply-to 1. | 1362 |
| f0993 | UiCompilationAccepted step 2 sends TokenizeUiRequest from SduiFrontend to SduiLexer via UiTokenizationCalls. | 1363 |
| f0994 | UiCompilationAccepted step 3 sends TokenizeUiResult from SduiLexer to SduiFrontend via UiTokenizationCalls reply-to 2. | 1364 |
| f0995 | UiCompilationAccepted step 4 sends BuildUiAstRequest from SduiFrontend to SduiParser via UiAstCalls. | 1365 |
| f0996 | UiCompilationAccepted step 5 sends BuildUiAstResult from SduiParser to SduiFrontend via UiAstCalls reply-to 4. | 1366 |
| f0997 | UiCompilationAccepted step 6 sends ValidateUiRequest from SduiFrontend to SduiValidator via UiValidationCalls. | 1367 |
| f0998 | UiCompilationAccepted step 7 sends ValidateUiResult from SduiValidator to SduiFrontend via UiValidationCalls reply-to 6. | 1368 |
| f0999 | UiCompilationAccepted step 8 sends NormalizeUiRequest from SduiFrontend to SduiNormalizer via UiNormalizationCalls. | 1369 |
| f1000 | UiCompilationAccepted step 9 sends NormalizeUiResult from SduiNormalizer to SduiFrontend via UiNormalizationCalls reply-to 8. | 1370 |
| f1001 | UiCompilationCalls upholds UiCompilationCallsProtocol. | 1371 |
| f1002 | UiCompilationCallsProtocol has completeness = closed. | 1372 |
| f1003 | UiCompilationCallsProtocol permits CompileUiRejected. | 1373 |
| f1004 | UiCompilationCallsProtocol permits CompileUiRequest. | 1374 |
| f1005 | UiCompilationCallsProtocol permits CompileUiResult. | 1375 |
| f1006 | UiDomainActions upholds UiDomainActionsProtocol. | 1376 |
| f1007 | UiDomainActionsProtocol has completeness = closed. | 1377 |
| f1008 | UiDomainActionsProtocol permits BoundActionRequest. | 1378 |
| f1009 | UiDomainActionsProtocol permits BoundActionResult. | 1379 |
| f1010 | UiGenerationChanged has-field NoticeGeneration. | 1380 |
| f1011 | UiGenerationContract defines UiGenerationChanged. | 1381 |
| f1012 | UiGenerationContract has completeness = closed. | 1382 |
| f1013 | UiGenerationContract has-field NoticeVersion. | 1383 |
| f1014 | UiGenerationEvents upholds UiGenerationEventsProtocol. | 1384 |
| f1015 | UiGenerationEventsProtocol has completeness = closed. | 1385 |
| f1016 | UiGenerationEventsProtocol permits UiGenerationNotices. | 1386 |
| f1017 | UiGenerationNotices from UiSessionState. | 1387 |
| f1018 | UiGenerationNotices upholds UiGenerationContract. | 1388 |
| f1019 | UiGenerationWire encodes UiGenerationChanged. | 1389 |
| f1020 | UiGenerationWire has bit-order = most-significant-first. | 1390 |
| f1021 | UiGenerationWire has byte-order = big-endian. | 1391 |
| f1022 | UiGenerationWire places NoticeGeneration at 16 bits 64. | 1392 |
| f1023 | UiGenerationWire places NoticeVersion at 0 bits 16. | 1393 |
| f1024 | UiModelReloadAccepted exercises EditRunningPrototype. | 1394 |
| f1025 | UiModelReloadAccepted has completeness = closed. | 1395 |
| f1026 | UiModelReloadAccepted illustrates G3M2CandidatePublication. | 1396 |
| f1027 | UiModelReloadAccepted illustrates G3M3CompatibleState. | 1397 |
| f1028 | UiModelReloadAccepted runs-in LiveEditing. | 1398 |
| f1029 | UiModelReloadAccepted step 1 sends ReloadRequest from SourceWatcher to ReloadCoordinator via ModelReloadCalls. | 1399 |
| f1030 | UiModelReloadAccepted step 2 sends CompileUiRequest from ReloadCoordinator to SduiFrontend via UiCompilationCalls. | 1400 |
| f1031 | UiModelReloadAccepted step 3 sends CompileUiResult from SduiFrontend to ReloadCoordinator via UiCompilationCalls reply-to 2. | 1401 |
| f1032 | UiModelReloadAccepted step 4 sends ReloadPublished from ReloadCoordinator to SourceWatcher via ModelReloadCalls reply-to 1. | 1402 |
| f1033 | UiModelReloadAccepted step 5 sends UiGenerationNotices variant UiGenerationChanged from SduiInstanceStore to FyneBackend via UiGenerationEvents. | 1403 |
| f1034 | UiModelReloadRejected exercises EditRunningPrototype. | 1404 |
| f1035 | UiModelReloadRejected has completeness = closed. | 1405 |
| f1036 | UiModelReloadRejected illustrates G3M2CandidatePublication. | 1406 |
| f1037 | UiModelReloadRejected runs-in LiveEditing. | 1407 |
| f1038 | UiModelReloadRejected step 1 sends ReloadRequest from SourceWatcher to ReloadCoordinator via ModelReloadCalls. | 1408 |
| f1039 | UiModelReloadRejected step 2 sends CompileUiRequest from ReloadCoordinator to SduiFrontend via UiCompilationCalls. | 1409 |
| f1040 | UiModelReloadRejected step 3 sends CompileUiRejected from SduiFrontend to ReloadCoordinator via UiCompilationCalls reply-to 2. | 1410 |
| f1041 | UiModelReloadRejected step 4 sends ReloadRejected from ReloadCoordinator to SourceWatcher via ModelReloadCalls reply-to 1. | 1411 |
| f1042 | UiNormalizationCalls upholds UiNormalizationCallsProtocol. | 1412 |
| f1043 | UiNormalizationCallsProtocol has completeness = closed. | 1413 |
| f1044 | UiNormalizationCallsProtocol permits NormalizeUiRequest. | 1414 |
| f1045 | UiNormalizationCallsProtocol permits NormalizeUiResult. | 1415 |
| f1046 | UiSessionRecord has completeness = closed. | 1416 |
| f1047 | UiSessionRecord has-field SessionDraft. | 1417 |
| f1048 | UiSessionRecord has-field SessionGeneration. | 1418 |
| f1049 | UiSessionState upholds UiSessionRecord. | 1419 |
| f1050 | UiStateReconciler consumes SduiModelPort. | 1420 |
| f1051 | UiStateReconciler consumes UiStatePort. | 1421 |
| f1052 | UiStateReconciler owns MatchCompatibleWidgets. | 1422 |
| f1053 | UiStateReconciler owns PreserveCompatibleUiState. | 1423 |
| f1054 | UiStateReconciler owns ResetIncompatibleUiState. | 1424 |
| f1055 | UiStateReconciler provides InteractiveSession. | 1425 |
| f1056 | UiTokenizationCalls upholds UiTokenizationCallsProtocol. | 1426 |
| f1057 | UiTokenizationCallsProtocol has completeness = closed. | 1427 |
| f1058 | UiTokenizationCallsProtocol permits TokenizeUiRequest. | 1428 |
| f1059 | UiTokenizationCallsProtocol permits TokenizeUiResult. | 1429 |
| f1060 | UiValidationCalls upholds UiValidationCallsProtocol. | 1430 |
| f1061 | UiValidationCallsProtocol has completeness = closed. | 1431 |
| f1062 | UiValidationCallsProtocol permits ValidateUiRequest. | 1432 |
| f1063 | UiValidationCallsProtocol permits ValidateUiResult. | 1433 |
| f1064 | UnboundLocalAction exercises PrototypeUserInterface. | 1434 |
| f1065 | UnboundLocalAction has completeness = closed. | 1435 |
| f1066 | UnboundLocalAction illustrates G2M3FyneInteractions. | 1436 |
| f1067 | UnboundLocalAction runs-in UiPreview. | 1437 |
| f1068 | UnboundLocalAction step 1 sends DomainActionRequest from FyneHost to GoDomainImplementation via GoDomainCalls. | 1438 |
| f1069 | UnboundLocalAction step 2 sends DomainActionResult from GoDomainImplementation to FyneHost via GoDomainCalls reply-to 1. | 1439 |
| f1070 | ValidateActionInput realizes ExecutableDesign. | 1440 |
| f1071 | ValidatePropertyBatch realizes InteractiveSession. | 1441 |
| f1072 | ValidateRelativeFormatting has repeatability = deterministic. | 1442 |
| f1073 | ValidateRelativeFormatting has state-retention = stateless. | 1443 |
| f1074 | ValidateRelativeFormatting realizes SduiSourceModel. | 1444 |
| f1075 | ValidateSdlProfile has repeatability = deterministic. | 1445 |
| f1076 | ValidateSdlProfile has state-retention = stateless. | 1446 |
| f1077 | ValidateSdlProfile realizes SdlSourceModel. | 1447 |
| f1078 | ValidateSdlStructure allocated-to CommandLineHost in mode SourceInspection. | 1448 |
| f1079 | ValidateSdlStructure contributes-to StructuralModelInspection. | 1449 |
| f1080 | ValidateSdlStructure has repeatability = deterministic. | 1450 |
| f1081 | ValidateSdlStructure has state-retention = stateless. | 1451 |
| f1082 | ValidateSdlStructure realizes SdlSourceModel. | 1452 |
| f1083 | ValidateSymbolicBindings has repeatability = deterministic. | 1453 |
| f1084 | ValidateSymbolicBindings has state-retention = stateless. | 1454 |
| f1085 | ValidateSymbolicBindings realizes SduiSourceModel. | 1455 |
| f1086 | ValidateUiEvent realizes InteractiveSession. | 1456 |
| f1087 | ValidateUiRequest has message-kind = request. | 1457 |
| f1088 | ValidateUiRequest upholds AstArtifactContract. | 1458 |
| f1089 | ValidateUiResult has message-kind = result. | 1459 |
| f1090 | ValidateUiResult replies-to ValidateUiRequest. | 1460 |
| f1091 | ValidateUiResult upholds ValidationOutcomeContract. | 1461 |
| f1092 | ValidateVisualResources realizes RichContent. | 1462 |
| f1093 | ValidateWidgetArguments allocated-to CommandLineHost in mode SourceInspection. | 1463 |
| f1094 | ValidateWidgetArguments contributes-to StructuralModelInspection. | 1464 |
| f1095 | ValidateWidgetArguments has repeatability = deterministic. | 1465 |
| f1096 | ValidateWidgetArguments has state-retention = stateless. | 1466 |
| f1097 | ValidateWidgetArguments realizes SduiSourceModel. | 1467 |
| f1098 | ValidationOutcomeContract has completeness = closed. | 1468 |
| f1099 | ValidationOutcomeContract has-field ModelDiagnostics. | 1469 |
| f1100 | ValidationOutcomeContract has-field ModelIsValid. | 1470 |
| f1101 | VerifyNativeBehaviorParity allocated-to CommandLineHost in mode NativeBuild. | 1471 |
| f1102 | VerifyNativeBehaviorParity contributes-to NativeGoAssembly. | 1472 |
| f1103 | WriteGeneratedArtifacts allocated-to CommandLineHost in mode StaticExport. | 1473 |
| f1104 | WriteGeneratedArtifacts contributes-to DesignDocumentation. | 1474 |
| f1105 | WriteGeneratedArtifacts realizes StaticDocumentation. | 1475 |

### Deklarasjonsregister

| Identitet | Type | Kildelinje |
| --- | --- | --- |
| ActionArguments | contract | 2 |
| ActionGeneration | field | 3 |
| ActionInputText | field | 4 |
| ActionOutcome | contract | 5 |
| ActionOutputText | field | 6 |
| ActionStatusCode | field | 7 |
| ActionSymbol | field | 8 |
| AllocateGeometry | functionality | 9 |
| ApplyPropertyBatch | functionality | 10 |
| AstArtifact | field | 11 |
| AstArtifactContract | contract | 12 |
| BindingReloadPort | interface | 13 |
| BoundActionAccepted | scenario | 14 |
| BoundActionRejected | scenario | 15 |
| BoundActionRequest | message | 16 |
| BoundActionResult | message | 17 |
| BoundExecution | mode | 18 |
| BoundInteraction | capability | 19 |
| BoundLiveEditing | mode | 20 |
| BuildGeneratedApplication | functionality | 21 |
| BuildGoRequest | message | 22 |
| BuildGoResult | message | 23 |
| BuildNativeProduct | usecase | 24 |
| BuildNativeRealization | activity | 25 |
| BuildPreparedFrame | functionality | 26 |
| BuildSdlAst | functionality | 27 |
| BuildSduiAst | functionality | 28 |
| BuildToolPort | interface | 29 |
| BuildUiAstRequest | message | 30 |
| BuildUiAstResult | message | 31 |
| CancelPendingActions | functionality | 32 |
| CheckDomainStateCompatibility | functionality | 33 |
| CheckExecutionCompleteness | functionality | 34 |
| CheckFunctionSignatures | functionality | 35 |
| CloseSdlInstance | functionality | 36 |
| CloseUiInstance | functionality | 37 |
| CoalesceSourceChanges | functionality | 38 |
| CommandLineHost | container | 39 |
| CompileSdlRequest | message | 40 |
| CompileSdlResult | message | 41 |
| CompileUiRejected | message | 42 |
| CompileUiRequest | message | 43 |
| CompileUiResult | message | 44 |
| ComposeHeadlessExport | functionality | 45 |
| ComposeInteractiveSession | functionality | 46 |
| ComposeMarkdownDocument | functionality | 47 |
| ComputeClipping | functionality | 48 |
| ConnectTypedWidgetHandles | functionality | 49 |
| ContentProviderPort | interface | 50 |
| ContentServices | unit | 51 |
| CoordinateSdlCompilation | functionality | 52 |
| CoordinateSduiCompilation | functionality | 53 |
| CorrelateActionResult | functionality | 54 |
| CorrelateUiResult | functionality | 55 |
| CreateSdlInstance | functionality | 56 |
| CreateUiInstance | functionality | 57 |
| DesignAuthor | actor | 58 |
| DesignDocumentation | feature | 59 |
| DesignReviewer | actor | 60 |
| DesignSourceArchive | database | 61 |
| DesignSourceDocuments | dataset | 62 |
| DesignSourceRecord | contract | 63 |
| DevelopmentReload | capability | 64 |
| DevelopmentTools | unit | 65 |
| DiagnosticPort | interface | 66 |
| DiagnosticReporter | unit | 67 |
| DiagramEnginePort | interface | 68 |
| DiagramPort | interface | 69 |
| DiagramProvider | unit | 70 |
| DisconnectBindings | functionality | 71 |
| DispatchUiEvent | functionality | 72 |
| DomainActionRequest | message | 73 |
| DomainActionResult | message | 74 |
| DomainBindingPort | interface | 75 |
| DomainFunctionPort | interface | 76 |
| DomainOperations | capability | 77 |
| DomainStateMigrator | unit | 78 |
| DomainStatePort | interface | 79 |
| EditRunningPrototype | usecase | 80 |
| ExecutableDesign | capability | 81 |
| ExecutionProfilePort | interface | 82 |
| ExpandUiDefinitions | functionality | 83 |
| ExportConsoleSnapshot | functionality | 84 |
| ExportSinkPort | interface | 85 |
| ExportSvgRequest | message | 86 |
| ExportSvgResult | message | 87 |
| ExportSvgSnapshot | functionality | 88 |
| ExportUiDocumentation | activity | 89 |
| ExportViewpointMarkdown | functionality | 90 |
| FileChangePort | interface | 91 |
| FramePresentationCalls | channel | 92 |
| FramePresentationCallsProtocol | contract | 93 |
| FyneBackend | unit | 94 |
| FyneHost | container | 95 |
| G1FrontendPort | activity | 96 |
| G1M1ParserAndAst | activity | 97 |
| G1M2ValidationAndNormalization | activity | 98 |
| G1M3Concept1AndDumps | activity | 99 |
| G2LayoutAndPresentation | activity | 100 |
| G2M1RelativeMeasurement | activity | 101 |
| G2M2SharedSvgGeometry | activity | 102 |
| G2M3FyneInteractions | activity | 103 |
| G2M4RichContent | activity | 104 |
| G3M1TypedUiSession | activity | 105 |
| G3M2CandidatePublication | activity | 106 |
| G3M3CompatibleState | activity | 107 |
| G3UiRuntimeAndReload | activity | 108 |
| G4M1SdlFrontend | activity | 109 |
| G4M2TypedExecution | activity | 110 |
| G4M3UiDomainBinding | activity | 111 |
| G4M4DomainReload | activity | 112 |
| G4SdlRuntimeAndBinding | activity | 113 |
| G5M1GeneratedGo | activity | 114 |
| G5M2BehaviorParity | activity | 115 |
| G5M3DocumentationExport | activity | 116 |
| G5M4RetirePython | activity | 117 |
| G5NativeGeneration | activity | 118 |
| GenerateBindingRegistration | functionality | 119 |
| GenerateGoRequest | message | 120 |
| GenerateGoResult | message | 121 |
| GenerateModelConstructors | functionality | 122 |
| GeneratedArtifactPort | interface | 123 |
| GeneratedGoContract | contract | 124 |
| GeneratedGoSources | field | 125 |
| GoBuildCalls | channel | 126 |
| GoBuildCallsProtocol | contract | 127 |
| GoBuildRunner | unit | 128 |
| GoCodeGenerator | unit | 129 |
| GoDomainCalls | channel | 130 |
| GoDomainCallsProtocol | contract | 131 |
| GoDomainImplementation | unit | 132 |
| GoGenerationCalls | channel | 133 |
| GoGenerationCallsProtocol | contract | 134 |
| HandleFocusAndTextInput | functionality | 135 |
| IdentifySourceRevision | functionality | 136 |
| InspectDesignSource | activity | 137 |
| InspectModels | usecase | 138 |
| InspectSdlSource | activity | 139 |
| InspectSduiSource | activity | 140 |
| InteractiveFramePrepared | scenario | 141 |
| InteractiveSession | capability | 142 |
| InteractiveUiPreview | feature | 143 |
| InvokeRegisteredFunction | functionality | 144 |
| KeepLastValidModels | functionality | 145 |
| LayoutArguments | contract | 146 |
| LayoutCalls | channel | 147 |
| LayoutCallsProtocol | contract | 148 |
| LayoutModelArtifact | field | 149 |
| LayoutRequest | message | 150 |
| LayoutResult | message | 151 |
| LayoutViewportHeight | field | 152 |
| LayoutViewportWidth | field | 153 |
| LiveEditing | mode | 154 |
| LiveModelReload | feature | 155 |
| ManageDomainState | functionality | 156 |
| ManageWidgetIdentities | functionality | 157 |
| MarkdownProvider | unit | 158 |
| MatchCompatibleWidgets | functionality | 159 |
| MeasureMarkdownContent | functionality | 160 |
| MeasureUiContent | functionality | 161 |
| MeasuredPresentation | capability | 162 |
| MeasurementPort | interface | 163 |
| MigrateOrResetDomainState | functionality | 164 |
| ModelDiagnostics | field | 165 |
| ModelIsValid | field | 166 |
| ModelReloadCalls | channel | 167 |
| ModelReloadCallsProtocol | contract | 168 |
| NativeBuild | mode | 169 |
| NativeBuildContract | contract | 170 |
| NativeBuildDiagnostics | field | 171 |
| NativeBuildSucceeded | field | 172 |
| NativeGoAssembly | feature | 173 |
| NativeInteraction | capability | 174 |
| NativeProgramBuilt | scenario | 175 |
| NativeRealization | capability | 176 |
| NativeUiActions | channel | 177 |
| NativeUiActionsProtocol | contract | 178 |
| NormalizeSdlModel | functionality | 179 |
| NormalizeUiRequest | message | 180 |
| NormalizeUiResult | message | 181 |
| NormalizedModelArtifact | field | 182 |
| NormalizedModelContract | contract | 183 |
| NoticeGeneration | field | 184 |
| NoticeVersion | field | 185 |
| ObserveSourceChanges | functionality | 186 |
| PerformDomainOperation | functionality | 187 |
| PrepareCandidateModels | functionality | 188 |
| PrepareDiagramResource | functionality | 189 |
| PrepareMarkdown | functionality | 190 |
| PreparedFrameArtifact | field | 191 |
| PreparedFrameContract | contract | 192 |
| PreparedFramePort | interface | 193 |
| PresentFrameRequest | message | 194 |
| PresentFrameResult | message | 195 |
| PresentationOutcome | contract | 196 |
| PresentationReady | field | 197 |
| PreserveCompatibleUiState | functionality | 198 |
| PreserveHandwrittenSources | functionality | 199 |
| PreserveSdlSourceMap | functionality | 200 |
| PreserveUiRegions | functionality | 201 |
| PreserveUiSourceMap | functionality | 202 |
| ProjectSdlViewpoints | functionality | 203 |
| ProjectUiGeneration | functionality | 204 |
| PrototypeUserInterface | usecase | 205 |
| PublishDesignDocumentation | usecase | 206 |
| PublishDomainUpdates | functionality | 207 |
| PublishModelGeneration | functionality | 208 |
| PublishPresentation | functionality | 209 |
| ReadBoundedSources | functionality | 210 |
| RealizeDesign | activity | 211 |
| ReconcileWidgets | functionality | 212 |
| RegisterDomainFunctions | functionality | 213 |
| RejectStaleUiEvent | functionality | 214 |
| ReleaseNativeWidgets | functionality | 215 |
| ReleaseVisualResources | functionality | 216 |
| ReloadArguments | contract | 217 |
| ReloadBoundModels | activity | 218 |
| ReloadCoordinator | unit | 219 |
| ReloadDesignSession | activity | 220 |
| ReloadDiagnostic | field | 221 |
| ReloadOutcome | contract | 222 |
| ReloadPort | interface | 223 |
| ReloadPublished | message | 224 |
| ReloadPublishedGeneration | field | 225 |
| ReloadRejected | message | 226 |
| ReloadRequest | message | 227 |
| ReloadSourceRevision | field | 228 |
| ReloadSourceText | field | 229 |
| ReloadUiModel | activity | 230 |
| ReportBindingDiagnostics | functionality | 231 |
| ReportSourceDiagnostics | functionality | 232 |
| ResetIncompatibleUiState | functionality | 233 |
| ResolveAncestorDimensions | functionality | 234 |
| ResolveCallbackSymbols | functionality | 235 |
| ResolveSdlSymbols | functionality | 236 |
| ResolveUiNames | functionality | 237 |
| ResourcePort | interface | 238 |
| ResourceStore | unit | 239 |
| RestartChangedGoProgram | functionality | 240 |
| RetirePreviousGeneration | functionality | 241 |
| RetireReplacedPythonEntryPoints | functionality | 242 |
| RevokeWidgetGenerations | functionality | 243 |
| RichContent | capability | 244 |
| RichDocument | mode | 245 |
| RouteDomainBindings | functionality | 246 |
| RunBoundUiAction | activity | 247 |
| RunDesignSession | activity | 248 |
| RunUnboundUiPreview | activity | 249 |
| ScheduleUiPublication | functionality | 250 |
| SdlActionCalls | channel | 251 |
| SdlActionCallsProtocol | contract | 252 |
| SdlActionRequest | message | 253 |
| SdlActionResult | message | 254 |
| SdlCompilationCalls | channel | 255 |
| SdlCompilationCallsProtocol | contract | 256 |
| SdlDispatcher | unit | 257 |
| SdlExecutionGate | unit | 258 |
| SdlExecutionPort | interface | 259 |
| SdlFrontend | unit | 260 |
| SdlFrontendPort | interface | 261 |
| SdlFunctionRegistry | unit | 262 |
| SdlLexer | unit | 263 |
| SdlLibrary | unit | 264 |
| SdlModelPort | interface | 265 |
| SdlModelReloadAccepted | scenario | 266 |
| SdlNormalizer | unit | 267 |
| SdlParser | unit | 268 |
| SdlReloadPort | interface | 269 |
| SdlRuntime | unit | 270 |
| SdlSourceModel | capability | 271 |
| SdlStateStore | unit | 272 |
| SdlUiBindingAdapter | unit | 273 |
| SdlValidator | unit | 274 |
| SdlViewpointGenerator | unit | 275 |
| SduiDispatcher | unit | 276 |
| SduiFrontend | unit | 277 |
| SduiFrontendPort | interface | 278 |
| SduiInstanceStore | unit | 279 |
| SduiLayout | unit | 280 |
| SduiLexer | unit | 281 |
| SduiLibrary | unit | 282 |
| SduiModelPort | interface | 283 |
| SduiNormalizer | unit | 284 |
| SduiParser | unit | 285 |
| SduiPresentation | unit | 286 |
| SduiPropertyStore | unit | 287 |
| SduiRuntime | unit | 288 |
| SduiSourceModel | capability | 289 |
| SduiValidator | unit | 290 |
| SessionDraft | field | 291 |
| SessionGeneration | field | 292 |
| SnapshotDomainState | functionality | 293 |
| SnapshotUiState | functionality | 294 |
| SourceDiagnostics | capability | 295 |
| SourceDocumentRevision | field | 296 |
| SourceDocumentText | field | 297 |
| SourceInputPort | interface | 298 |
| SourceInspection | mode | 299 |
| SourceLoader | unit | 300 |
| SourceLoading | capability | 301 |
| SourceSnapshotPort | interface | 302 |
| SourceWatcher | unit | 303 |
| StaticDocumentation | capability | 304 |
| StaticExport | mode | 305 |
| StaticFrameExported | scenario | 306 |
| StructuralModelInspection | feature | 307 |
| SvgDocumentContract | contract | 308 |
| SvgDocumentText | field | 309 |
| SvgExportCalls | channel | 310 |
| SvgExportCallsProtocol | contract | 311 |
| TokenArtifact | field | 312 |
| TokenArtifactContract | contract | 313 |
| TokenizeSdlSource | functionality | 314 |
| TokenizeSduiSource | functionality | 315 |
| TokenizeUiRequest | message | 316 |
| TokenizeUiResult | message | 317 |
| TraceViewpointFacts | functionality | 318 |
| TrackInputDraft | functionality | 319 |
| TryDomainInteraction | usecase | 320 |
| TypedDomainBinding | feature | 321 |
| UiActionRejected | message | 322 |
| UiActionRequest | message | 323 |
| UiActionResult | message | 324 |
| UiAstCalls | channel | 325 |
| UiAstCallsProtocol | contract | 326 |
| UiCompilationAccepted | scenario | 327 |
| UiCompilationCalls | channel | 328 |
| UiCompilationCallsProtocol | contract | 329 |
| UiDomainActions | channel | 330 |
| UiDomainActionsProtocol | contract | 331 |
| UiGenerationChanged | variant | 332 |
| UiGenerationContract | contract | 333 |
| UiGenerationEvents | channel | 334 |
| UiGenerationEventsProtocol | contract | 335 |
| UiGenerationNotices | datagram | 336 |
| UiGenerationWire | encoding | 337 |
| UiModelReloadAccepted | scenario | 338 |
| UiModelReloadRejected | scenario | 339 |
| UiNormalizationCalls | channel | 340 |
| UiNormalizationCallsProtocol | contract | 341 |
| UiPreview | mode | 342 |
| UiReloadPort | interface | 343 |
| UiSessionPort | interface | 344 |
| UiSessionRecord | contract | 345 |
| UiSessionState | dataset | 346 |
| UiSnapshotPort | interface | 347 |
| UiStatePort | interface | 348 |
| UiStateReconciler | unit | 349 |
| UiTokenizationCalls | channel | 350 |
| UiTokenizationCallsProtocol | contract | 351 |
| UiValidationCalls | channel | 352 |
| UiValidationCallsProtocol | contract | 353 |
| UnboundLocalAction | scenario | 354 |
| ValidateActionInput | functionality | 355 |
| ValidatePropertyBatch | functionality | 356 |
| ValidateRelativeFormatting | functionality | 357 |
| ValidateSdlProfile | functionality | 358 |
| ValidateSdlStructure | functionality | 359 |
| ValidateSymbolicBindings | functionality | 360 |
| ValidateUiEvent | functionality | 361 |
| ValidateUiRequest | message | 362 |
| ValidateUiResult | message | 363 |
| ValidateVisualResources | functionality | 364 |
| ValidateWidgetArguments | functionality | 365 |
| ValidationOutcomeContract | contract | 366 |
| VerifyNativeBehaviorParity | functionality | 367 |
| WidgetBackendPort | interface | 368 |
| WriteGeneratedArtifacts | functionality | 369 |
