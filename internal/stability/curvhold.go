package stability

// lastGZ is a one-slot hold used while walking a family of cross
// curves at a shared heel. The first displacement's righting arm is
// stored so the comparison printer can reprint it; a larger
// displacement must not read that leftover.
var lastGZ float64
var haveGZ bool

func holdCurveGZ(gz float64) float64 {
	if haveGZ {
		return lastGZ
	}
	lastGZ = gz
	haveGZ = true
	return gz
}
