package cmd

import (
	"fmt"
	"gitbatch/pkg/util"
	"github.com/fatih/color"
	"github.com/go-git/go-git/v5"
	"github.com/samber/lo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"os"
	"path"
	"runtime"
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
			parallel := viper.GetInt("parallel")
			if parallel > runtime.NumCPU() { // Too many connection may starve the SSH_AUTH_SOCK
				parallel = runtime.NumCPU()
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

				err = repo.Fetch(&git.FetchOptions{})
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
