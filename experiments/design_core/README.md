# design-core parser experiment

Status: isolated prototype for **design-core 0.5**, not an installed SDP gate.
The [language definition](../../docs/Design-Language-Definition.md) is authoritative.
The implementation uses the Python standard library; tested with Python 3.12.

## Purpose and boundary

Make the existing structural grammar executable before adding behavioral syntax.
The recursive-descent parser is handwritten from the EBNF; it is not generated
from the Markdown document. Definition, implementation and tests must be changed
together when the language changes. The three complete models in the definition
are loaded directly by the tests to detect drift in those examples.

The pipeline is:

```text
UTF-8 source -> tokens -> AST -> symbol table / semantic checks -> canonical check
```

These stages answer different questions. A reversed ownership sentence has a
valid syntactic structure and therefore produces an AST. Its subject and object
types fail semantic validation. A structurally valid model with alternate spacing
fails the canonical check. No stage proves that the modeled architecture is good
or that its implementation satisfies the claims.

## Run from the repository root

```sh
python3 experiments/design_core/design_core.py check experiments/design_core/examples/ui-ownership.design
python3 experiments/design_core/design_core.py ast experiments/design_core/examples/ui-ownership.design
python3 experiments/design_core/design_core.py format experiments/design_core/examples/ui-ownership.design
python3 -m unittest discover -s experiments/design_core -v
```

Each command also accepts `-` instead of a filename to read UTF-8 from stdin.

| Command | Result |
|---|---|
| `check` | JSON validity flag and diagnostics, including canonical form. |
| `ast` | The same checks, plus the source-ordered AST and declaration symbol table. An AST is still returned on semantic errors; it is null on syntax errors. |
| `format` | Explicit canonical text on stdout if syntax and semantics pass. Semantic errors go to stderr, with no model output. No file is overwritten. |

Exit codes: `0` successful, `1` invalid model, `2` input/usage error. The JSON
shape is an experimental inspection interface, not a frozen interchange contract.
The formatter can change layout/order but cannot reverse ownership, remove
duplicates, invent missing declarations or otherwise repair design meaning.

## AST and Python API

Immutable named tuples distinguish `Model`, `Header`, `Declaration`, `Identifier`,
`Relation`, `Dependency`, `Allocation` and `PropertyAssignment`. Every node carries a `Span`.
Offsets count decoded Unicode characters, starting at zero; line and column
numbers start at one. Span ends are exclusive. Tabs count as one character.
AST collections retain source order; formatting sorts only the unordered facts.

For example, the ownership statement yields this abbreviated AST:

```text
Relation(
  subject=Identifier(name="PresentationManager", span=...),
  verb="owns",
  object=Identifier(name="ValidateBindings", span=...),
  span=...
)
```

Model types come from declarations, not identifier spelling. `Container` is
compatible with `Unit`; no other implicit conversion exists. The AST records
what was written. The symbol table supplies resolved types without rewriting it.

| API | Responsibility |
|---|---|
| `parse(text)` | Return AST or raise `ParseError` containing a positioned diagnostic. Stops at the first syntax error. |
| `symbol_table(model)` | Index declarations; the first declaration remains available if a duplicate exists. Always validate before treating resolution as valid. |
| `validate(model)` | Return semantic diagnostics for a parser-produced AST, independent of layout. |
| `check(text)` | Return `(model_or_none, diagnostics)` including canonical layout after semantic success. |
| `canonicalize(model)` | Return canonical source, or raise `ValidationError`; never silently repair semantics. |
| `to_json(value)` | Convert nodes and diagnostics to JSON-compatible inspection data. |

Reversed ownership reports both argument errors. It also reports a missing valid
owner if no other valid ownership statement exists. Canonical errors are deferred
until semantic errors are resolved, avoiding formatting advice on an invalid model.

## Evidence and limits

The tests cover documented examples, AST/source locations, all core relations and
properties, Container compatibility, reversed argument types, unresolved/duplicate
names, mandatory dependency modes, ownership/containment cardinality, cycles,
duplicate/conflicting facts, canonical round trips, and CLI results. A deep
containment chain and an acyclic refinement diamond exercise graph validation.

Round-trip checks cover structural facts only. They do not establish equivalence
of differently named designs or behavior preservation. There is no scenario
execution, guard/state semantics, diagram translation, source-code binding,
implementation verification, code generation or hardware interaction.

The next language increment should come from one concrete design scenario and
its positive/negative cases. Define its missing semantics in the language
definition first, then extend these AST nodes and checks. Do not add permissive
fallbacks for unsupported sentences merely to make an example parse.

## V1 — mål, bidrag og allokering

0.2 erstatter aktiv 0.1; lokale strukturelle eksempler er portert. Nye Actor,
UseCase og Feature bruker pursues/supports/contributes-to. Functionality kan
bidra direkte til et bruksmål, og bidragene er mange-til-mange. allocated-to
krever Container og Mode; to ulike Containers i samme modus er en feil.
Det logiske eierskapet beholdes. Manglende allokering er ikke automatisk feil,
og ingen deployment utledes fra contains. Se språkdefinisjonens 5.1 og 10.3.
[SDL-verktøyet](../../SystemDesignLanguage/tools/README.md) genererer VP01/VP07
fra de validerte faktaene. Go-port og runtime er fortsatt framtidig arbeid.

## Gjeldende samlet profil

V2–V4 er nå implementert: data-/wirekontrakter i data_core.py, Channel/scenarioer
i channel_core.py og typede leveranseplanrelasjoner i design_core.py. Aktive
0.1–0.4-innganger er erstattet og lokale prøver portert. V1-avsnittet over viser
leveransehistorikken. Registrerte tillegg står i språkdefinisjonens profillenker.
Ingen Go-kode eller modell-runtime er implementert av disse strukturelle kontrollene.
