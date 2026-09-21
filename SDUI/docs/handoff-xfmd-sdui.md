# FOX/XFMD — tidligere handoff, avløst av Go/Fyne-retning

Oppdatert 2026-09-21. Dette er ikke lenger et aktivt implementasjonsoppdrag.
Eieren har valgt selvstendige SDL/SDUI-biblioteker i Go, første interaktive vert
Fyne og statisk SVG-eksport. [Målarkitektur](target-architecture.md) og
[PLAN-003](implementation-plan.md) erstatter tidligere P4–P6/Rust/C-ABI-instrukser.

XFMD kan fortsatt vise genererte SVG-bilder i Markdown. Det krever ikke en
SDUI-parser eller SDL-runtime i XFMD. Ekte interaktiv SDUI i XFMD er en eventuell
senere konsument, med egen avklart leveranse og nødvendig port.

Eksisterende FoxBoxUiOverlay/FoxBoxUiInput, lokal simulering, diagrammer og
Cairo/PDF-kode kan være referansemateriale. Bevar dem i sine repoer; ingen
worktree-opprydding, PR-endring, merge, push eller installasjon er bestilt her.
Historiske pinner og gjenbrukskilder står i [worktree-kartet](renderer-extraction-and-language-direction.md).
Deres nåstatus er ikke kontrollert på nytt i denne dokumentasjonsrunden.

Ikke start en ny FOX-backend, privat SDUI-parser eller C-ABI for å følge den gamle
handoffen. Første Go-leveranse ligger under SDUI/go og SystemDesignLanguage/go.
