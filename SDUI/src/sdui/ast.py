"""Immutable, source-preserving SDUI 0.2 syntax tree; no GUI/runtime imports."""
from dataclasses import asdict, dataclass
from typing import TypeAlias


@dataclass(frozen=True)
class Span:
    start: int  # UTF-8 byte offset, inclusive
    end: int  # UTF-8 byte offset, exclusive
    line: int  # one-based Unicode scalar line/column
    column: int


class SduiError(Exception):
    def __init__(self, code: str, message: str, span: Span):
        super().__init__(message)
        self.code, self.message, self.span = code, message, span

    def to_dict(self):
        return {"code": self.code, "message": self.message, "span": asdict(self.span)}


@dataclass(frozen=True)
class Literal:
    kind: str  # string, number, boolean, null
    value: str | float | bool | tuple[float, ...] | None
    span: Span


@dataclass(frozen=True)
class Reference:
    module: str
    object: str
    member: str
    span: Span


Value: TypeAlias = Literal | Reference


@dataclass(frozen=True)
class Argument:
    name: str | None
    value: Value
    span: Span


@dataclass(frozen=True)
class LayoutRule:
    name: str  # property name, arrow or ratio
    value: Literal
    spelling: str
    span: Span


@dataclass(frozen=True)
class Row:
    items: tuple['Node', ...]
    span: Span


@dataclass(frozen=True)
class Node:
    kind: str  # frame, group, markdown, widget, use
    name: str | None
    role: str | None
    rows: tuple[Row, ...]
    widget: str | None
    arguments: tuple[Argument, ...]
    text: Literal | None
    target: str | None
    variant: str | None
    layout: tuple[LayoutRule, ...]
    span: Span


@dataclass(frozen=True)
class ModuleRef:
    alias: str
    path: str
    span: Span


@dataclass(frozen=True)
class Definition:
    name: str
    root: Node
    span: Span


@dataclass(frozen=True)
class Connection:
    module: str
    object: str
    definition: str
    path: tuple[str, ...]
    span: Span


@dataclass(frozen=True)
class Document:
    profile: str
    references: tuple[ModuleRef, ...]
    definitions: tuple[Definition, ...]
    connections: tuple[Connection, ...]
    span: Span


def to_data(value):
    """Stable tagged JSON form. Source spelling/comments are not a concrete syntax tree."""
    if hasattr(value, '__dataclass_fields__'):
        return {"type": type(value).__name__, **{
            name: to_data(getattr(value, name)) for name in value.__dataclass_fields__}}
    if isinstance(value, tuple):
        return [to_data(item) for item in value]
    return value
