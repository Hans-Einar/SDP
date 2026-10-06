# FAM1-M1 — SDL Go familiarization evidence

Observed on 2026-10-07 (Europe/Oslo). [Session 0009](../../session-%230009--SDL_Go_familiarization.md)
and [KB-SDL-008](../../../KanBan/completed/%23008--SDL--Study--Go-code-familiarization.md)
own this study. No SDL or SDUI code was changed.

## Candidate and result

HEAD: `1d3ad16a258a94a64603e2465886ba0b15177c02`. Tracked SDL/go and SDUI/go files match HEAD.
The pre-existing untracked sourceinput package is included and has no tests.
[Candidate inventory](candidate.json) records SHA-256 identities for both module
trees, including fixtures. Other worktree changes belong to other workstreams.

Go: `go version go1.27.1 linux/amd64`, through ~/.local/bin/go to the existing
~/.local/share/sdp-toolchains/go1.27.1/go/bin/go installation. From SDL/go:
`GOMAXPROCS=2 go test -race -count=1 -p 2 ./...`.
Exit **0**; all **14 packages with tests passed**. [Raw output](go-test.log).
The count is derived from the 14 `ok` lines; an earlier progress message incorrectly
said 15. These tests ran without cached results, with the race detector enabled.

## Limits and disposition

MMDR was unset. TestRealSemanticBackend and TestClassDocuments skip their external
renderer checks under that condition. Actual renderer behavior and desktop-tag
Fyne programs were not verified. The untracked sourceinput package compiled but
has no tests. This is neither exhaustive review nor release/UI acceptance.

The bounded code familiarization, PATH repair and test-baseline task is complete.
Further product work requires a concrete selected scope.
