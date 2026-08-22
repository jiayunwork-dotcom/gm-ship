package hull

import "math"

// Resistance estimates the calm-water powering requirement of a hull using a
// simplified ITTC-style decomposition into frictional and residual (wave-making
// + viscous-pressure) components. All inputs are SI; speed is in m/s.

// FrictionalResistance returns the frictional resistance (N) from the ITTC-1957
// friction line:
//
//	R_f = 0.5 · ρ · V² · S · C_f,   C_f = 0.075 / (log10(Rn) - 2)²
//
// with Reynolds number Rn = V·L/ν and kinematic viscosity ν (default seawater).
func FrictionalResistance(rho, speed, wettedArea, length, viscosity float64) float64 {
	if viscosity <= 0 {
		viscosity = 1.19e-6 // seawater m²/s
	}
	if length <= 0 || speed <= 0 {
		return 0
	}
	rn := speed * length / viscosity
	if rn <= 1 {
		return 0
	}
	cf := 0.075 / math.Pow(math.Log10(rn)-2, 2)
	return 0.5 * rho * speed * speed * wettedArea * cf
}

// ResidualResistance models the residual (wave-making + form) drag with a
// Froude-number bump centred near the design speed. It is a curve-fit proxy, not
// a panel-method result.
func ResidualResistance(rho, speed, displacement, length float64) float64 {
	if length <= 0 {
		return 0
	}
	fn := speed / math.Sqrt(9.80665 * length)
	// simple hull form factor: grows with Fn² and mild peak near Fn≈0.3
	bump := 1.0 + 1.5*fn*fn - 2.0*fn*fn*fn
	if bump < 0.2 {
		bump = 0.2
	}
	return 0.5 * rho * speed * speed * displacement * 0.001 * bump
}

// TotalResistance sums the two components and returns the required thrust (N).
func TotalResistance(rho, speed, wettedArea, length, displacement, viscosity float64) float64 {
	return FrictionalResistance(rho, speed, wettedArea, length, viscosity) +
		ResidualResistance(rho, speed, displacement, length)
}

// EffectivePower returns the effective power (W) required: R · V.
func EffectivePower(rho, speed, wettedArea, length, displacement, viscosity float64) float64 {
	return TotalResistance(rho, speed, wettedArea, length, displacement, viscosity) * speed
}
