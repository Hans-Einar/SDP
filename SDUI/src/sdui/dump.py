"""Plain terminal-cell structural preview, not the future pixel/FOX layout engine."""
import math
import re
import unicodedata
from .ast import SduiError

MAX_CELLS = 2_000_000
FENCE = re.compile(r'^ {0,3}(`{3,}|~{3,})([^\r\n]*)$')


def markdown_lines(source):
    """Keep Markdown source. Omit Mermaid fenced blocks, respecting other fences."""
    fence, hidden = None, False
    lines = []
    for line in source.splitlines():
        match = FENCE.match(line)
        if fence:
            closing = re.fullmatch(r' {0,3}' + re.escape(fence[0]) +
                                   '{' + str(fence[1]) + r',}[ \t]*', line)
            if not hidden:
                lines.append(line)
            if closing:
                fence, hidden = None, False
        elif match and not (match[1][0] == '`' and '`' in match[2]):
            fence = (match[1][0], len(match[1]))
            info = match[2].strip().split()
            hidden = bool(info and info[0].lower() == 'mermaid')
            lines.append('[Mermaid utelatt]' if hidden else line)
        else:
            lines.append(line)
    return lines or ['']


def safe(text):
    # Raw Markdown means no Markdown interpretation, never terminal control execution.
    return ''.join(('    ' if c == '\t' else f'\\u{ord(c):04x}')
                   if unicodedata.category(c).startswith('C') else c for c in text)


def cell_width(text):
    return sum(0 if unicodedata.combining(c) else
               2 if unicodedata.east_asian_width(c) in {'W', 'F'} else 1 for c in text)


def wrap_line(text, width):
    lines, line, used = [], '', 0
    for char in safe(text):
        size = cell_width(char)
        if used + size > width:
            lines.append(line)
            line, used = '', 0
        line += char
        used += size
    return lines + [line]


def fit(text, width):
    line = wrap_line(text, width)[0]
    return line + ' ' * max(0, width - cell_width(line))


def horizontal_widths(items, available):
    # Structural dump: honor horizontal scale and fr, divide other cells equally.
    widths, weights = [], []
    for item in items:
        props = dict(item.layout)
        scale = props.get('scale-x', props.get('scale'))
        x = props.get('x', 'content')
        widths.append(max(3, int(available * scale)) if scale else None)
        weights.append(float(x[:-2]) if isinstance(x, str) and x.endswith('fr') else 1)
    remaining = available - sum(w or 0 for w in widths)
    flexible = [i for i, w in enumerate(widths) if w is None]
    if remaining < 3 * len(flexible):
        raise ValueError('Not enough terminal columns for this row')
    if not flexible and remaining < 0:
        raise ValueError('Relative widths exceed the available terminal columns')
    if flexible:
        extra = remaining - 3 * len(flexible)
        # Normalize before summation so even finite, very large fr values stay bounded.
        weight_scale = max(weights[i] for i in flexible)
        total = sum(weights[i] / weight_scale for i in flexible)
        fractions = {i: extra * (weights[i] / weight_scale) / total for i in flexible}
        for i in flexible:
            widths[i] = 3 + math.floor(fractions[i])
        for i in sorted(flexible, key=lambda i: (-(fractions[i] % 1), i))[:available-sum(widths)]:
            widths[i] += 1
    return widths


def gui_dump(root, columns=160):
    if not 20 <= columns <= 400:
        raise SduiError('dump-width', 'Columns must be between 20 and 400', root.span)
    budget = 0

    def render(node, width):
        nonlocal budget
        props = dict(node.layout)
        if props.get('visible') is False:
            return []
        if width < 3:
            raise SduiError('dump-space', 'Not enough columns for nested components', node.span)
        if node.kind == 'markdown':
            lines = [part for line in markdown_lines(node.text) for part in wrap_line(line, width)]
        elif node.kind == 'widget':
            args = dict(node.arguments)
            value = lambda key, default='': args[key].value if key in args else default
            if node.widget == 'button':
                content = '[ ' + value('label') + ' ]'
            elif node.widget == 'input':
                content = value('text') + ': [' + value('value') + ']'
            else:
                content = '[SVG plassholder: ' + value('label', node.path) + ']'
            lines = wrap_line(content, width)
        else:
            boxed = node.kind == 'frame' and node.variant == 'box'
            inner = width - 2 if boxed else width
            lines = []
            regions = dict(node.regions)
            if 'header' in regions:
                lines += render(regions['header'], inner)
            if 'body' in regions:
                lines += render(regions['body'], inner)
            for row in node.rows:
                visible = [item for item in row if dict(item.layout).get('visible') is not False]
                if not visible:
                    continue
                available = inner - (len(visible) - 1)
                try:
                    widths = horizontal_widths(visible, available)
                except (ValueError, OverflowError):
                    raise SduiError('dump-space', 'Row needs more columns; increase --columns', node.span)
                blocks = [render(item, w) for item, w in zip(visible, widths)]
                height = max((len(block) for block in blocks), default=0)
                for y in range(height):
                    lines.append(' '.join(fit(block[y] if y < len(block) else '', w)
                                          for block, w in zip(blocks, widths)))
            if 'footer' in regions:
                lines += render(regions['footer'], inner)
            if boxed:
                label = ' ' + node.path.rsplit('/', 1)[-1] + ' '
                if 'y' in props:
                    label += str(props['y']) + ' '
                border = '+' + fit(label, inner).replace(' ', '-') + '+'
                lines = [border] + ['|' + fit(line, inner) + '|' for line in lines] + ['+' + '-' * inner + '+']
        budget += sum(cell_width(line) for line in lines)
        if budget > MAX_CELLS:
            raise SduiError('dump-limit', 'Dump exceeds terminal cell budget', node.span)
        return lines
    heading = f'SDUI GUI dump | {safe(root.path)} | {columns} columns | structural preview'
    return heading + '\n' + '\n'.join(line.rstrip() for line in render(root, columns)) + '\n'
