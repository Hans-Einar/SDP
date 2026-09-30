# System and source-set contract — SSD1

Status: **recommended implementation design**, delivered by PLAN-SDP-0010.
This is not the released language specification or evidence that new syntax works.
The current implementation remains design-core/0.5. [Acceptance](Acceptance.md)
defines the tests required before adopting the proposed profiles.

## 1. Decision and bounded outcome

A System may have several authored files without changing declaration identity,
validation or navigation merely because a declaration moves between those files.
The first supported unit of compilation is **one System, one explicit source set,
one scope, one immutable validated result**. A directory is not an import.

Recommend design-core/0.6 for System syntax and sdl-source-set/0.1 for an input
manifest. These version labels are proposed profile versions, not release numbers.
Keep the existing design-core/0.5 single-file path and its frozen fixtures. This
is SDL compatibility; it does not revive SDUI 0.1 compatibility rejected by the owner.
Do not accept multi-file 0.5 by concatenation or infer a System from a folder name.

| Alternative | Assessment |
| --- | --- |
| Recursive filesystem discovery | Hidden files and checkout changes alter meaning; reject as authoritative membership. Optional inventory reporting can be added separately. |
| Textual include/preprocessor | Synthetic positions and order-dependent declarations make diagnostics and downstream provenance fragile; reject. |
| Explicit manifest + file ASTs + whole-System validation | Selected: reusable I/O-free frontend, deterministic membership, exact input revision, modest syntax extension. |
| Full packages, public exports and cross-System imports immediately | Useful destination, but needs contract visibility and dependency semantics beyond this delivery. Defer explicitly; do not simulate it with global names. |

The source-tree study's package/private/export recommendations remain future
scope. Here all declarations in a System share a scope. Neither file placement
nor a Containers directory creates private visibility. This bounded choice must
be stated in the new language profile; it is not proof of package isolation.

## 2. Input and language syntax

An explicit path selects either a standalone `.design` or a `.design-set.json`.
No directory guessing, nearest-manifest lookup or automatic profile upgrade.
The proposed manifest is strict UTF-8 JSON:

```json
{
  "profile": "sdl-source-set/0.1",
  "entry": "System.design",
  "sources": ["Features.design", "Containers/Frontend.design"]
}
```

The entry is included once automatically; sources lists only the remaining files.
Unknown/duplicate JSON keys, trailing values, absent fields, non-string members
and duplicate paths are errors. An empty sources array is legal. Entry and source
order carry no semantic meaning. Entry denotes where the System is declared, not
an execution starting point. The manifest does not repeat the System name.

Each file starts with `language design-core version 0.6.`. Each follows the
existing declaration-before-statement syntax. The entry contains exactly one
`system Name.` declaration. Other files contain none. A standalone 0.6 file has
exactly one System itself. Missing, multiple, misplaced or mixed-version Systems
fail before projections. A fragment can parse independently but cannot pass a
whole-System check by itself. Error output must explain that distinction.

The existing lower-case keywords and `[A-Z][A-Za-z0-9]*` identifiers remain.
Add exactly the `system` declaration kind and relation signature
`system contains container`. The existing unit/container relations remain.
Every declared container must have one direct System-membership relation, for
example `Frontend contains CommandProcess.`. This is namespace membership, not
runtime ownership or allocation. It is checked separately from unit containment:
a System membership edge cannot create a second runtime parent. Systems are not
unit subtypes: no allocating functionality to a System, no System owning a Unit,
no nested Systems, no System channel participation. A System with no containers
is valid for an early requirements model.

All other declarations belong to the declared System scope implicitly; membership
is distinct from `owns`, unit `contains`, and mode-dependent `allocated-to`.
No fabricated runtime container is needed for a library or shared contract.
System membership is displayed as such, never as a new runtime communication edge.

## 3. Scope, identity and public boundaries

Collect declarations from all members before semantic validation. A reference can
resolve forward or across files in the same set. Duplicate names remain errors,
including collisions with the System name; duplicate facts also remain errors.
Use existing relationship signatures, ownership, channel/data completeness and
cycle checks over the complete graph, not independently per file.

Logical identity is the pair `(System name, local name)`, with kind as checked
metadata. Moving an object between member files preserves identity; renaming the
System or local name is an identity change. The same local name in two independent
Systems is allowed. System names are only required to be unique within a consumer's
explicitly registered model collection; no global registry is introduced. Registration
model IDs continue to be routes, not implicit language namespace aliases.

Use structured identity fields in AST/model outputs. Do not introduce dotted names
into source just to print identity. Every occurrence retains its authored file
and span. Stable relation identity is the System identity plus the typed relation
and its operands/qualifiers, excluding file location; duplicate identical facts
are already rejected. Existing ordinal fact IDs may be retained for 0.5 output.
For 0.6, use deterministic IDs from canonical typed tuples, with collision checks;
scenario ordinals remain semantic numbers and must not become lexicographic order.

**Cross-System contract for this increment:** no imports, external declarations,
public/export modifier, qualified references or dependency manifests. An unresolved
name is an error, not an automatically external name. Do not merge independent
entries or an ecosystem Landscape into one scope. A local interface with an
external-owner README link remains a local boundary abstraction, not proof of
linking to that provider. Cross-System validation and public/private visibility
are explicitly unsupported capabilities, including in the migration outcome.

A later increment must define explicit dependency identities, public contract
closure, provider revisions, compatibility and source-versus-export precedence.
It must reject private access, stale exports and duplicate source/export identities.
The historical study describes that destination; this design does not silently
adopt it or require downloading other repositories.

## 4. Loading, limits and source positions

The loader owns disk access and captures bytes; the parser owns syntax; a pure
compile/check operation links and validates supplied file ASTs. No parser opens
files, resolves SDL/SDUI runtime modules, fetches URLs or executes code.

Manifest paths are normalized relative POSIX paths inside the manifest directory.
Reject absolute paths, backslashes, empty/dot/dot-dot segments, control characters,
colons, non-.design members and duplicate normalized names. A member must be a
regular file. Reject symlinked members or symlinked components beneath the input
root for this initial profile; resolve the root once. Detect hard-link aliases via
file identity, rather than accepting two paths to one input. No manifest nesting,
include graph or dependency recursion exists, so source inclusion cycles are
impossible by construction. Missing members fail; unlisted neighboring files are
ignored and never hashed or parsed. No directory crawl is part of a normal check.

Initial explicit limits: manifest 64 KiB, at most 128 files including entry,
2 MiB per source and 2 MiB total source bytes, at most 250,000 tokens across the
set, path length at most 1,024 UTF-8 bytes. Limits are enforced while reading and
tokenizing, not after allocating an arbitrarily large aggregate. Retain existing
semantic limits and bounded diagnostic collection. Reject the whole candidate on
any violated limit. These are conservative design choices to verify in S3/S4.

A SourceFile has root-relative path, exact bytes and SHA-256. SourceSpan combines
that path with the existing half-open byte offsets and one-based line/byte-column
values. CRLF/UTF-8 do not justify fabricated character columns. The manifest has
its own source identity and location for membership diagnostics. Never shift spans
into concatenated text. There is no meaningful single contiguous span over a set.

Diagnostics contain code, message, primary SourceSpan and optional related spans.
Duplicate declarations identify both occurrences; cross-file ownership/cycle errors
identify participating facts. A missing name points to its use. Missing files
point to the manifest member; file-read errors identify the requested path. Stable
ordering is path, byte offset, diagnostic code, then message, not map iteration.
No partial validated result is published on errors; parse-only AST output may
identify successfully parsed files but is explicitly not a validated model.

## 5. Go integration boundary

Proposed API roles (signatures to settle during S2, not installed symbols):

| Owner | Input → output | Responsibility |
| --- | --- | --- |
| sourceinput loader | explicit path → captured InputSet | Strict manifest, filesystem/size checks, exact bytes, paths and digest; no canonical synthetic source. |
| parser frontend | named bytes → FileAST; InputSet → checked SystemModel | Shared lexer/parser with version dispatch; declaration collection, resolution, validation and source-aware diagnostics. Remains I/O-free. |
| canonical formatter | checked model + original membership → formatted per-file outputs | Validate globally, format locally; preserve authored partition. |
| viewpoint projector | checked result → Views | No reload/reparse; preserve System identity, revision and provenance. |
| documents/broker | Views + query → revision-bound bundle | Reuse projector and existing lease/publication semantics; recheck all inputs before current delivery. |
| SDPTool facade | registered input or explicit preview path → same frontend/projector | No second parser or parallel manifest interpretation. |

A checked result carries language/input profile, System identity, immutable member
ASTs, symbols/relations, all source hashes, candidate revision and diagnostics.
Only a successful check can construct the validated representation used by views.
Immutability must be enforced by API ownership/copies, not only a comment on a
public Go struct. Consumers may not substitute a canonical string and reparse it.
Keep existing single-source public entrypoints as adapters with unchanged 0.5
results; add an explicit checked-model entrypoint for the new path.

The untracked `sourceinput` draft is not this implementation: it currently flattens
ASTs, shifts spans, returns canonical Source and keys origins by sentence text.
It also lacks System syntax and the strict path/aggregate/profile contract above.
Its bounded reads/output protection are ideas to reuse after review, not a reason
to adopt uncommitted files or preserve the synthetic-source interface.

## 6. Canonical form and revisions

Check complete semantics first. For 0.6 canonical form is defined **per file**:
header, declarations sorted by identifier, statements sorted by existing canonical
sentence order, LF and one final newline. Validation crosses files; canonicalization
does not move declarations/facts between them. Order-dependent runtime behavior is
not introduced; scenario rendering continues to order by numeric step ordinal.
Manifest canonical form uses fixed key order profile/entry/sources and a sorted
sources array. Whitespace alone can require formatting without changing meaning.

`sdl format set.design-set.json` should emit a versioned JSON file map containing
the manifest and all formatted sources on stdout. It does not concatenate sources
or overwrite several authored files implicitly. Existing standalone format output
remains plain source. Unknown output schemas/profile combinations fail explicitly.
AST set output includes per-file ASTs and resolved identities/provenance, not a
fictitious combined AST with one file span. Version the new set output envelope.

Candidate revision is SHA-256 over a domain-separated, length-prefixed sequence:
`SDL-source-revision/1`, input profile, language profile, entry path, manifest-byte
hash (empty for standalone), then sorted `(relative path, exact byte hash)` pairs.
Fields use UTF-8 with unsigned 64-bit big-endian length prefixes. Hashes are
lower-case hex ASCII. Standalone 0.6 uses logical path `@source`; display path is
separate provenance. 0.5 retains the existing hash-of-source-bytes convention.
Changing whitespace, any member, membership or manifest bytes invalidates revision.
Reordering manifest members changes revision but not semantic identity/facts.
Absolute checkout paths do not enter the revision. A file move preserves logical
identity but changes source provenance and candidate revision.

Cache keys additionally include frontend/profile implementation identity, projector,
query, renderer and output schema versions. Revision identifies input, not tool
correctness. Reuse one captured result for AST, tree and generation; never combine
facts from independently captured revisions.

Before publishing a current selection, re-read/re-hash the manifest and every
member with the same loader policy. Missing, changed or newly invalid input cancels
the candidate; old leased output can remain readable as its old revision. Do this
on cache hits too. Coalesced file watches are an optimization only. No busy-loop
retry: report source-changed and let the caller request again. A prepublication
check is an optimistic freshness boundary, not an atomic snapshot of a mutable
filesystem. Describe outputs as the captured revision checked at that boundary;
stronger cross-file editing transactions require an immutable snapshot facility.

## 7. Consumers and migration

| Surface | Required first delivery |
| --- | --- |
| sdl check/ast/format | Dispatch explicit source path; original diagnostics; per-file canonical checks; unchanged standalone 0.5 contract. stdin supports a standalone source only, not relative manifests. |
| sdl view/viewpoints | Use checked result; identical semantic facts and source-aware mappings; System membership visible without inventing runtime edges. |
| broker | Load full set, use aggregate revision in cache, recheck full set before delivery including cache hit; preserve cancellation, sequence and leases. |
| SDPTool registration/tree/select/preview | Add design-core/0.6 profile and explicit .design-set.json source support; registered System must equal declared System. Existing schema fields suffice; no automatic profile fallback. |
| standalone .design preview | Treat explicit .design as standalone; never silently load neighbors. A partial file produces a useful missing-System/unresolved-name diagnostic and points users to the registered set workflow, without guessing a manifest. |
| document provenance | Version source-map/bundle changes; record member paths/hashes and original spans for facts and declarations, including System memberships. Selection target revision must match the tree candidate. |
| snapshot / design attachment | Route design input through the same checked frontend if source sets are enabled there. Until migrated, advertise standalone only and reject set input before output rather than mislabeling support. |
| action/class tools and SDUI bridge | Keep current standalone profiles; source-set structural support does not turn design-core into an executable action language. |

SDPTool's registration `profile` stays the language profile; the source extension
selects the input profile. Its configured model ID retains routing identity.
Capabilities must distinguish single-file/set support; hardcoded 0.5 preview
provenance must become derived from the checked result. Fail older consumers clearly
rather than feeding them a flattened 0.5 model with discarded System information.
External XFMD application edits remain in its own backlog. Producer CLI/library
acceptance is mandatory here; actual native XFMD acceptance requires that workstream.

The first candidate fixture is a **copy** of the small Frontend ecosystem model,
split into System entry, Features and Containers files. Keep all original 0.5
facts and add only explicit System membership. Compare normalized semantic tuples,
not ordinal fact IDs or raw rendered diagram bytes. Keep its authoritative source
and navigation registration unchanged until the implementation proves the path.
KB-SDP-020 then decides migration of current sources and links. Do not migrate the
large SDUI model, all 17 ecosystem models or the 68-file MVP1 experiment here.
MVP1 uses experimental language constructs this extension does not implement.

## 8. Design disposition

Recommended for S2: plan the closed-System increment above, with original-source
provenance and producer/consumer revision tests in the same delivery. No new
mandatory owner gate is invented. Implementation authorization remains a separate
Session step. Public exports/cross-System linking remain explicit future scope of
KB-SDL-005; the full card cannot close while those obligations are silently omitted.
S2 must either retain that later milestone or record an explicit successor/split.
No external repository or installed process-template migration is selected.
