# Release Notes Policy

The canonical name is `RELEASE-NOTES.md`.

The first section is always:

```markdown
## [Unreleased]

Release-Date: unreleased
```

Use only non-empty categories from Added, Changed, Fixed, Deprecated, Removed,
Security and Migration. Every notable change enters `Unreleased` and should cite
stable SDP IDs when available.

At release preparation, move the selected entries into:

```markdown
## [X.Y.Z] - YYYY-MM-DD
```

Released sections are immutable. A historical correction is a new explicit
`release-notes-corrected` ledger event plus a clearly labelled correction entry;
it is never a silent edit. Validation compares released sections against a
baseline when one exists.

Toolkit migration impact must be explicit. A project installer never overwrites
a populated project release-notes file.
