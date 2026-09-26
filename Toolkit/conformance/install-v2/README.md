# Process installation conformance 2.0

Use Python 3, the Toolkit test requirements and PowerShell 7.4+ on Linux.
Set SDP_TEST_PWSH to the executable when it is not on PATH. Build SDPTool once
with its package.sh and set SDP_TEST_TOOL to that prebuilt executable.

~~~sh
python3 -m unittest discover -s Toolkit/tests -p test_process_profile.py -v
python3 -m unittest discover -s Toolkit/tests -p test_process_install.py -v
python3 Toolkit/conformance/install-v2/run_fault_matrix.py
~~~

The unit/integration fixtures test deterministic/no-write plans, legacy and
manual layouts, schema/chain rejection, links, FIFO/case conflicts, exact-plan
drift, forced managed backup, concurrent writers, post-failure edits, repeat
and installed consumer behavior. A missing PowerShell skips explicitly; it
does not pass installation evidence. The consumer test likewise needs the
prebuilt executable. CI supplies both and keeps the existing v1 suite separate.

run_fault_matrix.py divides one deterministic operation's step indices between
four independent temporary projects. Each selected step is interrupted after
backup, after atomic write/delete and after journal publication. It uses actual
process exit 97, including preparation and finalization; each fixture resumes
to completion and proves preserved history, unique events and no-change repeat.
This tests every operation step, not just one representative write. It does not
simulate hardware failure or claim whole-tree rollback.

verify_xfmd_snapshot.py accepts an explicit source checkout and a provenance
output path. It copies only the process area and incoming Markdown into a
temporary project, records source hashes/commit, installs there, verifies the
history prefix and checks that the original bytes/Git status did not change.
It must never upgrade the source checkout. Its fixture is evidence for the
observed baseline, not permission to roll out to a live consumer.

The authored configuration build and existing v1 conformance remain independent:
do not rewrite legacy expected plans to make v2 behavior pass. Windows v2
execution is experimental until separately verified.
