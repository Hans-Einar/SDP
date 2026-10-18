# Feature: BlueprintAssignments — mode RoutineHandling

[Viewpoint](index.md) · [Navigator](../../navigator.md)

Revision: `02ca3d1cf871998399f43803bbf47f6fa73321b38ce22301b196973963cd3480`.

## Feature: BlueprintAssignments — mode RoutineHandling

```mermaid
flowchart LR
    n_BlueprintAssignments["BlueprintAssignments (feature)"]
    n_BlueprintLifecycle["BlueprintLifecycle (unit)"]
    n_GovernBlueprintAssignment["GovernBlueprintAssignment (functionality)"]
    n_SdpToolHost["SdpToolHost (container)"]
    n_BlueprintLifecycle -->|owns| n_GovernBlueprintAssignment
    n_GovernBlueprintAssignment -->|allocated-to| n_SdpToolHost
    n_GovernBlueprintAssignment -->|contributes-to| n_BlueprintAssignments
```

Source facts: f0026, f0067, f0068.

