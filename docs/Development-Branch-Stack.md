# SDL/SDUI — fasebrancher og milepælcommits

Eierbeslutning 2026-09-22: én branch per fase, stablet på forrige fase, og
egne commits ved fullførte milepæler. Push til origin etter hver fullført fase
er autorisert i samme sesjon. Samlet PR skal ha `sdp-vNow` som base.
Utgangspunktet er `9ad432407004080dd7f4f0ab06d107523f4316fd`.

| Fasebranch | Forelder | Milepæler |
| --- | --- | --- |
| `sdl-sdui/phase-baseline` | `sdp-vNow` | B0-M1: samle eksisterende SDUI 0.2-prototype, eksempler, bevis og Go-retning; etablere Git-arbeidsmåten |
| `sdl/phase-v1-viewpoints` | `sdl-sdui/phase-baseline` | V1-M1: språk/AST/validering; V1-M2: verktøyprojeksjoner; V1-M3: felles modell, generert utskrift og checkpoint |
| `sdl/phase-v2-data-contracts` | `sdl/phase-v1-viewpoints` | V2-M1: data-/wireprofil; V2-M2: VP09/VP10; V2-M3: verifisert prøve og fasepush |
| `sdl/phase-v3-channels` | `sdl/phase-v2-data-contracts` | V3-M1: kontraktkontrollerte scenarioer; V3-M2: sekvens og MessageSet; V3-M3: renderbevis og fasepush |
| `sdl/phase-v4-integrated-design` | `sdl/phase-v3-channels` | V4-M1: G1–G5 som SDL; V4-M2: generert implementasjonsrapport; V4-M3: samlet eksport, checkpoint, fasepush og PR |
| `sdl/phase-g6-navigation-design` | `sdl/phase-v4-integrated-design` | G6-D1: navigasjon og scenarioer; G6-D2: navigator-only, abstraksjonsnivåer og notasjon; G6-M1–M6 fortsatt planlagt |
| `sdui/phase-g1-frontend` | `sdl/phase-g6-navigation-design` | G1-M1: parser/AST; G1-M2: lokal semantikk; G1-M3: porterte eksporter |

De første to fasene samler arbeid som allerede var utført lokalt. Dette er en
bevisst etterregistrering, ikke en påstand om at gamle økter hadde disse commit-
grensene. V0-generatoren og V1-utvidelsen var allerede utviklet i de samme filene;
de registreres samlet på V1-branchen. Ingen hypotetisk tidligere V0-kode bygges
opp for å skape en mer detaljert historie enn kildegrunnlaget tillater.

Baseline bevarer eksisterende SDUI-kode og dens dokumenterte retning. Enkelte
plan-/dokumentlenker peker fram til SDL-leveransen som følger på V1-branchen.
De tre V1-milepælene er henholdsvis språkgrunnlag, konsument og validert eksempel
med bevis; fasehodet er den samlede leveransen som skal gjennomgås.

Fra V2 følger arbeidet milepælene fortløpende. Før hver commit kontrolleres
avgrensningen, relevante tester og genererte artefakter. Fasebranch opprettes
før første endring i fasen; bare fullførte milepæler får leveransecommits.
V2/V3/V4s faglige omfang står i
[checkpointets implementasjonsplan](checkpoint%231/08-SDL-Viewpoints-and-Implementation-Status.md).
Go-arbeidet følger fortsatt G-fasene i [SDUI-planen](../SDUI/docs/implementation-plan.md).

Ved samlet PR brukes siste ferdige fasebranch som head og `sdp-vNow` som base.
Milepælcommits skal bevares ved integrasjon; ikke squash dem til én commit.
Branchene slettes eller historikken omskrives ikke som automatisk opprydding.

G2-implementasjonen fortsetter fra G1 på `sdui/phase-g2-presentation`.
Milepælbevis: `SDUI/go/evidence/G2.md`.

G3 fortsetter fra ferdig G2 (`6a4d968`) på `sdui/phase-g3-runtime`.

G4 fortsetter fra ferdig G3 (`d29a48f`) på `sdl/phase-g4-runtime`.

G6 implementeres fra ferdig G4 (`d5430f0`) på `sdl/phase-g6-navigation`,
før G5 fordi dokumentgenereringen i G5-M3 konsumerer G6-M1.

G6 er fullført på `1d52677`. G5 følger på `sdl-sdui/phase-g5-codegen`.
XFMDs nødvendige konsumentendringer er isolert i PR #38 i XFMD-repoet;
SDP- og XFMD-branchene er pushet, ikke merget.


| Implementasjonsfase | Pushet fasehode / branch |
| --- | --- |
| G1 | ca5aa91 — sdui/phase-g1-frontend |
| G2 | 6a4d968 — sdui/phase-g2-presentation |
| G3 | d29a48f — sdui/phase-g3-runtime |
| G4 | d5430f0 — sdl/phase-g4-runtime |
| G6 | 1d52677 — sdl/phase-g6-navigation |
| G5 | sdl-sdui/phase-g5-codegen; siste commit er G5-M4 |

G5 har M1 konstruktører, M2 kjøreparitet, M3 reproduserbar dokumentasjon og M4
portopprydding/checkpoint. Samlet PR bruker G5-hodet mot sdp-vNow. Opprinnelige
fasebrancher og milepælcommits beholdes. Ingen merge er en del av leveransen.
