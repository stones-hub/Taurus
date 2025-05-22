package main

import (
	"Taurus/config"
	"Taurus/pkg/consul"
	"Taurus/pkg/util"
	"log"

	"github.com/hashicorp/consul/api"
)

// implement configwatcher interface
type DefaultConfigWatcher struct {
}

// handle config change
func (w *DefaultConfigWatcher) OnChange(c *consul.ConsulClient, serviceName string, key string, value []byte) error {
	log.Printf("配置变更: %s, %s", key, string(value))
	// 更新配置
	// TODO 解析，修改当前内存的配置即可
	return nil
}

// implement ttlupdate interface
type DefaultTTLUpdater struct {
}

// update TTL
func (u *DefaultTTLUpdater) Update(c *consul.ConsulClient, checkID string) error {
	log.Printf("update TTL..")
	c.UpdateTTL(checkID, api.HealthPassing, "TTL update")
	return nil
}

// implement initkvconfig interface
type DefaultInitKVConfig struct {
}

// put config to KV
func (d *DefaultInitKVConfig) Put(c *consul.ConsulClient, serviceName string) error {
	c.PutKV(serviceName, "default", []byte(util.ToJsonString(config.Core)))
	return nil
}

func main() {
	_, cleanup, err := consul.Init(&consul.ServerConfig{}, &consul.ServiceConfig{}, new(DefaultConfigWatcher), new(DefaultTTLUpdater), new(DefaultInitKVConfig))
	if err != nil {
		log.Fatalf("Failed to initialize consul: %v", err)
	}
	// then server stop , the cleanup will be called
	defer cleanup()
}
