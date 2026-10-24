package httpapi

import (
	"encoding/json"
	"net/http"
)

func (s *server) healthzHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func (s *server) readyzHandler(w http.ResponseWriter, _ *http.Request) {
	if err := s.probe(); err != nil {
		writeError(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *server) printHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}
	if req.Text == "" {
		writeError(w, http.StatusBadRequest, "text must not be empty")
		return
	}

	s.withDevice(w, r, func(d Device) error {
		return d.PrintText(req.Text)
	})
}

func (s *server) clearHandler(w http.ResponseWriter, r *http.Request) {
	s.withDevice(w, r, func(d Device) error { return d.ClearDisplay() })
}

func (s *server) initHandler(w http.ResponseWriter, r *http.Request) {
	s.withDevice(w, r, func(d Device) error { return d.Initialize() })
}

func (s *server) cursorHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		X int `json:"x"`
		Y int `json:"y"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	s.withDevice(w, r, func(d Device) error { return d.MoveCursor(req.X, req.Y) })
}

func (s *server) backspaceHandler(w http.ResponseWriter, r *http.Request) {
	s.withDevice(w, r, func(d Device) error { return d.SendBackspace() })
}

func (s *server) linefeedHandler(w http.ResponseWriter, r *http.Request) {
	s.withDevice(w, r, func(d Device) error { return d.SendLineFeed() })
}

func (s *server) carriageReturnHandler(w http.ResponseWriter, r *http.Request) {
	s.withDevice(w, r, func(d Device) error { return d.SendCarriageReturn() })
}

func (s *server) deleteToEOLHandler(w http.ResponseWriter, r *http.Request) {
	s.withDevice(w, r, func(d Device) error { return d.DeleteToEndOfLine() })
}

// withDevice opens the display for the request's charset, runs fn, and
// closes it again.
func (s *server) withDevice(w http.ResponseWriter, r *http.Request, fn func(Device) error) {
	charset, err := parseCharset(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	d, err := s.open(charset)
	if err != nil {
		writeDeviceError(w, err)
		return
	}
	defer d.Close()

	if err := fn(d); err != nil {
		writeDeviceError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
