package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tendermint/tmlibs/common"

	pvm "github.com/tendermint/tendermint/types/priv_validator"

	kn "github.com/kooksee/usmint/node"
	"github.com/kooksee/usmint/app"
	"github.com/kooksee/usmint/proxy"
	"github.com/kooksee/usmint/server/web"
	"strings"
	"github.com/kooksee/usmint/cmn"
)

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

// NewRunNodeCmd returns the command that allows the CLI to start a
// node. It can be used with a custom PrivValidator and in-process ABCI application.
func NewRunNodeCmd() *cobra.Command {
	return AddNodeFlags(&cobra.Command{
		Use:   "node",
		Short: "Run the kchain node",
		RunE: func(cmd *cobra.Command, args []string) error {

			// 启动abci服务和tendermint节点
			n, err := kn.NewNode(
				config,
				pvm.LoadFilePV(config.PrivValidatorFile()),
				proxy.NewLocalClientCreator(app.New()),
				kn.DefaultGenesisDocProviderFunc(config),
				kn.DefaultDBProvider,
				logger,
			)
			cmn.MustNotErr("Failed to start node", err, n.Start())
			logger.Info("Started node", "nodeInfo", n.Switch().NodeInfo())

			// web server 启动
			{
				web.Init()
				addr := strings.Split(config.RPC.ListenAddress, ":")
				cmn.MustNotErr("web server 启动失败", web.Run(addr[len(addr)-1]))
			}

			common.TrapSignal(func() {
				cmn.MustNotErr("程序推出", n.Stop())
			})

			return nil
		},
	})
}
