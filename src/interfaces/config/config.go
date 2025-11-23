// Package config
package config

import "fmt"

type Verifiable interface {
	Verify() (bool, error)
}

type Config struct {
	GlobalConfig *GlobalConfig `yaml:"global"`
}

func NewConfig() *Config {
	return &Config{
		GlobalConfig: defaultGlobalConfig(),
	}
}

func (c *Config) Verify() (bool, error) {
	if c.GlobalConfig == nil {
		return false, fmt.Errorf("global config is nil")
	}
	if ok, err := c.GlobalConfig.Verify(); !ok {
		return false, err
	}
	return true, nil
}
