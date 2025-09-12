package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/dchest/safefile"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"github.com/webhippie/terrastate/pkg/config"
	"github.com/webhippie/terrastate/pkg/model"
)

// Lock is used to lock a specific state.
func Lock(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer handleMetrics(time.Now(), "lock", chi.URLParam(r, "*"))

		dir := path.Join(
			cfg.Server.Storage,
			filepath.Clean(
				chi.URLParam(r, "*"),
			),
		)

		full := path.Join(
			dir,
			"terraform.lock",
		)

		requested := model.LockInfo{}

		if err := json.NewDecoder(r.Body).Decode(&requested); err != nil {
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

			file, err := os.ReadFile(
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

			render.JSON(w, r, existing)
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
		render.Status(r, http.StatusOK)
	}
}
