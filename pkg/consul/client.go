package consul

import (
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/hashicorp/consul/api"
)

// ConsulServerConfig Consul服务端配置
type ServerConfig struct {
	Address   string
	Port      int
	Token     string // 用于认证的Token
	UseTLS    bool   // 是否使用TLS
	TLSConfig *api.TLSConfig
}

// ConsulClient Consul客户端
type ConsulClient struct {
	client *api.Client
	stop   chan struct{}
}

var Client *ConsulClient

// NewConsulClient 创建新的Consul客户端
func NewConsulClient(server *ServerConfig) (*ConsulClient, error) {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", server.Address, server.Port)
	config.Token = server.Token
	if server.UseTLS {
		config.TLSConfig = *server.TLSConfig
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("create consul client failed: %v", err)
	}

	// 测试连接, 获取所有服务
	s, _, err := client.Catalog().Services(nil)
	if err != nil {
		return nil, fmt.Errorf("test consul connection failed: %v", err)
	}
	log.Printf("consul services: %v", s)

	Client = &ConsulClient{
		client: client,
		stop:   make(chan struct{}),
	}

	return Client, nil
}

// RegisterService 注册服务，使用配置文件
// 根据配置文件中的健康检查类型动态创建健康检查
func (c *ConsulClient) Register(registration *api.AgentServiceRegistration) error {
	if err := c.client.Agent().ServiceRegister(registration); err != nil {
		log.Printf("Failed to register service: %v", err)
		return err
	}
	return nil
}

// deregister service 注销服务
func (c *ConsulClient) Deregister(serviceID string) error {
	if err := c.client.Agent().ServiceDeregister(serviceID); err != nil {
		log.Printf("Failed to deregister service: %v", err)
		return err
	}
	return nil
}

// discover service 发现服务
func (c *ConsulClient) Discover(serviceName string) (*api.ServiceEntry, error) {
	services, _, err := c.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		log.Printf("Failed to discover service: %v", err)
		return nil, err
	}

	if len(services) == 0 {
		return nil, fmt.Errorf("未找到服务: %s", serviceName)
	}
	// 随机选择一个服务
	service := services[rand.Intn(len(services))]
	return service, nil

}

// 向Consul写入KV配置
func (c *ConsulClient) PutKV(serviceName string, key string, config []byte) (*api.WriteMeta, error) {
	// 构建完整的key，格式：services/{serviceName}/config/{key}
	fullKey := fmt.Sprintf("services/%s/config/%s", serviceName, key)
	writeMeta, err := c.client.KV().Put(&api.KVPair{Key: fullKey, Value: config}, nil)
	if err != nil {
		log.Printf("Failed to put KV: %v", err)
		return nil, err
	}
	return writeMeta, nil
}

// 从Consul获取KV配置
func (c *ConsulClient) GetKV(serviceName string, key string) ([]byte, error) {
	// 构建完整的key，格式：services/{serviceName}/config/{key}
	fullKey := fmt.Sprintf("services/%s/config/%s", serviceName, key)
	pair, _, err := c.client.KV().Get(fullKey, nil)
	if err != nil {
		log.Printf("Failed to get KV: %v", err)
		return nil, err
	}
	if pair != nil {
		return pair.Value, nil
	}
	return nil, nil
}

// 列出服务的所有配置
func (c *ConsulClient) ListKV(serviceName string, waitIndex uint64) (api.KVPairs, *api.QueryMeta, error) {
	// 构建前缀，格式：services/{serviceName}/config/
	prefix := fmt.Sprintf("services/%s/config/", serviceName)
	opts := &api.QueryOptions{ // 设置opts，为了阻塞， 会等到配置发生变化才会返回
		WaitIndex: waitIndex,
		WaitTime:  time.Minute * 5, // 设置 5 分钟超时, 返回
	}
	pairs, meta, err := c.client.KV().List(prefix, opts)
	if err != nil {
		log.Printf("Failed to list KV: %v", err)
		return nil, nil, err
	}
	return pairs, meta, nil
}

// 删除服务的配置
func (c *ConsulClient) DeleteKV(serviceName string, key string) error {
	// 构建完整的key，格式：services/{serviceName}/config/{key}
	fullKey := fmt.Sprintf("services/%s/config/%s", serviceName, key)
	_, err := c.client.KV().Delete(fullKey, nil)
	if err != nil {
		log.Printf("Failed to delete KV: %v", err)
		return err
	}
	return nil
}

// UpdateTTL 更新TTL健康检查状态
// status : passing, warning, critical
func (c *ConsulClient) UpdateTTL(checkID, status, note string) error {
	switch status {
	case "passing":
		return c.client.Agent().PassTTL(checkID, note)
	case "warning":
		return c.client.Agent().WarnTTL(checkID, note)
	case "critical":
		return c.client.Agent().FailTTL(checkID, note)
	default:
		return fmt.Errorf("未知的TTL状态: %s", status)
	}
}

// Close 方法
func (c *ConsulClient) Close() error {
	close(c.stop)
	// 可以添加其他清理逻辑
	return nil
}

// 获取consul的leader
func (c *ConsulClient) GetLeader() (string, error) {
	leader, err := c.client.Status().Leader()
	if err != nil {
		return "", err
	}
	return leader, nil
}

// 获取consul所有节点
func (c *ConsulClient) GetNodes() ([]*api.Node, error) {
	nodes, _, err := c.client.Catalog().Nodes(nil)
	if err != nil {
		return nil, err
	}
	return nodes, nil
}

// 获取consul所有服务
// 返回值: map[服务名称] -> 服务标签
func (c *ConsulClient) GetServices() (map[string][]string, error) {
	services, _, err := c.client.Catalog().Services(nil)
	if err != nil {
		return nil, err
	}
	return services, nil
}
