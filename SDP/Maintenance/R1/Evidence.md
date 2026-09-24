# R1-verifikasjon

## R1-M2 — maler og prosjektområde

- 106 filplasseringer behandlet etter flyttematrise; sju nummererte rotmaler
  deduplisert etter innholdssammenligning uten avsluttende linjeskift.
- Toolkit: 80 unittest-prøver kjørt, 79 bestått. RepositoryValidation feiler bare
  på de samme ni gamle Issue #5-ID-ene som før flytting; baseline og etter-logg
  er byteidentiske. Det er ikke en grønn fullsuite.
- Install-v1 konformanspakke: 17 scenarioer validert. Forventede source-adresser
  er flyttet til Template; destinations og policy er bevart.
- PowerShell finnes ikke lokalt; native Windows-prøver er ikke utført her.
- Gamle bootstrapfiler er arkivert uten innholdsendring. Historiske ledgerlinjer
  beholdes og nye R1-hendelser appendes. Ingen eksisterende installasjon oppdateres.

Kommandoer fra repo-roten:

```sh
python3 Toolkit/scripts/validate_sdp.py
python3 -m unittest discover -s Toolkit/tests -p 'test_*.py'
python3 Toolkit/conformance/install-v1/run_conformance.py --validate-only
```

## R1-M3 — språk, dokumenter og innganger

- 216 ytterligere filplasseringer utført. Hele 322-raders migreringskartet
  kontrolleres av verify_structure.py; malduplikater er konsolidert, ikke nye kopier.
- SDL og SDUI: `go test -race ./...` bestod med Go 1.27.1 og registrert mmdr.
  Alle flyttede Go-kildefiler samt go.mod/go.sum er byteuendrede.
- 574 genererte manifestførte outputs under SDUI/design er hashkontrollert og
  uendret. Ingen viewpoints er håndredigert eller regenerert ved flytting.
- Historiske source-index/release-fingeravtrykk er byteuendrede; gamle
  implementasjonsledgerlinjer er bevart som prefiks. Arkivert bootstrap er uendret.
- Native launcher-prøve i isolert Xvfb bestod: installert sdl-design-symlink peker
  nå til SDL/scripts; hovedside/navigator åpnes, SVG genereres ved behov, ingen
  Go-bygg startes, manglende-verktøy-diagnose og opprydding ved lukking virker.
- Issue #5-bruksstudien bestod sin validator: 35 vurderte repoer, 17 rapporter.
- Lokal lenkekontroll omfatter fil-/katalogmål, ikke ankre eller eksterne URL-er;
  fryste testfixturer og utgått bootstrap er eksplisitt unntatt.
- KanBan-kjeder, metadata og plassering er kontrollert. Kort 001 forblir active
  med avgrenset restarbeid; 010 (innholdskonsolidering) og 011 (gamle ID-avvik)
  ligger i backlog. Ingen språkprofil eller underrepo er opprettet gjennom flyttingen.

Kommandoer:

```sh
python3 SDP/Maintenance/R1/verify_structure.py
python3 SDP/Studies/UsageAnalysis/validate_analysis.py
python3 SDP/Agents/KanBan/examples/verify_lineage.py
go -C SDL/go test -race ./...
go -C SDUI/go test -race ./...
```

Go og native renderer ble valgt med eksplisitte lokale stier i denne kjøringen.
Toolkit-status etter M3 sammenlignes med baseline; ni gamle ID-avvik skal være
hele restfeilsettet. Windows-installer kjøres av CI; ingen lokal Windows-verifikasjon.

Sluttkontroll: 322 filplasseringer, 574 genererte outputs og 1 887 lokale
Markdown-filmål bestod. KanBan: 15 kort, 21 hendelser, 130 lokale lenker.
CI-variant av Toolkit-validatoren med `--base-ref origin/main` ga nøyaktig samme
ni baselineavvik. Rå race-logger ligger i sdl-race.txt og sdui-race.txt.
