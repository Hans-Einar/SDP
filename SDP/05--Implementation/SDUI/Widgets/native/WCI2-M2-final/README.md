# WCI2-M2 final native evidence

Frozen binary SHA-256: `2c1db6d173c67ff873da4a7fdadb2a32722d3153705523622328e0d509be7806`.
The tested source inventory records 113 files at HEAD `a4f2c22435d91fe07935b8b9d5fcf46fafc6e36b` plus inventoried changes.
An independent reviewer rebuilt the binary byte-identically.

All 118 assertions in twelve completed scenarios passed. All twelve stderr files
are empty. `terminal-results.json` audits exactly one terminal result for each of
26 published openings, no unpublished receipts, and clean teardown with no pending
provider requests. Raw events, actual XTest input records and OS captures are retained.
Run `python3 ../audit_wci2_results.py */events.ndjson.gz` with the interrupted directory
excluded to reproduce the completed-scenario audit.

The nested scenario specifically proves a source-root sibling dialog follows its
actual nonmodal ParentSurface for canvas, size and focus. Native WM_DELETE_WINDOW
loss of that parent revokes both openings through the exact lifecycle target,
reason parent-closed and sequence zero. Explicit source Close/Cancel remains a
user operation with its triggering sequence. The coordinator clarified this existing
owner-lifetime boundary in the design; independent review approved the distinction.

The initial multi-scenario runner was terminated (exit 143) during accept-failures,
after eight completed scenarios. Its partial data is preserved separately and is
NOT counted as a passed run. Three remaining scenarios were rerun individually and
completed with exit zero. No product assertion failed in the interrupted run.
An earlier display-process loss before Ready likewise supplies no acceptance proof.

Full final phase SDUI race, selected SDL race and original full SDPTool suites
completed with exit zero. All 168 SDUI Go/module files match the original integrated
workspace exactly (`original-equivalence.json`). SDUI host race elapsed 109.795s.
The earlier 58 pane regression checks are reused from
[the predecessor archive](../WCI2-M2-candidate-22539389/README.md): their distinct
binary is identified there; the final delta only concerns auxiliary parent ordering
and native-owner closure. Independent review accepted that scope of reuse.

Environment: Go 1.27.1 linux/amd64, pinned Fyne 2.8.1, Xvfb :190 at 1600x1200,
XTest and OS captures. No window manager, decoration-click or Wayland coverage is
claimed. Native close sends actual WM_DELETE_WINDOW. XMoveWindow only arranges
nonmodal windows to make both surfaces reachable. Early pre-Ready minimum-size
rejection and deliberately injected failures are retained in raw logs; a blanket
claim of zero diagnostic events is not made. The predecessor preference EOF and
its clean repeat remain documented in the predecessor archive.
