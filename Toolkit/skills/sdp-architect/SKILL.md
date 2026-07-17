---
skillId: sdp-architect
skillVersion: 1.0.0
minimumToolkitVersion: 0.2.0
capabilities:
- sdp.architecture.design
- sdp.release.architecture
compatibilityNotes: Initial formal version for SDP Toolkit 0.2.0; includes additive Steering interaction-record guidance.
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

## Steering interaction preservation

When acting as a supervising Steering Group or issuing a material assignment to
a repository-local Master:

1. Prefer a project-owned `SDP/Steering/` directory with an interaction contract
   and stable interaction IDs.
2. Record the exact Steering prompt in the repository before or at the time it
   is issued. Chat history is transport, not authority.
3. Preserve the Master's complete raw response verbatim when it is received.
   Do not rewrite it into a cleaner retrospective account.
4. Add the Steering assessment in a separate section only after inspecting
   repository evidence such as commits, diffs, verification, review and
   traceability.
5. Keep prompt, raw response and assessment distinguishable and independently
   statused. Pending sections are allowed.
6. Make corrections additively and date them; do not silently alter historical
   prompts or responses.
7. Relate material interactions to Features, Releases, Sprints, Iterations,
   Slices, Refactors or Fixes where known.
8. Add truthful Ledger transitions when the installed SDP event contract
   supports them. Until a formal event contract exists, the Steering interaction
   record remains authoritative and must not invent unsupported completion
   claims.
9. Never treat the Master's summary as sufficient proof of completion. Review
   the actual repository state.
10. Do not expose secrets, credentials or unnecessary private chat content in a
    public repository; redact before recording and state that a redaction
    occurred.

Prefer stable machine-readable contracts over duplicated display strings. Do not
create abstractions without a current boundary or more than one real consumer.
