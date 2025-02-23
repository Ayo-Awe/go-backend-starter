package vcs

import (
	"fmt"
	"runtime/debug"
	"sync"
)

var (
	version   string // Version number, e.g., "1.0.0"
	commitSHA string // Git commit SHA
	buildTime string // Build time in a specific format

	// fullVersion is the cached full version string
	fullVersion     string
	fullVersionOnce sync.Once // Ensures fullVersion is calculated only once
)

// Version returns the full version string, ensuring it's calculated only once.
func Version() string {
	fullVersionOnce.Do(func() {
		fullVersion = buildFullVersion()
	})
	return fullVersion
}

// buildFullVersion constructs the full version string based on available metadata.
func buildFullVersion() string {
	if version == "" {
		return buildDevVersion()
	}
	return fmt.Sprintf("%s-%s-%s%s", version, commitSHA, buildTime, getModifiedSuffix())
}

// buildDevVersion constructs a version string for development builds.
func buildDevVersion() string {
	return fmt.Sprintf("dev%s", getModifiedSuffix())
}

// getModifiedSuffix checks if the build is from a modified source tree and returns a suffix if so.
func getModifiedSuffix() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		return ""
	}

	for _, setting := range bi.Settings {
		if setting.Key == "vcs.modified" && setting.Value == "true" {
			return "-dirty"
		}
	}
	return ""
}
