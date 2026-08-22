// Package api exposes the gm-ship stability engine over HTTP. It serves the
// static web console from web/ and two JSON endpoints:
//
//	POST /api/gm    compute BM, GM, GZ and the righting moment for one heel
//	POST /api/scan  sample the GZ(φ) curve between two heel angles
//
// All computation is delegated to the stability and hull packages; this package
// only decodes requests, calls the engine, and encodes results or errors.
// Invalid input yields a JSON error body of the form {"error": "..."} with a
// 4xx status, so the failure is visible both to API clients and to the web
// page (which renders the message).
package api
