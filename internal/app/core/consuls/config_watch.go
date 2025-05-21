package consuls

import (
	"Taurus/pkg/consul"
	"log"
)

// 实现configwatcher接口
type DefaultConfigWatcher struct {
}

// 处理配置变更
func (w *DefaultConfigWatcher) OnChange(c *consul.ConsulClient, serviceName string, key string, value []byte) error {
	log.Printf("配置变更: %s, %s", key, string(value))
	// 更新配置
	// TODO 解析，修改当前内存的配置即可
	return nil
}
