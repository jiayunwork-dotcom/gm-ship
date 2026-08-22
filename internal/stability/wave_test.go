package stability

import (
	"testing"
)

func TestResolveWaveStability(t *testing.T) {
	in := Input{Volume: 8000, KB: 3.0, KG: 4.0, IT: 12000, Density: DefaultDensity}
	ws := ResolveWaveStability(in, 2.0, 100.0)
	if ws.WaveInducedHeel <= 0 {
		t.Fatalf("wave-induced heel must be positive, got %v", ws.WaveInducedHeel)
	}
	if !ws.Safe {
		t.Fatalf("intact GM should leave positive residual GZ in mild wave")
	}
	if ws.ResidualGZ <= 0 {
		t.Fatalf("residual must be positive, got %v", ws.ResidualGZ)
	}
}

func TestResolveWaveStabilityInvalid(t *testing.T) {
	in := Input{Volume: 0, KB: 3.0, KG: 4.0, IT: 12000}
	ws := ResolveWaveStability(in, 2.0, 100.0)
	if ws.Safe {
		t.Fatalf("invalid input should be unsafe")
	}
}

func TestDynamicStability(t *testing.T) {
	in := Input{Volume: 8000, KB: 3.0, KG: 4.0, IT: 12000, Density: DefaultDensity}
	d := DynamicStability(in, 30.0)
	if d <= 0 {
		t.Fatalf("dynamic stability area must be positive, got %v", d)
	}
	if DynamicStability(in, 0) != 0 {
		t.Fatalf("zero range should yield 0")
	}
}

func TestWaveBendingMoment(t *testing.T) {
	m := WaveBendingMoment(3.0, 100, 5000)
	if m <= 0 {
		t.Fatalf("wave bending moment must be positive, got %v", m)
	}
	if WaveBendingMoment(3.0, 0, 5000) != 0 {
		t.Fatalf("zero length must yield 0")
	}
}
