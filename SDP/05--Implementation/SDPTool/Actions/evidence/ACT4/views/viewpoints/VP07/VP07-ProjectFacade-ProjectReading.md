# Feature: ProjectFacade — mode ProjectReading

[Viewpoint](index.md) · [Navigator](../../navigator.md)

Revision: `02ca3d1cf871998399f43803bbf47f6fa73321b38ce22301b196973963cd3480`.

## Feature: ProjectFacade — mode ProjectReading

```mermaid
flowchart LR
    n_DescribeRegisteredProject["DescribeRegisteredProject (functionality)"]
    n_ExistingNavigationFacade["ExistingNavigationFacade (unit)"]
    n_ProjectFacade["ProjectFacade (feature)"]
    n_SdpToolHost["SdpToolHost (container)"]
    n_DescribeRegisteredProject -->|allocated-to| n_SdpToolHost
    n_DescribeRegisteredProject -->|contributes-to| n_ProjectFacade
    n_ExistingNavigationFacade -->|owns| n_DescribeRegisteredProject
```

Source facts: f0042, f0043, f0057.

