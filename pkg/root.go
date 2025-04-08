package pkg

import (
	"github.com/spf13/cobra"
)

func NewCMD() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "CodeBeatCli",
		Short: "Command line interface used by CodeBeat vscode plugin.",
	}
	return cmd
}
