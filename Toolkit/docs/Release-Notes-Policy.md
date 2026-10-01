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

## Per-release generated logs

`sdptool release-log --all --output Releases` generates one Markdown file for each
versioned section. --check compares without writing; differing existing logs are
rejected. `--version X.Y.Z` selects one section (stdout unless --output is given).
Unreleased is never emitted as a published-version log. Canonical notes remain the
authored authority; generation does not imply release approval or publication.
Commit the generated log and use it as GitHub Release notes after authorization.
