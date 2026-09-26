# GIP implementation evidence

Linux x86_64; Go 1.27.1. Evidence identifies HEAD plus the milestone diff;
the introducing commit binds that diff. Live XFMD remains untouched.

## GIP-1-M1

2026-09-26T23:10:42Z. Executable closed records, five canonical golden fixtures, bounded strict JSON/YAML decoding, portable path/root inspection and receipt 3.0 reader delivered. Go package suite passes; management validator passes. No apply engine or signed distribution claimed.

## GIP-1-M2

2026-09-26T23:16:10Z. Direct packaged install/upgrade previews delivered for explicit development artifacts, with deterministic plans, manual adoption, preserving relocation/link rebasing, old/current/new managed comparison and no-change detection. Go race suite passes. Packaged preview produced six actions and left its project empty. Signed resolution and apply remain later milestones; therefore the broader SDL InstallationPlanSlice stays planned.

## GIP-2-M1

2026-09-26T23:20:20Z. Root-bound saved-plan apply, advisory lock, drift revalidation, byte backups, atomic replacements and journaled receipt/history/report publication delivered. Clean install and known upgrade tests pass with race detection. Actual packaged apply completed and discovery reads schema 3.0. Interruption matrix remains GIP-2-M2.

## GIP-2-M2

2026-09-26T23:24:02Z. Process-exit fault matrix passes every preparation/backup/write/checkpoint/completion boundary for clean, known-upgrade and manual-adoption cases. Reserved finalization bytes/IDs and ledger prefixes remain stable; edited targets and legacy pending operations block. Full Go race suite passes (installation package 48.3 seconds). Model apply/recovery activities updated and review regenerated through SDL. No power-loss or rollback guarantee.

## GIP-3-M1

2026-09-26T23:31:03Z. Shared stdlib bootstrap verifies Ed25519 test descriptors, platform binary digest/size and advertised protocol, with private immutable cache and offline checks. Engine independently verifies signed input, records test-signed versus production provenance and resolves previous descriptors by digest. Go profile explicitly reuses files-only legacy inventory without PowerShell prerequisites. Packaged signed install passes; bootstrap race tests pass. No production keys, stable catalog or published release exists.

### GIP-3 integration review and hardening

Independent gh-sdp review found a FIFO-open hang in shared bootstrap input reading.
Both bootstrap and engine now reject nonregular objects before opening, retaining
post-open bounds checks. A native FIFO regression and signed negative tests for
protocol, platform, binary size/path/digest and missing protocol advertisement pass.
HTTP/probe timeouts exist but are not claimed as separately exercised timing tests.

Full engine race suite passes after management-history predecessor/card validation,
cache publication hardening and incoming Markdown ancestor consistency fixes.
The updated version test explicitly expects receipt 3.0 and the now-implemented
installation protocol; this is the authorized compatibility addition, not a waived
legacy obligation. Clean project AGENTS.md is preserved automatically; unequal
AGENTS-project.md collisions still block. A fresh XFMD copied baseline preview/apply
also succeeded during integration; its reproducible paired trial belongs to GIP-4.

## GIP-3-M2

2026-09-26T23:42:57Z. Thin gh-sdp client implemented in its own repository with shared bootstrap pinned to fb79727. Independent review REV-SPS-003-002 approves client fca8480 plus clean engine fb79727 on Linux after resolving FIFO and evidence findings. Actual gh sdp/direct plan and apply parity, offline behavior, race and vet pass. Client implementation does not duplicate installation policy. Wider engine and XFMD evidence remain GIP-4.
