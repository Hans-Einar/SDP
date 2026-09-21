"""Bounded recursive descent for grammar/sdui-0.2.ebnf."""
from .ast import (Argument, Connection, Definition, Document, LayoutRule, Literal,
                  ModuleRef, Node, Reference, Row, SduiError, Span)
from .lexer import lex

RESERVED = {'sdui', 'ref', 'true', 'false', 'null', 'setHandle'}


class Parser:
    def __init__(self, source):
        self.tokens, self.source = lex(source), source.encode('utf-8')
        self.index, self.nodes = 0, 0

    @property
    def token(self):
        return self.tokens[self.index]

    def ahead(self, kind):
        return self.tokens[min(self.index + 1, len(self.tokens) - 1)].kind == kind

    def take(self, kind, value=None):
        token = self.token
        if token.kind != kind or (value is not None and token.value != value):
            self.fail(f'Expected {value or kind}, found {token.value!r}')
        self.index += 1
        return token

    def accept(self, kind):
        return self.take(kind) if self.token.kind == kind else None

    def fail(self, message, code='syntax'):
        raise SduiError(code, message, self.token.span)

    def identifier(self):
        if self.token.value in RESERVED:
            self.fail('Reserved word cannot be used as an identifier')
        token = self.take('ID')
        if '-' in token.value:
            raise SduiError('syntax', 'Hyphens are only allowed in property names', token.span)
        return token

    def span(self, start):
        return Span(start.start, self.tokens[self.index - 1].span.end, start.line, start.column)

    def parse(self):
        start = self.take('ID', 'sdui').span
        version = self.take('NUMBER')
        if self.source[version.span.start:version.span.end] != b'0.2':
            raise SduiError('version', 'Only exact profile sdui 0.2 is supported', version.span)
        self.take(';')
        refs, definitions, connections = [], [], []
        while self.token.value == 'ref':
            begin = self.take('ID').span
            self.take(':')
            alias = self.identifier().value
            path = self.take('STRING').value
            self.take(';')
            refs.append(ModuleRef(alias, path, self.span(begin)))
        while self.token.kind == 'ID' and self.ahead('='):
            begin = self.token.span
            name = self.identifier().value
            self.take('=')
            root = self.node(1)
            self.take(';')
            definitions.append(Definition(name, root, self.span(begin)))
        if not definitions:
            self.fail('At least one UI definition is required')
        while self.token.kind != 'EOF':
            begin = self.token.span
            module = self.identifier().value
            self.take('.')
            obj = self.identifier().value
            self.take('.')
            self.take('ID', 'setHandle')
            self.take('(')
            definition = self.identifier().value
            path = []
            while self.accept('.'):
                path.append(self.identifier().value)
            if not path:
                self.fail('setHandle needs a component instance path')
            self.take(')')
            self.take(';')
            connections.append(Connection(module, obj, definition, tuple(path), self.span(begin)))
        end = self.take('EOF').span.end
        return Document('sdui/0.2', tuple(refs), tuple(definitions), tuple(connections),
                        Span(start.start, end, start.line, start.column))

    def node(self, depth, parent=None):
        if depth > 64:
            self.fail('Component nesting exceeds 64', 'depth-limit')
        self.nodes += 1
        if self.nodes > 2048:
            self.fail('More than 2048 components', 'node-limit')
        start, name, role = self.token.span, None, None
        if self.token.kind == 'ID' and self.ahead('='):
            name = self.identifier().value
            self.take('=')
            if parent == 'frame' and name in {'header', 'body', 'footer'}:
                role, name = name, None
        rows, args, text, target, widget, variant = (), (), None, None, None, None
        if self.token.kind in {'[', '<'}:
            opening = self.token.kind
            self.take(opening)
            kind, close = ('frame', ']') if opening == '[' else ('group', '>')
            rows = self.rows(close, depth, kind)
            self.take(close)
        elif self.token.kind == 'STRING':
            kind, text = 'markdown', self.literal()
        elif self.token.kind == 'ID':
            target = self.identifier().value
            if self.accept('('):
                kind, widget, target = 'widget', target, None
                args = self.arguments()
                self.take(')')
            else:
                kind = 'use'
        else:
            self.fail('Expected frame, group, Markdown, widget call or component reference')
        if self.accept('*'):
            variant = self.identifier().value
        layout = self.layout() if self.token.kind == '{' else ()
        return Node(kind, name, role, rows, widget, args, text, target, variant, layout, self.span(start))

    def rows(self, close, depth, parent):
        rows = []
        if self.token.kind == close:
            return ()
        while True:
            begin = self.token.span
            items = [self.node(depth + 1, parent)]
            while self.accept(','):
                items.append(self.node(depth + 1, parent))
            rows.append(Row(tuple(items), self.span(begin)))
            if not self.accept(';') or self.token.kind == close:
                break
        return tuple(rows)

    def arguments(self):
        args = []
        if self.token.kind != ')':
            while True:
                start, name = self.token.span, None
                if self.token.kind == 'ID' and self.ahead('='):
                    name = self.identifier().value
                    self.take('=')
                args.append(Argument(name, self.value(), self.span(start)))
                if len(args) > 32:
                    self.fail('Widget has more than 32 arguments', 'argument-limit')
                if not self.accept(','):
                    break
        return tuple(args)

    def value(self):
        if self.token.kind == 'ID' and self.token.value not in {'true', 'false', 'null'}:
            start = self.token.span
            module = self.identifier().value
            self.take('.')
            obj = self.identifier().value
            self.take('.')
            self.take('@')
            member = self.identifier().value
            return Reference(module, obj, member, self.span(start))
        return self.literal()

    def literal(self):
        token = self.token
        if token.kind in {'STRING', 'NUMBER'}:
            self.take(token.kind)
            return Literal('string' if token.kind == 'STRING' else 'number', token.value, token.span)
        if token.kind == 'ID' and token.value in {'true', 'false', 'null'}:
            self.take('ID')
            return Literal('null' if token.value == 'null' else 'boolean',
                           {'true': True, 'false': False, 'null': None}[token.value], token.span)
        self.fail('Expected literal')

    def layout_value(self):
        start = self.token.span
        if self.accept('('):
            values = [self.take('NUMBER').value]
            while self.accept(','):
                values.append(self.take('NUMBER').value)
                if len(values) > 4:
                    self.fail('Padding tuple has at most four values')
            self.take(')')
            return Literal('tuple', tuple(values), self.span(start))
        if self.token.kind == 'ID' and self.token.value not in {'true', 'false', 'null'}:
            return Literal('enum', self.take('ID').value, start)
        value = self.literal()
        if value.kind == 'number' and self.token.value == 'fr':
            self.take('ID')
            return Literal('fr', value.value, self.span(start))
        return value

    def layout(self):
        self.take('{')
        rules = []
        if self.token.kind != '}':
            while True:
                start = self.token.span
                if self.token.kind == 'ID' and self.ahead('='):
                    name = self.take('ID').value
                    self.take('=')
                    value = self.layout_value()
                elif self.token.kind == 'NUMBER' and self.ahead(':'):
                    x = self.take('NUMBER').value
                    self.take(':')
                    y = self.take('NUMBER').value
                    name, value = 'ratio', Literal('ratio', (x, y), self.span(start))
                else:
                    if self.token.kind in {',', '}', 'EOF'}:
                        self.fail('Expected formatting rule')
                    while self.token.kind not in {',', '}', 'EOF'}:
                        self.index += 1
                    name = 'arrow'
                    raw = self.source[start.start:self.tokens[self.index-1].span.end].decode()
                    value = Literal('arrow', ''.join(raw.split()), self.span(start))
                spelling = self.source[start.start:self.tokens[self.index-1].span.end].decode()
                rules.append(LayoutRule(name, value, spelling, self.span(start)))
                if len(rules) > 32:
                    self.fail('More than 32 formatting rules', 'layout-limit')
                if not self.accept(','):
                    break
        self.take('}')
        return tuple(rules)


def parse(source: str) -> Document:
    return Parser(source).parse()
