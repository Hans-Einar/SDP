# Release Lifecycle And Gate

## Deterministic normal-release gate

A normal release candidate fails unless all of the following are true:

- included Slices or Fixes are complete and no incomplete active Slice is claimed
- required verification records pass
- blocking, high and medium review findings are resolved
- YAML, JSON, NDJSON and traceability relations parse and agree
- release notes are complete and the target section is frozen
- Toolkit, project, package and application versions agree where applicable
- the working tree is clean
- the release-preparation commit exists
- the intended tag does not already exist
- required migration notes exist
- the target is greater than the previous release
- immutable released history has not been changed silently
- governed record paths use canonical `.yaml` names and resolve to real files
- Slice/Fix review and verification, review-resolution, and release/migration links are reciprocal
- release record IDs, versions, state, dates, tags and commits agree
- each release event subject and tag matches the governed release identity

The gate emits evidence; it does not infer missing evidence.

## Variants

| Release kind | Additional or adjusted rule |
|---|---|
| Prerelease | Use SemVer prerelease syntax; clearly mark it non-final; all safety-critical checks still pass. |
| Patch/hotfix | Must be backward-compatible; a Fix record may supply the work contract, but the release remains a normal PATCH release. |
| Yanked | Never delete history; record reason, affected versions and replacement guidance, then mark state `yanked`. |
| Documentation-only | May omit product runtime tests only when no executable or public contract changes; documentation and link validation still pass. |
| SDP Toolkit | Installation-contract/plan schemas, installer migration fixtures, skill/manifest agreement, Framework compatibility and an extracted-source-archive check are mandatory. |
| Consuming project | Package/application version agreement and project-specific acceptance evidence are mandatory. |

## Truthful two-phase publication

### Phase A — prepare

1. Create a dedicated release-preparation Slice.
2. Freeze the selected release-note entries.
3. Set the manifest to `prerelease`; leave `gitTag` and `releaseCommit` null.
4. Complete verification, review and approval events.
5. Create a clean release-preparation commit.
6. Stop unless publication is explicitly authorized.

### Phase B — publish and reconcile

After authorization:

```bash
git tag -a vX.Y.Z <release-preparation-commit> -m "SDP Toolkit vX.Y.Z"
git push origin <release-branch>
git push origin vX.Y.Z
gh release create vX.Y.Z --verify-tag --notes-file <frozen-notes-file>
```

Only after those commands succeed, append `release-tag-created` and
`release-published` events and commit a small reconciliation change that sets the
manifest state to `released` and records the real tag and release commit.

This follow-up commit is intentionally after the tagged commit. A commit cannot
truthfully contain evidence that a future tag or GitHub Release already exists,
and it cannot contain its own SHA. The annotated tag and GitHub Release identify
the released source; the reconciliation commit records the completed transaction.

## Toolkit source-archive gate

Before publishing an SDP Toolkit release, verify the candidate from a normally
extracted GitHub source archive in a temporary path with no `.git` directory.
The installation manifest and every schema/source reference must resolve from
the archive root, the deterministic plan must not depend on Windows separators
or PowerShell interpretation, installation must succeed, and generated installed
facts must record `sourceCommit: null` when no trustworthy commit is available.

The normal GitHub source archive is the selected distribution artifact. Do not
add a custom release asset unless concrete verification shows that the normal
archive cannot satisfy the contract safely.
