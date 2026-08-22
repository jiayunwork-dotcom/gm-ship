package stability

import (
	"testing"
)

func TestSummarize(t *testing.T) {
	in := Input{Volume: 8000, KB: 3.0, KG: 4.0, IT: 12000, Density: DefaultDensity}
	r := Summarize(in)
	if r.Result.GM <= 0 {
		t.Fatalf("GM should be positive, got %v", r.Result.GM)
	}
	if !r.At30Pass {
		t.Fatalf("expected intact criteria to pass")
	}
	if r.Grade() == "unsafe" {
		t.Fatalf("well-founded case should not be unsafe")
	}
}

func TestSummarizeInvalid(t *testing.T) {
	in := Input{Volume: 0, KB: 3.0, KG: 4.0, IT: 12000}
	r := Summarize(in)
	if len(r.Warnings) == 0 {
		t.Fatalf("invalid input should produce a warning")
	}
	if r.Grade() != "unsafe" {
		t.Fatalf("invalid input should grade unsafe")
	}
}

func TestHeelAtCapacity(t *testing.T) {
	in := Input{Volume: 8000, KB: 3.0, KG: 4.0, IT: 12000, Density: DefaultDensity}
	r := Summarize(in)
	h := r.HeelAtCapacity(0.3)
	if h <= 0 || h > 90 {
		t.Fatalf("heel at capacity out of range: %v", h)
	}
}
