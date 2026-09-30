# SSD1-M1 design evidence

Baseline: 84fbdb34f7b87f2ef23547cc4b1684322916a342, design branch
sdp/sdl-source-sets/design. PLAN-SDP-0010 delivers a recommended contract and
acceptance design, not new parser/runtime functionality or owner approval of syntax.

## Inspected implementation

[Source manifest](evidence/source-manifest.json) records 15 exact file hashes.
The two untracked sourceinput files are explicitly marked inspected-only; they are
not included in the implementation candidate or modified by this plan.

| Observation | Evidence and design consequence |
| --- | --- |
| Parse accepts only 0.5; Check validates and checks canonical text | parser/parse.go: new profile dispatch and whole-set checking must precede per-file canonical checks. |
| Spans have no file identity; diagnostics have one span | parser/ast.go: new file-aware output and related spans required; no concatenated offset workaround. |
| System is absent; containers act as unit-compatible types | parser/vocabulary.go and validate.go: explicit System signature, membership semantics and separate validation needed. |
| Canonical hardcodes 0.5 | parser/ast.go: format profile must come from validated input, with an explicit multi-file output contract. |
| Views reparses text and hashes one string; facts get ordinal IDs | viewpoint/model.go: checked-result entrypoint needed to retain original provenance and coherent revisions. |
| Broker hashes/rechecks only one source file | broker/broker.go: aggregate input recheck must cover non-entry edits and cache hits while retaining sequence/lease handling. |
| SDPTool assumes source profile/revision in registration and preview | project.go, preview.go, tree.go: facade must consume the shared frontend and report actual profile; it must not invent a second loader. |
| Snapshot accepts one design string | snapshot/snapshot.go: require explicit migration or retain clearly limited standalone capability. |

## Executed baseline probes

[Probe output](evidence/probes.json) retains commands/results and the built binary hash.
Build from SDL/go:

```sh
GOMAXPROCS=2 /home/warloc/.local/share/sdp-toolchains/go1.27.1/go/bin/go build -p 2 -o /tmp/sdl-ssd1 ./cmd/sdl
```

From repository root:

```sh
/tmp/sdl-ssd1 check SDP/SDL/SDL/Frontend/System.design
/tmp/sdl-ssd1 format SDP/SDL/SDL/Frontend/System.design
/tmp/sdl-ssd1 check SDP/04--Design/SDL/SourceSets/examples/frontend/System.design
```

- Existing Frontend 0.5 check succeeds and format output is byte-identical.
- Proposed 0.6 entry fails with UNSUPPORTED_VERSION, as expected. This prevents a
  claim that the proposed examples already execute through the released parser.
- A multiset comparison of non-header lines in the baseline versus Features and
  Containers/Frontend confirms **all 57 declarations/statements are retained exactly**.
  System.design adds only the new header, System declaration and membership fact.
  This is an inventory comparison, not a substitute parser or semantic validator.
- All 15 inspected source hashes match, including the unchanged untracked draft.

One initial fixture-generation command used SDL/go as its working directory and
failed to find the baseline. Its accidental documentation directory was removed;
fixture generation was rerun with the absolute repository root. No product file
was modified. The successful build/probe evidence above uses the stated paths.

## Document and history checks

Management schema/history, Toolkit repository validation, affected local Markdown
links, historical ledger-prefix preservation and git diff --check are checked at
closeout. The plan/Session/card carry the outcome; these checks do not establish
new-profile implementation acceptance. See the plan delivery record for results.

CurrentIndex has no typed-plan active field; no fake Slice or release coordinate
is inserted. The Session and primary card point to the plan; the system design
record is linked from Traceability using its existing generic event envelope.
Management-only transitions remain exclusively in ProjectManagement.

## Limits and handoff

No Go source, language specification, installed template or live model registration
changed. No rendering/native XFMD test, independent reviewer, cross-System linking
or implementation acceptance is claimed. SS01–SS26 are future checks. The active
card still owns implementation and its explicit public-linking disposition.
The recommended closed-System boundary is a staged design choice, not a claim
that the entire historical workspace/export proposal has been fulfilled.
