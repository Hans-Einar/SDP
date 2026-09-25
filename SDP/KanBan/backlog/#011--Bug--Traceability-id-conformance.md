# Resolve older Traceability IDs against the Toolkit contract

| Field | Value |
| --- | --- |
| id | KB-SDP-011 |
| CardState | backlog |
| ScrumId | SCRUM-SDP-0001 |
| Systems | SDP |
| project | SDP |
| type | Bug |
| created | 2026-09-24T14:06:02Z |
| source | R1 baseline check; continuation of a known CI mismatch |
| next_review | At the next prioritization after R1 |

## Observation

RepositoryValidation and contracts CI fail on nine older Issue #5 IDs: SPR-SDP-005, four ITR-SDP-005-* and four SLC-SDP-005-*. The schema requires Sprint-, SPI- and SPS-style IDs. R1 baseline and post-migration checks report the same failures. [Baseline](../../Maintenance/R1/toolkit-baseline.txt).

## Scope and acceptance

Decide whether these are supported external IDs or data requiring migration; provide one consistent solution with alias/traceability rules. Do not rewrite append-only ledger history or old evidence to make CI pass. Follow incoming references in records, evidence, issues and schemas before changing anything. The full Toolkit validator and unittest suite must pass without exceptions that hide errors. Physical directory housekeeping did not fix this.

## Backlog review — 2026-09-25

Retained as a bounded compatibility/CI Bug. Repairing historical ID interpretation must not be hidden inside a future evidence redesign or rewrite append-only records.

Recorded 2026-09-25T01:41:26Z, Codex, EVT-KB-SDP-000079. CardState remains backlog.

## Scrum-0001 review

Keep the nine existing Toolkit ID/CI failures as a bounded compatibility bug. Coordinate new system prefixes without rewriting historical IDs or hiding failures.

EVT-KB-SDP-000118; next review at the next selection or relevant dependency delivery.

## SK1 baseline update

The actual baseline now has 38 errors; see [exact output](../../Maintenance/SK1/toolkit-baseline.txt).
They predate SK1 and include reciprocal-relation errors as well as ID formats.
This card retains repair ownership. EVT-KB-SDP-000159.
