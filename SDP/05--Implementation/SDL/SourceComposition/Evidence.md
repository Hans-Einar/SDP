# Source composition implementation evidence

PLAN-SDP-0012; baseline 32aa68a. Results will be appended at actual milestones.
No new-language implementation acceptance is claimed by plan creation.

## SSI1-M1 frontend

Starting commit 70aa48e plus milestone diff. Commands: GOMAXPROCS=2 go -C SDL/go
test -p 2 ./parser ./sourcegraph ./cmd/sdl (Go 1.27.1 local toolchain): pass.
CLI check examples/source-composition/System.design: valid design-core/0.6,
System Frontend, three original files, revision
eab3d03d82afa4fdb4627a1870be457029a936e0ad9c9d8e7f92543b14956862.
Late-root reuse, context isolation, original spans, source cycles/diamonds,
invalid membership/names/profiles, paths, limits, output protection and CLI
fragment/file-map operations are covered. Initial test fixture omitted the
existing required contract completeness property; adding the legitimate property
resolved that fixture failure without relaxing semantic validation.
Consumer migration and independent review are still pending.

## SSI2-M1 consumers and review corrections

Candidate: c090b28 plus the SSI2 milestone diff. Go 1.27.1, GOMAXPROCS=2,
`go -C SDL/go test -p 2 ./parser ./sourcegraph ./viewpoint ./documents ./broker ./cmd/sdl`
and `go -C SDPTool test -p 2 ./...`: passed. Additional sourcegraph, documents,
broker, CLI and SDPTool composition tests passed after addition.

All consumers use the graph revision and original-source provenance. Tests cover
registered System/profile mismatch across tree/select/ViewPlan, non-entry change
invalidation, old revision rejection retaining the previous document, dependency
output protection, broker cache invalidation and stale render rejection preserving
reader leases. The real Frontend pilot preserves every baseline fact; file moves
change revision/provenance while fact identities remain stable.

Independent Reviewer initially requested changes: diagnostic amplification,
aggregate-token bypass in pure compilation, missing opposite-edge cycle evidence,
and missing ViewPlan registration validation. All four were corrected with
regression tests. Diagnostic neighbourhood evidence is bounded to 100 diagnostics
and 8 related locations each; it is not exhaustive causal analysis. Repeated-source
warnings are capped at 100. The cache retains grammar token counts. Final review
and full verification remain SSI3 work.

## SSI3-M1 verification and independent review

Product candidate: 5fb5bccd9100751996e0f6ad6eb5cd46b98e953b, plus regenerated
actions_gen.go/manifest in examples/generatedmodel. The additional delta is 33
explicit empty Span.Source fields and the generated output hash, not new action
behavior. Rebuilt with the existing sdl-gen using edit-apt-cell.sdl/.sdui and
package generatedmodel; constructor equality, runtime application and generator
byte comparison pass. No generated source was hand-edited.

Go 1.27.1; Linux; GOMAXPROCS=2; test parallelism 2. Ran race-enabled tests for all
tracked headless package directories in SDL/go (desktop-only packages excluded by
Go build constraints; unrelated untracked sourceinput excluded), and SDPTool ./....
Initial SDL run passed every package except the stale generated constructor bytes;
after regeneration, codegen and examples/application/generatedmodel pass. Existing
parser cases include structural/action/class compatibility; runtime, bridge, reload,
snapshot, devhost, documents, reader and broker pass. SDPTool full race suite passes.
No graphics/native GUI verification is claimed.

Independent Reviewer, fresh context, read-only: initial changes-required result
and four fixes are described above. Re-review approved the bounded product scope
with independent uncached targeted SDL and complete SDPTool tests, plus original
adversarial token/cycle probes. At review, tracked product diff from 70aa48e had
SHA-256 f27c002e59394c71f73d4603da72b726c47ec4c71f53f9caaf082b9d8ff15498;
the two added tests had hashes e6cbfb115576d73b38212acea27be8f8eca6dccf80c041e5f0af432d1a61822a
(documents/sourcegraph_test.go) and 6091494b68153f27459e001ed34af251d171fd79b5fa26f86eb30df86f86e7f4
(SDPTool/sourcegraph_test.go). These were committed in 5fb5bcc. Reviewer separately
approved the regenerated constructor delta and profile documentation; independent
codegen/application tests passed. Profile document reviewed hash:
23c13a7ed401cb3bb2c70296e98e11c2fb3e4a31d3a29cf99c910f881f468f74.
This is technical review acceptance, not owner acceptance or release publication.

### Real author-to-document workflow

From SDL/go, `go run ./cmd/sdl viewpoints examples/source-composition/System.design
--output /tmp/sdl-session-0001-navigation --project frontend` produced a 16-file
navigation bundle with no pre-generated detail diagrams. `go run ./cmd/sdl view
examples/source-composition/System.design --uri sdl-view://frontend/VP02
--output /tmp/sdl-session-0001-selected` generated the selected architecture view
on demand. Both report revision
eab3d03d82afa4fdb4627a1870be457029a936e0ad9c9d8e7f92543b14956862.
The result shows Frontend membership and SdlCommandProcess decomposition directly
from the three checked source files; sources.json preserves file hashes and spans.
No diagrams were authored manually. Temporary outputs are reproducible previews,
not authoritative documentation or permanent evidence files.

| Obligation | Evidence / result |
| --- | --- |
| Original facts preserved | Real Frontend document test compares all baseline statements and declaration counts; only System and its membership added |
| Late root, shared cycles, isolated contexts | sourcegraph/model_test.go: syntax cache reuse, cold equivalence, unreachable cache rejection, originals unchanged |
| Invalid sources and bounded work | Missing/duplicate/profile/path/alias/model-cycle and aggregate byte/token cases; diagnostic stress regression |
| Stable identity and all-input revision | Source relocation/move, unlisted file, non-entry mutation and stable fact-ID tests |
| CLI consumers and formatting | sourcegraph_test.go checks, fragment state, explicit file-map, views/view revision equality and input protection |
| Registered navigation | SDPTool sourcegraph_test.go tree/select/ViewPlan registration rejection, old-tree failure preserving bundle |
| Stale render and cache behavior | Broker composition test invalidates cache for child changes, blocks stale delivery and retains earlier lease |
| Compatibility | Existing structural/action/class cases and headless runtime/generator/document suites pass within their existing profiles |
| Documentation / records | English profile and commands, Session roadmap and card/plan dispositions; management/Toolkit/link/history checks |

### Closeout boundary

KB-SDL-007 completes selected syntax reuse/context validation. KB-SDL-005 returns
to backlog with public cross-System linking still open. KB-SDP-020 remains backlog
for actual model migration; no authoritative model or registration was moved.
KB-SDP-042/043 process/timeline proposals remain separate. Session 0001 is complete
for its bounded source-composition goal. No native XFMD changes, full experimental
MVP1 acceptance, semantic incremental cache, merge or release occurred.

Final document checks: ProjectManagement validation passes (54 cards, 27 management
records, three lineage operations, 423 events); Toolkit repository validation and
git diff --check pass. Local file destinations checked in 25 changed/new Markdown
documents; both append-only ledger prefixes match baseline 32aa68a byte-for-byte.
Card moves update live links and index together; earlier event paths remain history.

| Milestone | Commit / phase branch |
| --- | --- |
| SSI0-M1 | 70aa48e — sdp/sdl-source-sets/implementation-plan |
| SSI1-M1 | c090b28 — sdp/sdl-source-sets/frontend |
| SSI2-M1 | 5fb5bcc — sdp/sdl-source-sets/consumers |
| SSI3-M1 | Commit containing this closeout — sdp/sdl-source-sets/verification |

Each phase stacks on its predecessor. Existing authorization covers phase pushes;
there is no merge or release in this plan. Untracked sourceinput and Node artifacts
were neither adopted nor staged.
