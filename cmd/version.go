package cmd

import (
	"github.com/spf13/cobra"

	"github.com/tendermint/tendermint/version"

	uv "github.com/kooksee/usmint/version"
)

// VersionCmd ...
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version info",
	Run: func(cmd *cobra.Command, args []string) {
		logger.Info("tendermint version", "version", version.Version)
		logger.Info("usmint version", "version", uv.Version)
		logger.Info("usmint commit version", "version", uv.GitCommit)
		logger.Info("usmint build version", "version", uv.BuildVersion)
	},
}
