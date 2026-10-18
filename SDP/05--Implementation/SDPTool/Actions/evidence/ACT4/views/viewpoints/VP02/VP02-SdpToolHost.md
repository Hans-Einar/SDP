# Logical decomposition: SdpToolHost

[Viewpoint](index.md) · [Navigator](../../navigator.md)

Revision: `02ca3d1cf871998399f43803bbf47f6fa73321b38ce22301b196973963cd3480`.

## Logical decomposition: SdpToolHost

```mermaid
flowchart LR
    n_ActionRegistry["ActionRegistry (unit)"]
    n_BlueprintAssessment["BlueprintAssessment (unit)"]
    n_BlueprintCatalogue["BlueprintCatalogue (unit)"]
    n_BlueprintFacade["BlueprintFacade (unit)"]
    n_BlueprintLifecycle["BlueprintLifecycle (unit)"]
    n_ContextFacade["ContextFacade (unit)"]
    n_ExistingInstallationFacade["ExistingInstallationFacade (unit)"]
    n_ExistingNavigationFacade["ExistingNavigationFacade (unit)"]
    n_JsonInvocationAdapter["JsonInvocationAdapter (unit)"]
    n_RoutineCoordinator["RoutineCoordinator (unit)"]
    n_RunRepository["RunRepository (unit)"]
    n_SdpToolHost["SdpToolHost (container)"]
    n_SdpToolHost -->|contains| n_ActionRegistry
    n_SdpToolHost -->|contains| n_BlueprintAssessment
    n_SdpToolHost -->|contains| n_BlueprintCatalogue
    n_SdpToolHost -->|contains| n_BlueprintFacade
    n_SdpToolHost -->|contains| n_BlueprintLifecycle
    n_SdpToolHost -->|contains| n_ContextFacade
    n_SdpToolHost -->|contains| n_ExistingInstallationFacade
    n_SdpToolHost -->|contains| n_ExistingNavigationFacade
    n_SdpToolHost -->|contains| n_JsonInvocationAdapter
    n_SdpToolHost -->|contains| n_RoutineCoordinator
    n_SdpToolHost -->|contains| n_RunRepository
```

Source facts: f0169, f0170, f0171, f0172, f0173, f0174, f0175, f0176, f0177, f0178, f0179.

