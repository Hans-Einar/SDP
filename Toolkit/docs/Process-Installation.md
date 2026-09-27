# Process installation through SDPTool

SDPTool is the sole installation and upgrade engine. The former profile-artifact
pipeline is retired. Current source configuration lives in
[SDPTool/profiles/five-phase.json](../../SDPTool/profiles/five-phase.json), with the
explicit [payload inventory](../../SDPTool/profiles/payload.json).

Go builds the release descriptor directly from the inventory. Published descriptors
are immutable and signed; the installed receipt binds the project to its exact
release and digest. A project-local manifest is not authoritative payload ownership.

Use `gh sdp install` for a new project and `gh sdp upgrade` for an installed project.
Both preview changes before applying a saved plan. Paths, ownership, signatures,
backups, journal recovery and predecessor checks are specified in
[installation records](../../SDPTool/install/Records.md).

Missing project-owned templates initialize. Existing project content and history
are preserved. Explicit adoption input is required when no authoritative predecessor
inventory is available. No folder shape or version label alone proves ownership.
New SDL guidance does not automatically move models or register navigation.
