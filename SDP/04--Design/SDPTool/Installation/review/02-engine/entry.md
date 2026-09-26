# VP02 — Architecture and logical decomposition

Revision: `c05cf8fcba4a1c7abb985b1a3814bf9690759dfa74d12cdf58d20a9c46d96330`.

Container/Unit and contains; library structure is not deployment allocation.

Selection: `{"viewpoint":"VP02","direction":"both","depth":2,"diagram":"VP02-SdpToolProcess"}`.

## Logical decomposition: SdpToolProcess

```mermaid
flowchart LR
    n_CardHistory["CardHistory (unit)"]
    n_InstallationBaselineInspector["InstallationBaselineInspector (unit)"]
    n_InstallationCoordinator["InstallationCoordinator (unit)"]
    n_InstallationExecutor["InstallationExecutor (unit)"]
    n_InstallationJournal["InstallationJournal (unit)"]
    n_InstallationPlanner["InstallationPlanner (unit)"]
    n_InstallationRecorder["InstallationRecorder (unit)"]
    n_InstallationReleaseResolver["InstallationReleaseResolver (unit)"]
    n_NavigationInventory["NavigationInventory (unit)"]
    n_PlanCoordinator["PlanCoordinator (unit)"]
    n_PreviewCoordinator["PreviewCoordinator (unit)"]
    n_ProjectContext["ProjectContext (unit)"]
    n_SdpToolProcess["SdpToolProcess (container)"]
    n_ViewerBridge["ViewerBridge (unit)"]
    n_SdpToolProcess -->|contains| n_CardHistory
    n_SdpToolProcess -->|contains| n_InstallationBaselineInspector
    n_SdpToolProcess -->|contains| n_InstallationCoordinator
    n_SdpToolProcess -->|contains| n_InstallationExecutor
    n_SdpToolProcess -->|contains| n_InstallationJournal
    n_SdpToolProcess -->|contains| n_InstallationPlanner
    n_SdpToolProcess -->|contains| n_InstallationRecorder
    n_SdpToolProcess -->|contains| n_InstallationReleaseResolver
    n_SdpToolProcess -->|contains| n_NavigationInventory
    n_SdpToolProcess -->|contains| n_PlanCoordinator
    n_SdpToolProcess -->|contains| n_PreviewCoordinator
    n_SdpToolProcess -->|contains| n_ProjectContext
    n_SdpToolProcess -->|contains| n_ViewerBridge
```

Source facts: f0526, f0527, f0528, f0529, f0530, f0531, f0532, f0533, f0534, f0535, f0536, f0537, f0538.


