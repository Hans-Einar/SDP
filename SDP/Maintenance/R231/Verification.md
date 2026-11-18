# SDP2.3.1 release qualification

| Field | Value |
| --- | --- |
| id | SDPTOOL-VER-R231 |
| Status | completed |
| Activity | ACTIVITY-SDP-0005 |
| Plan | MAINT-SDP-0020 |

**Completed and independently approved:** signed public2.3.1, actual local tool
upgrade and HEOS2.2.0-to2.3.1 upgrade with preserved owner work.

## Preparation and qualification checkpoints

Canonical candidate84d5e46012fdc45a3690bd0b26b617501fd1caa7 is clean in the release
worktree. Engine SHA25604f78532a4e7ae8653df5b1e623e1407312713a90043f2ac36e7f68e0dc0495b;
signed descriptor9c83225b3a5eae1762e5609965378b8c14d89fa676ae242856b323c1dd9facf6.
gh-sdp packaging sourceae7aa216c4343c6d4904078c8545cd1d2d8f7a64 republishes identical
engine bytes. Trusted Linux Ed25519 verification, all checksum entries and all73
payload content hashes pass independent inspection. Native helper SHA256 is
identical to the qualified2.3.0 payload; no repeat native GUI campaign is claimed.
Engine identity is explicit BuildRevision plus signed descriptor; embedded Go VCS
metadata is absent on this host. The complete descriptor bytes are release assets.

Evidence/: focused profile/install race tests, bootstrap tests, actual packaged
signed child,90 Toolkit tests, all9 release logs in source and Git-less extraction,
extension tests and all8 original signed predecessor upgrades pass. Each matrix
case checks owner-file preservation, managed conflict refusal and repeat no-op.
Existing2.3.0 broad/native evidence remains applicable to unchanged code. CI run
43544849624 uses synthetic merge8a552d9260cbd0efddeea16da04bde326cefd8c7; its tree
0493a4a18ce1766ecb4af3918c34335a9607d1d4 equals exact84d5e46 tree. Contracts pass;
the broad race suite has passed and the final packaging job remains pending. This is source-tree-equivalent CI,
not a claim that the CI build revision equals the published executable identity.

The initial b428a3a candidate omitted an immutable Session creation snapshot.
CI and independent review refused it. It is superseded, not published; the next
commit adds that exact snapshot and clean-source management validation passes.

## HEOS rehearsal and required explicit exception

Original HEOS is signed2.2.0 and contains existing dirty owner documentation.
The initial isolated-cache probe correctly required the original signed descriptor;
resolving the public2.2.0 through the signed installer populated that cache.
The installer then correctly rejected five malformed historical ledger fields.
A proposed copied-tree correction changes only two KanBan destination states from
in-progress to active, one management event from updated to started, and two
KanBan event types from completed to moved. Identities, dates, reasons, links and
actual recorded outcomes remain unchanged. Independent review confirms semantics.
The original live ledger remained byte-identical during review. The owner then
explicitly approved the exact five-field correction and archival exception. The
original is now permanently retained byte-for-byte in HEOS Maintenance evidence;
only those five fields were changed. The corrected live Framework still remains
2.2.0 pending public release and actual installer apply.

After only that proposed correction, copied HEOS upgrades with14 actions. All
existing owner files are byte-identical, corrected ledger bytes remain a prefix,
signed receipt identifies2.3.1, and repeat preview is a no-op. This proves the
proposed migration, not an applied live upgrade. The patch and full original are
retained outside the producer repository at /tmp/sdp-r231-final-heos.

## Gates recorded before publication — now complete

Final source-equivalent CI, public canonical/extension publication and downloaded
identity/signature verification, real public extension update/fresh install,
public-default Framework check, then actual HEOS upgrade under the now-approved historical exception. Independent review must observe actual outcomes before
claiming completion. Main merges, XFMD rebuilds and other projects are unselected.

## Additional extracted-source gate

Final reviewer required more than release-log checks for Git-less operation.
The actual git-archive84d5e46 tree passed Toolkit and management validation.
Its Go profile builder accepted explicit sourceCommit84d5e46 and the exact packaged
binary. The resulting development descriptor's files, binaries, profiles,
capabilities, protocol, predecessor list and sourceCommit exactly equal the signed
production candidate. Its development identity was not installed. A separate fresh
install using the actual production-signed descriptor produced the signed2.3.1
receipt at84d5e46. evidence/extracted-validation.log and extracted-install-result.json
retain results. No broad tests or source rebuild were needed for this added gate.

## Live rollout observed

Actual public-default HEOS preview/apply now succeeds with signed2.3.1 receipt and
repeat no-op. The local extension and registered standalone executable share the
exact published engine. The harness initially rejected its own overly strict
assumption that all14 actions had managed ownership, before any live apply.
Inspection confirmed9 managed writes plus5 initialize-if-missing schema files,
all absent in the prior tree. The corrected harness allowed only those explicit
ownership cases and resumed the saved apply; no installer change or weakened
product validation was involved. Installed bytes and owner preservation pass.
Final independent live review and durable closeout are next.

## Final observed delivery

All previously pending release and consumer gates are complete. CI43544849624
passes contracts and full race/package/signed-child jobs on the independently
verified identical source tree. Canonical and extension v2.3.1 are published;
all actual downloaded assets match API digests and reviewed local bytes. Independent
review verifies public production trust and exact source/engine identity.

Public ordinary extension upgrades from both0.2.3 and2.3.0, plus clean installation,
pass. The2.3.0 setup link resolves to the new engine before any new setup call.
Direct offline version/discovery/SDUI operations pass without gh or auth/network.
Public-default Framework install/apply and repeat no-op pass with preserved owner
files. The live local extension was then upgraded and registered through gh.

HEOS actually upgraded from signed2.2.0 to signed2.3.1 through public-default
`gh sdp` preview/apply. All9 managed-file updates and5 previously absent schema initializations match signed content;
all existing owner files are unchanged by the installer, corrected history remains
a byte prefix, and repeat preview is a no-op. The explicitly owner-approved
historical exception is archived and documented in HEOS Maintenance evidence.
Subsequent consumer README/Session maintenance is distinct from installer preservation.
Existing dirty HEOS work is retained and not committed as part of producer delivery.

No other project's Framework, main merge, XFMD rebuild or automatic Activity/PG
execution adoption is claimed. MAINT-SDP-0020 and Activity5 have no remaining
selected delivery obligation. Public and live receipts are in public-evidence/.

## Draft closeout record correction

The newly authored uncommitted Maintenance completion event initially linked a
Release and Activity, but this local profile allows management/card identities in
that field. Validation refused the draft. The independent reviewer confirmed that
correcting this new unpublished record is authorized closeout work; it is distinct
from changing pre-existing HEOS history. Only payload.links was corrected to the
actual preceding delivery plan/card. Original draft ledger bytes are retained
compressed with DraftLinkCorrection.json. Every prior committed ledger byte and
all actual Activity operation event bytes are unchanged. The operation receipts
retain their actual observed history revisions; they are not rewritten to conceal
the draft correction. Cross-domain release/Activity links remain in maintained
documents and Traceability. Final validation below applies to the corrected draft.
