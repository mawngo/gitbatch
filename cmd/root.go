package cmd

import (
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
)

func init() {
	cobra.OnInitialize(configure)
	rootCmd.AddCommand(&cloneGroupCmd)
	rootCmd.AddCommand(&fetchAllCmd)
	rootCmd.PersistentFlags().Int("parallel", 32, "Maximum parallel for each commands")
	cobra.CheckErr(viper.BindPFlag("parallel", rootCmd.Flag("parallel")))
}

var rootCmd = &cobra.Command{
	Use:   "gitbatch",
	Short: "Git batch operations",
	Long:  "Apply git command to all sub folder",
}

func configure() {
	viper.AddConfigPath(path.Join(lo.Must(os.UserHomeDir()), ".config"))
	viper.SetConfigType("json")
	viper.SetConfigName("gb")
	viper.AutomaticEnv()

	if err := viper.SafeWriteConfig(); err != nil {
		cobra.CheckErr(viper.ReadInConfig())
		cobra.CheckErr(viper.WriteConfig())
	}
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
