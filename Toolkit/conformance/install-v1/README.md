# Historical install-v1 contract data

The JSON scenario index and expected plans/failure records preserve the retired
installation contract for schema inspection and historical compatibility analysis.
There is no executable reference installer or scenario runner in the current tree.
Repository Python validators can inspect these records without running an old engine.

Current release verification uses [SDPTool installation tests](../../../SDPTool/install/Records.md)
for actual signed preview/apply, recovery, preservation and receipt behavior.
These historical expectations are not evidence that a current Go operation ran.
Published old assets and their original Git commits retain the prior implementation.

The retained expected-plan toolkitVersion is pinned to the final 1.0.0 baseline.
Current distribution version changes do not rewrite historical scenario bytes;
validation retains their schema, ordering, path and baseline-version checks.
