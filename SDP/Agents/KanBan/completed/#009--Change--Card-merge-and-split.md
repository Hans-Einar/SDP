# K3-M1: sporbar sammenslåing og splitting av kort

| Felt | Verdi |
| --- | --- |
| id | KB-SDP-009 |
| project | SDP |
| type | Change |
| created | 2026-09-23T23:57:37Z |
| source | owner-conversation-2026-09-24; XFMD conversation board |
| owner | Codex |

## Behov og grunnlag

Eieren ønsker at backlogkort kan slås sammen før aktivering eller splittes i
flere arbeidskort, med sporbarhet i ledgeren. Ny kunnskap må kunne presisere og
erstatte eldre forslag uten at informasjon eller restarbeid forsvinner.

XFMDs lokale KanBan-regler foreslår nytt samlekort, kildetabell, toveis lenker,
superseded for fullt erstattede kilder og reviewed for kilder med restarbeid.
Splitting nevnes også. Dette er nyttig gjenbruk av arbeidsmåte, uten at XFMD tar
inn SDP-prosessen. XFMDs filer var ucommittede ved lesing; ingen endres her.

Inspisert kilde: `/home/warloc/git/xfmd-sdl-navigation/Agents/KanBan/README.md`.
Registrering: 2026-09-23T23:57:37Z. SHA-256: `eb7488d631dc8c9c934a71da15fe3e3f52239989998596648e64154a8ff2c7c8`.
Det er lokal proveniens, ikke påstand om en publisert XFMD-kontrakt.

## Implementasjonsplan K3

Én fasebranch `sdp/phase-k3-card-lineage` fra K2 (`321e193`), én milepæl K3-M1:

1. Beskriv full/delvis sammenslåing og splitting, kilde-/målroller, bevaring av
   åpne spørsmål og aktivering med avgrenset, autorisert omfang.
2. Utvid ledger-payload med versjonert, typet opphav uten å skrive om gamle
   hendelser; bruk eksisterende created/moved/reviewed og stabile kort-ID-er.
3. Oppdater kortmal og grafidé; lag isolerte, eksplisitt fiktive eksempler og
   kontroller positive/negative tilfeller, replay og eksisterende tavler.

Dette leverer arbeidsmåte og datakontrakt, ikke automatisk filflytting, graf eller
sdptool-implementasjon. Ingen eksisterende produktforslag skal slås sammen uten
faglig vurdering. Kildene i XFMD er bare lest.

## Akseptanse og bevis

**K3-M1 levert 2026-09-24.** Opphavskontrakt, payload-schema 0.2, oppdatert
kortmal og grafpresisering er levert. Eksisterende payload-schema 0.1 og historiske
hendelser beholdes uendret. Nye kort opprettes ved merge/split; delvis overføring
beholder eksplisitt restarbeid. Ingen produksjonskort er slått sammen eller splittet.

`python3 SDP/Agents/KanBan/examples/verify_lineage.py` bestod: fiktive fulle og
delvise merge/split-forløp, historisk payload 0.1 og 15 negative tilfeller.
Prøvene avviser blant annet manglende deltakere/lenker, sprikende definisjoner,
gjenbrukte mål-ID-er, overlapp og ugyldig håndtering av restarbeid.

En separat engangskontroll av produksjonstavlene bestod: tre tavler, 21 statusmapper,
13 kort, fire Ref og 112 lokale Markdown-lenker. Schema/replay og filplassering
bestod også etter avslutning; samlet ledgerantall er 17 hendelser.
Gamle ledgerlinjer er kontrollert som bytebevarte prefikser mot foreldercommiten;
SDL/SDUI-ledgerne og payload 0.1 er uendret. `git diff --check` bestod.

Dette er en manuell prosess og testet datakontrakt. Ingen generell flyttekommando,
KanBan-graf eller kryssrepo-transaksjonsmotor er levert. XFMD-filer er bare lest;
forslaget der er brukt som lokalt, ucommittet grunnlag, ikke endret eller publisert.
