# Feature: DiscoverableActions — mode RoutineHandling

[Viewpoint](index.md) · [Navigator](../../navigator.md)

Revision: `02ca3d1cf871998399f43803bbf47f6fa73321b38ce22301b196973963cd3480`.

## Feature: DiscoverableActions — mode RoutineHandling

```mermaid
flowchart LR
    n_ActionRegistry["ActionRegistry (unit)"]
    n_DescribeSupportedActions["DescribeSupportedActions (functionality)"]
    n_DiscoverableActions["DiscoverableActions (feature)"]
    n_InvokeTypedAction["InvokeTypedAction (functionality)"]
    n_JsonInvocationAdapter["JsonInvocationAdapter (unit)"]
    n_SdpToolHost["SdpToolHost (container)"]
    n_ActionRegistry -->|owns| n_DescribeSupportedActions
    n_DescribeSupportedActions -->|contributes-to| n_DiscoverableActions
    n_InvokeTypedAction -->|allocated-to| n_SdpToolHost
    n_InvokeTypedAction -->|contributes-to| n_DiscoverableActions
    n_JsonInvocationAdapter -->|owns| n_InvokeTypedAction
```

Source facts: f0000, f0046, f0075, f0076, f0078.

