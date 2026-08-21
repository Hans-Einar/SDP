# SDP vNext provisional pilot contract

Status: **provisional / experimental / pilot-only**

This directory is the bounded output of
[Hans-Einar/SDP Issue #7](https://github.com/Hans-Einar/SDP/issues/7). It turns
the Steering-accepted Issue #5 study into a contract that real repositories can
pilot. It does not change the canonical Toolkit, freeze a future schema, or
claim compatibility with a later vNext release.

The words **MUST**, **MUST NOT**, **SHOULD**, and **MAY** in this directory state
pilot rules only.

## Reading order

1. [WorkflowContract.md](WorkflowContract.md) defines work ownership,
   cardinalities, evidence state, reopening, Steering, and release inclusion.
2. [WorkDomains.md](WorkDomains.md) defines scoped monorepo domains, stable
   identity, references, discovery, moves, and collision rules.
3. [IssueContract.md](IssueContract.md) is the reusable human-readable GitHub
   Issue contract.
4. [ConcurrentAssignments.md](ConcurrentAssignments.md) defines per-Issue
   assignments, reservation/refreeze, stale-base handling, and convergence.
5. [Profiles.md](Profiles.md) keeps ceremony proportional.
6. [Examples.md](Examples.md) walks through the machine-readable examples.
7. [HSX-Pilot-Handoff.md](HSX-Pilot-Handoff.md) is the concrete first real
   pilot handoff; it is evidence about HSX, not a change to HSX.
8. [OpenQuestions.md](OpenQuestions.md) records choices intentionally left for
   real pilots and later schema work.

`templates/` contains illustrative pilot-only record shapes. `examples/`
contains positive scenarios. `fixtures/negative/` contains deliberately invalid
cases and their required diagnostic codes. JSON is used as the machine-readable
pilot encoding because Python's standard library can parse it deterministically;
this does not decide the eventual canonical serialization. The existing
`Steering/Assignments/ISSUE-007.yaml` is JSON-compatible YAML for the same
reason.

Issue #7 dogfoods the model through `Studies/STU-007.json`, the compact
assignment, the single default-domain declaration, and
`Steering/Reservations/RSV-ISSUE-007-001.json`. The assignment binds the real
reservation file by canonical JSON SHA-256; the validator resolves the Study,
domain, paths, IDs, base, dependency, and convergence contract together.

## Validate

From the repository root, run:

```powershell
python SDP-vnext-pilot/validate_pilot.py
```

The validator uses only the Python standard library. It validates templates,
positive examples, representative negative fixtures, dogfood record/reservation
resolution, local Markdown links, required status markers, and trailing
whitespace.

## Boundary

This corpus does not modify or authorize changes to HSX, gh-sdp,
SDP-Analyzer, PR #4, canonical Toolkit contracts, Toolkit release state, tags,
or GitHub Releases. Tool behavior described here is a future consumer contract,
not an implementation in this Issue.
