"""Local 0.2 profile validation; never load or execute referenced SDL sources."""
from .ast import Document, Reference
from .formatting import fail, formatting

WIDGETS = {
    'button': ('label', {'label': 'string', 'callback': 'reference'}, {'label'}),
    'input': ('text', {'text': 'string', 'value': 'string', 'callback': 'reference'}, {'text'}),
    'svg': ('source', {'source': 'reference', 'label': 'string'}, {'source'}),
}


def widget_arguments(node):
    positional = WIDGETS[node.widget][0]
    return {arg.name or positional: arg.value for arg in node.arguments}


def validate(document: Document) -> None:
    modules, definitions = {}, {}
    for ref in document.references:
        if ref.alias in modules:
            fail('duplicate-module', f'Duplicate module alias {ref.alias}', ref)
        if not ref.path or any(ord(c) < 32 or ord(c) == 127 for c in ref.path):
            fail('module-path', 'Module path must be nonempty without control characters', ref)
        modules[ref.alias] = ref
    for definition in document.definitions:
        if definition.name in definitions or definition.name in modules:
            fail('duplicate-definition', f'Ambiguous definition name {definition.name}', definition)
        definitions[definition.name] = definition
    dependencies = {name: set() for name in definitions}
    for definition in document.definitions:
        names = set()

        def visit(node):
            # Unaliased references take their definition name as the instance name.
            name = node.name or (node.target if node.kind == 'use' else None)
            if name:
                if name in names:
                    fail('duplicate-node', f'Duplicate name {name} in {definition.name}', node)
                names.add(name)
            if node.kind == 'use':
                if node.target not in definitions:
                    fail('unknown-definition', f'Unknown component {node.target}', node)
                dependencies[definition.name].add(node.target)
                # Effective kind/layout checked on the resolved instance below.
            else:
                formatting(node)
            if node.variant and (node.variant not in {'b', 'box'} or node.kind != 'frame'):
                fail('variant', 'Only frame variants box and b are supported', node)
            if node.kind == 'widget':
                validate_widget(node, modules)
            roles, body = set(), False
            for row in node.rows:
                for child in row.items:
                    if node.kind == 'group' and child.kind == 'frame':
                        fail('group-content', 'Frames belong in frame bodies, not widget groups', child)
                    if child.role:
                        if child.role in roles:
                            fail('duplicate-region', f'Duplicate {child.role}', child)
                        roles.add(child.role)
                    else:
                        body = True
                    visit(child)
            if 'body' in roles and body:
                fail('body-conflict', 'Explicit body cannot be mixed with implicit body', node)
        visit(definition.root)
    done = set()

    def cycle(name, stack):
        if name in stack:
            fail('reference-cycle', f'Recursive component reference {name}', definitions[name])
        if len(stack) >= 64:
            fail('depth-limit', 'Reference chain exceeds 64', definitions[name])
        if name not in done:
            for target in sorted(dependencies[name]):
                cycle(target, stack + (name,))
            done.add(name)
    for name in definitions:
        cycle(name, ())
    # Expansion has a separate global budget; references cannot cause exponential work.
    from .normalize import resolve_all
    roots = resolve_all(document)
    widgets = {}

    def collect(instance):
        if instance.kind == 'widget' and not instance.path.rsplit('/', 1)[-1].startswith('$'):
            public_path = '/'.join(part for part in instance.path.split('/') if not part.startswith('$'))
            if public_path in widgets:
                fail('duplicate-instance', f'Ambiguous widget path {public_path}', instance)
            widgets[public_path] = instance
        for row in instance.rows:
            for child in row:
                collect(child)
        for _, region in instance.regions:
            collect(region)
    for root in roots.values():
        collect(root)
    objects, targets = set(), set()
    for connection in document.connections:
        if connection.module not in modules:
            fail('unknown-module', f'Unknown module {connection.module}', connection)
        if connection.definition not in roots:
            fail('unknown-definition', f'Unknown definition {connection.definition}', connection)
        target = '/'.join((connection.definition,) + connection.path)
        if target not in widgets:
            fail('unknown-widget', f'Unknown widget instance {target}', connection)
        obj = (connection.module, connection.object)
        if obj in objects or target in targets:
            fail('duplicate-connection', 'One connection per SDL object and widget', connection)
        objects.add(obj)
        targets.add(target)


def validate_widget(node, modules):
    if node.widget not in WIDGETS:
        fail('widget-kind', f'Unsupported widget kind {node.widget}; quoted text is Markdown', node)
    positional, fields, required = WIDGETS[node.widget]
    args, named_seen = {}, False
    for arg in node.arguments:
        if arg.name is None and (args or named_seen):
            fail('argument-order', 'Only one positional argument before named arguments', arg)
        key = arg.name or positional
        named_seen |= arg.name is not None
        if key in args:
            fail('duplicate-argument', f'Duplicate widget property {key}', arg)
        if key not in fields:
            fail('widget-property', f'Unknown {node.widget} property {key}', arg)
        value, expected = arg.value, fields[key]
        if isinstance(value, Reference):
            if expected != 'reference':
                fail('argument-type', f'{key} requires {expected}', arg)
            if value.module not in modules:
                fail('unknown-module', f'Unknown module {value.module}', value)
        elif expected == 'reference' or value.kind != expected:
            fail('argument-type', f'{key} requires {expected}', arg)
        args[key] = value
    if node.name is None and 'callback' in args:
        fail('widget-name', 'A callback requires a named widget', node)
    if required - args.keys():
        fail('missing-argument', f'Missing properties {sorted(required - args.keys())}', node)
