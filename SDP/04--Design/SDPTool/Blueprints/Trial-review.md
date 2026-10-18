# BP2-C independent review — 2026-10-07

Reviewer: independent agent /root/bp2c_review, read-only fresh context.
The reviewer read Trial-Assignment.md before assessing the temporary candidates.
The coordinator retains this summary of the actual returned review; it is not
owner acceptance or generated-bundle evidence.

## Disposition

Candidate A: approved for the bounded model-edit assignment. Only
Containers/MachineService.design changes. CalibrationProcessor owns
CalibrateMeasurement and is contained by MachineService. Inspection allocation,
reduction/publication ownership and AuthoritativeMachineState remain preserved.
The five other files are byte-identical to NOW.

Candidate B: changes required. High-severity violation at line 4:
BuckingWeb owns ReduceMachineState replaces MachineService ownership, violating
the explicit PRESERVE obligation. The calibration extraction itself is compliant.
This is a deliberate negative control, not an unresolved production defect.

Both pass the real SDL parser. Checked revisions:

- A: 29a5bfd5bacbb4ccb78054a62637a245572e2a8a043bac24c8e62080eb1f20b1
- B: e2b9a3ebfc1e74c466482456cd97b12591aec5435eb0c50463b272e62441edd6

The reviewer independently confirmed NOW matches d7aa3e5 and the binary SHA-256
matches retained 043c59c pilot evidence. Trial-evidence.json and
experiments/blueprint_mvp1/handoff_trial.py retain the reproducible specimens.

Producer-and-Handoff.md and PLAN-SDP-0020: approved as proposed design documents,
with no material boundary or integrity contradiction found.

## Limits

This establishes structural compliance and deliberate violation detection only.
It proves no runtime mathematics, sequencing, protobuf compatibility, code
conformance, production generator behavior or future autonomous compliance.
Publisher recovery remains an implementation verification obligation.
BP2-A owner disposition remains open. Production generated-bundle trials are BPI3.
