# VP02 — Architecture and logical decomposition

Revision: `7cb1cdf04b9bb1b9779aa570f9fdbeedc313873a8838d94908e7e6118829e936`.

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


