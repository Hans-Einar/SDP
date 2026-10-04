# MG3 filesystem experiment

This standalone Go module is design evidence for PLAN-SDP-0016 MG3-M1, not a
production ModelGovernance implementation or an SDPTool command. All test files
are created under Go test temporary directories. No project Git repository is
required by the store; one comparison test invokes `git merge-file` on three plain
files to measure a possible text-merge primitive. That is not a selected production
dependency. The module uses the existing yaml.v3 dependency version.

```sh
go test -race -count=1 -v ./...
go vet ./...
```

Thirteen top-level tests exercise file after-images/deletions, dirty-state recovery,
portable copies, fault-injected process interruption, lineage and promotion. One
merge test has five subcases. Commit payloads use typed YAML records, but the schema
is deliberately a subset (`mg-probe/0.1`), not the full proposed `sdp-model/0.1`.
Fixture identities are short labels, not production UUIDs. `head.yaml`, `sources/`
and the transaction primitive simplify the experiment and are not final artifact
layout choices. Do not copy this prototype into the product without its missing
controls. In particular it is not safe for untrusted metadata or concurrent writers.

The harness returns reconstructed file maps; it does not replace a user's live
source tree. Restore's logical preservation is tested, not transactional multi-file
replacement. Resume copies to a temporary validation area for simplicity; production
should not use recursive store copies as its commit verification mechanism.

[Proof and limitations](../../SDP/04--Design/SDPTool/ModelGovernance/Proof.md).
