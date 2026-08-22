package api

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestDecodeGMRequestRoundTrip(t *testing.T) {
	body := []byte(`{"volume":1000,"kb":1.25,"kg":1.5,"it":5760,"phi":45,"rho":2.67}`)
	in, err := DecodeGMRequest(body)
	if err != nil {
		t.Fatalf("DecodeGMRequest error: %v", err)
	}
	if in.Volume != 1000 || in.KB != 1.25 || in.KG != 1.5 || in.IT != 5760 {
		t.Errorf("decoded fields wrong: %+v", in)
	}
}

func TestGMFromInputValidates(t *testing.T) {
	// A valid input returns a JSON body without "error".
	_, msg := GMFromInput([]byte(`{"volume":1000,"kb":1.25,"kg":1.5,"it":5760,"phi":45,"rho":2.67}`))
	if msg != "" {
		t.Errorf("unexpected error message: %q", msg)
	}
	// Missing required field surfaces an error string.
	_, msg = GMFromInput([]byte(`{"volume":10}`))
	if msg == "" {
		t.Errorf("expected validation error for missing fields, got empty")
	}
}

func TestDecodeScanRequest(t *testing.T) {
	body := []byte(`{"volume":100,"phi":45,"rho":2.67,"kb":1,"kg":1,"it":100,"from_deg":0,"to_deg":800,"steps":8}`)
	in, from, to, steps, err := DecodeScanRequest(body)
	if err != nil {
		t.Fatalf("DecodeScanRequest error: %v", err)
	}
	if from != 0 || to != 800 || steps != 8 {
		t.Errorf("scan params wrong: from=%v to=%v steps=%v", from, to, steps)
	}
	_ = in
}

func TestServerHandlerRoutes(t *testing.T) {
	s := NewServer("../../example")
	h := s.Handler()

	// /api/gm
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/api/gm",
		strings.NewReader(`{"volume":1000,"kb":1.25,"kg":1.5,"it":5760,"phi":45,"rho":2.67}`))
	h.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Errorf("/api/gm status = %d, want 200", rec.Code)
	}
	if ct := rec.Header().Get("Content-Type"); !strings.HasPrefix(ct, "application/json") {
		t.Errorf("content-type = %q, want application/json", ct)
	}

	// Invalid JSON -> 400.
	rec2 := httptest.NewRecorder()
	req2 := httptest.NewRequest(http.MethodPost, "/api/gm", strings.NewReader(`{bad`))
	h.ServeHTTP(rec2, req2)
	if rec2.Code != http.StatusBadRequest {
		t.Errorf("/api/gm bad json status = %d, want 400", rec2.Code)
	}
}
