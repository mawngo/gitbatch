package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/mawngo/gitbatch/internal/util"
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
			RunWithGit(cmd, args, func(cmd *cobra.Command, file os.DirEntry, repo *git.Repository, auth transport.AuthMethod) {
				opt := git.FetchOptions{}
				if auth != nil {
					opt.Auth = auth
				}
				err := repo.Fetch(&opt)
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

func RunWithGit(cmd *cobra.Command, args []string, handler func(cmd *cobra.Command, file os.DirEntry, repo *git.Repository, auth transport.AuthMethod)) {
	workingDir := "."
	if len(args) > 0 {
		workingDir = args[0]
	}

	files := lo.Must(os.ReadDir(workingDir))
	parallel := viper.GetInt("parallel")

	if util.IsSSH() && parallel > 8 { // Too many connection may starve the SSH_AUTH_SOCK
		println("Reduce number of parallel to avoid SSH_AUTH_SOCK error")
		parallel = 8
	}

	auth := util.AskAuth()
	util.SplitParallel(parallel, files, func(file os.DirEntry) {
		if !file.IsDir() {
			return
		}
		repo, err := git.PlainOpen(path.Join(workingDir, file.Name()))

		if err != nil {
			if errors.Is(err, git.ErrRepositoryNotExists) {
				return
			}

			color.Yellow("Repo %s: %s", file.Name(), err)
			return
		}
		handler(cmd, file, repo, auth)
	})
}
