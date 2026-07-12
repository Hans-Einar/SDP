# Development Identity

## Canonical fields

Generated build metadata keeps these fields separately addressable:

- `releaseVersion`
- `releaseState`: `unreleased`, `prerelease`, `released` or `yanked`
- `sprintId` or `refactorId`
- `iterationId`
- `sliceId`
- `fixId`
- optional `revision`
- full and short Git commit SHA
- UTC build timestamp
- dirty-working-tree flag

The generated JSON schema is independent of React and Vite. Vite applications
may import the JSON; other systems may read it at build time, embed it in a
resource or expose it through an API.

Example UI text:

```text
0.8.0-dev · Sprint-026 / SPI-002 / SPS-003 · r1 · abc1234
```

`-dev`, `-prerelease` or `-yanked` text is mandatory whenever the build is not
a current official release.

## Revision rule

A revision is allowed only within the same already-planned Slice or Fix when a
small bounded correction is needed after verification or review. It must not:

- add a capability or requirement
- change architecture or a public contract
- broaden the expected file or component boundary materially
- hide work that deserves a new Slice or Fix record

Each revision has verification and an append-only ledger event. It resets when
the active Slice or Fix changes.
