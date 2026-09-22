"""Data and packet projections; no field widths, storage or sources inferred."""


def build_data_views(views, diagram_type):
    relations = views.relations
    for dataset in sorted(n for n, k in views.kinds.items() if k == 'dataset'):
        sources = [r for r in relations['from'] if r['object'] == dataset]
        families = {r['subject'] for r in sources}
        holders = [r for r in relations['holds'] if r['object'] == dataset]
        holder_names = {r['subject'] for r in holders}
        owners = [r for r in relations['owns'] if r['object'] in holder_names]
        contracts = [r for r in relations['upholds'] if r['subject'] in families | {dataset}]
        views.diagram('VP09-data-' + dataset, 'Dataopprinnelse og holder: ' + dataset,
                      sources + holders + owners + contracts, [dataset])
    for contract in sorted(n for n, k in views.kinds.items() if k == 'contract'):
        variants = [r for r in relations['defines'] if r['subject'] == contract]
        variant_names = {r['object'] for r in variants}
        fields = [r for r in relations['has-field'] if r['subject'] in variant_names | {contract}]
        permits = [r for r in relations['permits'] if r['subject'] == contract]
        views.diagram('VP09-contract-' + contract, 'Kontraktstruktur: ' + contract,
                      variants + fields + permits, [contract])
    for encoding in sorted(n for n, k in views.kinds.items() if k == 'encoding'):
        variant = next(r['object'] for r in relations['encodes'] if r['subject'] == encoding)
        contract = next(r['subject'] for r in relations['defines'] if r['object'] == variant)
        entries = sorted([r for r in views.facts if r['verb'] == 'places' and r['subject'] == encoding],
                         key=lambda r: r['offset'])
        fields = {r['field'] for r in entries}
        context = fields | {encoding, variant, contract}
        facts = [r['id'] for r in views.facts if r['subject'] in context]
        syntax = ['packet']
        elements = []
        for r in entries:
            first, last = r['offset'], r['offset'] + r['width'] - 1
            extent = str(first) if first == last else f'{first}-{last}'
            syntax.append(f'    {extent}: "{r["field"]}"')
            elements.append(dict(field=r['field'], first_bit=first, last_bit=last, fact=r['id']))
        props = {r['property']: r['value'] for r in views.facts if r['subject'] == encoding and r['verb'] == 'has'}
        title = f'Packet: {encoding} / {variant} — {props["byte-order"]}, {props["bit-order"]}'
        views.diagrams.append(diagram_type('VP10-' + encoding, title, fields, kind='packet',
                              source_facts=facts, elements=elements, syntax='\n'.join(syntax) + '\n'))
    if not any(k == 'dataset' for k in views.kinds.values()):
        views.gaps.append(dict(viewpoint='VP09', model_id='—', code='NO_DATASETS', message='Ingen Dataset deklarert.'))
    encoded = {r['object'] for r in relations['encodes']}
    for variant, kind in sorted(views.kinds.items()):
        if kind == 'variant' and variant not in encoded:
            views.gaps.append(dict(viewpoint='VP10', model_id=variant, code='NO_ENCODING',
                                   message='Ingen fast wirelayout; logisk variant gir ikke packet.'))
    if not any(k == 'encoding' for k in views.kinds.values()):
        views.gaps.append(dict(viewpoint='VP10', model_id='—', code='NO_ENCODINGS', message='Ingen Encoding deklarert.'))


def data_tables(views):
    if 'VP09' not in views.selected:
        return []
    lines = ['', '### VP09 — felt og kontraktegenskaper', '',
             '| Modellfaktum | Kilde-ID |', '| --- | --- |']
    for r in views.facts:
        if r['verb'] == 'has' and views.kinds[r['subject']] in ('field', 'contract', 'encoding'):
            lines.append(f'| {r["text"]} | {r["id"]} |')
    lines += ['', '### VP09 — projeksjonsansvar', '',
              '| Functionality | Dataset | Datagram-familie | Faktum |', '| --- | --- | --- | --- |']
    for r in views.facts:
        if r['verb'] == 'projects':
            lines.append(f'| {r["subject"]} | {r["dataset"]} | {r["datagram"]} | {r["id"]} |')
    return lines + ['']
