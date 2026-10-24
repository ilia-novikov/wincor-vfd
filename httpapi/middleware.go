package httpapi

import (
	"log"
	"net/http"
	"strings"
	"time"
)

func withAuth(next http.Handler, token string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if token == "" || !strings.HasPrefix(r.URL.Path, "/v1/") {
			next.ServeHTTP(w, r)
			return
		}

		if r.Header.Get("Authorization") != "Bearer "+token {
			writeError(w, http.StatusUnauthorized, "missing or invalid bearer token")
			return
		}

		next.ServeHTTP(w, r)
	})
}

type statusRecorder struct {
	http.ResponseWriter
	status int
}

func (r *statusRecorder) WriteHeader(status int) {
	r.status = status
	r.ResponseWriter.WriteHeader(status)
}

func withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		rec := &statusRecorder{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rec, r)

		log.Printf("%s %s -> %d (%s)", r.Method, r.URL.Path, rec.status, time.Since(start))
	})
}
