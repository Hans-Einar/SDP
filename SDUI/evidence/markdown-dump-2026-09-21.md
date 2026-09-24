# Static Markdown dump — verification, 2026-09-21

**30/30 tests passed**, including six new export tests.
[Log](markdown-dump-tests-2026-09-21.txt), [manifest](markdown-dump-manifest.json),
[generated dump](../examples/concept1-bucking.dump.md). This was the export evidence
following the frontend/text dump; older manifests describe pre-Markdown source.
The commands are historical, from SDUI:

```sh
PYTHONPATH=src python3 -m sdui examples/concept1-bucking.sdui --format markdown --entry bucking -o examples/concept1-bucking.dump.md
PYTHONPATH=src python3 -m unittest discover -s tests -v
```

Checks covered renderable Markdown, omitted Mermaid, retained other fences, open-fence
closure at widget boundaries, safe enclosing fences despite embedded backticks, hidden
nodes, region/row ordering and nesting. CLI requires profile validation, valid entry and
sufficient overview width; errors produce no partial stdout. Stored output matched a
fresh run byte-for-byte; git diff --check passed.

Separate markdown-it-py 4.2.0 in a temporary venv rendered CommonMark plus tables to HTML.
Assertions found localized heading Lengde, bold measurement, Sagtømmer cell, Cursor AV
code label and nested blockquotes; no button/input elements or Mermaid blocks. This was
structural Markdown→HTML checking, not visual XFMD/browser/PDF verification. No new
production dependency.

Plain Markdown has no general SDUI column/weight mechanism. The export therefore shows
a terminal overview followed by content in reading order. No interaction, callbacks,
new parser, Mermaid/XFMD changes, commit or push occurred in this delivery.
