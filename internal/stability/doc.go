// Package stability implements the small-angle (initial) intact stability of a
// floating body such as a ship or a barge.
//
// The model is the classical wall-sided, small-heel approximation:
//
//	BM = IT / ∇            (transverse metacentric radius)
//	GM = KB + BM - KG      (metacentric height, before free-surface correction)
//	GZ = GM · sin φ        (righting arm at heel angle φ)
//
// where
//
//	∇  is the displaced volume (m³),
//	KB is the vertical centre of buoyancy above the baseline (m),
//	KG is the vertical centre of gravity above the baseline (m),
//	IT is the transverse second moment of inertia of the waterplane (m⁴),
//	φ  is the heel angle, in degrees, converted to radians inside the maths.
//
// The righting moment is Δ·g·GZ with displacement mass Δ = ρ·∇. A free-surface
// correction subtracts i/∇ from GM, where i is the inertia of a liquid tank.
//
// All lengths are in metres and all angles are expressed in degrees at the
// public API boundary; the package converts degrees to radians internally so
// callers never have to remember the conversion.
package stability
