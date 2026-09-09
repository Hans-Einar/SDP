# Proportional profiles

Status: **pilot v0; not canonical**

Profiles scale durable evidence and safety gates. They do not weaken the core
rules: one intent owner, bounded Issue authority, stable identity, truthful
declared/observed/accepted state, proportionate verification, fresh independent
review, and a stop boundary. They also do not waive durable reserved inventory,
reserved-to-prospective materialization, permanent ID no-reuse, or allocation/execution
authority separation, complete authorized-Slice history, terminal aggregate
closure, strict JSON, semantic DAG, or active-write-surface rules.

| Profile | Required pilot surface | Add when |
|---|---|---|
| Minimal | One Issue; one Feature/Refactor/Fix or Study record; compact assignment; Slice contract unless a tiny Fix is itself the reviewed unit; exact candidate; proportionate verification; fresh review; concise handoff. Sprint, Iteration, CurrentIndex, Relations, and Ledger may be absent. | Small/low-risk repository or bounded change with one domain and little concurrency. |
| Standard | Minimal plus foundation/decision links, semantic Relations, current pointers, material transition history, formal verification/review records, early draft PR, and architecture revision gate. | Sustained product development, several work owners, or repeated assignments. |
| Complex/safety-critical | Standard plus independent bounded Studies and convergence/refreeze, explicit claim/evidence levels, hazard/security/physical gates, separate Architect/Verifier when risk warrants, failure/recovery/rollback evidence, stronger provenance and release gates. | Safety, hardware, destructive migration, regulated/security-critical or high-consequence work. |

## Scoped-monorepo is orthogonal

`scoped-domains` is a capability, not a fourth ceremony level. It may be added
to minimal, standard, or complex/safety-critical profiles. When enabled it
requires:

- explicit work-domain registry and immutable domain UIDs;
- scoped new IDs for a repository with multiple active domains;
- qualified cross-domain references;
- explicit shared-domain ownership;
- portable key/ID/path collision validation;
- reservation/refreeze for concurrent assignments; and
- tombstone/successor behavior for domain moves.

A one-domain minimal repository can use unscoped IDs. A multi-domain minimal
repository still obeys all scoped identity rules without inheriting a dense
Ledger or safety taxonomy. Conversely, one complex single-domain repository
may need full safety ceremony without scoped IDs.

Project-local product mechanics—device limits, ABI rules, credentials,
hardware thresholds, or repository-specific tests—remain local. The profile
defines their authority/evidence shape, not their content.

Study-only assignments and explicitly standalone zero-Slice Fix units may be
pathless when the Issue documents that they perform no implementation write
and the compact assignment mirrors a nonblank `pathlessReason`.
An active Feature/Refactor/non-standalone-Fix assignment and its one active
Slice may not use profile proportionality to omit all owned/shared write
surfaces or prohibit their own hosting repository.
