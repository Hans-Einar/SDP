# BPI3-M2b generated lifecycle trial

These are frozen outputs from the source-built SDPTool, not authored blueprint or
navigation status. The synthetic non-Git fixture was produced by
TestAssignmentLifecycleCLI, exported to a temporary directory, then copied under
/tmp/bpi3-m2b-consumer/SDP to verify project-relative source/evidence references.
Absolute open targets identify that observed temporary fixture, not portable
installation paths. No product project was modified.

- assignments.json: four separate attempts: draft, completed, canceled, superseded.
- discovery.json: the same assignments in state groups plus retained revisions.
- broad-race-timeout.log: unsuccessful 180-second-per-package broad attempt under
  host contention. This is retained as a failure to establish evidence, never a pass.

The exact canonical event fixture is also retained for pure replay and installer
regressions in SDPTool/blueprintstate/testdata/lifecycle.ndjson. Reproduce a fresh
non-Git fixture from SDPTool with the pinned Go toolchain on PATH:

```sh
SDP_BLUEPRINT_LIFECYCLE_EXPORT=/tmp/fresh-blueprint-trial/SDP \
  go test ./blueprints -run '^TestAssignmentLifecycleCLI$' -count=1 -v
go build -o /tmp/sdptool-lifecycle ./cmd/sdptool
/tmp/sdptool-lifecycle /tmp/fresh-blueprint-trial/SDP model assignment list --json
/tmp/sdptool-lifecycle /tmp/fresh-blueprint-trial discover --json
```

The export destination must not already exist. The test runs a real small Go
behavior check for its code receipt; this does not prove production Ponsse or XFMD
implementation conformance. A and subsequent synthetic attempt identities are local
fixture data. Lifecycle code never executes receipt commands.
