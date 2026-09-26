# Refactors

Use a refactor programme when the required change is broader than a normal
corrective sprint or code review.

Recommended structure inside each refactor folder:

```text
01--Mandate/
02--Study/
03--Requirements/
04--Architecture/
05--Design/
06--Implementation/
```

Execute refactors through vertical, behavior-preserving slices with explicit
baseline and verification evidence.
## Project-management integration

Use an optional KanBan CodeReview/Refactor card or select the work directly from
a Scrum. Record lifecycle in [ProjectManagement](../ProjectManagement/README.md),
with one scoped review/refactor record and links to input/output evidence.
Only actual system design/code changes and verification enter Traceability,
referencing that management work. A review's completion is not proof its findings
were fixed. Preserve old review/refactor IDs and records.
