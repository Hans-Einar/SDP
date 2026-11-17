# Implement DocumentHost prototype launch

| Field | Value |
| --- | --- |
| Session reference | SESSION-SDP-0026 |
| Status | planned |
| Activity | ACTIVITY-SDP-0007 |
| Phase | 04-Design |
| Phase contributions | 04-Design, 05-Implementation, 06-Verification |
| Plan | PLAN-SDP-0030 |

## Goal

Deliver ordinary prototype launch with explicit context and verified model

## Linked records

- plan: [PLAN-SDP-0030](../../../05--Implementation/SDUI/DocumentLaunch/Plan.md)

## T001 — Selected implementation and scoped preflight

Owner authorizes the proposed DocumentHost extension. Reopened Activity0007 for
implementation, retaining completed reproduction and Session0025. PLAN-SDP-0030
owns DL1/DL2; primary execution card KB-SDP-065. Preflight and worker write ownership
are recorded in the Plan (672fcb0 producer, clean 9b7af853 consumer). SDUI library
and existing DocumentHost contracts establish the implementation basis. Loaded SDP,
Architect, Planning, Worker and shared document/plan/model routines. A validated
WORK model defines the target before product changes. No separate blueprint is
needed for this bounded assignment; no publication or live replacement is selected.

| State | Step | Outcome |
| --- | --- | --- |
| on-going | DL1 | Explicit context, headless preflight and DocumentHost native route |
| planned | DL2 | Packaged native consumer checks, one review, model release and closeout |

One Sol/medium worker owns the prototype/native implementation; parent integrates
CLI/action/schema, tests, models and documents. Initial design adds --context and
--context-revision; watch observes source/context/resource identities. Headless
preflight must not claim actual native geometry. Next: finish DL1 and test the
same original TimberFrame model through the normal packaged launch path.

## DL1 source checkpoint

CLI/action context transport and truthful preflight are implemented. The bounded
worker delivered context parsing, resource snapshots and native DocumentHost/watch
composition. Targeted tests pass; one consolidated reviewer found no source
blockers. Exact packaged native qualification, final race evidence, model release
and closeout remain DL2 obligations. No product release is selected.
