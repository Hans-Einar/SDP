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
  instruks om å merge. Eierens tillegg 2026-09-22 autoriserer push til origin
  etter hver fullført fase i samme sesjon.

Branchstakken og håndteringen av tidligere ucommittet arbeid er beskrevet i
[utviklingshistorikken](SDP/Development-Branch-Stack.md). Les også relevante
underkatalogers AGENTS.md før endringer.

## Lokale fillenker i svar

Bruk rene absolutte `file:///`-URL-er, én per linje, uten Markdown-innpakking,
linjenummersuffiks eller terminalkontrollkoder. Utvid hjemmekatalogen og
prosentkod mellomrom og reserverte URL-tegn.

## KanBan for ideer og omfangsendringer

Les [KanBan-arbeidsmåten](SDP/Agents/KanBan/README.md) og relevant prosjekttavle
før nytt arbeid. SDP, SDL og SDUI har egne tavler registrert der. Nye ideer og
funn utenfor avtalt omfang registreres i et hovedkort, med Ref-kort ved behov,
før arbeidet skifter retning. Registrering er ikke autorisasjon til implementasjon.
Ved behandling/flytting oppdateres kort, append-only KanBan-ledger og lenker
samlet. Gjennomgå backlog/onHold ved milepælslutt og før neste fase.
KanBan dokumenterer behandlingen av forslag; faktisk implementasjon og
verifikasjon tilhører fortsatt Traceability. Ikke migrer gamle prosessområder
eller vedta nye språkregler bare fordi de er beskrevet i et backlogkort.
