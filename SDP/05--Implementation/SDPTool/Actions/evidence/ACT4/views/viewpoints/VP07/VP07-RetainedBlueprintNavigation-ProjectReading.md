# Feature: RetainedBlueprintNavigation — mode ProjectReading

[Viewpoint](index.md) · [Navigator](../../navigator.md)

Revision: `02ca3d1cf871998399f43803bbf47f6fa73321b38ce22301b196973963cd3480`.

## Feature: RetainedBlueprintNavigation — mode ProjectReading

```mermaid
flowchart LR
    n_BlueprintCatalogue["BlueprintCatalogue (unit)"]
    n_BrowseRetainedBlueprints["BrowseRetainedBlueprints (functionality)"]
    n_RetainedBlueprintNavigation["RetainedBlueprintNavigation (feature)"]
    n_SdpToolHost["SdpToolHost (container)"]
    n_VerifyRetainedBlueprints["VerifyRetainedBlueprints (functionality)"]
    n_BlueprintCatalogue -->|owns| n_BrowseRetainedBlueprints
    n_BlueprintCatalogue -->|owns| n_VerifyRetainedBlueprints
    n_BrowseRetainedBlueprints -->|allocated-to| n_SdpToolHost
    n_BrowseRetainedBlueprints -->|contributes-to| n_RetainedBlueprintNavigation
    n_VerifyRetainedBlueprints -->|allocated-to| n_SdpToolHost
    n_VerifyRetainedBlueprints -->|contributes-to| n_RetainedBlueprintNavigation
```

Source facts: f0015, f0016, f0034, f0035, f0193, f0194.

