# N12 — encoding inventory and explicit packet layouts

Session0007 T008 / PLAN-XFMD-0009 / KB028 / producer KB032. Owner continuation
selects VP10. BranchPolicy: sdp/xfmd-vp10-overview stacked on N11 6b80b40;
CommitPolicy: scoped N12 milestone. Phase push permitted; no merge/release/install.

## Selected design

Reuse the existing validated packet diagram per Encoding in the complete VP10
overview, with separate identities and preserved detail selections. A combined
packet would invent a layout across independent Encodings. Source-linked registers
show Encoding, Variant, closed Contract, explicit byte/bit order and total bit
length, followed by ordered field positions, widths, types, required presence and
header/variant ownership. Register every Variant, including those without Encoding;
missing layout is an explicit gap, never inferred from a logical Contract.
No fixed overview content cap or partial truncation. The language's validated
65536-bit encoding limit is unchanged. Empty Encoding/Variant inventories are
explicit. Reject mode/focus/diagram/level/relation/traversal filters on overview.

Byte/bit order are source declarations; positions are logical contract bit offsets.
Do not infer serialization implementation, executable codec, padding, variable
lengths, optional-field layout or byte ordering transformations. Only checked
closed fixed layouts are drawn. All relevant facts have source references;
legacy detail diagrams/tables retain their prior behavior. No authored model edits.

## Ownership and verification

SDL owns projection and context, SDPTool advertises overviewTarget via its existing
contract, XFMD selects/renders/opens source generically. No native product change
is anticipated. Verify multiple Encodings, both order settings, shared header,
non-byte-aligned fields, unencoded Variants, exact offsets/proofs, deterministic
selection, filter rejection, empty and more-than-24 packet inventories. Verify
invalid layouts remain rejected, composed sources/freshness, actions, full suites,
native tree overview/detail/source handling and PDF retention. Actual XFMD has no
Encodings; validate that truthful empty result, and use fixtures for packet output.
Concurrent-work preflight is recorded in XFMD N12-001. Existing root-discovery and
renderer/large-print residuals remain KB025/KB020. VP11 is next proposed.

## Delivered candidate and residuals

N12-candidate.json / XFMD N12-001 identify exact producer/native inputs. Full
SDL/SDPTool suites, vet and 12 consumer CTests pass. Native fixture verifies two
independent 29-bit Encodings, opposite byte/bit orders, shared 1-bit header, signed
12-bit field, bytes crossing a byte boundary, unencoded optional Variant and
original packet detail. All 34 distinct source spans are in bounds; six placement
spans match exact statements. Actual XFMD has no Encodings or Variants; its empty
page is verified natively and in PDF. Producer tests cover 32 packets, invalid
layouts, full provenance, filter rejection, determinism and detail preservation.

Visual inspection found existing narrow packet label clipping (1-bit PacketFlag)
and crowded adjacent 12/13 boundary labels, reproduced on the unchanged original
detail diagram. KB020 owns this renderer follow-up with screenshots/fixture. The
new field register retains complete names/positions, but this is not visual diagram
acceptance. PDF text extraction is not proof of visible label completeness. The
initial native detail test read the prior overview because both contained the same
title; waiting for changed document path/current render token fixes the test, and
all final native checks pass. No native product behavior change was needed.

This bounded projection is delivered; broader plan/cards and KB020/025 residuals
remain active. Stacked candidates await integration. VP11 is next proposed. No
merge/release/install, authored model edit, independent approval or large-print
acceptance. Shared SDP-vNow receives existing KB032/Integration handoff only.
