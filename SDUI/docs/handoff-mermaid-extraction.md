# Oppstart — Mermaid som diagrammotor og kilde for gjenbruk

Oppdatert 2026-09-21: Go/Fyne er valgt; [PLAN-003](implementation-plan.md) gjelder.
Historisk Mermaid-kartlegging nedenfor er fra 19.–20. september, ikke ny Git-kontroll. Dette erstatter overleveringen
som krevde et eget uttrekk med bevart BoxUI 0.1/BX-HOST-kompatibilitet.

## Oppdragets nåværende grense

Mermaid master på afab5e9 har ingen BoxUI-kode. Vi trenger **ingen endring i
Mermaid nå** for å starte ny SDUI-utvikling. Ikke opprett en cleanup-PR eller
legg SDUI-kode inn i master. Bevar vanlig Mermaid, treemap, måling og ruting.

SDUI-sesjonen i `/home/warloc/git/SDP-vNow` eier ny frontend, runtime og layout/
presentasjon. Den kan studere og porte egnede algoritmer/tester fra implementation-worktreet
til Go med dokumentert proveniens. Ingen Rust-avhengighet opprettes automatisk. Det kreves ikke en mellomliggende legacy-crate, gammel
JSON-parser, gammel wire-protokoll eller bakoverkompatibel geometri.

Les gjeldende AGENTS.md der de finnes, og:

- `/home/warloc/git/SDP-vNow/SDUI/docs/implementation-plan.md`
- `/home/warloc/git/SDP-vNow/SDUI/docs/layout-language-proposal.md`
- `/home/warloc/git/SDP-vNow/SDUI/docs/target-architecture.md`
- `/home/warloc/git/SDP-vNow/SDUI/docs/renderer-extraction-and-language-direction.md`

## Gjenbrukskilder og arbeid som skal bevares

- `mermaid-rs-renderer-boxui-implementation`, 61a85b6: implementasjonen XFMD
  bruker. Tekstmåling, SVG-escaping/presentasjon, kontrollkart, budsjett og
  tester er kandidater. Nye fr-/frame-regler erstatter gammel grow-semantikk.
- `mermaid-rs-renderer-boxui`, 4bfd179: parallell tidlig implementasjon fra
  samme base, ikke en nødvendig ekstra merge. Ikke bland disse parserfilene.
- `/tmp/xfmd-mrr-measurements`: 23 staged filer i siste kontroll; ingen
  sletting/nullstilling uten egen kartlegging. Ikke del av SDUI-oppryddingen.
- Den nummererte SDP-piloten med mandat ligger i BoxUI-worktreenes SDP-katalog.
  Den dokumenterer tidligere arbeid, ikke ny myndighet til å bevare legacy.

Gjenbruk dokumenteres med kildecommit og lisens. Gamle tester portes til nye
kontrakter der de fortsatt kontrollerer relevante egenskaper. Identiske gamle
AST-/wire-/SVG-bytes er ikke et krav. Testgeometri endres når semantikken endres.

## Eventuelt senere Mermaid-arbeid

Bare en konkret feil/manglende diagramkapabilitet avdekket av Markdown-provideren
skal gi et separat Mermaid-oppdrag. Den nye SDUI-kjernen skal ikke avhenge av
Mermaid som produksjonsbibliotek. Diagramtjenesten forblir en vertsavhengighet.

Ingen kildekode/worktrees er fjernet av kartleggingssesjonen. Ingen merge/push
eller installasjon er utført. Gamle worktrees er proveniens, ikke aktive språkveier.
