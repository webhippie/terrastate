package handler

import (
	"net/http"
	"os"
	"path"
	"path/filepath"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"github.com/webhippie/terrastate/pkg/config"
)

// Delete is used to purge a specific state.
func Delete(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer handleMetrics(time.Now(), "delete", chi.URLParam(r, "*"))

		dir := path.Join(
			cfg.Server.Storage,
			filepath.Clean(
				chi.URLParam(r, "*"),
			),
		)

		full := path.Join(
			dir,
			"terraform.tfstate",
		)

		if _, err := os.Stat(full); os.IsNotExist(err) {
			log.Error().
				Str("file", full).
				Msg("State file does not exist")

			http.Error(
				w,
				http.StatusText(http.StatusNotFound),
				http.StatusNotFound,
			)

			return
		}

		if err := os.Remove(full); err != nil {
			log.Error().
				Err(err).
				Str("file", full).
				Msg("Failed to delete state file")

			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)

			return
		}

		log.Info().
			Str("file", full).
			Msg("Successfully deleted state file")

		w.Header().Set("Content-Type", "application/json")
		render.Status(r, http.StatusOK)
	}
}
