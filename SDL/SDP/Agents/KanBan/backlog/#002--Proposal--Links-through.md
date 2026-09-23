---
id: KB-SDL-002
project: SDL
type: Proposal
created: 2026-09-23T18:25:13Z
source: owner-conversation-2026-09-23
next_review: 2026-09-30
---

# Svake links og visuelle stier gjennom noder

Registrert fra eierens samtale 2026-09-23. Tidspunktet er registreringstid,
ikke rekonstruert tidspunkt for tidligere diskusjoner. Status følger katalog/ledger.

## Avklart hensikt fra eieren

`links` er en beskrivende forbindelse; `uses`, `realizes` og andre spesifikke
relasjoner beholder sin egen semantikk. Svake forbindelser kan tegnes som stiplede
linjer med pil etter `to`/`from`. Referanser skal valideres, men forbindelsen er
ikke en forpliktelse, dataflyt, avhengighet eller bevis på implementasjon.

`through` angir én sammenhengende visuell sti via en ordnet liste av noder.
Stien går inn på en visuell port, tegnes synlig oppå noden og går ut på motstående
side før neste ledd. Dette er tilsiktet gjennomgang, ikke hindringsunngåelse eller
bare en linje bak en boks. Kan vise hvordan en Feature går gjennom arkitektur-
eller klasseblokker. Visuelle porter er ikke automatisk modellens Channel-porter.

## Foreløpig syntaksskisse

```text
Actor actor1 links to stakeholder1
Actor actor2 links through actor1 to stakeholder1
Feature feature1 links through container1, container2 to container3
```

Bevar én relasjon med kilde, retning, ordnede mellomledd og mål i AST/modell.
Ikke senk til uavhengige semantiske `uses`/`realizes`-kanter. Hele stien må kunne
kjennes igjen med etikett og sammenhengende markering, ikke farge alene.

## Åpent og neste arbeid

Avklar endelig grammatikk, `from` i kombinasjon med `through`, flere relasjoner,
identitet, tillatte gjentatte noder/sykler og diagnoser. Avklar portplassering ut
fra diagramretning, kollisjoner med tekst, kryssende stier og delvis skjulte noder
i viewpoints. En eksport som ikke kan tegne gjennomgang må oppgi begrensningen.
Undersøk eksisterende SDL-projektor/renderadapter før backendvalg eller endringer
i Mermaid-repoet; ingen rendererendring er autorisert av selve registreringen.

## Akseptanse for senere implementasjon

Parser/AST bevarer rekkefølge og kildeposisjoner; ukjente referanser får diagnose.
En figur med minst to mellomnoder viser innport, synlig passasje over hver node,
utport på motstående side og korrekt pilretning. Stien må overleve viewpoint-
utvalg med eksplisitt håndtering av utelatte noder. Status/dekningsberegning
ignorerer forbindelsen som oppfyllelsesbevis.
