package handler

import (
	"io/ioutil"
	"net/http"
	"path"

	"github.com/Unknwon/com"
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
		full := path.Join(
			config.Server.Storage,
			chi.URLParam(req, "*"),
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

		if com.IsFile(full) {
			err := safefile.WriteFile(full, content, 0644)

			if err != nil {
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
				"msg", "sucessfully updated state file",
				"file", full,
			)
		} else {
			err := safefile.WriteFile(full, content, 0644)

			if err != nil {
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
				"msg", "sucessfully created state file",
				"file", full,
			)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
