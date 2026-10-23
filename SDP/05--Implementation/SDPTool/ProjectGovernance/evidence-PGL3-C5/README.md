# PGL3-C5 evidence

`candidate.json` identifies baseline plus tested source/test/schema hashes.
`checks.json` and logs record full affected SDPTool race/vet results. The runner
sets `SDP_PGL3_C5_OLD_READER` to a controller built from a84d71c using `git archive`
of SDPTool, SDL/go and SDUI/go into a private cache. It does not change a checkout.
`TestCompiledReportAcknowledgement` exercises distinct current CLI processes and
asserts that this actual old binary refuses the new store without mutation.

`wire.json` captures the actual request, receipt and inbox before/after from that
compiled fixture. The four objects pass JSON Schema Draft 2020-12 validation using
`schema.json` definitions ReportAcknowledgementRequest/Result and inbox.schema.json.
`schema-checks.json` records those results. Temporary fixture paths/IDs are retained
as observations. No native model task, independent review or owner acceptance.

To reproduce: build the old controller in a private baseline archive; run the
SDPTool suite with `SDP_PGL3_C5_OLD_READER` pointing to it and `SDP_PGL3_C5_WIRE`
pointing to a writable private output file. `checks.py` preserves actual original
workstation paths; adjust them for a different host. Validate emitted JSON against
the candidate schemas. Original evidence files are records, not output destinations.
