# SDP skills — canonical source

Root Skills/ is the sole authored collection for the System Design Process.
Fourteen skills provide one method entrypoint, coordination/implementation/review
roles, [Planning](sdp-planning/SKILL.md) and focused analysis, architecture, evidence and release work. Start with
[sdp](sdp/SKILL.md); load only relevant roles. This replaces the former
Toolkit/skills and Toolkit/skills_v2 maintained collections.

The repository exposes relative per-role symlinks through .agents/skills.
Root AGENTS.md routes SDP work to the entrypoint. The Toolkit install manifest
copies this same collection, including shared references, into .codex/skills in
consuming projects. That existing destination was verified with Codex CLI 0.156.1;
it is not a second authored source. Do not install both discovery trees with
duplicate names. Other hosts need their own measured adapter; discovery alone
does not prove loading or compliant behavior.

[Migration and evidence](../SDP/Maintenance/SK1/Plan-and-Evidence.md) records adoption,
versions, tests and host limitations. Earlier draft evaluations are historical and
do not certify this collection. [Metadata contract](../Toolkit/docs/Skill-Versioning.md)
defines native fields and distribution facts. The Toolkit remains unreleased.
