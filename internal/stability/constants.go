package stability

import "math"

// Physical and modelling constants used throughout the stability calculations.
//
// They are exported so that callers (and the web frontend) can show the same
// numbers that the backend used, rather than re-deriving them independently.

const (
	// DefaultDensity is the mass density of seawater in kg/m³, used when a
	// request does not specify a water density. Fresh water is about 1000 and
	// tropical seawater can reach ~1025; the value only scales the righting
	// moment, never GM or GZ.
	DefaultDensity = 1025.0

	// Gravity is the standard acceleration of gravity in m/s² used to turn the
	// righting arm into a righting moment (Newtons).
	Gravity = 9.80665

	// SmallAngleMaxDeg is the conventional limit of the small-angle model.
	// For |φ| ≤ 10° the wall-sided approximation GZ = GM·sin φ is accepted
	// without comment; beyond it the same formula is still returned but the
	// result carries a warning because the linear/small-angle assumption is
	// stretched.
	SmallAngleMaxDeg = 10.0

	// DegToRad converts an angle expressed in degrees to radians. It is kept as
	// a named constant so that the conversion is applied in exactly one place
	// and can never be "forgotten" for one of the call sites.
	DegToRad = math.Pi / 180.0
)
