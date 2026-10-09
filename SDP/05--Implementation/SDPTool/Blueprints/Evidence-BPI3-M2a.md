# BPI3-M2a — Canonical history append foundation

Date: 2026-10-09. Plan PLAN-SDP-0020; card KB-SDP-050; Session0008 T011.
Candidate: baseline 7fddedb plus the source hashes in
[Pilot-BPI3-M2a.json](Pilot-BPI3-M2a.json). Tests were run on these uncommitted
source bytes, not retrospectively claimed against a later commit.

## Delivery and rationale

The inspected governance worktree at 4044f37 owns private operational snapshots
and reads selected canonical history subjects; it has no canonical ledger writer.
The M2a library establishes that shared prerequisite without copying its private
store or creating a second assignment ledger. M2b will adopt it for real lifecycle
commands and extend domain validators before publishing new event types.

SDPTool/projecthistory provides bounded Read and revision-bound Append. Existing
bytes and permissions survive; stale writes and ID conflicts are refused. Exact
retry remains idempotent even after later events. Linux writers lock a stable
sidecar, stage/fsync a complete stream, recheck, rename and sync its directory.
The SDL overview declares only its library interface, not a deployed container
or a consuming lifecycle unit. See the [assignment design](../../../04--Design/SDPTool/Blueprints/Assignment-Lifecycle.md).

## Verification

Go 1.27.1 linux/amd64 from the local pinned toolchain. From SDPTool:

- `go test -race ./projecthistory -count=1`: pass, 1.791 seconds.
- `go vet ./projecthistory`: pass.
- Independent reviewer repeated race tests (1.579 seconds) and vet: pass.
- SDL model check: pass after canonical formatter reordered the new declaration.
- `python3 SDP/ProjectManagement/validate.py`: pass (62 cards, 39 management
  records, 4 lineage operations, 515 events).
- `python3 Toolkit/scripts/validate_sdp.py`: pass.
- `git diff --check`: pass.

Tests exercise actual historical envelope reads without mutating repository
history; old-byte/permission preservation; duplicate retry after later events;
CAS and content conflicts; concurrent writers; malformed/torn/duplicate streams;
size and symlink rejection; cross-process lock refusal; forced process exits
before and after publication/sync; and retry after committed-but-unacknowledged
publication. These verify the reusable component, not a production assignment flow.

Independent review identified missing generic constraints for optional subjectId
and releaseId, plus invalid UTF-8 acceptance. Fixed with exact schema constraints,
an explicit UTF-8 check and negative fixtures. Reviewer approved corrected source
with no remaining material findings; hashes are in the machine record.

## Limits and next step

No command consumes the new library yet. M2 remains unfinished until assignment
operations, authority/evidence rules, canonical validator support and grouped
discovery pass integrated M2b acceptance. No actual blueprint lifecycle event was
appended to the project ledger in this milestone.

Mutation is Linux-only; other platforms return ErrUnsupported. Arbitrary/manual
writers do not share the lock; the final recheck detects observed interference
but cannot eliminate a later uncooperative write. Forced process termination
proves tested filesystem behavior, not universal physical power-loss durability.
No full regression or native XFMD/release acceptance is claimed. The previously
recorded installer race timeout remains open for final integration.
