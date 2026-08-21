# Repository Inventory

## Result

- **INFERENCE — fixed study method:** the owner's four-month “actual commits”
  criterion is measured in the exact inclusive interval
  `2026-04-20T21:08:52Z` through `2026-08-20T21:08:52Z`, anchored to Issue #5's
  creation time and using commit *committer* timestamps.
- **OBSERVED:** an authenticated owner listing and Git mirror capture completed
  at `2026-08-20T22:13:07Z` considered 35 repositories: 27 private and 8 public;
  3 are owned forks and none are archived.
- **OBSERVED:** all 35 repositories and their advertised refs were accessible.
  No repository was skipped because of access or an empty/missing default branch.
- **OBSERVED:** 17 repositories have one or more qualifying commits and 18 have
  none.
- **INFERENCE:** the 17 qualifying repositories are the exhaustive report set
  for this study under the owner-directed rule. `Hans-Einar/HSX` is included even
  though it was absent from Issue #5's preliminary `pushed_at` discovery list;
  its non-default advertised refs contain 104 qualifying commits.

## Derivation method

**OBSERVED:** GitHub CLI was authenticated as `Hans-Einar` with `repo` and
`read:org` scopes. The owner listing was obtained without `--source`, so owned
forks were not silently discarded:

```powershell
gh auth status
gh repo list Hans-Einar --limit 1000 `
  --json name,nameWithOwner,url,sshUrl,isPrivate,isArchived,isFork,defaultBranchRef
```

For each listed repository, a fresh filtered mirror was fetched. Filtering
omitted file trees and blobs, not refs or commit objects:

```powershell
git clone --mirror --filter=tree:0 `
  git@github.com:Hans-Einar/<repository>.git `
  <temporary-root>/<repository>.git
```

The classification enumerated the complete reachable commit set first and then
compared every committer timestamp after normalizing it to UTC. It did not rely
on Git's date-limited traversal pruning:

```powershell
git -C <mirror> for-each-ref --format='%(refname) %(objectname)'
git -C <mirror> log --all --format='%H%x09%cI'
git -C <mirror> rev-parse refs/heads/<advertised-default-branch>
git -C <mirror> for-each-ref --contains=<evidence-commit> --format='%(refname)'
```

The inclusive predicate was:

```powershell
$start = [DateTimeOffset]::Parse('2026-04-20T21:08:52Z')
$end   = [DateTimeOffset]::Parse('2026-08-20T21:08:52Z')
$utc   = ([DateTimeOffset]::Parse($committerIso8601)).ToUniversalTime()
$inScope = $utc -ge $start -and $utc -le $end
```

`Refs/commits` below is the count of captured advertised refs and the count of
unique commits reachable from them. `Study commit` is always the captured tip of
the advertised default branch. `First` and `last` are the chronologically first
and last qualifying reachable commits; their containing ref proves reachability.
When several refs contain the same commit, a branch was preferred over a tag and
a tag over a pull-request ref, followed by lexical order.

## In-scope repositories and report manifest

**OBSERVED:** every row has at least one qualifying commit. **OWNER DIRECTION:**
each row maps to exactly one later repository report. The activity evidence is
separate from the default-branch study commit, which may be older or newer than
the activity ref.

| Repository | Visibility | Default branch | Exact default-branch study commit (committer UTC) | Refs/commits | Window commits | First qualifying commit (UTC; containing ref) | Last qualifying commit (UTC; containing ref) | Planned report |
|---|---|---|---|---:|---:|---|---|---|
| `Hans-Einar/ActionCam` | private | `main` | `861a22d4b2336e1eabd24051b914c8d24219ee01` (`2026-08-20T18:32:55Z`) | 60 / 189 | 189 | `0d483beb29d5d6df201adc651a8c7ba83dc83d47` (`2026-08-12T22:37:08Z`; `refs/heads/agent/issue-12-camera-deploy`) | `d36196f1a0e2b3a82d74680a872f1db790a13011` (`2026-08-20T20:13:00Z`; `refs/pull/65/merge`) | `repositories/ActionCam.md` |
| `Hans-Einar/gh-sdp` | public | `main` | `32613734781bf39f2fce176db2acfb2284dfc92f` (`2026-07-15T06:23:24Z`) | 8 / 19 | 19 | `2dcc18423a10497678eb19ace8fd9668efda9551` (`2026-07-13T09:17:53Z`; `refs/heads/codex/phase-1-sdp-bootstrap-mandate`) | `bd3aa4d7713e9fe461016a06eb89a14ee77d8d03` (`2026-07-17T13:32:38Z`; `refs/pull/3/merge`) | `repositories/gh-sdp.md` |
| `Hans-Einar/GrassPhenology` | private | `master` | `8ba7567aaaf0524d86bdb4675a09232b41496f36` (`2026-07-14T07:37:07Z`) | 6 / 22 | 22 | `8dd3011c240ad349b9f7f16985245bf7d02bc082` (`2026-07-11T08:09:46Z`; `refs/heads/agent/add-radar-precipitation-wishlist`) | `f54d2a7a6e8338a9ad83ae679914b7c366b6c1c7` (`2026-07-22T20:15:06Z`; `refs/pull/2/merge`) | `repositories/GrassPhenology.md` |
| `Hans-Einar/HEOS` | private | `main` | `c41995ecc503684f4bb5200ae8f8da3665d57e1c` (`2026-04-25T20:30:33Z`) | 15 / 52 | 4 | `9b5cd82315dd6a06edf2c08a99cbcedc8f48fdcf` (`2026-04-24T13:15:39Z`; `refs/heads/Version_2.0`) | `8d4bddc04bace27d695f0faa0c320bba1b323390` (`2026-07-14T06:32:00Z`; `refs/heads/Version_1.0--Union`) | `repositories/HEOS.md` |
| `Hans-Einar/HSX` | public | `main` | `e374da88f4dd470bad2d8ec6a2a14f1ce367e40e` (`2025-11-09T17:55:00Z`) | 82 / 544 | 104 | `de84aafd732ec67f53f322481f75d77a17940a0a` (`2026-08-19T15:31:33Z`; `refs/heads/codex/dbg-da-001`) | `86aac69ac89f20d8ccce0fe83fd15af0daed8d9d` (`2026-08-20T20:55:36Z`; `refs/heads/codex/dbg-rf-002-003`) | `repositories/HSX.md` |
| `Hans-Einar/LogClassifier` | private | `main` | `a38f9bb061f8ab77e25b801c06915c754d746e74` (`2026-07-11T13:55:20Z`) | 1 / 2 | 2 | `a4811d6b6719ad70c8e0de6df5813b805b6430ba` (`2026-07-11T09:37:50Z`; `refs/heads/main`) | `a38f9bb061f8ab77e25b801c06915c754d746e74` (`2026-07-11T13:55:20Z`; `refs/heads/main`) | `repositories/LogClassifier.md` |
| `Hans-Einar/LogParser` | private | `main` | `197336b024aa5da8c773fd3076033ee76a099872` (`2026-07-13T21:17:20Z`) | 1 / 2 | 2 | `0d00ca68e341fe0b0787e3a4e37336faa5fc7a48` (`2026-07-13T21:17:18Z`; `refs/heads/main`) | `197336b024aa5da8c773fd3076033ee76a099872` (`2026-07-13T21:17:20Z`; `refs/heads/main`) | `repositories/LogParser.md` |
| `Hans-Einar/Lyndata` | private | `main` | `132a72f2c706ee6de32dae0f90491f5ac11d1aa3` (`2026-08-05T16:21:20Z`) | 18 / 167 | 165 | `2afe7d4ef36ef79155d27c9fdb519b223a916acb` (`2026-07-30T16:01:23Z`; `refs/heads/codex/boot-001-feature-first-sdp`) | `9ad04a848e5a2d2219c0f709db91abcdd1ad3fae` (`2026-08-07T11:49:16Z`; `refs/heads/codex/interactive-explorer-acceptance-revision`) | `repositories/Lyndata.md` |
| `Hans-Einar/map_tracker` | private | `master` | `22e03b49a27a7e888ee9dbdb97874e2038b43fb4` (`2025-12-31T22:14:55Z`) | 3 / 23 | 4 | `3367e01aec45057815cfcba6707f3ff71ae06f91` (`2026-05-09T13:41:34Z`; `refs/heads/V2`) | `19cb24d9ccf7d3d1d1ea8c5654de1e683e7574d5` (`2026-05-10T09:56:19Z`; `refs/heads/V2`) | `repositories/map_tracker.md` |
| `Hans-Einar/python` | private | `main` | `6bd5923ef6ab3041d853a66ed3b63ebbe489aa18` (`2026-06-19T22:23:24Z`) | 3 / 41 | 2 | `49b3c0a959b0cc5ae962c4e44dc015b683fee89a` (`2026-06-19T22:22:08Z`; `refs/heads/main`) | `6bd5923ef6ab3041d853a66ed3b63ebbe489aa18` (`2026-06-19T22:23:24Z`; `refs/heads/main`) | `repositories/python.md` |
| `Hans-Einar/RadarData` | private | `main` | `ca9129209eb2eb5d10854b7043f8255a69e90953` (`2026-07-26T09:34:38Z`) | 7 / 84 | 84 | `f2e0f918f46bcd5da50d459e90b578906e25de27` (`2026-07-25T10:57:27Z`; `refs/heads/agent/add-radardata-poc-concept`) | `d84a5718c5ddc307efb72cddb9963149b286d64b` (`2026-07-26T23:59:05Z`; `refs/pull/2/merge`) | `repositories/RadarData.md` |
| `Hans-Einar/SDP` | public | `main` | `e398ebaf3a4ace6a5d92fd9ce22736a7427a9e15` (`2026-07-17T12:59:57Z`) | 12 / 108 | 105 | `2eef0fa98a146ecfc345ffde5041ba2e03adae6e` (`2026-07-11T09:59:06Z`; `refs/heads/agent/sdp-release-versioning`) | `adcad55319205f8aa27ceeb8fe3f87620e8c2d5d` (`2026-07-17T15:27:23Z`; `refs/pull/4/merge`) | `repositories/SDP.md` |
| `Hans-Einar/SDP-Analyzer` | public | `main` | `632991a878100e8cd8c4efbb7d724edb3694d98a` (`2026-07-17T12:59:20Z`) | 1 / 30 | 30 | `122407e4f0cea8c32305d2fc3e7910d66e30422f` (`2026-07-11T13:11:10Z`; `refs/heads/main`) | `632991a878100e8cd8c4efbb7d724edb3694d98a` (`2026-07-17T12:59:20Z`; `refs/heads/main`) | `repositories/SDP-Analyzer.md` |
| `Hans-Einar/SharedUI` | private | `main` | `8cb61651107ff4fa23ad283504dc2f4ffc6e5e11` (`2026-06-29T13:24:43Z`) | 1 / 5 | 2 | `b0247bf042e2e239b33889e92ee38117147f3cbd` (`2026-06-28T09:36:34Z`; `refs/heads/main`) | `8cb61651107ff4fa23ad283504dc2f4ffc6e5e11` (`2026-06-29T13:24:43Z`; `refs/heads/main`) | `repositories/SharedUI.md` |
| `Hans-Einar/TerrainAnalyzer` | private | `main` | `77d71dc9bd248986409103e0e1fa1c37a3dfa21f` (`2026-08-19T17:53:15Z`) | 25 / 281 | 279 | `564fdca00e4f01e2984dc7023a8d4dd1d8ee4fa3` (`2026-06-21T19:12:05Z`; `refs/heads/agent/issue-12-parallelization-preparation`) | `7f3961c3527efa5ff08ec058e52b79447872e39f` (`2026-08-20T20:31:31Z`; `refs/heads/agent/issue-14-mesh-modeui`) | `repositories/TerrainAnalyzer.md` |
| `Hans-Einar/tplink` | private | `main` | `3da0f5ab8c2850e659b478718b6567023f8b8a7d` (`2026-08-10T18:19:58Z`) | 4 / 3 | 3 | `3da0f5ab8c2850e659b478718b6567023f8b8a7d` (`2026-08-10T18:19:58Z`; `refs/heads/agent/document-lte-tools`) | `63898b42a40db4d9aded384c5d0d2f3d046cbddd` (`2026-08-15T22:31:24Z`; `refs/pull/1/merge`) | `repositories/tplink.md` |
| `Hans-Einar/weight_app_flutter` | private | `master` | `538577830a2dbc13f603a2986f4a00f5afcc66a1` (`2026-07-14T07:30:41Z`) | 6 / 36 | 35 | `4e687ff4b54b62fc071ee26a8bca980053a57087` (`2026-04-25T20:21:44Z`; `refs/heads/master`) | `15455be0ed03f052790807ecd541e8d32a23a6ac` (`2026-07-14T19:08:09Z`; `refs/pull/2/merge`) | `repositories/weight_app_flutter.md` |

## Excluded repositories

**OBSERVED:** every excluded repository had zero reachable commits with a
committer timestamp in the fixed window. `Nearest earlier evidence` is the most
recent reachable commit before the start boundary and a captured ref containing
it; there were no excluded repositories whose nearest activity was only after
the end boundary.

| Repository | Visibility / flags | Default branch | Exact default-branch study commit (committer UTC) | Refs/commits | Window commits | Nearest earlier reachable commit (UTC; containing ref) | Evidence-backed exclusion |
|---|---|---|---|---:|---:|---|---|
| `Hans-Einar/7segmentDisplay` | private | `master` | `e6a804924a7afd998459f486bf2714d6b6d6bdfd` (`2023-05-14T10:41:50Z`) | 1 / 3 | 0 | `e6a804924a7afd998459f486bf2714d6b6d6bdfd` (`2023-05-14T10:41:50Z`; `refs/heads/master`) | No qualifying reachable commit. |
| `Hans-Einar/acanfd-stm32` | public; fork | `main` | `4c52989fd25402f06cf9aae08d89456e0feb712e` (`2025-03-16T18:04:12Z`) | 1 / 9 | 0 | `4c52989fd25402f06cf9aae08d89456e0feb712e` (`2025-03-16T18:04:12Z`; `refs/heads/main`) | No qualifying reachable commit. |
| `Hans-Einar/agro-crm` | private | `main` | `5cd3342264866d6956b4f81e9646715a0f180f94` (`2026-04-06T17:18:46Z`) | 2 / 51 | 0 | `5cd3342264866d6956b4f81e9646715a0f180f94` (`2026-04-06T17:18:46Z`; `refs/heads/main`) | No qualifying reachable commit; nearest activity predates the start by 14 days. |
| `Hans-Einar/CAD` | private | `master` | `4782bea5d0069edaf3d96b90771b20159a3d38ee` (`2025-02-05T11:59:33Z`) | 1 / 1 | 0 | `4782bea5d0069edaf3d96b90771b20159a3d38ee` (`2025-02-05T11:59:33Z`; `refs/heads/master`) | No qualifying reachable commit. |
| `Hans-Einar/CodexRemote` | public | `main` | `c3329fdc550e515166f805cfbb33425bbad73cb2` (`2026-04-03T09:38:54Z`) | 1 / 2 | 0 | `c3329fdc550e515166f805cfbb33425bbad73cb2` (`2026-04-03T09:38:54Z`; `refs/heads/main`) | No qualifying reachable commit. |
| `Hans-Einar/farmStatistics` | private | `main` | `bcba2e3d11e97a56a97739d809490ed8673b560d` (`2026-04-16T20:36:09Z`) | 1 / 4 | 0 | `bcba2e3d11e97a56a97739d809490ed8673b560d` (`2026-04-16T20:36:09Z`; `refs/heads/main`) | No qualifying reachable commit; nearest activity predates the start by four days. |
| `Hans-Einar/flutter_compass` | public; fork | `master` | `085b4c8cb42e485c1c26a38caaffa334fc50362a` (`2023-09-11T08:15:20Z`) | 2 / 126 | 0 | `59dc1771a0c0aec97d0a2a040ab3b8c84d550d42` (`2024-03-14T20:41:55Z`; `refs/heads/patch-1`) | No qualifying reachable commit. |
| `Hans-Einar/FT4232_CAN` | private | `main` | `2e82de4e0d96d32599fbf14edad33f39516d2536` (`2025-05-23T07:24:50Z`) | 1 / 4 | 0 | `2e82de4e0d96d32599fbf14edad33f39516d2536` (`2025-05-23T07:24:50Z`; `refs/heads/main`) | No qualifying reachable commit. |
| `Hans-Einar/html` | private | `master` | `451da27126a415fa9fd0b53f17731b975b3f16dd` (`2025-10-02T11:24:13Z`) | 3 / 14 | 0 | `451da27126a415fa9fd0b53f17731b975b3f16dd` (`2025-10-02T11:24:13Z`; `refs/heads/master`) | No qualifying reachable commit. |
| `Hans-Einar/hunter_observations` | private | `master` | `d2c6230f78f1a4861c2c2d65e2de06ba78c0a07e` (`2023-11-01T19:16:32Z`) | 1 / 2 | 0 | `d2c6230f78f1a4861c2c2d65e2de06ba78c0a07e` (`2023-11-01T19:16:32Z`; `refs/heads/master`) | No qualifying reachable commit. |
| `Hans-Einar/I2c_test` | private | `master` | `fe81c3d464d4ecf79fba0596c45ce70c21d04b68` (`2024-01-30T14:32:37Z`) | 3 / 27 | 0 | `0edcba3aec060be101b0bf263550daed6b1cafea` (`2024-06-14T08:51:13Z`; `refs/heads/version-2.0`) | No qualifying reachable commit. |
| `Hans-Einar/kicad` | private | `master` | `21803043ab31c4f6ddf2b0a9ed381a02aa54afd5` (`2026-02-05T11:22:12Z`) | 1 / 14 | 0 | `21803043ab31c4f6ddf2b0a9ed381a02aa54afd5` (`2026-02-05T11:22:12Z`; `refs/heads/master`) | No qualifying reachable commit. |
| `Hans-Einar/kinsim` | private | `main` | `20791c835a2882fde7b1acaf051a2041a75f6941` (`2026-03-13T19:26:56Z`) | 2 / 4 | 0 | `94d4c6ed75e8c84651b7278b665c1d973ddcaebe` (`2026-03-31T16:04:15Z`; `refs/heads/Kinsim-V2`) | No qualifying reachable commit. |
| `Hans-Einar/LedgerSandbox` | private | `main` | `85f0d56f23bb8465e40846d2491b6e64064694bb` (`2026-04-08T11:33:36Z`) | 1 / 1 | 0 | `85f0d56f23bb8465e40846d2491b6e64064694bb` (`2026-04-08T11:33:36Z`; `refs/heads/main`) | No qualifying reachable commit. |
| `Hans-Einar/mdb-runner` | private | `main` | `1745ac98ae20f4d2bb90581e0e0c555521c055d6` (`2024-10-01T19:10:52Z`) | 1 / 2 | 0 | `1745ac98ae20f4d2bb90581e0e0c555521c055d6` (`2024-10-01T19:10:52Z`; `refs/heads/main`) | No qualifying reachable commit. |
| `Hans-Einar/timberdata` | private | `main` | `36f641c9f14ceea597163cb7fc5e7725ca8844b4` (`2026-03-21T22:09:39Z`) | 1 / 10 | 0 | `36f641c9f14ceea597163cb7fc5e7725ca8844b4` (`2026-03-21T22:09:39Z`; `refs/heads/main`) | No qualifying reachable commit. |
| `Hans-Einar/tsm_locations_core` | private | `main` | `3e54b2f5ae61ee765f71dfaddbf970539adb8cb5` (`2025-12-26T17:32:40Z`) | 1 / 2 | 0 | `3e54b2f5ae61ee765f71dfaddbf970539adb8cb5` (`2025-12-26T17:32:40Z`; `refs/heads/main`) | No qualifying reachable commit. |
| `Hans-Einar/wakelock` | public; fork | `main` | `42389246ff277ec02f26c79a31b06445aaa245a4` (`2024-01-09T23:51:20Z`) | 57 / 214 | 0 | `c42b4217da8bc811d36991f137a669c4ef1361b8` (`2024-03-14T20:47:21Z`; `refs/heads/patch-1`) | No qualifying reachable commit. |

## Boundary and mutable-ref notes

- **OBSERVED:** no reachable commit in any repository had a committer timestamp
  exactly equal to either window boundary; the predicate nevertheless includes
  both boundaries.
- **OBSERVED:** `HSX` has a reachable commit at `2026-08-20T21:10:19Z`,
  `Lyndata` at `2026-08-20T21:40:15Z`, `TerrainAnalyzer` at
  `2026-08-20T21:49:50Z`, and `SDP` at `2026-08-20T22:05:48Z`. These commits are
  after the end boundary and were not counted. Each repository remains in scope
  because it also has qualifying commits.
- **OBSERVED:** default-branch study commits for `HSX` and `map_tracker` predate
  the window even though qualifying commits exist on other advertised refs. This
  is why default-branch history alone would not satisfy the contract.
- **OBSERVED:** several chronologically last activity commits are GitHub
  pull-request merge refs. They are counted because the fresh mirror captured
  those refs and the commits are reachable from them; their presence is not
  inferred from pull-request timestamps.
- **LIMITATION / OBSERVED:** advertised refs are mutable and GitHub may later
  delete branch or pull-request refs. The tables preserve the capture-time ref,
  full commit identity, committer timestamp, counts, and classification, but a
  later live re-run can differ if the owner-visible repository or ref set changes.
  No studied repository was modified and the temporary mirrors are not study
  deliverables.

## Verification checks

The inventory is ready for an independent verification pass when that pass:

1. authenticates as an account able to see the same owner repository set;
2. obtains 35 repositories with the same visibility/fork/archive breakdown;
3. fetches all advertised refs and enumerates all reachable commits before date
   filtering;
4. reproduces 17 included and 18 excluded repositories with the exact inclusive
   UTC predicate;
5. confirms each recorded default-branch study commit and evidence ref/commit
   relationship; and
6. confirms the 17 unique planned report paths exactly match the in-scope set.
