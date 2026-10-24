package httpapi

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/ilia-novikov/wincor-vfd/display"
)

func parseCharset(r *http.Request) (display.Charset, error) {
	switch c := r.URL.Query().Get("charset"); c {
	case "", "us":
		return display.CharsetUS, nil
	case "ru":
		return display.CharsetRU, nil
	default:
		return 0, fmt.Errorf("unknown charset %q, want %q or %q", c, "us", "ru")
	}
}

func writeDeviceError(w http.ResponseWriter, err error) {
	var notFound *display.DeviceNotFoundError
	var tooLarge *display.PayloadTooLargeError

	switch {
	case errors.As(err, &notFound):
		writeError(w, http.StatusServiceUnavailable, err.Error())
	case errors.As(err, &tooLarge):
		writeError(w, http.StatusBadRequest, err.Error())
	default:
		writeError(w, http.StatusBadGateway, err.Error())
	}
}

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
