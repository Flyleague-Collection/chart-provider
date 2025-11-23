// Package config
package config

import "fmt"

type ServerConfig struct {
	HttpServerConfig *HttpServerConfig `yaml:"http"`
}

func defaultServerConfig() *ServerConfig {
	return &ServerConfig{
		HttpServerConfig: defaultHttpServerConfig(),
	}
}

func (s *ServerConfig) Verify() (bool, error) {
	if s.HttpServerConfig == nil {
		return false, fmt.Errorf("http server config is nil")
	}
	if ok, err := s.HttpServerConfig.Verify(); !ok {
		return false, err
	}
	return true, nil
}
