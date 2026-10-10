# SDPTOOL-VER-RSP2 — release evidence

## Preparation and corrections

Release scope is MAINT-SDP-0015, authorized by owner Session0010 T011. Source
starts at RSP1 handoff 0652e174e4fbd126ca713c08ed5ec2465fa5e5f6. SDP 2.2.0 is
additive; Framework remains 2.0.0 and only Linux amd64 is published. gh-sdp 0.2.2
selects the new default without changing the thin client API. No XFMD source or
live project installation changes are selected.

Fresh independent review of RSP1 found terminal input could stop a background
runner and SIGTERM bypassed CLI cancellation. RSP2 transfers terminal foreground
ownership for interactive input, restores it afterwards and routes SIGTERM through
context cancellation. Real-executable Linux PTY and descendant-cleanup regressions
pass, alongside the original cancellation test. The reviewer reruns those checks
independently; the signed-package gate remains separate.

Repository validation, management validation and all 85 Toolkit tests pass after
aligning the maintained installation examples with 2.2.0. Original released notes
and descriptor digests are preserved. The payload inventory is unchanged; only
the managed SDP entry skill differs from 2.1.0 (previously accepted Session upkeep).

Exact clean package, production signature, predecessor trials, archive gate,
CI and final independent approval will be recorded below when observed.
