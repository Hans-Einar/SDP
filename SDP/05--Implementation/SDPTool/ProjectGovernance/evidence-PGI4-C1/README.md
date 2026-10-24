# PGI4-C1 retained verification

- candidate.json pins baseline plus changed source hashes and the controller binary.
- checks.py / checks.json / sdptool-*.log record full final race/vet verification.
- core-wire.json and cli-wire.json are actual test captures, not hand-authored wire fixtures.
- lab_trial.py / lab-summary.json / lab-wire.json record external tests and evidence
  handling on the C8 native result, in a private state fork. No new model is launched.
- schema-checks.json records JSON Schema validation of those three final captures.
- lab-delivery.json identifies the committed lab artifacts and read-only helper.
- manifest.json and final-record-checks.json bind artifacts, ledger prefixes and links.
- preliminary retains the earlier passing checks/trial before final refinements.

Run checks.py from the repository root with the stated local Go toolchain. The lab
runner deliberately refuses to overwrite its existing state fork; use an explicitly
selected fresh fork for another trial and retain the prior evidence. The original C8
store remains read-only. Native transcripts, operational stores and credentials are
not copied into repository evidence. These results are not independent review,
owner acceptance, canonical publication or full ProjectGovernance readiness.
