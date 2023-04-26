package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cloneGroupCmd = func() cobra.Command {
	var command = cobra.Command{
		Use:     "clonegroup [gitlab group id]",
		Aliases: []string{"cg"},
		Short:   "Clone all project in group",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			var host = askConfig("gitlab.host", "Enter Gitlab host")
			var token = askConfig("gitlab.token", "Enter Gitlab token")
			var url = fmt.Sprintf("https://%s/api/v4/groups/%s", host, args[0])

		},
	}

	command.Flags().String("token", "", "GitLab API Key")
	command.Flags().String("host", "", "GitLab host")
	lo.Must0(viper.BindPFlag("gitlab.token", command.Flag("token")))
	lo.Must0(viper.BindPFlag("gitlab.host", command.Flag("host")))
	return command
}()

func askConfig(key, message string) string {
	var value = viper.GetString(key)
	if value == "" {
		err := survey.AskOne(&survey.Input{Message: message}, &value, survey.WithValidator(survey.Required))
		cobra.CheckErr(err)
		viper.Set(key, value)
		lo.Must0(viper.WriteConfig())
	}
	return value
}
