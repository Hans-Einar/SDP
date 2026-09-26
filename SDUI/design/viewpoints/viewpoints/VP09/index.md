# VP09 — Dataset, Datagram and persistent Database

[Navigator](../../navigator.md)

Explicit holders, sources, contracts, variants, fields and projections.

- [Data origin and holder: DesignSourceDocuments](VP09-data-DesignSourceDocuments.md)
- [Data origin and holder: UiSessionState](VP09-data-UiSessionState.md)
- [Contract structure: ActionArguments](VP09-contract-ActionArguments.md)
- [Contract structure: ActionOutcome](VP09-contract-ActionOutcome.md)
- [Contract structure: AstArtifactContract](VP09-contract-AstArtifactContract.md)
- [Contract structure: DesignSourceRecord](VP09-contract-DesignSourceRecord.md)
- [Contract structure: FramePresentationCallsProtocol](VP09-contract-FramePresentationCallsProtocol.md)
- [Contract structure: GeneratedGoContract](VP09-contract-GeneratedGoContract.md)
- [Contract structure: GoBuildCallsProtocol](VP09-contract-GoBuildCallsProtocol.md)
- [Contract structure: GoDomainCallsProtocol](VP09-contract-GoDomainCallsProtocol.md)
- [Contract structure: GoGenerationCallsProtocol](VP09-contract-GoGenerationCallsProtocol.md)
- [Contract structure: LayoutArguments](VP09-contract-LayoutArguments.md)
- [Contract structure: LayoutCallsProtocol](VP09-contract-LayoutCallsProtocol.md)
- [Contract structure: ModelReloadCallsProtocol](VP09-contract-ModelReloadCallsProtocol.md)
- [Contract structure: NativeBuildContract](VP09-contract-NativeBuildContract.md)
- [Contract structure: NativeUiActionsProtocol](VP09-contract-NativeUiActionsProtocol.md)
- [Contract structure: NormalizedModelContract](VP09-contract-NormalizedModelContract.md)
- [Contract structure: PreparedFrameContract](VP09-contract-PreparedFrameContract.md)
- [Contract structure: PresentationOutcome](VP09-contract-PresentationOutcome.md)
- [Contract structure: ReloadArguments](VP09-contract-ReloadArguments.md)
- [Contract structure: ReloadOutcome](VP09-contract-ReloadOutcome.md)
- [Contract structure: SdlActionCallsProtocol](VP09-contract-SdlActionCallsProtocol.md)
- [Contract structure: SdlCompilationCallsProtocol](VP09-contract-SdlCompilationCallsProtocol.md)
- [Contract structure: SvgDocumentContract](VP09-contract-SvgDocumentContract.md)
- [Contract structure: SvgExportCallsProtocol](VP09-contract-SvgExportCallsProtocol.md)
- [Contract structure: TokenArtifactContract](VP09-contract-TokenArtifactContract.md)
- [Contract structure: UiAstCallsProtocol](VP09-contract-UiAstCallsProtocol.md)
- [Contract structure: UiCompilationCallsProtocol](VP09-contract-UiCompilationCallsProtocol.md)
- [Contract structure: UiDomainActionsProtocol](VP09-contract-UiDomainActionsProtocol.md)
- [Contract structure: UiGenerationContract](VP09-contract-UiGenerationContract.md)
- [Contract structure: UiGenerationEventsProtocol](VP09-contract-UiGenerationEventsProtocol.md)
- [Contract structure: UiNormalizationCallsProtocol](VP09-contract-UiNormalizationCallsProtocol.md)
- [Contract structure: UiSessionRecord](VP09-contract-UiSessionRecord.md)
- [Contract structure: UiTokenizationCallsProtocol](VP09-contract-UiTokenizationCallsProtocol.md)
- [Contract structure: UiValidationCallsProtocol](VP09-contract-UiValidationCallsProtocol.md)
- [Contract structure: ValidationOutcomeContract](VP09-contract-ValidationOutcomeContract.md)
- [Contract structure: ViewBundleContract](VP09-contract-ViewBundleContract.md)
- [Contract structure: ViewDisplayCallsProtocol](VP09-contract-ViewDisplayCallsProtocol.md)
- [Contract structure: ViewFailureContract](VP09-contract-ViewFailureContract.md)
- [Contract structure: ViewLaunchCallsProtocol](VP09-contract-ViewLaunchCallsProtocol.md)
- [Contract structure: ViewNavigationCallsProtocol](VP09-contract-ViewNavigationCallsProtocol.md)
- [Contract structure: ViewOpenContract](VP09-contract-ViewOpenContract.md)
- [Contract structure: ViewOpenedContract](VP09-contract-ViewOpenedContract.md)
- [Contract structure: ViewProjectionCallsProtocol](VP09-contract-ViewProjectionCallsProtocol.md)
- [Contract structure: ViewPublicationCallsProtocol](VP09-contract-ViewPublicationCallsProtocol.md)
- [Contract structure: ViewReferenceContract](VP09-contract-ViewReferenceContract.md)
- [Contract structure: ViewSelectionContract](VP09-contract-ViewSelectionContract.md)

## Fields and contract properties

Database means a persistent data source, not necessarily SQL.

| Fact | ID |
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

## Projection responsibilities

| Functionality | Dataset | Datagram | Fact |
| --- | --- | --- | --- |
| ProjectUiGeneration | UiSessionState | UiGenerationNotices | f0722 |

