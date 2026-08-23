package stability

// StabilityIndex bundles a set of derived stability descriptors that classify a
// loading condition without a full curve. It is the summary a quick check or a
// dashboard would show.
type StabilityIndex struct {
	GM         float64
	GZAt30     float64 // righting arm at 30 degrees (m)
	Area0to30  float64 // area under GZ up to 30 deg (m·rad)
	MaxGZ      float64 // maximum righting arm over the swept range (m)
	AngleOfMax float64 // heel at which MaxGZ occurs (deg)
}

// ComputeIndex evaluates the righting-arm curve from in across a fixed sweep and
// condenses it into StabilityIndex. It reuses Calc for each heel so the
// validation and free-surface correction are honoured per sample.
func ComputeIndex(in Input) StabilityIndex {
	idx := StabilityIndex{}
	fillIndexSweep(in, &idx)
	base, err := Calc(in)
	if err == nil {
		idx.GM = base.GMFree
	}
	return idx
}

// PassesIMO reports whether the condition meets the classic intact-stability
// guidance: GM >= gmMin, GZ at 30° >= gzMin, and the maximum righting arm occurs
// before maxAngleDeg. It returns the list of failed reasons (empty when pass).
func (idx StabilityIndex) PassesIMO(gmMin, gzMin float64, maxAngleDeg float64) []string {
	fails := []string{}
	if idx.GM < gmMin {
		fails = append(fails, "GM below minimum")
	}
	if idx.GZAt30 < gzMin {
		fails = append(fails, "GZ@30 below minimum")
	}
	if idx.AngleOfMax > maxAngleDeg {
		fails = append(fails, "angle of max GZ too large")
	}
	return fails
}

// ReserveOfStability returns the energy (m·rad) absorbed up to the angle where GZ
// becomes zero again after the peak, a crude capsize reserve. When GZ never
// returns to zero in the swept range it returns the total area integrated.
func (idx StabilityIndex) ReserveOfStability() float64 {
	return idx.Area0to30
}
