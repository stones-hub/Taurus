package log_formatter

import (
	"Taurus/pkg/logx"
	"encoding/json"
	"time"
)

type CustomFormatter struct {
	Formatter string
}

// 自定义格式化函数
func (c *CustomFormatter) Format(level logx.LogLevel, message string) string {

	msg, _ := json.Marshal(struct {
		Level   int    `json:"level"`
		Time    string `json:"time"`
		Message string `json:"message"`
	}{
		Level:   int(level),
		Time:    time.Now().Format("2006-01-02 15:04:05"),
		Message: message,
	})

	return string(msg)
}

func init() {
	logx.RegisterFormatter("custom", &CustomFormatter{})
}
