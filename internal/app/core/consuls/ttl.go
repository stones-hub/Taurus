package consuls

import (
	"Taurus/pkg/consul"

	"github.com/hashicorp/consul/api"
)

// 实现ttlupdate接口
type DefaultTTLUpdater struct {
}

// 更新TTL
func (u *DefaultTTLUpdater) Update(c *consul.ConsulClient, checkID string) error {
	c.UpdateTTL(checkID, api.HealthPassing, "TTL更新")
	return nil
}
