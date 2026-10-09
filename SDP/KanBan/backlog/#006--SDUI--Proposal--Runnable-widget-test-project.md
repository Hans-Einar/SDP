# Runnable SDL and SDUI widget test project

| Field | Value |
| --- | --- |
| id | KB-SDUI-006 |
| project | SDUI |
| type | Proposal |
| CardState | backlog |
| Systems | SDL, SDUI, SDPTOOL |
| created | 2026-10-09T02:15:01.427224+00:00 |
| source | Owner SESSION-SDP-0010 T007 release and test-project discussion |
| next_review | Before selecting the next SDPTool release candidate |

## Need and source

The owner asks whether to release SDPTool / gh sdp after PR #52 and whether an
SDP-bearing project containing SDL and SDUI can exercise all widgets, suggesting
experiments/ as its home. This registers the proposal; implementation and release
publication are not selected by this discussion.

## Observed baseline

- KB-SDUI-003 and its combined PR #52 are complete and merged. The implemented
  inventory uses Go/Fyne. Its local verification package is not a published release.
- Six connected Fyne applications exist under SDL/go/examples: collections,
  panes, commands, values, text and previews. They exercise real SDL bindings and
  providers, but are developer fixtures rather than one consumer project with an
  SDP directory and a documented manual all-family workflow.
- experiments/mvp1_sdl is a design-authoring corpus with an experimental SDL
  profile, not a running all-widget system. No existing experiment matching the
  requested project was found during this inspection.
- The new families run through DocumentHost. The standalone sdui-fyne helper's
  RuntimeView route rejects the new 0.3 families; merely putting all controls in
  a .sdui file does not make that route runnable.
- Current published versions inspected on GitHub are SDP v2.1.0 and gh-sdp
  v0.2.1. SDPTool/package.sh builds sdptool and its manifest/checksums; it does not
  itself package the separate Fyne applications. A new installer release alone
  is therefore not proof of a runnable full-widget installation.

## Proposed scope and acceptance

Use experiments/sdui_widget_lab/ for a standalone consumer test project with its
own SDP directory, discoverable SDL design and SDUI sources, and a Fyne application
with real typed handlers and providers. Reuse existing product libraries and
fixture scenarios. Do not introduce another runtime implementation.

Proposed acceptance:

1. One documented start command opens the application for manual testing.
2. A visible scenario catalogue and checklist cover every implemented widget
   family, including disabled/error/empty states and meaningful interaction.
3. SDL bindings change observable application state; collections, commands,
   dialogs, values, text/IME, panes and previews are exercised as applicable.
4. SDPTool can discover and inspect the project's models and UI sources.
5. Tests distinguish direct application launch from generic helper launch. Choose
   the launch contract explicitly before implementation; extending the generic
   helper is a separate design change, not an assumed effect of this project.
6. Before publication, test the selected release candidate's actual installation
   and project upgrade, plus the documented widget launch route and its packaged
   dependencies. Coordinate gh-sdp's pinned SDP version with that release.

## Boundaries and next action

Recommend selecting the runnable test project before publishing a release intended
to showcase the full widget inventory. First settle the consumer launch route and
bounded implementation plan. Release version/scope, packaging, independent review,
signing and publication remain separate release work. Do not reopen completed
KB-SDUI-003 merely to track this successor.

KB-SDUI-005 separately owns combined preview and launch-readiness integration.
Its local shared-tree facade-test mismatch, reported after widget delivery, must
be checked on any selected candidate that includes those changes; this observation
is not evidence that the reviewed PR #52 candidate failed.

## References

- [Session0010 T007](../../Sessions/session-%230010--SDUI_widgets.md)
- [Completed widget inventory](../completed/%23003--SDUI--Proposal--Capabilities-and-navigation-pilot.md)
- [Consumer and launcher boundary](../../05--Implementation/SDUI/Widgets/WCI4-consumer-preparation.md)
- [Connected preview application](../../../SDL/go/examples/previews/README.md)
- [MVP1 design corpus](../../../experiments/mvp1_sdl/README.md)
- [Release checklist](../../../Toolkit/docs/ReleaseChecklist.md)

## Worklog and revisions

| Time (RFC3339) | Actor / event | Work, finding or decision | Evidence / remaining work |
| --- | --- | --- | --- |
| 2026-10-09T02:15:01.427224+00:00 | codex; EVT-KB-SDUI-000031 | Registered owner proposal and release-readiness assessment | Select launch route and bounded plan; no implementation or publication claimed |
