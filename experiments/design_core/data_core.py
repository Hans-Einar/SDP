"""Bounded data contracts and fixed binary layout; no schema IO or execution."""
from collections import defaultdict, namedtuple
import re

Integer = namedtuple('Integer', 'value span')
Projection = namedtuple('Projection', 'subject dataset datagram span')
Placement = namedtuple('Placement', 'subject field offset width span')
KINDS = {'dataset', 'database', 'datagram', 'contract', 'variant', 'field', 'encoding'}
SIGNATURES = {
    'upholds': (('dataset', 'datagram'), 'contract'),
    'from': ('datagram', 'dataset'),
    'holds': (('unit', 'database'), 'dataset'),
    'defines': ('contract', 'variant'),
    'has-field': (('contract', 'variant'), 'field'),
    'encodes': ('encoding', 'variant'),
}
PROPERTIES = {
    'completeness': frozenset(('open', 'closed')),
    'value-type': frozenset(('unsigned', 'signed', 'boolean', 'text', 'bytes', 'decimal')),
    'presence': frozenset(('required', 'optional')),
    'byte-order': frozenset(('big-endian', 'little-endian')),
    'bit-order': frozenset(('most-significant-first', 'least-significant-first')),
}
PROPERTY_KINDS = {'completeness': 'contract', 'value-type': 'field', 'presence': 'field',
                  'byte-order': 'encoding', 'bit-order': 'encoding'}


def integer(parser):
    token = parser.current
    if not re.fullmatch(r'0|[1-9][0-9]{0,5}', token.text):
        parser.fail('Expected a canonical integer between 0 and 65536')
    parser.take()
    if int(token.text) > 65536:
        parser.fail('Integer exceeds fixed-layout limit', 'LAYOUT_LIMIT')
    return Integer(int(token.text), token.span)


def parse_statement(parser, subject, verb, start):
    target = parser.identifier()
    if verb == 'projects':
        parser.expect('into')
        datagram = parser.identifier()
        return Projection(subject, target, datagram, parser.finish(start))
    parser.expect('at')
    offset = integer(parser)
    parser.expect('bits')
    width = integer(parser)
    return Placement(subject, target, offset, width, parser.finish(start))


def sentence(s):
    if isinstance(s, Projection):
        return f'{s.subject.name} projects {s.dataset.name} into {s.datagram.name}.'
    return f'{s.subject.name} places {s.field.name} at {s.offset.value} bits {s.width.value}.'


def validate_data(model, symbols, core):
    # Invoked only after name/signature checking has succeeded.
    errors = []
    outgoing, incoming, props, placements = defaultdict(list), defaultdict(list), {}, defaultdict(list)
    for s in model.statements:
        if isinstance(s, core.Relation):
            outgoing[s.subject.name, s.verb].append(s.object.name)
            incoming[s.object.name, s.verb].append(s.subject.name)
        elif isinstance(s, core.PropertyAssignment):
            props[s.subject.name, s.property] = s.value
        elif isinstance(s, Placement):
            placements[s.subject.name].append(s)

    def error(code, name, message, span=None):
        errors.append(core.Diagnostic(code, name + ': ' + message, span or symbols[name].span))

    def one(name, verb, reverse=False):
        values = (incoming if reverse else outgoing)[name, verb]
        if len(values) != 1:
            error('DATA_CARDINALITY', name, f'exactly one {verb} relation required')
            return None
        return values[0]

    for name, d in symbols.items():
        kind = d.kind
        if kind in ('dataset', 'datagram'):
            contract = one(name, 'upholds')
            one(name, 'holds', True) if kind == 'dataset' else one(name, 'from')
            if contract and props.get((contract, 'completeness')) == 'closed':
                variants = outgoing[contract, 'defines']
                if kind == 'dataset' and (variants or not outgoing[contract, 'has-field']):
                    error('CONTRACT_SHAPE', name, 'closed Dataset requires record fields and no variants')
                if kind == 'datagram' and not variants:
                    error('CONTRACT_SHAPE', name, 'closed Datagram requires declared variants')
        if kind == 'variant':
            parent = one(name, 'defines', True)
            if parent and props.get((parent, 'completeness')) == 'closed' and not outgoing[name, 'has-field']:
                error('CONTRACT_SHAPE', name, 'closed variant requires payload fields')
        if kind == 'field':
            one(name, 'has-field', True)
        for prop, required_kind in PROPERTY_KINDS.items():
            if kind == required_kind and (name, prop) not in props:
                error('MISSING_CONTRACT_PROPERTY', name, 'missing ' + prop)

    for s in model.statements:
        if isinstance(s, Projection) and outgoing[s.datagram.name, 'from'] != [s.dataset.name]:
            error('PROJECTION_SOURCE_MISMATCH', s.subject.name,
                  'projection disagrees with Datagram source', s.span)

    for name, d in symbols.items():
        if d.kind != 'encoding':
            continue
        variant = one(name, 'encodes')
        parents = incoming[variant, 'defines'] if variant else []
        if len(parents) != 1:
            continue
        contract = parents[0]
        if props.get((contract, 'completeness')) != 'closed':
            error('OPEN_ENCODING', name, 'packet requires a closed contract')
        if not any(symbols[n].kind == 'datagram' for n in incoming[contract, 'upholds']):
            error('ENCODING_FAMILY', name, 'variant must belong to a Datagram contract')
        expected = set(outgoing[contract, 'has-field'] + outgoing[variant, 'has-field'])
        entries = sorted(placements[name], key=lambda s: (s.offset.value, s.field.name))
        seen = [s.field.name for s in entries]
        if set(seen) != expected or len(seen) != len(set(seen)) or not entries:
            error('ENCODING_FIELDS', name, 'place every header/variant field exactly once')
        cursor = 0
        for s in entries:
            field, width = s.field.name, s.width.value
            if s.offset.value != cursor or width < 1:
                error('ENCODING_RANGE', name, 'positive contiguous fields required; no overlap/gap', s.span)
            cursor = s.offset.value + width
            if cursor > 65536:
                error('LAYOUT_LIMIT', name, 'packet exceeds 65536 bits', s.span)
            value_type = props.get((field, 'value-type'))
            valid = ((value_type in ('unsigned', 'signed') and 1 <= width <= 64) or
                     (value_type == 'boolean' and width == 1) or
                     (value_type == 'bytes' and width > 0 and width % 8 == 0))
            if not valid or props.get((field, 'presence')) != 'required':
                error('ENCODING_TYPE', name, 'field needs a supported fixed-width required type', s.span)
    return errors
