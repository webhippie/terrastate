package handler

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/unrolled/render"
)

// Readyz is a simple ready check used by Docker and Kubernetes.
func Readyz(logger log.Logger, r *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		r.Text(w, http.StatusOK, http.StatusText(http.StatusOK))
	}
}
