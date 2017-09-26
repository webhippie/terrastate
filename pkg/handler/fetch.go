package handler

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/webhippie/terrastate/pkg/config"
)

// Fetch is used to fetch a specific state.
func Fetch(logger log.Logger) http.HandlerFunc {
	logger = log.WithPrefix(logger, "handler", "fetch")

	return func(w http.ResponseWriter, req *http.Request) {
		full := path.Join(
			config.Server.Storage,
			chi.URLParam(req, "*"),
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

		file, err := ioutil.ReadFile(
			full,
		)

		if err != nil {
			level.Info(logger).Log(
				"msg", "failed to read state file",
				"err", err,
			)

			http.Error(
				w,
				http.StatusText(http.StatusNoContent),
				http.StatusNoContent,
			)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		fmt.Fprintln(w, string(file))
	}
}
