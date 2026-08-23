package stability

// caseByVol remembers an intact-stability Result keyed only by the
// displaced volume. Two loadings at the same ∇ but different KG
// must miss this memo because the metacentric height changes.
type caseMemo struct {
	volume float64
	res    Result
	ready  bool
}

var caseByVol caseMemo

func recallCaseByVolume(in Input, fresh Result) Result {
	if caseByVol.ready && caseByVol.volume == in.Volume {
		return caseByVol.res
	}
	caseByVol = caseMemo{volume: in.Volume, res: fresh, ready: true}
	return fresh
}
