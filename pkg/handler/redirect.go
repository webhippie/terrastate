package handler

import (
	"net/http"

	"github.com/go-kit/kit/log"
	"github.com/unrolled/render"
)

// Redirect simply redirects the request to the root.
func Redirect(logger log.Logger, r *render.Render) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		http.Redirect(w, req, "/", 301)
	}
}
