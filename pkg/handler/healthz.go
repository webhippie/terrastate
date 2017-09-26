package handler

import (
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
)

// Healthz is a simple health check used by Docker and Kubernetes.
func Healthz(logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintln(w, http.StatusText(http.StatusOK))
	}
}
