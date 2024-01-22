package cmd

import (
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/spf13/cobra"
	"os"
)

var pullAllCmd = func() cobra.Command {
	var command = cobra.Command{
		Use:   "pull [dir]",
		Short: "Pull all project in directory",
		Args:  cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			RunWithGitWorkTree(cmd, args, func(cmd *cobra.Command, file os.DirEntry, worktree *git.Worktree, auth *http.BasicAuth) {
				opt := git.PullOptions{}
				if auth != nil {
					opt.Auth = auth
				}
				err := worktree.Pull(&opt)
				if err != nil {
					if errors.Is(err, git.NoErrAlreadyUpToDate) {
						return
					}
					fmt.Printf("Pull %s: %s\n", file.Name(), err)
					return
				}
				color.Green("Pull: %s", file.Name())
			})
		},
	}

	command.Flags().StringP("branch", "b", "master", "Specify branch to pull")

	return command
}()

func RunWithGitWorkTree(cmd *cobra.Command, args []string, handler func(cmd *cobra.Command, file os.DirEntry, worktree *git.Worktree, auth *http.BasicAuth)) {
	RunWithGit(cmd, args, func(cmd *cobra.Command, file os.DirEntry, repo *git.Repository, auth *http.BasicAuth) {
		worktree, err := repo.Worktree()
		if err != nil {
			color.Yellow("Worktree %s: %s", file.Name(), err)
			return
		}

		branchName, err := cmd.Flags().GetString("branch")
		if err == nil && branchName != "" {
			err = worktree.Checkout(&git.CheckoutOptions{Branch: plumbing.NewBranchReferenceName(branchName)})
			if err != nil {
				color.Yellow("Checkout %s %s: %s", branchName, file.Name(), err)
				return
			}
		}

		handler(cmd, file, worktree, auth)
	})
}
