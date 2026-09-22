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

Hot reload av modell og håndtering av pågående kall presiseres og prøves i G4-M4.
Go-funksjonsendringer krever vanlig Go-bygg og prosessrestart.
