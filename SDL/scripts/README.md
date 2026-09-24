# Åpne SDL-dokumentasjonen

Kjør `sdl-design` når kommandoen er lagt i PATH. Scriptet kan også kjøres direkte
fra denne katalogen som `./sdl-design`, eller fra repo-roten som
`./SDL/scripts/sdl-design`. Arbeidskatalogen ellers er irrelevant.

Scriptet bruker **ferdigbygde** programmer. Det kompilerer ikke Go, starter ikke
daemon og åpner ikke den forhåndsgenererte full-eksporten. Det gjør følgende:

1. Leser prosjektets `SDUI/design/architecture.design` med det registrerte SDL-verktøyet.
2. Genererer navigator og oversikter i en privat midlertidig katalog (14 Markdown-sider og ett manifest, ingen detaljdiagrammer).
3. Åpner hovedsiden `index.md` i XFMDs hovedpanel og `navigator.md` i sidepanelet.
4. Registrerer verktøy, kilde, prosjekt og Mermaid-renderer. Klikk på `sdl-view://` kjører SDL-verktøyet mot kilden slik den er ved klikket.
5. Fjerner oppstartsfilene når vinduet lukkes. XFMD eier levetiden til detaljvisningene.

Navigatoren er et snapshot fra oppstart; start på nytt for å oppdatere dens liste
etter strukturelle modellendringer. Detaljutvalg validerer alltid gjeldende SDL.
Ugyldig kilde gir en diagnose og beholder siste visning.

Standardprogrammer og valgfrie overstyringer:

| Variabel | Standard |
| --- | --- |
| SDP_SDL_TOOL | `$HOME/.local/lib/sdp/sdl` — ferdigbygd SDL CLI |
| SDP_XFMD | `xfmd-sdl-navigation/build/xfmd` i søsterrepoet |
| SDP_MMDR | `mermaid-rs-renderer/target/debug/mmdr` i søsterrepoet |
| SDP_SDL_SOURCE | Dette repoets `SDUI/design/architecture.design` |

`--help` viser bruken. Ett filargument velger en annen SDL-kilde. Manglende
programmer gir en konkret feil; det skjer ingen automatisk bygging/installasjon.
`--navigator` og `--sdl-tool` er XFMD-flagg som krever filnavn etter seg;
launch-scriptet fyller dem ut. Det vanlige installerte XFMD-programmet velges
ikke automatisk, siden det kan mangle navigasjonsintegrasjonen fra PR #38.

Full eksport under `SDUI/design/viewpoints` er det lagrede G5-verifikasjons-/
eksportresultatet. Den er ikke inngangen for daglig kildebasert browsing.
