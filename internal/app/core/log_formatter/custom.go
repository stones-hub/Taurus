package log_formatter

import (
	"Taurus/pkg/logx"
)

type CustomFormatter struct {
	Formatter string
}

// 自定义格式化函数
func (c *CustomFormatter) Format(level logx.LogLevel, message string) string {
	return message
}

func init() {
	logx.RegisterFormatter("custom", &CustomFormatter{})
}
