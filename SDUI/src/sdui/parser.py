"""Recursive descent implementation of grammar/sdui-0.1.ebnf."""
from .ast import (Argument, Box, Connection, Content, Definition, Document, Literal,
                  ModuleRef, Property, Reference, Row, SduiError, Span, Widget)
from .lexer import lex

RESERVED = {'sdui', 'ref', 'true', 'false', 'null', 'setHandle'}


class Parser:
    def __init__(self, source: str):
        self.tokens = lex(source)
        self.source = source.encode('utf-8')
        self.index = 0
        self.nodes = 0

    @property
    def token(self):
        return self.tokens[self.index]

    def take(self, kind, value=None):
        token = self.token
        if token.kind != kind or (value is not None and token.value != value):
            self.fail(f'Expected {value or kind}, found {token.value!r}')
        self.index += 1
        return token

    def accept(self, kind):
        if self.token.kind == kind:
            return self.take(kind)
        return None

    def fail(self, message, code='syntax'):
        raise SduiError(code, message, self.token.span)

    def identifier(self):
        if self.token.value in RESERVED:
            self.fail('Reserved word cannot be used as an identifier')
        return self.take('ID')

    def span(self, start):
        return Span(start.start, self.tokens[self.index - 1].span.end, start.line, start.column)

    def count(self):
        self.nodes += 1
        if self.nodes > 2048:
            self.fail('More than 2048 boxes/widgets', 'node-limit')

    def parse(self):
        start = self.take('ID', 'sdui').span
        version = self.take('NUMBER')
        if self.source[version.span.start:version.span.end] != b'0.1':
            raise SduiError('version', 'Only exact profile sdui 0.1 is supported', version.span)
        self.take(';')
        refs, definitions, connections = [], [], []
        while self.token.value == 'ref':
            begin = self.take('ID', 'ref').span
            self.take(':')
            alias = self.identifier().value
            path = self.take('STRING').value
            self.take(';')
            refs.append(ModuleRef(alias, path, self.span(begin)))
        while self.token.kind == 'ID' and self.tokens[self.index + 1].kind == '=':
            begin = self.token.span
            name = self.identifier().value
            self.take('=')
            self.take('{')
            root_name = self.identifier().value
            self.take('=')
            root = self.box(root_name, 1)
            self.take('}')
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
            self.take('.')
            widget = self.identifier().value
            self.take(')')
            self.take(';')
            connections.append(Connection(module, obj, definition, widget, self.span(begin)))
        # Span includes all trailing whitespace/comments, not an implicit AST node.
        end = self.take('EOF').span.end
        return Document('sdui/0.1', tuple(refs), tuple(definitions), tuple(connections),
                        Span(start.start, end, start.line, start.column))

    def box(self, name, depth):
        if depth > 64:
            self.fail('Box nesting exceeds 64', 'depth-limit')
        self.count()
        start = self.take('[').span
        properties, children, content = [], [], None
        body_started = False
        while self.token.kind != ']':
            if self.token.kind == '[':
                body_started = True
                children.append(self.box(None, depth + 1))
            elif self.token.kind == '{':
                body_started = True
                if content is not None:
                    self.fail('Only one widget content block per box')
                content = self.content()
            else:
                begin = self.token.span
                key = self.identifier().value
                self.take('=')
                if self.token.kind == '[':
                    body_started = True
                    children.append(self.box(key, depth + 1))
                else:
                    if body_started:
                        self.fail('Box properties must precede child boxes/content')
                    properties.append(Property(key, self.literal(), self.span(begin)))
            if not self.accept(','):
                break
            if self.token.kind == ']':
                self.fail('Trailing box comma is unsupported')
        self.take(']')
        return Box(name, tuple(properties), tuple(children), content, self.span(start))

    def content(self):
        start = self.take('{').span
        rows = []
        if self.token.kind == '}':
            self.fail('Widget content must not be empty')
        while True:
            begin = self.token.span
            widgets = [self.widget()]
            while self.accept(','):
                widgets.append(self.widget())
            rows.append(Row(tuple(widgets), self.span(begin)))
            if not self.accept(';') or self.token.kind == '}':
                break
        self.take('}')
        return Content(tuple(rows), self.span(start))

    def widget(self):
        self.count()
        start = self.token.span
        name = self.identifier().value
        self.take('=')
        kind = self.identifier().value
        self.take('(')
        arguments = []
        if self.token.kind != ')':
            while True:
                begin = self.token.span
                key = None
                if self.token.kind == 'ID' and self.tokens[self.index + 1].kind == '=':
                    key = self.identifier().value
                    self.take('=')
                arguments.append(Argument(key, self.value(), self.span(begin)))
                if len(arguments) > 32:
                    self.fail('Widget has more than 32 arguments', 'argument-limit')
                if not self.accept(','):
                    break
        self.take(')')
        return Widget(name, kind, tuple(arguments), self.span(start))

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
        if token.kind == 'STRING':
            self.take('STRING')
            return Literal('string', token.value, token.span)
        if token.kind == 'NUMBER':
            self.take('NUMBER')
            return Literal('number', token.value, token.span)
        if token.kind == 'ID' and token.value in {'true', 'false', 'null'}:
            self.take('ID')
            return Literal('null' if token.value == 'null' else 'boolean',
                           {'true': True, 'false': False, 'null': None}[token.value], token.span)
        self.fail('Expected literal')


def parse(source: str) -> Document:
    return Parser(source).parse()
