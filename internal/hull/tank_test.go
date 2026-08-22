package hull

import (
	"math"
	"testing"
)

func TestTankVolume(t *testing.T) {
	tk := Tank{Length: 10, Breadth: 4, Height: 2, Fill: 0.5}
	if math.Abs(tk.Volume()-40.0) > 1e-9 {
		t.Fatalf("half-full volume should be 40 m³, got %v", tk.Volume())
	}
	empty := Tank{Length: 10, Breadth: 4, Height: 2, Fill: 0}
	if math.Abs(empty.Volume()-0) > 1e-9 {
		t.Fatalf("empty volume should be 0")
	}
	full := Tank{Length: 10, Breadth: 4, Height: 2, Fill: 1}
	if math.Abs(full.Volume()-80.0) > 1e-9 {
		t.Fatalf("full volume should be 80 m³, got %v", full.Volume())
	}
}

func TestTankWeight(t *testing.T) {
	tk := Tank{Length:           10, Breadth: 4, Height: 2, Fill: 0.5, Density: 1000}
	// 50% fill -> 40 m³ -> 40000 kg.
	if math.Abs(tk.Weight()-40000.0) > 1e-6 {
		t.Fatalf("weight should be 40000, got %v", tk.Weight())
	}
}

func TestFreeSurfaceInertia(t *testing.T) {
	tk := Tank{Length: 10, Breadth: 4, Height: 2}
	i := tk.FreeSurfaceInertia()
	expected := 4.0 * 4.0 * 4.0 * 10.0 / 12.0 // 53.33
	if math.Abs(i-expected) > 1e-9 {
		t.Fatalf("free-surface inertia wrong, got %v", i)
	}
}

func TestFreeSurfaceTerm(t *testing.T) {
	tk := Tank{Length: 10, Breadth: 4, Height: 2}
	term := tk.FreeSurfaceTerm()
	if term <= 0 {
		t.Fatalf("free-surface term should be positive, got %v", term)
	}
	// Larger ship volume should reduce the correction.
	corr := tk.FreeSurfaceCorrection(2000)
	if corr >= tk.FreeSurfaceCorrection(1000) {
		t.Fatalf("larger ship volume should reduce correction")
	}
}

func TestSloshFrequency(t *testing.T) {
	tk := Tank{Length: 10, Breadth: 4, Height: 2, Fill: 0.5}
	f := tk.SloshFrequency()
	if f <= 0 {
		t.Fatalf("slosh frequency should be positive, got %v", f)
	}
	empty := Tank{Length: 10, Breadth: 4, Height: 2, Fill: 0}
	if empty.SloshFrequency() != 0 {
		t.Fatalf("empty tank slosh should be 0")
	}
}
