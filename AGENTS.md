# Samarbeid og Git-sporbarhet

Eierbeslutning 2026-09-22 for SDL/SDUI-utviklingen:

- Bruk én branch per fase og én egen commit per fullført milepæl.
- Opprett neste fasebranch fra siste commit på forrige fasebranch. Bevar
  fasebranchene slik at leveransene kan gjennomgås hver for seg.
- Målbranch for det samlede arbeidet er `sdp-vNow`, opprinnelig på
  `9ad432407004080dd7f4f0ab06d107523f4316fd`. Ikke utvikle direkte på denne.
- Bruk fase-/milepæl-ID i commitmeldingen, med konkret leveranse og relevant
  verifikasjon i committeksten. Ikke merk uferdig arbeid som levert.
- Kontroller Git-status før branching og staging. Stage bare arbeidet som
  tilhører milepælen; bevar annet lokalt arbeid. Unngå genererte cacher.
- Oppdater implementasjonsplan og bevis ved milepælen. Genererte SDL-viewpoints
  skal komme fra SDL-verktøyet og validerte modellfakta.
- Opprett samlet PR mot `sdp-vNow` når det avtalte arbeidet er klart. Opprettelse
  av lokale fasebrancher og milepælcommits er autorisert; dette er ikke en
  instruks om å merge eller publisere hver mellomleveranse.

Branchstakken og håndteringen av tidligere ucommittet arbeid er beskrevet i
[utviklingshistorikken](docs/Development-Branch-Stack.md). Les også relevante
underkatalogers AGENTS.md før endringer.

## Lokale fillenker i svar

Bruk rene absolutte `file:///`-URL-er, én per linje, uten Markdown-innpakking,
linjenummersuffiks eller terminalkontrollkoder. Utvid hjemmekatalogen og
prosentkod mellomrom og reserverte URL-tegn.
