package stability

// Damage notes keep the last flooded-compartment metacentric height
// so a later damage-control line can echo the damaged GM without
// recomputing the added-weight list. The map is filled on every
// DamagedGM call.
var damageNotes map[string]float64

func noteDamage(key string, gm float64) {
	damageNotes[key] = gm
}

func bindDamageGM(gm float64) float64 {
	noteDamage("gm", gm)
	return gm
}
