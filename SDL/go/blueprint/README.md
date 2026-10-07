# Structural blueprint analysis

BPI1 provides an I/O-free Go library, not a published blueprint command or an
assignment execution engine. Call Analyze with two validated sourcegraph Snapshots,
optional expected revisions, an authored Task and explicit resource limits.

The first policy accepts design-core/0.6 for the same System. It covers the closed
parser statement families: Relation, Dependency, Allocation, PropertyAssignment,
Projection, Placement, Participation and Step. Relation verbs have an explicit
allowlist. Policy v1 conservatively traverses them in both directions, stopping
propagation through System/mode labels. Selected contracts retain modeled payload,
fields, properties and participant context. Unknown future semantics fail closed.

Node IDs are System/kind/name. Fact identity includes canonical sentence and typed
endpoint identities; kind changes and renames are remove/add. Facts expose the typed
parser Statement with source spans removed from semantic data. NOW/TARGET origins
are retained separately, including duplicates. Relative source moves/formatting
change snapshot identity without becoming semantic changes.

## Task and result

Task ID and intent are required. Context and exclusions contain full node IDs.
Excluding any required conservative closure is a scope conflict. Each added/removed
declaration or fact needs an exact Permission to pass constraints. AllowedModelPaths
optionally narrows those permissions to exact relative source paths. AllowedCodePaths
is retained in task identity and validated for safe relative syntax; BPI1 does not
inspect code or enforce a filesystem sandbox.

Rules assert the presence/absence of a canonical fact separately in NOW and TARGET.
They require ID/provenance. Unknown assertions fail rather than silently pass.
Protection explicitly chooses boundary or subtree scope rooted in a node. Membership
follows owns/contains descendants in the union graph. Subtree protects all declarations
and incident facts. Boundary protects the root declaration and crossing facts,
allowing purely internal restructuring. It is a modeled structural boundary,
not a runtime/API compatibility proof. Exact fact Rules can add further constraints.
Conflicting permissions/protections are errors.

Analysis is always a diagnostic model view, never permission to execute an assignment.
ConstraintsPass describes only the authored structural constraints. Coverage identifies
modeled closure, unsupported SDUI semantics and unknown code/runtime conformance.
Missing channel peers are checked per side, including expected participation groups
from the other side and permitted messages with no participants.

Semantic data is deterministic and JSON-marshalable. Identity binds source revisions,
task values, limits and producer/policy versions. Lifecycle artifact UUIDs and complete
captured source inventories remain BPI2 adapter responsibilities; a graph revision
does not substitute for a full ModelGovernance artifact digest.

## Verification

Run from SDL/go:

    go test -race ./parser ./sourcegraph ./blueprint
    go vet ./blueprint

Tests use the real reduced MVP1 fixtures and parser/sourcegraph. Typed operand policy
tests are labeled separately. No Markdown publication, catalogue, assignment state,
GUI, behavior execution or Ponsse implementation is delivered here.
