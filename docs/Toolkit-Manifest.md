# Toolkit Manifest

`SDP.manifest.yaml` is the authoritative Toolkit release manifest and conforms
to `Toolkit/schemas/SDP-manifest.schema.json`. It owns schema version, Toolkit
version/state, Framework and AGENTS contract versions, skill versions,
capabilities, compatibility, release-note and migration paths, supported project
schemas and real publication identities when they exist.

`Toolkit/SDP-install.manifest.json` is a separate authority. It owns the
installable inventory, generators, destinations, ownership and update policies.
The two manifests must agree on Toolkit version and the complete ordered
capability list; Toolkit validation detects disagreement.

The `0.2.0` capability set includes the portable install and plan contracts,
Toolkit/project manifests, release metadata, skill metadata, reusable
CurrentIndex/Relations/generic Ledger/release-event schemas, project validation
and versioning.

During ordinary development, `releaseState` is `unreleased`, `releaseDate` is
`unreleased`, and `gitTag`/`releaseCommit` are null. Publication data is written
only after the corresponding tag and GitHub Release really exist.

A consuming project neither edits nor copies this root release manifest as its
own state. A conforming installer generates the smaller installed-facts manifest
under `SDP/Framework/` and creates a separate project-owned project manifest only
when missing.
