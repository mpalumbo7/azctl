package config

import (
	"github.com/spf13/cobra"
)

// NewCmdConfig represents the ad command
func NewCmdConfig() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "config",
		Short: "Manage the configuration file",
		Long: `Configuration options follow this order:
	1. If the --config flag is used, then only the options from that file are loaded.
	2. If the $AZCTLCONFIG environment variable is set, the value will be used to load the configuration file.
	3. Otherwise, the config file at ${HOME}/.azctl is loaded.`,
	}

	// Add subcommands
	cmd.AddCommand(NewCmdContext())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// adCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// adCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	return cmd
}
