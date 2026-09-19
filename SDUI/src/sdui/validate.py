"""Local SDUI profile checks. Never resolve or execute an external SDL module."""
from .ast import Box, Document, Literal, Reference, SduiError

PROPERTIES = {'title', 'axis', 'weight', 'min', 'max', 'gap', 'padding'}
# First positional argument is syntactic sugar for the named property below.
WIDGETS = {
    'text': ('text', {'text': 'string'}, {'text'}),
    'button': ('label', {'label': 'string', 'callback': 'reference'}, {'label'}),
    'input': ('text', {'text': 'string', 'value': 'string', 'callback': 'reference'}, {'text'}),
    'svg': ('source', {'source': 'reference', 'label': 'string'}, {'source'}),
}


def validate(document: Document) -> None:
    def fail(code, message, node):
        raise SduiError(code, message, node.span)

    modules = {}
    for ref in document.references:
        if ref.alias in modules:
            fail('duplicate-module', f'Duplicate module alias {ref.alias}', ref)
        if not ref.path or any(ord(c) < 32 or ord(c) == 127 for c in ref.path):
            fail('module-path', 'Module path must be nonempty and contain no control characters', ref)
        modules[ref.alias] = ref
    definitions = {}
    for definition in document.definitions:
        if definition.name in definitions or definition.name in modules:
            fail('duplicate-definition', f'Ambiguous definition name {definition.name}', definition)
        definitions[definition.name] = definition
    registries = {}
    for definition in document.definitions:
        names = {}

        def register(node):
            if node.name is not None:
                if node.name in names:
                    fail('duplicate-node', f'Duplicate name {node.name} in {definition.name}', node)
                names[node.name] = node

        def widget(node):
            register(node)
            if node.kind not in WIDGETS:
                fail('widget-kind', f'Unsupported widget kind {node.kind}', node)
            positional, fields, required = WIDGETS[node.kind]
            args, named_seen = {}, False
            for arg in node.arguments:
                if arg.name is None and (args or named_seen):
                    fail('argument-order', 'Only one positional argument, before named arguments', arg)
                key = arg.name or positional
                named_seen |= arg.name is not None
                if key in args:
                    fail('duplicate-argument', f'Duplicate widget property {key}', arg)
                if key not in fields:
                    fail('widget-property', f'Unknown {node.kind} property {key}', arg)
                expected = fields[key]
                value = arg.value
                if isinstance(value, Reference):
                    if expected != 'reference':
                        fail('argument-type', f'{key} requires {expected}', arg)
                    if value.module not in modules:
                        fail('unknown-module', f'Unknown module alias {value.module}', value)
                elif expected == 'reference' or value.kind != expected:
                    fail('argument-type', f'{key} requires {expected}', arg)
                args[key] = value
            if required - args.keys():
                fail('missing-argument', f'Missing widget properties: {sorted(required - args.keys())}', node)

        def box(node: Box):
            register(node)
            props = {}
            for prop in node.properties:
                if prop.name in props:
                    fail('duplicate-property', f'Duplicate box property {prop.name}', prop)
                if prop.name not in PROPERTIES:
                    fail('box-property', f'Unknown box property {prop.name}; use the assigned box name, not handle', prop)
                value = prop.value
                if prop.name in {'title', 'axis'}:
                    if value.kind != 'string':
                        fail('property-type', f'{prop.name} must be a string', prop)
                    if prop.name == 'axis' and value.value not in {'row', 'column'}:
                        fail('axis', 'axis must be row or column', prop)
                elif value.kind != 'number' or value.value < 0 or (prop.name == 'weight' and value.value == 0):
                    fail('property-type', f'{prop.name} requires a nonnegative number (weight > 0)', prop)
                props[prop.name] = value.value
            if 'min' in props and 'max' in props and props['min'] > props['max']:
                fail('size-range', 'min must not exceed max', node)
            if bool(node.children) == (node.content is not None):
                fail('box-body', 'A box requires either child boxes or one widget content block', node)
            if node.content is not None and 'axis' in props:
                fail('axis-scope', 'axis applies to child boxes; widget rows use semicolons', node)
            for child in node.children:
                box(child)
            if node.content:
                for row in node.content.rows:
                    for item in row.widgets:
                        widget(item)

        box(definition.root)
        registries[definition.name] = names
    objects, targets = set(), set()
    for connection in document.connections:
        if connection.module not in modules:
            fail('unknown-module', f'Unknown module alias {connection.module}', connection)
        if connection.definition not in registries:
            fail('unknown-definition', f'Unknown UI definition {connection.definition}', connection)
        node = registries[connection.definition].get(connection.widget)
        if node is None or isinstance(node, Box):
            fail('unknown-widget', f'Expected widget handle {connection.definition}.{connection.widget}', connection)
        obj = (connection.module, connection.object)
        target = (connection.definition, connection.widget)
        if obj in objects or target in targets:
            fail('duplicate-connection', 'Profile 0.1 permits one SDL object per widget and one widget per SDL object', connection)
        objects.add(obj)
        targets.add(target)
