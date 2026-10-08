# ModelGovernance concern entry

[ModelGovernance.design](ModelGovernance.design) is a composed design-core/0.6 entry
for the ModelGovernance feature of SDPTool. It is an independent concern view of
the same System/host, not another deployed system or an import of System.design.
The existing 0.5 routine overview remains authoritative for its own concern.
No directory-wide concatenation or migration is implied.

- Features/ModelGovernance.design owns actors, goals and functionality allocation.
- Containers/SdpToolHost/ModelGovernance.design owns internal units/library boundaries.
- Contracts/ModelGovernance.design owns logical calls and declared example scenarios.

All implementation activities remain planned. Libraries are interfaces/units, not
new runtime containers. Use the selected entry to resolve all fragments together.
[Contract](../../../04--Design/SDPTool/ModelGovernance/Contract.md) specifies the
filesystem behavior; [evidence](../../../04--Design/SDPTool/ModelGovernance/Evidence.md)
records successful validation and generated views. No application code delivered.

From repository root with Go available:

```sh
go -C SDL/go run ./cmd/sdl check ../../SDP/SDL/ProjectGovernance/SDPTool/ModelGovernance.design
```
