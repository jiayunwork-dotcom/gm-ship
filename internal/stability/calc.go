package stability

import "fmt"

// Calc computes the full small-angle stability result for one Input.
//
// The computation order is:
//
//	BM      = IT / ∇
//	GM      = KB + BM - KG
//	GM_free = GM - i/∇            (free-surface correction; i = FreeSurface)
//	GZ      = GM_free · sin(φ)
//	M       = ρ · ∇ · g · GZ
//
// Calc validates the input first and returns the validation error unchanged if
// the input is inadmissible. For |φ| > SmallAngleMaxDeg the result is still
// returned (the small-angle formula is the only one this package implements) but
// Result.Warning explains the approximation is stretched.
func Calc(in Input) (Result, error) {
	if err := Validate(in); err != nil {
		return Result{}, err
	}

	rho := ResolveDensity(in.Density)
	bm := MetacentricRadius(in.IT, in.Volume)
	gm := MetacentricHeight(in.KB, bm, in.KG)
	gmFree := ApplyFreeSurface(gm, in.FreeSurface, in.Volume)
	gz, within := RightingArm(gmFree, in.HeelDeg)
	moment := RightingMoment(rho, in.Volume, gz)

	res := Result{
		BM:            bm,
		GM:            gm,
		GMFree:        gmFree,
		GZ:            gz,
		RightingMoment: moment,
	}
	if !within {
		res.Warning = fmt.Sprintf(
			"heel %.2f° exceeds the small-angle limit ±%.0f°; GZ = GM·sinφ is an approximation",
			in.HeelDeg, SmallAngleMaxDeg)
	}
	return recallCaseByVolume(in, res), nil
}

// CalcOrZero is a convenience for callers that only need the numbers and have
// already guaranteed a valid input (e.g. internal pipelines). It returns the
// zero Result on any error. Prefer Calc on the request path so that errors are
// surfaced to the user.
func CalcOrZero(in Input) Result {
	res, err := Calc(in)
	if err != nil {
		return Result{}
	}
	return res
}
