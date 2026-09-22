# Generert implementasjonsplan

Status er modellens påstand, ikke bevis for kjøring.

[Navigator](navigator.md)

## BuildNativeRealization

Status: unspecified.

| Faktum | Kilde-ID |
| --- | --- |
| BuildNativeRealization refines RealizeDesign. | f0062 |

## ExportUiDocumentation

Status: unspecified.

| Faktum | Kilde-ID |
| --- | --- |
| ExportUiDocumentation refines InspectDesignSource. | f0234 |

## G1FrontendPort

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G1FrontendPort delivers StructuralModelInspection. | f0276 |
| G1FrontendPort has implementation-status = verified. | f0277 |
| G1M1ParserAndAst refines G1FrontendPort. | f0283 |
| G1M2ValidationAndNormalization refines G1FrontendPort. | f0294 |
| G1M3Concept1AndDumps refines G1FrontendPort. | f0299 |

## G1M1ParserAndAst

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G1M1ParserAndAst addresses BuildSduiAst. | f0278 |
| G1M1ParserAndAst addresses IdentifySourceRevision. | f0279 |
| G1M1ParserAndAst addresses ReadBoundedSources. | f0280 |
| G1M1ParserAndAst addresses TokenizeSduiSource. | f0281 |
| G1M1ParserAndAst has implementation-status = verified. | f0282 |
| G1M1ParserAndAst refines G1FrontendPort. | f0283 |
| G1M2ValidationAndNormalization depends-on G1M1ParserAndAst. | f0292 |
| UiCompilationAccepted illustrates G1M1ParserAndAst. | f1160 |
| SduiParser owns BuildSduiAst. | f1017 |
| SourceLoader owns IdentifySourceRevision. | f1085 |
| SourceLoader owns ReadBoundedSources. | f1086 |
| SduiLexer owns TokenizeSduiSource. | f0999 |

## G1M2ValidationAndNormalization

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G1M2ValidationAndNormalization addresses CoordinateSduiCompilation. | f0284 |
| G1M2ValidationAndNormalization addresses ExpandUiDefinitions. | f0285 |
| G1M2ValidationAndNormalization addresses PreserveUiRegions. | f0286 |
| G1M2ValidationAndNormalization addresses PreserveUiSourceMap. | f0287 |
| G1M2ValidationAndNormalization addresses ResolveUiNames. | f0288 |
| G1M2ValidationAndNormalization addresses ValidateRelativeFormatting. | f0289 |
| G1M2ValidationAndNormalization addresses ValidateSymbolicBindings. | f0290 |
| G1M2ValidationAndNormalization addresses ValidateWidgetArguments. | f0291 |
| G1M2ValidationAndNormalization depends-on G1M1ParserAndAst. | f0292 |
| G1M2ValidationAndNormalization has implementation-status = verified. | f0293 |
| G1M2ValidationAndNormalization refines G1FrontendPort. | f0294 |
| G1M3Concept1AndDumps depends-on G1M2ValidationAndNormalization. | f0297 |
| G2M1RelativeMeasurement depends-on G1M2ValidationAndNormalization. | f0306 |
| G3M1TypedUiSession depends-on G1M2ValidationAndNormalization. | f0344 |
| UiCompilationAccepted illustrates G1M2ValidationAndNormalization. | f1161 |
| SduiFrontend owns CoordinateSduiCompilation. | f0962 |
| SduiNormalizer owns ExpandUiDefinitions. | f1011 |
| SduiNormalizer owns PreserveUiRegions. | f1012 |
| SduiNormalizer owns PreserveUiSourceMap. | f1013 |
| SduiValidator owns ResolveUiNames. | f1043 |
| SduiValidator owns ValidateRelativeFormatting. | f1044 |
| SduiValidator owns ValidateSymbolicBindings. | f1045 |
| SduiValidator owns ValidateWidgetArguments. | f1046 |

## G1M3Concept1AndDumps

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G1M3Concept1AndDumps addresses ExportConsoleSnapshot. | f0295 |
| G1M3Concept1AndDumps addresses ReportSourceDiagnostics. | f0296 |
| G1M3Concept1AndDumps depends-on G1M2ValidationAndNormalization. | f0297 |
| G1M3Concept1AndDumps has implementation-status = verified. | f0298 |
| G1M3Concept1AndDumps refines G1FrontendPort. | f0299 |
| G5M4RetirePython depends-on G1M3Concept1AndDumps. | f0429 |
| SduiPresentation owns ExportConsoleSnapshot. | f1023 |
| DiagnosticReporter owns ReportSourceDiagnostics. | f0179 |

## G2LayoutAndPresentation

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G2LayoutAndPresentation delivers InteractiveUiPreview. | f0300 |
| G2LayoutAndPresentation has implementation-status = verified. | f0301 |
| G2M1RelativeMeasurement refines G2LayoutAndPresentation. | f0308 |
| G2M2SharedSvgGeometry refines G2LayoutAndPresentation. | f0313 |
| G2M3FyneInteractions refines G2LayoutAndPresentation. | f0322 |
| G2M4RichContent refines G2LayoutAndPresentation. | f0330 |

## G2M1RelativeMeasurement

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G2M1RelativeMeasurement addresses AllocateGeometry. | f0302 |
| G2M1RelativeMeasurement addresses ComputeClipping. | f0303 |
| G2M1RelativeMeasurement addresses MeasureUiContent. | f0304 |
| G2M1RelativeMeasurement addresses ResolveAncestorDimensions. | f0305 |
| G2M1RelativeMeasurement depends-on G1M2ValidationAndNormalization. | f0306 |
| G2M1RelativeMeasurement has implementation-status = verified. | f0307 |
| G2M1RelativeMeasurement refines G2LayoutAndPresentation. | f0308 |
| G2M2SharedSvgGeometry depends-on G2M1RelativeMeasurement. | f0311 |
| InteractiveFramePrepared illustrates G2M1RelativeMeasurement. | f0543 |
| SduiLayout owns AllocateGeometry. | f0989 |
| SduiLayout owns ComputeClipping. | f0991 |
| SduiLayout owns MeasureUiContent. | f0992 |
| SduiLayout owns ResolveAncestorDimensions. | f0993 |

## G2M2SharedSvgGeometry

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G2M2SharedSvgGeometry addresses BuildPreparedFrame. | f0309 |
| G2M2SharedSvgGeometry addresses ExportSvgSnapshot. | f0310 |
| G2M2SharedSvgGeometry depends-on G2M1RelativeMeasurement. | f0311 |
| G2M2SharedSvgGeometry has implementation-status = verified. | f0312 |
| G2M2SharedSvgGeometry refines G2LayoutAndPresentation. | f0313 |
| G2M3FyneInteractions depends-on G2M2SharedSvgGeometry. | f0320 |
| StaticFrameExported illustrates G2M2SharedSvgGeometry. | f1102 |
| SduiLayout owns BuildPreparedFrame. | f0990 |
| SduiPresentation owns ExportSvgSnapshot. | f1024 |

## G2M3FyneInteractions

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G2M3FyneInteractions addresses ComposeInteractiveSession. | f0314 |
| G2M3FyneInteractions addresses HandleFocusAndTextInput. | f0315 |
| G2M3FyneInteractions addresses PublishPresentation. | f0316 |
| G2M3FyneInteractions addresses ReconcileWidgets. | f0317 |
| G2M3FyneInteractions addresses ReleaseNativeWidgets. | f0318 |
| G2M3FyneInteractions addresses ScheduleUiPublication. | f0319 |
| G2M3FyneInteractions depends-on G2M2SharedSvgGeometry. | f0320 |
| G2M3FyneInteractions has implementation-status = verified. | f0321 |
| G2M3FyneInteractions refines G2LayoutAndPresentation. | f0322 |
| G2M4RichContent depends-on G2M3FyneInteractions. | f0328 |
| G3M1TypedUiSession depends-on G2M3FyneInteractions. | f0345 |
| InteractiveFramePrepared illustrates G2M3FyneInteractions. | f0544 |
| UnboundLocalAction illustrates G2M3FyneInteractions. | f1238 |
| FyneHost owns ComposeInteractiveSession. | f0267 |
| FyneBackend owns HandleFocusAndTextInput. | f0248 |
| FyneBackend owns PublishPresentation. | f0249 |
| FyneBackend owns ReconcileWidgets. | f0250 |
| FyneBackend owns ReleaseNativeWidgets. | f0251 |
| FyneHost owns ScheduleUiPublication. | f0268 |

## G2M4RichContent

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G2M4RichContent addresses MeasureMarkdownContent. | f0323 |
| G2M4RichContent addresses PrepareDiagramResource. | f0324 |
| G2M4RichContent addresses PrepareMarkdown. | f0325 |
| G2M4RichContent addresses ReleaseVisualResources. | f0326 |
| G2M4RichContent addresses ValidateVisualResources. | f0327 |
| G2M4RichContent depends-on G2M3FyneInteractions. | f0328 |
| G2M4RichContent has implementation-status = verified. | f0329 |
| G2M4RichContent refines G2LayoutAndPresentation. | f0330 |
| G5M3DocumentationExport depends-on G2M4RichContent. | f0423 |
| MarkdownProvider owns MeasureMarkdownContent. | f0595 |
| DiagramProvider owns PrepareDiagramResource. | f0183 |
| MarkdownProvider owns PrepareMarkdown. | f0596 |
| ResourceStore owns ReleaseVisualResources. | f0815 |
| ResourceStore owns ValidateVisualResources. | f0816 |

## G3M1TypedUiSession

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G3M1TypedUiSession addresses ApplyPropertyBatch. | f0331 |
| G3M1TypedUiSession addresses CloseUiInstance. | f0332 |
| G3M1TypedUiSession addresses CorrelateUiResult. | f0333 |
| G3M1TypedUiSession addresses CreateUiInstance. | f0334 |
| G3M1TypedUiSession addresses DispatchUiEvent. | f0335 |
| G3M1TypedUiSession addresses ManageWidgetIdentities. | f0336 |
| G3M1TypedUiSession addresses ProjectUiGeneration. | f0337 |
| G3M1TypedUiSession addresses RejectStaleUiEvent. | f0338 |
| G3M1TypedUiSession addresses RevokeWidgetGenerations. | f0339 |
| G3M1TypedUiSession addresses SnapshotUiState. | f0340 |
| G3M1TypedUiSession addresses TrackInputDraft. | f0341 |
| G3M1TypedUiSession addresses ValidatePropertyBatch. | f0342 |
| G3M1TypedUiSession addresses ValidateUiEvent. | f0343 |
| G3M1TypedUiSession depends-on G1M2ValidationAndNormalization. | f0344 |
| G3M1TypedUiSession depends-on G2M3FyneInteractions. | f0345 |
| G3M1TypedUiSession has implementation-status = verified. | f0346 |
| G3M1TypedUiSession refines G3UiRuntimeAndReload. | f0347 |
| G3M2CandidatePublication depends-on G3M1TypedUiSession. | f0354 |
| G4M3UiDomainBinding depends-on G3M1TypedUiSession. | f0396 |
| SduiPropertyStore owns ApplyPropertyBatch. | f1029 |
| SduiRuntime owns CloseUiInstance. | f1039 |
| SduiDispatcher owns CorrelateUiResult. | f0946 |
| SduiRuntime owns CreateUiInstance. | f1040 |
| SduiDispatcher owns DispatchUiEvent. | f0947 |
| SduiInstanceStore owns ManageWidgetIdentities. | f0979 |
| SduiInstanceStore owns ProjectUiGeneration. | f0980 |
| SduiDispatcher owns RejectStaleUiEvent. | f0948 |
| SduiInstanceStore owns RevokeWidgetGenerations. | f0981 |
| SduiInstanceStore owns SnapshotUiState. | f0982 |
| SduiPropertyStore owns TrackInputDraft. | f1030 |
| SduiPropertyStore owns ValidatePropertyBatch. | f1031 |
| SduiDispatcher owns ValidateUiEvent. | f0949 |

## G3M2CandidatePublication

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G3M2CandidatePublication addresses CoalesceSourceChanges. | f0348 |
| G3M2CandidatePublication addresses KeepLastValidModels. | f0349 |
| G3M2CandidatePublication addresses ObserveSourceChanges. | f0350 |
| G3M2CandidatePublication addresses PrepareCandidateModels. | f0351 |
| G3M2CandidatePublication addresses PublishModelGeneration. | f0352 |
| G3M2CandidatePublication addresses RetirePreviousGeneration. | f0353 |
| G3M2CandidatePublication depends-on G3M1TypedUiSession. | f0354 |
| G3M2CandidatePublication has implementation-status = verified. | f0355 |
| G3M2CandidatePublication refines G3UiRuntimeAndReload. | f0356 |
| G3M3CompatibleState depends-on G3M2CandidatePublication. | f0360 |
| UiModelReloadAccepted illustrates G3M2CandidatePublication. | f1198 |
| UiModelReloadRejected illustrates G3M2CandidatePublication. | f1208 |
| SourceWatcher owns CoalesceSourceChanges. | f1089 |
| ReloadCoordinator owns KeepLastValidModels. | f0760 |
| SourceWatcher owns ObserveSourceChanges. | f1090 |
| ReloadCoordinator owns PrepareCandidateModels. | f0761 |
| ReloadCoordinator owns PublishModelGeneration. | f0762 |
| ReloadCoordinator owns RetirePreviousGeneration. | f0763 |

## G3M3CompatibleState

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G3M3CompatibleState addresses MatchCompatibleWidgets. | f0357 |
| G3M3CompatibleState addresses PreserveCompatibleUiState. | f0358 |
| G3M3CompatibleState addresses ResetIncompatibleUiState. | f0359 |
| G3M3CompatibleState depends-on G3M2CandidatePublication. | f0360 |
| G3M3CompatibleState has implementation-status = verified. | f0361 |
| G3M3CompatibleState refines G3UiRuntimeAndReload. | f0362 |
| G4M4DomainReload depends-on G3M3CompatibleState. | f0403 |
| UiModelReloadAccepted illustrates G3M3CompatibleState. | f1199 |
| UiStateReconciler owns MatchCompatibleWidgets. | f1224 |
| UiStateReconciler owns PreserveCompatibleUiState. | f1225 |
| UiStateReconciler owns ResetIncompatibleUiState. | f1226 |

## G3UiRuntimeAndReload

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G3M1TypedUiSession refines G3UiRuntimeAndReload. | f0347 |
| G3M2CandidatePublication refines G3UiRuntimeAndReload. | f0356 |
| G3M3CompatibleState refines G3UiRuntimeAndReload. | f0362 |
| G3UiRuntimeAndReload delivers LiveModelReload. | f0363 |
| G3UiRuntimeAndReload has implementation-status = verified. | f0364 |

## G4M1SdlFrontend

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G4M1SdlFrontend addresses BuildSdlAst. | f0365 |
| G4M1SdlFrontend addresses CoordinateSdlCompilation. | f0366 |
| G4M1SdlFrontend addresses NormalizeSdlModel. | f0367 |
| G4M1SdlFrontend addresses PreserveSdlSourceMap. | f0368 |
| G4M1SdlFrontend addresses ResolveSdlSymbols. | f0369 |
| G4M1SdlFrontend addresses TokenizeSdlSource. | f0370 |
| G4M1SdlFrontend addresses ValidateSdlProfile. | f0371 |
| G4M1SdlFrontend addresses ValidateSdlStructure. | f0372 |
| G4M1SdlFrontend has implementation-status = verified. | f0373 |
| G4M1SdlFrontend refines G4SdlRuntimeAndBinding. | f0374 |
| G4M2TypedExecution depends-on G4M1SdlFrontend. | f0387 |
| G6M1StaticNavigation depends-on G4M1SdlFrontend. | f0444 |
| SdlModelReloadAccepted illustrates G4M1SdlFrontend. | f0884 |
| SdlParser owns BuildSdlAst. | f0894 |
| SdlFrontend owns CoordinateSdlCompilation. | f0867 |
| SdlNormalizer owns NormalizeSdlModel. | f0891 |
| SdlNormalizer owns PreserveSdlSourceMap. | f0892 |
| SdlValidator owns ResolveSdlSymbols. | f0923 |
| SdlLexer owns TokenizeSdlSource. | f0875 |
| SdlValidator owns ValidateSdlProfile. | f0925 |
| SdlValidator owns ValidateSdlStructure. | f0926 |

## G4M2TypedExecution

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G4M2TypedExecution addresses CancelPendingActions. | f0375 |
| G4M2TypedExecution addresses CheckExecutionCompleteness. | f0376 |
| G4M2TypedExecution addresses CheckFunctionSignatures. | f0377 |
| G4M2TypedExecution addresses CloseSdlInstance. | f0378 |
| G4M2TypedExecution addresses CorrelateActionResult. | f0379 |
| G4M2TypedExecution addresses CreateSdlInstance. | f0380 |
| G4M2TypedExecution addresses InvokeRegisteredFunction. | f0381 |
| G4M2TypedExecution addresses ManageDomainState. | f0382 |
| G4M2TypedExecution addresses PerformDomainOperation. | f0383 |
| G4M2TypedExecution addresses RegisterDomainFunctions. | f0384 |
| G4M2TypedExecution addresses SnapshotDomainState. | f0385 |
| G4M2TypedExecution addresses ValidateActionInput. | f0386 |
| G4M2TypedExecution depends-on G4M1SdlFrontend. | f0387 |
| G4M2TypedExecution has implementation-status = verified. | f0388 |
| G4M2TypedExecution refines G4SdlRuntimeAndBinding. | f0389 |
| G4M3UiDomainBinding depends-on G4M2TypedExecution. | f0397 |
| SdlDispatcher owns CancelPendingActions. | f0848 |
| SdlExecutionGate owns CheckExecutionCompleteness. | f0859 |
| SdlFunctionRegistry owns CheckFunctionSignatures. | f0872 |
| SdlRuntime owns CloseSdlInstance. | f0903 |
| SdlDispatcher owns CorrelateActionResult. | f0849 |
| SdlRuntime owns CreateSdlInstance. | f0904 |
| SdlDispatcher owns InvokeRegisteredFunction. | f0850 |
| SdlStateStore owns ManageDomainState. | f0907 |
| GoDomainImplementation owns PerformDomainOperation. | f0523 |
| SdlFunctionRegistry owns RegisterDomainFunctions. | f0873 |
| SdlStateStore owns SnapshotDomainState. | f0908 |
| SdlDispatcher owns ValidateActionInput. | f0851 |

## G4M3UiDomainBinding

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| BoundActionAccepted illustrates G4M3UiDomainBinding. | f0030 |
| BoundActionRejected illustrates G4M3UiDomainBinding. | f0043 |
| G4M3UiDomainBinding addresses ConnectTypedWidgetHandles. | f0390 |
| G4M3UiDomainBinding addresses DisconnectBindings. | f0391 |
| G4M3UiDomainBinding addresses PublishDomainUpdates. | f0392 |
| G4M3UiDomainBinding addresses ReportBindingDiagnostics. | f0393 |
| G4M3UiDomainBinding addresses ResolveCallbackSymbols. | f0394 |
| G4M3UiDomainBinding addresses RouteDomainBindings. | f0395 |
| G4M3UiDomainBinding depends-on G3M1TypedUiSession. | f0396 |
| G4M3UiDomainBinding depends-on G4M2TypedExecution. | f0397 |
| G4M3UiDomainBinding has implementation-status = verified. | f0398 |
| G4M3UiDomainBinding refines G4SdlRuntimeAndBinding. | f0399 |
| G4M4DomainReload depends-on G4M3UiDomainBinding. | f0404 |
| SdlUiBindingAdapter owns ConnectTypedWidgetHandles. | f0913 |
| SdlUiBindingAdapter owns DisconnectBindings. | f0914 |
| SdlUiBindingAdapter owns PublishDomainUpdates. | f0915 |
| DiagnosticReporter owns ReportBindingDiagnostics. | f0178 |
| SdlUiBindingAdapter owns ResolveCallbackSymbols. | f0916 |
| SdlUiBindingAdapter owns RouteDomainBindings. | f0917 |

## G4M4DomainReload

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G4M4DomainReload addresses CheckDomainStateCompatibility. | f0400 |
| G4M4DomainReload addresses MigrateOrResetDomainState. | f0401 |
| G4M4DomainReload addresses RestartChangedGoProgram. | f0402 |
| G4M4DomainReload depends-on G3M3CompatibleState. | f0403 |
| G4M4DomainReload depends-on G4M3UiDomainBinding. | f0404 |
| G4M4DomainReload has implementation-status = verified. | f0405 |
| G4M4DomainReload refines G4SdlRuntimeAndBinding. | f0406 |
| G5M1GeneratedGo depends-on G4M4DomainReload. | f0413 |
| SdlModelReloadAccepted illustrates G4M4DomainReload. | f0885 |
| DomainStateMigrator owns CheckDomainStateCompatibility. | f0214 |
| DomainStateMigrator owns MigrateOrResetDomainState. | f0215 |
| GoBuildRunner owns RestartChangedGoProgram. | f0506 |

## G4SdlRuntimeAndBinding

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G4M1SdlFrontend refines G4SdlRuntimeAndBinding. | f0374 |
| G4M2TypedExecution refines G4SdlRuntimeAndBinding. | f0389 |
| G4M3UiDomainBinding refines G4SdlRuntimeAndBinding. | f0399 |
| G4M4DomainReload refines G4SdlRuntimeAndBinding. | f0406 |
| G4SdlRuntimeAndBinding delivers TypedDomainBinding. | f0407 |
| G4SdlRuntimeAndBinding has implementation-status = verified. | f0408 |

## G5M1GeneratedGo

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G5M1GeneratedGo addresses BuildGeneratedApplication. | f0409 |
| G5M1GeneratedGo addresses GenerateBindingRegistration. | f0410 |
| G5M1GeneratedGo addresses GenerateModelConstructors. | f0411 |
| G5M1GeneratedGo addresses PreserveHandwrittenSources. | f0412 |
| G5M1GeneratedGo depends-on G4M4DomainReload. | f0413 |
| G5M1GeneratedGo has implementation-status = verified. | f0414 |
| G5M1GeneratedGo refines G5NativeGeneration. | f0415 |
| G5M2BehaviorParity depends-on G5M1GeneratedGo. | f0417 |
| NativeProgramBuilt illustrates G5M1GeneratedGo. | f0624 |
| GoBuildRunner owns BuildGeneratedApplication. | f0505 |
| GoCodeGenerator owns GenerateBindingRegistration. | f0513 |
| GoCodeGenerator owns GenerateModelConstructors. | f0514 |
| GoCodeGenerator owns PreserveHandwrittenSources. | f0515 |

## G5M2BehaviorParity

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G5M2BehaviorParity addresses VerifyNativeBehaviorParity. | f0416 |
| G5M2BehaviorParity depends-on G5M1GeneratedGo. | f0417 |
| G5M2BehaviorParity has implementation-status = verified. | f0418 |
| G5M2BehaviorParity refines G5NativeGeneration. | f0419 |
| G5M3DocumentationExport depends-on G5M2BehaviorParity. | f0424 |
| DevelopmentTools owns VerifyNativeBehaviorParity. | f0176 |

## G5M3DocumentationExport

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G5M3DocumentationExport addresses ComposeHeadlessExport. | f0420 |
| G5M3DocumentationExport addresses ComposeMarkdownDocument. | f0421 |
| G5M3DocumentationExport addresses WriteGeneratedArtifacts. | f0422 |
| G5M3DocumentationExport depends-on G2M4RichContent. | f0423 |
| G5M3DocumentationExport depends-on G5M2BehaviorParity. | f0424 |
| G5M3DocumentationExport depends-on G6M1StaticNavigation. | f0425 |
| G5M3DocumentationExport has implementation-status = verified. | f0426 |
| G5M3DocumentationExport refines G5NativeGeneration. | f0427 |
| G5M4RetirePython depends-on G5M3DocumentationExport. | f0430 |
| StaticFrameExported illustrates G5M3DocumentationExport. | f1103 |
| CommandLineHost owns ComposeHeadlessExport. | f0094 |
| SduiPresentation owns ComposeMarkdownDocument. | f1022 |
| CommandLineHost owns WriteGeneratedArtifacts. | f0095 |

## G5M4RetirePython

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G5M4RetirePython addresses RetireReplacedPythonEntryPoints. | f0428 |
| G5M4RetirePython depends-on G1M3Concept1AndDumps. | f0429 |
| G5M4RetirePython depends-on G5M3DocumentationExport. | f0430 |
| G5M4RetirePython has implementation-status = verified. | f0431 |
| G5M4RetirePython refines G5NativeGeneration. | f0432 |
| DevelopmentTools owns RetireReplacedPythonEntryPoints. | f0175 |

## G5NativeGeneration

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G5M1GeneratedGo refines G5NativeGeneration. | f0415 |
| G5M2BehaviorParity refines G5NativeGeneration. | f0419 |
| G5M3DocumentationExport refines G5NativeGeneration. | f0427 |
| G5M4RetirePython refines G5NativeGeneration. | f0432 |
| G5NativeGeneration delivers DesignDocumentation. | f0433 |
| G5NativeGeneration delivers NativeGoAssembly. | f0434 |
| G5NativeGeneration has implementation-status = verified. | f0435 |

## G6M1StaticNavigation

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G5M3DocumentationExport depends-on G6M1StaticNavigation. | f0425 |
| G6M1StaticNavigation addresses ComposeViewPackage. | f0436 |
| G6M1StaticNavigation addresses ExportModelInventories. | f0437 |
| G6M1StaticNavigation addresses ExportViewpointMarkdown. | f0438 |
| G6M1StaticNavigation addresses GenerateViewNavigation. | f0439 |
| G6M1StaticNavigation addresses IndexViewpointLevels. | f0440 |
| G6M1StaticNavigation addresses PreserveViewAnchors. | f0441 |
| G6M1StaticNavigation addresses ProjectSdlViewpoints. | f0442 |
| G6M1StaticNavigation addresses TraceViewpointFacts. | f0443 |
| G6M1StaticNavigation depends-on G4M1SdlFrontend. | f0444 |
| G6M1StaticNavigation has implementation-status = verified. | f0445 |
| G6M1StaticNavigation refines G6NavigableDocumentation. | f0446 |
| G6M2OnDemandViews depends-on G6M1StaticNavigation. | f0452 |
| G6M5SemanticNotation depends-on G6M1StaticNavigation. | f0471 |
| SelectedViewOpened illustrates G6M1StaticNavigation. | f1056 |
| SdlViewpointGenerator owns ComposeViewPackage. | f0929 |
| SdlViewpointGenerator owns ExportModelInventories. | f0930 |
| SdlViewpointGenerator owns ExportViewpointMarkdown. | f0931 |
| SdlViewpointGenerator owns GenerateViewNavigation. | f0932 |
| SdlViewpointGenerator owns IndexViewpointLevels. | f0933 |
| SdlViewpointGenerator owns PreserveViewAnchors. | f0934 |
| SdlViewpointGenerator owns ProjectSdlViewpoints. | f0936 |
| SdlViewpointGenerator owns TraceViewpointFacts. | f0939 |

## G6M2OnDemandViews

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G6M2OnDemandViews addresses KeyViewRevision. | f0447 |
| G6M2OnDemandViews addresses ProjectSelectedView. | f0448 |
| G6M2OnDemandViews addresses PublishViewBundle. | f0449 |
| G6M2OnDemandViews addresses SelectRelationshipViews. | f0450 |
| G6M2OnDemandViews addresses ValidateViewRequest. | f0451 |
| G6M2OnDemandViews depends-on G6M1StaticNavigation. | f0452 |
| G6M2OnDemandViews has implementation-status = verified. | f0453 |
| G6M2OnDemandViews refines G6NavigableDocumentation. | f0454 |
| G6M3XfmdNavigation depends-on G6M2OnDemandViews. | f0459 |
| InvalidViewSelectionRejected illustrates G6M2OnDemandViews. | f0553 |
| SelectedViewOpened illustrates G6M2OnDemandViews. | f1057 |
| ViewProjectionFailed illustrates G6M2OnDemandViews. | f1357 |
| DocumentBroker owns KeyViewRevision. | f0193 |
| SdlViewpointGenerator owns ProjectSelectedView. | f0937 |
| ViewArtifactStore owns PublishViewBundle. | f1282 |
| SdlViewpointGenerator owns SelectRelationshipViews. | f0938 |
| DocumentBroker owns ValidateViewRequest. | f0196 |

## G6M3XfmdNavigation

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G6M3XfmdNavigation addresses CaptureNavigationTarget. | f0455 |
| G6M3XfmdNavigation addresses DispatchViewOpen. | f0456 |
| G6M3XfmdNavigation addresses ResolveConfiguredViewer. | f0457 |
| G6M3XfmdNavigation addresses RouteDocumentToPane. | f0458 |
| G6M3XfmdNavigation depends-on G6M2OnDemandViews. | f0459 |
| G6M3XfmdNavigation has implementation-status = verified. | f0460 |
| G6M3XfmdNavigation refines G6NavigableDocumentation. | f0461 |
| G6M4SessionPublication depends-on G6M3XfmdNavigation. | f0466 |
| SelectedViewOpened illustrates G6M3XfmdNavigation. | f1058 |
| XfmdDocumentHost owns CaptureNavigationTarget. | f1416 |
| ViewerLaunchAdapter owns DispatchViewOpen. | f1406 |
| ViewerLaunchAdapter owns ResolveConfiguredViewer. | f1407 |
| XfmdDocumentHost owns RouteDocumentToPane. | f1417 |

## G6M4SessionPublication

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G6M4SessionPublication addresses EvictUnusedViewBundles. | f0462 |
| G6M4SessionPublication addresses RejectStaleViewResults. | f0463 |
| G6M4SessionPublication addresses RetainVisibleViewBundle. | f0464 |
| G6M4SessionPublication addresses ServeViewRequests. | f0465 |
| G6M4SessionPublication depends-on G6M3XfmdNavigation. | f0466 |
| G6M4SessionPublication has implementation-status = verified. | f0467 |
| G6M4SessionPublication refines G6NavigableDocumentation. | f0468 |
| InvalidViewSelectionRejected illustrates G6M4SessionPublication. | f0554 |
| SelectedViewOpened illustrates G6M4SessionPublication. | f1059 |
| ViewProjectionFailed illustrates G6M4SessionPublication. | f1358 |
| ViewArtifactStore owns EvictUnusedViewBundles. | f1281 |
| DocumentBroker owns RejectStaleViewResults. | f0194 |
| ViewArtifactStore owns RetainVisibleViewBundle. | f1283 |
| DocumentBroker owns ServeViewRequests. | f0195 |

## G6M5SemanticNotation

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G6M5SemanticNotation addresses ApplySemanticNotation. | f0469 |
| G6M5SemanticNotation addresses VerifyDiagramCapabilities. | f0470 |
| G6M5SemanticNotation depends-on G6M1StaticNavigation. | f0471 |
| G6M5SemanticNotation has implementation-status = verified. | f0472 |
| G6M5SemanticNotation refines G6NavigableDocumentation. | f0473 |
| G6M6ClassViews depends-on G6M5SemanticNotation. | f0476 |
| SdlViewpointGenerator owns ApplySemanticNotation. | f0928 |
| SdlViewpointGenerator owns VerifyDiagramCapabilities. | f0940 |

## G6M6ClassViews

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G6M6ClassViews addresses ProjectClassViews. | f0474 |
| G6M6ClassViews addresses ValidateClassRelations. | f0475 |
| G6M6ClassViews depends-on G6M5SemanticNotation. | f0476 |
| G6M6ClassViews has implementation-status = verified. | f0477 |
| G6M6ClassViews refines G6NavigableDocumentation. | f0478 |
| SdlViewpointGenerator owns ProjectClassViews. | f0935 |
| SdlValidator owns ValidateClassRelations. | f0924 |

## G6NavigableDocumentation

Status: verified.

| Faktum | Kilde-ID |
| --- | --- |
| G6M1StaticNavigation refines G6NavigableDocumentation. | f0446 |
| G6M2OnDemandViews refines G6NavigableDocumentation. | f0454 |
| G6M3XfmdNavigation refines G6NavigableDocumentation. | f0461 |
| G6M4SessionPublication refines G6NavigableDocumentation. | f0468 |
| G6M5SemanticNotation refines G6NavigableDocumentation. | f0473 |
| G6M6ClassViews refines G6NavigableDocumentation. | f0478 |
| G6NavigableDocumentation delivers NavigableDesignDocumentation. | f0479 |
| G6NavigableDocumentation delivers TypedDesignInspection. | f0480 |
| G6NavigableDocumentation has implementation-status = verified. | f0481 |

## InspectDesignSource

Status: unspecified.

| Faktum | Kilde-ID |
| --- | --- |
| ExportUiDocumentation refines InspectDesignSource. | f0234 |
| InspectSdlSource refines InspectDesignSource. | f0539 |
| InspectSduiSource refines InspectDesignSource. | f0540 |

## InspectSdlSource

Status: unspecified.

| Faktum | Kilde-ID |
| --- | --- |
| InspectSdlSource refines InspectDesignSource. | f0539 |

## InspectSduiSource

Status: unspecified.

| Faktum | Kilde-ID |
| --- | --- |
| InspectSduiSource refines InspectDesignSource. | f0540 |

## RealizeDesign

Status: unspecified.

| Faktum | Kilde-ID |
| --- | --- |
| BuildNativeRealization refines RealizeDesign. | f0062 |

## ReloadBoundModels

Status: unspecified.

| Faktum | Kilde-ID |
| --- | --- |
| ReloadBoundModels refines ReloadDesignSession. | f0752 |

## ReloadDesignSession

Status: unspecified.

| Faktum | Kilde-ID |
| --- | --- |
| ReloadBoundModels refines ReloadDesignSession. | f0752 |
| ReloadUiModel refines ReloadDesignSession. | f0795 |

## ReloadUiModel

Status: unspecified.

| Faktum | Kilde-ID |
| --- | --- |
| ReloadUiModel refines ReloadDesignSession. | f0795 |

## RunBoundUiAction

Status: unspecified.

| Faktum | Kilde-ID |
| --- | --- |
| RunBoundUiAction refines RunDesignSession. | f0830 |

## RunDesignSession

Status: unspecified.

| Faktum | Kilde-ID |
| --- | --- |
| RunBoundUiAction refines RunDesignSession. | f0830 |
| RunUnboundUiPreview refines RunDesignSession. | f0831 |

## RunUnboundUiPreview

Status: unspecified.

| Faktum | Kilde-ID |
| --- | --- |
| RunUnboundUiPreview refines RunDesignSession. | f0831 |

## Udekket modellansvar

