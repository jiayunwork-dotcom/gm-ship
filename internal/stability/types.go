package stability

// Input carries one intact-stability loading case. All fields use SI units.
//
// The heel angle is given in DEGREES at this boundary. The package converts it
// to radians internally; callers must not pre-convert. Density and free-surface
// inertia are optional and default to sensible values (see DefaultDensity and
// the free-surface helpers).
type Input struct {
	// Volume is the displaced volume ∇ in cubic metres. It must be strictly
	// positive.
	Volume float64

	// KB is the vertical centre of buoyancy above the baseline (m). It may be
	// negative only when BaselineDecl is true (i.e. a different reference was
	// declared, such as a raised baseline).
	KB float64

	// KG is the vertical centre of gravity above the baseline (m). It may be
	// negative only when BaselineDecl is true.
	KG float64

	// IT is the transverse second moment of inertia of the waterplane about the
	// centreline, in m⁴. It must be strictly positive. This is the transverse
	// value; the longitudinal value IL is a different quantity (see
	// LongitudinalMetacentricRadius).
	IT float64

	// HeelDeg is the heel angle in degrees. The small-angle model is valid for
	// |HeelDeg| ≤ SmallAngleMaxDeg; larger values still return the same
	// formula but with a warning attached to the Result.
	HeelDeg float64

	// Density is the mass density of the surrounding water ρ in kg/m³. A value
	// of 0 (the zero value) means "use DefaultDensity". Density only scales the
	// righting moment; it never changes GM or GZ.
	Density float64

	// FreeSurface is the free-surface inertia i of a liquid tank in m⁴. A value
	// of 0 means no free-surface effect. When positive it reduces the
	// effective metacentric height by i/∇.
	FreeSurface float64

	// BaselineDecl, when true, allows KB and/or KG to be negative. It exists
	// because some loading cases are quoted against a raised reference; by
	// default a negative vertical coordinate is treated as a data error.
	BaselineDecl bool
}

// Result holds the computed stability quantities for one Input.
type Result struct {
	// BM is the transverse metacentric radius IT/∇ in metres.
	BM float64

	// GM is the metacentric height KB + BM - KG, before any free-surface
	// correction, in metres.
	GM float64

	// GMFree is the metacentric height after the free-surface correction:
	// GM - FreeSurface/∇, in metres. GZ and the righting moment are derived
	// from GMFree.
	GMFree float64

	// GZ is the righting arm GMFree·sin(φ) in metres.
	GZ float64

	// RightingMoment is Δ·g·GZ in Newtons, with Δ = ρ·∇.
	RightingMoment float64

	// Warning is empty for |φ| ≤ SmallAngleMaxDeg. For larger heals it explains
	// that the small-angle approximation is being stretched. It is not an error.
	Warning string
}

// Point is one sample of the righting-arm curve GZ(φ).
type Point struct {
	// HeelDeg is the heel angle of the sample, in degrees.
	HeelDeg float64 `json:"phi_deg"`

	// GZ is the righting arm at that heel, in metres.
	GZ float64 `json:"gz"`
}
