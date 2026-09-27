# Shared MVP1 design components

This directory collects definitions used across MVP1 containers. It is not a
runtime container, a global state instance, or a new deployment boundary. Every
Unit retains its explicit SDL owner. All sources remain inside MVP1 for now.

| Source home | Responsibility | Reuse boundary |
| --- | --- | --- |
| Libraries/UIRuntime/ | Representation identity, composition lifetime, transport, delivery and command correlation | Generic mechanism candidates; per-application/session instances, never shared domain state |
| Libraries/BoxUI/ | Presentation validation/publication and abstract renderer adapter | Generic UI mechanism candidates; concrete Fyne/React implementations remain external to this corpus |
| Libraries/ServiceKit.design | Service sessions and command lifecycle | Cross-container infrastructure; not yet a separately versioned cross-project contract |
| Libraries/LabKit.design | Diagnostic recording and replay | Shared diagnostic infrastructure |
| Libraries/Contracts.design | ContractRegistry governance/binding generation | MVP1 contract ownership; not a second copy of authoritative wire schemas |
| Libraries/BuckingCore.design, TaperModels.design, StanfordClassic.design | Bucking, taper and codec responsibilities | Domain libraries shared inside MVP1; future extraction needs explicit contracts |
| UI/Boundary.design | Logical UI subsystem layers and cross-cutting constraints | MVP1-specific system composition, not a reusable library |
| UI/Views.design | ServiceBanner and UnhandledDiagnostics definitions | MVP1 common UI vocabulary; no singleton instance implied |

BuckingUI's projection and views live under Containers/BuckingUI. SimulatorUI's
projection and hidden-truth view live under Containers/SimulatorUI. Simulator-only
truth must not leak into normal machine/operator evidence. The UI subsystem still
groups both applications and shared mechanisms in the model; their separate
Container declarations are preserved in System.design.

A component earns cross-project reuse through a stable dependency boundary and
actual consumers, not by being put in a folder named Shared. Later extraction can
move selected library sources and introduce versioned references without copying
MVP1-specific policy. Keep the existing BoxUI model identity until an explicit
migration selects its relationship to the current SDUI runtime.

No SDUI cross-file component imports are supported in profile 0.2. Therefore each
application screen is self-contained and may reuse local named groups. Do not
invent a Shared.sdui import mechanism or duplicate the SDUI implementation here.
