# SDP Vertical Refactor

Use this skill for architecture-changing work that must preserve behavior.

## Procedure

1. Capture a baseline: responsibilities, state ownership, dependencies, user
   flows, tests, build and rendered behavior.
2. Define the target architecture and measurable exit criteria.
3. Plan vertical migration Slices, each leaving the application runnable.
4. Introduce the minimum contract required for the next real migration.
5. Move one complete workflow, including state, provider interaction, rendering
   and verification.
6. Keep compatibility adapters explicit and assign each a removal Slice.
7. Review behavior and coupling after every migration.
8. Remove obsolete paths only after all consumers have migrated.

## Guardrails

- Do not split files first and call it architecture.
- Do not mix visual redesign with structural migration unless required.
- Do not introduce abstractions solely for one hypothetical consumer.
- Reduced line count is not completion evidence.
