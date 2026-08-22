package stability

// MetacentricRadius returns the transverse metacentric radius BM.
//
//	BM = IT / ∇
//
// BM is the rise of the transverse metacentre above the centre of buoyancy. It
// grows when the waterplane is wide (large IT) and shrinks when the body is
// deeply submerged (large ∇). The argument order is (IT, volume) so that the
// numerator is the waterplane property and the denominator is the volume.
//
// Callers must ensure volume > 0; this function does not validate because it is
// also used internally after Validate has already passed. A non-positive volume
// would return a negative or infinite BM, which is why Validate guards it.
func MetacentricRadius(it, volume float64) float64 {
	return it / volume
}

// MetacentricHeight returns the uncorrected metacentric height GM.
//
//	GM = KB + BM - KG
//
// GM is the distance from the centre of gravity to the transverse metacentre.
// When GM > 0 the body is initially stable (it tends to return to upright);
// when GM < 0 it is initially unstable. The free-surface correction is applied
// separately by MetacentricHeightFree so that callers can inspect both values.
func MetacentricHeight(kb, bm, kg float64) float64 {
	return kb + bm - kg
}

// MetacentricHeightFree returns the metacentric height after the free-surface
// correction:
//
//	GM_free = KB + BM - KG - i/∇
//
// where freeSurface is i and volume is ∇. The correction always reduces GM
// (or leaves it unchanged when i = 0), because a sloshing liquid tank lowers the
// effective centre of gravity.
func MetacentricHeightFree(kb, bm, kg, freeSurface, volume float64) float64 {
	gm := kb + bm - kg
	return gm - freeSurface/volume
}
