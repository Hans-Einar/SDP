# Distribution And Upgrades

## Goals

- reuse the SDP method across existing projects
- avoid nested Git repositories
- never overwrite project-specific documents accidentally
- permit controlled updates to shared instructions and skills

## Ownership classes

### Toolkit-managed

Files copied from this repository and safe to replace when explicitly requested:

- `.codex/skills/sdp-*/SKILL.md`
- `SDP/Framework/*`

### Project-owned

Never overwritten by the default installer:

- `AGENTS.md`
- `SDP/AGENT-REMINDERS.md`
- Mandate, Study, Requirements, Architecture and Design
- Sprints, Refactors, CodeReview, Verification and Traceability
- project-specific Instructions

Template versions may be installed only when the destination does not exist.

## Recommended model

Keep one normal clone of the upstream toolkit outside all project repositories.
Run the installer against each project. The installed project remains a single
Git repository and records the copied files in its own history.

## Alternatives

### Git submodule

A submodule intentionally embeds another Git repository and pins a commit. This
is useful when exact independent versioning is desired, but adds operational
complexity and does not merge naturally with an existing project-owned `SDP/`
tree.

### Git subtree

A subtree imports another repository into a subdirectory without a nested
`.git`. It supports later pulls, but is awkward when shared files and
project-specific files must coexist and evolve independently.

### `git archive` or release ZIP

Good for a clean one-time copy without `.git`, but provides no automatic update
tracking.

## Version marker

The installer may write `SDP/Framework/VERSION` so projects can record which
upstream toolkit version was installed. Project-specific SDP documents should
remain authoritative even when the framework is upgraded.
