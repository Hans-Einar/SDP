# Source discovery and navigation across systems

| Field | Value |
| --- | --- |
| id | KB-SDP-046 |
| project | SDP |
| type | Change |
| CardState | backlog |
| created | 2026-10-01 |
| source | Owner report in Session 0002: SDL/SDUI absent despite source files in XFMD |

## Observed workflow gap

The owner expects discover/tree to recognize SDL and SDUI under SDP/SDL for both
single-system and multi-system projects. Current navigation.json has explicit
models and sdui arrays; empty arrays report absent even when files exist. The
new human-output layer does not change discovery semantics.

XFMD's registration has both arrays empty. Its source README identifies
SDP/SDL/XFMD/XfmdDesktop.design and SDUI/XfmdDesktop/Desktop.sdui as current inputs;
Navigation.design and Navigation.sdui remain historical reproduction fixtures.
A blind glob must not promote every .design/.sdui file into current authority.
model-map.json is local verification metadata, not a language entry manifest.

## Verified boundary

Temporary-copy registration of the two current XFMD inputs succeeds with the
unreleased SDPTool candidate: discover declares both capabilities; tree has 2454
nodes, SDL validated, SDUI available. The live XFMD project is unchanged.
Existing tests cover design-core/0.6 source composition. One compiled source graph
contains one System. Multiple registered models are supported, but tree chooses
one SDL model via --model/defaultModel; there is no aggregate systems tree.
Published engine 1.0.0 must not be credited with unreleased 0.6 support.

## Proposed next bounded work

Select a design/implementation plan for source-aware discovery and explicit
system/entry selection, then a unified multi-system navigation inventory. Respect
source-owned includes/contains, the flat System layout and ecosystem grouping;
distinguish roots, fragments, shared code, historical fixtures and SDUI entries.
Report unregistered/incompatible source candidates without claiming they are
absent or validated. Decide how existing explicit registrations interact with
convention-based discovery; avoid a new mandatory design-set manifest.

Acceptance should cover one and multiple Systems, nested ecosystem groups,
invalid/ambiguous candidates, source exclusions, stable IDs/revisions, default
selection, SDUI ownership and existing explicit registrations. Parsing/UI-runtime
integration and automatic native XFMD changes are not implied by this card.

## Provenance and handoff

[Session 0002](../../Sessions/session-%230002--SDPTool_output.md), T003, records the
owner report and read-only diagnosis. [KB045](../completed/%23045--Change--Portable-SDPTool-presentation.md)
owns output presentation only. This card is registered, not selected for execution.
The temporary evidence path is in the Session; canonical behavior is in the
[producer contract](../../../SDPTool/Contract.md).

## Owner refinement — automatic discovery, 2026-10-01

The owner rejects required manual model registration, whether edited by hand or
through a registration command. The proposed user operation is `gh sdp . discover`,
which should build the navigation inventory programmatically. This supersedes
manual registration as the desired workflow; it does not claim current code has
changed. ID, source and profile remain useful derived information.

### What the current file owns

The current producer contract and navigation.schema.json combine two concerns:
project recognition/bindings (project ID, process profile, manifest/board/plan paths,
default selection) and an explicit source inventory (models, SDUI, ownership label,
source and profile). Discover currently uses the file itself to recognize the SDP
area. Therefore deleting or overwriting it blindly would lose more than an index.
The installer currently initializes empty arrays and preserves project edits.

### Recommended design boundary

Make sources authoritative and navigation a derived result/cache. Discover should
recognize the SDP area from the supported project installation/structure, inspect
source headers and use language-owned parsing/validation to produce a deterministic
inventory for flat System folders and ecosystem/System folders. It may materialize
a generated navigation index; no user-maintained source list is required. Preserve
actual project facts in the existing project manifest/established conventions,
not another required hand-edited registration file.

Tree and select should use the same discovery service and refresh a missing/stale
index automatically, so remembering an explicit discover command is not another
prerequisite. A generated index must carry source/dependency revisions and producer
identity, be safely replaceable, and never suppress parse diagnostics by serving
old successful facts as current. In-memory use must remain possible for read-only
projects; cache location/lifetime and migration from current navigation.json are
design decisions still to be finalized. Retain installation-in-progress guards,
path containment, bounded traversal and stable client selections.

### Language version and compatibility

SDL already declares `language design-core version 0.5.` or `0.6.` in source;
SDUI declares `sdui 0.2;`. The authoritative parser exposes that information in its
AST. Read this declaration, validate with the supported profile and report the
result; do not duplicate the profile as manually maintained navigation metadata.
An unknown or invalid profile must be reported, not silently reinterpreted.

A keyword dictionary cannot establish compatibility: grammar combinations,
validation/identity rules, cross-file System constraints and semantics also matter.
AST inspection alone does not remove the need to choose a grammar first. Keep a
source's declared profile distinct from an optional compatibility assessment.
A future complies_with/minimum-profile field requires an explicitly defined and
actually executed checker for the relevant complete input graph, with diagnostic
scope and tool version. Do not infer 0.5 compliance just because a file declares
0.6 but contains no new keyword; do not rewrite its header automatically.

### Remaining source-authority question

XFMD places historical Navigation inputs beside current Desktop inputs. Both can
parse, so syntax cannot determine which is authoritative. Establish an explicit,
source-visible convention (for example current source roots versus an archive/
examples home, or supported source annotations) and report ambiguous candidates.
No such new annotation is adopted here. A 0.6 System graph can distinguish its
reachable fragments from roots; old 0.5 models and SDUI examples need a truthful
policy too. Migration must preserve historical evidence and existing links.

### Acceptance additions

- Adding/removing valid source roots updates navigation without editing a registry
  or remembering a separate registration command.
- Declared profile, validation result and unsupported-source diagnostics are
  distinguishable; compatibility is never guessed from word presence alone.
- One project can expose several Systems; entry selection remains deterministic
  without silently choosing an arbitrary file or flattening all Systems together.
- Existing project identity/bindings and legacy registration references receive an
  explicit migration, not silent loss under a generated-file overwrite.

This is recorded owner direction plus architectural recommendations. Product
implementation remains unstarted; a bounded plan and new execution Session should
own that successor rather than reopening completed presentation PLAN-SDP-0013.
