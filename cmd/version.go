package cmds

import (
	"fmt"

	"github.com/spf13/cobra"

	tv "github.com/tendermint/tendermint/version"
	"github.com/kooksee/usmint/version"
)

// VersionCmd ...
var VersionCmd = &cobra.Command{
	Use:     "version",
	Aliases: []string{"v", "ver"},
	Short:   "Show version info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tendermint version", tv.Version)
		fmt.Println("usmint version", version.Version)
		fmt.Println("usmint commit version", version.GitCommit)
		fmt.Println("usmint build version", version.BuildVersion)
	},
}
