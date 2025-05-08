// Package version provides build and version information variables.
package version

// Version indicates the current build version, injected by the build process.
var (
	Version      string
	GitTag       string
	GitCommit    string
	GitBranch    string
	GitTreeState string
	BuildDate    string
)
