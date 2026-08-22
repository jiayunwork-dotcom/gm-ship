package hull

import "math"

// Coefficients groups the classical hull-form coefficients that relate a real
// ship's volume and area to those of a rectangular block of the same length,
// breadth and draft. They are dimensionless shape descriptors used in early
// resistance, stability and powering estimates.
type Coefficients struct {
	Block      float64 // block (coefficient of fineness) Cb
	Midship    float64 // midship section coefficient Cm
	Prismatic  float64 // prismatic coefficient Cp (Cp = Cb/Cm)
	Waterplane float64 // waterplane area coefficient Cwp
}

// BlockCoefficient returns Cb = displaced volume / (L·B·T). A value near 1
// means a "full" form, near 0.5 a slender form.
func BlockCoefficient(volume, length, breadth, draft float64) float64 {
	if length <= 0 || breadth <= 0 || draft <= 0 {
		return 0
	}
	return volume / (length * breadth * draft)
}

// MidshipCoefficient returns Cm = midship area / (B·T) for the given midship
// section area (m²).
func MidshipCoefficient(midshipArea, breadth, draft float64) float64 {
	if breadth <= 0 || draft <= 0 {
		return 0
	}
	return midshipArea / (breadth * draft)
}

// PrismaticCoefficient returns Cp = Cb / Cm. It is derived, never stored.
func PrismaticCoefficient(cb, cm float64) float64 {
	if cm == 0 {
		return 0
	}
	return cb / cm
}

// WaterplaneCoefficient returns Cwp = waterplane area / (L·B).
func WaterplaneCoefficient(waterArea, length, breadth float64) float64 {
	if length <= 0 || breadth <= 0 {
		return 0
	}
	return waterArea / (length * breadth)
}

// CoefficientsFromShape derives the full set from the basic geometric inputs.
func CoefficientsFromShape(volume, length, breadth, draft, midshipArea, waterArea float64) Coefficients {
	return Coefficients{
		Block:      BlockCoefficient(volume, length,  breadth, draft),
		Midship:    MidshipCoefficient(midshipArea, breadth, draft),
		Prismatic:  PrismaticCoefficient(BlockCoefficient(volume, length, breadth, draft), MidshipCoefficient(midshipArea, breadth, draft)),
		Waterplane: WaterplaneCoefficient(waterArea, length, breadth),
	}
}

// LengthDisplacementRatio returns L / ∛∇, a slenderness indicator. Large
// values denote fine (fast) forms; small values denote blunt (slow) forms.
func LengthDisplacementRatio(length, volume float64) float64 {
	if volume <= 0 {
		return 0
	}
	return length / math.Cbrt(volume)
}

// SpeedLengthRatio returns the Froude number-like V / √(g·L) for a given speed.
func SpeedLengthRatio(speed, length float64) float64 {
	if length <= 0 {
		return 0
	}
	return speed / math.Sqrt(9.80665 * length)
}
