# VP02 — Architecture and logical decomposition

Revision: `dddd841c13f7b842c7f5e404aeb6ac2abc599ef1c382f54c21548aef53006417`.

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


