# Go-målekontrakt — G2

G2-M1 implementerer én måle-/plasseringsmodell i `go/layout`. SVG og native
verter får `Box` med rektangel, klipp, font, aktiv-status og instansbane. Ingen
callback kjøres ved måling. Kilden har ingen pikselmål.

`font` måles i absolutte logiske skjermenheter (DIP). `font=10` betyr ti slike
enheter både i SVG viewBox og verten. Enhetens fysiske størrelse følger vertens
DPI, ikke vindusstørrelsen. Referansemåleren bruker innbakt Go Regular gjennom
OpenType ved 72 DPI; Markdown kan velge relative overskriftsgrader. SVG må
bruke samme font eller oppgi eventuell fontsubstitusjon.

Rader endrer ikke ancestor for relative dimensjoner. Fr-spor fordeler rest etter
`clamp(lambda * weight,min,max)`; de legger ikke minimum oppå vekten. Skalerte
og innholdstilpassede søsken krympes ikke for å skjule overflow. Header/footer
reserveres før body; regioner bruker tilgjengelig bredde hvis ingen annen
x-policy er valgt. `items=stretch` lar en enkelt komponent i en rad fylle
bredden, og jevner høyden til komponenter i en flerkolonnerad.

En innholdstilpasset ancestor med relativt dimensjonerte barn på samme akse
avvises som `layout-dependency`. Breddebestemt ratio avleder høyden og trenger
ikke en høydebestemt ancestor. Min/max deformerer aldri ratio. Padding og gap
bruker nodens ancestor som referanse, slik språkforslaget beskriver.

Overflow `error` gir diagnose; `clip` klipper også trefftesting. Scroll krever
en vert med scrolltilstand og gir foreløpig eksplisitt `unsupported-scroll`.
Ingen usynlig implementasjon av scroll som clip. Målingen er begrenset til
200 000 rekursive operasjoner og viewport opp til 32 768 logiske enheter per
akse. Parserens egne kilde-/dybde-/utvidelsesgrenser gjelder først.

Prøvene dekker fr versus scale, minimum og maksimum, uavhengig font ved resize,
ratio versus contain, regionenes referanseområde, klippet trefftesting og den
faktiske Concept1-modellens seks bokser ved 1920 × 1200. Andre viewport-størrelser
kan avvise innhold som ikke passer; dette er forventet uten eksplisitt clip.
