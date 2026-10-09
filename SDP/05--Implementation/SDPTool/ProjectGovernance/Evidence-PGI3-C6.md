# PGI3-C6 — explicit skill invocation and exact correlation

PLAN-SDP-0018 / KB-SDP-038, 2026-10-09. Candidate: parent `922aff97db68941a618cfea69a13ea22a7deaecb` plus
[the recorded product/test file set](live-PGI3-C6.json), digest `b58ecaf59b38ab124e392fe9fd962dc4e47961750f043cfc45aecab7696869c2`.
Same-context implementation and verification, not independent review, completed
PGI3, owner acceptance or release.

## Delivered behavior

`run --skill PROJECT_RELATIVE_SKILL_MD` selects up to eight enabled project skills.
Fresh discovery, unique names/paths and local hashes are checked before recording
an intent. Exact owner text remains separate from outbound assignment text and
explicit name markers. Name/path/digest selections are immutable across prepare,
claim, uncertain and acknowledgment states. Local drift refuses before claim or
turn submission. Default text-only launches, steering, sandbox, scoped MCP and
no-delegation instructions retain their existing boundaries.

The adapter fingerprints ordered text and skill name/path inputs. Only absent or
empty text_elements is normalized. Unknown content, changed/missing/extra/reordered
skills, duplicate client messages and partial history cannot prove acceptance.
Recovery correlates the one original client ID and exact content without resending.
Offline provenance exposes invocation attempt/input identity, state and acceptance:
unconfirmed, response-acknowledged or host-message-observed. Compliance remains
not-observed. Inventory evidence retains its availability-only meaning.

## Actual workflow evidence

The preliminary read-only Codex 0.161.0 probe confirmed explicit skill items in
both completed notifications and full thread history. The
[official protocol reference](https://learn.chatgpt.com/docs/app-server) describes
text markers plus skill inputs. Returned content proves accepted input identity;
it does not attest exact consumed bytes or procedural compliance.

Two actual final-binary scenarios use a disposable project with maintained SDP
skills and discovery links. The Maintenance scenario sends explicit sdp and
sdp-master inputs, fixes Unicode boundary trimming, receives one exact report
approval and reaches verifying. Only the two permitted Go files change; driver
Go tests pass and other preimages/discovery links are preserved. The driver decision
is not a verbatim owner reply or completed-work acceptance.

The second scenario suppresses the real turn/start response after forwarding the
request. The agent finishes a read-only verification task, while the controller
retains uncertain input. A fresh controller reads history, matches exact text plus
both skill identities and acknowledges the original input. One attempt/thread/turn
is retained, with an explicit epoch gap and historical context. No second turn is
submitted. The routine remains running rather than being accepted from host
completion. All fixture preimages remain identical during this scenario.

## Checks, corrections and limits

Final affected race tests including the compiled controller CLI, vet, application
and installed-host schemas pass. Negative cases include disabled/ambiguous/external
selection, duplicate input, changed snapshots, changed skill identity, unsupported
metadata and partial history. Source hashes and actual binaries are retained in the
JSON alongside exact observed events, input records, driver scripts and outputs.
Management/Toolkit and final candidate/ledger/link checks are added after records.
No full-repository regression or publication claim is made.

Initial development failures were setup/assertion errors: a module-relative output
path prevented initial adapter file creation; an assertion omitted marker-separating
whitespace; a wire capture omitted the actual epoch/time supplied by its helper.
Each was corrected without weakening runtime requirements. Trial1 passed an earlier
candidate; trial2 repeats actual execution on final binaries after explicit invocation
attempt identity and complete-history guard refinements.

Native child/review mapping remains unverified under current no-delegation rules.
It is not simulated by changing the assigned role name. Next: resolve that PGI3
boundary and establish the PGI4 candidate/evidence/publication handoff; PGI5 still
requires independent integrated review. Current context snapshots cannot prove later
host consumption. No merge, push or release is selected; known private C1 event-ID
collisions remain for deliberate integration reconciliation.
