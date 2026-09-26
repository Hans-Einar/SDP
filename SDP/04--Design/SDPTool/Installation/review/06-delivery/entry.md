# VP06 — Activities and delivery plan

Revision: `c62878ccf713fe1fd51d9a7d14e90dd7f8a70dd5deb346af2d8fb9edae9969d2`.

refines, addresses, delivers and depends-on; status is an explicit source claim.

Selection: `{"viewpoint":"VP06","direction":"both","depth":2,"diagram":"VP06-dependencies"}`.

## Explicit activity dependencies

```mermaid
flowchart LR
    n_ConsumerValidationSlice["ConsumerValidationSlice (activity)"]
    n_GhSdpClientSlice["GhSdpClientSlice (activity)"]
    n_HistorySlice["HistorySlice (activity)"]
    n_InstallationApplySlice["InstallationApplySlice (activity)"]
    n_InstallationPlanSlice["InstallationPlanSlice (activity)"]
    n_InstallationRecoverySlice["InstallationRecoverySlice (activity)"]
    n_ManualAdoptionTrial["ManualAdoptionTrial (activity)"]
    n_NavigationSlice["NavigationSlice (activity)"]
    n_PlanAssistanceSlice["PlanAssistanceSlice (activity)"]
    n_ProjectDiscoverySlice["ProjectDiscoverySlice (activity)"]
    n_SavedFilePreview["SavedFilePreview (activity)"]
    n_UnsavedSourcePreview["UnsavedSourcePreview (activity)"]
    n_ConsumerValidationSlice -->|depends-on| n_NavigationSlice
    n_GhSdpClientSlice -->|depends-on| n_InstallationRecoverySlice
    n_HistorySlice -->|depends-on| n_ProjectDiscoverySlice
    n_InstallationApplySlice -->|depends-on| n_InstallationPlanSlice
    n_InstallationPlanSlice -->|depends-on| n_ProjectDiscoverySlice
    n_InstallationRecoverySlice -->|depends-on| n_InstallationApplySlice
    n_ManualAdoptionTrial -->|depends-on| n_GhSdpClientSlice
    n_NavigationSlice -->|depends-on| n_ProjectDiscoverySlice
    n_PlanAssistanceSlice -->|depends-on| n_ProjectDiscoverySlice
    n_UnsavedSourcePreview -->|depends-on| n_SavedFilePreview
```

Source facts: f0034, f0085, f0108, f0226, f0314, f0346, f0384, f0394, f0406, f0547.


