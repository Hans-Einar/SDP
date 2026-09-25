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
