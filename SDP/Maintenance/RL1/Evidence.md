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

## RL1-M2

Candidate: 6fc662d plus the RL1-M2 fence correction, final guide/template and
closeout. Upgrade-rehearsal.json identifies the exact generated candidate payload
hash. Final source identity is the RL1-M2 commit containing this document; there
is no production release identity. Go race-enabled SDPTool ./... passes; after
the fence correction, targeted regression and independent uncached full tests pass.

Original descriptors downloaded from published GitHub releases and verified by the
current Go installer through its compiled publisher key, not a local test key:

| Previous release | Verified descriptor SHA-256 |
| --- | --- |
| 0.2.0 | edc0c72101a437c6e12c40a081ef59ae41cf0db32bcaedb73824ec48495aaee5 |
| 0.2.1 | 66590e8e967ede6b36d8fa45cdbee1cd80f69505cca698b0e4a9bd960842735a |
| 1.0.0 | 767527e0d7f54866bab95f4642ffffb2a344024c1805f44b3d5ea9c024829125 |

For each predecessor: install into a fresh temporary project using --release and
--apply; optionally author a Session/guide/mandate; build the current Go descriptor
with --release development-rl1-final; upgrade with --artifact/--allow-unreleased
and saved-plan apply; inspect payloads/capability/history; repeat preview. All six
cases pass. New guide/template/checklist appear, owner contents and historical
ledger prefixes are preserved, and repeat reports noChange=true/zero actions.
The final candidate descriptor is local-development. This is not a production
signature, extracted-source release gate or live XFMD upgrade. No keys were read
into repository evidence, no live project installation changed.

Independent read-only Reviewer initially found a P2 fence issue: four-backtick
examples could close on three markers and produce an invented release heading.
Fix tracks opening marker and length, requires whitespace-only closing suffix,
and rejects unterminated fences. Backtick/tilde nested examples are covered. The
Reviewer independently ran uncached SDPTool tests and all three generated-log
checks and approved bounded release preparation after correction. It reviewed,
but did not independently repeat, the six install rehearsals. Generic/adopted
Session wording replaces local KB references; the planning link resolves through
the distributed ProjectManagement README. No product/runtime change beyond the
reviewed log fix followed. Full publication and native XFMD acceptance excluded.

[ReleaseChecklist](ReleaseChecklist.md) explicitly leaves freeze/version alignment,
clean exact production signing, CI/publication/download acceptance and client
selection pending. [Manual-upgrade](Manual-upgrade.md) gives future preview/apply
commands. Historical released notes are byte-for-byte unchanged; three generated
logs pass --check. Inventory tests ensure template files cannot disappear silently.

## Outcome

The selected preparation is complete. KB044 awaits a concrete owner decision on
next release integration/publication; this Maintenance does not impersonate that
approval. KB042 retains automated capture/identity work; manual adoption is now
implemented. SDP 1.1.0 is proposed; design-core 0.6 is a separate language version.
No tag, release, gh-sdp default change or XFMD upgrade was performed.

Final closeout checks pass: targeted release-log race tests after the correction;
ProjectManagement (55 cards, 28 management records, three lineage operations,
429 events), Toolkit against origin/main, changed Markdown destinations, unchanged
versioned notes, both append-only ledger prefixes and git diff --check. Unrelated
sourceinput/Node files remain untracked and excluded. KB044 is gate-review for the
concrete publication/version decision, while the preparation plan is complete.
