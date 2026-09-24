# SDUI 0.1 — verification, 2026-09-19

Historical evidence for retired SDUI 0.1. Hashes identify files at that date, not current
ported code. The next historical evidence is the [0.2 frontend/console](frontend-console-2026-09-21.md);
current Go evidence is under go/evidence. This was bounded parser evidence, not rendering,
SDL execution or XFMD integration. [Source manifest](source-manifest.json) identifies the
prototype independently of other local repository changes.

## Executed

Historical commands from the SDP worktree:

```sh
PYTHONPATH=SDUI/src python3 -m sdui SDUI/examples/main-page.sdui \
  -o SDUI/examples/main-page.ast.json
PYTHONPATH=SDUI/src python3 -m unittest discover -s SDUI/tests -v
```

**15/15 tests passed**; [log](parser-tests.txt). The manifest records the environment.
Tests use standard library/temporary files, no network or referenced SDL files.
Checks covered outer box, two grouped child boxes, SVG box, widget order/rows; distinct
callback/setHandle AST identity; Unicode/escapes/comments, CRLF and UTF-8 byte spans;
exact version, duplicate names, unknown module/definition/widget and invalid properties/types;
arbitrary code, misplaced @, trailing commas and empty rows; byte/token/depth/node/argument
limits; all minimal-fixture truncations and 300 deterministic noise strings producing
controlled results; reproducible AST; CLI streams, exit 2/3, invalid UTF-8, missing files
and source overwrite protection.

EBNF/parser were reviewed against the same constructs. EBNF is not compiled automatically;
this is not formal language-equivalence proof. Shared syntactic/semantic limits are documented.

## Limits

First diagnostic only, no editor recovery/formatter, comment retention, persistent handles
or external symbol validation. No geometry/readability test. No complete Markdown-file or
BoxUI-JSON input. Python 3.11 was the chosen minimum; only the manifest's actual version was
tested. Do not claim other environments.

## Preservation and delivery

Only SDUI was added to the SDL development worktree. Existing SystemDesignLanguage/README.md
changes and untracked SDL work documents were preserved. No XFMD/Mermaid code, installation,
branch switch or external publication occurred. At this milestone the work was local and
ready for review.
