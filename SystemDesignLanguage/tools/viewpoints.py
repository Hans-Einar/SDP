"""Deterministic structural viewpoints from validated design-core, without new syntax."""
from collections import defaultdict
from dataclasses import dataclass, field
from pathlib import Path
import sys

# Reuse the one existing SDL frontend until the planned Go port replaces it.
sys.path.insert(0, str(Path(__file__).resolve().parents[2] / 'experiments/design_core'))
import design_core as core
from goal_views import build_goal_views
from data_views import build_data_views, data_tables
from channel_views import build_channel_views, message_tables

CATALOG = [
 ('VP01', 'Bruksmål og sporbarhet', 'supported', 'pursues, supports og contributes-to; modellens omfang, uten oppdiktet System-grense.'),
 ('VP02', 'Arkitektur og logisk inndeling', 'supported', 'Container/Unit og contains; bibliotekstruktur er ikke en deployment-allokering.'),
 ('VP03', 'Ansvar og kapabiliteter over arkitekturen', 'supported', 'owns, realizes og provides. Capability er ikke Feature.'),
 ('VP04', 'Grensesnitt og samarbeid', 'supported', 'consumes viser bruk; ingen tilbyder, Channel eller kjørbar meldingsflyt utledes.'),
 ('VP05', 'Avhengigheter per modus', 'supported', 'requires in mode; modi har ingen implisitt arv.'),
 ('VP06', 'Aktivitetsdetaljering', 'supported', 'refines er detaljering, ikke rekkefølge eller tilstandsoverganger.'),
 ('VP07', 'Features over arkitekturen', 'supported', 'contributes-to, owns og eksplisitt allocated-to per modus. Uspesifisert allokering vises som hull.'),
 ('VP08', 'Channel-kontrakter og sekvenser', 'supported', 'Eksplisitte scenario-steg validert mot permits, deltakelse, modus og request/resultat-korrelasjon.'),
 ('VP09', 'Dataset, Datagram og persistent Database', 'supported', 'Eksplisitte holdere, kilde, kontrakter, varianter, felt og projeksjoner.'),
 ('VP10', 'Datagram-koding og packet', 'supported', 'Kun closed kontrakt med validert Encoding og eksplisitte bitplasseringer.'),
 ('VP11', 'Egenskaper, sporbarhet og modellhull', 'supported', 'Deklarasjoner og alle fakta med kildeposisjoner; støttegrenser beholdes.'),
]


@dataclass
class Diagram:
    ident: str
    title: str
    nodes: set = field(default_factory=set)
    edges: list = field(default_factory=list)
    kind: str = 'flowchart'
    source_facts: list = field(default_factory=list)
    elements: list = field(default_factory=list)
    syntax: str = ''

    def mermaid(self, kinds):
        if self.syntax:
            return self.syntax
        lines = ['flowchart LR']
        for name in sorted(self.nodes):
            lines.append(f'    n_{name}["{name} ({kinds[name]})"]')
        for subject, verb, obj, _ in sorted(self.edges):
            lines.append(f'    n_{subject} -->|{verb}| n_{obj}')
        return '\n'.join(lines) + '\n'


class Views:
    def __init__(self, source, selected=None):
        self.selected = set(selected) if selected is not None else {v[0] for v in CATALOG}
        if self.selected - {v[0] for v in CATALOG}:
            raise ValueError('Unknown viewpoint')
        self.model, diagnostics = core.check(source)
        if diagnostics:
            raise ValueError(core.to_json(diagnostics))
        self.kinds = {d.name.name: d.kind for d in self.model.declarations}
        self.relations = defaultdict(list)
        self.facts = []
        for index, s in enumerate(self.model.statements):
            record = {'id': f'f{index:04d}', 'node': type(s).__name__, 'line': s.span.line,
                      'span': core.to_json(s.span), 'text': core.sentence(s)}
            if isinstance(s, core.Relation):
                record.update(subject=s.subject.name, verb=s.verb, object=s.object.name)
                self.relations[s.verb].append(record)
            elif isinstance(s, core.Dependency):
                record.update(subject=s.subject.name, verb='requires', object=s.interface.name, mode=s.mode.name)
            elif isinstance(s, core.Allocation):
                record.update(subject=s.subject.name, verb='allocated-to', object=s.container.name, mode=s.mode.name)
                self.relations['allocated-to'].append(record)
            elif isinstance(s, core.Projection):
                record.update(subject=s.subject.name, verb='projects', dataset=s.dataset.name, datagram=s.datagram.name)
            elif isinstance(s, core.Placement):
                record.update(subject=s.subject.name, verb='places', field=s.field.name, offset=s.offset.value, width=s.width.value)
            elif isinstance(s, core.Participation):
                record.update(subject=s.subject.name, verb='uses', channel=s.channel.name, role=s.role,
                              message=s.message.name, mode=s.mode.name)
            elif isinstance(s, core.Step):
                record.update(subject=s.subject.name, verb='step', ordinal=s.ordinal.value,
                              message=s.message.name, variant=s.variant.name if s.variant else None,
                              sender=s.sender.name, receiver=s.receiver.name, channel=s.channel.name,
                              reply_to=s.reply_to.value if s.reply_to else None)
            else:
                record.update(subject=s.subject.name, verb='has', property=s.property, value=s.value)
            self.facts.append(record)
        self.diagrams = []
        self.gaps = []
        self.build_diagrams()
        self.diagrams = [d for d in self.diagrams if d.ident.split('-')[0] in self.selected]
        self.gaps = [g for g in self.gaps if g['viewpoint'] in self.selected]
        if 'VP08' not in self.selected:
            self.message_sets = []

    def diagram(self, ident, title, records, nodes=()):
        d = Diagram(ident, title, set(nodes))
        for record in records:
            d.nodes.update([record['subject'], record['object']])
            d.edges.append((record['subject'], record['verb'], record['object'], record['id']))
        if d.nodes:
            self.diagrams.append(d)

    def build_diagrams(self):
        build_goal_views(self)
        build_data_views(self, Diagram)
        build_channel_views(self, Diagram)
        children = {r['object'] for r in self.relations['contains']}
        roots = [name for name, kind in self.kinds.items() if kind in ('unit', 'container') and name not in children]
        self.diagram('VP02-roots', 'Arkitekturrøtter — ingen kobling/allokering er utledet', [], roots)
        for parent in sorted({r['subject'] for r in self.relations['contains']}):
            self.diagram('VP02-' + parent, 'Logisk inndeling: ' + parent,
                         [r for r in self.relations['contains'] if r['subject'] == parent])
        for cap in sorted(n for n, k in self.kinds.items() if k == 'capability'):
            contributions = [r for r in self.relations['realizes'] if r['object'] == cap]
            functions = {r['subject'] for r in contributions}
            owners = [r for r in self.relations['owns'] if r['object'] in functions]
            offers = [r for r in self.relations['provides'] if r['object'] == cap]
            self.diagram('VP03-' + cap, 'Bidrag til kapabilitet: ' + cap, contributions + owners, [cap])
            if offers:
                self.diagram('VP03-' + cap + '-offers', 'Tilbydere av kapabilitet: ' + cap, offers, [cap])
        for mode in sorted(n for n, k in self.kinds.items() if k == 'mode'):
            records = [r for r in self.facts if r['verb'] == 'requires' and r['mode'] == mode]
            self.diagram('VP05-' + mode, 'Nødvendige porter i modus: ' + mode, records)
        activities = [n for n, k in self.kinds.items() if k == 'activity']
        self.diagram('VP06-refinement', 'Detaljert aktivitet → overordnet aktivitet', self.relations['refines'], activities)

    def node_map(self, diagram):
        return {f'n_{d.name.name}': {'model_id': d.name.name, 'kind': d.kind,
                                   'span': core.to_json(d.span)}
                for d in self.model.declarations if d.name.name in diagram.nodes}

    def markdown(self, rendered=False):
        lines = ['# SDL-viewpoints — generert modellrapport', '',
                 'Generert fra validert SDL. Struktur og designpåstander, ikke observert kjøring.',
                 'Ingen håndskrevet arkitekturfakta er lagt til av generatoren.', '',
                 '| Viewpoint | Status | Grunnlag / mangel |', '| --- | --- | --- |']
        for ident, title, status, note in CATALOG:
            if ident not in self.selected:
                continue
            lines.append(f'| {ident} — {title} | {"Tilgjengelig" if status == "supported" else "Kan ikke genereres"} | {note} |')
        for prefix in ('VP01', 'VP02', 'VP03', 'VP05', 'VP06', 'VP07', 'VP08', 'VP09', 'VP10'):
            if prefix not in self.selected:
                continue
            title = next(title for ident, title, _, _ in CATALOG if ident == prefix)
            lines += ['', f'## {prefix} — {title}', '']
            if prefix == 'VP02':
                lines += ['Container er en erklært runtimegrense. Unit-røtter viser logisk struktur.',
                          'contains angir ikke deployment. Eksplisitt Functionality-allokering vises per modus i VP07.', '']
            if prefix == 'VP01':
                lines += ['Bruksmålskartene viser Actors, støttende Features og direkte Functionality-bidrag.',
                          'De etterfølgende Feature-kartene detaljerer bidragene med samme modellidentiteter.',
                          'Oppdelingen endrer ingen relasjoner og innfører ingen System-grense.', '']
            for d in self.diagrams:
                if not d.ident.startswith(prefix):
                    continue
                lines += [f'### {d.title}', '']
                if rendered:
                    lines += [f'![{d.title}](diagrams/{d.ident}.svg)', '']
                else:
                    lines += ['```mermaid', d.mermaid(self.kinds).rstrip(), '```', '']
                facts = ', '.join(sorted(set(d.source_facts + [edge[3] for edge in d.edges]))) or 'Kun deklarasjoner'
                lines += [f'Kildegrunnlag: {facts}.', '']
            gaps = [g for g in self.gaps if g['viewpoint'] == prefix]
            if gaps:
                lines += ['### Modellhull i dette utsnittet', '',
                          '| Identitet | Modus | Mangel |', '| --- | --- | --- |']
                for gap in gaps:
                    lines.append(f'| {gap["model_id"]} | {gap.get("mode") or "—"} | {gap["message"]} |')
                lines += ['']
        lines += data_tables(self)
        lines += message_tables(self)
        if 'VP04' in self.selected:
            lines += ['## VP04 — Grensesnittbruk', '',
                  'Tabellen dekker alle consumes-fakta. Portnavn er ikke Channel-kontrakter eller tilbyderkoblinger.', '',
                  '| Unit / Container | Interface | Faktum | Kildelinje |', '| --- | --- | --- | --- |']
            for r in self.relations['consumes']:
                lines.append(f'| {r["subject"]} | {r["object"]} | {r["id"]} | {r["line"]} |')
        if 'VP11' not in self.selected:
            return '\n'.join(lines) + '\n'
        lines += ['', '## VP11 — Egenskaper og fullstendig faktaregister', '',
                  'Registeret inkluderer alle fakta, også de som ikke har en egen tegning.', '',
                  '| ID | Utsagn | Kildelinje |', '| --- | --- | --- |']
        for r in self.facts:
            statement = r['text']
            lines.append(f'| {r["id"]} | {statement} | {r["line"]} |')
        lines += ['', '### Deklarasjonsregister', '', '| Identitet | Type | Kildelinje |', '| --- | --- | --- |']
        for d in self.model.declarations:
            lines.append(f'| {d.name.name} | {d.kind} | {d.span.line} |')
        return '\n'.join(lines) + '\n'
