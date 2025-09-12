package handler

import (
	"encoding/json"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"github.com/webhippie/terrastate/pkg/config"
	"github.com/webhippie/terrastate/pkg/model"
)

// Unlock is used to unlock a specific state.
func Unlock(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer handleMetrics(time.Now(), "unlock", chi.URLParam(r, "*"))

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

		if err := os.Remove(full); err != nil {
			log.Error().
				Err(err).
				Str("file", full).
				Msg("Failed to delete lock file")

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
			Msg("Successfully unlocked state")

		w.Header().Set("Content-Type", "application/json")
		render.Status(r, http.StatusOK)
	}
}
