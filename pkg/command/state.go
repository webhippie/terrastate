package command

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/dchest/safefile"
	"github.com/urfave/cli/v2"
	"github.com/webhippie/terrastate/pkg/config"
	"github.com/webhippie/terrastate/pkg/helper"
)

// StateCmd provides the sub-command releated to state.
func StateCmd(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:        "server",
		Usage:       "Read and update state files",
		Flags:       StateFlags(cfg),
		Subcommands: StateCmds(cfg),
	}
}

// StateFlags provides the flags for the state command.
func StateFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
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
	}
}

// StateCmds provides the sub-commands for state.
func StateCmds(cfg *config.Config) []*cli.Command {
	return []*cli.Command{
		{
			Name:      "list",
			Aliases:   []string{"ls"},
			Usage:     "List all states",
			ArgsUsage: " ",
			Flags:     StateListFlags(cfg),
			Action:    StateListAction(cfg),
		},
		{
			Name:      "show",
			Aliases:   []string{"read"},
			Usage:     "Show a state",
			ArgsUsage: "<state>",
			Flags:     StateShowFlags(cfg),
			Action:    StateShowAction(cfg),
		},
		{
			Name:      "encrypt",
			Aliases:   []string{},
			Usage:     "Encrypt a state",
			ArgsUsage: "<state>",
			Flags:     StateEncryptFlags(cfg),
			Action:    StateEncryptAction(cfg),
		},
		{
			Name:      "decrypt",
			Aliases:   []string{},
			Usage:     "Decrypt a state",
			ArgsUsage: "<state>",
			Flags:     StateDecryptFlags(cfg),
			Action:    StateDecryptAction(cfg),
		},
	}
}

// StateListFlags provides the flags for the state list command.
func StateListFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{}
}

// StateListAction provides the action that implements the state list command.
func StateListAction(cfg *config.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		var states []string

		err := filepath.Walk(cfg.Server.Storage, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "Failed to read %s: %q\n", path, err)
				return err
			}

			if info.Name() == "terraform.tfstate" {
				state := strings.TrimPrefix(
					strings.Replace(
						filepath.Dir(
							path,
						),
						cfg.Server.Storage,
						"",
						-1,
					),
					"/",
				)

				states = append(states, state)
			}

			return nil
		})

		if err != nil {
			return cli.Exit("Failed to list", 2)
		}

		fmt.Fprintln(os.Stdout, strings.Join(states, "\n"))
		return nil
	}
}

// StateShowFlags provides the flags for the state show command.
func StateShowFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{}
}

// StateShowAction provides the action that implements the state show command.
func StateShowAction(cfg *config.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		state := c.Args().Get(0)

		if state == "" {
			return cli.Exit("Missing state argument", 2)
		}

		full := path.Join(
			cfg.Server.Storage,
			state,
			"terraform.tfstate",
		)

		if _, err := os.Stat(full); os.IsNotExist(err) {
			return cli.Exit("State does not exist", 3)
		}

		file, err := ioutil.ReadFile(
			full,
		)

		if err != nil {
			return cli.Exit("Failed to read state", 4)
		}

		fmt.Fprintln(os.Stdout, string(file))
		return nil
	}
}

// StateEncryptFlags provides the flags for the state encrypt command.
func StateEncryptFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{}
}

// StateEncryptAction provides the action that implements the state encrypt command.
func StateEncryptAction(cfg *config.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		if cfg.General.Secret == "" {
			return cli.Exit("Missing encryption secret", 2)
		}

		state := c.Args().Get(0)

		if state == "" {
			return cli.Exit("Missing state argument", 3)
		}

		full := path.Join(
			cfg.Server.Storage,
			state,
			"terraform.tfstate",
		)

		if _, err := os.Stat(full); os.IsNotExist(err) {
			return cli.Exit("State does not exist", 4)
		}

		file, err := ioutil.ReadFile(
			full,
		)

		if err != nil {
			return cli.Exit("Failed to read state", 5)
		}

		encrypted, err := helper.Encrypt(file, []byte(cfg.General.Secret))

		if err != nil {
			return cli.Exit("Failed to encrypt state", 6)
		}

		if err := safefile.WriteFile(full, encrypted, 0644); err != nil {
			return cli.Exit("Failed to update file", 7)
		}

		fmt.Fprintln(os.Stderr, "Successfully encrypted state")
		return nil
	}
}

// StateDecryptFlags provides the flags for the state decrypt command.
func StateDecryptFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{}
}

// StateDecryptAction provides the action that implements the state decrypt command.
func StateDecryptAction(cfg *config.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		if cfg.General.Secret == "" {
			return cli.Exit("Missing encryption secret", 2)
		}

		state := c.Args().Get(0)

		if state == "" {
			return cli.Exit("Missing state argument", 3)
		}

		full := path.Join(
			cfg.Server.Storage,
			state,
			"terraform.tfstate",
		)

		if _, err := os.Stat(full); os.IsNotExist(err) {
			return cli.Exit("State does not exist", 4)
		}

		file, err := ioutil.ReadFile(
			full,
		)

		if err != nil {
			return cli.Exit("Failed to read state", 5)
		}

		decrypted, err := helper.Decrypt(file, []byte(cfg.General.Secret))

		if err != nil {
			return cli.Exit("Failed to decrypt state", 6)
		}

		if err := safefile.WriteFile(full, decrypted, 0644); err != nil {
			return cli.Exit("Failed to update file", 7)
		}

		fmt.Fprintln(os.Stderr, "Successfully encrypted state")
		return nil
	}
}
