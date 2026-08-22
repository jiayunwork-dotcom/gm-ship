package stability

import (
	"math"
	"testing"
)

func TestTrimFromDrafts(t *testing.T) {
	tr := TrimFromDrafts(5.0, 0.4)
	if math.Abs(tr.Value-0.4) > 1e-9 {
		t.Fatalf("trim magnitude should be 0.4, got %v", tr.Value)
	}
	if math.Abs(tr.ForwardDraft-5.2) > 1e-9 {
		t.Fatalf("forward draft should be 5.2, got %v", tr.ForwardDraft)
	}
	if math.Abs(tr.AftDraft-4.8) > 1e-9 {
		t.Fatalf("aft draft should be 4.8, got %v", tr.AftDraft)
	}
}

func TestTrimAngleRad(t *testing.T) {
	tr := Trim{Value: 1.0}
	ang := tr.TrimAngleRad(100.0)
	expected := math.Atan(0.01)
	if math.Abs(ang-expected) > 1e-9 {
		t.Fatalf("trim angle wrong, got %v vs %v", ang, expected)
	}
	if tr.TrimAngleRad(0) != 0 {
		t.Fatalf("zero LBP must yield 0")
	}
}

func TestTrimMomentToChange(t *testing.T) {
	m := TrimMomentToChange(8000*DefaultDensity, 200, 0.5, 100)
	if m <= 0 {
		t.Fatalf("trim moment must be positive, got %v", m)
	}
	if TrimMomentToChange(8000*DefaultDensity, 200, 0.5, 0) != 0 {
		t.Fatalf("zero LBP must yield 0")
	}
}

func TestTrimByMoment(t *testing.T) {
	tr := TrimByMoment(8000*DefaultDensity, 200, 8000*DefaultDensity*200*0.5/100, 100)
	if math.Abs(tr-0.5) > 1e-6 {
		t.Fatalf("trim should invert to 0.5, got %v", tr)
	}
}

func TestEffectiveSlope(t *testing.T) {
	tr := Trim{Value: 0.5}
	if math.Abs(tr.EffectiveSlope()-2.0) > 1e-9 {
		t.Fatalf("slope should be 2.0, got %v", tr.EffectiveSlope())
	}
	zero := Trim{Value: 0.0}
	if zero.EffectiveSlope() != 0 {
		t.Fatalf("zero trim slope must be 0")
	}
}
