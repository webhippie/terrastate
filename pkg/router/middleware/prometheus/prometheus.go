package prometheus

import (
	"net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// Handler initializes the prometheus middleware.
func Handler() http.HandlerFunc {
	h := promhttp.Handler()

	return func(w http.ResponseWriter, req *http.Request) {
		h.ServeHTTP(w, req)
	}
}
