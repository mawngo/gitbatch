package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/spf13/cobra"
	"os"
)

var pushAllCmd = func() cobra.Command {
	var command = cobra.Command{
		Use:   "push [dir]",
		Short: "Push all project in directory",
		Args:  cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			RunWithGit(cmd, args, func(_ *cobra.Command, file os.DirEntry, repo *git.Repository, auth transport.AuthMethod) {
				opt := git.PushOptions{}
				if auth != nil {
					opt.Auth = auth
				}
				err := repo.Push(&opt)
				if err != nil {
					if errors.Is(err, git.NoErrAlreadyUpToDate) {
						return
					}
					fmt.Printf("Push %s: %s\n", file.Name(), err)
					return
				}
				color.Green("Push: %s", file.Name())
			})
		},
	}

	command.Flags().StringP("branch", "b", "master", "Specify branch to push")

	return command
}()
