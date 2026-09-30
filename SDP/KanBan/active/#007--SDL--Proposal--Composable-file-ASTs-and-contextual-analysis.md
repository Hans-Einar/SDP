# Composable file ASTs and context-dependent semantic analysis

| Field | Value |
| --- | --- |
| id | KB-SDL-007 |
| project | SDL |
| type | Proposal |
| CardState | in-progress |
| Systems | SDL |
| created | 2026-09-30T10:36:36.093682+00:00 |
| source | Owner AST-forest discussion and capture request, 2026-09-30 |
| next_review | During Session 0001 S2 before selecting frontend implementation |
| tags | file-ast, source-graph, semantic-model, partial-analysis, incremental-analysis |

## Intent and Session

The owner wants independently parsed subtrees to be reusable when an enclosing
source/System is loaded later: a forest of small trees can join a larger model.
This supports component-first inspection rather than requiring the root to be
loaded first. Record the distinction between AST and decorated AST (DAST).

[Session 0001](../../Sessions/session-%230001--SDL_expansion.md) records the
conversation. [KB-SDL-005](%23005--SDL--Change--System-and-source-sets.md)
remains the source-composition delivery owner; this card captures a refinement,
not a competing parser or automatic split of KB005. Consider consolidation when
S2 selects its plan. [SSD2](../../04--Design/SDL/SourceComposition/Contract.md)
already proposes parse-once files, but starts discovery from a known entry.
This proposal adds late enclosing context and explicit partial-analysis behavior.

## Proposed model

- File AST: immutable representation of one file's written syntax and source spans.
- Source graph: includes/path dependencies, potentially shared and cyclic; no
  recursive copying or physical reparenting of a shared file AST.
- Semantic System model: resolved identities, typed membership/relations and
  validation in one explicitly selected System context.
- Decorations: symbol bindings, inferred/checked kinds, relation classification,
  context and diagnostic state. May live in side tables rather than mutable ASTs.
  “DAST” describes semantic enrichment, not a required separate file format.

An unchanged file AST can be reused under a later root. This does not justify
reusing semantic results across a changed namespace, dependency revision or
System: re-resolve/revalidate affected bindings. Cache syntax by source content,
parser/profile identity and required parse options; key semantic analysis by its
context and dependency revision as well. A physical file is not its global symbol
identity, and equal bytes in distinct files do not merge their declarations.

Keep partial states explicit: syntactically parsed, dependencies pending, partially
resolved and validated-in-context are different outcomes. An unresolved external
context is not a successful System build. Opening a fragment must not fabricate
System ownership. Root-relative paths from SSD2 require an explicit base context
before loading dependencies; parsing alone does not. Evaluate how the host supplies
that context without another authored manifest or implicit parent search.

## Required study/design questions and acceptance

1. Parse MachineService first; load MVP1 later; reuse its unchanged AST and bind
   membership. Compare semantics with root-first loading: same result.
2. A shared contract has two incoming references but one file AST. Source-cycle
   traversal terminates while illegal semantic containment still fails.
3. Missing enclosing context permits clearly labeled local inspection; complete
   check/generation rejects unresolved obligations and never presents partial data
   as a validated model. Define which local diagnostics are conclusive.
4. Later context introduces a duplicate name or changes a referenced declaration:
   old semantic decorations are invalidated; diagnostics retain original spans.
5. Reuse under two System contexts cannot leak bindings/identity from one into the
   other. Removing a source/root cannot leave stale membership edges.
6. Define the bounded first increment versus later dependency-aware recomputation.
   File-AST caching alone is not evidence of a complete incremental compiler,
   runtime hot reload or actual native XFMD fragment preview.

## References and authority

The owner endorsed documenting this direction. API/schema details and caching
strategy remain design work; registration does not authorize product implementation.
The previous explanatory turn cited
[Cornell semantic-analysis notes](https://www.cs.cornell.edu/courses/cs4120/2022sp/notes/semantic/)
for decorated AST terminology. That reference does not determine SDL's design.
No new executable language semantics or cross-System import capability is adopted.

## Worklog

Captured with the Session T008 discussion and T009 request. EVT-KB-SDL-000038.

## Execution

Owner selected PLAN-SDP-0012 for Session S2–S5; syntax reuse and contextual
semantic rebuilding are selected, not a fully incremental runtime compiler.
