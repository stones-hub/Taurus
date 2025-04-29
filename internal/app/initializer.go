package app

import (
	"Taurus/config"
	"Taurus/internal"
	"Taurus/pkg/db"
	"Taurus/pkg/loggerx"
	"Taurus/pkg/redisx"
	"Taurus/pkg/util"
	"Taurus/pkg/websocket"
	"fmt"
	"log"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/gorm/logger"
)

var GlobalInjector *internal.Injector

// Initialize calls the initialization functions of all modules
func Initialize(configPath string, env string) {
	// 1. 加载环境变量文件, 如果为空，则不加载
	err := godotenv.Load(env)
	if err != nil {
		log.Printf("Error loading .env file: %v\n", err.Error())
	}
	// 2. 加载应用配置文件
	config.LoadConfig(configPath)

	if config.AppConfig.PrintConfig {
		log.Println("Configuration:", util.ToJsonString(config.AppConfig))
	}

	// initialize logger
	loggerx.Initialize(loggerx.LoggerConfig{
		OutputType:  config.AppConfig.Logger.OutputType,
		LogFilePath: config.AppConfig.Logger.LogFilePath,
		MaxSize:     config.AppConfig.Logger.MaxSize,
		MaxBackups:  config.AppConfig.Logger.MaxBackups,
		MaxAge:      config.AppConfig.Logger.MaxAge,
		Compress:    config.AppConfig.Logger.Compress,
		Perfix:      config.AppConfig.Logger.Perfix,
		LogLevel:    parseCustomLoggerLevel(config.AppConfig.Logger.LogLevel),
	})

	// initialize database
	if config.AppConfig.DBEnable {
		for _, dbConfig := range config.AppConfig.Databases {
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
	if config.AppConfig.RedisEnable {
		redisx.InitRedis(redisx.RedisConfig{
			Addrs:        config.AppConfig.Redis.Addrs,
			Password:     config.AppConfig.Redis.Password,
			DB:           config.AppConfig.Redis.DB,
			PoolSize:     config.AppConfig.Redis.PoolSize,
			MinIdleConns: config.AppConfig.Redis.MinIdleConns,
			DialTimeout:  time.Duration(config.AppConfig.Redis.DialTimeout),
			ReadTimeout:  time.Duration(config.AppConfig.Redis.ReadTimeout),
			WriteTimeout: time.Duration(config.AppConfig.Redis.WriteTimeout),
			MaxRetries:   config.AppConfig.Redis.MaxRetries,
		})
		log.Println("Redis initialized successfully")
	}

	// initialize injector
	initializeInjector()

	// initialize websocket
	websocket.Initialize()
	log.Println("WebSocket initialized successfully")
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
func parseCustomLoggerLevel(level string) loggerx.LogLevel {
	switch level {
	case "none":
		return loggerx.LEVEL_NONE
	case "error":
		return loggerx.LEVEL_ERROR
	case "warn":
		return loggerx.LEVEL_WARN
	case "info":
		return loggerx.LEVEL_INFO
	case "debug":
		return loggerx.LEVEL_DEBUG
	default:
		log.Printf("Unknown log level '%s', defaulting to 'info'", level)
		return loggerx.LEVEL_NONE
	}
}

// initialize injector
func initializeInjector() {
	var (
		cleanup func()
		err     error
	)
	// initialize injector
	GlobalInjector, cleanup, err = internal.BuildInjector()
	if err != nil {
		log.Fatalf("Failed to build injector: %v", err)
	}
	defer cleanup()
}
