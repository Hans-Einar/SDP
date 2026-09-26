# Handoff — Mermaid as diagram engine and reuse source

Updated 2026-09-21: Go/Fyne is selected; [PLAN-003](implementation-plan.md) applies.
The mapping below dates from September 19–20, not a fresh Git inspection. It replaces
the handoff requiring extraction with BoxUI 0.1/BX-HOST compatibility.

## Current scope

Mermaid master at afab5e9 contains no BoxUI. **No Mermaid change is needed now** to
start SDUI development. Do not create a cleanup PR or add SDUI to master. Preserve
ordinary Mermaid, treemap, measurement and routing.

The SDUI session in /home/warloc/git/SDP-vNow owns frontend, runtime and layout/
presentation. Suitable algorithms/tests may be studied and ported from the
implementation worktree to Go with provenance. This creates no automatic Rust
dependency and requires no intermediate legacy crate, JSON parser, wire protocol
or backward-compatible geometry.

Read applicable AGENTS.md plus [plan](implementation-plan.md),
[layout proposal](layout-language-proposal.md), [architecture](target-architecture.md)
and [worktree mapping](renderer-extraction-and-language-direction.md).

## Reuse sources and preservation

- mermaid-rs-renderer-boxui-implementation, 61a85b6: XFMD's implementation source.
  Measurement, SVG escaping/presentation, control maps, budgets and tests are
  candidates; new fr/frame rules replace grow semantics.
- mermaid-rs-renderer-boxui, 4bfd179: parallel early implementation from the same
  base, not a required additional merge. Do not mix these parsers.
- /tmp/xfmd-mrr-measurements: 23 staged files at inspection; no deletion/reset
  without a separate inventory. Outside SDUI cleanup.
- Numbered SDP pilot and mandate in the BoxUI worktrees' SDP directories record
  earlier work, not new authority to retain legacy behavior.

Record source commit/license for reuse. Port tests that still check relevant
properties against new contracts. Identical historical AST/wire/SVG bytes are not
required; geometry tests change when semantics change.

## Possible later Mermaid work

Only a concrete diagram defect/capability gap found by the Markdown provider should
create a separate Mermaid task. The SDUI core must not depend on Mermaid as a
production library; diagram rendering remains a host dependency.

The mapping session removed no source/worktrees and performed no merge, push or
installation. Historical worktrees are provenance, not active language paths.
