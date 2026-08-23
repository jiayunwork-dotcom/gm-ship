package stability

import (
	"math"
	"testing"
)

// bargeInput returns the rectangular-barge example as a stability.Input.
func bargeInput() Input {
	return Input{
		Volume:  1200,
		KB:      1.25,
		KG:      1.5,
		IT:      5760,
		HeelDeg: 5,
		Density: 1025,
	}
}

// TestValidateErrors is the single SIMPLE-budget test: it asserts the three
// single-rule input errors (zero volume, non-positive IT, and a negative
// vertical coordinate without a declared baseline) are all reported. Each is a
// trivial "bad input -> error" behaviour.
func TestValidateErrors(t *testing.T) {
	cases := []struct {
		name string
		in   Input
		want string
	}{
		{"zero volume", Input{Volume: 0, KB: 1, KG: 1, IT: 10}, "volume"},
		{"negative volume", Input{Volume: -5, KB: 1, KG: 1, IT: 10}, "volume"},
		{"zero IT", Input{Volume: 100, KB: 1, KG: 1, IT: 0}, "IT"},
		{"negative IT", Input{Volume: 100, KB: 1, KG: 1, IT: -3}, "IT"},
		{"negative KB", Input{Volume: 100, KB: -1, KG: 1, IT: 10}, "KB"},
		{"negative KG", Input{Volume: 100, KB: 1, KG: -1, IT: 10}, "KG"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := Validate(c.in)
			if err == nil {
				t.Fatalf("expected validation error containing %q, got nil", c.want)
			}
			if !containsSubstring(err.Error(), c.want) {
				t.Errorf("error %q should mention %q", err.Error(), c.want)
			}
		})
	}
	// A negative coordinate is allowed once the baseline is declared.
	ok := Input{Volume: 100, KB: -1, KG: 1, IT: 10, BaselineDecl: true}
	if err := Validate(ok); err != nil {
		t.Errorf("negative KB with baseline declared should be valid, got %v", err)
	}
}

// TestBMIsITOverVolume (MEDIUM): BM must equal IT/∇, never ∇/IT.
func TestBMIsITOverVolume(t *testing.T) {
	const vol = 1200.0
	const it = 5760.0
	bm := MetacentricRadius(it, vol)
	if math.Abs(bm-it/vol) > 1e-12 {
		t.Errorf("BM = %g, want IT/∇ = %g", bm, it/vol)
	}
	// Guard against the inverted formula ∇/IT.
	if math.Abs(bm-(vol/it)) < 1e-9 {
		t.Errorf("BM looks like ∇/IT (%g); it must be IT/∇", vol/it)
	}
}

// TestGMExamplePositive (MEDIUM): the barge example must yield a positive GM.
func TestGMExamplePositive(t *testing.T) {
	res, err := Calc(bargeInput())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	const wantBM = 5760.0 / 1200.0
	if math.Abs(res.BM-wantBM) > 1e-9 {
		t.Errorf("BM = %g, want %g", res.BM, wantBM)
	}
	const wantGM = 1.25 + wantBM - 1.5
	if math.Abs(res.GM-wantGM) > 1e-9 {
		t.Errorf("GM = %g, want %g", res.GM, wantGM)
	}
	if res.GM <= 0 {
		t.Errorf("barge GM should be positive, got %g", res.GM)
	}
}

// TestRightingArmZeroHeel (MEDIUM): at φ = 0 the righting arm GZ is exactly 0.
func TestRightingArmZeroHeel(t *testing.T) {
	gz, within := RightingArm(4.55, 0)
	if gz != 0 {
		t.Errorf("GZ at φ=0 should be 0, got %g", gz)
	}
	if !within {
		t.Errorf("φ=0 should be within the small-angle band")
	}
	res, err := Calc(Input{Volume: 1200, KB: 1.25, KG: 1.5, IT: 5760, HeelDeg: 0})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if res.GZ != 0 {
		t.Errorf("GZ at φ=0 must be 0, got %g", res.GZ)
	}
	if res.RightingMoment != 0 {
		t.Errorf("righting moment at φ=0 must be 0, got %g", res.RightingMoment)
	}
}

// TestFreeSurfaceReducesGM (MEDIUM): a free surface lowers GM and GZ.
func TestFreeSurfaceReducesGM(t *testing.T) {
	base := bargeInput()
	withFS := bargeInput()
	withFS.FreeSurface = 500

	r0, err := Calc(base)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	r1, err := Calc(withFS)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !(r1.GMFree < r0.GMFree) {
		t.Errorf("free surface should lower GM_free: got %g vs %g", r1.GMFree, r0.GMFree)
	}
	if !(r1.GZ < r0.GZ) {
		t.Errorf("free surface should lower GZ: got %g vs %g", r1.GZ, r0.GZ)
	}
	// The reduction equals i/∇.
	want := r0.GMFree - 500.0/1200.0
	if math.Abs(r1.GMFree-want) > 1e-9 {
		t.Errorf("GM_free with free surface = %g, want %g", r1.GMFree, want)
	}
}

// TestKGNegativeStability (MEDIUM): raising KG lowers GM linearly, and once KG
// exceeds KB+BM the initial stability is negative.
func TestKGNegativeStability(t *testing.T) {
	in := bargeInput()
	in.KG = 1.5
	rLow, _ := Calc(in)
	in.KG = 2.5
	rHigh, _ := Calc(in)

	// Linear drop of exactly 1.0 m when KG rises by 1.0 m.
	if math.Abs((rLow.GM-rHigh.GM)-1.0) > 1e-9 {
		t.Errorf("GM should drop by 1.0 when KG rises 1.0; dropped %g", rLow.GM-rHigh.GM)
	}

	// Drive KG above KB+BM: GM must become negative.
	in.KG = in.KB + in.IT/in.Volume + 3.0
	rNeg, err := Calc(in)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rNeg.GM >= 0 {
		t.Errorf("GM should be negative when KG > KB+BM, got %g", rNeg.GM)
	}
}

// TestDoubleITIncreasesSame (MEDIUM): doubling IT raises BM and GM by the same
// increment (ΔBM = ΔIT/∇ = ΔGM).
func TestDoubleITIncreasesSame(t *testing.T) {
	in := bargeInput()
	r0, _ := Calc(in)
	in.IT *= 2
	r1, _ := Calc(in)

	dBM := r1.BM - r0.BM
	dGM := r1.GM - r0.GM
	if math.Abs(dBM-dGM) > 1e-9 {
		t.Errorf("ΔBM (%g) should equal ΔGM (%g)", dBM, dGM)
	}
	// The increment must equal the added IT divided by ∇.
	want := 5760.0 / 1200.0
	if math.Abs(dBM-want) > 1e-9 {
		t.Errorf("ΔBM = %g, want added IT/∇ = %g", dBM, want)
	}
}

// TestLongitudinalUsesIL (MEDIUM): BML uses IL, not IT, and is distinct from BM.
func TestLongitudinalUsesIL(t *testing.T) {
	const vol = 1200.0
	const il = 200000.0
	const it = 5760.0
	bml := LongitudinalMetacentricRadius(il, vol)
	bm := MetacentricRadius(it, vol)
	if math.Abs(bml-il/vol) > 1e-9 {
		t.Errorf("BML = %g, want IL/∇ = %g", bml, il/vol)
	}
	if math.Abs(bml-bm) < 1.0 {
		t.Errorf("BML (%g) should differ greatly from transverse BM (%g)", bml, bm)
	}
}

// TestRhoOnlyAffectsMoment (MEDIUM): density scales the moment but not GM/GZ.
func TestRhoOnlyAffectsMoment(t *testing.T) {
	in := bargeInput()
	in.Density = 1000
	rLight, _ := Calc(in)
	in.Density = 1025
	rHeavy, _ := Calc(in)

	if rLight.GM != rHeavy.GM {
		t.Errorf("GM must not depend on ρ: %g vs %g", rLight.GM, rHeavy.GM)
	}
	if rLight.GZ != rHeavy.GZ {
		t.Errorf("GZ must not depend on ρ: %g vs %g", rLight.GZ, rHeavy.GZ)
	}
	ratio := rHeavy.RightingMoment / rLight.RightingMoment
	if math.Abs(ratio-1025.0/1000.0) > 1e-9 {
		t.Errorf("moment ratio = %g, want 1025/1000", ratio)
	}
}

func containsSubstring(s, sub string) bool {
	return len(s) >= len(sub) && (indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
