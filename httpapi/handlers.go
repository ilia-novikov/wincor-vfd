package httpapi

import (
	"encoding/json"
	"net/http"

	"github.com/ilia-novikov/wincor-vfd/display"
)

func healthzHandler(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func readyzHandler(w http.ResponseWriter, _ *http.Request) {
	if err := display.Probe(); err != nil {
		writeError(w, http.StatusServiceUnavailable, err.Error())
		return
	}
	w.WriteHeader(http.StatusOK)
}

func printHandler(w http.ResponseWriter, r *http.Request) {
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

	withDisplay(w, r, func(d *display.Display) error {
		return d.PrintText(req.Text)
	})
}

func clearHandler(w http.ResponseWriter, r *http.Request) {
	withDisplay(w, r, func(d *display.Display) error { return d.ClearDisplay() })
}

func initHandler(w http.ResponseWriter, r *http.Request) {
	withDisplay(w, r, func(d *display.Display) error { return d.Initialize() })
}

func cursorHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		X int `json:"x"`
		Y int `json:"y"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid JSON body")
		return
	}

	withDisplay(w, r, func(d *display.Display) error { return d.MoveCursor(req.X, req.Y) })
}

func backspaceHandler(w http.ResponseWriter, r *http.Request) {
	withDisplay(w, r, func(d *display.Display) error { return d.SendBackspace() })
}

func linefeedHandler(w http.ResponseWriter, r *http.Request) {
	withDisplay(w, r, func(d *display.Display) error { return d.SendLineFeed() })
}

func carriageReturnHandler(w http.ResponseWriter, r *http.Request) {
	withDisplay(w, r, func(d *display.Display) error { return d.SendCarriageReturn() })
}

func deleteToEOLHandler(w http.ResponseWriter, r *http.Request) {
	withDisplay(w, r, func(d *display.Display) error { return d.DeleteToEndOfLine() })
}

// withDisplay opens the display for the request's charset, runs fn, and
// closes it again.
func withDisplay(w http.ResponseWriter, r *http.Request, fn func(*display.Display) error) {
	charset, err := parseCharset(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	d, err := display.Open(charset)
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
