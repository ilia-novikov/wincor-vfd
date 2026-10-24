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

	var h http.Handler = mux
	h = withAuth(h, token)
	h = withLogging(h)
	return h
}
