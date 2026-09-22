# Go-generering — G5-profil

Generatoren validerer hele kilden før den skriver Go. SDUI 0.2 genereres som
kildeposisjonert Document og valgt normalisert Root; SDL action-core 0.1 som
Program med records, actions og Go-symbolreferanser. Typede Go-litteraler
brukes, ikke innebygd kildekode som parses på nytt ved oppstart. Hvert konstruktør-
kall lager uavhengige modeller. Originale posisjoner og bindingsreferanser bevares.

Generert kode er modell-/koblingsdata, ikke en ny runtime. Verten registrerer
håndskrevne Go-funksjoner og eksplisitte bridge.Plan-feltkoblinger. Disse følger
samme validering som filbasert oppstart. Generering finner ikke på Go-domene-
funksjoner, lagring, channel-transport eller en oversettelse av design-core til
kjørbar atferd. Feil profil eller ukjent root avvises før publisering.

`sdl-gen -actions actions.sdl -ui page.sdui -entry page -package model -output DIR`
lager actions_gen.go, ui_gen.go og manifest.json med kildehasher/generatorversjon.
Publisering gjenbruker manifeststyrt atomisk dokumentpublisering; egne filer i
katalogen bevares, redigerte genererte filer avvises. Generert kode bygger mot
de samme Go-modulversjonene som generatoren. Go-kilde skal bygges normalt;
hot reload av språkmodeller og Go-rebuild/restart er forskjellige mekanismer.
