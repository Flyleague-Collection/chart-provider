// Package global
package global

import (
	"flag"
	"time"
)

var (
	NoLogs         = flag.Bool("no_logs", false, "Disable logging to file")
	ConfigFilePath = flag.String("config", "./config.yaml", "Path to configuration file")
	TokenCacheFile = flag.String("token_cache", "./token_cache", "Path to token cache file")
	RequestTimeout = flag.Duration("request_timeout", 30*time.Second, "Request timeout")
	GzipLevel      = flag.Int("gzip_level", 5, "GZip level")
)

const (
	AppVersion    = "0.1.0"
	ConfigVersion = "0.1.0"

	BeginYear = 2025

	SigningMethod = "HS512"

	DefaultFilePermissions     = 0644
	DefaultDirectoryPermission = 0755

	LogName = "MAIN"

	EnvNoLogs         = "NO_LOGS"
	EnvConfigFilePath = "CONFIG_FILE_PATH"
	EnvTokenCacheFile = "TOKEN_CACHE_FILE"
	EnvRequestTimeout = "REQUEST_TIMEOUT"
	EnvGzipLevel      = "GZIP_LEVEL"
)
