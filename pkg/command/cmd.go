package command

import (
	"os"

	"github.com/urfave/cli/v2"
	"github.com/webhippie/terrastate/pkg/config"
	"github.com/webhippie/terrastate/pkg/version"
)

// Run parses the command line arguments and executes the program.
func Run() error {
	cfg := config.Load()

	app := &cli.App{
		Name:     "terrastate",
		Version:  version.String,
		Usage:    "Terraform HTTP remote state storage",
		Authors:  AuthorList(),
		Flags:    GlobalFlags(cfg),
		Before:   GlobalBefore(cfg),
		Commands: GlobalCmds(cfg),
	}

	cli.HelpFlag = &cli.BoolFlag{
		Name:    "help",
		Aliases: []string{"h"},
		Usage:   "Show the help, so what you see now",
	}

	cli.VersionFlag = &cli.BoolFlag{
		Name:    "version",
		Aliases: []string{"v"},
		Usage:   "Print the current version of that tool",
	}

	return app.Run(os.Args)
}

// AuthorList defines the list of authors.
func AuthorList() []*cli.Author {
	return []*cli.Author{
		{
			Name:  "Thomas Boerger",
			Email: "thomas@webhippie.de",
		},
	}
}

// GlobalFlags defines the list of global flags.
func GlobalFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "log-level",
			Value:       "info",
			Usage:       "Set logging level",
			EnvVars:     []string{"TERRASTATE_LOG_LEVEL"},
			Destination: &cfg.Logs.Level,
		},
		&cli.BoolFlag{
			Name:        "log-pretty",
			Value:       true,
			Usage:       "Enable pretty logging",
			EnvVars:     []string{"TERRASTATE_LOG_PRETTY"},
			Destination: &cfg.Logs.Pretty,
		},
		&cli.BoolFlag{
			Name:        "log-color",
			Value:       true,
			Usage:       "Enable colored logging",
			EnvVars:     []string{"TERRASTATE_LOG_COLOR"},
			Destination: &cfg.Logs.Color,
		},
	}
}

// GlobalBefore defines a global hook for setup.
func GlobalBefore(cfg *config.Config) cli.BeforeFunc {
	return func(c *cli.Context) error {
		return setupLogger(cfg)
	}
}

// GlobalCmds defines the global commands.
func GlobalCmds(cfg *config.Config) []*cli.Command {
	return []*cli.Command{
		ServerCmd(cfg),
		HealthCmd(cfg),
		StateCmd(cfg),
	}
}
