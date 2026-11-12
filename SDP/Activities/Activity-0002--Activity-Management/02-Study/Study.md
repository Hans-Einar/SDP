# Existing study and implementation findings

The [need study](../../../02--Requirements/ProjectGovernance/Activities/Study.md)
and [selected design](../../../04--Design/SDPTool/ProjectGovernance/Activities/Design.md)
remain authoritative. Activity0001 completed the prerequisite architecture.

Initial source inspection identifies two necessary refinements: projecthistory's
transaction reservation does not itself publish multiple records, and linked
Session creation must retain old descriptor bytes so historical hashes remain
valid. Resolve these in the executable contract before dependent implementation.
Old projects retain Sessions navigation until Activity adoption; manual Activity
records must be read without renumbering or invented historical timestamps.
