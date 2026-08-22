package stability

// Validate checks that an Input describes a physically admissible loading case
// before any quantity is computed.
//
// The rules enforced here are exactly those stated for the small-angle model:
//
//   - Volume (∇) must be > 0. A zero or negative displaced volume is
//     meaningless, so it is rejected.
//   - IT must be > 0. A zero or negative waterplane inertia is rejected.
//   - KB or KG may be negative ONLY when BaselineDecl is true. A negative
//     vertical coordinate without a declared baseline is treated as a sign
//     error in the data and rejected.
//   - Density, when supplied (non-zero), must be positive.
//   - FreeSurface, when supplied (non-zero), must be non-negative.
//
// Validate returns a *StabilityError with Code "validate" on the first
// violation found; otherwise it returns nil.
func Validate(in Input) error {
	if in.Volume <= 0 {
		return newValidateError("displacement volume ∇ must be > 0, got %g", in.Volume)
	}
	if in.IT <= 0 {
		return newValidateError("waterplane transverse inertia IT must be > 0, got %g", in.IT)
	}
	if in.KB < 0 && !in.BaselineDecl {
		return newValidateError("KB is negative (%g) but no baseline was declared", in.KB)
	}
	if in.KG < 0 && !in.BaselineDecl {
		return newValidateError("KG is negative (%g) but no baseline was declared", in.KG)
	}
	if in.Density < 0 {
		return newValidateError("density ρ must be ≥ 0, got %g", in.Density)
	}
	if in.FreeSurface < 0 {
		return newValidateError("free-surface inertia i must be ≥ 0, got %g", in.FreeSurface)
	}
	return nil
}

// ValidateOrPanic is intended for callers that have already guaranteed a valid
// input (for example unit tests on derived quantities). It is NOT used on the
// request path; the HTTP and CLI layers must always use Validate and surface
// the error. It exists only to make the contract explicit where a panic would
// signal a programming error rather than bad user data.
func ValidateOrPanic(in Input) {
	if err := Validate(in); err != nil {
		panic(err)
	}
}
