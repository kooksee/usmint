package cmds

import (
	"fmt"

	"github.com/spf13/cobra"

	tv "github.com/tendermint/tendermint/version"
	"github.com/kooksee/usmint/version"
)

// VersionCmd ...
var VersionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show version info",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("tendermint version", tv.Version)
		fmt.Println("kchain version", version.Version)
		fmt.Println("kchain commit version", version.GitCommit)
		fmt.Println("kchain build version", version.BuildVersion)
	},
}
