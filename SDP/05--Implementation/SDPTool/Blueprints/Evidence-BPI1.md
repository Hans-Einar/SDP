# BPI1-M1 evidence — 2026-10-07

Candidate: base 3d265d1 plus the files in this milestone commit, branch
sdp/blueprint-implementation. Go 1.27.1, Linux. No product release.

## Delivered

SDL/go/blueprint is an I/O-free structural analyzer over existing checked
sourcegraph snapshots. It provides typed differences, retained per-side provenance,
conservative relation closure with explicit cuts, change permissions, scoped
protections, unknowns and coverage. Task permissions do not authorize code execution.
See the package README for the implemented contract and boundaries.

## Verification

- go test ./blueprint: pass.
- go test -race ./parser ./sourcegraph ./blueprint: pass.
- go vet ./blueprint: pass.
- Independent reviewer repeated blueprint tests with -count=1 -race and vet: pass.

Real parser-backed reduced MVP1 cases cover preservation/ownership drift, source
moves/formatting, permitted-message participants missing on one or both sides,
expected identity, limits, conflicting rules, explicit permissions and protection
scope. Additional cases cover System/profile mismatch, declaration kind/rename,
System boundary cuts, task path validation, canceled calls and unsupported policy.
The closure cycle fixture is an accepted directed diamond, cyclic when traversed
bidirectionally. Typed operand tests are unit evidence rather than parser workflows.

Initial tests failed on an incorrect fixture System ID, an invalid directed cycle
and a renamed file that broke the parser's path-target convention. These were
corrected without changing language validation.

## Independent review

Agent /root/bpi1_review, fresh context, read-only, reviewed owner intent/design
before code. Initial changes-required findings: missing-both-peer detection,
explicit change/protection contract, typed semantic output and coverage.
All four were addressed with tests. Final verdict: approved for bounded BPI1;
no additional material regression found in focused re-review.

Reviewed SHA-256:
- analyze.go: 39be5d7a367a66b548e96f0986f5a874c782810f620c9ac00c624bf79ee686b2
- analyze_test.go: 83ee92fef269503f9973f3d0675799a7b8f29a1bc359605ace989f08ee74903b

No production command, document publication/catalogue, assignment lifecycle,
full model-area provenance, runtime behavior, code-conformance or native XFMD
acceptance is proved. Those remain BPI2/BPI3 or separately owned features.
