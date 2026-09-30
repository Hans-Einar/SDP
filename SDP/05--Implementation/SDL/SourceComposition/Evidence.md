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
