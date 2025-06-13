// Copyright (c) 2025 Taurus Team. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Author: yelei
// Email: 61647649@qq.com
// Date: 2025-06-13

package consul

import (
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"time"

	"Taurus/pkg/util"

	"github.com/hashicorp/consul/api"
)

// ServiceConfig 服务配置
type ServiceConfig struct {
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

// ConfigWatcher 监听consul KV的变化, 处理接口
type ConfigWatcher interface {
	// 配置文件的 key value
	OnChange(c *ConsulClient, serviceName string, key string, value []byte) error
}

// TTLUpdater 定时处理TTL更新接口
type TTLUpdater interface {
	Update(client *ConsulClient, checkID string) error
}

type InitKVConfig interface { // 用于consul在初始化的时候，需要将配置文件初始化到kv中
	Put(c *ConsulClient, serviceName string) error
}

// Init 初始化Consul服务, server 是consul服务端配置, service 是服务配置, watcher 是配置变更监听, updater 是TTL更新
func Init(server *ServerConfig, service *ServiceConfig, watcher ConfigWatcher, updater TTLUpdater, initKV InitKVConfig) (*ConsulClient, func(), error) {
	// 创建客户端
	client, err := NewConsulClient(server)
	if err != nil {
		return nil, nil, fmt.Errorf("创建客户端失败: %v", err)
	}

	// 生成服务注册信息
	reg, err := buildRegistration(service)
	if err != nil {
		return nil, nil, fmt.Errorf("构建注册信息失败: %v", err)
	}

	log.Println("服务注册信息:", util.ToJsonString(reg))

	// 注册服务
	err = client.Register(reg)
	if err != nil {
		return nil, nil, fmt.Errorf("注册服务失败: %v", err)
	}

	// 验证服务是否成功注册
	maxRetries := 3
	for i := 0; i < maxRetries; i++ {
		services, err := client.GetServices()
		if err != nil {
			log.Printf("获取服务列表失败: %v", err)
			continue
		}
		if _, ok := services[service.Name]; ok {
			break
		}
		if i == maxRetries-1 {
			return nil, nil, fmt.Errorf("服务注册验证失败: 服务 %s 未在 Consul 中注册", service.Name)
		}
		time.Sleep(time.Second)
	}

	// 初始化配置到KV
	if initKV != nil {
		err = initKV.Put(client, service.Name)
		if err != nil {
			return nil, nil, fmt.Errorf("初始化配置到KV失败: %v", err)
		}
	}

	// TTL更新
	if service.Check.Type == "ttl" && updater != nil {
		go startTTL(client, service.Check.CheckID, service.Check.CheckTTL.TTL, updater)
	}

	// 配置监听
	if watcher != nil {
		go watchConfig(client, service.Name, watcher)
	}

	// 返回注销函数
	cleanup := func() {
		// 验证服务是否存在
		services, err := client.GetServices()
		if err != nil {
			log.Printf("获取服务列表失败: %v", err)
		} else if _, ok := services[service.Name]; !ok {
			log.Printf("服务 %s 未在 Consul 中注册，无需注销", service.Name)
			close(client.stop)
			return
		}

		// 尝试注销服务
		maxRetries := 3
		for i := 0; i < maxRetries; i++ {
			err := client.Deregister(service.ID)
			if err == nil {
				log.Printf("服务 %s 注销成功", service.Name)
				break
			}
			log.Printf("注销服务失败(尝试 %d/%d): %v", i+1, maxRetries, err)
			if i < maxRetries-1 {
				time.Sleep(time.Second)
			}
		}

		close(client.stop)
		log.Println("Clean consul client and service success !")
	}

	return client, cleanup, nil
}

// 构建服务注册信息
func buildRegistration(cfg *ServiceConfig) (*api.AgentServiceRegistration, error) {
	reg := &api.AgentServiceRegistration{
		Kind:    api.ServiceKind(cfg.Kind),
		ID:      cfg.ID,
		Name:    cfg.Name,
		Tags:    cfg.Tags,
		Port:    cfg.Port,
		Address: cfg.Address,
		// Namespace: cfg.Namespace, 不支持, 默认使用default
	}

	// 设置地理位置
	if cfg.Locality.Region != "" || cfg.Locality.Zone != "" {
		reg.Locality = &api.Locality{
			Region: cfg.Locality.Region,
			Zone:   cfg.Locality.Zone,
		}
	}

	// 创建健康检查
	check, err := buildHealthCheck(cfg)
	if err != nil {
		return nil, err
	}
	reg.Check = check

	return reg, nil
}

// 创建健康检查
func buildHealthCheck(cfg *ServiceConfig) (*api.AgentServiceCheck, error) {
	check := &api.AgentServiceCheck{
		CheckID:                        cfg.Check.CheckID,
		Name:                           cfg.Check.Name,
		Notes:                          cfg.Check.Notes,
		Status:                         cfg.Check.Status,
		DeregisterCriticalServiceAfter: cfg.Check.DeregisterCriticalServiceAfter,
		SuccessBeforePassing:           cfg.Check.SuccessBeforePassing,
		FailuresBeforeWarning:          cfg.Check.FailuresBeforeWarning,
		FailuresBeforeCritical:         cfg.Check.FailuresBeforeCritical,
	}

	// 根据类型设置参数
	switch cfg.Check.Type {
	case "http":
		check.HTTP = cfg.Check.CheckHTTP.HTTP
		check.Method = cfg.Check.CheckHTTP.Method
		check.Interval = cfg.Check.CheckHTTP.Interval
		check.Timeout = cfg.Check.CheckHTTP.Timeout

		if len(cfg.Check.CheckHTTP.Header) > 0 {
			check.Header = make(map[string][]string)
			for k, v := range cfg.Check.CheckHTTP.Header {
				check.Header[k] = []string{v}
			}
		}
		check.Body = cfg.Check.CheckHTTP.Body

	case "tcp":
		check.TCP = cfg.Check.CheckTCP.TCP
		check.TCPUseTLS = cfg.Check.CheckTCP.TCPUseTLS
		check.Interval = cfg.Check.CheckTCP.Interval
		check.Timeout = cfg.Check.CheckTCP.Timeout
		check.TLSServerName = cfg.Check.CheckTCP.TLSServerName
		check.TLSSkipVerify = cfg.Check.CheckTCP.TLSSkipVerify

	case "grpc":
		check.GRPC = cfg.Check.CheckGRPC.GRPC
		check.GRPCUseTLS = cfg.Check.CheckGRPC.GRPCUseTLS
		check.Interval = cfg.Check.CheckGRPC.Interval
		check.Timeout = cfg.Check.CheckGRPC.Timeout
		check.TLSServerName = cfg.Check.CheckGRPC.TLSServerName
		check.TLSSkipVerify = cfg.Check.CheckGRPC.TLSSkipVerify

	case "ttl":
		check.TTL = cfg.Check.CheckTTL.TTL

	case "shell":
		check.Shell = cfg.Check.CheckShell.Shell
		check.Args = cfg.Check.CheckShell.Args
		check.Interval = cfg.Check.CheckShell.Interval
		check.Timeout = cfg.Check.CheckShell.Timeout
		check.DockerContainerID = cfg.Check.CheckShell.DockerContainerID

	default:
		return nil, fmt.Errorf("不支持的健康检查类型: %s", cfg.Check.Type)
	}

	return check, nil
}

// 启动TTL更新
func startTTL(c *ConsulClient, checkID string, ttl string, updater TTLUpdater) {
	// 解析 TTL 时间
	duration, err := time.ParseDuration(ttl)
	if err != nil {
		log.Printf("解析 TTL 时间失败: %v", err)
		return
	}

	// 设置更新间隔为 TTL 的一半，最小 1 秒
	updateInterval := time.Duration(math.Max(float64(duration/2), float64(time.Second)))
	ticker := time.NewTicker(updateInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			if err := updater.Update(c, checkID); err != nil {
				log.Printf("更新 TTL 失败: %v", err)
			}
		case <-c.stop:
			log.Printf("TTL 更新已停止")
			return
		}
	}
}

// 监听配置变更
func watchConfig(c *ConsulClient, serviceName string, watcher ConfigWatcher) {
	var lastIndex uint64
	var retryCount int
	const maxRetryInterval = time.Minute // 最大重试间隔为1分钟

	for {
		select {
		case <-c.stop:
			log.Printf("配置监听已停止")
			return
		default:
			// 使用 client 的 List 方法获取配置，注意这里会阻塞, 除非配置index跟lastIndex不一致
			pairs, meta, err := c.ListKV(serviceName, lastIndex)
			if err != nil {
				log.Printf("监听配置失败: %v", err)
				// 使用指数退避策略计算重试间隔, 重试间隔按照 2^n 秒递增（1s, 2s, 4s, 8s, 16s, 32s...）
				retryInterval := time.Duration(math.Min(
					float64(time.Second*time.Duration(1<<uint(retryCount))),
					float64(maxRetryInterval),
				))
				time.Sleep(retryInterval)
				retryCount++
				continue
			}

			// 成功获取数据后重置重试计数
			retryCount = 0

			// 如果索引值未变化，说明配置未发生变更，继续监听
			if meta.LastIndex == lastIndex {
				continue
			}
			lastIndex = meta.LastIndex

			// 处理配置变更
			for _, pair := range pairs {
				if err := watcher.OnChange(c, serviceName, pair.Key, pair.Value); err != nil {
					log.Printf("处理配置变更失败: %v", err)
				}
			}
		}
	}
}

// 服务调用
func CallService(ServerName string, request *http.Request) (interface{}, error) {
	// 发现服务
	service, err := Client.Discover(ServerName)
	if err != nil {
		return nil, fmt.Errorf("服务发现失败: %v", err)
	}

	// 构建请求
	request.URL.Host = fmt.Sprintf("%s:%d", service.Service.Address, service.Service.Port)
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %v", err)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %v", err)
	}
	defer response.Body.Close()

	return body, nil
}
