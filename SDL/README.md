# SystemDesignLanguage

The [documentation map](docs/README.md) distinguishes active profiles, studies and historical snapshots. [SDL KanBan](SDP/Agents/KanBan/README.md) tracks further language proposals. R1 moved sources from SystemDesignLanguage to SDL, retaining the Go module identity `github.com/Hans-Einar/SDP/SystemDesignLanguage/go`. Separate repository/module extraction has not occurred.

The SDL tool is implemented in Go and generates Markdown/Mermaid/SVG viewpoints from validated model facts. G1–G6 are delivered within explicit profile boundaries. The [dated checkpoint #1 implementation snapshot](../SDP/History/checkpoint-1/11-Go-Implementation-and-Navigation.md) distinguishes this from broader candidates and future Ponsse production logic.

- [Go tools and commands](go/README.md): check/AST, viewpoints, dynamic selection, document service, runtime and Go generation.
- [SDL model of the SDUI system](../SDUI/design/README.md) and [generated development plan](../SDUI/design/viewpoints/implementation.md).
- [Shared phase plan](../SDUI/docs/implementation-plan.md) and [Git traceability](../SDP/Development-Branch-Stack.md).

Design-core 0.5 describes structure, use cases/features, allocation, data contracts, Channels, scenarios and delivery-plan facts. Action-core 0.1 runs typed actions through registered handwritten Go functions. Class-core 0.1 describes explicit classes, roles, multiplicity and aggregation/composition; it is not an object runtime. No old Python design-core or viewpoint projector remains active.

The [MVP1 corpus](../experiments/mvp1_sdl/README.md) remains a separate 66-file candidate exercise, not executable design-core or completed machine logic. EditAptCell is a labeled simulation, not a Ponsse backend. [Mandate/study](Mandate-and-Study.md), [research](research/README.md) and [SDP #10](https://github.com/Hans-Einar/SDP/issues/10) retain their original authority; they do not change canonical Toolkit/SDP contracts.
