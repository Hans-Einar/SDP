# L1 verification

## L1-M1 — English policy and KanBan

Translated all 17 cards and the shared/SDL/SDUI workflow entry points, card
history, merge/split contract and template. Root and SDUI agent instructions and
the Toolkit-distributed agent template now require English documentation.

The translation preserves proposal status, technical boundaries, IDs, dates,
old test claims and source references. Historical ledger bytes remain unchanged;
new reviewed events record translation. Code snippets retain syntax; illustrative
narrative strings in the requirements sketch and lineage example are English.

Checks: three-board schema/replay/location validation and append-only prefix
comparison against R2 `918fa46`; R1 migration/link checks; `git diff --check`.
The boards contain 17 cards and 47 events after translation reviews. R1 checks
322 file dispositions, 574 generated artifacts and 1975 local Markdown targets
before adding this evidence file. These checks do not validate external URLs or
claim a new runtime delivery.

L1-M2/M3 remain in progress. No machine-translated draft is accepted without
review; the local draft model is not reliable enough for technical contracts.
