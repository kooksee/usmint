package version

import "fmt"

// Major version component of the current release
const Major = 0

// Minor version component of the current release
const Minor = 5

// Fix version component of the current release
const Fix = 3

var (
	// Version is the full version string
	Version = fmt.Sprintf("%d.%d.%d", Major, Minor, Fix)

	// GitCommit is set with --ldflags "-X main.gitCommit=$(git rev-parse --short HEAD)"
	GitCommit string
)

func init() {
	if GitCommit != "" {
		Version += "-" + GitCommit
	}
}
