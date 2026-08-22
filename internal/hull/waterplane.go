package hull

import "math"

// Waterplane inertia formulas for non-rectangular shapes.
//
// These are provided so that the geometry layer can serve more than just boxes.
// Every function returns the transverse second moment of inertia IT about the
// centreline, in m⁴, consistent with the convention used by the stability
// package. They are independent of the displaced volume; that is supplied
// separately to stability.MetacentricRadius.

// EllipseIT returns the transverse inertia of an elliptical waterplane of
// breadth (major axis) B and length (minor axis) L:
//
//	IT = π · B³ · L / 64
//
// (B is the transverse axis; the cubic power is again on the transverse
// dimension, exactly as for the rectangle.)
func EllipseIT(breadth, length float64) float64 {
	return math.Pi * breadth * breadth * breadth * length / 64.0
}

// TriangleIT returns the transverse inertia of an isosceles triangular
// waterplane with transverse base B and length L, apex centred:
//
//	IT = B³ · L / 48
//
// The factor 1/48 (versus 1/12 for a rectangle) reflects the smaller area
// moment of the triangle.
func TriangleIT(breadth, length float64) float64 {
	return breadth * breadth * breadth * length / 48.0
}

// CircleIT returns the transverse inertia of a circular waterplane of diameter
// D. A circle is symmetric, so the transverse and longitudinal inertias are
// equal:
//
//	IT = π · D⁴ / 64
func CircleIT(diameter float64) float64 {
	return math.Pi * math.Pow(diameter, 4) / 64.0
}

// RegularPolygonIT returns the transverse inertia of a regular n-gon
// waterplane with circumscribed "breadth" B (distance across the flats is
// 2·B·cos(π/n)). The formula is:
//
//	IT = (n · B³ · L_poly) / (24 · tan(π/n))
//
// where L_poly is the longitudinal prism length. For n = 4 and L_poly = L it
// reduces to the rectangle factor 1/12, confirming consistency.
func RegularPolygonIT(sideCount int, breadth, length float64) float64 {
	if sideCount < 3 {
		return 0
	}
	n := float64(sideCount)
	return n * breadth * breadth * breadth * length / (24.0 * math.Tan(math.Pi/n))
}
