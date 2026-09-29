# Ecosystem model verification — E3

Scope: [PLAN-SDP-0008](Plan.md), [catalog](../../SDL/Catalog.md), KB-SDP-039.
[Evidence.json](Evidence.json) pins each checked model, registration, verification
recipe, Go tool binaries and the unchanged tracked implementation basis. It is
structural/model evidence, not a release record or runtime acceptance.

## Delivered model basis

Three ecosystems contain 17 bounded system models. Landscape.design adds one
cross-system integration-context model. Together they contain 658 declarations,
1,535 facts and 22 declared scenarios. Each source is independently accepted by
the existing design-core 0.5 Go checker; no includes, new parser, concatenated
source fragments or experimental sourceinput code is used.

Existing code and proposed responsibilities are identified in system READMEs and
activity status facts. After review, SDUI models gained explicit implemented
capability/planned packaging activities so generated delivery views also carry
that distinction. Valid scenarios remain illustrative paths, not executed tests.

## Reproduced checks

The final integrated run used Go 1.27.1 and [verify.sh](verify.sh), building the
current tracked SDL and SDPTool commands. For all 18 entries it performed:

1. Semantic/canonical source checking and AST export.
2. Canonical formatting and rechecking of that output.
3. SDPTool discovery and model-specific navigation tree generation.
4. Revision-bound selection using a target returned by the tree, with a published
   entry.md checked on disk.
5. Static Markdown/Mermaid generation for VP01, VP02, VP06 and VP08.

All passed. The selected exports contain 148 diagrams and 674 bundle files in a
disposable output area, not committed derived documents. The real SDL toolkit
produced the diagram bodies. Workers additionally generated complete static
viewpoints for their system models before the final SDUI maturity annotations.
The final run separately checked both old detailed model trees and verified that
old registration entries, default, project identity and all SDUI entries remain
unchanged. Detailed legacy model bytes remain unchanged.

Reproduction from repository root:

```sh
SDP_GO=/absolute/path/to/go bash SDP/03--Architecture/Ecosystems/verify.sh /tmp/new-empty-ecosystem-review
python3 SDP/ProjectManagement/validate.py
git diff --check
```

The selected Go executable must be compatible with the repository modules.
The output must be an empty directory outside the repository; it is never
silently replaced. The generated index is a convenience wrapper linking toolkit
output. It does not add an ecosystem UI or URI handler to XFMD.

## Review

A fresh independent Reviewer approved the initial model, source-mapping and
navigation candidate without material findings. It independently ran all 18
entries through Go/SDPTool, checked 277 local links and all seven SDL worker
model hashes, and inspected selected runtime, generation, host, Markdown and
board boundaries against their implementation. Its suggested maturity improvement
for SDUI was applied as described above and the integrated checks were rerun.
Final delta and lifecycle consistency are reviewed separately before closeout.

## Limits and retained work

No process classifier, routine engine, MCP adapter, app-server client or KanBanTUI
was implemented or launched. No account/authentication, native GUI or domain
execution test was run. Independent package/release extraction, actual cross-file
System semantics and complete ABI/wire contracts are future work. Templates,
released artifacts and external repositories remain unchanged.

The earlier detailed models stay authoritative for their existing behavior. This
catalog establishes tool boundaries, not a lossless migration or a complete
blueprint/source-code conformance system. The management ledger records this
architecture delivery separately from future product implementation. Unrelated
sourceinput and Node workspace files remain excluded.
