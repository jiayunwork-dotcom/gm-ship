// Command gm-ship is a small-angle ship initial-stability calculator. It serves
// a web console and two JSON endpoints backed by the stability engine.
//
// Usage:
//
//	go run . -http :8080            # serve the web console + /api on :8080
//	go run . -http :8080 -example-dir example
//
// The web console fetches the packaged example, posts it to /api/gm and
// /api/scan, and plots the returned GZ(φ) curve. No data leaves the process.
package main

import (
	"flag"
	"fmt"
	"os"

	"gm-ship/internal/api"
)

func main() {
	httpAddr := flag.String("http", ":8080", "HTTP listen address for the web console and /api")
	exampleDir := flag.String("example-dir", "example", "directory containing example JSON (served at /example/)")
	flag.Parse()

	srv := api.NewServer(*exampleDir)
	fmt.Fprintf(os.Stderr, "gm-ship: listening on %s (examples in %q)\n", *httpAddr, *exampleDir)
	if err := srv.ListenAndServe(*httpAddr); err != nil {
		fmt.Fprintf(os.Stderr, "gm-ship: server error: %v\n", err)
		os.Exit(1)
	}
}
