package cmd

import (
	"github.com/spf13/cobra"
	"github.com/codeskyblue/go-sh"
	"strings"
	"time"
	"text/template"
	"path/filepath"
	"bytes"
	"io/ioutil"
	"os"

	kts "kchain/types"
)

const versionTpl = `package version
const Version = "v1.3.0"
const BuildVersion = "{{.BuildVersion}}"
const GitCommit = "{{.GitCommit}}"
`

func BuildFlags(cmd *cobra.Command) *cobra.Command {

	gitCommit, err := sh.Command("git", "rev-parse", "--short=8", "HEAD").Output()
	if err != nil {
		panic(err.Error())
	}

	buf := bytes.NewBuffer(make([]byte, 1024*16))
	template.Must(template.New("version").Parse(versionTpl)).Execute(buf, map[string]string{
		"GitCommit":    strings.TrimSpace(string(gitCommit)),
		"BuildVersion": time.Now().String(),
	})

	if err := ioutil.WriteFile(filepath.Join("version", "version.go"), kts.BytesTrimSpace(buf.Bytes()), 0755); err != nil {
		panic(err.Error())
	}

	filedir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		panic(err.Error())
	}

	cmd.AddCommand(
		&cobra.Command{
			Use:   "local",
			Short: "编译应用",
			RunE: func(cmd *cobra.Command, args []string) error {

				sh.NewSession().
					SetEnv("GOBIN", filedir).
					Command("go", "install", "main.go").
					Run()

				return nil
			},
		},
		&cobra.Command{
			Use:   "linux",
			Short: "交叉编译成linux应用",
			RunE: func(cmd *cobra.Command, args []string) error {

				sh.NewSession().
					SetEnv("GOBIN", filedir).
					SetEnv("CGO_ENABLED", "0").
					SetEnv("GOOS", "linux").
					SetEnv("GOARCH", "amd64").
					Command("go", "install", "main.go").
					Run()

				return nil
			},
		},
	)
	return cmd
}

// NewRunNodeCmd returns the command that allows the CLI to start a
// node. It can be used with a custom PrivValidator and in-process ABCI application.
func NewBuildCmd() *cobra.Command {
	return BuildFlags(&cobra.Command{
		Use:   "build",
		Short: "app build",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	})
}
