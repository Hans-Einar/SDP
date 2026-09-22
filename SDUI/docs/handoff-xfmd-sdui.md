# FOX/XFMD — tidligere handoff, avløst av Go/Fyne-retning

Oppdatert 2026-09-22. Den gamle FOX-widget-handoffen er avløst.
G6s dokumentnavigasjon er implementert separat i [XFMD PR #38](https://github.com/Hans-Einar/xfmd/pull/38);
se [G6-bevis](../../SystemDesignLanguage/go/evidence/G6.md).
Eieren har valgt selvstendige SDL/SDUI-biblioteker i Go, første interaktive vert
Fyne og statisk SVG-eksport. [Målarkitektur](target-architecture.md) og
[PLAN-003](implementation-plan.md) erstatter tidligere P4–P6/Rust/C-ABI-instrukser.

XFMD kan fortsatt vise genererte SVG-bilder i Markdown. Det krever ikke en
SDUI-parser eller SDL-runtime i XFMD. Ekte interaktiv SDUI i XFMD er en eventuell
senere konsument, med egen avklart leveranse og nødvendig port.

Eksisterende FoxBoxUiOverlay/FoxBoxUiInput, lokal simulering, diagrammer og
Cairo/PDF-kode kan være referansemateriale. Bevar dem i sine repoer; ingen
opprydding eller merge av det gamle BoxUI-worktreet inngår i G6-leveransen.
Historiske pinner og gjenbrukskilder står i [worktree-kartet](renderer-extraction-and-language-direction.md).
G6-implementasjonen ble basert på kontrollert ren XFMD-main c245fd9 i eget
xfmd-sdl-navigation-worktree. Phase 049/050 er pushet; PR #38 er ikke merget.

Ikke start en ny FOX-backend, privat SDUI-parser eller C-ABI for å følge den gamle
handoffen. Første Go-leveranse ligger under SDUI/go og SystemDesignLanguage/go.
