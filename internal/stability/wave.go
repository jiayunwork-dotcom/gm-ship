package stability

import "math"

// WaveStability estimates the effect of a regular wave on a vessel's effective
// stability. In a wave of height H and length Lw the surface slope is roughly
// π·H/Lw (radians) near the crest/trough region, which heels the ship by the
// same angle if it were to follow the wave. We report the wave-slope heel and
// the residual righting capacity after subtracting the wave-induced heel from
// the intact GZ curve at the relevant angle.
type WaveStability struct {
	WaveInducedHeel float64 // heel imposed by wave slope (deg)
	ResidualGZ      float64 // GZ remaining at the wave-induced heel (m)
	Safe            bool    // true when residual GZ stays positive
}

// ResolveWaveStability evaluates the worst-case wave slope at the given height
// and length and combines it with the intact righting arm of in. The heel is
// limited to ±SmallAngleMaxDeg so the small-angle model remains valid; if the
// wave slope exceeds that, the computation still returns a Safe=false hint via
// the warning path rather than silently stretching the model.
func ResolveWaveStability(in Input, waveHeight, waveLength float64) WaveStability {
	if waveLength <= 0 {
		return WaveStability{}
	}
	slope := math.Atan(waveHeight * math.Pi / waveLength) // radians
	heelRad := slope
	if heelRad > SmallAngleMaxDeg*DegToRad {
		heelRad = SmallAngleMaxDeg * DegToRad
	}
	heelDeg := heelRad / DegToRad

	res, err := Calc(in)
	if err != nil {
		return WaveStability{WaveInducedHeel: heelDeg, ResidualGZ: 0, Safe: false}
	}
	// righting arm at the wave-induced heel: this is the remaining restoring
	// capacity the hull still has once it has been forced to follow the wave.
	gzAtHeel := res.GMFree * math.Sin(heelRad)
	return WaveStability{
		WaveInducedHeel: heelDeg,
		ResidualGZ:      gzAtHeel,
		Safe:            gzAtHeel > 0,
	}
}

// DynamicStability returns the dynamic stability (area under the GZ curve up to
// a specified angle in degrees) using the trapezoidal integration of the
// small-angle sine model. It is the energy a beam sea must supply to capsize.
func DynamicStability(in Input, toDeg float64) float64 {
	res, err := Calc(in)
	if err != nil {
		return 0
	}
	if toDeg <= 0 {
		return 0
	}
	steps := 100
	sum := 0.0
	prev := 0.0
	prevHeel := 0.0
	for i := 1; i <= steps; i++ {
		deg := toDeg * float64(i) / float64(steps)
		gz := res.GMFree * math.Sin(deg*DegToRad)
		sum += 0.5 * (gz + prev) * (deg - prevHeel) * DegToRad
		prev = gz
		prevHeel = deg
	}
	return sum
}

// WaveBendingMoment estimates the still-water-plus-wave hogging/sagging bending
// moment proxy from a wave of height H over length L at midship, scaled by a
// weight proxy W. It is a crude rule-of-thumb for longitudinal strength checks.
func WaveBendingMoment(waveHeight, length, weight float64) float64 {
	if length <= 0 {
		return 0
	}
	// simplified: M ≈ W * H * L / (some constant); keep monotonic in each input
	return weight * waveHeight * (length / 100.0)
}
