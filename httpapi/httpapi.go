package httpapi

import (
	"net/http"
)

// New builds the HTTP handler for the display REST API.
func New() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", healthzHandler)
	mux.HandleFunc("GET /readyz", readyzHandler)

	return mux
}
