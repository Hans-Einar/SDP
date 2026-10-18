# WCI2-M1 actual native acceptance

Frozen binary SHA-256:
`f1de87cefd039f84321c145d02f4b79bbd219916a30b76dfaa3a2007600239cc`.
All four runs used this binary, isolated Xvfb :189 and actual XTest pointer/keys.
No window manager or new text/IME behavior is claimed by these pane tests.

| Variant | Passing checks |
| --- | --- |
| panes | 31 |
| split | 11 |
| vertical | 9 |
| lifecycle | 7 |

Reproduce with `native/verify_wci2.py --binary BINARY --display DISPLAY --out OUT`
and `--variant panes`, `--variant split`, `--variant split --vertical` or
`--variant lifecycle`. Build the desktop-tagged SDL examples/panes/cmd/native
fixture from the candidate specified by Evidence-WCI2-M1.md.

Per variant, actions.ndjson records actual input/control barriers and assertions;
events.ndjson.gz contains the complete unmodified fixture stream. OS screenshots
show final native painting. Every stderr.txt is empty; teardown passed once in
each run. inventory.json hashes these artifacts and final suite logs.

The initial horizontal final trial requested owner state before the native driver
processed the final drag motion. Its log showed the correct .05 keyboard step
from the actual final position. The harness now waits for geometry matching the
full 70-pixel drag before issuing the key. This strengthens the input barrier;
no product bytes or expectations changed. See Pilot-WCI2.md for earlier findings.

Panes explicitly proves selected-header focus without an action, disabled native
and runtime selection agreement, and retained domain succeeded/unknown when UI
publication fails. Split variants prove positive-minimum collapse/restore,
relative keyboard/drag geometry and compatible reload. Lifecycle proves provider
revocation on page departure, late-result rejection and explicit fresh recovery.
