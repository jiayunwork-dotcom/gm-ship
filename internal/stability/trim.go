package stability

import "math"

// Trim describes the longitudinal inclination of a vessel: the difference between
// the forward and aft drafts. A positive trim means the bow is lower (by convention
// used here); the magnitude is what matters for longitudinal stability checks.
type Trim struct {
	Value        float64 // trim magnitude (m)
	ForwardDraft float64 // draft at the forward perpendicular (m)
	AftDraft     float64 // draft at the aft perpendicular (m)
}

// TrimFromDrafts returns the trim magnitude and the two end drafts given a mean
// draft and a signed trim (bow-down positive). The two end drafts split the mean
// symmetrically about amidships.
func TrimFromDrafts(meanDraft, trimSigned float64) Trim {
	return Trim{
		Value:        math.Abs(trimSigned),
		ForwardDraft: meanDraft + trimSigned/2,
		AftDraft:     meanDraft - trimSigned/2,
	}
}

// TrimAngleRad returns the longitudinal inclination angle in radians, computed
// from the trim and the length between perpendiculars LBP.
func (t Trim) TrimAngleRad(lengthBetweenPerpendiculars float64) float64 {
	if lengthBetweenPerpendiculars <= 0 {
		return 0
	}
	return math.Atan(t.Value / lengthBetweenPerpendiculars)
}

// TrimMomentToChange converts a desired change of trim (m) into the longitudinal
// trimming moment needed, using the long form: moment = (Δ·GML·δtrim)/LBP, with
// Δ = ρ·∇. It mirrors the transverse righting relation but along the ship axis.
func TrimMomentToChange(displacement, gml, deltaTrim, lengthBetweenPerpendiculars float64) float64 {
	if lengthBetweenPerpendiculars <= 0 {
		return 0
	}
	return (displacement * gml * deltaTrim) / lengthBetweenPerpendiculars
}

// TrimByMoment returns the resulting trim (m) when a longitudinal moment is
// applied to a vessel with displacement Δ and longitudinal metacentric height
// GML over length LBP. It is the inverse of TrimMomentToChange.
func TrimByMoment(displacement, gml, moment, lengthBetweenPerpendiculars float64) float64 {
	if displacement <= 0 || gml <= 0 || lengthBetweenPerpendiculars <= 0 {
		return 0
	}
	return (moment * lengthBetweenPerpendiculars) / (displacement * gml)
}

// MomentToShiftPerUnitTrimmingArm returns how much longitudinal weight shift is
// needed per metre of trim given the LCG change, i.e. the "moment to change
// trim 1 cm" proxy expressed in metres of trim per (kg·m).
func (t Trim) EffectiveSlope() float64 {
	if t.Value == 0 {
		return 0
	}
	return 1.0 / t.Value
}
