# Design-core — port fullført

Python-eksperimentet for design-core 0.5 er erstattet av
[Go-parseren](../../SDL/go/README.md) i G4/G5.
[Språkdefinisjonen](../../SDL/docs/studies/Design-Language-Definition.md) og
[eksemplet](examples/ui-ownership.design) er bevart. 151 fryste parsertilfeller
ligger under SystemDesignLanguage/go/parser/testdata; 170 diagramprojeksjoner
med original kilde ligger under go/viewpoint/testdata.

Ingen aktiv Python-CLI eller fallback. Historisk kode finnes i Git før G5-M4.
Fra SDP-roten: `go -C SystemDesignLanguage/go run ./cmd/sdl check
../../experiments/design_core/examples/ui-ownership.design`.
