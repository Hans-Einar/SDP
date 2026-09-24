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
