# SDPTool runnable SDUI programs

| Field | Value |
| --- | --- |
| id | PLAN-SDP-0023 |
| project | SDP |
| state | completed |
| PlanType | ImplementationPlan |
| BranchPolicy | current |
| CommitPolicy | milestone |
| source | Owner Session0010 T010; KB-SDP-051 |

## Outcome and authority

The owner assigns discovery and launch of runnable SDUI programs to SDPTool.
`gh sdp . discover` must retain source discovery and additionally describe runnable
applications. The widget lab is the first real consumer. XFMD consumes this
contract; changes to XFMD, a main merge and publication are not selected here.

## Analysis and design boundary

Current discover inventories source files; XFMD directly launches sdui-fyne while
make run starts the lab's connected host with its Go registry/providers. Neither
source parsing nor generic prototype launch supplies those application services.
Implement explicit project-owned SDP/programs.json declarations, separate from
SDUI grammar and generated source inventory. Each program has a stable ID, label,
source/entry and argument-vector command. Run in the project root. Discovery is
read-only and never probes commands by executing them. An explicit run selects
one ID and starts its declared application; preserve argument boundaries and
propagate child exit status. Invalid declarations/missing commands are diagnosed
without concealing existing source discovery. Validate containment and selected
source/entry before launch. Optional expected revision rejects stale selections.

Owner T010 supersedes the former blanket prohibition on executable registration
only for this explicit application declaration/Run boundary. Ordinary SDL/SDUI
source, preview metadata and rendering still cannot register commands. Do not
infer executable commands from Makefiles or put project-specific Go functions in
SDPTool. A generic compiled-plugin loader would add an unnecessary runtime.

## Git and milestones

Use isolated working branch sdp/runnable-programs at main merge 04f88ff because the
shared tree has concurrent modifications. Keep canonical management/Session in
SDP-vNow; preserve all unrelated changes. Consumer changes stay on lab/initial.
Commits per milestone; no merge or published release is inferred.

| Milestone | Acceptance | State |
| --- | --- | --- |
| RSP1-M1 | Explicit declaration, source-preserving discovery and foreground program run with invalid/stale/exit/cancellation tests | completed |
| RSP1-M2 | Widget lab declaration and actual gh-sdp development-candidate discover/run through real Fyne window; documentation and exact evidence | completed |

## Verification

Run focused and existing SDPTool suites on the isolated candidate; test missing,
invalid, duplicate, escaping and unsupported declarations, entry association,
read-only discovery, exact argv/cwd, error/exit propagation and stale selection.
Verify actual gh extension routing with a locally signed development descriptor;
keep released bootstrap defaults unchanged. A local candidate is not a release.
Record real native launch/interaction, including SDL -> Go result, and limitations.
Same-context inspection is not independent review. [Evidence](Evidence.md) records
passing full race suites and actual gh-sdp/native discovery/run on 9963159.
All selected milestones are complete; merge/release and XFMD UI remain separate.
