package config

import (
	"encoding/json"
	"fmt"
	"os"

	configpkg "github.com/mpalumbo7/azctl/pkg/config"
	"github.com/mpalumbo7/azctl/pkg/print"
	printpkg "github.com/mpalumbo7/azctl/pkg/print"

	"github.com/spf13/cobra"
)

// NewCmdGetCurrentContext intializes a new command
func NewCmdGetCurrentContext() *cobra.Command {
	// cmd represents the aduser command
	cmd := &cobra.Command{
		Use:   "get [context name goes here]",
		Short: "Gets a list of all context, or a specific context",
		Long:  ``,
		Run: func(cmd *cobra.Command, args []string) {
			execute(cmd, args)
		},
	}
	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	return cmd
}

// columns to be printed
var (
	columns = []string{"Name", "SubscriptionID", "User"}
)

func execute(cmd *cobra.Command, args []string) {
	cfg, err := configpkg.GetConfig()
	if err != nil {
		fmt.Printf("Error retrieving configuration: %s\n", err)
		os.Exit(1)
	}

	// show selected contexts

	rows := []map[string]interface{}{}
	for _, cxt := range cfg.Contexts {
		// marshal struct to json, then back to mapstructure
		var cxtmap map[string]interface{}
		cxtjson, _ := json.Marshal(cxt)
		json.Unmarshal(cxtjson, &cxtmap)

		if len(args) != 0 {
			for _, arg := range args {
				if arg == cxt.Name {
					rows = append(rows, cxtmap)
					break
				}
			}
		} else {
			rows = append(rows, cxtmap)
		}
	}

	out := printpkg.GetNewTabWriter(os.Stdout)
	defer out.Flush()
	print.WriteTable(columns, rows, out, print.NewTableOptions())
}
