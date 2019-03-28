package cmd

import (
	"fmt"
	"os"

	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	ad "github.com/mpalumbo7/azctl/cmd/ad"
	config "github.com/mpalumbo7/azctl/cmd/config"
	"github.com/mpalumbo7/azctl/cmd/resourcegroups"
)

const (
	defaultCfgFile     string = ".azctl"
	envVariableCfgFile string = "AZCTLCONFIG"
)

var (
	cfgFile string
	context string
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "azctl",
	Short: "A utility for managing Azure via REST APIs",
	Long:  `A command line utility for interfacing with Azure via REST APIs.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// sub-commands
	rootCmd.AddCommand(ad.NewCmdAd())
	rootCmd.AddCommand(config.NewCmdConfig())
	rootCmd.AddCommand(resourcegroups.NewCmdResourceGroups())

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.azctl)")
	rootCmd.PersistentFlags().StringVar(&context, "context", "", "the context to use in commands (default is defined in config file)")
	viper.BindPFlag("context", rootCmd.PersistentFlags().Lookup("context"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	//rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Check environment variable
		if viper.GetString(envVariableCfgFile) != "" {
			// Use config file from the flag.
			viper.SetConfigFile(viper.GetString(envVariableCfgFile))
		} else {
			// Find home directory.
			home, err := homedir.Dir()
			if err != nil {
				fmt.Println(err)
				os.Exit(1)
			}

			// Search config in home directory with name ".azctl" (without extension).
			viper.AddConfigPath(home)
			viper.SetConfigName(defaultCfgFile)
		}
	}

	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	err := viper.ReadInConfig()
	if err != nil {
		// Handle errors reading the config file
		fmt.Println(err)
		os.Exit(1)
	}
}
