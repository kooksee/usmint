package main

import (
	"os"
	"github.com/tendermint/tendermint/libs/cli"

	"github.com/tendermint/tendermint/cmd/tendermint/commands"
	"github.com/tendermint/tendermint/privval"
	"github.com/tendermint/tendermint/config"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/proxy"
	"github.com/kooksee/usmint/app"
	"github.com/kooksee/usmint/cmn"
	"github.com/kooksee/usmint/cmd"
	"github.com/kooksee/usmint/reactors"
	"time"
	"github.com/tendermint/tendermint/p2p"
	"fmt"
	"github.com/tendermint/tendermint/node"
	"bytes"
)

func ff(s *p2p.Switch, kr *reactors.KReactor, logger log.Logger, node2 *node.Node) {
	for {

		if !bytes.Contains(node2.NodeInfo().Channels, []byte{0x60}) {
			nf := node2.Switch().NodeInfo()
			nf.Channels = append(nf.Channels, kr.ChId)
			node2.Switch().SetNodeInfo(nf)
		}

		fmt.Println(node2.NodeInfo().Channels.Bytes())
		fmt.Println(node2.Switch().NumPeers())
		node2.Switch().Broadcast(kr.ChId, []byte("hello kr"))
		logger.Error("test sent")
		time.Sleep(time.Second * 2)
	}
}

func DefaultNewNode(config *config.Config, logger log.Logger) (*node.Node, error) {
	// init cmn
	cmn.InitLog(logger)

	// init config
	cmn.InitCfg(config)

	n, err := node.NewNode(config,
		privval.LoadOrGenFilePV(config.PrivValidatorFile()),
		proxy.NewLocalClientCreator(app.New()),
		node.DefaultGenesisDocProviderFunc(config),
		node.DefaultDBProvider,
		node.DefaultMetricsProvider(config.Instrumentation),
		logger,
	)

	// 获得node
	cmn.InitNode(n)

	kr := reactors.NewKReactor()
	kr.SetLogger(logger.With("module", "kr"))
	n.Switch().AddReactor(kr.Name, kr)
	fmt.Println(kr.ChId)

	go ff(n.Switch(), kr, logger, n)
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

	if err := cli.PrepareBaseCmd(rootCmd, "K", os.ExpandEnv("$PWD/kdata")).Execute(); err != nil {
		panic(err)
	}
}
