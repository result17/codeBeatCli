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

	flags.BoolP("version", "v", false, "Print CodeBeatCli version, and exit.")
	flags.Bool("dlog", false, "set debugger logger level")

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
