# SDP project templates

[sdp-root](sdp-root/README.md) is the canonical template for a new SDP directory.
Its [document guide](sdp-root/SDP-DOCUMENT-GUIDE.md) explains process documents;
its [SDL guide](sdp-root/SDL/README.md) defines System-owned model source homes.
The five-phase layout is current; no parallel Template/profiles tree is maintained.

Templates contain neutral project-owned initial content, never live project cards,
history or models. Managed instructions/Framework payloads live in Toolkit/payload;
canonical skills live in Skills. Repository lifecycle records under SDP are not
installation seeds.

## Distribution and upgrades

[Toolkit/profiles/five-phase.json](../Toolkit/profiles/five-phase.json) is the
explicit current inventory. The Go descriptor builder uses its files only and
freezes the selected bytes into a release descriptor. Adding a file here without
an inventory entry does not distribute it. The profile identifier and configuration
filenames remain stable; relocating authoring templates does not migrate consumers.

[Go installation records](../SDPTool/install/Records.md) define preview, apply,
recovery and ownership. `gh sdp` consumes a selected published release, not this
checkout. Templates initialize missing project files; upgrades preserve existing
project prose. Reconcile local documentation explicitly when adopting newer rules.

## Legacy

[legacy/install-v1](legacy/install-v1/README.md) retains both old template roots
for the legacy installation contract and compatibility checks. It is not an
alternative recommendation for new projects. That inventory now points to the
archive, while the current inventory points only to sdp-root. Historical evidence
and already published release bytes remain unchanged.
