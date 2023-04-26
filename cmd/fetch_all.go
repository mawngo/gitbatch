package cmd

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"os"
	"path"
)

var fetchAllCmd = func() cobra.Command {
	var command = cobra.Command{
		Use:   "fetch [dir]",
		Short: "Fetch all project in directory",
		Long:  "Fetch all project in directory. Only support ssh cloned project.",
		Args:  cobra.RangeArgs(0, 1),
		Run: func(cmd *cobra.Command, args []string) {
			workingDir := "."
			if len(args) > 0 {
				workingDir = args[0]
			}

			files := lo.Must(os.ReadDir(workingDir))
			for _, file := range files {
				if !file.IsDir() {
					continue
				}
				repo, err := git.PlainOpen(path.Join(workingDir, file.Name()))

				if err != nil {
					color.Yellow("Fetch %s: %s", file.Name(), err)
					continue
				}

				err = repo.Fetch(&git.FetchOptions{})
				if err != nil {
					fmt.Printf("Fetch %s: %s\n", file.Name(), err)
					continue
				}
				color.Green("Fetch: %s", file.Name())
			}
		},
	}
	return command
}()
