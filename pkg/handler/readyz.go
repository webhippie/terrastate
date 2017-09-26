package handler

import (
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
)

// Readyz is a simple ready check used by Docker and Kubernetes.
func Readyz(logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintln(w, http.StatusText(http.StatusOK))
	}
}
