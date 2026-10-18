# SDPTool runnable SDUI programs

## SDPTOOL-REQ-RSP1 — discover and execute project applications

Owner Session0010 T010 assigns program discovery and launch to SDPTool. A consumer
must be able to run gh sdp . discover in sdui_widget_lab and receive runnable SDUI
programs alongside the existing SDUI sources. Selecting a program starts the same
connected application as its project start command, including SDL routines, Go
functions and providers. XFMD consumes this capability rather than owning a
project-specific registry or inferring commands from a source file.

Acceptance: retain source identities/results; describe explicit programs and
blocked/invalid reasons; discovery never executes code; explicit ID selection
starts the declared argv in the project root and exposes failures; stale selected
revisions reject before launch. The widget lab proves the actual native path.
A published engine/default gh-sdp release and XFMD UI changes are separate scopes.

Design: SDPTOOL-DES-RSP1, ../../SDPTool/Programs.md.
Execution: PLAN-SDP-0023, KB-SDP-051, Session0010 S7.
