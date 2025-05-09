// Package version provides version information for the application.
package version

var (
	// Version represents the current version of the application.
	Version string
	// GitTag is the latest git tag from which this build was created.
	GitTag string
	// GitCommit is the git commit hash for this build.
	GitCommit string
	// GitBranch is the git branch from which this build was created.
	GitBranch string
	// GitTreeState indicates whether the git tree was clean or had uncommitted changes.
	GitTreeState string
	// BuildDate is the date when this binary was built.
	BuildDate string
)
