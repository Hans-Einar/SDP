# VP01 — Use cases and traceability

Revision: `7cb1cdf04b9bb1b9779aa570f9fdbeedc313873a8838d94908e7e6118829e936`.

pursues, supports and contributes-to; modeled scope without an invented System boundary.

Selection: `{"viewpoint":"VP01","focus":"ProjectMaintainer","relations":["pursues","supports"],"direction":"both","depth":2}`.

## Use case: AdoptManualProcess

```mermaid
flowchart LR
    n_AdoptManualProcess["AdoptManualProcess (usecase)"]
    n_ProjectMaintainer["ProjectMaintainer (actor)"]
    n_SdpTool["SdpTool (feature)"]
    n_ProjectMaintainer -->|pursues| n_AdoptManualProcess
    n_SdpTool -->|supports| n_AdoptManualProcess
```

Source facts: f0435, f0506.

## Use case: InstallProjectProcess

```mermaid
flowchart LR
    n_InstallProjectProcess["InstallProjectProcess (usecase)"]
    n_ProjectMaintainer["ProjectMaintainer (actor)"]
    n_SdpTool["SdpTool (feature)"]
    n_ProjectMaintainer -->|pursues| n_InstallProjectProcess
    n_SdpTool -->|supports| n_InstallProjectProcess
```

Source facts: f0436, f0509.

## Use case: RecoverProjectUpgrade

```mermaid
flowchart LR
    n_ProjectMaintainer["ProjectMaintainer (actor)"]
    n_RecoverProjectUpgrade["RecoverProjectUpgrade (usecase)"]
    n_SdpTool["SdpTool (feature)"]
    n_ProjectMaintainer -->|pursues| n_RecoverProjectUpgrade
    n_SdpTool -->|supports| n_RecoverProjectUpgrade
```

Source facts: f0437, f0511.

## Use case: UpgradeProjectProcess

```mermaid
flowchart LR
    n_ProjectMaintainer["ProjectMaintainer (actor)"]
    n_SdpTool["SdpTool (feature)"]
    n_UpgradeProjectProcess["UpgradeProjectProcess (usecase)"]
    n_ProjectMaintainer -->|pursues| n_UpgradeProjectProcess
    n_SdpTool -->|supports| n_UpgradeProjectProcess
```

Source facts: f0438, f0513.


