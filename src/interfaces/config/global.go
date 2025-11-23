// Package config
package config

import (
	"chart-provider/src/interfaces/global"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type LogConfig struct {
	Level      string `yaml:"level"`
	Path       string `yaml:"path"`
	Rotate     bool   `yaml:"rotate"`
	MaxSize    int    `yaml:"max_size"`
	MaxAge     int    `yaml:"max_age"`
	MaxBackups int    `yaml:"max_backups"`
	Compress   bool   `yaml:"compress"`
	LocalTime  bool   `yaml:"local_time"`
}

func defaultLogConfig() *LogConfig {
	return &LogConfig{
		Level:      "info",
		Path:       "./logs/logs.log",
		Rotate:     true,
		MaxSize:    2,
		MaxAge:     28,
		MaxBackups: 30,
		Compress:   true,
		LocalTime:  true,
	}
}

var levels = []string{"debug", "info", "warn", "error", "fatal"}

func (l *LogConfig) Verify() (bool, error) {
	if l.Level == "" {
		return false, fmt.Errorf("log level is empty")
	}
	level := strings.ToLower(l.Level)
	if !slices.Contains(levels, level) {
		return false, fmt.Errorf("log level is invalid")
	}
	if l.Path == "" {
		return false, fmt.Errorf("log path is empty")
	}
	if err := os.MkdirAll(filepath.Dir(l.Path), global.DefaultDirectoryPermission); err != nil {
		return false, fmt.Errorf("create log directory failed: %s", err)
	}
	if l.Rotate {
		if l.MaxSize <= 0 {
			return false, fmt.Errorf("log max size is invalid")
		}
		if l.MaxAge <= 0 {
			return false, fmt.Errorf("log max age is invalid")
		}
		if l.MaxBackups <= 0 {
			return false, fmt.Errorf("log max backups is invalid")
		}
	}
	return true, nil
}

type GlobalConfig struct {
	Name      string     `yaml:"name"`
	Version   string     `yaml:"version"`
	LogConfig *LogConfig `yaml:"log"`
}

func defaultGlobalConfig() *GlobalConfig {
	return &GlobalConfig{
		Name:      "metar-provider",
		Version:   global.ConfigVersion,
		LogConfig: defaultLogConfig(),
	}
}

func (g *GlobalConfig) Verify() (bool, error) {
	if g.Name == "" {
		return false, fmt.Errorf("global name is empty")
	}
	if g.Version == "" {
		return false, fmt.Errorf("global version is empty")
	}
	configVersion, err := global.NewVersion(g.Version)
	if err != nil {
		return false, fmt.Errorf("global version is invalid: %s", err)
	}
	targetConfigVersion, _ := global.NewVersion(global.ConfigVersion)
	if targetConfigVersion.CheckVersion(configVersion) != global.AllMatch {
		return false, fmt.Errorf("config version mismatch, expected %s, got %s", global.ConfigVersion, g.Version)
	}
	if g.LogConfig == nil {
		return false, fmt.Errorf("log config is empty")
	}
	if ok, err := g.LogConfig.Verify(); !ok {
		return false, err
	}
	return true, nil
}
