package consuls

import (
	"Taurus/config"
	"Taurus/pkg/consul"
	"Taurus/pkg/util"
)

type DefaultInitKVConfig struct {
}

// Put 初始化配置到KV
func (d *DefaultInitKVConfig) Put(c *consul.ConsulClient, serviceName string) error {
	c.PutKV(serviceName, "default", []byte(util.ToJsonString(config.Core)))
	return nil
}
