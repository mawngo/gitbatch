package clone

import (
	"encoding/json"
	"fmt"
	"gitbatch/pkg/util"
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

var cloneGitLabCommand = func() cobra.Command {
	var command = cobra.Command{
		Use:   "gitlab [group id] [dir]",
		Short: "Clone all project in gitlab group",
		Args:  cobra.RangeArgs(1, 2),
		Run:   executeCloneGitLabCommand,
	}
	bindCloneGitlabFlags(&command)
	return command
}()

func bindCloneGitlabFlags(command *cobra.Command) {
	command.Flags().String("token", "", "GitLab API Key")
	command.Flags().String("host", "", "GitLab host")
	command.Flags().Int("max", 10000, "Maximum number of projects to clone")
	lo.Must0(viper.BindPFlag("gitlab.token", command.Flag("token")))
	lo.Must0(viper.BindPFlag("gitlab.host", command.Flag("host")))
}

func executeCloneGitLabCommand(cmd *cobra.Command, args []string) {
	host := util.AskConfig("gitlab.host", "Enter Gitlab host")
	token := util.AskConfig("gitlab.token", "Enter Gitlab token")
	size := lo.Must(cmd.Flags().GetInt("max"))

	res := lo.Must(http.Get(fmt.Sprintf("https://%s/api/v4/groups/%s?private_token=%s&per_page=%d",
		host, args[0], token, size)))

	if res.StatusCode != 200 {
		fmt.Printf("Error status code %d when fetch project infomation\n", res.StatusCode)
		fmt.Printf("API: https://<%s>/api/v4/groups/<%s>?per_page=<%d>\n", host, args[0], size)
		return
	}
	var group map[string]any
	cobra.CheckErr(json.Unmarshal(lo.Must(io.ReadAll(res.Body)), &group))
	if message, ok := group["message"]; ok {
		println(message)
		return
	}

	workingDir := "."
	if len(args) > 1 {
		workingDir = args[1]
	}

	projects := group["projects"].([]any)
	fmt.Printf("Total %d projects\n", len(projects))

	ssh := lo.Must(cmd.Flags().GetString("user")) == "@ssh"
	util.SplitParallel(viper.GetInt("parallel"), projects, func(p any) {
		project := p.(map[string]any)
		url := lo.Ternary(ssh, project["ssh_url_to_repo"], project["http_url_to_repo"]).(string)
		name := project["path"].(string)

		dir := path.Join(workingDir, name)
		if _, err := os.Stat(dir); err == nil {
			color.Blue("Skip %s already exist.\n", name)
			return
		}

		color.Green("Cloning %s...\n", name)
		opt := &git.CloneOptions{
			URL: url,
		}
		if viper.GetInt("parallel") == 1 {
			opt.Progress = os.Stdout
		}
		_, err := git.PlainClone(dir, false, opt)
		if err != nil {
			color.Red("Error %s\n", err)
		}
	})
}
