package stability

import (
	"math"
	"testing"
)

func TestDamagedGMZeroCase(t *testing.T) {
	base := Result{GM: 1.2, GMFree: 1.1}
	dmg := DamageCase{Volume: 200, Arm: 2.0, Density: 1000}
	newGM, listDeg := DamagedGM(base, dmg, DefaultDensity)
	if newGM <= 0 {
		t.Fatalf("DamagedGM returned non-positive GM: %v", newGM)
	}
	if listDeg < 0 || listDeg >= 90 {
		t.Fatalf("unexpected list angle: %v", listDeg)
	}
	// With a valid arm and weight, the model should produce a finite heel.
	if math.IsInf(listDeg, 1) {
		t.Fatalf("expected finite list angle, got +Inf")
	}
}

func TestDamagedGMNoDamage(t *testing.T) {
	base := Result{GM: 1.0, GMFree: 0.95}
	dmg := DamageCase{Volume: 0, Arm: 0, Density: 1000}
	newGM, _ := DamagedGM(base, dmg, DefaultDensity)
	if math.Abs(newGM-base.GM) > 1e-9 {
		t.Fatalf("no damage should keep GM, got %v", newGM)
	}
}

func TestFloodableLength(t *testing.T) {
	base := Result{GM: 1.0, GMFree: 0.9}
	fl := FloodableLength(base, 20.0, 8.0, 7.0, DefaultDensity)
	if fl <= 0 || fl > 30 {
		t.Fatalf("floodable length out of range: %v", fl)
	}
}

func TestPermissibleHeel(t *testing.T) {
	ok, margin := PermissibleHeel(4.0, 7.0)
	if !ok {
		t.Fatalf("heel within limit should be permissible")
	}
	if margin < 0 {
		t.Fatalf("margin should be positive when within limit")
	}
	ok2, _ := PermissibleHeel(10.0, 7.0)
	if ok2 {
		t.Fatalf("heel beyond limit should be impermissible")
	}
}
