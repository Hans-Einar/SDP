# SDP Examples

A mature consuming project may contain:

```text
Project/
├── AGENTS.md                         Toolkit-managed
├── AGENTS-project.md                 project-owned
├── .codex/skills/sdp-*/SKILL.md      Toolkit-managed, versioned
└── SDP/
    ├── SDP-project.manifest.yaml     project-owned release/work state
    ├── RELEASE-NOTES.md              project-owned history
    ├── Framework/
    │   └── installed-toolkit.manifest.yaml
    ├── 01--Mandate/ ... 07--Implementation/
    ├── Sprints/ Refactors/ Fixes/ Releases/
    ├── CodeReview/ Verification/
    └── Traceability/
        ├── CurrentIndex.yaml
        ├── Relations.yaml
        └── Ledger.ndjson
```

Repository examples include:

- `installed-toolkit.manifest.example.yaml`
- `build-identity.example.json`
- `release-events.ndjson.example`

They illustrate contracts and are not records of real published releases.
