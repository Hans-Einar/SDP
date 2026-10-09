# SDPTOOL-VER-RSP1 — program discovery and native execution

## Candidate and result

Runtime candidate `99631595eb2ccf5d1e8cc52be5ce30a02b492341`, clean branch sdp/runnable-programs based on main
merge 04f88ff. Consumer declaration/evidence commit 75b6dca in sdui_widget_lab.
RSP1-M1 and M2 completed at 2026-10-09T10:18:54.830417+00:00. [Source inventory](evidence/source-hashes.json)
and [binary build identity](evidence/engine-build-info.txt) identify the actual
candidate; later records do not change the tested implementation.

Full final `go test -race -count=1 ./...` in SDPTool passed: root 49.321s,
install 64.147s, model 34.244s, presentation 1.029s, tools/profile 47.721s.
[Raw log](evidence/sdptool-race.log). Focused program race tests also passed in
2.378s. The tests cover read-only discovery/source preservation, duplicate and
invalid catalog, missing/escaping executable/source, nonexistent entry, stale
source/declaration/removal, exact argv/cwd, exit code, rejected JSON invocation
and Unix cancellation of owned descendants. Existing source/installation suites
passed; no concurrent shared-tree code was included in the isolated candidate.

## Actual installed gh-sdp integration

Actual gh-sdp v0.2.1 bootstrapped the locally test-signed development descriptor
and invoked SDPTool 2.1.0-dev.rsp1. It was not a fake wrapper or replacement gh
extension. The descriptor selects only this development executable; no install,
upgrade or publication occurred. Default signed-release bootstrap remains intact.

From the lab root, with its generated local development environment selected:

```sh
source .build/sdptool-programs/env.sh
gh sdp . discover --json
gh sdp . discover
DISPLAY=:193 FYNE_SCALE=1 gh sdp . run --program widget-lab
```

[Discovery](evidence/discover.json) retains the validated Lab.sdui source and
returns widget-lab as runnable with source, page entry, revision and explicit Run
arguments. The manifest command is [make, run]. [Run log](evidence/run.log) proves
the existing dependency/build/application start path was used. The real native
window opened on isolated Xvfb :193; [initial pixels](evidence/start.png) show
waiting state. Actual X11 click (65,56) then produced the
[SDL Run -> GoRun -> SDUI result](evidence/run-result.png). Normal native window
close returned 0 through application, make, SDPTool and gh-sdp.
[Receipt](evidence/result.json) records the five integration observations.

Development binary and environment are available locally under the lab's ignored
.build/sdptool-programs directory. The reusable test-descriptor generator is
[evidence/dev-descriptor.go](evidence/dev-descriptor.go): after building sdptool
there, run it with distribution directory and candidate commit arguments. It
creates descriptor/signature/public key, then discards its private test key.
This helper is verification tooling, not a production signing or release path.

## Scope and limitations

SDPTOOL-REQ-RSP1 / SDPTOOL-DES-RSP1 are verified to their bounded contract.
Discovery readiness means validated source/entry and command availability;
application-specific build and providers remain the runner's responsibility.
Program revision covers declaration and primary SDUI source, not all transitive
inputs. Existing source navigation remains unchanged; clients use the returned
program revision/run target separately. No new SDL/SDUI grammar, XFMD UI change,
released engine/default gh-sdp update or independent reviewer approval is claimed.
LAB1's 26 native checks remain their original evidence; this integration checks
the newly introduced launch route rather than relabeling those checks as rerun.

Canonical management/Session closeout lives in SDP-vNow; the working branch keeps
matching selected-plan and evidence copies. Shared unrelated changes are preserved.

Record checks passed: canonical management replay/schema/placement and isolated
Toolkit validation. Candidate product source hashes remain unchanged at closeout.
