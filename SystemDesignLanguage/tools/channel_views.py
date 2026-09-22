"""Sequences and derived MessageSets from validated Channel/scenario facts."""
import re


def label(name):
    return re.sub(r'(?<=[a-z0-9])(?=[A-Z])', ' ', name)


def build_channel_views(views, diagram_type):
    facts, rel = views.facts, views.relations
    participations = [r for r in facts if r['verb'] == 'uses']
    props = {(r['subject'], r['property']): r['value'] for r in facts if r['verb'] == 'has'}
    contracts = {r['subject']: r['object'] for r in rel['upholds']}
    modes = {r['subject']: r['object'] for r in rel['runs-in']}
    views.message_sets = []
    for channel in sorted(n for n, k in views.kinds.items() if k == 'channel'):
        contract = contracts[channel]
        allowed = [r for r in rel['permits'] if r['subject'] == contract]
        bindings = [r for r in participations if r['channel'] == channel]
        for mode in sorted({r['mode'] for r in bindings}):
            for permit in allowed:
                message = permit['object']
                members = [r for r in bindings if r['mode'] == mode and r['message'] == message]
                senders = sorted({r['subject'] for r in members if r['role'] == 'sender'})
                receivers = sorted({r['subject'] for r in members if r['role'] == 'receiver'})
                source_facts = [r['id'] for r in rel['upholds'] if r['subject'] in (channel, message)]
                source_facts += [permit['id']] + [r['id'] for r in members]
                variants = sorted(r['object'] for r in rel['defines'] if r['subject'] == contracts[message])
                views.message_sets.append(dict(channel=channel, mode=mode, message=message,
                    contract=contract, payload_contract=contracts[message], variants=variants,
                    senders=senders, receivers=receivers, source_facts=sorted(set(source_facts))))
                if not senders or not receivers:
                    views.gaps.append(dict(viewpoint='VP08', model_id=channel, mode=mode,
                        code='INCOMPLETE_PARTICIPATION', message='Ufullstendig deltakelse for ' + message + '.'))
        if not bindings:
            views.gaps.append(dict(viewpoint='VP08', model_id=channel, code='NO_PARTICIPANTS', message='Ingen eksplisitt deltakelse.'))
    for scenario in sorted(n for n, k in views.kinds.items() if k == 'scenario'):
        mode = modes[scenario]
        steps = sorted([r for r in facts if r['verb'] == 'step' and r['subject'] == scenario], key=lambda r: r['ordinal'])
        nodes = {r[role] for r in steps for role in ('sender', 'receiver')}
        order = list(dict.fromkeys(r[role] for r in steps for role in ('sender', 'receiver')))
        lines = ['sequenceDiagram'] + [f'    participant n_{n} as {label(n)}' for n in order]
        source_facts = {r['id'] for r in facts if r['subject'] == scenario}
        elements = []
        for s in steps:
            channel, message = s['channel'], s['message']
            bindings = [r for r in participations if r['channel'] == channel and r['message'] == message and r['mode'] == mode
                        and (r['subject'], r['role']) in ((s['sender'], 'sender'), (s['receiver'], 'receiver'))]
            permits = [r for r in rel['permits'] if r['subject'] == contracts[channel] and r['object'] == message]
            governing = [r for r in rel['upholds'] if r['subject'] in (channel, message)]
            type_facts = [r for r in facts if r['subject'] == message and r['verb'] in ('has', 'replies-to')]
            proof = sorted({r['id'] for r in bindings + permits + governing + type_facts} | {s['id']})
            source_facts.update(proof)
            arrow = '-->>' if props.get((message, 'message-kind')) == 'result' else '->>'
            caption = f'{s["ordinal"]}: {label(message)}' + (f' / {label(s["variant"])}' if s['variant'] else '')
            caption += f' ({label(channel)})' + (f' reply-to {s["reply_to"]}' if s['reply_to'] else '')
            lines.append(f'    n_{s["sender"]}{arrow}n_{s["receiver"]}: {caption}')
            elements.append(dict(ordinal=s['ordinal'], sender=s['sender'], receiver=s['receiver'],
                message=message, variant=s['variant'], channel=channel, reply_to=s['reply_to'], fact=s['id'], proof=proof))
        views.diagrams.append(diagram_type('VP08-' + scenario, f'Scenario: {scenario} — modus {mode}', nodes,
            kind='sequence', source_facts=sorted(source_facts), elements=elements, syntax='\n'.join(lines) + '\n'))
    if not modes:
        views.gaps.append(dict(viewpoint='VP08', model_id='—', code='NO_SCENARIOS', message='Ingen deklarerte scenario-steg; rekkefølge utledes ikke.'))


def message_tables(views):
    if 'VP08' not in views.selected:
        return []
    lines = ['', '### VP08 — avledet MessageSet per Channel og modus', '',
             'Generert fra permits og deltakelse, ikke en separat authored modell. Tom deltakelse er et hull.', '',
             '| Channel | Mode | Message / Datagram | Sender | Receiver | Kilde-ID-er |', '| --- | --- | --- | --- | --- | --- |']
    for r in views.message_sets:
        lines.append('| ' + ' | '.join([r['channel'], r['mode'], r['message'],
            ', '.join(r['senders']) or '—', ', '.join(r['receivers']) or '—', ', '.join(r['source_facts'])]) + ' |')
    return lines + ['']
