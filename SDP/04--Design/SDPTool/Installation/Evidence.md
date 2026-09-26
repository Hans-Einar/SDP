# IPD design evidence

Evidence distinguishes parser/model consistency, generated projections and authored
contract review from installer execution. No Go installation or live XFMD upgrade
is performed. The current tool was built from SDL/go at planning baseline db7d82d
with the existing Go 1.27.1 linux/amd64 toolchain into /tmp/ipd-sdl.

## IPD-1-M1 — installation responsibility model

The canonical model adds the project maintainer, four install/upgrade/adoption/
recovery use cases, SDPTool installation units, the external gh-sdp process and
external release artifact service. The latter is a collaborating service boundary,
not a new SDP-owned server. Five new delivery activities are explicitly planned.
Library/package boundaries are not modeled as deployed containers.

Commands: `sdl format`, `sdl check`, `sdl ast` on the canonical SDPTool.design.
The first draft exposed declaration-before-fact and canonical-order requirements;
the formatter normalized it, then check/AST passed with no diagnostics and 90
symbols. A set comparison preserves every statement in the db7d82d model.
The AST is parser output, not a manually constructed JSON model. Detailed channels,
records and scenarios belong to IPD-2; no protocol implementation is inferred.


## IPD-2-M1 — contracts and explicit exchanges

Working candidate: 9607fdc plus IPD-2 diff. Canonical check/AST passes with 188
symbols; every IPD-1 statement is preserved. The extended model contains twelve
channels, closed request/result protocols, five scenarios and four datasets.
The first draft reused field declarations across contracts; SDL correctly rejected
that cardinality. Contract-scoped field names fixed the model without parser changes.
No AST or generated diagram is handwritten. Signatures, JSON validation, path safety
and mutation semantics remain authored design requirements, not parser proofs.

Contract/scenario review records preview-only default, root-bound apply, explicit
manual adoption, signed descriptor design, thin binary delegation, incompatible
legacy journals, forward-only recovery and no fixed packet encoding. The first
implementation milestone must freeze bounded wire schemas and conformance cases.
Actual XFMD read-only status was refreshed to b95a4bb (clean); no files were changed.
Independent reviewer/owner acceptance is not claimed.
