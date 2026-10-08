# WCI3-M2 pinned native input dependency — verified SDL-root selection

Selected after scalar delivery ea49991 and handoff d674274. The local dependency
is SDUI/third_party/glfw, upstream Go module github.com/go-gl/glfw/v3.4/glfw at
v0.1.0-pre.1.0.20260707082822-2a407d02d01a. All 145 upstream source/license files
are preserved; only glfw/src/x11_window.c has the reviewed one-condition X11 fix.
The policy directory retains the reversible patch and upstream/selected hashes.

Coordinator and independent reviewer verified every selected hash and exact patch
reversal to upstream bytes. Reviewer independently recomputed the upstream ZIP h1
against both go.sum files and checked the untouched global module cache. The
verification script grants write permission only to its temporary reverse-patch
copy, accommodating read-only cache-derived directory modes.

Both SDUI/go and SDL/go explicitly replace the module with the same selected local
source; GOWORK=off go list -mod=readonly -m confirms resolution. No transitive
replacement is assumed. The inspected XFMD helper builder uses SDUI/go as its
build root; generated native applications in another module need an explicit
matching replacement. SDPTool does not itself compile Fyne. No installed helper,
module cache, owner display or external application source was modified.

Independent reviewer 01a11854-523d-7c11-adf4-f622b19b68c6 approved the selected provenance and conditional scope. Actual SDL-root build information and configured native IME composition/commit/cancel now pass; see Evidence-WCI3-M2 and its raw archive. WCI4 must still prove the staged SDUI-root helper/fixture recipes and include upstream notices/native libraries. No publication is inferred.
