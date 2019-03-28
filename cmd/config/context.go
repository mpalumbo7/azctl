package config

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewCmdContext represents the ad command
func NewCmdContext() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "context",
		Short: "Manage the contexts in the configuration file",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Configuration stuff...")
		},
	}

	// Add subcommands
	cmd.AddCommand(NewCmdGetCurrentContext())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// adCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// adCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	return cmd
}
