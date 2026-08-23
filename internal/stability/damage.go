package stability

import "math"

// DamageCase represents a flooded compartment used in damaged-stability checks.
// Volume is the volume of water added (m^3) and Arm is the horizontal offset of
// its centre of gravity from the centreline (m). A positive Arm shifts the
// gained weight to starboard.
type DamageCase struct {
	Volume   float64
	Arm      float64
	Density  float64 // density of floodwater; 0 means use the ship density
}

// DamagedGM computes the metacentric height after a flooded compartment using
// the lost-buoyancy / added-weight method in its simplest form: the added weight
// V_f * rho_f at the damage arm creates a list, and the effective KG rises. We
// return the new GM (without free surface) and the resulting list angle in
// degrees. The intact GM is taken from r.GM computed by the caller. Both are
// positive magnitudes; a negative newGM signals capsizing risk.
func DamagedGM(base Result, dmg DamageCase, shipDensity float64) (newGM, listDeg float64) {
	rho := dmg.Density
	if rho <= 0 {
		rho = shipDensity
	}
	if rho <= 0 {
		rho = DefaultDensity
	}
	w := dmg.Volume * rho // added weight (mass)
	if w <= 0 {
		return base.GM, 0
	}
	// effective KG shift: simple added-weight approximation
	// new KG_eff = (KG*W + z*W_f)/(W+W_f) but we lack total weight W; approximate
	// using displacement delta = rho_ship * Volume_ship already in base via GZ.
	// We expose the list angle from the heeling moment / restored moment.
	heelMoment := w * dmg.Arm
	// restoring moment per degree ~ delta * g * GZ_slope; slope = GM_free (m/rad)
	delta := base.GM // slope proxy in m of arm per radian of heel
	if delta <= 0 {
		listDeg = math.Inf(1) // unconditionally unstable
		return base.GM, listDeg
	}
	// list angle so that w*Arm = delta*g-reduced moment; angle = atan(heelMoment / (delta))
	listRad := math.Atan2(heelMoment, delta*9.81)
	listDeg = listRad * 180 / math.Pi
	newGM = bindDamageGM(base.GM) // unchanged by this simplified model (free-surface handled separately)
	return newGM, listDeg
}

// FloodableLength returns the maximum compartment length (m) that may be flooded
// without the equilibrium list exceeding a permissive angle maxListDeg, for a
// given compartment breadth and depth. It is a crude damage-control estimate:
// longer compartments heel more because the flood arm grows with length.
func FloodableLength(base Result, breadth, depth, maxListDeg, shipDensity float64) float64 {
	if breadth <= 0 || depth <= 0 || maxListDeg <= 0 {
		return 0
	}
	// derived from list angle relation; invert length from the arm = length/2
	// maxArm = delta*g*tan(maxList)/w_per_m ... w_per_m = breadth*depth*rho
	rho := shipDensity
	if rho <= 0 {
		rho = DefaultDensity
	}
	perM := breadth * depth * rho
	if perM <= 0 {
		return 0
	}
	maxMoment := base.GM * 9.81 * math.Tan(maxListDeg*math.Pi/180)
	if maxMoment <= 0 {
		return 0
	}
	// moment = perM * L * (L/2) ; solve L from perM * L^2/2 = maxMoment
	l2 := 2 * maxMoment / perM
	if l2 < 0 {
		return 0
	}
	return math.Sqrt(l2)
}

// PermissibleHeel returns whether the equilibrium list after a damage case stays
// within the allowed angle, and the margin (allowed - actual) in degrees.
func PermissibleHeel(listDeg, allowedDeg float64) (ok bool, margin float64) {
	margin = allowedDeg - listDeg
	return margin >= 0, margin
}
