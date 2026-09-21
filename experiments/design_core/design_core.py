"""Experimental parser and structural validator for design-core 0.2.

The language definition in docs/Design-Language-Definition.md is authoritative.
No input is executed. Formatting writes to stdout, never to the input file.
"""

import argparse
import bisect
from collections import namedtuple
import json
from pathlib import Path
import re
import sys


Span = namedtuple("Span", "start end line column end_line end_column")
Token = namedtuple("Token", "text span")
Identifier = namedtuple("Identifier", "name span")
Header = namedtuple("Header", "language version span")
Declaration = namedtuple("Declaration", "kind name span")
Relation = namedtuple("Relation", "subject verb object span")
Dependency = namedtuple("Dependency", "subject interface mode span")
Allocation = namedtuple("Allocation", "subject container mode span")
PropertyAssignment = namedtuple("PropertyAssignment", "subject property value span")
Model = namedtuple("Model", "header declarations statements span")
Diagnostic = namedtuple("Diagnostic", "code message span")

KINDS = frozenset(("unit", "container", "functionality", "capability",
                   "interface", "activity", "mode", "actor", "usecase", "feature"))
SIGNATURES = {
    "contains": ("unit", "unit"),
    "owns": ("unit", "functionality"),
    "realizes": ("functionality", "capability"),
    "provides": ("unit", "capability"),
    "consumes": ("unit", "interface"),
    "refines": ("activity", "activity"),
    "pursues": ("actor", "usecase"),
    "supports": ("feature", "usecase"),
    "contributes-to": ("functionality", ("feature", "usecase")),
}
PROPERTIES = {
    "state-retention": frozenset(("stateful", "stateless")),
    "repeatability": frozenset(("deterministic", "nondeterministic")),
}
VALUES = frozenset().union(*PROPERTIES.values())
TOKEN = re.compile(r"[A-Za-z][A-Za-z0-9-]*|[0-9]+(?:\.[0-9]+)*|[.=]")
NAME = re.compile(r"[A-Z][A-Za-z0-9]*\Z")


class ParseError(ValueError):
    def __init__(self, diagnostic):
        super().__init__(diagnostic.message)
        self.diagnostic = diagnostic


class ValidationError(ValueError):
    def __init__(self, diagnostics):
        super().__init__("Model has semantic errors")
        self.diagnostics = diagnostics


class Source:
    def __init__(self, text):
        self.text = text
        self.lines = [0] + [i + 1 for i, c in enumerate(text) if c == "\n"]

    def span(self, start, end):
        first = bisect.bisect_right(self.lines, start) - 1
        last = bisect.bisect_right(self.lines, end) - 1
        return Span(start, end, first + 1, start - self.lines[first] + 1,
                    last + 1, end - self.lines[last] + 1)

    def tokens(self):
        result = []
        offset = 0
        while offset < len(self.text):
            if self.text[offset] in " \t\r\n":
                offset += 1
                continue
            match = TOKEN.match(self.text, offset)
            if not match:
                raise ParseError(Diagnostic(
                    "UNSUPPORTED_SYNTAX", "Unsupported character {!r}".format(
                        self.text[offset]), self.span(offset, offset + 1)))
            result.append(Token(match.group(), self.span(offset, match.end())))
            offset = match.end()
        result.append(Token("", self.span(offset, offset)))
        return result


class Parser:
    def __init__(self, text):
        self.source = Source(text)
        self.tokens = self.source.tokens()
        self.index = 0

    @property
    def current(self):
        return self.tokens[self.index]

    def fail(self, message, code="UNSUPPORTED_SYNTAX"):
        raise ParseError(Diagnostic(code, message, self.current.span))

    def take(self):
        token = self.current
        self.index += 1
        return token

    def expect(self, text):
        if self.current.text != text:
            self.fail("Expected {!r}, received {!r}".format(
                text, self.current.text or "end of input"))
        return self.take()

    def choice(self, choices):
        if self.current.text not in choices:
            self.fail("Expected one of {}; received {!r}".format(
                ", ".join(sorted(choices)), self.current.text or "end of input"))
        return self.take().text

    def identifier(self):
        if not NAME.fullmatch(self.current.text):
            self.fail("Expected an ASCII identifier matching [A-Z][A-Za-z0-9]*")
        token = self.take()
        return Identifier(token.text, token.span)

    def finish(self, start):
        end = self.expect(".").span.end
        return self.source.span(start, end)

    def model(self):
        start = self.expect("language").span.start
        self.expect("design-core")
        self.expect("version")
        if self.current.text != "0.2":
            self.fail("Only design-core version 0.2 is supported", "UNSUPPORTED_VERSION")
        self.take()
        header = Header("design-core", "0.2", self.finish(start))
        declarations = []
        while self.current.text in KINDS:
            token = self.take()
            name = self.identifier()
            declarations.append(Declaration(token.text, name, self.finish(token.span.start)))
        statements = []
        while self.current.text:
            subject = self.identifier()
            start = subject.span.start
            verb = self.choice(set(SIGNATURES) | {"requires", "allocated-to", "has"})
            if verb in ("requires", "allocated-to"):
                target = self.identifier()
                self.expect("in")
                self.expect("mode")
                mode = self.identifier()
                node = Dependency if verb == "requires" else Allocation
                statement = node(subject, target, mode, self.finish(start))
            elif verb == "has":
                prop = self.choice(PROPERTIES)
                self.expect("=")
                value = self.choice(VALUES)
                statement = PropertyAssignment(subject, prop, value, self.finish(start))
            else:
                obj = self.identifier()
                statement = Relation(subject, verb, obj, self.finish(start))
            statements.append(statement)
        return Model(header, tuple(declarations), tuple(statements),
                     self.source.span(header.span.start, self.tokens[self.index - 1].span.end))


def parse(text):
    """Build an immutable, source-ordered AST; do not infer or check model types."""
    return Parser(text).model()


def symbol_table(model):
    """Keep the first declaration for diagnostics; duplicates remain errors."""
    symbols = {}
    for declaration in model.declarations:
        symbols.setdefault(declaration.name.name, declaration)
    return symbols


def sentence(statement):
    if isinstance(statement, Relation):
        return "{} {} {}.".format(statement.subject.name, statement.verb, statement.object.name)
    if isinstance(statement, Dependency):
        return "{} requires {} in mode {}.".format(
            statement.subject.name, statement.interface.name, statement.mode.name)
    if isinstance(statement, Allocation):
        return "{} allocated-to {} in mode {}.".format(
            statement.subject.name, statement.container.name, statement.mode.name)
    return "{} has {} = {}.".format(statement.subject.name, statement.property, statement.value)


def cycle_diagnostics(edges, relation):
    """Iterative DFS avoids Python recursion limits on deep structural models."""
    graph = {}
    for source, target, span in edges:
        graph.setdefault(source, []).append((target, span))
    colors = {}
    diagnostics = []
    for root in sorted(graph):
        if colors.get(root):
            continue
        colors[root] = 1
        stack = [(root, iter(graph[root]))]
        while stack:
            node, children = stack[-1]
            child = next(children, None)
            if child is None:
                colors[node] = 2
                stack.pop()
                continue
            target, span = child
            if colors.get(target) == 1:
                diagnostics.append(Diagnostic("STRUCTURE_CYCLE",
                    "{} edge {} -> {} closes a cycle".format(relation, node, target), span))
            elif not colors.get(target):
                colors[target] = 1
                stack.append((target, iter(graph.get(target, ()))))
    return diagnostics


def validate(model):
    """Check symbols, typed arguments and structural constraints, not behavior."""
    diagnostics = []
    symbols = symbol_table(model)
    declared = set()
    for declaration in model.declarations:
        name = declaration.name.name
        if name in declared:
            diagnostics.append(Diagnostic("DUPLICATE_DECLARATION",
                "{} is already declared".format(name), declaration.name.span))
        declared.add(name)

    def check(identifier, expected, role):
        declaration = symbols.get(identifier.name)
        if declaration is None:
            diagnostics.append(Diagnostic("UNDECLARED_NAME",
                "{} is not declared".format(identifier.name), identifier.span))
            return False
        actual = declaration.kind
        allowed = (expected,) if isinstance(expected, str) else expected
        if actual not in allowed and not (actual == "container" and "unit" in allowed):
            diagnostics.append(Diagnostic(role + "_TYPE_MISMATCH",
                "{} {} expects {}, received {}".format(
                    role.lower(), identifier.name, " or ".join(allowed), actual), identifier.span))
            return False
        return True

    facts = set()
    properties = {}
    owners = {}
    parents = {}
    allocations = {}
    edges = {"contains": [], "refines": []}
    for statement in model.statements:
        fact = sentence(statement)
        if fact in facts:
            diagnostics.append(Diagnostic("DUPLICATE_FACT", "Repeated fact: " + fact, statement.span))
        facts.add(fact)
        if isinstance(statement, PropertyAssignment):
            check(statement.subject, "functionality", "PROPERTY")
            if statement.value not in PROPERTIES[statement.property]:
                diagnostics.append(Diagnostic("PROPERTY_TYPE_MISMATCH",
                    "{} accepts {}".format(statement.property,
                        ", ".join(sorted(PROPERTIES[statement.property]))), statement.span))
            key = (statement.subject.name, statement.property)
            if key in properties and properties[key] != statement.value:
                diagnostics.append(Diagnostic("PROPERTY_CONFLICT",
                    "{} has conflicting {} values".format(*key), statement.span))
            properties.setdefault(key, statement.value)
        elif isinstance(statement, Dependency):
            check(statement.subject, "capability", "SUBJECT")
            check(statement.interface, "interface", "OBJECT")
            check(statement.mode, "mode", "QUALIFIER")
        elif isinstance(statement, Allocation):
            subject_ok = check(statement.subject, "functionality", "SUBJECT")
            target_ok = check(statement.container, "container", "OBJECT")
            mode_ok = check(statement.mode, "mode", "QUALIFIER")
            if subject_ok and target_ok and mode_ok:
                key = (statement.subject.name, statement.mode.name)
                target = statement.container.name
                if key in allocations and allocations[key] != target:
                    diagnostics.append(Diagnostic("ALLOCATION_CARDINALITY",
                        "{} has multiple Containers in mode {}".format(*key), statement.span))
                allocations.setdefault(key, target)
        else:
            subject_type, object_type = SIGNATURES[statement.verb]
            subject_ok = check(statement.subject, subject_type, "SUBJECT")
            object_ok = check(statement.object, object_type, "OBJECT")
            if not (subject_ok and object_ok):
                continue
            subject, obj = statement.subject.name, statement.object.name
            if statement.verb in edges:
                edges[statement.verb].append((subject, obj, statement.span))
            if statement.verb in ("owns", "contains"):
                mapping = owners if statement.verb == "owns" else parents
                code = "OWNERSHIP_CARDINALITY" if statement.verb == "owns" else "CONTAINMENT_CARDINALITY"
                sources = mapping.setdefault(obj, set())
                if sources and subject not in sources:
                    diagnostics.append(Diagnostic(code,
                        "{} has multiple immediate {}".format(obj,
                            "owners" if statement.verb == "owns" else "parents"), statement.span))
                sources.add(subject)
    for name, declaration in symbols.items():
        if declaration.kind == "functionality" and not owners.get(name):
            diagnostics.append(Diagnostic("OWNERSHIP_CARDINALITY",
                "{} has no valid immediate owner".format(name), declaration.name.span))
    for relation, graph_edges in edges.items():
        diagnostics.extend(cycle_diagnostics(graph_edges, relation))
    return diagnostics


def canonicalize(model):
    """Serialize valid structural facts; never repair an invalid design."""
    diagnostics = validate(model)
    if diagnostics:
        raise ValidationError(diagnostics)
    lines = ["language design-core version 0.2."]
    lines.extend("{} {}.".format(d.kind, d.name.name)
                 for d in sorted(model.declarations, key=lambda d: d.name.name))
    lines.extend(sorted(sentence(s) for s in model.statements))
    return "\n".join(lines) + "\n"


def check(text):
    """Return AST and diagnostics. Canonical layout is checked after semantics."""
    try:
        model = parse(text)
    except ParseError as error:
        return None, [error.diagnostic]
    diagnostics = validate(model)
    if not diagnostics:
        canonical = canonicalize(model)
        if text != canonical:
            offset = next((i for i, pair in enumerate(zip(text, canonical))
                           if pair[0] != pair[1]), min(len(text), len(canonical)))
            diagnostics.append(Diagnostic("NONCANONICAL_FORM",
                "Use canonical ordering, spacing and LF lines with one final newline",
                Source(text).span(offset, min(offset + 1, len(text)))))
    return model, diagnostics


def to_json(value):
    if hasattr(value, "_fields"):
        result = {key: to_json(getattr(value, key)) for key in value._fields}
        result["node"] = type(value).__name__
        return result
    if isinstance(value, (tuple, list)):
        return [to_json(item) for item in value]
    if isinstance(value, dict):
        return {key: to_json(item) for key, item in value.items()}
    return value


def main(argv=None):
    cli = argparse.ArgumentParser(description=__doc__)
    cli.add_argument("command", choices=("ast", "check", "format"))
    cli.add_argument("file", help="UTF-8 model file, or - for stdin")
    args = cli.parse_args(argv)
    try:
        raw = sys.stdin.buffer.read() if args.file == "-" else Path(args.file).read_bytes()
        text = raw.decode("utf-8")
    except (OSError, UnicodeError) as error:
        print(json.dumps({"error": "INPUT_ERROR", "message": str(error)}), file=sys.stderr)
        return 2
    if args.command == "format":
        try:
            output = canonicalize(parse(text))
        except ParseError as error:
            diagnostics = [error.diagnostic]
        except ValidationError as error:
            diagnostics = error.diagnostics
        else:
            sys.stdout.write(output)
            return 0
        print(json.dumps(to_json(diagnostics), indent=2), file=sys.stderr)
        return 1
    model, diagnostics = check(text)
    report = {"valid": not diagnostics, "diagnostics": to_json(diagnostics)}
    if args.command == "ast":
        report["ast"] = to_json(model)
        report["symbols"] = to_json(symbol_table(model)) if model is not None else {}
    print(json.dumps(report, indent=2))
    return 1 if diagnostics else 0


if __name__ == "__main__":
    sys.exit(main())
