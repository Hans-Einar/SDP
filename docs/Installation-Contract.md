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
use normalized forward-slash relative paths. Sources, governing schemas and
exclusions resolve from the repository root; destinations resolve from the
consuming project root.

A conforming client uses one portable segment subset for every manifest source,
destination, governing schema and exclusion. Paths use `/`, are relative and
normalized, and contain no empty, `.` or `..` segment or trailing slash. Segments
contain no control character, tilde, or Windows-invalid `< > : " | ? *` character, do
not end in a dot or space, are not Windows device aliases (including `CON`,
`NUL`, `COM1`-`COM9`, `LPT1`-`LPT9`, superscript-digit `COM`/`LPT` aliases,
`CLOCK$`, `CONIN$` and `CONOUT$`, with or without an extension), and are never
`.git` in any case. Collision checks use a case-folded, physical-alias-normalized
key so Windows-equivalent destinations cannot coexist. A destination also may
not be a case-insensitive ancestor or descendant of another destination, and
every existing destination ancestor must be a directory before a plan can be
applicable.

Lexical containment is insufficient. Before reading or writing, a client must
either resolve every existing source/target ancestor and prove that its physical
path remains under the intended root, or reject the link/reparse-point chain.
The PowerShell reference rejects reparse points and symlinks in the repository
root, project root, backup root and every existing source/destination ancestor,
then repeats the check immediately before mutation. On Windows it opens each
existing root and ancestor and compares the OS-provided volume serial and
128-bit file ID. Source/project and source/backup tree overlap is rejected in
both ancestor directions across local, UNC, `\\?\` drive/UNC and available 8.3
aliases. Failure to acquire a required identity is a preflight failure. Distinct
same-volume siblings remain valid when both ancestor chains reach the same
visible namespace root. Different share roots on the same device are rejected
when neither chain exposes enough ancestry to prove separation. Before removing
an extended prefix, PowerShell 5.1 accepts only canonical `\\?\X:\...` and
`\\?\UNC\server\share\...` spellings. It rejects mixed or forward-slash device
prefixes, forward-slash separators, empty/doubled/trailing segments (except the
required drive-root separator), `.` or `..` segments and segments ending in a
space or dot. Other device namespaces fail closed. Valid long-name and available
8.3 segments remain supported.

A client also rejects duplicate entry IDs or destinations, missing sources,
unsupported schema versions and paths outside their permitted roots. It never
interprets a manifest path relative to the process working directory.

The source need not be a Git checkout. When no trustworthy repository HEAD is
available, generated installed facts record `sourceCommit: null`. A directory or
archive name is not commit evidence. In a Git checkout, a non-null
`sourceCommit` records the available `HEAD` baseline. If the checkout is dirty,
that value is not an attestation that the installed bytes equal the named
commit.

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

The v1 JSON contract is closed-world. Every object is checked for its exact
required and permitted properties, arrays obey their declared minima and
uniqueness, and policy combinations are validated before any target file is
inspected. Copied Toolkit-managed sources are restricted to `Toolkit/payload/`
or `Toolkit/skills/`; copied project-owned sources are restricted to
`Toolkit/project-templates/`. Governing schemas are restricted to
`Toolkit/schemas/` and use these canonical capability pairings:

| Schema | Capability |
|---|---|
| `current-index.schema.json` | `sdp.traceability.current-index.v1` |
| `relations.schema.json` | `sdp.traceability.relations.v1` |
| `ledger-event.schema.json` | `sdp.traceability.ledger-events.v1` |
| project/installed manifest schemas | `sdp.manifest.v1` |
| release/Fix record schemas | `sdp.release.v1` |

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
version reinstall. Its source commit is the available repository HEAD when Git
identity is trustworthy and otherwise null; checkout cleanliness is not encoded,
so the field alone does not attest installed-byte equality. Installed-manifest
equality is semantic: mapping order and
equivalent YAML quoting do not cause a refresh or backup, while changed facts,
ordered capabilities, source commit or Toolkit build identity do. `empty-ledger`
emits an empty NDJSON file.

## Deterministic plan

`-PlanJson` emits the read-only JSON contract defined by
`Toolkit/schemas/SDP-install-plan.schema.json`. It is derived at runtime from the
installation manifest, target observations and the initialize/force options; it
is not a second inventory.

The plan has schema and manifest versions, old/new Toolkit versions, options,
the required constant
`orderingPolicy: "migration-first-manifest-order-v1"`, `canApply`, and an
ordered action list. Actions are `create`, `replace`,
`preserve`, `unchanged`, `backup`, `migrate`, `generate` or `block`. Every action
has a stable sequence, entry ID, nullable `source`, `generator` and
`targetSource`, nullable `targetSourceSha256` and `destinationPrecondition`,
project-relative destination, ownership, stable reason, mutation flag and
relevant versions. Ordinary entry actions identify exactly one manifest source
or generator and set the two migration precondition fields to null. A `migrate`
action instead identifies the existing project file in `targetSource`, records
the lowercase SHA-256 of its exact bytes in `targetSourceSha256`, and requires
`destinationPrecondition: "absent"`; apply writes those exact bytes to the exact
planned destination. A block identifies no content source.

Canonical v1 action order is normative:

1. A blocked plan contains exactly one `block` action at sequence 1.
2. Target-to-target migrations precede every ordinary manifest-entry action.
3. Migration rules are processed in their declared migration-policy order.
   Contract v1 declares one AGENTS migration rule, so it can emit at most one
   migration action.
4. Selected ordinary entries are processed in the exact array order of
   `Toolkit/SDP-install.manifest.json`.
5. A client must not independently sort by destination, entry ID, ownership,
   source, action or any private implementation detail.
6. When an entry requires a backup, `backup` is immediately followed by its
   matching `replace` or `generate`; both actions have identical entry ID,
   source/generator, destination, ownership and old/new version identities.
   No unrelated action may split the pair.
7. An entry producing one action emits that action at its manifest position.
8. Sequence values are assigned after this ordering is complete and are
   contiguous integers beginning at 1.
9. Repeated planning against the same source, options and target bytes produces
   semantically identical ordered actions.

The v1 reason vocabulary and decisions are normative:

| Reason | Action / mutation | Decision condition |
|---|---|---|
| `missing-target` | `create` / yes | Copied entry destination is absent. |
| `missing-generated-target` | `generate` / yes | Generated entry destination is absent. |
| `content-matches` | `unchanged` / no | Existing copied bytes or generated semantic facts match. |
| `missing-only-content` | `preserve` / no | Existing project-owned missing-only content differs. |
| `managed-content-differs` | `preserve` / no | Managed content differs but refresh/force policy does not authorize replacement. |
| `backup-before-replace` | `backup` / yes | A changed managed entry requires backup before refresh. |
| `refresh-managed-content` | `replace` / yes | Policy authorizes replacement of a changed copied managed entry. |
| `refresh-generated-content` | `generate` / yes | Policy authorizes regeneration of changed managed facts. |
| `migrate-existing-agents` | `migrate` / yes | Legacy `AGENTS.md` differs and `AGENTS-project.md` is absent. |
| `preserve-existing-agents-conflict` | `migrate` / yes | Both AGENTS files exist; legacy rules move to `AGENTS-project.migration-sha256-<content-hash>.md`. |
| `malformed-project-manifest` | `block` / no | Strict root-YAML preflight fails for the project manifest. |
| `unsupported-project-schema` | `block` / no | The root project schema version is unsupported. |
| `malformed-installed-manifest` | `block` / no | Strict installed-manifest parsing or facts validation fails. |
| `unsupported-installed-schema` | `block` / no | The root installed schema version is unsupported. |
| `downgrade-blocked` | `block` / no | Installed SemVer precedence is greater than the source Toolkit. |

Semantic plan validation is required in addition to JSON Schema validation.
Sequences start at one and are contiguous; every action old/new version agrees
with the top level; ordinary action entry/source/generator/destination/ownership
facts agree with the installation manifest; the two AGENTS migration shapes are
exact; canonical ordering and backup adjacency are enforced; and `canApply` is
true exactly when there is no block action.

### AGENTS exact-byte migration

For `preserve-existing-agents-conflict`, the deterministic destination is
`AGENTS-project.migration-sha256-<sha256>.md`, where `<sha256>` is lowercase
hexadecimal SHA-256 over the exact legacy `AGENTS.md` bytes. Clients do not
decode the file or normalize line endings before hashing or comparison.

- If the destination is absent, emit exactly one migration-first `migrate`
  action with `targetSource: "AGENTS.md"`, the exact-byte hash,
  `destinationPrecondition: "absent"`, project-owned ownership, the
  deterministic destination and reason `preserve-existing-agents-conflict`.
- If it is a regular file with identical bytes, preservation is already
  complete. Do not overwrite it and do not emit a migration action; continue
  ordinary planning for the managed `AGENTS.md` entry.
- If it is a regular file with different bytes, planning fails before mutation
  with `agents-migration-destination-content-mismatch`.
- If it is a directory, link, symlink, reparse point or other unsupported
  object, planning fails closed with
  `agents-migration-destination-unsupported-object`.

The migration action captures both source hash and absent-destination
precondition. Immediately before apply, the reference installer rechecks that
the source is still an ordinary file with the same exact-byte hash and that the
destination is still absent. A changed source fails with
`agents-migration-source-changed`; a destination that appeared or changed fails
with `agents-migration-destination-changed`. Creation uses exclusive
create-new behavior: no alternate path is selected and project-owned content is
never replaced. Because migration is first, either it completes or fails before
ordinary manifest-entry mutations. A later failure elsewhere can still leave
this already completed preservation file and any earlier planned operations in
place; installation is not a general transaction.

### Applicable, blocked and fatal outcomes

An applicable plan has `canApply: true` and no block action. A blocked plan is a
valid plan with `canApply: false` and exactly one canonical block action at
sequence 1; it represents supported, target-derived decisions such as an
unsupported installed/project schema or downgrade. A fatal pre-plan failure
does not emit a plan because the contract or preservation-sensitive filesystem
facts cannot be represented safely. The closed v1 fatal vocabulary is:

- `install-manifest-invalid`
- `agents-migration-destination-content-mismatch`
- `agents-migration-destination-unsupported-object`
- `agents-migration-source-changed`
- `agents-migration-destination-changed`

Raw localized PowerShell exception text is not part of the portable contract.

Plan mode performs no target mutation: it creates no directory, file, backup or
installed timestamp. Volatile timestamps, absolute extraction paths and
generated backup-directory names are outside the portable plan contract. Human
`-Preview` output is also mutation-free but is not the conformance format.

## Shared conformance package

`Toolkit/conformance/install-v1/` is the checked-in, language-neutral authority
for portable scenarios. `scenarios.json`, validated by
`scenario-index.schema.json`, declares source mode, options, portable before
state, expected outcome and important preservation assertions. Outcomes under
`expected/` are either complete plan JSON conforming structurally and
semantically to the plan schema, or a two-field fatal object containing
`kind: "fatal"` and one closed `failureClass`.

The package covers empty/default, initialize, repeat, repeat-initialize, legacy
AGENTS preservation and collision cases, upgrade, force, project-owned
preservation, source archive without `.git`, unsupported schemas, downgrade and
malformed-manifest behavior. It contains no absolute or temporary paths and a
consumer may parse and test it without PowerShell. The Python reference harness
can additionally materialize every scenario and compare the PowerShell
reference implementation when PowerShell is available.

Committed outcomes are reviewed contract authorities. Normal validation only
compares against them and never regenerates them from the implementation.
Maintainers may explicitly write candidate outcomes with
`run_conformance.py --write-candidates`, then must review and commit the diff.

## External installer conformance

An external client such as future `gh-sdp` must:

1. locate and schema-validate the manifest before target mutation;
2. verify Toolkit/capability agreement and reject unsupported schemas or a
   Toolkit downgrade using full SemVer 2.0 precedence (build metadata changes
   identity but not precedence);
3. select only explicit entries and honor every ownership and policy field;
4. preserve project-owned content under normal, forced and repeated runs;
5. expose a mutation-free preview and a plan conforming to the plan schema;
6. consume `Toolkit/conformance/install-v1/scenarios.json` directly and produce
   the checked-in normalized plans or failure classes for every scenario,
   including exact canonical order and exact-byte preservation;
7. prove absence of live Toolkit state, especially
   `SDP/Releases/REL-0.2.0.yaml`.

This contract does not prescribe private Go structures or implement `gh-sdp`.
Conformance concerns public inputs, outputs and observable filesystem behavior.
