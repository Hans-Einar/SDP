# MVP1 blueprint input pilot — BP2-A

This authored specimen translates a small part of the experimental MVP1 model into
supported design-core/0.6. It does not replace or port that model and authorizes no
Ponsse code change. Each snapshot contains a complete six-file composed model.

NOW is a model-only baseline. TARGET extracts CalibrateMeasurement into a new
CalibrationProcessor Unit within MachineService. Reduction/publication ownership,
allocation, neighboring container ownership and the selected contract stay unchanged.

## Provenance and limits

Source material under experiments/mvp1_sdl/SDL/MVP1:

- Containers/MachineService/Domain.design: calibration, authoritative reduction and publication.
- Scenarios/InspectMachine.design: inspection flow through MachineService, BuckingWeb and BuckingUI.
- Contracts/Machine.design: independent subscriptions, evidence and no simulator hidden truth.

The original sources use sdl-mvp1-exercise/0.1. The pilot's BaselineRequest,
BaselinePayload, CalibrationState and request/reply projection are authored
simplifications, not an exact protobuf schema or a claim of wire compatibility.
Target transport, simulator, command permissions, P1000 parsing, sequencing,
calibration math, browser channel and UI behavior are omitted and remain unknown
for impact completeness. No SDUI change is proposed by this extraction.

Use the detailed contract in SDP/04--Design/SDPTool/Blueprints/Contract.md.

## Reproduction

Build `sdptool` and `sdl` from commit 043c59c, then run from repository root:

```sh
python3 experiments/blueprint_mvp1/probe.py --sdptool /path/to/sdptool --sdl /path/to/sdl --output /tmp/bp2-fresh-run
```

The output directory must not exist. The probe validates both snapshots, creates a
model-only release, makes a dirty WORK from it, captures a preliminary source view,
freezes a candidate and compares source digests. It invokes SDL viewpoint generation
for both models. These are ordinary source-generated views, not a blueprint diff or
an annotation renderer. The retained evidence records binaries, source hashes and
actual artifact identities; repeat runs get new UUIDs/timestamps.

## BP2-B selection experiment

```sh
BP2_SDL=/path/to/sdl PYTHONDONTWRITEBYTECODE=1 python3 -m unittest discover -s experiments/blueprint_mvp1 -p 'test_selection_probe.py' -v
python3 experiments/blueprint_mvp1/selection_probe.py --sdl /path/to/sdl --output /tmp/selection.json
```

The Python probe consumes real toolkit AST JSON; it is design evidence only.
Production belongs in Go. Selection-and-Evidence.md defines the proposed complete
contract and explicitly distinguishes it from the experiment's smaller coverage.
