# SDUI — komposisjonsavklaring 2026-09-21

Dokumentarbeid, ingen parser/runtime/renderer-endring.

Eierens kanoniske hjørner ^</>^/v</>v erstatter forrige forslag. Ratio med
horisontalt strekk styrer full bredde og avledet høyde, ikke contain-fit.
Header/body/footer og navngitt mainBody-referanse er bevart og konkretisert
i komposisjonsforslaget. Senere eierpresisering samme dato: font=10/12 er
absolutt tekststørrelse og resize skalerer ikke innholdet. Den konkrete native
fontenheten gjenstår. Formatering etter hver komponent står før separator;
komma fortsetter horisontalt og semikolon starter neste rad. Krav R21 og
P1/P2/P5-planen registrerer parser-, måle- og GUI-kontrollene som gjenstår.
Ytterligere presisering: flere nivåer av `<>` danner grupper med egen layout.
Et designeksempel viser tre nivåer og lokal separatorvirkning. R13 og P1/P2
krever bevart gruppestruktur, også når en eksplisitt gruppe bare har ett barn.

Kontroll: relative Markdown-lenker, balanserte kodegjerder, fravær av pikselmål
i nye eksempler og aritmetikken W=1, H=0.25, ratio=16/9 → 1×0.5625 med
høydeoverflow. Dette er dokument-/aritmetikkontroll, ikke en layoutimplementasjon.
Ingen ny GUI-/PDF-kjøring. Ingen commit/push eller endringer i eksterne repoer.
