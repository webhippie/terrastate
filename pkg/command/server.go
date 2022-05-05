package command

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/oklog/run"
	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/webhippie/terrastate/pkg/config"
	"github.com/webhippie/terrastate/pkg/router"
)

// ServerCmd provides the sub-command releated to server.
func ServerCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:   "server",
		Usage:  "Start integrated server",
		Flags:  ServerFlags(cfg),
		Action: ServerAction(cfg),
	}
}

// ServerFlags provides the flags for the server command.
func ServerFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "metrics-addr",
			Value:       "0.0.0.0:8081",
			Usage:       "Address to bind the metrics",
			EnvVars:     []string{"TERRASTATE_METRICS_ADDR"},
			Destination: &cfg.Metrics.Addr,
		},
		&cli.StringFlag{
			Name:        "metrics-token",
			Value:       "",
			Usage:       "Token to make metrics secure",
			EnvVars:     []string{"TERRASTATE_METRICS_TOKEN"},
			Destination: &cfg.Metrics.Token,
		},
		&cli.StringFlag{
			Name:        "server-addr",
			Value:       "0.0.0.0:8080",
			Usage:       "Address to bind the server",
			EnvVars:     []string{"TERRASTATE_SERVER_ADDR"},
			Destination: &cfg.Server.Addr,
		},
		&cli.BoolFlag{
			Name:        "server-pprof",
			Value:       false,
			Usage:       "Enable pprof debugging",
			EnvVars:     []string{"TERRASTATE_SERVER_PPROF"},
			Destination: &cfg.Server.Pprof,
		},
		&cli.StringFlag{
			Name:        "server-root",
			Value:       "/",
			Usage:       "Root path of the server",
			EnvVars:     []string{"TERRASTATE_SERVER_ROOT"},
			Destination: &cfg.Server.Root,
		},
		&cli.StringFlag{
			Name:        "server-host",
			Value:       "http://localhost:8080",
			Usage:       "External access to server",
			EnvVars:     []string{"TERRASTATE_SERVER_HOST"},
			Destination: &cfg.Server.Host,
		},
		&cli.StringFlag{
			Name:        "server-cert",
			Value:       "",
			Usage:       "Path to ssl cert",
			EnvVars:     []string{"TERRASTATE_SERVER_CERT"},
			Destination: &cfg.Server.Cert,
		},
		&cli.StringFlag{
			Name:        "server-key",
			Value:       "",
			Usage:       "Path to ssl key",
			EnvVars:     []string{"TERRASTATE_SERVER_KEY"},
			Destination: &cfg.Server.Key,
		},
		&cli.BoolFlag{
			Name:        "strict-curves",
			Value:       false,
			Usage:       "Use strict ssl curves",
			EnvVars:     []string{"TERRASTATE_STRICT_CURVES"},
			Destination: &cfg.Server.StrictCurves,
		},
		&cli.BoolFlag{
			Name:        "strict-ciphers",
			Value:       false,
			Usage:       "Use strict ssl ciphers",
			EnvVars:     []string{"TERRASTATE_STRICT_CIPHERS"},
			Destination: &cfg.Server.StrictCiphers,
		},
		&cli.StringFlag{
			Name:        "storage-path",
			Value:       "storage/",
			Usage:       "Folder for storing certs and misc files",
			EnvVars:     []string{"TERRASTATE_SERVER_STORAGE"},
			Destination: &cfg.Server.Storage,
		},
		&cli.StringFlag{
			Name:        "encryption-secret",
			Value:       "",
			Usage:       "Secret for file encryption",
			EnvVars:     []string{"TERRASTATE_ENCRYPTION_SECRET"},
			Destination: &cfg.General.Secret,
		},
		&cli.StringFlag{
			Name:        "general-username",
			Value:       "",
			Usage:       "Username for basic auth",
			EnvVars:     []string{"TERRASTATE_GENERAL_USERNAME"},
			Destination: &cfg.General.Username,
		},
		&cli.StringFlag{
			Name:        "general-password",
			Value:       "",
			Usage:       "Password for basic auth",
			EnvVars:     []string{"TERRASTATE_GENERAL_PASSWORD"},
			Destination: &cfg.General.Password,
		},
	}
}

// ServerAction provides the action that implements the server command.
func ServerAction(cfg *config.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		var gr run.Group

		if cfg.Server.Cert != "" && cfg.Server.Key != "" {
			cert, err := tls.LoadX509KeyPair(
				cfg.Server.Cert,
				cfg.Server.Key,
			)

			if err != nil {
				log.Info().
					Err(err).
					Msg("Failed to load certificates")

				return err
			}

			server := &http.Server{
				Addr:         cfg.Server.Addr,
				Handler:      router.Load(cfg),
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
				TLSConfig: &tls.Config{
					PreferServerCipherSuites: true,
					MinVersion:               tls.VersionTLS12,
					CurvePreferences:         router.Curves(cfg),
					CipherSuites:             router.Ciphers(cfg),
					Certificates:             []tls.Certificate{cert},
				},
			}

			gr.Add(func() error {
				log.Info().
					Str("addr", cfg.Server.Addr).
					Msg("Starting HTTPS server")

				return server.ListenAndServeTLS("", "")
			}, func(reason error) {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()

				if err := server.Shutdown(ctx); err != nil {
					log.Error().
						Err(err).
						Msg("Failed to shutdown HTTPS gracefully")

					return
				}

				log.Info().
					Err(reason).
					Msg("Shutdown HTTPS gracefully")
			})
		} else {
			server := &http.Server{
				Addr:         cfg.Server.Addr,
				Handler:      router.Load(cfg),
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
			}

			gr.Add(func() error {
				log.Info().
					Str("addr", cfg.Server.Addr).
					Msg("Starting HTTP server")

				return server.ListenAndServe()
			}, func(reason error) {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()

				if err := server.Shutdown(ctx); err != nil {
					log.Error().
						Err(err).
						Msg("Failed to shutdown HTTP gracefully")

					return
				}

				log.Info().
					Err(reason).
					Msg("Shutdown HTTP gracefully")
			})
		}

		{
			server := &http.Server{
				Addr:         cfg.Metrics.Addr,
				Handler:      router.Metrics(cfg),
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
			}

			gr.Add(func() error {
				log.Info().
					Str("addr", cfg.Metrics.Addr).
					Msg("Starting metrics server")

				return server.ListenAndServe()
			}, func(reason error) {
				ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
				defer cancel()

				if err := server.Shutdown(ctx); err != nil {
					log.Error().
						Err(err).
						Msg("Failed to shutdown metrics gracefully")

					return
				}

				log.Info().
					Err(reason).
					Msg("Shutdown metrics gracefully")
			})
		}

		{
			stop := make(chan os.Signal, 1)

			gr.Add(func() error {
				signal.Notify(stop, os.Interrupt)

				<-stop

				return nil
			}, func(err error) {
				close(stop)
			})
		}

		return gr.Run()
	}
}
