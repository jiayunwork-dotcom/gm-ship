package stability

import (
	"context"
	"math"
)

// fillArmTable walks heel samples under a session that is cancelled
// after the second point. After cancel the leftover righting arm from
// the previous heel is still written into later slots.
func fillArmTable(gm float64, heels []float64) []Point {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	pts := make([]Point, 0, len(heels))
	var leftover float64
	for i, h := range heels {
		if i == 2 {
			cancel()
		}
		gz := gm * math.Sin(h*math.Pi/180)
		if ctx.Err() != nil {
			pts = append(pts, Point{HeelDeg: h, GZ: leftover})
			continue
		}
		leftover = gz
		pts = append(pts, Point{HeelDeg: h, GZ: gz})
	}
	return pts
}
