"""Bounded lexer. Quotes/comments are deliberately independent of Python syntax."""
import math
import re
from dataclasses import dataclass
from .ast import SduiError, Span

MAX_BYTES = 262144
MAX_TOKENS = 50000
_IDENTIFIER = re.compile(r'[A-Za-z_][A-Za-z0-9_]*(?:-[A-Za-z_][A-Za-z0-9_]*)*')
_NUMBER = re.compile(r'-?(?:0|[1-9][0-9]*)(?:\.[0-9]+)?(?:[eE][+-]?[0-9]+)?')


@dataclass(frozen=True)
class Token:
    kind: str
    value: str | float
    span: Span


def lex(source: str) -> tuple[Token, ...]:
    try:
        size = len(source.encode('utf-8'))
    except UnicodeEncodeError as exc:
        raise SduiError('encoding', 'Source must contain Unicode scalar values', Span(0, 0, 1, 1)) from exc
    if size > MAX_BYTES:
        raise SduiError('source-limit', 'Source exceeds 256 KiB', Span(0, 0, 1, 1))
    # Positions remain linear-time, including Unicode and CRLF documents.
    positions = []
    byte, line, col = 0, 1, 1
    for char in source:
        positions.append((byte, line, col))
        byte += len(char.encode('utf-8'))
        if char == '\n':
            line, col = line + 1, 1
        else:
            col += 1
    positions.append((byte, line, col))

    def span(start, end):
        b, l, c = positions[start]
        return Span(b, positions[end][0], l, c)

    def error(message, start, end):
        raise SduiError('lexical', message, span(start, end))

    tokens, i = [], 0
    while i < len(source):
        c = source[i]
        if c in ' \t\r\n':
            i += 1
            continue
        if c == '#':
            while i < len(source) and source[i] not in '\r\n':
                i += 1
            continue
        start = i
        if source.startswith('"""', i):
            end = source.find('"""', i + 3)
            if end < 0:
                error('Unterminated multiline Markdown string', start, len(source))
            value = source[i + 3:end]
            if any((ord(ch) < 32 and ch not in '\n\r\t') or ord(ch) == 127 for ch in value):
                error('Control character in multiline string', start, end)
            i = end + 3
            tokens.append(Token('STRING', value, span(start, i)))
        elif c in "\"'":
            quote, chars = c, []
            i += 1
            while i < len(source) and source[i] != quote:
                c = source[i]
                if ord(c) < 32 or ord(c) == 127:
                    error('Raw control character in string; use an escape', i, i + 1)
                if c == '\\':
                    i += 1
                    if i == len(source):
                        error('Unfinished string escape', start, i)
                    c = source[i]
                    escapes = {'n': '\n', 'r': '\r', 't': '\t', '\\': '\\', '"': '"', "'": "'"}
                    if c == 'u':
                        raw = source[i + 1:i + 5]
                        if len(raw) != 4 or not re.fullmatch('[0-9a-fA-F]{4}', raw):
                            error('Expected four hexadecimal digits after \\u', i, min(i + 5, len(source)))
                        number = int(raw, 16)
                        if number == 0 or 0xD800 <= number <= 0xDFFF:
                            error('NUL and surrogate escapes are unsupported', i, i + 5)
                        chars.append(chr(number))
                        i += 4
                    elif c in escapes:
                        chars.append(escapes[c])
                    else:
                        error('Unknown string escape', i - 1, i + 1)
                else:
                    chars.append(c)
                i += 1
            if i == len(source):
                error('Unterminated string', start, i)
            i += 1
            tokens.append(Token('STRING', ''.join(chars), span(start, i)))
        elif match := _IDENTIFIER.match(source, i):
            i = match.end()
            tokens.append(Token('ID', match.group(), span(start, i)))
        elif match := _NUMBER.match(source, i):
            i = match.end()
            value = float(match.group())
            if not math.isfinite(value):
                error('Number must be finite', start, i)
            tokens.append(Token('NUMBER', value, span(start, i)))
        elif c in '[]{}(),;=.:@<>*^|-':
            i += 1
            tokens.append(Token(c, c, span(start, i)))
        else:
            error(f'Unexpected character {c!r}', i, i + 1)
        if len(tokens) > MAX_TOKENS:
            raise SduiError('token-limit', 'More than 50000 tokens', span(start, i))
    tokens.append(Token('EOF', '', span(i, i)))
    return tuple(tokens)
