# SDP Examples

A mature consuming project may contain:

```text
Project/
|-- AGENTS.md                         Toolkit-managed
|-- AGENTS-project.md                 project-owned
|-- .codex/skills/sdp-*/SKILL.md      Toolkit-managed, versioned
`-- SDP/
    |-- SDP-project.manifest.yaml     project-owned release/work state
    |-- RELEASE-NOTES.md              project-owned history
    |-- Framework/
    |   `-- installed-toolkit.manifest.yaml  generated Toolkit-managed facts
    |-- 01--Mandate/ ... 07--Implementation/
    |-- Sprints/ Refactors/ Fixes/ Releases/
    |-- CodeReview/ Verification/
    `-- Traceability/
        |-- CurrentIndex.yaml
        |-- Relations.yaml
        `-- Ledger.ndjson
```

Repository examples include:

- `installed-toolkit.manifest.example.yaml` — installed facts from an archive,
  including truthful null source commit;
- `install-plan.example.json` — deterministic, target-relative plan actions;
- `ledger-events.ndjson.example` — generic work/review/verification events and a
  namespaced project extension;
- `release-events.ndjson.example` — events that also obey the stricter release
  specialization;
- `build-identity.example.json` — generated development identity.

They illustrate schemas and are not records of real work or published releases.
Plan actions include `targetSource` for deterministic migration and every
action/reason pair follows the normative installation decision table; apply
must execute those exact planned fields. The example declares
`migration-first-manifest-order-v1`; migration precedes ordinary entries,
ordinary entries retain installation-manifest array order, backup/mutation
pairs are adjacent, and sequence values are assigned last. Migration actions
also carry the exact-byte `targetSourceSha256` and
`destinationPrecondition`; ordinary actions set both fields to null.

The authoritative cross-client examples are the 17 scenarios under
`Toolkit/conformance/install-v1/`, not this one illustrative plan. A Go or other
independent implementation can parse `scenarios.json`, materialize its portable
before-state declarations, and compare with committed plan JSON or fatal
failure JSON without executing PowerShell.
Neutral installation seeds live in `Toolkit/project-templates/`; none contains
the Toolkit repository's `REL-0.2.0` record, active Sprint, Ledger, review or
verification history. The release-event example uses an illustrative release ID
and is never an installation source.
