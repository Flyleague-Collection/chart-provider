// Package global
package global

import (
	"flag"
)

var (
	NoLogs         = flag.Bool("no_logs", false, "Disable logging to file")
	ConfigFilePath = flag.String("config", "./config.yaml", "Path to configuration file")
)

const (
	AppVersion    = "0.1.0"
	ConfigVersion = "0.1.0"

	BeginYear = 2025

	DefaultFilePermissions     = 0644
	DefaultDirectoryPermission = 0755

	LogName = "MAIN"

	EnvNoLogs         = "NO_LOGS"
	EnvConfigFilePath = "CONFIG_FILE_PATH"
)
