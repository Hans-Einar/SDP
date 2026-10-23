# Combined SDUI preview and launch readiness for XFMD

| Field | Value |
| --- | --- |
| id | KB-SDUI-005 |
| project | SDUI |
| type | Proposal |
| CardState | gate-review |
| Systems | SDUI, SDPTOOL |
| created | 2026-10-06T21:50:48.729784+00:00 |
| source | Owner XFMD Session 0003 T001, 2026-10-06 |
| next_review | Review delivered combined preview/readiness and XFMD U4-001 evidence |

## Request and verified baseline

XFMD should show one generated page containing a composition diagram and the
existing text-based SDUI render, plus an explicit Fyne launch action when usable.
Installed gh sdp sdui-preview was exercised on XFMD's EditorWorkspace and already
returns presentation.Markdown output (layout dump + static content). This is a
refinement of an existing facade, not a request for a second parser/text renderer.

## Producer work to select

- Extend SDUI presentation with source-structure projection and compose it with
  the existing text output in an SDPTool generated-view bundle. Preserve static
  output labels; no callbacks execute while browsing. Recommend grouped frame/
  region/ordered-row composition, with bindings/reuse distinct from containment.
- Preserve source declaration IDs/spans and normalized instance identity, reuse,
  hidden nodes, overrides and unsupported constructs. Explain differences between
  source structure and visible text render. Return revision-aware source links
  that XFMD can open and highlight exactly; define position units and missing/
  stale provenance explicitly. Do not guess spans after normalization.
- Define discoverable entry/host/profile and binding readiness plus a typed launch
  contract. A valid frame is not proof of runnable SDL-connected behavior. Distinguish
  static preview, runnable prototype and connected runtime, with unavailable reasons.
  Reuse KB-SDUI-003 preflight/capability work rather than implementing it twice.
- Fyne normally owns a separate window and may use Fyne-backed controls. XFMD's
  historical custom-widget question remains an upstream boundary decision, not a
  requirement to replace working controls. Embedding inside FOX is an optional
  feasibility experiment, with no supported bridge established by this capture.
- Reuse existing SDL runtime/binding contracts and let the caller own generated
  resources. Do not introduce arbitrary total-entry caps; coordinate existing
  producer size/count limits with XFMD #025.

## Related cards and acceptance

[KB-SDUI-003](../completed/%23003--SDUI--Proposal--Capabilities-and-navigation-pilot.md) retains
capabilities/pilot ownership; [KB-SDUI-004](../backlog/%23004--SDUI--Bug--Text-and-Markdown-fidelity.md)
retains text fidelity. Neither is superseded. KB-SDP-017 retains the facade's
broader work. This card owns the combined projection and consumer readiness handoff.

Verify multiple/reused frames, hidden instances, source spans, deterministic
projection, text-output fidelity, missing host/invalid profile/unbound callbacks,
explicit launch/lifecycle failures and actual native consumer compatibility.
No parser duplication, new widget language or application parity is selected.
The owner authorizes backlog coordination here, not implementation/publication.

External consumer: KB-XFMD-027 and toolbar KB-XFMD-029 in xfmd-sdl-navigation.
The substantive cross-project turn is recorded in external SESSION-XFMD-0003;
no activity is attributed to SDP's unrelated current governance Session.

file:///home/warloc/git/xfmd-sdl-navigation/SDP/Sessions/session-%230003--SDL-navigation-and-SDUI.md

## XFMD 0.7 producer execution — 2026-10-07

Owner XFMD Session 0004 T002 now authorizes completing the SDUI integration.
The earlier capture-only boundary is superseded for the bounded producer work
in SDP/05--Implementation/SDUI/XFMD-Preview/Plan.md, under PLAN-XFMD-0010.
No SDPTool/governance release, full capability pilot or text-fidelity rewrite.

## Delivered in XFMD 0.7 — 2026-10-07

Combined mindmaps/text with UTF-8 source provenance, revision-bound readiness
and labelled local Fyne prototype are implemented. Parser grammar is unchanged;
SDUI race and facade tests pass; real XFMD/Fyne/PDF inspection passes. Private
helpers are installed with XFMD; the owner's gh extension is unchanged. No SDL
callback execution, embedding, full capability pilot or fidelity completion.
Card moves to gate-review; no independent approval or upstream release claimed.

file:///home/warloc/git/xfmd-sdl-navigation/SDP/05--Implementation/SDUI/evidence/U4-001/README.md
