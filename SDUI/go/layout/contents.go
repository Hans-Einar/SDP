package layout

import (
	"math"

	"github.com/Hans-Einar/SDP/SDUI/go/parser"
)

// contents uses the explicit parent's body as the scale reference. Synthetic
// rows are only track allocations and never become an ancestor reference.
func (e *Engine) contents(n *parser.Instance, inner, ancestor Size, font float64, definiteH bool) ([]placement, Size, error) {
	gx, gy := gaps(n, ancestor)
	placements := []placement{}
	top, bottom, maxW := 0., 0., 0.
	for _, role := range []string{"header", "footer"} {
		c := n.Region(role)
		if c == nil || !visible(c) {
			continue
		}
		s, err := e.desired(c, inner, inner, assigned{x: !explicit(c, "x")}, font)
		if err != nil {
			return nil, Size{}, err
		}
		y := 0.
		if role == "header" {
			top = s.H + gy
		} else {
			bottom = s.H + gy
			y = inner.H - s.H
		}
		x := alignment(inner.W, s.W, choice(c, "align-x", "start"))
		placements = append(placements, placement{c, Rect{x, y, s.W, s.H}, inner, assigned{}})
		maxW = math.Max(maxW, s.W)
	}
	body := Size{inner.W, math.Max(0, inner.H-top-bottom)}
	rows := n.Rows
	if c := n.Region("body"); c != nil {
		rows = [][]*parser.Instance{{c}}
	}
	filtered := [][]*parser.Instance{}
	for _, row := range rows {
		v := []*parser.Instance{}
		for _, c := range row {
			if visible(c) {
				v = append(v, c)
			}
		}
		if len(v) > 0 {
			filtered = append(filtered, v)
		}
	}
	rows = filtered
	if choice(n, "wrap", "none") == "wrap" && len(rows) == 1 {
		wrapped := [][]*parser.Instance{}
		current := []*parser.Instance{}
		used := 0.
		for _, c := range rows[0] {
			s, err := e.desired(c, body, body, assigned{}, font)
			if err != nil {
				return nil, Size{}, err
			}
			next := s.W
			if len(current) > 0 {
				next += gx
			}
			if len(current) > 0 && used+next > body.W {
				wrapped = append(wrapped, current)
				current = nil
				used = 0
				next = s.W
			}
			current = append(current, c)
			used += next
		}
		if len(current) > 0 {
			wrapped = append(wrapped, current)
		}
		rows = wrapped
	}
	rowPlaces := make([][]placement, len(rows))
	tracks := make([]track, len(rows))
	for i, row := range rows {
		crossHeight := 0.0
		if definiteH && len(rows) == 1 {
			crossHeight = body.H
		}
		places, height, width, err := e.row(n, row, body, gx, font, crossHeight, false)
		if err != nil {
			return nil, Size{}, err
		}
		rowPlaces[i] = places
		maxW = math.Max(maxW, width)
		tracks[i] = track{size: height}
		if definiteH && len(row) == 1 && weight(row[0], "y") > 0 {
			lo, hi := bounds(row[0], "y", body.H)
			tracks[i] = track{weight: weight(row[0], "y"), min: lo, max: hi}
		}
	}
	heights := distribute(tracks, body.H, gy)
	used := gy * float64(max(0, len(rows)-1))
	for _, h := range heights {
		used += h
	}
	extra := math.Max(0, body.H-used)
	y := top
	if definiteH {
		y += alignment(body.H, used, choice(n, "justify", "start"))
		if choice(n, "justify", "start") == "between" && len(rows) > 1 {
			gy += extra / float64(len(rows)-1)
		}
	}
	for i, row := range rows {
		places := rowPlaces[i]
		if tracks[i].weight > 0 {
			var err error
			places, _, _, err = e.row(n, row, body, gx, font, heights[i], true)
			if err != nil {
				return nil, Size{}, err
			}
		}
		for _, p := range places {
			p.rect.Y += y
			placements = append(placements, p)
		}
		y += heights[i] + gy
	}
	naturalH := top + used + bottom
	if len(rows) == 0 {
		naturalH = math.Max(0, top+bottom-gy)
	}
	// In content-height measurement, footer follows the body, not the temporary bound.
	if !definiteH {
		for i := range placements {
			if placements[i].node == n.Region("footer") {
				placements[i].rect.Y = naturalH - placements[i].rect.H
			}
		}
	}
	return placements, Size{maxW, naturalH}, nil
}
func (e *Engine) row(owner *parser.Instance, row []*parser.Instance, ref Size, gap, font, height float64, forceH bool) ([]placement, float64, float64, error) {
	tracks := make([]track, len(row))
	for i, c := range row {
		lo, hi := bounds(c, "x", ref.W)
		w := weight(c, "x")
		_, hasRatio := c.Layout["ratio"]
		if len(row) == 1 && !hasRatio && choice(owner, "items", "start") == "stretch" && !explicit(c, "x") {
			w = 1
		}
		if w > 0 {
			tracks[i] = track{weight: w, min: lo, max: hi}
			continue
		}
		s, err := e.desired(c, ref, ref, assigned{}, font)
		if err != nil {
			return nil, 0, 0, err
		}
		tracks[i] = track{size: s.W}
	}
	widths := distribute(tracks, ref.W, gap)
	sizes := make([]Size, len(row))
	rowH := height
	used := gap * float64(max(0, len(row)-1))
	for i, c := range row {
		slot := Size{widths[i], ref.H}
		if forceH {
			slot.H = height
		}
		s, err := e.desired(c, ref, slot, assigned{x: tracks[i].weight > 0, y: forceH}, font)
		if err != nil {
			return nil, 0, 0, err
		}
		sizes[i] = s
		rowH = math.Max(rowH, s.H)
		used += s.W
	}
	x := alignment(ref.W, used, choice(owner, "justify", "start"))
	if choice(owner, "justify", "start") == "between" && len(row) > 1 {
		gap += math.Max(0, ref.W-used) / float64(len(row)-1)
	}
	out := []placement{}
	for i, c := range row {
		s := sizes[i]
		align := choice(c, "align-y", choice(owner, "items", "start"))
		_, hasRatio := c.Layout["ratio"]
		if align == "stretch" && !hasRatio && !explicit(c, "y") {
			s.H = rowH
		}
		y := alignment(rowH, s.H, align)
		out = append(out, placement{c, Rect{x, y, s.W, s.H}, ref, assigned{x: tracks[i].weight > 0, y: forceH}})
		x += s.W + gap
	}
	return out, rowH, used, nil
}
