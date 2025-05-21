package config

// Config holds the application configuration
type Config struct {
	Version       string `json:"version" yaml:"version" toml:"version"`                   // 版本
	AppName       string `json:"app_name" yaml:"app_name" toml:"app_name"`                // 应用名称
	AppHost       string `json:"app_host" yaml:"app_host" toml:"app_host"`                // 应用主机
	AppPort       int    `json:"app_port" yaml:"app_port" toml:"app_port"`                // 应用端口
	Authorization string `json:"authorization" yaml:"authorization" toml:"authorization"` // app授权码

	PrintEnable     bool `json:"print_enable" yaml:"print_enable" toml:"print_enable"`             // 是否打印配置
	DBEnable        bool `json:"db_enable" yaml:"db_enable" toml:"db_enable"`                      // 是否启用数据库
	RedisEnable     bool `json:"redis_enable" yaml:"redis_enable" toml:"redis_enable"`             // 是否启用redis
	CronEnable      bool `json:"cron_enable" yaml:"cron_enable" toml:"cron_enable"`                // 是否启用cron
	TemplatesEnable bool `json:"templates_enable" yaml:"templates_enable" toml:"templates_enable"` // 是否启用模板
	WebsocketEnable bool `json:"websocket_enable" yaml:"websocket_enable" toml:"websocket_enable"` // 是否启用websocket
	MCPEnable       bool `json:"mcp_enable" yaml:"mcp_enable" toml:"mcp_enable"`                   // 是否启用mcp
	ConsulEnable    bool `json:"consul_enable" yaml:"consul_enable" toml:"consul_enable"`          // 是否启用consul
	GRPCEnable      bool `json:"grpc_enable" yaml:"grpc_enable" toml:"grpc_enable"`                // 是否启用grpc

	MCP struct {
		Transport string `json:"transport" yaml:"transport" toml:"transport"` // 传输方式，可选值：sse, streamable_http, stdio
		Mode      string `json:"mode" yaml:"mode" toml:"mode"`                // 模式，可选值：stateless, stateful
	} `json:"mcp" yaml:"mcp" toml:"mcp"`

	WebSocket struct {
		Handler string `json:"handler" yaml:"handler" toml:"handler"` // 处理方式，可选值：default,
		Path    string `json:"path" yaml:"path" toml:"path"`          // 路径
	} `json:"websocket" yaml:"websocket" toml:"websocket"`

	Templates []struct {
		Name string `json:"name" yaml:"name" toml:"name"` // 模板名称
		Path string `json:"path" yaml:"path" toml:"path"` // 模板路径
	} `json:"templates" yaml:"templates" toml:"templates"`

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

	Loggers []struct {
		Name        string `json:"name" yaml:"name" toml:"name"`                            // 日志名称
		Perfix      string `json:"perfix" yaml:"perfix" toml:"perfix"`                      // 日志前缀
		LogLevel    string `json:"log_level" yaml:"log_level" toml:"log_level"`             // 日志等级
		OutputType  string `json:"output_type" yaml:"output_type" toml:"output_type"`       // 输出类型（console/file）
		LogFilePath string `json:"log_file_path" yaml:"log_file_path" toml:"log_file_path"` // 日志文件路径, 支持相对路径和绝对路径
		MaxSize     int    `json:"max_size" yaml:"max_size" toml:"max_size"`                // 单个日志文件的最大大小（单位：MB）
		MaxBackups  int    `json:"max_backups" yaml:"max_backups" toml:"max_backups"`       // 保留的旧日志文件的最大数量
		MaxAge      int    `json:"max_age" yaml:"max_age" toml:"max_age"`                   // 日志文件的最大保存天数
		Compress    bool   `json:"compress" yaml:"compress" toml:"compress"`                // 是否压缩旧日志文件
		Formatter   string `json:"formatter" yaml:"formatter" toml:"formatter"`             // 自定义日志格式化函数的名称
	} `json:"loggers" yaml:"loggers" toml:"loggers"`

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

	Consul struct {
		Server  ConsulServer  `json:"server" yaml:"server" toml:"server"`
		Service ConsulService `json:"service" yaml:"service" toml:"service"`
	} `json:"consul" yaml:"consul" toml:"consul"`

	GRPC struct {
		Address  string `json:"address" yaml:"address" toml:"address"`       // grpc地址
		MaxConns int    `json:"max_conns" yaml:"max_conns" toml:"max_conns"` // grpc最大连接数
		TLS      struct {
			Enabled bool   `json:"enabled" yaml:"enabled" toml:"enabled"` // 是否启用TLS
			Cert    string `json:"cert" yaml:"cert" toml:"cert"`          // 证书文件
			Key     string `json:"key" yaml:"key" toml:"key"`             // 证书密钥文件
		} `json:"tls" yaml:"tls" toml:"tls"`
		Keepalive struct {
			Enabled               bool `json:"enabled" yaml:"enabled" toml:"enabled"`                                                    // 是否启用keepalive
			MaxConnectionIdle     int  `json:"max_connection_idle" yaml:"max_connection_idle" toml:"max_connection_idle"`                // 空闲连接最长保持时间 单位: 分钟
			MaxConnectionAge      int  `json:"max_connection_age" yaml:"max_connection_age" toml:"max_connection_age"`                   // 连接在接收到关闭信号，还能保持的时间 单位: 分钟
			MaxConnectionAgeGrace int  `json:"max_connection_age_grace" yaml:"max_connection_age_grace" toml:"max_connection_age_grace"` // MaxConnectionAgeGrace是MaxConnectionAge之后的一个附加周期, 过了这个周期强制关闭 单位: 秒
			Time                  int  `json:"time" yaml:"time" toml:"time"`                                                             // 健康检查间隔 单位: 小时
			Timeout               int  `json:"timeout" yaml:"timeout" toml:"timeout"`                                                    // 健康检查超时 单位: 秒
		} `json:"keepalive" yaml:"keepalive" toml:"keepalive"`
	} `json:"grpc" yaml:"grpc" toml:"grpc"`
}

type ConsulServer struct {
	Address   string `json:"address" yaml:"address" toml:"address"` // consul服务端地址
	Port      int    `json:"port" yaml:"port" toml:"port"`          // consul服务端端口
	Token     string `json:"token" yaml:"token" toml:"token"`       // consul服务端token
	UseTLS    bool   `json:"use_tls" yaml:"use_tls" toml:"use_tls"` // 是否使用TLS
	TLSConfig struct {
		Address            string `json:"address" yaml:"address" toml:"address"`                                        // 证书地址
		Port               int    `json:"port" yaml:"port" toml:"port"`                                                 // 证书端口
		CAFile             string `json:"ca_file" yaml:"ca_file" toml:"ca_file"`                                        // 证书文件
		CertFile           string `json:"cert_file" yaml:"cert_file" toml:"cert_file"`                                  // 证书文件
		KeyFile            string `json:"key_file" yaml:"key_file" toml:"key_file"`                                     // 证书文件
		InsecureSkipVerify bool   `json:"insecure_skip_verify" yaml:"insecure_skip_verify" toml:"insecure_skip_verify"` // 是否忽略证书验证
	} `json:"tls_config" yaml:"tls_config" toml:"tls_config"`
}

type ConsulService struct {
	Kind      string   `json:"kind" yaml:"kind" toml:"kind"`                // 服务类型
	ID        string   `json:"id" yaml:"id" toml:"id"`                      // 服务ID
	Name      string   `json:"name" yaml:"name" toml:"name"`                // 服务名称
	Tags      []string `json:"tags" yaml:"tags" toml:"tags"`                // 服务标签
	Port      int      `json:"port" yaml:"port" toml:"port"`                // 服务端口
	Address   string   `json:"address" yaml:"address" toml:"address"`       // 服务地址
	Namespace string   `json:"namespace" yaml:"namespace" toml:"namespace"` // 服务命名空间
	Locality  struct {
		Region string `json:"region" yaml:"region" toml:"region"` // 服务所在区域
		Zone   string `json:"zone" yaml:"zone" toml:"zone"`       // 服务所在区域
	} `json:"locality" yaml:"locality" toml:"locality"`
	Check struct {
		Type                           string `json:"type" yaml:"type" toml:"type"`                                                                                        // 健康检查类型
		CheckID                        string `json:"check_id" yaml:"check_id" toml:"check_id"`                                                                            // 健康检查ID
		Name                           string `json:"name" yaml:"name" toml:"name"`                                                                                        // 健康检查名称
		Notes                          string `json:"notes" yaml:"notes" toml:"notes"`                                                                                     // 健康检查备注
		Status                         string `json:"status" yaml:"status" toml:"status"`                                                                                  // 健康检查状态
		SuccessBeforePassing           int    `json:"success_before_passing" yaml:"success_before_passing" toml:"success_before_passing"`                                  // 连续成功次数
		FailuresBeforeWarning          int    `json:"failures_before_warning" yaml:"failures_before_warning" toml:"failures_before_warning"`                               // 连续失败次数
		FailuresBeforeCritical         int    `json:"failures_before_critical" yaml:"failures_before_critical" toml:"failures_before_critical"`                            // 连续失败次数
		DeregisterCriticalServiceAfter string `json:"deregister_critical_service_after" yaml:"deregister_critical_service_after" toml:"deregister_critical_service_after"` // 连续失败次数
		CheckTTL                       struct {
			TTL string `json:"ttl" yaml:"ttl" toml:"ttl"` // 健康检查TTL
		} `json:"check_ttl" yaml:"check_ttl" toml:"check_ttl"` // 健康检查TTL
		CheckShell struct {
			Shell             string   `json:"shell" yaml:"shell" toml:"shell"`                                           // 健康检查shell
			Args              []string `json:"args" yaml:"args" toml:"args"`                                              // 健康检查args
			DockerContainerID string   `json:"docker_container_id" yaml:"docker_container_id" toml:"docker_container_id"` // 健康检查docker容器ID
			Interval          string   `json:"interval" yaml:"interval" toml:"interval"`                                  // 健康检查间隔
			Timeout           string   `json:"timeout" yaml:"timeout" toml:"timeout"`                                     // 健康检查超时
		} `json:"check_shell" yaml:"check_shell" toml:"check_shell"` // 健康检查shell
		CheckHTTP struct {
			HTTP     string            `json:"http" yaml:"http" toml:"http"`             // 健康检查http
			Method   string            `json:"method" yaml:"method" toml:"method"`       // 健康检查method
			Header   map[string]string `json:"header" yaml:"header" toml:"header"`       // 健康检查header
			Body     string            `json:"body" yaml:"body" toml:"body"`             // 健康检查body
			Interval string            `json:"interval" yaml:"interval" toml:"interval"` // 健康检查间隔
			Timeout  string            `json:"timeout" yaml:"timeout" toml:"timeout"`    // 健康检查超时
		} `json:"check_http" yaml:"check_http" toml:"check_http"` // 健康检查http
		CheckTCP struct {
			TCP           string `json:"tcp" yaml:"tcp" toml:"tcp"`                                     // 健康检查tcp
			TCPUseTLS     bool   `json:"tcp_use_tls" yaml:"tcp_use_tls" toml:"tcp_use_tls"`             // 健康检查是否使用TLS
			TLSServerName string `json:"tls_server_name" yaml:"tls_server_name" toml:"tls_server_name"` // 健康检查TLS服务器名称
			TLSSkipVerify bool   `json:"tls_skip_verify" yaml:"tls_skip_verify" toml:"tls_skip_verify"` // 健康检查是否跳过TLS证书验证
			Interval      string `json:"interval" yaml:"interval" toml:"interval"`                      // 健康检查间隔
			Timeout       string `json:"timeout" yaml:"timeout" toml:"timeout"`                         // 健康检查超时
		} `json:"check_tcp" yaml:"check_tcp" toml:"check_tcp"` // 健康检查tcp
		CheckGRPC struct {
			GRPC          string `json:"grpc" yaml:"grpc" toml:"grpc"`                                  // 健康检查grpc
			GRPCUseTLS    bool   `json:"grpc_use_tls" yaml:"grpc_use_tls" toml:"grpc_use_tls"`          // 健康检查是否使用TLS
			TLSServerName string `json:"tls_server_name" yaml:"tls_server_name" toml:"tls_server_name"` // 健康检查TLS服务器名称
			TLSSkipVerify bool   `json:"tls_skip_verify" yaml:"tls_skip_verify" toml:"tls_skip_verify"` // 健康检查是否跳过TLS证书验证
			Interval      string `json:"interval" yaml:"interval" toml:"interval"`                      // 健康检查间隔
			Timeout       string `json:"timeout" yaml:"timeout" toml:"timeout"`                         // 健康检查超时
		} `json:"check_grpc" yaml:"check_grpc" toml:"check_grpc"` // 健康检查grpc
	} `json:"check" yaml:"check" toml:"check"`
}

// global configuration instance
var Core Config
