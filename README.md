# SDP Draft Project

Status: working draft  
Created: 2026-06-23

This folder is the central draft home for the Standard Document Procedure
used across local projects under `C:\Users\hanse\GIT`.

It was created because no top-level `C:\Users\hanse\GIT\SDP` project existed
yet, while several repositories already contain repo-local SDP folders.

## Source Repositories Inspected

Existing SDP folders observed:

- `C:\Users\hanse\GIT\farmStatistics\SDP`
- `C:\Users\hanse\GIT\tsm_locations_core\SDP`
- `C:\Users\hanse\GIT\weight_app_flutter\SDP`
- `C:\Users\hanse\GIT\SharedUI\SDP`
- `C:\Users\hanse\GIT\HEOS\SDP`
- `C:\Users\hanse\GIT\TerrainAnalyzer\SDP`
- additional older/lightweight SDP folders in other local repositories

The newest and richest local operating instructions currently appear in
`farmStatistics\SDP`, with useful related process notes in
`tsm_locations_core\SDP`, `SharedUI\SDP`, and `HEOS\SDP`.

## Documents

- `DraftStandardDocumentProcedure.md`: the current cross-project SDP draft.
- `TieredDesignAndImplementation.md`: the Tier concept added from the
  TerrainAnalyzer design-analysis work.

## Current Principle

Design horizontally. Implement vertically.

Horizontal design documents describe architectural layers and contracts.
Vertical implementation Tiers cut through those layers to deliver coherent,
verified capabilities sprint by sprint.
