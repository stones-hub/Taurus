package logx

import (
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"

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
type Config struct {
	Name        string   // 日志名称
	Perfix      string   // 日志前缀, outputType = file下有用
	LogLevel    LogLevel // 日志等级
	OutputType  string   // 输出类型（console/file）
	LogFilePath string   // 日志文件路径, 支持相对路径和绝对路径, outputType = file下有用
	MaxSize     int      // 单个日志文件的最大大小（单位：MB） outputType = file下有用
	MaxBackups  int      // 保留的旧日志文件的最大数量 outputType = file下有用
	MaxAge      int      // 日志文件的最大保存天数 outputType = file下有用
	Compress    bool     // 是否压缩旧日志文件 outputType = file下有用
	Formatter   string   // 自定义日志格式化函数的名称 outputType = file下有用
}

// Logger 定义日志工具
type Logger struct {
	config Config
	logger *log.Logger
	writer io.Writer
}

// 定义新的类型
// LoggerMap 封装了 map[string]*Logger

type LoggerMap map[string]*Logger

// 为 LoggerMap 添加方法
func (lm LoggerMap) Info(name string, format string, a ...any) {
	if _, ok := lm[name]; !ok {
		log.Printf("[Warning] Logger %s not found, use default", name)
		lm["default"].Info(format, a...)
	} else {
		lm[name].Info(format, a...)
	}
}

func (lm LoggerMap) Debug(name string, format string, a ...any) {
	if _, ok := lm[name]; !ok {
		log.Printf("[Warning] Logger %s not found, use default", name)
		lm["default"].Debug(format, a...)
	} else {
		lm[name].Debug(format, a...)
	}
}

func (lm LoggerMap) Warn(name string, format string, a ...any) {
	if _, ok := lm[name]; !ok {
		log.Printf("[Warning] Logger %s not found, use default", name)
		lm["default"].Warn(format, a...)
	} else {
		lm[name].Warn(format, a...)
	}
}

func (lm LoggerMap) Error(name string, format string, a ...any) {
	if _, ok := lm[name]; !ok {
		log.Printf("[Warning] Logger %s not found, use default", name)
		lm["default"].Error(format, a...)
	} else {
		lm[name].Error(format, a...)
	}
}

func (lm LoggerMap) Fatal(name string, format string, a ...any) {
	if _, ok := lm[name]; !ok {
		lm["default"].Fatal(format, a...)
	} else {
		lm[name].Fatal(format, a...)
	}
}

// 将 Core 定义为 LoggerMap 类型的实例
var Core = LoggerMap{}

func Initialize(configs []Config) {
	for _, config := range configs {
		if _, ok := Core[config.Name]; ok {
			log.Printf("[Warning] Logger %s already exists", config.Name)
			continue
		}
		Core[config.Name] = new(config)
	}
}

// initialize 初始化日志工具（默认）
func new(c Config) *Logger {
	var (
		logFilePath string
		// 配置日志文件轮转
		writer  io.Writer
		baseDIR string
		err     error
		logger  *log.Logger
	)

	if c.OutputType == "file" {
		if filepath.IsAbs(c.LogFilePath) {
			logFilePath = c.LogFilePath // 绝对路径
		} else {
			baseDIR, err = os.Getwd()
			if err != nil {
				log.Fatalf("Failed to get current working directory: %v", err)
			}

			// 相对路径基于 logs 目录, 相对路径
			logFilePath = filepath.Join(baseDIR, c.LogFilePath)
		}

		logDir := filepath.Dir(logFilePath)
		if err := os.MkdirAll(logDir, os.ModePerm); err != nil {
			log.Fatalf("Failed to create log directory: %v", err)
		}

		writer = &lumberjack.Logger{
			Filename:   logFilePath,
			MaxSize:    c.MaxSize,
			MaxBackups: c.MaxBackups,
			MaxAge:     c.MaxAge,
			Compress:   c.Compress,
		}
		logger = log.New(writer, c.Perfix, 0)
	} else {
		writer = os.Stdout
		logger = log.New(writer, "", 0)
	}

	return &Logger{
		config: c,
		logger: logger,
		writer: writer,
	}
}

// logWithLevel 根据日志等级输出日志
func (l *Logger) logWithLevel(level LogLevel, message string) {
	// 日志等级过滤
	if level < l.config.LogLevel {
		return
	}

	_, file, line, ok := runtime.Caller(3)
	if !ok {
		file = "unknown"
		line = 0
	}

	// 格式化日志内容
	formattedMessage := GetFormatter(l.config.Formatter).Format(level, file, line, message)

	// 如果是控制台输出，添加颜色
	if l.config.OutputType == "console" {
		color := levelColors[level]
		reset := "\033[0m"
		formattedMessage = color + formattedMessage + reset
	}

	// 打印日志
	l.logger.Println(formattedMessage)
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
func (l *Logger) SetFormatter(formatter string) {
	l.config.Formatter = formatter
}
