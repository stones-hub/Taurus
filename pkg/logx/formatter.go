package logx

import (
	"fmt"
	"log"
	"runtime"
	"strings"
	"time"
)

// Formatter 定义格式化函数接口
type Formatter interface {
	Format(level LogLevel, message string) string
}

// 注册表，用于存储用户注册的格式化函数
var formatterRegistry = make(map[string]Formatter)

// RegisterFormatter 注册格式化函数
func RegisterFormatter(name string, formatter Formatter) {
	if _, exists := formatterRegistry[name]; exists {
		log.Printf("Formatter %s already registered", name)
	}
	formatterRegistry[name] = formatter
}

// GetFormatter 获取格式化函数
func GetFormatter(name string) Formatter {
	if formatter, exists := formatterRegistry[name]; exists {
		return formatter
	}
	return defaultFormatter{} // 返回默认格式化函数
}

// 默认格式化函数实现
type defaultFormatter struct{}

func (f defaultFormatter) Format(level LogLevel, message string) string {
	// 获取当前时间
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	// 获取调用者信息
	file, line := GetCallerInfo()
	caller := fmt.Sprintf("%s:%d", file, line)
	return fmt.Sprintf("[%s] [%s] [%s] : %s", timestamp, caller, GetLevelSTR(level), message)
}

// 根据日志等级获取日志等级字符串
func GetLevelSTR(level LogLevel) string {
	switch level {
	case LEVEL_DEBUG:
		return LEVEL_DEBUG_STR
	case LEVEL_INFO:
		return LEVEL_INFO_STR
	case LEVEL_WARN:
		return LEVEL_WARN_STR
	case LEVEL_ERROR:
		return LEVEL_ERROR_STR
	case LEVEL_FATAL:
		return LEVEL_FATAL_STR
	default:
		return ""
	}
}

// GetCallerInfo 获取真实调用者的文件路径和行号
func GetCallerInfo() (string, int) {
	for skip := 0; skip < 15; skip++ {
		_, file, line, ok := runtime.Caller(skip)
		if !ok {
			break
		}
		// 如果文件路径不包含 pkg/logx，说明已经找到了真实调用者
		if !strings.Contains(file, "pkg/logx") {
			return file, line
		}
	}
	return "", 0
}
