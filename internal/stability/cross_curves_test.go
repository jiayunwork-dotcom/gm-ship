package stability

import (
	"testing"
)

func TestCrossCurveSweep(t *testing.T) {
	samples := CrossCurveSweep(15000, 6.0, []float64{5000, 10000, 15000}, 30.0)
	if len(samples) != 3 {
		t.Fatalf("expected 3 samples, got %d", len(samples))
	}
	for _, s := range samples {
		if s.GZ <= 0 {
			t.Fatalf("GZ should be positive for 30deg, got %v", s.GZ)
		}
	}
	// Larger displacement yields a smaller metacentric radius BM = IT/∇ and thus
	// a smaller GZ at fixed KG; verify the trend reverses.
	if samples[2].GZ >= samples[0].GZ {
		t.Fatalf("larger displacement should reduce GZ, got %v vs %v", samples[2].GZ, samples[0].GZ)
	}
}

func TestRightingArmTable(t *testing.T) {
	pts := RightingArmTable(1.0, []float64{0, 10, 20, 30})
	if len(pts) != 4 {
		t.Fatalf("expected 4 points, got %d", len(pts))
	}
	if pts[0].GZ != 0 {
		t.Fatalf("GZ at 0 deg must be 0, got %v", pts[0].GZ)
	}
	// GZ should increase with heel in the small-angle region.
	if pts[3].GZ <= pts[1].GZ {
		t.Fatalf("GZ should increase with heel, got %v vs %v", pts[3].GZ, pts[1].GZ)
	}
}

func TestMaxRightingArmAngle(t *testing.T) {
	// With a knee angle, the modelled peak sits at the knee.
	a := MaxRightingArmAngle(1.0, 30.0)
	if a != 30 {
		t.Fatalf("expected peak at knee 30, got %v", a)
	}
	// With no knee, the pure sin model peaks at 90.
	a2 := MaxRightingArmAngle(1.0, 0)
	if a2 != 90 {
		t.Fatalf("expected peak at 90 without knee, got %v", a2)
	}
}

func TestHeelForGZ(t *testing.T) {
	h := HeelForGZ(1.0, 0.5)
	if h <= 0 || h > 90 {
		t.Fatalf("heel for target GZ out of range: %v", h)
	}
	h2 := HeelForGZ(0.0, 0.5)
	if h2 != 0 {
		t.Fatalf("zero GM should yield 0, got %v", h2)
	}
}
