package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Handler initializes the prometheus middleware.
func Handler(token string) http.HandlerFunc {
	h := promhttp.Handler()

	return func(w http.ResponseWriter, r *http.Request) {
		if token == "" {
			h.ServeHTTP(w, r)
			return
		}

		header := r.Header.Get("Authorization")

		if header == "" {
			http.Error(w, "Invalid or missing token", http.StatusUnauthorized)
			return
		}

		if header != "Bearer "+token {
			http.Error(w, "Invalid or missing token", http.StatusUnauthorized)
			return
		}

		h.ServeHTTP(w, r)
	}
}
