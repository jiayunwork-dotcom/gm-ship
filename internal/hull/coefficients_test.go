package hull

import (
	"math"
	"testing"
)

func TestBlockCoefficient(t *testing.T) {
	// A box exactly fills the block: Cb = 1.
	cb := BlockCoefficient(100, 10, 5, 2)
	if math.Abs(cb-1.0) > 1e-9 {
		t.Fatalf("block coefficient of box should be 1, got %v", cb)
	}
	if BlockCoefficient(10, 0, 5, 2) != 0 {
		t.Fatalf("zero length must yield 0")
	}
}

func TestMidshipCoefficient(t *testing.T) {
	cm := MidshipCoefficient(8, 4, 2)
	if math.Abs(cm-1.0) > 1e-9 {
		t.Fatalf("full rectangular section Cm should be 1, got %v", cm)
	}
}

func TestPrismaticCoefficient(t *testing.T) {
	cp := PrismaticCoefficient(0.7, 0.9)
	if math.Abs(cp-0.7/0.9) > 1e-9 {
		t.Fatalf("prismatic should be cb/cm, got %v", cp)
	}
	if PrismaticCoefficient(0.7, 0) == 0 {
		// zero Cm -> 0, no panic expected
	} else {
		t.Fatalf("zero Cm should yield 0")
	}
}

func TestWaterplaneCoefficient(t *testing.T) {
	cwp := WaterplaneCoefficient(40, 10, 5)
	if math.Abs(cwp-0.8) > 1e-9 {
		t.Fatalf("Cwp should be 0.8, got %v", cwp)
	}
}

func TestCoefficientsFromShape(t *testing.T) {
	c := CoefficientsFromShape(100, 10, 5, 2, 10, 40)
	if c.Block <= 0 || c.Midship <= 0 || c.Waterplane <= 0 {
		t.Fatalf("all coefficients should be positive, got %+v", c)
	}
	if c.Prismatic <= 0 {
		t.Fatalf("prismatic should be positive, got %v", c.Prismatic)
	}
}

func TestLengthDisplacementRatio(t *testing.T) {
	r := LengthDisplacementRatio(100, 1000)
	if r <= 0 {
		t.Fatalf("L/∇^(1/3) should be positive, got %v", r)
	}
}

func TestSpeedLengthRatio(t *testing.T) {
	// V=√(gL) gives ~1.0
	v := SpeedLengthRatio(math.Sqrt(9.80665*100), 100)
	if math.Abs(v-1.0) > 1e-9 {
		t.Fatalf("expected ~1.0, got %v", v)
	}
}
