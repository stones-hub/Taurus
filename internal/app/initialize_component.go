package app

import (
	"Taurus/config"
	"Taurus/internal"
	"Taurus/pkg/consul"
	"Taurus/pkg/cron"
	"Taurus/pkg/db"
	"Taurus/pkg/grpc/server"
	"Taurus/pkg/logx"
	"Taurus/pkg/mcp"
	"Taurus/pkg/middleware"
	"Taurus/pkg/redisx"
	"Taurus/pkg/router"
	"Taurus/pkg/templates"
	"Taurus/pkg/wsocket"
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"time"

	"Taurus/internal/app/core/consuls"
	_ "Taurus/internal/app/core/crons"         // 引入crons，注册crons包下的所有的定时任务
	_ "Taurus/internal/app/core/log_formatter" // 引入log_formatter，注册log_formatter包下的所有的日志格式化器
	_ "Taurus/internal/app/core/ws_handler"    // 引入ws_handler，注册ws_handler包下的所有的websocket处理器

	// 引入mcps包下的所有的提示词、资源、工具
	_ "Taurus/internal/app/core/mcps/prompts"   // 引入prompts，注册prompts包下的所有的提示词
	_ "Taurus/internal/app/core/mcps/resources" // 引入resources，注册resources包下的所有的资源
	_ "Taurus/internal/app/core/mcps/tools"     // 引入tools，注册tools包下的所有的工具

	// 引入 gRPC 包下的所有的中间件、服务
	_ "Taurus/internal/controller/gRPC/mid"     // 引入mid，注册mid包下的所有的中间件
	_ "Taurus/internal/controller/gRPC/service" // 引入service，注册service包下的所有的服务

	"github.com/hashicorp/consul/api"
	"google.golang.org/grpc/keepalive"
	"gorm.io/gorm/logger"
)

var (
	Cleanup []func() = make([]func(), 0)
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
	log.Println("\033[1;32m🔗 -> Log initialized successfully\033[0m")
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
		log.Println("\033[1;32m🔗 -> Database all initialized successfully\033[0m")
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
		log.Println("\033[1;32m🔗 -> Redis initialized successfully\033[0m")
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
		log.Println("\033[1;32m🔗 -> Templates initialized successfully\033[0m")
	}
}

// InitializeCron initialize cron
func InitializeCron() {
	if config.Core.CronEnable {
		cron.Core.Start()
		log.Println("\033[1;32m🔗 -> Cron initialized successfully\033[0m")
	}
}

// InitializeWebsocket initialize websocket
func InitializeWebsocket() {
	// initialize websocket
	if config.Core.WebsocketEnable {
		wsocket.Initialize()

		router.AddRouter(router.Router{
			Path: "/ws",
			Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wsocket.HandleWebSocket(w, r, wsocket.GetHandler(config.Core.WebSocket.Handler).Handle)
			}),
			Middleware: []router.MiddlewareFunc{
				middleware.ErrorHandlerMiddleware,
				middleware.TraceMiddleware,
			},
		})
		log.Println("\033[1;32m🔗 -> Websocket initialized successfully\033[0m")
	}
}

// InitializeMCP initialize mcp, but need to register tools, prompts, resources, resource templates
func InitializeMCP() {
	// 注意：stdio 模式下，需要手动启动 server，其他模式下，server 会自动启动， 不建议在http服务器上使用stdio模式，如果需要，可以依据工具函数，自行构建main函数
	if config.Core.MCPEnable && config.Core.MCP.Transport != mcp.TransportStdio {
		server, _, err := mcp.NewMCPServer(config.Core.AppName, config.Core.Version, config.Core.MCP.Transport, config.Core.MCP.Mode)
		if err != nil {
			log.Fatalf("Failed to initialize mcp server: %v", err)
		}
		// register handler for mcp server
		server.RegisterHandler(mcp.MCPHandler)
		log.Println("\033[1;32m🔗 -> MCP initialized successfully\033[0m")
	}
}

// InitializeInjector initialize injector
func InitializeInjector() {
	var (
		err     error
		cleanup func()
	)
	// initialize injector
	internal.Core, cleanup, err = internal.BuildInjector()
	if err != nil {
		log.Fatalf("Failed to build injector: %v", err)
	}

	Cleanup = append(Cleanup, func() {
		cleanup()
		log.Println("\033[1;32m🔗 -> Injector initialized successfully\033[0m")

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
		log.Printf("%s🔗 -> Clean up cron redis db components successfully. %s\n", Green, Reset)
	})
}

// InitializegRPC initialize grpc
func InitializegRPC() {
	// initialize grpc
	if config.Core.GRPCEnable {
		opts := []server.ServerOption{
			server.WithAddress(config.Core.GRPC.Address),
			server.WithMaxConns(config.Core.GRPC.MaxConns),
		}

		for _, middleware := range server.GetServiceMiddleware() {
			opts = append(opts, server.WithUnaryMiddleware(middleware))
		}

		for _, streamMiddleware := range server.GetServiceStreamMiddleware() {
			opts = append(opts, server.WithStreamMiddleware(streamMiddleware))
		}

		for _, interceptor := range server.GetServiceInterceptor() {
			opts = append(opts, server.WithUnaryInterceptor(interceptor))
		}

		for _, streamInterceptor := range server.GetServiceStreamInterceptor() {
			opts = append(opts, server.WithStreamInterceptor(streamInterceptor))
		}

		if config.Core.GRPC.TLS.Enabled {
			cert, err := tls.LoadX509KeyPair(config.Core.GRPC.TLS.Cert, config.Core.GRPC.TLS.Key)
			if err != nil {
				log.Fatalf("Failed to load TLS certificate: %v", err)
			}
			opts = append(opts, server.WithTLS(&tls.Config{
				Certificates: []tls.Certificate{cert},
				MinVersion:   tls.VersionTLS12,
			}))
		}

		if config.Core.GRPC.Keepalive.Enabled {
			opts = append(opts, server.WithKeepAlive(&keepalive.ServerParameters{
				Time:                  time.Duration(config.Core.GRPC.Keepalive.Time) * time.Hour,
				Timeout:               time.Duration(config.Core.GRPC.Keepalive.Timeout) * time.Second,
				MaxConnectionIdle:     time.Duration(config.Core.GRPC.Keepalive.MaxConnectionIdle) * time.Minute,
				MaxConnectionAge:      time.Duration(config.Core.GRPC.Keepalive.MaxConnectionAge) * time.Minute,
				MaxConnectionAgeGrace: time.Duration(config.Core.GRPC.Keepalive.MaxConnectionAgeGrace) * time.Second,
			}))
		}

		s, cleanup, err := server.NewServer(opts...)
		if err != nil {
			log.Fatalf("Failed to initialize gRPC server: %v", err)
		}
		Cleanup = append(Cleanup, cleanup)

		// 遍历所有注册的服务注册
		for _, service := range server.GetRegisteredServices() {
			service.RegisterService(s.Server())
		}

		go func() {
			err := s.Start()
			if err != nil {
				log.Fatalf("Failed to start gRPC server: %v", err)
			}
		}()
		log.Println("\033[1;32m🔗 -> gRPC initialized successfully\033[0m")
	}
}

// InitializeConsul initialize consul
func InitializeConsul() {
	if config.Core.ConsulEnable {
		serverConfig, serviceConfig := buildConsulConfig(config.Core.Consul.Server, config.Core.Consul.Service)
		_, cleanup, err := consul.Init(serverConfig, serviceConfig, new(consuls.DefaultConfigWatcher), new(consuls.DefaultTTLUpdater), new(consuls.DefaultInitKVConfig))
		if err != nil {
			log.Fatalf("Failed to initialize consul: %v", err)
		}
		Cleanup = append(Cleanup, cleanup)
		log.Println("\033[1;32m🔗 -> Consul initialized successfully\033[0m")
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

// BuildConfig 构建Consul配置
func buildConsulConfig(server config.ConsulServer, service config.ConsulService) (*consul.ServerConfig, *consul.ServiceConfig) {
	serverConfig := &consul.ServerConfig{
		Address: server.Address,
		Port:    server.Port,
		Token:   server.Token,
		UseTLS:  server.UseTLS,
		TLSConfig: &api.TLSConfig{
			Address:            server.TLSConfig.Address,
			CAFile:             server.TLSConfig.CAFile,
			CertFile:           server.TLSConfig.CertFile,
			KeyFile:            server.TLSConfig.KeyFile,
			InsecureSkipVerify: server.TLSConfig.InsecureSkipVerify,
		},
	}

	serviceConfig := &consul.ServiceConfig{
		Kind:      service.Kind,
		ID:        service.ID,
		Name:      service.Name,
		Tags:      service.Tags,
		Port:      service.Port,
		Address:   service.Address,
		Namespace: service.Namespace,
		Locality: struct {
			Region string `json:"region" yaml:"region" toml:"region"`
			Zone   string `json:"zone" yaml:"zone" toml:"zone"`
		}{
			Region: service.Locality.Region,
			Zone:   service.Locality.Zone,
		},
		Check: struct {
			Type                           string `json:"type" yaml:"type" toml:"type"`
			CheckID                        string `json:"check_id" yaml:"check_id" toml:"check_id"`
			Name                           string `json:"name" yaml:"name" toml:"name"`
			Notes                          string `json:"notes" yaml:"notes" toml:"notes"`
			Status                         string `json:"status" yaml:"status" toml:"status"`
			SuccessBeforePassing           int    `json:"success_before_passing" yaml:"success_before_passing" toml:"success_before_passing"`
			FailuresBeforeWarning          int    `json:"failures_before_warning" yaml:"failures_before_warning" toml:"failures_before_warning"`
			FailuresBeforeCritical         int    `json:"failures_before_critical" yaml:"failures_before_critical" toml:"failures_before_critical"`
			DeregisterCriticalServiceAfter string `json:"deregister_critical_service_after" yaml:"deregister_critical_service_after" toml:"deregister_critical_service_after"`
			CheckTTL                       struct {
				TTL string `json:"ttl" yaml:"ttl" toml:"ttl"`
			} `json:"check_ttl" yaml:"check_ttl" toml:"check_ttl"`
			CheckShell struct {
				Shell             string   `json:"shell" yaml:"shell" toml:"shell"`
				Args              []string `json:"args" yaml:"args" toml:"args"`
				DockerContainerID string   `json:"docker_container_id" yaml:"docker_container_id" toml:"docker_container_id"`
				Interval          string   `json:"interval" yaml:"interval" toml:"interval"`
				Timeout           string   `json:"timeout" yaml:"timeout" toml:"timeout"`
			} `json:"check_shell" yaml:"check_shell" toml:"check_shell"`
			CheckHTTP struct {
				HTTP     string            `json:"http" yaml:"http" toml:"http"`
				Method   string            `json:"method" yaml:"method" toml:"method"`
				Header   map[string]string `json:"header" yaml:"header" toml:"header"`
				Body     string            `json:"body" yaml:"body" toml:"body"`
				Interval string            `json:"interval" yaml:"interval" toml:"interval"`
				Timeout  string            `json:"timeout" yaml:"timeout" toml:"timeout"`
			} `json:"check_http" yaml:"check_http" toml:"check_http"`
			CheckTCP struct {
				TCP           string `json:"tcp" yaml:"tcp" toml:"tcp"`
				TCPUseTLS     bool   `json:"tcp_use_tls" yaml:"tcp_use_tls" toml:"tcp_use_tls"`
				TLSServerName string `json:"tls_server_name" yaml:"tls_server_name" toml:"tls_server_name"`
				TLSSkipVerify bool   `json:"tls_skip_verify" yaml:"tls_skip_verify" toml:"tls_skip_verify"`
				Interval      string `json:"interval" yaml:"interval" toml:"interval"`
				Timeout       string `json:"timeout" yaml:"timeout" toml:"timeout"`
			} `json:"check_tcp" yaml:"check_tcp" toml:"check_tcp"`
			CheckGRPC struct {
				GRPC          string `json:"grpc" yaml:"grpc" toml:"grpc"`
				GRPCUseTLS    bool   `json:"grpc_use_tls" yaml:"grpc_use_tls" toml:"grpc_use_tls"`
				TLSServerName string `json:"tls_server_name" yaml:"tls_server_name" toml:"tls_server_name"`
				TLSSkipVerify bool   `json:"tls_skip_verify" yaml:"tls_skip_verify" toml:"tls_skip_verify"`
				Interval      string `json:"interval" yaml:"interval" toml:"interval"`
				Timeout       string `json:"timeout" yaml:"timeout" toml:"timeout"`
			} `json:"check_grpc" yaml:"check_grpc" toml:"check_grpc"`
		}{
			Type:                           service.Check.Type,
			CheckID:                        service.Check.CheckID,
			Name:                           service.Check.Name,
			Notes:                          service.Check.Notes,
			Status:                         service.Check.Status,
			SuccessBeforePassing:           service.Check.SuccessBeforePassing,
			FailuresBeforeWarning:          service.Check.FailuresBeforeWarning,
			FailuresBeforeCritical:         service.Check.FailuresBeforeCritical,
			DeregisterCriticalServiceAfter: service.Check.DeregisterCriticalServiceAfter,
			CheckTTL: struct {
				TTL string `json:"ttl" yaml:"ttl" toml:"ttl"`
			}{
				TTL: service.Check.CheckTTL.TTL,
			},
			CheckShell: struct {
				Shell             string   `json:"shell" yaml:"shell" toml:"shell"`
				Args              []string `json:"args" yaml:"args" toml:"args"`
				DockerContainerID string   `json:"docker_container_id" yaml:"docker_container_id" toml:"docker_container_id"`
				Interval          string   `json:"interval" yaml:"interval" toml:"interval"`
				Timeout           string   `json:"timeout" yaml:"timeout" toml:"timeout"`
			}{
				Shell:             service.Check.CheckShell.Shell,
				Args:              service.Check.CheckShell.Args,
				DockerContainerID: service.Check.CheckShell.DockerContainerID,
				Interval:          service.Check.CheckShell.Interval,
				Timeout:           service.Check.CheckShell.Timeout,
			},
			CheckHTTP: struct {
				HTTP     string            `json:"http" yaml:"http" toml:"http"`
				Method   string            `json:"method" yaml:"method" toml:"method"`
				Header   map[string]string `json:"header" yaml:"header" toml:"header"`
				Body     string            `json:"body" yaml:"body" toml:"body"`
				Interval string            `json:"interval" yaml:"interval" toml:"interval"`
				Timeout  string            `json:"timeout" yaml:"timeout" toml:"timeout"`
			}{
				HTTP:     service.Check.CheckHTTP.HTTP,
				Method:   service.Check.CheckHTTP.Method,
				Header:   service.Check.CheckHTTP.Header,
				Body:     service.Check.CheckHTTP.Body,
				Interval: service.Check.CheckHTTP.Interval,
				Timeout:  service.Check.CheckHTTP.Timeout,
			},
			CheckTCP: struct {
				TCP           string `json:"tcp" yaml:"tcp" toml:"tcp"`
				TCPUseTLS     bool   `json:"tcp_use_tls" yaml:"tcp_use_tls" toml:"tcp_use_tls"`
				TLSServerName string `json:"tls_server_name" yaml:"tls_server_name" toml:"tls_server_name"`
				TLSSkipVerify bool   `json:"tls_skip_verify" yaml:"tls_skip_verify" toml:"tls_skip_verify"`
				Interval      string `json:"interval" yaml:"interval" toml:"interval"`
				Timeout       string `json:"timeout" yaml:"timeout" toml:"timeout"`
			}{
				TCP:           service.Check.CheckTCP.TCP,
				TCPUseTLS:     service.Check.CheckTCP.TCPUseTLS,
				TLSServerName: service.Check.CheckTCP.TLSServerName,
				TLSSkipVerify: service.Check.CheckTCP.TLSSkipVerify,
				Interval:      service.Check.CheckTCP.Interval,
				Timeout:       service.Check.CheckTCP.Timeout,
			},
			CheckGRPC: struct {
				GRPC          string `json:"grpc" yaml:"grpc" toml:"grpc"`
				GRPCUseTLS    bool   `json:"grpc_use_tls" yaml:"grpc_use_tls" toml:"grpc_use_tls"`
				TLSServerName string `json:"tls_server_name" yaml:"tls_server_name" toml:"tls_server_name"`
				TLSSkipVerify bool   `json:"tls_skip_verify" yaml:"tls_skip_verify" toml:"tls_skip_verify"`
				Interval      string `json:"interval" yaml:"interval" toml:"interval"`
				Timeout       string `json:"timeout" yaml:"timeout" toml:"timeout"`
			}{
				GRPC:          service.Check.CheckGRPC.GRPC,
				GRPCUseTLS:    service.Check.CheckGRPC.GRPCUseTLS,
				TLSServerName: service.Check.CheckGRPC.TLSServerName,
				TLSSkipVerify: service.Check.CheckGRPC.TLSSkipVerify,
				Interval:      service.Check.CheckGRPC.Interval,
				Timeout:       service.Check.CheckGRPC.Timeout,
			},
		},
	}

	return serverConfig, serviceConfig
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
