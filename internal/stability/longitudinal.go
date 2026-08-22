package stability

// Longitudinal stability helpers.
//
// The longitudinal metacentric radius is a DIFFERENT quantity from the
// transverse one and must never be confused with it:
//
//	BML = IL / ∇
//
// where IL is the longitudinal second moment of inertia of the waterplane
// (about a transverse axis through the centre of flotation), not the transverse
// IT. Because a ship is usually much longer than it is wide, IL ≫ IT and
// therefore BML ≫ BM. Mixing the two up would produce a wildly wrong GM if BML
// were ever used where BM belongs. The small-angle righting model here is
// explicitly transverse, so BML is provided for completeness and for
// cross-checks but is NOT fed into GM.

// LongitudinalMetacentricRadius returns the longitudinal metacentric radius:
//
//	BML = IL / ∇
//
// with longitudinalInertia = IL (m⁴) and volume = ∇ (m³). It uses IL, not IT.
func LongitudinalMetacentricRadius(longitudinalInertia, volume float64) float64 {
	return longitudinalInertia / volume
}

// LongitudinalMetacentricHeight returns the longitudinal metacentric height
// about the same baseline convention:
//
//	GML = KB + BML - KG
//
// Provided for symmetry with the transverse MetacentricHeight and for tests
// that assert the two are computed from different inertias.
func LongitudinalMetacentricHeight(kb, bml, kg float64) float64 {
	return kb + bml - kg
}
