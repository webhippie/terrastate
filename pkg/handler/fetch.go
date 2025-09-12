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
	"github.com/webhippie/terrastate/pkg/helper"
)

// Fetch is used to fetch a specific state.
func Fetch(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer handleMetrics(time.Now(), "fetch", chi.URLParam(r, "*"))

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

		file, err := os.ReadFile(
			full,
		)

		if err != nil {
			log.Error().
				Err(err).
				Str("file", full).
				Msg("Failed to read state file")

			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)

			return
		}

		if cfg.Encryption.Secret != "" {
			decrypted, err := helper.Decrypt(file, []byte(cfg.Encryption.Secret))

			if err != nil {
				log.Info().
					Err(err).
					Str("file", full).
					Msg("Failed to decrypt the state")

				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)

				return
			}

			file = decrypted
		}

		render.Status(r, http.StatusOK)
		render.JSON(w, r, file)
	}
}
