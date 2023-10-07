package version

import (
	"fmt"
	"runtime"
)

var (
	Version        = "0.0.0"
	GitCommit      = ""
	ProjectURL     = "github.com/cfanbo/kubectr"
	BuildMeta      = ""
	runtimeVersion = runtime.Version()
)

func FullVersion() string {
	result := fmt.Sprintf("Version: %s\n", Version)
	if GitCommit != "" {
		result += fmt.Sprintf("GitCommit: %s\n", GitCommit)
	}
	result += fmt.Sprintf("runtimeVersion: %s\nProjectURL: %s", runtimeVersion, ProjectURL)

	return result
}
