# PGI3-C4 — executing-turn reattachment checkpoint

PLAN-SDP-0018 / KB-SDP-038, 2026-10-08. Candidate: parent `95f821fcf8c7f3327fe627c5ccc69f436b8f3dac` plus
[the exact product/test file set](live-PGI3-C4.json), digest `b4cb4cc97cd3227a16c56bb29d82ec36f1b183a6a033a8da0a6555e5f7aa69b4`.
This is a same-context implementation and verification checkpoint, not independent
review, owner acceptance, full PGI3 delivery or a release.

## Delivered behavior

An opt-in private Linux Unix WebSocket connection leaves the separately launched
Codex 0.161.0 server alive when its controller closes or dies. The new attach
command rejoins only the mapped active turn on the same measured host, without
turn/start, source/permission overrides or uncertain-input resends. The observation
epoch changes and retains a gap; the original bounded agent capability remains
valid. Ordinary history recovery retains its prior capability fencing behavior.

Endpoint identity binds the private socket and peer process lifetime. Source drift,
a replacement server, old observations, unknown mappings, incomplete history and
unresolved outbound inputs/approval closure refuse. New approval requests require
new explicit answers. Terminal history replays completed content without inventing
completion for partial items. See the [controller instructions](../../../../SDPTool/governancecontroller/README.md).

## Measured workflow

A preliminary actual host probe confirmed that an executing command survives
connection closure and that thread/resume observes the same turn's completion.
Raw proxy JSONL did not initialize; idle threads without a persisted turn did not
resume. Neither result was promoted to a success claim. The installed schemas and
[official app-server reference](https://learn.chatgpt.com/docs/app-server) guided
the measured transport and subscription boundary.

Two initial controller starts refused before model submission: this host creates
a direct socket alias into a private per-user runtime directory. The implementation
now accepts that measured form only after checking both parents, the canonical
socket and peer identity. Indirect aliases and public paths still refuse.

The final binary completed two disposable Maintenance journeys. One controller
explicitly detached; the other was killed with SIGKILL while the mapped turn was
in progress. Each new controller attached to exactly one existing run, attempt,
thread and turn, submitted in-scope steering and accepted exactly one observed
report request as a driver scenario. The original assignment-bound MCP report
succeeded and the routine reached verifying. No new turn was submitted and neither
routine was completed or owner-accepted. The interrupted observation interval stays
visible as gap=true. The binary reproduced identically after the documentation and
direct-dependency annotation updates; later test additions do not change runtime
sources. Every declared excluded file retained its preimage hash; only the two
permitted Go files changed, and separate driver invocations of go test passed.
Those test reruns share implementation context and are not an independent review.

## Evidence and remaining work

Affected-package race tests, compiled core CLI, vet, nine captured application
request/results, twelve installed-host protocol cases and 54 live wire records
pass. Protocol fixtures explicitly check thread/resume has only threadId and no
new dispatch, plus unloaded/wrong-turn/summary refusal. Core fixtures check legacy
mapping refusal, source drift, unresolved approval/input, old observer fencing and
preservation of only the original capability. Socket tests cover private ownership,
direct/indirect aliases, replacement identity, frame boundaries and version refusal.
No full-repository regression or publication/installability claim is made.

The accompanying JSON retains candidate/binary hashes, actual identities/events,
pre/postimages, sanitized host MCP results, outputs and the disposable driver.
Management/Toolkit and final ledger/link integrity checks are recorded separately
there after record updates. Historical ledger prefixes and concurrent shared work
remain untouched. Known C1 event-ID collisions still require deliberate integration
reconciliation; this checkpoint does not authorize a merge.

Next is applicable role/skill/context provenance within current delegation rules.
The disposable fixture does not install SDP skills; the model explicitly reported
that absence. Instructions to read skills are not evidence of actual loading.
Native child attribution and independent review remain unverified. PGI4 still owns
checked evidence/publication; PGI5 owns independent integrated verification.
Unresolved approval closure and server failure remain explicit recovery boundaries.
