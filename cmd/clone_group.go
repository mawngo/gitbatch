package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"path"
)

var cloneGroupCmd = func() cobra.Command {
	var command = cobra.Command{
		Use:     "clonegroup [gitlab group id]",
		Aliases: []string{"cg"},
		Short:   "Clone all project in group",
		Args:    cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			host := askConfig("gitlab.host", "Enter Gitlab host")
			token := askConfig("gitlab.token", "Enter Gitlab token")
			size := lo.Must(cmd.Flags().GetInt("max"))

			res := lo.Must(http.Get(fmt.Sprintf("https://%s/api/v4/groups/%s?private_token=%s&per_page=%d",
				host, args[0], token, size)))

			var group map[string]any
			cobra.CheckErr(json.Unmarshal(lo.Must(io.ReadAll(res.Body)), &group))
			if message, ok := group["message"]; ok {
				println(message)
			}

			projects := group["projects"].([]any)
			fmt.Printf("Total %d projects\n", len(projects))
			for _, project := range projects {
				project := project.(map[string]any)
				ssh := lo.Must(cmd.Flags().GetBool("ssh"))
				url := lo.Ternary(ssh, project["ssh_url_to_repo"], project["http_url_to_repo"]).(string)
				name := project["path"].(string)

				dir := path.Join(lo.Must(os.Getwd()), name)
				if _, err := os.Stat(dir); err == nil {
					color.Blue("Skip %s already exist.\n", name)
					continue
				}

				color.Green("Cloning %s...\n", name)
				_, err := git.PlainClone(dir, false, &git.CloneOptions{
					URL:      url,
					Progress: os.Stdout,
				})
				if err != nil {
					color.Red("Error %s\n", err)
				}
			}
		},
	}

	command.Flags().String("token", "", "GitLab API Key")
	command.Flags().String("host", "", "GitLab host")
	command.Flags().Bool("ssh", true, "Clone using ssh")
	command.Flags().Int("max", 10000, "Maximum number of projects to clone")
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
