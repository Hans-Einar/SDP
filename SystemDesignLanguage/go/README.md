# SDL i Go

G4-M1 leverer parser, kildeposisjonert AST, symbol-/typekontroll, data-/wire- og
Channel-/scenariovalidering og kanonisk form for **design-core 0.5**.
Ingen struktursetning utføres. Checkpoint-kandidater og hele MVP1-korpuset er
ikke del av profilen. G4-M2 leverer også den eksplisitte kjøreprofilen
[action-core 0.1](../../docs/SDL-Executable-Action-Profile.md), parser og runtime
med typede records og registrerte Go-funksjoner.

Modul: `github.com/Hans-Einar/SDP/SystemDesignLanguage/go`, Go 1.26 som felles
baseline. Verifisert med Go 1.27.1. Strukturkjernen bruker bare standardbiblioteket.

Fra denne katalogen:

```sh
go test -race ./...
go run ./cmd/sdl check ../../SDUI/design/architecture.design
go run ./cmd/sdl ast ../../SDUI/design/architecture.design
go run ./cmd/sdl format ../../SDUI/design/architecture.design
```

`parser.Data` beholder Python-profilens navngitte JSON-AST-form, inkludert start-
og sluttposisjoner. 151 porttilfeller og full SDUI-modell har sammenligningsbevis.
[Faktiske bevis](evidence/G4.md), [felles plan](../../SDUI/docs/implementation-plan.md),
[språkdefinisjon](../../docs/Design-Language-Definition.md).

Python-frontenden og viewpoint-generatoren er erstattet i G5-M4. Fryste
portfixturer beholdes som historisk orakel; det finnes ingen fallback.

Kjøreprofilkontroll: `go run ./cmd/sdl action-check examples/echo.sdl`.
`runtime.New` krever eksplisitt signaturregistrering; ingen kilde kjøres som Go.

G4-M3 har en typet SDUI-port i `bridge/`. Enkel Echo og en eksplisitt simulert
EditAptCell går gjennom begge runtimene. `go run ./cmd/sdl-simulate` skriver
korrelert hendelsesspor; `go run -tags desktop ./cmd/sdl-demo` viser den native
Fyne-prototypen. `examples/simulation` er separat håndskrevet Go-domenelogikk.

G4-M4: demoen følger begge kildefiler; ugyldig kilde beholder siste gyldige
modell og feilstatus. SDL-modelreload beholder Go-eid domenestate og avviser
eldre hendelser. Endret Go-kode krever bygg/restart:

```sh
go run ./cmd/sdl-dev -root . -package ./cmd/sdl-demo -tags desktop
```

Utviklingsverten beholder kjørende prosess ved byggfeil. Vellykket bygg starter
ny prosess; vedvarende domenetilstand over prosessrestart krever egen lagring.

G6-M1 porter alle 11 strukturelle viewpoints med uendret kildegrunnlag:

```sh
go run ./cmd/sdl viewpoints ../../SDUI/design/architecture.design --output /tmp/sdl-navigation --project sdui-design
go run ./cmd/sdl viewpoints ../../SDUI/design/architecture.design --output /tmp/sdl-static --format static --monolithic
```

Valgfritt `--renderer /absolutt/sti/til/mmdr` lager SVG i statisk eksport.
`--viewpoint VP02,VP08` avgrenser eksporten. Navigator er standard og renderer
ingen detaljer. Dens handlingslenker krever den registrerte XFMD-leseradapteren.

Et utvalg ved behov:

```sh
go run ./cmd/sdl view ../../SDUI/design/architecture.design --uri 'sdl-view://sdui-design/VP02?focus=SduiFrontend&relations=contains&direction=out&depth=1&level=A2' --output /tmp/sdl-selected
```

`--renderer` er valgfri også her. Hele pakken publiseres før `entry.md` tilbys.
CLI-en løser kildefilen eksplisitt; en vertsadapter må registrere prosjekt-ID.

Valgfri dokumenttjeneste (Linux):

```sh
go run ./cmd/sdl-viewsd -source ../../SDUI/design/architecture.design -project sdui-design -xfmd /absolute/xfmd -renderer /absolute/mmdr
go run ./cmd/sdl-view-request -socket /private/sdl/views.sock -uri 'sdl-view://sdui-design/VP02?diagram=VP02-roots' -window design-one -client terminal -sequence 1 -open
```

Tjenesten skriver valgt socket ved start. Bruk ny sekvens per klient/vindu/panel.
Uten `-open` returneres en lesbar pakke og lease. `-release TOKEN` frigjør den;
`-sweep` fjerner bare frigjorte pakker. Leases overlever daemonkrasj; ved krasjet
leser kreves eksplisitt release. XFMD fase 050 frigjør automatisk ved bytte/lukking.


G5 — generer og bygg modeller sammen med håndskrevet Go-domene:

```sh
go run ./cmd/sdl-gen -actions examples/edit-apt-cell.sdl -ui examples/edit-apt-cell.sdui -output examples/generatedmodel -package generatedmodel
go run ./cmd/sdl-compiled
go run -tags desktop ./cmd/sdl-compiled-fyne
go run ./cmd/sdl-document -ui examples/edit-apt-cell.sdui -state examples/accepted-state.json -design ../../SDUI/design/architecture.design -output /tmp/sdl-ui-document
```

Genererte konstruktører åpner ikke kildefiler. -values på sdl-compiled/sdl-simulate
velger simulert hendelsesserie. sdl-document utfører ingen callbacks; -state angir
aksepterte widgetverdier, label, enabled/visible, og manifestet eier bare genererte
filer. Full layout eksporteres som SVG og Markdown med provenance og SDL-navigator.

Klasseprofilen er eksplisitt, uten utledning fra contains/owns:

```sh
go run ./cmd/sdl class-check examples/runtime-classes.sdl
go run ./cmd/sdl class-view examples/runtime-classes.sdl --output /tmp/sdl-classes --renderer /absolute/mmdr
```

[Go-genereringsprofil](../../SDUI/docs/go-generation.md), [G5-bevis](evidence/G5.md),
[class-core](../../docs/SDL-Class-Profile.md), [samlet status](../../docs/checkpoint%231/11-Go-Implementation-and-Navigation.md).
