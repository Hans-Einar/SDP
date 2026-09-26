# SDL class-core 0.1

Explicit G6-M6 class-design profile; it does not reinterpret design-core facts. Uses the shared SDL lexer/source positions. All declarations precede facts, sorted by name; facts sort by canonical text.

`class` declares a class. `attribute Name as Type` and `operation Name returns Type` define attributes and parameterless operation signatures. Type is text/integer/boolean or a declared class name. These are signatures, not implemented functions or an executable object runtime.

`association` declares a relationship between two explicit classes, with role names and multiplicity at both ends. `multiplicity 0 to many` means 0..*; `1 to 1` means 1. Numeric bounds must be 0–65536; many is the only open upper bound. Each end gives the number of objects at that end per object at the opposite end.

Each association requires exactly one `links` and one `ownership`:

- none: ordinary association, without ownership/lifetime semantics.
- aggregation: first end is the whole; shared parts, without implied cascading deletion.
- composition: first end is the whole; a part has at most one whole at a time and participates in its lifetime. The first end's upper multiplicity must therefore be 1.

Composition between class types must also be acyclic. This is a strict, bounded model rule, not a claim to cover all UML models. Roles must differ and be unique among a class's association ends. Attribute/operation names are unique within each compartment. Limits: 128 classes, 256 associations and 64 members per class. Undeclared types, invalid ranges, repeated/incomplete relations and unjustified composition are rejected.

Mermaid classDiagram places an open aggregation diamond or filled composition diamond at the first end. Ordinary associations use a line. Both role names and multiplicities accompany the diagram. Contains/owns/allocation are not translated into these relations. Class profiles do not automatically join the action-core runtime.
