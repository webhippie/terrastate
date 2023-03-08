package command

import (
	"fmt"
	"net/http"
	"os"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	healthCmd = &cobra.Command{
		Use:   "health",
		Short: "Perform health checks",
		Run:   healthAction,
		Args:  cobra.NoArgs,
	}
)

func init() {
	rootCmd.AddCommand(healthCmd)

	healthCmd.PersistentFlags().String("metrics-addr", defaultMetricsAddr, "Address to bind the metrics")
	viper.SetDefault("metrics.addr", defaultMetricsAddr)
	viper.BindPFlag("metrics.addr", healthCmd.PersistentFlags().Lookup("metrics-addr"))
}

func healthAction(_ *cobra.Command, _ []string) {
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

		os.Exit(1)
	}

	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Error().
			Int("code", resp.StatusCode).
			Msg("health seems to be in bad state")

		os.Exit(1)
	}
}
