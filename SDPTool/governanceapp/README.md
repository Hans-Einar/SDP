# Governance application operations

Operations() is the shared inventory for existing local governance CLI operations
and Session actions. Each entry binds a typed domain handler, strict decoder,
input kind and optional public action metadata. Both adapters call Execute directly;
this package must not import root sdptool, CLI, presentation, host or MCP adapters.

Write marks local caller-attribution eligibility; it never grants domain authority.
Governance methods retain their own caller, revision, stage, claim and recovery
checks. Session MCP applies its separate launch-bound capability rules directly to
the same domain service; it is not an unrestricted local CLI adapter.

Preserve externally published action metadata and CLI grammar. New functionality
requires its own selected Activity. To change an existing operation, update its
single definition and the corresponding compatibility/behavior tests. Never route
an action through CLI argv or decode presentation output as an internal result.
The root architecture tests enforce this direction and use violating specimens.

The public Invoke response carries heterogeneous results. Session results now
retain domain Go types; their JSON wire contracts remain unchanged. Consumers of
Invoke should use the result's documented wire schema rather than rely on the
previous incidental map produced by internal JSON reparsing.
