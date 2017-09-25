package handler

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/unrolled/render"
)

// Root just acts as repsonse to the public root.
func Root(logger log.Logger, r *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		r.Text(w, http.StatusOK, "Terrastate")
	}
}
