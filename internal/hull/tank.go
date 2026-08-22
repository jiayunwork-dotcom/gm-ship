package hull

import "math"

// Tank represents a liquid-filled compartment with a rectangular footprint used
// for free-surface and sloshing estimates. All dimensions are metres.
type Tank struct {
	Length  float64 // along the ship
	Breadth float64 // across the ship
	Height  float64 // vertical
	Fill    float64 // fill fraction 0..1
	Density float64 // liquid density; 0 means use water default
}

// Volume returns the contained液体 volume at the current fill fraction.
func (t Tank) Volume() float64 {
	if t.Fill < 0 {
		return 0
	}
	if t.Fill > 1 {
		t.Fill = 1
	}
	return t.Length * t.Breadth * t.Height * t.Fill
}

// FreeSurfaceInertia returns the transverse free-surface moment of inertia i of
// the liquid surface, B³·L/12 regardless of fill (surface area only depends on
// the footprint). Units m⁴.
func (t Tank) FreeSurfaceInertia() float64 {
	return t.Breadth * t.Breadth * t.Breadth * t.Length / 12.0
}

// SloshFrequency returns an estimate of the first sloshing natural frequency
// (rad/s) using the shallow-water relation ω² = (π/L)^2 · g · h, where h is the
// filled depth. Returns 0 when the tank is empty.
func (t Tank) SloshFrequency() float64 {
	h := t.Height * t.fillClamped()
	if h <= 0 {
		return 0
	}
	return math.Sqrt((math.Pi/math.Max(t.Length, 1e-6)) * (math.Pi/math.Max(t.Length, 1e-6)) * 9.81 * h)
}

func (t Tank) fillClamped() float64 {
	f := t.Fill
	if f < 0 {
		return 0
	}
	if f > 1 {
		return 1
	}
	return f
}

// Weight returns the mass (kg) of the liquid in the tank.
func (t Tank) Weight() float64 {
	rho := t.Density
	if rho <= 0 {
		rho = 1000
	}
	return t.Volume() * rho
}

// FreeSurfaceCorrection returns the GM reduction i/∇ for a given displacement
// volume (m³) of the ship. A larger ship volume diminishes the penalty.
func (t Tank) FreeSurfaceCorrection(shipVolume float64) float64 {
	if shipVolume <= 0 {
		return 0
	}
	return t.FreeSurfaceInertia() / shipVolume
}

// FreeSurfaceTerm returns the free-surface inertia of the tank, exposed for
// testing and for callers that build a combined correction across several tanks.
func (t Tank) FreeSurfaceTerm() float64 {
	return t.FreeSurfaceInertia()
}
