# VP09 — Dataset, Datagram og persistent Database

[Navigator](../../navigator.md)

Eksplisitte holdere, kilde, kontrakter, varianter, felt og projeksjoner.

- [Dataopprinnelse og holder: DesignSourceDocuments](VP09-data-DesignSourceDocuments.md)
- [Dataopprinnelse og holder: UiSessionState](VP09-data-UiSessionState.md)
- [Kontraktstruktur: ActionArguments](VP09-contract-ActionArguments.md)
- [Kontraktstruktur: ActionOutcome](VP09-contract-ActionOutcome.md)
- [Kontraktstruktur: AstArtifactContract](VP09-contract-AstArtifactContract.md)
- [Kontraktstruktur: DesignSourceRecord](VP09-contract-DesignSourceRecord.md)
- [Kontraktstruktur: FramePresentationCallsProtocol](VP09-contract-FramePresentationCallsProtocol.md)
- [Kontraktstruktur: GeneratedGoContract](VP09-contract-GeneratedGoContract.md)
- [Kontraktstruktur: GoBuildCallsProtocol](VP09-contract-GoBuildCallsProtocol.md)
- [Kontraktstruktur: GoDomainCallsProtocol](VP09-contract-GoDomainCallsProtocol.md)
- [Kontraktstruktur: GoGenerationCallsProtocol](VP09-contract-GoGenerationCallsProtocol.md)
- [Kontraktstruktur: LayoutArguments](VP09-contract-LayoutArguments.md)
- [Kontraktstruktur: LayoutCallsProtocol](VP09-contract-LayoutCallsProtocol.md)
- [Kontraktstruktur: ModelReloadCallsProtocol](VP09-contract-ModelReloadCallsProtocol.md)
- [Kontraktstruktur: NativeBuildContract](VP09-contract-NativeBuildContract.md)
- [Kontraktstruktur: NativeUiActionsProtocol](VP09-contract-NativeUiActionsProtocol.md)
- [Kontraktstruktur: NormalizedModelContract](VP09-contract-NormalizedModelContract.md)
- [Kontraktstruktur: PreparedFrameContract](VP09-contract-PreparedFrameContract.md)
- [Kontraktstruktur: PresentationOutcome](VP09-contract-PresentationOutcome.md)
- [Kontraktstruktur: ReloadArguments](VP09-contract-ReloadArguments.md)
- [Kontraktstruktur: ReloadOutcome](VP09-contract-ReloadOutcome.md)
- [Kontraktstruktur: SdlActionCallsProtocol](VP09-contract-SdlActionCallsProtocol.md)
- [Kontraktstruktur: SdlCompilationCallsProtocol](VP09-contract-SdlCompilationCallsProtocol.md)
- [Kontraktstruktur: SvgDocumentContract](VP09-contract-SvgDocumentContract.md)
- [Kontraktstruktur: SvgExportCallsProtocol](VP09-contract-SvgExportCallsProtocol.md)
- [Kontraktstruktur: TokenArtifactContract](VP09-contract-TokenArtifactContract.md)
- [Kontraktstruktur: UiAstCallsProtocol](VP09-contract-UiAstCallsProtocol.md)
- [Kontraktstruktur: UiCompilationCallsProtocol](VP09-contract-UiCompilationCallsProtocol.md)
- [Kontraktstruktur: UiDomainActionsProtocol](VP09-contract-UiDomainActionsProtocol.md)
- [Kontraktstruktur: UiGenerationContract](VP09-contract-UiGenerationContract.md)
- [Kontraktstruktur: UiGenerationEventsProtocol](VP09-contract-UiGenerationEventsProtocol.md)
- [Kontraktstruktur: UiNormalizationCallsProtocol](VP09-contract-UiNormalizationCallsProtocol.md)
- [Kontraktstruktur: UiSessionRecord](VP09-contract-UiSessionRecord.md)
- [Kontraktstruktur: UiTokenizationCallsProtocol](VP09-contract-UiTokenizationCallsProtocol.md)
- [Kontraktstruktur: UiValidationCallsProtocol](VP09-contract-UiValidationCallsProtocol.md)
- [Kontraktstruktur: ValidationOutcomeContract](VP09-contract-ValidationOutcomeContract.md)
- [Kontraktstruktur: ViewBundleContract](VP09-contract-ViewBundleContract.md)
- [Kontraktstruktur: ViewDisplayCallsProtocol](VP09-contract-ViewDisplayCallsProtocol.md)
- [Kontraktstruktur: ViewFailureContract](VP09-contract-ViewFailureContract.md)
- [Kontraktstruktur: ViewLaunchCallsProtocol](VP09-contract-ViewLaunchCallsProtocol.md)
- [Kontraktstruktur: ViewNavigationCallsProtocol](VP09-contract-ViewNavigationCallsProtocol.md)
- [Kontraktstruktur: ViewOpenContract](VP09-contract-ViewOpenContract.md)
- [Kontraktstruktur: ViewOpenedContract](VP09-contract-ViewOpenedContract.md)
- [Kontraktstruktur: ViewProjectionCallsProtocol](VP09-contract-ViewProjectionCallsProtocol.md)
- [Kontraktstruktur: ViewPublicationCallsProtocol](VP09-contract-ViewPublicationCallsProtocol.md)
- [Kontraktstruktur: ViewReferenceContract](VP09-contract-ViewReferenceContract.md)
- [Kontraktstruktur: ViewSelectionContract](VP09-contract-ViewSelectionContract.md)

## Felt og kontraktegenskaper

Database betyr persistent datakilde, ikke nødvendigvis SQL.

| Faktum | ID |
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
| AstArtifact has presence = required. | f0024 |
| AstArtifact has value-type = bytes. | f0025 |
| AstArtifactContract has completeness = closed. | f0026 |
| DesignSourceRecord has completeness = closed. | f0158 |
| FailedViewRequestId has presence = required. | f0240 |
| FailedViewRequestId has value-type = text. | f0241 |
| FramePresentationCallsProtocol has completeness = closed. | f0243 |
| GeneratedGoContract has completeness = closed. | f0495 |
| GeneratedGoSources has presence = required. | f0497 |
| GeneratedGoSources has value-type = bytes. | f0498 |
| GoBuildCallsProtocol has completeness = closed. | f0500 |
| GoDomainCallsProtocol has completeness = closed. | f0520 |
| GoGenerationCallsProtocol has completeness = closed. | f0530 |
| LayoutArguments has completeness = closed. | f0568 |
| LayoutCallsProtocol has completeness = closed. | f0573 |
| LayoutModelArtifact has presence = required. | f0576 |
| LayoutModelArtifact has value-type = bytes. | f0577 |
| LayoutViewportHeight has presence = required. | f0583 |
| LayoutViewportHeight has value-type = decimal. | f0584 |
| LayoutViewportWidth has presence = required. | f0585 |
| LayoutViewportWidth has value-type = decimal. | f0586 |
| ModelDiagnostics has presence = optional. | f0604 |
| ModelDiagnostics has value-type = text. | f0605 |
| ModelIsValid has presence = required. | f0606 |
| ModelIsValid has value-type = boolean. | f0607 |
| ModelReloadCallsProtocol has completeness = closed. | f0609 |
| NativeBuildContract has completeness = closed. | f0613 |
| NativeBuildDiagnostics has presence = optional. | f0616 |
| NativeBuildDiagnostics has value-type = text. | f0617 |
| NativeBuildSucceeded has presence = required. | f0618 |
| NativeBuildSucceeded has value-type = boolean. | f0619 |
| NativeUiActionsProtocol has completeness = closed. | f0633 |
| NormalizedModelArtifact has presence = required. | f0646 |
| NormalizedModelArtifact has value-type = bytes. | f0647 |
| NormalizedModelContract has completeness = closed. | f0648 |
| NoticeGeneration has presence = required. | f0650 |
| NoticeGeneration has value-type = unsigned. | f0651 |
| NoticeVersion has presence = required. | f0652 |
| NoticeVersion has value-type = unsigned. | f0653 |
| OpenViewConsumerId has presence = required. | f0657 |
| OpenViewConsumerId has value-type = text. | f0658 |
| OpenViewEntryPath has presence = required. | f0659 |
| OpenViewEntryPath has value-type = text. | f0660 |
| OpenViewLeaseId has presence = required. | f0661 |
| OpenViewLeaseId has value-type = text. | f0662 |
| OpenViewPaneId has presence = required. | f0663 |
| OpenViewPaneId has value-type = text. | f0664 |
| OpenViewRequestId has presence = required. | f0665 |
| OpenViewRequestId has value-type = text. | f0666 |
| OpenViewRevision has presence = required. | f0667 |
| OpenViewRevision has value-type = text. | f0668 |
| OpenViewWindowId has presence = required. | f0669 |
| OpenViewWindowId has value-type = text. | f0670 |
| OpenedViewRequestId has presence = required. | f0671 |
| OpenedViewRequestId has value-type = text. | f0672 |
| OpenedViewRevision has presence = required. | f0673 |
| OpenedViewRevision has value-type = text. | f0674 |
| PreparedFrameArtifact has presence = required. | f0681 |
| PreparedFrameArtifact has value-type = bytes. | f0682 |
| PreparedFrameContract has completeness = closed. | f0683 |
| PresentationOutcome has completeness = closed. | f0690 |
| PresentationReady has presence = required. | f0692 |
| PresentationReady has value-type = boolean. | f0693 |
| ReloadArguments has completeness = closed. | f0749 |
| ReloadDiagnostic has presence = optional. | f0776 |
| ReloadDiagnostic has value-type = text. | f0777 |
| ReloadOutcome has completeness = closed. | f0778 |
| ReloadPublishedGeneration has presence = optional. | f0784 |
| ReloadPublishedGeneration has value-type = unsigned. | f0785 |
| ReloadSourceRevision has presence = required. | f0791 |
| ReloadSourceRevision has value-type = unsigned. | f0792 |
| ReloadSourceText has presence = required. | f0793 |
| ReloadSourceText has value-type = text. | f0794 |
| SdlActionCallsProtocol has completeness = closed. | f0834 |
| SdlCompilationCallsProtocol has completeness = closed. | f0843 |
| SessionDraft has presence = optional. | f1073 |
| SessionDraft has value-type = text. | f1074 |
| SessionGeneration has presence = required. | f1075 |
| SessionGeneration has value-type = unsigned. | f1076 |
| SourceDocumentRevision has presence = required. | f1079 |
| SourceDocumentRevision has value-type = unsigned. | f1080 |
| SourceDocumentText has presence = required. | f1081 |
| SourceDocumentText has value-type = text. | f1082 |
| SvgDocumentContract has completeness = closed. | f1111 |
| SvgDocumentText has presence = required. | f1113 |
| SvgDocumentText has value-type = text. | f1114 |
| SvgExportCallsProtocol has completeness = closed. | f1116 |
| TokenArtifact has presence = required. | f1119 |
| TokenArtifact has value-type = bytes. | f1120 |
| TokenArtifactContract has completeness = closed. | f1121 |
| UiAstCallsProtocol has completeness = closed. | f1155 |
| UiCompilationCallsProtocol has completeness = closed. | f1174 |
| UiDomainActionsProtocol has completeness = closed. | f1179 |
| UiGenerationContract has completeness = closed. | f1184 |
| UiGenerationEventsProtocol has completeness = closed. | f1187 |
| UiGenerationWire has bit-order = most-significant-first. | f1192 |
| UiGenerationWire has byte-order = big-endian. | f1193 |
| UiNormalizationCallsProtocol has completeness = closed. | f1215 |
| UiSessionRecord has completeness = closed. | f1218 |
| UiTokenizationCallsProtocol has completeness = closed. | f1229 |
| UiValidationCallsProtocol has completeness = closed. | f1233 |
| ValidationOutcomeContract has completeness = closed. | f1274 |
| ViewBundleBytes has presence = required. | f1286 |
| ViewBundleBytes has value-type = bytes. | f1287 |
| ViewBundleContract has completeness = closed. | f1288 |
| ViewBundleRevision has presence = required. | f1294 |
| ViewBundleRevision has value-type = text. | f1295 |
| ViewConsumerId has presence = required. | f1296 |
| ViewConsumerId has value-type = text. | f1297 |
| ViewDisplayCallsProtocol has completeness = closed. | f1299 |
| ViewEntryPath has presence = required. | f1306 |
| ViewEntryPath has value-type = text. | f1307 |
| ViewFailureCode has presence = required. | f1308 |
| ViewFailureCode has value-type = text. | f1309 |
| ViewFailureContract has completeness = closed. | f1310 |
| ViewFailureDiagnostic has presence = required. | f1314 |
| ViewFailureDiagnostic has value-type = text. | f1315 |
| ViewLaunchCallsProtocol has completeness = closed. | f1317 |
| ViewLeaseId has presence = required. | f1323 |
| ViewLeaseId has value-type = text. | f1324 |
| ViewManifestPath has presence = required. | f1325 |
| ViewManifestPath has value-type = text. | f1326 |
| ViewNavigationCallsProtocol has completeness = closed. | f1328 |
| ViewOpenContract has completeness = closed. | f1332 |
| ViewOpenedContract has completeness = closed. | f1340 |
| ViewPaneId has presence = required. | f1346 |
| ViewPaneId has value-type = text. | f1347 |
| ViewProjectId has presence = required. | f1348 |
| ViewProjectId has value-type = text. | f1349 |
| ViewProjectionCallsProtocol has completeness = closed. | f1351 |
| ViewPublicationCallsProtocol has completeness = closed. | f1368 |
| ViewPublishedRevision has presence = required. | f1371 |
| ViewPublishedRevision has value-type = text. | f1372 |
| ViewReferenceContract has completeness = closed. | f1373 |
| ViewRequestId has presence = required. | f1381 |
| ViewRequestId has value-type = text. | f1382 |
| ViewSelectionContract has completeness = closed. | f1386 |
| ViewSelector has presence = required. | f1394 |
| ViewSelector has value-type = text. | f1395 |
| ViewSourceRevision has presence = required. | f1399 |
| ViewSourceRevision has value-type = text. | f1400 |
| ViewWindowId has presence = required. | f1404 |
| ViewWindowId has value-type = text. | f1405 |

## Projeksjonsansvar

| Functionality | Dataset | Datagram | Faktum |
| --- | --- | --- | --- |
| ProjectUiGeneration | UiSessionState | UiGenerationNotices | f0722 |

