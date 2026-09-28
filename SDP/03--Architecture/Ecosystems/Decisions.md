# Ecosystem and system boundary decisions

Authority: owner request 2026-09-29; execution [PLAN-SDP-0008](Plan.md).
This is the selected source organization and an initial architecture proposal,
not approval or implementation of every proposed binary/API.

## E01 — Ecosystems group systems

Use SDP/SDL/ProjectGovernance, SDP/SDL/SDL and SDP/SDL/SDUI. Each contains
system-owned source folders. No Ecosystems wrapper is needed: SDP/SDL already
identifies the model source root. Keeping SDL and SDUI as three coarse products
would hide the independently useful parsers, compilers and service entry points
the owner wants to describe. Treating every package as a system would instead
confuse implementation organization with product boundaries.

A system is a coherent public tool/service responsibility that may expose a
binary and library. A container is an actual or explicitly proposed runtime
boundary inside a system. A library remains a unit with one owner, usable by
other systems without creating a separate daemon. Runtime type does not imply
one OS process per model in the current implementation. Ecosystem membership
is index metadata, not a `contains` relation between deployment containers.

## E02 — Current products and target boundaries differ

Current `sdl` combines structural parsing/checking, action checking, class
checking and viewpoints. Current `sdui` combines parsing and static presentation.
Splitting their design responsibilities into systems does not claim extraction
or new installed executables. System READMEs map the actual commands/packages;
future standalone binary names are proposals. Shared parsers must not be copied.

SDPTool remains the project facade and owner of proposed shared process rules.
The routine engine is a responsibility inside its core, with persistent state
and versioned routines; MCP and the app-server client do not own rival ledgers
or transition rules. KanBan CLI exists today; KanBanTUI is a proposed read-only
consumer first. Keep feature maturity explicit rather than inferring it from
valid SDL syntax. The worker/reviewer role distinction is not an OS boundary.

## E03 — Real parser support, bounded models

The tracked Go parser implements design-core 0.5 single-file models. It does not
implement System declarations, cross-file imports or ecosystem linking. Model
files therefore contain actual supported containers, units, capabilities,
interfaces and behavior. The enclosing directory/README names their System.
A dependency on another system is an interface in the local model, with a
README/catalog link to its owner; there is no hidden global symbol resolver.
Use one compact entry per system now. Do not introduce a new preprocessor,
concatenate fragments, or adopt the unrelated untracked sourceinput draft.

These entries are capability/boundary proposals grounded in inspected code and
RGS2 studies. Existing large models retain detailed/historical authority until
KB-SDP-020 and KB-SDL-005 can migrate facts, locations and consumers together.
Conflicting historical planned statuses are not proof of current implementation.
The new source maps state observed implementation and deliberately deferred work.

## E04 — Runnable navigation now, hierarchy later

Register complete new models in SDP/navigation.json using existing schema 1.0.
Preserve old IDs/default and SDUI registrations. IDs for new boundary models are
separate from `sdptool` and `sdl-sdui`, whose existing consumers keep their meaning.
The catalog groups them by ecosystem for humans; SDPTool can select each model
and produce its tree and Markdown on demand. Do not claim the current consumer
has native ecosystem groups merely because the directory does.

A reproducible catalog verifier invokes the real Go tools to check sources,
produce AST and generate viewpoint packages in an explicitly selected output
folder. Generated output is disposable and carries source/tool identity.
The index may link to generated tool output, but never substitute hand-authored
Mermaid for SDL-derived diagrams.

## E05 — External systems and distributed templates

Codex app-server, XFMD, gh-sdp, Mermaid rendering and Ponsse remain collaborators,
not newly owned products under this source tree. The models expose their ports
and limitations without requiring their worktrees to exist. Template scripts,
release signing and CLI helpers are mapped as support tooling; a folder or
script alone is not a new software system. Template/sdp-root remains a released
installation concern. This local structure does not silently migrate consumers
or change installed project conventions.

## Next refinement

Use these models to choose the first bounded routine-governance DesignPlan and
the source-set input contract. Register exact cross-system contract identities,
version compatibility and traceability queries when those mechanisms are selected.
The catalog is useful now for navigation and critique; parsing does not prove
execution, complete threat modeling, API compatibility or blueprint correctness.
