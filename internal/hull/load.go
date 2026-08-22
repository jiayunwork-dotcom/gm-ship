package hull

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"gm-ship/internal/stability"
)

// BargeSpec is the on-disk representation of the rectangular-barge example. It
// keeps the raw dimensions alongside the pre-computed volume/KB/KG/IT so that
// the example file is self-describing and can be independently checked.
type BargeSpec struct {
	Name    string  `json:"name"`
	Breadth float64 `json:"breadth"`
	Length  float64 `json:"length"`
	Draft   float64 `json:"draft"`
	Volume  float64 `json:"volume"`
	KB      float64 `json:"kb"`
	KG      float64 `json:"kg"`
	IT      float64 `json:"it"`
	Rho     float64 `json:"rho"`
}

// ToInput converts a BargeSpec into a stability.Input. Heel defaults to 0 and
// there is no free-surface effect in the base example; callers may extend it.
func (s BargeSpec) ToInput() stability.Input {
	return stability.Input{
		Volume:  s.Volume,
		KB:      s.KB,
		KG:      s.KG,
		IT:      s.IT,
		HeelDeg: 0,
		Density: s.Rho,
	}
}

// LoadBargeFile reads a barge example JSON file from an explicit path and
// returns the corresponding stability.Input. A missing or malformed file yields
// an error rather than a partial input.
func LoadBargeFile(path string) (stability.Input, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return stability.Input{}, fmt.Errorf("read barge example %q: %w", path, err)
	}
	var spec BargeSpec
	if err := json.Unmarshal(data, &spec); err != nil {
		return stability.Input{}, fmt.Errorf("parse barge example %q: %w", path, err)
	}
	if err := stability.Validate(spec.ToInput()); err != nil {
		return stability.Input{}, fmt.Errorf("barge example %q invalid: %w", path, err)
	}
	return spec.ToInput(), nil
}

// findBargeUp walks up from the current working directory looking for
// example/barge.json. Because the repo root is an ancestor of every package
// directory, this finds the packaged example whether the process is run from
// the repo root (production) or from a package directory (tests).
func findBargeUp() string {
	dir, err := os.Getwd()
	if err != nil {
		return "example/barge.json"
	}
	for {
		cand := filepath.Join(dir, "example", "barge.json")
		if info, statErr := os.Stat(cand); statErr == nil && !info.IsDir() {
			return cand
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "example/barge.json"
}

// FindBargePath returns the first existing candidate path for the barge
// example, searching an explicit example directory (if given) and then walking
// up from the current working directory. This keeps the loader working
// regardless of where the process was started.
func FindBargePath(exampleDir string) string {
	if exampleDir != "" {
		cand := filepath.Join(exampleDir, "barge.json")
		if info, err := os.Stat(cand); err == nil && !info.IsDir() {
			return cand
		}
	}
	return findBargeUp()
}

// LoadBarge loads the packaged rectangular-barge example, searching the usual
// locations when exampleDir is empty.
func LoadBarge(exampleDir string) (stability.Input, error) {
	return LoadBargeFile(FindBargePath(exampleDir))
}
