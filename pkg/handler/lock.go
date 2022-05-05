package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/dchest/safefile"
	"github.com/go-chi/chi/v5"
	"github.com/rs/zerolog/log"
	"github.com/webhippie/terrastate/pkg/config"
	"github.com/webhippie/terrastate/pkg/model"
)

// Lock is used to lock a specific state.
func Lock(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer handleMetrics(time.Now(), "lock", chi.URLParam(req, "*"))

		dir := strings.Replace(
			path.Join(
				cfg.Server.Storage,
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
			log.Error().
				Err(err).
				Msg("Failed to parse body")

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
				log.Error().
					Err(err).
					Str("file", full).
					Msg("Failed to read lock file")

				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)

				return
			}

			if err := json.Unmarshal(file, &existing); err != nil {
				log.Error().
					Err(err).
					Str("file", full).
					Msg("Failed to parse lock file")

				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)

				return
			}

			log.Info().
				Str("existing", existing.ID).
				Str("requested", requested.ID).
				Msg("Lock file already exists")

			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusLocked)

			json.NewEncoder(w).Encode(existing)
			return
		}

		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Error().
				Err(err).
				Str("dir", dir).
				Msg("Failed to create lock dir")

			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)

			return
		}

		marshaled, _ := json.Marshal(requested)

		if err := safefile.WriteFile(full, marshaled, 0644); err != nil {
			log.Error().
				Err(err).
				Str("file", full).
				Msg("Failed to write lock file")

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
