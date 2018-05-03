package main

import (
	"context"
	"crypto/tls"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/oklog/run"
	"github.com/rs/zerolog/log"
	"github.com/webhippie/terrastate/pkg/config"
	"github.com/webhippie/terrastate/pkg/router"
	"gopkg.in/urfave/cli.v2"
)

var (
	defaultAddr = "0.0.0.0:8080"
	privateAddr = "127.0.0.1:8081"
)

// Server provides the sub-command to start the server.
func Server(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:   "server",
		Usage:  "start the integrated server",
		Flags:  serverFlags(cfg),
		Before: serverBefore(cfg),
		Action: serverAction(cfg),
	}
}

func serverFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "private-addr",
			Value:       privateAddr,
			Usage:       "address for metrics and health",
			EnvVars:     []string{"TERRASTATE_PRIVATE_ADDR"},
			Destination: &cfg.Server.Private,
		},
		&cli.StringFlag{
			Name:        "server-addr",
			Value:       defaultAddr,
			Usage:       "address to bind the server",
			EnvVars:     []string{"TERRASTATE_SERVER_ADDR"},
			Destination: &cfg.Server.Public,
		},
		&cli.StringFlag{
			Name:        "server-root",
			Value:       "/",
			Usage:       "root path of the proxy",
			EnvVars:     []string{"TERRASTATE_SERVER_ROOT"},
			Destination: &cfg.Server.Root,
		},
		&cli.StringFlag{
			Name:        "server-host",
			Value:       "http://localhost:8080",
			Usage:       "external access to server",
			EnvVars:     []string{"TERRASTATE_SERVER_HOST"},
			Destination: &cfg.Server.Host,
		},
		&cli.StringFlag{
			Name:        "server-cert",
			Value:       "",
			Usage:       "path to ssl cert",
			EnvVars:     []string{"TERRASTATE_SERVER_CERT"},
			Destination: &cfg.Server.Cert,
		},
		&cli.StringFlag{
			Name:        "server-key",
			Value:       "",
			Usage:       "path to ssl key",
			EnvVars:     []string{"TERRASTATE_SERVER_KEY"},
			Destination: &cfg.Server.Key,
		},
		&cli.BoolFlag{
			Name:        "strict-curves",
			Value:       false,
			Usage:       "use strict ssl curves",
			EnvVars:     []string{"TERRASTATE_STRICT_CURVES"},
			Destination: &cfg.Server.StrictCurves,
		},
		&cli.BoolFlag{
			Name:        "strict-ciphers",
			Value:       false,
			Usage:       "use strict ssl ciphers",
			EnvVars:     []string{"TERRASTATE_STRICT_CIPHERS"},
			Destination: &cfg.Server.StrictCiphers,
		},
		&cli.StringFlag{
			Name:        "templates-path",
			Value:       "",
			Usage:       "path to custom templates",
			EnvVars:     []string{"TERRASTATE_SERVER_TEMPLATES"},
			Destination: &cfg.Server.Templates,
		},
		&cli.StringFlag{
			Name:        "assets-path",
			Value:       "",
			Usage:       "path to custom assets",
			EnvVars:     []string{"TERRASTATE_SERVER_ASSETS"},
			Destination: &cfg.Server.Assets,
		},
		&cli.StringFlag{
			Name:        "storage-path",
			Value:       "storage/",
			Usage:       "folder for storing certs and misc files",
			EnvVars:     []string{"TERRASTATE_SERVER_STORAGE"},
			Destination: &cfg.Server.Storage,
		},
		&cli.StringFlag{
			Name:        "general-username",
			Value:       "",
			Usage:       "username for basic auth",
			EnvVars:     []string{"TERRASTATE_GENERAL_USERNAME"},
			Destination: &cfg.General.Username,
		},
		&cli.StringFlag{
			Name:        "general-password",
			Value:       "",
			Usage:       "password for basic auth",
			EnvVars:     []string{"TERRASTATE_GENERAL_PASSWORD"},
			Destination: &cfg.General.Password,
		},
		&cli.StringFlag{
			Name:        "encryption-secret",
			Value:       "",
			Usage:       "secret for file encryption",
			EnvVars:     []string{"TERRASTATE_ENCRYPTION_SECRET"},
			Destination: &cfg.General.Secret,
		},
	}
}

func serverBefore(cfg *config.Config) cli.BeforeFunc {
	return func(c *cli.Context) error {
		return nil
	}
}

func serverAction(cfg *config.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		var gr run.Group

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

		if cfg.Server.Cert != "" && cfg.Server.Key != "" {
			cert, err := tls.LoadX509KeyPair(
				cfg.Server.Cert,
				cfg.Server.Key,
			)

			if err != nil {
				log.Info().
					Err(err).
					Msg("failed to load certificates")

				return err
			}

			{
				server := &http.Server{
					Addr:         cfg.Server.Public,
					Handler:      router.Load(cfg),
					ReadTimeout:  5 * time.Second,
					WriteTimeout: 10 * time.Second,
					TLSConfig: &tls.Config{
						PreferServerCipherSuites: true,
						MinVersion:               tls.VersionTLS12,
						CurvePreferences:         curves(cfg),
						CipherSuites:             ciphers(cfg),
						Certificates:             []tls.Certificate{cert},
					},
				}

				gr.Add(func() error {
					log.Info().
						Str("addr", cfg.Server.Public).
						Msg("starting https server")

					return server.ListenAndServeTLS("", "")
				}, func(reason error) {
					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					if err := server.Shutdown(ctx); err != nil {
						log.Info().
							Err(err).
							Msg("failed to stop https server gracefully")

						return
					}

					log.Info().
						Err(reason).
						Msg("https server stopped gracefully")
				})
			}

			return gr.Run()
		}

		{
			server := &http.Server{
				Addr:         cfg.Server.Public,
				Handler:      router.Load(cfg),
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
			}

			gr.Add(func() error {
				log.Info().
					Str("addr", cfg.Server.Public).
					Msg("starting http server")

				return server.ListenAndServe()
			}, func(reason error) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				if err := server.Shutdown(ctx); err != nil {
					log.Info().
						Err(err).
						Msg("failed to stop http server gracefully")

					return
				}

				log.Info().
					Err(reason).
					Msg("http server stopped gracefully")
			})
		}

		{
			server := &http.Server{
				Addr:         cfg.Server.Private,
				Handler:      router.Status(cfg),
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
			}

			gr.Add(func() error {
				log.Info().
					Str("addr", cfg.Server.Private).
					Msg("starting status server")

				return server.ListenAndServe()
			}, func(reason error) {
				ctx, cancel := context.WithTimeout(context.Background(), time.Second)
				defer cancel()

				if err := server.Shutdown(ctx); err != nil {
					log.Info().
						Err(err).
						Msg("failed to stop status server gracefully")

					return
				}

				log.Info().
					Err(reason).
					Msg("status server stopped gracefully")
			})
		}

		return gr.Run()
	}
}

func curves(cfg *config.Config) []tls.CurveID {
	if cfg.Server.StrictCurves {
		return []tls.CurveID{
			tls.CurveP521,
			tls.CurveP384,
			tls.CurveP256,
		}
	}

	return nil
}

func ciphers(cfg *config.Config) []uint16 {
	if cfg.Server.StrictCiphers {
		return []uint16{
			tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
			tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
			tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		}
	}

	return nil
}
