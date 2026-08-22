package api

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// errorBody is the JSON shape returned when a request cannot be served. It
// always carries a human-readable "error" field so the web page can show it.
type errorBody struct {
	Error string `json:"error"`
}

// gmBody is the JSON shape returned by POST /api/gm.
type gmBody struct {
	BM            float64 `json:"bm"`
	GM            float64 `json:"gm"`
	GMFree        float64 `json:"gm_free"`
	GZ            float64 `json:"gz"`
	RightingMoment float64 `json:"righting_moment_n"`
	Warning       string  `json:"warning,omitempty"`
}

// scanBody is the JSON shape returned by POST /api/scan.
type scanBody struct {
	Points []pointBody `json:"points"`
}

// pointBody is one sample of the GZ(φ) curve.
type pointBody struct {
	HeelDeg float64 `json:"phi_deg"`
	GZ      float64 `json:"gz"`
}

// writeJSON encodes v as JSON with the given status code.
func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		// The header is already sent; nothing more we can do cleanly.
		fmt.Fprintf(w, "\n")
	}
}

// writeError writes a JSON error body with the given HTTP status.
func writeError(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, errorBody{Error: msg})
}
