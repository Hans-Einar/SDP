# Activity0001 qualification

Product candidate: **527cc9008bdf8f567f8b201b665b2960865518f1**, branch
sdp/sdptool-architecture-migration. Baseline: d0bf6dc37027bbb75f5f96fecc169e681a3a6c9b
(product bytes unchanged from a96affd). Source is frozen during final qualification;
subsequent edits are documentation/evidence only. No new public capability, Action,
Activity automation or browsing behavior was introduced.

Status: technical qualification passed; all required Linux/native gates pass.
Owner disposition is pending.
Owner acceptance, accepted model release, merge and publication are separate.

## Evidence and provenance

M0/ retains original source SHA-256s, 258 exported declarations and the 28-action
catalogue. M5/source-candidate.json pins all candidate Go files. All 258 original
declaration signatures remain (M5/library-compatibility.json); the added
governanceapp package supplies application composition, not a new public action.

| Check | Result and retained artifact |
| --- | --- |
| Baseline full SDPTool module, go test -count=1 -timeout 30m ./... | PASS; M5/baseline-tests.log, isolated baseline worktree |
| Initial migrated module, same full command | PASS; M5/migration-tests.log; started before final facade refinement, superseded for final-candidate proof |
| Frozen final candidate, same full command | PASS; M5/pinned-tests.log, source frozen at 527cc90 |
| Final root package, go test . -count=1 | PASS; M5/root-final.log, includes architecture/negative specimens, catalogue, routing, preview, program and adapter tests |
| Final compiled Session CLI/MCP journeys | PASS; M5/compiled-final.log; go test -count=1 -timeout 10m ./governance -run 'TestSession.*Compiled\|TestSessionMCP\|TestSessionHostMCPExplicitGrant' |
| Final race checks | PASS; M5/race-final.log; go test -race -count=1 ./governanceapp ./governancemcp ./projecthistory |
| Final go vet ./... | PASS; M5/vet-final.log is empty on success |
| Nested bootstrap and gh-sdp go test ./... | PASS; M5/bootstrap.log, M5/gh-sdp.log; source unchanged by migration |
| Compiled baseline/final CLI comparison | PASS, 187 cases, zero stdout/stderr/exit differences, no state written; M5/cli-parity.json and reproducible compare_cli.py |
| Native XFMD producer | Integration, navigation, overview, deferred and source-link tests PASS with final binary; M5/native-without-helpers.json |
| Native SDUI workflow | PASS on baseline and final with required helpers; M5/native-qualified.json. Initial missing-helper timeouts retained in the earlier logs |
| Independent product and qualification review | PASS at exact candidate; M5/Review.md, including independent negative-specimen/compatibility checks and final evidence closeout |

Native tests use the existing XFMD build at 8893a32 plus its unrelated local
RPM/process changes, and its tools/run_with_xvfb.py (isolated X server/home). They
exercise the selected actual producer binary, not a fake producer. XFMD files are
not edited. Missing SDUI helpers are built from unchanged SDP/SDUI source with
go build -tags desktop into /tmp/activity1-native/libexec/xfmd; the existing native
test executable is copied beside that private helper installation. No installation
or change to the user's desktop is performed. This setup is necessary because the
test temporarily unsets its helper override and verifies default executable lookup.

Go environment: go1.26.10, linux/amd64. Other native platforms are not qualified
here. The earlier uncompleted/canceled broad attempts are not passing evidence;
only the completed logs listed above count. Tests report behavior and recovery on
this candidate, not a mathematical proof of every runtime execution.

## Inventory, state ownership and removal closure

[CodeMapping.md](../04-Design/CodeMapping.md) maps every TARGET unit to actual code.
M5/inventory-closure.json reconciles all 60 operation rows, all 51 controller paths,
28 MCP names across six launch modes, library and external boundaries. There are
no unclassified operations or migration-only adapters. All rows are reconciled to this candidate and the completed evidence.

| Invariant | Evidence |
| --- | --- |
| Typed internal calls; no domain dependency on CLI/presentation | Architecture AST checks plus deliberately invalid specimens; direct Session action/CLI parity and independent review |
| One operation definition/policy per capability | governanceapp registry/input tests; frozen catalogue; one root CLI table; registered attribution/availability checks |
| One durable owner per state | Existing governance, model, Blueprint assignment, installer and shared projecthistory ownership retained; domain recovery/idempotency suites and race checks |
| Authority/cancellation semantics preserved | Existing domain and MCP capability tests, compiled journeys, strict request decoding, busy/stale/output-cap tests; no post-commit cancellation reinterpretation |
| Existing producer contracts retained | Root producer tests, compiled wire comparison and native consumer gates |
| All replaced paths removed | CodeMapping removal matrix and negative checks; reviewer found no transition shim or duplicate caller/dispatch/framing authority |

MCP capability registration, controller child-process dispatch, external gh-sdp
forwarding and SDUI/renderer process execution are permanent explicit boundaries.
They are not deferred migration adapters. No state-store merger or new authorization
model was selected. Unrelated future work remains KB061/KB060.

## Model and process claims

NOW/TARGET WORK and the retained Blueprint are pinned in
[Blueprints.md](../04-Design/Blueprints.md). Canonical parsing and structural task
constraints pass. The generated bundle's generic CODE_MAPPING unknown remains
honest; authored mapping, executable checks and independent review supply the
additional evidence. This trial demonstrates a usable bounded worker handoff,
not automatic translation from SDL to verified Go or acceptance of a model release.

[Bootstrap evidence](Bootstrap.md) applies to manual Activity records. The shared
management ledger remains append-only. Legacy-domain validation is filtered, with
separate Activity schema/hash/causal/single-active checks; the existing unmodified
validator is not claimed to support Activity; M5/records-validation.json records
the checks. M5/manifest.json pins retained evidence bytes. Activity0001 remains active pending
the actual owner disposition; Activities0002/0003 are not created or started.
