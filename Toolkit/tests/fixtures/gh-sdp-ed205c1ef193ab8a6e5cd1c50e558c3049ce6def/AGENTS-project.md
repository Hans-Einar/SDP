# gh-sdp Project Instructions

`gh-sdp` is the official GitHub CLI extension project for the future `gh sdp`
command. `Hans-Einar/SDP` is the canonical, independent SDP Toolkit repository;
do not modify that repository without a separate, explicit authorization.

- Follow this repository's project-local `SDP/` records for all work.
- Keep the `gh-sdp` client version distinct from every SDP Toolkit version and
  describe release/support state only from repository evidence.
- Support claims with repository evidence and applicable tests, then require a
  fresh independent review of the actual change.
- In later lifecycle phases, address Windows, Linux, and macOS behavior together
  with preservation, provenance, path-safety, and other security concerns.
- `SPS-001` authorizes governance, Mandate, traceability, and verification
  records only. It authorizes no production implementation, build, packaging,
  or release workflow.
- At installed source commit `bc110bb5fd60009ba67015cf640ad6ddbfe1b04b`,
  `-InitializeProjectStructure` is a one-time bootstrap operation: always preview
  it and do not apply a repeat that proposes the Toolkit-specific
  `SDP/Releases/REL-0.2.0.yaml` in this project.

The managed `AGENTS.md` remains authoritative for role boundaries, Slice/Fix
discipline, review, verification, release gates, and stop conditions.
