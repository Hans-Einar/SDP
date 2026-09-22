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

### Bruksmål: BrowseDesignViews

```mermaid
flowchart LR
    n_BrowseDesignViews["BrowseDesignViews (usecase)"]
    n_DesignAuthor["DesignAuthor (actor)"]
    n_DesignReviewer["DesignReviewer (actor)"]
    n_NavigableDesignDocumentation["NavigableDesignDocumentation (feature)"]
    n_DesignAuthor -->|pursues| n_BrowseDesignViews
    n_DesignReviewer -->|pursues| n_BrowseDesignViews
    n_NavigableDesignDocumentation -->|supports| n_BrowseDesignViews
```

Kildegrunnlag: f0142, f0150, f0617.

### Bruksmål: BuildNativeProduct

```mermaid
flowchart LR
    n_BuildNativeProduct["BuildNativeProduct (usecase)"]
    n_DesignAuthor["DesignAuthor (actor)"]
    n_NativeGoAssembly["NativeGoAssembly (feature)"]
    n_DesignAuthor -->|pursues| n_BuildNativeProduct
    n_NativeGoAssembly -->|supports| n_BuildNativeProduct
```

Kildegrunnlag: f0143, f0600.

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

Kildegrunnlag: f0144, f0567, f1078.

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

Kildegrunnlag: f0145, f0151, f1079, f1107.

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

Kildegrunnlag: f0146, f0152, f0530.

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

Kildegrunnlag: f0147, f0149, f0153.

### Bruksmål: TryDomainInteraction

```mermaid
flowchart LR
    n_DesignAuthor["DesignAuthor (actor)"]
    n_TryDomainInteraction["TryDomainInteraction (usecase)"]
    n_TypedDomainBinding["TypedDomainBinding (feature)"]
    n_DesignAuthor -->|pursues| n_TryDomainInteraction
    n_TypedDomainBinding -->|supports| n_TryDomainInteraction
```

Kildegrunnlag: f0148, f1112.

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

Kildegrunnlag: f0124, f0228, f0233, f0693, f1106, f1377.

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

Kildegrunnlag: f0018, f0062, f0121, f0516, f0717.

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

Kildegrunnlag: f0542, f0635, f0657, f0675, f0699, f0707, f0777.

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

Kildegrunnlag: f0053, f0465, f0473, f0678, f0801, f1243.

### Functionality-bidrag til Feature: NavigableDesignDocumentation

```mermaid
flowchart LR
    n_CaptureNavigationTarget["CaptureNavigationTarget (functionality)"]
    n_ComposeViewPackage["ComposeViewPackage (functionality)"]
    n_DispatchViewOpen["DispatchViewOpen (functionality)"]
    n_EvictUnusedViewBundles["EvictUnusedViewBundles (functionality)"]
    n_GenerateViewNavigation["GenerateViewNavigation (functionality)"]
    n_KeyViewRevision["KeyViewRevision (functionality)"]
    n_NavigableDesignDocumentation["NavigableDesignDocumentation (feature)"]
    n_PreserveViewAnchors["PreserveViewAnchors (functionality)"]
    n_ProjectSelectedView["ProjectSelectedView (functionality)"]
    n_PublishViewBundle["PublishViewBundle (functionality)"]
    n_RejectStaleViewResults["RejectStaleViewResults (functionality)"]
    n_ResolveConfiguredViewer["ResolveConfiguredViewer (functionality)"]
    n_RetainVisibleViewBundle["RetainVisibleViewBundle (functionality)"]
    n_RouteDocumentToPane["RouteDocumentToPane (functionality)"]
    n_ServeViewRequests["ServeViewRequests (functionality)"]
    n_ValidateViewRequest["ValidateViewRequest (functionality)"]
    n_CaptureNavigationTarget -->|contributes-to| n_NavigableDesignDocumentation
    n_ComposeViewPackage -->|contributes-to| n_NavigableDesignDocumentation
    n_DispatchViewOpen -->|contributes-to| n_NavigableDesignDocumentation
    n_EvictUnusedViewBundles -->|contributes-to| n_NavigableDesignDocumentation
    n_GenerateViewNavigation -->|contributes-to| n_NavigableDesignDocumentation
    n_KeyViewRevision -->|contributes-to| n_NavigableDesignDocumentation
    n_PreserveViewAnchors -->|contributes-to| n_NavigableDesignDocumentation
    n_ProjectSelectedView -->|contributes-to| n_NavigableDesignDocumentation
    n_PublishViewBundle -->|contributes-to| n_NavigableDesignDocumentation
    n_RejectStaleViewResults -->|contributes-to| n_NavigableDesignDocumentation
    n_ResolveConfiguredViewer -->|contributes-to| n_NavigableDesignDocumentation
    n_RetainVisibleViewBundle -->|contributes-to| n_NavigableDesignDocumentation
    n_RouteDocumentToPane -->|contributes-to| n_NavigableDesignDocumentation
    n_ServeViewRequests -->|contributes-to| n_NavigableDesignDocumentation
    n_ValidateViewRequest -->|contributes-to| n_NavigableDesignDocumentation
```

Kildegrunnlag: f0081, f0127, f0188, f0216, f0476, f0545, f0690, f0697, f0712, f0724, f0786, f0798, f0806, f1041, f1232.

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

Kildegrunnlag: f0065, f0070, f0778, f1218, f1235.

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

Kildegrunnlag: f0130, f0185, f0539, f0704, f0783.


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
    n_ViewServiceHost["ViewServiceHost (container)"]
    n_XfmdDocumentHost["XfmdDocumentHost (container)"]
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

Kildegrunnlag: f0132, f0133, f0134.

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

Kildegrunnlag: f0167, f0168, f0169, f0170, f0171, f0172.

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

Kildegrunnlag: f0841, f0842, f0843, f0844.

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

Kildegrunnlag: f0855, f0856, f0857.

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

Kildegrunnlag: f0876, f0877, f0878, f0879, f0880.

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

Kildegrunnlag: f0929, f0930, f0931, f0932.

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

Kildegrunnlag: f0974, f0975, f0976, f0977.

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

Kildegrunnlag: f1006, f1007, f1008, f1009.

### Logisk inndeling: ViewServiceHost

```mermaid
flowchart LR
    n_DocumentBroker["DocumentBroker (unit)"]
    n_ViewArtifactStore["ViewArtifactStore (unit)"]
    n_ViewServiceHost["ViewServiceHost (container)"]
    n_ViewerLaunchAdapter["ViewerLaunchAdapter (unit)"]
    n_ViewServiceHost -->|contains| n_DocumentBroker
    n_ViewServiceHost -->|contains| n_ViewArtifactStore
    n_ViewServiceHost -->|contains| n_ViewerLaunchAdapter
```

Kildegrunnlag: f1359, f1360, f1361.


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

Kildegrunnlag: f0131, f0183, f0705, f0784, f0807, f0891, f0892, f0893, f0894, f0895.

### Tilbydere av kapabilitet: BoundInteraction

```mermaid
flowchart LR
    n_BoundInteraction["BoundInteraction (capability)"]
    n_SdlUiBindingAdapter["SdlUiBindingAdapter (unit)"]
    n_SdlUiBindingAdapter -->|provides| n_BoundInteraction
```

Kildegrunnlag: f0896.

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

Kildegrunnlag: f0087, f0543, f0636, f0658, f0709, f0738, f0739, f0740, f0741, f0799, f1058, f1059.

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

Kildegrunnlag: f0175, f0742, f1060.

### Bidrag til kapabilitet: DomainOperations

```mermaid
flowchart LR
    n_DomainOperations["DomainOperations (capability)"]
    n_GoDomainImplementation["GoDomainImplementation (unit)"]
    n_PerformDomainOperation["PerformDomainOperation (functionality)"]
    n_GoDomainImplementation -->|owns| n_PerformDomainOperation
    n_PerformDomainOperation -->|realizes| n_DomainOperations
```

Kildegrunnlag: f0505, f0655.

### Tilbydere av kapabilitet: DomainOperations

```mermaid
flowchart LR
    n_DomainOperations["DomainOperations (capability)"]
    n_GoDomainImplementation["GoDomainImplementation (unit)"]
    n_GoDomainImplementation -->|provides| n_DomainOperations
```

Kildegrunnlag: f0506.

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

Kildegrunnlag: f0079, f0082, f0083, f0084, f0085, f0138, f0140, f0212, f0213, f0540, f0569, f0583, f0721, f0826, f0827, f0828, f0829, f0837, f0850, f0851, f0881, f0882, f0885, f0886, f1046, f1209.

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

Kildegrunnlag: f0214, f0830, f0838, f0852, f0858, f0883, f0887.

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

Kildegrunnlag: f0021, f0086, f0139, f0141, f0186, f0571, f0578, f0676, f0722, f0780, f0802, f0917, f0918, f0919, f0920, f0950, f0952, f0953, f1000, f1001, f1002, f1010, f1011, f1047, f1111, f1191, f1192, f1193, f1210, f1225.

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

Kildegrunnlag: f0921, f0954, f0978, f1003, f1012, f1194.

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

Kildegrunnlag: f0019, f0063, f0128, f0580, f0781, f0960, f0961, f0962, f0963, f0964.

### Tilbydere av kapabilitet: MeasuredPresentation

```mermaid
flowchart LR
    n_MeasuredPresentation["MeasuredPresentation (capability)"]
    n_SduiLayout["SduiLayout (unit)"]
    n_SduiLibrary["SduiLibrary (unit)"]
    n_SduiLayout -->|provides| n_MeasuredPresentation
    n_SduiLibrary -->|provides| n_MeasuredPresentation
```

Kildegrunnlag: f0965, f0979.

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

Kildegrunnlag: f0122, f0244, f0245, f0246, f0247, f0263, f0264, f0517, f0710, f0719, f0725, f0810.

### Tilbydere av kapabilitet: NativeInteraction

```mermaid
flowchart LR
    n_FyneBackend["FyneBackend (unit)"]
    n_FyneHost["FyneHost (container)"]
    n_NativeInteraction["NativeInteraction (capability)"]
    n_FyneBackend -->|provides| n_NativeInteraction
    n_FyneHost -->|provides| n_NativeInteraction
```

Kildegrunnlag: f0248, f0265.

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

Kildegrunnlag: f0054, f0466, f0474, f0487, f0488, f0495, f0496, f0497, f0679, f0796.

### Tilbydere av kapabilitet: NativeRealization

```mermaid
flowchart LR
    n_GoBuildRunner["GoBuildRunner (unit)"]
    n_GoCodeGenerator["GoCodeGenerator (unit)"]
    n_NativeRealization["NativeRealization (capability)"]
    n_GoBuildRunner -->|provides| n_NativeRealization
    n_GoCodeGenerator -->|provides| n_NativeRealization
```

Kildegrunnlag: f0489, f0498.

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

Kildegrunnlag: f0181, f0575, f0576, f0579, f0659, f0660, f0726, f0793, f0794, f1233.

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

Kildegrunnlag: f0135, f0182, f0577, f0795.

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

Kildegrunnlag: f0068, f0136, f0620, f0682, f0789, f0845, f0853, f0869, f0870, f0872, f0901, f0902, f0903, f1094, f1216, f1221.

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

Kildegrunnlag: f0846, f0854, f0859, f0871, f0873, f0904.

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

Kildegrunnlag: f0073, f0137, f0220, f0685, f0688, f0792, f0933, f0970, f0982, f0983, f0984, f0988, f1014, f1015, f1016, f1017, f1097, f1213, f1224, f1238.

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

Kildegrunnlag: f0934, f0971, f0980, f0985, f0989, f1018.

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

Kildegrunnlag: f0176, f0177, f0774, f0779.

### Tilbydere av kapabilitet: SourceDiagnostics

```mermaid
flowchart LR
    n_DiagnosticReporter["DiagnosticReporter (unit)"]
    n_SourceDiagnostics["SourceDiagnostics (capability)"]
    n_DiagnosticReporter -->|provides| n_SourceDiagnostics
```

Kildegrunnlag: f0178.

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

Kildegrunnlag: f0518, f0715, f1054, f1055.

### Tilbydere av kapabilitet: SourceLoading

```mermaid
flowchart LR
    n_SourceLoader["SourceLoader (unit)"]
    n_SourceLoading["SourceLoading (capability)"]
    n_SourceLoader -->|provides| n_SourceLoading
```

Kildegrunnlag: f1056.

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

Kildegrunnlag: f0092, f0093, f0119, f0125, f0221, f0229, f0993, f0994, f0995, f1378.

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

Kildegrunnlag: f0094, f0981, f0996.


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

Kildegrunnlag: f0050, f0051, f0217.

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

Kildegrunnlag: f0159, f0160, f0162, f0163, f0165.

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

Kildegrunnlag: f0161, f0164, f0166.

### Nødvendige porter i modus: NativeBuild

```mermaid
flowchart LR
    n_BuildToolPort["BuildToolPort (interface)"]
    n_GeneratedArtifactPort["GeneratedArtifactPort (interface)"]
    n_NativeRealization["NativeRealization (capability)"]
    n_NativeRealization -->|requires| n_BuildToolPort
    n_NativeRealization -->|requires| n_GeneratedArtifactPort
```

Kildegrunnlag: f0610, f0611.

### Nødvendige porter i modus: RichDocument

```mermaid
flowchart LR
    n_ContentProviderPort["ContentProviderPort (interface)"]
    n_DiagramEnginePort["DiagramEnginePort (interface)"]
    n_RichContent["RichContent (capability)"]
    n_RichContent -->|requires| n_ContentProviderPort
    n_RichContent -->|requires| n_DiagramEnginePort
```

Kildegrunnlag: f0803, f0804.

### Nødvendige porter i modus: SourceInspection

```mermaid
flowchart LR
    n_SdlSourceModel["SdlSourceModel (capability)"]
    n_SduiSourceModel["SduiSourceModel (capability)"]
    n_SourceSnapshotPort["SourceSnapshotPort (interface)"]
    n_SdlSourceModel -->|requires| n_SourceSnapshotPort
    n_SduiSourceModel -->|requires| n_SourceSnapshotPort
```

Kildegrunnlag: f0884, f1013.

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

Kildegrunnlag: f0581, f1067, f1068.

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

Kildegrunnlag: f0582, f0601.


## VP06 — Aktiviteter og leveranseplan

### Aktivitetsrøtter

```mermaid
flowchart LR
    n_G1FrontendPort["G1FrontendPort (activity)"]
    n_G2LayoutAndPresentation["G2LayoutAndPresentation (activity)"]
    n_G3UiRuntimeAndReload["G3UiRuntimeAndReload (activity)"]
    n_G4SdlRuntimeAndBinding["G4SdlRuntimeAndBinding (activity)"]
    n_G5NativeGeneration["G5NativeGeneration (activity)"]
    n_G6NavigableDocumentation["G6NavigableDocumentation (activity)"]
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

Kildegrunnlag: f0272, f0279, f0290, f0295.

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

Kildegrunnlag: f0296, f0304, f0309, f0318, f0326.

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

Kildegrunnlag: f0343, f0352, f0358, f0359.

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

Kildegrunnlag: f0370, f0385, f0395, f0402, f0403.

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

Kildegrunnlag: f0411, f0415, f0423, f0428, f0429, f0430.

### Aktivitetsinndeling: G6NavigableDocumentation

```mermaid
flowchart LR
    n_G6M1StaticNavigation["G6M1StaticNavigation (activity)"]
    n_G6M2OnDemandViews["G6M2OnDemandViews (activity)"]
    n_G6M3XfmdNavigation["G6M3XfmdNavigation (activity)"]
    n_G6M4SessionPublication["G6M4SessionPublication (activity)"]
    n_G6NavigableDocumentation["G6NavigableDocumentation (activity)"]
    n_NavigableDesignDocumentation["NavigableDesignDocumentation (feature)"]
    n_G6M1StaticNavigation -->|refines| n_G6NavigableDocumentation
    n_G6M2OnDemandViews -->|refines| n_G6NavigableDocumentation
    n_G6M3XfmdNavigation -->|refines| n_G6NavigableDocumentation
    n_G6M4SessionPublication -->|refines| n_G6NavigableDocumentation
    n_G6NavigableDocumentation -->|delivers| n_NavigableDesignDocumentation
```

Kildegrunnlag: f0440, f0447, f0454, f0461, f0462.

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

Kildegrunnlag: f0230, f0519, f0520.

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

Kildegrunnlag: f0730, f0773.

### Aktivitetsinndeling: RunDesignSession

```mermaid
flowchart LR
    n_RunBoundUiAction["RunBoundUiAction (activity)"]
    n_RunDesignSession["RunDesignSession (activity)"]
    n_RunUnboundUiPreview["RunUnboundUiPreview (activity)"]
    n_RunBoundUiAction -->|refines| n_RunDesignSession
    n_RunUnboundUiPreview -->|refines| n_RunDesignSession
```

Kildegrunnlag: f0808, f0809.

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

Kildegrunnlag: f0274, f0275, f0276, f0277, f0970, f0988, f1054, f1055.

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

Kildegrunnlag: f0280, f0281, f0282, f0283, f0284, f0285, f0286, f0287, f0933, f0982, f0983, f0984, f1014, f1015, f1016, f1017.

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

Kildegrunnlag: f0177, f0291, f0292, f0994.

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

Kildegrunnlag: f0298, f0299, f0300, f0301, f0960, f0962, f0963, f0964.

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

Kildegrunnlag: f0305, f0306, f0961, f0995.

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

Kildegrunnlag: f0244, f0245, f0246, f0247, f0263, f0264, f0310, f0311, f0312, f0313, f0314, f0315.

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

Kildegrunnlag: f0181, f0319, f0320, f0321, f0322, f0323, f0575, f0576, f0793, f0794.

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

Kildegrunnlag: f0327, f0328, f0329, f0330, f0331, f0332, f0333, f0334, f0335, f0336, f0337, f0338, f0339, f0917, f0918, f0919, f0920, f0950, f0951, f0952, f0953, f1000, f1001, f1002, f1010, f1011.

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

Kildegrunnlag: f0344, f0345, f0346, f0347, f0348, f0349, f0738, f0739, f0740, f0741, f1058, f1059.

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

Kildegrunnlag: f0353, f0354, f0355, f1191, f1192, f1193.

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

Kildegrunnlag: f0361, f0362, f0363, f0364, f0365, f0366, f0367, f0368, f0845, f0853, f0869, f0870, f0872, f0901, f0902, f0903.

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

Kildegrunnlag: f0371, f0372, f0373, f0374, f0375, f0376, f0377, f0378, f0379, f0380, f0381, f0382, f0505, f0826, f0827, f0828, f0829, f0837, f0850, f0851, f0881, f0882, f0885, f0886.

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

Kildegrunnlag: f0176, f0386, f0387, f0388, f0389, f0390, f0391, f0891, f0892, f0893, f0894, f0895.

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

Kildegrunnlag: f0212, f0213, f0396, f0397, f0398, f0488.

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

Kildegrunnlag: f0405, f0406, f0407, f0408, f0487, f0495, f0496, f0497.

### Planlagt ansvar: G5M2BehaviorParity

```mermaid
flowchart LR
    n_DevelopmentTools["DevelopmentTools (unit)"]
    n_G5M2BehaviorParity["G5M2BehaviorParity (activity)"]
    n_VerifyNativeBehaviorParity["VerifyNativeBehaviorParity (functionality)"]
    n_DevelopmentTools -->|owns| n_VerifyNativeBehaviorParity
    n_G5M2BehaviorParity -->|addresses| n_VerifyNativeBehaviorParity
```

Kildegrunnlag: f0174, f0412.

### Planlagt ansvar: G5M3DocumentationExport

```mermaid
flowchart LR
    n_CommandLineHost["CommandLineHost (container)"]
    n_ComposeHeadlessExport["ComposeHeadlessExport (functionality)"]
    n_ComposeMarkdownDocument["ComposeMarkdownDocument (functionality)"]
    n_G5M3DocumentationExport["G5M3DocumentationExport (activity)"]
    n_SduiPresentation["SduiPresentation (unit)"]
    n_WriteGeneratedArtifacts["WriteGeneratedArtifacts (functionality)"]
    n_CommandLineHost -->|owns| n_ComposeHeadlessExport
    n_CommandLineHost -->|owns| n_WriteGeneratedArtifacts
    n_G5M3DocumentationExport -->|addresses| n_ComposeHeadlessExport
    n_G5M3DocumentationExport -->|addresses| n_ComposeMarkdownDocument
    n_G5M3DocumentationExport -->|addresses| n_WriteGeneratedArtifacts
    n_SduiPresentation -->|owns| n_ComposeMarkdownDocument
```

Kildegrunnlag: f0092, f0093, f0416, f0417, f0418, f0993.

### Planlagt ansvar: G5M4RetirePython

```mermaid
flowchart LR
    n_DevelopmentTools["DevelopmentTools (unit)"]
    n_G5M4RetirePython["G5M4RetirePython (activity)"]
    n_RetireReplacedPythonEntryPoints["RetireReplacedPythonEntryPoints (functionality)"]
    n_DevelopmentTools -->|owns| n_RetireReplacedPythonEntryPoints
    n_G5M4RetirePython -->|addresses| n_RetireReplacedPythonEntryPoints
```

Kildegrunnlag: f0173, f0424.

### Planlagt ansvar: G6M1StaticNavigation

```mermaid
flowchart LR
    n_ComposeViewPackage["ComposeViewPackage (functionality)"]
    n_ExportViewpointMarkdown["ExportViewpointMarkdown (functionality)"]
    n_G6M1StaticNavigation["G6M1StaticNavigation (activity)"]
    n_GenerateViewNavigation["GenerateViewNavigation (functionality)"]
    n_PreserveViewAnchors["PreserveViewAnchors (functionality)"]
    n_ProjectSdlViewpoints["ProjectSdlViewpoints (functionality)"]
    n_SdlViewpointGenerator["SdlViewpointGenerator (unit)"]
    n_TraceViewpointFacts["TraceViewpointFacts (functionality)"]
    n_G6M1StaticNavigation -->|addresses| n_ComposeViewPackage
    n_G6M1StaticNavigation -->|addresses| n_ExportViewpointMarkdown
    n_G6M1StaticNavigation -->|addresses| n_GenerateViewNavigation
    n_G6M1StaticNavigation -->|addresses| n_PreserveViewAnchors
    n_G6M1StaticNavigation -->|addresses| n_ProjectSdlViewpoints
    n_G6M1StaticNavigation -->|addresses| n_TraceViewpointFacts
    n_SdlViewpointGenerator -->|owns| n_ComposeViewPackage
    n_SdlViewpointGenerator -->|owns| n_ExportViewpointMarkdown
    n_SdlViewpointGenerator -->|owns| n_GenerateViewNavigation
    n_SdlViewpointGenerator -->|owns| n_PreserveViewAnchors
    n_SdlViewpointGenerator -->|owns| n_ProjectSdlViewpoints
    n_SdlViewpointGenerator -->|owns| n_TraceViewpointFacts
```

Kildegrunnlag: f0432, f0433, f0434, f0435, f0436, f0437, f0905, f0906, f0907, f0908, f0909, f0911.

### Planlagt ansvar: G6M2OnDemandViews

```mermaid
flowchart LR
    n_DocumentBroker["DocumentBroker (unit)"]
    n_G6M2OnDemandViews["G6M2OnDemandViews (activity)"]
    n_KeyViewRevision["KeyViewRevision (functionality)"]
    n_ProjectSelectedView["ProjectSelectedView (functionality)"]
    n_PublishViewBundle["PublishViewBundle (functionality)"]
    n_SdlViewpointGenerator["SdlViewpointGenerator (unit)"]
    n_ValidateViewRequest["ValidateViewRequest (functionality)"]
    n_ViewArtifactStore["ViewArtifactStore (unit)"]
    n_DocumentBroker -->|owns| n_KeyViewRevision
    n_DocumentBroker -->|owns| n_ValidateViewRequest
    n_G6M2OnDemandViews -->|addresses| n_KeyViewRevision
    n_G6M2OnDemandViews -->|addresses| n_ProjectSelectedView
    n_G6M2OnDemandViews -->|addresses| n_PublishViewBundle
    n_G6M2OnDemandViews -->|addresses| n_ValidateViewRequest
    n_SdlViewpointGenerator -->|owns| n_ProjectSelectedView
    n_ViewArtifactStore -->|owns| n_PublishViewBundle
```

Kildegrunnlag: f0191, f0194, f0441, f0442, f0443, f0444, f0910, f1245.

### Planlagt ansvar: G6M3XfmdNavigation

```mermaid
flowchart LR
    n_CaptureNavigationTarget["CaptureNavigationTarget (functionality)"]
    n_DispatchViewOpen["DispatchViewOpen (functionality)"]
    n_G6M3XfmdNavigation["G6M3XfmdNavigation (activity)"]
    n_ResolveConfiguredViewer["ResolveConfiguredViewer (functionality)"]
    n_RouteDocumentToPane["RouteDocumentToPane (functionality)"]
    n_ViewerLaunchAdapter["ViewerLaunchAdapter (unit)"]
    n_XfmdDocumentHost["XfmdDocumentHost (container)"]
    n_G6M3XfmdNavigation -->|addresses| n_CaptureNavigationTarget
    n_G6M3XfmdNavigation -->|addresses| n_DispatchViewOpen
    n_G6M3XfmdNavigation -->|addresses| n_ResolveConfiguredViewer
    n_G6M3XfmdNavigation -->|addresses| n_RouteDocumentToPane
    n_ViewerLaunchAdapter -->|owns| n_DispatchViewOpen
    n_ViewerLaunchAdapter -->|owns| n_ResolveConfiguredViewer
    n_XfmdDocumentHost -->|owns| n_CaptureNavigationTarget
    n_XfmdDocumentHost -->|owns| n_RouteDocumentToPane
```

Kildegrunnlag: f0448, f0449, f0450, f0451, f1369, f1370, f1379, f1380.

### Planlagt ansvar: G6M4SessionPublication

```mermaid
flowchart LR
    n_DocumentBroker["DocumentBroker (unit)"]
    n_EvictUnusedViewBundles["EvictUnusedViewBundles (functionality)"]
    n_G6M4SessionPublication["G6M4SessionPublication (activity)"]
    n_RejectStaleViewResults["RejectStaleViewResults (functionality)"]
    n_RetainVisibleViewBundle["RetainVisibleViewBundle (functionality)"]
    n_ServeViewRequests["ServeViewRequests (functionality)"]
    n_ViewArtifactStore["ViewArtifactStore (unit)"]
    n_DocumentBroker -->|owns| n_RejectStaleViewResults
    n_DocumentBroker -->|owns| n_ServeViewRequests
    n_G6M4SessionPublication -->|addresses| n_EvictUnusedViewBundles
    n_G6M4SessionPublication -->|addresses| n_RejectStaleViewResults
    n_G6M4SessionPublication -->|addresses| n_RetainVisibleViewBundle
    n_G6M4SessionPublication -->|addresses| n_ServeViewRequests
    n_ViewArtifactStore -->|owns| n_EvictUnusedViewBundles
    n_ViewArtifactStore -->|owns| n_RetainVisibleViewBundle
```

Kildegrunnlag: f0192, f0193, f0455, f0456, f0457, f0458, f1244, f1246.

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
    n_G6M1StaticNavigation["G6M1StaticNavigation (activity)"]
    n_G6M2OnDemandViews["G6M2OnDemandViews (activity)"]
    n_G6M3XfmdNavigation["G6M3XfmdNavigation (activity)"]
    n_G6M4SessionPublication["G6M4SessionPublication (activity)"]
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
    n_G5M3DocumentationExport -->|depends-on| n_G6M1StaticNavigation
    n_G5M4RetirePython -->|depends-on| n_G1M3Concept1AndDumps
    n_G5M4RetirePython -->|depends-on| n_G5M3DocumentationExport
    n_G6M1StaticNavigation -->|depends-on| n_G4M1SdlFrontend
    n_G6M2OnDemandViews -->|depends-on| n_G6M1StaticNavigation
    n_G6M3XfmdNavigation -->|depends-on| n_G6M2OnDemandViews
    n_G6M4SessionPublication -->|depends-on| n_G6M3XfmdNavigation
```

Kildegrunnlag: f0288, f0293, f0302, f0307, f0316, f0324, f0340, f0341, f0350, f0356, f0383, f0392, f0393, f0399, f0400, f0409, f0413, f0419, f0420, f0421, f0425, f0426, f0438, f0445, f0452, f0459.


## VP07 — Features over arkitekturen

### Feature: DesignDocumentation — modus DocumentBrowsing

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
    n_ViewServiceHost["ViewServiceHost (container)"]
    n_WriteGeneratedArtifacts["WriteGeneratedArtifacts (functionality)"]
    n_CommandLineHost -->|owns| n_WriteGeneratedArtifacts
    n_ComposeMarkdownDocument -->|contributes-to| n_DesignDocumentation
    n_ExportSvgSnapshot -->|contributes-to| n_DesignDocumentation
    n_ExportViewpointMarkdown -->|allocated-to| n_ViewServiceHost
    n_ExportViewpointMarkdown -->|contributes-to| n_DesignDocumentation
    n_ProjectSdlViewpoints -->|allocated-to| n_ViewServiceHost
    n_ProjectSdlViewpoints -->|contributes-to| n_DesignDocumentation
    n_SdlViewpointGenerator -->|owns| n_ExportViewpointMarkdown
    n_SdlViewpointGenerator -->|owns| n_ProjectSdlViewpoints
    n_SdlViewpointGenerator -->|owns| n_TraceViewpointFacts
    n_SduiPresentation -->|owns| n_ComposeMarkdownDocument
    n_SduiPresentation -->|owns| n_ExportSvgSnapshot
    n_TraceViewpointFacts -->|allocated-to| n_ViewServiceHost
    n_TraceViewpointFacts -->|contributes-to| n_DesignDocumentation
    n_WriteGeneratedArtifacts -->|contributes-to| n_DesignDocumentation
```

Kildegrunnlag: f0093, f0124, f0228, f0232, f0233, f0692, f0693, f0906, f0909, f0911, f0993, f0995, f1105, f1106, f1377.

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

Kildegrunnlag: f0093, f0124, f0228, f0233, f0693, f0906, f0909, f0911, f0993, f0995, f1103, f1106, f1377.

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

Kildegrunnlag: f0093, f0123, f0124, f0227, f0228, f0231, f0233, f0691, f0693, f0906, f0909, f0911, f0993, f0995, f1104, f1106, f1376, f1377.

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

Kildegrunnlag: f0017, f0018, f0061, f0062, f0120, f0121, f0244, f0246, f0263, f0515, f0516, f0716, f0717, f0960, f0961.

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

Kildegrunnlag: f0177, f0541, f0542, f0634, f0635, f0656, f0657, f0674, f0675, f0698, f0699, f0706, f0707, f0738, f0739, f0740, f0776, f0777, f0951, f1059, f1192.

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

Kildegrunnlag: f0177, f0542, f0635, f0657, f0675, f0699, f0707, f0738, f0739, f0740, f0775, f0777, f0951, f1059, f1192.

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

Kildegrunnlag: f0052, f0053, f0173, f0174, f0464, f0465, f0472, f0473, f0487, f0495, f0496, f0497, f0677, f0678, f0800, f0801, f1242, f1243.

### Feature: NavigableDesignDocumentation — modus DocumentBrowsing

```mermaid
flowchart LR
    n_CaptureNavigationTarget["CaptureNavigationTarget (functionality)"]
    n_ComposeViewPackage["ComposeViewPackage (functionality)"]
    n_DispatchViewOpen["DispatchViewOpen (functionality)"]
    n_DocumentBroker["DocumentBroker (unit)"]
    n_EvictUnusedViewBundles["EvictUnusedViewBundles (functionality)"]
    n_GenerateViewNavigation["GenerateViewNavigation (functionality)"]
    n_KeyViewRevision["KeyViewRevision (functionality)"]
    n_NavigableDesignDocumentation["NavigableDesignDocumentation (feature)"]
    n_PreserveViewAnchors["PreserveViewAnchors (functionality)"]
    n_ProjectSelectedView["ProjectSelectedView (functionality)"]
    n_PublishViewBundle["PublishViewBundle (functionality)"]
    n_RejectStaleViewResults["RejectStaleViewResults (functionality)"]
    n_ResolveConfiguredViewer["ResolveConfiguredViewer (functionality)"]
    n_RetainVisibleViewBundle["RetainVisibleViewBundle (functionality)"]
    n_RouteDocumentToPane["RouteDocumentToPane (functionality)"]
    n_SdlViewpointGenerator["SdlViewpointGenerator (unit)"]
    n_ServeViewRequests["ServeViewRequests (functionality)"]
    n_ValidateViewRequest["ValidateViewRequest (functionality)"]
    n_ViewArtifactStore["ViewArtifactStore (unit)"]
    n_ViewServiceHost["ViewServiceHost (container)"]
    n_ViewerLaunchAdapter["ViewerLaunchAdapter (unit)"]
    n_XfmdDocumentHost["XfmdDocumentHost (container)"]
    n_CaptureNavigationTarget -->|allocated-to| n_XfmdDocumentHost
    n_CaptureNavigationTarget -->|contributes-to| n_NavigableDesignDocumentation
    n_ComposeViewPackage -->|allocated-to| n_ViewServiceHost
    n_ComposeViewPackage -->|contributes-to| n_NavigableDesignDocumentation
    n_DispatchViewOpen -->|allocated-to| n_ViewServiceHost
    n_DispatchViewOpen -->|contributes-to| n_NavigableDesignDocumentation
    n_DocumentBroker -->|owns| n_KeyViewRevision
    n_DocumentBroker -->|owns| n_RejectStaleViewResults
    n_DocumentBroker -->|owns| n_ServeViewRequests
    n_DocumentBroker -->|owns| n_ValidateViewRequest
    n_EvictUnusedViewBundles -->|allocated-to| n_ViewServiceHost
    n_EvictUnusedViewBundles -->|contributes-to| n_NavigableDesignDocumentation
    n_GenerateViewNavigation -->|allocated-to| n_ViewServiceHost
    n_GenerateViewNavigation -->|contributes-to| n_NavigableDesignDocumentation
    n_KeyViewRevision -->|allocated-to| n_ViewServiceHost
    n_KeyViewRevision -->|contributes-to| n_NavigableDesignDocumentation
    n_PreserveViewAnchors -->|allocated-to| n_ViewServiceHost
    n_PreserveViewAnchors -->|contributes-to| n_NavigableDesignDocumentation
    n_ProjectSelectedView -->|allocated-to| n_ViewServiceHost
    n_ProjectSelectedView -->|contributes-to| n_NavigableDesignDocumentation
    n_PublishViewBundle -->|allocated-to| n_ViewServiceHost
    n_PublishViewBundle -->|contributes-to| n_NavigableDesignDocumentation
    n_RejectStaleViewResults -->|allocated-to| n_ViewServiceHost
    n_RejectStaleViewResults -->|contributes-to| n_NavigableDesignDocumentation
    n_ResolveConfiguredViewer -->|allocated-to| n_ViewServiceHost
    n_ResolveConfiguredViewer -->|contributes-to| n_NavigableDesignDocumentation
    n_RetainVisibleViewBundle -->|allocated-to| n_ViewServiceHost
    n_RetainVisibleViewBundle -->|contributes-to| n_NavigableDesignDocumentation
    n_RouteDocumentToPane -->|allocated-to| n_XfmdDocumentHost
    n_RouteDocumentToPane -->|contributes-to| n_NavigableDesignDocumentation
    n_SdlViewpointGenerator -->|owns| n_ComposeViewPackage
    n_SdlViewpointGenerator -->|owns| n_GenerateViewNavigation
    n_SdlViewpointGenerator -->|owns| n_PreserveViewAnchors
    n_SdlViewpointGenerator -->|owns| n_ProjectSelectedView
    n_ServeViewRequests -->|allocated-to| n_ViewServiceHost
    n_ServeViewRequests -->|contributes-to| n_NavigableDesignDocumentation
    n_ValidateViewRequest -->|allocated-to| n_ViewServiceHost
    n_ValidateViewRequest -->|contributes-to| n_NavigableDesignDocumentation
    n_ViewArtifactStore -->|owns| n_EvictUnusedViewBundles
    n_ViewArtifactStore -->|owns| n_PublishViewBundle
    n_ViewArtifactStore -->|owns| n_RetainVisibleViewBundle
    n_ViewerLaunchAdapter -->|owns| n_DispatchViewOpen
    n_ViewerLaunchAdapter -->|owns| n_ResolveConfiguredViewer
    n_XfmdDocumentHost -->|owns| n_CaptureNavigationTarget
    n_XfmdDocumentHost -->|owns| n_RouteDocumentToPane
```

Kildegrunnlag: f0080, f0081, f0126, f0127, f0187, f0188, f0191, f0192, f0193, f0194, f0215, f0216, f0475, f0476, f0544, f0545, f0689, f0690, f0696, f0697, f0711, f0712, f0723, f0724, f0785, f0786, f0797, f0798, f0805, f0806, f0905, f0907, f0908, f0910, f1040, f1041, f1231, f1232, f1244, f1245, f1246, f1369, f1370, f1379, f1380.

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

Kildegrunnlag: f0065, f0070, f0177, f0776, f0778, f0872, f0903, f0988, f1017, f1218, f1235.

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

Kildegrunnlag: f0064, f0065, f0069, f0070, f0177, f0775, f0778, f0872, f0903, f0988, f1017, f1217, f1218, f1234, f1235.

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

Kildegrunnlag: f0129, f0130, f0184, f0185, f0538, f0539, f0703, f0704, f0782, f0783, f0828, f0891, f0893, f0894, f0918.

### Modellhull i dette utsnittet

| Identitet | Modus | Mangel |
| --- | --- | --- |
| ComposeMarkdownDocument | DocumentBrowsing | Container-allokering er uspesifisert for bidrag til DesignDocumentation. |
| ExportSvgSnapshot | DocumentBrowsing | Container-allokering er uspesifisert for bidrag til DesignDocumentation. |
| WriteGeneratedArtifacts | DocumentBrowsing | Container-allokering er uspesifisert for bidrag til DesignDocumentation. |
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

Kildegrunnlag: f0026, f0027, f0028, f0029, f0030, f0031, f0032, f0033, f0034, f0035, f0036, f0037, f0038, f0045, f0046, f0047, f0048, f0049, f0205, f0206, f0207, f0208, f0209, f0252, f0253, f0254, f0501, f0503, f0504, f0507, f0509, f0612, f0615, f0616, f0811, f0813, f0814, f0815, f0816, f0817, f0818, f0819, f0831, f0832, f0833, f0834, f0897, f0898, f0899, f0900, f0922, f0924, f0925, f0926, f0955, f1116, f1117, f1118, f1119, f1120, f1145, f1147, f1148, f1153, f1155, f1157.

### Scenario: BoundActionRejected — modus BoundExecution

```mermaid
sequenceDiagram
    participant n_FyneBackend as Fyne Backend
    participant n_SduiDispatcher as Sdui Dispatcher
    n_FyneBackend->>n_SduiDispatcher: 1: Ui Action Request (Native Ui Actions)
    n_SduiDispatcher-->>n_FyneBackend: 2: Ui Action Rejected (Native Ui Actions) reply-to 1
```

Kildegrunnlag: f0039, f0040, f0041, f0042, f0043, f0044, f0251, f0253, f0612, f0614, f0615, f0922, f0923, f1113, f1114, f1115, f1116, f1117.

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

Kildegrunnlag: f0238, f0240, f0241, f0249, f0250, f0266, f0267, f0270, f0271, f0521, f0522, f0523, f0524, f0525, f0526, f0527, f0528, f0529, f0552, f0554, f0555, f0558, f0559, f0560, f0561, f0562, f0665, f0666, f0667, f0668, f0669, f0967, f0969.

### Scenario: InvalidViewSelectionRejected — modus DocumentBrowsing

```mermaid
sequenceDiagram
    participant n_XfmdDocumentHost as Xfmd Document Host
    participant n_DocumentBroker as Document Broker
    n_XfmdDocumentHost->>n_DocumentBroker: 1: Select View Request (View Navigation Calls)
    n_DocumentBroker-->>n_XfmdDocumentHost: 2: View Request Rejected (View Navigation Calls) reply-to 1
```

Kildegrunnlag: f0197, f0199, f0531, f0532, f0533, f0534, f0535, f0536, f0537, f1021, f1022, f1290, f1292, f1294, f1346, f1347, f1348, f1385, f1386.

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

Kildegrunnlag: f0055, f0056, f0057, f0058, f0059, f0095, f0096, f0097, f0098, f0467, f0468, f0469, f0470, f0471, f0481, f0483, f0484, f0490, f0491, f0499, f0500, f0511, f0513, f0514, f0602, f0603, f0604, f0605, f0606, f0607, f0608, f0609.

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

Kildegrunnlag: f0106, f0107, f0108, f0109, f0110, f0588, f0590, f0592, f0743, f0745, f0749, f0750, f0759, f0760, f0761, f0767, f0768, f0820, f0822, f0823, f0847, f0848, f0860, f0861, f0862, f0863, f0864, f0865, f0866, f0867, f0868, f1061, f1065.

### Scenario: SelectedViewOpened — modus DocumentBrowsing

```mermaid
sequenceDiagram
    participant n_XfmdDocumentHost as Xfmd Document Host
    participant n_DocumentBroker as Document Broker
    participant n_SdlViewpointGenerator as Sdl Viewpoint Generator
    participant n_ViewArtifactStore as View Artifact Store
    participant n_ViewerLaunchAdapter as Viewer Launch Adapter
    n_XfmdDocumentHost->>n_DocumentBroker: 1: Select View Request (View Navigation Calls)
    n_DocumentBroker->>n_SdlViewpointGenerator: 2: Project View Request (View Projection Calls)
    n_SdlViewpointGenerator-->>n_DocumentBroker: 3: View Bundle Result (View Projection Calls) reply-to 2
    n_DocumentBroker->>n_ViewArtifactStore: 4: Publish View Request (View Publication Calls)
    n_ViewArtifactStore-->>n_DocumentBroker: 5: View Reference Result (View Publication Calls) reply-to 4
    n_DocumentBroker->>n_ViewerLaunchAdapter: 6: Launch View Request (View Launch Calls)
    n_ViewerLaunchAdapter->>n_XfmdDocumentHost: 7: Display View Request (View Display Calls)
    n_XfmdDocumentHost-->>n_ViewerLaunchAdapter: 8: View Display Result (View Display Calls) reply-to 7
    n_ViewerLaunchAdapter-->>n_DocumentBroker: 9: View Launch Result (View Launch Calls) reply-to 6
    n_DocumentBroker-->>n_XfmdDocumentHost: 10: View Opened Result (View Navigation Calls) reply-to 1
```

Kildegrunnlag: f0189, f0190, f0195, f0196, f0197, f0198, f0200, f0202, f0203, f0204, f0546, f0547, f0701, f0702, f0713, f0714, f0912, f0913, f1021, f1022, f1023, f1024, f1025, f1026, f1027, f1028, f1029, f1030, f1031, f1032, f1033, f1034, f1035, f1036, f1037, f1038, f1039, f1247, f1248, f1254, f1255, f1256, f1261, f1263, f1264, f1266, f1267, f1268, f1279, f1281, f1282, f1283, f1284, f1285, f1290, f1292, f1293, f1306, f1307, f1308, f1313, f1315, f1316, f1330, f1332, f1333, f1341, f1342, f1343, f1371, f1373, f1374, f1375, f1381, f1382, f1384, f1386.

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

Kildegrunnlag: f0099, f0100, f0101, f0102, f0222, f0223, f0224, f0225, f0226, f0552, f0554, f0555, f0558, f0559, f0560, f0561, f0562, f0966, f0968, f0997, f0998, f1069, f1070, f1071, f1072, f1073, f1074, f1075, f1076, f1077, f1084, f1086, f1087.

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

Kildegrunnlag: f0074, f0075, f0076, f0077, f0078, f0104, f0105, f0114, f0115, f0116, f0117, f0118, f0621, f0622, f0623, f0624, f0625, f0935, f0936, f0938, f0942, f0943, f0944, f0945, f0946, f0947, f0948, f0972, f0973, f0986, f0987, f0990, f0991, f1019, f1020, f1098, f1099, f1100, f1101, f1102, f1121, f1123, f1124, f1125, f1126, f1127, f1128, f1129, f1130, f1131, f1132, f1133, f1134, f1135, f1136, f1137, f1138, f1139, f1140, f1143, f1144, f1181, f1183, f1184, f1195, f1197, f1198, f1199, f1201, f1202, f1226, f1227, f1228, f1229, f1230.

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

Kildegrunnlag: f0114, f0115, f0116, f0117, f0118, f0255, f0588, f0590, f0592, f0744, f0746, f0752, f0753, f0759, f0760, f0761, f0767, f0768, f0937, f0941, f0956, f1062, f1066, f1140, f1143, f1144, f1153, f1155, f1157, f1163, f1164, f1165, f1166, f1167, f1168, f1169, f1170, f1171, f1172.

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

Kildegrunnlag: f0111, f0112, f0113, f0114, f0115, f0588, f0591, f0592, f0744, f0748, f0751, f0753, f0764, f0765, f0766, f0767, f0768, f0937, f0939, f1064, f1066, f1140, f1142, f1143, f1173, f1174, f1175, f1176, f1177, f1178, f1179, f1180.

### Scenario: UnboundLocalAction — modus UiPreview

```mermaid
sequenceDiagram
    participant n_FyneHost as Fyne Host
    participant n_GoDomainImplementation as Go Domain Implementation
    n_FyneHost->>n_GoDomainImplementation: 1: Domain Action Request (Go Domain Calls)
    n_GoDomainImplementation-->>n_FyneHost: 2: Domain Action Result (Go Domain Calls) reply-to 1
```

Kildegrunnlag: f0205, f0206, f0207, f0208, f0209, f0268, f0269, f0501, f0503, f0504, f0508, f0510, f1203, f1204, f1205, f1206, f1207, f1208.

### Scenario: ViewProjectionFailed — modus DocumentBrowsing

```mermaid
sequenceDiagram
    participant n_XfmdDocumentHost as Xfmd Document Host
    participant n_DocumentBroker as Document Broker
    participant n_SdlViewpointGenerator as Sdl Viewpoint Generator
    n_XfmdDocumentHost->>n_DocumentBroker: 1: Select View Request (View Navigation Calls)
    n_DocumentBroker->>n_SdlViewpointGenerator: 2: Project View Request (View Projection Calls)
    n_SdlViewpointGenerator-->>n_DocumentBroker: 3: View Projection Rejected (View Projection Calls) reply-to 2
    n_DocumentBroker-->>n_XfmdDocumentHost: 4: View Request Rejected (View Navigation Calls) reply-to 1
```

Kildegrunnlag: f0197, f0199, f0201, f0202, f0701, f0702, f0912, f0914, f1021, f1022, f1290, f1292, f1294, f1313, f1315, f1317, f1318, f1319, f1320, f1321, f1322, f1323, f1324, f1325, f1326, f1327, f1328, f1329, f1346, f1347, f1348, f1385, f1386.


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

Kildegrunnlag: f0154, f0155, f1053.

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

Kildegrunnlag: f0949, f1156, f1157, f1188.

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

Kildegrunnlag: f0157, f0158.

### Kontraktstruktur: FramePresentationCallsProtocol

```mermaid
flowchart LR
    n_FramePresentationCallsProtocol["FramePresentationCallsProtocol (contract)"]
    n_PresentFrameRequest["PresentFrameRequest (message)"]
    n_PresentFrameResult["PresentFrameResult (message)"]
    n_FramePresentationCallsProtocol -->|permits| n_PresentFrameRequest
    n_FramePresentationCallsProtocol -->|permits| n_PresentFrameResult
```

Kildegrunnlag: f0240, f0241.

### Kontraktstruktur: GeneratedGoContract

```mermaid
flowchart LR
    n_GeneratedGoContract["GeneratedGoContract (contract)"]
    n_GeneratedGoSources["GeneratedGoSources (field)"]
    n_GeneratedGoContract -->|has-field| n_GeneratedGoSources
```

Kildegrunnlag: f0478.

### Kontraktstruktur: GoBuildCallsProtocol

```mermaid
flowchart LR
    n_BuildGoRequest["BuildGoRequest (message)"]
    n_BuildGoResult["BuildGoResult (message)"]
    n_GoBuildCallsProtocol["GoBuildCallsProtocol (contract)"]
    n_GoBuildCallsProtocol -->|permits| n_BuildGoRequest
    n_GoBuildCallsProtocol -->|permits| n_BuildGoResult
```

Kildegrunnlag: f0483, f0484.

### Kontraktstruktur: GoDomainCallsProtocol

```mermaid
flowchart LR
    n_DomainActionRequest["DomainActionRequest (message)"]
    n_DomainActionResult["DomainActionResult (message)"]
    n_GoDomainCallsProtocol["GoDomainCallsProtocol (contract)"]
    n_GoDomainCallsProtocol -->|permits| n_DomainActionRequest
    n_GoDomainCallsProtocol -->|permits| n_DomainActionResult
```

Kildegrunnlag: f0503, f0504.

### Kontraktstruktur: GoGenerationCallsProtocol

```mermaid
flowchart LR
    n_GenerateGoRequest["GenerateGoRequest (message)"]
    n_GenerateGoResult["GenerateGoResult (message)"]
    n_GoGenerationCallsProtocol["GoGenerationCallsProtocol (contract)"]
    n_GoGenerationCallsProtocol -->|permits| n_GenerateGoRequest
    n_GoGenerationCallsProtocol -->|permits| n_GenerateGoResult
```

Kildegrunnlag: f0513, f0514.

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

Kildegrunnlag: f0549, f0550, f0551.

### Kontraktstruktur: LayoutCallsProtocol

```mermaid
flowchart LR
    n_LayoutCallsProtocol["LayoutCallsProtocol (contract)"]
    n_LayoutRequest["LayoutRequest (message)"]
    n_LayoutResult["LayoutResult (message)"]
    n_LayoutCallsProtocol -->|permits| n_LayoutRequest
    n_LayoutCallsProtocol -->|permits| n_LayoutResult
```

Kildegrunnlag: f0554, f0555.

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

Kildegrunnlag: f0590, f0591, f0592.

### Kontraktstruktur: NativeBuildContract

```mermaid
flowchart LR
    n_NativeBuildContract["NativeBuildContract (contract)"]
    n_NativeBuildDiagnostics["NativeBuildDiagnostics (field)"]
    n_NativeBuildSucceeded["NativeBuildSucceeded (field)"]
    n_NativeBuildContract -->|has-field| n_NativeBuildDiagnostics
    n_NativeBuildContract -->|has-field| n_NativeBuildSucceeded
```

Kildegrunnlag: f0594, f0595.

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

Kildegrunnlag: f0614, f0615, f0616.

### Kontraktstruktur: NormalizedModelContract

```mermaid
flowchart LR
    n_NormalizedModelArtifact["NormalizedModelArtifact (field)"]
    n_NormalizedModelContract["NormalizedModelContract (contract)"]
    n_NormalizedModelContract -->|has-field| n_NormalizedModelArtifact
```

Kildegrunnlag: f0629.

### Kontraktstruktur: PreparedFrameContract

```mermaid
flowchart LR
    n_PreparedFrameArtifact["PreparedFrameArtifact (field)"]
    n_PreparedFrameContract["PreparedFrameContract (contract)"]
    n_PreparedFrameContract -->|has-field| n_PreparedFrameArtifact
```

Kildegrunnlag: f0664.

### Kontraktstruktur: PresentationOutcome

```mermaid
flowchart LR
    n_PresentationOutcome["PresentationOutcome (contract)"]
    n_PresentationReady["PresentationReady (field)"]
    n_PresentationOutcome -->|has-field| n_PresentationReady
```

Kildegrunnlag: f0671.

### Kontraktstruktur: ReloadArguments

```mermaid
flowchart LR
    n_ReloadArguments["ReloadArguments (contract)"]
    n_ReloadSourceRevision["ReloadSourceRevision (field)"]
    n_ReloadSourceText["ReloadSourceText (field)"]
    n_ReloadArguments -->|has-field| n_ReloadSourceRevision
    n_ReloadArguments -->|has-field| n_ReloadSourceText
```

Kildegrunnlag: f0728, f0729.

### Kontraktstruktur: ReloadOutcome

```mermaid
flowchart LR
    n_ReloadDiagnostic["ReloadDiagnostic (field)"]
    n_ReloadOutcome["ReloadOutcome (contract)"]
    n_ReloadPublishedGeneration["ReloadPublishedGeneration (field)"]
    n_ReloadOutcome -->|has-field| n_ReloadDiagnostic
    n_ReloadOutcome -->|has-field| n_ReloadPublishedGeneration
```

Kildegrunnlag: f0757, f0758.

### Kontraktstruktur: SdlActionCallsProtocol

```mermaid
flowchart LR
    n_SdlActionCallsProtocol["SdlActionCallsProtocol (contract)"]
    n_SdlActionRequest["SdlActionRequest (message)"]
    n_SdlActionResult["SdlActionResult (message)"]
    n_SdlActionCallsProtocol -->|permits| n_SdlActionRequest
    n_SdlActionCallsProtocol -->|permits| n_SdlActionResult
```

Kildegrunnlag: f0813, f0814.

### Kontraktstruktur: SdlCompilationCallsProtocol

```mermaid
flowchart LR
    n_CompileSdlRequest["CompileSdlRequest (message)"]
    n_CompileSdlResult["CompileSdlResult (message)"]
    n_SdlCompilationCallsProtocol["SdlCompilationCallsProtocol (contract)"]
    n_SdlCompilationCallsProtocol -->|permits| n_CompileSdlRequest
    n_SdlCompilationCallsProtocol -->|permits| n_CompileSdlResult
```

Kildegrunnlag: f0822, f0823.

### Kontraktstruktur: SvgDocumentContract

```mermaid
flowchart LR
    n_SvgDocumentContract["SvgDocumentContract (contract)"]
    n_SvgDocumentText["SvgDocumentText (field)"]
    n_SvgDocumentContract -->|has-field| n_SvgDocumentText
```

Kildegrunnlag: f1081.

### Kontraktstruktur: SvgExportCallsProtocol

```mermaid
flowchart LR
    n_ExportSvgRequest["ExportSvgRequest (message)"]
    n_ExportSvgResult["ExportSvgResult (message)"]
    n_SvgExportCallsProtocol["SvgExportCallsProtocol (contract)"]
    n_SvgExportCallsProtocol -->|permits| n_ExportSvgRequest
    n_SvgExportCallsProtocol -->|permits| n_ExportSvgResult
```

Kildegrunnlag: f1086, f1087.

### Kontraktstruktur: TokenArtifactContract

```mermaid
flowchart LR
    n_TokenArtifact["TokenArtifact (field)"]
    n_TokenArtifactContract["TokenArtifactContract (contract)"]
    n_TokenArtifactContract -->|has-field| n_TokenArtifact
```

Kildegrunnlag: f1091.

### Kontraktstruktur: UiAstCallsProtocol

```mermaid
flowchart LR
    n_BuildUiAstRequest["BuildUiAstRequest (message)"]
    n_BuildUiAstResult["BuildUiAstResult (message)"]
    n_UiAstCallsProtocol["UiAstCallsProtocol (contract)"]
    n_UiAstCallsProtocol -->|permits| n_BuildUiAstRequest
    n_UiAstCallsProtocol -->|permits| n_BuildUiAstResult
```

Kildegrunnlag: f1123, f1124.

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

Kildegrunnlag: f1142, f1143, f1144.

### Kontraktstruktur: UiDomainActionsProtocol

```mermaid
flowchart LR
    n_BoundActionRequest["BoundActionRequest (message)"]
    n_BoundActionResult["BoundActionResult (message)"]
    n_UiDomainActionsProtocol["UiDomainActionsProtocol (contract)"]
    n_UiDomainActionsProtocol -->|permits| n_BoundActionRequest
    n_UiDomainActionsProtocol -->|permits| n_BoundActionResult
```

Kildegrunnlag: f1147, f1148.

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

Kildegrunnlag: f1149, f1150, f1152.

### Kontraktstruktur: UiGenerationEventsProtocol

```mermaid
flowchart LR
    n_UiGenerationEventsProtocol["UiGenerationEventsProtocol (contract)"]
    n_UiGenerationNotices["UiGenerationNotices (datagram)"]
    n_UiGenerationEventsProtocol -->|permits| n_UiGenerationNotices
```

Kildegrunnlag: f1155.

### Kontraktstruktur: UiNormalizationCallsProtocol

```mermaid
flowchart LR
    n_NormalizeUiRequest["NormalizeUiRequest (message)"]
    n_NormalizeUiResult["NormalizeUiResult (message)"]
    n_UiNormalizationCallsProtocol["UiNormalizationCallsProtocol (contract)"]
    n_UiNormalizationCallsProtocol -->|permits| n_NormalizeUiRequest
    n_UiNormalizationCallsProtocol -->|permits| n_NormalizeUiResult
```

Kildegrunnlag: f1183, f1184.

### Kontraktstruktur: UiSessionRecord

```mermaid
flowchart LR
    n_SessionDraft["SessionDraft (field)"]
    n_SessionGeneration["SessionGeneration (field)"]
    n_UiSessionRecord["UiSessionRecord (contract)"]
    n_UiSessionRecord -->|has-field| n_SessionDraft
    n_UiSessionRecord -->|has-field| n_SessionGeneration
```

Kildegrunnlag: f1186, f1187.

### Kontraktstruktur: UiTokenizationCallsProtocol

```mermaid
flowchart LR
    n_TokenizeUiRequest["TokenizeUiRequest (message)"]
    n_TokenizeUiResult["TokenizeUiResult (message)"]
    n_UiTokenizationCallsProtocol["UiTokenizationCallsProtocol (contract)"]
    n_UiTokenizationCallsProtocol -->|permits| n_TokenizeUiRequest
    n_UiTokenizationCallsProtocol -->|permits| n_TokenizeUiResult
```

Kildegrunnlag: f1197, f1198.

### Kontraktstruktur: UiValidationCallsProtocol

```mermaid
flowchart LR
    n_UiValidationCallsProtocol["UiValidationCallsProtocol (contract)"]
    n_ValidateUiRequest["ValidateUiRequest (message)"]
    n_ValidateUiResult["ValidateUiResult (message)"]
    n_UiValidationCallsProtocol -->|permits| n_ValidateUiRequest
    n_UiValidationCallsProtocol -->|permits| n_ValidateUiResult
```

Kildegrunnlag: f1201, f1202.

### Kontraktstruktur: ValidationOutcomeContract

```mermaid
flowchart LR
    n_ModelDiagnostics["ModelDiagnostics (field)"]
    n_ModelIsValid["ModelIsValid (field)"]
    n_ValidationOutcomeContract["ValidationOutcomeContract (contract)"]
    n_ValidationOutcomeContract -->|has-field| n_ModelDiagnostics
    n_ValidationOutcomeContract -->|has-field| n_ModelIsValid
```

Kildegrunnlag: f1240, f1241.

### Kontraktstruktur: ViewBundleContract

```mermaid
flowchart LR
    n_ViewBundleBytes["ViewBundleBytes (field)"]
    n_ViewBundleContract["ViewBundleContract (contract)"]
    n_ViewBundleRevision["ViewBundleRevision (field)"]
    n_ViewBundleContract -->|has-field| n_ViewBundleBytes
    n_ViewBundleContract -->|has-field| n_ViewBundleRevision
```

Kildegrunnlag: f1252, f1253.

### Kontraktstruktur: ViewDisplayCallsProtocol

```mermaid
flowchart LR
    n_DisplayViewRequest["DisplayViewRequest (message)"]
    n_ViewDisplayCallsProtocol["ViewDisplayCallsProtocol (contract)"]
    n_ViewDisplayResult["ViewDisplayResult (message)"]
    n_ViewTargetUnavailable["ViewTargetUnavailable (message)"]
    n_ViewDisplayCallsProtocol -->|permits| n_DisplayViewRequest
    n_ViewDisplayCallsProtocol -->|permits| n_ViewDisplayResult
    n_ViewDisplayCallsProtocol -->|permits| n_ViewTargetUnavailable
```

Kildegrunnlag: f1263, f1264, f1265.

### Kontraktstruktur: ViewFailureContract

```mermaid
flowchart LR
    n_FailedViewRequestId["FailedViewRequestId (field)"]
    n_ViewFailureCode["ViewFailureCode (field)"]
    n_ViewFailureContract["ViewFailureContract (contract)"]
    n_ViewFailureDiagnostic["ViewFailureDiagnostic (field)"]
    n_ViewFailureContract -->|has-field| n_FailedViewRequestId
    n_ViewFailureContract -->|has-field| n_ViewFailureCode
    n_ViewFailureContract -->|has-field| n_ViewFailureDiagnostic
```

Kildegrunnlag: f1274, f1275, f1276.

### Kontraktstruktur: ViewLaunchCallsProtocol

```mermaid
flowchart LR
    n_LaunchViewRequest["LaunchViewRequest (message)"]
    n_ViewLaunchCallsProtocol["ViewLaunchCallsProtocol (contract)"]
    n_ViewLaunchResult["ViewLaunchResult (message)"]
    n_ViewLaunchCallsProtocol -->|permits| n_LaunchViewRequest
    n_ViewLaunchCallsProtocol -->|permits| n_ViewLaunchResult
```

Kildegrunnlag: f1281, f1282.

### Kontraktstruktur: ViewNavigationCallsProtocol

```mermaid
flowchart LR
    n_SelectViewRequest["SelectViewRequest (message)"]
    n_ViewNavigationCallsProtocol["ViewNavigationCallsProtocol (contract)"]
    n_ViewOpenedResult["ViewOpenedResult (message)"]
    n_ViewRequestRejected["ViewRequestRejected (message)"]
    n_ViewNavigationCallsProtocol -->|permits| n_SelectViewRequest
    n_ViewNavigationCallsProtocol -->|permits| n_ViewOpenedResult
    n_ViewNavigationCallsProtocol -->|permits| n_ViewRequestRejected
```

Kildegrunnlag: f1292, f1293, f1294.

### Kontraktstruktur: ViewOpenContract

```mermaid
flowchart LR
    n_OpenViewConsumerId["OpenViewConsumerId (field)"]
    n_OpenViewEntryPath["OpenViewEntryPath (field)"]
    n_OpenViewLeaseId["OpenViewLeaseId (field)"]
    n_OpenViewPaneId["OpenViewPaneId (field)"]
    n_OpenViewRequestId["OpenViewRequestId (field)"]
    n_OpenViewRevision["OpenViewRevision (field)"]
    n_OpenViewWindowId["OpenViewWindowId (field)"]
    n_ViewOpenContract["ViewOpenContract (contract)"]
    n_ViewOpenContract -->|has-field| n_OpenViewConsumerId
    n_ViewOpenContract -->|has-field| n_OpenViewEntryPath
    n_ViewOpenContract -->|has-field| n_OpenViewLeaseId
    n_ViewOpenContract -->|has-field| n_OpenViewPaneId
    n_ViewOpenContract -->|has-field| n_OpenViewRequestId
    n_ViewOpenContract -->|has-field| n_OpenViewRevision
    n_ViewOpenContract -->|has-field| n_OpenViewWindowId
```

Kildegrunnlag: f1296, f1297, f1298, f1299, f1300, f1301, f1302.

### Kontraktstruktur: ViewOpenedContract

```mermaid
flowchart LR
    n_OpenedViewRequestId["OpenedViewRequestId (field)"]
    n_OpenedViewRevision["OpenedViewRevision (field)"]
    n_ViewOpenedContract["ViewOpenedContract (contract)"]
    n_ViewOpenedContract -->|has-field| n_OpenedViewRequestId
    n_ViewOpenedContract -->|has-field| n_OpenedViewRevision
```

Kildegrunnlag: f1304, f1305.

### Kontraktstruktur: ViewProjectionCallsProtocol

```mermaid
flowchart LR
    n_ProjectViewRequest["ProjectViewRequest (message)"]
    n_ViewBundleResult["ViewBundleResult (message)"]
    n_ViewProjectionCallsProtocol["ViewProjectionCallsProtocol (contract)"]
    n_ViewProjectionRejected["ViewProjectionRejected (message)"]
    n_ViewProjectionCallsProtocol -->|permits| n_ProjectViewRequest
    n_ViewProjectionCallsProtocol -->|permits| n_ViewBundleResult
    n_ViewProjectionCallsProtocol -->|permits| n_ViewProjectionRejected
```

Kildegrunnlag: f1315, f1316, f1317.

### Kontraktstruktur: ViewPublicationCallsProtocol

```mermaid
flowchart LR
    n_PublishViewRequest["PublishViewRequest (message)"]
    n_ViewPublicationCallsProtocol["ViewPublicationCallsProtocol (contract)"]
    n_ViewReferenceResult["ViewReferenceResult (message)"]
    n_ViewPublicationCallsProtocol -->|permits| n_PublishViewRequest
    n_ViewPublicationCallsProtocol -->|permits| n_ViewReferenceResult
```

Kildegrunnlag: f1332, f1333.

### Kontraktstruktur: ViewReferenceContract

```mermaid
flowchart LR
    n_ViewEntryPath["ViewEntryPath (field)"]
    n_ViewLeaseId["ViewLeaseId (field)"]
    n_ViewManifestPath["ViewManifestPath (field)"]
    n_ViewPublishedRevision["ViewPublishedRevision (field)"]
    n_ViewReferenceContract["ViewReferenceContract (contract)"]
    n_ViewReferenceContract -->|has-field| n_ViewEntryPath
    n_ViewReferenceContract -->|has-field| n_ViewLeaseId
    n_ViewReferenceContract -->|has-field| n_ViewManifestPath
    n_ViewReferenceContract -->|has-field| n_ViewPublishedRevision
```

Kildegrunnlag: f1337, f1338, f1339, f1340.

### Kontraktstruktur: ViewSelectionContract

```mermaid
flowchart LR
    n_ViewConsumerId["ViewConsumerId (field)"]
    n_ViewPaneId["ViewPaneId (field)"]
    n_ViewProjectId["ViewProjectId (field)"]
    n_ViewRequestId["ViewRequestId (field)"]
    n_ViewSelectionContract["ViewSelectionContract (contract)"]
    n_ViewSelector["ViewSelector (field)"]
    n_ViewSourceRevision["ViewSourceRevision (field)"]
    n_ViewWindowId["ViewWindowId (field)"]
    n_ViewSelectionContract -->|has-field| n_ViewConsumerId
    n_ViewSelectionContract -->|has-field| n_ViewPaneId
    n_ViewSelectionContract -->|has-field| n_ViewProjectId
    n_ViewSelectionContract -->|has-field| n_ViewRequestId
    n_ViewSelectionContract -->|has-field| n_ViewSelector
    n_ViewSelectionContract -->|has-field| n_ViewSourceRevision
    n_ViewSelectionContract -->|has-field| n_ViewWindowId
```

Kildegrunnlag: f1350, f1351, f1352, f1353, f1354, f1355, f1356.


## VP10 — Datagram-koding og packet

### Packet: UiGenerationWire / UiGenerationChanged — big-endian, most-significant-first

```mermaid
packet
    0-15: "NoticeVersion"
    16-79: "NoticeGeneration"
```

Kildegrunnlag: f0630, f0631, f0632, f0633, f1149, f1150, f1151, f1152, f1158, f1159, f1160, f1161, f1162.


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
| DesignSourceRecord has completeness = closed. | f0156 |
| FailedViewRequestId has presence = required. | f0236 |
| FailedViewRequestId has value-type = text. | f0237 |
| FramePresentationCallsProtocol has completeness = closed. | f0239 |
| GeneratedGoContract has completeness = closed. | f0477 |
| GeneratedGoSources has presence = required. | f0479 |
| GeneratedGoSources has value-type = bytes. | f0480 |
| GoBuildCallsProtocol has completeness = closed. | f0482 |
| GoDomainCallsProtocol has completeness = closed. | f0502 |
| GoGenerationCallsProtocol has completeness = closed. | f0512 |
| LayoutArguments has completeness = closed. | f0548 |
| LayoutCallsProtocol has completeness = closed. | f0553 |
| LayoutModelArtifact has presence = required. | f0556 |
| LayoutModelArtifact has value-type = bytes. | f0557 |
| LayoutViewportHeight has presence = required. | f0563 |
| LayoutViewportHeight has value-type = decimal. | f0564 |
| LayoutViewportWidth has presence = required. | f0565 |
| LayoutViewportWidth has value-type = decimal. | f0566 |
| ModelDiagnostics has presence = optional. | f0584 |
| ModelDiagnostics has value-type = text. | f0585 |
| ModelIsValid has presence = required. | f0586 |
| ModelIsValid has value-type = boolean. | f0587 |
| ModelReloadCallsProtocol has completeness = closed. | f0589 |
| NativeBuildContract has completeness = closed. | f0593 |
| NativeBuildDiagnostics has presence = optional. | f0596 |
| NativeBuildDiagnostics has value-type = text. | f0597 |
| NativeBuildSucceeded has presence = required. | f0598 |
| NativeBuildSucceeded has value-type = boolean. | f0599 |
| NativeUiActionsProtocol has completeness = closed. | f0613 |
| NormalizedModelArtifact has presence = required. | f0626 |
| NormalizedModelArtifact has value-type = bytes. | f0627 |
| NormalizedModelContract has completeness = closed. | f0628 |
| NoticeGeneration has presence = required. | f0630 |
| NoticeGeneration has value-type = unsigned. | f0631 |
| NoticeVersion has presence = required. | f0632 |
| NoticeVersion has value-type = unsigned. | f0633 |
| OpenViewConsumerId has presence = required. | f0637 |
| OpenViewConsumerId has value-type = text. | f0638 |
| OpenViewEntryPath has presence = required. | f0639 |
| OpenViewEntryPath has value-type = text. | f0640 |
| OpenViewLeaseId has presence = required. | f0641 |
| OpenViewLeaseId has value-type = text. | f0642 |
| OpenViewPaneId has presence = required. | f0643 |
| OpenViewPaneId has value-type = text. | f0644 |
| OpenViewRequestId has presence = required. | f0645 |
| OpenViewRequestId has value-type = text. | f0646 |
| OpenViewRevision has presence = required. | f0647 |
| OpenViewRevision has value-type = text. | f0648 |
| OpenViewWindowId has presence = required. | f0649 |
| OpenViewWindowId has value-type = text. | f0650 |
| OpenedViewRequestId has presence = required. | f0651 |
| OpenedViewRequestId has value-type = text. | f0652 |
| OpenedViewRevision has presence = required. | f0653 |
| OpenedViewRevision has value-type = text. | f0654 |
| PreparedFrameArtifact has presence = required. | f0661 |
| PreparedFrameArtifact has value-type = bytes. | f0662 |
| PreparedFrameContract has completeness = closed. | f0663 |
| PresentationOutcome has completeness = closed. | f0670 |
| PresentationReady has presence = required. | f0672 |
| PresentationReady has value-type = boolean. | f0673 |
| ReloadArguments has completeness = closed. | f0727 |
| ReloadDiagnostic has presence = optional. | f0754 |
| ReloadDiagnostic has value-type = text. | f0755 |
| ReloadOutcome has completeness = closed. | f0756 |
| ReloadPublishedGeneration has presence = optional. | f0762 |
| ReloadPublishedGeneration has value-type = unsigned. | f0763 |
| ReloadSourceRevision has presence = required. | f0769 |
| ReloadSourceRevision has value-type = unsigned. | f0770 |
| ReloadSourceText has presence = required. | f0771 |
| ReloadSourceText has value-type = text. | f0772 |
| SdlActionCallsProtocol has completeness = closed. | f0812 |
| SdlCompilationCallsProtocol has completeness = closed. | f0821 |
| SessionDraft has presence = optional. | f1042 |
| SessionDraft has value-type = text. | f1043 |
| SessionGeneration has presence = required. | f1044 |
| SessionGeneration has value-type = unsigned. | f1045 |
| SourceDocumentRevision has presence = required. | f1048 |
| SourceDocumentRevision has value-type = unsigned. | f1049 |
| SourceDocumentText has presence = required. | f1050 |
| SourceDocumentText has value-type = text. | f1051 |
| SvgDocumentContract has completeness = closed. | f1080 |
| SvgDocumentText has presence = required. | f1082 |
| SvgDocumentText has value-type = text. | f1083 |
| SvgExportCallsProtocol has completeness = closed. | f1085 |
| TokenArtifact has presence = required. | f1088 |
| TokenArtifact has value-type = bytes. | f1089 |
| TokenArtifactContract has completeness = closed. | f1090 |
| UiAstCallsProtocol has completeness = closed. | f1122 |
| UiCompilationCallsProtocol has completeness = closed. | f1141 |
| UiDomainActionsProtocol has completeness = closed. | f1146 |
| UiGenerationContract has completeness = closed. | f1151 |
| UiGenerationEventsProtocol has completeness = closed. | f1154 |
| UiGenerationWire has bit-order = most-significant-first. | f1159 |
| UiGenerationWire has byte-order = big-endian. | f1160 |
| UiNormalizationCallsProtocol has completeness = closed. | f1182 |
| UiSessionRecord has completeness = closed. | f1185 |
| UiTokenizationCallsProtocol has completeness = closed. | f1196 |
| UiValidationCallsProtocol has completeness = closed. | f1200 |
| ValidationOutcomeContract has completeness = closed. | f1239 |
| ViewBundleBytes has presence = required. | f1249 |
| ViewBundleBytes has value-type = bytes. | f1250 |
| ViewBundleContract has completeness = closed. | f1251 |
| ViewBundleRevision has presence = required. | f1257 |
| ViewBundleRevision has value-type = text. | f1258 |
| ViewConsumerId has presence = required. | f1259 |
| ViewConsumerId has value-type = text. | f1260 |
| ViewDisplayCallsProtocol has completeness = closed. | f1262 |
| ViewEntryPath has presence = required. | f1269 |
| ViewEntryPath has value-type = text. | f1270 |
| ViewFailureCode has presence = required. | f1271 |
| ViewFailureCode has value-type = text. | f1272 |
| ViewFailureContract has completeness = closed. | f1273 |
| ViewFailureDiagnostic has presence = required. | f1277 |
| ViewFailureDiagnostic has value-type = text. | f1278 |
| ViewLaunchCallsProtocol has completeness = closed. | f1280 |
| ViewLeaseId has presence = required. | f1286 |
| ViewLeaseId has value-type = text. | f1287 |
| ViewManifestPath has presence = required. | f1288 |
| ViewManifestPath has value-type = text. | f1289 |
| ViewNavigationCallsProtocol has completeness = closed. | f1291 |
| ViewOpenContract has completeness = closed. | f1295 |
| ViewOpenedContract has completeness = closed. | f1303 |
| ViewPaneId has presence = required. | f1309 |
| ViewPaneId has value-type = text. | f1310 |
| ViewProjectId has presence = required. | f1311 |
| ViewProjectId has value-type = text. | f1312 |
| ViewProjectionCallsProtocol has completeness = closed. | f1314 |
| ViewPublicationCallsProtocol has completeness = closed. | f1331 |
| ViewPublishedRevision has presence = required. | f1334 |
| ViewPublishedRevision has value-type = text. | f1335 |
| ViewReferenceContract has completeness = closed. | f1336 |
| ViewRequestId has presence = required. | f1344 |
| ViewRequestId has value-type = text. | f1345 |
| ViewSelectionContract has completeness = closed. | f1349 |
| ViewSelector has presence = required. | f1357 |
| ViewSelector has value-type = text. | f1358 |
| ViewSourceRevision has presence = required. | f1362 |
| ViewSourceRevision has value-type = text. | f1363 |
| ViewWindowId has presence = required. | f1367 |
| ViewWindowId has value-type = text. | f1368 |

### VP09 — projeksjonsansvar

| Functionality | Dataset | Datagram-familie | Faktum |
| --- | --- | --- | --- |
| ProjectUiGeneration | UiSessionState | UiGenerationNotices | f0700 |


### VP08 — avledet MessageSet per Channel og modus

Generert fra permits og deltakelse, ikke en separat authored modell. Tom deltakelse er et hull.

| Channel | Mode | Message / Datagram | Sender | Receiver | Kilde-ID-er |
| --- | --- | --- | --- | --- | --- |
| FramePresentationCalls | UiPreview | PresentFrameRequest | FyneHost | FyneBackend | f0238, f0240, f0249, f0267, f0666 |
| FramePresentationCalls | UiPreview | PresentFrameResult | FyneBackend | FyneHost | f0238, f0241, f0250, f0266, f0669 |
| GoBuildCalls | NativeBuild | BuildGoRequest | CommandLineHost | GoBuildRunner | f0056, f0096, f0481, f0483, f0490 |
| GoBuildCalls | NativeBuild | BuildGoResult | GoBuildRunner | CommandLineHost | f0059, f0095, f0481, f0484, f0491 |
| GoDomainCalls | BoundExecution | DomainActionRequest | SdlDispatcher | GoDomainImplementation | f0206, f0501, f0503, f0507, f0832 |
| GoDomainCalls | BoundExecution | DomainActionResult | GoDomainImplementation | SdlDispatcher | f0209, f0501, f0504, f0509, f0831 |
| GoDomainCalls | UiPreview | DomainActionRequest | FyneHost | GoDomainImplementation | f0206, f0269, f0501, f0503, f0508 |
| GoDomainCalls | UiPreview | DomainActionResult | GoDomainImplementation | FyneHost | f0209, f0268, f0501, f0504, f0510 |
| GoGenerationCalls | NativeBuild | GenerateGoRequest | CommandLineHost | GoCodeGenerator | f0098, f0468, f0499, f0511, f0513 |
| GoGenerationCalls | NativeBuild | GenerateGoResult | GoCodeGenerator | CommandLineHost | f0097, f0471, f0500, f0511, f0514 |
| LayoutCalls | StaticExport | LayoutRequest | CommandLineHost | SduiLayout | f0100, f0552, f0554, f0559, f0966 |
| LayoutCalls | StaticExport | LayoutResult | SduiLayout | CommandLineHost | f0099, f0552, f0555, f0562, f0968 |
| LayoutCalls | UiPreview | LayoutRequest | FyneHost | SduiLayout | f0271, f0552, f0554, f0559, f0967 |
| LayoutCalls | UiPreview | LayoutResult | SduiLayout | FyneHost | f0270, f0552, f0555, f0562, f0969 |
| ModelReloadCalls | BoundLiveEditing | ReloadPublished | ReloadCoordinator | SourceWatcher | f0588, f0590, f0745, f0761, f1061 |
| ModelReloadCalls | BoundLiveEditing | ReloadRejected | ReloadCoordinator | SourceWatcher | f0588, f0591, f0747, f0766, f1063 |
| ModelReloadCalls | BoundLiveEditing | ReloadRequest | SourceWatcher | ReloadCoordinator | f0588, f0592, f0743, f0768, f1065 |
| ModelReloadCalls | LiveEditing | ReloadPublished | ReloadCoordinator | SourceWatcher | f0588, f0590, f0746, f0761, f1062 |
| ModelReloadCalls | LiveEditing | ReloadRejected | ReloadCoordinator | SourceWatcher | f0588, f0591, f0748, f0766, f1064 |
| ModelReloadCalls | LiveEditing | ReloadRequest | SourceWatcher | ReloadCoordinator | f0588, f0592, f0744, f0768, f1066 |
| NativeUiActions | BoundExecution | UiActionRejected | SduiDispatcher | FyneBackend | f0251, f0612, f0614, f0923, f1115 |
| NativeUiActions | BoundExecution | UiActionRequest | FyneBackend | SduiDispatcher | f0253, f0612, f0615, f0922, f1117 |
| NativeUiActions | BoundExecution | UiActionResult | SduiDispatcher | FyneBackend | f0252, f0612, f0616, f0924, f1120 |
| SdlActionCalls | BoundExecution | SdlActionRequest | SdlUiBindingAdapter | SdlDispatcher | f0811, f0813, f0816, f0833, f0898 |
| SdlActionCalls | BoundExecution | SdlActionResult | SdlDispatcher | SdlUiBindingAdapter | f0811, f0814, f0819, f0834, f0897 |
| SdlCompilationCalls | BoundLiveEditing | CompileSdlRequest | ReloadCoordinator | SdlFrontend | f0107, f0750, f0820, f0822, f0847 |
| SdlCompilationCalls | BoundLiveEditing | CompileSdlResult | SdlFrontend | ReloadCoordinator | f0110, f0749, f0820, f0823, f0848 |
| SvgExportCalls | StaticExport | ExportSvgRequest | CommandLineHost | SduiPresentation | f0102, f0223, f0997, f1084, f1086 |
| SvgExportCalls | StaticExport | ExportSvgResult | SduiPresentation | CommandLineHost | f0101, f0226, f0998, f1084, f1087 |
| UiAstCalls | SourceInspection | BuildUiAstRequest | SduiFrontend | SduiParser | f0075, f0936, f0990, f1121, f1123 |
| UiAstCalls | SourceInspection | BuildUiAstResult | SduiParser | SduiFrontend | f0078, f0935, f0991, f1121, f1124 |
| UiCompilationCalls | LiveEditing | CompileUiRejected | SduiFrontend | ReloadCoordinator | f0113, f0751, f0939, f1140, f1142 |
| UiCompilationCalls | LiveEditing | CompileUiRequest | ReloadCoordinator | SduiFrontend | f0115, f0753, f0937, f1140, f1143 |
| UiCompilationCalls | LiveEditing | CompileUiResult | SduiFrontend | ReloadCoordinator | f0118, f0752, f0941, f1140, f1144 |
| UiCompilationCalls | SourceInspection | CompileUiRejected | SduiFrontend | CommandLineHost | f0103, f0113, f0940, f1140, f1142 |
| UiCompilationCalls | SourceInspection | CompileUiRequest | CommandLineHost | SduiFrontend | f0105, f0115, f0938, f1140, f1143 |
| UiCompilationCalls | SourceInspection | CompileUiResult | SduiFrontend | CommandLineHost | f0104, f0118, f0942, f1140, f1144 |
| UiDomainActions | BoundExecution | BoundActionRequest | SduiDispatcher | SdlUiBindingAdapter | f0046, f0899, f0926, f1145, f1147 |
| UiDomainActions | BoundExecution | BoundActionResult | SdlUiBindingAdapter | SduiDispatcher | f0049, f0900, f0925, f1145, f1148 |
| UiGenerationEvents | BoundExecution | UiGenerationNotices | SduiInstanceStore | FyneBackend | f0254, f0955, f1153, f1155, f1157 |
| UiGenerationEvents | LiveEditing | UiGenerationNotices | SduiInstanceStore | FyneBackend | f0255, f0956, f1153, f1155, f1157 |
| UiNormalizationCalls | SourceInspection | NormalizeUiRequest | SduiFrontend | SduiNormalizer | f0622, f0944, f0986, f1181, f1183 |
| UiNormalizationCalls | SourceInspection | NormalizeUiResult | SduiNormalizer | SduiFrontend | f0625, f0943, f0987, f1181, f1184 |
| UiTokenizationCalls | SourceInspection | TokenizeUiRequest | SduiFrontend | SduiLexer | f0946, f0972, f1099, f1195, f1197 |
| UiTokenizationCalls | SourceInspection | TokenizeUiResult | SduiLexer | SduiFrontend | f0945, f0973, f1102, f1195, f1198 |
| UiValidationCalls | SourceInspection | ValidateUiRequest | SduiFrontend | SduiValidator | f0948, f1019, f1199, f1201, f1227 |
| UiValidationCalls | SourceInspection | ValidateUiResult | SduiValidator | SduiFrontend | f0947, f1020, f1199, f1202, f1230 |
| ViewDisplayCalls | DocumentBrowsing | DisplayViewRequest | ViewerLaunchAdapter | XfmdDocumentHost | f0190, f1261, f1263, f1373, f1381 |
| ViewDisplayCalls | DocumentBrowsing | ViewDisplayResult | XfmdDocumentHost | ViewerLaunchAdapter | f1261, f1264, f1268, f1371, f1382 |
| ViewDisplayCalls | DocumentBrowsing | ViewTargetUnavailable | XfmdDocumentHost | ViewerLaunchAdapter | f1261, f1265, f1366, f1372, f1383 |
| ViewLaunchCalls | DocumentBrowsing | LaunchViewRequest | DocumentBroker | ViewerLaunchAdapter | f0196, f0547, f1279, f1281, f1374 |
| ViewLaunchCalls | DocumentBrowsing | ViewLaunchResult | ViewerLaunchAdapter | DocumentBroker | f0195, f1279, f1282, f1285, f1375 |
| ViewNavigationCalls | DocumentBrowsing | SelectViewRequest | XfmdDocumentHost | DocumentBroker | f0197, f1022, f1290, f1292, f1386 |
| ViewNavigationCalls | DocumentBrowsing | ViewOpenedResult | DocumentBroker | XfmdDocumentHost | f0198, f1290, f1293, f1308, f1384 |
| ViewNavigationCalls | DocumentBrowsing | ViewRequestRejected | DocumentBroker | XfmdDocumentHost | f0199, f1290, f1294, f1348, f1385 |
| ViewProjectionCalls | DocumentBrowsing | ProjectViewRequest | DocumentBroker | SdlViewpointGenerator | f0202, f0702, f0912, f1313, f1315 |
| ViewProjectionCalls | DocumentBrowsing | ViewBundleResult | SdlViewpointGenerator | DocumentBroker | f0200, f0913, f1256, f1313, f1316 |
| ViewProjectionCalls | DocumentBrowsing | ViewProjectionRejected | SdlViewpointGenerator | DocumentBroker | f0201, f0914, f1313, f1317, f1329 |
| ViewPublicationCalls | DocumentBrowsing | PublishViewRequest | DocumentBroker | ViewArtifactStore | f0204, f0714, f1247, f1330, f1332 |
| ViewPublicationCalls | DocumentBrowsing | ViewReferenceResult | ViewArtifactStore | DocumentBroker | f0203, f1248, f1330, f1333, f1343 |

## VP04 — Grensesnittbruk

Tabellen dekker alle consumes-fakta. Portnavn er ikke Channel-kontrakter eller tilbyderkoblinger.

| Unit / Container | Interface | Faktum | Kildelinje |
| --- | --- | --- | --- |
| CommandLineHost | ExportSinkPort | f0088 | 543 |
| CommandLineHost | PreparedFramePort | f0089 | 544 |
| CommandLineHost | SduiFrontendPort | f0090 | 545 |
| CommandLineHost | SourceSnapshotPort | f0091 | 546 |
| DiagramProvider | DiagramEnginePort | f0179 | 634 |
| DiagramProvider | ResourcePort | f0180 | 635 |
| DomainStateMigrator | DomainStatePort | f0210 | 665 |
| DomainStateMigrator | SdlModelPort | f0211 | 666 |
| FyneBackend | PreparedFramePort | f0242 | 697 |
| FyneBackend | UiSessionPort | f0243 | 698 |
| FyneHost | DomainBindingPort | f0256 | 711 |
| FyneHost | ReloadPort | f0257 | 712 |
| FyneHost | SdlFrontendPort | f0258 | 713 |
| FyneHost | SduiFrontendPort | f0259 | 714 |
| FyneHost | SourceSnapshotPort | f0260 | 715 |
| FyneHost | UiSessionPort | f0261 | 716 |
| FyneHost | WidgetBackendPort | f0262 | 717 |
| GoBuildRunner | BuildToolPort | f0485 | 940 |
| GoBuildRunner | GeneratedArtifactPort | f0486 | 941 |
| GoCodeGenerator | ExecutionProfilePort | f0492 | 947 |
| GoCodeGenerator | SdlModelPort | f0493 | 948 |
| GoCodeGenerator | SduiModelPort | f0494 | 949 |
| MarkdownProvider | DiagramPort | f0572 | 1027 |
| MarkdownProvider | MeasurementPort | f0573 | 1028 |
| MarkdownProvider | ResourcePort | f0574 | 1029 |
| ReloadCoordinator | BindingReloadPort | f0731 | 1186 |
| ReloadCoordinator | DiagnosticPort | f0732 | 1187 |
| ReloadCoordinator | SdlFrontendPort | f0733 | 1188 |
| ReloadCoordinator | SdlReloadPort | f0734 | 1189 |
| ReloadCoordinator | SduiFrontendPort | f0735 | 1190 |
| ReloadCoordinator | SourceSnapshotPort | f0736 | 1191 |
| ReloadCoordinator | UiReloadPort | f0737 | 1192 |
| SdlDispatcher | DomainFunctionPort | f0824 | 1279 |
| SdlDispatcher | DomainStatePort | f0825 | 1280 |
| SdlExecutionGate | DiagnosticPort | f0835 | 1290 |
| SdlExecutionGate | SdlModelPort | f0836 | 1291 |
| SdlFrontend | DiagnosticPort | f0839 | 1294 |
| SdlFrontend | SourceSnapshotPort | f0840 | 1295 |
| SdlFunctionRegistry | DomainFunctionPort | f0849 | 1304 |
| SdlRuntime | DomainFunctionPort | f0874 | 1329 |
| SdlRuntime | SdlModelPort | f0875 | 1330 |
| SdlUiBindingAdapter | DiagnosticPort | f0888 | 1343 |
| SdlUiBindingAdapter | SdlExecutionPort | f0889 | 1344 |
| SdlUiBindingAdapter | UiSessionPort | f0890 | 1345 |
| SduiDispatcher | DomainBindingPort | f0915 | 1370 |
| SduiDispatcher | UiStatePort | f0916 | 1371 |
| SduiFrontend | DiagnosticPort | f0927 | 1382 |
| SduiFrontend | SourceSnapshotPort | f0928 | 1383 |
| SduiLayout | ContentProviderPort | f0957 | 1412 |
| SduiLayout | MeasurementPort | f0958 | 1413 |
| SduiLayout | UiSnapshotPort | f0959 | 1414 |
| SduiPresentation | PreparedFramePort | f0992 | 1447 |
| SduiPropertyStore | UiStatePort | f0999 | 1454 |
| SduiRuntime | DomainBindingPort | f1004 | 1459 |
| SduiRuntime | SduiModelPort | f1005 | 1460 |
| SourceLoader | SourceInputPort | f1052 | 1507 |
| SourceWatcher | FileChangePort | f1057 | 1512 |
| UiStateReconciler | SduiModelPort | f1189 | 1644 |
| UiStateReconciler | UiStatePort | f1190 | 1645 |

## VP11 — Egenskaper og fullstendig faktaregister

Registeret inkluderer alle fakta, også de som ikke har en egen tegning.

| ID | Utsagn | Kildelinje |
| --- | --- | --- |
| f0000 | ActionArguments has completeness = closed. | 455 |
| f0001 | ActionArguments has-field ActionGeneration. | 456 |
| f0002 | ActionArguments has-field ActionInputText. | 457 |
| f0003 | ActionArguments has-field ActionSymbol. | 458 |
| f0004 | ActionGeneration has presence = required. | 459 |
| f0005 | ActionGeneration has value-type = unsigned. | 460 |
| f0006 | ActionInputText has presence = required. | 461 |
| f0007 | ActionInputText has value-type = text. | 462 |
| f0008 | ActionOutcome has completeness = closed. | 463 |
| f0009 | ActionOutcome has-field ActionOutputText. | 464 |
| f0010 | ActionOutcome has-field ActionStatusCode. | 465 |
| f0011 | ActionOutputText has presence = required. | 466 |
| f0012 | ActionOutputText has value-type = text. | 467 |
| f0013 | ActionStatusCode has presence = required. | 468 |
| f0014 | ActionStatusCode has value-type = unsigned. | 469 |
| f0015 | ActionSymbol has presence = required. | 470 |
| f0016 | ActionSymbol has value-type = text. | 471 |
| f0017 | AllocateGeometry allocated-to FyneHost in mode UiPreview. | 472 |
| f0018 | AllocateGeometry contributes-to InteractiveUiPreview. | 473 |
| f0019 | AllocateGeometry realizes MeasuredPresentation. | 474 |
| f0020 | ApplyPropertyBatch has state-retention = stateful. | 475 |
| f0021 | ApplyPropertyBatch realizes InteractiveSession. | 476 |
| f0022 | AstArtifact has presence = required. | 477 |
| f0023 | AstArtifact has value-type = bytes. | 478 |
| f0024 | AstArtifactContract has completeness = closed. | 479 |
| f0025 | AstArtifactContract has-field AstArtifact. | 480 |
| f0026 | BoundActionAccepted exercises TryDomainInteraction. | 481 |
| f0027 | BoundActionAccepted has completeness = closed. | 482 |
| f0028 | BoundActionAccepted illustrates G4M3UiDomainBinding. | 483 |
| f0029 | BoundActionAccepted runs-in BoundExecution. | 484 |
| f0030 | BoundActionAccepted step 1 sends UiActionRequest from FyneBackend to SduiDispatcher via NativeUiActions. | 485 |
| f0031 | BoundActionAccepted step 2 sends BoundActionRequest from SduiDispatcher to SdlUiBindingAdapter via UiDomainActions. | 486 |
| f0032 | BoundActionAccepted step 3 sends SdlActionRequest from SdlUiBindingAdapter to SdlDispatcher via SdlActionCalls. | 487 |
| f0033 | BoundActionAccepted step 4 sends DomainActionRequest from SdlDispatcher to GoDomainImplementation via GoDomainCalls. | 488 |
| f0034 | BoundActionAccepted step 5 sends DomainActionResult from GoDomainImplementation to SdlDispatcher via GoDomainCalls reply-to 4. | 489 |
| f0035 | BoundActionAccepted step 6 sends SdlActionResult from SdlDispatcher to SdlUiBindingAdapter via SdlActionCalls reply-to 3. | 490 |
| f0036 | BoundActionAccepted step 7 sends BoundActionResult from SdlUiBindingAdapter to SduiDispatcher via UiDomainActions reply-to 2. | 491 |
| f0037 | BoundActionAccepted step 8 sends UiActionResult from SduiDispatcher to FyneBackend via NativeUiActions reply-to 1. | 492 |
| f0038 | BoundActionAccepted step 9 sends UiGenerationNotices variant UiGenerationChanged from SduiInstanceStore to FyneBackend via UiGenerationEvents. | 493 |
| f0039 | BoundActionRejected exercises TryDomainInteraction. | 494 |
| f0040 | BoundActionRejected has completeness = closed. | 495 |
| f0041 | BoundActionRejected illustrates G4M3UiDomainBinding. | 496 |
| f0042 | BoundActionRejected runs-in BoundExecution. | 497 |
| f0043 | BoundActionRejected step 1 sends UiActionRequest from FyneBackend to SduiDispatcher via NativeUiActions. | 498 |
| f0044 | BoundActionRejected step 2 sends UiActionRejected from SduiDispatcher to FyneBackend via NativeUiActions reply-to 1. | 499 |
| f0045 | BoundActionRequest has message-kind = request. | 500 |
| f0046 | BoundActionRequest upholds ActionArguments. | 501 |
| f0047 | BoundActionResult has message-kind = result. | 502 |
| f0048 | BoundActionResult replies-to BoundActionRequest. | 503 |
| f0049 | BoundActionResult upholds ActionOutcome. | 504 |
| f0050 | BoundInteraction requires SdlExecutionPort in mode BoundExecution. | 505 |
| f0051 | BoundInteraction requires UiSessionPort in mode BoundExecution. | 506 |
| f0052 | BuildGeneratedApplication allocated-to CommandLineHost in mode NativeBuild. | 507 |
| f0053 | BuildGeneratedApplication contributes-to NativeGoAssembly. | 508 |
| f0054 | BuildGeneratedApplication realizes NativeRealization. | 509 |
| f0055 | BuildGoRequest has message-kind = request. | 510 |
| f0056 | BuildGoRequest upholds GeneratedGoContract. | 511 |
| f0057 | BuildGoResult has message-kind = result. | 512 |
| f0058 | BuildGoResult replies-to BuildGoRequest. | 513 |
| f0059 | BuildGoResult upholds NativeBuildContract. | 514 |
| f0060 | BuildNativeRealization refines RealizeDesign. | 515 |
| f0061 | BuildPreparedFrame allocated-to FyneHost in mode UiPreview. | 516 |
| f0062 | BuildPreparedFrame contributes-to InteractiveUiPreview. | 517 |
| f0063 | BuildPreparedFrame realizes MeasuredPresentation. | 518 |
| f0064 | BuildSdlAst allocated-to CommandLineHost in mode SourceInspection. | 519 |
| f0065 | BuildSdlAst contributes-to StructuralModelInspection. | 520 |
| f0066 | BuildSdlAst has repeatability = deterministic. | 521 |
| f0067 | BuildSdlAst has state-retention = stateless. | 522 |
| f0068 | BuildSdlAst realizes SdlSourceModel. | 523 |
| f0069 | BuildSduiAst allocated-to CommandLineHost in mode SourceInspection. | 524 |
| f0070 | BuildSduiAst contributes-to StructuralModelInspection. | 525 |
| f0071 | BuildSduiAst has repeatability = deterministic. | 526 |
| f0072 | BuildSduiAst has state-retention = stateless. | 527 |
| f0073 | BuildSduiAst realizes SduiSourceModel. | 528 |
| f0074 | BuildUiAstRequest has message-kind = request. | 529 |
| f0075 | BuildUiAstRequest upholds TokenArtifactContract. | 530 |
| f0076 | BuildUiAstResult has message-kind = result. | 531 |
| f0077 | BuildUiAstResult replies-to BuildUiAstRequest. | 532 |
| f0078 | BuildUiAstResult upholds AstArtifactContract. | 533 |
| f0079 | CancelPendingActions realizes ExecutableDesign. | 534 |
| f0080 | CaptureNavigationTarget allocated-to XfmdDocumentHost in mode DocumentBrowsing. | 535 |
| f0081 | CaptureNavigationTarget contributes-to NavigableDesignDocumentation. | 536 |
| f0082 | CheckDomainStateCompatibility realizes ExecutableDesign. | 537 |
| f0083 | CheckExecutionCompleteness realizes ExecutableDesign. | 538 |
| f0084 | CheckFunctionSignatures realizes ExecutableDesign. | 539 |
| f0085 | CloseSdlInstance realizes ExecutableDesign. | 540 |
| f0086 | CloseUiInstance realizes InteractiveSession. | 541 |
| f0087 | CoalesceSourceChanges realizes DevelopmentReload. | 542 |
| f0088 | CommandLineHost consumes ExportSinkPort. | 543 |
| f0089 | CommandLineHost consumes PreparedFramePort. | 544 |
| f0090 | CommandLineHost consumes SduiFrontendPort. | 545 |
| f0091 | CommandLineHost consumes SourceSnapshotPort. | 546 |
| f0092 | CommandLineHost owns ComposeHeadlessExport. | 547 |
| f0093 | CommandLineHost owns WriteGeneratedArtifacts. | 548 |
| f0094 | CommandLineHost provides StaticDocumentation. | 549 |
| f0095 | CommandLineHost uses GoBuildCalls as receiver of BuildGoResult in mode NativeBuild. | 550 |
| f0096 | CommandLineHost uses GoBuildCalls as sender of BuildGoRequest in mode NativeBuild. | 551 |
| f0097 | CommandLineHost uses GoGenerationCalls as receiver of GenerateGoResult in mode NativeBuild. | 552 |
| f0098 | CommandLineHost uses GoGenerationCalls as sender of GenerateGoRequest in mode NativeBuild. | 553 |
| f0099 | CommandLineHost uses LayoutCalls as receiver of LayoutResult in mode StaticExport. | 554 |
| f0100 | CommandLineHost uses LayoutCalls as sender of LayoutRequest in mode StaticExport. | 555 |
| f0101 | CommandLineHost uses SvgExportCalls as receiver of ExportSvgResult in mode StaticExport. | 556 |
| f0102 | CommandLineHost uses SvgExportCalls as sender of ExportSvgRequest in mode StaticExport. | 557 |
| f0103 | CommandLineHost uses UiCompilationCalls as receiver of CompileUiRejected in mode SourceInspection. | 558 |
| f0104 | CommandLineHost uses UiCompilationCalls as receiver of CompileUiResult in mode SourceInspection. | 559 |
| f0105 | CommandLineHost uses UiCompilationCalls as sender of CompileUiRequest in mode SourceInspection. | 560 |
| f0106 | CompileSdlRequest has message-kind = request. | 561 |
| f0107 | CompileSdlRequest upholds ReloadArguments. | 562 |
| f0108 | CompileSdlResult has message-kind = result. | 563 |
| f0109 | CompileSdlResult replies-to CompileSdlRequest. | 564 |
| f0110 | CompileSdlResult upholds ReloadOutcome. | 565 |
| f0111 | CompileUiRejected has message-kind = result. | 566 |
| f0112 | CompileUiRejected replies-to CompileUiRequest. | 567 |
| f0113 | CompileUiRejected upholds ReloadOutcome. | 568 |
| f0114 | CompileUiRequest has message-kind = request. | 569 |
| f0115 | CompileUiRequest upholds ReloadArguments. | 570 |
| f0116 | CompileUiResult has message-kind = result. | 571 |
| f0117 | CompileUiResult replies-to CompileUiRequest. | 572 |
| f0118 | CompileUiResult upholds ReloadOutcome. | 573 |
| f0119 | ComposeHeadlessExport realizes StaticDocumentation. | 574 |
| f0120 | ComposeInteractiveSession allocated-to FyneHost in mode UiPreview. | 575 |
| f0121 | ComposeInteractiveSession contributes-to InteractiveUiPreview. | 576 |
| f0122 | ComposeInteractiveSession realizes NativeInteraction. | 577 |
| f0123 | ComposeMarkdownDocument allocated-to CommandLineHost in mode StaticExport. | 578 |
| f0124 | ComposeMarkdownDocument contributes-to DesignDocumentation. | 579 |
| f0125 | ComposeMarkdownDocument realizes StaticDocumentation. | 580 |
| f0126 | ComposeViewPackage allocated-to ViewServiceHost in mode DocumentBrowsing. | 581 |
| f0127 | ComposeViewPackage contributes-to NavigableDesignDocumentation. | 582 |
| f0128 | ComputeClipping realizes MeasuredPresentation. | 583 |
| f0129 | ConnectTypedWidgetHandles allocated-to FyneHost in mode BoundExecution. | 584 |
| f0130 | ConnectTypedWidgetHandles contributes-to TypedDomainBinding. | 585 |
| f0131 | ConnectTypedWidgetHandles realizes BoundInteraction. | 586 |
| f0132 | ContentServices contains DiagramProvider. | 587 |
| f0133 | ContentServices contains MarkdownProvider. | 588 |
| f0134 | ContentServices contains ResourceStore. | 589 |
| f0135 | ContentServices provides RichContent. | 590 |
| f0136 | CoordinateSdlCompilation realizes SdlSourceModel. | 591 |
| f0137 | CoordinateSduiCompilation realizes SduiSourceModel. | 592 |
| f0138 | CorrelateActionResult realizes ExecutableDesign. | 593 |
| f0139 | CorrelateUiResult realizes InteractiveSession. | 594 |
| f0140 | CreateSdlInstance realizes ExecutableDesign. | 595 |
| f0141 | CreateUiInstance realizes InteractiveSession. | 596 |
| f0142 | DesignAuthor pursues BrowseDesignViews. | 597 |
| f0143 | DesignAuthor pursues BuildNativeProduct. | 598 |
| f0144 | DesignAuthor pursues EditRunningPrototype. | 599 |
| f0145 | DesignAuthor pursues InspectModels. | 600 |
| f0146 | DesignAuthor pursues PrototypeUserInterface. | 601 |
| f0147 | DesignAuthor pursues PublishDesignDocumentation. | 602 |
| f0148 | DesignAuthor pursues TryDomainInteraction. | 603 |
| f0149 | DesignDocumentation supports PublishDesignDocumentation. | 604 |
| f0150 | DesignReviewer pursues BrowseDesignViews. | 605 |
| f0151 | DesignReviewer pursues InspectModels. | 606 |
| f0152 | DesignReviewer pursues PrototypeUserInterface. | 607 |
| f0153 | DesignReviewer pursues PublishDesignDocumentation. | 608 |
| f0154 | DesignSourceArchive holds DesignSourceDocuments. | 609 |
| f0155 | DesignSourceDocuments upholds DesignSourceRecord. | 610 |
| f0156 | DesignSourceRecord has completeness = closed. | 611 |
| f0157 | DesignSourceRecord has-field SourceDocumentRevision. | 612 |
| f0158 | DesignSourceRecord has-field SourceDocumentText. | 613 |
| f0159 | DevelopmentReload requires BindingReloadPort in mode BoundLiveEditing. | 614 |
| f0160 | DevelopmentReload requires FileChangePort in mode BoundLiveEditing. | 615 |
| f0161 | DevelopmentReload requires FileChangePort in mode LiveEditing. | 616 |
| f0162 | DevelopmentReload requires SdlReloadPort in mode BoundLiveEditing. | 617 |
| f0163 | DevelopmentReload requires SourceSnapshotPort in mode BoundLiveEditing. | 618 |
| f0164 | DevelopmentReload requires SourceSnapshotPort in mode LiveEditing. | 619 |
| f0165 | DevelopmentReload requires UiReloadPort in mode BoundLiveEditing. | 620 |
| f0166 | DevelopmentReload requires UiReloadPort in mode LiveEditing. | 621 |
| f0167 | DevelopmentTools contains DiagnosticReporter. | 622 |
| f0168 | DevelopmentTools contains GoBuildRunner. | 623 |
| f0169 | DevelopmentTools contains GoCodeGenerator. | 624 |
| f0170 | DevelopmentTools contains ReloadCoordinator. | 625 |
| f0171 | DevelopmentTools contains SourceLoader. | 626 |
| f0172 | DevelopmentTools contains SourceWatcher. | 627 |
| f0173 | DevelopmentTools owns RetireReplacedPythonEntryPoints. | 628 |
| f0174 | DevelopmentTools owns VerifyNativeBehaviorParity. | 629 |
| f0175 | DevelopmentTools provides DevelopmentReload. | 630 |
| f0176 | DiagnosticReporter owns ReportBindingDiagnostics. | 631 |
| f0177 | DiagnosticReporter owns ReportSourceDiagnostics. | 632 |
| f0178 | DiagnosticReporter provides SourceDiagnostics. | 633 |
| f0179 | DiagramProvider consumes DiagramEnginePort. | 634 |
| f0180 | DiagramProvider consumes ResourcePort. | 635 |
| f0181 | DiagramProvider owns PrepareDiagramResource. | 636 |
| f0182 | DiagramProvider provides RichContent. | 637 |
| f0183 | DisconnectBindings realizes BoundInteraction. | 638 |
| f0184 | DispatchUiEvent allocated-to FyneHost in mode BoundExecution. | 639 |
| f0185 | DispatchUiEvent contributes-to TypedDomainBinding. | 640 |
| f0186 | DispatchUiEvent realizes InteractiveSession. | 641 |
| f0187 | DispatchViewOpen allocated-to ViewServiceHost in mode DocumentBrowsing. | 642 |
| f0188 | DispatchViewOpen contributes-to NavigableDesignDocumentation. | 643 |
| f0189 | DisplayViewRequest has message-kind = request. | 644 |
| f0190 | DisplayViewRequest upholds ViewOpenContract. | 645 |
| f0191 | DocumentBroker owns KeyViewRevision. | 646 |
| f0192 | DocumentBroker owns RejectStaleViewResults. | 647 |
| f0193 | DocumentBroker owns ServeViewRequests. | 648 |
| f0194 | DocumentBroker owns ValidateViewRequest. | 649 |
| f0195 | DocumentBroker uses ViewLaunchCalls as receiver of ViewLaunchResult in mode DocumentBrowsing. | 650 |
| f0196 | DocumentBroker uses ViewLaunchCalls as sender of LaunchViewRequest in mode DocumentBrowsing. | 651 |
| f0197 | DocumentBroker uses ViewNavigationCalls as receiver of SelectViewRequest in mode DocumentBrowsing. | 652 |
| f0198 | DocumentBroker uses ViewNavigationCalls as sender of ViewOpenedResult in mode DocumentBrowsing. | 653 |
| f0199 | DocumentBroker uses ViewNavigationCalls as sender of ViewRequestRejected in mode DocumentBrowsing. | 654 |
| f0200 | DocumentBroker uses ViewProjectionCalls as receiver of ViewBundleResult in mode DocumentBrowsing. | 655 |
| f0201 | DocumentBroker uses ViewProjectionCalls as receiver of ViewProjectionRejected in mode DocumentBrowsing. | 656 |
| f0202 | DocumentBroker uses ViewProjectionCalls as sender of ProjectViewRequest in mode DocumentBrowsing. | 657 |
| f0203 | DocumentBroker uses ViewPublicationCalls as receiver of ViewReferenceResult in mode DocumentBrowsing. | 658 |
| f0204 | DocumentBroker uses ViewPublicationCalls as sender of PublishViewRequest in mode DocumentBrowsing. | 659 |
| f0205 | DomainActionRequest has message-kind = request. | 660 |
| f0206 | DomainActionRequest upholds ActionArguments. | 661 |
| f0207 | DomainActionResult has message-kind = result. | 662 |
| f0208 | DomainActionResult replies-to DomainActionRequest. | 663 |
| f0209 | DomainActionResult upholds ActionOutcome. | 664 |
| f0210 | DomainStateMigrator consumes DomainStatePort. | 665 |
| f0211 | DomainStateMigrator consumes SdlModelPort. | 666 |
| f0212 | DomainStateMigrator owns CheckDomainStateCompatibility. | 667 |
| f0213 | DomainStateMigrator owns MigrateOrResetDomainState. | 668 |
| f0214 | DomainStateMigrator provides ExecutableDesign. | 669 |
| f0215 | EvictUnusedViewBundles allocated-to ViewServiceHost in mode DocumentBrowsing. | 670 |
| f0216 | EvictUnusedViewBundles contributes-to NavigableDesignDocumentation. | 671 |
| f0217 | ExecutableDesign requires DomainFunctionPort in mode BoundExecution. | 672 |
| f0218 | ExpandUiDefinitions has repeatability = deterministic. | 673 |
| f0219 | ExpandUiDefinitions has state-retention = stateless. | 674 |
| f0220 | ExpandUiDefinitions realizes SduiSourceModel. | 675 |
| f0221 | ExportConsoleSnapshot realizes StaticDocumentation. | 676 |
| f0222 | ExportSvgRequest has message-kind = request. | 677 |
| f0223 | ExportSvgRequest upholds PreparedFrameContract. | 678 |
| f0224 | ExportSvgResult has message-kind = result. | 679 |
| f0225 | ExportSvgResult replies-to ExportSvgRequest. | 680 |
| f0226 | ExportSvgResult upholds SvgDocumentContract. | 681 |
| f0227 | ExportSvgSnapshot allocated-to CommandLineHost in mode StaticExport. | 682 |
| f0228 | ExportSvgSnapshot contributes-to DesignDocumentation. | 683 |
| f0229 | ExportSvgSnapshot realizes StaticDocumentation. | 684 |
| f0230 | ExportUiDocumentation refines InspectDesignSource. | 685 |
| f0231 | ExportViewpointMarkdown allocated-to CommandLineHost in mode StaticExport. | 686 |
| f0232 | ExportViewpointMarkdown allocated-to ViewServiceHost in mode DocumentBrowsing. | 687 |
| f0233 | ExportViewpointMarkdown contributes-to DesignDocumentation. | 688 |
| f0234 | ExportViewpointMarkdown has repeatability = deterministic. | 689 |
| f0235 | ExportViewpointMarkdown has state-retention = stateless. | 690 |
| f0236 | FailedViewRequestId has presence = required. | 691 |
| f0237 | FailedViewRequestId has value-type = text. | 692 |
| f0238 | FramePresentationCalls upholds FramePresentationCallsProtocol. | 693 |
| f0239 | FramePresentationCallsProtocol has completeness = closed. | 694 |
| f0240 | FramePresentationCallsProtocol permits PresentFrameRequest. | 695 |
| f0241 | FramePresentationCallsProtocol permits PresentFrameResult. | 696 |
| f0242 | FyneBackend consumes PreparedFramePort. | 697 |
| f0243 | FyneBackend consumes UiSessionPort. | 698 |
| f0244 | FyneBackend owns HandleFocusAndTextInput. | 699 |
| f0245 | FyneBackend owns PublishPresentation. | 700 |
| f0246 | FyneBackend owns ReconcileWidgets. | 701 |
| f0247 | FyneBackend owns ReleaseNativeWidgets. | 702 |
| f0248 | FyneBackend provides NativeInteraction. | 703 |
| f0249 | FyneBackend uses FramePresentationCalls as receiver of PresentFrameRequest in mode UiPreview. | 704 |
| f0250 | FyneBackend uses FramePresentationCalls as sender of PresentFrameResult in mode UiPreview. | 705 |
| f0251 | FyneBackend uses NativeUiActions as receiver of UiActionRejected in mode BoundExecution. | 706 |
| f0252 | FyneBackend uses NativeUiActions as receiver of UiActionResult in mode BoundExecution. | 707 |
| f0253 | FyneBackend uses NativeUiActions as sender of UiActionRequest in mode BoundExecution. | 708 |
| f0254 | FyneBackend uses UiGenerationEvents as receiver of UiGenerationNotices in mode BoundExecution. | 709 |
| f0255 | FyneBackend uses UiGenerationEvents as receiver of UiGenerationNotices in mode LiveEditing. | 710 |
| f0256 | FyneHost consumes DomainBindingPort. | 711 |
| f0257 | FyneHost consumes ReloadPort. | 712 |
| f0258 | FyneHost consumes SdlFrontendPort. | 713 |
| f0259 | FyneHost consumes SduiFrontendPort. | 714 |
| f0260 | FyneHost consumes SourceSnapshotPort. | 715 |
| f0261 | FyneHost consumes UiSessionPort. | 716 |
| f0262 | FyneHost consumes WidgetBackendPort. | 717 |
| f0263 | FyneHost owns ComposeInteractiveSession. | 718 |
| f0264 | FyneHost owns ScheduleUiPublication. | 719 |
| f0265 | FyneHost provides NativeInteraction. | 720 |
| f0266 | FyneHost uses FramePresentationCalls as receiver of PresentFrameResult in mode UiPreview. | 721 |
| f0267 | FyneHost uses FramePresentationCalls as sender of PresentFrameRequest in mode UiPreview. | 722 |
| f0268 | FyneHost uses GoDomainCalls as receiver of DomainActionResult in mode UiPreview. | 723 |
| f0269 | FyneHost uses GoDomainCalls as sender of DomainActionRequest in mode UiPreview. | 724 |
| f0270 | FyneHost uses LayoutCalls as receiver of LayoutResult in mode UiPreview. | 725 |
| f0271 | FyneHost uses LayoutCalls as sender of LayoutRequest in mode UiPreview. | 726 |
| f0272 | G1FrontendPort delivers StructuralModelInspection. | 727 |
| f0273 | G1FrontendPort has implementation-status = planned. | 728 |
| f0274 | G1M1ParserAndAst addresses BuildSduiAst. | 729 |
| f0275 | G1M1ParserAndAst addresses IdentifySourceRevision. | 730 |
| f0276 | G1M1ParserAndAst addresses ReadBoundedSources. | 731 |
| f0277 | G1M1ParserAndAst addresses TokenizeSduiSource. | 732 |
| f0278 | G1M1ParserAndAst has implementation-status = planned. | 733 |
| f0279 | G1M1ParserAndAst refines G1FrontendPort. | 734 |
| f0280 | G1M2ValidationAndNormalization addresses CoordinateSduiCompilation. | 735 |
| f0281 | G1M2ValidationAndNormalization addresses ExpandUiDefinitions. | 736 |
| f0282 | G1M2ValidationAndNormalization addresses PreserveUiRegions. | 737 |
| f0283 | G1M2ValidationAndNormalization addresses PreserveUiSourceMap. | 738 |
| f0284 | G1M2ValidationAndNormalization addresses ResolveUiNames. | 739 |
| f0285 | G1M2ValidationAndNormalization addresses ValidateRelativeFormatting. | 740 |
| f0286 | G1M2ValidationAndNormalization addresses ValidateSymbolicBindings. | 741 |
| f0287 | G1M2ValidationAndNormalization addresses ValidateWidgetArguments. | 742 |
| f0288 | G1M2ValidationAndNormalization depends-on G1M1ParserAndAst. | 743 |
| f0289 | G1M2ValidationAndNormalization has implementation-status = planned. | 744 |
| f0290 | G1M2ValidationAndNormalization refines G1FrontendPort. | 745 |
| f0291 | G1M3Concept1AndDumps addresses ExportConsoleSnapshot. | 746 |
| f0292 | G1M3Concept1AndDumps addresses ReportSourceDiagnostics. | 747 |
| f0293 | G1M3Concept1AndDumps depends-on G1M2ValidationAndNormalization. | 748 |
| f0294 | G1M3Concept1AndDumps has implementation-status = planned. | 749 |
| f0295 | G1M3Concept1AndDumps refines G1FrontendPort. | 750 |
| f0296 | G2LayoutAndPresentation delivers InteractiveUiPreview. | 751 |
| f0297 | G2LayoutAndPresentation has implementation-status = planned. | 752 |
| f0298 | G2M1RelativeMeasurement addresses AllocateGeometry. | 753 |
| f0299 | G2M1RelativeMeasurement addresses ComputeClipping. | 754 |
| f0300 | G2M1RelativeMeasurement addresses MeasureUiContent. | 755 |
| f0301 | G2M1RelativeMeasurement addresses ResolveAncestorDimensions. | 756 |
| f0302 | G2M1RelativeMeasurement depends-on G1M2ValidationAndNormalization. | 757 |
| f0303 | G2M1RelativeMeasurement has implementation-status = planned. | 758 |
| f0304 | G2M1RelativeMeasurement refines G2LayoutAndPresentation. | 759 |
| f0305 | G2M2SharedSvgGeometry addresses BuildPreparedFrame. | 760 |
| f0306 | G2M2SharedSvgGeometry addresses ExportSvgSnapshot. | 761 |
| f0307 | G2M2SharedSvgGeometry depends-on G2M1RelativeMeasurement. | 762 |
| f0308 | G2M2SharedSvgGeometry has implementation-status = planned. | 763 |
| f0309 | G2M2SharedSvgGeometry refines G2LayoutAndPresentation. | 764 |
| f0310 | G2M3FyneInteractions addresses ComposeInteractiveSession. | 765 |
| f0311 | G2M3FyneInteractions addresses HandleFocusAndTextInput. | 766 |
| f0312 | G2M3FyneInteractions addresses PublishPresentation. | 767 |
| f0313 | G2M3FyneInteractions addresses ReconcileWidgets. | 768 |
| f0314 | G2M3FyneInteractions addresses ReleaseNativeWidgets. | 769 |
| f0315 | G2M3FyneInteractions addresses ScheduleUiPublication. | 770 |
| f0316 | G2M3FyneInteractions depends-on G2M2SharedSvgGeometry. | 771 |
| f0317 | G2M3FyneInteractions has implementation-status = planned. | 772 |
| f0318 | G2M3FyneInteractions refines G2LayoutAndPresentation. | 773 |
| f0319 | G2M4RichContent addresses MeasureMarkdownContent. | 774 |
| f0320 | G2M4RichContent addresses PrepareDiagramResource. | 775 |
| f0321 | G2M4RichContent addresses PrepareMarkdown. | 776 |
| f0322 | G2M4RichContent addresses ReleaseVisualResources. | 777 |
| f0323 | G2M4RichContent addresses ValidateVisualResources. | 778 |
| f0324 | G2M4RichContent depends-on G2M3FyneInteractions. | 779 |
| f0325 | G2M4RichContent has implementation-status = planned. | 780 |
| f0326 | G2M4RichContent refines G2LayoutAndPresentation. | 781 |
| f0327 | G3M1TypedUiSession addresses ApplyPropertyBatch. | 782 |
| f0328 | G3M1TypedUiSession addresses CloseUiInstance. | 783 |
| f0329 | G3M1TypedUiSession addresses CorrelateUiResult. | 784 |
| f0330 | G3M1TypedUiSession addresses CreateUiInstance. | 785 |
| f0331 | G3M1TypedUiSession addresses DispatchUiEvent. | 786 |
| f0332 | G3M1TypedUiSession addresses ManageWidgetIdentities. | 787 |
| f0333 | G3M1TypedUiSession addresses ProjectUiGeneration. | 788 |
| f0334 | G3M1TypedUiSession addresses RejectStaleUiEvent. | 789 |
| f0335 | G3M1TypedUiSession addresses RevokeWidgetGenerations. | 790 |
| f0336 | G3M1TypedUiSession addresses SnapshotUiState. | 791 |
| f0337 | G3M1TypedUiSession addresses TrackInputDraft. | 792 |
| f0338 | G3M1TypedUiSession addresses ValidatePropertyBatch. | 793 |
| f0339 | G3M1TypedUiSession addresses ValidateUiEvent. | 794 |
| f0340 | G3M1TypedUiSession depends-on G1M2ValidationAndNormalization. | 795 |
| f0341 | G3M1TypedUiSession depends-on G2M3FyneInteractions. | 796 |
| f0342 | G3M1TypedUiSession has implementation-status = planned. | 797 |
| f0343 | G3M1TypedUiSession refines G3UiRuntimeAndReload. | 798 |
| f0344 | G3M2CandidatePublication addresses CoalesceSourceChanges. | 799 |
| f0345 | G3M2CandidatePublication addresses KeepLastValidModels. | 800 |
| f0346 | G3M2CandidatePublication addresses ObserveSourceChanges. | 801 |
| f0347 | G3M2CandidatePublication addresses PrepareCandidateModels. | 802 |
| f0348 | G3M2CandidatePublication addresses PublishModelGeneration. | 803 |
| f0349 | G3M2CandidatePublication addresses RetirePreviousGeneration. | 804 |
| f0350 | G3M2CandidatePublication depends-on G3M1TypedUiSession. | 805 |
| f0351 | G3M2CandidatePublication has implementation-status = planned. | 806 |
| f0352 | G3M2CandidatePublication refines G3UiRuntimeAndReload. | 807 |
| f0353 | G3M3CompatibleState addresses MatchCompatibleWidgets. | 808 |
| f0354 | G3M3CompatibleState addresses PreserveCompatibleUiState. | 809 |
| f0355 | G3M3CompatibleState addresses ResetIncompatibleUiState. | 810 |
| f0356 | G3M3CompatibleState depends-on G3M2CandidatePublication. | 811 |
| f0357 | G3M3CompatibleState has implementation-status = planned. | 812 |
| f0358 | G3M3CompatibleState refines G3UiRuntimeAndReload. | 813 |
| f0359 | G3UiRuntimeAndReload delivers LiveModelReload. | 814 |
| f0360 | G3UiRuntimeAndReload has implementation-status = planned. | 815 |
| f0361 | G4M1SdlFrontend addresses BuildSdlAst. | 816 |
| f0362 | G4M1SdlFrontend addresses CoordinateSdlCompilation. | 817 |
| f0363 | G4M1SdlFrontend addresses NormalizeSdlModel. | 818 |
| f0364 | G4M1SdlFrontend addresses PreserveSdlSourceMap. | 819 |
| f0365 | G4M1SdlFrontend addresses ResolveSdlSymbols. | 820 |
| f0366 | G4M1SdlFrontend addresses TokenizeSdlSource. | 821 |
| f0367 | G4M1SdlFrontend addresses ValidateSdlProfile. | 822 |
| f0368 | G4M1SdlFrontend addresses ValidateSdlStructure. | 823 |
| f0369 | G4M1SdlFrontend has implementation-status = planned. | 824 |
| f0370 | G4M1SdlFrontend refines G4SdlRuntimeAndBinding. | 825 |
| f0371 | G4M2TypedExecution addresses CancelPendingActions. | 826 |
| f0372 | G4M2TypedExecution addresses CheckExecutionCompleteness. | 827 |
| f0373 | G4M2TypedExecution addresses CheckFunctionSignatures. | 828 |
| f0374 | G4M2TypedExecution addresses CloseSdlInstance. | 829 |
| f0375 | G4M2TypedExecution addresses CorrelateActionResult. | 830 |
| f0376 | G4M2TypedExecution addresses CreateSdlInstance. | 831 |
| f0377 | G4M2TypedExecution addresses InvokeRegisteredFunction. | 832 |
| f0378 | G4M2TypedExecution addresses ManageDomainState. | 833 |
| f0379 | G4M2TypedExecution addresses PerformDomainOperation. | 834 |
| f0380 | G4M2TypedExecution addresses RegisterDomainFunctions. | 835 |
| f0381 | G4M2TypedExecution addresses SnapshotDomainState. | 836 |
| f0382 | G4M2TypedExecution addresses ValidateActionInput. | 837 |
| f0383 | G4M2TypedExecution depends-on G4M1SdlFrontend. | 838 |
| f0384 | G4M2TypedExecution has implementation-status = planned. | 839 |
| f0385 | G4M2TypedExecution refines G4SdlRuntimeAndBinding. | 840 |
| f0386 | G4M3UiDomainBinding addresses ConnectTypedWidgetHandles. | 841 |
| f0387 | G4M3UiDomainBinding addresses DisconnectBindings. | 842 |
| f0388 | G4M3UiDomainBinding addresses PublishDomainUpdates. | 843 |
| f0389 | G4M3UiDomainBinding addresses ReportBindingDiagnostics. | 844 |
| f0390 | G4M3UiDomainBinding addresses ResolveCallbackSymbols. | 845 |
| f0391 | G4M3UiDomainBinding addresses RouteDomainBindings. | 846 |
| f0392 | G4M3UiDomainBinding depends-on G3M1TypedUiSession. | 847 |
| f0393 | G4M3UiDomainBinding depends-on G4M2TypedExecution. | 848 |
| f0394 | G4M3UiDomainBinding has implementation-status = planned. | 849 |
| f0395 | G4M3UiDomainBinding refines G4SdlRuntimeAndBinding. | 850 |
| f0396 | G4M4DomainReload addresses CheckDomainStateCompatibility. | 851 |
| f0397 | G4M4DomainReload addresses MigrateOrResetDomainState. | 852 |
| f0398 | G4M4DomainReload addresses RestartChangedGoProgram. | 853 |
| f0399 | G4M4DomainReload depends-on G3M3CompatibleState. | 854 |
| f0400 | G4M4DomainReload depends-on G4M3UiDomainBinding. | 855 |
| f0401 | G4M4DomainReload has implementation-status = planned. | 856 |
| f0402 | G4M4DomainReload refines G4SdlRuntimeAndBinding. | 857 |
| f0403 | G4SdlRuntimeAndBinding delivers TypedDomainBinding. | 858 |
| f0404 | G4SdlRuntimeAndBinding has implementation-status = planned. | 859 |
| f0405 | G5M1GeneratedGo addresses BuildGeneratedApplication. | 860 |
| f0406 | G5M1GeneratedGo addresses GenerateBindingRegistration. | 861 |
| f0407 | G5M1GeneratedGo addresses GenerateModelConstructors. | 862 |
| f0408 | G5M1GeneratedGo addresses PreserveHandwrittenSources. | 863 |
| f0409 | G5M1GeneratedGo depends-on G4M4DomainReload. | 864 |
| f0410 | G5M1GeneratedGo has implementation-status = planned. | 865 |
| f0411 | G5M1GeneratedGo refines G5NativeGeneration. | 866 |
| f0412 | G5M2BehaviorParity addresses VerifyNativeBehaviorParity. | 867 |
| f0413 | G5M2BehaviorParity depends-on G5M1GeneratedGo. | 868 |
| f0414 | G5M2BehaviorParity has implementation-status = planned. | 869 |
| f0415 | G5M2BehaviorParity refines G5NativeGeneration. | 870 |
| f0416 | G5M3DocumentationExport addresses ComposeHeadlessExport. | 871 |
| f0417 | G5M3DocumentationExport addresses ComposeMarkdownDocument. | 872 |
| f0418 | G5M3DocumentationExport addresses WriteGeneratedArtifacts. | 873 |
| f0419 | G5M3DocumentationExport depends-on G2M4RichContent. | 874 |
| f0420 | G5M3DocumentationExport depends-on G5M2BehaviorParity. | 875 |
| f0421 | G5M3DocumentationExport depends-on G6M1StaticNavigation. | 876 |
| f0422 | G5M3DocumentationExport has implementation-status = planned. | 877 |
| f0423 | G5M3DocumentationExport refines G5NativeGeneration. | 878 |
| f0424 | G5M4RetirePython addresses RetireReplacedPythonEntryPoints. | 879 |
| f0425 | G5M4RetirePython depends-on G1M3Concept1AndDumps. | 880 |
| f0426 | G5M4RetirePython depends-on G5M3DocumentationExport. | 881 |
| f0427 | G5M4RetirePython has implementation-status = planned. | 882 |
| f0428 | G5M4RetirePython refines G5NativeGeneration. | 883 |
| f0429 | G5NativeGeneration delivers DesignDocumentation. | 884 |
| f0430 | G5NativeGeneration delivers NativeGoAssembly. | 885 |
| f0431 | G5NativeGeneration has implementation-status = planned. | 886 |
| f0432 | G6M1StaticNavigation addresses ComposeViewPackage. | 887 |
| f0433 | G6M1StaticNavigation addresses ExportViewpointMarkdown. | 888 |
| f0434 | G6M1StaticNavigation addresses GenerateViewNavigation. | 889 |
| f0435 | G6M1StaticNavigation addresses PreserveViewAnchors. | 890 |
| f0436 | G6M1StaticNavigation addresses ProjectSdlViewpoints. | 891 |
| f0437 | G6M1StaticNavigation addresses TraceViewpointFacts. | 892 |
| f0438 | G6M1StaticNavigation depends-on G4M1SdlFrontend. | 893 |
| f0439 | G6M1StaticNavigation has implementation-status = planned. | 894 |
| f0440 | G6M1StaticNavigation refines G6NavigableDocumentation. | 895 |
| f0441 | G6M2OnDemandViews addresses KeyViewRevision. | 896 |
| f0442 | G6M2OnDemandViews addresses ProjectSelectedView. | 897 |
| f0443 | G6M2OnDemandViews addresses PublishViewBundle. | 898 |
| f0444 | G6M2OnDemandViews addresses ValidateViewRequest. | 899 |
| f0445 | G6M2OnDemandViews depends-on G6M1StaticNavigation. | 900 |
| f0446 | G6M2OnDemandViews has implementation-status = planned. | 901 |
| f0447 | G6M2OnDemandViews refines G6NavigableDocumentation. | 902 |
| f0448 | G6M3XfmdNavigation addresses CaptureNavigationTarget. | 903 |
| f0449 | G6M3XfmdNavigation addresses DispatchViewOpen. | 904 |
| f0450 | G6M3XfmdNavigation addresses ResolveConfiguredViewer. | 905 |
| f0451 | G6M3XfmdNavigation addresses RouteDocumentToPane. | 906 |
| f0452 | G6M3XfmdNavigation depends-on G6M2OnDemandViews. | 907 |
| f0453 | G6M3XfmdNavigation has implementation-status = planned. | 908 |
| f0454 | G6M3XfmdNavigation refines G6NavigableDocumentation. | 909 |
| f0455 | G6M4SessionPublication addresses EvictUnusedViewBundles. | 910 |
| f0456 | G6M4SessionPublication addresses RejectStaleViewResults. | 911 |
| f0457 | G6M4SessionPublication addresses RetainVisibleViewBundle. | 912 |
| f0458 | G6M4SessionPublication addresses ServeViewRequests. | 913 |
| f0459 | G6M4SessionPublication depends-on G6M3XfmdNavigation. | 914 |
| f0460 | G6M4SessionPublication has implementation-status = planned. | 915 |
| f0461 | G6M4SessionPublication refines G6NavigableDocumentation. | 916 |
| f0462 | G6NavigableDocumentation delivers NavigableDesignDocumentation. | 917 |
| f0463 | G6NavigableDocumentation has implementation-status = planned. | 918 |
| f0464 | GenerateBindingRegistration allocated-to CommandLineHost in mode NativeBuild. | 919 |
| f0465 | GenerateBindingRegistration contributes-to NativeGoAssembly. | 920 |
| f0466 | GenerateBindingRegistration realizes NativeRealization. | 921 |
| f0467 | GenerateGoRequest has message-kind = request. | 922 |
| f0468 | GenerateGoRequest upholds NormalizedModelContract. | 923 |
| f0469 | GenerateGoResult has message-kind = result. | 924 |
| f0470 | GenerateGoResult replies-to GenerateGoRequest. | 925 |
| f0471 | GenerateGoResult upholds GeneratedGoContract. | 926 |
| f0472 | GenerateModelConstructors allocated-to CommandLineHost in mode NativeBuild. | 927 |
| f0473 | GenerateModelConstructors contributes-to NativeGoAssembly. | 928 |
| f0474 | GenerateModelConstructors realizes NativeRealization. | 929 |
| f0475 | GenerateViewNavigation allocated-to ViewServiceHost in mode DocumentBrowsing. | 930 |
| f0476 | GenerateViewNavigation contributes-to NavigableDesignDocumentation. | 931 |
| f0477 | GeneratedGoContract has completeness = closed. | 932 |
| f0478 | GeneratedGoContract has-field GeneratedGoSources. | 933 |
| f0479 | GeneratedGoSources has presence = required. | 934 |
| f0480 | GeneratedGoSources has value-type = bytes. | 935 |
| f0481 | GoBuildCalls upholds GoBuildCallsProtocol. | 936 |
| f0482 | GoBuildCallsProtocol has completeness = closed. | 937 |
| f0483 | GoBuildCallsProtocol permits BuildGoRequest. | 938 |
| f0484 | GoBuildCallsProtocol permits BuildGoResult. | 939 |
| f0485 | GoBuildRunner consumes BuildToolPort. | 940 |
| f0486 | GoBuildRunner consumes GeneratedArtifactPort. | 941 |
| f0487 | GoBuildRunner owns BuildGeneratedApplication. | 942 |
| f0488 | GoBuildRunner owns RestartChangedGoProgram. | 943 |
| f0489 | GoBuildRunner provides NativeRealization. | 944 |
| f0490 | GoBuildRunner uses GoBuildCalls as receiver of BuildGoRequest in mode NativeBuild. | 945 |
| f0491 | GoBuildRunner uses GoBuildCalls as sender of BuildGoResult in mode NativeBuild. | 946 |
| f0492 | GoCodeGenerator consumes ExecutionProfilePort. | 947 |
| f0493 | GoCodeGenerator consumes SdlModelPort. | 948 |
| f0494 | GoCodeGenerator consumes SduiModelPort. | 949 |
| f0495 | GoCodeGenerator owns GenerateBindingRegistration. | 950 |
| f0496 | GoCodeGenerator owns GenerateModelConstructors. | 951 |
| f0497 | GoCodeGenerator owns PreserveHandwrittenSources. | 952 |
| f0498 | GoCodeGenerator provides NativeRealization. | 953 |
| f0499 | GoCodeGenerator uses GoGenerationCalls as receiver of GenerateGoRequest in mode NativeBuild. | 954 |
| f0500 | GoCodeGenerator uses GoGenerationCalls as sender of GenerateGoResult in mode NativeBuild. | 955 |
| f0501 | GoDomainCalls upholds GoDomainCallsProtocol. | 956 |
| f0502 | GoDomainCallsProtocol has completeness = closed. | 957 |
| f0503 | GoDomainCallsProtocol permits DomainActionRequest. | 958 |
| f0504 | GoDomainCallsProtocol permits DomainActionResult. | 959 |
| f0505 | GoDomainImplementation owns PerformDomainOperation. | 960 |
| f0506 | GoDomainImplementation provides DomainOperations. | 961 |
| f0507 | GoDomainImplementation uses GoDomainCalls as receiver of DomainActionRequest in mode BoundExecution. | 962 |
| f0508 | GoDomainImplementation uses GoDomainCalls as receiver of DomainActionRequest in mode UiPreview. | 963 |
| f0509 | GoDomainImplementation uses GoDomainCalls as sender of DomainActionResult in mode BoundExecution. | 964 |
| f0510 | GoDomainImplementation uses GoDomainCalls as sender of DomainActionResult in mode UiPreview. | 965 |
| f0511 | GoGenerationCalls upholds GoGenerationCallsProtocol. | 966 |
| f0512 | GoGenerationCallsProtocol has completeness = closed. | 967 |
| f0513 | GoGenerationCallsProtocol permits GenerateGoRequest. | 968 |
| f0514 | GoGenerationCallsProtocol permits GenerateGoResult. | 969 |
| f0515 | HandleFocusAndTextInput allocated-to FyneHost in mode UiPreview. | 970 |
| f0516 | HandleFocusAndTextInput contributes-to InteractiveUiPreview. | 971 |
| f0517 | HandleFocusAndTextInput realizes NativeInteraction. | 972 |
| f0518 | IdentifySourceRevision realizes SourceLoading. | 973 |
| f0519 | InspectSdlSource refines InspectDesignSource. | 974 |
| f0520 | InspectSduiSource refines InspectDesignSource. | 975 |
| f0521 | InteractiveFramePrepared exercises PrototypeUserInterface. | 976 |
| f0522 | InteractiveFramePrepared has completeness = closed. | 977 |
| f0523 | InteractiveFramePrepared illustrates G2M1RelativeMeasurement. | 978 |
| f0524 | InteractiveFramePrepared illustrates G2M3FyneInteractions. | 979 |
| f0525 | InteractiveFramePrepared runs-in UiPreview. | 980 |
| f0526 | InteractiveFramePrepared step 1 sends LayoutRequest from FyneHost to SduiLayout via LayoutCalls. | 981 |
| f0527 | InteractiveFramePrepared step 2 sends LayoutResult from SduiLayout to FyneHost via LayoutCalls reply-to 1. | 982 |
| f0528 | InteractiveFramePrepared step 3 sends PresentFrameRequest from FyneHost to FyneBackend via FramePresentationCalls. | 983 |
| f0529 | InteractiveFramePrepared step 4 sends PresentFrameResult from FyneBackend to FyneHost via FramePresentationCalls reply-to 3. | 984 |
| f0530 | InteractiveUiPreview supports PrototypeUserInterface. | 985 |
| f0531 | InvalidViewSelectionRejected exercises BrowseDesignViews. | 986 |
| f0532 | InvalidViewSelectionRejected has completeness = closed. | 987 |
| f0533 | InvalidViewSelectionRejected illustrates G6M2OnDemandViews. | 988 |
| f0534 | InvalidViewSelectionRejected illustrates G6M4SessionPublication. | 989 |
| f0535 | InvalidViewSelectionRejected runs-in DocumentBrowsing. | 990 |
| f0536 | InvalidViewSelectionRejected step 1 sends SelectViewRequest from XfmdDocumentHost to DocumentBroker via ViewNavigationCalls. | 991 |
| f0537 | InvalidViewSelectionRejected step 2 sends ViewRequestRejected from DocumentBroker to XfmdDocumentHost via ViewNavigationCalls reply-to 1. | 992 |
| f0538 | InvokeRegisteredFunction allocated-to FyneHost in mode BoundExecution. | 993 |
| f0539 | InvokeRegisteredFunction contributes-to TypedDomainBinding. | 994 |
| f0540 | InvokeRegisteredFunction realizes ExecutableDesign. | 995 |
| f0541 | KeepLastValidModels allocated-to FyneHost in mode LiveEditing. | 996 |
| f0542 | KeepLastValidModels contributes-to LiveModelReload. | 997 |
| f0543 | KeepLastValidModels realizes DevelopmentReload. | 998 |
| f0544 | KeyViewRevision allocated-to ViewServiceHost in mode DocumentBrowsing. | 999 |
| f0545 | KeyViewRevision contributes-to NavigableDesignDocumentation. | 1000 |
| f0546 | LaunchViewRequest has message-kind = request. | 1001 |
| f0547 | LaunchViewRequest upholds ViewOpenContract. | 1002 |
| f0548 | LayoutArguments has completeness = closed. | 1003 |
| f0549 | LayoutArguments has-field LayoutModelArtifact. | 1004 |
| f0550 | LayoutArguments has-field LayoutViewportHeight. | 1005 |
| f0551 | LayoutArguments has-field LayoutViewportWidth. | 1006 |
| f0552 | LayoutCalls upholds LayoutCallsProtocol. | 1007 |
| f0553 | LayoutCallsProtocol has completeness = closed. | 1008 |
| f0554 | LayoutCallsProtocol permits LayoutRequest. | 1009 |
| f0555 | LayoutCallsProtocol permits LayoutResult. | 1010 |
| f0556 | LayoutModelArtifact has presence = required. | 1011 |
| f0557 | LayoutModelArtifact has value-type = bytes. | 1012 |
| f0558 | LayoutRequest has message-kind = request. | 1013 |
| f0559 | LayoutRequest upholds LayoutArguments. | 1014 |
| f0560 | LayoutResult has message-kind = result. | 1015 |
| f0561 | LayoutResult replies-to LayoutRequest. | 1016 |
| f0562 | LayoutResult upholds PreparedFrameContract. | 1017 |
| f0563 | LayoutViewportHeight has presence = required. | 1018 |
| f0564 | LayoutViewportHeight has value-type = decimal. | 1019 |
| f0565 | LayoutViewportWidth has presence = required. | 1020 |
| f0566 | LayoutViewportWidth has value-type = decimal. | 1021 |
| f0567 | LiveModelReload supports EditRunningPrototype. | 1022 |
| f0568 | ManageDomainState has state-retention = stateful. | 1023 |
| f0569 | ManageDomainState realizes ExecutableDesign. | 1024 |
| f0570 | ManageWidgetIdentities has state-retention = stateful. | 1025 |
| f0571 | ManageWidgetIdentities realizes InteractiveSession. | 1026 |
| f0572 | MarkdownProvider consumes DiagramPort. | 1027 |
| f0573 | MarkdownProvider consumes MeasurementPort. | 1028 |
| f0574 | MarkdownProvider consumes ResourcePort. | 1029 |
| f0575 | MarkdownProvider owns MeasureMarkdownContent. | 1030 |
| f0576 | MarkdownProvider owns PrepareMarkdown. | 1031 |
| f0577 | MarkdownProvider provides RichContent. | 1032 |
| f0578 | MatchCompatibleWidgets realizes InteractiveSession. | 1033 |
| f0579 | MeasureMarkdownContent realizes RichContent. | 1034 |
| f0580 | MeasureUiContent realizes MeasuredPresentation. | 1035 |
| f0581 | MeasuredPresentation requires MeasurementPort in mode StaticExport. | 1036 |
| f0582 | MeasuredPresentation requires MeasurementPort in mode UiPreview. | 1037 |
| f0583 | MigrateOrResetDomainState realizes ExecutableDesign. | 1038 |
| f0584 | ModelDiagnostics has presence = optional. | 1039 |
| f0585 | ModelDiagnostics has value-type = text. | 1040 |
| f0586 | ModelIsValid has presence = required. | 1041 |
| f0587 | ModelIsValid has value-type = boolean. | 1042 |
| f0588 | ModelReloadCalls upholds ModelReloadCallsProtocol. | 1043 |
| f0589 | ModelReloadCallsProtocol has completeness = closed. | 1044 |
| f0590 | ModelReloadCallsProtocol permits ReloadPublished. | 1045 |
| f0591 | ModelReloadCallsProtocol permits ReloadRejected. | 1046 |
| f0592 | ModelReloadCallsProtocol permits ReloadRequest. | 1047 |
| f0593 | NativeBuildContract has completeness = closed. | 1048 |
| f0594 | NativeBuildContract has-field NativeBuildDiagnostics. | 1049 |
| f0595 | NativeBuildContract has-field NativeBuildSucceeded. | 1050 |
| f0596 | NativeBuildDiagnostics has presence = optional. | 1051 |
| f0597 | NativeBuildDiagnostics has value-type = text. | 1052 |
| f0598 | NativeBuildSucceeded has presence = required. | 1053 |
| f0599 | NativeBuildSucceeded has value-type = boolean. | 1054 |
| f0600 | NativeGoAssembly supports BuildNativeProduct. | 1055 |
| f0601 | NativeInteraction requires WidgetBackendPort in mode UiPreview. | 1056 |
| f0602 | NativeProgramBuilt exercises BuildNativeProduct. | 1057 |
| f0603 | NativeProgramBuilt has completeness = closed. | 1058 |
| f0604 | NativeProgramBuilt illustrates G5M1GeneratedGo. | 1059 |
| f0605 | NativeProgramBuilt runs-in NativeBuild. | 1060 |
| f0606 | NativeProgramBuilt step 1 sends GenerateGoRequest from CommandLineHost to GoCodeGenerator via GoGenerationCalls. | 1061 |
| f0607 | NativeProgramBuilt step 2 sends GenerateGoResult from GoCodeGenerator to CommandLineHost via GoGenerationCalls reply-to 1. | 1062 |
| f0608 | NativeProgramBuilt step 3 sends BuildGoRequest from CommandLineHost to GoBuildRunner via GoBuildCalls. | 1063 |
| f0609 | NativeProgramBuilt step 4 sends BuildGoResult from GoBuildRunner to CommandLineHost via GoBuildCalls reply-to 3. | 1064 |
| f0610 | NativeRealization requires BuildToolPort in mode NativeBuild. | 1065 |
| f0611 | NativeRealization requires GeneratedArtifactPort in mode NativeBuild. | 1066 |
| f0612 | NativeUiActions upholds NativeUiActionsProtocol. | 1067 |
| f0613 | NativeUiActionsProtocol has completeness = closed. | 1068 |
| f0614 | NativeUiActionsProtocol permits UiActionRejected. | 1069 |
| f0615 | NativeUiActionsProtocol permits UiActionRequest. | 1070 |
| f0616 | NativeUiActionsProtocol permits UiActionResult. | 1071 |
| f0617 | NavigableDesignDocumentation supports BrowseDesignViews. | 1072 |
| f0618 | NormalizeSdlModel has repeatability = deterministic. | 1073 |
| f0619 | NormalizeSdlModel has state-retention = stateless. | 1074 |
| f0620 | NormalizeSdlModel realizes SdlSourceModel. | 1075 |
| f0621 | NormalizeUiRequest has message-kind = request. | 1076 |
| f0622 | NormalizeUiRequest upholds AstArtifactContract. | 1077 |
| f0623 | NormalizeUiResult has message-kind = result. | 1078 |
| f0624 | NormalizeUiResult replies-to NormalizeUiRequest. | 1079 |
| f0625 | NormalizeUiResult upholds NormalizedModelContract. | 1080 |
| f0626 | NormalizedModelArtifact has presence = required. | 1081 |
| f0627 | NormalizedModelArtifact has value-type = bytes. | 1082 |
| f0628 | NormalizedModelContract has completeness = closed. | 1083 |
| f0629 | NormalizedModelContract has-field NormalizedModelArtifact. | 1084 |
| f0630 | NoticeGeneration has presence = required. | 1085 |
| f0631 | NoticeGeneration has value-type = unsigned. | 1086 |
| f0632 | NoticeVersion has presence = required. | 1087 |
| f0633 | NoticeVersion has value-type = unsigned. | 1088 |
| f0634 | ObserveSourceChanges allocated-to FyneHost in mode LiveEditing. | 1089 |
| f0635 | ObserveSourceChanges contributes-to LiveModelReload. | 1090 |
| f0636 | ObserveSourceChanges realizes DevelopmentReload. | 1091 |
| f0637 | OpenViewConsumerId has presence = required. | 1092 |
| f0638 | OpenViewConsumerId has value-type = text. | 1093 |
| f0639 | OpenViewEntryPath has presence = required. | 1094 |
| f0640 | OpenViewEntryPath has value-type = text. | 1095 |
| f0641 | OpenViewLeaseId has presence = required. | 1096 |
| f0642 | OpenViewLeaseId has value-type = text. | 1097 |
| f0643 | OpenViewPaneId has presence = required. | 1098 |
| f0644 | OpenViewPaneId has value-type = text. | 1099 |
| f0645 | OpenViewRequestId has presence = required. | 1100 |
| f0646 | OpenViewRequestId has value-type = text. | 1101 |
| f0647 | OpenViewRevision has presence = required. | 1102 |
| f0648 | OpenViewRevision has value-type = text. | 1103 |
| f0649 | OpenViewWindowId has presence = required. | 1104 |
| f0650 | OpenViewWindowId has value-type = text. | 1105 |
| f0651 | OpenedViewRequestId has presence = required. | 1106 |
| f0652 | OpenedViewRequestId has value-type = text. | 1107 |
| f0653 | OpenedViewRevision has presence = required. | 1108 |
| f0654 | OpenedViewRevision has value-type = text. | 1109 |
| f0655 | PerformDomainOperation realizes DomainOperations. | 1110 |
| f0656 | PrepareCandidateModels allocated-to FyneHost in mode LiveEditing. | 1111 |
| f0657 | PrepareCandidateModels contributes-to LiveModelReload. | 1112 |
| f0658 | PrepareCandidateModels realizes DevelopmentReload. | 1113 |
| f0659 | PrepareDiagramResource realizes RichContent. | 1114 |
| f0660 | PrepareMarkdown realizes RichContent. | 1115 |
| f0661 | PreparedFrameArtifact has presence = required. | 1116 |
| f0662 | PreparedFrameArtifact has value-type = bytes. | 1117 |
| f0663 | PreparedFrameContract has completeness = closed. | 1118 |
| f0664 | PreparedFrameContract has-field PreparedFrameArtifact. | 1119 |
| f0665 | PresentFrameRequest has message-kind = request. | 1120 |
| f0666 | PresentFrameRequest upholds PreparedFrameContract. | 1121 |
| f0667 | PresentFrameResult has message-kind = result. | 1122 |
| f0668 | PresentFrameResult replies-to PresentFrameRequest. | 1123 |
| f0669 | PresentFrameResult upholds PresentationOutcome. | 1124 |
| f0670 | PresentationOutcome has completeness = closed. | 1125 |
| f0671 | PresentationOutcome has-field PresentationReady. | 1126 |
| f0672 | PresentationReady has presence = required. | 1127 |
| f0673 | PresentationReady has value-type = boolean. | 1128 |
| f0674 | PreserveCompatibleUiState allocated-to FyneHost in mode LiveEditing. | 1129 |
| f0675 | PreserveCompatibleUiState contributes-to LiveModelReload. | 1130 |
| f0676 | PreserveCompatibleUiState realizes InteractiveSession. | 1131 |
| f0677 | PreserveHandwrittenSources allocated-to CommandLineHost in mode NativeBuild. | 1132 |
| f0678 | PreserveHandwrittenSources contributes-to NativeGoAssembly. | 1133 |
| f0679 | PreserveHandwrittenSources realizes NativeRealization. | 1134 |
| f0680 | PreserveSdlSourceMap has repeatability = deterministic. | 1135 |
| f0681 | PreserveSdlSourceMap has state-retention = stateless. | 1136 |
| f0682 | PreserveSdlSourceMap realizes SdlSourceModel. | 1137 |
| f0683 | PreserveUiRegions has repeatability = deterministic. | 1138 |
| f0684 | PreserveUiRegions has state-retention = stateless. | 1139 |
| f0685 | PreserveUiRegions realizes SduiSourceModel. | 1140 |
| f0686 | PreserveUiSourceMap has repeatability = deterministic. | 1141 |
| f0687 | PreserveUiSourceMap has state-retention = stateless. | 1142 |
| f0688 | PreserveUiSourceMap realizes SduiSourceModel. | 1143 |
| f0689 | PreserveViewAnchors allocated-to ViewServiceHost in mode DocumentBrowsing. | 1144 |
| f0690 | PreserveViewAnchors contributes-to NavigableDesignDocumentation. | 1145 |
| f0691 | ProjectSdlViewpoints allocated-to CommandLineHost in mode StaticExport. | 1146 |
| f0692 | ProjectSdlViewpoints allocated-to ViewServiceHost in mode DocumentBrowsing. | 1147 |
| f0693 | ProjectSdlViewpoints contributes-to DesignDocumentation. | 1148 |
| f0694 | ProjectSdlViewpoints has repeatability = deterministic. | 1149 |
| f0695 | ProjectSdlViewpoints has state-retention = stateless. | 1150 |
| f0696 | ProjectSelectedView allocated-to ViewServiceHost in mode DocumentBrowsing. | 1151 |
| f0697 | ProjectSelectedView contributes-to NavigableDesignDocumentation. | 1152 |
| f0698 | ProjectUiGeneration allocated-to FyneHost in mode LiveEditing. | 1153 |
| f0699 | ProjectUiGeneration contributes-to LiveModelReload. | 1154 |
| f0700 | ProjectUiGeneration projects UiSessionState into UiGenerationNotices. | 1155 |
| f0701 | ProjectViewRequest has message-kind = request. | 1156 |
| f0702 | ProjectViewRequest upholds ViewSelectionContract. | 1157 |
| f0703 | PublishDomainUpdates allocated-to FyneHost in mode BoundExecution. | 1158 |
| f0704 | PublishDomainUpdates contributes-to TypedDomainBinding. | 1159 |
| f0705 | PublishDomainUpdates realizes BoundInteraction. | 1160 |
| f0706 | PublishModelGeneration allocated-to FyneHost in mode LiveEditing. | 1161 |
| f0707 | PublishModelGeneration contributes-to LiveModelReload. | 1162 |
| f0708 | PublishModelGeneration has state-retention = stateful. | 1163 |
| f0709 | PublishModelGeneration realizes DevelopmentReload. | 1164 |
| f0710 | PublishPresentation realizes NativeInteraction. | 1165 |
| f0711 | PublishViewBundle allocated-to ViewServiceHost in mode DocumentBrowsing. | 1166 |
| f0712 | PublishViewBundle contributes-to NavigableDesignDocumentation. | 1167 |
| f0713 | PublishViewRequest has message-kind = request. | 1168 |
| f0714 | PublishViewRequest upholds ViewBundleContract. | 1169 |
| f0715 | ReadBoundedSources realizes SourceLoading. | 1170 |
| f0716 | ReconcileWidgets allocated-to FyneHost in mode UiPreview. | 1171 |
| f0717 | ReconcileWidgets contributes-to InteractiveUiPreview. | 1172 |
| f0718 | ReconcileWidgets has state-retention = stateful. | 1173 |
| f0719 | ReconcileWidgets realizes NativeInteraction. | 1174 |
| f0720 | RegisterDomainFunctions has state-retention = stateful. | 1175 |
| f0721 | RegisterDomainFunctions realizes ExecutableDesign. | 1176 |
| f0722 | RejectStaleUiEvent realizes InteractiveSession. | 1177 |
| f0723 | RejectStaleViewResults allocated-to ViewServiceHost in mode DocumentBrowsing. | 1178 |
| f0724 | RejectStaleViewResults contributes-to NavigableDesignDocumentation. | 1179 |
| f0725 | ReleaseNativeWidgets realizes NativeInteraction. | 1180 |
| f0726 | ReleaseVisualResources realizes RichContent. | 1181 |
| f0727 | ReloadArguments has completeness = closed. | 1182 |
| f0728 | ReloadArguments has-field ReloadSourceRevision. | 1183 |
| f0729 | ReloadArguments has-field ReloadSourceText. | 1184 |
| f0730 | ReloadBoundModels refines ReloadDesignSession. | 1185 |
| f0731 | ReloadCoordinator consumes BindingReloadPort. | 1186 |
| f0732 | ReloadCoordinator consumes DiagnosticPort. | 1187 |
| f0733 | ReloadCoordinator consumes SdlFrontendPort. | 1188 |
| f0734 | ReloadCoordinator consumes SdlReloadPort. | 1189 |
| f0735 | ReloadCoordinator consumes SduiFrontendPort. | 1190 |
| f0736 | ReloadCoordinator consumes SourceSnapshotPort. | 1191 |
| f0737 | ReloadCoordinator consumes UiReloadPort. | 1192 |
| f0738 | ReloadCoordinator owns KeepLastValidModels. | 1193 |
| f0739 | ReloadCoordinator owns PrepareCandidateModels. | 1194 |
| f0740 | ReloadCoordinator owns PublishModelGeneration. | 1195 |
| f0741 | ReloadCoordinator owns RetirePreviousGeneration. | 1196 |
| f0742 | ReloadCoordinator provides DevelopmentReload. | 1197 |
| f0743 | ReloadCoordinator uses ModelReloadCalls as receiver of ReloadRequest in mode BoundLiveEditing. | 1198 |
| f0744 | ReloadCoordinator uses ModelReloadCalls as receiver of ReloadRequest in mode LiveEditing. | 1199 |
| f0745 | ReloadCoordinator uses ModelReloadCalls as sender of ReloadPublished in mode BoundLiveEditing. | 1200 |
| f0746 | ReloadCoordinator uses ModelReloadCalls as sender of ReloadPublished in mode LiveEditing. | 1201 |
| f0747 | ReloadCoordinator uses ModelReloadCalls as sender of ReloadRejected in mode BoundLiveEditing. | 1202 |
| f0748 | ReloadCoordinator uses ModelReloadCalls as sender of ReloadRejected in mode LiveEditing. | 1203 |
| f0749 | ReloadCoordinator uses SdlCompilationCalls as receiver of CompileSdlResult in mode BoundLiveEditing. | 1204 |
| f0750 | ReloadCoordinator uses SdlCompilationCalls as sender of CompileSdlRequest in mode BoundLiveEditing. | 1205 |
| f0751 | ReloadCoordinator uses UiCompilationCalls as receiver of CompileUiRejected in mode LiveEditing. | 1206 |
| f0752 | ReloadCoordinator uses UiCompilationCalls as receiver of CompileUiResult in mode LiveEditing. | 1207 |
| f0753 | ReloadCoordinator uses UiCompilationCalls as sender of CompileUiRequest in mode LiveEditing. | 1208 |
| f0754 | ReloadDiagnostic has presence = optional. | 1209 |
| f0755 | ReloadDiagnostic has value-type = text. | 1210 |
| f0756 | ReloadOutcome has completeness = closed. | 1211 |
| f0757 | ReloadOutcome has-field ReloadDiagnostic. | 1212 |
| f0758 | ReloadOutcome has-field ReloadPublishedGeneration. | 1213 |
| f0759 | ReloadPublished has message-kind = result. | 1214 |
| f0760 | ReloadPublished replies-to ReloadRequest. | 1215 |
| f0761 | ReloadPublished upholds ReloadOutcome. | 1216 |
| f0762 | ReloadPublishedGeneration has presence = optional. | 1217 |
| f0763 | ReloadPublishedGeneration has value-type = unsigned. | 1218 |
| f0764 | ReloadRejected has message-kind = result. | 1219 |
| f0765 | ReloadRejected replies-to ReloadRequest. | 1220 |
| f0766 | ReloadRejected upholds ReloadOutcome. | 1221 |
| f0767 | ReloadRequest has message-kind = request. | 1222 |
| f0768 | ReloadRequest upholds ReloadArguments. | 1223 |
| f0769 | ReloadSourceRevision has presence = required. | 1224 |
| f0770 | ReloadSourceRevision has value-type = unsigned. | 1225 |
| f0771 | ReloadSourceText has presence = required. | 1226 |
| f0772 | ReloadSourceText has value-type = text. | 1227 |
| f0773 | ReloadUiModel refines ReloadDesignSession. | 1228 |
| f0774 | ReportBindingDiagnostics realizes SourceDiagnostics. | 1229 |
| f0775 | ReportSourceDiagnostics allocated-to CommandLineHost in mode SourceInspection. | 1230 |
| f0776 | ReportSourceDiagnostics allocated-to FyneHost in mode LiveEditing. | 1231 |
| f0777 | ReportSourceDiagnostics contributes-to LiveModelReload. | 1232 |
| f0778 | ReportSourceDiagnostics contributes-to StructuralModelInspection. | 1233 |
| f0779 | ReportSourceDiagnostics realizes SourceDiagnostics. | 1234 |
| f0780 | ResetIncompatibleUiState realizes InteractiveSession. | 1235 |
| f0781 | ResolveAncestorDimensions realizes MeasuredPresentation. | 1236 |
| f0782 | ResolveCallbackSymbols allocated-to FyneHost in mode BoundExecution. | 1237 |
| f0783 | ResolveCallbackSymbols contributes-to TypedDomainBinding. | 1238 |
| f0784 | ResolveCallbackSymbols realizes BoundInteraction. | 1239 |
| f0785 | ResolveConfiguredViewer allocated-to ViewServiceHost in mode DocumentBrowsing. | 1240 |
| f0786 | ResolveConfiguredViewer contributes-to NavigableDesignDocumentation. | 1241 |
| f0787 | ResolveSdlSymbols has repeatability = deterministic. | 1242 |
| f0788 | ResolveSdlSymbols has state-retention = stateless. | 1243 |
| f0789 | ResolveSdlSymbols realizes SdlSourceModel. | 1244 |
| f0790 | ResolveUiNames has repeatability = deterministic. | 1245 |
| f0791 | ResolveUiNames has state-retention = stateless. | 1246 |
| f0792 | ResolveUiNames realizes SduiSourceModel. | 1247 |
| f0793 | ResourceStore owns ReleaseVisualResources. | 1248 |
| f0794 | ResourceStore owns ValidateVisualResources. | 1249 |
| f0795 | ResourceStore provides RichContent. | 1250 |
| f0796 | RestartChangedGoProgram realizes NativeRealization. | 1251 |
| f0797 | RetainVisibleViewBundle allocated-to ViewServiceHost in mode DocumentBrowsing. | 1252 |
| f0798 | RetainVisibleViewBundle contributes-to NavigableDesignDocumentation. | 1253 |
| f0799 | RetirePreviousGeneration realizes DevelopmentReload. | 1254 |
| f0800 | RetireReplacedPythonEntryPoints allocated-to CommandLineHost in mode NativeBuild. | 1255 |
| f0801 | RetireReplacedPythonEntryPoints contributes-to NativeGoAssembly. | 1256 |
| f0802 | RevokeWidgetGenerations realizes InteractiveSession. | 1257 |
| f0803 | RichContent requires ContentProviderPort in mode RichDocument. | 1258 |
| f0804 | RichContent requires DiagramEnginePort in mode RichDocument. | 1259 |
| f0805 | RouteDocumentToPane allocated-to XfmdDocumentHost in mode DocumentBrowsing. | 1260 |
| f0806 | RouteDocumentToPane contributes-to NavigableDesignDocumentation. | 1261 |
| f0807 | RouteDomainBindings realizes BoundInteraction. | 1262 |
| f0808 | RunBoundUiAction refines RunDesignSession. | 1263 |
| f0809 | RunUnboundUiPreview refines RunDesignSession. | 1264 |
| f0810 | ScheduleUiPublication realizes NativeInteraction. | 1265 |
| f0811 | SdlActionCalls upholds SdlActionCallsProtocol. | 1266 |
| f0812 | SdlActionCallsProtocol has completeness = closed. | 1267 |
| f0813 | SdlActionCallsProtocol permits SdlActionRequest. | 1268 |
| f0814 | SdlActionCallsProtocol permits SdlActionResult. | 1269 |
| f0815 | SdlActionRequest has message-kind = request. | 1270 |
| f0816 | SdlActionRequest upholds ActionArguments. | 1271 |
| f0817 | SdlActionResult has message-kind = result. | 1272 |
| f0818 | SdlActionResult replies-to SdlActionRequest. | 1273 |
| f0819 | SdlActionResult upholds ActionOutcome. | 1274 |
| f0820 | SdlCompilationCalls upholds SdlCompilationCallsProtocol. | 1275 |
| f0821 | SdlCompilationCallsProtocol has completeness = closed. | 1276 |
| f0822 | SdlCompilationCallsProtocol permits CompileSdlRequest. | 1277 |
| f0823 | SdlCompilationCallsProtocol permits CompileSdlResult. | 1278 |
| f0824 | SdlDispatcher consumes DomainFunctionPort. | 1279 |
| f0825 | SdlDispatcher consumes DomainStatePort. | 1280 |
| f0826 | SdlDispatcher owns CancelPendingActions. | 1281 |
| f0827 | SdlDispatcher owns CorrelateActionResult. | 1282 |
| f0828 | SdlDispatcher owns InvokeRegisteredFunction. | 1283 |
| f0829 | SdlDispatcher owns ValidateActionInput. | 1284 |
| f0830 | SdlDispatcher provides ExecutableDesign. | 1285 |
| f0831 | SdlDispatcher uses GoDomainCalls as receiver of DomainActionResult in mode BoundExecution. | 1286 |
| f0832 | SdlDispatcher uses GoDomainCalls as sender of DomainActionRequest in mode BoundExecution. | 1287 |
| f0833 | SdlDispatcher uses SdlActionCalls as receiver of SdlActionRequest in mode BoundExecution. | 1288 |
| f0834 | SdlDispatcher uses SdlActionCalls as sender of SdlActionResult in mode BoundExecution. | 1289 |
| f0835 | SdlExecutionGate consumes DiagnosticPort. | 1290 |
| f0836 | SdlExecutionGate consumes SdlModelPort. | 1291 |
| f0837 | SdlExecutionGate owns CheckExecutionCompleteness. | 1292 |
| f0838 | SdlExecutionGate provides ExecutableDesign. | 1293 |
| f0839 | SdlFrontend consumes DiagnosticPort. | 1294 |
| f0840 | SdlFrontend consumes SourceSnapshotPort. | 1295 |
| f0841 | SdlFrontend contains SdlLexer. | 1296 |
| f0842 | SdlFrontend contains SdlNormalizer. | 1297 |
| f0843 | SdlFrontend contains SdlParser. | 1298 |
| f0844 | SdlFrontend contains SdlValidator. | 1299 |
| f0845 | SdlFrontend owns CoordinateSdlCompilation. | 1300 |
| f0846 | SdlFrontend provides SdlSourceModel. | 1301 |
| f0847 | SdlFrontend uses SdlCompilationCalls as receiver of CompileSdlRequest in mode BoundLiveEditing. | 1302 |
| f0848 | SdlFrontend uses SdlCompilationCalls as sender of CompileSdlResult in mode BoundLiveEditing. | 1303 |
| f0849 | SdlFunctionRegistry consumes DomainFunctionPort. | 1304 |
| f0850 | SdlFunctionRegistry owns CheckFunctionSignatures. | 1305 |
| f0851 | SdlFunctionRegistry owns RegisterDomainFunctions. | 1306 |
| f0852 | SdlFunctionRegistry provides ExecutableDesign. | 1307 |
| f0853 | SdlLexer owns TokenizeSdlSource. | 1308 |
| f0854 | SdlLexer provides SdlSourceModel. | 1309 |
| f0855 | SdlLibrary contains SdlFrontend. | 1310 |
| f0856 | SdlLibrary contains SdlRuntime. | 1311 |
| f0857 | SdlLibrary contains SdlViewpointGenerator. | 1312 |
| f0858 | SdlLibrary provides ExecutableDesign. | 1313 |
| f0859 | SdlLibrary provides SdlSourceModel. | 1314 |
| f0860 | SdlModelReloadAccepted exercises EditRunningPrototype. | 1315 |
| f0861 | SdlModelReloadAccepted has completeness = closed. | 1316 |
| f0862 | SdlModelReloadAccepted illustrates G4M1SdlFrontend. | 1317 |
| f0863 | SdlModelReloadAccepted illustrates G4M4DomainReload. | 1318 |
| f0864 | SdlModelReloadAccepted runs-in BoundLiveEditing. | 1319 |
| f0865 | SdlModelReloadAccepted step 1 sends ReloadRequest from SourceWatcher to ReloadCoordinator via ModelReloadCalls. | 1320 |
| f0866 | SdlModelReloadAccepted step 2 sends CompileSdlRequest from ReloadCoordinator to SdlFrontend via SdlCompilationCalls. | 1321 |
| f0867 | SdlModelReloadAccepted step 3 sends CompileSdlResult from SdlFrontend to ReloadCoordinator via SdlCompilationCalls reply-to 2. | 1322 |
| f0868 | SdlModelReloadAccepted step 4 sends ReloadPublished from ReloadCoordinator to SourceWatcher via ModelReloadCalls reply-to 1. | 1323 |
| f0869 | SdlNormalizer owns NormalizeSdlModel. | 1324 |
| f0870 | SdlNormalizer owns PreserveSdlSourceMap. | 1325 |
| f0871 | SdlNormalizer provides SdlSourceModel. | 1326 |
| f0872 | SdlParser owns BuildSdlAst. | 1327 |
| f0873 | SdlParser provides SdlSourceModel. | 1328 |
| f0874 | SdlRuntime consumes DomainFunctionPort. | 1329 |
| f0875 | SdlRuntime consumes SdlModelPort. | 1330 |
| f0876 | SdlRuntime contains DomainStateMigrator. | 1331 |
| f0877 | SdlRuntime contains SdlDispatcher. | 1332 |
| f0878 | SdlRuntime contains SdlExecutionGate. | 1333 |
| f0879 | SdlRuntime contains SdlFunctionRegistry. | 1334 |
| f0880 | SdlRuntime contains SdlStateStore. | 1335 |
| f0881 | SdlRuntime owns CloseSdlInstance. | 1336 |
| f0882 | SdlRuntime owns CreateSdlInstance. | 1337 |
| f0883 | SdlRuntime provides ExecutableDesign. | 1338 |
| f0884 | SdlSourceModel requires SourceSnapshotPort in mode SourceInspection. | 1339 |
| f0885 | SdlStateStore owns ManageDomainState. | 1340 |
| f0886 | SdlStateStore owns SnapshotDomainState. | 1341 |
| f0887 | SdlStateStore provides ExecutableDesign. | 1342 |
| f0888 | SdlUiBindingAdapter consumes DiagnosticPort. | 1343 |
| f0889 | SdlUiBindingAdapter consumes SdlExecutionPort. | 1344 |
| f0890 | SdlUiBindingAdapter consumes UiSessionPort. | 1345 |
| f0891 | SdlUiBindingAdapter owns ConnectTypedWidgetHandles. | 1346 |
| f0892 | SdlUiBindingAdapter owns DisconnectBindings. | 1347 |
| f0893 | SdlUiBindingAdapter owns PublishDomainUpdates. | 1348 |
| f0894 | SdlUiBindingAdapter owns ResolveCallbackSymbols. | 1349 |
| f0895 | SdlUiBindingAdapter owns RouteDomainBindings. | 1350 |
| f0896 | SdlUiBindingAdapter provides BoundInteraction. | 1351 |
| f0897 | SdlUiBindingAdapter uses SdlActionCalls as receiver of SdlActionResult in mode BoundExecution. | 1352 |
| f0898 | SdlUiBindingAdapter uses SdlActionCalls as sender of SdlActionRequest in mode BoundExecution. | 1353 |
| f0899 | SdlUiBindingAdapter uses UiDomainActions as receiver of BoundActionRequest in mode BoundExecution. | 1354 |
| f0900 | SdlUiBindingAdapter uses UiDomainActions as sender of BoundActionResult in mode BoundExecution. | 1355 |
| f0901 | SdlValidator owns ResolveSdlSymbols. | 1356 |
| f0902 | SdlValidator owns ValidateSdlProfile. | 1357 |
| f0903 | SdlValidator owns ValidateSdlStructure. | 1358 |
| f0904 | SdlValidator provides SdlSourceModel. | 1359 |
| f0905 | SdlViewpointGenerator owns ComposeViewPackage. | 1360 |
| f0906 | SdlViewpointGenerator owns ExportViewpointMarkdown. | 1361 |
| f0907 | SdlViewpointGenerator owns GenerateViewNavigation. | 1362 |
| f0908 | SdlViewpointGenerator owns PreserveViewAnchors. | 1363 |
| f0909 | SdlViewpointGenerator owns ProjectSdlViewpoints. | 1364 |
| f0910 | SdlViewpointGenerator owns ProjectSelectedView. | 1365 |
| f0911 | SdlViewpointGenerator owns TraceViewpointFacts. | 1366 |
| f0912 | SdlViewpointGenerator uses ViewProjectionCalls as receiver of ProjectViewRequest in mode DocumentBrowsing. | 1367 |
| f0913 | SdlViewpointGenerator uses ViewProjectionCalls as sender of ViewBundleResult in mode DocumentBrowsing. | 1368 |
| f0914 | SdlViewpointGenerator uses ViewProjectionCalls as sender of ViewProjectionRejected in mode DocumentBrowsing. | 1369 |
| f0915 | SduiDispatcher consumes DomainBindingPort. | 1370 |
| f0916 | SduiDispatcher consumes UiStatePort. | 1371 |
| f0917 | SduiDispatcher owns CorrelateUiResult. | 1372 |
| f0918 | SduiDispatcher owns DispatchUiEvent. | 1373 |
| f0919 | SduiDispatcher owns RejectStaleUiEvent. | 1374 |
| f0920 | SduiDispatcher owns ValidateUiEvent. | 1375 |
| f0921 | SduiDispatcher provides InteractiveSession. | 1376 |
| f0922 | SduiDispatcher uses NativeUiActions as receiver of UiActionRequest in mode BoundExecution. | 1377 |
| f0923 | SduiDispatcher uses NativeUiActions as sender of UiActionRejected in mode BoundExecution. | 1378 |
| f0924 | SduiDispatcher uses NativeUiActions as sender of UiActionResult in mode BoundExecution. | 1379 |
| f0925 | SduiDispatcher uses UiDomainActions as receiver of BoundActionResult in mode BoundExecution. | 1380 |
| f0926 | SduiDispatcher uses UiDomainActions as sender of BoundActionRequest in mode BoundExecution. | 1381 |
| f0927 | SduiFrontend consumes DiagnosticPort. | 1382 |
| f0928 | SduiFrontend consumes SourceSnapshotPort. | 1383 |
| f0929 | SduiFrontend contains SduiLexer. | 1384 |
| f0930 | SduiFrontend contains SduiNormalizer. | 1385 |
| f0931 | SduiFrontend contains SduiParser. | 1386 |
| f0932 | SduiFrontend contains SduiValidator. | 1387 |
| f0933 | SduiFrontend owns CoordinateSduiCompilation. | 1388 |
| f0934 | SduiFrontend provides SduiSourceModel. | 1389 |
| f0935 | SduiFrontend uses UiAstCalls as receiver of BuildUiAstResult in mode SourceInspection. | 1390 |
| f0936 | SduiFrontend uses UiAstCalls as sender of BuildUiAstRequest in mode SourceInspection. | 1391 |
| f0937 | SduiFrontend uses UiCompilationCalls as receiver of CompileUiRequest in mode LiveEditing. | 1392 |
| f0938 | SduiFrontend uses UiCompilationCalls as receiver of CompileUiRequest in mode SourceInspection. | 1393 |
| f0939 | SduiFrontend uses UiCompilationCalls as sender of CompileUiRejected in mode LiveEditing. | 1394 |
| f0940 | SduiFrontend uses UiCompilationCalls as sender of CompileUiRejected in mode SourceInspection. | 1395 |
| f0941 | SduiFrontend uses UiCompilationCalls as sender of CompileUiResult in mode LiveEditing. | 1396 |
| f0942 | SduiFrontend uses UiCompilationCalls as sender of CompileUiResult in mode SourceInspection. | 1397 |
| f0943 | SduiFrontend uses UiNormalizationCalls as receiver of NormalizeUiResult in mode SourceInspection. | 1398 |
| f0944 | SduiFrontend uses UiNormalizationCalls as sender of NormalizeUiRequest in mode SourceInspection. | 1399 |
| f0945 | SduiFrontend uses UiTokenizationCalls as receiver of TokenizeUiResult in mode SourceInspection. | 1400 |
| f0946 | SduiFrontend uses UiTokenizationCalls as sender of TokenizeUiRequest in mode SourceInspection. | 1401 |
| f0947 | SduiFrontend uses UiValidationCalls as receiver of ValidateUiResult in mode SourceInspection. | 1402 |
| f0948 | SduiFrontend uses UiValidationCalls as sender of ValidateUiRequest in mode SourceInspection. | 1403 |
| f0949 | SduiInstanceStore holds UiSessionState. | 1404 |
| f0950 | SduiInstanceStore owns ManageWidgetIdentities. | 1405 |
| f0951 | SduiInstanceStore owns ProjectUiGeneration. | 1406 |
| f0952 | SduiInstanceStore owns RevokeWidgetGenerations. | 1407 |
| f0953 | SduiInstanceStore owns SnapshotUiState. | 1408 |
| f0954 | SduiInstanceStore provides InteractiveSession. | 1409 |
| f0955 | SduiInstanceStore uses UiGenerationEvents as sender of UiGenerationNotices in mode BoundExecution. | 1410 |
| f0956 | SduiInstanceStore uses UiGenerationEvents as sender of UiGenerationNotices in mode LiveEditing. | 1411 |
| f0957 | SduiLayout consumes ContentProviderPort. | 1412 |
| f0958 | SduiLayout consumes MeasurementPort. | 1413 |
| f0959 | SduiLayout consumes UiSnapshotPort. | 1414 |
| f0960 | SduiLayout owns AllocateGeometry. | 1415 |
| f0961 | SduiLayout owns BuildPreparedFrame. | 1416 |
| f0962 | SduiLayout owns ComputeClipping. | 1417 |
| f0963 | SduiLayout owns MeasureUiContent. | 1418 |
| f0964 | SduiLayout owns ResolveAncestorDimensions. | 1419 |
| f0965 | SduiLayout provides MeasuredPresentation. | 1420 |
| f0966 | SduiLayout uses LayoutCalls as receiver of LayoutRequest in mode StaticExport. | 1421 |
| f0967 | SduiLayout uses LayoutCalls as receiver of LayoutRequest in mode UiPreview. | 1422 |
| f0968 | SduiLayout uses LayoutCalls as sender of LayoutResult in mode StaticExport. | 1423 |
| f0969 | SduiLayout uses LayoutCalls as sender of LayoutResult in mode UiPreview. | 1424 |
| f0970 | SduiLexer owns TokenizeSduiSource. | 1425 |
| f0971 | SduiLexer provides SduiSourceModel. | 1426 |
| f0972 | SduiLexer uses UiTokenizationCalls as receiver of TokenizeUiRequest in mode SourceInspection. | 1427 |
| f0973 | SduiLexer uses UiTokenizationCalls as sender of TokenizeUiResult in mode SourceInspection. | 1428 |
| f0974 | SduiLibrary contains SduiFrontend. | 1429 |
| f0975 | SduiLibrary contains SduiLayout. | 1430 |
| f0976 | SduiLibrary contains SduiPresentation. | 1431 |
| f0977 | SduiLibrary contains SduiRuntime. | 1432 |
| f0978 | SduiLibrary provides InteractiveSession. | 1433 |
| f0979 | SduiLibrary provides MeasuredPresentation. | 1434 |
| f0980 | SduiLibrary provides SduiSourceModel. | 1435 |
| f0981 | SduiLibrary provides StaticDocumentation. | 1436 |
| f0982 | SduiNormalizer owns ExpandUiDefinitions. | 1437 |
| f0983 | SduiNormalizer owns PreserveUiRegions. | 1438 |
| f0984 | SduiNormalizer owns PreserveUiSourceMap. | 1439 |
| f0985 | SduiNormalizer provides SduiSourceModel. | 1440 |
| f0986 | SduiNormalizer uses UiNormalizationCalls as receiver of NormalizeUiRequest in mode SourceInspection. | 1441 |
| f0987 | SduiNormalizer uses UiNormalizationCalls as sender of NormalizeUiResult in mode SourceInspection. | 1442 |
| f0988 | SduiParser owns BuildSduiAst. | 1443 |
| f0989 | SduiParser provides SduiSourceModel. | 1444 |
| f0990 | SduiParser uses UiAstCalls as receiver of BuildUiAstRequest in mode SourceInspection. | 1445 |
| f0991 | SduiParser uses UiAstCalls as sender of BuildUiAstResult in mode SourceInspection. | 1446 |
| f0992 | SduiPresentation consumes PreparedFramePort. | 1447 |
| f0993 | SduiPresentation owns ComposeMarkdownDocument. | 1448 |
| f0994 | SduiPresentation owns ExportConsoleSnapshot. | 1449 |
| f0995 | SduiPresentation owns ExportSvgSnapshot. | 1450 |
| f0996 | SduiPresentation provides StaticDocumentation. | 1451 |
| f0997 | SduiPresentation uses SvgExportCalls as receiver of ExportSvgRequest in mode StaticExport. | 1452 |
| f0998 | SduiPresentation uses SvgExportCalls as sender of ExportSvgResult in mode StaticExport. | 1453 |
| f0999 | SduiPropertyStore consumes UiStatePort. | 1454 |
| f1000 | SduiPropertyStore owns ApplyPropertyBatch. | 1455 |
| f1001 | SduiPropertyStore owns TrackInputDraft. | 1456 |
| f1002 | SduiPropertyStore owns ValidatePropertyBatch. | 1457 |
| f1003 | SduiPropertyStore provides InteractiveSession. | 1458 |
| f1004 | SduiRuntime consumes DomainBindingPort. | 1459 |
| f1005 | SduiRuntime consumes SduiModelPort. | 1460 |
| f1006 | SduiRuntime contains SduiDispatcher. | 1461 |
| f1007 | SduiRuntime contains SduiInstanceStore. | 1462 |
| f1008 | SduiRuntime contains SduiPropertyStore. | 1463 |
| f1009 | SduiRuntime contains UiStateReconciler. | 1464 |
| f1010 | SduiRuntime owns CloseUiInstance. | 1465 |
| f1011 | SduiRuntime owns CreateUiInstance. | 1466 |
| f1012 | SduiRuntime provides InteractiveSession. | 1467 |
| f1013 | SduiSourceModel requires SourceSnapshotPort in mode SourceInspection. | 1468 |
| f1014 | SduiValidator owns ResolveUiNames. | 1469 |
| f1015 | SduiValidator owns ValidateRelativeFormatting. | 1470 |
| f1016 | SduiValidator owns ValidateSymbolicBindings. | 1471 |
| f1017 | SduiValidator owns ValidateWidgetArguments. | 1472 |
| f1018 | SduiValidator provides SduiSourceModel. | 1473 |
| f1019 | SduiValidator uses UiValidationCalls as receiver of ValidateUiRequest in mode SourceInspection. | 1474 |
| f1020 | SduiValidator uses UiValidationCalls as sender of ValidateUiResult in mode SourceInspection. | 1475 |
| f1021 | SelectViewRequest has message-kind = request. | 1476 |
| f1022 | SelectViewRequest upholds ViewSelectionContract. | 1477 |
| f1023 | SelectedViewOpened exercises BrowseDesignViews. | 1478 |
| f1024 | SelectedViewOpened has completeness = closed. | 1479 |
| f1025 | SelectedViewOpened illustrates G6M1StaticNavigation. | 1480 |
| f1026 | SelectedViewOpened illustrates G6M2OnDemandViews. | 1481 |
| f1027 | SelectedViewOpened illustrates G6M3XfmdNavigation. | 1482 |
| f1028 | SelectedViewOpened illustrates G6M4SessionPublication. | 1483 |
| f1029 | SelectedViewOpened runs-in DocumentBrowsing. | 1484 |
| f1030 | SelectedViewOpened step 1 sends SelectViewRequest from XfmdDocumentHost to DocumentBroker via ViewNavigationCalls. | 1485 |
| f1031 | SelectedViewOpened step 10 sends ViewOpenedResult from DocumentBroker to XfmdDocumentHost via ViewNavigationCalls reply-to 1. | 1486 |
| f1032 | SelectedViewOpened step 2 sends ProjectViewRequest from DocumentBroker to SdlViewpointGenerator via ViewProjectionCalls. | 1487 |
| f1033 | SelectedViewOpened step 3 sends ViewBundleResult from SdlViewpointGenerator to DocumentBroker via ViewProjectionCalls reply-to 2. | 1488 |
| f1034 | SelectedViewOpened step 4 sends PublishViewRequest from DocumentBroker to ViewArtifactStore via ViewPublicationCalls. | 1489 |
| f1035 | SelectedViewOpened step 5 sends ViewReferenceResult from ViewArtifactStore to DocumentBroker via ViewPublicationCalls reply-to 4. | 1490 |
| f1036 | SelectedViewOpened step 6 sends LaunchViewRequest from DocumentBroker to ViewerLaunchAdapter via ViewLaunchCalls. | 1491 |
| f1037 | SelectedViewOpened step 7 sends DisplayViewRequest from ViewerLaunchAdapter to XfmdDocumentHost via ViewDisplayCalls. | 1492 |
| f1038 | SelectedViewOpened step 8 sends ViewDisplayResult from XfmdDocumentHost to ViewerLaunchAdapter via ViewDisplayCalls reply-to 7. | 1493 |
| f1039 | SelectedViewOpened step 9 sends ViewLaunchResult from ViewerLaunchAdapter to DocumentBroker via ViewLaunchCalls reply-to 6. | 1494 |
| f1040 | ServeViewRequests allocated-to ViewServiceHost in mode DocumentBrowsing. | 1495 |
| f1041 | ServeViewRequests contributes-to NavigableDesignDocumentation. | 1496 |
| f1042 | SessionDraft has presence = optional. | 1497 |
| f1043 | SessionDraft has value-type = text. | 1498 |
| f1044 | SessionGeneration has presence = required. | 1499 |
| f1045 | SessionGeneration has value-type = unsigned. | 1500 |
| f1046 | SnapshotDomainState realizes ExecutableDesign. | 1501 |
| f1047 | SnapshotUiState realizes InteractiveSession. | 1502 |
| f1048 | SourceDocumentRevision has presence = required. | 1503 |
| f1049 | SourceDocumentRevision has value-type = unsigned. | 1504 |
| f1050 | SourceDocumentText has presence = required. | 1505 |
| f1051 | SourceDocumentText has value-type = text. | 1506 |
| f1052 | SourceLoader consumes SourceInputPort. | 1507 |
| f1053 | SourceLoader owns DesignSourceArchive. | 1508 |
| f1054 | SourceLoader owns IdentifySourceRevision. | 1509 |
| f1055 | SourceLoader owns ReadBoundedSources. | 1510 |
| f1056 | SourceLoader provides SourceLoading. | 1511 |
| f1057 | SourceWatcher consumes FileChangePort. | 1512 |
| f1058 | SourceWatcher owns CoalesceSourceChanges. | 1513 |
| f1059 | SourceWatcher owns ObserveSourceChanges. | 1514 |
| f1060 | SourceWatcher provides DevelopmentReload. | 1515 |
| f1061 | SourceWatcher uses ModelReloadCalls as receiver of ReloadPublished in mode BoundLiveEditing. | 1516 |
| f1062 | SourceWatcher uses ModelReloadCalls as receiver of ReloadPublished in mode LiveEditing. | 1517 |
| f1063 | SourceWatcher uses ModelReloadCalls as receiver of ReloadRejected in mode BoundLiveEditing. | 1518 |
| f1064 | SourceWatcher uses ModelReloadCalls as receiver of ReloadRejected in mode LiveEditing. | 1519 |
| f1065 | SourceWatcher uses ModelReloadCalls as sender of ReloadRequest in mode BoundLiveEditing. | 1520 |
| f1066 | SourceWatcher uses ModelReloadCalls as sender of ReloadRequest in mode LiveEditing. | 1521 |
| f1067 | StaticDocumentation requires ExportSinkPort in mode StaticExport. | 1522 |
| f1068 | StaticDocumentation requires PreparedFramePort in mode StaticExport. | 1523 |
| f1069 | StaticFrameExported exercises PublishDesignDocumentation. | 1524 |
| f1070 | StaticFrameExported has completeness = closed. | 1525 |
| f1071 | StaticFrameExported illustrates G2M2SharedSvgGeometry. | 1526 |
| f1072 | StaticFrameExported illustrates G5M3DocumentationExport. | 1527 |
| f1073 | StaticFrameExported runs-in StaticExport. | 1528 |
| f1074 | StaticFrameExported step 1 sends LayoutRequest from CommandLineHost to SduiLayout via LayoutCalls. | 1529 |
| f1075 | StaticFrameExported step 2 sends LayoutResult from SduiLayout to CommandLineHost via LayoutCalls reply-to 1. | 1530 |
| f1076 | StaticFrameExported step 3 sends ExportSvgRequest from CommandLineHost to SduiPresentation via SvgExportCalls. | 1531 |
| f1077 | StaticFrameExported step 4 sends ExportSvgResult from SduiPresentation to CommandLineHost via SvgExportCalls reply-to 3. | 1532 |
| f1078 | StructuralModelInspection supports EditRunningPrototype. | 1533 |
| f1079 | StructuralModelInspection supports InspectModels. | 1534 |
| f1080 | SvgDocumentContract has completeness = closed. | 1535 |
| f1081 | SvgDocumentContract has-field SvgDocumentText. | 1536 |
| f1082 | SvgDocumentText has presence = required. | 1537 |
| f1083 | SvgDocumentText has value-type = text. | 1538 |
| f1084 | SvgExportCalls upholds SvgExportCallsProtocol. | 1539 |
| f1085 | SvgExportCallsProtocol has completeness = closed. | 1540 |
| f1086 | SvgExportCallsProtocol permits ExportSvgRequest. | 1541 |
| f1087 | SvgExportCallsProtocol permits ExportSvgResult. | 1542 |
| f1088 | TokenArtifact has presence = required. | 1543 |
| f1089 | TokenArtifact has value-type = bytes. | 1544 |
| f1090 | TokenArtifactContract has completeness = closed. | 1545 |
| f1091 | TokenArtifactContract has-field TokenArtifact. | 1546 |
| f1092 | TokenizeSdlSource has repeatability = deterministic. | 1547 |
| f1093 | TokenizeSdlSource has state-retention = stateless. | 1548 |
| f1094 | TokenizeSdlSource realizes SdlSourceModel. | 1549 |
| f1095 | TokenizeSduiSource has repeatability = deterministic. | 1550 |
| f1096 | TokenizeSduiSource has state-retention = stateless. | 1551 |
| f1097 | TokenizeSduiSource realizes SduiSourceModel. | 1552 |
| f1098 | TokenizeUiRequest has message-kind = request. | 1553 |
| f1099 | TokenizeUiRequest upholds ReloadArguments. | 1554 |
| f1100 | TokenizeUiResult has message-kind = result. | 1555 |
| f1101 | TokenizeUiResult replies-to TokenizeUiRequest. | 1556 |
| f1102 | TokenizeUiResult upholds TokenArtifactContract. | 1557 |
| f1103 | TraceViewpointFacts allocated-to CommandLineHost in mode SourceInspection. | 1558 |
| f1104 | TraceViewpointFacts allocated-to CommandLineHost in mode StaticExport. | 1559 |
| f1105 | TraceViewpointFacts allocated-to ViewServiceHost in mode DocumentBrowsing. | 1560 |
| f1106 | TraceViewpointFacts contributes-to DesignDocumentation. | 1561 |
| f1107 | TraceViewpointFacts contributes-to InspectModels. | 1562 |
| f1108 | TraceViewpointFacts has repeatability = deterministic. | 1563 |
| f1109 | TraceViewpointFacts has state-retention = stateless. | 1564 |
| f1110 | TrackInputDraft has state-retention = stateful. | 1565 |
| f1111 | TrackInputDraft realizes InteractiveSession. | 1566 |
| f1112 | TypedDomainBinding supports TryDomainInteraction. | 1567 |
| f1113 | UiActionRejected has message-kind = result. | 1568 |
| f1114 | UiActionRejected replies-to UiActionRequest. | 1569 |
| f1115 | UiActionRejected upholds ActionOutcome. | 1570 |
| f1116 | UiActionRequest has message-kind = request. | 1571 |
| f1117 | UiActionRequest upholds ActionArguments. | 1572 |
| f1118 | UiActionResult has message-kind = result. | 1573 |
| f1119 | UiActionResult replies-to UiActionRequest. | 1574 |
| f1120 | UiActionResult upholds ActionOutcome. | 1575 |
| f1121 | UiAstCalls upholds UiAstCallsProtocol. | 1576 |
| f1122 | UiAstCallsProtocol has completeness = closed. | 1577 |
| f1123 | UiAstCallsProtocol permits BuildUiAstRequest. | 1578 |
| f1124 | UiAstCallsProtocol permits BuildUiAstResult. | 1579 |
| f1125 | UiCompilationAccepted exercises InspectModels. | 1580 |
| f1126 | UiCompilationAccepted has completeness = closed. | 1581 |
| f1127 | UiCompilationAccepted illustrates G1M1ParserAndAst. | 1582 |
| f1128 | UiCompilationAccepted illustrates G1M2ValidationAndNormalization. | 1583 |
| f1129 | UiCompilationAccepted runs-in SourceInspection. | 1584 |
| f1130 | UiCompilationAccepted step 1 sends CompileUiRequest from CommandLineHost to SduiFrontend via UiCompilationCalls. | 1585 |
| f1131 | UiCompilationAccepted step 10 sends CompileUiResult from SduiFrontend to CommandLineHost via UiCompilationCalls reply-to 1. | 1586 |
| f1132 | UiCompilationAccepted step 2 sends TokenizeUiRequest from SduiFrontend to SduiLexer via UiTokenizationCalls. | 1587 |
| f1133 | UiCompilationAccepted step 3 sends TokenizeUiResult from SduiLexer to SduiFrontend via UiTokenizationCalls reply-to 2. | 1588 |
| f1134 | UiCompilationAccepted step 4 sends BuildUiAstRequest from SduiFrontend to SduiParser via UiAstCalls. | 1589 |
| f1135 | UiCompilationAccepted step 5 sends BuildUiAstResult from SduiParser to SduiFrontend via UiAstCalls reply-to 4. | 1590 |
| f1136 | UiCompilationAccepted step 6 sends ValidateUiRequest from SduiFrontend to SduiValidator via UiValidationCalls. | 1591 |
| f1137 | UiCompilationAccepted step 7 sends ValidateUiResult from SduiValidator to SduiFrontend via UiValidationCalls reply-to 6. | 1592 |
| f1138 | UiCompilationAccepted step 8 sends NormalizeUiRequest from SduiFrontend to SduiNormalizer via UiNormalizationCalls. | 1593 |
| f1139 | UiCompilationAccepted step 9 sends NormalizeUiResult from SduiNormalizer to SduiFrontend via UiNormalizationCalls reply-to 8. | 1594 |
| f1140 | UiCompilationCalls upholds UiCompilationCallsProtocol. | 1595 |
| f1141 | UiCompilationCallsProtocol has completeness = closed. | 1596 |
| f1142 | UiCompilationCallsProtocol permits CompileUiRejected. | 1597 |
| f1143 | UiCompilationCallsProtocol permits CompileUiRequest. | 1598 |
| f1144 | UiCompilationCallsProtocol permits CompileUiResult. | 1599 |
| f1145 | UiDomainActions upholds UiDomainActionsProtocol. | 1600 |
| f1146 | UiDomainActionsProtocol has completeness = closed. | 1601 |
| f1147 | UiDomainActionsProtocol permits BoundActionRequest. | 1602 |
| f1148 | UiDomainActionsProtocol permits BoundActionResult. | 1603 |
| f1149 | UiGenerationChanged has-field NoticeGeneration. | 1604 |
| f1150 | UiGenerationContract defines UiGenerationChanged. | 1605 |
| f1151 | UiGenerationContract has completeness = closed. | 1606 |
| f1152 | UiGenerationContract has-field NoticeVersion. | 1607 |
| f1153 | UiGenerationEvents upholds UiGenerationEventsProtocol. | 1608 |
| f1154 | UiGenerationEventsProtocol has completeness = closed. | 1609 |
| f1155 | UiGenerationEventsProtocol permits UiGenerationNotices. | 1610 |
| f1156 | UiGenerationNotices from UiSessionState. | 1611 |
| f1157 | UiGenerationNotices upholds UiGenerationContract. | 1612 |
| f1158 | UiGenerationWire encodes UiGenerationChanged. | 1613 |
| f1159 | UiGenerationWire has bit-order = most-significant-first. | 1614 |
| f1160 | UiGenerationWire has byte-order = big-endian. | 1615 |
| f1161 | UiGenerationWire places NoticeGeneration at 16 bits 64. | 1616 |
| f1162 | UiGenerationWire places NoticeVersion at 0 bits 16. | 1617 |
| f1163 | UiModelReloadAccepted exercises EditRunningPrototype. | 1618 |
| f1164 | UiModelReloadAccepted has completeness = closed. | 1619 |
| f1165 | UiModelReloadAccepted illustrates G3M2CandidatePublication. | 1620 |
| f1166 | UiModelReloadAccepted illustrates G3M3CompatibleState. | 1621 |
| f1167 | UiModelReloadAccepted runs-in LiveEditing. | 1622 |
| f1168 | UiModelReloadAccepted step 1 sends ReloadRequest from SourceWatcher to ReloadCoordinator via ModelReloadCalls. | 1623 |
| f1169 | UiModelReloadAccepted step 2 sends CompileUiRequest from ReloadCoordinator to SduiFrontend via UiCompilationCalls. | 1624 |
| f1170 | UiModelReloadAccepted step 3 sends CompileUiResult from SduiFrontend to ReloadCoordinator via UiCompilationCalls reply-to 2. | 1625 |
| f1171 | UiModelReloadAccepted step 4 sends ReloadPublished from ReloadCoordinator to SourceWatcher via ModelReloadCalls reply-to 1. | 1626 |
| f1172 | UiModelReloadAccepted step 5 sends UiGenerationNotices variant UiGenerationChanged from SduiInstanceStore to FyneBackend via UiGenerationEvents. | 1627 |
| f1173 | UiModelReloadRejected exercises EditRunningPrototype. | 1628 |
| f1174 | UiModelReloadRejected has completeness = closed. | 1629 |
| f1175 | UiModelReloadRejected illustrates G3M2CandidatePublication. | 1630 |
| f1176 | UiModelReloadRejected runs-in LiveEditing. | 1631 |
| f1177 | UiModelReloadRejected step 1 sends ReloadRequest from SourceWatcher to ReloadCoordinator via ModelReloadCalls. | 1632 |
| f1178 | UiModelReloadRejected step 2 sends CompileUiRequest from ReloadCoordinator to SduiFrontend via UiCompilationCalls. | 1633 |
| f1179 | UiModelReloadRejected step 3 sends CompileUiRejected from SduiFrontend to ReloadCoordinator via UiCompilationCalls reply-to 2. | 1634 |
| f1180 | UiModelReloadRejected step 4 sends ReloadRejected from ReloadCoordinator to SourceWatcher via ModelReloadCalls reply-to 1. | 1635 |
| f1181 | UiNormalizationCalls upholds UiNormalizationCallsProtocol. | 1636 |
| f1182 | UiNormalizationCallsProtocol has completeness = closed. | 1637 |
| f1183 | UiNormalizationCallsProtocol permits NormalizeUiRequest. | 1638 |
| f1184 | UiNormalizationCallsProtocol permits NormalizeUiResult. | 1639 |
| f1185 | UiSessionRecord has completeness = closed. | 1640 |
| f1186 | UiSessionRecord has-field SessionDraft. | 1641 |
| f1187 | UiSessionRecord has-field SessionGeneration. | 1642 |
| f1188 | UiSessionState upholds UiSessionRecord. | 1643 |
| f1189 | UiStateReconciler consumes SduiModelPort. | 1644 |
| f1190 | UiStateReconciler consumes UiStatePort. | 1645 |
| f1191 | UiStateReconciler owns MatchCompatibleWidgets. | 1646 |
| f1192 | UiStateReconciler owns PreserveCompatibleUiState. | 1647 |
| f1193 | UiStateReconciler owns ResetIncompatibleUiState. | 1648 |
| f1194 | UiStateReconciler provides InteractiveSession. | 1649 |
| f1195 | UiTokenizationCalls upholds UiTokenizationCallsProtocol. | 1650 |
| f1196 | UiTokenizationCallsProtocol has completeness = closed. | 1651 |
| f1197 | UiTokenizationCallsProtocol permits TokenizeUiRequest. | 1652 |
| f1198 | UiTokenizationCallsProtocol permits TokenizeUiResult. | 1653 |
| f1199 | UiValidationCalls upholds UiValidationCallsProtocol. | 1654 |
| f1200 | UiValidationCallsProtocol has completeness = closed. | 1655 |
| f1201 | UiValidationCallsProtocol permits ValidateUiRequest. | 1656 |
| f1202 | UiValidationCallsProtocol permits ValidateUiResult. | 1657 |
| f1203 | UnboundLocalAction exercises PrototypeUserInterface. | 1658 |
| f1204 | UnboundLocalAction has completeness = closed. | 1659 |
| f1205 | UnboundLocalAction illustrates G2M3FyneInteractions. | 1660 |
| f1206 | UnboundLocalAction runs-in UiPreview. | 1661 |
| f1207 | UnboundLocalAction step 1 sends DomainActionRequest from FyneHost to GoDomainImplementation via GoDomainCalls. | 1662 |
| f1208 | UnboundLocalAction step 2 sends DomainActionResult from GoDomainImplementation to FyneHost via GoDomainCalls reply-to 1. | 1663 |
| f1209 | ValidateActionInput realizes ExecutableDesign. | 1664 |
| f1210 | ValidatePropertyBatch realizes InteractiveSession. | 1665 |
| f1211 | ValidateRelativeFormatting has repeatability = deterministic. | 1666 |
| f1212 | ValidateRelativeFormatting has state-retention = stateless. | 1667 |
| f1213 | ValidateRelativeFormatting realizes SduiSourceModel. | 1668 |
| f1214 | ValidateSdlProfile has repeatability = deterministic. | 1669 |
| f1215 | ValidateSdlProfile has state-retention = stateless. | 1670 |
| f1216 | ValidateSdlProfile realizes SdlSourceModel. | 1671 |
| f1217 | ValidateSdlStructure allocated-to CommandLineHost in mode SourceInspection. | 1672 |
| f1218 | ValidateSdlStructure contributes-to StructuralModelInspection. | 1673 |
| f1219 | ValidateSdlStructure has repeatability = deterministic. | 1674 |
| f1220 | ValidateSdlStructure has state-retention = stateless. | 1675 |
| f1221 | ValidateSdlStructure realizes SdlSourceModel. | 1676 |
| f1222 | ValidateSymbolicBindings has repeatability = deterministic. | 1677 |
| f1223 | ValidateSymbolicBindings has state-retention = stateless. | 1678 |
| f1224 | ValidateSymbolicBindings realizes SduiSourceModel. | 1679 |
| f1225 | ValidateUiEvent realizes InteractiveSession. | 1680 |
| f1226 | ValidateUiRequest has message-kind = request. | 1681 |
| f1227 | ValidateUiRequest upholds AstArtifactContract. | 1682 |
| f1228 | ValidateUiResult has message-kind = result. | 1683 |
| f1229 | ValidateUiResult replies-to ValidateUiRequest. | 1684 |
| f1230 | ValidateUiResult upholds ValidationOutcomeContract. | 1685 |
| f1231 | ValidateViewRequest allocated-to ViewServiceHost in mode DocumentBrowsing. | 1686 |
| f1232 | ValidateViewRequest contributes-to NavigableDesignDocumentation. | 1687 |
| f1233 | ValidateVisualResources realizes RichContent. | 1688 |
| f1234 | ValidateWidgetArguments allocated-to CommandLineHost in mode SourceInspection. | 1689 |
| f1235 | ValidateWidgetArguments contributes-to StructuralModelInspection. | 1690 |
| f1236 | ValidateWidgetArguments has repeatability = deterministic. | 1691 |
| f1237 | ValidateWidgetArguments has state-retention = stateless. | 1692 |
| f1238 | ValidateWidgetArguments realizes SduiSourceModel. | 1693 |
| f1239 | ValidationOutcomeContract has completeness = closed. | 1694 |
| f1240 | ValidationOutcomeContract has-field ModelDiagnostics. | 1695 |
| f1241 | ValidationOutcomeContract has-field ModelIsValid. | 1696 |
| f1242 | VerifyNativeBehaviorParity allocated-to CommandLineHost in mode NativeBuild. | 1697 |
| f1243 | VerifyNativeBehaviorParity contributes-to NativeGoAssembly. | 1698 |
| f1244 | ViewArtifactStore owns EvictUnusedViewBundles. | 1699 |
| f1245 | ViewArtifactStore owns PublishViewBundle. | 1700 |
| f1246 | ViewArtifactStore owns RetainVisibleViewBundle. | 1701 |
| f1247 | ViewArtifactStore uses ViewPublicationCalls as receiver of PublishViewRequest in mode DocumentBrowsing. | 1702 |
| f1248 | ViewArtifactStore uses ViewPublicationCalls as sender of ViewReferenceResult in mode DocumentBrowsing. | 1703 |
| f1249 | ViewBundleBytes has presence = required. | 1704 |
| f1250 | ViewBundleBytes has value-type = bytes. | 1705 |
| f1251 | ViewBundleContract has completeness = closed. | 1706 |
| f1252 | ViewBundleContract has-field ViewBundleBytes. | 1707 |
| f1253 | ViewBundleContract has-field ViewBundleRevision. | 1708 |
| f1254 | ViewBundleResult has message-kind = result. | 1709 |
| f1255 | ViewBundleResult replies-to ProjectViewRequest. | 1710 |
| f1256 | ViewBundleResult upholds ViewBundleContract. | 1711 |
| f1257 | ViewBundleRevision has presence = required. | 1712 |
| f1258 | ViewBundleRevision has value-type = text. | 1713 |
| f1259 | ViewConsumerId has presence = required. | 1714 |
| f1260 | ViewConsumerId has value-type = text. | 1715 |
| f1261 | ViewDisplayCalls upholds ViewDisplayCallsProtocol. | 1716 |
| f1262 | ViewDisplayCallsProtocol has completeness = closed. | 1717 |
| f1263 | ViewDisplayCallsProtocol permits DisplayViewRequest. | 1718 |
| f1264 | ViewDisplayCallsProtocol permits ViewDisplayResult. | 1719 |
| f1265 | ViewDisplayCallsProtocol permits ViewTargetUnavailable. | 1720 |
| f1266 | ViewDisplayResult has message-kind = result. | 1721 |
| f1267 | ViewDisplayResult replies-to DisplayViewRequest. | 1722 |
| f1268 | ViewDisplayResult upholds ViewOpenedContract. | 1723 |
| f1269 | ViewEntryPath has presence = required. | 1724 |
| f1270 | ViewEntryPath has value-type = text. | 1725 |
| f1271 | ViewFailureCode has presence = required. | 1726 |
| f1272 | ViewFailureCode has value-type = text. | 1727 |
| f1273 | ViewFailureContract has completeness = closed. | 1728 |
| f1274 | ViewFailureContract has-field FailedViewRequestId. | 1729 |
| f1275 | ViewFailureContract has-field ViewFailureCode. | 1730 |
| f1276 | ViewFailureContract has-field ViewFailureDiagnostic. | 1731 |
| f1277 | ViewFailureDiagnostic has presence = required. | 1732 |
| f1278 | ViewFailureDiagnostic has value-type = text. | 1733 |
| f1279 | ViewLaunchCalls upholds ViewLaunchCallsProtocol. | 1734 |
| f1280 | ViewLaunchCallsProtocol has completeness = closed. | 1735 |
| f1281 | ViewLaunchCallsProtocol permits LaunchViewRequest. | 1736 |
| f1282 | ViewLaunchCallsProtocol permits ViewLaunchResult. | 1737 |
| f1283 | ViewLaunchResult has message-kind = result. | 1738 |
| f1284 | ViewLaunchResult replies-to LaunchViewRequest. | 1739 |
| f1285 | ViewLaunchResult upholds ViewOpenedContract. | 1740 |
| f1286 | ViewLeaseId has presence = required. | 1741 |
| f1287 | ViewLeaseId has value-type = text. | 1742 |
| f1288 | ViewManifestPath has presence = required. | 1743 |
| f1289 | ViewManifestPath has value-type = text. | 1744 |
| f1290 | ViewNavigationCalls upholds ViewNavigationCallsProtocol. | 1745 |
| f1291 | ViewNavigationCallsProtocol has completeness = closed. | 1746 |
| f1292 | ViewNavigationCallsProtocol permits SelectViewRequest. | 1747 |
| f1293 | ViewNavigationCallsProtocol permits ViewOpenedResult. | 1748 |
| f1294 | ViewNavigationCallsProtocol permits ViewRequestRejected. | 1749 |
| f1295 | ViewOpenContract has completeness = closed. | 1750 |
| f1296 | ViewOpenContract has-field OpenViewConsumerId. | 1751 |
| f1297 | ViewOpenContract has-field OpenViewEntryPath. | 1752 |
| f1298 | ViewOpenContract has-field OpenViewLeaseId. | 1753 |
| f1299 | ViewOpenContract has-field OpenViewPaneId. | 1754 |
| f1300 | ViewOpenContract has-field OpenViewRequestId. | 1755 |
| f1301 | ViewOpenContract has-field OpenViewRevision. | 1756 |
| f1302 | ViewOpenContract has-field OpenViewWindowId. | 1757 |
| f1303 | ViewOpenedContract has completeness = closed. | 1758 |
| f1304 | ViewOpenedContract has-field OpenedViewRequestId. | 1759 |
| f1305 | ViewOpenedContract has-field OpenedViewRevision. | 1760 |
| f1306 | ViewOpenedResult has message-kind = result. | 1761 |
| f1307 | ViewOpenedResult replies-to SelectViewRequest. | 1762 |
| f1308 | ViewOpenedResult upholds ViewOpenedContract. | 1763 |
| f1309 | ViewPaneId has presence = required. | 1764 |
| f1310 | ViewPaneId has value-type = text. | 1765 |
| f1311 | ViewProjectId has presence = required. | 1766 |
| f1312 | ViewProjectId has value-type = text. | 1767 |
| f1313 | ViewProjectionCalls upholds ViewProjectionCallsProtocol. | 1768 |
| f1314 | ViewProjectionCallsProtocol has completeness = closed. | 1769 |
| f1315 | ViewProjectionCallsProtocol permits ProjectViewRequest. | 1770 |
| f1316 | ViewProjectionCallsProtocol permits ViewBundleResult. | 1771 |
| f1317 | ViewProjectionCallsProtocol permits ViewProjectionRejected. | 1772 |
| f1318 | ViewProjectionFailed exercises BrowseDesignViews. | 1773 |
| f1319 | ViewProjectionFailed has completeness = closed. | 1774 |
| f1320 | ViewProjectionFailed illustrates G6M2OnDemandViews. | 1775 |
| f1321 | ViewProjectionFailed illustrates G6M4SessionPublication. | 1776 |
| f1322 | ViewProjectionFailed runs-in DocumentBrowsing. | 1777 |
| f1323 | ViewProjectionFailed step 1 sends SelectViewRequest from XfmdDocumentHost to DocumentBroker via ViewNavigationCalls. | 1778 |
| f1324 | ViewProjectionFailed step 2 sends ProjectViewRequest from DocumentBroker to SdlViewpointGenerator via ViewProjectionCalls. | 1779 |
| f1325 | ViewProjectionFailed step 3 sends ViewProjectionRejected from SdlViewpointGenerator to DocumentBroker via ViewProjectionCalls reply-to 2. | 1780 |
| f1326 | ViewProjectionFailed step 4 sends ViewRequestRejected from DocumentBroker to XfmdDocumentHost via ViewNavigationCalls reply-to 1. | 1781 |
| f1327 | ViewProjectionRejected has message-kind = result. | 1782 |
| f1328 | ViewProjectionRejected replies-to ProjectViewRequest. | 1783 |
| f1329 | ViewProjectionRejected upholds ViewFailureContract. | 1784 |
| f1330 | ViewPublicationCalls upholds ViewPublicationCallsProtocol. | 1785 |
| f1331 | ViewPublicationCallsProtocol has completeness = closed. | 1786 |
| f1332 | ViewPublicationCallsProtocol permits PublishViewRequest. | 1787 |
| f1333 | ViewPublicationCallsProtocol permits ViewReferenceResult. | 1788 |
| f1334 | ViewPublishedRevision has presence = required. | 1789 |
| f1335 | ViewPublishedRevision has value-type = text. | 1790 |
| f1336 | ViewReferenceContract has completeness = closed. | 1791 |
| f1337 | ViewReferenceContract has-field ViewEntryPath. | 1792 |
| f1338 | ViewReferenceContract has-field ViewLeaseId. | 1793 |
| f1339 | ViewReferenceContract has-field ViewManifestPath. | 1794 |
| f1340 | ViewReferenceContract has-field ViewPublishedRevision. | 1795 |
| f1341 | ViewReferenceResult has message-kind = result. | 1796 |
| f1342 | ViewReferenceResult replies-to PublishViewRequest. | 1797 |
| f1343 | ViewReferenceResult upholds ViewReferenceContract. | 1798 |
| f1344 | ViewRequestId has presence = required. | 1799 |
| f1345 | ViewRequestId has value-type = text. | 1800 |
| f1346 | ViewRequestRejected has message-kind = result. | 1801 |
| f1347 | ViewRequestRejected replies-to SelectViewRequest. | 1802 |
| f1348 | ViewRequestRejected upholds ViewFailureContract. | 1803 |
| f1349 | ViewSelectionContract has completeness = closed. | 1804 |
| f1350 | ViewSelectionContract has-field ViewConsumerId. | 1805 |
| f1351 | ViewSelectionContract has-field ViewPaneId. | 1806 |
| f1352 | ViewSelectionContract has-field ViewProjectId. | 1807 |
| f1353 | ViewSelectionContract has-field ViewRequestId. | 1808 |
| f1354 | ViewSelectionContract has-field ViewSelector. | 1809 |
| f1355 | ViewSelectionContract has-field ViewSourceRevision. | 1810 |
| f1356 | ViewSelectionContract has-field ViewWindowId. | 1811 |
| f1357 | ViewSelector has presence = required. | 1812 |
| f1358 | ViewSelector has value-type = text. | 1813 |
| f1359 | ViewServiceHost contains DocumentBroker. | 1814 |
| f1360 | ViewServiceHost contains ViewArtifactStore. | 1815 |
| f1361 | ViewServiceHost contains ViewerLaunchAdapter. | 1816 |
| f1362 | ViewSourceRevision has presence = required. | 1817 |
| f1363 | ViewSourceRevision has value-type = text. | 1818 |
| f1364 | ViewTargetUnavailable has message-kind = result. | 1819 |
| f1365 | ViewTargetUnavailable replies-to DisplayViewRequest. | 1820 |
| f1366 | ViewTargetUnavailable upholds ViewFailureContract. | 1821 |
| f1367 | ViewWindowId has presence = required. | 1822 |
| f1368 | ViewWindowId has value-type = text. | 1823 |
| f1369 | ViewerLaunchAdapter owns DispatchViewOpen. | 1824 |
| f1370 | ViewerLaunchAdapter owns ResolveConfiguredViewer. | 1825 |
| f1371 | ViewerLaunchAdapter uses ViewDisplayCalls as receiver of ViewDisplayResult in mode DocumentBrowsing. | 1826 |
| f1372 | ViewerLaunchAdapter uses ViewDisplayCalls as receiver of ViewTargetUnavailable in mode DocumentBrowsing. | 1827 |
| f1373 | ViewerLaunchAdapter uses ViewDisplayCalls as sender of DisplayViewRequest in mode DocumentBrowsing. | 1828 |
| f1374 | ViewerLaunchAdapter uses ViewLaunchCalls as receiver of LaunchViewRequest in mode DocumentBrowsing. | 1829 |
| f1375 | ViewerLaunchAdapter uses ViewLaunchCalls as sender of ViewLaunchResult in mode DocumentBrowsing. | 1830 |
| f1376 | WriteGeneratedArtifacts allocated-to CommandLineHost in mode StaticExport. | 1831 |
| f1377 | WriteGeneratedArtifacts contributes-to DesignDocumentation. | 1832 |
| f1378 | WriteGeneratedArtifacts realizes StaticDocumentation. | 1833 |
| f1379 | XfmdDocumentHost owns CaptureNavigationTarget. | 1834 |
| f1380 | XfmdDocumentHost owns RouteDocumentToPane. | 1835 |
| f1381 | XfmdDocumentHost uses ViewDisplayCalls as receiver of DisplayViewRequest in mode DocumentBrowsing. | 1836 |
| f1382 | XfmdDocumentHost uses ViewDisplayCalls as sender of ViewDisplayResult in mode DocumentBrowsing. | 1837 |
| f1383 | XfmdDocumentHost uses ViewDisplayCalls as sender of ViewTargetUnavailable in mode DocumentBrowsing. | 1838 |
| f1384 | XfmdDocumentHost uses ViewNavigationCalls as receiver of ViewOpenedResult in mode DocumentBrowsing. | 1839 |
| f1385 | XfmdDocumentHost uses ViewNavigationCalls as receiver of ViewRequestRejected in mode DocumentBrowsing. | 1840 |
| f1386 | XfmdDocumentHost uses ViewNavigationCalls as sender of SelectViewRequest in mode DocumentBrowsing. | 1841 |

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
| BrowseDesignViews | usecase | 21 |
| BuildGeneratedApplication | functionality | 22 |
| BuildGoRequest | message | 23 |
| BuildGoResult | message | 24 |
| BuildNativeProduct | usecase | 25 |
| BuildNativeRealization | activity | 26 |
| BuildPreparedFrame | functionality | 27 |
| BuildSdlAst | functionality | 28 |
| BuildSduiAst | functionality | 29 |
| BuildToolPort | interface | 30 |
| BuildUiAstRequest | message | 31 |
| BuildUiAstResult | message | 32 |
| CancelPendingActions | functionality | 33 |
| CaptureNavigationTarget | functionality | 34 |
| CheckDomainStateCompatibility | functionality | 35 |
| CheckExecutionCompleteness | functionality | 36 |
| CheckFunctionSignatures | functionality | 37 |
| CloseSdlInstance | functionality | 38 |
| CloseUiInstance | functionality | 39 |
| CoalesceSourceChanges | functionality | 40 |
| CommandLineHost | container | 41 |
| CompileSdlRequest | message | 42 |
| CompileSdlResult | message | 43 |
| CompileUiRejected | message | 44 |
| CompileUiRequest | message | 45 |
| CompileUiResult | message | 46 |
| ComposeHeadlessExport | functionality | 47 |
| ComposeInteractiveSession | functionality | 48 |
| ComposeMarkdownDocument | functionality | 49 |
| ComposeViewPackage | functionality | 50 |
| ComputeClipping | functionality | 51 |
| ConnectTypedWidgetHandles | functionality | 52 |
| ContentProviderPort | interface | 53 |
| ContentServices | unit | 54 |
| CoordinateSdlCompilation | functionality | 55 |
| CoordinateSduiCompilation | functionality | 56 |
| CorrelateActionResult | functionality | 57 |
| CorrelateUiResult | functionality | 58 |
| CreateSdlInstance | functionality | 59 |
| CreateUiInstance | functionality | 60 |
| DesignAuthor | actor | 61 |
| DesignDocumentation | feature | 62 |
| DesignReviewer | actor | 63 |
| DesignSourceArchive | database | 64 |
| DesignSourceDocuments | dataset | 65 |
| DesignSourceRecord | contract | 66 |
| DevelopmentReload | capability | 67 |
| DevelopmentTools | unit | 68 |
| DiagnosticPort | interface | 69 |
| DiagnosticReporter | unit | 70 |
| DiagramEnginePort | interface | 71 |
| DiagramPort | interface | 72 |
| DiagramProvider | unit | 73 |
| DisconnectBindings | functionality | 74 |
| DispatchUiEvent | functionality | 75 |
| DispatchViewOpen | functionality | 76 |
| DisplayViewRequest | message | 77 |
| DocumentBroker | unit | 78 |
| DocumentBrowsing | mode | 79 |
| DomainActionRequest | message | 80 |
| DomainActionResult | message | 81 |
| DomainBindingPort | interface | 82 |
| DomainFunctionPort | interface | 83 |
| DomainOperations | capability | 84 |
| DomainStateMigrator | unit | 85 |
| DomainStatePort | interface | 86 |
| EditRunningPrototype | usecase | 87 |
| EvictUnusedViewBundles | functionality | 88 |
| ExecutableDesign | capability | 89 |
| ExecutionProfilePort | interface | 90 |
| ExpandUiDefinitions | functionality | 91 |
| ExportConsoleSnapshot | functionality | 92 |
| ExportSinkPort | interface | 93 |
| ExportSvgRequest | message | 94 |
| ExportSvgResult | message | 95 |
| ExportSvgSnapshot | functionality | 96 |
| ExportUiDocumentation | activity | 97 |
| ExportViewpointMarkdown | functionality | 98 |
| FailedViewRequestId | field | 99 |
| FileChangePort | interface | 100 |
| FramePresentationCalls | channel | 101 |
| FramePresentationCallsProtocol | contract | 102 |
| FyneBackend | unit | 103 |
| FyneHost | container | 104 |
| G1FrontendPort | activity | 105 |
| G1M1ParserAndAst | activity | 106 |
| G1M2ValidationAndNormalization | activity | 107 |
| G1M3Concept1AndDumps | activity | 108 |
| G2LayoutAndPresentation | activity | 109 |
| G2M1RelativeMeasurement | activity | 110 |
| G2M2SharedSvgGeometry | activity | 111 |
| G2M3FyneInteractions | activity | 112 |
| G2M4RichContent | activity | 113 |
| G3M1TypedUiSession | activity | 114 |
| G3M2CandidatePublication | activity | 115 |
| G3M3CompatibleState | activity | 116 |
| G3UiRuntimeAndReload | activity | 117 |
| G4M1SdlFrontend | activity | 118 |
| G4M2TypedExecution | activity | 119 |
| G4M3UiDomainBinding | activity | 120 |
| G4M4DomainReload | activity | 121 |
| G4SdlRuntimeAndBinding | activity | 122 |
| G5M1GeneratedGo | activity | 123 |
| G5M2BehaviorParity | activity | 124 |
| G5M3DocumentationExport | activity | 125 |
| G5M4RetirePython | activity | 126 |
| G5NativeGeneration | activity | 127 |
| G6M1StaticNavigation | activity | 128 |
| G6M2OnDemandViews | activity | 129 |
| G6M3XfmdNavigation | activity | 130 |
| G6M4SessionPublication | activity | 131 |
| G6NavigableDocumentation | activity | 132 |
| GenerateBindingRegistration | functionality | 133 |
| GenerateGoRequest | message | 134 |
| GenerateGoResult | message | 135 |
| GenerateModelConstructors | functionality | 136 |
| GenerateViewNavigation | functionality | 137 |
| GeneratedArtifactPort | interface | 138 |
| GeneratedGoContract | contract | 139 |
| GeneratedGoSources | field | 140 |
| GoBuildCalls | channel | 141 |
| GoBuildCallsProtocol | contract | 142 |
| GoBuildRunner | unit | 143 |
| GoCodeGenerator | unit | 144 |
| GoDomainCalls | channel | 145 |
| GoDomainCallsProtocol | contract | 146 |
| GoDomainImplementation | unit | 147 |
| GoGenerationCalls | channel | 148 |
| GoGenerationCallsProtocol | contract | 149 |
| HandleFocusAndTextInput | functionality | 150 |
| IdentifySourceRevision | functionality | 151 |
| InspectDesignSource | activity | 152 |
| InspectModels | usecase | 153 |
| InspectSdlSource | activity | 154 |
| InspectSduiSource | activity | 155 |
| InteractiveFramePrepared | scenario | 156 |
| InteractiveSession | capability | 157 |
| InteractiveUiPreview | feature | 158 |
| InvalidViewSelectionRejected | scenario | 159 |
| InvokeRegisteredFunction | functionality | 160 |
| KeepLastValidModels | functionality | 161 |
| KeyViewRevision | functionality | 162 |
| LaunchViewRequest | message | 163 |
| LayoutArguments | contract | 164 |
| LayoutCalls | channel | 165 |
| LayoutCallsProtocol | contract | 166 |
| LayoutModelArtifact | field | 167 |
| LayoutRequest | message | 168 |
| LayoutResult | message | 169 |
| LayoutViewportHeight | field | 170 |
| LayoutViewportWidth | field | 171 |
| LiveEditing | mode | 172 |
| LiveModelReload | feature | 173 |
| ManageDomainState | functionality | 174 |
| ManageWidgetIdentities | functionality | 175 |
| MarkdownProvider | unit | 176 |
| MatchCompatibleWidgets | functionality | 177 |
| MeasureMarkdownContent | functionality | 178 |
| MeasureUiContent | functionality | 179 |
| MeasuredPresentation | capability | 180 |
| MeasurementPort | interface | 181 |
| MigrateOrResetDomainState | functionality | 182 |
| ModelDiagnostics | field | 183 |
| ModelIsValid | field | 184 |
| ModelReloadCalls | channel | 185 |
| ModelReloadCallsProtocol | contract | 186 |
| NativeBuild | mode | 187 |
| NativeBuildContract | contract | 188 |
| NativeBuildDiagnostics | field | 189 |
| NativeBuildSucceeded | field | 190 |
| NativeGoAssembly | feature | 191 |
| NativeInteraction | capability | 192 |
| NativeProgramBuilt | scenario | 193 |
| NativeRealization | capability | 194 |
| NativeUiActions | channel | 195 |
| NativeUiActionsProtocol | contract | 196 |
| NavigableDesignDocumentation | feature | 197 |
| NormalizeSdlModel | functionality | 198 |
| NormalizeUiRequest | message | 199 |
| NormalizeUiResult | message | 200 |
| NormalizedModelArtifact | field | 201 |
| NormalizedModelContract | contract | 202 |
| NoticeGeneration | field | 203 |
| NoticeVersion | field | 204 |
| ObserveSourceChanges | functionality | 205 |
| OpenViewConsumerId | field | 206 |
| OpenViewEntryPath | field | 207 |
| OpenViewLeaseId | field | 208 |
| OpenViewPaneId | field | 209 |
| OpenViewRequestId | field | 210 |
| OpenViewRevision | field | 211 |
| OpenViewWindowId | field | 212 |
| OpenedViewRequestId | field | 213 |
| OpenedViewRevision | field | 214 |
| PerformDomainOperation | functionality | 215 |
| PrepareCandidateModels | functionality | 216 |
| PrepareDiagramResource | functionality | 217 |
| PrepareMarkdown | functionality | 218 |
| PreparedFrameArtifact | field | 219 |
| PreparedFrameContract | contract | 220 |
| PreparedFramePort | interface | 221 |
| PresentFrameRequest | message | 222 |
| PresentFrameResult | message | 223 |
| PresentationOutcome | contract | 224 |
| PresentationReady | field | 225 |
| PreserveCompatibleUiState | functionality | 226 |
| PreserveHandwrittenSources | functionality | 227 |
| PreserveSdlSourceMap | functionality | 228 |
| PreserveUiRegions | functionality | 229 |
| PreserveUiSourceMap | functionality | 230 |
| PreserveViewAnchors | functionality | 231 |
| ProjectSdlViewpoints | functionality | 232 |
| ProjectSelectedView | functionality | 233 |
| ProjectUiGeneration | functionality | 234 |
| ProjectViewRequest | message | 235 |
| PrototypeUserInterface | usecase | 236 |
| PublishDesignDocumentation | usecase | 237 |
| PublishDomainUpdates | functionality | 238 |
| PublishModelGeneration | functionality | 239 |
| PublishPresentation | functionality | 240 |
| PublishViewBundle | functionality | 241 |
| PublishViewRequest | message | 242 |
| ReadBoundedSources | functionality | 243 |
| RealizeDesign | activity | 244 |
| ReconcileWidgets | functionality | 245 |
| RegisterDomainFunctions | functionality | 246 |
| RejectStaleUiEvent | functionality | 247 |
| RejectStaleViewResults | functionality | 248 |
| ReleaseNativeWidgets | functionality | 249 |
| ReleaseVisualResources | functionality | 250 |
| ReloadArguments | contract | 251 |
| ReloadBoundModels | activity | 252 |
| ReloadCoordinator | unit | 253 |
| ReloadDesignSession | activity | 254 |
| ReloadDiagnostic | field | 255 |
| ReloadOutcome | contract | 256 |
| ReloadPort | interface | 257 |
| ReloadPublished | message | 258 |
| ReloadPublishedGeneration | field | 259 |
| ReloadRejected | message | 260 |
| ReloadRequest | message | 261 |
| ReloadSourceRevision | field | 262 |
| ReloadSourceText | field | 263 |
| ReloadUiModel | activity | 264 |
| ReportBindingDiagnostics | functionality | 265 |
| ReportSourceDiagnostics | functionality | 266 |
| ResetIncompatibleUiState | functionality | 267 |
| ResolveAncestorDimensions | functionality | 268 |
| ResolveCallbackSymbols | functionality | 269 |
| ResolveConfiguredViewer | functionality | 270 |
| ResolveSdlSymbols | functionality | 271 |
| ResolveUiNames | functionality | 272 |
| ResourcePort | interface | 273 |
| ResourceStore | unit | 274 |
| RestartChangedGoProgram | functionality | 275 |
| RetainVisibleViewBundle | functionality | 276 |
| RetirePreviousGeneration | functionality | 277 |
| RetireReplacedPythonEntryPoints | functionality | 278 |
| RevokeWidgetGenerations | functionality | 279 |
| RichContent | capability | 280 |
| RichDocument | mode | 281 |
| RouteDocumentToPane | functionality | 282 |
| RouteDomainBindings | functionality | 283 |
| RunBoundUiAction | activity | 284 |
| RunDesignSession | activity | 285 |
| RunUnboundUiPreview | activity | 286 |
| ScheduleUiPublication | functionality | 287 |
| SdlActionCalls | channel | 288 |
| SdlActionCallsProtocol | contract | 289 |
| SdlActionRequest | message | 290 |
| SdlActionResult | message | 291 |
| SdlCompilationCalls | channel | 292 |
| SdlCompilationCallsProtocol | contract | 293 |
| SdlDispatcher | unit | 294 |
| SdlExecutionGate | unit | 295 |
| SdlExecutionPort | interface | 296 |
| SdlFrontend | unit | 297 |
| SdlFrontendPort | interface | 298 |
| SdlFunctionRegistry | unit | 299 |
| SdlLexer | unit | 300 |
| SdlLibrary | unit | 301 |
| SdlModelPort | interface | 302 |
| SdlModelReloadAccepted | scenario | 303 |
| SdlNormalizer | unit | 304 |
| SdlParser | unit | 305 |
| SdlReloadPort | interface | 306 |
| SdlRuntime | unit | 307 |
| SdlSourceModel | capability | 308 |
| SdlStateStore | unit | 309 |
| SdlUiBindingAdapter | unit | 310 |
| SdlValidator | unit | 311 |
| SdlViewpointGenerator | unit | 312 |
| SduiDispatcher | unit | 313 |
| SduiFrontend | unit | 314 |
| SduiFrontendPort | interface | 315 |
| SduiInstanceStore | unit | 316 |
| SduiLayout | unit | 317 |
| SduiLexer | unit | 318 |
| SduiLibrary | unit | 319 |
| SduiModelPort | interface | 320 |
| SduiNormalizer | unit | 321 |
| SduiParser | unit | 322 |
| SduiPresentation | unit | 323 |
| SduiPropertyStore | unit | 324 |
| SduiRuntime | unit | 325 |
| SduiSourceModel | capability | 326 |
| SduiValidator | unit | 327 |
| SelectViewRequest | message | 328 |
| SelectedViewOpened | scenario | 329 |
| ServeViewRequests | functionality | 330 |
| SessionDraft | field | 331 |
| SessionGeneration | field | 332 |
| SnapshotDomainState | functionality | 333 |
| SnapshotUiState | functionality | 334 |
| SourceDiagnostics | capability | 335 |
| SourceDocumentRevision | field | 336 |
| SourceDocumentText | field | 337 |
| SourceInputPort | interface | 338 |
| SourceInspection | mode | 339 |
| SourceLoader | unit | 340 |
| SourceLoading | capability | 341 |
| SourceSnapshotPort | interface | 342 |
| SourceWatcher | unit | 343 |
| StaticDocumentation | capability | 344 |
| StaticExport | mode | 345 |
| StaticFrameExported | scenario | 346 |
| StructuralModelInspection | feature | 347 |
| SvgDocumentContract | contract | 348 |
| SvgDocumentText | field | 349 |
| SvgExportCalls | channel | 350 |
| SvgExportCallsProtocol | contract | 351 |
| TokenArtifact | field | 352 |
| TokenArtifactContract | contract | 353 |
| TokenizeSdlSource | functionality | 354 |
| TokenizeSduiSource | functionality | 355 |
| TokenizeUiRequest | message | 356 |
| TokenizeUiResult | message | 357 |
| TraceViewpointFacts | functionality | 358 |
| TrackInputDraft | functionality | 359 |
| TryDomainInteraction | usecase | 360 |
| TypedDomainBinding | feature | 361 |
| UiActionRejected | message | 362 |
| UiActionRequest | message | 363 |
| UiActionResult | message | 364 |
| UiAstCalls | channel | 365 |
| UiAstCallsProtocol | contract | 366 |
| UiCompilationAccepted | scenario | 367 |
| UiCompilationCalls | channel | 368 |
| UiCompilationCallsProtocol | contract | 369 |
| UiDomainActions | channel | 370 |
| UiDomainActionsProtocol | contract | 371 |
| UiGenerationChanged | variant | 372 |
| UiGenerationContract | contract | 373 |
| UiGenerationEvents | channel | 374 |
| UiGenerationEventsProtocol | contract | 375 |
| UiGenerationNotices | datagram | 376 |
| UiGenerationWire | encoding | 377 |
| UiModelReloadAccepted | scenario | 378 |
| UiModelReloadRejected | scenario | 379 |
| UiNormalizationCalls | channel | 380 |
| UiNormalizationCallsProtocol | contract | 381 |
| UiPreview | mode | 382 |
| UiReloadPort | interface | 383 |
| UiSessionPort | interface | 384 |
| UiSessionRecord | contract | 385 |
| UiSessionState | dataset | 386 |
| UiSnapshotPort | interface | 387 |
| UiStatePort | interface | 388 |
| UiStateReconciler | unit | 389 |
| UiTokenizationCalls | channel | 390 |
| UiTokenizationCallsProtocol | contract | 391 |
| UiValidationCalls | channel | 392 |
| UiValidationCallsProtocol | contract | 393 |
| UnboundLocalAction | scenario | 394 |
| ValidateActionInput | functionality | 395 |
| ValidatePropertyBatch | functionality | 396 |
| ValidateRelativeFormatting | functionality | 397 |
| ValidateSdlProfile | functionality | 398 |
| ValidateSdlStructure | functionality | 399 |
| ValidateSymbolicBindings | functionality | 400 |
| ValidateUiEvent | functionality | 401 |
| ValidateUiRequest | message | 402 |
| ValidateUiResult | message | 403 |
| ValidateViewRequest | functionality | 404 |
| ValidateVisualResources | functionality | 405 |
| ValidateWidgetArguments | functionality | 406 |
| ValidationOutcomeContract | contract | 407 |
| VerifyNativeBehaviorParity | functionality | 408 |
| ViewArtifactStore | unit | 409 |
| ViewBundleBytes | field | 410 |
| ViewBundleContract | contract | 411 |
| ViewBundleResult | message | 412 |
| ViewBundleRevision | field | 413 |
| ViewConsumerId | field | 414 |
| ViewDisplayCalls | channel | 415 |
| ViewDisplayCallsProtocol | contract | 416 |
| ViewDisplayResult | message | 417 |
| ViewEntryPath | field | 418 |
| ViewFailureCode | field | 419 |
| ViewFailureContract | contract | 420 |
| ViewFailureDiagnostic | field | 421 |
| ViewLaunchCalls | channel | 422 |
| ViewLaunchCallsProtocol | contract | 423 |
| ViewLaunchResult | message | 424 |
| ViewLeaseId | field | 425 |
| ViewManifestPath | field | 426 |
| ViewNavigationCalls | channel | 427 |
| ViewNavigationCallsProtocol | contract | 428 |
| ViewOpenContract | contract | 429 |
| ViewOpenedContract | contract | 430 |
| ViewOpenedResult | message | 431 |
| ViewPaneId | field | 432 |
| ViewProjectId | field | 433 |
| ViewProjectionCalls | channel | 434 |
| ViewProjectionCallsProtocol | contract | 435 |
| ViewProjectionFailed | scenario | 436 |
| ViewProjectionRejected | message | 437 |
| ViewPublicationCalls | channel | 438 |
| ViewPublicationCallsProtocol | contract | 439 |
| ViewPublishedRevision | field | 440 |
| ViewReferenceContract | contract | 441 |
| ViewReferenceResult | message | 442 |
| ViewRequestId | field | 443 |
| ViewRequestRejected | message | 444 |
| ViewSelectionContract | contract | 445 |
| ViewSelector | field | 446 |
| ViewServiceHost | container | 447 |
| ViewSourceRevision | field | 448 |
| ViewTargetUnavailable | message | 449 |
| ViewWindowId | field | 450 |
| ViewerLaunchAdapter | unit | 451 |
| WidgetBackendPort | interface | 452 |
| WriteGeneratedArtifacts | functionality | 453 |
| XfmdDocumentHost | container | 454 |
