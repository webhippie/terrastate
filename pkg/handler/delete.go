package handler

import (
	"net/http"
	"os"
	"path"

	"github.com/Unknwon/com"
	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/webhippie/terrastate/pkg/config"
)

// Delete is used to purge a specific state.
func Delete(logger log.Logger) http.HandlerFunc {
	logger = log.WithPrefix(logger, "handler", "delete")

	return func(w http.ResponseWriter, req *http.Request) {
		full := path.Join(
			config.Server.Storage,
			chi.URLParam(req, "*"),
			"terraform.tfstate",
		)

		if !com.IsFile(full) {
			level.Info(logger).Log(
				"msg", "state file does not exist",
				"file", full,
			)

			http.Error(
				w,
				http.StatusText(http.StatusNotFound),
				http.StatusNotFound,
			)

			return
		}

		err := os.Remove(full)

		if err != nil {
			level.Info(logger).Log(
				"msg", "failed to delete state file",
				"err", err,
			)

			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)

			return
		}

		level.Info(logger).Log(
			"msg", "successfully deleted state file",
			"file", full,
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
