"""Activity delivery views and a review document generated from the same SDL facts."""


def build_plan_views(views):
    rel = views.relations
    refinements = rel['refines']
    children = {r['subject'] for r in refinements}
    views.diagram('VP06-roots', 'Aktivitetsrøtter', [],
                  [n for n, k in views.kinds.items() if k == 'activity' and n not in children])
    for parent in sorted({r['object'] for r in refinements}):
        children = [r for r in refinements if r['object'] == parent]
        activities = {r['subject'] for r in children} | {parent}
        deliveries = [r for r in rel['delivers'] if r['subject'] in activities]
        views.diagram('VP06-detail-' + parent, 'Aktivitetsinndeling: ' + parent, children + deliveries)
    for parent in sorted({r['subject'] for r in rel['addresses']}):
        addresses = [r for r in rel['addresses'] if r['subject'] == parent]
        functions = {r['object'] for r in addresses}
        owners = [r for r in rel['owns'] if r['object'] in functions]
        views.diagram('VP06-work-' + parent, 'Planlagt ansvar: ' + parent, addresses + owners)
    views.diagram('VP06-dependencies', 'Eksplisitte aktivitetsavhengigheter', rel['depends-on'])


def implementation_markdown(views, rendered=False, source_link=None):
    rel = views.relations
    statuses = {r['subject']: (r['value'], r['id']) for r in views.facts
                if r['verb'] == 'has' and r['property'] == 'implementation-status'}
    owners = {r['object']: r for r in rel['owns']}
    addressed = {r['object'] for r in rel['addresses']}
    pending = sorted(n for n, kind in views.kinds.items() if kind == 'functionality' and n not in addressed)
    work = {r['subject'] for r in rel['addresses']} | set(statuses)
    parents = sorted({r['object'] for r in rel['refines'] if r['subject'] in work})
    grouped = {r['subject'] for r in rel['refines'] if r['object'] in parents}
    parents += sorted(work - grouped - set(parents))
    diagrams = {d.ident: d for d in views.diagrams}
    lines = ['# SDL — generert implementasjonsplan', '',
             'Denne rapporten er generert av SDL-verktøyet fra den validerte designkilden.',
             'Status er modellens påstand. Planlagt Go-/UI-arbeid er ikke implementert av denne rapporten.', '',
             (f'[SDL-kilde]({source_link}) · ' if source_link else '') +
             ('[Alle viewpoints](printout.md)' if rendered else '[Alle viewpoints](viewpoints.md)'), '',
             f'{len(addressed)} Functionality-er har eksplisitt milepæl-/ansvarskobling; {len(pending)} mangler.', '']

    def illustrate(ident):
        d = diagrams.get(ident)
        if d is None:
            return
        if rendered:
            lines.extend([f'![{d.title}](diagrams/{d.ident}.svg)', ''])
        else:
            lines.extend(['```mermaid', d.mermaid(views.kinds).rstrip(), '```', ''])

    shown_scenarios = set()
    for parent in parents:
        state, state_fact = statuses.get(parent, ('unspecified', '—'))
        lines += [f'## {parent}', '', f'Status: **{state}**. Kilde: {state_fact}.', '']
        deliveries = [r for r in rel['delivers'] if r['subject'] == parent]
        if deliveries:
            lines += ['Leveransebidrag: ' + ', '.join(f'{r["object"]} ({r["id"]})' for r in deliveries) + '.', '']
        illustrate('VP06-detail-' + parent)
        children = [r for r in rel['refines'] if r['object'] == parent]
        milestones = sorted({r['subject'] for r in children}) or [parent]
        for milestone in milestones:
            state, state_fact = statuses.get(milestone, ('unspecified', '—'))
            lines += [f'### {milestone}', '', f'Status: **{state}**. Kilde: {state_fact}.', '']
            dependencies = [r for r in rel['depends-on'] if r['subject'] == milestone]
            lines += ['Forutsetninger: ' + (', '.join(f'{r["object"]} ({r["id"]})' for r in dependencies) or 'Ingen eksplisitt deklarert') + '.', '',
                      '| Ansvar | Logisk eier | Kilde-ID-er |', '| --- | --- | --- |']
            for r in rel['addresses']:
                if r['subject'] == milestone:
                    owner = owners[r['object']]
                    lines.append(f'| {r["object"]} | {owner["subject"]} | {r["id"]}, {owner["id"]} |')
            lines += ['']
            for r in rel['illustrates']:
                if r['object'] == milestone:
                    scenario = r['subject']
                    lines += [f'Eksempelbane: {scenario}. Kobling: {r["id"]}.', '']
                    if 'VP08-' + scenario not in diagrams:
                        lines += ['Scenariofiguren er ikke med i valgt viewpoint-utvalg.', '']
                    elif scenario not in shown_scenarios:
                        illustrate('VP08-' + scenario)
                        shown_scenarios.add(scenario)
                    else:
                        lines += [f'[Scenariofigur](diagrams/VP08-{scenario}.svg)' if rendered else
                                  f'[Scenariokilde](diagrams/VP08-{scenario}.mmd)', '']
    lines += ['## Udekket modellansvar', '', ', '.join(pending) if pending else
              'Ingen deklarerte Functionality-er mangler addresses-kobling. Dette beviser ikke full kravdekning.', '',
              'Runtime-semantikk, full AST-/ABI-schema og fysisk layoutmåling må fortsatt realiseres og testes i implementasjonsfasene.', '']
    return '\n'.join(lines)
