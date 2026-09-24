# SDUI ↔ SDL — implementert runtimegrense

Oppdatert 2026-09-24 etter G3/G4. SDUI 0.2 og SDL action-core 0.1 har separate
Go-runtimer med eksplisitt bridge. Ingen C-ABI, FOX-peker eller automatisk
SDL-loader inngår. Dette dokumentet samler grensen; de detaljerte kontraktene
har ett hjem hver:

| Ansvar | Aktiv kontrakt / implementasjon | Bevis |
| --- | --- | --- |
| UI-session, handles, draft/accepted og propertybatch | [SDUI-runtime](../go/runtime/README.md) | [G3](../go/evidence/G3.md) |
| SDL-handlinger, records og Go-registrering | [action-core 0.1](../../SDL/docs/profiles/SDL-Executable-Action-Profile.md) | [G4](../../SDL/go/evidence/G4.md) |
| Feltkobling og resultatport | [bridge](../../SDL/go/bridge/bind.go), action-core-kontraktens SDUI-port | [G4-M3](../../SDL/go/evidence/G4.md#g4-m3) |
| UI-/SDL-modellreload og Go-restart | [kjørekommandoer](../../SDL/go/README.md), UI-runtime | [G4-M4](../../SDL/go/evidence/G4.md#g4-m4) |
| Genererte modeller mot samme runtime | [Go-generering](go-generation.md) | [G5](../../SDL/go/evidence/G5.md) |

## Oppkobling

SDUI-parseren lagrer `ref`, callbacks og `setHandle` som data. Verten registrerer
modulalias → SDL-runtime, Go-funksjonenes signaturer og en typet `bridge.Plan`.
En ref-sti åpnes ikke av parser eller bridge. Handlingsformen er
`module.Action.@invoke`; `module.Action.setHandle(page.input)` peker på
resultatmottakeren. Navn alene aktiverer ikke domenekjøring.

Bindingsplanen angir nøyaktig én inputkilde per felt: widgetdraft, eventverdi,
typet literal eller navngitt kontekstverdi. Signaturer, alias, action og
resultatwidget valideres før handlers installeres. Den leverte resultatporten
skriver et text-felt til et input-handle; det er ikke en generell widget-API.

## Identitet, eierskap og oppdatering

UI-handles inneholder session, instansbane, generation og widgettype. De er ikke
native pekere. SDUI eier accepted/draft, fokus og UI-revisjoner; den registrerte
Go-funksjonen eier domenestate og domenerevisjon. En ekstern verdioppdatering
skal ikke overskrive dirty draft stille. UI-batcher valideres samlet.

Brukerhendelser går gjennom kontroll av handle/revisjon, binding og enabled.
Enter committer draft; Escape gjenoppretter accepted. Programmatisk oppdatering
utløser ikke samme bruker-callback igjen. Fyne-verten publiserer på UI-tråden.
SDL validerer input/output-records og korrelerer kall; avviste eller dupliserte
kommandoer gir observerbart utfall. En allerede utført Go-domenehandling rulles
ikke tilbake ved senere UI-feil, og feilen utløser ikke automatisk replay.

Kompatibel modellreload bevarer state etter pakkekontraktene; slettede eller
inkompatible widgets får ugyldige gamle handles. Ugyldig kandidat beholder siste
gyldige modell. SDL-reload validerer registreringer og bevarer Go-eid domenestate.
Endret Go-kode krever bygg/restart og bevarer ikke automatisk minnestate.

## Grenser og tidligere forslag

Design-core-fakta blir ikke kjørbare gjennom denne koblingen. EditAptCell er
en eksplisitt simulering; dette er ingen maskin-/Ponsse-runtime, distribuert
transport eller exactly-once-garanti. Generell asynkron widget-API, FOX-/C-ABI
og kjørbar SVG-produsent er ikke levert av G3/G4. SVG-widgeten er plassholder.

Det opprinnelige SDUI-RUNTIME-001-forslaget finnes i Git før R2-M2. Det er
erstattet som aktiv kontrakt av pakkekontraktene ovenfor; eldre foreslåtte
payloads og metoder skal ikke leses som et ekstra støttet API.
