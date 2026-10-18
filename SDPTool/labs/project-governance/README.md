# ProjectGovernance owner lab

This is a local, release-independent **owner trial**, not the finished PG interface.
The real PG core, controller, scoped MCP adapter and Codex host perform the work.
The small Python console only prepares requests and renders existing responses.
It does not implement a competing run state machine or accept work automatically.

## Start here

From the prepared lab directory:

```sh
./lab
./lab run
```

The first command shows the candidate, project, goal, current run, selected Session,
step, card, plan, authority and writable/excluded paths. No model starts until you
run the second command and type `start`. You can inspect the documents beneath
`project/SDP` first. The project has minimal pilot coverage records; it is **not**
a full released SDP installation. The actual maintained skill sources are copied
from the pinned candidate, with project-local discovery links.

The project is a tiny Go label command. Its defect is deliberate: it trims ASCII
spaces but leaves other surrounding Unicode whitespace. Try it before and after:

```sh
cd project
go run ./cmd/label $'\t  Keep  Case\u00a0'
go test ./...
```

The initial tests pass but miss the defect. The selected task changes only
normalize.go and normalize_test.go. The command, unrelated file, process records
and skills must be preserved. This first project is small enough to inspect in full;
a later selected task in sdui_widget_lab checks whether the workflow transfers to
an existing consumer. This trial alone does not establish that broader fit.

## While the agent works

The console shows actual agent commentary/final messages, tool completion events,
run state and supported per-request decisions. Raw controller output is retained
in runtime/events.ndjson, exact forwarded commands in runtime/commands.ndjson and
diagnostics in runtime/controller.err. Output is a readable stream, not a dashboard
claiming uninterrupted live state or a full tool-output viewer.

Type one command and press Enter:

| Command | Meaning |
| --- | --- |
| /status | Read the current recorded run/attempt state |
| /steer Keep interior tabs unchanged and include a test | Add direction within the existing assignment |
| /interrupt Stop this attempt so I can inspect it | Request interruption; inspect the observed outcome |
| /scope I also want a project-wide label registry | Record a broader request and stop for replanning |
| /accept, /decline, /cancel | Answer only the exact displayed pending request, once |
| /quit | Stop this controller; inspect/recover the outcome before further work |

There is no automatic approval. The report prompt displays its actual transition
arguments; answering it does not accept the completed work. Only the supported SDP
MCP decision is exposed. Changing the assignment, continuing a stopped run and
replanning remain low-level workflows in this candidate: inspect the result with
the development agent rather than editing state or creating duplicate runs to
bypass uncertainty. EOF sends neither consent nor cancellation. This lab uses the
owned stdio host; it does not demonstrate the separate-server detach/attach path.

## After the turn

```sh
./lab status
./lab provenance
./lab result
./lab feedback
```

Result shows the real Git diff, runs project tests, compares baseline file hashes
and discovery links, and lists untracked files. These are local checks, not independent
review. `verifying` means the agent reported its result and checking remains open;
it does not mean the card is complete. Current PG does not publish the canonical
Session/card outcome or supply a finished evidence/acceptance gate.

Record observations in FEEDBACK.md or discuss them with the development agent, who
can maintain that record. Useful prompts:

- Before start, could you identify the goal, selected work and allowed changes?
- During execution, could you tell who was acting and what they were doing?
- Did steering, interruption and the report decision behave as you expected?
- At the end, could you distinguish changed code, reported tests and unaccepted work?
- Which part was confusing because of this console, and which because of the workflow?

Record expectation, actual experience and candidate/run identity. Classify each
observation as implementation defect, workflow mismatch, new requirement or understood
limitation. The owner and development agent review the observations, revise the
canonical design/plan where needed, then select the next increment. Passing technical
tests does not substitute for owner experience or imply owner acceptance.

`./lab recover` reads mapped host history and reconciles observations; it never
resends the task. It can explicitly refuse incomplete or ambiguous history. Do not
reset this lab after a failure: retain the files and diagnostics for investigation.
Create another destination for another clean trial.

## Candidate and installation isolation

candidate.json records source commit, source archive hash, Go version, binary hashes
and the copied harness hashes. The console refuses changed binaries. `go version -m`
on a binary may show the module's existing development version; candidate.json owns
this trial's exact source identity. The host requirement is Codex 0.161.0 using the
existing managed ChatGPT account. No API-key fallback or account reconfiguration is
provided. The three binaries live only in candidate/bin; stable gh-sdp, SDPTool,
PATH and the widget lab are not replaced.

Project files are in project. The operational state is outside that project in
runtime/state, bound to this exact checkout path. Do not copy the state to another
project, overwrite the candidate during a run, or treat a copied lab as a new run.
A fresh preparation creates fresh identities and an initial project Git branch.
No remote, push, merge or release is involved.

## Reproduce preparation from the development source

```sh
python3 SDPTool/labs/project-governance/prepare.py --source /path/to/SDP-project-governance --revision COMMIT --destination /path/to/new-pg-lab --go /path/to/go
```

The destination must not exist. Preparation exports committed SDPTool/SDL/SDUI/Skills
sources, builds all three binaries with buildvcs=false, records their source identity,
and creates a ready trial run. It does not start Codex. The harness source may be a
newer development version than the selected engine; both identities are recorded.
An interrupted preparation leaves its directory for inspection and refuses reuse.

The action catalogue and JSON invoke adapter belong to KB-SDP-053, now integrated
on main. This lab uses the standalone PG commands; it does not add PG Actions or
require an XFMD change. A long-running owner conversation still needs streaming events and control;
a one-request JSON adapter alone is not that interaction surface.

## Assignment preparation trial — PGL2-C1

Build the development controller from SDPTool, then use a fresh destination and a
pre-PGL2 SDPTool executable to exercise compatibility refusal:

```sh
python3 preparation_trial.py create --destination /absolute/new-trial --controller /absolute/new-sdp-controller --sdptool /absolute/pre-PGL2-sdptool
python3 preparation_trial.py inspect --destination /absolute/new-trial
```

The script refuses an existing destination. It creates synthetic canonical coverage
and file artifacts, binds a separate store, prepares disjoint SAD/Master assignments,
repeats requests through fresh subprocesses, refuses overlapping scope, then changes
TARGET and verifies Master inspection becomes stale. Finally it checks that the old
binary refuses the upgraded store without changing HEAD. Actual requests/results
are in evidence.json; no model task is submitted and no design acceptance is inferred.
The source artifacts are fixtures, not ModelGovernance or Blueprint API evidence.
Existing owner-lab project/runtime and candidate remain untouched.

## Provider-backed design check — PGL2-C2b

Build the development binaries from SDPTool, without replacing installed tools:

```sh
go build -o /tmp/pg-controller ./cmd/sdp-controller
go build -o /tmp/pg-sdptool ./cmd/sdptool
python3 labs/project-governance/design_trial.py --destination /absolute/new-design-trial --controller /tmp/pg-controller --sdptool /tmp/pg-sdptool --source ../experiments/blueprint_mvp1/NOW
python3 labs/project-governance/design_view.py --preview /absolute/new-design-trial
```

The separate project uses real ModelGovernance and retained Blueprint services with
a synthetic assignment. It proves repeatable current checks, uncommitted TARGET
drift with unchanged metadata, no writes during inspection and exact restoration.
It leaves TARGET deliberately changed, so the readable view reports blocked.
Use `--json` on the view for the full result. Each preview retains binaries,
requests, results and the selected project/runtime. Existing owner trials are untouched.

The underlying command is `sdp-controller assignment design-check --project ROOT
--state STATE --input REQUEST`. It requires the prepared assignment revision/generation
and canonical Blueprint assignment/event. Status current means a consistent review
basis; it never means owner approval or permission to launch. A blocked result is a
successful inspection with conditions, not a process failure. Invalid requests fail
as usual. The provider's bounded live-check budget can return unknown and block a
check in larger projects; this slice preserves that limitation.
