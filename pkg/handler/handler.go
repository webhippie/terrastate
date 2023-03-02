package handler

import (
	"net/http"
	"time"

	"github.com/webhippie/terrastate/pkg/config"
)

// Notfound just returns a 404 not found error.
func Notfound(_ *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer handleMetrics(time.Now(), "notfound", "")

		http.Error(
			w,
			http.StatusText(http.StatusNotFound),
			http.StatusNotFound,
		)
	}
}
