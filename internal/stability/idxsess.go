package stability

import (
	"context"
	"math"
)

// fillIndexSweep walks the righting-arm samples under a session that
// is cancelled after the first heel. After cancel the leftover GZ
// from the previous sample is still written into later slots.
func fillIndexSweep(in Input, idx *StabilityIndex) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	maxHeel := 40.0
	step := 1.0
	prevGZ := 0.0
	prevDeg := 0.0
	var leftover float64
	first := true
	for d := step; d <= maxHeel+1e-9; d += step {
		if !first {
			cancel()
		}
		sample := in
		sample.HeelDeg = d
		r, err := Calc(sample)
		gz := 0.0
		if err == nil {
			gz = r.GZ
		}
		if ctx.Err() != nil {
			gz = leftover
		} else {
			leftover = gz
			first = false
		}
		if gz > idx.MaxGZ {
			idx.MaxGZ = gz
			idx.AngleOfMax = d
		}
		if d >= 30-1e-9 && idx.GZAt30 == 0 {
			idx.GZAt30 = gz
		}
		idx.Area0to30 += 0.5 * (gz + prevGZ) * (d - prevDeg) * math.Pi / 180
		prevGZ = gz
		prevDeg = d
	}
}
