# Stakeholders, actors, user stories og lesbar SDL

| Felt | Verdi |
| --- | --- |
| id | KB-SDL-001 |
| project | SDL |
| type | Proposal |
| created | 2026-09-23T18:25:13Z |
| source | owner-conversation-2026-09-23 |
| next_review | 2026-09-30 |

Registrert fra eierens samtale 2026-09-23. Tidspunktet er registreringstid,
ikke rekonstruert tidspunkt for tidligere diskusjoner. Status følger katalog/ledger.

## Eierens forslag

Start kravarbeidet direkte i SDL, med `01--Actors.design`, deretter
`02--UserStories.design` og `03--UseCases.design`. Stakeholders listes først;
actors kan kobles til stakeholders, men kobling er ikke obligatorisk.
Prosessens plassering/studier eies av [KB-SDP-001 — Prosjektstruktur, Template og studier per fase](../../../../../SDP/Agents/KanBan/active/%23001--Proposal--Project-structure.md).

Stakeholder har interesser i systemet. Actor er en rolle eller ekstern enhet
som samhandler med det. Skill menneskelig bruker fra maskin/programvaresystem:
P1000 er en maskinaktør, ikke en User. Avklar eksplisitte aktørtyper og faste
visuelle kjennetegn; ikke bestem type ut fra navnet eller tegn alle som mennesker.

UserStory har identitet, fritekst og koblinger til actors og senere use cases.
Teksten er informasjonsgrunnlag, ikke automatisk kjørbar logikk eller utledede krav.
Foreslå engelsk stavemåte `UserStory` (eierens `userstorry` er samtaleskisse).

## Syntaksskisser — ikke gyldiggjort som implementert SDL

```text
Stakeholder stakeholder1
Actor actor1 links to stakeholder1
UserStory story1 links to actor1
story1.story = "Som operatør ønsker jeg ..."

UserStory story2 links to actor1
.story = "Som operatør ønsker jeg ..."
.comment = "Beskrivelse som skal vises i dokumentasjonen."
```

Eieren ønsker både eksplisitt objektegenskap og vurdering av en kortform med
innledende punktum. Anbefaling til avklaring: kortformen gjelder bare nærmeste
foregående deklarasjon i samme blokk, ikke vilkårlig tidligere bruk av en referanse.
Avklar regler ved import, blokkgrense, blanklinjer, blandede utsagn og filslutt.

`comment` skal være dokumenterbart beskrivelsesinnhold, adskilt fra kommentarer
som bare er for kildekoden. Keyword kontra egenskap er ikke besluttet.
Flerlinjet tekst (quotes/braces/brackets) og import som eierens skisse
`#include 01--Actors.design` må sammenholdes med faktisk grammatikk før valg.
Ingen include-preprosessor eller friteksttolkning innføres ved å skrive dette kortet.
Svake forbindelser og visuelle stier eies av [KB-SDL-002 — Svake links og visuelle stier gjennom noder](%23002--Proposal--Links-through.md).

## Neste arbeid og akseptanse

Undersøk aktiv parser/profil før grammatikkforslag. Avklar ord, casing, aktørtyper,
multiline/import og kortformsbinding. Deretter versjonert profil, AST/kildeposisjoner,
validering, negative eksempler og dokumentprojeksjoner. En lesbar kravpakke skal
bevare tekst nøyaktig og vise stakeholder–actor–story–use-case-sporbarhet uten
at en løs forbindelse tolkes som kravoppfyllelse.

[SystemDesignLanguage/README.md](../../../../../SystemDesignLanguage/README.md)
