# SDUI runtime in Go

G3-M1 implements a synchronous, host-owned UI session with typed handles/events and atomic property batches. Marshal background calls onto the owning UI goroutine first. Runtime stores no native pointers and does not execute SDL.

`New` accepts a validated instance tree. `Bind` explicitly registers Go handlers; `Dispatch` rejects stale/duplicate events and unbound controls. `Draft` edits locally without callbacks; `Commit` sends current draft to its handler. `Apply` validates complete batches before publication. External value updates conflict with dirty fields; explicit acceptance must match current draft. `Revert` provides Escape behavior. `SnapshotRoot` creates an independent presentation model; programmatic updates are not user events. `Close` revokes the session.

G3-M2/M3 add hash-based file watching in `../reload`, last-valid-model handling and native Fyne adaptation. Parse/profile validation and measurement precede candidate publication. `Reload` preserves value/draft/focus/handle for compatible named paths/types. New defaults affect new instances; new labels/enabled/visible apply immediately. Anonymous widgets receive new generations. Changed callback references do not reuse old handlers. Type changes/deletion revoke handles.

CLI watches by default (`-watch=false` disables it). Input edits update draft; Enter commits, Escape restores accepted values. Native callbacks capture their model revision. Old callbacks/results after teardown/reload cause no new action. The adapter mutes programmatic `SetText`, preventing property updates from becoming domain events. Session still contains no SDL loader. [Evidence](../evidence/G3.md).
