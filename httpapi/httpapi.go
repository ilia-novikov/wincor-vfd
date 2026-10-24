// Package httpapi exposes the display over a REST API. Each request opens
// and closes the underlying HID device itself; the server never holds the
// handle open between requests.
package httpapi

import (
	"net/http"

	"github.com/ilia-novikov/wincor-vfd/display"
)

type server struct {
	open  OpenFunc
	probe ProbeFunc
	token string
}

// New builds the HTTP handler for the display REST API. If token is
// non-empty, all /v1/* routes require an "Authorization: Bearer <token>" header
func New(token string) http.Handler {
	return newServer(token, defaultOpener, display.Probe)
}

func defaultOpener(charset display.Charset) (Device, error) {
	return display.Open(charset)
}

func newServer(token string, open OpenFunc, probe ProbeFunc) http.Handler {
	s := &server{open: open, probe: probe, token: token}

	mux := http.NewServeMux()

	mux.HandleFunc("GET /healthz", s.healthzHandler)
	mux.HandleFunc("GET /readyz", s.readyzHandler)

	mux.HandleFunc("POST /v1/print", s.printHandler)
	mux.HandleFunc("POST /v1/clear", s.clearHandler)
	mux.HandleFunc("POST /v1/init", s.initHandler)
	mux.HandleFunc("PUT /v1/cursor", s.cursorHandler)
	mux.HandleFunc("POST /v1/backspace", s.backspaceHandler)
	mux.HandleFunc("POST /v1/linefeed", s.linefeedHandler)
	mux.HandleFunc("POST /v1/carriage-return", s.carriageReturnHandler)
	mux.HandleFunc("POST /v1/delete-to-eol", s.deleteToEOLHandler)

	var h http.Handler = mux
	h = withAuth(h, token)
	h = withLogging(h)
	return h
}
