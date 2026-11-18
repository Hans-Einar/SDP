# Independent SDP2.3.1 release review

Status: changes required before publication.

Candidate b428a3aa00237e8fb1dfd635d2f2f39ce0106553, inspected in clean detached /tmp/sdp-r231-source against integrated main cef86e4 and released v2.3.0. Loaded applicable AGENTS.md, SDP entrypoint, Reviewer, Verifier, Release and document-workflow reference; reviewed R231 plan, release lifecycle, checklist, signing contract and Activity5 Session22. Scope is independent release review only; no product or lifecycle edits.

## Finding R231-1 — High — missing immutable Session snapshot

The exact clean candidate fails `python SDP/ProjectManagement/validate.py` with `ValueError: Activity record unavailable: Activities/.sessions/c921df5a36832a92b55ad9e39279e3b88779adcfa34fee121768a904cc49fb9b.md`. The Activity Session creation history references absent immutable evidence. The same candidate CI43544653325 contracts job fails at that validator and skips subsequent Toolkit tests. This violates the normal release gate requiring agreement of management records and passing required validation. Add the original referenced immutable snapshot without rewriting historical event bytes; rebuild/sign the resulting exact candidate and obtain relevant validation/qualification before publication.

## Inspected evidence

The executable changes are bounded to build version/default descriptor URL and one additive predecessor digest/test. Runtime/native implementation, distributed payload and Skills sources are unchanged from v2.3.0, supporting reuse of unchanged native behavior qualification. Package manifest records 2.3.1 and exact b428a3a revision. All five SHA256SUMS entries pass. Independently verified the exact descriptor Ed25519 signature with OpenSSL using the documented Linux public key/keyId 2c947457cffe015f870168cb332242c49606fdea378a13d875579544c2989c03. Descriptor digest b181442eaac987ff47082821710b982286dc38b1a502f7587893915e9051b8e1 pins this source and engine SHA256 6f7bafcd545f9135dbcc3713abe1aeed5750771cc28bc1e4da95b955ea35f5b3.

Inspected /tmp/sdp-r231-verify-predecessors.py, predecessor logs and saved plans/receipts for all eight 0.2.0–2.3.0 original signed descriptors. The script exercises initial install, managed-file conflict rejection, 2.3.1 preview/apply, owner instructions/Session/SDL preservation and repeated no-op with production trust and cleared test-key overrides. Local profile/install and bootstrap logs pass.

## Pending review boundary

Candidate correction, applicable successful CI (including checkout identity distinction), extracted-source gate, packaged signed child, copied HEOS rehearsal/preservation, extension exact-engine repack, final reciprocal verification/review links and populated checklist remain to inspect before a publication-ready verdict. Actual publication and live HEOS upgrade require later observed evidence; they are not approved or claimed here.

## Replacement candidate review

Candidate 84d5e46012fdc45a3690bd0b26b617501fd1caa7 supersedes b428a3a. Its sole additional source change adds the exact missing immutable Session snapshot. Clean detached checkout and successful management validation independently observed. Finding R231-1 is resolved for this candidate.

Independently verified final package in /tmp/sdp-r231-final-package: all five checksum entries pass; OpenSSL verification against the documented Linux production public key succeeds; descriptor SHA256 9c83225b3a5eae1762e5609965378b8c14d89fa676ae242856b323c1dd9facf6 identifies 84d5e46012fdc45a3690bd0b26b617501fd1caa7 and engine SHA256 04f78532a4e7ae8653df5b1e623e1407312713a90043f2ac36e7f68e0dc0495b. All 73 embedded payload file hashes agree with content. Executable manifest reports exact revision and2.3.1.

Final predecessor evidence /tmp/sdp-r231-final-predecessors passes all eight original signed releases. Independently checked each saved install plan embeds signed original descriptor bytes whose digest belongs to target upgradesFrom; every upgraded receipt pins the final descriptor, upgrades to2.3.1, and repeat plan reports noChange. Final matrix supersedes first-candidate artifact results.

CI43544849624 contracts job passes management validation and90 Toolkit tests. Checkout logs identify synthetic merge 8a552d9260cbd0efddeea16da04bde326cefd8c7 merging candidate84d5e46 into unchanged main cef86e4; headSha alone must not be described as exact checkout. Go installation job remains pending at this observation. No new product finding; publication readiness remains conditional on outstanding applicable CI and release/consumer evidence listed above.

## Consumer qualification and HEOS historical-data exception

Reviewed final extension distribution source ae7aa216c4343c6d4904078c8545cd1d2d8f7a64: changed files contain release documentation/records only. Local package /tmp/sdp-r231-extension-local manifest identifies final canonical descriptor and engine. Independently compared gh-sdp-linux-amd64 with /tmp/sdp-r231-final-package/sdptool: bytes identical. Focused extension tests pass. Public distribution verification remains subsequent evidence.

Original live HEOS upgrade correctly refuses invalid management history (`history`, `card path/state mismatch`, exit3). This is pre-existing consumer data, not a new installer defect. Independently confirmed live ledger equals /tmp/sdp-r231-final-heos/original-ledger.ndjson byte-for-byte. Proposed correction preserves all64 event identities, timestamps, reasons, paths, links and ordering; exactly five scalar fields change: two KanBan payload.to values in-progress→active, management eventType updated→started for its existing planned→active transition, and two KanBan eventTypes completed→moved retaining destination completed. These corrections match the installed history grammar and stated historical transitions.

The corrected copy result /tmp/sdp-r231-final-heos/corrected-result.json identifies the final signed candidate and descriptor. Its fresh preview has14 applicable actions; receipt is signed2.3.1; repeated preview is noChange. Independently compared every pre-existing owner snapshot path excluding declared installer-managed paths, receipt and ledger: zero unexpected changes. Corrected ledger is an exact prefix of post-upgrade ledger. This establishes upgrade readiness for the specifically corrected baseline, not the original invalid baseline.

Live historical edits remain subject to the owner's explicit append-only exception and byte-identical archive/documented correction. No such permission or live application is inferred by this review. Candidate has no unresolved product findings; release readiness remains conditional on applicable successful installation CI, remaining extracted-source/packaged/public qualification and truthful final checklist/review links. Live HEOS completion additionally requires that owner exception, freshly reviewed live plan and observed preservation evidence.

The CI synthetic merge and release candidate source trees were independently compared: both tree0493a4a18ce1766ecb4af3918c34335a9607d1d4. CI source-equivalence therefore holds while commit/build identities remain distinct.

## Authorized HEOS historical correction

Owner explicitly approved the exact five-field correction with archived-original preservation and subsequent HEOS upgrade, per current owner instruction supplied to reviewer. Independently inspected live branch heos/sdp-upgrade-2.3.1, correction metadata and archive. SDP/Maintenance/SDP-Upgrade-20261114/evidence/Ledger.original.ndjson is byte-identical to the original independently reviewed ledger snapshot. Current live ledger begins with the exact proposed five-field corrected ledger. This resolves the prior permission boundary for that specific historical exception. It does not authorize additional rewrites.

Actual Framework remains2.2.0 at this observation; live2.3.1 upgrade and public CI/distribution qualification remain pending. Manual Session recording under the existing guide does not constitute Activity admission adoption.

## Final CI observation

Independently queried CI43544849624 after completion: overall success, contracts success and go-installation-linux success. All workflow steps passed, including complete race suites, all release-log checks, package build, actual packaged signed child, descriptor build and retired-entrypoint rejection. Source-tree equivalence remains as established above. Clean release source84d5e46 and clean extensionae7aa21 unchanged. Native helper digest50cdfabb3ab213eecdada79c0ca661528bb19a2e097ccb357bd285b4bda520f8 exactly equals qualified2.3.0 helper.

The remaining prepublication clarification is the mandatory extracted-source operation evidence: current extracted log proves only nine release logs. Release-Lifecycle additionally requires archive validation, explicit exact-sourceCommit Go descriptor build and installed receipt identity. Reviewer requested the narrow evidence or check before issuing the final publication-ready disposition. This is an evidence gap, not an identified product defect. Public asset verification and live HEOS operation remain postpublication obligations.

## Publication readiness disposition — approved

The narrow extracted-source gate is now satisfied. Independently inspected /tmp/sdp-r231-extracted-validation.log (Toolkit and management validators pass), Git-less extraction, extracted descriptor and actual signed installation receipt. The extracted descriptor pins exact84d5e46012fdc45a3690bd0b26b617501fd1caa7; comparison with production descriptor differs only in intentionally distinct development-extracted release identity. All payloads, binaries, capabilities, predecessor support, profiles and protocol match. Actual install uses the production signed descriptor, with signed2.3.1 receipt and correct exact sourceCommit/descriptor digest; development descriptor is not installed. Independently compared extracted files against git archive of exact candidate.

Approved for authorized publication: canonical source84d5e46012fdc45a3690bd0b26b617501fd1caa7, extension sourceae7aa216c4343c6d4904078c8545cd1d2d8f7a64 and previously verified final local signed package. R231-1 resolved; no remaining prepublication product or evidence finding. Complete applicable durable checklist/review references with this observed disposition. This independent approval does not authorize main merges or expand existing owner authority.

Publication itself, downloaded public canonical/extension identity and signature verification, actual extension update/fresh install and live HEOS preview/apply/preservation remain subsequent checks. Do not claim release rollout complete before those observations.

## Published distribution and consumer review

Independently queried live GitHub annotated tag objects: canonical v2.3.1 tag71754eea1dafb6e021f74fb392cb9066ee8389a6 targets exact84d5e46012fdc45a3690bd0b26b617501fd1caa7; extension v2.3.1 tag0c7af0e573ad7d41cafbab98e8969c5aef4b9251 targets exactae7aa216c4343c6d4904078c8545cd1d2d8f7a64. Live releases are public final releases. Independently compared all six downloaded canonical assets and seven downloaded extension assets to their live API sizes/digests and qualified originals. All match; public extension engine is canonical exact bytes. Independently reverified downloaded descriptor production Ed25519 signature and all canonical SHA256SUMS.

Inspected consumer scripts, saved commands and results at /tmp/sdp-r231-public-consumers and /tmp/sdp-r231-public-framework. Normal public upgrades from0.2.3 and2.3.0 and fresh public install produce exact2.3.1 engine/revision. The2.3.0 update preserves existing setup registration before new invocation. Direct operations succeed in isolated network namespaces with gh absent from PATH, absent credentials and empty cache; owner project remains unchanged. Public-default Framework install uses no release override or test key and records signed2.3.1 final descriptor receipt; repeat preview is no-op and owner file preserved.

Published distribution and disposable consumer qualification approved. Actual live extension update/registration and HEOS preview/apply/preservation remain final rollout observations; no main merge or other-project upgrade is implied.

## Final combined delivery disposition — approved

Actual rollout evidence /tmp/sdp-r231-live and current live files independently inspected. Live gh extension upgrade/setup, gh version and direct version commands succeed with2.3.1/exact84d5e46. Actual /home/warloc/.local/bin/sdptool resolves to the registered gh extension engine and hashes to04f78532a4e7ae8653df5b1e623e1407312713a90043f2ac36e7f68e0dc0495b.

HEOS public-default fresh preview/apply succeeds from signed2.2.0 to signed2.3.1. Actual receipt pins exact source84d5e46012fdc45a3690bd0b26b617501fd1caa7 and descriptor9c83225b3a5eae1762e5609965378b8c14d89fa676ae242856b323c1dd9facf6. All14 actions are justified: nine managed payload writes and five initialize-if-missing schema writes whose paths are absent in the recorded before snapshot. Independently compared all actual managed/action bytes against signed descriptor: all match. The harness correction clarified action ownership before application; it did not weaken product validation or change the installer.

Independently compared existing owner files to actual pre-upgrade snapshot: zero unexpected changes. Approved corrected ledger remains an exact prefix; original archived ledger remains byte-identical. Repeat preview reports noChange. Actual command log records successful default preview/apply, registration and repeat checks. Installing the new Activity framework/schema files does not establish execution-admission adoption.

Final combined publication and first-project delivery approved with no unresolved review findings. Canonical and extension publication, disposable public consumers, actual local extension registration and live HEOS upgrade have now been observed. Durable lifecycle/checklist reconciliation may record these outcomes truthfully; owner acceptance and any main merges remain distinct and are not inferred by independent review.

## Closeout draft-link correction disposition

Reviewer advised that correcting newly authored, uncommitted invalid EVT-PM-SDP-000267 metadata falls within authorized closeout work and requires no additional owner exception. This differs from altering pre-existing HEOS history. Applicable rules preserve historical records and prohibit rewriting history to conceal corrections; they do not explicitly freeze a failed unpublished draft.

The advised boundary is: retain the invalid draft and actual operation receipts; correct only payload.links to meaningful existing management/card IDs (PLAN-SDP-0029 and KB-SDP-063); document the correction; preserve committed ledger prefix bytes and actual Activity operation events unchanged; rerun validation. Cross-domain release/Activity relationships remain in ordinary documents/Traceability. No validator or product behavior change is warranted.

Coordinator reports this correction applied exactly, with original full current-turn ledger archived as gzip and DraftLinkCorrection.json under Maintenance/R231/evidence. This report records the scoped advice and reported disposition; final local validation is still running at this observation. Earlier delivery approval remains valid, with lifecycle closeout contingent on successful final record validation.
