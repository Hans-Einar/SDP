"""Resolve component uses into bounded instance trees, independent of GUI toolkits."""
from dataclasses import dataclass, replace
from .ast import Literal, Reference, Span
from .formatting import check_combination, fail, formatting
from .validate import widget_arguments


@dataclass(frozen=True)
class Instance:
    kind: str
    path: str
    variant: str | None
    widget: str | None
    text: str | None
    arguments: tuple[tuple[str, Literal | Reference], ...]
    layout: tuple[tuple[str, object], ...]
    rows: tuple[tuple['Instance', ...], ...]
    regions: tuple[tuple[str, 'Instance'], ...]
    span: Span


def resolve_all(document):
    definitions = {definition.name: definition.root for definition in document.definitions}
    count = 0

    def expand(node, path, depth):
        nonlocal count
        count += 1
        if count > 8192:
            fail('expansion-limit', 'More than 8192 expanded components', node)
        if depth > 64:
            fail('depth-limit', 'Expanded component depth exceeds 64', node)
        if node.kind == 'use':
            if node.target not in definitions:
                fail('unknown-definition', f'Unknown component {node.target}', node)
            instance = expand(definitions[node.target], path, depth + 1)
            # Validate suffix against the effective component, not the use-site node.
            effective = replace(node, kind=instance.kind,
                                rows=instance.rows)
            overlay = formatting(effective)
            props = dict(instance.layout)
            props.update(overlay)
            check_combination(props, instance.kind, effective)
            check_wrap_children(props, instance.rows, node)
            return replace(instance, layout=tuple(props.items()))
        props = formatting(node)
        rows, regions = [], []
        for r, row in enumerate(node.rows):
            items = []
            for c, child in enumerate(row.items):
                child_name = child.role or child.name or child.target
                if child_name:
                    child_path = path + '/' + child_name
                else:
                    child_path = f'{path}/$r{r}c{c}'
                instance = expand(child, child_path, depth + 1)
                if node.kind == 'group' and instance.kind == 'frame':
                    fail('group-content', 'Frame reference cannot be placed in a widget group', child)
                if child.role:
                    regions.append((child.role, instance))
                else:
                    items.append(instance)
            if items:
                rows.append(tuple(items))
        check_wrap_children(props, rows, node)
        return Instance(node.kind, path, 'box' if node.variant else None,
                        node.widget, node.text.value if node.text else None,
                        tuple(widget_arguments(node).items()) if node.kind == 'widget' else (),
                        tuple(props.items()), tuple(rows), tuple(regions), node.span)
    return {name: expand(node, name, 1) for name, node in definitions.items()}


def normalize(document):
    from .validate import validate
    validate(document)
    return resolve_all(document)


def check_wrap_children(props, rows, node):
    if props.get('wrap') == 'wrap':
        for row in rows:
            for child in row:
                x = dict(child.layout).get('x', '')
                if x == 'fill' or str(x).endswith('fr'):
                    fail('wrap-layout', 'Wrapped children cannot use horizontal fill/fr', node)
