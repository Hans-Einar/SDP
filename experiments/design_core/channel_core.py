"""Mode-scoped messaging and finite scenario witnesses; never execute messages."""
from collections import defaultdict, namedtuple
from data_core import integer

Participation = namedtuple('Participation', 'subject channel role message mode span')
Step = namedtuple('Step', 'subject ordinal message variant sender receiver channel reply_to span')
KINDS = {'channel', 'message', 'scenario'}
SIGNATURES = {'permits': ('contract', ('message', 'datagram')),
              'replies-to': ('message', 'message'),
              'runs-in': ('scenario', 'mode'), 'exercises': ('scenario', 'usecase')}
PROPERTIES = {'message-kind': frozenset(('request', 'result', 'event'))}


def parse_statement(p, subject, verb, start):
    if verb == 'uses':
        channel = p.identifier()
        p.expect('as')
        role = p.choice(('sender', 'receiver'))
        p.expect('of')
        message = p.identifier()
        p.expect('in'); p.expect('mode')
        mode = p.identifier()
        return Participation(subject, channel, role, message, mode, p.finish(start))
    ordinal = integer(p)
    p.expect('sends')
    message = p.identifier()
    variant = None
    if p.current.text == 'variant':
        p.take(); variant = p.identifier()
    p.expect('from'); sender = p.identifier()
    p.expect('to'); receiver = p.identifier()
    p.expect('via'); channel = p.identifier()
    reply = None
    if p.current.text == 'reply-to':
        p.take(); reply = integer(p)
    return Step(subject, ordinal, message, variant, sender, receiver, channel, reply, p.finish(start))


def sentence(s):
    if isinstance(s, Participation):
        return f'{s.subject.name} uses {s.channel.name} as {s.role} of {s.message.name} in mode {s.mode.name}.'
    text = f'{s.subject.name} step {s.ordinal.value} sends {s.message.name}'
    if s.variant:
        text += ' variant ' + s.variant.name
    text += f' from {s.sender.name} to {s.receiver.name} via {s.channel.name}'
    if s.reply_to:
        text += ' reply-to ' + str(s.reply_to.value)
    return text + '.'


def validate_channels(model, symbols, core):
    errors = []
    rel, props, participation, steps = defaultdict(list), {}, set(), defaultdict(list)
    for s in model.statements:
        if isinstance(s, core.Relation):
            rel[s.subject.name, s.verb].append(s.object.name)
        elif isinstance(s, core.PropertyAssignment):
            props[s.subject.name, s.property] = s.value
        elif isinstance(s, Participation):
            participation.add((s.subject.name, s.channel.name, s.role, s.message.name, s.mode.name))
        elif isinstance(s, Step):
            steps[s.subject.name].append(s)

    def error(code, name, message, span=None):
        errors.append(core.Diagnostic(code, name + ': ' + message, span or symbols[name].span))

    def one(name, verb):
        values = rel[name, verb]
        if len(values) != 1:
            error('CHANNEL_CARDINALITY', name, 'exactly one ' + verb + ' relation required')
            return None
        return values[0]

    contracts, modes = {}, {}
    for name, d in symbols.items():
        if d.kind in ('message', 'channel'):
            contract = one(name, 'upholds')
            if contract:
                contracts[name] = contract
                fields, variants, permits = rel[contract, 'has-field'], rel[contract, 'defines'], rel[contract, 'permits']
                if d.kind == 'channel' and (fields or variants or (props.get((contract, 'completeness')) == 'closed' and not permits)):
                    error('CHANNEL_CONTRACT_SHAPE', name, 'Channel contract requires permits, without payload fields/variants')
                if d.kind == 'message' and (variants or permits or (props.get((contract, 'completeness')) == 'closed' and not fields)):
                    error('MESSAGE_CONTRACT_SHAPE', name, 'Message contract requires record fields, without permits/variants')
            if d.kind == 'message':
                kind = props.get((name, 'message-kind'))
                if kind is None:
                    error('MISSING_CONTRACT_PROPERTY', name, 'missing message-kind')
                replies = rel[name, 'replies-to']
                if kind == 'result':
                    if len(replies) != 1 or props.get((replies[0], 'message-kind')) != 'request':
                        error('REPLY_TYPE', name, 'result requires exactly one request message type')
                elif replies:
                    error('REPLY_TYPE', name, 'only result types can declare replies-to')
        elif d.kind == 'scenario':
            modes[name] = one(name, 'runs-in')
            if (name, 'completeness') not in props:
                error('MISSING_CONTRACT_PROPERTY', name, 'missing completeness')
            if not steps[name]:
                error('EMPTY_SCENARIO', name, 'scenario requires explicit steps')
    # Contract shape cannot be reused as both payload and messaging permissions.
    for (name, verb), values in list(rel.items()):
        if verb == 'permits' and values and (rel[name, 'has-field'] or rel[name, 'defines']):
            error('CHANNEL_CONTRACT_SHAPE', name, 'permissions cannot be mixed with payload structure')
    for s in model.statements:
        if isinstance(s, Participation):
            contract = contracts.get(s.channel.name)
            if contract and s.message.name not in rel[contract, 'permits']:
                error('MESSAGE_NOT_PERMITTED', s.subject.name, 'Channel does not permit message', s.span)

    for scenario, items in steps.items():
        ordered = sorted(items, key=lambda s: s.ordinal.value)
        if [s.ordinal.value for s in ordered] != list(range(1, len(items) + 1)):
            error('STEP_ORDER', scenario, 'step indexes must be unique and contiguous from 1')
        prior, answered = {}, set()
        for s in ordered:
            message, channel = s.message.name, s.channel.name
            channel_contract = contracts.get(channel)
            payload_contracts = rel[message, 'upholds']
            if channel_contract and message not in rel[channel_contract, 'permits']:
                error('MESSAGE_NOT_PERMITTED', scenario, 'step is outside Channel contract', s.span)
            for contract in ([channel_contract] if channel_contract else []) + payload_contracts:
                if props.get((contract, 'completeness')) != 'closed':
                    error('OPEN_SCENARIO_CONTRACT', scenario, 'steps require closed contracts', s.span)
            for participant, role in [(s.sender.name, 'sender'), (s.receiver.name, 'receiver')]:
                key = (participant, channel, role, message, modes.get(scenario))
                if key not in participation:
                    error('PARTICIPATION_MISMATCH', scenario, 'missing ' + role + ' in scenario Mode', s.span)
            if symbols[message].kind == 'datagram':
                variants = {v for c in payload_contracts for v in rel[c, 'defines']}
                if not s.variant or s.variant.name not in variants:
                    error('VARIANT_MISMATCH', scenario, 'Datagram needs its own contracted variant', s.span)
            elif s.variant:
                error('VARIANT_MISMATCH', scenario, 'record Message cannot select variant', s.span)
            kind = props.get((message, 'message-kind'), 'event')
            if kind == 'result':
                target = s.reply_to.value if s.reply_to else None
                request = prior.get(target)
                if (not request or target in answered or
                    props.get((request.message.name, 'message-kind')) != 'request' or
                    rel[message, 'replies-to'] != [request.message.name] or
                    (request.channel.name, request.sender.name, request.receiver.name) !=
                    (channel, s.receiver.name, s.sender.name)):
                    error('CORRELATION_MISMATCH', scenario, 'result must reverse one earlier matching request', s.span)
                else:
                    answered.add(target)
            elif s.reply_to:
                error('CORRELATION_MISMATCH', scenario, 'only results can use reply-to', s.span)
            prior[s.ordinal.value] = s
        if props.get((scenario, 'completeness')) == 'closed':
            pending = [i for i, s in prior.items() if props.get((s.message.name, 'message-kind')) == 'request' and i not in answered]
            if pending:
                error('INCOMPLETE_SCENARIO', scenario, 'closed scenario has unanswered requests')
    return errors
