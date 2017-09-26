package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/webhippie/terrastate/pkg/config"
	"github.com/webhippie/terrastate/pkg/model"
)

// Unlock is used to unlock a specific state.
func Unlock(logger log.Logger) http.HandlerFunc {
	logger = log.WithPrefix(logger, "handler", "unlock")

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

		if err := os.Remove(full); err != nil {
			level.Info(logger).Log(
				"msg", "failed to delete lock file",
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
			"msg", "successfully unlocked state",
			"existing", existing.ID,
			"requested", requested.ID,
		)

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
