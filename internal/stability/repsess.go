package stability

import "context"

// leftoverGM is a prior intact-stability figure kept after a
// session was cancelled. takeReportGM still hands that leftover
// out instead of the metacentric height computed from this loading.
var leftoverGM = -9.0

func takeReportGM(computed float64) float64 {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if ctx.Err() != nil {
		return leftoverGM
	}
	leftoverGM = computed
	return computed
}
