# Runnable SDL and SDUI widget test project

| Field | Value |
| --- | --- |
| id | KB-SDUI-006 |
| project | SDUI |
| type | Proposal |
| CardState | completed |
| Systems | SDL, SDUI, SDPTOOL |
| created | 2026-10-09T02:15:01.427224+00:00 |
| source | Owner SESSION-SDP-0010 T007 release and test-project discussion |
| next_review | Before selecting the next SDPTool release candidate |

## Need and source

The owner asks whether to release SDPTool / gh sdp after PR #52 and whether an
SDP-bearing project containing SDL and SDUI can exercise all widgets, suggesting
experiments/ as its home. T007 registered the proposal without selecting implementation or release.
T008 subsequently authorized the independent lab and full SDL/SDUI design below.

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

## Original proposal and refined acceptance

The original proposal suggested experiments/sdui_widget_lab/ for a consumer with its
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

## Release boundary

The runnable project and launch contract were selected and delivered in T008.
Future publication still requires its own selected release scope and actual
candidate install/upgrade rehearsal (original criterion 6). Release version/scope,
packaging, independent review,
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

## Owner execution selection — Session0010 T008

The owner authorizes creating the lab with a valid SDP folder installed via gh sdp
and asks about repository placement. Installer inspection and a nested-root preview
confirm both placements work. The owner explicitly selected the recommended independent local repo.
The later steering requires full SDL/SDUI source and interpreted Go calls.
No remote is created.

Execution plan: PLAN-LAB-0001; local Ref KB-LAB-001.
file:///home/warloc/git/sdui_widget_lab/SDP/05--Implementation/WidgetLab/Plan.md

The signed SDP 2.1.0 installation through gh sdp and PLAN-LAB-0001 are completed.
Future publication and release-candidate install/upgrade rehearsal remain separate
release work; no future candidate or publication is claimed here.

- 2026-10-09T02:22:44.836480+00:00: EVT-KB-SDUI-000032; backlog -> active/in-progress under T008 authority.

## Delivery — LAB1-M3

Independent local repository: /home/warloc/git/sdui_widget_lab, branch lab/initial.
Runtime candidate db206bc91615b6a37474c80a1306a53bde80d502; documentation/evidence
closeout 0946342. Signed installed SDP 2.1.0 is separate from the clean pinned
widget-library dependency 04f88ff918e7c25c9fa1883f061ef20b8c3923a2.

Complete composed design-core System, source-owned SDUI application and 17 executable
SDL routines call independently typed named Go functions. Tests prove actual UI
dispatch, changing SDL invokes or SDUI ref without rebuilding, and early rejection
of incompatible bindings. Final race suite and canonical source check passed.
Real native Fyne/X11 input passed 26 checks, including typed return values, modal
Cancel/Accept, nonmodal editing and explicit recovery after a requested Go failure.
Screenshots and candidate hashes are retained in the consumer evidence.

file:///home/warloc/git/sdui_widget_lab/SDP/05--Implementation/WidgetLab/Evidence.md

All widget families are available through make run and the manual checklist.
This is bounded native smoke, not every interaction combination or independent
review. Installed SDPTool discovers the project and validates composed SDL, but
its released parser reports SDUI 0.3 unsupported; the lab uses pinned 0.3 libraries.
Current-release upgrade preview is no-change/applicable without conflicts. No new
release, generic launcher migration or remote publication is claimed. Original
criterion 6 is retained as a future publication obligation outside the selected lab.

- 2026-10-09T03:04:59.384968+00:00: EVT-KB-SDUI-000033; active -> completed. LAB1 and local Ref KB-LAB-001 completed with WIDGETLAB-VER-001 evidence; no selected implementation remains.
