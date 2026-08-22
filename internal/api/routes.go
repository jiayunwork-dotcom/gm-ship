package api

import "net/http"

// routes wires the endpoints and static file serving. The API paths are
// registered first so they take precedence over the catch-all root handler.
//
//	POST /api/gm     -> handleGM
//	POST /api/scan   -> handleScan
//	/example/*       -> static files from ExampleDir (e.g. barge.json)
//	/                -> static web console from WebDir
func (s *Server) routes() {
	s.mux.HandleFunc("/api/gm", s.handleGM)
	s.mux.HandleFunc("/api/scan", s.handleScan)
	s.mux.Handle(
		"/example/",
		http.StripPrefix("/example/", http.FileServer(http.Dir(s.ExampleDir))),
	)
	s.mux.Handle("/", http.FileServer(http.Dir(s.WebDir)))
}
