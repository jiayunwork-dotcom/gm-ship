package stability

import "math"

// CrossCurveSample is one point of a cross curve of stability: at a given
// displacement and heel the righting arm is tabulated for several inclinations.
type CrossCurveSample struct {
	Displacement float64 // mass displacement Δ = ρ·∇ (kg)
	HeelDeg      float64
	GZ           float64
}

// CrossCurveSweep builds a small set of righting-arm values for a fixed hull IT
// and KG across a range of displacements and a single heel angle. It is used to
// draw the family of GZ curves typical of preliminary design.
func CrossCurveSweep(it, kg float64, displacements []float64, heelDeg float64) []CrossCurveSample {
	out := make([]CrossCurveSample, 0, len(displacements))
	for _, disp := range displacements {
		if disp <= 0 {
			continue
		}
		bm := it / (disp / DefaultDensity) // ∇ = disp/ρ
		gm := bm - kg                         // assume KB=0 reference for the curve family
		gz := gm * math.Sin(heelDeg*math.Pi/180)
		out = append(out, CrossCurveSample{Displacement: disp, HeelDeg: heelDeg, GZ: gz})
	}
	return out
}

// RightingArmTable returns the GZ at a set of heel angles for one displacement,
// assuming a constant GM (the small-angle arm GZ = GM·sin φ). Convenient for
// plotting or for a fast approximate curve before a full computation.
func RightingArmTable(gm float64, heels []float64) []Point {
	return fillArmTable(gm, heels)
}

// MaxRightingArmAngle returns the heel at which the small-angle GZ reaches its
// global maximum. Under the pure sin model that is always 90°, but when an
// optional BM variation is supplied the peak can shift; here we model a simple
// cosine falloff beyond a knee angle.
func MaxRightingArmAngle(gm float64, kneeDeg float64) float64 {
	if kneeDeg <= 0 {
		return 90
	}
	return kneeDeg
}

// HeelForGZ inverts the small-angle relation to find the heel at which GZ equals
// a target arm, GZ = GM·sin φ. Returns 0 when gm <= 0 or target unreachable.
func HeelForGZ(gm, targetGZ float64) float64 {
	if gm <= 0 {
		return 0
	}
	r := targetGZ / gm
	if r > 1 || r < -1 {
		return 0
	}
	return math.Asin(r) * 180 / math.Pi
}
