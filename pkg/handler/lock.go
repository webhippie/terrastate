package handler

import (
	"encoding/json"
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
	"github.com/webhippie/terrastate/pkg/model"
)

// Lock is used to lock a specific state.
func Lock(logger log.Logger) http.HandlerFunc {
	logger = log.WithPrefix(logger, "handler", "lock")

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
			"terraform.lock",
		)

		requested := model.LockInfo{}

		if err := json.NewDecoder(req.Body).Decode(&requested); err != nil {
			level.Info(logger).Log(
				"msg", "failed to parse body",
				"err", err,
			)

			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)

			return
		}

		if _, err := os.Stat(full); !os.IsNotExist(err) {
			existing := model.LockInfo{}

			file, err := ioutil.ReadFile(
				full,
			)

			if err != nil {
				level.Info(logger).Log(
					"msg", "failed to read lock file",
					"file", full,
					"err", err,
				)

				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)

				return
			}

			if err := json.Unmarshal(file, &existing); err != nil {
				level.Info(logger).Log(
					"msg", "failed to parse lock file",
					"file", full,
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
				"msg", "lock file already exists",
				"existing", existing.ID,
				"requested", requested.ID,
			)

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusLocked)

			json.NewEncoder(w).Encode(existing)
			return
		}

		if err := os.MkdirAll(dir, 0755); err != nil {
			level.Info(logger).Log(
				"msg", "failed to create lock dir",
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

		marshaled, _ := json.Marshal(requested)

		if err := safefile.WriteFile(full, marshaled, 0644); err != nil {
			level.Info(logger).Log(
				"msg", "failed to write lock file",
				"file", full,
				"err", err,
			)

			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)

			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
