package layout

import "math"

type track struct{ size, weight, min, max float64 }

// distribute solves size=clamp(lambda*weight,min,max). Fixed content and scale
// tracks reserve space first; they never silently shrink when the row overflows.
func distribute(t []track, available, gap float64) []float64 {
	rest := math.Max(0, available-gap*float64(max(0, len(t)-1)))
	out := make([]float64, len(t))
	maxWeight := 0.
	for i, v := range t {
		if v.weight == 0 {
			out[i] = v.size
			rest -= v.size
		} else {
			maxWeight = math.Max(maxWeight, v.weight)
		}
	}
	if maxWeight == 0 {
		return out
	}
	rest = math.Max(0, rest)
	sum := func(lambda float64) float64 {
		s := 0.
		for _, v := range t {
			if v.weight > 0 {
				s += clamp(lambda*(v.weight/maxWeight), v.min, v.max)
			}
		}
		return s
	}
	lo, hi := 0., math.Max(1, rest)
	for j := 0; j < 64 && sum(hi) < rest; j++ {
		hi *= 2
	}
	for j := 0; j < 80; j++ {
		mid := (lo + hi) / 2
		if sum(mid) > rest {
			hi = mid
		} else {
			lo = mid
		}
	}
	for i, v := range t {
		if v.weight > 0 {
			out[i] = clamp(lo*(v.weight/maxWeight), v.min, v.max)
		}
	}
	return out
}
