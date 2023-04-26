package cmd

import (
	"fmt"
	"gitbatch/pkg/util"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
)

var fetchAllCmd = func() cobra.Command {
	var command = cobra.Command{
		Use:   "fetch [dir]",
		Short: "Fetch all project in directory",
		Args:  cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			workingDir := "."
			if len(args) > 0 {
				workingDir = args[0]
			}

			files := lo.Must(os.ReadDir(workingDir))
			parallel := viper.GetInt("parallel")
			user := lo.Must(cmd.Flags().GetString("user"))
			ssh := user == "@ssh"

			if ssh && parallel > 8 { // Too many connection may starve the SSH_AUTH_SOCK
				println("Reduce number of parallel to avoid SSH_AUTH_SOCK error")
				parallel = 8
			}

			password := ""
			if !ssh {
				password = lo.Ternary(ssh, "", util.AskPassword("Enter gitlab password (this won't be saved)"))
			}

			util.SplitParallel(parallel, files, func(file os.DirEntry) {
				if !file.IsDir() {
					return
				}
				repo, err := git.PlainOpen(path.Join(workingDir, file.Name()))

				if err != nil {
					color.Yellow("Fetch %s: %s", file.Name(), err)
					return
				}

				opt := git.FetchOptions{}
				if !ssh {
					opt.Auth = &http.BasicAuth{Username: user, Password: password}
				}
				err = repo.Fetch(&opt)
				if err != nil {
					fmt.Printf("Fetch %s: %s\n", file.Name(), err)
					return
				}
				color.Green("Fetch: %s", file.Name())
			})
		},
	}
	return command
}()
