# Go generation — G5 profile

The generator validates the whole source before writing Go. SDUI 0.2 becomes a source-positioned Document and selected normalized Root; SDL action-core 0.1 becomes a Program with records, actions and Go symbol references. It emits typed Go literals, not embedded source reparsed at startup. Each constructor call creates independent models. Original positions and binding references are preserved.

Generated code contains model/binding data, not another runtime. The host registers handwritten Go functions and explicit bridge.Plan field bindings, using the same validation as file-based startup. Generation invents no domain functions, storage, Channel transport or executable interpretation of design-core. Wrong profiles and unknown roots are rejected before publication.

`sdl-gen -actions actions.sdl -ui page.sdui -entry page -package model -output DIR` creates actions_gen.go, ui_gen.go and manifest.json with source hashes/generator version. Publication reuses manifest-managed atomic document publication: preserve user-owned files and reject modified generated files. Generated code builds against the generator's Go module versions. Build Go source normally; language-model hot reload and Go rebuild/restart are distinct mechanisms.
