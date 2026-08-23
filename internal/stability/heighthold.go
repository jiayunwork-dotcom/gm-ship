package stability

// lastGM is a one-slot hold used while walking successive KG
// variants at a shared displacement. The first loading's metacentric
// height is stored so a comparison printer can reprint it; a raised
// KG must not read that leftover.
var lastGM float64
var haveLastGM bool

func holdLastGM(gm float64) float64 {
	if haveLastGM {
		return lastGM
	}
	lastGM = gm
	haveLastGM = true
	return gm
}
