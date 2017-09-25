package handler

import (
	"io/ioutil"
	"net/http"
	"path"

	"github.com/go-chi/chi"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/webhippie/terrastate/pkg/config"
)

// Unlock is used to unlock a specific state.
func Unlock(logger log.Logger) http.HandlerFunc {
	logger = log.WithPrefix(logger, "handler", "unlock")

	return func(w http.ResponseWriter, req *http.Request) {
		full := path.Join(
			config.Server.Storage,
			chi.URLParam(req, "*"),
			"terraform.tfstate",
		)

		// TODO: handle unlock requests
		body, _ := ioutil.ReadAll(req.Body)
		level.Info(logger).Log(
			"msg", "unlock",
			"file", full,
			"body", body,
		)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte{})
	}
}
