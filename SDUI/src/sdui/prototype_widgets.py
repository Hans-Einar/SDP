"""Small preview-only button/input backend. No runtime, domain dispatch or GUI imports."""
from dataclasses import dataclass
from html import escape
import math
from .ast import SduiError
from .dump import safe

PALETTE = {
    'ink': '#20302a', 'muted': '#65776d', 'line': '#9aaa9d',
    'paper': '#ffffff', 'face': '#f3f6f0', 'accent': '#426e54',
    'pressed': '#d6e4d5', 'disabled': '#ecefe9', 'focus': '#3975b0',
}


@dataclass(frozen=True)
class PrototypeWidget:
    path: str
    kind: str
    label: str
    value: str
    enabled: bool


def widget_from_instance(node):
    if node.kind != 'widget' or node.widget not in {'button', 'input'}:
        raise SduiError('prototype-widget', 'Prototype supports only button/input', node.span)
    args = dict(node.arguments)
    label = args['label' if node.widget == 'button' else 'text'].value
    value = args['value'].value if 'value' in args else ''
    # References/callbacks are deliberately not transferred to the preview backend.
    return PrototypeWidget(node.path, node.widget, label, value,
                           dict(node.layout).get('enabled', True))


def prototype_widgets(root):
    """Index visible, supported widgets by instance path; never read an SDL source."""
    result = {}

    def visit(node, inherited_enabled=True):
        props = dict(node.layout)
        if props.get('visible') is False:
            return
        enabled = inherited_enabled and props.get('enabled', True)
        if node.kind == 'widget' and node.widget in {'button', 'input'}:
            widget = widget_from_instance(node)
            result[widget.path] = PrototypeWidget(widget.path, widget.kind, widget.label,
                                                   widget.value, enabled)
        for row in node.rows:
            for child in row:
                visit(child, enabled)
        for _, child in node.regions:
            visit(child, enabled)
    visit(root)
    return result


def svg_widget(widget, x, y, width, *, pressed=False, focused=False):
    """SVG logical units; fixed preview metrics, not native SDUI font/layout rules."""
    if not all(math.isfinite(n) for n in (x, y, width)) or width < 40:
        raise ValueError('Widget geometry requires finite coordinates and width >= 40')
    disabled = not widget.enabled
    ink = PALETTE['muted' if disabled else 'ink']
    label, value = escape(safe(widget.label)), escape(safe(widget.value))
    path = escape(widget.path, quote=True)
    font = 'font-family="sans-serif" font-size="14"'
    attrs = f'data-widget="{widget.kind}" data-path="{path}"'
    # Nested SVG establishes a clipping viewport without requiring global clip IDs.
    if widget.kind == 'button':
        height = 36
        down = pressed and not disabled
        fill = PALETTE['disabled' if disabled else 'pressed' if down else 'face']
        dy = 2 if down else 0
        lines = [f'<svg x="{x}" y="{y}" width="{width}" height="40" {attrs}>',
                 f'<title>{label}</title>',
                 f'<rect x="1" y="4" width="{width-2}" height="33" rx="5" fill="{PALETTE["line"]}"/>',
                 f'<rect x="1" y="{1+dy}" width="{width-2}" height="33" rx="5" fill="{fill}" stroke="{PALETTE["line"]}"/>',
                 f'<text x="{width/2}" y="{23+dy}" text-anchor="middle" fill="{ink}" {font}>{label}</text>']
    elif widget.kind == 'input':
        height = 66
        fill = PALETTE['disabled' if disabled else 'paper']
        stroke = PALETTE['focus' if focused and not disabled else 'line']
        lines = [f'<svg x="{x}" y="{y}" width="{width}" height="{height}" {attrs}>',
                 f'<title>{label}: {value}</title>',
                 f'<text x="0" y="15" fill="{ink}" {font}>{label}</text>',
                 f'<rect x="1" y="24" width="{width-2}" height="36" rx="4" fill="{fill}" stroke="{stroke}" stroke-width="{2 if focused and not disabled else 1}"/>',
                 f'<svg x="10" y="25" width="{width-20}" height="34"><text x="0" y="23" fill="{ink}" {font}>{value}</text></svg>']
    else:
        raise ValueError('Unsupported prototype widget')
    if focused and not disabled and widget.kind == 'button':
        lines.append(f'<rect x="3" y="3" width="{width-6}" height="31" rx="4" fill="none" stroke="{PALETTE["focus"]}" stroke-dasharray="3 2"/>')
    lines.append('</svg>')
    return ''.join(lines), height


def html_widget(widget, index):
    """Editable HTML stand-in; no form submission, link or domain callback."""
    label, value = escape(widget.label), escape(widget.value, quote=True)
    ident = f'widget-{index}'
    disabled = '' if widget.enabled else ' disabled'
    if widget.kind == 'button':
        return f'<button type="button" class="sdui-button"{disabled}>{label}</button>'
    if widget.kind == 'input':
        return (f'<label class="sdui-field" for="{ident}"><span>{label}</span>'
                f'<input id="{ident}" type="text" value="{value}" autocomplete="off"{disabled}>'
                f'<span class="print-value" aria-hidden="true">{escape(widget.value)}</span></label>')
    raise ValueError('Unsupported prototype widget')
