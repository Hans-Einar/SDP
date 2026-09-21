# Concept1 — referanse for relative flater

Kontrollert 2026-09-20 mot lokal Ponsse HEAD 882ad7c. Leste filer hadde ingen
rapporterte lokale endringer; andre lokale endringer ble bevart. Ingen GUI-kjøring.
Relative stier nedenfor gjelder /home/warloc/git/ponsse.

## Faktisk implementasjon

| Kilde under Concept1 | Observasjon |
| --- | --- |
| shared/ui-box/model.mjs | layoutTrackTemplate mapper positive weight til CSS fr. |
| shared/ui-box/UILayout.jsx | Grupper bruker CSS Grid i rad/kolonne; ID kobles til separat React-innhold. aspectRatio overføres som CSS-variabel. |
| apps/operator-ui/src/ui/operator-layout.mjs | Ratio 16/9, vertikale vekter 15/45/40 og midtfelt 25/50/25. |
| shared/ui-box/fixed-aspect-model.mjs | 4:3 og 16:9 har interne designflater 1600×1200 og 1600×900. Fit bruker min(viewportWidth/designWidth, viewportHeight/designHeight) og sentrerer. |
| shared/ui-box/FixedAspectViewport.jsx | Følger hovedvinduets resize og legger én samlet scale-transform på designflaten; geometry-authority=outer-only. |
| shared/ui-box/UIAspectRegion.jsx | Underregion har scaling-authority=none; ingen egen nested scaler. |
| shared/ui-box/ui-box.css | Flater fyller 100%; gap/padding/header og kontrollstørrelser bruker fortsatt CSS px. |
| apps/operator-ui/src/ui/operator-ui.css | Overstyrer ui-layout-frame aspect-ratio til auto; ytre viewport har viktig ansvar. |

Dette er relativ sporfordeling og samlet skalering av en intern designflate.
Det er ikke en eksisterende implementasjon av vilkårlig ancestor-relativ
skalering på hver SDUI-frame, og ikke en internt pikselfri løsning.

## Overføring til SDUI

Ta med vekter, sideforhold, vertseid viewport og én konsistent transform for
bilde/kontroller/input. Barn skal ikke måle hovedvinduet på egen hånd.
Ikke kopier faste 1600-baserte mål, root-only-skalering eller CSS px til kildespråket.

Scale-x/scale-y/scale bestemmer frameutstrekning relativt til nærmeste kildeancestor.
Eierpresisering 2026-09-21: rollen til FixedAspectViewport uttrykkes med en
ytterframe `{16:9,<->}`. Den bruker hele bredden og avleder høyden, selv om
vinduet er for lavt. Dette skiller seg fra Concept1s min-baserte contain-fit.
SDUI skal dessuten beholde absolutt fontstørrelse ved resize og gjøre ny layout/
tekstombryting. Concept1s samlede innholdstransform skal derfor ikke kopieres.
Felles native fontenhet konkretiseres før måleimplementasjonen.

## Regneeksempler for G2

Normaliserte referanseenheter; ingen pixeldimensjoner i SDUI-kilden.

| Referanse | Egenskaper | Forventet frameutstrekning |
| --- | --- | --- |
| 1 × 1 | 16:9, scale-x=0.8 | 0.8 × 0.45 |
| 1 × 1 | 16:9, scale-y=0.5 | 8/9 × 0.5 |
| 1 × 0.75 | 16:9 uten scale, foreslått contain | 1 × 0.5625 |
| 1 × 0.25 | 16:9, <-> | 1 × 0.5625: høydeoverflow, bredden beholdes |
| 1 × 1 | 4:3, scale-x=0.6 | 0.6 × 0.45 |
| Mor 0.8 × 0.45 | Barn uten ratio, scale=0.5 | 0.4 × 0.225; innholdet skaleres ikke av frame-scale |
| 1 × 1 | 16:9, scale=0.8 | Avvis: begge akser eksplisitt styrt |
| 1 × 1 | 16:9, scale-y=1 | 16/9 × 1: overflow, ingen stille deformering |

Dette er testorakler, ikke resultater fra en implementert SDUI-layoutmotor.
Syntetiske rader må ikke endre ancestor-referanse; nye eksplisitte foreldre gjør det.
