# RP3-GO — Remove the retired installation toolchain

Owner correction and SDP 1.0.0 selection: 2026-09-28. Previous 0.2.2 candidate
checks remain historical and do not govern this release.

## Removed execution paths

All tracked .ps1 files, including the legacy bootstrap copy, are removed. The
old profile builder, generated artifact, shell-dependent Python integration and
recovery runners, and old skill verification runners are removed as well. CI has
no jobs or command paths invoking that runtime. Current authored inventory now
lives in SDPTool/profiles/payload.json; strict Go decoding consumes it directly.

Legacy template copies and language-neutral old JSON schemas/expectations remain
inspectable data. Historical prose and immutable release assets retain provenance;
links to removed scripts use their exact pre-retirement Git revision. The L1
preservation check explicitly records the one owner-retired archived executable,
without weakening checks on other frozen records.

## Current obligations and coverage

| Obligation | Go evidence |
| --- | --- |
| Strict records, path/collision checks | records_test.go |
| Deterministic no-write preview, explicit adoption, history checks | plan_test.go |
| Apply, repeated no-op, known upgrades, drift/lock rejection | execute_test.go |
| Process-exit matrix and recovery edit rejection | recovery_test.go |
| Packaged child and signed distribution | signed_test.go and bootstrap tests |
| XFMD-shaped adoption and preserved local data | xfmd_trial_test.go |
| Complete canonical template inventory, reproducible descriptor, predecessor hashes | tools/profile tests |

The removed runtime checks exercised the retired engine. They are not silently
counted as Go results; the current Go suite exercises actual preview/apply/recovery.
Fresh signed upgrades from published 0.2.0 and 0.2.1 and actual XFMD upgrade must
be rerun with the final 1.0.0 descriptor before closeout.

Local checks use GOMAXPROCS=2 and go test -p 1 to bound load. The full current Go
suite passes, including process-exit recovery; 85 Python data/document tests pass
without runtime-dependent skips. Go profile tests verify that new canonical
templates cannot be omitted from release payloads and descriptors are reproducible.
Exact candidate, CI, independent review and publication evidence follow separately.
