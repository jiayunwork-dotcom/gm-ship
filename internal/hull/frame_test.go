package hull

import (
	"math"
	"testing"
)

func TestFrameArea(t *testing.T) {
	f := Frame{
		Station:  0,
		Heights:  []float64{0, 2, 4},
		Breadths: []float64{0, 6, 10},
	}
	// trapezoids: (0+6)/2*2 + (6+10)/2*2 = 6 + 16 = 22
	if math.Abs(f.Area()-22.0) > 1e-9 {
		t.Fatalf("area should be 22, got %v", f.Area())
	}
}

func TestFrameCentroidHeight(t *testing.T) {
	f := Frame{
		Station:  0,
		Heights:  []float64{0, 2, 4},
		Breadths: []float64{0, 6, 10},
	}
	ch := f.CentroidHeight()
	if ch <= 0 || ch >= 4 {
		t.Fatalf("centroid height out of range: %v", ch)
	}
}

func TestFrameSectionModulus(t *testing.T) {
	f := Frame{
		Station:  0,
		Heights:  []float64{0, 2, 4},
		Breadths: []float64{0, 6, 10},
	}
	sm := f.SectionModulus()
	if sm <= 0 {
		t.Fatalf("section modulus should be positive, got %v", sm)
	}
}

func TestFrameInvalid(t *testing.T) {
	f := Frame{Station: 0, Heights: []float64{1}, Breadths: []float64{1}}
	if f.Area() != 0 {
		t.Fatalf("single-offset frame area must be 0")
	}
	if f.CentroidHeight() != 0 {
		t.Fatalf("single-offset centroid must be 0")
	}
}
