package version

var (
	Version = "0.0.1"
	// git hash should be filled by:
	// go build -ldflags="-X github.com/keystone-coin/coind/version.GitHash=xxxx"
	GitHash   = "dev snapshot"
	BuildDate string
)
