package handler

import (
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/zerolog/log"
	"github.com/webhippie/terrastate/pkg/config"
)

// Delete is used to purge a specific state.
func Delete(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		defer handleMetrics(time.Now(), "delete", chi.URLParam(req, "*"))

		dir := strings.Replace(
			path.Join(
				cfg.Server.Storage,
				chi.URLParam(req, "*"),
			),
			"../", "", -1,
		)

		full := path.Join(
			dir,
			"terraform.tfstate",
		)

		if _, err := os.Stat(full); os.IsNotExist(err) {
			log.Info().
				Str("file", full).
				Msg("state file does not exist")

			http.Error(
				w,
				http.StatusText(http.StatusNotFound),
				http.StatusNotFound,
			)

			return
		}

		if err := os.Remove(full); err != nil {
			log.Info().
				Err(err).
				Str("file", full).
				Msg("failed to delete state file")

			http.Error(
				w,
				http.StatusText(http.StatusInternalServerError),
				http.StatusInternalServerError,
			)

			return
		}

		log.Info().
			Str("file", full).
			Msg("successfully deleted state file")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}
