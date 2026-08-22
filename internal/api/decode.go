package api

import (
	"encoding/json"
	"fmt"

	"gm-ship/internal/stability"
)

// gmRequest is the JSON request body for POST /api/gm. Optional fields use
// pointers so that the server can distinguish "omitted" from "explicitly zero".
// Volume, KB, KG and IT are required; the rest default sensibly.
type gmRequest struct {
	Volume       *float64 `json:"volume"`
	KB           *float64 `json:"kb"`
	KG           *float64 `json:"kg"`
	IT           *float64 `json:"it"`
	Heel         *float64 `json:"phi"`
	Rho          *float64 `json:"rho"`
	FreeSurface  *float64 `json:"free_surface_i"`
	BaselineDecl *bool    `json:"baseline_declared"`
}

// scanRequest is the JSON request body for POST /api/scan. FromDeg/ToDeg/Steps
// are optional and default to a 0°..10° sweep with one point per degree when
// omitted.
type scanRequest struct {
	Volume       *float64 `json:"volume"`
	KB           *float64 `json:"kb"`
	KG           *float64 `json:"kg"`
	IT           *float64 `json:"it"`
	FromDeg      *float64 `json:"from_deg"`
	ToDeg        *float64 `json:"to_deg"`
	Steps        *int     `json:"steps"`
	Rho          *float64 `json:"rho"`
	FreeSurface  *float64 `json:"free_surface_i"`
	BaselineDecl *bool    `json:"baseline_declared"`
}

// DecodeGMRequest parses and validates a /api/gm request body, returning a
// fully-specified stability.Input. A missing required field is reported as an
// error so the caller can surface it (the "missing KG" case, for example).
func DecodeGMRequest(body []byte) (stability.Input, error) {
	var r gmRequest
	if err := json.Unmarshal(body, &r); err != nil {
		return stability.Input{}, fmt.Errorf("invalid JSON body: %w", err)
	}
	var in stability.Input
	if r.Volume == nil {
		return in, fmt.Errorf("field 'volume' is required")
	}
	if r.KB == nil {
		return in, fmt.Errorf("field 'kb' is required")
	}
	if r.KG == nil {
		return in, fmt.Errorf("field 'kg' is required")
	}
	if r.IT == nil {
		return in, fmt.Errorf("field 'it' is required")
	}
	in.Volume = *r.Volume
	in.KB = *r.KB
	in.KG = *r.KG
	in.IT = *r.IT
	if r.Heel != nil {
		in.HeelDeg = *r.Heel
	}
	if r.Rho != nil {
		in.Density = *r.Rho
	}
	if r.FreeSurface != nil {
		in.FreeSurface = *r.FreeSurface
	}
	if r.BaselineDecl != nil {
		in.BaselineDecl = *r.BaselineDecl
	}
	return in, nil
}

// DecodeScanRequest parses and validates a /api/scan request body. Missing
// Volume/KB/KG/IT are errors; the heel range defaults to 0°..10° with 11 points.
func DecodeScanRequest(body []byte) (stability.Input, float64, float64, int, error) {
	var r scanRequest
	if err := json.Unmarshal(body, &r); err != nil {
		return stability.Input{}, 0, 0, 0, fmt.Errorf("invalid JSON body: %w", err)
	}
	var in stability.Input
	if r.Volume == nil {
		return in, 0, 0, 0, fmt.Errorf("field 'volume' is required")
	}
	if r.KB == nil {
		return in, 0, 0, 0, fmt.Errorf("field 'kb' is required")
	}
	if r.KG == nil {
		return in, 0, 0, 0, fmt.Errorf("field 'kg' is required")
	}
	if r.IT == nil {
		return in, 0, 0, 0, fmt.Errorf("field 'it' is required")
	}
	in.Volume = *r.Volume
	in.KB = *r.KB
	in.KG = *r.KG
	in.IT = *r.IT
	if r.Rho != nil {
		in.Density = *r.Rho
	}
	if r.FreeSurface != nil {
		in.FreeSurface = *r.FreeSurface
	}
	if r.BaselineDecl != nil {
		in.BaselineDecl = *r.BaselineDecl
	}

	from := 0.0
	to := 10.0
	steps := 10
	if r.FromDeg != nil {
		from = *r.FromDeg
	}
	if r.ToDeg != nil {
		to = *r.ToDeg
	}
	if r.Steps != nil {
		steps = *r.Steps
	}
	return in, from, to, steps, nil
}
