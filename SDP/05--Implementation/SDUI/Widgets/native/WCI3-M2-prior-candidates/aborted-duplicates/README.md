# Intentionally aborted duplicate test attempts

Crossed asynchronous coordinator/worker handoffs briefly started duplicate SDL
regression groups. Both initial runs and a worker retry were explicitly stopped
by coordinator instructions; their partial logs are diagnostic records, not test
failures or acceptance. A worker filename says CANCELLED-BY-OWNER: owner there
means the coordinating agent, not a user request to stop this card. The sole final
coordinator run is retained in WCI3-M2-final/suites/sdl-current-group and passed.
