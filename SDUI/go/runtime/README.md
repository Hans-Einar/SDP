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

G3-M2/M3 leverer hashbasert filwatcher i `../reload`, siste gyldige modell og
native Fyne-adapter. Nye kandidater blir parse-/profilvalidert og målt før
publisering. `Reload` beholder verdi/draft/fokus og handle for samme navngitte
bane/type. Nye defaults gjelder nye instanser; etikett og enabled/visible fra
ny kilde gjelder straks. Anonyme widgets får nye generasjoner. Endret callback-
referanse gjenbruker ikke gammel handler. Typebytte/sletting revokerer handle.

CLI overvåker filen som standard (`-watch=false` slår av). Inputendring er draft;
Enter committer, Escape gjenoppretter akseptert verdi. Native callbacks fanger
modellrevisjonen de ble opprettet for. Gamle callbacks og resultater etter
teardown eller reload utfører ingen ny handling. Programmatisk `SetText` er
mutet hos adapteren, slik at en propertyoppdatering ikke blir et domene-event.
Session har fortsatt ingen SDL-loader. [Bevis](../evidence/G3.md).
