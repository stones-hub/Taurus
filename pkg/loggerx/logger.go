package loggerx

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/natefinch/lumberjack"
)

// LogLevel 定义日志等级
type LogLevel int

const (
	LEVEL_DEBUG LogLevel = iota
	LEVEL_INFO
	LEVEL_WARN
	LEVEL_ERROR
	LEVEL_FATAL
	LEVEL_NONE // 不输出日志等级
)

const (
	LEVEL_DEBUG_STR = "DEBUG"
	LEVEL_INFO_STR  = "INFO"
	LEVEL_WARN_STR  = "WARN"
	LEVEL_ERROR_STR = "ERROR"
	LEVEL_FATAL_STR = "FATAL"
	LEVEL_NONE_STR  = ""
)

// levelColors 定义不同日志等级的颜色，仅适用于控制台输出
var levelColors = map[LogLevel]string{
	LEVEL_DEBUG: "\033[36m", // 青色
	LEVEL_INFO:  "\033[32m", // 绿色
	LEVEL_WARN:  "\033[33m", // 黄色
	LEVEL_ERROR: "\033[31m", // 红色
	LEVEL_FATAL: "\033[35m", // 紫色
}

// LoggerConfig 定义日志配置
type LoggerConfig struct {
	Perfix       string                        // 日志前缀, outputType = file下有用
	LogLevel     LogLevel                      // 日志等级
	OutputType   string                        // 输出类型（console/file）
	LogFilePath  string                        // 日志文件路径, 支持相对路径和绝对路径, outputType = file下有用
	MaxSize      int                           // 单个日志文件的最大大小（单位：MB） outputType = file下有用
	MaxBackups   int                           // 保留的旧日志文件的最大数量 outputType = file下有用
	MaxAge       int                           // 日志文件的最大保存天数 outputType = file下有用
	Compress     bool                          // 是否压缩旧日志文件 outputType = file下有用
	CustomFormat func(LogLevel, string) string // 自定义日志格式化函数 outputType = file下有用
}

// Logger 定义日志工具
type Logger struct {
	config LoggerConfig
	logger *log.Logger
	writer io.Writer
}

var DefaultLogger *Logger

// NewLogger 创建一个新的日志工具
func Initialize(config LoggerConfig) *Logger {
	var (
		logFilePath string
		// 配置日志文件轮转
		writer  io.Writer
		baseDIR string
		err     error
		logger  *log.Logger
	)

	if config.OutputType == "file" {
		if filepath.IsAbs(config.LogFilePath) {
			logFilePath = config.LogFilePath // 绝对路径
		} else {
			baseDIR, err = os.Getwd()
			if err != nil {
				log.Fatalf("Failed to get current working directory: %v", err)
			}

			// 相对路径基于 logs 目录, 相对路径
			logFilePath = filepath.Join(baseDIR, config.LogFilePath)
		}

		logDir := filepath.Dir(logFilePath)
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			log.Fatalf("Failed to create log directory: %v", err)
		}

		writer = &lumberjack.Logger{
			Filename:   logFilePath,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}
		logger = log.New(writer, config.Perfix, 0)
	} else {
		writer = os.Stdout
		logger = log.New(writer, "", 0)
	}

	// 如果 FormatFunc 为空，设置默认格式化函数
	if config.CustomFormat == nil {
		config.CustomFormat = defaultFormatFunc
	}

	DefaultLogger = &Logger{
		config: config,
		logger: logger,
		writer: writer,
	}

	return DefaultLogger
}

// defaultFormatFunc 默认日志格式化函数
func defaultFormatFunc(level LogLevel, message string) string {
	// 获取当前时间
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	// 获取调用者信息
	_, file, line, _ := runtime.Caller(3)
	caller := fmt.Sprintf("%s:%d", filepath.Base(file), line)
	return fmt.Sprintf("[%s] [%s] [%s] : %s", timestamp, caller, getLevelSTR(level), message)
}

// logWithLevel 根据日志等级输出日志
func (l *Logger) logWithLevel(level LogLevel, message string) {
	// 日志等级过滤
	if level < l.config.LogLevel {
		return
	}

	// 格式化日志内容
	formattedMessage := l.config.CustomFormat(level, message)

	// 如果是控制台输出，添加颜色
	if l.config.OutputType == "console" {
		color := levelColors[level]
		reset := "\033[0m"
		formattedMessage = color + formattedMessage + reset
	}

	// 打印日志
	l.logger.Println(formattedMessage)
}

func getLevelSTR(level LogLevel) string {
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

// Debug 输出 Debug 级别日志
func (l *Logger) Debug(format string, a ...any) {
	l.logWithLevel(LEVEL_DEBUG, fmt.Sprintf(format, a...))
}

// Info 输出 Info 级别日志
func (l *Logger) Info(format string, a ...any) {
	l.logWithLevel(LEVEL_INFO, fmt.Sprintf(format, a...))
}

// Warn 输出 Warn 级别日志
func (l *Logger) Warn(format string, a ...any) {
	l.logWithLevel(LEVEL_WARN, fmt.Sprintf(format, a...))
}

// Error 输出 Error 级别日志
func (l *Logger) Error(format string, a ...any) {
	l.logWithLevel(LEVEL_ERROR, fmt.Sprintf(format, a...))
}

// Fatal 输出 Fatal 级别日志并终止程序
func (l *Logger) Fatal(format string, a ...any) {
	l.logWithLevel(LEVEL_FATAL, fmt.Sprintf(format, a...))
	os.Exit(1)
}

// SetLogLevel 设置日志等级函数
func (l *Logger) SetLogLevel(level LogLevel) {
	l.config.LogLevel = level
}

// 设置日志输出格式函数
func (l *Logger) SetFormatFunc(formatFunc func(LogLevel, string) string) {
	l.config.CustomFormat = formatFunc
}
