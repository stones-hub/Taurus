package consul

import (
	"fmt"
	"log"

	"github.com/hashicorp/consul/api"
)

// ConfigChangeHandler 配置变更处理接口
type ConfigChangeHandler interface {
	Handle(key string, value []byte) error
}

// TTLUpdater TTL更新处理接口
type TTLUpdater interface {
	Update(client *ConsulClient, checkID string) error
}

// ServiceCaller 服务调用接口
type ServiceCaller interface {
	CallService(serviceName string, args interface{}) (interface{}, error)
}

// ConsulServerConfig Consul服务端配置
type ConsulServerConfig struct {
	Address   string
	Port      int
	Token     string // 用于认证的Token
	UseTLS    bool   // 是否使用TLS
	TLSConfig *api.TLSConfig
}

// ConsulClient Consul客户端
type ConsulClient struct {
	client *api.Client
}

// NewConsulClient 创建新的Consul客户端
func NewConsulClient(consulServerConfig *ConsulServerConfig) (*ConsulClient, error) {
	config := api.DefaultConfig()
	config.Address = fmt.Sprintf("%s:%d", consulServerConfig.Address, consulServerConfig.Port)
	config.Token = consulServerConfig.Token
	if consulServerConfig.UseTLS {
		config.TLSConfig = *consulServerConfig.TLSConfig
	}

	client, err := api.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("create consul client failed: %v", err)
	}

	return &ConsulClient{
		client: client,
	}, nil
}

// RegisterService 注册服务，使用配置文件
// 根据配置文件中的健康检查类型动态创建健康检查
func (c *ConsulClient) RegisterService(registration *api.AgentServiceRegistration) error {
	if err := c.client.Agent().ServiceRegister(registration); err != nil {
		log.Printf("Failed to register service: %v", err)
		return err
	}
	return nil
}

// deregister service 注销服务
func (c *ConsulClient) DeregisterService(serviceID string) error {
	if err := c.client.Agent().ServiceDeregister(serviceID); err != nil {
		log.Printf("Failed to deregister service: %v", err)
		return err
	}
	return nil
}

// discover service 发现服务
func (c *ConsulClient) DiscoverService(serviceName string) ([]*api.ServiceEntry, error) {
	services, _, err := c.client.Health().Service(serviceName, "", true, nil)
	if err != nil {
		log.Printf("Failed to discover service: %v", err)
		return nil, err
	}
	return services, nil
}

// 向Consul写入KV配置
func (c *ConsulClient) Put(key string, config []byte) (*api.WriteMeta, error) {
	writeMeta, err := c.client.KV().Put(&api.KVPair{Key: key, Value: config}, nil)
	if err != nil {
		log.Printf("Failed to put KV: %v", err)
		return nil, err
	}
	return writeMeta, nil
}

// 从Consul获取KV配置
func (c *ConsulClient) Get(key string) ([]byte, error) {
	pair, _, err := c.client.KV().Get(key, nil)
	if err != nil {
		log.Printf("Failed to get KV: %v", err)
		return nil, err
	}
	if pair != nil {
		return pair.Value, nil
	}
	return nil, nil
}
