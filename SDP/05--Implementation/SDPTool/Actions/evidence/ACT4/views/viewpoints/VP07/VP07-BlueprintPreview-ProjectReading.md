# Feature: BlueprintPreview — mode ProjectReading

[Viewpoint](index.md) · [Navigator](../../navigator.md)

Revision: `02ca3d1cf871998399f43803bbf47f6fa73321b38ce22301b196973963cd3480`.

## Feature: BlueprintPreview — mode ProjectReading

```mermaid
flowchart LR
    n_BlueprintFacade["BlueprintFacade (unit)"]
    n_BlueprintPreview["BlueprintPreview (feature)"]
    n_GenerateBlueprintPreview["GenerateBlueprintPreview (functionality)"]
    n_SdpToolHost["SdpToolHost (container)"]
    n_BlueprintFacade -->|owns| n_GenerateBlueprintPreview
    n_GenerateBlueprintPreview -->|allocated-to| n_SdpToolHost
    n_GenerateBlueprintPreview -->|contributes-to| n_BlueprintPreview
```

Source facts: f0023, f0064, f0065.

