# R1 verification

## R1-M2 — templates and project area

- Addressed 106 file locations per migration map; deduplicated seven numbered root templates after comparing content without trailing newlines.
- Toolkit: 80 unittest cases run, 79 passed. RepositoryValidation fails only on the same nine old Issue #5 IDs as before migration; baseline/post-check logs are byte-identical. The full suite is not green.
- Install-v1 conformance: 17 scenarios validated. Expected source addresses moved to Template; destinations/policies preserved.
- PowerShell is unavailable locally; no native Windows tests ran here.
- Old bootstrap files archived without content changes. Historical ledger lines preserved, new R1 events appended. No existing installation updated.

Commands from repository root:

```sh
python3 Toolkit/scripts/validate_sdp.py
python3 -m unittest discover -s Toolkit/tests -p 'test_*.py'
python3 Toolkit/conformance/install-v1/run_conformance.py --validate-only
```

## R1-M3 — languages, documents and entry points

- Moved 216 additional files. verify_structure.py checks all 322 migration entries; duplicate templates consolidated, not copied again.
- SDL/SDUI: `go test -race ./...` passed with Go 1.27.1 and registered mmdr. Moved Go files and go.mod/go.sum are byte-unchanged.
- Hash-checked 574 manifest-listed SDUI/design outputs, unchanged. No viewpoints hand-edited or regenerated during migration.
- Historical source-index/release fingerprints byte-unchanged; old implementation-ledger lines preserved as a prefix. Archived bootstrap unchanged.
- Native launcher trial in isolated Xvfb passed: installed sdl-design symlink now targets SDL/scripts; main page/navigator open, SVG generates on demand, no Go build starts, missing-tool diagnostics and close-time cleanup work.
- Issue #5 usage-study validator passed: 35 assessed repositories, 17 reports.
- Local link checks cover file/directory targets, not anchors/external URLs; frozen fixtures and retired bootstrap explicitly excluded.
- KanBan chains, metadata and locations checked. Card 001 remains active with bounded remaining work; 010 (consolidation) and 011 (old IDs) are backlog. Migration created no language profile or subrepository.

Commands:

```sh
python3 SDP/Maintenance/R1/verify_structure.py
python3 SDP/Studies/UsageAnalysis/validate_analysis.py
python3 SDP/Agents/KanBan/examples/verify_lineage.py
go -C SDL/go test -race ./...
go -C SDUI/go test -race ./...
```

Go and native renderer were selected through explicit local paths. Post-M3 Toolkit status is compared with baseline; the nine old ID mismatches must be the entire remaining failure set. CI runs Windows installer tests; there is no local Windows verification.

Final checks passed: 322 file locations, 574 generated outputs and 1887 local Markdown targets. KanBan: 15 cards, 21 events, 130 local links. Toolkit's CI validator mode with `--base-ref origin/main` returned exactly the same nine baseline mismatches. Raw race logs: sdl-race.txt and sdui-race.txt.
