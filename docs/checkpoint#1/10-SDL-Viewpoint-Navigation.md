# Checkpoint #1 — G6 dokumentnavigasjon

Dato: 2026-09-22. Eierbestilling: del viewpoints i navigerbare dokumenter og
planlegg generering ved klikk, med navigator og hoveddokument i XFMD.
**Design er lagt inn; G6-M1–M6 er ikke implementert.**

Les [designkontrakten](../SDL-Navigable-Viewpoints-Design.md) for katalogforslag,
URI-/adaptergrense, vindusvalg, publisering og levetid. Se
[den genererte planen](../../SDUI/design/viewpoints/implementation.md) for G6s
ansvar, eiere, avhengigheter og scenarioer fra [SDL-kilden](../../SDUI/design/architecture.design).
[Samlet utskrift](../../SDUI/design/viewpoints/printout.md) er regenerert av SDL-verktøyet.

De første fire milepælene følger statisk navigasjon → behovsgenerering → XFMD-paneler →
valgfri bakgrunnstjeneste. G6 kan starte etter G4-M1s strukturelle SDL-port og
trenger ikke domenekjøring eller native Go-generering. G5-M3 gjenbruker eksporten. G6-D2 legger til M5 for fast diagramnotasjon
og M6 for senere klassediagrammer med eksplisitt språk-/relasjonsprofil.
[Nivå- og notasjonsdesignet](../SDL-Viewpoint-Levels-and-Notation.md) beskriver
navigator/overview som egen eksportform, A0–A5, typekataloger og Mode/State.

Navigatoren velger et utsnitt av SDL-modellen. Full rapport og enkeltsider
skal bruke samme prosjektor, ikke tolkning av ferdig Markdown. Generering holder
seg i SDL-verktøyet; vindu/panel og lenkehandler tilhører XFMD. Dagens eksport-
format er uendret. Fyne er fortsatt første vert for interaktive SDUI-widgets.

Valgt leser er en lokalt konfigurert adapter-ID, ikke en vilkårlig kommando i
Markdown. Det konkrete målvinduet fanges ved klikket. --active og URI-skjemaet er
forslag til XFMD-kontrakt, ikke eksisterende funksjonalitet. Filbaserte lesere
får små immutable dokumentpakker med bilder i en privat runtime-katalog;
ren RAM-overføring er en senere adaptermulighet. Ingen daemon startes her.

SDL 0.5 uttrykker plan, ansvarsdeling og typede meldinger uten grammatikkendring.
SelectedViewOpened viser den planlagte banen; InvalidViewSelectionRejected og
ViewProjectionFailed viser avvisning før åpning. Modellen utfører ikke
protokollen eller kontrollerer alle feltverdier, leases og ressursgrenser.

Verifisert: 121 tester består (61 SDL, 24 verktøy, 36 SDUI). Modellen har
463 deklarasjoner, 1424 fakta, 116 ansvar og 24 planlagte milepæler. Verktøyet
genererer 170 SVG-diagrammer, med byte-identisk reeksport og kontrollerte
kildekoblinger. G6-fasegrafen og åpningssekvensen er visuelt kontrollert.
VP07 rapporterer 18 manglende modusallokeringer; tre nye gjelder delte
dokumentasjonsansvar i DocumentBrowsing. Ingen full deployment hevdes.
Verifikasjon registreres i [maskinrapporten](../../SystemDesignLanguage/tools/verification.json).
Reproduser fra repoets rot:

```sh
python3 SystemDesignLanguage/tools/verify_design.py --phase G6-D2 --renderer /home/warloc/git/mermaid-rs-renderer/target/debug/mmdr
```

[Branchstakken](../Development-Branch-Stack.md) registrerer G6-D1/D2 som en egen
designleveranse oppå V4. Dette er ikke fullføring av G6s implementasjonsfase.
