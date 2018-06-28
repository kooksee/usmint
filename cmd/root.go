package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	k "github.com/kooksee/usmint/config"
	"github.com/tendermint/tmlibs/cli"
	tmflags "github.com/tendermint/tmlibs/cli/flags"
	"github.com/tendermint/tmlibs/log"
	cfg "github.com/tendermint/tendermint/config"
)

var (
	config = k.DefaultCfg()
	logger = log.NewTMLogger(log.NewSyncWriter(os.Stdout))
)

func init() {
	registerFlagsRootCmd(RootCmd)
}

func registerFlagsRootCmd(cmd *cobra.Command) {
	cmd.PersistentFlags().String("log_level", config.LogLevel, "Log level")
}

// ParseConfig retrieves the default environment configuration,
// sets up the Tendermint root and ensures that the root exists
func ParseConfig() (*k.Config, error) {
	conf := k.DefaultCfg()
	err := viper.Unmarshal(conf)
	if err != nil {
		return nil, err
	}
	conf.SetRoot(conf.RootDir)
	cfg.EnsureRoot(conf.RootDir)
	return conf, err
}

// RootCmd is the root command for Tendermint core.
var RootCmd = &cobra.Command{
	Use:   "mint",
	Short: "mint core",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) (err error) {

		if cmd.Name() == VersionCmd.Name() {
			return nil
		}

		config, err = ParseConfig()
		if err != nil {
			return err
		}

		logger, err = tmflags.ParseLogLevel(config.LogLevel, logger, cfg.DefaultLogLevel())
		if err != nil {
			return err
		}
		if viper.GetBool(cli.TraceFlag) {
			logger = log.NewTracingLogger(logger)
		}

		k.SetLog(logger)

		logger = logger.With("module", "main")
		return nil
	},
}
