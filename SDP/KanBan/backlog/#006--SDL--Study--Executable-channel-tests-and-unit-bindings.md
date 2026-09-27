# Study executable SDL Channel tests and real Unit bindings

| Field | Value |
| --- | --- |
| id | KB-SDL-006 |
| project | SDL |
| type | Study |
| CardState | backlog |
| Systems | SDL, SDPTOOL |
| created | 2026-09-26T10:12:02Z |
| source | Owner clarification 2026-09-26 after BP1 blueprint study |
| next_review | When selecting PLAN-SDP-0001 or the next SDL runtime study |

## Owner intent

Use SDL to drive tests through Channels: supply data satisfying the input contract,
invoke actual Unit code and verify processed outputs against the output contract
and expected behavior. Assignment bundles carry NOW and TARGET models so a reviewer
can exercise baseline code before work and resulting code afterward. Preserve the
surrounding system's unchanged contracts as regression obligations.

The owner also proposes exploring a companion language for Unit behavior and
future code generation. Go is a candidate because the runtime and intended system
implementations use Go; this is an open option, not an adopted Go interpreter.

## Existing implementation to inspect and reuse

[action-core](../../../SDL/docs/profiles/SDL-Executable-Action-Profile.md) already
has typed input/output records, symbolic Go function bindings and runtime execution.
Inspect SDL/go/runtime/engine.go, parser/actions.go, codegen/generate.go and existing
runtime/codegen tests. Structural Channel/scenario contracts currently validate
design facts; their automatic mapping to executable actions is not implemented.
The existing generator emits model constructors and leaves domain handlers handwritten.

## Questions and bounded study output

- Define explicit Channel/message/Unit-to-handler bindings, input/output conversion,
  ownership and correlation. Compare in-process compiled Go handlers, process-boundary
  adapters, existing action-core and a companion interpreter only where needed.
- Separate interface validity from correct results: define authored scenarios,
  expected outputs, failure assertions and state preconditions. Do not infer a
  domain oracle merely from a Dataset/Datagram shape. Include invalid-input tests.
- Specify state isolation/reset, determinism, clocks, cancellation, ordering,
  errors and asynchronous channels where relevant; do not imply all tests are
  isolated unit tests when they exercise several collaborating components.
- Reuse real implementation through an explicit binding; distinguish a simulation,
  test double and production adapter. No machine/Ponsse side effects during research.
- Design NOW/TARGET model and baseline/result code identities, test selection,
  accepted pre-existing failures and retained-neighbor regression checks.
- Coordinate code identity tags with KB-SDP-004 and PLAN-SDP-0001. Tags locate code;
  executable evidence tests behavior. Neither alone proves total conformance.
- Explore code generation for registration/adapters or test scaffolds without
  claiming SDL can generate unspecified domain behavior. Go runtime implementation
  does not require the companion language to be Go or require interpreting Go source.

Deliver a Study and proposed phased plan using one real bounded Unit/Channel flow,
explicit expected results and a faulty-result negative case. Inspect capabilities
before creating another runtime. Record unsupported Channel features and avoid a
silent conversion from structural design-core into execution semantics.

## Scope and relationship

This is registered, unselected research. No parser/runtime, code-tag syntax,
companion language, source interpreter or generated production code is authorized
by registration. Coordinate the
[blueprint DesignPlan](../../04--Design/SDPTool/Blueprints/Plan.md) and
[Traceability contract](%23004--Proposal--Design-traceability.md).
The [owner clarification](../../04--Design/SDPTool/Blueprints/Study.md#owner-clarification-after-bp1--2026-09-26)
is the primary source for the combined direction.

## Worklog

2026-09-26T10:12:02Z: Registered in backlog; EVT-KB-SDL-000027.
