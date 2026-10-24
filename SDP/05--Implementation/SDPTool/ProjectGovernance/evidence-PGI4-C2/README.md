# PGI4-C2 retained evidence

- candidate.json pins the baseline, changed source hashes and final native binaries.
- checks.py/checks.json and sdptool-*.log record the final full race/vet suite.
- core-wire.json/controller-wire.json are actual service/controller test captures.
- native_trial.py runs actual managed-host positive/negative reviews in fresh state
  forks and restores deliberate source changes. It refuses to reset existing forks.
- native-summary.json/native-wire.json and protocol NDJSON retain bounded actual
  identities/results; protocol logs exclude private reasoning, raw tool outputs and auth.
- rework_native_negative.py extends the completed negative native trial through
  attributed evidence and guarded rework, restoring source afterward.
- schema-checks.json validates the four final capture files; old-reader.json records
  actual 0.10-reader refusal without state mutation.
- preliminary preserves earlier passing trials before final refinements.
- lab-delivery.json, manifest.json and final-record-checks.json pin delivery,
  artifacts, ledger prefixes and project validation.

The runners use explicit workstation paths. A fresh rerun needs newly selected state
forks and retained artifacts; never reset an uncertain review or overwrite native
history to manufacture a pass. Native test reviews concern pg_lab only. They are not
independent review of PG, owner acceptance, canonical publication or full readiness.
