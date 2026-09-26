# VP02 — Architecture and logical decomposition

Revision: `c62878ccf713fe1fd51d9a7d14e90dd7f8a70dd5deb346af2d8fb9edae9969d2`.

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


