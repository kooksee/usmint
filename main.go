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
	"github.com/kooksee/usmint/cmn"
	"github.com/kooksee/usmint/cmd"
)

func DefaultNewNode(config *config.Config, logger log.Logger) (*node.Node, error) {
	// init cmn
	cmn.InitLog(logger)

	// init config
	cmn.InitCfg(config)

	return node.NewNode(config,
		privval.LoadOrGenFilePV(config.PrivValidatorFile()),
		proxy.NewLocalClientCreator(app.New()),
		node.DefaultGenesisDocProviderFunc(config),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(config.Instrumentation),
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
		cmds.VersionCmd,
	)

	// Create & start node
	rootCmd.AddCommand(commands.NewRunNodeCmd(DefaultNewNode))

	if err := cli.PrepareBaseCmd(rootCmd, "K", os.ExpandEnv("$PWD/kdata")).Execute(); err != nil {
		panic(err)
	}
}
