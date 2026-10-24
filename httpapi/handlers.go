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

func writeError(w http.ResponseWriter, status int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
