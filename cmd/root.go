package cmd

import (
	"github.com/lana-toolbox/gitbatch/cmd/clone"
	"github.com/lana-toolbox/gitbatch/internal/util"
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
	rootCmd.AddCommand(&pullAllCmd)
	rootCmd.AddCommand(&pushAllCmd)
	rootCmd.PersistentFlags().Int("parallel", 32, "Maximum parallel for each commands")
	rootCmd.PersistentFlags().String("token", "", "Host token")
	rootCmd.PersistentFlags().StringP("mode", "m", "gitlab", "Host mode")
	rootCmd.PersistentFlags().StringP("user", "u", "@ssh", "Auth user name [<user>, @ssh]")
	cobra.EnableCommandSorting = false

	cobra.CheckErr(viper.BindPFlag("parallel", rootCmd.PersistentFlags().Lookup("parallel")))
	cobra.CheckErr(viper.BindPFlag("mode", rootCmd.PersistentFlags().Lookup("mode")))
	cobra.CheckErr(viper.BindPFlag("user", rootCmd.PersistentFlags().Lookup("user")))
}

var rootCmd = &cobra.Command{
	Use:   "github.com/lana-toolbox/gitbatch",
	Short: "Git batch operations",
	Long:  "Apply git command to all sub folder",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if t, err := cmd.Flags().GetString("token"); err == nil && t != "" {
			viper.Set(util.GetMode()+".token", t)
		}
		return nil
	},
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
