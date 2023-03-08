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
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/webhippie/terrastate/pkg/router"
)

var (
	serverCmd = &cobra.Command{
		Use:   "server",
		Short: "Start integrated server",
		Run:   serverAction,
		Args:  cobra.NoArgs,
	}

	defaultMetricsAddr         = "0.0.0.0:8081"
	defaultServerAddr          = "0.0.0.0:8080"
	defaultServerPprof         = false
	defaultServerRoot          = "/"
	defaultServerHost          = "http://localhost:8080"
	defaultServerCert          = ""
	defaultServerKey           = ""
	defaultServerStrictCurves  = false
	defaultServerStrictCiphers = false
	defaultServerStorage       = "storage/"
	defaultEncryptionSecret    = ""
	defaultAccessUsername      = ""
	defaultAccessPassword      = ""
)

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.PersistentFlags().String("metrics-addr", defaultMetricsAddr, "Address to bind the metrics")
	viper.SetDefault("metrics.addr", defaultMetricsAddr)
	viper.BindPFlag("metrics.addr", serverCmd.PersistentFlags().Lookup("metrics-addr"))

	serverCmd.PersistentFlags().String("metrics-token", "", "Token to make metrics secure")
	viper.SetDefault("metrics.token", "")
	viper.BindPFlag("metrics.token", serverCmd.PersistentFlags().Lookup("metrics-token"))

	serverCmd.PersistentFlags().String("server-addr", defaultServerAddr, "Address to bind the server")
	viper.SetDefault("server.addr", defaultServerAddr)
	viper.BindPFlag("server.addr", serverCmd.PersistentFlags().Lookup("server-addr"))

	serverCmd.PersistentFlags().Bool("server-pprof", defaultServerPprof, "Enable pprof debugging")
	viper.SetDefault("server.pprof", defaultServerPprof)
	viper.BindPFlag("server.pprof", serverCmd.PersistentFlags().Lookup("server-pprof"))

	serverCmd.PersistentFlags().String("server-root", defaultServerRoot, "Root path of the server")
	viper.SetDefault("server.root", defaultServerRoot)
	viper.BindPFlag("server.root", serverCmd.PersistentFlags().Lookup("server-root"))

	serverCmd.PersistentFlags().String("server-host", defaultServerHost, "External access to server")
	viper.SetDefault("server.host", defaultServerHost)
	viper.BindPFlag("server.host", serverCmd.PersistentFlags().Lookup("server-host"))

	serverCmd.PersistentFlags().String("server-cert", defaultServerCert, "Path to cert for SSL encryption")
	viper.SetDefault("server.cert", defaultServerCert)
	viper.BindPFlag("server.cert", serverCmd.PersistentFlags().Lookup("server-cert"))

	serverCmd.PersistentFlags().String("server-key", defaultServerKey, "Path to key for SSL encryption")
	viper.SetDefault("server.key", defaultServerKey)
	viper.BindPFlag("server.key", serverCmd.PersistentFlags().Lookup("server-key"))

	serverCmd.PersistentFlags().Bool("strict-curves", defaultServerStrictCurves, "Use strict SSL curves")
	viper.SetDefault("server.strict_curves", defaultServerStrictCurves)
	viper.BindPFlag("server.strict_curves", serverCmd.PersistentFlags().Lookup("strict-curves"))

	serverCmd.PersistentFlags().Bool("strict-ciphers", defaultServerStrictCiphers, "Use strict SSL ciphers")
	viper.SetDefault("server.strict_ciphers", defaultServerStrictCiphers)
	viper.BindPFlag("server.strict_ciphers", serverCmd.PersistentFlags().Lookup("strict-ciphers"))

	serverCmd.PersistentFlags().String("storage-path", defaultServerStorage, "Folder for storing the states")
	viper.SetDefault("server.storage", defaultServerStorage)
	viper.BindPFlag("server.storage", serverCmd.PersistentFlags().Lookup("server-storage"))

	serverCmd.PersistentFlags().String("encryption-secret", defaultEncryptionSecret, "Secret for file encryption")
	viper.SetDefault("encryption.secret", defaultEncryptionSecret)
	viper.BindPFlag("encryption.secret", serverCmd.PersistentFlags().Lookup("encryption-secret"))

	serverCmd.PersistentFlags().String("general-username", defaultAccessUsername, "Username for basic auth")
	viper.SetDefault("access.username", defaultAccessUsername)
	viper.BindPFlag("access.username", serverCmd.PersistentFlags().Lookup("general-username"))

	serverCmd.PersistentFlags().String("general-password", defaultAccessPassword, "Password for basic auth")
	viper.SetDefault("access.password", defaultAccessPassword)
	viper.BindPFlag("access.password", serverCmd.PersistentFlags().Lookup("general-password"))
}

func serverAction(_ *cobra.Command, _ []string) {
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

			os.Exit(1)
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

	if err := gr.Run(); err != nil {
		os.Exit(1)
	}
}
