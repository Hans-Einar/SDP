# VP02 — Architecture and logical decomposition

Revision: `c05cf8fcba4a1c7abb985b1a3814bf9690759dfa74d12cdf58d20a9c46d96330`.

Container/Unit and contains; library structure is not deployment allocation.

Selection: `{"viewpoint":"VP02","direction":"both","depth":2,"diagram":"VP02-GhSdpProcess"}`.

## Logical decomposition: GhSdpProcess

```mermaid
flowchart LR
    n_GhSdpLauncher["GhSdpLauncher (unit)"]
    n_GhSdpProcess["GhSdpProcess (container)"]
    n_GhSdpProcess -->|contains| n_GhSdpLauncher
```

Source facts: f0104.


