#!/usr/bin/env python3
"""Build bounded prototype fixtures. Recorded treemap geometry, not a general layout engine."""
import json
from pathlib import Path
import sys
from html import escape

ROOT = Path(__file__).resolve().parents[1]
sys.path.insert(0, str(ROOT / 'src'))
from sdui import parse
from sdui.normalize import normalize
from sdui.prototype_widgets import prototype_widgets, svg_widget, PALETTE
from sdui.prototype_html import widget_gallery_html


class Scene:
    def __init__(self, width, height, title):
        self.parts = [f'<svg xmlns="http://www.w3.org/2000/svg" width="{width}" height="{height}" viewBox="0 0 {width} {height}">',
                      f'<title>{escape(title)}</title>',
                      f'<rect width="{width}" height="{height}" fill="#edf1e8"/>']

    def rect(self, x, y, w, h, fill='#faffef', stroke='#9aaa9d', radius=5):
        self.parts.append(f'<rect x="{x}" y="{y}" width="{w}" height="{h}" rx="{radius}" fill="{fill}" stroke="{stroke}"/>')

    def text(self, x, y, text, size=14, weight='normal', fill=PALETTE['ink']):
        self.parts.append(f'<text x="{x}" y="{y}" font-family="sans-serif" font-size="{size}" font-weight="{weight}" fill="{fill}">{escape(text)}</text>')

    def widget(self, widget, x, y, width, **state):
        self.parts.append(svg_widget(widget, x, y, width, **state)[0])

    def finish(self):
        return '\n'.join(self.parts + ['</svg>']) + '\n'


def named(widgets):
    result = {}
    for widget in widgets.values():
        key = widget.path.rsplit('/', 1)[-1]
        if key in result:
            raise ValueError('Fixture requires unique widget names: ' + key)
        result[key] = widget
    return result


def build():
    examples = ROOT / 'examples'
    controls = prototype_widgets(normalize(parse((examples/'prototype-controls.sdui').read_text()))['controls'])
    gallery = named(controls)
    bucking = named(prototype_widgets(normalize(parse((examples/'concept1-bucking.sdui').read_text()))['bucking']))
    # Use the measured treemap row regions; suppress redundant group headings to free widget space.
    layout = json.loads((ROOT/'evidence/treemap-probe/layout.json').read_text())
    rects = {node['id']: node for node in layout['nodes']}
    scene = Scene(1600, 900, 'Concept1 med statiske SDUI prototype-widgets')
    scene.text(16, 25, 'Aptering · SDUI-prototype', 19, 'bold')
    scene.text(1070, 24, 'Eksempeldata · tegnede kontroller · uten domenekobling', 13)

    def panel(leaf, row, title, heading_controls=()):
        a, b = rects[f'treemap_{leaf}'], rects[f'treemap_{row}']
        x, y, w, h = a['x'], b['y'], a['width'], b['height']
        scene.rect(x, y, w, h)
        scene.parts.append(f'<path d="M {x} {y+48} H {x+w}" stroke="#c9d2c2"/>')
        scene.text(x+14, y+31, title, 17, 'bold')
        controls_width = sum(width for _, width in heading_controls) + max(0, len(heading_controls)-1)*8
        control_x = x+w-14-controls_width
        for name, width in heading_controls:
            scene.widget(bucking[name], control_x, y+6, width)
            control_x += width+8
        return x+14, y+54, w-28, h-68

    x, y, w, _ = panel(2, 1, 'Lengde', [('lengthCursor', 126), ('lengthDelta', 44)])
    scene.text(x, y+24, '412,0 cm', 28, 'bold')
    scene.text(x, y+51, 'Måling · aggregat', 14)
    x, y, w, _ = panel(3, 1, 'Diameter', [
        ('diameterCursor', 120), ('diameterDelta', 44),
        ('barkBasis', 100), ('barkSettings', 190)])
    scene.text(x, y+24, '28,4 cm', 28, 'bold')
    scene.text(x, y+51, 'Måling · O/B', 14)

    x, y, w, _ = panel(5, 4, 'Treslag / Sortiment')
    scene.text(x, y+20, 'Gran · Sagtømmer', 21, 'bold')
    scene.text(x, y+50, 'Siste knapp · GRAN')
    scene.widget(gallery['operator'], x, y+80, w)
    scene.widget(gallery['stemId'], x, y+168, w, focused=True)
    scene.text(x, y+274, 'Inndatafeltene er prototypeinnhold.', 12)

    x, y, w, _ = panel(6, 4, 'Apteringsforslag')
    scene.widget(bucking['optimize'], x, y, 168)
    scene.widget(bucking['columns'], x+180, y, 126)
    scene.text(x, y+68, 'Neste kapp · eksempeldata', 14)
    cols = [0, 46, 292, 447]
    for i, row in enumerate([['Nr', 'Sortiment', 'Lengde', 'Topp-Ø'], ['1', 'Sagtømmer', '430 cm', '26 cm'], ['2', 'Massevirke', '310 cm', '18 cm']]):
        scene.rect(x, y+84+i*36, w, 36, '#f0f4e9' if i==0 else '#fff', '#c9d2c2', 0)
        for dx, text in zip(cols, row):
            scene.text(x+dx+10, y+107+i*36, text, 14, 'bold' if i==0 else 'normal')
    scene.widget(bucking['canonical'], x, y+224, 146, pressed=True)
    scene.widget(bucking['alternative'], x+160, y+224, 158)
    scene.text(x, y+290, 'Canonical er tegnet i trykktilstand · ingen plan er valgt.', 13)

    x, y, w, _ = panel(7, 4, 'Stammen i aggregatet')
    for i, line in enumerate(['Aktuell stamme', 'Låste feil · 0', 'Nåposisjon · 412,0 cm', 'Maks kvistet · 412,0 cm', 'Produsert / kappet · 0 cm']):
        scene.text(x, y+20+i*29, line, 14, 'bold' if i==1 else 'normal')
    scene.widget(gallery['note'], x, y+187, w)
    scene.text(x, y+286, 'Merknad er et prototypefelt.', 12)

    x, y, w, h = panel(9, 8, 'Stammeforløp', [('taperNor', 136), ('mixed', 116)])
    scene.text(x, y+24, 'Predikert · Målt · Plan · Rest · Kvalitet 1/2 · Underkjent · Utkast', 14)
    scene.rect(x, y+65, w, 112, '#f1f5e9', '#bac6b1')
    scene.text(x+20, y+113, 'Stammeprofil — plassholder for egen visuell provider', 19, 'bold')
    scene.text(x+20, y+143, 'Denne prøven tegner widgetutseende; ingen maskin- eller SDL-funksjon kjøres.')
    scene.text(x, y+h-33, 'FØLGER AGGREGAT     Lengde 412,0 cm     Relativ lengde 0,0 cm', 15)
    scene.text(x, y+h-7, 'Predikert Ø O/B 28,4 cm     Profilvolum 0,312 m³     Toppsylinder 0,260 m³', 14)
    (examples/'concept1-bucking.widgets.svg').write_text(scene.finish())

    sheet = Scene(1020, 570, 'SDUI button og input – statiske utseender')
    sheet.text(28, 40, 'SDUI · begrenset prototypebibliotek', 25, 'bold')
    sheet.text(28, 72, 'Samme widgetbeskrivelse · SVG for Markdown · lokal HTML for utfylling', 15)
    for i, (key, caption) in enumerate([('ok','Normal'), ('pressed','Trykket'), ('cancel','Fokus'), ('disabled','Deaktivert')]):
        x = 28+i*246
        sheet.text(x, 119, caption, 14)
        sheet.widget(gallery[key], x, 135, 222, pressed=key=='pressed', focused=key=='cancel')
    sheet.widget(gallery['operator'], 28, 227, 456)
    sheet.widget(gallery['stemId'], 526, 227, 456, focused=True)
    sheet.widget(gallery['note'], 28, 340, 954)
    sheet.text(28, 461, 'SVG-bildet er statisk. Åpne HTML-demoen for å skrive og prøve knapper.', 16)
    sheet.text(28, 500, 'Utskrift bruker feltverdiene som står i HTML-demoen når du skriver ut.', 15)
    (examples/'prototype-controls.svg').write_text(sheet.finish())
    (examples/'prototype-controls.html').write_text(widget_gallery_html(list(controls.values())))
    print('Built concept1-bucking.widgets.svg, prototype-controls.svg and prototype-controls.html')


if __name__ == '__main__':
    build()
