# Distribution and upgrades

The current distribution is a signed Go SDPTool executable plus immutable release
descriptor. gh-sdp is a thin local bootstrap/client: it verifies/downloads the
selected release, then invokes the compiled tool. GitHub does not run the user's
installation on a remote machine.

The descriptor is built from SDPTool/profiles/five-phase.json and
SDPTool/profiles/payload.json. It binds source identity, destination paths, content
hashes, ownership, supported predecessor digests and platform binary hashes.
Installed receipts record version and descriptor identity; authoritative inventories
remain publisher-owned. Published bytes must not be overwritten.

Upgrade preview resolves the signed predecessor and inspects the actual target.
Applying the saved plan verifies its inputs and records backups, journal, receipt
and Maintenance report. Managed files and initialize-if-missing project templates
have different policies. Existing project-owned prose is preserved. An unknown
manual baseline requires explicit adoption rather than guessed overwrite rights.

Read [installation records](../../SDPTool/install/Records.md) and
[migration commands](Installer-Migration.md). Releasing the payload and changing
gh-sdp's immutable default are separate deliveries.
