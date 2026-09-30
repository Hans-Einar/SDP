# SDL source-owned composition — SSD2 design record

Implementation follow-up: [design-core 0.6](../../../../SDL/docs/profiles/SDL-Source-Composition-Profile.md)
and [PLAN-SDP-0012 evidence](../../../05--Implementation/SDL/SourceComposition/Evidence.md)
now define delivered scope. The proposal language below records the pre-implementation
design; compact System declaration-plus-membership was deferred.

Owner direction, 2026-09-30: inclusion belongs in SDL source; no separate authored
.design-set.json. This document supersedes the external-manifest recommendation
in [SSD1](../SourceSets/Contract.md), whose original body/evidence are historical.
[PLAN-SDP-0011](Plan.md) records the correction. Syntax below is proposed, not
implemented in design-core/0.5. Use design-core/0.6 only as a candidate profile label.

## What was already decided

The earlier workspace proposal selected an explicit source entry with one declared
System. The MVP1 experiment already uses `includes "path.design".`. The tracked Go
parser still does not implement either mechanism. SSD1 introduced the extra JSON
file as an agent design choice, not as a prior owner decision. That choice is now
withdrawn. Source membership must be discoverable by following the SDL source.

## 1. Source entry and path-addressed membership

Recommended canonical spelling, retaining current lower-case keywords:

```text
language design-core version 0.6.
system MVP1.
includes "Features.design".
MVP1 contains Containers/MachineService.
MVP1 contains Channels/Measurements.
MVP1 contains Contracts/Measurement.
```

`Containers/MachineService` means load `Containers/MachineService.design` and
reference its declared `MachineService`. That file must declare that object once:

```text
language design-core version 0.6.
container MachineService.
```

The path-containing relation supplies both a source dependency and a typed model
membership fact. It does not redeclare the container. Thus there is no second list
to keep synchronized and no duplicate declaration in the entry file. The equivalent
expanded form is useful when the source filename differs from the object name:

```text
includes "Containers/MachineService.design".
MVP1 contains MachineService.
```

These forms yield the same linked fact. The shorthand's path is a locator, not
object identity. A declaration may move while retaining its name: update the path,
or use explicit includes plus the symbolic relation when the filename differs.
An arbitrary `Boundary.design` is handled by explicit includes; do not introduce
a directory-to-Boundary fallback or another export/primary-object registry.

The owner's compact `system MVP1 contains Containers/MachineService.` can be a
syntax shortcut for the System declaration plus that relation. Recommend defining
one canonical form first (the two sentences above); shorthand is a spelling choice
to settle before parser work. Upper-case `System` in the owner's sketch is not
silently adopted as a second keyword spelling. `Container/` or `Containers/` is
an exact directory name, not a magical singular/plural alias; use the established
Containers convention for new source organization.

A path operand contains at least one slash, consists of slash-separated identifier
segments and omits the `.design` suffix. A bare name is always a symbol: it never
causes a search for Name.design. For unusual filenames or a root-level file use
quoted includes. No globs, implicit scans, macro substitution or conditional source
selection in this increment. A file may contain several declarations and may
include further files. Features.design needs no fake container merely to be loaded.

## 2. Paths, scope and membership

All paths are relative to the selected entry file's directory (the System source
root), including includes inside child files. This avoids hidden changes in meaning
when a child file is moved. Use forward slashes and exact case. Append .design
exactly once for shorthand; includes uses the explicit .design filename. The same
limits/root-confinement checks from SSD1 apply, excluding manifest-only limits.
No absolute paths, parent traversal, symlinked components beneath the root, aliases
of one file under different paths, or non-regular sources. Load-once identity is
normalized root-relative path plus verified filesystem identity, not content hash:
two distinct files with equal bytes still contain duplicate declarations.

System membership is distinct from runtime ownership. Recommend broadening the
proposed `System contains` signature to any currently supported declaration kind
except System, covering the owner's channels/contracts as well as containers,
features and requirements-era actors/use cases. This is a deliberate change from
the earlier System×Container-only proposal; its type rule does not become the
unrestricted Unit×anything relation. Existing unit containment/ownership/allocation
signatures and cycle rules remain. A System is still not a Unit subtype, executable
participant or runtime allocation target.

Every loaded declaration is in the entry System's namespace. Explicit System
contains facts select modeled membership relationships, not a duplicate source
inventory. Do not require a second contains statement for every declaration loaded
through includes. An explicit membership fact remains meaningful even if both
objects were declared in the same file. Projections distinguish System membership
from runtime Unit containment; they do not manufacture deployment edges.

There is exactly one System declaration, in the selected entry. Loading another
System is an error, not cross-System import. File paths do not create namespaces
or private visibility; names remain unique across this System. The resolved kind
comes from the declaration, not from the folder label: a Contracts directory does
not turn a container declaration into a contract. Unresolved names and conflicting
declarations remain errors. Cross-System exports, visibility and dependency
versioning retain their explicit later scope from SSD1/KB-SDL-005.

## 3. A preprocessing phase without text pasting

Use the existing grammar frontend to parse each source once into a FileAST that
also records includes and path operands. An I/O-owning resolver follows those
references. There is no second regex parser for dependency discovery and no
replacement of source text before it reaches the real parser.

Per compilation, keep a table keyed by resolved file identity with states unseen,
loading, loaded and failed. Mark loading before following dependencies:

1. Parse the selected entry, retain original bytes, file identity and all spans.
2. For each include/path operand, record a dependency edge at the written location.
3. If unseen, load and parse the target. If loading or loaded, reuse its entry;
   never paste or append its declarations again. Retain every edge for diagnostics.
4. Continue until all reachable files have been parsed or an error occurs.
5. Collect all declarations, resolve path targets and names, then validate the
   complete typed model. Publish no validated candidate when errors remain.

Proposed inclusion-cycle policy: allow a cycle in the **source dependency graph**.
A→B→A terminates because A is already loading. Declaration resolution happens
after graph discovery; no source file executes or initializes another. This also
handles the common diamond A→B→Shared and A→C→Shared without author-written guards.
Self-inclusion is a redundant edge and produces a nonfatal diagnostic. Repeated
includes can similarly be reported without duplicating AST nodes.

This policy does not legalize a semantic cycle such as UnitA contains UnitB and
UnitB contains UnitA. Dependency traversal and model constraints are different
checks. Two distinct files declaring MachineService are a real duplicate-definition
error, even when their text matches. Report both declaration locations and useful
include paths. A cycle containing a missing or invalid file still fails the entire
candidate; load-once never hides a failure.

No pragma once, header guards, forward-declaration files or author-managed include
order are required. This is a proposed compiler guarantee, not a statement that
current SDL already has it. Future executable initialization, if ever added to this
structural profile, would need separate cycle/ordering semantics.

## 4. Provenance, formatting and consumers

Keep all SSD1 obligations for original file/byte spans, immutable checked models,
complete-graph validation, stable semantic identity, deterministic projection and
stale-candidate rejection. Inclusion is a typed source directive; it is not a fact
to show as a runtime channel. A path relation retains both its written path span
and its linked target declaration. Forward references are resolved after loading.

Canonical file sections: profile header, declarations, includes sorted by normalized
path, ordinary facts sorted by existing canonical rules. A combined System shortcut,
if selected, formats to separate declaration/fact sentences. Formatting preserves
file partition and explicit directives; it does not move every include to the root.
It must not silently replace a path operand with a symbol and lose its load edge.

A new 0.6 candidate revision hashes a domain-separated record of language/profile
identity, entry path and sorted reachable (relative path, exact bytes hash) pairs.
Use SSD1's length-prefix encoding, with a new domain `SDL-source-graph-revision/1`;
there is no manifest hash or separate input-profile file. Dependency edges are
already represented in source bytes. Changing any reachable file or inclusion
invalidates the revision. Unreferenced neighbors are excluded. Moving the whole
checkout preserves revision; moving a member changes provenance/revision but not
its declared semantic identity. Build metadata inventories may be generated from
this graph, but are disposable outputs and never required authored input.

Before publishing current output, resolve/capture the graph again and compare
its revision, including on cache hits. Missing/newly invalid dependencies fail
freshness. Preserve the existing optimistic freshness limit: this does not create
an atomic filesystem snapshot. Keep tool/projector/renderer identity in cache keys.

`sdl check System.design`, AST, viewpoint generation, broker and SDPTool consume
the same reachable graph. SDPTool registration points to the .design entry and
checks the declared System/profile; no .design-set.json registration is needed.
0.5 retains its current single-file meaning. Explicit .design preview in 0.6
follows source dependencies; it does not require a special source-set command.
Opening a fragment without a System cannot silently guess its parent System.
A future fragment-preview context must explicitly identify the entry if desired.

For formatting, retain stdout-only behavior: normal single-source formatting emits
text; multi-file formatting requires an explicit `--file-map` mode that emits a
versioned map of formatted files without rewriting sources. Without this mode,
reject multi-file formatting with guidance instead of changing stdout shape based
on the graph. This option is proposed and requires implementation acceptance.
No new action/class/SDUI source inclusion or native XFMD implementation is implied.

## 5. Replacement acceptance cases

SSD1's semantic/provenance/consumer cases remain useful. Replace all manifest-only
cases and the no-includes assumptions with these required future checks:

| Case | Required result |
| --- | --- |
| Entry contains Containers/MachineService | Load exact file; resolve its declared MachineService; one typed System membership fact; original relation/target spans. |
| Expanded includes + symbolic contains | Same linked identities/facts as shorthand; different authored bytes/revision are expected. |
| Missing path or file missing the expected declaration | Error at original path operand with include chain; no fabricated declaration or basename fallback. |
| A and B both include Shared | One Shared AST/declaration set, two retained source dependency edges; no guards required. |
| A includes B includes A | Traversal terminates; validate after collection; no duplicate insertion. Invalid model facts still fail. |
| Distinct files declare the same name | Error at both definitions, including when byte-identical. |
| Two authored identical membership facts | Existing duplicate-fact diagnostic; load-once does not silently deduplicate independently written facts. |
| Files refer forward/across an include cycle | Resolve after collection; truly missing symbols fail. |
| System contains channel/contract versus Unit contains contract | System membership passes proposed typing; invalid Unit relation remains rejected. |
| Change a transitively included child during generation | Current delivery rejected; cache hit does not bypass full-graph freshness. |
| Add unreferenced .design beside inputs | No semantic/revision change and no implicit inclusion. |
| Include another System, exceed aggregate limits, use escaping/aliased path | Whole-candidate rejection, bounded reads, useful original-source diagnostics. |
| Format nested includes and path relations, then recompile | Same graph membership/facts, retained partition, deterministic file map. |
| Check/tree/select from one registered entry | Same checked model and graph revision; no external membership list; original-source evidence. |

All are proposed tests, not passing implementation evidence. No current model is
migrated. The existing MVP1 inclusion spelling is evidence of intent, not evidence
that its other experimental constructs become supported automatically.
