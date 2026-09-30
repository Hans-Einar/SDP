# Proposed Frontend source-set fixture

Design-only fixture for [SSD1](../../Contract.md). Current Go rejects version 0.6;
this fixture does not demonstrate working new language support.

The tracked source SDP/SDL/SDL/Frontend/System.design is copied by partitioning
its existing declarations and facts. Features holds actors, use cases, features
and activities and facts whose subjects are those objects. Containers/Frontend
holds the remaining declarations/facts. System adds Frontend and its explicit
membership of SdlCommandProcess. Every existing baseline statement is retained.
A file may refer to names declared in another member. File paths are organizational,
not new visibility or ownership boundaries. No real source is moved or registered.

[Acceptance](../../Acceptance.md) defines negative mutations and expected results.
Do not run a private preprocessor to make this fixture pass the current parser.
