# Feature: GovernedProjectOperations — mode RoutineHandling

[Viewpoint](index.md) · [Navigator](../../navigator.md)

Revision: `02ca3d1cf871998399f43803bbf47f6fa73321b38ce22301b196973963cd3480`.

## Feature: GovernedProjectOperations — mode RoutineHandling

```mermaid
flowchart LR
    n_CheckTransitionAuthority["CheckTransitionAuthority (functionality)"]
    n_ContextFacade["ContextFacade (unit)"]
    n_GovernedProjectOperations["GovernedProjectOperations (feature)"]
    n_ReconcileRunRevision["ReconcileRunRevision (functionality)"]
    n_ResolveApplicableRoutine["ResolveApplicableRoutine (functionality)"]
    n_ResolveCheckoutAuthority["ResolveCheckoutAuthority (functionality)"]
    n_RoutineCoordinator["RoutineCoordinator (unit)"]
    n_RunRepository["RunRepository (unit)"]
    n_SdpToolHost["SdpToolHost (container)"]
    n_CheckTransitionAuthority -->|allocated-to| n_SdpToolHost
    n_CheckTransitionAuthority -->|contributes-to| n_GovernedProjectOperations
    n_ContextFacade -->|owns| n_ResolveCheckoutAuthority
    n_ReconcileRunRevision -->|allocated-to| n_SdpToolHost
    n_ReconcileRunRevision -->|contributes-to| n_GovernedProjectOperations
    n_ResolveApplicableRoutine -->|allocated-to| n_SdpToolHost
    n_ResolveApplicableRoutine -->|contributes-to| n_GovernedProjectOperations
    n_ResolveCheckoutAuthority -->|allocated-to| n_SdpToolHost
    n_ResolveCheckoutAuthority -->|contributes-to| n_GovernedProjectOperations
    n_RoutineCoordinator -->|owns| n_CheckTransitionAuthority
    n_RoutineCoordinator -->|owns| n_ResolveApplicableRoutine
    n_RunRepository -->|owns| n_ReconcileRunRevision
```

Source facts: f0037, f0038, f0041, f0086, f0087, f0091, f0092, f0094, f0095, f0113, f0114, f0166.

