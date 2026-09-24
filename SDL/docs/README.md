# SDL — dokumentoversikt og autoritet

SDL-koden ligger nå i SDL/go. R1 har flyttet dokumentene etter ansvar; det er
ikke en ny språkversjon. [Go-README](../go/README.md) og de aktive parserprofilene
angir kjørbart omfang. [SDL KanBan](../SDP/Agents/KanBan/README.md) inneholder
forslag som ennå ikke er språkregler.

## Aktive profiler og grenser

Design-core 0.5, action-core 0.1 og class-core 0.1 er separate avgrensede profiler.

- [Strukturell kjerne](studies/Design-Language-Definition.md): §§1–11 og registrerte tillegg; §12 og videre er åpent arbeid/kandidater.
- [Data-/wirekontrakter](profiles/SDL-Data-Contract-Profile.md).
- [Channels og scenarioer](profiles/SDL-Channel-Scenario-Profile.md).
- [Leveranseplanfakta](profiles/SDL-Delivery-Plan-Profile.md).
- [Kjørbare handlinger](profiles/SDL-Executable-Action-Profile.md) og [grammatikk](../grammar/action-core-0.1.ebnf).
- [Klasser](profiles/SDL-Class-Profile.md) og [grammatikk](../grammar/class-core-0.1.ebnf).
- [Symbolprofil](profiles/SDL-Symbol-Profile.md).

Profilfilene skiller aktive regler fra daterte V2–V4-leveransebevis. Gamle
testtall og planned-status gjelder de navngitte milepælene. Nåværende
[Go-kommandoer](../go/README.md) og [implementasjonsplan](../../SDUI/docs/implementation-plan.md)
leder til senere G-fasebevis; en modellpåstand er ikke alene verifikasjon.

## Integrasjon, studier og historikk

[Go-portoversikt](integration/SDL-Go-Port-Inventory.md),
[navigasjonsdesign](integration/SDL-Navigable-Viewpoints-Design.md) og
[nivåer/notasjon](integration/SDL-Viewpoint-Levels-and-Notation.md) dokumenterer
gjenbruk og opprinnelig designgrunnlag. G6 er implementert innen avgrensede
profiler; de daterte designtekstene har ikke selvstendig myndighet til å utvide dem.
[G6-bevis](../go/evidence/G6.md) og [launcher](../scripts/README.md) viser levert omfang.

[Design Language Definition](studies/Design-Language-Definition.md) inneholder
både aktiv kjerne og tydelig avgrensede forslag; katalogen studies/ alene avgjør
ikke autoriteten. De øvrige studiene utforsker bredere begreper enn aktive profiler.
[Research](../research/README.md) sammenligner andre designspråk.
[MVP1-korpuset](../../experiments/mvp1_sdl/README.md) er en kandidatøvelse.
[Checkpoint #1](../../SDP/History/checkpoint-1/README.md) er et felles datert
SDP/SDL/SDUI-snapshot, ikke en samlet aktiv SDL-spesifikasjon.

Alle tidligere docs-filer finnes i [dokumentkartet](../../SDP/Maintenance/R1/Documentation-index.md).
