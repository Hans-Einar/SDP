# M0 operation and consumer inventory

Baseline code: d0bf6dc37027bbb75f5f96fecc169e681a3a6c9b. Exact Go file digests and
all 258 exported function/method declarations are retained in ../06-Verification/M0.
[Operations.json](Operations.json) enumerates facade/model/action operations and
current/target owners; action catalogue has exactly 28 entries and is frozen as
wire-compatibility evidence. No Activity actions are added by this migration.

## Complete family disposition

Root CLI has 14 routing tokens: actions, invoke, governance, run, model, release-log,
install, upgrade, preview, discover, view, tree, select, sdui-preview, plus help and
version options. Model syntax/variants are enumerated in Operations.json. All
22 Session actions also have governance CLI paths; six non-Session governance
operations remain CLI-only. The current Session action argv/output bridge and
central attribution allowlist are migrated. Other services already use typed Go
calls and retain ownership. Model, install, Blueprint and PG stores remain separate
bounded state owners, with projecthistory owning shared management persistence.

Controller, MCP and external-consumer inventory is in [Boundaries.md](Boundaries.md).
It includes 51 controller command paths, nested stdin controls and six MCP launch
configurations (28 unique tools). Controller direct service calls stay; MCP's
repeated transport framing/retry code is consolidated with distinct error policies.
MCP launch/instance/skill/owner restrictions remain adapter-specific requirements,
not replaced by local action caller policy.

## Preserved external behavior and defects

Keep CLI grammar/flags, exit codes, error/result envelopes, catalogue IDs/schema/
metadata revision, model heads/recovery, ledger locking/idempotency, PG leases and
read-only guarantees. Preserve root-relative versus cwd-relative file resolution.
Source links, deferred navigation and all current viewpoint content stay unchanged.
Existing legacy navigation limitations and manual Activity reader limitations are
known future-work constraints, not reclassified as delivered fixes.

Existing native tests exercise source/overview/deferred navigation via the actual
producer. Verify them on the candidate when locally available; source inspection
or fixture-only tests do not prove native acceptance. External repositories remain
read-only. All Go subpackages are either explicit adapters, domain/state owners,
external language libraries or presentation/bootstrap utilities; no unclassified
package is hidden behind a generic rewrite category.
