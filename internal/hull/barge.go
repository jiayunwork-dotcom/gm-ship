package hull

// RectangularBargeIT returns the transverse second moment of inertia of a
// rectangular waterplane about its centreline:
//
//	IT = B³ · L / 12
//
// with breadth B and length L. The cubic dependence on B is why making a barge
// wider is so effective at raising BM: doubling B multiplies IT (and therefore
// BM) by 8, all else equal.
func RectangularBargeIT(breadth, length float64) float64 {
	return breadth * breadth * breadth * length / 12.0
}

// RectangularBargeDisplacement returns the displaced volume of a rectangular
// box barge of constant section:
//
//	∇ = B · L · T
//
// with breadth B, length L and draft T.
func RectangularBargeDisplacement(breadth, length, draft float64) float64 {
	return breadth * length * draft
}

// RectangularKB returns the vertical centre of buoyancy of a rectangular box
// barge, measured from the keel (baseline):
//
//	KB = T / 2
//
// because the submerged volume is a uniform prism whose centroid sits at half
// the draft.
func RectangularKB(draft float64) float64 {
	return draft / 2.0
}

// RectangularKGStable is a convenience that returns a typical, positively
// stable KG for a box barge given its draft and a desired GM. It is used by the
// example generator; real loading cases supply KG directly. It inverts
// GM = KB + BM - KG for KG:
//
//	KG = KB + BM - GM
//
// with BM = IT/∇.
func RectangularKGStable(breadth, length, draft, wantGM float64) float64 {
	vol := RectangularBargeDisplacement(breadth, length, draft)
	it := RectangularBargeIT(breadth, length)
	bm := it / vol
	kb := RectangularKB(draft)
	return kb + bm - wantGM
}
