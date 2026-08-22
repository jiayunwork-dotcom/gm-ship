package stability

import "math"

// RightingArm returns the small-angle righting arm GZ and whether the heel is
// within the small-angle limit.
//
//	GZ = GM_free · sin(φ_rad),   φ_rad = heelDeg · π / 180
//
// The second return value, withinSmallAngle, is true when |heelDeg| ≤
// SmallAngleMaxDeg. It lets callers attach a warning without re-checking the
// limit. The formula always uses the SINE of the angle in radians; it never
// uses the angle itself, and it never treats the degree number as radians.
//
// This is the single place where the degree→radian conversion happens, which is
// deliberate: every other function passes heel angles around in degrees and
// only here are they converted. Keeping the conversion here prevents the classic
// mistake of feeding degrees straight into math.Sin.
func RightingArm(gmFree, heelDeg float64) (gz float64, withinSmallAngle bool) {
	rad := heelDeg * DegToRad
	gz = gmFree * math.Sin(rad)
	return gz, math.Abs(heelDeg) <= SmallAngleMaxDeg
}

// RightingMoment returns the righting moment in Newtons:
//
//	M = Δ · g · GZ = ρ · ∇ · g · GZ
//
// where density is ρ (kg/m³), volume is ∇ (m³), Gravity is g, and gz is the
// righting arm (m). The moment scales linearly with density and with the
// righting arm, but it is independent of GM itself except through GZ.
func RightingMoment(density, volume, gz float64) float64 {
	return density * volume * Gravity * gz
}

// ResolveDensity returns the effective water density, substituting
// DefaultDensity when the supplied value is zero (the zero value meaning
// "unspecified").
func ResolveDensity(density float64) float64 {
	if density == 0 {
		return DefaultDensity
	}
	return density
}
