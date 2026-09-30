# RL1 evidence

Baseline a0aa805; Go 1.27.1; Linux; GOMAXPROCS=2; -p 2. Scope is release preparation,
not publication or a live XFMD upgrade. Latest GitHub release observed is v1.0.0.

## RL1-M1

Go SDPTool ./... tests pass, including deterministic release-log extraction/check,
unknown/duplicate versions, date errors, fenced example headings, preservation of
existing differing logs, inventory completeness and real current-payload installation.
The predecessor inventory assertion initially expected two published releases; it
now asserts all three original digests after adding the verified 1.0.0 descriptor.
No predecessor was removed or weakened. Three historical per-release logs are
machine-generated from unchanged versioned RELEASE-NOTES.md sections; --check passes.

Manual Sessions distribute only README and blank template, initialize-if-missing.
The receipt advertises sdp.sessions.manual.v1. Managed AGENTS and ReleaseChecklist
point agents to the workflow. Existing Session conversations are never payloads.
The installer test preserves project-owned guide/session contents and checks no-op.

ProjectManagement validation: 55 cards, 28 management records, three lineage
operations, 425 events before milestone completion; Toolkit validator against
origin/main passes. Git diff whitespace checks pass. New public command is
`sdptool [ROOT] release-log --all --output Releases [--check]` or --version X.Y.Z.
No release identity is fabricated by generating a log.

## RL1-M2 (in progress)

Upgrade-rehearsal.json records actual temporary-project install/apply/upgrade.
Original 1.0.0 descriptor downloaded from GitHub, verified through compiled publisher
trust, SHA-256 767527e0d7f54866bab95f4642ffffb2a344024c1805f44b3d5ea9c024829125.
Candidate uses explicit local-development provenance, not a production signature.
Fresh and custom-content projects pass; four/three actions respectively, preserved
owner documents/history prefix, installed Sessions/checklist and repeat no-op.
Further supported predecessor checks and independent review follow before closeout.
