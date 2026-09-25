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
