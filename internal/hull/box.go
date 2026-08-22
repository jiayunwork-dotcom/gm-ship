package hull

import "gm-ship/internal/stability"

// Box describes a rectangular (wall-sided) floating body of constant section,
// such as a barge or a pontoon. It carries the principal dimensions and derives
// the quantities the stability package needs.
type Box struct {
	// Breadth B (m), the transverse extent of the waterplane.
	Breadth float64
	// Length L (m), the longitudinal extent of the waterplane.
	Length float64
	// Draft T (m), the submerged depth from the baseline (keel).
	Draft float64
	// KG (m), the vertical centre of gravity above the baseline.
	KG float64
	// Density rho (kg/m³) of the surrounding water; 0 means seawater default.
	Density float64
}

// Volume returns the displaced volume ∇ = B·L·T.
func (b Box) Volume() float64 {
	return RectangularBargeDisplacement(b.Breadth, b.Length, b.Draft)
}

// KB returns the vertical centre of buoyancy KB = T/2.
func (b Box) KB() float64 {
	return RectangularKB(b.Draft)
}

// IT returns the transverse waterplane inertia IT = B³·L/12.
func (b Box) IT() float64 {
	return RectangularBargeIT(b.Breadth, b.Length)
}

// BM returns the transverse metacentric radius BM = IT/∇.
func (b Box) BM() float64 {
	return b.IT() / b.Volume()
}

// GM returns the uncorrected metacentric height GM = KB + BM - KG.
func (b Box) GM() float64 {
	return b.KB() + b.BM() - b.KG
}

// ToInput builds a stability.Input from the box. Heel defaults to 0 and there is
// no free-surface effect; callers may override those fields afterwards.
func (b Box) ToInput() stability.Input {
	return stability.Input{
		Volume:  b.Volume(),
		KB:      b.KB(),
		KG:      b.KG,
		IT:      b.IT(),
		HeelDeg: 0,
		Density: b.Density,
	}
}
