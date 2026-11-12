# Activity delivery requirements

This Activity implements the existing [design and acceptance cases](../../../04--Design/SDPTool/ProjectGovernance/Activities/Design.md).
Current system sources stay authoritative. Requirements are grouped by observable
outcome, not by Go file: transactional Activity/Session creation and retry; checked
lifecycle and one active Activity; legacy/manual compatibility; producer-owned
Activity navigation; lifecycle intervals and attributed commit/ledger markers;
Activity-scoped PG admission; installer and native consumer compatibility.

PLAN-SDP-0027 must account for each group. No unfinished group can be relabeled as
complete merely because create works. Blueprint/code mapping and integration tests
must supplement SDL facts where the model cannot prove runtime behavior.
