// Package config
package config

type ManagerInterface interface {
	Init() error
	GetConfig() *Config
	SaveConfig() error
}
