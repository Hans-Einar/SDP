# VP09 — Dataset, Datagram and persistent Database

Revision: `c05cf8fcba4a1c7abb985b1a3814bf9690759dfa74d12cdf58d20a9c46d96330`.

Explicit holders, sources, contracts, variants, fields and projections.

Selection: `{"viewpoint":"VP09","direction":"both","depth":2,"diagram":"VP09-data-InstalledProcessRecords"}`.

## Data origin and holder: InstalledProcessRecords

```mermaid
flowchart LR
    n_InstallReceiptRecord["InstallReceiptRecord (contract)"]
    n_InstallationRecorder["InstallationRecorder (unit)"]
    n_InstallationStore["InstallationStore (database)"]
    n_InstalledProcessRecords["InstalledProcessRecords (dataset)"]
    n_InstallationRecorder -->|owns| n_InstallationStore
    n_InstallationStore -->|holds| n_InstalledProcessRecords
    n_InstalledProcessRecords -->|upholds| n_InstallReceiptRecord
```

Source facts: f0332, f0374, f0375.

## Fields and contract properties

Database means a persistent data source, not necessarily SQL.

| Fact | ID |
| --- | --- |
| BaselineInspectionCallsProtocol has completeness = closed. | f0015 |
| InstallApplyArguments has completeness = closed. | f0119 |
| InstallApplyArgumentsPlanDigest has presence = required. | f0122 |
| InstallApplyArgumentsPlanDigest has value-type = text. | f0123 |
| InstallApplyArgumentsProjectRoot has presence = required. | f0124 |
| InstallApplyArgumentsProjectRoot has value-type = text. | f0125 |
| InstallJournalRecord has completeness = closed. | f0126 |
| InstallJournalRecordOperationIdentity has presence = required. | f0131 |
| InstallJournalRecordOperationIdentity has value-type = text. | f0132 |
| InstallJournalRecordPlanDigest has presence = required. | f0133 |
| InstallJournalRecordPlanDigest has value-type = text. | f0134 |
| InstallJournalRecordRecordBytes has presence = required. | f0135 |
| InstallJournalRecordRecordBytes has value-type = bytes. | f0136 |
| InstallJournalRecordResultStatus has presence = required. | f0137 |
| InstallJournalRecordResultStatus has value-type = text. | f0138 |
| InstallOutcome has completeness = closed. | f0139 |
| InstallOutcomeRecordBytes has presence = required. | f0142 |
| InstallOutcomeRecordBytes has value-type = bytes. | f0143 |
| InstallOutcomeResultStatus has presence = required. | f0144 |
| InstallOutcomeResultStatus has value-type = text. | f0145 |
| InstallPlanRecord has completeness = closed. | f0146 |
| InstallPlanRecordPlanDigest has presence = required. | f0150 |
| InstallPlanRecordPlanDigest has value-type = text. | f0151 |
| InstallPlanRecordProjectRoot has presence = required. | f0152 |
| InstallPlanRecordProjectRoot has value-type = text. | f0153 |
| InstallPlanRecordRecordBytes has presence = required. | f0154 |
| InstallPlanRecordRecordBytes has value-type = bytes. | f0155 |
| InstallPlanningArguments has completeness = closed. | f0156 |
| InstallPlanningArgumentsAdoptionDigest has presence = required. | f0161 |
| InstallPlanningArgumentsAdoptionDigest has value-type = text. | f0162 |
| InstallPlanningArgumentsArtifactDigest has presence = required. | f0163 |
| InstallPlanningArgumentsArtifactDigest has value-type = text. | f0164 |
| InstallPlanningArgumentsProjectRoot has presence = required. | f0165 |
| InstallPlanningArgumentsProjectRoot has value-type = text. | f0166 |
| InstallPlanningArgumentsReleaseIdentity has presence = required. | f0167 |
| InstallPlanningArgumentsReleaseIdentity has value-type = text. | f0168 |
| InstallReceiptRecord has completeness = closed. | f0169 |
| InstallReceiptRecordArtifactDigest has presence = required. | f0174 |
| InstallReceiptRecordArtifactDigest has value-type = text. | f0175 |
| InstallReceiptRecordRecordBytes has presence = required. | f0176 |
| InstallReceiptRecordRecordBytes has value-type = bytes. | f0177 |
| InstallReceiptRecordReleaseIdentity has presence = required. | f0178 |
| InstallReceiptRecordReleaseIdentity has value-type = text. | f0179 |
| InstallReceiptRecordSchemaIdentity has presence = required. | f0180 |
| InstallReceiptRecordSchemaIdentity has value-type = text. | f0181 |
| InstallReleaseArguments has completeness = closed. | f0182 |
| InstallReleaseArgumentsReleaseIdentity has presence = required. | f0184 |
| InstallReleaseArgumentsReleaseIdentity has value-type = text. | f0185 |
| InstallReleaseRecord has completeness = closed. | f0186 |
| InstallReleaseRecordArtifactDigest has presence = required. | f0190 |
| InstallReleaseRecordArtifactDigest has value-type = text. | f0191 |
| InstallReleaseRecordRecordBytes has presence = required. | f0192 |
| InstallReleaseRecordRecordBytes has value-type = bytes. | f0193 |
| InstallReleaseRecordReleaseIdentity has presence = required. | f0194 |
| InstallReleaseRecordReleaseIdentity has value-type = text. | f0195 |
| InstallResumeArguments has completeness = closed. | f0196 |
| InstallResumeArgumentsOperationIdentity has presence = required. | f0199 |
| InstallResumeArgumentsOperationIdentity has value-type = text. | f0200 |
| InstallResumeArgumentsProjectRoot has presence = required. | f0201 |
| InstallResumeArgumentsProjectRoot has value-type = text. | f0202 |
| InstallationApplyCallsProtocol has completeness = closed. | f0218 |
| InstallationExecutionCallsProtocol has completeness = closed. | f0261 |
| InstallationJournalCallsProtocol has completeness = closed. | f0287 |
| InstallationPlanCallsProtocol has completeness = closed. | f0292 |
| InstallationPlanningCallsProtocol has completeness = closed. | f0323 |
| InstallationReceiptCallsProtocol has completeness = closed. | f0328 |
| InstallationRecoveryCallsProtocol has completeness = closed. | f0340 |
| InstallationResumeCallsProtocol has completeness = closed. | f0356 |
| ReleaseArtifactCallsProtocol has completeness = closed. | f0458 |
| ReleaseResolutionCallsProtocol has completeness = closed. | f0470 |
| SdpToolBootstrapCallsProtocol has completeness = closed. | f0516 |

## Projection responsibilities

| Functionality | Dataset | Datagram | Fact |
| --- | --- | --- | --- |

