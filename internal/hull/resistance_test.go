package hull

import (
	"testing"
)

func TestFrictionalResistance(t *testing.T) {
	rf := FrictionalResistance(1025, 5, 200, 50, 1.19e-6)
	if rf <= 0 {
		t.Fatalf("frictional resistance must be positive, got %v", rf)
	}
	if FrictionalResistance(1025, 0, 200, 50, 1.19e-6) != 0 {
		t.Fatalf("zero speed must yield 0")
	}
	if FrictionalResistance(1025, 5, 200, 0, 1.19e-6) != 0 {
		t.Fatalf("zero length must yield 0")
	}
}

func TestResidualResistance(t *testing.T) {
	rr := ResidualResistance(1025, 4, 1000, 50)
	if rr <= 0 {
		t.Fatalf("residual resistance must be positive, got %v", rr)
	}
	if ResidualResistance(1025, 4, 1000, 0) != 0 {
		t.Fatalf("zero length must yield 0")
	}
}

func TestTotalResistance(t *testing.T) {
	rt := TotalResistance(1025, 5, 200, 50, 1000, 1.19e-6)
	if rt <= 0 {
		t.Fatalf("total resistance must be positive, got %v", rt)
	}
}

func TestEffectivePower(t *testing.T) {
	p := EffectivePower(1025, 5, 200, 50, 1000, 1.19e-6)
	if p <= 0 {
		t.Fatalf("effective power must be positive, got %v", p)
	}
	// Power must scale with speed (R·V): doubling speed more than doubles power.
	p2 := EffectivePower(1025, 10, 200, 50, 1000, 1.19e-6)
	if p2 <= p {
		t.Fatalf("higher speed should need more power, got %v vs %v", p2, p)
	}
}
