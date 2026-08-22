package api

import (
	"io"
	"net/http"

	"gm-ship/internal/stability"
)

// handleScan implements POST /api/scan.
//
// Request JSON: {volume, kb, kg, it, from_deg?, to_deg?, steps?, rho?,
// free_surface_i?, baseline_declared?}.
// Response JSON: {points: [{phi_deg, gz}, ...]}.
//
// The points are produced entirely by stability.ScanGZ; the web chart plots
// exactly these numbers, so the GZ–φ curve always reflects the real formula and
// is never a hardcoded sine drawn in the browser.
func (s *Server) handleScan(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		writeError(w, http.StatusMethodNotAllowed, "POST is required for /api/scan")
		return
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		writeError(w, http.StatusBadRequest, "could not read request body")
		return
	}
	in, from, to, steps, derr := DecodeScanRequest(body)
	if derr != nil {
		writeError(w, http.StatusBadRequest, derr.Error())
		return
	}
	pts, serr := stability.ScanGZ(in, from, to, steps)
	if serr != nil {
		writeError(w, http.StatusBadRequest, serr.Error())
		return
	}
	out := scanBody{Points: make([]pointBody, 0, len(pts))}
	for _, p := range pts {
		out.Points = append(out.Points, pointBody{HeelDeg: p.HeelDeg, GZ: p.GZ})
	}
	writeJSON(w, http.StatusOK, out)
}

// ScanFromInput mirrors GMFromInput for the scan endpoint: it decodes a raw
// body and returns the encoded scanBody plus an error string (empty on
// success). Tests use it to assert the produced GZ(φ) follows GZ = GM·sin φ.
func ScanFromInput(body []byte) (scanBody, string) {
	in, from, to, steps, derr := DecodeScanRequest(body)
	if derr != nil {
		return scanBody{}, derr.Error()
	}
	pts, serr := stability.ScanGZ(in, from, to, steps)
	if serr != nil {
		return scanBody{}, serr.Error()
	}
	out := scanBody{Points: make([]pointBody, 0, len(pts))}
	for _, p := range pts {
		out.Points = append(out.Points, pointBody{HeelDeg: p.HeelDeg, GZ: p.GZ})
	}
	return out, ""
}
