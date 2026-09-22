# SDL — Channel-kontrakter og deklarerte scenarioer

V3, 2026-09-22: design-core 0.4 erstatter aktiv 0.3. Dette er strukturell
protokollkontroll og generering av sekvensdiagrammer, ikke en runtime.
Det videre workspace-/Channel-forslaget i språkdefinisjonen er fortsatt en studie.

Nye typer: `channel`, `message`, `scenario`. Channel er logisk samarbeid,
og deltakere kan være Units eller Containers. Bibliotekkall blir ikke automatisk
nettverkstrafikk. Scenario beskriver en eksplisitt ordnet eksempelbane.

| Utsagn | Betydning |
| --- | --- |
| `C upholds K.` | Channel → Contract, nøyaktig én kontrakt |
| `K permits M.` | Contract → Message/Datagram; tillatte meldingstyper/familier |
| `M upholds K.` | Message → record-Contract; ingen Dataset kreves for en handling |
| `M has message-kind = request/result/event.` | Meldingsrolle; obligatorisk for Message |
| `Result replies-to Request.` | Message → Message; én request-type per result-type |
| `P uses C as sender of M in mode Mode.` | Unit × Channel × sender/receiver × Message/Datagram × Mode |
| `S runs-in Mode.` | Scenario → Mode, nøyaktig én |
| `S exercises Goal.` | Scenario → UseCase, eksplisitt sporbarhet; ingen kjørbarhetsbevis |
| `S has completeness = closed/open.` | Closed krever svar på alle requests; open kan vise en avgrenset prefix |

```text
S step 1 sends Request from Client to Server via Service.
S step 2 sends Result from Server to Client via Service reply-to 1.
S step 3 sends Notices variant Changed from Server to Observer via Events.
```

Stegindekser er eksplisitte positive heltall, unike og sammenhengende fra 1.
Kanonisk tekstsortering gir ikke tidsrekkefølge; generatoren bruker stegindeksene.
Datagram krever en variant som tilhører familiens kontrakt. Vanlig Message skal
ikke ha variant. Alle brukte kontrakter må være closed før sekvensen kan bygges.
Open modeller kan fortsatt deklarere Channel/Message uten kjørbart scenario.

En step må finnes blant Channel-kontraktens permits og samsvare med sender- og
receiver-deltakelse for meldingstypen i scenarioets eksplisitte Mode. Ingen arv
av deltakere, kontrakter eller modi. Selvsend er tillatt som et uttrykkelig lokalt
kall. Flere receivers krever separate steps; ingen broadcast utledes.

Result krever reply-to til en tidligere request i samme scenario, på samme
Channel og med reverserte deltakere, og result-typen må ha riktig replies-to.
Én requeststep får høyst ett result. Request/event/Datagram kan ikke bruke
reply-to. Dette kontrollerer eksempelbanens korrelasjon, ikke runtime-ID-er,
leveringsgarantier, tid, reentrancy, transaksjoner eller faktisk utføring.

Channel-kontraktens form er permits, uten record-felt/varianter. Message bruker
record-felt. Et Contract kan ikke samtidig være Channel-kontrakt og payloadform.
Alle Message-resulttyper krever replies-to; øvrige message-kinds kan ikke ha det.
MessageSet er kun avledet eksport per Channel/Mode fra permits og deltakelse;
ingen selvstendig authored MessageSet-syntaks innføres.

EBNF-tillegg:

```text
participation = identifier, "uses", identifier, "as", ("sender" | "receiver"),
                "of", identifier, "in", "mode", identifier, "." ;
step = identifier, "step", integer, "sends", identifier,
       [ "variant", identifier ], "from", identifier, "to", identifier,
       "via", identifier, [ "reply-to", integer ], "." ;
```

AST: `Participation(subject, channel, role, message, mode, span)` og
`Step(subject, ordinal, message, variant, sender, receiver, channel, reply_to, span)`.
Identifier-/Integer-noder har kildeposisjoner. Sekvenspil og MessageSet-rad får
kilde-ID-er til step, deltakelse, governing kontrakt og permits.

Alternativer modelleres som navngitte separate scenarioer, for eksempel godtatt
eller avvist input. Branch/loop/parallel-syntaks utsettes til semantikken er
definert. Et diagram er en avtalt eksempelbane, ikke hele tillatte protokollen.

Milepæler: V3-M1 språk og negative kontraktprøver; V3-M2 VP08/MessageSet;
V3-M3 request/resultat- og reload-eksempler med generert renderbevis.

V3-M1 verifisert: 59 parser-/kontrakttester består, inkludert request/resultat,
modus, tillatelser, korrelasjon, stegorden og obligatorisk Datagram-variant.

V3-M2: VP08 og avledet MessageSet er implementert. 22 verktøytester består;
numerisk stegorden, korrelasjon og pilens kontraktgrunnlag kontrolleres.
Sekvenssyntaks følger [Mermaid sequence](https://mermaid.js.org/syntax/sequenceDiagram.html);
rendererens grafplassering påvirker ikke meldingsretning eller rekkefølge.
