package api

import (
	"io"
	"net/http"

	"gm-ship/internal/stability"
)

// handleGM implements POST /api/gm.
//
// Request JSON: {volume, kb, kg, it, phi?, rho?, free_surface_i?,
// baseline_declared?}.
// Response JSON: {bm, gm, gm_free, gz, righting_moment_n, warning?}.
// On any validation or compute error the response is {"error": "..."} with
// status 400.
func (s *Server) handleGM(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST is required for /api/gm")
		return
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "could not read request body")
		return
	}
	in, derr := DecodeGMRequest(body)
	if derr != nil {
		writeError(w, http.StatusBadRequest, derr.Error())
		return
	}
	res, cerr := stability.Calc(in)
	if cerr != nil {
		writeError(w, http.StatusBadRequest, cerr.Error())
		return
	}
	writeJSON(w, http.StatusOK, gmBody{
		BM:             res.BM,
		GM:             res.GM,
		GMFree:         res.GMFree,
		GZ:             res.GZ,
		RightingMoment: res.RightingMoment,
		Warning:        res.Warning,
	})
}

// GMFromInput is a thin helper used by tests and by any future non-HTTP caller:
// it decodes a raw body and returns the encoded gmBody plus an error string
// (empty on success). It isolates the decode→calc→encode pipeline so the HTTP
// handler stays a one-liner.
func GMFromInput(body []byte) (gmBody, string) {
	in, derr := DecodeGMRequest(body)
	if derr != nil {
		return gmBody{}, derr.Error()
	}
	res, cerr := stability.Calc(in)
	if cerr != nil {
		return gmBody{}, cerr.Error()
	}
	return gmBody{
		BM:             res.BM,
		GM:             res.GM,
		GMFree:         res.GMFree,
		GZ:             res.GZ,
		RightingMoment: res.RightingMoment,
		Warning:        res.Warning,
	}, ""
}
