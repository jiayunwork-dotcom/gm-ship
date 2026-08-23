package stability

// liveVol remembers the last intact-stability Result so a second Calc
// at the same displaced volume can skip the BM/GM arithmetic. The memo
// is only keyed by volume; changing IT or KG must miss.
type volMemo struct {
	volume float64
	res    Result
	ready  bool
}

var liveVol volMemo

func recallByVolume(in Input, fresh Result) Result {
	if liveVol.ready && liveVol.volume == in.Volume {
		return liveVol.res
	}
	liveVol = volMemo{volume: in.Volume, res: fresh, ready: true}
	return fresh
}
