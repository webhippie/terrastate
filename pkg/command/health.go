package command

import (
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"github.com/urfave/cli/v2"
	"github.com/webhippie/terrastate/pkg/config"
)

// HealthCmd provides the sub-command releated to health.
func HealthCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:   "health",
		Usage:  "Perform health checks",
		Flags:  HealthFlags(cfg),
		Action: HealthAction(cfg),
	}
}

// HealthFlags provides the flags for the health command.
func HealthFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "metrics-addr",
			Value:       "0.0.0.0:8081",
			Usage:       "Address to bind the metrics",
			EnvVars:     []string{"TERRASTATE_METRICS_ADDR"},
			Destination: &cfg.Metrics.Addr,
		},
	}
}

// HealthAction provides the action that implements the health command.
func HealthAction(cfg *config.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		resp, err := http.Get(
			fmt.Sprintf(
				"http://%s/healthz",
				cfg.Metrics.Addr,
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
				Int("code", resp.StatusCode).
				Msg("health seems to be in bad state")

			return err
		}

		return nil
	}
}
