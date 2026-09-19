"""Parse a standalone SDUI file or stdin into a versioned JSON AST."""
import argparse
import json
import sys
from pathlib import Path
from . import SduiError, parse, to_data, validate
from .ast import Span
from .lexer import MAX_BYTES


def main():
    cli = argparse.ArgumentParser(description=__doc__)
    cli.add_argument('source', help='UTF-8 .sdui file, or - for stdin')
    cli.add_argument('-o', '--output', help='JSON output file (stdout by default)')
    cli.add_argument('--syntax-only', action='store_true', help='Skip local profile validation')
    args = cli.parse_args()
    try:
        if args.output and args.source != '-' and Path(args.output).resolve() == Path(args.source).resolve():
            raise OSError('Output must not overwrite source')
        if args.source == '-':
            raw = sys.stdin.buffer.read(MAX_BYTES + 1)
        else:
            with open(args.source, 'rb') as stream:
                raw = stream.read(MAX_BYTES + 1)
        if len(raw) > MAX_BYTES:
            raise SduiError('source-limit', 'Source exceeds 256 KiB', Span(0, 0, 1, 1))
        try:
            source = raw.decode('utf-8')
        except UnicodeDecodeError as exc:
            prefix = raw[:exc.start].decode('utf-8')
            raise SduiError('encoding', 'Invalid UTF-8', Span(exc.start, exc.end,
                prefix.count('\n') + 1, len(prefix.rsplit('\n', 1)[-1]) + 1)) from exc
        tree = parse(source)
        if not args.syntax_only:
            validate(tree)
        result = {'astFormat': 'sdui-ast/0.1', 'validation': 'syntax-only' if args.syntax_only else 'local-profile',
                  'document': to_data(tree)}
        text = json.dumps(result, ensure_ascii=False, indent=2, allow_nan=False) + '\n'
        if args.output:
            Path(args.output).write_text(text, encoding='utf-8')
        else:
            sys.stdout.write(text)
        return 0
    except SduiError as exc:
        print(json.dumps({'source': args.source, 'error': exc.to_dict()}, ensure_ascii=False), file=sys.stderr)
        return 2
    except OSError as exc:
        print(json.dumps({'source': args.source, 'error': {'code': 'io', 'message': str(exc)}}), file=sys.stderr)
        return 3


if __name__ == '__main__':
    sys.exit(main())
