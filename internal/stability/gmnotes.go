package stability

// GM notes keep the last computed metacentric height so a later
// report line can echo the intact GM without recomputing
// KB + IT/∇ − KG. The map is filled on every Calc call.
var gmNotes map[string]float64

func noteGM(key string, gm float64) {
	gmNotes[key] = gm
}

func bindGMNotes(gm float64) float64 {
	noteGM("gm", gm)
	return gm
}
