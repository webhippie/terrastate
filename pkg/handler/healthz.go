package handler

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/unrolled/render"
)

// Healthz is a simple health check used by Docker and Kubernetes.
func Healthz(logger log.Logger, r *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		r.Text(w, http.StatusOK, http.StatusText(http.StatusOK))
	}
}
