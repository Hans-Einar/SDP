# Example Project SDP Tree

The active SDP document tree belongs inside each project, not at the root of the
upstream toolkit repository.

A mature project may use:

```text
Project/
├── AGENTS.md
├── .codex/
│   └── skills/
│       ├── sdp-master/SKILL.md
│       ├── sdp-worker/SKILL.md
│       ├── sdp-reviewer/SKILL.md
│       ├── sdp-architect/SKILL.md
│       ├── sdp-traceability/SKILL.md
│       └── sdp-vertical-refactor/SKILL.md
└── SDP/
    ├── AGENT-REMINDERS.md
    ├── Framework/
    ├── Instructions/
    ├── 01--Mandate/
    ├── 02--Study/
    ├── 03--Requirements/
    ├── 04--Architecture/
    ├── 05--DesignAnalysis/
    ├── 06--Design/
    ├── 07--Implementation/
    ├── Sprints/
    ├── Refactors/
    ├── CodeReview/
    ├── Verification/
    └── Traceability/
        ├── CurrentIndex.yaml
        ├── Relations.yaml
        └── Ledger.ndjson
```

This is an example, not a payload to copy blindly. Existing project structures
may use singular `Refactor` or different numbered folders. The installer adds
shared guidance and skills without creating or replacing these project-owned
folders.
