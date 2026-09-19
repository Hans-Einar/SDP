# SDUI 0.1 — språk og semantisk profil

**ID:** SDUI-LANG-001. EBNF i [sdui-0.1.ebnf](../grammar/sdui-0.1.ebnf) beskriver
syntaks. Dette dokumentet definerer leksikalske regler og lokal profilvalidering.
Språket er et konkret prototypevalg, ikke en del av vedtatt SDL.

## 1. Dokument og kilde

Et dokument begynner med eksakt `sdui 0.1;`, så null/flere `ref:`-deklarasjoner,
en/flere UI-definisjoner og til slutt null/flere `setHandle`-oppkoblinger.
Rekkefølgen er bevisst; ingen moduldeklarasjoner etter definisjonene.
Parseren leser hele dokumentet til EOF og avviser resttekst.

Kilden er UTF-8. Identifikatorer er ASCII `[A-Za-z_][A-Za-z0-9_]*` og er
case-sensitive. `sdui`, `ref`, `true`, `false`, `null`, `setHandle` er reserverte.
Både enkle og doble anførselstegn støttes. `#` starter kommentar utenfor strenger.
Strenger tillater `\n`, `\r`, `\t`, `\\`, `\"`, `\'` og `\uXXXX` for en Unicode-
skalar uten NUL/surrogater. Supplerende tegn kan skrives direkte som UTF-8.
Rå kontrolltegn U+0000–001F og U+007F er ikke tillatt i strenger.

Tall har JSON-lignende desimalsyntaks, uten innledende pluss eller ledende nuller,
og tolkes som endelig binary64. Det er ingen heltallseksakthetsgaranti over 2^53.
Headeren må staves `0.1`, ikke en numerisk ekvivalent som `1e-1`.
Whitespace er mellomrom, tab, CR og LF. Diagnoselinjer økes ved LF (også i CRLF),
kolonner teller Unicode-skalarer fra 1; tab teller én kolonne.

## 2. Bokser og identitet

```text
Page = {
  main = [title="Main", axis="row",
    left = [weight=1, {ok=button("OK")}],
    right = [weight=2, {help=text("Explanation"); field=input("Value:")}]
  ]
};
```

Hver definisjon har nøyaktig én navngitt rotboks. Alle navngitte bokser og widgets
har unike navn innenfor definisjonen. To definisjoner kan gjenbruke widgetnavn.
En anonym gruppe `[...]` er tillatt, men eksponerer ikke noe eksternt handle.
`handle=` støttes ikke: tilordningsnavnet er identiteten. Properties står før innhold.

En boks inneholder enten underbokser eller én widgetblokk, aldri begge. Tomme
bokser og tomme widgetblokker er ikke gyldige i profilen. Widgetblokken bruker
komma for samme rad og semikolon for neste rad. Siste semikolon før `}` er valgfritt;
tomme rader og trailing comma avvises. Semikolon skifter ikke rad mellom bokser:
bruk en gruppe med eksplisitt akse, slik hovedeksemplet viser.

| Boksproperty | Lokal kontroll | Foreslått betydning ved senere layout |
| --- | --- | --- |
| title | string | Synlig tittel på boksen; tom streng utelater titteltekst |
| axis | "row" / "column", bare for underbokser | Fordeling av underbokser; default column |
| weight | endelig tall > 0 | Relativ andel av tilgjengelig hovedakse; default 1 |
| min / max | endelige tall ≥ 0, min ≤ max | Begrensning på boksens hovedaksestørrelse |
| gap / padding | endelige tall ≥ 0 | Mellomrom / innvendig margin i logiske px; foreslått default 8 |

Layoututførelsen er ikke implementert. Planlagt algoritme fordeler tilgjengelig
plass etter vekt, håndterer grenser eksplisitt og gir diagnose ved uforenlige
minimumskrav. Dette er ikke dagens BoxUI `grow` (minimum først, deretter restplass).
Rottittel/padding, tekstmåling og overflyt må inngå i layoutakseptansen senere.

## 3. Widgets

Widgets har navnet `id = kind(...)`. Første argument kan være posisjonelt og
oversettes semantisk som angitt nedenfor. Deretter bare navngitte argumenter.
Ingen dupliserte argumenter; det posisjonelle og tilsvarende navngitte kan ikke kombineres.

| Kind | Første argument | Andre tillatte egenskaper | Krav |
| --- | --- | --- | --- |
| text | text: string | ingen | text |
| button | label: string | callback: referanse | label |
| input | text: string | value: string, callback: referanse | text |
| svg | source: referanse | label: string | source |

Input `text` er forklarende etikett, ikke feltverdien. `value` er en deklarert
startverdi; en senere runtime-kontrakt må styre initialisering mot domeneeide verdier.
En button/input uten callback er syntaktisk gyldig, men skal være unbound/deaktivert
som handling i en runtime. Bool/null/tall kan representeres i AST; profilens
widgetegenskaper godtar foreløpig bare typene i tabellen.

`svg` erklærer en ressursprodusent. Parseren kjører den ikke. Normal SVG-validering,
ressursbudsjett og inert presentasjon hører til en framtidig utførelsesgrense.
Ukjente widgettyper/egenskaper avvises ved profilvalidering. `--syntax-only` bevarer
slike navn, men merker ikke resultatet som godkjent profil eller kjørbart UI.

## 4. SDL-referanser og oppkobling

```text
ref: sdlFile "some_SDL_file.sdl";
# widgetargument: callback=sdlFile.input1_sdl.@callback
# etter UI-definisjonene:
sdlFile.input1_sdl.setHandle(BoxUIDefinition.input1_boxui);
```

En referanse har formen `module.object.@member`. Aliaset må være deklarert.
Objekt/medlem kontrolleres først av en framtidig SDL-adapter. Filbanen er en
ikke-tom streng uten kontrolltegn; parseren undersøker verken filsystem, URL,
filtype eller ekstern symbolsikt. Det er ikke Markdown-import.

`setHandle` er en deklarativ Connection-node, ikke et generelt metodekall.
Målet må være en widget i en eksisterende definisjon, ikke en boks. Profil 0.1
tillater høyst én kobling per SDL-objekt og per widget i dokumentet. Callbacks
kan dele mottaker; en callback krever ikke setHandle hvis den ikke trenger å
oppdatere widgeten. Oppslag bruker eksakt definisjonsnavn, også store/små bokstaver.

Uttrykk som `input1_sdl.value(...)`, kontrollflyt og vilkårlig funksjonskjøring
hører til SDL/runtime, ikke til dette språket. Callback-signaturer, utførelse og
rettigheter avgjøres ikke av at parseren aksepterer en referanse.

## 5. Markdown og kompatibilitet

En framtidig vertsadapter kan trekke ut et `sdui`-gjerde og sende kroppen til
parseren med kildekartlegging. Denne CLI-en leser bare selvstendig SDUI, ikke
hele Markdown-filer. Ingen eksisterende `boxui 0.1`-JSON aksepteres som SDUI uten
eksplisitt adapter. Concept1s `group/axis/weight/box-id` skal få en separat
normalisering med dokumentert bevaring; React-komponentinnhold oversettes ikke automatisk.
