package cmd

import (
	"gitbatch/cmd/clone"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
)

func init() {
	cobra.OnInitialize(configure)
	rootCmd.AddCommand(lo.ToPtr(clone.NewCloneCommand()))
	rootCmd.AddCommand(&fetchAllCmd)
	rootCmd.PersistentFlags().Int("parallel", 32, "Maximum parallel for each commands")
	rootCmd.PersistentFlags().StringP("user", "u", "@ssh", "Auth user name (set to @ssh to auth using ssh)")
	cobra.EnableCommandSorting = false

	cobra.CheckErr(viper.BindPFlag("parallel", rootCmd.PersistentFlags().Lookup("parallel")))
	cobra.CheckErr(viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user")))
}

var rootCmd = &cobra.Command{
	Use:   "gitbatch",
	Short: "Git batch operations",
	Long:  "Apply git command to all sub folder",
	PersistentPostRun: func(cmd *cobra.Command, args []string) {
		cobra.CheckErr(viper.WriteConfig())
	},
}

func configure() {
	viper.AddConfigPath(path.Join(lo.Must(os.UserHomeDir()), ".config"))
	viper.SetConfigType("json")
	viper.SetConfigName("gb")

	if err := viper.SafeWriteConfig(); err != nil {
		cobra.CheckErr(viper.ReadInConfig())
	}
}

func Execute() {
	cobra.CheckErr(rootCmd.Execute())
}
