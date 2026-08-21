# Reusable GitHub Issue contract

Status: **pilot v0; not canonical**

The Issue is the complete durable human-readable assignment. An Issue Master
must be able to start from this text, repository instructions, and referenced
evidence without recovering a chat transcript. Later owner/Steering comments
are append-only amendments; they are not silently folded into the original
body.

The template deliberately does not copy mutable facts such as current Issue/PR
state, current branch head, check status, review status, merge commit, labels,
assignees, or timestamps. Those are derived from GitHub at an `observedAt`
instant. Stable URLs, an exact starting baseline, the intended delivery branch,
and the early draft-PR policy are durable authored terms.

Stable GitHub identities use the exact pilot v0 forms defined in
[WorkDomains.md](WorkDomains.md): canonical repository, Issue, PR, and
Issue-comment URLs on `https://github.com` with positive non-zero decimal
numbers and no query, trailing slash, or transport alias. The assignment Issue
and planned PR belong to the registry repository hosting the primary work
domain. Owner/repository comparison is case-insensitive for identity and
collision safety.

Branch names use a conservative portable Git subset: nonblank ASCII letters,
digits, `.`, `_`, `-`, and `/`; no whitespace/control characters, leading
`-`/`.`/slash, trailing slash/dot, empty or dot-prefixed component, `.lock`
component, `..`, `//`, `@{`, `~`, `^`, `:`, `?`, `*`, `[`, or backslash. This
is intentionally narrower than every name Git might technically accept.

## Copyable template

```markdown
# <Feature | Refactor | Fix | Study>: <bounded outcome>

## Authority and work identity

- Primary work reference: `<domainUid> / <recordId>`
- Work type: `<feature | refactor | fix | study>`
- Related work references: `<qualified domainUid + recordId list>`
- Governing requirements/decisions: `<stable links or IDs>`

## Goal and why now

<One observable bounded outcome and its value/urgency.>

## Authoritative starting point

- Repository: `<stable repository URL>`
- Integration branch: `<branch>`
- Exact integration-base commit: `<40-character SHA>`
- Prerequisite Issue/decision/candidate references: `<stable links/IDs>`
- Evidence caveat: `<default truth versus named unmerged evidence>`

## Scope and outcome

- In scope: `<areas and behavior>`
- Expected owned paths: `<portable reservations>`
- Shared touchpoints: `<explicit paths, owner, and convergence gate>`
- Candidate Slice IDs: `<pre-reserved IDs>`
- Observable completion: `<what can be run/seen>`

## Non-goals and prohibited actions

- `<explicit exclusions, repositories, paths, releases, migrations>`

## Invariants and refinement authority

- `<product, architecture, compatibility, security, physical-safety rules>`
- Allowed local requirements/design refinement: `<bounded rule>`
- Must escalate/refreeze when: `<architecture threshold or contradiction>`

## Roles and concurrency

- Issue Master duties and direct-implementation boundary: `<rule>`
- Worker/Architect/Verifier/Reviewer separation: `<proportionate rule>`
- Reservation-set/common-base reference: `<stable record + digest>`
- Dependencies/conflicts/convergence/merge order: `<explicit graph>`

## Delivery contract

- Dedicated branch: `<planned branch>`
- Early draft PR: required, targeting `<integration branch>`
- Commit/publish expectations: `<remote-resolvable exact candidate>`
- Default-branch promotion authority: `<who/what; normally not this Issue>`

## Verification and independent review

- Automated checks: `<commands and expected results>`
- Real/manual evidence: `<required, not applicable, or explicitly pending>`
- Scope/protected-path checks: `<commands/rules>`
- Independent review challenge: `<semantics, risk, exact head>`
- Maximum unresolved severity: `<normally Low or none>`
- Changed-head re-verification/re-review: required after material change

## Steering, stop, and reporting

- Steering may: `<accept / changes required / block / split / cancel>`
- Stop/escalate on: `<missing authority, stale base, collision, scope change>`
- Completion signal: `<exact signal>`
- Final report: `<exact head/PR, validation, review, limitations, handoff>`
- Do not start: `<next Issue or downstream work>`
```

## Assignment mirror

The branch adds one compact
`Steering/Assignments/ISSUE-<number>.yaml|json` record before broad work. It
stores its own portable `source`, the stable authority URL, primary qualified
work reference, complete `authorizedSlices` history with immutable accepted
candidates, `0..1` currently executing `activeSlices`, declared
baseline/delivery targets, reservations, dependency/convergence contract,
required evidence, and stop condition. It links back to the Issue for scope,
non-goals, and prose; it does not duplicate them.

The current mirror and every bound current reservation use their exact pilot
v0 `schemaVersion`, exact `kind`, and `experimental: true`; marker omission is
not a compatibility mode. Its normalized owned/shared write surfaces are also
disjoint from normalized `boundaries.prohibitedPaths`. Active implementation
also cannot prohibit the hosting repository and cannot be pathless; Study-only
and explicitly standalone zero-Slice units may document a pathless boundary
with a nonblank assignment `pathlessReason`.

An amendment comment is recorded as a stable URL under `authority.amendments`.
If it changes scope, baseline, reservation, or merge order, the Master increments
the assignment revision and performs the refreeze protocol before work resumes.
The new revision links a durable history snapshot containing the prior exact
candidate and reservation digest; a changed digest cannot reuse the prior
revision number. Revision 1 omits predecessor/refreeze fields. Revision N
retains the complete `1..N-1` snapshot chain, each snapshot links its immediate
predecessor, and the repository driver reconstructs the exact historical
assignment bytes from each declared Git candidate. Missing repository
resolution, paths, Git objects, or matching Issue/source/revision/reservation
context blocks validation.

Each declared candidate is an exact strict ancestor of the validated Git
`HEAD`, and revision candidates are strictly chronological ancestors of one
another. Existing-but-disconnected roots, sibling commits, descendants,
duplicates, reversals, and the current commit are not assignment history.

The driver also loads and strictly parses the exact historical reservation file
from each historical assignment, recomputes its digest, and compares actual
revision semantics rather than trusting optional snapshot pointers. Canonical
Issue, assignment source, and primary work reference are immutable;
`authorizedSlices` is append-only and every historical non-null accepted
candidate remains unchanged. Exactly one canonical Issue row in that historical
reservation must equal the exact assignment projection for work, authorized/
active Slices, IDs, paths, and edges; historical base, merge order, and
convergence must also agree. Coordinating altered bytes, digests, and snapshot
pointers does not relax this equality. Only the active subset and a duly
refrozen reservation/base may change. Duplicate JSON member names are rejected
before any canonical hash. Compatibility is revision-bound: revision 1 alone
accepts the exact early set-level/string-ID shape, revisions 2–4 alone accept
the exact missing-`authorizedSlices` form, and revision 5+ requires current row
fields. Partial, mixed, hybrid, or late legacy forms cannot satisfy the gate.

Preserved terminal assignments may reference different historical reservation
sets and bases. Only the set group containing a nonterminal assignment is the
current activation epoch; current collision/DAG/order/convergence checks do not
misclassify accepted historical paths as live writers. More than one such
current group is incoherent and blocks validation.
Epoch identity is always the pair `(reservation-set ID, canonical digest)`, so
terminal epochs may retain one stable set ID with different immutable digests.
Every pair resolves exactly one object; an exact duplicate pair is ambiguous.
Every typed set is structurally and semantically valid even before an
assignment references it. Such an unbound set is a preparatory epoch, not a
detached syntax object: its Issues, domains, primary records, inventory
authority, IDs, authorized reservations, paths, and shared owners cohere with
the hosting repository, and its claims are collision-checked against all
nonterminal bound epochs and other preparatory epochs.

`boundaries.prohibitedRepositories` uses canonical GitHub repository URLs;
every prohibited/owned/shared/local-source path uses the common portable path
contract; and each shared touchpoint has a nonblank `allowedMutation`.
