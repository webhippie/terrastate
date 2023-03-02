package command

import (
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/dchest/safefile"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"github.com/webhippie/terrastate/pkg/helper"
)

var (
	stateCmd = &cobra.Command{
		Use:   "state",
		Short: "Read and update state",
	}

	stateListCmd = &cobra.Command{
		Use:   "list",
		Short: "List all states",
		Run:   stateListAction,
		Args:  cobra.NoArgs,
	}

	stateShowCmd = &cobra.Command{
		Use:   "show <state>",
		Short: "Show a state",
		Run:   stateShowAction,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("missing state argument")
			}

			return nil
		},
	}

	stateEncryptCmd = &cobra.Command{
		Use:   "encrypt <state>",
		Short: "Encrypt a state",
		Run:   stateEncryptAction,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("missing state argument")
			}

			return nil
		},
	}

	stateDecryptCmd = &cobra.Command{
		Use:   "decrypt <state>",
		Short: "Decrypt a state",
		Run:   stateDecryptAction,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) != 1 {
				return fmt.Errorf("missing state argument")
			}

			return nil
		},
	}
)

func init() {
	rootCmd.AddCommand(stateCmd)

	stateCmd.PersistentFlags().String("storage-path", defaultServerStorage, "Folder for storing the states")
	viper.SetDefault("server.storage", defaultServerStorage)
	viper.BindPFlag("server.storage", stateCmd.PersistentFlags().Lookup("server-storage"))

	stateCmd.PersistentFlags().String("encryption-secret", defaultEncryptionSecret, "Secret for file encryption")
	viper.SetDefault("encryption.secret", defaultEncryptionSecret)
	viper.BindPFlag("encryption.secret", stateCmd.PersistentFlags().Lookup("encryption-secret"))

	stateCmd.AddCommand(stateListCmd)
	stateCmd.AddCommand(stateShowCmd)
	stateCmd.AddCommand(stateEncryptCmd)
	stateCmd.AddCommand(stateDecryptCmd)
}

func stateListAction(_ *cobra.Command, _ []string) {
	var (
		states []string
	)

	err := filepath.Walk(cfg.Server.Storage, func(path string, info os.FileInfo, err error) error {
		if err != nil {
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
		cobra.CheckErr("Failed to list")
	}

	if len(states) > 0 {
		fmt.Fprintln(os.Stdout, strings.Join(states, "\n"))
	}
}

func stateShowAction(_ *cobra.Command, args []string) {
	state := args[0]

	full := path.Join(
		cfg.Server.Storage,
		state,
		"terraform.tfstate",
	)

	if _, err := os.Stat(full); os.IsNotExist(err) {
		cobra.CheckErr("State does not exist")
	}

	file, err := os.ReadFile(
		full,
	)

	if err != nil {
		cobra.CheckErr("Failed to read state")
	}

	fmt.Fprintln(os.Stdout, string(file))
}

func stateEncryptAction(_ *cobra.Command, args []string) {
	if cfg.Encryption.Secret == "" {
		cobra.CheckErr("Missing encryption secret")
	}

	state := args[0]

	full := path.Join(
		cfg.Server.Storage,
		state,
		"terraform.tfstate",
	)

	if _, err := os.Stat(full); os.IsNotExist(err) {
		cobra.CheckErr("State does not exist")
	}

	file, err := os.ReadFile(
		full,
	)

	if err != nil {
		cobra.CheckErr("Failed to read state")
	}

	encrypted, err := helper.Encrypt(
		file,
		[]byte(cfg.Encryption.Secret),
	)

	if err != nil {
		cobra.CheckErr("Failed to encrypt state")
	}

	if err := safefile.WriteFile(full, encrypted, 0644); err != nil {
		cobra.CheckErr("Failed to update file")
	}

	fmt.Fprintln(os.Stderr, "Successfully encrypted state")
}

func stateDecryptAction(_ *cobra.Command, args []string) {
	if cfg.Encryption.Secret == "" {
		cobra.CheckErr("Missing encryption secret")
	}

	state := args[0]

	full := path.Join(
		cfg.Server.Storage,
		state,
		"terraform.tfstate",
	)

	if _, err := os.Stat(full); os.IsNotExist(err) {
		cobra.CheckErr("State does not exist")
	}

	file, err := os.ReadFile(
		full,
	)

	if err != nil {
		cobra.CheckErr("Failed to read state")
	}

	decrypted, err := helper.Decrypt(
		file,
		[]byte(cfg.Encryption.Secret),
	)

	if err != nil {
		cobra.CheckErr("Failed to decrypt state")
	}

	if err := safefile.WriteFile(full, decrypted, 0644); err != nil {
		cobra.CheckErr("Failed to update file")
	}

	fmt.Fprintln(os.Stderr, "Successfully encrypted state")
}
