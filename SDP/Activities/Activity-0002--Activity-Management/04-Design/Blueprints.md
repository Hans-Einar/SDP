# Pinned Activity implementation handoff

TARGET: b3e45cab-c629-490b-9192-3ed2b4628387:00002; canonical parser revision
35c1b76910eac0af13b767b344ef0873cc031ada67f0fbdc69b165b547984874.
NOW: 73a3ef6f-6869-4b83-a6d8-993d2a8c4bda:00001, modeling the absence of Activity
services in this concern, not the absence of the wider SDPTool system.

[Retained bundle](../../../Blueprints/f75501b41f7bafdfbbeca5c5ba3e10f23ab44c0b68681e5fafb1673255b6187d/c805b96624525af1c308ddb7ff1beee94896f7a09328811f55b6f64b84c5e4d5/index.md).
Blueprint revision: 830b6836603b55b150e04d6f9018f636ce8332c144e7b7d01ee91349725b25ec.
Retained revision: c805b96624525af1c308ddb7ff1beee94896f7a09328811f55b6f64b84c5e4d5.
Structural constraints pass. Code conformance/runtime dependencies remain unknown;
SDUI semantics is unsupported by this analyzer. These are honest limits, not gates
waived by a green parse. WORKs are preliminary, not accepted model releases.

[ImplementationContract.md](ImplementationContract.md) supplies transaction,
compatibility, caller, recovery and consumer semantics. PLAN-SDP-0027 owns acceptance.
Task/code paths scope the SDP producer implementation; XFMD consumer changes require
their own checked worktree and source mapping, not a claim of SDP Blueprint coverage.

The implementation refinement adds Template and SDP/ProjectManagement to allowed
code paths because installed schemas and read-only validators are maintained there.
Regeneration retains the original A0 bundle unchanged (retained revision9bfcb52f…)
and uses the same NOW/TARGET checkpoints. This is a scope-source correction within
the selected installer/readers requirement, not a changed product need.
