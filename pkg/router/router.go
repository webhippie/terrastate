package router

import (
	"net/http"
	"strings"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-kit/kit/log"
	"github.com/unrolled/render"
	"github.com/webhippie/terrastate/pkg/config"
	"github.com/webhippie/terrastate/pkg/handler"
	"github.com/webhippie/terrastate/pkg/router/middleware/basicauth"
	"github.com/webhippie/terrastate/pkg/router/middleware/header"
	"github.com/webhippie/terrastate/pkg/router/middleware/prometheus"
	"github.com/webhippie/terrastate/pkg/router/middleware/requests"
)

// Load initializes the routing of the application.
func Load(logger log.Logger) http.Handler {
	r := render.New(render.Options{
		IsDevelopment: strings.ToLower(config.LogLevel) == "debug",
	})

	mux := chi.NewRouter()

	mux.Use(requests.Requests(logger))

	mux.Use(middleware.Timeout(60 * time.Second))
	mux.Use(middleware.RealIP)

	mux.Use(header.Version)
	mux.Use(header.Cache)
	mux.Use(header.Secure)
	mux.Use(header.Options)

	mux.NotFound(handler.Redirect(logger, r))

	mux.Route("/", func(root chi.Router) {
		if config.Server.Prometheus {
			root.Get("/metrics", prometheus.Handler())
		}

		if config.Server.Pprof {
			root.Mount("/debug", middleware.Profiler())
		}

		root.Get("/", handler.Root(logger, r))
		root.Get("/healthz", handler.Healthz(logger, r))
		root.Get("/readyz", handler.Readyz(logger, r))

		root.Route("/remote", func(state chi.Router) {
			state.Use(basicauth.Basicauth)

			state.Lock("/*", handler.Lock(logger))
			state.Unlock("/*", handler.Unlock(logger))
			state.Get("/*", handler.Fetch(logger))
			state.Post("/*", handler.Update(logger))
			state.Delete("/*", handler.Delete(logger))
		})
	})

	return mux
}
