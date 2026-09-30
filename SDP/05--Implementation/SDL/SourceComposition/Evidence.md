# Source composition implementation evidence

PLAN-SDP-0012; baseline 32aa68a. Results will be appended at actual milestones.
No new-language implementation acceptance is claimed by plan creation.

## SSI1-M1 frontend

Starting commit 70aa48e plus milestone diff. Commands: GOMAXPROCS=2 go -C SDL/go
test -p 2 ./parser ./sourcegraph ./cmd/sdl (Go 1.27.1 local toolchain): pass.
CLI check examples/source-composition/System.design: valid design-core/0.6,
System Frontend, three original files, revision
eab3d03d82afa4fdb4627a1870be457029a936e0ad9c9d8e7f92543b14956862.
Late-root reuse, context isolation, original spans, source cycles/diamonds,
invalid membership/names/profiles, paths, limits, output protection and CLI
fragment/file-map operations are covered. Initial test fixture omitted the
existing required contract completeness property; adding the legitimate property
resolved that fixture failure without relaxing semantic validation.
Consumer migration and independent review are still pending.

## SSI2-M1 consumers and review corrections

Candidate: c090b28 plus the SSI2 milestone diff. Go 1.27.1, GOMAXPROCS=2,
`go -C SDL/go test -p 2 ./parser ./sourcegraph ./viewpoint ./documents ./broker ./cmd/sdl`
and `go -C SDPTool test -p 2 ./...`: passed. Additional sourcegraph, documents,
broker, CLI and SDPTool composition tests passed after addition.

All consumers use the graph revision and original-source provenance. Tests cover
registered System/profile mismatch across tree/select/ViewPlan, non-entry change
invalidation, old revision rejection retaining the previous document, dependency
output protection, broker cache invalidation and stale render rejection preserving
reader leases. The real Frontend pilot preserves every baseline fact; file moves
change revision/provenance while fact identities remain stable.

Independent Reviewer initially requested changes: diagnostic amplification,
aggregate-token bypass in pure compilation, missing opposite-edge cycle evidence,
and missing ViewPlan registration validation. All four were corrected with
regression tests. Diagnostic neighbourhood evidence is bounded to 100 diagnostics
and 8 related locations each; it is not exhaustive causal analysis. Repeated-source
warnings are capped at 100. The cache retains grammar token counts. Final review
and full verification remain SSI3 work.
