"""Immutable, source-preserving SDUI 0.1 syntax tree; no GUI/runtime imports."""
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
    value: str | float | bool | None
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
class Widget:
    name: str
    kind: str
    arguments: tuple[Argument, ...]
    span: Span


@dataclass(frozen=True)
class Row:
    widgets: tuple[Widget, ...]
    span: Span


@dataclass(frozen=True)
class Content:
    rows: tuple[Row, ...]
    span: Span


@dataclass(frozen=True)
class Property:
    name: str
    value: Literal
    span: Span


@dataclass(frozen=True)
class Box:
    name: str | None
    properties: tuple[Property, ...]
    children: tuple['Box', ...]
    content: Content | None
    span: Span


@dataclass(frozen=True)
class ModuleRef:
    alias: str
    path: str
    span: Span


@dataclass(frozen=True)
class Definition:
    name: str
    root: Box
    span: Span


@dataclass(frozen=True)
class Connection:
    module: str
    object: str
    definition: str
    widget: str
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
