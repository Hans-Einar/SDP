# SDL — avgrenset data- og wireprofil

Data-/wireprofilen ble innført i V2 (0.3) og inngår nå i **design-core 0.5**
i Go. Reglene nedenfor er aktive; den bredere datastudien er en kandidat.
Strukturfakta skrives separat. Dette implementerer ikke binær serialisering.
[Go-innganger](../../go/README.md) og [G4-bevis](../../go/evidence/G4.md).

Nye deklarasjonstyper: `dataset`, `database`, `datagram`, `contract`, `variant`,
`field`, `encoding`. Alle bruker `kind Name.` og én felles navnetabell.

| Utsagn | Signatur og regel |
| --- | --- |
| `D upholds C.` | Dataset/Datagram → Contract; nøyaktig én kontrakt per datadefinisjon |
| `M from D.` | Datagram → Dataset; nøyaktig én logisk kilde |
| `U owns B.` | Unit → Database; én umiddelbar ansvarlig eier, som for Functionality |
| `U holds D.` | Unit/Database → Dataset; én umiddelbar holder i denne definisjonsprofilen |
| `C defines V.` | Contract → Variant; én kontrakteier per variant |
| `C has-field F.` | Contract/Variant → Field; én umiddelbar felteier |
| `P projects D into M.` | Functionality × Dataset × Datagram; må samsvare med M sin kilde |
| `E encodes V.` | Encoding → Variant; nøyaktig én variant per encoding |
| `E places F at 0 bits 16.` | Encoding × Field × offset × bitbredde; eksplisitte, ikke-negative heltall, positiv bredde |

`Contract has completeness = open/closed.` er obligatorisk. Open kan beskrive
ufullstendig arbeid, men gir ingen packet. Closed Dataset-kontrakt har minst ett
felt og ingen varianter. Closed Datagram-kontrakt har minst én variant; felles
felt på kontrakten er headerfelt, og hver variant har minst ett payloadfelt.
Samme kontrakt kan gjenbrukes av flere Datasets/familier når formkravene passer.

Hvert Field må angi `value-type = unsigned/signed/boolean/text/bytes/decimal`
og `presence = required/optional`. Fravær er ikke null, tom tekst, false eller
nulltall. Dette er formkontrakter; enheter, verdidomener, felttransformasjoner,
versjonsforhandling og utføring av projeksjoner er ennå ikke implementert.

Database betyr stedet persistent data kan hentes ved behov. SQL, prosessgrense,
transaksjon og bestemt lagringsmotor utledes ikke. Dataset kan være transient
hos en Unit; persistens krever en eksplisitt Database-holder. Definisjonene er
ikke kjørbare instanser eller en flerinstans/deploymentmodell.

Encoding krever closed kontrakt og eksplisitt `byte-order = big-endian/little-endian`
og `bit-order = most-significant-first/least-significant-first`. Alle header- og
variantfelt må plasseres nøyaktig én gang, uten hull eller overlapp, fra offset 0.
Bitposisjonene er fortløpende kontraktposisjoner; etikettene er ikke encoderkode.
Bool krever 1 bit; signed/unsigned høyst 64; bytes et multiplum av 8. Text,
decimal og optional felt avvises i denne faste packet-profilen. Inntil 65536 bits
støttes per layout. Utvidelser må definere padding, variable felt og presence-
encoding eksplisitt; generatoren gjetter ingen av disse.

EBNF-tillegg til kjerneprofilen:

```text
projection = identifier, "projects", identifier, "into", identifier, "." ;
placement = identifier, "places", identifier, "at", integer, "bits", integer, "." ;
integer = "0" | nonzeroDigit, { digit } ;
```

AST: `Projection(subject, dataset, datagram, span)` og
`Placement(subject, field, offset, width, span)`; tall har egen `Integer(value, span)`.
Andre utsagn bruker eksisterende Relation/PropertyAssignment. VP09 gir data-
og kontraktkart, VP10 gir packet for en validert Encoding. Kildekartet tar med
også projeksjonens tredje argument og hver eksplisitte feltplassering.

## Historisk V2-leveranse — 2026-09-22

Tall og status nedenfor gjelder denne milepælen før Go-porten. De er ikke
nykjørte tester eller dagens samlede implementasjonsstatus.

Milepæler: V2-M1 språk/validering; V2-M2 generiske VP09/VP10 og tester;
V2-M3 portert eksempel, regenerering og verifikasjon før fasepush.

V2-M1 verifisert: 49 parser-/validator-/CLI-tester består. Senere milepæler
leverer viewpoints og regenererte fellesartefakter. 0.3 erstatter aktiv 0.2.

V2-M2: VP09/VP10 og kildekart er implementert; 17 verktøytester består.
Packet bruker eksplisitte bitområder etter [Mermaid packet](https://mermaid.js.org/syntax/packet.html),
verifiseres også mot lokal mmdr, og er ikke en binær serialiseringsimplementasjon.

V2-M3 levert: 227 deklarasjoner, 509 fakta og 73 SVG-diagrammer. 49 SDL-,
17 viewpoint- og 36 SDUI-tester består; 15 delvise allokeringshull bevares.
Packet- og dataeksemplet er generert av verktøyet fra den felles SDL-kilden.
