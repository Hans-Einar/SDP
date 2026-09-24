# SDL — dokumentoversikt og autoritet

SDL-koden ligger nå i SDL/go. R1 har flyttet dokumentene etter ansvar; det er
ikke en ny språkversjon. [Go-README](../go/README.md) og de aktive parserprofilene
angir kjørbart omfang. [SDL KanBan](../SDP/Agents/KanBan/README.md) inneholder
forslag som ennå ikke er språkregler.

## Aktive profiler og grenser

Design-core 0.5, action-core 0.1 og class-core 0.1 er separate avgrensede profiler.

- [Data-/wirekontrakter](profiles/SDL-Data-Contract-Profile.md).
- [Channels og scenarioer](profiles/SDL-Channel-Scenario-Profile.md).
- [Leveranseplanfakta](profiles/SDL-Delivery-Plan-Profile.md).
- [Kjørbare handlinger](profiles/SDL-Executable-Action-Profile.md) og [grammatikk](../grammar/action-core-0.1.ebnf).
- [Klasser](profiles/SDL-Class-Profile.md) og [grammatikk](../grammar/class-core-0.1.ebnf).
- [Symbolprofil](profiles/SDL-Symbol-Profile.md).

Profilfilene inneholder også daterte leveranseeksempler. Eksempelvis beskriver
leveranseplanprofilens V4-avsnitt eldre planned-status; det overstyrer ikke senere
G-fasebevis eller gjeldende modell. Slike avsnitt konsolideres i KB-SDP-010.

## Integrasjon, studier og historikk

[Go-portoversikt](integration/SDL-Go-Port-Inventory.md),
[navigasjonsdesign](integration/SDL-Navigable-Viewpoints-Design.md) og
[nivåer/notasjon](integration/SDL-Viewpoint-Levels-and-Notation.md) dokumenterer
retning og grunnlag. Der eldre tekster sier «planlagt», se faktisk implementasjon
og begrensninger i Go-README og G-fasebevis før status tolkes.

[Design Language Definition](studies/Design-Language-Definition.md) og resten av
studies/ er utforskning med bredere begreper enn den implementerte profilen.
[Research](../research/README.md) sammenligner andre designspråk.
[MVP1-korpuset](../../experiments/mvp1_sdl/README.md) er en kandidatøvelse.
[Checkpoint #1](../../SDP/History/checkpoint-1/README.md) er et felles datert
SDP/SDL/SDUI-snapshot, ikke en samlet aktiv SDL-spesifikasjon.

Alle tidligere docs-filer finnes i [dokumentkartet](../../SDP/Maintenance/R1/Documentation-index.md).
