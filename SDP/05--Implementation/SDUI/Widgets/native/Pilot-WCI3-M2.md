# WCI3-M2 native pilot history

Evolving candidates, not final stage acceptance. Main uses actual SDL text fixture,
Xvfb :190 at 1600x1200, XTest and external X11 clipboard. Configured IBus/XIM
composition is separate from pasted Unicode and remains required.

Coordinator pilot1 binary SHA-256
`0cff7b7bcd6a9e259edbb7b6c5d5fc3c822166750e9cf2c00b3bffd268c30015`,
/tmp/wci3-m2-text-pilot1. Built with GOWORK=off, -mod=readonly and desktop tag from
SDL/go during implementation; not a frozen all-file candidate. Build information
selects ../../SDUI/third_party/glfw. Worker concurrent pilot d4cdefe4 differs;
no byte-equivalence or final source freeze is claimed.

Single-line workflow passes six actual native assertions: explicit false opts into
typed policy, exact Unicode paste changes draft only, Tab does not Commit, Return
saves once, Escape restores accepted text, and teardown closes session/providers.
Original-resolution OS capture inspected. Further multiline, history, paste refusal,
read-only, dialogs, failures, retention, constraints, scroll and IME remain underway.

- Pilot1 multiline passes exact paste, ordinary Enter, Primary+Enter, copy, cut
  and Undo, then actual Ctrl+Shift+Z fails to redo. Pinned driver emits a custom
  shortcut while Entry registers only the named Redo form. Host adds a bounded
  exact Primary+Shift+Z mapping with a focused regression; actual rerun remains
  required. This pilot is not a passed multiline workflow.
- Two subsequent cold launches timed out before ready under the original 15-second
  harness startup bound while host/text race suites were running. Both emitted an
  initial accepted state and a transient pane-minimum diagnostic; the diagnostic
  is not proven to cause the timeout. Startup-only bound is now 45 seconds, with
  interaction waits unchanged. Original failed logs are retained.

## Resolved pilot findings and harness assumptions

Actual Ctrl+Shift+Z now passes the complete eight-check multiline workflow after
the bounded host mapping. The final worker/coordinator build is SHA256
`0078f72eec7fe61bd41a1b090bb7ca2f26ba48c3a1e54c4eff6bd240c010f001`.
Detached measurement now uses policy/placeholder/font with empty text: independent
short/long Unicode/CRLF and font-size checks prove the scrolling Entry minimum
is text-independent. Live Entries retain complete text. This removes unnecessary
font/layout work without changing runtime data or adding a geometry cache.

Three failed assertions exposed harness assumptions, not product defects:

- Fyne groups an initial paste and following typing into one native undo action.
  The history workflow now creates independent actions with distinct real pastes.
- Pinned Linux Fyne does not implement Ctrl+Home/End. Native select-all plus
  Left/Right reaches the multiline start/end; plain End is used for single-line.
- The intentional stale-publication fixture changes the preview label before
  trying its stale candidate. That legitimate publication advances StateRevision;
  the corrected assertion permits exactly that increment while preserving all
  other field values and revisions.

A 45-second startup retry admitted the earlier candidate but its Load completion
arrived after the unchanged 15-second interaction bound. The raw final accepted
Unicode state remains retained. Concurrent CPU/font work is not a proven cause.
The earlier text race suite exceeded its unchanged 600-second limit; the frozen
optimized candidate completes in 508.408 seconds without a timeout waiver.

Final-build pilots pass read-only copy, native history, CR/LF paste refusal,
required-empty retention, constraint rollback, modal forms, failed action/reload
preservation, hidden page retention and native scrolling. Strengthened scrolling
adds actual narrow-window wrapping and native scrollbar dragging, nine checks.
Configured real IBus/XIM passes seven checks, including consumed Return/Escape
and zero preedit SDL actions. Fixture stderr is empty; isolated D-Bus/portal
permission diagnostics are environment evidence, not a clean-environment claim.
The configured XIM style did not visibly draw marked preedit/candidate panels,
so those visuals are not claimed. Fresh final matrix evidence remains separate.

## Late shared-command/tab receiver defect

After the 0078 native matrix passed, fixture preparation inspection found that
interactionHandler did not include ExpectedDraftRevision on an extended input
TextResult update. Independent actual SDL command/tab overlays reproduce one
domain call followed by field-conflict and unchanged UI state. This is a product
defect, unlike the earlier harness assumptions. M2 approval was suspended and
the full candidate retained in WCI3-M2-before-interaction-fix. The correction must
preserve legacy receivers, post-Execute conflict checks, domain outcome and no
automatic replay; a new actual command Load variant covers the cross-family route.

Post-correction binary 2579df7a exposed one context-menu orchestration failure:
immediate End/Up/Return after right-click reached the Entry, producing SaveSingle
once instead of invoking Paste. The harness had no popup-ready observation. It
now observes actual Entry focus loss to the native popup, captures its six visible
items and sends navigation without forcing X focus. The same-binary corrected
workflow passes with zero SDL calls and preserved Undo. Independent review accepts
this readiness barrier; delay versus forced-focus causation was not isolated. The
failed raw run remains diagnostic evidence and is excluded from final acceptance.
