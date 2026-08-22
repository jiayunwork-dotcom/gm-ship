package hull

import "math"

// Displacement groups the weight/displacement descriptors of a loaded ship. All
// mass quantities are in kilograms, volumes in cubic metres.
type Displacement struct {
	LightShip   float64 // mass of the empty ship
	Deadweight  float64 // cargo + fuel + stores that can be loaded
	Consumables float64 // fuel/water that will be consumed during the voyage
}

// TotalMass returns the all-up mass: light ship plus everything loadable plus
// the consumables already on board at the start of the voyage.
func (d Displacement) TotalMass() float64 {
	return d.LightShip + d.Deadweight + d.Consumables
}

// DisplacementVolume returns the submerged volume required to float the loaded
// mass, given water density rho (kg/m³).
func (d Displacement) DisplacementVolume(rho float64) float64 {
	if rho <= 0 {
		rho = 1025.0
	}
	if rho == 0 {
		return 0
	}
	return d.TotalMass() / rho
}

// CargoMassToDraft estimates the mean draft produced by a cargo mass using the
// waterplane area A (m²) as a linear approximation near the design draft:
//
//	δT ≈ Δm / (ρ · A)
//
// Useful for quick "how much will I sink" estimates before a full integration.
func (d Displacement) CargoMassToDraft(cargoMass, waterplaneArea, rho float64) float64 {
	if waterplaneArea <= 0 {
		return 0
	}
	if rho <= 0 {
		rho = 1025.0
	}
	return cargoMass / (rho * waterplaneArea)
}

// TonsPerCentimetreImmersion returns the mass (kg) needed to change the mean
// draft by one centimetre, again via the waterplane area:
//
//	TPC = ρ · A · 0.01
//
// It is the longitudinal/daily-loading planning quantity used by cargo officers.
func TonsPerCentimetreImmersion(waterplaneArea, rho float64) float64 {
	if waterplaneArea <= 0 {
		return 0
	}
	if rho <= 0 {
		rho = 1025.0
	}
	return rho * waterplaneArea * 0.01
}

// DeadweightFraction returns DWT / (DWT + light ship), a measure of how much of
// the designed mass is revenue-generating versus the hull itself.
func (d Displacement) DeadweightFraction() float64 {
	total := d.TotalMass()
	if total <= 0 {
		return 0
	}
	return d.Deadweight / total
}

// VolumetricUtilisation returns the fraction of deadweight that the given cargo
// volume represents, using a cargo density (kg/m³). Values near 1 mean the ship
// is volume-limited (could carry more by mass but ran out of hold space).
func (d Displacement) VolumetricUtilisation(cargoVolume, cargoDensity float64) float64 {
	if d.Deadweight <= 0 {
		return 0
	}
	massInHold := cargoVolume * cargoDensity
	if d.Deadweight < massInHold {
		return 1.0
	}
	return massInHold / d.Deadweight
}

// Freeboard is the vertical distance from the waterline to the deck edge. With a
// known draft and depth, it is simply depth - draft; this helper guards against
// negative values (submerged deck) by returning 0 in that case.
func Freeboard(depth, draft float64) float64 {
	fb := depth - draft
	if fb < 0 {
		return 0
	}
	return fb
}

// BlockCoefficientFromDisplacement recovers Cb when only mass, density, L, B, T
// are known, by converting mass to volume first. Convenience wrapper.
func BlockCoefficientFromDisplacement(mass, rho, length, breadth, draft float64) float64 {
	if rho <= 0 {
		rho = 1025.0
	}
	vol := mass / rho
	return BlockCoefficient(vol, length, breadth, draft)
}

// ReserveBuoyancyFraction returns the ratio of the above-water volume to total
// volume, a survivability indicator. With deck submergence it approaches 0.
func ReserveBuoyancyFraction(depth, draft, breadth, length float64) float64 {
	if depth <= 0 || breadth <= 0 || length <= 0 {
		return 0
	}
	total := breadth * length * depth
	if total == 0 {
		return 0
	}
	submerged := breadth * length * draft
	return (total - submerged) / total
}

// HeelMassShift returns the transverse shift of the centre of gravity (m) when a
// mass m is moved a horizontal distance d across the deck. Used for quick
// list estimates before invoking the full stability model.
func HeelMassShift(mass, distance float64) float64 {
	return mass * distance
}

// ListFromShift estimates the list angle (deg) from a transverse weight shift,
// using the metacentric height GM and displacement Δ: tan(θ) = (m·d)/(Δ·GM).
func ListFromShift(mass, distance, displacement, gm float64) float64 {
	if displacement <= 0 || gm <= 0 {
		return 0
	}
	ratio := (mass * distance) / (displacement * gm)
	return math.Atan(ratio) * 180 / math.Pi
}
