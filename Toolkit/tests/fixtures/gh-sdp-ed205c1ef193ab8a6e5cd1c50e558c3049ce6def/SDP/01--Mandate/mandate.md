# MAN-001 — gh-sdp GitHub CLI Extension Mandate

Status: proposed — ready for Steering Group assessment
Lifecycle level: Mandate only
Completed Phase 1 work: `Sprint-001` / `SPI-001` / `SPS-001`
Proposed client release: `REL-0.1.0` (`gh-sdp` `0.1.0`, unreleased)
Proposed first Toolkit compatibility target: SDP Toolkit `0.2.0` (unreleased)

Identifiers in this document follow the project-local convention in
`SDP/Instructions/StableIdentifiers.md`; they are not claimed as an upstream SDP
format.

## Problem

Using SDP in a consuming repository currently requires people and automation to
find the correct Toolkit source or release, understand its installer contract,
invoke platform-specific tooling, and distinguish Toolkit-managed files from
project-owned records. Repeating that work separately across repositories makes
installation and inspection difficult, makes safe updates harder to reason
about, and creates a risk that a client silently diverges from the canonical
installer's preservation and compatibility rules.

GitHub CLI users need one discoverable, project-oriented interface that can
install, inspect, preview, and update SDP without obscuring which product version
is involved or putting existing repository content at risk.

## Users and Stakeholders

- Primary users are maintainers, contributors, and automation operators who use
  GitHub CLI and want SDP in a project repository.
- Consuming-project owners are stakeholders because their project-owned content,
  history, local policy, and release state must remain authoritative.
- `Hans-Einar/SDP` maintainers are stakeholders because the extension must honor
  the canonical Toolkit's versioned installer and compatibility rules without
  creating a competing rule set.
- Security reviewers, release reviewers, and the Steering Group are stakeholders
  in provenance, platform behavior, evidence, lifecycle decisions, and truthful
  release claims.

## Purpose and Desired Outcomes

The purpose of `gh-sdp` is to become an official GitHub CLI extension, installed
and invoked as `gh sdp`, that provides a safe, comprehensible front door to the
canonical SDP Toolkit for a selected project.

- **MAN-OUT-001 — Discoverable command surface.** Users can use one extension to
  install, inspect, preview, update, identify, and learn about SDP from a project
  context.
- **MAN-OUT-002 — Safe project-local operation.** Installation and update
  behavior protects project-owned records and makes intended changes visible
  before they are applied.
- **MAN-OUT-003 — Truthful version and release selection.** Client and Toolkit
  identities remain separate, and stable Toolkit releases are selected by
  default without overstating availability or support.
- **MAN-OUT-004 — Cross-platform experience.** Windows, Linux, and macOS users
  receive an appropriate native experience; Linux and macOS users are not forced
  to install or invoke PowerShell.
- **MAN-OUT-005 — Canonical rule alignment.** Extension behavior remains
  demonstrably aligned with versioned `Hans-Einar/SDP` installer, ownership, and
  compatibility rules instead of allowing an independently copied rule set to
  drift.
- **MAN-OUT-006 — Secure and reliable operation.** The extension treats remote
  artifacts, repository paths, existing content, credentials, interruption, and
  failure as trust and recovery concerns rather than happy-path details.

## Future High-Level Command Surface

The intended high-level surface is directional, not an implemented or approved
requirements contract:

| Future command | Mandate-level intent |
| --- | --- |
| `gh sdp install` | Install an appropriate SDP Toolkit into the selected project. |
| `gh sdp status` | Report project installation, compatibility, and relevant client/Toolkit identities. |
| `gh sdp preview` | Show the changes an install or update would propose without mutating the project. |
| `gh sdp update` | Move an existing project installation to an eligible Toolkit release safely. |
| `gh sdp version` | Identify the extension and, where relevant, the installed/selected Toolkit separately. |
| `gh sdp help` | Explain commands, selection, safety boundaries, and evidence in GitHub CLI form. |

Project selection may later receive explicit syntax, but the default project is
the current working directory. Detailed flags, exit codes, output formats,
network behavior, and non-interactive semantics belong to later Study and
Requirements work.

## Product and Repository Relationship

`Hans-Einar/gh-sdp` is the client/extension project. `Hans-Einar/SDP` is a
separate, canonical upstream repository that owns the SDP Toolkit, its installer
rules, schemas, and release evidence. The extension must consume an evidenced,
versioned Toolkit contract and must not recast its own interpretation as
canonical. Any change to the upstream repository requires separate authority and
its own SDP work.

The default source is a stable Toolkit release from `Hans-Einar/SDP`. Use of an
unreleased revision, prerelease, alternate repository, or offline source is a
later policy question and must never be mislabeled as the stable default.
Conformance evidence must guard against drift between extension behavior and the
selected Toolkit's installation and preservation rules.

## Version and Lifecycle Truth

The proposed `gh-sdp` development target is client version `0.1.0`, represented
by `REL-0.1.0`. It is unreleased. The proposed first supported compatibility
target is SDP Toolkit `0.2.0`; at the inspected canonical source revision it is
also unreleased and has no evidenced tag or GitHub Release. These numbers belong
to different products and do not advance together automatically.

This document establishes intent only. Phase 1 preparation passed fresh
independent review and `MAN-001` is proposed for Steering Group assessment. No
Study, Requirements, Architecture, Design Analysis, Design, or Implementation
decision is approved; no product capability, compatibility promise, tag, or
release exists as a result of `SPS-001`.

## Boundaries and Explicit Exclusions

- **MAN-BND-001 — Mandate-only phase.** `SPS-001` creates governance, Mandate,
  traceability, and verification records only. Product code, a Go module,
  commands, handlers, installer/downloader/archive/migration logic, tests of
  product behavior, builds, packages, and workflows are excluded.
- **MAN-BND-002 — Independent upstream.** This project does not modify, tag,
  release, or publish `Hans-Einar/SDP` without a separately authorized task.
- **MAN-BND-003 — Separate identities.** A `gh-sdp` version never stands in for
  a Toolkit version. Proposed or source-state versions are not described as
  published or supported without real repository evidence.
- **MAN-BND-004 — Project ownership.** Existing project-owned content must be
  preserved. No install or update may silently replace, delete, or reinterpret
  it merely because Toolkit-managed content also exists in the project.
- **MAN-BND-005 — Stable-by-default Toolkit.** Normal selection resolves an
  eligible stable Toolkit release from the canonical upstream; use of any other
  source or release state must be explicit and truthful.
- **MAN-BND-006 — Cross-platform without forced PowerShell.** The product intent
  covers Windows, Linux, and macOS. Linux and macOS operation cannot require
  PowerShell; the later architecture may choose suitable host-native or portable
  mechanisms.
- **MAN-BND-007 — Project selection.** Commands act on a selected project and
  default to the current working directory. Detailed selection syntax and path
  semantics remain open for later lifecycle work.
- **MAN-BND-008 — No premature solution lock-in.** Command names and product
  outcomes are in scope, but transport, archive, migration, storage, update,
  authentication, and implementation architecture are not decided here.

Publishing `gh-sdp` `0.1.0`, declaring Toolkit `0.2.0` support, and beginning
Phase 2 are explicitly outside this Mandate Slice.

## Security and Reliability Objectives

At Mandate level, future product decisions must collectively:

- establish Toolkit source, version, provenance, and integrity before applying
  remote content, and prefer stable canonical releases by default;
- constrain reads and writes to the intended project and safe destinations,
  account for archive/path-link hazards, avoid unnecessary privilege, and avoid
  exposing credentials or sensitive repository data;
- preserve project ownership and provide a deterministic preview whose material
  decisions agree with apply behavior;
- fail closed on unsupported, ambiguous, malformed, or incompatible state, with
  actionable diagnostics and no false success;
- make repeated operations predictable and make interrupted or failed changes
  recoverable, with evidence appropriate to the risk; and
- verify cross-platform conformance to the selected Toolkit contract so safety
  and ownership rules cannot drift between the extension and canonical
  installer.

The precise controls, threat model, trust mechanism, recovery model, and test
matrix are later Study, Requirements, and Architecture work.

## Success Criteria

These criteria describe future Mandate success and are not claims that `SPS-001`
has delivered product behavior:

- **MAN-SC-001.** A GitHub CLI user can discover and invoke the six named command
  intents against an explicitly selected project or, by default, the current
  working directory.
- **MAN-SC-002.** Preview, install, and update preserve project-owned content and
  demonstrate predictable, idempotent, and recoverable behavior with clear
  failure reporting.
- **MAN-SC-003.** Stable canonical Toolkit releases are selected by default, and
  every status/version/release claim distinguishes the client, selected Toolkit,
  installed Toolkit, and their actual release states.
- **MAN-SC-004.** Supported behavior is verified on Windows, Linux, and macOS,
  without a PowerShell dependency for Linux or macOS users.
- **MAN-SC-005.** Automated conformance evidence detects drift from each
  supported Toolkit release's installer, ownership, and compatibility contract.
- **MAN-SC-006.** Security assessment and deterministic tests cover provenance,
  integrity, unsafe paths or links, credentials, privilege, interruption,
  rollback/recovery, and incompatible or malformed input before support is
  declared.
- **MAN-SC-007.** Each lifecycle and release claim is backed by current
  project-local SDP traceability, repository evidence, applicable tests, and a
  fresh independent review.

## Assumptions and Constraints

### Assumptions

- **MAN-ASM-001.** GitHub CLI continues to provide a viable extension mechanism
  for the intended users and platforms.
- **MAN-ASM-002.** The canonical Toolkit will provide stable, versioned,
  machine-consumable release evidence and an installation contract suitable for
  client conformance before compatibility is declared.
- **MAN-ASM-003.** A user selecting a project can grant the minimum filesystem
  access needed for that project; broader privileges should not be required.
- **MAN-ASM-004.** A consuming repository's project-local SDP manifest and
  Toolkit-managed facts can remain the authoritative basis for detecting its
  installed state.

### Constraints

- The extension identity is `gh-sdp`, invoked through GitHub CLI as `gh sdp`.
- The default project is the current working directory.
- Client and Toolkit versions and release states are independent.
- `Hans-Einar/SDP` is the default canonical Toolkit upstream.
- Project-owned content preservation, stable-by-default selection,
  cross-platform intent, and no forced PowerShell on Linux/macOS are product
  constraints, not optional enhancements.
- Repository evidence, deterministic verification, and fresh independent review
  govern support and release claims.
- Go is the preferred implementation direction because a portable extension
  binary may fit the cross-platform goal, but this is not an architecture
  decision. Study and Architecture must evaluate and either confirm or replace
  that direction.

## Unresolved Strategic Questions

- **MAN-Q-001.** What canonical release assets, metadata, checksums, signatures,
  or attestations will establish Toolkit provenance and integrity?
- **MAN-Q-002.** How will client/Toolkit compatibility, migrations, rollback,
  interrupted operations, and version skew be modeled?
- **MAN-Q-003.** Which offline, alternate-upstream, fork, prerelease, or private
  repository use cases should be supported, and how will they be made explicit?
- **MAN-Q-004.** Which GitHub authentication, API, caching, and rate-limit
  behaviors are necessary without exposing credentials or coupling basic local
  inspection to network availability?
- **MAN-Q-005.** Does Study confirm Go as the best implementation direction, or
  does another approach better satisfy portability, extension distribution,
  security, and canonical-rule reuse?
- **MAN-Q-006.** What exact command flags, output contracts, exit codes,
  interactive/non-interactive behavior, and accessibility expectations belong in
  Requirements?
- **MAN-Q-007.** What filesystem, archive, link, locking/concurrency, backup, and
  recovery semantics are required on each supported platform?

These questions intentionally remain open and are inputs to later lifecycle
work, not permission to begin it during `SPS-001`.
