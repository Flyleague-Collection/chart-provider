// Package main
package main

import (
	"chart-provider/src/interfaces/global"
	"flag"
	"fmt"
	"time"

	cleanerImpl "chart-provider/src/cleaner"
	configImpl "chart-provider/src/config"
	loggerImpl "chart-provider/src/logger"
)

func main() {

	flag.Parse()

	global.CheckBoolEnv(global.EnvNoLogs, global.NoLogs)
	global.CheckStringEnv(global.EnvConfigFilePath, global.ConfigFilePath)

	configManager := configImpl.NewManager()
	if err := configManager.Init(); err != nil {
		fmt.Printf("fail to initialize configuration file: %v", err)
		return
	}

	applicationConfig := configManager.GetConfig()
	logger := loggerImpl.NewLogger()
	logger.Init(
		applicationConfig.GlobalConfig.LogConfig.Path,
		global.LogName,
		applicationConfig.GlobalConfig.LogConfig.Level,
		applicationConfig.GlobalConfig.LogConfig,
	)

	logger.Info(" _____ _           _   _____             _   _")
	logger.Info("|     | |_ ___ ___| |_|  _  |___ ___ _ _|_|_| |___ ___")
	logger.Info("|   --|   | .'|  _|  _|   __|  _| . | | | | . | -_|  _|")
	logger.Info("|_____|_|_|__,|_| |_| |__|  |_| |___|\\_/|_|___|___|_|")
	logger.Infof("                     Copyright © %d-%d Half_nothing", global.BeginYear, time.Now().Year())
	logger.Infof("                                   ChartProvider v%s", global.AppVersion)

	cleaner := cleanerImpl.NewCleaner(logger)
	cleaner.Init()

	//applicationContent := content.NewApplicationContentBuilder().
	//	SetConfigManager(configManager).
	//	SetCleaner(cleaner).
	//	SetLogger(logger).
	//	Build()

	cleaner.Wait()
}
