# SK1 migration contract

The owner selected execution before the installer Scrum. Source baseline:
aabb359; canonical implementation remains Toolkit/skills until SK1-M2 switches
all current consumers together. SK1-M1 has prepared the thirteen replacement
entrypoints in the draft location, without adding a discovery adapter yet.

| Consumer | Atomic SK1-M2 change and removal condition |
| --- | --- |
| Toolkit/skills and Toolkit/skills_v2 | Move prepared roles/references to root Skills/; remove both maintained old collections once their manifest and references migrate |
| Root AGENTS and native catalog | Route SDP work through Skills/sdp; add thirteen relative .agents/skills symlinks, no duplicate source files |
| Install manifest and source allowlists | Read root Skills/; explicitly inventory all thirteen roles and shared references |
| Installed copies and generated AGENTS | Retain .codex/skills destination; install complete adapted collection and route through sdp; do not install a second discovery tree |
| Metadata validator and tests | Accept native name/description plus string-valued metadata fields; preserve support for old installed v1 metadata without allowing new v2 claims to bypass native validation |
| Version authorities | Ten replacement roles 2.0.0 because procedure/contract changes; three new roles 1.0.0; AGENTS contract 2.0.0; Toolkit 0.2.0 remains unreleased |
| Same-Toolkit-version skill upgrade | Require the existing explicit ForceManagedFiles option if declared skill versions/inventory change; fail before any mutation without it, preserve changed old files via existing backup policy when forced |
| Conformance fixtures/docs | Rebuild expected plans explicitly after reviewing changed source/inventory facts; preserve frozen consuming-project fixtures and historical evidence |
| Dated evaluation/probe | Retain original JSON and evaluation narrative; make the probe use baseline Git bytes when old source paths disappear |

A fresh Codex CLI 0.156.1 probe also discovers a native-metadata skill in a
.codex/skills directory. This measured compatibility permits preserving current
install destinations. Root .agents/skills remains the documented repository
adapter. Installed copies are generated from root sources, never maintained
separately. This test does not claim support on every host/version.

New metadata retains skillId, skillVersion, minimumToolkitVersion,
compatibilityNotes and capabilities under native metadata. All values are strings;
capabilities is comma-separated and parsed into the existing validator's list.
Native name equals skillId. This is recorded as sdp.skill-metadata.v2 alongside
v1 support; it is not a newly released Toolkit. The existing install-v1 path and
plan schema remain adequate; reject an unsupported same-version transition with
an actionable preflight failure rather than generating untrue installed facts.

The broader installer Scrum owns portable process-profile migration, automatic
adoption of manually initialized projects and future destination changes. No
actual consumer repository is upgraded during SK1. Test clean/repeat installation,
forced old-skill upgrade, unforced preservation/failure, external-link containment
and installed catalog discovery in disposable projects.

Baseline on aabb359: 80 Python tests ran, 79 passed; the repository-validation
test reports 38 pre-existing Traceability errors (IDs and reciprocal relations).
The full validator reports the same 38; this supersedes older nine-error counts.
Keep that exact baseline visible under KB-SDP-011, not a blanket all-tests-pass
claim and not permission to change unrelated relations.
