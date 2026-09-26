# SDL — executable action profile

G4-M2, 2026-09-22. `action-core 0.1` is an explicit, bounded execution profile alongside structural `design-core 0.5`. Structural Functionality, Mode and Channel facts receive no implicit execution semantics. Profiles share the SDL lexer, identifier rules and source positions; each has a closed grammar.

```text
language action-core version 0.1.
action Echo.
record EchoInput.
record EchoOutput.
Echo invokes GoEcho.
Echo returns EchoOutput.
Echo takes EchoInput.
EchoInput field Value as text.
EchoOutput field Value as text.
```

An Action requires exactly one `takes`, `returns` and `invokes`. The first two reference declared Records; `invokes` names an explicitly registered Go function. It is a symbolic reference, never Go source, shell or dynamic code import. A function may be reused when both record signatures match. Each Record has 1–32 locally named fields. Fields are required and typed `text`, `integer` (signed 64-bit) or `boolean`. Reject extra/missing fields and wrong scalar types. No null/default or implicit conversion. Text must be valid UTF-8 without NUL, at most 32 KiB.

Sort declarations by name, then all facts by complete ASCII sentence. Canonical source contains no comments; documentation surrounds the model. Initial limits: 128 Actions, 128 Records, 32 fields per Record. Source/token limits follow the shared SDL frontend. Unknown statements are rejected.

The Go registry provides input/output signatures and a `func(context.Context, Record) (Record, error)` function. Validate the whole model and required registrations before creating runtime. Validate input before invoking handlers and results before publication. Copy record values at boundaries. Errors are observable outcomes; runtime does not infer domain rules.

Registered Go implementations own domain state. Keep UI draft, UI revision, command identity and accepted domain result separate. This is not a distributed messaging system or a guarantee layer for external side effects. G4-M3 includes a bounded EditAptCell trial: stable cell identity, draft, expected domain revision, explicit rejection and correlated accepted result. Label it a simulation; it performs no Ponsse/machine actions.

Model hot reload and in-flight call handling are implemented/tested in [G4-M4](../../go/evidence/G4.md#g4-m4). Changing Go functions requires a normal Go build and process restart.

## SDUI bridge — G4-M3

SDUI refs are symbolic. The host supplies an explicit alias→SDL-runtime table; neither bridge nor parser opens ref paths. Executable callback: `module.Action.@invoke`. `module.Action.setHandle(page.input)` names the result receiver. The initial result port writes a `text` field to an input handle.

A typed Go binding plan identifies each input field's source: widget draft, event value, typed literal or named context value. Exactly one source per field. The adapter explicitly converts text input for integer/boolean SDL fields; runtime performs no implicit conversion. Validate all bindings before installing handlers. Reject unknown modules/members/actions, missing field sources and wrong result widgets. Link maps contain source positions for SDUI widgets and SDL actions.

Domain revision is separate binding context, never the UI value revision. An accepted SDL result may update that context and the UI value. UI property batches still validate atomically; presentation failure does not roll back an already executed Go domain action and must not trigger automatic retries. `CurrentInvocation(ctx)` exposes action, model revision and sequence ID to Go functions. The profile has one ordered command source per Engine; concurrent senders must share a sequence owner. No distributed exactly-once guarantee is implied.

[EditAptCell sources and handwritten simulation](../../go/examples) demonstrate stable cell identity, separate domain revision and explicit edit rejection. The older MVP1 scenario is a requirements reference, not loaded executable code.
