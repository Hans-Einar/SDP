# MAINT-SDP-0003 implementation evidence

## IU1-M2

Candidate: IU1-M1 9a5c79a plus the profile build/template/schema diff.
Python 3 on Linux. build_process_profile.py emits a reproducible self-contained
artifact; two test groups cover byte reproducibility and ten negative inputs.
The 1.0 installed schema is unchanged; 2.0 has its own closed schema.
New neutral templates are excluded from the legacy v1 inventory explicitly.

The full Python suite exposed a stale expected scenario list: SK1 had added two
conformance scenarios without updating that set. Added their existing IDs;
no expectation about installation behavior was weakened. The other failure is
the 38 pre-existing Traceability compatibility errors recorded under SK1/KB011.

## IU2-M1

Candidate: IU1 89981fd plus versioned PowerShell planning diff.
Five integration test groups pass with PowerShell 7.6.6: clean deterministic
read-only plan, legacy phase/archive link mapping, manual XFMD-shaped history
transfer, mixed-layout/link conflicts and corrupt artifacts. Shared installed
facts parsing remains closed and explicit; legacy defaults stay schema 1.0.
Fixture provenance is documented in Process-Installation.md. No live consumer
files were changed. Apply/resume remains IU2-M2 work.

The owner added typed plans and a Planning skill during execution. Registered
KB-SDP-029 without altering this Maintenance's already selected Git policy.

## IU2-M2

Delivered journaled apply/forward resume with exact-plan comparison, exclusive
writer lock, byte backups, post-failure drift checks and retry-safe report/facts/
history finalization. Thirteen integration groups passed on a frozen candidate
(307.662s, PowerShell 7.6.6/Linux); the final reserved-path and balanced/escaped
Markdown refinements also passed focused tests. Full legacy PowerShell fixture
suite passed. The exhaustive forced-exit matrix remains an IU3 completion gate,
not a completed check here.

Independent review by iu_review required four corrections: history schema/chain
validation, broader Markdown rebasing, FIFO rejection and reserved artifact
paths. All were corrected; the reviewer approved the IU2 scope after five
focused tests. Final reviewed hashes:

- Process-Install.ps1: 3af5f6e066fa07ea4ba59bc4d7245c3c7761b0e66f87aae1f51f3ed78117e2cb
- Install-SDP.ps1: 2f79d1a1d34af856fbe64feedaf944daefedeba8938dcf0e59ef20e24122021b
- build_process_profile.py: e7a35d15ae6ee7d3851758c36a71c6f84042b9fe28e88b9cd57970a994fe680f
- artifact: 47a5363f77408a9a99803cd370ba7bbe7f2b3d8e3bdce377f9a5aa1d4890ab11

The distributed management payload schema generalizes the project namespace;
this repository keeps its stricter local SDP namespace validator. Installation
allocates target-project Maintenance/event IDs from existing history/documents.
Backlog review: KB014 still owns standalone KanBan packaging, KB018 the wider
Toolkit audit, KB029 the new plan/Planning-skill proposal. None is silently
marked implemented by this installation milestone.
