# SDL — navigable viewpoints and on-demand generation

**Original design:** G6-D1/D2, 2026-09-22. **Status clarification, 2026-09-24:** G6-M1–M6 are implemented; G7 launches prebuilt tools. Use [Go commands](../../go/README.md), [G6 evidence](../../go/evidence/G6.md) and the [launcher](../../scripts/README.md) for actual behavior, URI/IPC and configuration. Normal browsing requires no daemon; an optional broker is delivered. The XFMD consumer belongs to separate phase work identified in G6 evidence; arbitrary installed XFMD versions may lack its link handler.

The remainder preserves design foundations, including proposed directories/protocol sketches, not exact current CLI/wire contracts. The [shared plan](../../../SDUI/docs/implementation-plan.md) owns milestones; the [SDL model](../../../SDUI/design/architecture.design) owns responsibilities/scenarios.

## One source, multiple documents

**G6-D2 clarification:** [levels/notation design](SDL-Viewpoint-Levels-and-Notation.md) recommends navigator/overview without prerendered details for development. Static packages are an equivalent export choice. A0–A5/type inventories organize menus; G6-M5/M6 plan consistent symbols and later class diagrams.

SDL source → validated model → selected viewpoint/projection → document sections → Markdown with Mermaid/SVG. Full reports, single pages and navigators share projectors/source maps. Navigators describe selections rather than scrape headings from large generated Markdown files. A future content model with stable section IDs can reuse sections across outputs.

Proposed static package, not current output format:

```text
design-docs/
  index.md
  navigator.md
  implementation.md
  viewpoints/
    index.md
    VP02/
      index.md
      architecture.md
    VP06/
      index.md
      G1FrontendPort.md
      G6NavigableDocumentation.md
    VP08/
      index.md
      SelectedViewOpened.md
  assets/
    <content-hash>.svg
  manifest.json
```

Pages include titles, back/parent links, source revision and related objects. Stable IDs, not heading text, determine anchors/filenames. Escape Markdown/URI names and validate before using them as paths. Image paths are relative to referring pages. Manifests own generated files; preserve user notes on re-export.

Combined viewpoints.md/printout.md is optional, not produced at every click. Ordinary relative links work without an SDL service in static export; navigator-only creates no detail files.

## Navigator links and host adapters

Candidate notation from the design stage, not standard Markdown behavior or an exact implemented URI contract:

```markdown
[G1 — frontend](sdl-view://sdui-design/VP06?focus=G1FrontendPort&target=main&consumer=xfmd)
[Static G1 page](../../../docs/viewpoints/VP06/G1FrontendPort.md)
```

XFMD translates the first link into typed ViewSelection. Resolve project IDs through configured projects; validate viewpoint, focus, optional mode and profile against the selected revision. Invalid selections report diagnostics and retain the displayed document.

`consumer=xfmd` selects a locally registered adapter with fixed executable/arguments, not shell or arbitrary post-command. This permits selecting readers under a common contract. Ordinary readers use static links. Final URI/CLI forms are established with XFMD implementation; these are proposals.

The host adds request ID, session, exact window and target panel. `target=main` preserves the side navigator; future `target=navigation` may replace navigation pages. History/scroll belong to each panel. G6 handles documents; Fyne remains the first executable-widget host.

## From click to document

1. XfmdDocumentHost captures window/panel and sends SelectViewRequest.
2. DocumentBroker validates selection and chooses a consistent source snapshot.
3. SdlViewpointGenerator builds only selected documents/resources.
4. ViewArtifactStore publishes a complete package and supplies reference/lease.
5. ViewerLaunchAdapter sends the reference to the selected reader.
6. Reader acknowledges opening; broker answers the original request.

SelectedViewOpened models this sequence for generated VP08. InvalidViewSelectionRejected and ViewProjectionFailed show pre-open errors. Test publication failure, closed windows and cancellation too; scenarios do not claim complete failure coverage.

Proposed host boundary: `OpenDocument(requestId, entry, revision, lease, windowId, pane)`. SDL contract fields describe argument boundaries, not executed field/transport checks. Implementation defines Go types/IPC version. Projectors import neither XFMD nor process launchers; broker/launch adapter coordinate.

`xfmd --active <entry>` was desired behavior. The inspected `src/application/main.cpp` then documented only file/directory, `--help` and `--version`; new window/IPC support required XFMD changes. Capture target on request receipt; navigator clicks use sender-window ID. Later focus changes must not redirect results. Proposed adapter: `--window <id> --pane main`. Closed targets fail the same request; opening new windows is explicit host policy. Panel clicks need no global focus detection.

## Memory, publication and lifetime

A daemon can reuse ASTs, projections and resources but is unnecessary for initial single requests. G6-M2 may use one CLI process; G6-M4 adds user-owned local IPC service.

Naming a Go-memory object does not create a readable file. File-based readers need actual packages. Initially publish small on-demand packages under `$XDG_RUNTIME_DIR/sdl/views/<session>/<revision>/<selection>/` with entry.md, SVGs and manifest. These are temporary files, not full repository exports. Pure memory transfer needs a separate streaming/IPC consumer and is deferred; no FUSE/pseudo-file required.

[XDG](https://specifications.freedesktop.org/basedir/latest/) specifies user-owned runtime directories, mode 0700 and session lifetime, not guaranteed RAM storage. Do not use global /run/SDL. If unavailable, explicitly fall back to private temporary storage and report location. Bound package/total-cache sizes; export large archives elsewhere explicitly.

Publish immutable revision directories atomically before OpenDocument, including required relative images/link targets. Never display partial files. Cache keys include source/import hashes, profile, projector version, selection and renderer/theme. Source changes invalidate affected results; unknown dependencies require full invalidation.

Reader leases retain packages while resources may load/reload; open acknowledgment alone is insufficient for deletion. Release on document switch/close. Generic readers without leases retain packages until explicit cleanup/session end. At capacity reject requests rather than remove visible resources. Test timeouts, crashes, quotas and cleanup in G6-M4.

A newer click supersedes older results for the same window/panel. Check request ID/revision before opening. Source/render failures retain the last valid view and report diagnostics. Navigator/main documents identify their revisions. Navigation executes no SDL domain actions.

## Responsibilities and acceptance

| Boundary | Owner |
| --- | --- |
| Viewpoint/section selection, Markdown, navigation, source maps | SDL tool / SdlViewpointGenerator |
| Revisions, queue, stale-result rejection, cache keys | DocumentBroker |
| Packages, atomic publication, leases, cleanup | ViewArtifactStore |
| Registered reader selection/open protocol | ViewerLaunchAdapter |
| Two Markdown panels, window/panel IDs, routing, acknowledgments | XFMD / XfmdDocumentHost |

Original acceptance: M1 static directories/navigation without broken links/duplicate projection; M2 on-demand selection equivalent to full export, preserving last view on errors; M3 real XFMD panel clicks, multiple windows, focus changes/closed targets; M4 cache hits/invalidation, rapid clicks, crash/restart, leases, limits/cleanup. All six milestones were subsequently delivered within explicit profiles; see introductory evidence.

G4-M1's Go structural frontend is prerequisite to G6-M1. Python was port evidence and later removed in G5-M4. G6 needs neither SDL domain execution, Go generation nor SDUI layout. Move viewpoint port from G5-M3 to G6-M1; G5-M3 consumes its export. G6 names a workstream, not a requirement to finish G1–G5 first.
