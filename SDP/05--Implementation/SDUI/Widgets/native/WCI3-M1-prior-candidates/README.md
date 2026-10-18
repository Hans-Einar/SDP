# Superseded WCI3-M1 candidates

These retained records are not final acceptance. Pilot 78344a92 failed an actual
numeric step-button click, corrected before final candidate. Candidate 0e92 failed
slider observer recapture and next-key gesture after focus loss. Initial reentry
run falsely passed on an intermediate snapshot; the serialized owner-state barrier
reproduces the failure. The first focus run queried focus too early; its corrected
barrier run proves the missing next Commit. Neither harness artifact is a product
pass. Final 3c1abba7 runs use the hardened harness and pass both regressions.

Pilot history and candidate binary identities are in ../Pilot-WCI3-M1.md. Raw event
bytes are losslessly compressed. Screenshots are actual X11 captures. No caches,
configuration or executables are included.
