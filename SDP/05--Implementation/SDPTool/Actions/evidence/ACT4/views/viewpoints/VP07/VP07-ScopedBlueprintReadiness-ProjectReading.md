# Feature: ScopedBlueprintReadiness — mode ProjectReading

[Viewpoint](index.md) · [Navigator](../../navigator.md)

Revision: `02ca3d1cf871998399f43803bbf47f6fa73321b38ce22301b196973963cd3480`.

## Feature: ScopedBlueprintReadiness — mode ProjectReading

```mermaid
flowchart LR
    n_AssessBlueprintEvidence["AssessBlueprintEvidence (functionality)"]
    n_BlueprintAssessment["BlueprintAssessment (unit)"]
    n_ScopedBlueprintReadiness["ScopedBlueprintReadiness (feature)"]
    n_SdpToolHost["SdpToolHost (container)"]
    n_AssessBlueprintEvidence -->|allocated-to| n_SdpToolHost
    n_AssessBlueprintEvidence -->|contributes-to| n_ScopedBlueprintReadiness
    n_BlueprintAssessment -->|owns| n_AssessBlueprintEvidence
```

Source facts: f0005, f0006, f0008.

