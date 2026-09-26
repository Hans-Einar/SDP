# VP08 — Channel contracts and sequences

[Navigator](../../navigator.md)

Explicit scenario steps validated against permits, participation, mode and request/result correlation.

- [Scenario: BoundActionAccepted — mode BoundExecution](VP08-BoundActionAccepted.md)
- [Scenario: BoundActionRejected — mode BoundExecution](VP08-BoundActionRejected.md)
- [Scenario: InteractiveFramePrepared — mode UiPreview](VP08-InteractiveFramePrepared.md)
- [Scenario: InvalidViewSelectionRejected — mode DocumentBrowsing](VP08-InvalidViewSelectionRejected.md)
- [Scenario: NativeProgramBuilt — mode NativeBuild](VP08-NativeProgramBuilt.md)
- [Scenario: SdlModelReloadAccepted — mode BoundLiveEditing](VP08-SdlModelReloadAccepted.md)
- [Scenario: SelectedViewOpened — mode DocumentBrowsing](VP08-SelectedViewOpened.md)
- [Scenario: StaticFrameExported — mode StaticExport](VP08-StaticFrameExported.md)
- [Scenario: UiCompilationAccepted — mode SourceInspection](VP08-UiCompilationAccepted.md)
- [Scenario: UiModelReloadAccepted — mode LiveEditing](VP08-UiModelReloadAccepted.md)
- [Scenario: UiModelReloadRejected — mode LiveEditing](VP08-UiModelReloadRejected.md)
- [Scenario: UnboundLocalAction — mode UiPreview](VP08-UnboundLocalAction.md)
- [Scenario: ViewProjectionFailed — mode DocumentBrowsing](VP08-ViewProjectionFailed.md)

## Derived MessageSet

| Channel | Mode | Datagram | Sender | Receiver | Source IDs |
| --- | --- | --- | --- | --- | --- |
| FramePresentationCalls | UiPreview | PresentFrameRequest | FyneHost | FyneBackend | f0242, f0244, f0253, f0271, f0686 |
| FramePresentationCalls | UiPreview | PresentFrameResult | FyneBackend | FyneHost | f0242, f0245, f0254, f0270, f0689 |
| GoBuildCalls | NativeBuild | BuildGoRequest | CommandLineHost | GoBuildRunner | f0058, f0098, f0499, f0501, f0508 |
| GoBuildCalls | NativeBuild | BuildGoResult | GoBuildRunner | CommandLineHost | f0061, f0097, f0499, f0502, f0509 |
| GoDomainCalls | BoundExecution | DomainActionRequest | SdlDispatcher | GoDomainImplementation | f0208, f0519, f0521, f0525, f0854 |
| GoDomainCalls | BoundExecution | DomainActionResult | GoDomainImplementation | SdlDispatcher | f0211, f0519, f0522, f0527, f0853 |
| GoDomainCalls | UiPreview | DomainActionRequest | FyneHost | GoDomainImplementation | f0208, f0273, f0519, f0521, f0526 |
| GoDomainCalls | UiPreview | DomainActionResult | GoDomainImplementation | FyneHost | f0211, f0272, f0519, f0522, f0528 |
| GoGenerationCalls | NativeBuild | GenerateGoRequest | CommandLineHost | GoCodeGenerator | f0100, f0486, f0517, f0529, f0531 |
| GoGenerationCalls | NativeBuild | GenerateGoResult | GoCodeGenerator | CommandLineHost | f0099, f0489, f0518, f0529, f0532 |
| LayoutCalls | StaticExport | LayoutRequest | CommandLineHost | SduiLayout | f0102, f0572, f0574, f0579, f0995 |
| LayoutCalls | StaticExport | LayoutResult | SduiLayout | CommandLineHost | f0101, f0572, f0575, f0582, f0997 |
| LayoutCalls | UiPreview | LayoutRequest | FyneHost | SduiLayout | f0275, f0572, f0574, f0579, f0996 |
| LayoutCalls | UiPreview | LayoutResult | SduiLayout | FyneHost | f0274, f0572, f0575, f0582, f0998 |
| ModelReloadCalls | BoundLiveEditing | ReloadPublished | ReloadCoordinator | SourceWatcher | f0608, f0610, f0767, f0783, f1092 |
| ModelReloadCalls | BoundLiveEditing | ReloadRejected | ReloadCoordinator | SourceWatcher | f0608, f0611, f0769, f0788, f1094 |
| ModelReloadCalls | BoundLiveEditing | ReloadRequest | SourceWatcher | ReloadCoordinator | f0608, f0612, f0765, f0790, f1096 |
| ModelReloadCalls | LiveEditing | ReloadPublished | ReloadCoordinator | SourceWatcher | f0608, f0610, f0768, f0783, f1093 |
| ModelReloadCalls | LiveEditing | ReloadRejected | ReloadCoordinator | SourceWatcher | f0608, f0611, f0770, f0788, f1095 |
| ModelReloadCalls | LiveEditing | ReloadRequest | SourceWatcher | ReloadCoordinator | f0608, f0612, f0766, f0790, f1097 |
| NativeUiActions | BoundExecution | UiActionRejected | SduiDispatcher | FyneBackend | f0255, f0632, f0634, f0952, f1148 |
| NativeUiActions | BoundExecution | UiActionRequest | FyneBackend | SduiDispatcher | f0257, f0632, f0635, f0951, f1150 |
| NativeUiActions | BoundExecution | UiActionResult | SduiDispatcher | FyneBackend | f0256, f0632, f0636, f0953, f1153 |
| SdlActionCalls | BoundExecution | SdlActionRequest | SdlUiBindingAdapter | SdlDispatcher | f0833, f0835, f0838, f0855, f0920 |
| SdlActionCalls | BoundExecution | SdlActionResult | SdlDispatcher | SdlUiBindingAdapter | f0833, f0836, f0841, f0856, f0919 |
| SdlCompilationCalls | BoundLiveEditing | CompileSdlRequest | ReloadCoordinator | SdlFrontend | f0109, f0772, f0842, f0844, f0869 |
| SdlCompilationCalls | BoundLiveEditing | CompileSdlResult | SdlFrontend | ReloadCoordinator | f0112, f0771, f0842, f0845, f0870 |
| SvgExportCalls | StaticExport | ExportSvgRequest | CommandLineHost | SduiPresentation | f0104, f0227, f1026, f1115, f1117 |
| SvgExportCalls | StaticExport | ExportSvgResult | SduiPresentation | CommandLineHost | f0103, f0230, f1027, f1115, f1118 |
| UiAstCalls | SourceInspection | BuildUiAstRequest | SduiFrontend | SduiParser | f0077, f0965, f1019, f1154, f1156 |
| UiAstCalls | SourceInspection | BuildUiAstResult | SduiParser | SduiFrontend | f0080, f0964, f1020, f1154, f1157 |
| UiCompilationCalls | LiveEditing | CompileUiRejected | SduiFrontend | ReloadCoordinator | f0115, f0773, f0968, f1173, f1175 |
| UiCompilationCalls | LiveEditing | CompileUiRequest | ReloadCoordinator | SduiFrontend | f0117, f0775, f0966, f1173, f1176 |
| UiCompilationCalls | LiveEditing | CompileUiResult | SduiFrontend | ReloadCoordinator | f0120, f0774, f0970, f1173, f1177 |
| UiCompilationCalls | SourceInspection | CompileUiRejected | SduiFrontend | CommandLineHost | f0105, f0115, f0969, f1173, f1175 |
| UiCompilationCalls | SourceInspection | CompileUiRequest | CommandLineHost | SduiFrontend | f0107, f0117, f0967, f1173, f1176 |
| UiCompilationCalls | SourceInspection | CompileUiResult | SduiFrontend | CommandLineHost | f0106, f0120, f0971, f1173, f1177 |
| UiDomainActions | BoundExecution | BoundActionRequest | SduiDispatcher | SdlUiBindingAdapter | f0048, f0921, f0955, f1178, f1180 |
| UiDomainActions | BoundExecution | BoundActionResult | SdlUiBindingAdapter | SduiDispatcher | f0051, f0922, f0954, f1178, f1181 |
| UiGenerationEvents | BoundExecution | UiGenerationNotices | SduiInstanceStore | FyneBackend | f0258, f0984, f1186, f1188, f1190 |
| UiGenerationEvents | LiveEditing | UiGenerationNotices | SduiInstanceStore | FyneBackend | f0259, f0985, f1186, f1188, f1190 |
| UiNormalizationCalls | SourceInspection | NormalizeUiRequest | SduiFrontend | SduiNormalizer | f0642, f0973, f1015, f1214, f1216 |
| UiNormalizationCalls | SourceInspection | NormalizeUiResult | SduiNormalizer | SduiFrontend | f0645, f0972, f1016, f1214, f1217 |
| UiTokenizationCalls | SourceInspection | TokenizeUiRequest | SduiFrontend | SduiLexer | f0975, f1001, f1130, f1228, f1230 |
| UiTokenizationCalls | SourceInspection | TokenizeUiResult | SduiLexer | SduiFrontend | f0974, f1002, f1133, f1228, f1231 |
| UiValidationCalls | SourceInspection | ValidateUiRequest | SduiFrontend | SduiValidator | f0977, f1048, f1232, f1234, f1262 |
| UiValidationCalls | SourceInspection | ValidateUiResult | SduiValidator | SduiFrontend | f0976, f1049, f1232, f1235, f1265 |
| ViewDisplayCalls | DocumentBrowsing | DisplayViewRequest | ViewerLaunchAdapter | XfmdDocumentHost | f0192, f1298, f1300, f1410, f1418 |
| ViewDisplayCalls | DocumentBrowsing | ViewDisplayResult | XfmdDocumentHost | ViewerLaunchAdapter | f1298, f1301, f1305, f1408, f1419 |
| ViewDisplayCalls | DocumentBrowsing | ViewTargetUnavailable | XfmdDocumentHost | ViewerLaunchAdapter | f1298, f1302, f1403, f1409, f1420 |
| ViewLaunchCalls | DocumentBrowsing | LaunchViewRequest | DocumentBroker | ViewerLaunchAdapter | f0198, f0567, f1316, f1318, f1411 |
| ViewLaunchCalls | DocumentBrowsing | ViewLaunchResult | ViewerLaunchAdapter | DocumentBroker | f0197, f1316, f1319, f1322, f1412 |
| ViewNavigationCalls | DocumentBrowsing | SelectViewRequest | XfmdDocumentHost | DocumentBroker | f0199, f1053, f1327, f1329, f1423 |
| ViewNavigationCalls | DocumentBrowsing | ViewOpenedResult | DocumentBroker | XfmdDocumentHost | f0200, f1327, f1330, f1345, f1421 |
| ViewNavigationCalls | DocumentBrowsing | ViewRequestRejected | DocumentBroker | XfmdDocumentHost | f0201, f1327, f1331, f1385, f1422 |
| ViewProjectionCalls | DocumentBrowsing | ProjectViewRequest | DocumentBroker | SdlViewpointGenerator | f0204, f0724, f0941, f1350, f1352 |
| ViewProjectionCalls | DocumentBrowsing | ViewBundleResult | SdlViewpointGenerator | DocumentBroker | f0202, f0942, f1293, f1350, f1353 |
| ViewProjectionCalls | DocumentBrowsing | ViewProjectionRejected | SdlViewpointGenerator | DocumentBroker | f0203, f0943, f1350, f1354, f1366 |
| ViewPublicationCalls | DocumentBrowsing | PublishViewRequest | DocumentBroker | ViewArtifactStore | f0206, f0736, f1284, f1367, f1369 |
| ViewPublicationCalls | DocumentBrowsing | ViewReferenceResult | ViewArtifactStore | DocumentBroker | f0205, f1285, f1367, f1370, f1380 |

