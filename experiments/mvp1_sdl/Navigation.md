# MVP1 navigation handoff

Goal: open the MVP1 System in XFMD, browse its model facts and generate selected
documents on demand. This pilot reorganizes authored inputs and adds real SDUI
0.2 sources; it does not yet make the experimental SDL corpus a supported model.

## Available now

- [System source](SDL/MVP1/System.design) explicitly lists all 67 companion design
  files; [shared guidance](SDL/MVP1/Shared/README.md) explains component placement.
- Three SDUI sources are registered in this repository's SDP/navigation.json:
  mvp1-operator, mvp1-apt and mvp1-simulator. SDPTool can discover their frames and
  generate structural Markdown on demand through the released sdui-preview service.
- [The SVG gallery](preview/index.md) is a small generated visual review snapshot.
  SVG layout comes from the Go SDUI exporter; the released SDPTool facade currently
  exposes structural Markdown, not this SVG presentation.

From this repository root, with the current gh-sdp extension installed:

```sh
gh sdp . tree
gh sdp . sdui-preview --model mvp1-operator --entry page --output /tmp/mvp1-operator-review
```

The consumer should pass the source revision returned by tree using --revision,
keep the result bundle while displayed, and discard obsolete asynchronous replies.
The command above is a manual invocation without an expected revision. The
[SDPTool contract](../../SDPTool/Contract.md) is the authoritative consumer API.
No XFMD application code or Ponsse installation changes are part of this pilot.

## Required before full SDL navigation

1. KB-SDL-005: define/implement one declared System and explicit source sets in
   the existing Go frontend, preserving per-file diagnostics and stable identities.
2. Resolve which constructs from sdl-mvp1-exercise/0.1 enter a supported profile
   or receive an explicit semantics-preserving migration. The corpus includes
   layers, state, governance and behavioral descriptions beyond design-core/0.5;
   merely loading 68 files is insufficient. Keep unsupported facts visible as gaps.
3. Register the complete validated MVP1 model with an aggregate revision over
   all declared inputs. Changing any participating file must invalidate affected
   navigation/generation; no stale single-entry-file hash.
4. Expose model-derived groups for Features, Scenarios, Containers, Contracts,
   Libraries, state and supported governance facts. Container nodes lead to their
   capabilities, owned Units/Functionalities, channel participation and related
   scenarios. Shared source folders must not become invented runtime containers.
5. Bind the separate SDUI sources to their application owners as navigation
   metadata with defined semantics; do not infer ownership from directory names.
   Generate selected document/resource bundles using existing SDPTool services.

Keep the visual folder tree and semantic model tree distinct. Cross-container
Features/Scenarios reference identities instead of duplicating declarations.
A viewer may offer source browsing as well as viewpoints, with clear source versus
projection labels. No handwritten substitute viewpoint or silently reduced aggregate
model should be used as evidence that full MVP1 navigation is implemented.

## First integrated acceptance path

Choose BuckingUI → OperatorInteraction → InspectMachine → MachineService and its
channel contract; generate that detail, open its defining source, then open the
operator SDUI screen. Repeat for SimulatorUI → SimulatorWireAction without leaking
simulator truth into the operator context. Edit an included design source and
verify invalidation; a bad edit must retain the last valid displayed document
with a diagnostic. This provides practical feedback before expanding the entire
SDP process. KB-SDP-032 owns owner feedback; KB-SDP-020 owns wider model migration.
