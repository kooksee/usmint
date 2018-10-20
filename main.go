package main

import (
	"os"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/tendermint/tendermint/cmd/tendermint/commands"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/kooksee/usmint/app"
	"github.com/kooksee/usmint/cmn"
	"github.com/kooksee/usmint/cmd"
	"github.com/tendermint/tendermint/proxy"
	"path/filepath"
	"github.com/kooksee/usmint/node"
	"github.com/kooksee/usmint/usdb"
)

func DefaultNewNode(config *config.Config, logger log.Logger) (*node.Node, error) {
	h, err := os.Hostname()
	if err != nil {
		cmn.MustNotErr("init hostname error", err)
	}

	logger = logger.With("service", "usmint", "hostname", h)

	// init cmn
	cmn.InitLog(logger)

	// init config
	cmn.InitCfg(config)

	// 初始化db
	cmn.InitAppDb(filepath.Join(config.DBDir(), "mint_app.db"))

	usdb.Name = "txs"
	usdb.Init(logger)

	n, err := node.NewNode(config,
		privval.LoadOrGenFilePV(config.PrivValidatorFile()),
		proxy.NewLocalClientCreator(app.New(logger)),
		node.DefaultGenesisDocProviderFunc(config),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(config.Instrumentation),
		logger,
	)

	// 获得node
	node.InitNode(n)

	return n, err
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

	if err := cli.PrepareBaseCmd(rootCmd, "MINT", os.ExpandEnv("$PWD/kdata")).Execute(); err != nil {
		panic(err)
	}
}
