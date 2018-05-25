package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/tendermint/tendermint/proxy"
	"github.com/kooksee/kchain/cfg"
	cmn "github.com/tendermint/tmlibs/common"
	kn "github.com/tendermint/tendermint/node"
	"github.com/kooksee/kchain/app"
	pvm "github.com/tendermint/tendermint/types/priv_validator"
	"github.com/gin-gonic/gin/json"
	"github.com/kooksee/kchain/reactors"
	"time"
)

var kcfg = cfg.GetConfig()

// AddNodeFlags exposes some common configuration options on the command-line
// These are exposed for convenience of commands embedding a tendermint node
func AddNodeFlags(cmd *cobra.Command) *cobra.Command {
	// bind flags
	cmd.Flags().String("moniker", config.Moniker, "Node Name")

	// priv val flags
	cmd.Flags().String("priv_validator_laddr", config.PrivValidatorListenAddr, "Socket address to listen on for connections from external priv_validator process")

	// node flags
	cmd.Flags().Bool("fast_sync", config.FastSync, "Fast blockchain syncing")

	// abci flags
	cmd.Flags().String("proxy_app", config.ProxyApp, "Proxy app address, or 'nilapp' or 'kvstore' for local testing.")
	cmd.Flags().String("abci", config.ABCI, "Specify abci transport (socket | grpc)")

	// rpc flags
	cmd.Flags().String("rpc.laddr", config.RPC.ListenAddress, "RPC listen address. Port required")
	cmd.Flags().String("rpc.grpc_laddr", config.RPC.GRPCListenAddress, "GRPC listen address (BroadcastTx only). Port required")
	cmd.Flags().Bool("rpc.unsafe", config.RPC.Unsafe, "Enabled unsafe rpc methods")

	// p2p flags
	cmd.Flags().String("p2p.laddr", config.P2P.ListenAddress, "Node listen address. (0.0.0.0:0 means any interface, any port)")
	cmd.Flags().String("p2p.seeds", config.P2P.Seeds, "Comma-delimited ID@host:port seed nodes")
	cmd.Flags().String("p2p.persistent_peers", config.P2P.PersistentPeers, "Comma-delimited ID@host:port persistent peers")
	cmd.Flags().Bool("p2p.skip_upnp", config.P2P.SkipUPNP, "Skip UPNP configuration")
	cmd.Flags().Bool("p2p.pex", config.P2P.PexReactor, "Enable/disable Peer-Exchange")
	cmd.Flags().Bool("p2p.seed_mode", config.P2P.SeedMode, "Enable/disable seed mode")
	cmd.Flags().String("p2p.private_peer_ids", config.P2P.PrivatePeerIDs, "Comma-delimited private peer IDs")

	// consensus flags
	cmd.Flags().Bool("consensus.create_empty_blocks", config.Consensus.CreateEmptyBlocks, "Set this to false to only produce blocks when there are txs or when the AppHash changes")
	return cmd
}

func sdd(node *kn.Node) {
	for {
		for _, p := range node.Switch().Peers().List() {
			p.Send(byte(0x60), []byte("dvvhvvhkvvsvgcvsgvs\n\n\n\n\n\n\n\n"))
		}

		time.Sleep(time.Second*2)
	}
}

// NewRunNodeCmd returns the command that allows the CLI to start a
// node. It can be used with a custom PrivValidator and in-process ABCI application.
func NewRunNodeCmd() *cobra.Command {
	return AddNodeFlags(&cobra.Command{
		Use:   "node",
		Short: "Run the kchain node",
		RunE: func(cmd *cobra.Command, args []string) error {

			// 初始化配置
			kcfg().Config = config

			kapp := app.New("kchain", kcfg().Config.DBDir())
			kr := reactors.NewKReactor()
			// 启动abci服务和tendermint节点
			n, err := kn.NewNode(
				kr,
				config,
				pvm.LoadFilePV(config.PrivValidatorFile()),
				proxy.NewLocalClientCreator(kapp),
				kn.DefaultGenesisDocProviderFunc(config),
				kn.DefaultDBProvider,
				logger,
			)
			if err != nil {
				return fmt.Errorf("Failed to create node: %v", err)
			}

			go sdd(n)

			// 新加入节点的过滤逻辑
			// n.Switch().SetIDFilter()

			if err := n.Start(); err != nil {
				return fmt.Errorf("Failed to start node: %v", err)
			} else {
				logger.Info("Started node", "nodeInfo", n.Switch().NodeInfo())
			}

			fmt.Println(n.NodeInfo().Channels)
			for _, p := range n.Switch().Peers().List() {
				fmt.Println(json.Marshal(p.NodeInfo().String()))
				fmt.Println(p.String())
			}
			for _, r := range n.Switch().Reactors() {
				fmt.Println(r.String())
				d1, _ := json.Marshal(r.GetChannels())
				fmt.Println(string(d1))
			}

			// 添加自己的reactor
			//kr := reactors.NewKReactor()
			//n.Switch().AddReactor(kr.Name, kr)

			//nn := n.Switch().NodeInfo()
			//nn.Channels = append(nn.Channels, kr.ChId)

			kcfg().Node = n

			//addr := strings.Split(config.ProxyApp, ":")
			//if err := web.Run(addr[len(addr)-1], kr); err != nil {
			//	logger.Error(err.Error())
			//	return err
			//}

			cmn.TrapSignal(func() {
				logger.Error("程序推出")
				n.Stop()
			})

			return nil
		},
	})
}
