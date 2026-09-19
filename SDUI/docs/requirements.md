# SDUI — krav og sporbarhet

**ID:** SDUI-REQ-001 · Arbeidsprofil 0.1. «Implementert/testet» gjelder bare oppgitt
del, ikke UI-rendering eller SDL-runtime. Testnavn viser inn i tests/test_language.py.

| ID | Krav | Eier / bevis | Status |
| --- | --- | --- | --- |
| SDUI-R01 | Navngitte nestede bokser, egen rot, egenskaper og underbokser | parser.Box; owner_example_preserves_boxes_rows_and_bindings | Implementert/testet struktur |
| SDUI-R02 | Widgets på samme rad og eksplisitt radskille | parser.content; owner_example, static_comments_quotes_and_unicode | Implementert/testet |
| SDUI-R03 | Bevar modulalias, path og module.object.@member | ast.ModuleRef/Reference; owner_example og local_symbol_errors | Implementert/testet; ingen SDL-resolusjon |
| SDUI-R04 | setHandle blir data med målwidget, ikke kodekjøring | ast.Connection, validate; local_symbol_errors, no_arbitrary_code | Implementert/testet lokalt |
| SDUI-R05 | Stabilt navneoppslag per definisjon uten redundant handle | validate; local_symbol_errors, anonymous_group_and_reusable_names | Implementert/testet lokalt |
| SDUI-R06 | Versjonert EBNF, AST og feil med kildeposisjoner | grammar, ast.Span; version_is_exact, utf8_byte_spans_and_crlf | Implementert/testet; ingen formell grammatikkekvivalens |
| SDUI-R07 | Ukjent/motstridende profilsemantikk må ikke slippes gjennom stille | validate; widget_contracts, box_properties_and_layout_constraints | Implementert/testet |
| SDUI-R08 | Avgrenset ressursbruk og ingen eksekvering under parsing | lexer/parser/CLI; bounded_work, truncations, cli_json_exit_codes | Implementert/testet; ikke hard sanntid |
| SDUI-R09 | Standalone parsing og JSON-AST uten GUI-/Mermaid-avhengighet | CLI, standardbibliotek; cli_json_exit_codes | Implementert/testet |
| SDUI-R10 | Samme normaliserte UI-modell fra SDUI og Concept1-layout | Framtidig normalisering/adapter | Planlagt, ikke implementert |
| SDUI-R11 | Synlig ytre boks, tittel, vekter og widgetlayout i XFMD/Markdown | Framtidig UI-kjerne og vertsadapter | Planlagt, ikke implementert |
| SDUI-R12 | SDL-objekt kan motta handle og oppdatere text/value | Foreslått runtime-kontrakt | Planlagt, ikke implementert |

Testnavnene i tabellen er forkortet for lesbarhet; filen inneholder fulle navn.
[Verifikasjonsrapporten](../evidence/verification.md) registrerer faktisk kjøring.
Ingen semantisk validator påstår at en ekstern SDL-fil eller funksjon finnes.
