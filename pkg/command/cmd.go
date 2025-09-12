package command

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/webhippie/terrastate/pkg/config"
	"github.com/webhippie/terrastate/pkg/version"
)

var (
	rootCmd = &cobra.Command{
		Use:           "terrastate",
		Short:         "Terraform HTTP remote state storage",
		Version:       version.String,
		SilenceErrors: false,
		SilenceUsage:  true,

		PersistentPreRunE: func(_ *cobra.Command, _ []string) error {
			return setupLogger()
		},

		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}

	cfg *config.Config
)

func init() {
	cfg = config.Load()
	cobra.OnInitialize(setupConfig)

	rootCmd.PersistentFlags().BoolP("help", "h", false, "Show the help, so what you see now")
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Print the current version of that tool")

	rootCmd.PersistentFlags().String("config-file", "", "Path to optional config file")
	_ = viper.BindPFlag("config.file", rootCmd.PersistentFlags().Lookup("config-file"))

	rootCmd.PersistentFlags().String("log-level", "info", "Set logging level")
	viper.SetDefault("log.level", "info")
	_ = viper.BindPFlag("log.level", rootCmd.PersistentFlags().Lookup("log-level"))

	rootCmd.PersistentFlags().Bool("log-pretty", true, "Enable pretty logging")
	viper.SetDefault("log.pretty", true)
	_ = viper.BindPFlag("log.pretty", rootCmd.PersistentFlags().Lookup("log-pretty"))

	rootCmd.PersistentFlags().Bool("log-color", true, "Enable colored logging")
	viper.SetDefault("log.color", true)
	_ = viper.BindPFlag("log.color", rootCmd.PersistentFlags().Lookup("log-color"))
}

// Run parses the command line arguments and executes the program.
func Run() error {
	return rootCmd.Execute()
}
