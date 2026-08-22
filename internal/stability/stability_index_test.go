package stability

import (
	"math"
	"testing"
)

func TestComputeIndexBasic(t *testing.T) {
	in := Input{
		Volume:     8000,
		KB:         3.0,
		KG:         4.0,
		IT:         12000,
		HeelDeg:    0,
		Density:    DefaultDensity,
		FreeSurface: 0,
	}
	idx := ComputeIndex(in)
	if idx.GM <= 0 {
		t.Fatalf("GM should be positive, got %v", idx.GM)
	}
	if idx.MaxGZ <= 0 {
		t.Fatalf("MaxGZ should be positive, got %v", idx.MaxGZ)
	}
	if idx.GZAt30 <= 0 {
		t.Fatalf("GZ@30 should be positive, got %v", idx.GZAt30)
	}
	if idx.Area0to30 <= 0 {
		t.Fatalf("Area0to30 should be positive, got %v", idx.Area0to30)
	}
}

func TestPassesIMO(t *testing.T) {
	in := Input{Volume: 8000, KB: 3.0, KG: 4.0, IT: 12000, Density: DefaultDensity}
	idx := ComputeIndex(in)
	fails := idx.PassesIMO(0.15, 0.1, 45.0)
	if len(fails) != 0 {
		t.Fatalf("expected pass, got failures: %v", fails)
	}
	fails2 := idx.PassesIMO(5.0, 0.1, 35.0)
	if len(fails2) == 0 {
		t.Fatalf("expected GM failure for high threshold")
	}
}

func TestReserveOfStability(t *testing.T) {
	in := Input{Volume: 8000, KB: 3.0, KG: 4.0, IT: 1200, Density: DefaultDensity}
	idx := ComputeIndex(in)
	if idx.ReserveOfStability() != idx.Area0to30 {
		t.Fatalf("ReserveOfStability should equal Area0to30")
	}
}

func TestStabilityIndexMonotonicMax(t *testing.T) {
	in := Input{Volume: 9000, KB: 4.0, KG: 5.0, IT: 15000, Density: DefaultDensity}
	idx := ComputeIndex(in)
	if idx.AngleOfMax < 20 || idx.AngleOfMax > 40 {
		t.Fatalf("angle of max GZ unexpected: %v", idx.AngleOfMax)
	}
	if math.IsNaN(idx.MaxGZ) {
		t.Fatalf("MaxGZ is NaN")
	}
}
