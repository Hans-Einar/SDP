# PGC1 retained evidence

Code candidate 98b449c, Go 1.27.1 linux/amd64. final-checks.json binds commands,
exit codes, logs and final controller/MCP binary hashes. All final checks pass.
Later closeout changes only docs/management/evidence; runtime source is unchanged.

The real controller probe performs initialization/account-mode inspection and idle
MCP discovery against installed Codex 0.164.0; modelSubmitted is false. No account
payload, token or unfiltered Codex stderr is retained. Future-version tests use a
recorded fake protocol host and private socket fixtures, not real newer releases.

schema-cases contains actual test-emitted objects; final-schema-validation.json
records 14 host and 46 application validations and exported-schema hashes. Runtime
captures identify temporary fixture paths; these are observed records, not current
project registrations. final_checks.py is the exact execution harness with original
workstation paths/preconditions, not a portable installer. Restore its original
cache location before reproducing it; do not execute it inside this evidence tree.

lab-candidate.json and lab-preparation-summary.json identify a freshly prepared
private fixture, not the owner's existing lab or an actual model task. The run is
ready, not running or accepted. Existing labs and stable signed installations were
not replaced. No independent review or full current-host workflow claim.
