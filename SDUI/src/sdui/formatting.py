"""Typed local formatting checks and canonical properties; no viewport geometry."""
from .ast import SduiError

ARROWS = {
    '^<': {'align-x': 'start', 'align-y': 'start'},
    '-^': {'align-x': 'center', 'align-y': 'start'},
    '>^': {'align-x': 'end', 'align-y': 'start'},
    '<-': {'align-x': 'start', 'align-y': 'center'},
    '--': {'align-x': 'center', 'align-y': 'center'},
    '->': {'align-x': 'end', 'align-y': 'center'},
    'v<': {'align-x': 'start', 'align-y': 'end'},
    '-v': {'align-x': 'center', 'align-y': 'end'},
    '>v': {'align-x': 'end', 'align-y': 'end'},
    '<->': {'x': 'fill'}, '^|v': {'y': 'fill'},
    '>-<': {'x': 'content'}, '>|<': {'y': 'content'},
}
ALIASES = {key[::-1]: key for key in list(ARROWS)[:9]}
ENUMS = {
    'align-x': {'start', 'center', 'end'}, 'align-y': {'start', 'center', 'end'},
    'justify': {'start', 'center', 'end', 'between'},
    'items': {'start', 'center', 'end', 'stretch'},
    'overflow-x': {'error', 'clip', 'scroll'}, 'overflow-y': {'error', 'clip', 'scroll'},
    'wrap': {'none', 'wrap'},
}
RELATIVE = {'scale', 'scale-x', 'scale-y', 'min-x', 'min-y', 'max-x', 'max-y',
            'gap', 'gap-x', 'gap-y', 'padding'}


def fail(code, message, node):
    raise SduiError(code, message, node.span)


def formatting(node, kind=None):
    result = {}
    for rule in node.layout:
        key, literal = rule.name, rule.value
        value = literal.value
        if key == 'arrow':
            canonical = ALIASES.get(value, value)
            if canonical not in ARROWS:
                fail('layout-arrow', f'Unknown layout shape {value!r}', rule)
            additions = ARROWS[canonical]
        else:
            if key in ENUMS:
                valid = literal.kind == 'enum' and value in ENUMS[key]
            elif key in RELATIVE or key == 'font':
                valid = literal.kind == 'number' and value >= 0
                if key in {'scale', 'scale-x', 'scale-y', 'font'}:
                    valid = valid and value > 0
                if key == 'padding' and literal.kind == 'tuple':
                    valid = len(value) == 4 and all(x >= 0 for x in value)
            elif key in {'x', 'y'}:
                valid = (literal.kind == 'enum' and value in {'content', 'fill'} or
                         literal.kind == 'fr' and value > 0)
                if literal.kind == 'fr':
                    value = f'{value:g}fr'
            elif key == 'ratio':
                valid = literal.kind == 'ratio' and all(x > 0 for x in value)
            elif key in {'enabled', 'visible'}:
                valid = literal.kind == 'boolean'
            else:
                fail('layout-property', f'Unknown formatting property {key}', rule)
            if not valid:
                fail('layout-type', f'Invalid value for {key}', rule)
            additions = {key: value}
        for name, value in additions.items():
            if name in result:
                fail('layout-conflict', f'Multiple rules for {name}', rule)
            result[name] = value
    check_combination(result, kind or node.kind, node)
    return result


def check_combination(props, kind, node):
    for axis in ('x', 'y'):
        choices = [key for key in ('scale', f'scale-{axis}', axis) if key in props]
        if len(choices) > 1:
            fail('layout-conflict', f'Competing {axis} sizes: {choices}', node)
        low, high = props.get(f'min-{axis}', 0), props.get(f'max-{axis}', float('inf'))
        if low > high:
            fail('size-range', f'min-{axis} exceeds max-{axis}', node)
    if 'ratio' in props:
        if kind != 'frame':
            fail('ratio-scope', 'Aspect ratio applies to frames', node)
        axes = [axis for axis in ('x', 'y') if f'scale-{axis}' in props or axis in props]
        if 'scale' in props or len(axes) > 1:
            fail('ratio-axis', 'Aspect ratio permits only one controlling axis', node)
        if any(props.get(axis, 'fill') != 'fill' for axis in axes):
            fail('ratio-axis', 'With ratio use one scale axis or fill, not content/fr', node)
    if 'wrap' in props and kind != 'group':
        fail('wrap-scope', 'Wrap applies to widget groups', node)
    if props.get('wrap') == 'wrap':
        if len(node.rows) != 1 or props.get('x') == 'fill' or str(props.get('x', '')).endswith('fr'):
            fail('wrap-layout', 'Wrap needs one explicit row without horizontal fill/fr', node)
