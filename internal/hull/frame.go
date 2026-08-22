package hull

// Frame is a transverse section of a hull described by half-breadths at several
// waterlines. It is the classic stations-and-offsets representation used in
// hydrostatic tables. Breadths are measured from the centreline (full, not half).
type Frame struct {
	Station float64   // longitudinal position (m)
	Heights []float64 // vertical positions of offsets (m, ascending)
	Breadths []float64 // full breadths at those heights (m)
}

// Area computes the submerged section area (m²) using the trapezoidal rule on
// the offsets between successive heights. It is the transverse area at one frame.
func (f Frame) Area() float64 {
	n := len(f.Heights)
	if n < 2 || len(f.Breadths) != n {
		return 0
	}
	area := 0.0
	for i := 1; i < n; i++ {
		h0 := f.Heights[i-1]
		h1 := f.Heights[i]
		b0 := f.Breadths[i-1]
		b1 := f.Breadths[i]
		area += 0.5 * (b0 + b1) * (h1 - h0)
	}
	return area
}

// CentroidHeight returns the vertical centroid of the section area above the
// baseline (m), using the trapezoidal centroid formula. Used to estimate the
// centre of buoyancy per station.
func (f Frame) CentroidHeight() float64 {
	n := len(f.Heights)
	if n < 2 || len(f.Breadths) != n {
		return 0
	}
	numer := 0.0
	denom := 0.0
	for i := 1; i < n; i++ {
		h0 := f.Heights[i-1]
		h1 := f.Heights[i]
		b0 := f.Breadths[i-1]
		b1 := f.Breadths[i]
		dh := h1 - h0
		// trapezoid centroid height (average of the two heights)
		hc := (h0 + h1) / 2
		a := 0.5 * (b0 + b1) * dh
		numer += hc * a
		denom += a
	}
	if denom == 0 {
		return 0
	}
	return numer / denom
}

// SectionModulus returns the section modulus (m³) for bending, approximated as
// I/c where I is the second moment of area about the centroid and c is the
// distance to the extreme fibre (max height). A crude but monotonic proxy.
func (f Frame) SectionModulus() float64 {
	n := len(f.Heights)
	if n < 2 || len(f.Breadths) != n {
		return 0
	}
	c := f.Heights[n-1] // topmost offset
	if c == 0 {
		return 0
	}
	yc := f.CentroidHeight()
	i := 0.0
	for i2 := 1; i2 < n; i2++ {
		h0 := f.Heights[i2-1]
		h1 := f.Heights[i2]
		b0 := f.Breadths[i2-1]
		b1 := f.Breadths[i2]
		dh := h1 - h0
		yc0 := (h0 + h1) / 2
		i += (b0 + b1) / 2 * dh * (yc0-yc) * (yc0-yc)
	}
	return i / c
}
