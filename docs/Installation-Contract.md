# Installation Contract

## Authorities

Two manifests have different, non-overlapping authority:

- root `SDP.manifest.yaml` owns Toolkit release, version and capability facts;
- `Toolkit/SDP-install.manifest.json` owns the complete installable inventory,
  destinations, ownership and update policies.

The installation manifest conforms to
`Toolkit/schemas/SDP-install-manifest.schema.json`. It is the source of truth for
PowerShell and external clients. Executable code must not contain a second file
inventory that can drift silently.

## Archive and path resolution

In a clone or extracted GitHub source archive, locate the directory containing
`Toolkit/SDP-install.manifest.json`. The manifest's
`sources.repositoryRoot: ".."` resolves from the manifest directory to that
archive root. All entry sources, destinations, governing schemas and exclusions
use normalized forward-slash repository-relative paths.

A conforming client rejects absolute paths, backslashes, empty or dot segments,
traversal, duplicate entry IDs or destinations, missing sources, unsupported
schema versions and paths that resolve outside their permitted roots. It must
not interpret paths relative to the process working directory.

The source need not be a Git checkout. When no trustworthy repository HEAD is
available, generated installed facts record `sourceCommit: null`. A directory or
archive name is not commit evidence.

## Entry model

Every installable target is one explicit entry; v1 has no glob or directory
expansion. Each entry declares:

- stable `id`, `copied` or `generated` kind, and source or generator;
- normalized destination and `toolkit-managed` or `project-owned` ownership;
- `default` or `initialize-only` selection;
- `always` or `missing-only` installation;
- `always`, `upgrade-or-force` or `never` refresh;
- `before-replace`, `migration-aware` or `none` backup;
- `replace-managed` or `preserve` force behavior;
- any migration policy and governing schema/capability.

Project-owned entries are always missing-only, never refreshed, never backed up
for replacement and preserved under force. Initialize-only entries are
project-owned. Toolkit-managed files may refresh only as their entry permits.
The managed `AGENTS.md` entry additionally preserves pre-SDP instructions through
the declared AGENTS migration policy.

## Source and ownership classes

| Class | Canonical source | Installed behavior |
|---|---|---|
| Repository instance records | root `Sprints/`, `Releases/`, `Traceability/`, reviews and related live folders | Excluded; never installed |
| Toolkit-managed copies | `Toolkit/payload/`, `Toolkit/skills/` | Refreshed and backed up according to entry policy |
| Neutral project seeds | `Toolkit/project-templates/` | Project-owned immediately; created only when missing |
| Generated managed facts | `installed-toolkit-manifest` generator | Regenerated from declared facts; source commit may be null |
| Generated project history | `empty-ledger` generator | Creates an empty Ledger only when missing; never replaces history |

The legacy root `payload/` and `skills/` trees and the specifically listed legacy
project-owned seed copies under `Toolkit/payload/` are excluded. Managed
Framework templates remain installable only where an explicit entry lists them.
No unlisted source file becomes installable merely because it is in the archive.

Default entries create the project manifest, neutral release notes, neutral
CurrentIndex and Relations, an empty Ledger, installed facts and required agent,
Framework and skill files. `-InitializeProjectStructure` additionally selects
only neutral lifecycle and operating README/document-guide seeds. Neither mode
copies `Releases/REL-0.2.0.yaml`, active Sprints/Refactors, populated Toolkit
release notes, Ledger history, review or verification evidence.

## Generated facts

The manifest declares both v1 generators. `installed-toolkit-manifest` supplies
the static Toolkit, Framework, AGENTS contract, installer, skill and capability
facts. Its installation timestamp is UTC and is preserved on a same-Toolkit-
version reinstall. Its source commit is repository HEAD when trustworthy and
otherwise null. `empty-ledger` emits an empty NDJSON file.

## Deterministic plan

`-PlanJson` emits the read-only JSON contract defined by
`Toolkit/schemas/SDP-install-plan.schema.json`. It is derived at runtime from the
installation manifest, target observations and the initialize/force options; it
is not a second inventory.

The plan has schema and manifest versions, old/new Toolkit versions, options,
`canApply`, and an ordered action list. Actions are `create`, `replace`,
`preserve`, `unchanged`, `backup`, `generate`, `warn` or `block`. Every action
has a stable sequence, entry ID, nullable source and generator fields,
project-relative destination, ownership, stable reason, mutation flag and
relevant versions. Non-warning/non-blocking entry actions identify exactly one
source or generator. A plan with `canApply: false` contains a `block` action.

Plan mode performs no target mutation: it creates no directory, file, backup or
installed timestamp. Volatile timestamps, absolute extraction paths and
generated backup-directory names are outside the portable plan contract. Human
`-Preview` output is also mutation-free but is not the conformance format.

## External installer conformance

An external client such as future `gh-sdp` must:

1. locate and schema-validate the manifest before target mutation;
2. verify Toolkit/capability agreement and reject unsupported schemas or a
   Toolkit downgrade;
3. select only explicit entries and honor every ownership and policy field;
4. preserve project-owned content under normal, forced and repeated runs;
5. expose a mutation-free preview and a plan conforming to the plan schema;
6. produce equivalent normalized actions for shared empty, legacy, upgrade,
   force, initialize, repeat-initialize, archive and error fixtures;
7. prove absence of live Toolkit state, especially
   `SDP/Releases/REL-0.2.0.yaml`.

This contract does not prescribe private Go structures or implement `gh-sdp`.
Conformance concerns public inputs, outputs and observable filesystem behavior.
