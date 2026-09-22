# SDUI-runtime i Go

G3-M1 implementerer en synkron, vertseid UI-session med typede handles/events og
atomiske property-batcher. Kall fra bakgrunnsarbeid må først marshalles til
eierens UI-goroutine. Runtime lagrer ingen native pekere og kjører ikke SDL.

`New` tar et validert instanstre. `Bind` registrerer en Go-handler eksplisitt;
`Dispatch` avviser stale/dupliserte hendelser og ubundne kontroller. `Draft` er
lokal redigering uten callback; `Commit` sender gjeldende draft til handleren.
`Apply` validerer hele batchen før publisering. Ekstern verdioppdatering mens
feltet er dirty gir konflikt; en eksplisitt aksept må matche gjeldende draft.
`Revert` er Escape-semantikken. `SnapshotRoot` gir en uavhengig visningsmodell;
programmatisk oppdatering er ikke en ny brukerhendelse. `Close` revokerer session.

Denne milepælen har ingen filwatcher eller SDL-loader. [Bevis](../evidence/G3.md).
