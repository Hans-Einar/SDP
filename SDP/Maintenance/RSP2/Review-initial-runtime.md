**CHANGES_REQUIRED** — runtime candidate `99631595eb2ccf5d1e8cc52be5ce30a02b492341`, compared with `04f88ff918e7c25c9fa1883f061ef20b8c3923a2`.

Two process-lifecycle findings:

- **P2 — Terminal input can hang the runner.** Line 12 creates a separate process group without transferring terminal foreground ownership, while the launch code inherits terminal stdin. A runner reading stdin receives `SIGTTIN` and stops; SDPTool keeps waiting. A file-free PTY probe reproduced this mechanism. Implement terminal-aware foreground handling and add a PTY regression test.
  
file:///tmp/sdp-runnable-programs/SDPTool/program_process_unix.go

- **P2 — SIGTERM leaves the launched application running.** The CLI subscribes only to `os.Interrupt` at line 11. Terminating SDPTool with SIGTERM therefore bypasses context cancellation and the new process-group cleanup. Route Unix termination signals through cancellation and verify descendant cleanup through the actual executable.

file:///tmp/sdp-runnable-programs/SDPTool/cmd/sdptool/main.go

Inspected the requirements, design, plan, contract, complete runtime diff, surrounding discovery/path/CLI code, new tests, full race-suite log, and gh-sdp/native integration evidence—including screenshots and normal-close receipt. Read-only discovery, source preservation, literal argv/cwd, containment, revision checks, and normal exit propagation otherwise align with the bounded contract.

All 14 recorded source hashes matched the candidate; the saved executable’s hash matched its receipt. Independently rerunning that executable’s discovery returned five sources and the expected runnable widget-lab revision. Later checkout commits contain no runtime changes.

Focused Go tests could not rerun: the read-only sandbox prevents creating Go’s build directory. The PTY probe establishes the underlying mechanism, not an end-to-end candidate test.

No files or lifecycle records were changed. Final release manifest/signing/package review remains separate and pending.