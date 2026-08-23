package stability

// ScanGZ samples the righting-arm curve GZ(φ) between two heel angles.
//
// The returned points are evenly spaced from fromDeg to toDeg inclusive, with
// exactly steps intervals (so steps+1 points). Each point is computed with the
// same small-angle model as Calc:
//
//	GZ(φ) = (KB + IT/∇ - KG - i/∇) · sin(φ_rad)
//
// ScanGZ validates the loading case first, so an inadmissible input yields an
// error rather than a misleading curve. steps must be ≥ 1; a non-positive step
// count is a programming error and returns a compute error.
//
// The scan is the backend source of truth for the web chart: the SVG plots
// exactly these returned points, so the curve always reflects the real
// formula rather than a hardcoded sine drawn in the browser.
func ScanGZ(in Input, fromDeg, toDeg float64, steps int) ([]Point, error) {
	if err := Validate(in); err != nil {
		return nil, err
	}
	if steps <= 0 {
		return nil, newComputeError("scan step count must be ≥ 1, got %d", steps)
	}

	gmFree := MetacentricHeightFree(in.KB, MetacentricRadius(in.IT, in.Volume), in.KG, in.FreeSurface, in.Volume)

	heels := make([]float64, 0, steps+1)
	for i := 0; i <= steps; i++ {
		h := fromDeg + (toDeg-fromDeg)*float64(i)/float64(steps)
		heels = append(heels, h)
	}
	return fillArmTable(gmFree, heels), nil
}

// ScanGZDefault samples 0°..maxDeg with a sensible default resolution. It is a
// small convenience used by the example loader and by quick experiments; the
// explicit ScanGZ is preferred on the request path so the caller controls the
// range and resolution.
func ScanGZDefault(in Input, maxDeg float64) ([]Point, error) {
	return ScanGZ(in, 0, maxDeg, int(maxDeg)+1)
}
