# BPI2-M1 — diagnostic publication evidence

Execution: 2026-10-07/08. Baseline 84497bf plus this milestone's recorded files.
Branch sdp/blueprint-implementation; exact reviewed/compiled identities and observed
pilot output are retained in Pilot-BPI2-M1.json. No release or main integration.

## Delivery

The compiled SDPTool command captures NOW/TARGET artifacts, parses reachable SDL
from owned bytes and produces Markdown/Mermaid, typed JSON and complete captured
source trees. Source links resolve within the bundle. Entry, compiler/executable,
artifact metadata, source inventories, task bytes and analysis/policy identify the
output. Sources are rechecked immediately before publication. All bundles remain
diagnostic previews; failed obligations are visible, not assignment permission.

Task input is strict JSON (assignment.json), replacing the design's proposed YAML.
The common Go human/JSON presenter and owned-file publisher are reused.
The existing source model now includes BlueprintFacade and diagnostic responsibilities;
its canonical design-core/0.5 check passes. It adds no new deployed container.

## Verification

Targeted generation/compiled CLI tests pass without Git or SDP installation.
They cover repeatable payloads, source/task edits before publish, cancellation,
invalid SDL, source/output overlap including symlink aliases, literal authored text,
unchanged unmanaged notes, edited generated files and explicit provenance.
Shared publisher tests cover staged write rejection, failure of either rename,
successful rollback, and retained backup diagnostics when rollback itself fails.
No physical power-loss or arbitrary-editor atomicity guarantee is claimed.

SDL blueprint/documents/sourcegraph/parser race suites pass. Final focused tests
and vet were repeated after fixes. Broad SDPTool suite results are recorded at
closeout below; tests begun before review fixes are distinguished from the focused
final-candidate reruns.

A real compiled CLI generated the reduced MVP1 calibration blueprint, including
preserved reduction ownership. The selected changes.mmd rendered successfully with
the existing mmdr into SVG, rasterized with rsvg-convert and visually inspected:
change labels and all six nodes are readable. This is renderer compatibility
evidence, not native XFMD GUI acceptance or a claim of polished layout.
The full context graph rendering was stopped after several minutes during concurrent
machine load; it is retained as Mermaid but is not claimed visually verified.

## Independent review

Fresh read-only agent /root/bpi2_review reviewed the producer contract and candidate.
Initial changes-required findings:
1. Symlinked model-area alias could bypass lexical overlap and write into sources.
2. Explicit entrypoint and parser/compiler identity were absent from output.
A later finding identified missing PRESERVE labels for slash-containing protection IDs.
All were fixed and regression-tested. Final verdict: approved for bounded BPI2-M1.
The reviewer confirmed the original exploit no longer writes into sources. Parenthesis
filename links were investigated and passed without an unnecessary escaping patch.

## Remaining scope

BPI2-M2 owns immutable retained catalogue revisions, discovery/tree projection and
the consumer fixture. BPI3 owns readiness, assignment records and state groups.
No assignment lifecycle, full behavior evidence, external project change or release
is delivered by this milestone.


## Broad-suite result

SDPTool root, blueprints, model, presentation and tools/profile passed the broad
race run. The unchanged installer package exceeded the default ten-minute limit
in TestProcessExitRecoveryMatrix; this is an incomplete installer race check,
not a full-suite pass or a proven functional regression. Raw output is retained
in Tests-SDPTool-BPI2-M1.txt. Concurrent machine load was high; the timeout's cause
has not been independently isolated. No installer code was changed. Its race
verification remains outstanding for final integration. The focused final-candidate
blueprint tests, source-symlink regression, CLI tests and shared publisher failure
tests passed. Tests-SDL-BPI2-M1.txt records the SDL package results.
