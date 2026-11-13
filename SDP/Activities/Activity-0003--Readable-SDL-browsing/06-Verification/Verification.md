# Activity3 verification

| Field | Value |
| --- | --- |
| id | SDPTOOL-VER-READING-0003 |

## Qualified product

Independent reviewer approved producer `5802355a8af6a89f9ea3fc035198c32d46d35567`
and native browsing implementation `d530aa7` for code/product qualification.
That initial product disposition was followed by final release qualification below.

Producer ownership and acceptance are defined in ../04-Design/ImplementationContract.md
and canonical SDP/02--Requirements/SDPTool.md (REQ-READ-001 through REQ-READ-008).
All eleven viewpoints have ordered unique collection routes, meaningful landings,
chapter collections, contextual details, readable statement links and separate
registers. Systems own their explicit and qualified dependency closures; orphan
sources render local syntax with honest context limitations. Mode, source revision,
exact source spans and producer/consumer authority remain checked.

## Retained regression evidence

Logs are in [B4](B4/). Commands ran on AlmaLinux10.2 amd64, Go1.26.10.
This is layered regression evidence, not a claim that every command ran on final HEAD.

| Check | Input and result |
| --- | --- |
| SDPTool `go test ./... -count=1 -timeout30m` | Initial integrated candidate: root and all packages except activities passed; two real-history test fixtures assumed only Activity1/2 existed |
| SDPTool `go test ./activities -count=1` | Fixture isolation repaired in99ea777; full package passed11.065s; original failed log retained |
| SDL `go test ./... -count=1` | Integrated reading implementation passed |
| SDL `go test -race ./viewpoint ./documents -count=1` | Final structural panel implementation passed; unchanged by later SDPTool-only guards/membership |
| SDPTool root targeted race (Reading, Deferred, Activity projection/expansion/overview, Source, Overview) | Passed185.809s |
| `go vet ./...` in SDPTool; `go vet ./viewpoint ./documents` in SDL | Passed, empty logs retained |
| SDL `go test ./documents -count=1` | Final aggregate/partitioned structure passed2.262s |
| SDPTool final source bounds checks | 7370e29 passed, including128files/32MiB limits |
| SDPTool `go test . -run 'TestReading\|TestDeferred' -count=1` | Final5802355 passed16.108s; regression follows nested qualified paths through canonical sourcegraph.Dependencies |
| Actual HEOS producer selection matrix | 7370e29: Micro87, Mini75, Base75 selections;237 total passed, retained snapshots and no visible machine IDs |

The final5802355 delta only fixes qualified-path source membership. HEOS explicit
include inputs and rendering were unchanged, so native and matrix evidence reuse
is justified. Test fixture repairs preserve the original pre-automation history
bytes rather than weakening lifecycle checks. Compatibility tests explicitly
request retained old overview/diagram URIs instead of asserting removed tree rows.

## Native workflows and independent inspection

XFMD commitd530aa7 retains native fixtures, logs, screenshots, SVGs and actual PDFs
under SDP/05--Implementation/SDL-Navigation/evidence/Reading-001. Native regression
passed90/90. Additional fixture/Micro/Base checks opened all11 landings, compared
advertised/admitted visual counts, activated before expansion, expanded without
selection, followed exact source spans, went Back and rejected stale siblings.

Independent review inspected final Kernel FileApi page4, Base collaboration page2,
Micro channel page4 and Contracts RegistryBindings page13 at normal A4 proportions.
The original dense PDF was rejected and replaced with bounded channel/mode panels,
readable labels and adjacent payload tables. All actual bindings are preserved
exactly once by regression. Unit containment, own-scope attachment and mixed-mode
proof filtering findings have passing corrections. Product approval does not
claim owner personal review or a public package upgrade.

## Remaining release gate

Exact0.10 RPM candidate, pinned producer archive/provenance, package probes, server
signatures, public download validation, signed public predecessor DNF upgrade and
fresh installation are pending. Record actual outputs in ReleaseChecklist and the
XFMD release evidence before completing B5 or Activity3. No main merge is authorized.

## Final scaffold correction

During release reconciliation, the Activity README creation template was found to
violate the owner's explicit generic-directory-guide requirement (Activity Design
sections86–92/457). Producerc3463b78d53ebf984b0f553ebe93a556fb3ac6a8 embeds the
standard guide with navigation for all six phases and root Sessions, preserving
Activity-specific need only in Mandate.md. The existing create/retry regression
now compares two distinct Activities and verifies separate mandate preservation;
full activities tests pass8.459s (B4/scaffold-final.log). Existing documents are not
silently overwritten by the command; Activity3's guide was explicitly reconciled.
This delta leaves qualified reading/native behavior unchanged. The private RPM
producer pin must usec3463b7; independent delta review approved this correction
and independently reran the extended creation/retry regression.

## Delivered application and final disposition

[ReleaseChecklist](ReleaseChecklist.md) and [independent review](B5/Review.md)
record the completed delivery. Native source8546e956d1a941ece9195d335b7196abea905289
and private producer c3463b78d53ebf984b0f553ebe93a556fb3ac6a8 ship together in signed
XFMD0.10.0-76.el10_2. Container89/89 tests, installed workflows, local upgrade and
fresh-install checks pass. Public signed metadata/package verification, actual
DNF upgrade from0.9.0-71 and fresh public install pass the same native workflows.
The independent reviewer approved final public qualification with no findings.

The owner authorized autonomous completion/publication and had no time for personal
review. Closure records the objectively verified authorized delivery; it does not
claim that personal review, a main merge, an accepted SDL model release or a public
SDP Toolkit installation release occurred. Initial pending/checkpoint sections
above are historical, superseded by this final disposition. No implementation or
publication obligation remains in PLAN-SDP-0028.
