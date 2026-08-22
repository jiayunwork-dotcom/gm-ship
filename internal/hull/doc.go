// Package hull provides the geometry of simple floating bodies (barges and
// boxes) and loads the packaged example cases into the stability package's
// Input type.
//
// The hull package is deliberately small: it only turns shape dimensions into
// the waterplane inertia IT and displaced volume ∇ that the stability package
// needs. It does not itself compute BM/GM/GZ — that is stability's job. Keeping
// the geometry separate makes the "double the breadth → IT becomes 8×" and
// "double IT → BM, GM rise by the same increment" rules easy to test in
// isolation.
package hull
