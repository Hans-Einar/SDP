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
absent or validated. Replace explicit registrations with source/directory discovery; migrate existing
non-source project facts without introducing a replacement registration manifest.

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

### Selected direction — remove navigation.json

Owner correction, 2026-10-01: remove navigation.json as the navigation/discovery
contract, not merely replace manual editing with a generated version of the same
required file. Sources and directories define what can be browsed. This supersedes
the earlier recommendation to materialize a generated project navigation index.
Current code still uses registration; removal has not been implemented.

Discover inspects the selected project's SDP/source structure and obtains language
metadata from source headers and validated ASTs. Tree and select use the same
service directly. No registry file, registration command, prior discover run,
manual current-model list or replacement sidecar is required for navigation.
A file that exists in the browsable source tree must not disappear because it is
absent from a second list. Report parse/validation errors on affected entries;
one invalid source must not make unrelated sources invisible.

Directory navigation and semantic navigation are complementary views of the same
sources: actual files/folders provide the browsing hierarchy; language-owned
System/includes/contains relations provide semantic structure and entry context.
Multiple roots can be shown for selection rather than silently electing one.
Internal, disposable in-memory indexing may avoid repeated work; it is not a new
project document or authority and must not change fresh discovery results.
Persistent cache design is not required by this delivery and is not selected.

Actual non-source project facts already captured by installation/project manifests
retain those owners. Assess old default/plan/board/identity bindings during
migration; do not retain navigation.json merely to preserve the old arrangement,
and do not relocate the same manually maintained source list into another file.
Retain installation-in-progress guards, path containment, bounded traversal,
source revisions and stable client selections in the new implementation.

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

### Correction — historical files do not block discovery

The earlier recommendation treated distinguishing current versus historical XFMD
files as a prerequisite for navigation. The owner rejects that prerequisite.
If Desktop and Navigation sources are present in the same browsable directory,
show both. Their successful parsing does not establish that either is current,
but navigation need not establish currentness to expose the files.

Archival organization is an ordinary source-housekeeping decision when appropriate,
not an additional register/annotation users must maintain for discovery to work.
Do not hide or move XFMD files from this workstream. Do not invent a required
current/historical source keyword. When semantic compilation needs a root, use
source-declared composition and the user's selection; report ambiguity rather
than inventing a single authoritative project-wide root. Navigation must remain
available while such a semantic selection is unresolved.

### Acceptance additions

- Adding/removing valid source roots updates navigation without editing a registry
  or remembering a separate registration command.
- Declared profile, validation result and unsupported-source diagnostics are
  distinguishable; compatibility is never guessed from word presence alone.
- One project can expose several Systems; entry selection remains deterministic
  without silently choosing an arbitrary file or flattening all Systems together.
- Discovery/navigation works without navigation.json, and an old copy cannot hide
  newly added sources or override their declared language profile.
- Two independently parseable files in one directory are both browsable without
  current/historical labels, registry entries or an elected default root.
- Existing project identity/bindings and legacy registration references receive an
  explicit migration; no replacement manually maintained source list is introduced.

This is recorded owner direction plus architectural recommendations. Product
implementation remains unstarted; a bounded plan and new execution Session should
own that successor rather than reopening completed presentation PLAN-SDP-0013.

## Selected discovery and viewer ownership — 2026-10-01

Owner decision: discover inspects the SDP area and returns discovered information.
It is a finite request/response operation, not a watcher or daemon, and does not
write a project navigation.json. Its scope is the SDP area, including supported
process/navigation information as well as SDL/SDUI sources; it is not only a source
file search. Discovery reports observations/diagnostics, not whole-project proof.

A viewer requests `gh sdp . discover --json` (explicit machine mode under the new
output contract), keeps the returned navigation data in its own memory buffer,
and presents that result. It owns filesystem watching beneath the selected SDP
area and issues another discovery request on changes or manual Refresh. No
persistent sidecar or prior registration command is required. The human command
without --json uses the existing presentation layer on the same result.

Derived/recommended refresh details for the implementation plan:

- Watch creation, deletion, rename and content changes, including new directories.
  Coalesce bursts of saves rather than starting a discovery for each event.
- A change during an in-flight request schedules a subsequent refresh. Track
  request/project identity so a late response cannot overwrite a newer buffer or
  a different project's navigation. Buffer replacement is performed as one update.
- A failed refresh remains visible as an error; a retained prior buffer must be
  marked stale. Per-entry source diagnostics remain browsable rather than hiding
  unrelated entries. A filesystem scan is not claimed to be an atomic snapshot.
- Manual Refresh works independently of watcher availability. Rendering output
  must not create a self-sustaining refresh loop; select output locations and
  event filtering deliberately without hiding authored-source changes.

These mechanics are implementation recommendations derived from the selected
ownership, not claims that watcher behavior exists. Native XFMD integration belongs
to its own workstream; SDPTool owns the bounded discovery response and its tests.
