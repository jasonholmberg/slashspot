package config

// These vars are set at build time by the go linker
var (
	// Version - the version
	Version = "undefined"

	// BuildTime - the build time
	BuildTime = "undefined"

	// GitHash - the git commit hash of this build
	GitHash = "undefined"
)
