// Package content
package content

import (
	"chart-provider/src/interfaces/cleaner"
	"chart-provider/src/interfaces/config"
	"chart-provider/src/interfaces/logger"
)

// ApplicationContent 应用程序上下文结构体，包含所有核心组件的接口
type ApplicationContent struct {
	configManager config.ManagerInterface // 配置管理器
	cleaner       cleaner.Interface       // 清理器
	logger        logger.Interface        // 日志
}

func (app *ApplicationContent) ConfigManager() config.ManagerInterface {
	return app.configManager
}

func (app *ApplicationContent) Cleaner() cleaner.Interface { return app.cleaner }

func (app *ApplicationContent) Logger() logger.Interface { return app.logger }
