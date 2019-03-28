package resourcegroups

import (
	"fmt"
	"os"

	"github.com/spf13/cast"

	armpkg "github.com/mpalumbo7/azctl/pkg/arm"
	clientpkg "github.com/mpalumbo7/azctl/pkg/client"
	configpkg "github.com/mpalumbo7/azctl/pkg/config"
	printpkg "github.com/mpalumbo7/azctl/pkg/print"

	"github.com/spf13/cobra"
)

// NewCmdGetResourceGroups gets a list of resource groups
func NewCmdGetResourceGroups() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get [command goes here]",
		Short: "Get resource groups",
		Long:  ``,
		Args:  cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			execute(cmd, args)
		},
	}

	// Add subcommands

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	//cmd.PersistentFlags().String("subscriptionId", "", "The subscription for this ")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// adCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	return cmd
}

// columns to be printed
var (
	// left out `id` column, as it is too long and wraps on screen
	columns = []string{"name", "location"}
)

func execute(cmd *cobra.Command, args []string) {
	// global variable context
	// TODO: Make this load --context flag
	cxt, err := configpkg.GetCurrentContext()
	if err != nil {
		fmt.Printf("Error retrieving context: %s\n", err)
		os.Exit(1)
	}

	clt, err := clientpkg.NewClient(cxt)
	if err != nil {
		fmt.Printf("Error initiating ARM client: %s\n", err)
		os.Exit(1)
	}

	groups, err := armpkg.GetResourceGroups(clt, cxt.SubscriptionID)
	if err != nil {
		fmt.Printf("Error retrieving resource groups: %s\n", err)
		os.Exit(1)
	}

	rawList := groups["value"].([]interface{})
	list := make([]map[string]interface{}, len(rawList))
	for i, grp := range rawList {
		grpmap := cast.ToStringMap(grp)
		list[i] = grpmap
	}

	out := printpkg.GetNewTabWriter(os.Stdout)
	defer out.Flush()
	printpkg.WriteTable(columns, list, out, printpkg.NewTableOptions())
}
