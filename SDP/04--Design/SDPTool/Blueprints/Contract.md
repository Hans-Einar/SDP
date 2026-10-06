# BP2-A — semantic blueprint contract candidate

2026-10-06, PLAN-SDP-0001 / KB050 / Session0008. This is a proposed producer
contract and executable input pilot, not a shipped blueprint command. The owner
selected continuing the blueprint work; the particular calibration extraction is
an agent-selected design specimen, not an approved Ponsse implementation change.

## Purpose and input authority

A blueprint makes an intended change understandable in its affected system context.
Use two complete, parser-validated source snapshots: NOW and TARGET. NOW means the
selected baseline, not automatically the implemented truth. Both identify artifact
UUID/head, source digest, entrypoint, language profile and compiler identity. A WORK
TARGET is permitted and labeled preliminary; generating a view never freezes it.
The same captured bytes must drive analysis, validation and output. Later edits do
not silently alter a retained bundle; stale inputs require explicit regeneration.

The authored assignment supplies intent, permitted edits, preserved obligations and
named evidence. Generated facts supply declarations, relations, changes, provenance
and selected context. Observed code/test evidence is a third category. Never infer
an owner requirement, code mapping or implemented status merely from an SDL edge.
No new navigation registry, central VCS or source-registration file is introduced.

## Pilot and protected context

Use the six-file NOW/TARGET specimen in experiments/blueprint_mvp1. It is a reduced
translation from the experimental MVP1 model, with explicit omissions in its README.
The production parser does not accept sdl-mvp1-exercise/0.1. Rewriting only its header
would falsely claim equivalent semantics. This pilot instead uses design-core/0.6.

Intent: move CalibrateMeasurement ownership into CalibrationProcessor, a Unit inside
MachineService, while retaining allocation to MachineService in Inspection mode.
The modeled contract with BuckingWeb remains byte-for-byte unchanged. Do not move
ReduceMachineState or PublishMachineBaseline into the web/UI containers. Preserve
BuckingWeb's adaptation and BuckingUI's projection ownership. Calibration fallback
(raw-only/unknown), sequencing and actual protobuf behavior remain prose obligations
from the experiment; this structural slice does not prove them.

The model has no internal call edge from calibration to reduction. The producer must
expose that missing behavioral dependency as unknown rather than claiming it found
all runtime impact. BuckingUI is retained as authored contextual scope even though
this reduced graph lacks its browser channel. Code paths/tests for the proposed
extraction are deliberately unmapped; no fabricated Go assignment is generated.

## Bundle surface and annotations

The proposed bundle contains complete NOW/TARGET source copies, input manifest,
assignment narrative, semantic change inventory, context diagrams, obligations,
verification references and limitations. One authored task is authoritative;
generated files identify their origins and are not a second task specification.

| Annotation | Meaning | Authority |
| --- | --- | --- |
| CHANGE | Permitted design/code work boundary | Authored intent, checked against model delta |
| PRESERVE | Relation, contract or behavior must still hold | Authored rule plus explicit verification method |
| CONTEXT | Included to explain a boundary or dependency | Selection provenance or authored inclusion |
| UNKNOWN | Missing mapping, provider, dependency or evidence | Observed gap; never implicitly satisfied |
| ADDED / REMOVED / MODIFIED | Model fact difference | Generated from checked NOW/TARGET facts |

Use labels as well as color. Deleted edges/nodes must remain visible from NOW.
A preserved boundary can surround changed internals; do not mark the entire enclosing
container immutable when the task explicitly permits work inside it. Source renames
are initially remove/add unless explicitly mapped; no heuristic name similarity
may silently establish identity. Display-only source formatting differences are
separate from semantic changes, while still changing source provenance hashes.

## Verification semantics

Baseline checks are pinned to NOW and the corresponding code revision when one is
known; target checks are pinned to TARGET and the candidate code revision. A check
record states command, inputs, execution status, output/evidence digest and scope.
Missing code correspondence is UNKNOWN. A parser pass is structural evidence only.
The pilot's model-only release must never be labeled implemented.

Negative specimen for the next phase: move ReduceMachineState to BuckingWeb. It can
remain syntactically valid, yet violates a PRESERVE rule. Another counterexample
changes only Go code while leaving both models untouched; semantic diff alone cannot
detect that drift. Independent code review and mapped behavioral checks remain needed.

## Next design milestone

BP2-B specifies deterministic selection over the NOW/TARGET union, relation-specific
traversal, explicit expansion boundaries, source-linked reasons, removed-edge
neighbors, cycles, missing providers, overflow and stale-input failures. The producer
must not turn an arbitrary radius into a proof of complete impact. BP2-C defines
library/facade ownership and a bounded implementation plan, then an authorized
worker/reviewer trial. General code-tag mapping remains KB004; channel execution
remains KB-SDL-006. No product generator is implemented by this document.
