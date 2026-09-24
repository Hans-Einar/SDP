# K2-M1: visible metadata in KanBan cards

| Field | Value |
| --- | --- |
| id | KB-SDP-008 |
| project | SDP |
| type | Change |
| created | 2026-09-23T22:22:08Z |
| source | owner-conversation-2026-09-24 |

## Need and scope

The owner reports that VS Code hides YAML frontmatter while XFMD renders it as headings. Metadata should be readable as a normal Markdown table.

## Implementation plan K2

Owner: Codex. Phase branch `sdp/phase-k2-readable-metadata` from K1 (`bb3728c`). K2-M1 converts metadata in every existing card and the template into one visible table immediately below the title. Preserve field names/values, content, identities and ledger history. Update the format description; no viewer/parser changes.

## Verification and outcome

Check preservation of metadata and other card content, absence of frontmatter, and valid tables, links and ledger histories. **K2-M1 delivered on 2026-09-24.** Twelve cards and the template have visible metadata tables. Before/after comparison confirmed identical field names/values in all 13 documents and unchanged content in all 12 cards. Template instructions now describe table rows. This card's outcome was then finalized.

Checks passed for three boards, 21 status directories, 12 cards, four Ref cards and 101 local Markdown links, including JSON schema and ledger replay against file locations. After closure, ledgers contain 14 events. `git diff --check` passed. This verifies documentation/structure; native VS Code and XFMD viewing was not tested. Viewer code and Markdown parsers are unchanged.
