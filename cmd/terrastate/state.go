package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/dchest/safefile"
	"github.com/webhippie/terrastate/pkg/config"
	"github.com/webhippie/terrastate/pkg/helper"
	"gopkg.in/urfave/cli.v2"
)

// State provides the sub-command to access states.
func State(cfg *config.Config) *cli.Command {
	return &cli.Command{
		Name:        "state",
		Usage:       "read and update state files",
		Flags:       stateFlags(cfg),
		Subcommands: stateCommands(cfg),
	}
}

func stateFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{
		&cli.StringFlag{
			Name:        "storage-path",
			Value:       "storage/",
			Usage:       "folder for storing certs and misc files",
			EnvVars:     []string{"TERRASTATE_SERVER_STORAGE"},
			Destination: &cfg.Server.Storage,
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

func stateCommands(cfg *config.Config) []*cli.Command {
	return []*cli.Command{
		{
			Name:      "list",
			Aliases:   []string{"ls"},
			Usage:     "list all states",
			ArgsUsage: " ",
			Flags:     stateListFlags(cfg),
			Action:    stateListAction(cfg),
		},
		{
			Name:      "show",
			Aliases:   []string{"read"},
			Usage:     "show a state",
			ArgsUsage: "<state>",
			Flags:     stateShowFlags(cfg),
			Action:    stateShowAction(cfg),
		},
		{
			Name:      "encrypt",
			Aliases:   []string{},
			Usage:     "encrypt a state",
			ArgsUsage: "<state>",
			Flags:     stateEncryptFlags(cfg),
			Action:    stateEncryptAction(cfg),
		},
		{
			Name:      "decrypt",
			Aliases:   []string{},
			Usage:     "decrypt a state",
			ArgsUsage: "<state>",
			Flags:     stateDecryptFlags(cfg),
			Action:    stateDecryptAction(cfg),
		},
	}
}

func stateListFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{}
}

func stateListAction(cfg *config.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		var states []string

		err := filepath.Walk(cfg.Server.Storage, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Fprintf(os.Stderr, "failed to read %s: %q\n", path, err)
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
			return cli.Exit("failed to list", 2)
		}

		fmt.Fprintln(os.Stdout, strings.Join(states, "\n"))
		return nil
	}
}

func stateShowFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{}
}

func stateShowAction(cfg *config.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		state := c.Args().Get(0)

		if state == "" {
			return cli.Exit("missing state argument", 2)
		}

		full := path.Join(
			cfg.Server.Storage,
			state,
			"terraform.tfstate",
		)

		if _, err := os.Stat(full); os.IsNotExist(err) {
			return cli.Exit("state does not exist", 3)
		}

		file, err := ioutil.ReadFile(
			full,
		)

		if err != nil {
			return cli.Exit("failed to read state", 4)
		}

		fmt.Fprintln(os.Stdout, string(file))
		return nil
	}
}

func stateEncryptFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{}
}

func stateEncryptAction(cfg *config.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		if cfg.General.Secret == "" {
			return cli.Exit("missing encryption secret", 2)
		}

		state := c.Args().Get(0)

		if state == "" {
			return cli.Exit("missing state argument", 3)
		}

		full := path.Join(
			cfg.Server.Storage,
			state,
			"terraform.tfstate",
		)

		if _, err := os.Stat(full); os.IsNotExist(err) {
			return cli.Exit("state does not exist", 4)
		}

		file, err := ioutil.ReadFile(
			full,
		)

		if err != nil {
			return cli.Exit("failed to read state", 5)
		}

		encrypted, err := helper.Encrypt(file, []byte(cfg.General.Secret))

		if err != nil {
			return cli.Exit("failed to encrypt state", 6)
		}

		if err := safefile.WriteFile(full, encrypted, 0644); err != nil {
			return cli.Exit("failed to update file", 7)
		}

		fmt.Fprintln(os.Stderr, "successfully encrypted state")
		return nil
	}
}

func stateDecryptFlags(cfg *config.Config) []cli.Flag {
	return []cli.Flag{}
}

func stateDecryptAction(cfg *config.Config) cli.ActionFunc {
	return func(c *cli.Context) error {
		if cfg.General.Secret == "" {
			return cli.Exit("missing encryption secret", 2)
		}

		state := c.Args().Get(0)

		if state == "" {
			return cli.Exit("missing state argument", 3)
		}

		full := path.Join(
			cfg.Server.Storage,
			state,
			"terraform.tfstate",
		)

		if _, err := os.Stat(full); os.IsNotExist(err) {
			return cli.Exit("state does not exist", 4)
		}

		file, err := ioutil.ReadFile(
			full,
		)

		if err != nil {
			return cli.Exit("failed to read state", 5)
		}

		decrypted, err := helper.Decrypt(file, []byte(cfg.General.Secret))

		if err != nil {
			return cli.Exit("failed to decrypt state", 6)
		}

		if err := safefile.WriteFile(full, decrypted, 0644); err != nil {
			return cli.Exit("failed to update file", 7)
		}

		fmt.Fprintln(os.Stderr, "successfully encrypted state")
		return nil
	}
}
