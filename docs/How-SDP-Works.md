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
├── public release identity (SemVer)
└── Sprint or Refactor
    └── Iteration
        └── Slice or bounded Fix
```

A revision is only a small correction inside the same planned Slice/Fix. It is
not a route around planning.

## Roles

The Master coordinates active work and release target, delegates implementation,
requests fresh verification/review, records only truthful events and stops at the
work boundary. Workers implement one contract. Reviewers inspect the actual diff
and evidence independently. Architects own long-term boundaries and
compatibility.

Task skills add concrete release, versioning, auditing, verification,
traceability and refactor procedures.

## Versions and manifests

- Toolkit release facts: root `SDP.manifest.yaml`
- Installed Toolkit facts: `SDP/Framework/installed-toolkit.manifest.yaml`
- Project release/work state: `SDP/SDP-project.manifest.yaml`
- Dynamic Git/build identity: generated JSON

Released software and Toolkits use SemVer. Sprint/Iteration/Slice/Fix identities
remain separately addressable. Unreleased builds must not present themselves as
released.

## Release notes and history

All notable changes enter the editable `Unreleased` section. At release
preparation, selected entries move into a dated version section. Released
sections are immutable; corrections are explicit events, never silent edits.

## Traceability

`CurrentIndex.yaml` records actual current state. `Relations.yaml` connects stable
IDs. `Ledger.ndjson` records append-only transitions, including release events.
Publication events require real tag/commit/GitHub Release identities.

## Verification and release gate

Completion requires evidence appropriate to the work. A release gate checks
included work, verification, review findings, schemas, traceability, notes,
version agreement, migration, clean state, tag uniqueness and immutable history.
Missing evidence is failure, not an inference.

## Truthful two-phase publication

A release-preparation commit cannot claim a future GitHub Release or contain its
own SHA. SDP therefore prepares and approves first, publishes only after explicit
human authorization, then adds a small reconciliation commit recording the real
tag and GitHub Release.

## Installation and local adaptation

The installer refreshes clearly managed AGENTS/Framework/skill files while
preserving project-owned documents. Missing manifests, notes and templates are
added safely. Supported old installations migrate additively; unsupported schemas
stop without mutation.

Projects may extend SDP, but local rules must preserve stable IDs, truthful
traceability, review independence, evidence requirements and publication safety.
