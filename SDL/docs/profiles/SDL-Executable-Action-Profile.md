# SDL — kjørbar handlingsprofil

G4-M2, 2026-09-22. `action-core 0.1` er en eksplisitt, avgrenset kjøreprofil
ved siden av strukturell `design-core 0.5`. Strukturelle Functionality-, Mode-
og Channel-fakta får ingen skjult utføringsbetydning. Profilene deler SDL-lexer,
identifikatorregler og kildeposisjoner; hver har sin lukkede grammatikk.

```text
language action-core version 0.1.
action Echo.
record EchoInput.
record EchoOutput.
Echo invokes GoEcho.
Echo returns EchoOutput.
Echo takes EchoInput.
EchoInput field Value as text.
EchoOutput field Value as text.
```

En Action må ha nøyaktig én `takes`, `returns` og `invokes`. De to første
refererer deklarerte Records; `invokes` navngir én eksplisitt registrert Go-
funksjon. Den er en symbolreferanse, aldri Go-kildetekst, shell eller dynamisk
kodeimport. Samme funksjon kan gjenbrukes når begge recordsignaturer er like.
Hver Record har 1–32 navngitte felt; feltnavn er lokale for recorden. Felt er
obligatoriske og typede `text`, `integer` (signert 64-bit) eller `boolean`.
Ekstra/manglende felt og feil scalar-type avvises. Ingen null/default eller
implisitt konvertering. Tekst er gyldig UTF-8, uten NUL, maksimalt 32 KiB.

Deklarasjoner sorteres på navn, deretter alle fakta på full ASCII-setning.
Ingen kommentarer i kanonisk kilde; dokumentasjon ligger rundt modellen.
Maksimalt 128 Actions, 128 Records og 32 felt per Record støttes først.
Kilde-/tokengrenser følger den felles SDL-frontenden. Ukjente utsagn avvises.

Go-registeret angir input-/outputsignatur og en funksjon
`func(context.Context, Record) (Record, error)`. Hele modellen og alle nødvendige
registreringer valideres før runtime opprettes. Kallet validerer record før
handleren og resultat før publisering. Recordverdier kopieres ved grensen.
Feil er observerbare resultater; runtime gjetter ikke domeneregler.

Domene-state eies av den registrerte Go-implementasjonen. UI-draft, UI-revisjon,
command-identitet og akseptert domeneresultat holdes atskilt. Denne profilen er
ikke et distribuert meldingssystem eller et garantilag for eksterne sideeffekter.
En avgrenset EditAptCell-prøve følger G4-M3: stabil celleidentitet, draft, forventet
domenerevisjon, eksplisitt avvisning og korrelert godkjent resultat. Den skal
merkes som simulering og utfører ingen Ponsse-/maskinhandlinger.

Hot reload av modell og håndtering av pågående kall er implementert og prøvd i
[G4-M4](../../go/evidence/G4.md#g4-m4).
Go-funksjonsendringer krever vanlig Go-bygg og prosessrestart.

## SDUI-port — G4-M3

SDUI-refene er symbolske. Verten leverer en eksplisitt alias→SDL-runtime-tabell;
bridge eller parser åpner ikke ref-stien. Kjørbar callback er
`module.Action.@invoke`. `module.Action.setHandle(page.input)` navngir resultat-
mottakeren. Første resultatport skriver et `text`-felt til et input-handle.

En typet Go-bindingsplan sier hvor hvert inputfelt kommer fra: widgetens draft,
eventverdi, typed literal eller navngitt kontekstverdi. Det finnes nøyaktig én
kilde per felt. Tekstinput konverteres eksplisitt i adapteren hvis SDL-feltet er
integer/boolean; runtime gjør ingen implisitt konvertering. Hele bindingssettet
valideres før noen handler installeres. Ukjent modul/member/action, manglende
feltkilde og feil resultatwidget avvises. Linkoversikten har kildeposisjoner til
både SDUI-widgeten og SDL-handlingen.

Domenerevisjon er egen bindingskontekst, aldri UI-verdirevisjonen. Et godkjent
SDL-resultat kan oppdatere denne konteksten og UI-verdien. UI-propertybatchen
valideres fortsatt samlet; en presentasjonsfeil ruller ikke tilbake en allerede
utført Go-domenehandling og må ikke føre til automatisk gjentakelse.
`CurrentInvocation(ctx)` gir Go-funksjonen action, modellrevisjon og sekvens-ID.
Profilen har én ordnet kommandokilde per Engine; flere samtidige avsendere må
koordineres av samme sekvenseier. Ingen distribuert exactly-once-garanti utledes.

[EditAptCell-kildene og håndskrevet simulering](../../go/examples)
viser stabil celleidentitet, separat domenerevisjon og eksplisitt avvist edit.
Den eldre MVP1-scenariofilen er en kravreferanse, ikke innlest kjørbar kode.
