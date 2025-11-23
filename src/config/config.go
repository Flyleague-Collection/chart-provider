// Package config
package config

import (
	"chart-provider/src/interfaces/config"
	"chart-provider/src/interfaces/global"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// readConfig 读取并解析配置文件
func readConfig() (*config.Config, error) {
	c := config.NewConfig()

	file, err := os.OpenFile(*global.ConfigFilePath, os.O_RDONLY, global.DefaultFilePermissions)
	defer func(file *os.File) { _ = file.Close() }(file)

	if err != nil {
		if err := saveConfig(c); err != nil {
			return nil, fmt.Errorf("fail to save configuration file while creating configuration file, %s", err)
		} else {
			return c, nil
		}
	}

	decoder := yaml.NewDecoder(file)

	if err := decoder.Decode(c); err != nil {
		return nil, err
	}

	return c, nil
}

// saveConfig 将配置保存到配置文件
func saveConfig(c *config.Config) error {
	file, err := os.OpenFile(*global.ConfigFilePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, global.DefaultFilePermissions)
	defer func(file *os.File) { _ = file.Close() }(file)
	if err != nil {
		return err
	}

	encoder := yaml.NewEncoder(file)
	defer func(encoder *yaml.Encoder) { _ = encoder.Close() }(encoder)
	return encoder.Encode(c)
}

// Manager 配置管理器结构体
type Manager struct {
	config *config.Config
}

// NewManager 创建一个新的配置管理器实例
func NewManager() *Manager {
	manager := &Manager{
		config: nil,
	}
	return manager
}

func (m *Manager) Init() error {
	c, err := readConfig()
	if err != nil {
		return err
	}
	if ok, err := c.Verify(); !ok {
		return err
	}
	m.config = c
	return nil
}

func (m *Manager) GetConfig() *config.Config {
	return m.config
}

func (m *Manager) SaveConfig() error {
	return saveConfig(m.config)
}
