package handler

import (
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/webhippie/terrastate/pkg/config"
)

// Delete is used to purge a specific state.
func Delete(logger log.Logger) http.HandlerFunc {
	logger = log.WithPrefix(logger, "handler", "delete")

	return func(w http.ResponseWriter, req *http.Request) {
		dir := strings.Replace(
			path.Join(
				config.Server.Storage,
				chi.URLParam(req, "*"),
			),
			"../", "", -1,
		)

		full := path.Join(
			dir,
			"terraform.tfstate",
		)

		if _, err := os.Stat(full); os.IsNotExist(err) {
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

		if err := os.Remove(full); err != nil {
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
