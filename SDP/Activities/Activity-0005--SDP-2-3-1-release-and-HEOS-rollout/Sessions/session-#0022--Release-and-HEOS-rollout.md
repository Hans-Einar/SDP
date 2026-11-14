# Release and HEOS rollout

| Field | Value |
| --- | --- |
| Session reference | SESSION-SDP-0022 |
| Status | planned |
| Activity | ACTIVITY-SDP-0005 |
| Phase | 05-Implementation |
| Phase contributions | 05-Implementation, 06-Verification |
| Plan | MAINT-SDP-0020 |

## Goal

Qualify and publish SDP2.3.1 and upgrade HEOS through the signed installer

## Linked records

- plan: [MAINT-SDP-0020](../../../Maintenance/R231/Plan.md)

## Release preflight and route

Owner selects a new installable release and HEOS first. Main cef86e4 contains
completed Activity1–4. New Activity5 owns this bounded rollout. Loaded SDP,
Release, Versioning, Verifier, Planning and Traceability routines. Existing HEOS
dirty browsing documents/history are preserved; all rehearsal uses a copied tree.
MAINT-SDP-0020 owns the milestones. M1 qualification is current; M2 public delivery
and live HEOS verification follow. No other project or XFMD package is selected.
Version2.3.1 is an additive predecessor-support/documentation patch with unchanged
runtime semantics; no language/wire version bump. See Maintenance/R231/Plan.md.
