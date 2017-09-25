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

// Lock is used to lock a specific state.
func Lock(logger log.Logger) http.HandlerFunc {
	logger = log.WithPrefix(logger, "handler", "lock")

	return func(w http.ResponseWriter, req *http.Request) {
		full := path.Join(
			config.Server.Storage,
			chi.URLParam(req, "*"),
			"terraform.tfstate",
		)

		// TODO: handle lock requests
		body, _ := ioutil.ReadAll(req.Body)
		level.Info(logger).Log(
			"msg", "lock",
			"file", full,
			"body", body,
		)

		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte{})
	}
}
