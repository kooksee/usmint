package main

import (
	"os"

	"github.com/kooksee/usmint/cmd"
	"github.com/tendermint/tmlibs/cli"
)

func main() {
	rootCmd := cmd.RootCmd
	rootCmd.AddCommand(
		cmd.GenValidatorCmd,
		cmd.InitFilesCmd,
		cmd.ProbeUpnpCmd,
		cmd.LiteCmd,
		cmd.ReplayCmd,
		cmd.ReplayConsoleCmd,
		cmd.ResetAllCmd,
		cmd.ResetPrivValidatorCmd,
		cmd.ShowValidatorCmd,
		cmd.TestnetFilesCmd,
		cmd.ShowNodeIDCmd,
		cmd.GenNodeKeyCmd,
		cmd.VersionCmd,

		cmd.NewRunNodeCmd(),
	)

	err := cli.PrepareBaseCmd(rootCmd, "U", os.ExpandEnv("$PWD/kdata")).Execute()
	if err != nil {
		panic(err)
	}
}
