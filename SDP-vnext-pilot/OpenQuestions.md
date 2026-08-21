# Explicit pilot questions

Status: unresolved by design; none prevents the first pilot

The pilot contract makes enough choices for deterministic use, while deferring
schema details that need evidence:

1. Should every multi-domain repository always require scoped IDs, as pilot v0
   does for new work, or can a declared default domain continue minting
   unscoped IDs after additional domains are added?
2. Is a UUID URN the best permanent domain UID representation, or should a
   future registry use another opaque identifier while preserving the same
   non-derived semantics?
3. Should the domain registry become a capability section in the project
   manifest, a dedicated canonical record, or a manifest pointer to one or
   several registries?
4. Should canonical schemas store the full issued-ID inventory, an append-only
   allocation ledger, a high-water mark plus exceptions, or a digest-backed
   generated projection? Pilot v0 uses structured prospective/legacy-preserved
   issued members plus per-Issue qualified reservations and enforces no-reuse.
5. What generated `CurrentAssignments` snapshot/envelope is most useful without
   turning observed GitHub state into committed authored truth?
6. Which material transitions deserve the future general Ledger vocabulary,
   and how should commit-time versus event-time validity be represented?
7. Does a rare destination key collision need a read-only mount alias
   mechanism, or should moves remain blocked until the repository topology
   avoids it? Any answer must preserve canonical IDs and absolute references.
8. After how many successful profiles/repository moves should optional fields,
   status vocabulary, and serialization be frozen into canonical schemas?

Pilot reports MUST record friction or evidence for these questions rather than
silently inventing local answers.
