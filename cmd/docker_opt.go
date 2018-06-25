package cmd

import (
	"github.com/spf13/cobra"
	"github.com/codeskyblue/go-sh"
	"fmt"

	"kchain/version"
)

func DockerOptFlags(cmd *cobra.Command) *cobra.Command {
	imagesPrefix := "registry.cn-hangzhou.aliyuncs.com/yuanben/"
	imageName := fmt.Sprintf("kchain:%s_%s", version.Version, version.GitCommit)
	cmd.AddCommand(
		&cobra.Command{
			Use:   "build",
			Short: "Run the docker build",
			RunE: func(cmd *cobra.Command, args []string) error {

				sh.Command("docker", "build", "-t", imageName, ".").Run()
				sh.Command("docker", "tag", imageName, imagesPrefix+imageName).Run()
				sh.Command("docker", "tag", imageName, imagesPrefix+"kchain").Run()

				return nil
			},
		},
		&cobra.Command{
			Use:   "push",
			Short: "Run the docker push",
			RunE: func(cmd *cobra.Command, args []string) error {
				sh.Command("docker", "push", imagesPrefix+imageName).Run()
				sh.Command("docker", "push", imagesPrefix+"kchain").Run()
				return nil
			},
		},
		&cobra.Command{
			Use:   "rm_none",
			Short: "Run the docker rm_none",
			RunE: func(cmd *cobra.Command, args []string) error {
				sh.Command("sudo", "docker", "images", "").
					Command("grep", "none").
					Command("awk", "{print $3}").
					Command("xargs", "docker", "rmi", "-f").Run()

				sh.Command("sudo", "docker", "images").Run()
				return nil
			},
		},
	)
	return cmd
}

// NewRunNodeCmd returns the command that allows the CLI to start a
// node. It can be used with a custom PrivValidator and in-process ABCI application.
func NewDockerOptCmd() *cobra.Command {
	return DockerOptFlags(&cobra.Command{
		Use:   "docker",
		Short: "Run the docker opt",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	})
}
