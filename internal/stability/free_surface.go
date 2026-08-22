package stability

// Free-surface correction helpers.
//
// A liquid tank that is not fully pressed (e.g. a partially filled fuel or
// ballast tank) behaves like a pendulum: as the ship heels, the liquid shifts
// to the low side and effectively raises the centre of gravity. The standard
// correction subtracts i/∇ from GM, where i is the second moment of inertia of
// the free surface of the liquid about its own centreline.

// FreeSurfaceCorrection returns the reduction of GM caused by a free surface:
//
//	δGM = i / ∇
//
// with freeSurface = i (m⁴) and volume = ∇ (m³). It is always non-negative for
// a non-negative i, so the correction never increases GM.
func FreeSurfaceCorrection(freeSurface, volume float64) float64 {
	return freeSurface / volume
}

// ApplyFreeSurface subtracts the free-surface correction from an uncorrected
// metacentric height:
//
//	GM_free = GM - i/∇
//
// It is a thin, named wrapper so the intent is obvious at call sites and the
// relationship to MetacentricHeightFree is explicit.
func ApplyFreeSurface(gm, freeSurface, volume float64) float64 {
	return gm - FreeSurfaceCorrection(freeSurface, volume)
}

// FreeSurfaceReducesGM reports whether a given free-surface inertia actually
// lowers the metacentric height. It is true whenever i > 0; it exists mainly so
// that the web UI can label the correction and so that tests can assert the
// direction of the effect without duplicating the formula.
func FreeSurfaceReducesGM(freeSurface float64) bool {
	return freeSurface > 0
}
