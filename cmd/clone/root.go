package clone

import (
	"github.com/spf13/cobra"
)

func NewCloneCommand() cobra.Command {
	var command = cobra.Command{
		Use:     "clone [gitlab group id] [dir]",
		Aliases: []string{"cg"},
		Short:   "Clone all project in group (alias for 'clone gitlab')",
		Args:    cobra.RangeArgs(1, 2),
		Run:     executeCloneGitLabCommand,
	}

	bindCloneGitlabFlags(&command)
	command.AddCommand(&cloneGitLabCommand)
	return command
}
