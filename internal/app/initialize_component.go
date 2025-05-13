package app

import (
	"Taurus/config"
	"Taurus/internal"
	"Taurus/pkg/cron"
	"Taurus/pkg/db"
	"Taurus/pkg/logx"
	"Taurus/pkg/redisx"
	"Taurus/pkg/templates"
	"Taurus/pkg/websocket"
	"fmt"
	"log"
	"time"

	_ "Taurus/internal/controller/crons" // 没有依赖的包， 包体内的init是不会被执行的的; 所以导入
	_ "Taurus/internal/log_formatter"    // 没有依赖的包， 包体内的init是不会被执行的的; 所以导入

	"gorm.io/gorm/logger"
)

var (
	GlobalInjector *internal.Injector
	Cleanup        func()
)

// InitialzeLog initialize logger
func InitialzeLog() {
	// initialize logger
	logConfigs := make([]logx.Config, 0)
	for _, c := range config.Core.Loggers {
		logConfigs = append(logConfigs, logx.Config{
			Name:        c.Name,
			Perfix:      c.Perfix,
			LogLevel:    parseLevel(c.LogLevel),
			OutputType:  c.OutputType,
			LogFilePath: c.LogFilePath,
			MaxSize:     c.MaxSize,
			MaxBackups:  c.MaxBackups,
			MaxAge:      c.MaxAge,
			Compress:    c.Compress,
			Formatter:   c.Formatter,
		})
	}
	logx.Initialize(logConfigs)
}

// InitializeDB initialize database
func InitializeDB() {
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
				LogLevel:      parseDbLevel(dbConfig.Logger.LogLevel),
				SlowThreshold: time.Duration(dbConfig.Logger.SlowThreshold),
			}
			customLogger := db.NewDBCustomLogger(loggerConfig)

			// 初始化数据库
			db.InitDB(dbConfig.Name, dbConfig.Type, dsn, customLogger, dbConfig.MaxRetries, dbConfig.Delay)
			log.Printf("Database '%s' initialized successfully", dbConfig.Name)
		}
	}

}

// InitializeRedis initialize redis
func InitializeRedis() {
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
}

// InitializeTemplates initialize templates
func InitializeTemplates() {
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
}

// InitializeCron initialize cron
func InitializeCron() {
	if config.Core.CronEnable {
		cron.Core.Start()
	}
}

// InitializeWebsocket initialize websocket
func InitializeWebsocket() {
	// initialize websocket
	if config.Core.WebsocketEnable {
		websocket.Initialize()
	}
}

// InitializeInjector initialize injector
func InitializeInjector() {
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
		cleanup()

		if cron.Core != nil {
			cron.Core.Stop()
		}

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

// ParseLogLevel converts a string log level to gorm's logger.LogLevel
func parseDbLevel(level string) logger.LogLevel {
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
func parseLevel(level string) logx.LogLevel {
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

/*
- http.Dir(path):
	职责: 创建一个文件系统，决定了文件系统的位置（path的位置），规定path的地方才可以被访问, 其他地方不可以被访问

- http.FileServer(http.Dir(path)):
	职责: 创建一个文件服务器，解析http请求，并从文件系统中(http.Dir(path))读取文件，返回给客户端

- http.StripPrefix(prefix, handler):
	职责: 从请求路径中移除指定的前缀(prefix)后，将请求转发给下一个处理器(handler)

注意：其实http.FileServer(http.Dir(path)) 就可以直接解决静态文件的问题，但是有弊端，比如:
1. http.Handle("/", http.FileServer(http.Dir("/project/static"))) ,
	当访问 /css/style.css 时，FileServer 会查找 /project/static/css/style.css 但问题是，这会把整个网站根目录映射到静态文件目录，所有URL都会尝试查找静态文件
2. http.Handle("/static/", http.FileServer(http.Dir("/project/static")))
	当访问 /static/css/style.css 时，FileServer 会查找 /project/static/static/css/style.css 多了一个static
所以让三者配合起来一起使用,让您可以更灵活地组织网站结构，而不受文件系统结构的限制

- http.Redirect(w, r, url, code):
	职责: 重定向请求, 支持的重定向状态有以下几种
	1. http.StatusMovedPermanently (301) 永久重定向，告诉浏览器和搜索引擎该资源已永久移动
	2. http.StatusFound (302) 临时重定向，告诉浏览器和搜索引擎该资源已临时移动
	3. http.StatusSeeOther (303) "查看其他"，通常用于POST请求后重定向到GET页面
	4. http.StatusTemporaryRedirect (307) 临时重定向，保留原请求方法
	5. http.StatusPermanentRedirect (308) 永久重定向，保留原请求方法

	注意：
	1. 重定向状态码(code) 是可选的，默认是 http.StatusFound (302)
*/

/*
执行顺序:
导入包全局变量执行 -> 导入包的init函数执行 -> 当前包的全局变量执行 -> 当前包的init函数执行
*/
