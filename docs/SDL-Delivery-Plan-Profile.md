# SDL — sporbar leveranseplan

V4, 2026-09-22. design-core 0.5 erstatter aktiv 0.4 med tre avgrensede
planrelasjoner og en Activity-egenskap. SDL-modellen kan nå beskrive hvilke
ansvar en planlagt utviklingsaktivitet skal levere, uten å påstå at de er utført.

| Utsagn | Signatur og betydning |
| --- | --- |
| `A addresses F.` | Activity → Functionality; aktiviteten dekker dette implementasjonsansvaret |
| `A delivers F.` | Activity → Feature; planlagt leveransebidrag, ikke bevis på ferdig Feature |
| `A depends-on B.` | Activity → Activity; eksplisitt forutsetning, ingen selvreferanse/syklus |
| `A has implementation-status = planned/implemented/verified.` | Activity-egenskap; en kildepåstand, ikke en status verktøyet utleder fra tester |

G-fasene og deres milepæler er Activity-identiteter. Milepælen `refines` sin fase.
Dette er en visningskonvensjon i eksemplet; generatoren kjenner ikke G-navn eller
antall faser. Den bruker relasjonene, ikke prefiks eller en innebygd plan.
`refines` gir ikke arvet status eller avhengighet. Tidsplan, varighet, ressurser
og automatisk utføring av arbeidsaktiviteter inngår ikke i profilen.

VP06 viser aktivitetsinndeling, leveranser, ansvar og en egen avhengighetsgraf.
`implementation.md` genereres fra de samme faktaene, med kilde-ID-er og eierkobling.
Rapporten markerer udekkede Functionality-er; den oppretter ingen milepæler for dem.
All G1–G5-status i SDL/SDUI-designet er planned. V2–V4s verktøyleveranser og bevis
føres i checkpoint og Git, og skal ikke forveksles med Go-implementasjon.

V4-milepæler: M1 sporbar fase-/ansvarsmodell; M2 generert leserapport og komplette
parser/presentasjon/binding/reload-scenarioer; M3 samlet verifikasjon, checkpoint,
fasepush og PR for gjennomgang mot sdp-vNow.

V4-M1 verifisert: 61 parsertester består. Modellen angir 5 planlagte G-faser,
18 milepæler og ansvarskoblinger for alle 94 Functionality-er. Ingen status er
oppgradert til implemented/verified for den kommende Go-koden.
