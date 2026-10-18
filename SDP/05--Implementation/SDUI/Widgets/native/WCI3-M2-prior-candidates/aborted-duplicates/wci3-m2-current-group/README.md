# Aborted duplicate coordinator invocation

A crossed asynchronous handoff started the same current SDL regression group twice.
The coordinator stopped only its own process tree (go test PID 210212) after
identifying the earlier worker run (209945). This diagnostic invocation is not a
failed product test or acceptance evidence. Its captured exit is retained; the
worker run in /tmp/wci3-m2-current-group-regression is the selected proof.
