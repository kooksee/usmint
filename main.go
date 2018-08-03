package main

import (
	"os"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/tendermint/tendermint/cmd/tendermint/commands"
	"github.com/tendermint/tendermint/node"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/proxy"
	"github.com/kooksee/usmint/app"
)

func DefaultNewNode(config *config.Config, logger log.Logger) (*node.Node, error) {
	return node.NewNode(config,
		privval.LoadOrGenFilePV(config.PrivValidatorFile()),
		proxy.NewLocalClientCreator(app.New(logger)),
		node.DefaultGenesisDocProviderFunc(config),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider,
		logger,
	)
}

func main() {
	rootCmd := commands.RootCmd
	rootCmd.AddCommand(
		commands.GenValidatorCmd,
		commands.InitFilesCmd,
		commands.ProbeUpnpCmd,
		commands.LiteCmd,
		commands.ReplayCmd,
		commands.ReplayConsoleCmd,
		commands.ResetAllCmd,
		commands.ResetPrivValidatorCmd,
		commands.ShowValidatorCmd,
		commands.TestnetFilesCmd,
		commands.ShowNodeIDCmd,
		commands.GenNodeKeyCmd,
		commands.VersionCmd,
	)

	// DefaultNewNode function
	nodeFunc := node.DefaultNewNode

	// Create & start node
	rootCmd.AddCommand(commands.NewRunNodeCmd(nodeFunc))

	if err := cli.PrepareBaseCmd(rootCmd, "K", os.ExpandEnv("$PWD/kdata")).Execute(); err != nil {
		panic(err)
	}
}
