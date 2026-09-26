# VP04 — Interfaces and collaboration

[Navigator](../../navigator.md)

consumes records use; no provider, Channel or executable message flow is inferred.


## Interface use

No Channel or provider is inferred.

| Consumer | Interface | Fact | Line |
| --- | --- | --- | --- |
| CommandLineHost | ExportSinkPort | f0090 | 555 |
| CommandLineHost | PreparedFramePort | f0091 | 556 |
| CommandLineHost | SduiFrontendPort | f0092 | 557 |
| CommandLineHost | SourceSnapshotPort | f0093 | 558 |
| DiagramProvider | DiagramEnginePort | f0181 | 646 |
| DiagramProvider | ResourcePort | f0182 | 647 |
| DomainStateMigrator | DomainStatePort | f0212 | 677 |
| DomainStateMigrator | SdlModelPort | f0213 | 678 |
| FyneBackend | PreparedFramePort | f0246 | 711 |
| FyneBackend | UiSessionPort | f0247 | 712 |
| FyneHost | DomainBindingPort | f0260 | 725 |
| FyneHost | ReloadPort | f0261 | 726 |
| FyneHost | SdlFrontendPort | f0262 | 727 |
| FyneHost | SduiFrontendPort | f0263 | 728 |
| FyneHost | SourceSnapshotPort | f0264 | 729 |
| FyneHost | UiSessionPort | f0265 | 730 |
| FyneHost | WidgetBackendPort | f0266 | 731 |
| GoBuildRunner | BuildToolPort | f0503 | 968 |
| GoBuildRunner | GeneratedArtifactPort | f0504 | 969 |
| GoCodeGenerator | ExecutionProfilePort | f0510 | 975 |
| GoCodeGenerator | SdlModelPort | f0511 | 976 |
| GoCodeGenerator | SduiModelPort | f0512 | 977 |
| MarkdownProvider | DiagramPort | f0592 | 1057 |
| MarkdownProvider | MeasurementPort | f0593 | 1058 |
| MarkdownProvider | ResourcePort | f0594 | 1059 |
| ReloadCoordinator | BindingReloadPort | f0753 | 1218 |
| ReloadCoordinator | DiagnosticPort | f0754 | 1219 |
| ReloadCoordinator | SdlFrontendPort | f0755 | 1220 |
| ReloadCoordinator | SdlReloadPort | f0756 | 1221 |
| ReloadCoordinator | SduiFrontendPort | f0757 | 1222 |
| ReloadCoordinator | SourceSnapshotPort | f0758 | 1223 |
| ReloadCoordinator | UiReloadPort | f0759 | 1224 |
| SdlDispatcher | DomainFunctionPort | f0846 | 1311 |
| SdlDispatcher | DomainStatePort | f0847 | 1312 |
| SdlExecutionGate | DiagnosticPort | f0857 | 1322 |
| SdlExecutionGate | SdlModelPort | f0858 | 1323 |
| SdlFrontend | DiagnosticPort | f0861 | 1326 |
| SdlFrontend | SourceSnapshotPort | f0862 | 1327 |
| SdlFunctionRegistry | DomainFunctionPort | f0871 | 1336 |
| SdlRuntime | DomainFunctionPort | f0896 | 1361 |
| SdlRuntime | SdlModelPort | f0897 | 1362 |
| SdlUiBindingAdapter | DiagnosticPort | f0910 | 1375 |
| SdlUiBindingAdapter | SdlExecutionPort | f0911 | 1376 |
| SdlUiBindingAdapter | UiSessionPort | f0912 | 1377 |
| SduiDispatcher | DomainBindingPort | f0944 | 1409 |
| SduiDispatcher | UiStatePort | f0945 | 1410 |
| SduiFrontend | DiagnosticPort | f0956 | 1421 |
| SduiFrontend | SourceSnapshotPort | f0957 | 1422 |
| SduiLayout | ContentProviderPort | f0986 | 1451 |
| SduiLayout | MeasurementPort | f0987 | 1452 |
| SduiLayout | UiSnapshotPort | f0988 | 1453 |
| SduiPresentation | PreparedFramePort | f1021 | 1486 |
| SduiPropertyStore | UiStatePort | f1028 | 1493 |
| SduiRuntime | DomainBindingPort | f1033 | 1498 |
| SduiRuntime | SduiModelPort | f1034 | 1499 |
| SourceLoader | SourceInputPort | f1083 | 1548 |
| SourceWatcher | FileChangePort | f1088 | 1553 |
| UiStateReconciler | SduiModelPort | f1222 | 1687 |
| UiStateReconciler | UiStatePort | f1223 | 1688 |

