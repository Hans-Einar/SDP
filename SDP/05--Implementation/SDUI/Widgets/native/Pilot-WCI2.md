# WCI2-M1 native pilot — evolving candidate

First runnable fixture binary was produced during lane integration. These runs
are development feedback, not final-source acceptance.

The panes pilot passed actual pointer/header-key callbacks, same-page silence,
native input draft retention, hidden-page tree offset retention, silent
programmatic/fallback selection, zero eligible pages, and real SDL error/malformed
result/newer-draft conflicts. The first recovery assertion correctly failed while
the deliberately dirtied receiver remained dirty. The harness now reaches that
receiver with native Tab and reverts it with Escape before independently arming
resource failure. The corrected panes run passed 25 checks, including final
resource rejection, fresh gesture recovery, six failed reload stages and valid
reload, with clean teardown. No product guard was weakened.

The first split pilot's arrow assertion compared the first intermediate drag
snapshot with the later key result. Recorded offsets proved the arrow correctly
added .05 to the final drag position. The harness now requests the final owner
snapshot after the completed gesture before calculating its next expected value.
This is evidence-harness correction, not a native divider defect.

Final hashes, complete split/provider checks, visual inspection and independent
review remain pending. Initial window resize can emit a native-minimum diagnostic
before the requested size is applied; the last valid presentation remains intact.

Expanded horizontal split pilot passed eleven checks, including admission at
220x600 while collapsed, rejected restore retaining zero extent, then successful
restore at 850x600. The page/provider lifecycle pilot passed seven checks: leaving
a page cancels its provider, late delivery is ignored, R starts a fresh request
on reveal, and hide/show retains state with silent fallback and explicit activation.

The first lifecycle input used transient AppTabs inspector geometry immediately
after native header item recreation; the new button was not yet laid out. The
harness now waits for valid native header hit geometry before clicking. The host
owner is also avoiding unchanged header item recreation on unrelated mutations.
No provider cancellation failure was reproduced after targeting the actual header.

The vertical split variant passed nine checks. Native error inspection exposed
discarded InteractionResult metadata; the host correction preserves domain
succeeded/unknown in its typed error. A further actual pointer probe from a dirty
input to its already selected tab confirmed that the header did not acquire focus.
The panes harness now retains that regression. Both fixes require final binary
rerun and independent review; the earlier 52 pilot checks do not close M1.
