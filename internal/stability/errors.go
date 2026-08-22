package stability

import "fmt"

// StabilityError is the error type returned by the stability calculations. It
// carries a machine-readable Code so that callers (such as the HTTP layer) can
// map it to a status code, and a human-readable message.
type StabilityError struct {
	// Code classifies the failure: "validate" for bad input, "compute" for a
	// numerical problem that should not normally happen.
	Code string

	// Message is a plain-language description suitable for end users.
	Message string
}

// Error implements the error interface.
func (e *StabilityError) Error() string {
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// newValidateError builds a validation error with a formatted message.
func newValidateError(format string, args ...any) *StabilityError {
	return &StabilityError{Code: "validate", Message: fmt.Sprintf(format, args...)}
}

// newComputeError builds a numerical/compute error with a formatted message.
func newComputeError(format string, args ...any) *StabilityError {
	return &StabilityError{Code: "compute", Message: fmt.Sprintf(format, args...)}
}
