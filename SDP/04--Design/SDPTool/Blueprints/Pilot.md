# BP1 worked assignment — unsaved design preview

This is an authored study specimen, not an authorized implementation assignment.
Its model extracts come from the existing SDL toolkit. Narrative intent, contract
selection, code mapping and the protected-owner assertion below were authored for
this study and must not be presented as automatically inferred SDL facts.

[Study](Study.md) · [Generated extracts](Model-extracts.md) · [Evidence](Evidence.json)

## Assignment entry

| Field | Study value |
| --- | --- |
| Work anchor | UnsavedSourcePreview / AcceptSourceSnapshot |
| Requirement | REQ-SDPTOOL-004, preserve source and session identity |
| Current design | SDP/03--Architecture/SDPTool.design at main 2ff71c4a033aa1393e02657777e85f5fb67ccefd |
| Model revision | 44adae9604a5f555c3ac4e00eb77c5cc2a28da2979ed7e84d3e411252a0b9a35 |
| Proposed target revision | f5f88b57d67fd722bcf54c8ded0f46f905f749b5052780ed0f36706b0b3db8bf |
| Existing behavior | SDPTool reads a saved file; checks expected revision and changes during generation; publishes a complete caller-owned bundle |
| Desired behavior to design | Accept an immutable editor snapshot with source identity and revision, using the same validation/projection/publishing behavior; retain saved-file compatibility |
| Non-goals | XFMD implementation, callbacks/domain execution, source-set imports, new renderer, general assignment engine |
| Evidence status | Parser/projection experiment and current saved-file regression tests only; no unsaved preview implementation or worker trial |

## Required context before implementation

The system purpose is to preview supported design source without executing it and
to preserve model provenance. DesignAuthor and DocumentViewer pursue PreviewDesign;
SdpTool supports it. PreviewCoordinator owns AcceptSourceSnapshot and
BuildDesignPreview and consumes SdlProjectionPort. ViewerBridge owns
DeliverViewBundle and consumes ViewerDeliveryPort. These are structural model
facts, not a generated Go call graph or proof that each planned responsibility has
an implementation. UnsavedSourcePreview depends on SavedFilePreview and remains
planned; the target does not change that status.

Protected context comes from the [producer contract](../../../../Toolkit/SDPTool/Contract.md)
and [requirements](../../../02--Requirements/SDPTool.md):

- Parse and validate the complete snapshot through SDL; no second parser.
- Keep delivery ownership with ViewerBridge. Accepting a source snapshot does not
  authorize taking over viewer routing or changing the document consumer protocol.
- Preserve the last valid result on invalid source or failed generation. Source
  bytes must not be overwritten as a shortcut for pretending an unsaved buffer is saved.
- Source revision and consumer request identity serve different purposes; the host
  must reject obsolete replies even when source hashes match. Cancellation alone
  does not prove an old result cannot arrive.
- Keep outputs distinct per request and caller-owned until resources are released.
  Source/output containment and bounded-resource rules still apply.

## Observed code and write-scope proposal

| Design responsibility | Observed implementation at baseline | Proposed treatment |
| --- | --- | --- |
| Saved source validation and bundle creation | Toolkit/SDPTool/preview.go: Preview, loadModel, readSource | Reuse and separate snapshot input from saved-file acquisition; exact API is still to design |
| Project selection | Toolkit/SDPTool/selection.go: SelectProject | Keep explicit registration and expected revision; no parent/project inference |
| Language and projections | SDL/go/parser, SDL/go/viewpoint, SDL/go/documents | Dependencies; no parser/projection semantics change in this example |
| Host request ordering and window behavior | Consumer contract, XFMD-owned implementation | Protected neighboring boundary; no XFMD edits in SDP assignment |
| Design-to-code mapping | This inspected table | Authored observation, not an SDL fact or a stable mapping schema |

A future authorized assignment should name exact allowed implementation files and
shared interfaces after API design. This specimen does not grant blanket writes
to the paths above. Changes to bundle publication or the consumer contract need an
explicit reviewed amendment; errors in the current mapping remain reportable gaps.

## Proposed model delta

The full target is generated reproducibly from the complete baseline, with one new
declaration and five statements, then canonicalized and checked by the Go parser:

```text
functionality ValidateSnapshotRevision.
PreviewCoordinator owns ValidateSnapshotRevision.
ValidateSnapshotRevision allocated-to SdpToolProcess in mode StandalonePreview.
ValidateSnapshotRevision contributes-to SdpTool.
ValidateSnapshotRevision realizes DesignPreview.
UnsavedSourcePreview addresses ValidateSnapshotRevision.
```

This fragment explains the authored proposal; it is not a standalone model or
parser input. No implementation-status changes are made. The stronger prose about
revision identity, resource lifetime and failure behavior still lives in the
contract; merely declaring ValidateSnapshotRevision does not encode these rules.

## Context selection and exclusions

Mandatory entry: assignment and protected behavior above. Mandatory model views:
VP03 around PreviewCoordinator (depth 1), and VP03 around SdpToolProcess (depth 2),
both showing DesignPreview. VP06 around UnsavedSourcePreview supplies plan context
in the evidence file. Generated diagrams contain only validated source facts.

The broader neighborhood is intentionally inspected before restricting the display
to DesignPreview. PlanAssistance and CardHistory are not changed; retain links to
the complete model instead of printing their detailed diagrams in this assignment.
SDUI rendering, packet diagrams, Ponsse domain code and class/runtime profiles are
outside this change. No missing Channel or provider is invented to draw a flow.
Snapshot request/reply details are an acknowledged contract-design gap.

These exclusions are authored judgments for this example. A future selector must
show why boundaries are safe and retain cross-boundary contracts; a graph depth
alone is not evidence of completeness.

## Drift challenge and review questions

Counterexample: a worker puts all preview behavior into PreviewCoordinator and
moves DeliverViewBundle from ViewerBridge to it. The source remains valid SDL and
might look convenient from a local task perspective. The current and target
container views plus the explicit protected-owner assertion expose the violation.

The study probe evaluates that one assertion against toolkit-generated facts:
current and target pass, ownership-drift fails. No production drift checker exists.
If the worker changes only Go and leaves SDL unchanged, this assertion passes;
code review and integrated tests are still required. A structure-only guard also
cannot prove revision checking is actually reached on every response path.

Reviewer questions:

1. Is the same complete snapshot used for validation, rendering and revision? Can
   an editor edit/save race publish a result under the wrong identity?
2. Does a failed refresh retain the prior valid preview and correctly report stale
   status, including cancellation and reordered equal-hash replies?
3. Are source provenance and resource ownership preserved without saving dirty data?
4. Does actual ownership/call flow still respect the delivery boundary? If it must
   change, was the assignment amended instead of silently weakening the invariant?
5. Are all accepted inputs, errors and consumers covered by appropriate evidence?

Future acceptance includes saved-file compatibility, unsaved invalid/valid source,
source identity mismatch, late responses, cancellation, output conflict and failed
refresh. The existing tests below are a baseline, not the unsaved-input acceptance.

## Reproduce the observed evidence

From repository root, with a Go toolchain supporting the module's declared version:

```sh
python3 SDP/04--Design/SDPTool/Blueprints/probe.py --go /path/to/go --output /tmp/bp1-fresh-output
cd Toolkit/SDPTool
go test ./... -run 'TestPreviewFreshSourceAndPreservation|TestFailedRefreshKeepsPreviousBundle|TestRelationshipReferencesAndSelectionRevision' -count=1
```

The output directory must not exist and must be outside the repository. The probe
requires the inspected parser/projection/document/CLI files to match the pinned
baseline; it intentionally refuses a changed implementation rather than attributing
new output to the old baseline. It reconstructs three complete models and generates
nine selected bundles in temporary storage. Only compact JSON evidence and six
Markdown diagram extracts are retained here; model snapshots are reproducible from
Git and the probe. No renderer, desktop GUI, daemon or network service is needed.

Observed with Go 1.27.1: all three models validate, the ownership guard rejects only
the drift model, and the three named existing regression tests pass. An initial
probe draft appended a declaration after statements and correctly failed parsing;
the final probe puts declarations first and uses the canonicalizer. Two clean final
runs are compared byte-for-byte for the retained evidence and Markdown. These are
research checks, not a claim of complete blueprint generation or drift prevention.
