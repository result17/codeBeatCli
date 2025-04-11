package cmd

import (
	"errors"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/result17/codeBeatCli/pkg/exitcode"
)

func NewCMD() *cobra.Command {
	v := viper.New()
	cmd := &cobra.Command{
		Use:   "CodeBeatCli",
		Short: "Command line interface used by CodeBeat vscode plugin.",
		RunE: func(cmd *cobra.Command, _ []string) error {
			if err := RunE(cmd, v); err != nil {
				var errexitcode exitcode.Err

				if errors.As(err, &errexitcode) {
					os.Exit(errexitcode.Code)
				}

				os.Exit(exitcode.ErrGeneric)
			}

			os.Exit(exitcode.Success)

			return nil
		},
	}
	setFlags(cmd, v)
	return cmd
}

func setFlags(cmd *cobra.Command, v *viper.Viper) {
	flags := cmd.Flags()
	flags.String("alternate-project", "", "Alternate project name.(Optional)")
	flags.String("config", "", "Plugin config file.(Optional)")
	flags.BoolP("version", "v", false, "Print CodeBeatCli version, and exit.")
	flags.Bool("dlog", false, "Set debugger logger level")
	flags.Int("cursorpos", 0, "Cursor position in the current file for the heartbeat.(Optional)")
	flags.Int("lineno", 0, "Current line number int the file.")
	flags.Int(
		"lines-in-file",
		0,
		"The total line count of file for the heartbeat.",
	)
	flags.String(
		"entity",
		"",
		"Absolute path to file for the heartbeat.",
	)
	flags.String("log-file", "", "Plugin log file absolute path.(Optional)")
	flags.String("project-folder", "", "Absolute path to project folder.(Optional)")
	flags.String("log-filer", "", "Absolute path to plugin log file.(Optional)")
	flags.Int64("time", 0, "Unix epoch timeStamp. Uses current time by default.")
	flags.String("plugin", "", "Text editor plugin name and version")
	err := v.BindPFlags(flags)
	if err != nil {
		log.Fatalf("failed to bind cobra flags to viper: %s", err)
	}
}

func Execute() {
	if err := NewCMD().Execute(); err != nil {
		log.Fatalf("failed to run codeBeatCli: %s", err)
	}
}
