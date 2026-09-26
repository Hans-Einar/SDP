# SDPTool producer review

| Field | Value |
| --- | --- |
| id | REVIEW-SDP-0001 |
| project | SDP |
| state | completed |
| source | KB-SDP-026; SPR-SDP-0001 |
| reviewer | Codex self-review, not independent review or owner approval |

## Scope

Review the executable producer contract against saved-preview, project selection,
configured plan viewing, general SDL tree/selected generation and KanBan/SDUI
journeys. Native XFMD code/visual acceptance remains external.

## Findings and disposition

| Finding | Resolution / evidence |
| --- | --- |
| Optional null/empty metadata paths could differ from the JSON schema | Reject them; regression cases for null, empty and boolean values |
| Registration needed an explicit initial model for view ip | Optional defaultModel validated against SDL entries; repo selects sdptool and tests reject unknown defaults |
| Tree construction repeatedly scanned facts and lacked an input-size guard | Index facts by ID; bound declarations/facts and combined inventory; retain reference/depth limits |
| CLI lacked a concise help entry | Added --help/help; machine commands retain JSON responses |
| Last valid preview must survive parse/render errors, not just stale selection | Dedicated failure-refresh test preserves prior entry bytes |

The executable consumer fixture passed discovery → inventory → selected bundle →
stale rejection → regeneration → configured viewer invocation through the compiled
CLI. Additional race tests and go vet passed. Real repository trees contain all
11 viewpoints (459 nodes for SDPTool, 3,616 for the shared SDL/SDUI model at review
time), with no unavailable optional services. The actual mmdr produced Mermaid/SVG
for SDPTool's structural model. Concept1 page structural Markdown export succeeded.
Node counts are observations, not fixed compatibility assertions.

Go toolchain: 1.27.1. mmdr executable SHA-256:
18b2574fcf6f5bbfa60ea96a97d45e0fc29211b150a7fc40842ad4e497cca8e3.
The tested candidate is the source in this milestone commit; uncommitted unrelated
sourceinput code is excluded and not imported by the facade.

## Review limits and remaining work

This is a self-review with executable producer/consumer-harness evidence, not an
independent review or owner acceptance. Native XFMD T4-M2, unsaved P0-M2 and advanced
T5 services stay with KB-SDP-017. The initial SDUI service is structural Markdown;
no new runtime, SVG widget backend or interactive controls are implied. Installation
facts are declared, not certified by the facade; full installer/profile migration
belongs to KB-SDP-028. Skills activation remains KB-SDP-027.
