# Draft Standard Document Procedure

Status: working draft  
Audience: agents first, humans second

## Purpose

The Standard Document Procedure, or SDP, is a repository-local documentation
and execution system. Its purpose is to reduce drift between discussion,
requirements, design intent, implementation, verification, and handoff.

Every project can adapt the SDP locally, but the common rule is that repository
documents are the source of truth. Chat memory is not enough.

## Baseline Folder Structure

Recommended SDP folder structure:

- `01--Mandate`
- `02--Study`
- `03--Requirements`
- `04--Architecture`
- `05--DesignAnalysis`
- `06--Design`
- `07--Implementation`
- `Traceability`
- `Verification`
- `Sprints`
- `CodeReview`
- `Refactor`
- `Instructions`

The numbered folders move from intent to implementation readiness:

- mandate explains why the project exists
- study captures investigation and tradeoffs
- requirements define numbered, stable requirements
- architecture defines major system boundaries
- design-analysis maps horizontal layers, contracts, and Tier fan-out
- design records chosen detailed designs
- implementation defines ordered delivery plans and slices

## Traceability Layer

The traceability layer should normally contain:

- `CurrentIndex.yaml`
- `Relations.yaml`
- `Ledger.ndjson`

The current index names the active sprint, iteration, slice, and stable IDs.
Relations map requirements, design, slices, review, and verification.
The ledger is append-only history. It should parse as NDJSON.

## Sprints And Slices

There should normally be one active sprint and one active iteration.

A slice contract should state:

- goal
- why now
- files or modules expected to change
- invariants
- non-goals
- traceability IDs
- verification evidence
- completion signal

Product-code implementation should be delegated to worker agents. Review
should be performed by separate reviewer agents. The Master agent coordinates
documents, traceability, verification, and handoff.

## Design And Implementation Rule

Design horizontally. Implement vertically.

Horizontal design separates concerns by architectural layer:

- frontend shell
- backend API
- identity/auth
- persistence
- artifact service
- workers
- visualization
- domain analysis
- external integrations
- verification

Vertical implementation delivers Tiers. A Tier cuts through the layers needed
to produce one coherent capability.

## Verification

A slice is not done until the repository records:

- what changed
- which IDs it touched
- review outcome or why review was not required
- verification evidence
- ledger entries
- handoff state when useful

Documentation-only slices may use file existence checks, cross-reference
checks, YAML/JSON parsing, and ledger validation.

Product-code slices should use real build/test/manual evidence appropriate to
the change.

## Living Document Rule

When a repeated failure mode appears, add a concrete rule here or to the
repository-local SDP instructions.
