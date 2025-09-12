package handler

import (
	"io"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/dchest/safefile"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/render"
	"github.com/rs/zerolog/log"
	"github.com/webhippie/terrastate/pkg/config"
	"github.com/webhippie/terrastate/pkg/helper"
)

// Update is used to update a specific state.
func Update(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer handleMetrics(time.Now(), "update", chi.URLParam(r, "*"))

		dir := strings.ReplaceAll(
			path.Join(
				cfg.Server.Storage,
				chi.URLParam(r, "*"),
			),
			"../", "",
		)

		full := path.Join(
			dir,
			"terraform.tfstate",
		)

		content, err := io.ReadAll(r.Body)

		if err != nil {
			log.Error().
				Err(err).
				Msg("Failed to load request body")

			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)

			return
		}

		if err := os.MkdirAll(dir, 0755); err != nil {
			log.Error().
				Err(err).
				Str("dir", dir).
				Msg("Failed to create state dir")

			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)

			return
		}

		if cfg.Encryption.Secret != "" {
			encrypted, err := helper.Encrypt(content, []byte(cfg.Encryption.Secret))

			if err != nil {
				log.Error().
					Err(err).
					Str("file", full).
					Msg("Failed to encrypt the state")

				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)

				return
			}

			content = encrypted
		}

		if _, err := os.Stat(full); os.IsNotExist(err) {
			if err := safefile.WriteFile(full, content, 0644); err != nil {
				log.Error().
					Err(err).
					Str("file", full).
					Msg("Failed to create state file")

				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)

				return
			}

			log.Info().
				Str("file", full).
				Msg("Successfully created state file")
		} else {
			if err := safefile.WriteFile(full, content, 0644); err != nil {
				log.Error().
					Err(err).
					Str("file", full).
					Msg("Failed to update state file")

				http.Error(
					w,
					http.StatusText(http.StatusInternalServerError),
					http.StatusInternalServerError,
				)

				return
			}

			log.Info().
				Str("file", full).
				Msg("Successfully updated state file")
		}

		w.Header().Set("Content-Type", "application/json")
		render.Status(r, http.StatusOK)
	}
}
