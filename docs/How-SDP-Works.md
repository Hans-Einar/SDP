# How SDP Works

## Source of truth

The repository is authoritative. Chat history and agent memory are temporary.
Material decisions, evidence and state transitions belong in repository records.

## Document progression

```text
Mandate -> Study -> Requirements -> Architecture
        -> Design Analysis -> Design -> Implementation
```

Documents remain living records, but later work must not silently contradict
approved earlier contracts.

## Horizontal design, vertical delivery

Architecture and design describe horizontal layers and contracts. Delivery is
organized as vertical Sprints/Refactors, Iterations and Slices that produce a
small coherent capability across the required layers.

```text
Project
|-- public release identity (SemVer)
`-- Sprint or Refactor
    `-- Iteration
        `-- Slice or bounded Fix
```

A revision is only a small correction inside the same planned Slice/Fix. It is
not a route around planning.

## Roles

The Master coordinates active work and release target, delegates implementation,
requests fresh verification/review, records only truthful events and stops at the
work boundary. Workers implement one contract. Reviewers inspect the actual diff
and evidence independently. Architects own long-term boundaries and
compatibility.

## Versions and manifests

- Toolkit release facts: root `SDP.manifest.yaml`;
- installable inventory and policies: `Toolkit/SDP-install.manifest.json`;
- installed Toolkit facts: `SDP/Framework/installed-toolkit.manifest.yaml`;
- project release/work state: `SDP/SDP-project.manifest.yaml`;
- dynamic Git/build identity: generated JSON.

Released software and Toolkits use SemVer. Work identities remain separately
addressable. Unreleased builds must not present themselves as released.

## Release notes and traceability

Notable changes enter editable `Unreleased`. At release preparation selected
entries move into a dated version section. Released sections are immutable;
corrections are explicit events.

CurrentIndex records actual current state, Relations connects stable IDs and the
Ledger records append-only transitions. Publication events require real
tag/commit/GitHub Release identities.

## Verification and publication

Completion requires evidence appropriate to the work. A release gate checks
included work, verification, review findings, schemas, traceability, notes,
version agreement, migration, clean state, tag uniqueness and immutable history.
Missing evidence is failure, not an inference.

Release preparation cannot claim a future GitHub Release or contain its own SHA.
SDP therefore prepares and approves first, publishes only after explicit human
authorization, then reconciles the real tag and Release identity.

## Installation and local adaptation

The canonical JSON installation manifest drives managed refresh, missing-only
project seeds and generation. Neutral project templates live physically under
`Toolkit/project-templates/`. The Toolkit repository's root lifecycle, release,
review, verification and traceability records are live instance state and are
never template inputs.

PowerShell and external clients must preserve project-owned files, reject unsafe
or unsupported contracts before mutation and expose a mutation-free plan.
Projects may extend SDP, but local rules must preserve stable IDs, truthful
traceability, review independence, evidence and publication safety.
