---
skillId: sdp-architect
skillVersion: 1.0.0
minimumToolkitVersion: 0.2.0
capabilities:
- sdp.architecture.design
- sdp.release.architecture
compatibilityNotes: Initial formal version for SDP Toolkit 0.2.0.
---

# SDP Architect

Use this skill for long-term method or product structure, not product-code
implementation.

1. Read Mandate, Study, Requirements, Architecture, manifests and active
   traceability.
2. Separate public release contracts from internal development coordinates.
3. Define boundaries, ownership, data formats, compatibility and migration rules.
4. Record alternatives, tradeoffs, assumptions and future Analyzer contracts.
5. Update the appropriate SDP documents and relations.
6. Produce implementation-ready boundaries without implementing product code.
7. Stop with decisions, open questions and the recommended next Slice.

Prefer stable machine-readable contracts over duplicated display strings. Do not
create abstractions without a current boundary or more than one real consumer.
