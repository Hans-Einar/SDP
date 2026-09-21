# Treemap som statisk UI-skisse — 2026-09-21

Eierens hensikt er et samlet boksoppsett med innhold, tilsvarende tekstdumpen.
Den tidligere Markdown-rapporten med ASCII og separate sitatblokker oppfyller
ikke dette visuelle målet.

Kontroll av lokal mermaid-rs-renderer-kode viste at treemap fordeler horisontalt
ved partallsdybde og vertikalt ved oddetallsdybde. En rot med tre undergrupper
og blader i disse gir ønsket rad-/kolonnemønster. Foreldre summerer barnas vekter.
Prøvens blader 7.5+7.5 / 11.25+22.5+11.25 / 40 gir radene 15/45/40.

Den eksisterende lokale mmdr 0.3.1-binæren ble brukt, uten bygg eller kodeendringer
i rendererrepoet. Binærhash og kontroller finnes i verification.json. Repoets
HEAD var afab5e9; dette er ikke bevis på hvilken kilde binæren opprinnelig ble bygget fra.

Kjørte operasjoner: mmdr med view.mmd og config.json → SVG + layout.json;
rsvg-convert → PNG; PNG inspisert visuelt. Assertions på layout.json kontrollerte
radvekter, midtkolonner, like toppkolonner og ikke-tomme etiketter i alle seks blader.
SVG-en ligger i ../../examples/concept1-bucking.treemap.svg og vises fra tilhørende .md.

Viktig funn: lengre toppetiketter ble helt borte fordi de målte tekstblokkene
ikke fikk plass. Kortere etiketter ga seks synlige bokser med tekst. Vektene
blir også skrevet som tall i bladene. Metadata/padding gjør at de indre boksenes
høyder ikke er eksakt 15/45/40 av hele rotflaten.

Dette er en manuelt avgrenset prøve, ingen generell AST→treemap-eksport og ingen
full Markdown-rendering inni boksene. Ingen garanti om identisk plassering i
andre Mermaid-motorer. XFMD dokumenterer lokal SVG-bildestøtte, men denne økten
kjørte ikke XFMD GUI. Markdown bruker derfor forhåndsrendret SVG. Ingen runtime,
interaktive kontroller, nye parserregler, commit eller push.
