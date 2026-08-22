package api

import (
	"net/http"
	"time"
)

// Server bundles the HTTP routing for gm-ship. It serves the static web
// console from WebDir and the example files from ExampleDir, and registers the
// two JSON endpoints. The actual maths live in the stability and hull packages.
type Server struct {
	// ExampleDir is the directory containing the packaged example JSON files
	// (served under /example/ and used as the default load location).
	ExampleDir string

	// WebDir is the directory containing the static web console (index.html
	// and its assets), served from the root path.
	WebDir string

	mux *http.ServeMux
}

// NewServer constructs a Server with sensible defaults: the web console in
// "web" and example files in "example".
func NewServer(exampleDir string) *Server {
	s := &Server{
		ExampleDir: exampleDir,
		WebDir:     "web",
		mux:        http.NewServeMux(),
	}
	s.routes()
	return s
}

// Handler returns the underlying http.Handler (used by tests via httptest).
func (s *Server) Handler() http.Handler {
	return s.mux
}

// ListenAndServe starts the HTTP server on addr. A short read-header timeout is
// set so a stalled client cannot tie up the server indefinitely.
func (s *Server) ListenAndServe(addr string) error {
	srv := &http.Server{
		Addr:              addr,
		Handler:           s.mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	return srv.ListenAndServe()
}
