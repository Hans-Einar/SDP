# SDP 2.4.0 release qualification

| Field | Value |
| --- | --- |
| id | SDPTOOL-VER-R240 |
| Status | completed |
| Activity | ACTIVITY-SDP-0007 |
| Plan | MAINT-SDP-0023 |

Canonical candidate: 52888aaf96c657c0d6656b104c9594e06f99d46f.
Extension packaging candidate: f7763acf27900c0fb51fd28d3f2b6edc8a40a5f7.
Engine SHA256: 2034c184c7443d7cb90bba6a30926ed0c7f99ae12063407c5ec0f282e415b9f9.
Signed descriptor SHA256: 7184efb8b5a4b1fa9280f61f5e9c839b4d6dc14ee58173a81478a7eaaa815908.

## Candidate gates

Completed PLAN-SDP-0030 implementation and independent review apply to unchanged
product sources. LaunchWorkflow model0.1.0 is released. This release changes
version/default/predecessor metadata only, adding the original signed2.3.2 digest.
Minor2.4.0 reflects new context/DocumentHost capability; legacy0.2 remains supported.
Check semantics and explicit-context obligations are described in frozen notes.

Profile/install and bootstrap race tests, 90 Toolkit tests, exact extracted-source
validation and generated release logs pass. Archive descriptor differs only in its
intentional unsigned development identity. Actual packaged signed-child test passes.
All ten original signed predecessors install and upgrade with owner documents,
Sessions and SDL preserved; managed conflicts refuse and repeated upgrade is a no-op.
Payload ownership, schemas and distributed file inventory are unchanged.

Local extension package uses the byte-identical production-signed engine/descriptor.
No private signing material is retained. Exact CI and one combined release review pass; final public outcomes are below. No live extension replacement, real Framework migration or main merge.

## Exact signed-package native outcome

The actual2.4.0 executable opens TimberFrame Analysis in isolated Xvfb with explicit
choice/SVG context. Observed native ready/resource reload/source rejection/recovery,
nonblank painting, no startup minimum diagnostic and normal window close exit0.
Initial/recovered images are pixel-identical. Detailed state assertions remain in
unchanged DocumentLaunch tests; no domain execution or new platform coverage claimed.

## Final publication outcome

Required CI44287504558 passed for exact source52888aa. Both annotatedv2.4.0 tags
and public releases identify the reviewed candidates. Downloaded assets match local
packages and API digests. Extension packaging uses the actual public immutable
canonical descriptor and redistributes its engine byte-for-byte.

Isolated ordinary `gh extension upgrade sdp` from2.3.2 and fresh installation pass.
Existing setup registration resolves to the new engine before repeating setup;
disconnected direct operations work, and project documents are unchanged. Public
default Framework installation produces a signed2.4.0 receipt and repeat no-op.
No live extension, real project Framework or main branch was changed. Public
receipts and reproducible scripts are retained in public-evidence/.
