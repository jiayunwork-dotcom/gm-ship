package stability

// StabilityReport aggregates the intact and damage-relevant descriptors for one
// loading case so a single call can drive a dashboard or a printed summary. It
// deliberately separates the "raw" Calc result from the derived index so the two
// can be audited independently.
type StabilityReport struct {
	Input        Input
	Result       Result
	Index        StabilityIndex
	At30Pass     bool
	Warnings     []string
}

// Summarize builds a StabilityReport for in, running Calc and ComputeIndex and
// applying the standard intact-stability criteria (GM ≥ 0.15 m, GZ@30 ≥ 0.1 m,
// peak GZ ≤ 45°). Validation errors are captured into Warnings rather than
// returned, matching the "report, don't abort" convention of a summary view.
func Summarize(in Input) StabilityReport {
	r := StabilityReport{Input: in}
	res, err := Calc(in)
	if err != nil {
		r.Warnings = append(r.Warnings, "validation: "+err.Error())
		return r
	}
	r.Result = res
	r.Index = ComputeIndex(in)
	fails := r.Index.PassesIMO(0.15, 0.1, 45.0)
	r.At30Pass = len(fails) == 0
	r.Warnings = fails
	return r
}

// Grade returns a one-word classification: "safe", "marginal" or "unsafe" based
// on GM and the intact criteria. It lets a UI badge the report quickly.
func (r StabilityReport) Grade() string {
	if len(r.Warnings) > 0 {
		return "unsafe"
	}
	if r.Result.GM < 0.3 {
		return "marginal"
	}
	return "safe"
}

// HeelAtCapacity returns the heel at which the righting arm equals the supplied
// target (a capacity limit), or 0 if the target is unreachable. It is handy for
// "how far can I heel before losing X" questions.
func (r StabilityReport) HeelAtCapacity(targetGZ float64) float64 {
	return HeelForGZ(r.Result.GMFree, targetGZ)
}
