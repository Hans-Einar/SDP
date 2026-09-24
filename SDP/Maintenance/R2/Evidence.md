# R2 — kontrollbevis

## R2-M1

R1-baseline og dagens dokumentinnganger er lest; ti konkrete konflikter er
registrert i Findings. KanBan schema/replay/filplassering og lokale lenker
består: 16 kort, 27 hendelser og 1907 lokale Markdown-filmål. Historiske
ledgerprefikser, 574 genererte artefakter og R1s migreringsgrenser er bevart.
`git diff --check` består. Ingen kode eller genererte modeller er endret.

## R2-M2

Dokumentendringer kontrollert mot SDL parser/vocabulary, bridge.Plan og
SDUI Handle, samt eksisterende G3/G4/G6/G5-bevis. Ingen endring i grammatikk,
Go-kilde, moduler, installasjonsmanifest eller genererte viewpoints.

Kjørte CLI-stikkprøver med lokal Go 1.27.1 fra repoets rot:

```sh
go -C SDL/go run ./cmd/sdl check ../../SDUI/design/architecture.design
go -C SDL/go run ./cmd/sdl action-check examples/echo.sdl
go -C SDL/go run ./cmd/sdl class-check examples/runtime-classes.sdl
go -C SDUI/go run ./cmd/sdui ../examples/concept1-bucking.sdui --format svg --entry bucking -o /tmp/concept1.svg
```

Alle tre SDL-profiler ga valid=true. SVG-kommandoen ble kjørt med unik midlertidig
utkatalog: gyldig XML/SVG, viewBox 0 0 1920 1200, 402669 bytes. Dette er strukturell
CLI-verifikasjon, ingen ny visuell eller native GUI-prøve. G-fasenes gamle
testtall er ikke nykjørt. Full Go-regresjon er ikke gjentatt for kun dokumentendringer.

R1-kontrollen validerer lokale Markdown-filmål, 574 uendrede genererte artefakter,
historiske fingeravtrykk, Go-kilder og installasjonspolicy. KanBan kontrolleres
med K3s schema/replay-funksjon og fysisk kortplassering på tre tavler; ledgerens
prefikser fra K4 beholdes byte for byte. Avsluttende kontrolltall føres under.

Begrensninger: lenkekontrollen tester ikke ankre eller eksterne URL-er. Ingen ny
påstand om gjeldende status i andre repoer/PR-er er lagt inn; G6s eksterne bevis
er datert. De ni eldre Toolkit-ID-avvikene i KB-SDP-011 er ikke endret eller
retestet her. Fase-/malprofil og større kandidatsemantikk forblir åpne.

Avsluttende kontroller bestod: 16 kort / 29 hendelser på tre tavler;
1973 lokale Markdown-filmål; alle 17 faglige SDUI-docs indeksert. Git-diff mot
K4 bekrefter uendrede History-, Go-, design-, Template- og installasjonsfiler.
Kun Markdown og append-only ledgerhendelser inngår i R2. `git diff --check` består.
