package cmd

import (
	"fmt"
	"github.com/mawngo/gitbatch/cmd/clone"
	"github.com/mawngo/gitbatch/internal/util"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
)

func init() {
	cobra.EnableCommandSorting = false
}

type CLI struct {
	command *cobra.Command
}

// NewCLI create new CLI instance and setup application config.
func NewCLI() *CLI {
	command := cobra.Command{
		Use:   "gitbatch",
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
	cobra.OnInitialize(configure)
	command.AddCommand(lo.ToPtr(clone.NewCloneCommand()))
	command.AddCommand(&fetchAllCmd)
	command.AddCommand(&pullAllCmd)
	command.AddCommand(&pushAllCmd)
	command.PersistentFlags().Int("parallel", 32, "Maximum parallel for each commands")
	command.PersistentFlags().String("token", "", "Host token")
	command.PersistentFlags().StringP("mode", "m", "gitlab", "Host mode")
	command.PersistentFlags().StringP("user", "u", "@ssh", "Auth user name [<user>, @ssh]")
	command.Flags().SortFlags = false

	cobra.CheckErr(viper.BindPFlag("parallel", command.PersistentFlags().Lookup("parallel")))
	cobra.CheckErr(viper.BindPFlag("mode", command.PersistentFlags().Lookup("mode")))
	cobra.CheckErr(viper.BindPFlag("user", command.PersistentFlags().Lookup("user")))
	return &CLI{&command}
}

func configure() {
	viper.AddConfigPath(path.Join(lo.Must(os.UserHomeDir()), ".config"))
	viper.SetConfigType("json")
	viper.SetConfigName("gb")

	if err := viper.SafeWriteConfig(); err != nil {
		cobra.CheckErr(viper.ReadInConfig())
	}
}

func (cli *CLI) Execute() {
	if err := cli.command.Execute(); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, err)
	}
}
