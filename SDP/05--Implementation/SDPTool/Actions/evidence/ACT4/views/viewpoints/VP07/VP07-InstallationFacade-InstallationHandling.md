# Feature: InstallationFacade — mode InstallationHandling

[Viewpoint](index.md) · [Navigator](../../navigator.md)

Revision: `02ca3d1cf871998399f43803bbf47f6fa73321b38ce22301b196973963cd3480`.

## Feature: InstallationFacade — mode InstallationHandling

```mermaid
flowchart LR
    n_ExistingInstallationFacade["ExistingInstallationFacade (unit)"]
    n_InstallationFacade["InstallationFacade (feature)"]
    n_InvokeReviewedInstallation["InvokeReviewedInstallation (functionality)"]
    n_SdpToolHost["SdpToolHost (container)"]
    n_ExistingInstallationFacade -->|owns| n_InvokeReviewedInstallation
    n_InvokeReviewedInstallation -->|allocated-to| n_SdpToolHost
    n_InvokeReviewedInstallation -->|contributes-to| n_InstallationFacade
```

Source facts: f0055, f0072, f0073.

