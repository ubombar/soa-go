package main

import (
	"errors"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/ubombar/soa/internal/add"
	"github.com/ubombar/soa/internal/log"
	"github.com/ubombar/soa/internal/sync"
)

var logger = log.GlobalLogger

func main() {
	vaultDir := os.Getenv("SOA_DIR")

	logger := log.GlobalLogger
	rootCmd := &cobra.Command{
		Use:               "soa",
		Short:             "State of the art manager",
		Long:              "State of the Art File Manager, used for managing Zettelkasten notes",
		PersistentPreRunE: rootCmdPersistentPreRunE,
		Run:               rootCmd,
	}
	rootCmd.PersistentFlags().Bool("debug", false, "enable debug messages")
	rootCmd.PersistentFlags().String("vault-dir", vaultDir, "vault dir")

	// add other commands
	rootCmd.AddCommand(add.AddCmd())
	rootCmd.AddCommand(sync.SyncCmd())

	// bind variables to viper
	viper.BindPFlags(rootCmd.PersistentFlags())

	if err := rootCmd.Execute(); err != nil {
		logger.Fatalf("There was an error while running the command: %v", err)
	}
}

func rootCmd(cmd *cobra.Command, args []string) {
	cmd.Help()
}

func rootCmdPersistentPreRunE(cmd *cobra.Command, args []string) error {
	debug := viper.GetBool("debug")
	vaultDir := viper.GetString("vault-dir")
	if debug {
		logger.SetLevel(logrus.DebugLevel)
	} else {
		logger.SetLevel(logrus.InfoLevel)
	}

	if vaultDir == "" {
		return errors.New("vault-dir flag is not given and SAO_DIR env variable is not set")
	}

	return nil
}
