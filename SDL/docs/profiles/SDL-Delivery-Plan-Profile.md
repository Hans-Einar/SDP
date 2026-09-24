# SDL — sporbar leveranseplan

Leveranseplanprofilen ble innført i V4 og inngår i **design-core 0.5** i Go.
Fire planrelasjoner og en Activity-egenskap beskriver hvilke ansvar en aktivitet
skal levere. En slik kildepåstand beviser ikke at arbeidet er utført.
[Go-innganger](../../go/README.md) og [G4-bevis](../../go/evidence/G4.md).

| Utsagn | Signatur og betydning |
| --- | --- |
| `A addresses F.` | Activity → Functionality; aktiviteten dekker dette implementasjonsansvaret |
| `A delivers F.` | Activity → Feature; planlagt leveransebidrag, ikke bevis på ferdig Feature |
| `A depends-on B.` | Activity → Activity; eksplisitt forutsetning, ingen selvreferanse/syklus |
| `S illustrates A.` | Scenario → Activity; eksplisitt eksempelbane knyttet til milepælen |
| `A has implementation-status = planned/implemented/verified.` | Activity-egenskap; en kildepåstand, ikke en status verktøyet utleder fra tester |

G-fasene og deres milepæler er Activity-identiteter. Milepælen `refines` sin fase.
Dette er en visningskonvensjon i eksemplet; generatoren kjenner ikke G-navn eller
antall faser. Den bruker relasjonene, ikke prefiks eller en innebygd plan.
`refines` gir ikke arvet status eller avhengighet. Tidsplan, varighet, ressurser
og automatisk utføring av arbeidsaktiviteter inngår ikke i profilen.

VP06 viser aktivitetsinndeling, leveranser, ansvar og en egen avhengighetsgraf.
`implementation.md` genereres fra de samme faktaene, med kilde-ID-er og eierkobling.
Rapporten markerer udekkede Functionality-er; den oppretter ingen milepæler for dem.
G1–G6-leveransene er siden gjennomført innen profilgrensene. Gjeldende
[designkilde](../../../SDUI/design/architecture.design) eier modellens eksplisitte
statuspåstander; [faseplanen](../../../SDUI/docs/implementation-plan.md) peker på
implementasjonsbevis. V4s opprinnelige planned-snapshot beholdes som historie under.

## Historisk V4-leveranse — 2026-09-22

Tall og status nedenfor gjelder denne milepælen før Go-porten. De er ikke
nykjørte tester eller dagens samlede implementasjonsstatus.

V4-milepæler: M1 sporbar fase-/ansvarsmodell; M2 generert leserapport og komplette
parser/presentasjon/binding/reload-scenarioer; M3 samlet verifikasjon, checkpoint,
fasepush og PR for gjennomgang mot sdp-vNow.

V4-M1 verifisert: 61 parsertester består. Modellen angir 5 planlagte G-faser,
18 milepæler og ansvarskoblinger for alle 94 Functionality-er. Ingen status er
oppgradert til implemented/verified for den kommende Go-koden.

V4-M2 levert: generisk implementation.md fra fase-/milepælrelasjoner, med
ansvarseiere og scenariofigurer. 24 verktøytester består. Modellens ti scenarioer
dekker kompilering, interaktiv/statisk presentasjon, lokal Go-handling, SDL-binding,
UI-/SDL-reload og native bygg. Dette er designbaner, ikke kjørte Go-programmer.

V4-M3 verifisert: 61 SDL-parsertester, 24 verktøytester og 36 SDUI-tester
består. 142 SVG-diagrammer fra 368 deklarasjoner og 1106 fakta er kontrollert
for kildekobling og byte-identisk reeksport. Femten ufullstendige modusallokeringer
rapporteres fortsatt eksplisitt. Fasegraf, parser- og bindingssekvens er visuelt
stikkprøvekontrollert; ingen Go-runtime eller fysisk print er verifisert.
