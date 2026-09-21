# Aptering — SDUI med prototype-widgets

Knapper og inputfelt tegnet i den samme boksinndelingen som treemap-prøven.
Bildet er statisk; kontrollenes etiketter og feltverdier hentes gjennom SDUI-parseren.
Plasseringen er en avgrenset referansekomposisjon, ikke en generell layoutmotor.
Overskriftsradene inneholder Cursor AV/Δ for Lengde og Diameter, O/B · U/B og
BarkNOR-kalibrering for Diameter, samt taperNOR/Mixed for Stammeforløp.

![Apteringsflaten med tegnede knapper og inputfelt](concept1-bucking.widgets.svg)

Operatør, Stamme-ID og Merknad er lagt til som prototypeinnhold. De er ikke en
påstand om eksisterende felter i Ponsse. Øvrige verdier er også eksempeldata.

## Det begrensede widgetbiblioteket

![Normal, trykket, fokusert og deaktivert knapp; tekstfelt](prototype-controls.svg)

## Prøv utfylling og knappetrykk

Åpne [den lokale HTML-demoen](prototype-controls.html) i en nettleser.
Du kan skrive i feltene, trykke knapper med mus eller tastatur og velge
«Skriv ut utfylte verdier». Utskriften bruker de aktuelle feltverdiene og bryter
lange verdier over flere linjer. Nettleserens utskriftsdialog kan også lagre PDF.
Feltene tilbakestilles til eksempelverdiene når siden lastes på nytt.

HTML-demoen viser widgetsettet; den er ikke en interaktiv utgave av hele
apteringsflaten. Knappene gir bare lokal visuell respons og en trykkteller.

[SDUI-kilde for widgetprøven](prototype-controls.sdui) ·
[SDUI-kilde for apteringsflaten](concept1-bucking.sdui) ·
[Omfang og videreføring til FOX](../docs/prototype-widgets.md)
