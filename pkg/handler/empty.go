package handler

import (
	"fmt"
	"net/http"

	"github.com/go-kit/kit/log"
)

// Root just acts as response to the public root.
func Root(logger log.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintln(w, "Terrastate")
	}
}
