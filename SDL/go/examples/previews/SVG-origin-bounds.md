# SVG origin regression: expected native green bounds

Bounded fixture delta in `/tmp/wci4-svg-origin-fix` only. No product API,
SDL, module, host or default-source change. Resources are chosen by the existing
`resource Hero SLOT` + `reload` condition commands; ReportSVG and DetailSVG
accept the same slots. Main owns native execution and final matching builds.

All resources declare dimensions 400×200. Solid green is exactly `#238248`,
RGB (35,130,72), opacity 1, with no stroke. The root-transform resource contains
only the specified green rectangle; its remaining area is transparent.

| Slot | Root viewBox | Root transform | Source rectangle (x,y,width,height) |
| --- | --- | --- | --- |
| negative-origin (unchanged) | -40 -20 400 200 | identity | -40,-20,400,200 |
| positive-origin | 40 20 400 200 | identity | 40,20,400,200 |
| mixed-origin | -40 20 400 200 | identity | -40,20,400,200 |
| root-transform | 40 20 400 200 | scale(0.5) | 120,60,160,80 |

Use the selected SVG2 composition **S × Rroot × Torigin × Cchildren**.
Here Torigin subtracts the declared viewBox origin; Cchildren is identity.
For root-transform, Torigin maps opposite source corners (120,60), (280,140)
to (80,40), (240,120); Rroot then maps them to (40,20), (120,60).
The root transform and origin translation do not commute: reversing them would
produce (20,10), (100,50), a visibly different rectangle. This slot deliberately
detects that error independently of root-transform-free origin translation.

Let the host's **already fitted image rectangle** from Inspect.previews be
I=(X,Y,W,H), with W/H=2. This excludes any caption band and letterbox allocation.
Let s=W/400=H/200. Expected un-clipped painted bounds, as continuous edge positions:

- negative-origin, positive-origin, mixed-origin: [X,Y,X+400s,Y+200s].
- root-transform: [X+40s,Y+20s,X+120s,Y+60s], width80s and height40s;
  equivalently [X+0.1W,Y+0.1H,X+0.3W,Y+0.3H].

| Uniform scale s | Image dimensions | Full-fill green bounds relative to I origin | Root-transform green bounds relative to I origin |
| --- | --- | --- | --- |
| 0.5 | 200×100 | [0,0,200,100] | [20,10,60,30] |
| 1 | 400×200 | [0,0,400,200] | [40,20,120,60] |
| 2 | 800×400 | [0,0,800,400] | [80,40,240,120] |

Bounds notation is [left,top,right,bottom], not x/y/width/height. The formula
applies at arbitrary positive scale, including fractional fit sizes. Intersect
the painted bounds with the actual shared clip and owning canvas/window bounds
after fitting; clipping must never recenter or rescale artwork. Convert canvas
units to OS pixels using the actual canvas/window device scale. Add the owning
window origin only once. For modal/nonmodal surfaces use that surface's actual
canvas/image rectangle, not the main window's coordinates.

At integer-aligned physical edges, interior pixels have the exact green color.
Treat the outer antialiased edge separately (for example, allow one physical
pixel at the boundary); test interior points and outside margins as well as the
observed green bounding box. A wholly clipped image has no visible green area.
No host rasterization or native screenshot was run by this fixture worker;
resource validation alone does not prove any of these native pixel outcomes.
