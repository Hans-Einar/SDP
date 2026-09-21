"""Static Markdown document from the same instance tree as the terminal dump."""
import re
from .ast import SduiError
from .dump import FENCE, gui_dump, markdown_lines, safe

MAX_OUTPUT_BYTES = 2_000_000


def code_span(value):
    """Literal widget labels must not become links, HTML or Markdown controls."""
    text = safe(str(value)).replace('\n', ' ')
    ticks = '`' * (1 + max((len(run) for run in re.findall(r'`+', text)), default=0))
    return ticks + ' ' + text + ' ' + ticks


def quoted(lines):
    return ['> ' + line if line else '>' for line in lines]


def content_lines(source):
    """Omit diagrams and close unfinished code fences at the widget boundary."""
    lines = [safe(line) for line in markdown_lines(source)]
    fence = None
    for line in lines:
        if fence:
            if re.fullmatch(r' {0,3}' + re.escape(fence[0]) +
                            '{' + str(fence[1]) + r',}[ \t]*', line):
                fence = None
        else:
            match = FENCE.match(line)
            if match and not (match[1][0] == '`' and '`' in match[2]):
                fence = (match[1][0], len(match[1]))
    if fence:
        lines.append(fence[0] * fence[1])
    return lines


def markdown_dump(root, columns=160):
    """Portable Markdown: fenced layout overview, then rendered source content.

    Blockquotes preserve component nesting without HTML/CSS or Mermaid support.
    Row captions preserve order; Markdown itself does not implement SDUI geometry.
    """
    overview = gui_dump(root, columns).rstrip('\n')
    fence = '`' * max(3, 1 + max((len(run) for run in re.findall(r'`+', overview)), default=0))
    budget = 0

    def render(node):
        nonlocal budget
        if dict(node.layout).get('visible') is False:
            return []
        if node.kind == 'markdown':
            lines = quoted(content_lines(node.text))
        elif node.kind == 'widget':
            args = dict(node.arguments)
            value = lambda key, default='': args[key].value if key in args else default
            if node.widget == 'button':
                text = '**Knapp:** ' + code_span(value('label'))
            elif node.widget == 'input':
                text = '**Inndata:** ' + code_span(value('text')) + ' — ' + code_span(value('value'))
            else:
                text = '**SVG-plassholder:** ' + code_span(value('label', node.path))
            lines = [text]
        else:
            kind = 'BoxUI-frame' if node.variant == 'box' else 'Frame' if node.kind == 'frame' else 'Gruppe'
            lines = ['**' + kind + ':** ' + code_span(node.path), '']
            regions = dict(node.regions)
            for role in ('header', 'body'):
                if role in regions and dict(regions[role].layout).get('visible') is not False:
                    lines += ['**' + role + '**', ''] + render(regions[role]) + ['']
            visible_rows = [[child for child in row if dict(child.layout).get('visible') is not False]
                            for row in node.rows]
            visible_rows = [row for row in visible_rows if row]
            for index, row in enumerate(visible_rows, 1):
                if len(visible_rows) > 1 or len(row) > 1:
                    noun = 'komponent' if len(row) == 1 else 'komponenter'
                    lines += [f'**Rad {index} · {len(row)} {noun} fra venstre mot høyre**', '']
                for child in row:
                    lines += render(child) + ['', '---', '']
            if 'footer' in regions and dict(regions['footer'].layout).get('visible') is not False:
                lines += ['**footer**', ''] + render(regions['footer']) + ['']
            lines = quoted(lines)
        budget += sum(len(line.encode('utf-8')) + 1 for line in lines)
        if budget > MAX_OUTPUT_BYTES:
            raise SduiError('markdown-limit', 'Markdown dump exceeds output budget', node.span)
        return lines

    lines = [
        '# SDUI — ' + code_span(root.path), '',
        'Statisk GUI-dump. Knapper og felt er tekstetiketter; ingen callbacks kjøres.', '',
        '## Layoutoversikt', '',
        'Rad-/kolonnestruktur i terminalceller. Høydene følger innholdet; dette er ikke målt GUI-geometri.', '',
        fence + 'text', overview, fence, '',
        '## Innhold', '',
        'Markdown gjengis som innhold. Nestede sitatblokker viser grupper og frames. '
        'Horisontale søsken står i leserekkefølge her; plasseringen vises i oversikten. '
        'Mermaid-diagrammer er utelatt.', '',
    ] + render(root)
    text = '\n'.join(lines).rstrip() + '\n'
    if len(text.encode('utf-8')) > MAX_OUTPUT_BYTES:
        raise SduiError('markdown-limit', 'Markdown dump exceeds output budget', root.span)
    return text
