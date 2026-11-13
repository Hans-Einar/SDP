# Shared runtime distribution requirements

Authority: Activity4 mandate and owner selection to execute PLAN-SDP-0029.
This is the authoritative concern for the shared runtime; project installation,
model, Activity and source-preview contracts retain their existing authority.

| Obligation | Required outcome | Implementation / proof |
| --- | --- | --- |
| One update owner | gh extension installation/upgrade replaces the exact canonical SDPTool engine; no runtime engine selector or XFMD copy | gh-sdp build-time verified republisher; actual gh binary-install/upgrade and byte-identity gate |
| Local entrypoint | Explicit one-time setup creates a stable direct entrypoint and diagnoses conflicts or shell shadowing | runtime_setup.go; conflict/repeat/transient tests; actual gh registration |
| Offline consumer | XFMD calls the registered engine directly without gh, credentials or network; rejects incompatible capabilities | native preflight and isolated network-namespace application tests |
| Complete helper ownership | Source preview/check and native SDUI prototype belong to the same update; no separately fetched helper | SDUI preview library, embedded native_host.go payload; actual registered-engine SDUI GUI probe |
| Shared operations | CLI and typed actions reuse one domain handler; host operations never weaken project guards | action_runtime.go, runtime_cli.go, action policy/legacy compatibility tests |
| Data preservation | Engine updates do not mutate projects; explicit Framework upgrades retain signed descriptor, receipt, conflict and owner-document guarantees | real gh no-project-mutation check; seven original signed predecessor trials |
| Honest recovery | Interrupted gh updates never silently select another engine; documented repair uses the supported update owner | actual interrupted-transfer, fail-closed and reinstall gate |
| Ordered delivery | Canonical engine and extension become available before XFMD removes its old bundled runtime | release checklist and downloaded public asset/DNF verification |

The native prototype needs OS GUI libraries; these are normal platform packages,
not another SDP engine update route. Release qualification is Linux amd64. A
protocol identifier, model checkpoint or calendar status is not runtime evidence.
