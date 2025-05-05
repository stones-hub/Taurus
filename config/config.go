package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
	"regexp"

	"github.com/BurntSushi/toml"
	"gopkg.in/yaml.v3"
)

// Config holds the application configuration
type Config struct {
	Version       string `json:"version" yaml:"version" toml:"version"`                   // 版本
	AppName       string `json:"app_name" yaml:"app_name" toml:"app_name"`                // 应用名称
	AppHost       string `json:"app_host" yaml:"app_host" toml:"app_host"`                // 应用主机
	AppPort       int    `json:"app_port" yaml:"app_port" toml:"app_port"`                // 应用端口
	Authorization string `json:"authorization" yaml:"authorization" toml:"authorization"` // app授权码
	PrintConfig   bool   `json:"print_config" yaml:"print_config" toml:"print_config"`    // 是否打印配置
	DBEnable      bool   `json:"db_enable" yaml:"db_enable" toml:"db_enable"`             // 是否启用数据库
	RedisEnable   bool   `json:"redis_enable" yaml:"redis_enable" toml:"redis_enable"`    // 是否启用redis
	CronEnable    bool   `json:"cron_enable" yaml:"cron_enable" toml:"cron_enable"`       // 是否启用cron

	// 支持多个数据库配置
	Databases []struct {
		Name       string `json:"name" yaml:"name" toml:"name"`                      // 数据库名称
		Type       string `json:"type" yaml:"type" toml:"type"`                      // 数据库类型 (postgres, mysql, sqlite)
		Host       string `json:"host" yaml:"host" toml:"host"`                      // 数据库主机
		Port       int    `json:"port" yaml:"port" toml:"port"`                      // 数据库端口
		User       string `json:"user" yaml:"user" toml:"user"`                      // 数据库用户名
		Password   string `json:"password" yaml:"password" toml:"password"`          // 数据库密码
		DBName     string `json:"dbname" yaml:"dbname" toml:"dbname"`                // 数据库名称
		SSLMode    string `json:"sslmode" yaml:"sslmode" toml:"sslmode"`             // SSL 模式 (仅适用于 PostgreSQL)
		DSN        string `json:"dsn" yaml:"dsn" toml:"dsn"`                         // 可选，直接提供完整的 DSN 字符串
		MaxRetries int    `json:"max_retries" yaml:"max_retries" toml:"max_retries"` // 最大重试次数
		Delay      int    `json:"delay" yaml:"delay" toml:"delay"`                   // 重试延迟时间 秒

		// 日志配置
		Logger struct {
			LogFilePath   string `json:"log_file_path" yaml:"log_file_path" toml:"log_file_path"`    // 日志文件路径（为空时输出到控制台）
			MaxSize       int    `json:"max_size" yaml:"max_size" toml:"max_size"`                   // 单个日志文件的最大大小（单位：MB）
			MaxBackups    int    `json:"max_backups" yaml:"max_backups" toml:"max_backups"`          // 保留的旧日志文件的最大数量
			MaxAge        int    `json:"max_age" yaml:"max_age" toml:"max_age"`                      // 日志文件的最大保存天数
			Compress      bool   `json:"compress" yaml:"compress" toml:"compress"`                   // 是否压缩旧日志文件
			LogLevel      string `json:"log_level" yaml:"log_level" toml:"log_level"`                // 日志等级 (silent, error, warn, info)
			SlowThreshold int    `json:"slow_threshold" yaml:"slow_threshold" toml:"slow_threshold"` // 慢查询阈值（单位：毫秒）
		} `json:"logger" yaml:"logger" toml:"logger"`
	} `json:"databases" yaml:"databases" toml:"databases"`

	Logger struct {
		Perfix      string `json:"perfix" yaml:"perfix" toml:"perfix"`                      // 日志前缀
		LogLevel    string `json:"log_level" yaml:"log_level" toml:"log_level"`             // 日志等级
		OutputType  string `json:"output_type" yaml:"output_type" toml:"output_type"`       // 输出类型（console/file）
		LogFilePath string `json:"log_file_path" yaml:"log_file_path" toml:"log_file_path"` // 日志文件路径, 支持相对路径和绝对路径
		MaxSize     int    `json:"max_size" yaml:"max_size" toml:"max_size"`                // 单个日志文件的最大大小（单位：MB）
		MaxBackups  int    `json:"max_backups" yaml:"max_backups" toml:"max_backups"`       // 保留的旧日志文件的最大数量
		MaxAge      int    `json:"max_age" yaml:"max_age" toml:"max_age"`                   // 日志文件的最大保存天数
		Compress    bool   `json:"compress" yaml:"compress" toml:"compress"`                // 是否压缩旧日志文件
	} `json:"logger" yaml:"logger" toml:"logger"`

	Redis struct {
		Addrs        []string `json:"addrs" yaml:"addrs" toml:"addrs"`
		Password     string   `json:"password" yaml:"password" toml:"password"`
		DB           int      `json:"db" yaml:"db" toml:"db"`
		PoolSize     int      `json:"pool_size" yaml:"pool_size" toml:"pool_size"`
		MinIdleConns int      `json:"min_idle_conns" yaml:"min_idle_conns" toml:"min_idle_conns"`
		DialTimeout  int      `json:"dial_timeout" yaml:"dial_timeout" toml:"dial_timeout"`
		ReadTimeout  int      `json:"read_timeout" yaml:"read_timeout" toml:"read_timeout"`
		WriteTimeout int      `json:"write_timeout" yaml:"write_timeout" toml:"write_timeout"`
		MaxRetries   int      `json:"max_retries" yaml:"max_retries" toml:"max_retries"`
	} `json:"redis" yaml:"redis" toml:"redis"`
}

// AppConfig is the global configuration instance
var AppConfig Config

// LoadConfig reads and parses configuration files from a directory or a single file
func LoadConfig(path string) {
	info, err := os.Stat(path)
	if err != nil {
		log.Fatalf("Failed to access config path: %v\n", err)
	}

	if info.IsDir() {
		// Recursively load all configuration files in the directory
		err := filepath.Walk(path, func(filePath string, fileInfo os.FileInfo, err error) error {
			if err != nil {
				log.Printf("Error accessing file %s: %v\n", filePath, err)
				return nil
			}

			// Skip directories
			if fileInfo.IsDir() {
				return nil
			}

			loadConfigFile(filePath)
			return nil
		})
		if err != nil {
			log.Fatalf("Failed to walk through config directory: %v\n", err)
		}
	} else {
		// Load a single configuration file
		loadConfigFile(path)
	}

	log.Println("Configuration loaded successfully")
}

// loadConfigFile loads a single configuration file based on its extension
func loadConfigFile(filePath string) {
	ext := filepath.Ext(filePath)
	data, err := os.ReadFile(filePath)
	if err != nil {
		log.Printf("Failed to open config file: %v\n", err)
		return
	}
	// Replace placeholders with environment variables
	content := replacePlaceholders(string(data))

	switch ext {
	case ".json":
		err = json.Unmarshal([]byte(content), &AppConfig)
		if err != nil {
			log.Printf("Failed to parse JSON config file: %s; error: %v\n", filePath, err)
		}
	case ".yaml", ".yml":
		err = yaml.Unmarshal([]byte(content), &AppConfig)
		if err != nil {
			log.Printf("Failed to parse YAML config file: %s; error: %v\n", filePath, err)
		}
	case ".toml":
		_, err = toml.Decode(content, &AppConfig)
		if err != nil {
			log.Printf("Failed to parse TOML config file: %s; error: %v\n", filePath, err)
		}
	default:
		log.Printf("Unsupported config file format: %s\n", filePath)
	}
}

// replacePlaceholders replaces placeholders in the config content with environment variables
func replacePlaceholders(content string) string {
	re := regexp.MustCompile(`\$\{(\w+):([^}]+)\}`)
	return re.ReplaceAllStringFunc(content, func(match string) string {
		parts := re.FindStringSubmatch(match)
		if len(parts) == 3 {
			envVar := parts[1]
			defaultValue := parts[2]
			if value, exists := os.LookupEnv(envVar); exists {
				return value
			}
			return defaultValue
		}
		return match
	})
}
