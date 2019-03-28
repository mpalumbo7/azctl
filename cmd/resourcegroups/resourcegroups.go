package resourcegroups

import (
	"fmt"

	"github.com/spf13/cobra"
)

// NewCmdResourceGroups represents the ad command
func NewCmdResourceGroups() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "resourcegroups",
		Aliases: []string{"rg", "groups"},
		Short:   "Manage resource groups",
		Long:    ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Resource group stuff...")
		},
	}

	// Add subcommands
	cmd.AddCommand(NewCmdGetResourceGroups())

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//cmd.PersistentFlags().String("group", "", "The resource group to operate on")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// adCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	return cmd
}
