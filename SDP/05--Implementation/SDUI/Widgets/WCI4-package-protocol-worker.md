# WCI4 packaged-protocol Worker handoff

Status: harness complete and frozen. Candidate1 protocol pilot PASS: eight cases,
50 real executable invocations. This is prior-candidate evidence; the coordinator
reports a pending host parent-hide preview-detachment correction. Rebuild the final
matching package after that correction and rerun this unchanged harness.

## Authorized writes and preservation

Only WCI4-package-protocol.py and this new root report were written in the clone.
SDP Worker/document workflow reused. No product, module, canonical/Session/PM,
branch or commit changes; no broad tests or helper rebuilds. The previous frontend
report remains byte-identical at SHA-256
69a270c7e446c05d7dc2ab75dac4b3283cdd7b773744c932309c0ffadbd6ade8;
all its 23 frozen frontend path hashes were rechecked unchanged after this task.
Main owns integration, final package/native execution, codegen evidence association
and archive publication. This task adds no OS-native, IME or connected-host claim.

## Reproducible command

Run with Python 3 (stdlib only), replacing the final package/output paths. Output
must be an empty directory outside the read-only package; existing evidence is never
overwritten. All executable invocations use argv arrays, preserving spaces/Unicode.

```sh
python3 /tmp/sdp-sdui-widgets/WCI4-package-protocol.py \
  --package /path/to/final-matching-package \
  --out '/tmp/wci4 final package protocol evidence' \
  --source /tmp/sdp-sdui-widgets/SDUI/go/parser/testdata/wci4-all-families.sdui
```

Required layout: bin/sdptool, bin/sdui,
libexec/xfmd/sdui-preview and libexec/xfmd/sdui-fyne. Missing executables yield
INCOMPLETE/exit 2, with their absence recorded; no missing-binary PASS.
An explicit --pilot-helpers option runs just helper cases and cannot produce final
PASS (INCOMPLETE or FAIL). No automatic build, dependency install or module fallback.
The --source default is the parser fixture adjacent to the script in this checkout;
provide it explicitly when archiving/running elsewhere. Its exact bytes are copied
to the output and hashed, not regenerated or reduced to a separate invented model.

Each command retains full argv/cwd, expected outcome, actual exit, duration, timeout,
stdout/stderr bytes and hashes, plus executable identity and unchanged-binary checks.
A 30-second timeout fails acceptance. DISPLAY/WAYLAND_DISPLAY are removed. The
native executable is invoked only with stale revisions or unsupported .3 sources;
no .2 positive native launch or -check flag is invented for sdui-fyne.
summary.json contains all receipts, fixture/harness/binary hashes, case results and
an artifact inventory. The exact executed harness is copied as harness.py. Exit 0
means this bounded protocol suite passed, not final WCI4 acceptance/publication.

## Checked actual seams

| Protocol case | Assertions |
| --- | --- |
| 0.2 and 0.3 discovery/UIPreview | Exact source profiles and sdptool/0.2 protocol; one discovered source; revision-pinned Main target; published response identity; spaced/Unicode source paths; UIPreview bytes equal actual sdui Markdown command; stale source preserves prior bundle. |
| 0.2 and 0.3 Combined | Source bytes/hash, artifact hashes, exact legacy envelope keys, .3 metadata, immutable snapshot; every source link has a valid UTF-8 byte range and matching line/column; multibyte content/reuse; stale request preserves artifacts. |
| Complete family fixture | Copy/hash actual phase parser fixture; Combined exposes every family path including hidden/reused content; AST/source tags and all widget kinds; structural dump/Markdown explicitly report unsupplied resources/provider data. |
| Legacy static export | Both profiles accept basic AST/dump/SVG; use ASCII glyphs for positive legacy static SVG separately from UTF-8 structural proof. |
| Standalone protocols | Positive .2 sdui-preview -check reports local prototype with no SDL runtime and no new legacy profile key; preview/native stale revisions reject. Native and preview helpers reject hidden .3 collection/pane/interaction/value/choice/text/SVG/Markdown adapters with the actual sourced diagnostic. |
| Atomic unsupported export | All-family transient controls plus separate hidden opt-in SVG/Markdown reject publicly with exact diagnostic codes/path/span, empty stdout and unchanged prior output bytes. |

No codegen command exists in the four-binary package. summary.json explicitly
records codegen NOT_PROVIDED, with the required separate exact-candidate SDUI
codegen suite and actual emitted-constructor execution. Existing frontend handoff
identifies that evidence; main must associate its final matching-candidate run.
The harness neither fabricates a CLI nor runs Go/compiler tests from an unrelated
checkout. It also does not prove package build roots, dependencies/licenses,
installation, connected document-host behavior, GUI/IME or publication.

## Actual validation and candidate1 evidence

Python syntax compilation and --help passed. An earlier helper-only rehearsal
correctly failed stale pre-WCI4 binaries at markdown/description schema admission;
no product change was made to conceal this. Its final harness receipt is
/tmp/wci4-package-protocol-helper-pilot-2/summary.json (21 commands).

Coordinator then authorized candidate1 full protocol pilot:

```sh
python3 WCI4-package-protocol.py \
  --package /tmp/wci4-package-candidate1 \
  --out /tmp/wci4-package-protocol-candidate1-pilot-1 \
  --source SDUI/go/parser/testdata/wci4-all-families.sdui
```

Result: PASS, exit 0, eight cases/50 commands. Full evidence directory:
/tmp/wci4-package-protocol-candidate1-pilot-1

- Harness SHA-256: 6620f594bc3ca76ac8a08f872d51310dadb012ebe532d9780c8c99a2206f4072
- Source fixture SHA-256: 1681e0610e41df77917740052b81b5b1630a8c36baad9388e1ab33fc2b1c87e0
- Candidate1 summary SHA-256: 26e3c95e54291009d89af26e81b2372c43939084d3dfb3e2facb9b3fe4d64e7d

| Candidate1 executable | SHA-256 |
| --- | --- |
| sdptool | `6537cf23381394bca2c2c2de10180cdda6406698ee74efbef005d72d221e8756` |
| sdui | `7789ef39c0e5b484fe569b0e7f42066871eeec1ef73a9e0c29012c8726e944ad` |
| preview | `b0e11827124f52c5f6b384d615a2d2bbbd72e151b23cceabbcfa1ff93c422ae7` |
| native | `a1a101d9502c2442d064d1b779bb8c1e5828f5ac7f38dd8f74e450fcc315ca51` |

Candidate1 is not the final package: the upcoming host correction changes the
candidate even though these source/helper routes are unaffected. The exact final
package must pass a new run with a new empty output directory. Keep both pilot and
final receipts; do not rename this pilot as final evidence. No remaining harness
failure or product change is known/claimed by this lane.
