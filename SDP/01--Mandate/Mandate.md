# Mandate — develop SDP using SDP

Project: SDP. Owner: Hans Einar. Current local application: 2026-09-25.
This records the owner's current assignment and its bounded application, not a
reconstruction or replacement of every historical SDP mandate.

Develop SDP as the process and tools for designing systems in SDL, prototyping
interfaces in SDUI, planning implementation and tracing design to actual evidence.
Use this repository's own SDP area to design and deliver SDP product features.

The owner accepted the five-phase structure in KB-SDP-001 and explicitly requested
its first real use here. Apply 01--Mandate, 02--Requirements, 03--Architecture,
04--Design and 05--Implementation, with optional studies within the relevant phase.
Keep common KanBan, Traceability, Verification and other operational records at
project level. Phase homes and SDL A0–A5 abstractions have different purposes.

## First feature: SDPTool

Create sdptool as one SDP product feature, with code in Toolkit/SDPTool and design
in this shared project area. Do not create a separate SDP project, board or ledger
under Toolkit/SDPTool. Its first candidate delivery is a facade for saved .design
preview, reusing existing Go SDL projection and Rust layout. Further scope includes
project recognition, general navigation and later process/traceability services.

Stakeholders: the Project Owner supplies scope and priorities; design authors and
reviewers consume model/document tools; Toolkit and language maintainers own reusable
contracts; XFMD integrators consume navigation/preview services. XFMD owns its
native UI implementation and continues its own development process.

The feature's [requirements](../02--Requirements/SDPTool.md),
[active card](../Agents/KanBan/active/%23017--Proposal--sdptool-and-project-navigation.md)
and [implementation plan](../05--Implementation/SDPTool.md) record acceptance and
current work. The [pipeline study](../../Toolkit/SDPTool/Navigation-and-Design-Preview.md)
provides reuse evidence. No additional mandate study is needed for this explicit
placement decision; unresolved API and discovery details belong in Design.

Applying the local directory profile does not migrate distributed templates,
select a discovery descriptor or make the future sdptool resolver recognize it.
Those have separate implementation/compatibility requirements.
