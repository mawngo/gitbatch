package cmd

import (
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
)

func init() {
	cobra.OnInitialize(configure)
	rootCmd.AddCommand(&cloneGroupCmd)
	rootCmd.AddCommand(&fetchAllCmd)
}

var rootCmd = &cobra.Command{
	Use:   "gitbatch",
	Short: "Git batch operations",
	Long:  "Apply git command to all sub folder",
}

func configure() {
	home, err := os.UserHomeDir()
	cobra.CheckErr(err)
	viper.AddConfigPath(path.Join(home, ".config"))
	viper.SetConfigType("json")
	viper.SetConfigName("gb")
	viper.AutomaticEnv()
	_ = viper.SafeWriteConfig()
	cobra.CheckErr(viper.ReadInConfig())
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
