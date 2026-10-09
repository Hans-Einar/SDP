# Owner development trials — PGI3-U1

Owner direction, Session0007 T013 (2026-10-09): establish shared understanding through
use before extending the implementation. The owner selected a dedicated PG lab and
wants both project/work overview and interaction during an actual task. This local
practice is part of PLAN-SDP-0018; it is not a new installed global SDP procedure.

## Why this changes the next step

C1–C6 establish bounded technical behavior. They do not establish that the emerging
workflow matches the owner's needs. Documents, model claims and passing agent-run
tests cannot replace this feedback. The immediate next step becomes owner use and
reconciliation of expectations; native mapping and PGI4/PGI5 remain unfinished.
Do not interpret preparation of the lab as owner acceptance of the candidate.

## Repeatable loop

1. Select a short user operation and the uncertainty to learn from it. Write what
   the owner should be able to understand and do, including realistic limitations.
2. Freeze an engine source commit and build private binaries. Record commit,
   source/binary identities, host/toolchain and adapter identity. Keep the stable
   installation untouched. Select a fresh trial project or a bounded existing
   consumer assignment; do not attach an agent to an unselected active task.
3. Prepare valid current work and separate operational state. Expose goal, Session,
   card/plan, authority, allowed changes and current state before execution. Record
   whether the project is a minimal pilot or a full installation.
4. The owner starts and uses the workflow, steers and responds to actual decisions.
   Retain visible events, exact forwarded input, diff and test results. Label any
   fixture/replay checks separately; never fabricate an owner experience result.
5. Record expected versus experienced behavior, confusion and desired changes.
   Separate implementation defect, workflow mismatch, newly clarified requirement
   and known limitation. A console presentation issue and a core workflow issue
   can have different fixes.
6. Review observations together and update the governing design/plan/card/Session.
   Agree the next bounded change from this evidence, then create a fresh candidate
   and repeat. Broader proposals go to their owning cards. Successful tests alone
   do not authorize progressing past an unresolved owner workflow decision.

## Current prepared surface

The [lab instructions](../../../../SDPTool/labs/project-governance/README.md) and
[source preparation](../../../../SDPTool/labs/project-governance/prepare.py) define
the runnable surface. The local lab home is `/home/warloc/git/pg_lab`. It has a
readable terminal overview, an explicit start, steering/interruption/scope controls,
per-request decisions, provenance, actual diff/tests and a feedback sheet.
The engine candidate is C6 commit 34b345e; the lab's manifest records the full hash
and separate harness identity. Preparation starts no Codex task.

The first operation is a small runnable Go label command with a Unicode boundary
trimming defect. Only its function and tests may change. This supplies a transparent
first experience, not proof of suitability for a large existing project. The next
consumer trial can use sdui_widget_lab after choosing a real bounded task with its
owner; its existing task, pinned library source, installed Framework and Git work
must remain intact. The current assignment does not modify that project.

## Current architecture versus future integration

The PG services already live in SDPTool Go packages. `sdptool PROJECT governance`
exposes read/routing/run operations, several with JSON request files and JSON output.
`sdp-controller` owns the persistent host conversation and control stream;
`sdp-mcp` exposes scoped agent operations through the same core. They are separate
adapters today, not three separate implementations awaiting consolidation.

KB-SDP-051 in the blueprint workstream owns the proposed common action catalogue
and JSON invoke adapter. Its inspected handoff is
`/tmp/sdp-blueprint-implementation/SDPTool/XFMD-Blueprint-Integration.md`.
Future discrete PG operations should align with that shared contract rather than
add a parallel generic dispatcher. The running conversation also requires streaming
observations, commands, cancellation and decisions; a single subprocess request/response
is not the whole interactive client. No XFMD integration is claimed by this lab.

## Owner trial questions and disposition

Can the owner locate current work, understand its boundaries, see who is acting,
steer safely, understand a decision request, stop work and assess what remains?
The lab's FEEDBACK.md retains the answers with candidate/run identity. At preparation,
these questions are **not yet tested by the owner**. Further feature expansion waits
for this first feedback discussion; technical inspection/fixes needed to make the
trial usable can proceed. This is a sequencing correction, not a claim that PGI3
or the Session is complete.
