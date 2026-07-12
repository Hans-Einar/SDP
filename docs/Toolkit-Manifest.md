# Toolkit Manifest

`SDP.manifest.yaml` is the authoritative Toolkit release manifest and conforms
to `Toolkit/schemas/SDP-manifest.schema.json`.

Toolkit-owned fields include schema version, Toolkit version and state,
Framework and AGENTS contract versions, skill versions, capabilities,
compatibility, release-note path, migration metadata, supported project schemas,
and real publication identities when they exist.

During ordinary development, `releaseState` is `unreleased`, `releaseDate` is
`unreleased`, and `gitTag`/`releaseCommit` are null. Publication data is written
only by the post-publication reconciliation step.

A consuming project does not edit or copy this root manifest as its own state.
The installer writes a smaller generated installed-facts manifest under
`SDP/Framework/`.
