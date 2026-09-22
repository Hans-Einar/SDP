# Avgrenset Markdown-provider — G2-M4

Go-provideren parser Markdown én gang før måling. Måling og SVG bruker samme
blokker, tabellceller og diagramdimensjoner. Den brukes av SVG-CLI og Fyne-CLI;
konsolldumpens rå Markdown er fortsatt en annen, bevisst eksportform.

Støttet presentasjonsprofil:

- Overskrifter i seks grader, avsnitt, linjebrudd og enkle listepunkter.
- Inline strong/emphasis, kode og lenker beholder lesbart innhold uten Markdown-
  markørene. Denne første profilen gjengir ikke egne fontvekter, kursiv eller
  lenkenavigasjon; den påstår ikke full CommonMark-presentasjon.
- Inngjerdet/innrykket kode som literal tekst, og GFM-tabeller med cellegrenser,
  tekstombryting og samme mål i SVG som i layout.
- Mermaid-fences blir separate, identifiserte diagramressurser. Uten registrert
  renderer vises en uttrykkelig plassholder; de forsvinner ikke stille.

HTML og bilder avvises. Ingen URL lastes automatisk. Markdown-widgeten er
begrenset til 32 KiB, 256 blokker og åtte diagrammer. Eksterne ressurser,
animasjon, CSS og DOM-scripting er ikke del av profilen.

`--mermaid-renderer /absolutt/sti/til/mmdr` registrerer en lokal renderer for
SVG-eksport. Kilden går på stdin, aldri gjennom shell. Første verifiserte
kapabilitet er `flowchart`/`graph`; andre diagramtyper gir profilfeil. Kilden er
maksimalt 12 KiB, renderer får fem sekunder og maksimalt 4 MiB utdata. SVG må ha
endelig viewBox og være uten aktive eller eksterne elementer/referanser.
Initialiseringsdirektiver er ikke tillatt. Parseren importerer ingen renderer.

`--resources DIR` skriver innholdsnavngitte `.svg`-filer separat; den samlede
SVG-en har også data-ressurser og kan vises alene. Eksisterende ressurs med
samme navn må ha identiske bytes. Publisering med katalogtransaksjon, cache-
revisjoner og lease følger G6; M4 er ingen dokumentdaemon.

Fyne-CLI bruker Markdown-provideren uten ekstern diagramrenderer i denne fasen.
Diagrammer viser derfor den eksplisitte plassholderen der. Dette unngår å hevde
støtte for alle Mermaid-SVG-funksjoner i Fynes SVG-rasteriserer uten prøvebevis.
Den eksterne flowchart-ressursen er verifisert i SVG-eksport via librsvg.
