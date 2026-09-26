# Protocol Buffers

[Catalogue](README.md) · Category: **Data/service IDL** · Research: **2026-09-10**

## Purpose and abstraction level

Typed messages with binary serialization and code generation across implementation languages. RPC tools can use service descriptions.

## Model mechanisms

Identity: packages, message/type names and field numbers. Relations: type references/services. Contracts: data types/wire format; agree semantic units separately. Views: external documentation tools. Extension: options/generator plugins.

## Strengths and limitations — our assessment

**Strength:** Field identities and documented evolution rules illustrate controlled contract change.

**Limitation:** Wire compatibility differs from semantic preservation. Switching meters to centimeters can be a serious defect without changing types.

## History, change and transitions

Do not reuse deleted field numbers; reserved protects removed fields. This supports contract evolution, not complete responsibility history.

## Illustrative example

Proto3 message. Meters are an explicit field-name convention, not a universal SDP unit rule. The example has not been parser/runtime tested.

```proto
syntax = "proto3";
package measurements.v1;
message Length {
  double value_m = 1;
  reserved 2;
}
```

## Tools, maintenance and terms

Proto3 guide available; other editions exist, and proto3 is not claimed to be latest. Repository not archived; metadata returned NOASSERTION, so check actual LICENSE/NOTICE/runtime terms.

## Primary sources

All sources checked on 2026-09-10; see the catalogue methodology for evidence and licensing limits. Translation does not refresh these dated findings.

- [Primary documentation](https://protobuf.dev/programming-guides/proto3/)
- [Official repository; metadata checked through GitHub API](https://github.com/protocolbuffers/protobuf)
