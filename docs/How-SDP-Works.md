# How SDP Works

## 1. Source of truth

The repository is authoritative. Chat history, agent memory and local reasoning
are temporary. Decisions that matter must be written into repository documents.

## 2. Document progression

A normal project progresses through:

```text
Mandate -> Study -> Requirements -> Architecture
        -> Design Analysis -> Design -> Implementation
```

The progression is not a waterfall. Documents remain living records, but later
work must not silently contradict earlier decisions.

## 3. Horizontal design, vertical implementation

Architecture and design describe horizontal layers and contracts. Delivery is
organized as vertical Tiers, Sprints, Iterations and Slices that produce a small
working capability across the required layers.

## 4. Execution hierarchy

```text
Project
└── Sprint or Refactor
    └── Iteration
        └── Slice
```

A Slice is the smallest committed unit of work. It should define goal, scope,
invariants, non-goals, expected files, verification and completion signal.

## 5. Agent roles

### Master

Identifies the active work, maintains traceability, writes or refines slice
contracts, delegates implementation, requests independent review, verifies the
result and stops at the slice boundary.

### Worker

Implements one bounded slice. It must not broaden scope or redesign architecture
without returning the discovery to the Master.

### Reviewer

Uses a fresh context. It checks correctness, architecture, requirements,
traceability and evidence rather than trusting the Worker summary.

### Architect

Studies long-term structure and updates Mandate, Study, Requirements,
Architecture, Design and refactor planning. It does not implement product code.

## 6. Traceability

Recommended files:

- `Traceability/CurrentIndex.yaml`
- `Traceability/Relations.yaml`
- `Traceability/Ledger.ndjson`

`CurrentIndex.yaml` tells agents what is active. `Relations.yaml` connects stable
IDs. `Ledger.ndjson` records append-only events.

## 7. Verification

Completion requires evidence appropriate to the project. Typical product-code
checks include typecheck, tests, build and rendered or manual verification.
Documentation work should still verify links, IDs, YAML/JSON parsing and file
consistency.

## 8. Stop discipline

Agents do not automatically continue into the next Slice. The Master integrates
the result, records evidence, requests review and stops for acceptance.

## 9. Local adaptation

SDP defines a common operating model, not one rigid folder tree. Projects may
extend it, but repository-local instructions must identify deviations and keep
stable IDs and traceability intact.
