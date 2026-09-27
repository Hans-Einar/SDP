# Toolkit Manifest

`SDP.manifest.yaml` is the authoritative Toolkit release manifest and conforms
to `Toolkit/schemas/SDP-manifest.schema.json`. It owns schema version, Toolkit
version/state, Framework and AGENTS contract versions, skill versions,
capabilities, compatibility, release-note and migration paths, supported project
schemas and real publication identities when they exist.

The current payload inventory is `SDPTool/profiles/payload.json`, selected by
`SDPTool/profiles/five-phase.json`. The Go builder emits the authoritative signed
release descriptor. `Toolkit/SDP-install.manifest.json` is retained only as a
schema-regression contract for the retired interface; current installation never
reads it. Root release identity, historical schema examples and actual published
Go descriptors have distinct roles.

The `0.2.0` capability set includes the portable install and plan contracts,
Toolkit/project manifests, release metadata, skill metadata, reusable
CurrentIndex/Relations/generic Ledger/release-event schemas, project validation
and versioning.

During ordinary development, `releaseState` is `unreleased`, `releaseDate` is
`unreleased`, and `gitTag`/`releaseCommit` are null. Publication data is written
only after the corresponding tag and GitHub Release really exist. Released or
yanked records require a real date, tag and commit; the tag must be exactly
`v<toolkitVersion>`. Versions use the full SemVer 2.0 form, including optional
prerelease and build identifiers.

A consuming project neither edits nor copies this root release manifest as its
own state. A conforming installer generates the smaller installed-facts manifest
under `SDP/Framework/` and creates a separate project-owned project manifest only
when missing.
