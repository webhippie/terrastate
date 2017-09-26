package handler

import (
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/dchest/safefile"
	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/webhippie/terrastate/pkg/config"
)

// Update is used to update a specific state.
func Update(logger log.Logger) http.HandlerFunc {
	logger = log.WithPrefix(logger, "handler", "update")

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

		// TODO: ID param can be a lock id
		level.Info(logger).Log(
			"msg", "debugging",
			"id", req.URL.Query().Get("ID"),
		)

		content, err := ioutil.ReadAll(req.Body)

		if err != nil {
			level.Info(logger).Log(
				"msg", "failed to load request body",
				"err", err,
			)

			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)

			return
		}

		if err := os.MkdirAll(dir, 0755); err != nil {
			level.Info(logger).Log(
				"msg", "failed to create state dir",
				"dir", dir,
				"err", err,
			)

			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)

			return
		}

		if _, err := os.Stat(full); os.IsNotExist(err) {
			if err := safefile.WriteFile(full, content, 0644); err != nil {
				level.Info(logger).Log(
					"msg", "failed to create state file",
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
				"msg", "successfully created state file",
				"file", full,
			)
		} else {
			if err := safefile.WriteFile(full, content, 0644); err != nil {
				level.Info(logger).Log(
					"msg", "failed to update state file",
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
				"msg", "successfully updated state file",
				"file", full,
			)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
