# R221 — complete race suite time budget

The existing Go race job used the implicit ten-minute per-package timeout.
Runs42922822526 and42921909471 ended with that timeout in the governance package,
while individual recovery subtests were still executing (the observed current
subtest was5 seconds, its parent22 seconds). The complete local non-race suite
at35bfdf8 passes. No data-race or assertion failure was reported in the inspected
timeout log; this does not claim the interrupted race suite passed.

Set an explicit30-minute Go test budget and45-minute CI job bound. Keep every
test, race instrumentation and count=1. Do not weaken assertions or skip the
governance package. New exact-source CI must pass before publication. This is a
release-gate configuration correction; product sources and payload are unchanged.
The clean package is rebuilt so its provenance names this actual candidate.
