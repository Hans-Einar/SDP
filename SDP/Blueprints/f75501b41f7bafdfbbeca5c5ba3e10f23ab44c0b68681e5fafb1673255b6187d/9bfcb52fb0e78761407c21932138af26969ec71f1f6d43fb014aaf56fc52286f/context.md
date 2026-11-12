# Modeled context

[Overview](index.md)

CHANGE is added/removed; unchanged nodes are CONTEXT, not permission to edit.
PRESERVE obligations are listed separately. Runtime and code coverage remain UNKNOWN.

- added ActivityService uses ActivityCommands as sender of ImplementationRefused in mode ActivityManagement. [TARGET Contracts.design line 118](sources/TARGET/Contracts.design)
- added ActivityEventStream upholds ActivityEventShape. [TARGET Contracts.design line 88](sources/TARGET/Contracts.design)
- added SessionObjective has presence = required. [TARGET Contracts.design line 235](sources/TARGET/Contracts.design)
- added ActivityRecordRepository owns ActivityDocumentStore. [TARGET Responsibilities.design line 51](sources/TARGET/Responsibilities.design)
- added RetainDescriptorRevisions contributes-to NeedScopedWork. [TARGET Responsibilities.design line 101](sources/TARGET/Responsibilities.design)
- added ActivityCommandsDelivery addresses AdoptActivityRecords. [TARGET Responsibilities.design line 28](sources/TARGET/Responsibilities.design)
- added SessionCreated has message-kind = result. [TARGET Contracts.design line 229](sources/TARGET/Contracts.design)
- added RecordEffectShape has-field PlannedPostimages. [TARGET Contracts.design line 218](sources/TARGET/Contracts.design)
- added RecordEffectShape has completeness = closed. [TARGET Contracts.design line 215](sources/TARGET/Contracts.design)
- added ExternalCoordinatorCaller uses ActivityCommands as sender of RequestImplementationAdmission in mode ActivityManagement. [TARGET Contracts.design line 195](sources/TARGET/Contracts.design)
- added OperationIdentity has value-type = text. [TARGET Contracts.design line 209](sources/TARGET/Contracts.design)
- added RecordEffectResultShape has-field EffectDisposition. [TARGET Contracts.design line 213](sources/TARGET/Contracts.design)
- added ExpectedHistoryRevision has value-type = text. [TARGET Contracts.design line 167](sources/TARGET/Contracts.design)
- added OutOfScopeImplementationIsRefused has completeness = closed. [TARGET Scenarios.design line 16](sources/TARGET/Scenarios.design)
- added ActivityDocumentStore holds ActivityDocuments. [TARGET Responsibilities.design line 42](sources/TARGET/Responsibilities.design)
- added AdoptActivityRecords realizes ActivityRecordManagement. [TARGET Responsibilities.design line 66](sources/TARGET/Responsibilities.design)
- added NeedScopedWork supports ControlActivityLifecycle. [TARGET Features.design line 13](sources/TARGET/Features.design)
- added RequestExecutionQuiescence allocated-to SdpToolHost in mode ActivityManagement. [TARGET Responsibilities.design line 96](sources/TARGET/Responsibilities.design)
- added PlannedPostimages has presence = required. [TARGET Contracts.design line 210](sources/TARGET/Contracts.design)
- added ActivityCommandProtocol permits CreateActivityCommand. [TARGET Contracts.design line 61](sources/TARGET/Contracts.design)
- added ExternalCliCaller uses ActivityCommands as receiver of ActivityCreated in mode ActivityManagement. [TARGET Contracts.design line 170](sources/TARGET/Contracts.design)
- added StaleCreationIsRefused step 3 sends ActivityRecordsConflicted from ActivityRecordRepository to ActivityService via ActivityRecordEffects reply-to 2. [TARGET Scenarios.design line 39](sources/TARGET/Scenarios.design)
- added MandateReference has value-type = text. [TARGET Contracts.design line 207](sources/TARGET/Contracts.design)
- added RecoverActivityRecordOperation realizes ActivityRecordManagement. [TARGET Responsibilities.design line 95](sources/TARGET/Responsibilities.design)
- added ExternalCliCaller uses ActivityCommands as sender of CreateSessionCommand in mode ActivityManagement. [TARGET Contracts.design line 180](sources/TARGET/Contracts.design)
- added CheckImplementationScope contributes-to NeedScopedWork. [TARGET Responsibilities.design line 69](sources/TARGET/Responsibilities.design)
- added StaleCreationIsRefused has completeness = closed. [TARGET Scenarios.design line 35](sources/TARGET/Scenarios.design)
- added OwnerCreatesPlannedActivity step 3 sends ActivityRecordsCommitted from ActivityRecordRepository to ActivityService via ActivityRecordEffects reply-to 2. [TARGET Scenarios.design line 32](sources/TARGET/Scenarios.design)
- added RecordReference has value-type = text. [TARGET Contracts.design line 220](sources/TARGET/Contracts.design)
- added ActivityViewReturned upholds ActivityResultShape. [TARGET Contracts.design line 140](sources/TARGET/Contracts.design)
- added ExternalCoordinatorCaller uses ActivityCommands as receiver of ActivityCreated in mode ActivityManagement. [TARGET Contracts.design line 184](sources/TARGET/Contracts.design)
- added ProjectActivityTimeline allocated-to SdpToolHost in mode ActivityManagement. [TARGET Responsibilities.design line 84](sources/TARGET/Responsibilities.design)
- added OperationIdentity has presence = required. [TARGET Contracts.design line 208](sources/TARGET/Contracts.design)
- added EffectRevision has value-type = text. [TARGET Contracts.design line 159](sources/TARGET/Contracts.design)
- added SessionCreationRefused upholds ActivityResultShape. [TARGET Contracts.design line 234](sources/TARGET/Contracts.design)
- added ExternalCliCaller uses ActivityCommands as receiver of ActivityStateChanged in mode ActivityManagement. [TARGET Contracts.design line 173](sources/TARGET/Contracts.design)
- added CoordinatorCreatesSessionDocument has completeness = closed. [TARGET Scenarios.design line 9](sources/TARGET/Scenarios.design)
- added CreatePlannedActivity realizes ActivityRecordManagement. [TARGET Responsibilities.design line 79](sources/TARGET/Responsibilities.design)
- added EventTransition has value-type = text. [TARGET Contracts.design line 165](sources/TARGET/Contracts.design)
- added ActivityCommandProtocol permits ActivityCreated. [TARGET Contracts.design line 56](sources/TARGET/Contracts.design)
- added ActivitySessionShape has-field SessionPlanReference. [TARGET Contracts.design line 128](sources/TARGET/Contracts.design)
- added ActivityRecordsConflicted has message-kind = result. [TARGET Contracts.design line 102](sources/TARGET/Contracts.design)
- added ActivityRecordRepository owns RetainDescriptorRevisions. [TARGET Responsibilities.design line 55](sources/TARGET/Responsibilities.design)
- added EventTime has value-type = text. [TARGET Contracts.design line 163](sources/TARGET/Contracts.design)
- added ChangeActivityStateCommand upholds ActivityCommandShape. [TARGET Contracts.design line 147](sources/TARGET/Contracts.design)
- added ActivityRecordRepository owns PublishActivityRecordSet. [TARGET Responsibilities.design line 53](sources/TARGET/Responsibilities.design)
- added CreateActivityCommand upholds ActivityCommandShape. [TARGET Contracts.design line 149](sources/TARGET/Contracts.design)
- added ResultRevision has value-type = text. [TARGET Contracts.design line 226](sources/TARGET/Contracts.design)
- added RecordEffectShape has-field ExpectedPreimages. [TARGET Contracts.design line 217](sources/TARGET/Contracts.design)
- added EffectRevision has presence = required. [TARGET Contracts.design line 158](sources/TARGET/Contracts.design)
- added SessionCreationRefused has message-kind = result. [TARGET Contracts.design line 232](sources/TARGET/Contracts.design)
- added SdpToolHost contains ActivityService. [TARGET Responsibilities.design line 106](sources/TARGET/Responsibilities.design)
- added DescriptorReference has value-type = text. [TARGET Contracts.design line 153](sources/TARGET/Contracts.design)
- added ActivityProjectionService uses ActivityViews as sender of ActivityViewReturned in mode ActivityManagement. [TARGET Contracts.design line 90](sources/TARGET/Contracts.design)
- added SessionPlanReference has presence = required. [TARGET Contracts.design line 239](sources/TARGET/Contracts.design)
- added ProjectActivityTimeline has state-retention = stateless. [TARGET Responsibilities.design line 86](sources/TARGET/Responsibilities.design)
- added CreatePlannedActivity allocated-to SdpToolHost in mode ActivityManagement. [TARGET Responsibilities.design line 76](sources/TARGET/Responsibilities.design)
- added OwnerCreatesPlannedActivity step 2 sends ApplyActivityRecordSet from ActivityService to ActivityRecordRepository via ActivityRecordEffects. [TARGET Scenarios.design line 31](sources/TARGET/Scenarios.design)
- added RetainDescriptorRevisions has state-retention = stateful. [TARGET Responsibilities.design line 102](sources/TARGET/Responsibilities.design)
- added DescriptorReference has presence = required. [TARGET Contracts.design line 152](sources/TARGET/Contracts.design)
- added ActivityService uses ActivityRecordEffects as sender of ApplyActivityRecordSet in mode ActivityManagement. [TARGET Contracts.design line 123](sources/TARGET/Contracts.design)
- added ActivityCommandProtocol permits ImplementationAdmitted. [TARGET Contracts.design line 63](sources/TARGET/Contracts.design)
- added ActivityService uses ActivityRecordEffects as receiver of ActivityRecordsCommitted in mode ActivityManagement. [TARGET Contracts.design line 121](sources/TARGET/Contracts.design)
- added ImplementationAdmitted has message-kind = result. [TARGET Contracts.design line 198](sources/TARGET/Contracts.design)
- added ActivityCommandProtocol permits SessionCreated. [TARGET Contracts.design line 66](sources/TARGET/Contracts.design)
- added ActivityRecordRepository owns CanonicalManagementHistory. [TARGET Responsibilities.design line 52](sources/TARGET/Responsibilities.design)
- added ActivityRecordRepository owns RecoverActivityRecordOperation. [TARGET Responsibilities.design line 54](sources/TARGET/Responsibilities.design)
- added OwnerCreatesPlannedActivity runs-in ActivityManagement. [TARGET Scenarios.design line 29](sources/TARGET/Scenarios.design)
- added ActivityRecordsCommitted has message-kind = result. [TARGET Contracts.design line 99](sources/TARGET/Contracts.design)
- added RequestExecutionQuiescence contributes-to NeedScopedWork. [TARGET Responsibilities.design line 97](sources/TARGET/Responsibilities.design)
- added PublishActivityRecordSet allocated-to SdpToolHost in mode ActivityManagement. [TARGET Responsibilities.design line 88](sources/TARGET/Responsibilities.design)
- added PublishActivityRecordSet contributes-to NeedScopedWork. [TARGET Responsibilities.design line 89](sources/TARGET/Responsibilities.design)
- added ChangeActivityStateCommand has message-kind = request. [TARGET Contracts.design line 146](sources/TARGET/Contracts.design)
- added ImplementationRefused has message-kind = result. [TARGET Contracts.design line 201](sources/TARGET/Contracts.design)
- added ProjectOwner pursues PrepareActivitySession. [TARGET Features.design line 20](sources/TARGET/Features.design)
- added StaleCreationIsRefused runs-in ActivityManagement. [TARGET Scenarios.design line 36](sources/TARGET/Scenarios.design)
- added ProjectOwner pursues InspectActivityTimeline. [TARGET Features.design line 19](sources/TARGET/Features.design)
- added ExternalViewer uses ActivityViews as sender of InspectActivityCommand in mode ActivityManagement. [TARGET Contracts.design line 197](sources/TARGET/Contracts.design)
- added EventEvidence has presence = required. [TARGET Contracts.design line 160](sources/TARGET/Contracts.design)
- added ActivityViewReturned has message-kind = result. [TARGET Contracts.design line 138](sources/TARGET/Contracts.design)
- added ActivityViewProtocol has completeness = closed. [TARGET Contracts.design line 135](sources/TARGET/Contracts.design)
- added AdoptActivityRecords allocated-to SdpToolHost in mode ActivityManagement. [TARGET Responsibilities.design line 63](sources/TARGET/Responsibilities.design)
- added RequestExecutionQuiescence has state-retention = stateful. [TARGET Responsibilities.design line 98](sources/TARGET/Responsibilities.design)
- added ActivityCreated upholds ActivityResultShape. [TARGET Contracts.design line 76](sources/TARGET/Contracts.design)
- added ActivityDocumentShape has completeness = closed. [TARGET Contracts.design line 80](sources/TARGET/Contracts.design)
- added RetainDescriptorRevisions realizes ActivityRecordManagement. [TARGET Responsibilities.design line 103](sources/TARGET/Responsibilities.design)
- added ActivityCommandProtocol permits ImplementationRefused. [TARGET Contracts.design line 64](sources/TARGET/Contracts.design)
- added CreateSessionCommand upholds ActivityCommandShape. [TARGET Contracts.design line 151](sources/TARGET/Contracts.design)
- added EffectIdentity has value-type = text. [TARGET Contracts.design line 157](sources/TARGET/Contracts.design)
- added ApplyActivityRecordSet has message-kind = request. [TARGET Contracts.design line 142](sources/TARGET/Contracts.design)
- added SdpToolHost provides ActivityRecordManagement. [TARGET Responsibilities.design line 108](sources/TARGET/Responsibilities.design)
- added ExternalCliCaller uses ActivityCommands as receiver of SessionCreated in mode ActivityManagement. [TARGET Contracts.design line 176](sources/TARGET/Contracts.design)
- added ActivityViewProtocol permits InspectActivityCommand. [TARGET Contracts.design line 137](sources/TARGET/Contracts.design)
- added ViewerReadsActivityOverview step 1 sends InspectActivityCommand from ExternalViewer to ActivityProjectionService via ActivityViews. [TARGET Scenarios.design line 44](sources/TARGET/Scenarios.design)
- added StaleCreationIsRefused step 4 sends ActivityCreationRefused from ActivityService to ExternalCliCaller via ActivityCommands reply-to 1. [TARGET Scenarios.design line 40](sources/TARGET/Scenarios.design)
- added RequestExecutionQuiescence realizes ActivityRecordManagement. [TARGET Responsibilities.design line 99](sources/TARGET/Responsibilities.design)
- added OwnerCreatesPlannedActivity step 4 sends ActivityCreated from ActivityService to ExternalCliCaller via ActivityCommands reply-to 1. [TARGET Scenarios.design line 33](sources/TARGET/Scenarios.design)
- added ActivityService owns AdoptActivityRecords. [TARGET Responsibilities.design line 57](sources/TARGET/Responsibilities.design)
- added ActivityEventShape has-field EventEvidence. [TARGET Contracts.design line 85](sources/TARGET/Contracts.design)
- added PlannedPostimages has value-type = text. [TARGET Contracts.design line 211](sources/TARGET/Contracts.design)
- added ExternalCliCaller uses ActivityCommands as receiver of ActivityCreationRefused in mode ActivityManagement. [TARGET Contracts.design line 171](sources/TARGET/Contracts.design)
- added ActivityCommandShape has-field AuthoredIntent. [TARGET Contracts.design line 69](sources/TARGET/Contracts.design)
- added ActivityService owns RequestExecutionQuiescence. [TARGET Responsibilities.design line 61](sources/TARGET/Responsibilities.design)
- added ViewerReadsActivityOverview has completeness = closed. [TARGET Scenarios.design line 42](sources/TARGET/Scenarios.design)
- added CreatePlannedActivity has state-retention = stateful. [TARGET Responsibilities.design line 78](sources/TARGET/Responsibilities.design)
- added AuthoredIntent has presence = required. [TARGET Contracts.design line 144](sources/TARGET/Contracts.design)
- added SessionCreated upholds ActivityResultShape. [TARGET Contracts.design line 231](sources/TARGET/Contracts.design)
- added ActivityCommandProtocol has completeness = closed. [TARGET Contracts.design line 55](sources/TARGET/Contracts.design)
- added ActivityService uses ActivityCommands as receiver of ChangeActivityStateCommand in mode ActivityManagement. [TARGET Contracts.design line 109](sources/TARGET/Contracts.design)
- added OutOfScopeImplementationIsRefused exercises AdmitScopedImplementation. [TARGET Scenarios.design line 15](sources/TARGET/Scenarios.design)
- added CreatePlannedActivity contributes-to NeedScopedWork. [TARGET Responsibilities.design line 77](sources/TARGET/Responsibilities.design)
- added ViewerReadsActivityOverview exercises BrowseActivityDocuments. [TARGET Scenarios.design line 41](sources/TARGET/Scenarios.design)
- added ExternalCoordinatorCaller uses ActivityCommands as receiver of ActivityCreationRefused in mode ActivityManagement. [TARGET Contracts.design line 185](sources/TARGET/Contracts.design)
- added AdoptActivityRecords has state-retention = stateful. [TARGET Responsibilities.design line 65](sources/TARGET/Responsibilities.design)
- added ActivityResultShape has-field ResultDisposition. [TARGET Contracts.design line 107](sources/TARGET/Contracts.design)
- added SdpToolHost contains ActivityRecordRepository. [TARGET Responsibilities.design line 105](sources/TARGET/Responsibilities.design)
- added ReadableActivityHistory supports BrowseActivityDocuments. [TARGET Features.design line 21](sources/TARGET/Features.design)
- added RecordReference has presence = required. [TARGET Contracts.design line 219](sources/TARGET/Contracts.design)
- added InspectActivityCommand upholds ActivityCommandShape. [TARGET Contracts.design line 205](sources/TARGET/Contracts.design)
- added SdpToolHost provides ActivityReading. [TARGET Responsibilities.design line 107](sources/TARGET/Responsibilities.design)
- added PublishActivityRecordSet has state-retention = stateful. [TARGET Responsibilities.design line 90](sources/TARGET/Responsibilities.design)
- added InspectActivityCommand has message-kind = request. [TARGET Contracts.design line 204](sources/TARGET/Contracts.design)
- added ActivityService uses ActivityCommands as sender of SessionCreationRefused in mode ActivityManagement. [TARGET Contracts.design line 120](sources/TARGET/Contracts.design)
- added SelectedActivity has value-type = text. [TARGET Contracts.design line 228](sources/TARGET/Contracts.design)
- added ActivityReading requires CanonicalProjectHistoryLibrary in mode ActivityManagement. [TARGET Responsibilities.design line 47](sources/TARGET/Responsibilities.design)
- added ActivityCommandsDelivery delivers NeedScopedWork. [TARGET Responsibilities.design line 39](sources/TARGET/Responsibilities.design)
- added ActivityService uses ActivityCommands as receiver of CreateSessionCommand in mode ActivityManagement. [TARGET Contracts.design line 111](sources/TARGET/Contracts.design)
- added ActivityCommandsDelivery addresses CreateActivitySession. [TARGET Responsibilities.design line 30](sources/TARGET/Responsibilities.design)
- added ValidateActivityTransition realizes ActivityRecordManagement. [TARGET Responsibilities.design line 112](sources/TARGET/Responsibilities.design)
- added ImplementationAdmitted upholds ActivityResultShape. [TARGET Contracts.design line 200](sources/TARGET/Contracts.design)
- added ActivityProjectionService owns ProjectActivityTimeline. [TARGET Responsibilities.design line 46](sources/TARGET/Responsibilities.design)
- added ActivityDocuments upholds ActivityDocumentShape. [TARGET Contracts.design line 83](sources/TARGET/Contracts.design)
- added ActivityCreationRefused has message-kind = result. [TARGET Contracts.design line 77](sources/TARGET/Contracts.design)
- added ActivityCommandProtocol permits ChangeActivityStateCommand. [TARGET Contracts.design line 60](sources/TARGET/Contracts.design)
- added ActivityViews upholds ActivityViewProtocol. [TARGET Contracts.design line 141](sources/TARGET/Contracts.design)
- added ActivityService uses ActivityCommands as sender of ActivityCreationRefused in mode ActivityManagement. [TARGET Contracts.design line 114](sources/TARGET/Contracts.design)
- added ActivitySessionShape has completeness = closed. [TARGET Contracts.design line 125](sources/TARGET/Contracts.design)
- added ActivityCommandShape has-field ExpectedHistoryRevision. [TARGET Contracts.design line 70](sources/TARGET/Contracts.design)
- added ReadableActivityHistory supports InspectActivityTimeline. [TARGET Features.design line 22](sources/TARGET/Features.design)
- added ActivityEventShape has completeness = closed. [TARGET Contracts.design line 84](sources/TARGET/Contracts.design)
- added CoordinatorCreatesSessionDocument step 2 sends ApplyActivityRecordSet from ActivityService to ActivityRecordRepository via ActivityRecordEffects. [TARGET Scenarios.design line 12](sources/TARGET/Scenarios.design)
- added ActivityService uses ActivityCommands as receiver of RequestImplementationAdmission in mode ActivityManagement. [TARGET Contracts.design line 112](sources/TARGET/Contracts.design)
- added ExternalCliCaller uses ActivityCommands as sender of CreateActivityCommand in mode ActivityManagement. [TARGET Contracts.design line 179](sources/TARGET/Contracts.design)
- added EventTime has presence = required. [TARGET Contracts.design line 162](sources/TARGET/Contracts.design)
- added OwnerChangesActivityState runs-in ActivityManagement. [TARGET Scenarios.design line 22](sources/TARGET/Scenarios.design)
- added ActivityCommandsDelivery delivers ReadableActivityHistory. [TARGET Responsibilities.design line 40](sources/TARGET/Responsibilities.design)
- added SdpToolHost contains ActivityProjectionService. [TARGET Responsibilities.design line 104](sources/TARGET/Responsibilities.design)
- added ActivityDocumentShape has-field MandateReference. [TARGET Contracts.design line 82](sources/TARGET/Contracts.design)
- added ActivityDocumentStore holds ActivitySessionDocuments. [TARGET Responsibilities.design line 43](sources/TARGET/Responsibilities.design)
- added ActivityCommandProtocol permits ActivityStateChangeRefused. [TARGET Contracts.design line 58](sources/TARGET/Contracts.design)
- added ActivityRecordsCommitted replies-to ApplyActivityRecordSet. [TARGET Contracts.design line 100](sources/TARGET/Contracts.design)
- added ImplementationRefused replies-to RequestImplementationAdmission. [TARGET Contracts.design line 202](sources/TARGET/Contracts.design)
- added ActivityCommandsDelivery addresses ValidateActivityTransition. [TARGET Responsibilities.design line 38](sources/TARGET/Responsibilities.design)
- added EventTransition has presence = required. [TARGET Contracts.design line 164](sources/TARGET/Contracts.design)
- added RecordEffectResultShape has-field EffectRevision. [TARGET Contracts.design line 214](sources/TARGET/Contracts.design)
- added ActivityService owns CreatePlannedActivity. [TARGET Responsibilities.design line 60](sources/TARGET/Responsibilities.design)
- added CoordinatorCreatesSessionDocument step 4 sends SessionCreated from ActivityService to ExternalCoordinatorCaller via ActivityCommands reply-to 1. [TARGET Scenarios.design line 14](sources/TARGET/Scenarios.design)
- added ActivityCommandProtocol permits CreateSessionCommand. [TARGET Contracts.design line 62](sources/TARGET/Contracts.design)
- added ActivityRecordRepository uses ActivityRecordEffects as sender of ActivityRecordsCommitted in mode ActivityManagement. [TARGET Contracts.design line 97](sources/TARGET/Contracts.design)
- added ApplyActivityRecordSet upholds RecordEffectShape. [TARGET Contracts.design line 143](sources/TARGET/Contracts.design)
- added ActivityCreated replies-to CreateActivityCommand. [TARGET Contracts.design line 75](sources/TARGET/Contracts.design)
- added ValidateActivityTransition allocated-to SdpToolHost in mode ActivityManagement. [TARGET Responsibilities.design line 109](sources/TARGET/Responsibilities.design)
- added StaleCreationIsRefused exercises CaptureActivityNeed. [TARGET Scenarios.design line 34](sources/TARGET/Scenarios.design)
- added StaleCreationIsRefused step 2 sends ApplyActivityRecordSet from ActivityService to ActivityRecordRepository via ActivityRecordEffects. [TARGET Scenarios.design line 38](sources/TARGET/Scenarios.design)
- added RecoverActivityRecordOperation has state-retention = stateful. [TARGET Responsibilities.design line 94](sources/TARGET/Responsibilities.design)
- added ExternalCliCaller uses ActivityCommands as receiver of ImplementationRefused in mode ActivityManagement. [TARGET Contracts.design line 175](sources/TARGET/Contracts.design)
- added NeedScopedWork supports PrepareActivitySession. [TARGET Features.design line 14](sources/TARGET/Features.design)
- added ActivityProjectionService owns ProjectActivityNavigation. [TARGET Responsibilities.design line 45](sources/TARGET/Responsibilities.design)
- added ActivityService uses ActivityRecordEffects as receiver of ActivityRecordsConflicted in mode ActivityManagement. [TARGET Contracts.design line 122](sources/TARGET/Contracts.design)
- added ActivityService owns CheckImplementationScope. [TARGET Responsibilities.design line 58](sources/TARGET/Responsibilities.design)
- added OwnerCreatesPlannedActivity step 1 sends CreateActivityCommand from ExternalCliCaller to ActivityService via ActivityCommands. [TARGET Scenarios.design line 30](sources/TARGET/Scenarios.design)
- added OwnerCreatesPlannedActivity has completeness = closed. [TARGET Scenarios.design line 28](sources/TARGET/Scenarios.design)
- added SessionPhase has value-type = text. [TARGET Contracts.design line 238](sources/TARGET/Contracts.design)
- added ActivityStateChanged upholds ActivityResultShape. [TARGET Contracts.design line 134](sources/TARGET/Contracts.design)
- added ActivityRecordRepository uses ActivityRecordEffects as receiver of ApplyActivityRecordSet in mode ActivityManagement. [TARGET Contracts.design line 96](sources/TARGET/Contracts.design)
- added ActivityCommandsDelivery addresses CreatePlannedActivity. [TARGET Responsibilities.design line 31](sources/TARGET/Responsibilities.design)
- added ActivityStateChangeRefused replies-to ChangeActivityStateCommand. [TARGET Contracts.design line 130](sources/TARGET/Contracts.design)
- added OutOfScopeImplementationIsRefused step 2 sends ImplementationRefused from ActivityService to ExternalCoordinatorCaller via ActivityCommands reply-to 1. [TARGET Scenarios.design line 19](sources/TARGET/Scenarios.design)
- added RequestImplementationAdmission upholds ActivityCommandShape. [TARGET Contracts.design line 222](sources/TARGET/Contracts.design)
- added ActivitySessionDocuments upholds ActivitySessionShape. [TARGET Contracts.design line 124](sources/TARGET/Contracts.design)
- added ExternalCliCaller uses ActivityCommands as receiver of ImplementationAdmitted in mode ActivityManagement. [TARGET Contracts.design line 174](sources/TARGET/Contracts.design)
- added ActivityService uses ActivityCommands as sender of ImplementationAdmitted in mode ActivityManagement. [TARGET Contracts.design line 117](sources/TARGET/Contracts.design)
- added CoordinatorCreatesSessionDocument runs-in ActivityManagement. [TARGET Scenarios.design line 10](sources/TARGET/Scenarios.design)
- added AuthoredIntent has value-type = text. [TARGET Contracts.design line 145](sources/TARGET/Contracts.design)
- added ActivityRecordProtocol permits ActivityRecordsConflicted. [TARGET Contracts.design line 94](sources/TARGET/Contracts.design)
- added ExternalCliCaller uses ActivityCommands as sender of ChangeActivityStateCommand in mode ActivityManagement. [TARGET Contracts.design line 178](sources/TARGET/Contracts.design)
- added CanonicalManagementHistory holds ActivityEventStream. [TARGET Responsibilities.design line 67](sources/TARGET/Responsibilities.design)
- added ExternalCoordinatorCaller uses ActivityCommands as receiver of ActivityStateChangeRefused in mode ActivityManagement. [TARGET Contracts.design line 186](sources/TARGET/Contracts.design)
- added ActivityService uses ActivityCommands as sender of ActivityStateChanged in mode ActivityManagement. [TARGET Contracts.design line 116](sources/TARGET/Contracts.design)
- added ActivityEventShape has-field EventTransition. [TARGET Contracts.design line 87](sources/TARGET/Contracts.design)
- added ActivityCommandShape has completeness = closed. [TARGET Contracts.design line 68](sources/TARGET/Contracts.design)
- added OwnerChangesActivityState step 1 sends ChangeActivityStateCommand from ExternalCliCaller to ActivityService via ActivityCommands. [TARGET Scenarios.design line 23](sources/TARGET/Scenarios.design)
- added RecordEffectShape has-field EffectIdentity. [TARGET Contracts.design line 216](sources/TARGET/Contracts.design)
- added CheckImplementationScope allocated-to SdpToolHost in mode ActivityManagement. [TARGET Responsibilities.design line 68](sources/TARGET/Responsibilities.design)
- added ProjectOwner pursues CaptureActivityNeed. [TARGET Features.design line 17](sources/TARGET/Features.design)
- added ActivityCommands upholds ActivityCommandProtocol. [TARGET Contracts.design line 73](sources/TARGET/Contracts.design)
- added ActivityDocumentShape has-field DescriptorReference. [TARGET Contracts.design line 81](sources/TARGET/Contracts.design)
- added ExternalCoordinatorCaller uses ActivityCommands as sender of CreateSessionCommand in mode ActivityManagement. [TARGET Contracts.design line 194](sources/TARGET/Contracts.design)
- added ExpectedPreimages has presence = required. [TARGET Contracts.design line 168](sources/TARGET/Contracts.design)
- added ActivityCommandsDelivery addresses ProjectActivityTimeline. [TARGET Responsibilities.design line 33](sources/TARGET/Responsibilities.design)
- added ActivityCommandsDelivery addresses RetainDescriptorRevisions. [TARGET Responsibilities.design line 37](sources/TARGET/Responsibilities.design)
- added NeedScopedWork supports AdmitScopedImplementation. [TARGET Features.design line 11](sources/TARGET/Features.design)
- added ActivityCommandProtocol permits ActivityStateChanged. [TARGET Contracts.design line 59](sources/TARGET/Contracts.design)
- added ActivityService owns CreateActivitySession. [TARGET Responsibilities.design line 59](sources/TARGET/Responsibilities.design)
- added ExternalCoordinatorCaller uses ActivityCommands as sender of CreateActivityCommand in mode ActivityManagement. [TARGET Contracts.design line 193](sources/TARGET/Contracts.design)
- added CreateActivityCommand has message-kind = request. [TARGET Contracts.design line 148](sources/TARGET/Contracts.design)
- added ProjectActivityNavigation realizes ActivityReading. [TARGET Responsibilities.design line 83](sources/TARGET/Responsibilities.design)
- added PublishActivityRecordSet realizes ActivityRecordManagement. [TARGET Responsibilities.design line 91](sources/TARGET/Responsibilities.design)
- added ExternalCliCaller uses ActivityViews as sender of InspectActivityCommand in mode ActivityManagement. [TARGET Contracts.design line 183](sources/TARGET/Contracts.design)
- added RetainDescriptorRevisions allocated-to SdpToolHost in mode ActivityManagement. [TARGET Responsibilities.design line 100](sources/TARGET/Responsibilities.design)
- added ExpectedHistoryRevision has presence = required. [TARGET Contracts.design line 166](sources/TARGET/Contracts.design)
- added SessionObjective has value-type = text. [TARGET Contracts.design line 236](sources/TARGET/Contracts.design)
- added OwnerChangesActivityState step 3 sends ActivityRecordsCommitted from ActivityRecordRepository to ActivityService via ActivityRecordEffects reply-to 2. [TARGET Scenarios.design line 25](sources/TARGET/Scenarios.design)
- added ResultDisposition has presence = required. [TARGET Contracts.design line 223](sources/TARGET/Contracts.design)
- added ExternalCliCaller uses ActivityCommands as receiver of SessionCreationRefused in mode ActivityManagement. [TARGET Contracts.design line 177](sources/TARGET/Contracts.design)
- added ActivityRecordRepository consumes CanonicalProjectHistoryLibrary. [TARGET Responsibilities.design line 50](sources/TARGET/Responsibilities.design)
- added ActivitySessionShape has-field SessionObjective. [TARGET Contracts.design line 126](sources/TARGET/Contracts.design)
- added SessionCreated replies-to CreateSessionCommand. [TARGET Contracts.design line 230](sources/TARGET/Contracts.design)
- added NeedScopedWork supports CaptureActivityNeed. [TARGET Features.design line 12](sources/TARGET/Features.design)
- added SessionPhase has presence = required. [TARGET Contracts.design line 237](sources/TARGET/Contracts.design)
- added ActivityViewProtocol permits ActivityViewReturned. [TARGET Contracts.design line 136](sources/TARGET/Contracts.design)
- added OwnerChangesActivityState step 4 sends ActivityStateChanged from ActivityService to ExternalCliCaller via ActivityCommands reply-to 1. [TARGET Scenarios.design line 26](sources/TARGET/Scenarios.design)
- added CreateActivitySession contributes-to NeedScopedWork. [TARGET Responsibilities.design line 73](sources/TARGET/Responsibilities.design)
- added ActivityProjectionService uses ActivityViews as receiver of InspectActivityCommand in mode ActivityManagement. [TARGET Contracts.design line 89](sources/TARGET/Contracts.design)
- added ActivityResultShape has-field ResultRevision. [TARGET Contracts.design line 108](sources/TARGET/Contracts.design)
- added ProjectOwner pursues BrowseActivityDocuments. [TARGET Features.design line 16](sources/TARGET/Features.design)
- added ProjectOwner pursues ControlActivityLifecycle. [TARGET Features.design line 18](sources/TARGET/Features.design)
- added OutOfScopeImplementationIsRefused step 1 sends RequestImplementationAdmission from ExternalCoordinatorCaller to ActivityService via ActivityCommands. [TARGET Scenarios.design line 18](sources/TARGET/Scenarios.design)
- added SessionPlanReference has value-type = text. [TARGET Contracts.design line 240](sources/TARGET/Contracts.design)
- added ActivityCreationRefused upholds ActivityResultShape. [TARGET Contracts.design line 79](sources/TARGET/Contracts.design)
- added ViewerReadsActivityOverview step 2 sends ActivityViewReturned from ActivityProjectionService to ExternalViewer via ActivityViews reply-to 1. [TARGET Scenarios.design line 45](sources/TARGET/Scenarios.design)
- added ResultDisposition has value-type = text. [TARGET Contracts.design line 224](sources/TARGET/Contracts.design)
- added ActivityViewReturned replies-to InspectActivityCommand. [TARGET Contracts.design line 139](sources/TARGET/Contracts.design)
- added CoordinatorCreatesSessionDocument step 3 sends ActivityRecordsCommitted from ActivityRecordRepository to ActivityService via ActivityRecordEffects reply-to 2. [TARGET Scenarios.design line 13](sources/TARGET/Scenarios.design)
- added EffectIdentity has presence = required. [TARGET Contracts.design line 156](sources/TARGET/Contracts.design)
- added ImplementationRefused upholds ActivityResultShape. [TARGET Contracts.design line 203](sources/TARGET/Contracts.design)
- added EffectDisposition has presence = required. [TARGET Contracts.design line 154](sources/TARGET/Contracts.design)
- added OwnerChangesActivityState has completeness = closed. [TARGET Scenarios.design line 21](sources/TARGET/Scenarios.design)
- added ActivityRecordManagement requires CanonicalProjectHistoryLibrary in mode ActivityManagement. [TARGET Responsibilities.design line 49](sources/TARGET/Responsibilities.design)
- added ExpectedPreimages has value-type = text. [TARGET Contracts.design line 169](sources/TARGET/Contracts.design)
- added ActivityRecordsCommitted upholds RecordEffectResultShape. [TARGET Contracts.design line 101](sources/TARGET/Contracts.design)
- added EventEvidence has value-type = text. [TARGET Contracts.design line 161](sources/TARGET/Contracts.design)
- added AdoptActivityRecords contributes-to NeedScopedWork. [TARGET Responsibilities.design line 64](sources/TARGET/Responsibilities.design)
- added ActivitySessionShape has-field SessionPhase. [TARGET Contracts.design line 127](sources/TARGET/Contracts.design)
- added ActivityRecordsConflicted replies-to ApplyActivityRecordSet. [TARGET Contracts.design line 103](sources/TARGET/Contracts.design)
- unchanged SDPTool contains SdpToolHost. [NOW System.design line 4](sources/NOW/System.design) [TARGET System.design line 7](sources/TARGET/System.design)
- added ActivityRecordProtocol has completeness = closed. [TARGET Contracts.design line 92](sources/TARGET/Contracts.design)
- added ActivityResultShape has completeness = closed. [TARGET Contracts.design line 105](sources/TARGET/Contracts.design)
- added RecoverActivityRecordOperation allocated-to SdpToolHost in mode ActivityManagement. [TARGET Responsibilities.design line 92](sources/TARGET/Responsibilities.design)
- added ActivityCommandsDelivery addresses PublishActivityRecordSet. [TARGET Responsibilities.design line 34](sources/TARGET/Responsibilities.design)
- added CoordinatorCreatesSessionDocument step 1 sends CreateSessionCommand from ExternalCoordinatorCaller to ActivityService via ActivityCommands. [TARGET Scenarios.design line 11](sources/TARGET/Scenarios.design)
- added ExternalCoordinatorCaller uses ActivityCommands as receiver of ActivityStateChanged in mode ActivityManagement. [TARGET Contracts.design line 187](sources/TARGET/Contracts.design)
- added ActivityRecordProtocol permits ApplyActivityRecordSet. [TARGET Contracts.design line 95](sources/TARGET/Contracts.design)
- added ActivityService consumes ActivityExecutionDispositionPort. [TARGET Responsibilities.design line 56](sources/TARGET/Responsibilities.design)
- added ProjectActivityNavigation allocated-to SdpToolHost in mode ActivityManagement. [TARGET Responsibilities.design line 80](sources/TARGET/Responsibilities.design)
- added RecoverActivityRecordOperation contributes-to NeedScopedWork. [TARGET Responsibilities.design line 93](sources/TARGET/Responsibilities.design)
- added CreateActivitySession has state-retention = stateful. [TARGET Responsibilities.design line 74](sources/TARGET/Responsibilities.design)
- added ValidateActivityTransition has state-retention = stateless. [TARGET Responsibilities.design line 111](sources/TARGET/Responsibilities.design)
- added ActivityService uses ActivityCommands as sender of ActivityCreated in mode ActivityManagement. [TARGET Contracts.design line 113](sources/TARGET/Contracts.design)
- added ViewerReadsActivityOverview runs-in ActivityManagement. [TARGET Scenarios.design line 43](sources/TARGET/Scenarios.design)
- added ActivityStateChanged has message-kind = result. [TARGET Contracts.design line 132](sources/TARGET/Contracts.design)
- added ActivityCommandsDelivery addresses RequestExecutionQuiescence. [TARGET Responsibilities.design line 36](sources/TARGET/Responsibilities.design)
- added ActivityCommandShape has-field SelectedActivity. [TARGET Contracts.design line 72](sources/TARGET/Contracts.design)
- added SessionCreationRefused replies-to CreateSessionCommand. [TARGET Contracts.design line 233](sources/TARGET/Contracts.design)
- added ActivityRecordRepository uses ActivityRecordEffects as sender of ActivityRecordsConflicted in mode ActivityManagement. [TARGET Contracts.design line 98](sources/TARGET/Contracts.design)
- added ActivityCommandsDelivery addresses RecoverActivityRecordOperation. [TARGET Responsibilities.design line 35](sources/TARGET/Responsibilities.design)
- added ActivityCommandsDelivery addresses ProjectActivityNavigation. [TARGET Responsibilities.design line 32](sources/TARGET/Responsibilities.design)
- added CreateSessionCommand has message-kind = request. [TARGET Contracts.design line 150](sources/TARGET/Contracts.design)
- added ActivityService uses ActivityCommands as sender of SessionCreated in mode ActivityManagement. [TARGET Contracts.design line 119](sources/TARGET/Contracts.design)
- added ProjectActivityTimeline realizes ActivityReading. [TARGET Responsibilities.design line 87](sources/TARGET/Responsibilities.design)
- added ExternalCliCaller uses ActivityCommands as receiver of ActivityStateChangeRefused in mode ActivityManagement. [TARGET Contracts.design line 172](sources/TARGET/Contracts.design)
- added ExternalCliCaller uses ActivityViews as receiver of ActivityViewReturned in mode ActivityManagement. [TARGET Contracts.design line 182](sources/TARGET/Contracts.design)
- added ActivityStateChangeRefused has message-kind = result. [TARGET Contracts.design line 129](sources/TARGET/Contracts.design)
- added ActivityRecordManagement requires ActivityExecutionDispositionPort in mode ActivityManagement. [TARGET Responsibilities.design line 48](sources/TARGET/Responsibilities.design)
- added ExternalCoordinatorCaller uses ActivityCommands as receiver of SessionCreated in mode ActivityManagement. [TARGET Contracts.design line 190](sources/TARGET/Contracts.design)
- added ActivityCommandShape has-field OperationIdentity. [TARGET Contracts.design line 71](sources/TARGET/Contracts.design)
- added ImplementationAdmitted replies-to RequestImplementationAdmission. [TARGET Contracts.design line 199](sources/TARGET/Contracts.design)
- added ResultRevision has presence = required. [TARGET Contracts.design line 225](sources/TARGET/Contracts.design)
- added ActivityService owns ValidateActivityTransition. [TARGET Responsibilities.design line 62](sources/TARGET/Responsibilities.design)
- added ActivityCreationRefused replies-to CreateActivityCommand. [TARGET Contracts.design line 78](sources/TARGET/Contracts.design)
- added ExternalCoordinatorCaller uses ActivityCommands as sender of ChangeActivityStateCommand in mode ActivityManagement. [TARGET Contracts.design line 192](sources/TARGET/Contracts.design)
- added ActivityResultShape has-field RecordReference. [TARGET Contracts.design line 106](sources/TARGET/Contracts.design)
- added StaleCreationIsRefused step 1 sends CreateActivityCommand from ExternalCliCaller to ActivityService via ActivityCommands. [TARGET Scenarios.design line 37](sources/TARGET/Scenarios.design)
- added ActivityCreated has message-kind = result. [TARGET Contracts.design line 74](sources/TARGET/Contracts.design)
- added SelectedActivity has presence = required. [TARGET Contracts.design line 227](sources/TARGET/Contracts.design)
- added ActivityCommandProtocol permits RequestImplementationAdmission. [TARGET Contracts.design line 65](sources/TARGET/Contracts.design)
- added ExternalCoordinatorCaller uses ActivityCommands as receiver of ImplementationAdmitted in mode ActivityManagement. [TARGET Contracts.design line 188](sources/TARGET/Contracts.design)
- added OwnerChangesActivityState step 2 sends ApplyActivityRecordSet from ActivityService to ActivityRecordRepository via ActivityRecordEffects. [TARGET Scenarios.design line 24](sources/TARGET/Scenarios.design)
- added CreateActivitySession realizes ActivityRecordManagement. [TARGET Responsibilities.design line 75](sources/TARGET/Responsibilities.design)
- added ActivityEventShape has-field EventTime. [TARGET Contracts.design line 86](sources/TARGET/Contracts.design)
- added ProjectActivityNavigation has state-retention = stateless. [TARGET Responsibilities.design line 82](sources/TARGET/Responsibilities.design)
- added ActivityProjectionService consumes CanonicalProjectHistoryLibrary. [TARGET Responsibilities.design line 44](sources/TARGET/Responsibilities.design)
- added ProjectActivityTimeline contributes-to ReadableActivityHistory. [TARGET Responsibilities.design line 85](sources/TARGET/Responsibilities.design)
- added EffectDisposition has value-type = text. [TARGET Contracts.design line 155](sources/TARGET/Contracts.design)
- added ActivityCommandProtocol permits SessionCreationRefused. [TARGET Contracts.design line 67](sources/TARGET/Contracts.design)
- added CoordinatorCreatesSessionDocument exercises PrepareActivitySession. [TARGET Scenarios.design line 8](sources/TARGET/Scenarios.design)
- added CheckImplementationScope realizes ActivityRecordManagement. [TARGET Responsibilities.design line 71](sources/TARGET/Responsibilities.design)
- added MandateReference has presence = required. [TARGET Contracts.design line 206](sources/TARGET/Contracts.design)
- added OwnerChangesActivityState exercises ControlActivityLifecycle. [TARGET Scenarios.design line 20](sources/TARGET/Scenarios.design)
- added ActivityCommandProtocol permits ActivityCreationRefused. [TARGET Contracts.design line 57](sources/TARGET/Contracts.design)
- added ActivityService uses ActivityCommands as sender of ActivityStateChangeRefused in mode ActivityManagement. [TARGET Contracts.design line 115](sources/TARGET/Contracts.design)
- added ExternalViewer uses ActivityViews as receiver of ActivityViewReturned in mode ActivityManagement. [TARGET Contracts.design line 196](sources/TARGET/Contracts.design)
- added ActivityRecordEffects upholds ActivityRecordProtocol. [TARGET Contracts.design line 91](sources/TARGET/Contracts.design)
- added ActivityRecordsConflicted upholds RecordEffectResultShape. [TARGET Contracts.design line 104](sources/TARGET/Contracts.design)
- added RequestImplementationAdmission has message-kind = request. [TARGET Contracts.design line 221](sources/TARGET/Contracts.design)
- added OwnerCreatesPlannedActivity exercises CaptureActivityNeed. [TARGET Scenarios.design line 27](sources/TARGET/Scenarios.design)
- added ValidateActivityTransition contributes-to NeedScopedWork. [TARGET Responsibilities.design line 110](sources/TARGET/Responsibilities.design)
- added ActivityRecordProtocol permits ActivityRecordsCommitted. [TARGET Contracts.design line 93](sources/TARGET/Contracts.design)
- added ActivityStateChangeRefused upholds ActivityResultShape. [TARGET Contracts.design line 131](sources/TARGET/Contracts.design)
- added CreateActivitySession allocated-to SdpToolHost in mode ActivityManagement. [TARGET Responsibilities.design line 72](sources/TARGET/Responsibilities.design)
- added ActivityStateChanged replies-to ChangeActivityStateCommand. [TARGET Contracts.design line 133](sources/TARGET/Contracts.design)
- added ActivityCommandsDelivery addresses CheckImplementationScope. [TARGET Responsibilities.design line 29](sources/TARGET/Responsibilities.design)
- added RecordEffectResultShape has completeness = closed. [TARGET Contracts.design line 212](sources/TARGET/Contracts.design)
- added ActivityCommandsDelivery has implementation-status = planned. [TARGET Responsibilities.design line 41](sources/TARGET/Responsibilities.design)
- added ExternalCliCaller uses ActivityCommands as sender of RequestImplementationAdmission in mode ActivityManagement. [TARGET Contracts.design line 181](sources/TARGET/Contracts.design)
- added ProjectOwner pursues AdmitScopedImplementation. [TARGET Features.design line 15](sources/TARGET/Features.design)
- added ProjectActivityNavigation contributes-to ReadableActivityHistory. [TARGET Responsibilities.design line 81](sources/TARGET/Responsibilities.design)
- added ExternalCoordinatorCaller uses ActivityCommands as receiver of SessionCreationRefused in mode ActivityManagement. [TARGET Contracts.design line 191](sources/TARGET/Contracts.design)
- added CheckImplementationScope has state-retention = stateless. [TARGET Responsibilities.design line 70](sources/TARGET/Responsibilities.design)
- added ExternalCoordinatorCaller uses ActivityCommands as receiver of ImplementationRefused in mode ActivityManagement. [TARGET Contracts.design line 189](sources/TARGET/Contracts.design)
- added ActivityService uses ActivityCommands as receiver of CreateActivityCommand in mode ActivityManagement. [TARGET Contracts.design line 110](sources/TARGET/Contracts.design)
- added OutOfScopeImplementationIsRefused runs-in ActivityManagement. [TARGET Scenarios.design line 17](sources/TARGET/Scenarios.design)

```mermaid
flowchart LR
    f0["ActivityService uses ActivityCommands as sender of ImplementationRefused in mode ActivityManagement. (CHANGE added)"]
    f1["ActivityEventStream upholds ActivityEventShape. (CHANGE added)"]
    f10["OperationIdentity has value-type = text. (CHANGE added)"]
    f100["ExternalCliCaller uses ActivityCommands as receiver of ActivityCreationRefused in mode ActivityManagement. (CHANGE added)"]
    f101["ActivityCommandShape has-field AuthoredIntent. (CHANGE added)"]
    f102["ActivityService owns RequestExecutionQuiescence. (CHANGE added)"]
    f103["ViewerReadsActivityOverview has completeness = closed. (CHANGE added)"]
    f104["CreatePlannedActivity has state-retention = stateful. (CHANGE added)"]
    f105["AuthoredIntent has presence = required. (CHANGE added)"]
    f106["SessionCreated upholds ActivityResultShape. (CHANGE added)"]
    f107["ActivityCommandProtocol has completeness = closed. (CHANGE added)"]
    f108["ActivityService uses ActivityCommands as receiver of ChangeActivityStateCommand in mode ActivityManagement. (CHANGE added)"]
    f109["OutOfScopeImplementationIsRefused exercises AdmitScopedImplementation. (CHANGE added)"]
    f11["RecordEffectResultShape has-field EffectDisposition. (CHANGE added)"]
    f110["CreatePlannedActivity contributes-to NeedScopedWork. (CHANGE added)"]
    f111["ViewerReadsActivityOverview exercises BrowseActivityDocuments. (CHANGE added)"]
    f112["ExternalCoordinatorCaller uses ActivityCommands as receiver of ActivityCreationRefused in mode ActivityManagement. (CHANGE added)"]
    f113["AdoptActivityRecords has state-retention = stateful. (CHANGE added)"]
    f114["ActivityResultShape has-field ResultDisposition. (CHANGE added)"]
    f115["SdpToolHost contains ActivityRecordRepository. (CHANGE added)"]
    f116["ReadableActivityHistory supports BrowseActivityDocuments. (CHANGE added)"]
    f117["RecordReference has presence = required. (CHANGE added)"]
    f118["InspectActivityCommand upholds ActivityCommandShape. (CHANGE added)"]
    f119["SdpToolHost provides ActivityReading. (CHANGE added)"]
    f12["ExpectedHistoryRevision has value-type = text. (CHANGE added)"]
    f120["PublishActivityRecordSet has state-retention = stateful. (CHANGE added)"]
    f121["InspectActivityCommand has message-kind = request. (CHANGE added)"]
    f122["ActivityService uses ActivityCommands as sender of SessionCreationRefused in mode ActivityManagement. (CHANGE added)"]
    f123["SelectedActivity has value-type = text. (CHANGE added)"]
    f124["ActivityReading requires CanonicalProjectHistoryLibrary in mode ActivityManagement. (CHANGE added)"]
    f125["ActivityCommandsDelivery delivers NeedScopedWork. (CHANGE added)"]
    f126["ActivityService uses ActivityCommands as receiver of CreateSessionCommand in mode ActivityManagement. (CHANGE added)"]
    f127["ActivityCommandsDelivery addresses CreateActivitySession. (CHANGE added)"]
    f128["ValidateActivityTransition realizes ActivityRecordManagement. (CHANGE added)"]
    f129["ImplementationAdmitted upholds ActivityResultShape. (CHANGE added)"]
    f13["OutOfScopeImplementationIsRefused has completeness = closed. (CHANGE added)"]
    f130["ActivityProjectionService owns ProjectActivityTimeline. (CHANGE added)"]
    f131["ActivityDocuments upholds ActivityDocumentShape. (CHANGE added)"]
    f132["ActivityCreationRefused has message-kind = result. (CHANGE added)"]
    f133["ActivityCommandProtocol permits ChangeActivityStateCommand. (CHANGE added)"]
    f134["ActivityViews upholds ActivityViewProtocol. (CHANGE added)"]
    f135["ActivityService uses ActivityCommands as sender of ActivityCreationRefused in mode ActivityManagement. (CHANGE added)"]
    f136["ActivitySessionShape has completeness = closed. (CHANGE added)"]
    f137["ActivityCommandShape has-field ExpectedHistoryRevision. (CHANGE added)"]
    f138["ReadableActivityHistory supports InspectActivityTimeline. (CHANGE added)"]
    f139["ActivityEventShape has completeness = closed. (CHANGE added)"]
    f14["ActivityDocumentStore holds ActivityDocuments. (CHANGE added)"]
    f140["CoordinatorCreatesSessionDocument step 2 sends ApplyActivityRecordSet from ActivityService to ActivityRecordRepository via ActivityRecordEffects. (CHANGE added)"]
    f141["ActivityService uses ActivityCommands as receiver of RequestImplementationAdmission in mode ActivityManagement. (CHANGE added)"]
    f142["ExternalCliCaller uses ActivityCommands as sender of CreateActivityCommand in mode ActivityManagement. (CHANGE added)"]
    f143["EventTime has presence = required. (CHANGE added)"]
    f144["OwnerChangesActivityState runs-in ActivityManagement. (CHANGE added)"]
    f145["ActivityCommandsDelivery delivers ReadableActivityHistory. (CHANGE added)"]
    f146["SdpToolHost contains ActivityProjectionService. (CHANGE added)"]
    f147["ActivityDocumentShape has-field MandateReference. (CHANGE added)"]
    f148["ActivityDocumentStore holds ActivitySessionDocuments. (CHANGE added)"]
    f149["ActivityCommandProtocol permits ActivityStateChangeRefused. (CHANGE added)"]
    f15["AdoptActivityRecords realizes ActivityRecordManagement. (CHANGE added)"]
    f150["ActivityRecordsCommitted replies-to ApplyActivityRecordSet. (CHANGE added)"]
    f151["ImplementationRefused replies-to RequestImplementationAdmission. (CHANGE added)"]
    f152["ActivityCommandsDelivery addresses ValidateActivityTransition. (CHANGE added)"]
    f153["EventTransition has presence = required. (CHANGE added)"]
    f154["RecordEffectResultShape has-field EffectRevision. (CHANGE added)"]
    f155["ActivityService owns CreatePlannedActivity. (CHANGE added)"]
    f156["CoordinatorCreatesSessionDocument step 4 sends SessionCreated from ActivityService to ExternalCoordinatorCaller via ActivityCommands reply-to 1. (CHANGE added)"]
    f157["ActivityCommandProtocol permits CreateSessionCommand. (CHANGE added)"]
    f158["ActivityRecordRepository uses ActivityRecordEffects as sender of ActivityRecordsCommitted in mode ActivityManagement. (CHANGE added)"]
    f159["ApplyActivityRecordSet upholds RecordEffectShape. (CHANGE added)"]
    f16["NeedScopedWork supports ControlActivityLifecycle. (CHANGE added)"]
    f160["ActivityCreated replies-to CreateActivityCommand. (CHANGE added)"]
    f161["ValidateActivityTransition allocated-to SdpToolHost in mode ActivityManagement. (CHANGE added)"]
    f162["StaleCreationIsRefused exercises CaptureActivityNeed. (CHANGE added)"]
    f163["StaleCreationIsRefused step 2 sends ApplyActivityRecordSet from ActivityService to ActivityRecordRepository via ActivityRecordEffects. (CHANGE added)"]
    f164["RecoverActivityRecordOperation has state-retention = stateful. (CHANGE added)"]
    f165["ExternalCliCaller uses ActivityCommands as receiver of ImplementationRefused in mode ActivityManagement. (CHANGE added)"]
    f166["NeedScopedWork supports PrepareActivitySession. (CHANGE added)"]
    f167["ActivityProjectionService owns ProjectActivityNavigation. (CHANGE added)"]
    f168["ActivityService uses ActivityRecordEffects as receiver of ActivityRecordsConflicted in mode ActivityManagement. (CHANGE added)"]
    f169["ActivityService owns CheckImplementationScope. (CHANGE added)"]
    f17["RequestExecutionQuiescence allocated-to SdpToolHost in mode ActivityManagement. (CHANGE added)"]
    f170["OwnerCreatesPlannedActivity step 1 sends CreateActivityCommand from ExternalCliCaller to ActivityService via ActivityCommands. (CHANGE added)"]
    f171["OwnerCreatesPlannedActivity has completeness = closed. (CHANGE added)"]
    f172["SessionPhase has value-type = text. (CHANGE added)"]
    f173["ActivityStateChanged upholds ActivityResultShape. (CHANGE added)"]
    f174["ActivityRecordRepository uses ActivityRecordEffects as receiver of ApplyActivityRecordSet in mode ActivityManagement. (CHANGE added)"]
    f175["ActivityCommandsDelivery addresses CreatePlannedActivity. (CHANGE added)"]
    f176["ActivityStateChangeRefused replies-to ChangeActivityStateCommand. (CHANGE added)"]
    f177["OutOfScopeImplementationIsRefused step 2 sends ImplementationRefused from ActivityService to ExternalCoordinatorCaller via ActivityCommands reply-to 1. (CHANGE added)"]
    f178["RequestImplementationAdmission upholds ActivityCommandShape. (CHANGE added)"]
    f179["ActivitySessionDocuments upholds ActivitySessionShape. (CHANGE added)"]
    f18["PlannedPostimages has presence = required. (CHANGE added)"]
    f180["ExternalCliCaller uses ActivityCommands as receiver of ImplementationAdmitted in mode ActivityManagement. (CHANGE added)"]
    f181["ActivityService uses ActivityCommands as sender of ImplementationAdmitted in mode ActivityManagement. (CHANGE added)"]
    f182["CoordinatorCreatesSessionDocument runs-in ActivityManagement. (CHANGE added)"]
    f183["AuthoredIntent has value-type = text. (CHANGE added)"]
    f184["ActivityRecordProtocol permits ActivityRecordsConflicted. (CHANGE added)"]
    f185["ExternalCliCaller uses ActivityCommands as sender of ChangeActivityStateCommand in mode ActivityManagement. (CHANGE added)"]
    f186["CanonicalManagementHistory holds ActivityEventStream. (CHANGE added)"]
    f187["ExternalCoordinatorCaller uses ActivityCommands as receiver of ActivityStateChangeRefused in mode ActivityManagement. (CHANGE added)"]
    f188["ActivityService uses ActivityCommands as sender of ActivityStateChanged in mode ActivityManagement. (CHANGE added)"]
    f189["ActivityEventShape has-field EventTransition. (CHANGE added)"]
    f19["ActivityCommandProtocol permits CreateActivityCommand. (CHANGE added)"]
    f190["ActivityCommandShape has completeness = closed. (CHANGE added)"]
    f191["OwnerChangesActivityState step 1 sends ChangeActivityStateCommand from ExternalCliCaller to ActivityService via ActivityCommands. (CHANGE added)"]
    f192["RecordEffectShape has-field EffectIdentity. (CHANGE added)"]
    f193["CheckImplementationScope allocated-to SdpToolHost in mode ActivityManagement. (CHANGE added)"]
    f194["ProjectOwner pursues CaptureActivityNeed. (CHANGE added)"]
    f195["ActivityCommands upholds ActivityCommandProtocol. (CHANGE added)"]
    f196["ActivityDocumentShape has-field DescriptorReference. (CHANGE added)"]
    f197["ExternalCoordinatorCaller uses ActivityCommands as sender of CreateSessionCommand in mode ActivityManagement. (CHANGE added)"]
    f198["ExpectedPreimages has presence = required. (CHANGE added)"]
    f199["ActivityCommandsDelivery addresses ProjectActivityTimeline. (CHANGE added)"]
    f2["SessionObjective has presence = required. (CHANGE added)"]
    f20["ExternalCliCaller uses ActivityCommands as receiver of ActivityCreated in mode ActivityManagement. (CHANGE added)"]
    f200["ActivityCommandsDelivery addresses RetainDescriptorRevisions. (CHANGE added)"]
    f201["NeedScopedWork supports AdmitScopedImplementation. (CHANGE added)"]
    f202["ActivityCommandProtocol permits ActivityStateChanged. (CHANGE added)"]
    f203["ActivityService owns CreateActivitySession. (CHANGE added)"]
    f204["ExternalCoordinatorCaller uses ActivityCommands as sender of CreateActivityCommand in mode ActivityManagement. (CHANGE added)"]
    f205["CreateActivityCommand has message-kind = request. (CHANGE added)"]
    f206["ProjectActivityNavigation realizes ActivityReading. (CHANGE added)"]
    f207["PublishActivityRecordSet realizes ActivityRecordManagement. (CHANGE added)"]
    f208["ExternalCliCaller uses ActivityViews as sender of InspectActivityCommand in mode ActivityManagement. (CHANGE added)"]
    f209["RetainDescriptorRevisions allocated-to SdpToolHost in mode ActivityManagement. (CHANGE added)"]
    f21["StaleCreationIsRefused step 3 sends ActivityRecordsConflicted from ActivityRecordRepository to ActivityService via ActivityRecordEffects reply-to 2. (CHANGE added)"]
    f210["ExpectedHistoryRevision has presence = required. (CHANGE added)"]
    f211["SessionObjective has value-type = text. (CHANGE added)"]
    f212["OwnerChangesActivityState step 3 sends ActivityRecordsCommitted from ActivityRecordRepository to ActivityService via ActivityRecordEffects reply-to 2. (CHANGE added)"]
    f213["ResultDisposition has presence = required. (CHANGE added)"]
    f214["ExternalCliCaller uses ActivityCommands as receiver of SessionCreationRefused in mode ActivityManagement. (CHANGE added)"]
    f215["ActivityRecordRepository consumes CanonicalProjectHistoryLibrary. (CHANGE added)"]
    f216["ActivitySessionShape has-field SessionObjective. (CHANGE added)"]
    f217["SessionCreated replies-to CreateSessionCommand. (CHANGE added)"]
    f218["NeedScopedWork supports CaptureActivityNeed. (CHANGE added)"]
    f219["SessionPhase has presence = required. (CHANGE added)"]
    f22["MandateReference has value-type = text. (CHANGE added)"]
    f220["ActivityViewProtocol permits ActivityViewReturned. (CHANGE added)"]
    f221["OwnerChangesActivityState step 4 sends ActivityStateChanged from ActivityService to ExternalCliCaller via ActivityCommands reply-to 1. (CHANGE added)"]
    f222["CreateActivitySession contributes-to NeedScopedWork. (CHANGE added)"]
    f223["ActivityProjectionService uses ActivityViews as receiver of InspectActivityCommand in mode ActivityManagement. (CHANGE added)"]
    f224["ActivityResultShape has-field ResultRevision. (CHANGE added)"]
    f225["ProjectOwner pursues BrowseActivityDocuments. (CHANGE added)"]
    f226["ProjectOwner pursues ControlActivityLifecycle. (CHANGE added)"]
    f227["OutOfScopeImplementationIsRefused step 1 sends RequestImplementationAdmission from ExternalCoordinatorCaller to ActivityService via ActivityCommands. (CHANGE added)"]
    f228["SessionPlanReference has value-type = text. (CHANGE added)"]
    f229["ActivityCreationRefused upholds ActivityResultShape. (CHANGE added)"]
    f23["RecoverActivityRecordOperation realizes ActivityRecordManagement. (CHANGE added)"]
    f230["ViewerReadsActivityOverview step 2 sends ActivityViewReturned from ActivityProjectionService to ExternalViewer via ActivityViews reply-to 1. (CHANGE added)"]
    f231["ResultDisposition has value-type = text. (CHANGE added)"]
    f232["ActivityViewReturned replies-to InspectActivityCommand. (CHANGE added)"]
    f233["CoordinatorCreatesSessionDocument step 3 sends ActivityRecordsCommitted from ActivityRecordRepository to ActivityService via ActivityRecordEffects reply-to 2. (CHANGE added)"]
    f234["EffectIdentity has presence = required. (CHANGE added)"]
    f235["ImplementationRefused upholds ActivityResultShape. (CHANGE added)"]
    f236["EffectDisposition has presence = required. (CHANGE added)"]
    f237["OwnerChangesActivityState has completeness = closed. (CHANGE added)"]
    f238["ActivityRecordManagement requires CanonicalProjectHistoryLibrary in mode ActivityManagement. (CHANGE added)"]
    f239["ExpectedPreimages has value-type = text. (CHANGE added)"]
    f24["ExternalCliCaller uses ActivityCommands as sender of CreateSessionCommand in mode ActivityManagement. (CHANGE added)"]
    f240["ActivityRecordsCommitted upholds RecordEffectResultShape. (CHANGE added)"]
    f241["EventEvidence has value-type = text. (CHANGE added)"]
    f242["AdoptActivityRecords contributes-to NeedScopedWork. (CHANGE added)"]
    f243["ActivitySessionShape has-field SessionPhase. (CHANGE added)"]
    f244["ActivityRecordsConflicted replies-to ApplyActivityRecordSet. (CHANGE added)"]
    f245["SDPTool contains SdpToolHost. (PRESERVE unchanged)"]
    f246["ActivityRecordProtocol has completeness = closed. (CHANGE added)"]
    f247["ActivityResultShape has completeness = closed. (CHANGE added)"]
    f248["RecoverActivityRecordOperation allocated-to SdpToolHost in mode ActivityManagement. (CHANGE added)"]
    f249["ActivityCommandsDelivery addresses PublishActivityRecordSet. (CHANGE added)"]
    f25["CheckImplementationScope contributes-to NeedScopedWork. (CHANGE added)"]
    f250["CoordinatorCreatesSessionDocument step 1 sends CreateSessionCommand from ExternalCoordinatorCaller to ActivityService via ActivityCommands. (CHANGE added)"]
    f251["ExternalCoordinatorCaller uses ActivityCommands as receiver of ActivityStateChanged in mode ActivityManagement. (CHANGE added)"]
    f252["ActivityRecordProtocol permits ApplyActivityRecordSet. (CHANGE added)"]
    f253["ActivityService consumes ActivityExecutionDispositionPort. (CHANGE added)"]
    f254["ProjectActivityNavigation allocated-to SdpToolHost in mode ActivityManagement. (CHANGE added)"]
    f255["RecoverActivityRecordOperation contributes-to NeedScopedWork. (CHANGE added)"]
    f256["CreateActivitySession has state-retention = stateful. (CHANGE added)"]
    f257["ValidateActivityTransition has state-retention = stateless. (CHANGE added)"]
    f258["ActivityService uses ActivityCommands as sender of ActivityCreated in mode ActivityManagement. (CHANGE added)"]
    f259["ViewerReadsActivityOverview runs-in ActivityManagement. (CHANGE added)"]
    f26["StaleCreationIsRefused has completeness = closed. (CHANGE added)"]
    f260["ActivityStateChanged has message-kind = result. (CHANGE added)"]
    f261["ActivityCommandsDelivery addresses RequestExecutionQuiescence. (CHANGE added)"]
    f262["ActivityCommandShape has-field SelectedActivity. (CHANGE added)"]
    f263["SessionCreationRefused replies-to CreateSessionCommand. (CHANGE added)"]
    f264["ActivityRecordRepository uses ActivityRecordEffects as sender of ActivityRecordsConflicted in mode ActivityManagement. (CHANGE added)"]
    f265["ActivityCommandsDelivery addresses RecoverActivityRecordOperation. (CHANGE added)"]
    f266["ActivityCommandsDelivery addresses ProjectActivityNavigation. (CHANGE added)"]
    f267["CreateSessionCommand has message-kind = request. (CHANGE added)"]
    f268["ActivityService uses ActivityCommands as sender of SessionCreated in mode ActivityManagement. (CHANGE added)"]
    f269["ProjectActivityTimeline realizes ActivityReading. (CHANGE added)"]
    f27["OwnerCreatesPlannedActivity step 3 sends ActivityRecordsCommitted from ActivityRecordRepository to ActivityService via ActivityRecordEffects reply-to 2. (CHANGE added)"]
    f270["ExternalCliCaller uses ActivityCommands as receiver of ActivityStateChangeRefused in mode ActivityManagement. (CHANGE added)"]
    f271["ExternalCliCaller uses ActivityViews as receiver of ActivityViewReturned in mode ActivityManagement. (CHANGE added)"]
    f272["ActivityStateChangeRefused has message-kind = result. (CHANGE added)"]
    f273["ActivityRecordManagement requires ActivityExecutionDispositionPort in mode ActivityManagement. (CHANGE added)"]
    f274["ExternalCoordinatorCaller uses ActivityCommands as receiver of SessionCreated in mode ActivityManagement. (CHANGE added)"]
    f275["ActivityCommandShape has-field OperationIdentity. (CHANGE added)"]
    f276["ImplementationAdmitted replies-to RequestImplementationAdmission. (CHANGE added)"]
    f277["ResultRevision has presence = required. (CHANGE added)"]
    f278["ActivityService owns ValidateActivityTransition. (CHANGE added)"]
    f279["ActivityCreationRefused replies-to CreateActivityCommand. (CHANGE added)"]
    f28["RecordReference has value-type = text. (CHANGE added)"]
    f280["ExternalCoordinatorCaller uses ActivityCommands as sender of ChangeActivityStateCommand in mode ActivityManagement. (CHANGE added)"]
    f281["ActivityResultShape has-field RecordReference. (CHANGE added)"]
    f282["StaleCreationIsRefused step 1 sends CreateActivityCommand from ExternalCliCaller to ActivityService via ActivityCommands. (CHANGE added)"]
    f283["ActivityCreated has message-kind = result. (CHANGE added)"]
    f284["SelectedActivity has presence = required. (CHANGE added)"]
    f285["ActivityCommandProtocol permits RequestImplementationAdmission. (CHANGE added)"]
    f286["ExternalCoordinatorCaller uses ActivityCommands as receiver of ImplementationAdmitted in mode ActivityManagement. (CHANGE added)"]
    f287["OwnerChangesActivityState step 2 sends ApplyActivityRecordSet from ActivityService to ActivityRecordRepository via ActivityRecordEffects. (CHANGE added)"]
    f288["CreateActivitySession realizes ActivityRecordManagement. (CHANGE added)"]
    f289["ActivityEventShape has-field EventTime. (CHANGE added)"]
    f29["ActivityViewReturned upholds ActivityResultShape. (CHANGE added)"]
    f290["ProjectActivityNavigation has state-retention = stateless. (CHANGE added)"]
    f291["ActivityProjectionService consumes CanonicalProjectHistoryLibrary. (CHANGE added)"]
    f292["ProjectActivityTimeline contributes-to ReadableActivityHistory. (CHANGE added)"]
    f293["EffectDisposition has value-type = text. (CHANGE added)"]
    f294["ActivityCommandProtocol permits SessionCreationRefused. (CHANGE added)"]
    f295["CoordinatorCreatesSessionDocument exercises PrepareActivitySession. (CHANGE added)"]
    f296["CheckImplementationScope realizes ActivityRecordManagement. (CHANGE added)"]
    f297["MandateReference has presence = required. (CHANGE added)"]
    f298["OwnerChangesActivityState exercises ControlActivityLifecycle. (CHANGE added)"]
    f299["ActivityCommandProtocol permits ActivityCreationRefused. (CHANGE added)"]
    f3["ActivityRecordRepository owns ActivityDocumentStore. (CHANGE added)"]
    f30["ExternalCoordinatorCaller uses ActivityCommands as receiver of ActivityCreated in mode ActivityManagement. (CHANGE added)"]
    f300["ActivityService uses ActivityCommands as sender of ActivityStateChangeRefused in mode ActivityManagement. (CHANGE added)"]
    f301["ExternalViewer uses ActivityViews as receiver of ActivityViewReturned in mode ActivityManagement. (CHANGE added)"]
    f302["ActivityRecordEffects upholds ActivityRecordProtocol. (CHANGE added)"]
    f303["ActivityRecordsConflicted upholds RecordEffectResultShape. (CHANGE added)"]
    f304["RequestImplementationAdmission has message-kind = request. (CHANGE added)"]
    f305["OwnerCreatesPlannedActivity exercises CaptureActivityNeed. (CHANGE added)"]
    f306["ValidateActivityTransition contributes-to NeedScopedWork. (CHANGE added)"]
    f307["ActivityRecordProtocol permits ActivityRecordsCommitted. (CHANGE added)"]
    f308["ActivityStateChangeRefused upholds ActivityResultShape. (CHANGE added)"]
    f309["CreateActivitySession allocated-to SdpToolHost in mode ActivityManagement. (CHANGE added)"]
    f31["ProjectActivityTimeline allocated-to SdpToolHost in mode ActivityManagement. (CHANGE added)"]
    f310["ActivityStateChanged replies-to ChangeActivityStateCommand. (CHANGE added)"]
    f311["ActivityCommandsDelivery addresses CheckImplementationScope. (CHANGE added)"]
    f312["RecordEffectResultShape has completeness = closed. (CHANGE added)"]
    f313["ActivityCommandsDelivery has implementation-status = planned. (CHANGE added)"]
    f314["ExternalCliCaller uses ActivityCommands as sender of RequestImplementationAdmission in mode ActivityManagement. (CHANGE added)"]
    f315["ProjectOwner pursues AdmitScopedImplementation. (CHANGE added)"]
    f316["ProjectActivityNavigation contributes-to ReadableActivityHistory. (CHANGE added)"]
    f317["ExternalCoordinatorCaller uses ActivityCommands as receiver of SessionCreationRefused in mode ActivityManagement. (CHANGE added)"]
    f318["CheckImplementationScope has state-retention = stateless. (CHANGE added)"]
    f319["ExternalCoordinatorCaller uses ActivityCommands as receiver of ImplementationRefused in mode ActivityManagement. (CHANGE added)"]
    f32["OperationIdentity has presence = required. (CHANGE added)"]
    f320["ActivityService uses ActivityCommands as receiver of CreateActivityCommand in mode ActivityManagement. (CHANGE added)"]
    f321["OutOfScopeImplementationIsRefused runs-in ActivityManagement. (CHANGE added)"]
    f33["EffectRevision has value-type = text. (CHANGE added)"]
    f34["SessionCreationRefused upholds ActivityResultShape. (CHANGE added)"]
    f35["ExternalCliCaller uses ActivityCommands as receiver of ActivityStateChanged in mode ActivityManagement. (CHANGE added)"]
    f36["CoordinatorCreatesSessionDocument has completeness = closed. (CHANGE added)"]
    f37["CreatePlannedActivity realizes ActivityRecordManagement. (CHANGE added)"]
    f38["EventTransition has value-type = text. (CHANGE added)"]
    f39["ActivityCommandProtocol permits ActivityCreated. (CHANGE added)"]
    f4["RetainDescriptorRevisions contributes-to NeedScopedWork. (CHANGE added)"]
    f40["ActivitySessionShape has-field SessionPlanReference. (CHANGE added)"]
    f41["ActivityRecordsConflicted has message-kind = result. (CHANGE added)"]
    f42["ActivityRecordRepository owns RetainDescriptorRevisions. (CHANGE added)"]
    f43["EventTime has value-type = text. (CHANGE added)"]
    f44["ChangeActivityStateCommand upholds ActivityCommandShape. (CHANGE added)"]
    f45["ActivityRecordRepository owns PublishActivityRecordSet. (CHANGE added)"]
    f46["CreateActivityCommand upholds ActivityCommandShape. (CHANGE added)"]
    f47["ResultRevision has value-type = text. (CHANGE added)"]
    f48["RecordEffectShape has-field ExpectedPreimages. (CHANGE added)"]
    f49["EffectRevision has presence = required. (CHANGE added)"]
    f5["ActivityCommandsDelivery addresses AdoptActivityRecords. (CHANGE added)"]
    f50["SessionCreationRefused has message-kind = result. (CHANGE added)"]
    f51["SdpToolHost contains ActivityService. (CHANGE added)"]
    f52["DescriptorReference has value-type = text. (CHANGE added)"]
    f53["ActivityProjectionService uses ActivityViews as sender of ActivityViewReturned in mode ActivityManagement. (CHANGE added)"]
    f54["SessionPlanReference has presence = required. (CHANGE added)"]
    f55["ProjectActivityTimeline has state-retention = stateless. (CHANGE added)"]
    f56["CreatePlannedActivity allocated-to SdpToolHost in mode ActivityManagement. (CHANGE added)"]
    f57["OwnerCreatesPlannedActivity step 2 sends ApplyActivityRecordSet from ActivityService to ActivityRecordRepository via ActivityRecordEffects. (CHANGE added)"]
    f58["RetainDescriptorRevisions has state-retention = stateful. (CHANGE added)"]
    f59["DescriptorReference has presence = required. (CHANGE added)"]
    f6["SessionCreated has message-kind = result. (CHANGE added)"]
    f60["ActivityService uses ActivityRecordEffects as sender of ApplyActivityRecordSet in mode ActivityManagement. (CHANGE added)"]
    f61["ActivityCommandProtocol permits ImplementationAdmitted. (CHANGE added)"]
    f62["ActivityService uses ActivityRecordEffects as receiver of ActivityRecordsCommitted in mode ActivityManagement. (CHANGE added)"]
    f63["ImplementationAdmitted has message-kind = result. (CHANGE added)"]
    f64["ActivityCommandProtocol permits SessionCreated. (CHANGE added)"]
    f65["ActivityRecordRepository owns CanonicalManagementHistory. (CHANGE added)"]
    f66["ActivityRecordRepository owns RecoverActivityRecordOperation. (CHANGE added)"]
    f67["OwnerCreatesPlannedActivity runs-in ActivityManagement. (CHANGE added)"]
    f68["ActivityRecordsCommitted has message-kind = result. (CHANGE added)"]
    f69["RequestExecutionQuiescence contributes-to NeedScopedWork. (CHANGE added)"]
    f7["RecordEffectShape has-field PlannedPostimages. (CHANGE added)"]
    f70["PublishActivityRecordSet allocated-to SdpToolHost in mode ActivityManagement. (CHANGE added)"]
    f71["PublishActivityRecordSet contributes-to NeedScopedWork. (CHANGE added)"]
    f72["ChangeActivityStateCommand has message-kind = request. (CHANGE added)"]
    f73["ImplementationRefused has message-kind = result. (CHANGE added)"]
    f74["ProjectOwner pursues PrepareActivitySession. (CHANGE added)"]
    f75["StaleCreationIsRefused runs-in ActivityManagement. (CHANGE added)"]
    f76["ProjectOwner pursues InspectActivityTimeline. (CHANGE added)"]
    f77["ExternalViewer uses ActivityViews as sender of InspectActivityCommand in mode ActivityManagement. (CHANGE added)"]
    f78["EventEvidence has presence = required. (CHANGE added)"]
    f79["ActivityViewReturned has message-kind = result. (CHANGE added)"]
    f8["RecordEffectShape has completeness = closed. (CHANGE added)"]
    f80["ActivityViewProtocol has completeness = closed. (CHANGE added)"]
    f81["AdoptActivityRecords allocated-to SdpToolHost in mode ActivityManagement. (CHANGE added)"]
    f82["RequestExecutionQuiescence has state-retention = stateful. (CHANGE added)"]
    f83["ActivityCreated upholds ActivityResultShape. (CHANGE added)"]
    f84["ActivityDocumentShape has completeness = closed. (CHANGE added)"]
    f85["RetainDescriptorRevisions realizes ActivityRecordManagement. (CHANGE added)"]
    f86["ActivityCommandProtocol permits ImplementationRefused. (CHANGE added)"]
    f87["CreateSessionCommand upholds ActivityCommandShape. (CHANGE added)"]
    f88["EffectIdentity has value-type = text. (CHANGE added)"]
    f89["ApplyActivityRecordSet has message-kind = request. (CHANGE added)"]
    f9["ExternalCoordinatorCaller uses ActivityCommands as sender of RequestImplementationAdmission in mode ActivityManagement. (CHANGE added)"]
    f90["SdpToolHost provides ActivityRecordManagement. (CHANGE added)"]
    f91["ExternalCliCaller uses ActivityCommands as receiver of SessionCreated in mode ActivityManagement. (CHANGE added)"]
    f92["ActivityViewProtocol permits InspectActivityCommand. (CHANGE added)"]
    f93["ViewerReadsActivityOverview step 1 sends InspectActivityCommand from ExternalViewer to ActivityProjectionService via ActivityViews. (CHANGE added)"]
    f94["StaleCreationIsRefused step 4 sends ActivityCreationRefused from ActivityService to ExternalCliCaller via ActivityCommands reply-to 1. (CHANGE added)"]
    f95["RequestExecutionQuiescence realizes ActivityRecordManagement. (CHANGE added)"]
    f96["OwnerCreatesPlannedActivity step 4 sends ActivityCreated from ActivityService to ExternalCliCaller via ActivityCommands reply-to 1. (CHANGE added)"]
    f97["ActivityService owns AdoptActivityRecords. (CHANGE added)"]
    f98["ActivityEventShape has-field EventEvidence. (CHANGE added)"]
    f99["PlannedPostimages has value-type = text. (CHANGE added)"]
    n0["ActivityCommandsDelivery (activity CHANGE added)"]
    n1["ProjectOwner (actor CHANGE added)"]
    n10["ActivityDocumentShape (contract CHANGE added)"]
    n11["ActivityEventShape (contract CHANGE added)"]
    n12["ActivityRecordProtocol (contract CHANGE added)"]
    n13["ActivityResultShape (contract CHANGE added)"]
    n14["ActivitySessionShape (contract CHANGE added)"]
    n15["ActivityViewProtocol (contract CHANGE added)"]
    n16["RecordEffectResultShape (contract CHANGE added)"]
    n17["RecordEffectShape (contract CHANGE added)"]
    n18["ActivityDocumentStore (database CHANGE added)"]
    n19["CanonicalManagementHistory (database CHANGE added)"]
    n2["ActivityReading (capability CHANGE added)"]
    n20["ActivityDocuments (dataset CHANGE added)"]
    n21["ActivityEventStream (dataset CHANGE added)"]
    n22["ActivitySessionDocuments (dataset CHANGE added)"]
    n23["NeedScopedWork (feature CHANGE added)"]
    n24["ReadableActivityHistory (feature CHANGE added)"]
    n25["AuthoredIntent (field CHANGE added)"]
    n26["DescriptorReference (field CHANGE added)"]
    n27["EffectDisposition (field CHANGE added)"]
    n28["EffectIdentity (field CHANGE added)"]
    n29["EffectRevision (field CHANGE added)"]
    n3["ActivityRecordManagement (capability CHANGE added)"]
    n30["EventEvidence (field CHANGE added)"]
    n31["EventTime (field CHANGE added)"]
    n32["EventTransition (field CHANGE added)"]
    n33["ExpectedHistoryRevision (field CHANGE added)"]
    n34["ExpectedPreimages (field CHANGE added)"]
    n35["MandateReference (field CHANGE added)"]
    n36["OperationIdentity (field CHANGE added)"]
    n37["PlannedPostimages (field CHANGE added)"]
    n38["RecordReference (field CHANGE added)"]
    n39["ResultDisposition (field CHANGE added)"]
    n4["ActivityCommands (channel CHANGE added)"]
    n40["ResultRevision (field CHANGE added)"]
    n41["SelectedActivity (field CHANGE added)"]
    n42["SessionObjective (field CHANGE added)"]
    n43["SessionPhase (field CHANGE added)"]
    n44["SessionPlanReference (field CHANGE added)"]
    n45["AdoptActivityRecords (functionality CHANGE added)"]
    n46["CheckImplementationScope (functionality CHANGE added)"]
    n47["CreateActivitySession (functionality CHANGE added)"]
    n48["CreatePlannedActivity (functionality CHANGE added)"]
    n49["ProjectActivityNavigation (functionality CHANGE added)"]
    n5["ActivityRecordEffects (channel CHANGE added)"]
    n50["ProjectActivityTimeline (functionality CHANGE added)"]
    n51["PublishActivityRecordSet (functionality CHANGE added)"]
    n52["RecoverActivityRecordOperation (functionality CHANGE added)"]
    n53["RequestExecutionQuiescence (functionality CHANGE added)"]
    n54["RetainDescriptorRevisions (functionality CHANGE added)"]
    n55["ValidateActivityTransition (functionality CHANGE added)"]
    n56["ActivityExecutionDispositionPort (interface CHANGE added)"]
    n57["CanonicalProjectHistoryLibrary (interface CHANGE added)"]
    n58["ActivityCreated (message CHANGE added)"]
    n59["ActivityCreationRefused (message CHANGE added)"]
    n6["ActivityViews (channel CHANGE added)"]
    n60["ActivityRecordsCommitted (message CHANGE added)"]
    n61["ActivityRecordsConflicted (message CHANGE added)"]
    n62["ActivityStateChangeRefused (message CHANGE added)"]
    n63["ActivityStateChanged (message CHANGE added)"]
    n64["ActivityViewReturned (message CHANGE added)"]
    n65["ApplyActivityRecordSet (message CHANGE added)"]
    n66["ChangeActivityStateCommand (message CHANGE added)"]
    n67["CreateActivityCommand (message CHANGE added)"]
    n68["CreateSessionCommand (message CHANGE added)"]
    n69["ImplementationAdmitted (message CHANGE added)"]
    n7["SdpToolHost (container CONTEXT)"]
    n70["ImplementationRefused (message CHANGE added)"]
    n71["InspectActivityCommand (message CHANGE added)"]
    n72["RequestImplementationAdmission (message CHANGE added)"]
    n73["SessionCreated (message CHANGE added)"]
    n74["SessionCreationRefused (message CHANGE added)"]
    n75["ActivityManagement (mode CHANGE added)"]
    n76["CoordinatorCreatesSessionDocument (scenario CHANGE added)"]
    n77["OutOfScopeImplementationIsRefused (scenario CHANGE added)"]
    n78["OwnerChangesActivityState (scenario CHANGE added)"]
    n79["OwnerCreatesPlannedActivity (scenario CHANGE added)"]
    n8["ActivityCommandProtocol (contract CHANGE added)"]
    n80["StaleCreationIsRefused (scenario CHANGE added)"]
    n81["ViewerReadsActivityOverview (scenario CHANGE added)"]
    n82["SDPTool (system CONTEXT)"]
    n83["ActivityProjectionService (unit CHANGE added)"]
    n84["ActivityRecordRepository (unit CHANGE added)"]
    n85["ActivityService (unit CHANGE added)"]
    n86["ExternalCliCaller (unit CHANGE added)"]
    n87["ExternalCoordinatorCaller (unit CHANGE added)"]
    n88["ExternalViewer (unit CHANGE added)"]
    n89["AdmitScopedImplementation (usecase CHANGE added)"]
    n9["ActivityCommandShape (contract CHANGE added)"]
    n90["BrowseActivityDocuments (usecase CHANGE added)"]
    n91["CaptureActivityNeed (usecase CHANGE added)"]
    n92["ControlActivityLifecycle (usecase CHANGE added)"]
    n93["InspectActivityTimeline (usecase CHANGE added)"]
    n94["PrepareActivitySession (usecase CHANGE added)"]
    n4 -->|operand| f0
    n70 -->|operand| f0
    n75 -->|operand| f0
    n85 -->|operand| f0
    n11 -->|operand| f1
    n21 -->|operand| f1
    n42 -->|operand| f2
    n18 -->|operand| f3
    n84 -->|operand| f3
    n23 -->|operand| f4
    n54 -->|operand| f4
    n0 -->|operand| f5
    n45 -->|operand| f5
    n73 -->|operand| f6
    n17 -->|operand| f7
    n37 -->|operand| f7
    n17 -->|operand| f8
    n4 -->|operand| f9
    n72 -->|operand| f9
    n75 -->|operand| f9
    n87 -->|operand| f9
    n36 -->|operand| f10
    n16 -->|operand| f11
    n27 -->|operand| f11
    n33 -->|operand| f12
    n77 -->|operand| f13
    n18 -->|operand| f14
    n20 -->|operand| f14
    n3 -->|operand| f15
    n45 -->|operand| f15
    n23 -->|operand| f16
    n92 -->|operand| f16
    n7 -->|operand| f17
    n53 -->|operand| f17
    n75 -->|operand| f17
    n37 -->|operand| f18
    n8 -->|operand| f19
    n67 -->|operand| f19
    n4 -->|operand| f20
    n58 -->|operand| f20
    n75 -->|operand| f20
    n86 -->|operand| f20
    n5 -->|operand| f21
    n61 -->|operand| f21
    n80 -->|operand| f21
    n84 -->|operand| f21
    n85 -->|operand| f21
    n35 -->|operand| f22
    n3 -->|operand| f23
    n52 -->|operand| f23
    n4 -->|operand| f24
    n68 -->|operand| f24
    n75 -->|operand| f24
    n86 -->|operand| f24
    n23 -->|operand| f25
    n46 -->|operand| f25
    n80 -->|operand| f26
    n5 -->|operand| f27
    n60 -->|operand| f27
    n79 -->|operand| f27
    n84 -->|operand| f27
    n85 -->|operand| f27
    n38 -->|operand| f28
    n13 -->|operand| f29
    n64 -->|operand| f29
    n4 -->|operand| f30
    n58 -->|operand| f30
    n75 -->|operand| f30
    n87 -->|operand| f30
    n7 -->|operand| f31
    n50 -->|operand| f31
    n75 -->|operand| f31
    n36 -->|operand| f32
    n29 -->|operand| f33
    n13 -->|operand| f34
    n74 -->|operand| f34
    n4 -->|operand| f35
    n63 -->|operand| f35
    n75 -->|operand| f35
    n86 -->|operand| f35
    n76 -->|operand| f36
    n3 -->|operand| f37
    n48 -->|operand| f37
    n32 -->|operand| f38
    n8 -->|operand| f39
    n58 -->|operand| f39
    n14 -->|operand| f40
    n44 -->|operand| f40
    n61 -->|operand| f41
    n54 -->|operand| f42
    n84 -->|operand| f42
    n31 -->|operand| f43
    n9 -->|operand| f44
    n66 -->|operand| f44
    n51 -->|operand| f45
    n84 -->|operand| f45
    n9 -->|operand| f46
    n67 -->|operand| f46
    n40 -->|operand| f47
    n17 -->|operand| f48
    n34 -->|operand| f48
    n29 -->|operand| f49
    n74 -->|operand| f50
    n7 -->|operand| f51
    n85 -->|operand| f51
    n26 -->|operand| f52
    n6 -->|operand| f53
    n64 -->|operand| f53
    n75 -->|operand| f53
    n83 -->|operand| f53
    n44 -->|operand| f54
    n50 -->|operand| f55
    n7 -->|operand| f56
    n48 -->|operand| f56
    n75 -->|operand| f56
    n5 -->|operand| f57
    n65 -->|operand| f57
    n79 -->|operand| f57
    n84 -->|operand| f57
    n85 -->|operand| f57
    n54 -->|operand| f58
    n26 -->|operand| f59
    n5 -->|operand| f60
    n65 -->|operand| f60
    n75 -->|operand| f60
    n85 -->|operand| f60
    n8 -->|operand| f61
    n69 -->|operand| f61
    n5 -->|operand| f62
    n60 -->|operand| f62
    n75 -->|operand| f62
    n85 -->|operand| f62
    n69 -->|operand| f63
    n8 -->|operand| f64
    n73 -->|operand| f64
    n19 -->|operand| f65
    n84 -->|operand| f65
    n52 -->|operand| f66
    n84 -->|operand| f66
    n75 -->|operand| f67
    n79 -->|operand| f67
    n60 -->|operand| f68
    n23 -->|operand| f69
    n53 -->|operand| f69
    n7 -->|operand| f70
    n51 -->|operand| f70
    n75 -->|operand| f70
    n23 -->|operand| f71
    n51 -->|operand| f71
    n66 -->|operand| f72
    n70 -->|operand| f73
    n1 -->|operand| f74
    n94 -->|operand| f74
    n75 -->|operand| f75
    n80 -->|operand| f75
    n1 -->|operand| f76
    n93 -->|operand| f76
    n6 -->|operand| f77
    n71 -->|operand| f77
    n75 -->|operand| f77
    n88 -->|operand| f77
    n30 -->|operand| f78
    n64 -->|operand| f79
    n15 -->|operand| f80
    n7 -->|operand| f81
    n45 -->|operand| f81
    n75 -->|operand| f81
    n53 -->|operand| f82
    n13 -->|operand| f83
    n58 -->|operand| f83
    n10 -->|operand| f84
    n3 -->|operand| f85
    n54 -->|operand| f85
    n8 -->|operand| f86
    n70 -->|operand| f86
    n9 -->|operand| f87
    n68 -->|operand| f87
    n28 -->|operand| f88
    n65 -->|operand| f89
    n3 -->|operand| f90
    n7 -->|operand| f90
    n4 -->|operand| f91
    n73 -->|operand| f91
    n75 -->|operand| f91
    n86 -->|operand| f91
    n15 -->|operand| f92
    n71 -->|operand| f92
    n6 -->|operand| f93
    n71 -->|operand| f93
    n81 -->|operand| f93
    n83 -->|operand| f93
    n88 -->|operand| f93
    n4 -->|operand| f94
    n59 -->|operand| f94
    n80 -->|operand| f94
    n85 -->|operand| f94
    n86 -->|operand| f94
    n3 -->|operand| f95
    n53 -->|operand| f95
    n4 -->|operand| f96
    n58 -->|operand| f96
    n79 -->|operand| f96
    n85 -->|operand| f96
    n86 -->|operand| f96
    n45 -->|operand| f97
    n85 -->|operand| f97
    n11 -->|operand| f98
    n30 -->|operand| f98
    n37 -->|operand| f99
    n4 -->|operand| f100
    n59 -->|operand| f100
    n75 -->|operand| f100
    n86 -->|operand| f100
    n9 -->|operand| f101
    n25 -->|operand| f101
    n53 -->|operand| f102
    n85 -->|operand| f102
    n81 -->|operand| f103
    n48 -->|operand| f104
    n25 -->|operand| f105
    n13 -->|operand| f106
    n73 -->|operand| f106
    n8 -->|operand| f107
    n4 -->|operand| f108
    n66 -->|operand| f108
    n75 -->|operand| f108
    n85 -->|operand| f108
    n77 -->|operand| f109
    n89 -->|operand| f109
    n23 -->|operand| f110
    n48 -->|operand| f110
    n81 -->|operand| f111
    n90 -->|operand| f111
    n4 -->|operand| f112
    n59 -->|operand| f112
    n75 -->|operand| f112
    n87 -->|operand| f112
    n45 -->|operand| f113
    n13 -->|operand| f114
    n39 -->|operand| f114
    n7 -->|operand| f115
    n84 -->|operand| f115
    n24 -->|operand| f116
    n90 -->|operand| f116
    n38 -->|operand| f117
    n9 -->|operand| f118
    n71 -->|operand| f118
    n2 -->|operand| f119
    n7 -->|operand| f119
    n51 -->|operand| f120
    n71 -->|operand| f121
    n4 -->|operand| f122
    n74 -->|operand| f122
    n75 -->|operand| f122
    n85 -->|operand| f122
    n41 -->|operand| f123
    n2 -->|operand| f124
    n57 -->|operand| f124
    n75 -->|operand| f124
    n0 -->|operand| f125
    n23 -->|operand| f125
    n4 -->|operand| f126
    n68 -->|operand| f126
    n75 -->|operand| f126
    n85 -->|operand| f126
    n0 -->|operand| f127
    n47 -->|operand| f127
    n3 -->|operand| f128
    n55 -->|operand| f128
    n13 -->|operand| f129
    n69 -->|operand| f129
    n50 -->|operand| f130
    n83 -->|operand| f130
    n10 -->|operand| f131
    n20 -->|operand| f131
    n59 -->|operand| f132
    n8 -->|operand| f133
    n66 -->|operand| f133
    n6 -->|operand| f134
    n15 -->|operand| f134
    n4 -->|operand| f135
    n59 -->|operand| f135
    n75 -->|operand| f135
    n85 -->|operand| f135
    n14 -->|operand| f136
    n9 -->|operand| f137
    n33 -->|operand| f137
    n24 -->|operand| f138
    n93 -->|operand| f138
    n11 -->|operand| f139
    n5 -->|operand| f140
    n65 -->|operand| f140
    n76 -->|operand| f140
    n84 -->|operand| f140
    n85 -->|operand| f140
    n4 -->|operand| f141
    n72 -->|operand| f141
    n75 -->|operand| f141
    n85 -->|operand| f141
    n4 -->|operand| f142
    n67 -->|operand| f142
    n75 -->|operand| f142
    n86 -->|operand| f142
    n31 -->|operand| f143
    n75 -->|operand| f144
    n78 -->|operand| f144
    n0 -->|operand| f145
    n24 -->|operand| f145
    n7 -->|operand| f146
    n83 -->|operand| f146
    n10 -->|operand| f147
    n35 -->|operand| f147
    n18 -->|operand| f148
    n22 -->|operand| f148
    n8 -->|operand| f149
    n62 -->|operand| f149
    n60 -->|operand| f150
    n65 -->|operand| f150
    n70 -->|operand| f151
    n72 -->|operand| f151
    n0 -->|operand| f152
    n55 -->|operand| f152
    n32 -->|operand| f153
    n16 -->|operand| f154
    n29 -->|operand| f154
    n48 -->|operand| f155
    n85 -->|operand| f155
    n4 -->|operand| f156
    n73 -->|operand| f156
    n76 -->|operand| f156
    n85 -->|operand| f156
    n87 -->|operand| f156
    n8 -->|operand| f157
    n68 -->|operand| f157
    n5 -->|operand| f158
    n60 -->|operand| f158
    n75 -->|operand| f158
    n84 -->|operand| f158
    n17 -->|operand| f159
    n65 -->|operand| f159
    n58 -->|operand| f160
    n67 -->|operand| f160
    n7 -->|operand| f161
    n55 -->|operand| f161
    n75 -->|operand| f161
    n80 -->|operand| f162
    n91 -->|operand| f162
    n5 -->|operand| f163
    n65 -->|operand| f163
    n80 -->|operand| f163
    n84 -->|operand| f163
    n85 -->|operand| f163
    n52 -->|operand| f164
    n4 -->|operand| f165
    n70 -->|operand| f165
    n75 -->|operand| f165
    n86 -->|operand| f165
    n23 -->|operand| f166
    n94 -->|operand| f166
    n49 -->|operand| f167
    n83 -->|operand| f167
    n5 -->|operand| f168
    n61 -->|operand| f168
    n75 -->|operand| f168
    n85 -->|operand| f168
    n46 -->|operand| f169
    n85 -->|operand| f169
    n4 -->|operand| f170
    n67 -->|operand| f170
    n79 -->|operand| f170
    n85 -->|operand| f170
    n86 -->|operand| f170
    n79 -->|operand| f171
    n43 -->|operand| f172
    n13 -->|operand| f173
    n63 -->|operand| f173
    n5 -->|operand| f174
    n65 -->|operand| f174
    n75 -->|operand| f174
    n84 -->|operand| f174
    n0 -->|operand| f175
    n48 -->|operand| f175
    n62 -->|operand| f176
    n66 -->|operand| f176
    n4 -->|operand| f177
    n70 -->|operand| f177
    n77 -->|operand| f177
    n85 -->|operand| f177
    n87 -->|operand| f177
    n9 -->|operand| f178
    n72 -->|operand| f178
    n14 -->|operand| f179
    n22 -->|operand| f179
    n4 -->|operand| f180
    n69 -->|operand| f180
    n75 -->|operand| f180
    n86 -->|operand| f180
    n4 -->|operand| f181
    n69 -->|operand| f181
    n75 -->|operand| f181
    n85 -->|operand| f181
    n75 -->|operand| f182
    n76 -->|operand| f182
    n25 -->|operand| f183
    n12 -->|operand| f184
    n61 -->|operand| f184
    n4 -->|operand| f185
    n66 -->|operand| f185
    n75 -->|operand| f185
    n86 -->|operand| f185
    n19 -->|operand| f186
    n21 -->|operand| f186
    n4 -->|operand| f187
    n62 -->|operand| f187
    n75 -->|operand| f187
    n87 -->|operand| f187
    n4 -->|operand| f188
    n63 -->|operand| f188
    n75 -->|operand| f188
    n85 -->|operand| f188
    n11 -->|operand| f189
    n32 -->|operand| f189
    n9 -->|operand| f190
    n4 -->|operand| f191
    n66 -->|operand| f191
    n78 -->|operand| f191
    n85 -->|operand| f191
    n86 -->|operand| f191
    n17 -->|operand| f192
    n28 -->|operand| f192
    n7 -->|operand| f193
    n46 -->|operand| f193
    n75 -->|operand| f193
    n1 -->|operand| f194
    n91 -->|operand| f194
    n4 -->|operand| f195
    n8 -->|operand| f195
    n10 -->|operand| f196
    n26 -->|operand| f196
    n4 -->|operand| f197
    n68 -->|operand| f197
    n75 -->|operand| f197
    n87 -->|operand| f197
    n34 -->|operand| f198
    n0 -->|operand| f199
    n50 -->|operand| f199
    n0 -->|operand| f200
    n54 -->|operand| f200
    n23 -->|operand| f201
    n89 -->|operand| f201
    n8 -->|operand| f202
    n63 -->|operand| f202
    n47 -->|operand| f203
    n85 -->|operand| f203
    n4 -->|operand| f204
    n67 -->|operand| f204
    n75 -->|operand| f204
    n87 -->|operand| f204
    n67 -->|operand| f205
    n2 -->|operand| f206
    n49 -->|operand| f206
    n3 -->|operand| f207
    n51 -->|operand| f207
    n6 -->|operand| f208
    n71 -->|operand| f208
    n75 -->|operand| f208
    n86 -->|operand| f208
    n7 -->|operand| f209
    n54 -->|operand| f209
    n75 -->|operand| f209
    n33 -->|operand| f210
    n42 -->|operand| f211
    n5 -->|operand| f212
    n60 -->|operand| f212
    n78 -->|operand| f212
    n84 -->|operand| f212
    n85 -->|operand| f212
    n39 -->|operand| f213
    n4 -->|operand| f214
    n74 -->|operand| f214
    n75 -->|operand| f214
    n86 -->|operand| f214
    n57 -->|operand| f215
    n84 -->|operand| f215
    n14 -->|operand| f216
    n42 -->|operand| f216
    n68 -->|operand| f217
    n73 -->|operand| f217
    n23 -->|operand| f218
    n91 -->|operand| f218
    n43 -->|operand| f219
    n15 -->|operand| f220
    n64 -->|operand| f220
    n4 -->|operand| f221
    n63 -->|operand| f221
    n78 -->|operand| f221
    n85 -->|operand| f221
    n86 -->|operand| f221
    n23 -->|operand| f222
    n47 -->|operand| f222
    n6 -->|operand| f223
    n71 -->|operand| f223
    n75 -->|operand| f223
    n83 -->|operand| f223
    n13 -->|operand| f224
    n40 -->|operand| f224
    n1 -->|operand| f225
    n90 -->|operand| f225
    n1 -->|operand| f226
    n92 -->|operand| f226
    n4 -->|operand| f227
    n72 -->|operand| f227
    n77 -->|operand| f227
    n85 -->|operand| f227
    n87 -->|operand| f227
    n44 -->|operand| f228
    n13 -->|operand| f229
    n59 -->|operand| f229
    n6 -->|operand| f230
    n64 -->|operand| f230
    n81 -->|operand| f230
    n83 -->|operand| f230
    n88 -->|operand| f230
    n39 -->|operand| f231
    n64 -->|operand| f232
    n71 -->|operand| f232
    n5 -->|operand| f233
    n60 -->|operand| f233
    n76 -->|operand| f233
    n84 -->|operand| f233
    n85 -->|operand| f233
    n28 -->|operand| f234
    n13 -->|operand| f235
    n70 -->|operand| f235
    n27 -->|operand| f236
    n78 -->|operand| f237
    n3 -->|operand| f238
    n57 -->|operand| f238
    n75 -->|operand| f238
    n34 -->|operand| f239
    n16 -->|operand| f240
    n60 -->|operand| f240
    n30 -->|operand| f241
    n23 -->|operand| f242
    n45 -->|operand| f242
    n14 -->|operand| f243
    n43 -->|operand| f243
    n61 -->|operand| f244
    n65 -->|operand| f244
    n7 -->|operand| f245
    n82 -->|operand| f245
    n12 -->|operand| f246
    n13 -->|operand| f247
    n7 -->|operand| f248
    n52 -->|operand| f248
    n75 -->|operand| f248
    n0 -->|operand| f249
    n51 -->|operand| f249
    n4 -->|operand| f250
    n68 -->|operand| f250
    n76 -->|operand| f250
    n85 -->|operand| f250
    n87 -->|operand| f250
    n4 -->|operand| f251
    n63 -->|operand| f251
    n75 -->|operand| f251
    n87 -->|operand| f251
    n12 -->|operand| f252
    n65 -->|operand| f252
    n56 -->|operand| f253
    n85 -->|operand| f253
    n7 -->|operand| f254
    n49 -->|operand| f254
    n75 -->|operand| f254
    n23 -->|operand| f255
    n52 -->|operand| f255
    n47 -->|operand| f256
    n55 -->|operand| f257
    n4 -->|operand| f258
    n58 -->|operand| f258
    n75 -->|operand| f258
    n85 -->|operand| f258
    n75 -->|operand| f259
    n81 -->|operand| f259
    n63 -->|operand| f260
    n0 -->|operand| f261
    n53 -->|operand| f261
    n9 -->|operand| f262
    n41 -->|operand| f262
    n68 -->|operand| f263
    n74 -->|operand| f263
    n5 -->|operand| f264
    n61 -->|operand| f264
    n75 -->|operand| f264
    n84 -->|operand| f264
    n0 -->|operand| f265
    n52 -->|operand| f265
    n0 -->|operand| f266
    n49 -->|operand| f266
    n68 -->|operand| f267
    n4 -->|operand| f268
    n73 -->|operand| f268
    n75 -->|operand| f268
    n85 -->|operand| f268
    n2 -->|operand| f269
    n50 -->|operand| f269
    n4 -->|operand| f270
    n62 -->|operand| f270
    n75 -->|operand| f270
    n86 -->|operand| f270
    n6 -->|operand| f271
    n64 -->|operand| f271
    n75 -->|operand| f271
    n86 -->|operand| f271
    n62 -->|operand| f272
    n3 -->|operand| f273
    n56 -->|operand| f273
    n75 -->|operand| f273
    n4 -->|operand| f274
    n73 -->|operand| f274
    n75 -->|operand| f274
    n87 -->|operand| f274
    n9 -->|operand| f275
    n36 -->|operand| f275
    n69 -->|operand| f276
    n72 -->|operand| f276
    n40 -->|operand| f277
    n55 -->|operand| f278
    n85 -->|operand| f278
    n59 -->|operand| f279
    n67 -->|operand| f279
    n4 -->|operand| f280
    n66 -->|operand| f280
    n75 -->|operand| f280
    n87 -->|operand| f280
    n13 -->|operand| f281
    n38 -->|operand| f281
    n4 -->|operand| f282
    n67 -->|operand| f282
    n80 -->|operand| f282
    n85 -->|operand| f282
    n86 -->|operand| f282
    n58 -->|operand| f283
    n41 -->|operand| f284
    n8 -->|operand| f285
    n72 -->|operand| f285
    n4 -->|operand| f286
    n69 -->|operand| f286
    n75 -->|operand| f286
    n87 -->|operand| f286
    n5 -->|operand| f287
    n65 -->|operand| f287
    n78 -->|operand| f287
    n84 -->|operand| f287
    n85 -->|operand| f287
    n3 -->|operand| f288
    n47 -->|operand| f288
    n11 -->|operand| f289
    n31 -->|operand| f289
    n49 -->|operand| f290
    n57 -->|operand| f291
    n83 -->|operand| f291
    n24 -->|operand| f292
    n50 -->|operand| f292
    n27 -->|operand| f293
    n8 -->|operand| f294
    n74 -->|operand| f294
    n76 -->|operand| f295
    n94 -->|operand| f295
    n3 -->|operand| f296
    n46 -->|operand| f296
    n35 -->|operand| f297
    n78 -->|operand| f298
    n92 -->|operand| f298
    n8 -->|operand| f299
    n59 -->|operand| f299
    n4 -->|operand| f300
    n62 -->|operand| f300
    n75 -->|operand| f300
    n85 -->|operand| f300
    n6 -->|operand| f301
    n64 -->|operand| f301
    n75 -->|operand| f301
    n88 -->|operand| f301
    n5 -->|operand| f302
    n12 -->|operand| f302
    n16 -->|operand| f303
    n61 -->|operand| f303
    n72 -->|operand| f304
    n79 -->|operand| f305
    n91 -->|operand| f305
    n23 -->|operand| f306
    n55 -->|operand| f306
    n12 -->|operand| f307
    n60 -->|operand| f307
    n13 -->|operand| f308
    n62 -->|operand| f308
    n7 -->|operand| f309
    n47 -->|operand| f309
    n75 -->|operand| f309
    n63 -->|operand| f310
    n66 -->|operand| f310
    n0 -->|operand| f311
    n46 -->|operand| f311
    n16 -->|operand| f312
    n0 -->|operand| f313
    n4 -->|operand| f314
    n72 -->|operand| f314
    n75 -->|operand| f314
    n86 -->|operand| f314
    n1 -->|operand| f315
    n89 -->|operand| f315
    n24 -->|operand| f316
    n49 -->|operand| f316
    n4 -->|operand| f317
    n74 -->|operand| f317
    n75 -->|operand| f317
    n87 -->|operand| f317
    n46 -->|operand| f318
    n4 -->|operand| f319
    n70 -->|operand| f319
    n75 -->|operand| f319
    n87 -->|operand| f319
    n4 -->|operand| f320
    n67 -->|operand| f320
    n75 -->|operand| f320
    n85 -->|operand| f320
    n75 -->|operand| f321
    n77 -->|operand| f321
```

## Cuts and exclusions

