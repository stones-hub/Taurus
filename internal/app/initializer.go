package app

import (
	"Taurus/config"
	"Taurus/internal"
	"Taurus/pkg/cron"
	"Taurus/pkg/db"
	"Taurus/pkg/logx"
	"Taurus/pkg/mcp/mcp_server"
	"Taurus/pkg/redisx"
	"Taurus/pkg/templates"
	"Taurus/pkg/util"
	"Taurus/pkg/websocket"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"time"

	"github.com/BurntSushi/toml"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"gorm.io/gorm/logger"
)

var (
	GlobalInjector *internal.Injector
	Cleanup        func()
)

// initialize calls the initialization functions of all modules
func initialize(configPath string, env string) {
	// initialize environment variables, if empty, do not load
	err := godotenv.Load(env)
	if err != nil {
		log.Printf("Error loading .env file: %v\n", err.Error())
	}

	// load application configuration file
	log.Printf("Loading application configuration file: %s", configPath)
	loadConfig(configPath)

	// print application configuration
	if config.Core.PrintConfig {
		log.Println("Configuration:", util.ToJsonString(config.Core))
	}

	// initialize logger
	logx.Initialize(logx.LoggerConfig{
		OutputType:  config.Core.Logger.OutputType,
		LogFilePath: config.Core.Logger.LogFilePath,
		MaxSize:     config.Core.Logger.MaxSize,
		MaxBackups:  config.Core.Logger.MaxBackups,
		MaxAge:      config.Core.Logger.MaxAge,
		Compress:    config.Core.Logger.Compress,
		Perfix:      config.Core.Logger.Perfix,
		LogLevel:    parseCustomLoggerLevel(config.Core.Logger.LogLevel),
	})

	// initialize database
	if config.Core.DBEnable {
		for _, dbConfig := range config.Core.Databases {
			// 构造 DSN
			dsn := dbConfig.DSN
			if dsn == "" {
				switch dbConfig.Type {
				case "postgres":
					dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
						dbConfig.Host, dbConfig.Port, dbConfig.User, dbConfig.Password, dbConfig.DBName, dbConfig.SSLMode)
				case "mysql":
					dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
						dbConfig.User, dbConfig.Password, dbConfig.Host, dbConfig.Port, dbConfig.DBName)
				default:
					log.Fatalf("Unsupported database type: %s", dbConfig.Type)
				}
			}

			// 创建自定义日志器
			loggerConfig := db.DBLoggerConfig{
				LogFilePath:   dbConfig.Logger.LogFilePath,
				MaxSize:       dbConfig.Logger.MaxSize,
				MaxBackups:    dbConfig.Logger.MaxBackups,
				MaxAge:        dbConfig.Logger.MaxAge,
				Compress:      dbConfig.Logger.Compress,
				LogLevel:      parseDBLogLevel(dbConfig.Logger.LogLevel),
				SlowThreshold: time.Duration(dbConfig.Logger.SlowThreshold),
			}
			customLogger := db.NewDBCustomLogger(loggerConfig)

			// 初始化数据库
			db.InitDB(dbConfig.Name, dbConfig.Type, dsn, customLogger, dbConfig.MaxRetries, dbConfig.Delay)
			log.Printf("Database '%s' initialized successfully", dbConfig.Name)
		}
	}

	// initialize redis
	if config.Core.RedisEnable {
		redisx.InitRedis(redisx.RedisConfig{
			Addrs:        config.Core.Redis.Addrs,
			Password:     config.Core.Redis.Password,
			DB:           config.Core.Redis.DB,
			PoolSize:     config.Core.Redis.PoolSize,
			MinIdleConns: config.Core.Redis.MinIdleConns,
			DialTimeout:  time.Duration(config.Core.Redis.DialTimeout),
			ReadTimeout:  time.Duration(config.Core.Redis.ReadTimeout),
			WriteTimeout: time.Duration(config.Core.Redis.WriteTimeout),
			MaxRetries:   config.Core.Redis.MaxRetries,
		})
		log.Println("Redis initialized successfully")
	}

	// initialize templates
	if config.Core.TemplatesEnable {
		tmplConfigs := make([]templates.TemplateConfig, 0)
		for _, tmplConf := range config.Core.Templates {
			tmplConfig := templates.TemplateConfig{
				Name: tmplConf.Name,
				Path: tmplConf.Path,
			}
			tmplConfigs = append(tmplConfigs, tmplConfig)
		}
		templates.InitTemplates(tmplConfigs)
	}

	// initialize cron
	if config.Core.CronEnable {
		cron.Core.Start()
	}

	// initialize websocket
	if config.Core.WebsocketEnable {
		websocket.Initialize()
	}

	// initialize mcp server
	if config.Core.MCPEnable {
		mcp_server.InitializeServer(&mcp_server.ServerConfig{
			Name:        config.Core.MCP.Name,
			Version:     config.Core.MCP.Version,
			Addr:        config.Core.MCP.Addr,
			Transport:   config.Core.MCP.Transport,
			Subscribe:   config.Core.MCP.Resource.Subscribe,
			ListChanged: config.Core.MCP.Resource.ListChanged,
			Prompt:      config.Core.MCP.Prompt,
			Tool:        config.Core.MCP.Tool,
		})
		mcp_server.Core.ListenAndServe()
	}

	// initialize injector (internal module initialization)
	initializeInjector()
}

// loadConfig reads and parses configuration files from a directory or a single file
func loadConfig(path string) {
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
		err = json.Unmarshal([]byte(content), &config.Core)
		if err != nil {
			log.Printf("Failed to parse JSON config file: %s; error: %v\n", filePath, err)
		}
	case ".yaml", ".yml":
		err = yaml.Unmarshal([]byte(content), &config.Core)
		if err != nil {
			log.Printf("Failed to parse YAML config file: %s; error: %v\n", filePath, err)
		}
	case ".toml":
		_, err = toml.Decode(content, &config.Core)
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

// ParseLogLevel converts a string log level to gorm's logger.LogLevel
func parseDBLogLevel(level string) logger.LogLevel {
	switch level {
	case "silent":
		return logger.Silent
	case "error":
		return logger.Error
	case "warn":
		return logger.Warn
	case "info":
		return logger.Info
	default:
		log.Printf("Unknown log level '%s', defaulting to 'info'", level)
		return logger.Info
	}
}

// none(无效) error（错误）、warn（警告）、info（信息）、debug（调试）
func parseCustomLoggerLevel(level string) logx.LogLevel {
	switch level {
	case "none":
		return logx.LEVEL_NONE
	case "error":
		return logx.LEVEL_ERROR
	case "warn":
		return logx.LEVEL_WARN
	case "info":
		return logx.LEVEL_INFO
	case "debug":
		return logx.LEVEL_DEBUG
	default:
		log.Printf("Unknown log level '%s', defaulting to 'info'", level)
		return logx.LEVEL_NONE
	}
}

// initialize injector
func initializeInjector() {
	var (
		err     error
		cleanup func()
	)
	// initialize injector
	GlobalInjector, cleanup, err = internal.BuildInjector()
	if err != nil {
		log.Fatalf("Failed to build injector: %v", err)
	}
	Cleanup = func() {
		mcp_server.Core.Shutdown()
		cleanup()
		if redisx.Redis != nil {
			err = redisx.Redis.Close()
			if err != nil {
				log.Printf("Failed to close redis: %v", err)
			} else {
				log.Println("Redis closed successfully")
			}
		}
		db.CloseDB()
	}
}
