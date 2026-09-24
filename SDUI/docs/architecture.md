# SDUI — implementert Go-arkitektur

Oppdatert 2026-09-22. Én aktiv frontend for SDUI 0.2. Kjørbare moduler og
kommandoer står i [Go-området](../go/README.md); felles design er
[beskrevet i SDL](../design/README.md).

| Pakke | Ansvar |
| --- | --- |
| go/parser | Lexer, recursive descent, AST/spans, lokale regler og normalisering |
| go/layout | Én målt geometri; relative dimensjoner, rader, wrap, ratio og klipp |
| go/markdown | Avgrenset Goldmark-innhold; måling og registrert Mermaid-provider |
| go/svg | Statisk eksport fra samme geometri, Go Regular-glypher og native kontrollutseende |
| go/presentation | Strukturell konsoll-/Markdown-dump og avgrenset kontrollgalleri |
| go/runtime | Sesjon, handles, accepted/draft, typed events og atomiske propertybatches |
| go/reload | Validerte kandidater og bevaring av kompatibel state/identitet |
| go/host/fynehost | Native input/button, fokus/tastatur og UI-trådpublisering |
| go/codegen | Uavhengige typede Go-konstruktører for Document/Root |
| go/cmd | CLI og native komposisjon; filtilgang skjer her, ikke i parseren |

Parse → Normalize/Compile produserer Document og ekspanderte Instance-trær.
AST-JSON er tagget sdui-ast/0.2. Span er halvlukkede UTF-8-byteområder med
ettbaserte Unicode-linje-/kolonneposisjoner. Node bevarer grupper, rader,
regioner og suffix-formatering. Runtime kopierer inputmodeller; Go-structs er
ikke språklig immutable. Ikke muter en modell som er i bruk av en vert.

Eksplisitte widgetnavn gir offentlige instansbaner; anonyme segmenter har
syntetiske navn. Session/Path/Generation/Kind identifiserer handles.
Reload bevarer kompatible navngitte instanser, accepted/draft og fokus;
sletting/typebytte invaliderer tidligere handles. Kilde- og bindingsfeil beholder
siste gyldige modell. Kildewatchere leverer kandidater gjennom fyne.Do.

Layout får vertens tilgjengelige område; kilde har ingen pikselbredde/-høyde.
Barn bruker sin nærmeste kildeancestor, og `{16:9,<->}` avleder høyde fra fylt
bredde. Fonten forblir logiske DIP ved resize. SVG og Fyne deler rektangler,
tekstmål og klipp, men native kontroller har vertens rasterisering/tema.

SDL-adapteren ligger i SDL/go/bridge. Parserens ref/callback/
setHandle er data; komposisjonen registrerer SDL-modul, Go-funksjon og typed
bridge.Plan. Action-core 0.1 gir eksplisitt avgrenset kjøring. En allerede
akseptert Go-domenetransaksjon kan ikke rulles tilbake hvis senere UI-publisering
feiler; det gjøres ingen automatisk replay. Se [runtimekontrakten](runtime-contract.md).

Grenser og reelle prøver: [G1](../go/evidence/G1.md), [G2](../go/evidence/G2.md),
[G3](../go/evidence/G3.md), [G4](../../SDL/go/evidence/G4.md),
[G5](../../SDL/go/evidence/G5.md). Ingen alternativ parser i
SVG/Fyne/XFMD, ingen Rust-uttrekkscrate eller obligatorisk C-ABI.
