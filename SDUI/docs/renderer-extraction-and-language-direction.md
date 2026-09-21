# SDUI — worktrees, uttrekk og ny språkretning

**Historisk worktree-kart:** Kontroller og råd nedenfor er datert 19.–20. september.
Fra 21. september er Go/Fyne og selvstendig SVG valgt; se
[målarkitekturen](target-architecture.md). Ingen kilde-/worktree-status er
kontrollert på nytt i dokumentasjonsoppdateringen. FOX/C-ABI er ikke aktive leveranser.

Kartlagt 2026-09-19, korrigert 2026-09-20. Status: arbeidsoversikt og kildeinventar.
Dette notatet bevarer eierens nye innspill etter opprydding i chat-/worktree-
organiseringen. Dokumentendringen implementerer ikke ny kode. Eieren har avvist krav om gammel SDUI-/BoxUI-kompatibilitet.
Versjonsnumre, nye biblioteknavn og API-navn nedenfor er forslag der annet ikke er sagt.

## 1. Eierens nye innspill

- SDUI utvikles som et eget språk for UI-prototyper, med XFMD som første vert.
- BoxUI skal skilles fra Mermaid-rendereren. Eieren foreslår en egen
  `sdui-fox-renderer` og spør hvordan lagdelingen bør være.
- Rendererløsningen skal bruke SDUI-parserimplementasjonen, ikke få en kopi av
  parseren. `libsdui` er foreslått biblioteksnavn.
- SDUI og SDL trenger runtime-samspill over en definert ABI, og verten må kunne
  koble SDUI-widgetinstanser til native FOX-widgets med callbacks og oppdateringer.
- Fritekst på en SDUI-side skal være Markdown. En frame skal kunne inneholde
  full Markdown, inkludert Mermaid-diagrammer; dette er ønsket dekning, ikke
  allerede implementert støtte.
- `[]` skal betegne en generell frame. BoxUI blir et valgt uttrykk/implementasjon
  av en frame. `<>` foreslås for widgetinnhold, og `{}` reserveres for layout.
- `*box` foreslås som valg av alternativ komponentimplementasjon; `*b` som
  forkortelse når den er entydig. Dette betyr ikke repetisjon eller multiplikasjon.

Eierens eksempel, bevart som skisse (ikke gyldig SDUI 0.1):

```text
[heading='boxui heading', [<"# gyldig markdown"; button('OK') > ]*box, [ ]*b ; [ ]*b ]*b
```

## 2. Faktisk Git-status

Alle opplysningene er et kontrollert øyeblikksbilde; nye sesjoner må kontrollere
status før endringer. Kartleggingen brukte worktree-listene, status, logg,
merge-base, differ og konkret kildekode. Ingen branch er byttet eller slettet.

SDP-worktreet `/home/warloc/git/SDP-vNow` var rent på `sdp-vNow`, `9ad4324`, før
dette notatet. SDL/MVP1 er bevart i `e051fe2`, SDUI i `3becdb1`.

### Mermaid-rendererens worktrees

| Absolutt katalog | Branch / HEAD | Innhold og status |
| --- | --- | --- |
| `/home/warloc/git/mermaid-rs-renderer` | `master`, `afab5e9` | Ren. Mermaid-utvidelser, måling/ruting og diagramarbeid. Ingen BoxUI-modul på denne branchen. |
| `/home/warloc/git/mermaid-rs-renderer-boxui` | `phase/boxui-045-core`, `4bfd179` | Ren. Felles mandat/design pluss en tidlig BoxUI-kjerne: modell, JSON-parser, målt layout, SVG og seks dokumenterte kjernetester. |
| `/home/warloc/git/mermaid-rs-renderer-boxui-implementation` | `phase/boxui-046-implementation`, `61a85b6` | Ren ved kontrollen. R1–R4-implementasjonen som XFMD faktisk er pinnet til: strict kilde-/prepare-parsing, typed modell, snapshots, geometri, SVG, kontrollkart, child-SVG-komposisjon og grenser. Ingen FOX eller SDL-runtime. |
| `/tmp/xfmd-mrr-measurements` | detached, `3726ccb` | 23 staged filer, 784 tilføyde linjer: forhåndsmålte tekster og kooperativt tidsbudsjett/checkpoints i diagramlayout. Separat fra BoxUI. Bevares. |

De to BoxUI-grenene har felles base `7076cac`; core har én egen commit og
implementation fire egne commits. **Implementation er ikke en fortsettelse oppå
core.** Ikke bland parser-/layoutfilene fra dem. Produsentens overlevering og
XFMDs integrasjonsdokument sier eksplisitt at implementation er integrasjonskilden.
`feature/boxui-extension` peker på designbasen `7076cac`, uten eget registrert worktree.

Målearbeid er senere representert i Mermaid-historikken ved `38b6018` og utvidet
ved `3aad3fb`. Den staged patchen i `/tmp` har ikke samme stabile patch-ID som
`38b6018`; det er ikke bevist at hele den midlertidige tilstanden er redundant.
Ingen opprydding av den er utført eller anbefalt uten en egen differansekontroll.

`gh pr list --state open` for Hans-Einar/mermaid-rs-renderer returnerte ingen åpne
PR-er ved kontrollen. Dette sier ikke at grenene er integrert. Lokal Git viser
at BoxUI ikke finnes i `master`.

### XFMDs worktrees

| Absolutt katalog | Branch / HEAD | Innhold og status |
| --- | --- | --- |
| `/home/warloc/git/xfmd` | `main`, `c245fd9` | Ren, uten BoxUI-integrasjonen. Mermaid-pinnen er `589517a`. |
| `/home/warloc/git/xfmd-boxui` | `sprint/003/phase/048-boxui-verification`, `a7495b8` | Ren. Native input/knapper, lokal syntetisk aktivitet, sesjon/ledger, statisk PDF og BoxUI-fence. Mermaid-pinnen er hele `61a85b6`. PR #37 er åpen og ikke merget. |
| `/tmp/xfmd-p31-lifecycle` | `phase/p31-mermaid-svg`, `35e70b7` | Ren. Eldre SVG/Cairo-integrasjon og Pango-opprydding; HEAD er allerede stamfar til `main` og BoxUI-grenen. Ingen nye lokale endringer funnet. |

PR: https://github.com/Hans-Einar/xfmd/pull/37

Det er derfor ikke nødvendig å reversere BoxUI ut av Mermaid `master` nå.
Uttrekket skal ta vare på implementation-koden og flytte XFMDs avhengighet til
en selvstendig produsent. Historiske branches/worktrees kan beholdes som proveniens.

## 3. Konkrete kodegrenser for uttrekket

Relative stier i denne tabellen gjelder navngitt worktree, ikke SDP.

| Nåværende eier | Kilder | Anbefalt disposisjon |
| --- | --- | --- |
| Mermaid implementation | `src/boxui/{model,parse,validate,frame,layout,svg,embedded}.rs`, `mod.rs` | Gjenbruk egnede deler i ny FOX-uavhengig SDUI-kjerne. Gammel JSON-parser og wiremodell skal ikke bli en kompatibilitetsvei. |
| Mermaid implementation | `tests/boxui.rs`, `examples/boxui_*`, `SDP/06--Container-Design/contracts`, fixtures og bevis | Bevar proveniens/lisens og port relevante tester; gamle kontrakter er historisk sammenligningsgrunnlag. |
| Mermaid implementation | `src/lib.rs`, Cargo-manifest/-lock | BoxUI er lagt til via `pub mod boxui` og `roxmltree`. Selvstendig crate trenger sine faktiske avhengigheter; den trenger ikke hele Mermaid-craten for kjernen. |
| XFMD BoxUI | `src/interpreter/mermaid/rust/src/boxui.rs` | Dagens tynne parseadapter kaller `mermaid_rs_renderer::boxui`. Flytt avhengigheten til ny produsent. |
| XFMD BoxUI | `src/renderer/diagram/rust/src/boxui.rs` | Dagens prepareadapter kaller samme produsent og låner tekstmåling/cancellation. Skill UI-adapteren fra Mermaid-adapteren. |
| XFMD BoxUI | `src/application/composition/mermaid/src/boxui.rs` | Eier C-ABI parse/prepare/free, bufferfrigjøring og panic-grense. Er ikke implementert i rendererforken. Bevar funksjonen under overgangen. |
| XFMD BoxUI | `src/contracts/boxui`, `src/interpreter/boxui`, `src/renderer/boxui` | Modell, fence-ekstraksjon, vertsadaptere og plassering. Migrer koordinert med produsenten. |
| XFMD BoxUI | `src/application/boxui/BoxUiSession.*`, `BoxUiPreparation.*` | Sesjonsidentitet, snapshots, ledger og forberedelse. Vurder generisk runtime-uttrekk separat; kopier ikke all applikasjonskoordinering inn i layoutmotoren. |
| XFMD BoxUI | `src/application/boxui/SyntheticActivity.*` | Lokal simulering/testadapter. Skal fortsatt merkes som simulering, ikke SDL-runtime. |
| XFMD BoxUI | `src/application/adapters/FoxBoxUiOverlay.*`, `FoxBoxUiInput.*` | Gjenbruksgrunnlag for FOX-backend. I dag koblet til FoxRenderHost og BoxUiSession; krever en eksplisitt vertsgrense før det blir et selvstendig bibliotek. |

Rust-kjernens moduler bruker egne `super`-moduler og standard-/støttebiblioteker;
parent-prepare parser ikke Mermaid. Verten leverer allerede forberedte child-SVG-er.
Eksempler og integrasjonstester kan fortsatt bruke Mermaid som testkonsument.
Vanlig Mermaid og standard treemap skal beholde sine eiere og regresjoner.

## 4. Gjeldende retning etter eierens korrigering 2026-09-20

Tidligere anbefaling om et selvstendig bibliotek med bevart BoxUI 0.1/BX-HOST
som første leveranse er trukket tilbake. Vi trenger ingen kompatibilitetsvei;
gjenbruk egnet kode direkte i ny SDUI-modell, port eksempler og fjern erstattede
kjøreveier. Mermaid master trenger ingen SDUI-endring. Historiske worktrees
skal ikke forveksles med aktive produktavhengigheter.

[Implementasjonsplanen](implementation-plan.md) eier phases/milestones og opprydding.
[Målarkitekturen](target-architecture.md) eier nye porter og bibliotekansvar.
[Layoutforslaget](layout-language-proposal.md) konkretiserer frames, widgets,
Markdown, modifiere og symbolene eieren skisserte. Disse er planlagt ny profil,
ikke støtte i dagens parser. Ingen gamle ABI-/wireformer fryses av dette notatet.

## 5. SDP-piloten

Den nummererte strukturen ble pilotert i BoxUI-prosjektet, og finnes i begge
Mermaid BoxUI-worktreenes SDP-katalog. Mandat: 01--Mandate/01-01--Mandate.md.
Struktur: 00--Project/00-02--SDP-Structure.md. Manifest: sdp-project.json.
Den tidligere intensjonen om BoxUI som Mermaid-utvidelse er nå erstattet av
selvstendig SDUI-retning; pilotens historiske mandat er ikke en ny begrensning.

## 8. Verifikasjonsgrense

Dette notatet bygger på faktisk kilde-/Git-kontroll. Den foregående gjennomgangen
i samme sesjon kjørte SDUI 15/15, design-core 31/31 og MVP1-inventarkontroll;
disse testene etablerer ikke det nye språket eller runtime/rendering.
I denne kartleggingen ble følgende kjørt fra renderer-implementation på `61a85b6`:

```sh
cargo test --offline --locked --no-default-features --test boxui
```

**18/18 tester bestod.** Bygget ga ni dead-code-advarsler fra ordinære
diagramrutingmoduler. Testene omfatter blant annet deterministisk frame,
kontrollgeometri, child-komposisjon, ressurser og kansellering. De beviser ikke
at en ny selvstendig crate finnes eller fungerer. Worktreet forble Git-rent.

Ny uttrekksleveranse må kjøre sine egne regresjoner og identifisere konsumentpin.
Ingen ny GUI-/PDF-verifikasjon eller installasjon inngår i kartleggingen.
