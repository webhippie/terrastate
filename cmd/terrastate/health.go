package main

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/webhippie/terrastate/pkg/config"
	"gopkg.in/urfave/cli.v2"
)

// Health provides the sub-command to perform a health check.
func Health(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:   "health",
		Usage:  "perform health checks for service",
		Flags:  healthFlags(cfg),
		Action: healthAction(cfg),
	}
}

func healthFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "private-addr",
			Value:       privateAddr,
			Usage:       "address for metrics and health",
			EnvVars:     []string{"TERRASTATE_PRIVATE_ADDR"},
			Destination: &cfg.Server.Private,
		},
	}
}

func healthAction(cfg *config.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		resp, err := http.Get(
			fmt.Sprintf(
				"http://%s/healthz",
				cfg.Server.Private,
			),
		)

		if err != nil {
			log.Error().
				Err(err).
				Msg("failed to request health check")

			return err
		}

		defer resp.Body.Close()

		if resp.StatusCode != 200 {
			log.Error().
				Err(err).
				Msg("health seems to be in a bad state")

			return err
		}

		return nil
	}
}
