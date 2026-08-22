package hull

import (
	"math"
	"testing"

	"gm-ship/internal/stability"
)

func TestRectangularBargeIT(t *testing.T) {
	// For a rectangular waterplane, IT = L*B^3 / 12.
	got := RectangularBargeIT(4, 10)
	want := 10 * math.Pow(4, 3) / 12
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("RectangularBargeIT(4,10) = %g, want %g", got, want)
	}
}

func TestRectangularKB(t *testing.T) {
	// KB for a rectangular box at draft d is d/2.
	if got := RectangularKB(2.5); math.Abs(got-1.25) > 1e-9 {
		t.Errorf("RectangularKB(2.5) = %g, want 1.25", got)
	}
}

func TestRectangularKGStable(t *testing.T) {
	// Box with breadth 4, length 10, draft 2.5, target GM 4.55.
	got := RectangularKGStable(4, 10, 2.5, 4.55)
	// Must be a valid (positive) KG and produce the desired GM with the box.
	b := Box{Breadth: 4, Length: 10, Draft: 2.5, KG: got}
	if math.Abs(b.GM()-4.55) > 1e-6 {
		t.Errorf("KGStable produced GM = %g, want 4.55", b.GM())
	}
}

func TestBoxMethods(t *testing.T) {
	b := Box{Breadth: 4, Length: 10, Draft: 2.5, KG: 1.5}
	if math.Abs(b.Volume()-100) > 1e-9 {
		t.Errorf("Volume = %g, want 100", b.Volume())
	}
	if math.Abs(b.KB()-1.25) > 1e-9 {
		t.Errorf("KB = %g, want 1.25", b.KB())
	}
	if math.Abs(b.IT()-53.33333333) > 1e-6 {
		t.Errorf("IT = %g, want 53.333", b.IT())
	}
	in := b.ToInput()
	if err := stability.Validate(in); err != nil {
		t.Errorf("Box.ToInput should be valid: %v", err)
	}
}

func TestEllipseITVersusRectangular(t *testing.T) {
	// Ellipse waterplane: IT = π·B³·L/64 with B the transverse (cubic) axis.
	got := EllipseIT(4, 10) // breadth=4, length=10
	want := math.Pi * math.Pow(4.0, 3) * 10.0 / 64
	if math.Abs(got-want) > 1e-9 {
		t.Errorf("EllipseIT = %g, want %g", got, want)
	}
}

func TestLoadBargeFileMissingErrors(t *testing.T) {
	_, err := LoadBargeFile("/nonexistent/file.json")
	if err == nil {
		t.Errorf("LoadBargeFile should error on missing file")
	}
}
