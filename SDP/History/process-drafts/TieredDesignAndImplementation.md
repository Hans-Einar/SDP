# Tiered Design And Implementation

Status: working draft  
Origin: TerrainAnalyzer design-analysis work, 2026-06-23

## Purpose

This document defines the Tier concept for SDP-driven work.

Tiers help bridge horizontal design and vertical implementation. They make it
possible to design a system by layers while still implementing end-to-end
capabilities sprint by sprint.

## Definitions

Layer:

- a horizontal design boundary
- owns responsibilities, classes, data records, contracts, and verification
  expectations
- examples: frontend shell, backend API, persistence, workers, artifact store

Contract:

- a communication interface between layers
- may be an API, DTO, schema boundary, artifact manifest, event contract, file
  format, or verification evidence format

Tier:

- a vertical implementation capability
- cuts through several layers
- should produce a small coherent working result
- should normally be implemented by one sprint or by several slices inside one
  sprint

## Why Tiers

Without Tiers, teams often implement one horizontal layer at a time and end up
with broad unfinished scaffolding. Tiers keep delivery grounded in a working
capability.

Examples:

- not just "create database tables"
- instead "create project from staged coverage and reload it"

- not just "install viewer library"
- instead "load a self-hosted terrain fixture in 2D and 3D with coordinate
  picking"

## Tier Contract Template

Each Tier should state:

- Tier ID and name
- goal
- user-visible or operator-visible capability
- requirements implemented
- layer designs touched
- communication contracts used
- expected files/modules
- invariants
- non-goals
- verification evidence
- review expectations
- completion signal

## Recommended Tier Fan-Out

The exact Tier names differ by project, but a common sequence is:

1. runtime skeleton and verification harness
2. user/project/domain spine
3. job or workflow spine
4. first real visualization or primary user workflow
5. real backend processing
6. cache/offline or reliability improvements
7. external integration
8. analysis/reporting workflows
9. advanced domain extensions

Do not start with every layer's full implementation. Start with the smallest
vertical path that can be verified honestly.

## How Tiers Relate To SDP Folders

`05--DesignAnalysis`:

- identifies layers
- identifies contracts
- proposes Tier fan-out

`06--Design`:

- creates one detailed design per layer
- explains how each layer participates in Tiers

`07--Implementation`:

- turns Tiers into implementation strategy
- defines expected sprint order

`Sprints`:

- execute one Tier or a bounded part of a Tier
- contain slice contracts, implementation notes, review state, verification,
  and handoff

`Traceability`:

- maps requirements to layers, contracts, Tiers, slices, reviews, and
  verification evidence

## Agent Delegation Pattern

During design:

- spawn one sub-agent per layer document
- keep shared traceability integration with the Master agent
- review layer documents independently

During implementation:

- spawn worker agents by vertical slice, not by isolated layer
- each worker should own a bounded set of files/modules
- reviewers should be separate fresh agents
- the Master agent integrates outcomes into SDP traceability and handoff
