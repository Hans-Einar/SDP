# WR2 observed verification artifacts

Source candidate: 4174c68, Go 1.27.1, Linux amd64. Current Go/runtime sources remain
unchanged; later changes are scoped documentation/management records. checks.json
records exact commands, durations, logs and binary hashes. Logs are retained here;
their absolute cache paths identify the original execution locations.

Probe scripts are retained as the exact executed harnesses with their original
workstation paths and preconditions, not as general installation utilities. The
bound worker harness uses a scripted host binding; it submits no model task and
records no owner approval. codex-mcp-readcalls.json is actual 0.164.0 host-mediated
read calls, not the pinned PG controller. runtime-probes.json retains that
controller's expected version refusal separately.

provider-trial.json comes from the existing lab script and checked reduced specimen
experiments/blueprint_mvp1/NOW. This is an isolated regression fixture, not a new
Ponsse pilot or changes to the owner's pilot model. provider-input-attempts.json
retains two refused preliminary inputs, with unchanged failed directories in the
private cache. Neither refusal was bypassed or labeled a passing trial.

No raw account response, authentication material or unfiltered Codex stderr is
retained. Independent review and new native model execution are not established
by these records. See ../../Execution-Evidence.md for result-level limitations.
