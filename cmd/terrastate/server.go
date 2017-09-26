package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
	"github.com/oklog/oklog/pkg/group"
	"github.com/webhippie/terrastate/pkg/config"
	"github.com/webhippie/terrastate/pkg/router"
	"gopkg.in/urfave/cli.v2"
)

// Server provides the sub-command to start the server.
func Server() *cli.Command {
	return &cli.Command{
		Name:  "server",
		Usage: "Start the integrated server",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "host",
				Value:       "http://localhost:8080",
				Usage:       "external access to server",
				EnvVars:     []string{"TERRASTATE_HOST"},
				Destination: &config.Server.Host,
			},
			&cli.StringFlag{
				Name:        "addr",
				Value:       "0.0.0.0:8080",
				Usage:       "address to bind the server",
				EnvVars:     []string{"TERRASTATE_ADDR"},
				Destination: &config.Server.Addr,
			},
			&cli.BoolFlag{
				Name:        "pprof",
				Value:       false,
				Usage:       "enable pprof debugging server",
				EnvVars:     []string{"TERRASTATE_PPROF"},
				Destination: &config.Server.Pprof,
			},
			&cli.BoolFlag{
				Name:        "prometheus",
				Value:       false,
				Usage:       "enable prometheus exporter",
				EnvVars:     []string{"TERRASTATE_PROMETHEUS"},
				Destination: &config.Server.Prometheus,
			},
			&cli.StringFlag{
				Name:        "cert",
				Value:       "",
				Usage:       "path to ssl cert",
				EnvVars:     []string{"TERRASTATE_CERT"},
				Destination: &config.Server.Cert,
			},
			&cli.StringFlag{
				Name:        "key",
				Value:       "",
				Usage:       "path to ssl key",
				EnvVars:     []string{"TERRASTATE_KEY"},
				Destination: &config.Server.Key,
			},
			&cli.StringFlag{
				Name:        "storage",
				Value:       "storage/",
				Usage:       "folder for storing states",
				EnvVars:     []string{"TERRASTATE_STORAGE"},
				Destination: &config.Server.Storage,
			},
			&cli.StringFlag{
				Name:        "username",
				Value:       "",
				Usage:       "username for basic auth",
				EnvVars:     []string{"TERRASTATE_USERNAME"},
				Destination: &config.General.Username,
			},
			&cli.StringFlag{
				Name:        "password",
				Value:       "",
				Usage:       "password for basic auth",
				EnvVars:     []string{"TERRASTATE_PASSWORD"},
				Destination: &config.General.Password,
			},
		},
		Before: func(c *cli.Context) error {
			return nil
		},
		Action: func(c *cli.Context) error {
			logger := log.NewLogfmtLogger(log.NewSyncWriter(os.Stdout))

			switch strings.ToLower(config.LogLevel) {
			case "debug":
				logger = level.NewFilter(logger, level.AllowDebug())
			case "warn":
				logger = level.NewFilter(logger, level.AllowWarn())
			case "error":
				logger = level.NewFilter(logger, level.AllowError())
			default:
				logger = level.NewFilter(logger, level.AllowInfo())
			}

			logger = log.WithPrefix(logger,
				"app", c.App.Name,
				"ts", log.DefaultTimestampUTC,
			)

			cfg, err := ssl()

			if err != nil {
				level.Error(logger).Log(
					"msg", "failed to load certificates",
					"err", err,
				)

				return err
			}

			server := &http.Server{
				Addr:         config.Server.Addr,
				Handler:      router.Load(logger),
				ReadTimeout:  5 * time.Second,
				WriteTimeout: 10 * time.Second,
				TLSConfig:    cfg,
			}

			var gr group.Group

			{
				gr.Add(func() error {
					level.Info(logger).Log(
						"msg", "starting server",
						"addr", config.Server.Addr,
					)

					if server.TLSConfig != nil {
						return server.ListenAndServeTLS("", "")
					}

					return server.ListenAndServe()
				}, func(err error) {
					if err != nil {
						level.Error(logger).Log(
							"msg", "failed to start server",
							"err", err,
						)

						return
					}

					ctx, cancel := context.WithTimeout(context.Background(), time.Second)
					defer cancel()

					if err := server.Shutdown(ctx); err != nil {
						level.Error(logger).Log(
							"msg", "failed to shutdown server gracefully",
							"err", err,
						)

						return
					}

					level.Info(logger).Log(
						"msg", "server shutdown gracefully",
					)
				})
			}

			return gr.Run()
		},
	}
}

func curves() []tls.CurveID {
	return []tls.CurveID{
		tls.CurveP521,
		tls.CurveP384,
		tls.CurveP256,
	}
}

func ciphers() []uint16 {
	return []uint16{
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
	}
}

func ssl() (*tls.Config, error) {
	if config.Server.Cert != "" && config.Server.Key != "" {
		cert, err := tls.LoadX509KeyPair(
			config.Server.Cert,
			config.Server.Key,
		)

		if err != nil {
			return nil, fmt.Errorf("failed to load certificates. %s", err)
		}

		return &tls.Config{
			PreferServerCipherSuites: true,
			MinVersion:               tls.VersionTLS12,
			CurvePreferences:         curves(),
			CipherSuites:             ciphers(),
			Certificates:             []tls.Certificate{cert},
		}, nil
	}

	return nil, nil
}
