package cmd

import (
	"fmt"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

var (
	// Used for flags.
	cfgFile string

	rootCmd = &cobra.Command{
		Use:   "simplify-terminations-fip-data",
		Short: "Tool to collect sector data for FIP proposal",
		Run: func(cmd *cobra.Command, args []string) {
		},
	}
)

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", "mainnet.env", "config file (default is ./mainnet.env)")

	viper.BindEnv("lotus_addr")
	viper.BindEnv("lotus_token")

	rootCmd.Flags().VisitAll(func(f *pflag.Flag) {
		// Environment variables can't have dashes in them, so bind them to their equivalent
		// keys with underscores, e.g. --favorite-color to STING_FAVORITE_COLOR
		if strings.Contains(f.Name, "-") {
			envVarSuffix := strings.ToUpper(strings.ReplaceAll(f.Name, "-", "_"))
			viper.BindEnv(f.Name, envVarSuffix)
		}

		// Apply the viper config value to the flag when the flag is not set and viper has a value
		if !f.Changed && viper.IsSet(f.Name) {
			val := viper.Get(f.Name)
			rootCmd.Flags().Set(f.Name, fmt.Sprintf("%v", val))
		}
	})
}

func initConfig() {
	viper.SetConfigName(cfgFile)
	viper.AddConfigPath(".")
	viper.SetConfigType("env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		fmt.Println("No config file found", err)
	}
}
