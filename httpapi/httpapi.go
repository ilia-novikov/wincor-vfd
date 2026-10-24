package httpapi

import (
	"net/http"
)

// New builds the HTTP handler for the display REST API. If token is
// non-empty, all /v1/* routes require an "Authorization: Bearer <token>" header
func New(token string) http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", healthzHandler)
	mux.HandleFunc("GET /readyz", readyzHandler)

	mux.HandleFunc("POST /v1/print", printHandler)
	mux.HandleFunc("POST /v1/clear", clearHandler)
	mux.HandleFunc("POST /v1/init", initHandler)
	mux.HandleFunc("PUT /v1/cursor", cursorHandler)
	mux.HandleFunc("POST /v1/backspace", backspaceHandler)
	mux.HandleFunc("POST /v1/linefeed", linefeedHandler)
	mux.HandleFunc("POST /v1/carriage-return", carriageReturnHandler)
	mux.HandleFunc("POST /v1/delete-to-eol", deleteToEOLHandler)

	var h http.Handler = mux
	h = withAuth(h, token)
	h = withLogging(h)
	return h
}
